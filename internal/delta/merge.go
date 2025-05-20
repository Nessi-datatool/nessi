package delta

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/parquet"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
	"github.com/google/uuid"
)

// MergeCondition represents a condition for the merge operation
type MergeCondition struct {
	LeftColumn  string
	RightColumn string
	Operator    string
}

// MergeAction represents an action to take when a condition is met
type MergeAction struct {
	Type      string            // "update", "delete", "insert"
	Condition *MergeCondition   // Optional condition for the action
	Values    map[string]string // Column values for update/insert
}

// MergeOptions defines options for the merge operation
type MergeOptions struct {
	SourceTable    string
	Conditions     []MergeCondition
	Actions        []MergeAction
	WhenMatched    []MergeAction
	WhenNotMatched []MergeAction
}

// MergeStats represents statistics for a merge operation
type MergeStats struct {
	NumSourceRows     int64
	NumTargetRows     int64
	NumMatchedRows    int64
	NumNotMatchedRows int64
	NumUpdatedRows    int64
	NumDeletedRows    int64
	NumInsertedRows   int64
	DurationSeconds   float64
}

// Merge performs a merge operation on the Delta table
func (c *DeltaConnector) Merge(ctx context.Context, options MergeOptions) (*MergeStats, error) {
	stats := &MergeStats{}

	// Get source table schema
	sourceSchema, err := c.getSourceSchema(options.SourceTable)
	if err != nil {
		return nil, fmt.Errorf("failed to get source schema: %w", err)
	}

	// Get target table schema
	targetSchema := c.GetSchema()
	if targetSchema == nil {
		return nil, fmt.Errorf("target schema not found")
	}

	// Validate schemas
	if err := c.validateMergeSchemas(sourceSchema, targetSchema); err != nil {
		return nil, fmt.Errorf("invalid schemas: %w", err)
	}

	// Create merge plan
	plan, err := c.createMergePlan(options, sourceSchema, targetSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create merge plan: %w", err)
	}

	// Execute merge plan
	if err := c.executeMergePlan(ctx, plan, stats); err != nil {
		return nil, fmt.Errorf("failed to execute merge plan: %w", err)
	}

	return stats, nil
}

// getSourceSchema gets the schema of the source table
func (c *DeltaConnector) getSourceSchema(sourceTable string) (*arrow.Schema, error) {
	// Create a new connector for the source table
	sourceConnector := &DeltaConnector{
		TablePath: sourceTable,
		Config:    c.Config,
	}

	// Initialize the source connector
	if err := sourceConnector.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize source connector: %w", err)
	}

	// Get the schema from the source table
	schema := sourceConnector.GetSchema()
	if schema == nil {
		return nil, fmt.Errorf("failed to get source schema")
	}

	return schema, nil
}

// validateMergeSchemas validates the source and target schemas for a merge
func (c *DeltaConnector) validateMergeSchemas(source, target *arrow.Schema) error {
	// Safety check for nil schemas
	if source == nil {
		return fmt.Errorf("source schema is nil")
	}
	if target == nil {
		return fmt.Errorf("target schema is nil")
	}

	// Safety check for empty fields
	if len(source.Fields()) == 0 {
		return fmt.Errorf("source schema has no fields")
	}
	if len(target.Fields()) == 0 {
		return fmt.Errorf("target schema has no fields")
	}

	// Check if all source columns exist in target
	for _, sourceField := range source.Fields() {
		found := false
		for _, targetField := range target.Fields() {
			if targetField.Name == sourceField.Name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("source column %s not found in target", sourceField.Name)
		}
	}

	// Check if MergeOptions is initialized
	if c.MergeOptions == nil {
		// Not all operations might need this validation
		return nil
	}

	// First check matched actions
	for _, action := range append(c.MergeOptions.WhenMatched, c.MergeOptions.WhenNotMatched...) {
		if action.Values == nil {
			continue // Skip if no values to check
		}

		for col := range action.Values {
			found := false
			for _, field := range target.Fields() {
				if field.Name == col {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("target column %s not found", col)
			}
		}
	}

	return nil
}

// createMergePlan creates a plan for the merge operation
func (c *DeltaConnector) createMergePlan(options MergeOptions, sourceSchema, targetSchema *arrow.Schema) (*MergePlan, error) {
	plan := &MergePlan{
		SourceSchema:   sourceSchema,
		TargetSchema:   targetSchema,
		Conditions:     options.Conditions,
		Actions:        options.Actions,
		WhenMatched:    options.WhenMatched,
		WhenNotMatched: options.WhenNotMatched,
	}

	// Validate conditions
	for _, condition := range options.Conditions {
		if err := c.validateMergeCondition(condition, sourceSchema, targetSchema); err != nil {
			return nil, fmt.Errorf("invalid condition: %w", err)
		}
	}

	// Validate actions
	for _, action := range options.Actions {
		if err := c.validateMergeAction(action, sourceSchema, targetSchema); err != nil {
			return nil, fmt.Errorf("invalid action: %w", err)
		}
	}

	return plan, nil
}

// MergePlan represents a plan for a merge operation
type MergePlan struct {
	SourceSchema   *arrow.Schema
	TargetSchema   *arrow.Schema
	Conditions     []MergeCondition
	Actions        []MergeAction
	WhenMatched    []MergeAction
	WhenNotMatched []MergeAction
}

// validateMergeCondition validates a merge condition
func (c *DeltaConnector) validateMergeCondition(condition MergeCondition, sourceSchema, targetSchema *arrow.Schema) error {
	// Check if columns exist
	sourceFound := false
	targetFound := false

	for _, field := range sourceSchema.Fields() {
		if field.Name == condition.LeftColumn {
			sourceFound = true
			break
		}
	}

	for _, field := range targetSchema.Fields() {
		if field.Name == condition.RightColumn {
			targetFound = true
			break
		}
	}

	if !sourceFound {
		return fmt.Errorf("source column %s not found", condition.LeftColumn)
	}
	if !targetFound {
		return fmt.Errorf("target column %s not found", condition.RightColumn)
	}

	// Validate operator
	switch strings.ToUpper(condition.Operator) {
	case "=", "!=", ">", "<", ">=", "<=":
		// Valid operators
	default:
		return fmt.Errorf("invalid operator: %s", condition.Operator)
	}

	return nil
}

// validateMergeAction validates a merge action
func (c *DeltaConnector) validateMergeAction(action MergeAction, sourceSchema, targetSchema *arrow.Schema) error {
	// Validate action type
	switch strings.ToLower(action.Type) {
	case "update", "delete", "insert":
		// Valid action types
	default:
		return fmt.Errorf("invalid action type: %s", action.Type)
	}

	// Validate condition if present
	if action.Condition != nil {
		if err := c.validateMergeCondition(*action.Condition, sourceSchema, targetSchema); err != nil {
			return fmt.Errorf("invalid action condition: %w", err)
		}
	}

	// Validate values for update/insert
	if action.Type == "update" || action.Type == "insert" {
		for col := range action.Values {
			found := false
			for _, field := range targetSchema.Fields() {
				if field.Name == col {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("target column %s not found", col)
			}
		}
	}

	return nil
}

// executeMergePlan executes a merge plan
func (c *DeltaConnector) executeMergePlan(ctx context.Context, plan *MergePlan, stats *MergeStats) error {
	startTime := time.Now()

	// Reset merge stats to avoid double counting
	stats.NumMatchedRows = 0
	stats.NumNotMatchedRows = 0
	stats.NumUpdatedRows = 0
	stats.NumInsertedRows = 0

	// Read source data
	sourceData, err := c.readSourceData(ctx, plan.SourceSchema)
	if err != nil {
		return fmt.Errorf("failed to read source data: %w", err)
	}
	stats.NumSourceRows = sourceData.NumRows()

	// Read target data
	targetData, err := c.readTargetData(ctx, plan.TargetSchema)
	if err != nil {
		return fmt.Errorf("failed to read target data: %w", err)
	}
	stats.NumTargetRows = targetData.NumRows()

	// Create result builders
	resultBuilders := make([]array.Builder, len(plan.TargetSchema.Fields()))
	for i, field := range plan.TargetSchema.Fields() {
		resultBuilders[i] = array.NewBuilder(memory.DefaultAllocator, field.Type)
		defer resultBuilders[i].Release()
	}

	// Track which source rows are matched to avoid double counting
	matchedSourceRows := make(map[int64]bool)
	// Track which target rows are matched to avoid double counting
	matchedTargetRows := make(map[int64]bool)

	// First, identify all matches based on the conditions
	matches := make(map[int64]int64) // sourceRow -> targetRow
	for i := int64(0); i < sourceData.NumRows(); i++ {
		sourceRow := sourceData.NewSlice(i, i+1)
		defer sourceRow.Release()

		for j := int64(0); j < targetData.NumRows(); j++ {
			if matchedTargetRows[j] {
				continue
			}

			targetRow := targetData.NewSlice(j, j+1)
			defer targetRow.Release()

			if c.rowsMatch(sourceRow, targetRow, plan.Conditions) {
				// Record that we matched these rows
				matches[i] = j
				matchedSourceRows[i] = true
				matchedTargetRows[j] = true
				stats.NumMatchedRows++
				break
			}
		}
	}

	// Process all matched rows first
	for srcIdx, tgtIdx := range matches {
		sourceRow := sourceData.NewSlice(srcIdx, srcIdx+1)
		defer sourceRow.Release()
		targetRow := targetData.NewSlice(tgtIdx, tgtIdx+1)
		defer targetRow.Release()

		// Apply all matched actions (typically updates)
		updated := false
		for _, action := range plan.WhenMatched {
			if action.Type == "update" {
				// Apply the update action
				if err := c.applyUpdateAction(sourceRow, targetRow, action, resultBuilders); err != nil {
					return fmt.Errorf("failed to apply update action for match (%d,%d): %w", srcIdx, tgtIdx, err)
				}
				updated = true
			}
		}

		if updated {
			stats.NumUpdatedRows++
		} else {
			// No update action found, just copy the target row
			for i, builder := range resultBuilders {
				if err := c.copyValue(builder, targetRow.Column(i), 0); err != nil {
					return fmt.Errorf("failed to copy matched target row: %w", err)
				}
			}
		}
	}

	// Process unmatched target rows (keep them)
	for j := int64(0); j < targetData.NumRows(); j++ {
		if !matchedTargetRows[j] {
			targetRow := targetData.NewSlice(j, j+1)
			defer targetRow.Release()

			// Copy unmatched target rows
			for i, builder := range resultBuilders {
				if err := c.copyValue(builder, targetRow.Column(i), 0); err != nil {
					return fmt.Errorf("failed to copy unmatched target row: %w", err)
				}
			}
		}
	}

	// Process unmatched source rows (insert them if specified)
	for i := int64(0); i < sourceData.NumRows(); i++ {
		if !matchedSourceRows[i] {
			stats.NumNotMatchedRows++
			sourceRow := sourceData.NewSlice(i, i+1)
			defer sourceRow.Release()

			// Do we have any WhenNotMatched actions?
			if len(plan.WhenNotMatched) > 0 {
				// Apply the relevant action (typically insert)
				inserted := false
				for _, action := range plan.WhenNotMatched {
					if action.Type == "insert" {
						// Basic insert action - copy all values from source
						if action.Values == nil || len(action.Values) == 0 {
							for fieldIdx, field := range sourceRow.Schema().Fields() {
								// Find the corresponding target field
								targetFieldIdx := -1
								for idx, tgtField := range plan.TargetSchema.Fields() {
									if tgtField.Name == field.Name {
										targetFieldIdx = idx
										break
									}
								}

								if targetFieldIdx >= 0 {
									// Copy the value from source
									if err := c.copyValue(resultBuilders[targetFieldIdx], sourceRow.Column(fieldIdx), 0); err != nil {
										return fmt.Errorf("failed to copy source value for insert: %w", err)
									}
								} else {
									// Field not found in target, skip it
									continue
								}
							}
						} else {
							// Insert with specific values
							for fieldIdx, field := range plan.TargetSchema.Fields() {
								if value, ok := action.Values[field.Name]; ok {
									// Use the specified value
									if err := c.setValue(resultBuilders[fieldIdx], value, sourceRow, nil); err != nil {
										return fmt.Errorf("failed to set value for insert: %w", err)
									}
								} else {
									// Find the field in source to copy
									sourceFieldIdx := -1
									for idx, srcField := range sourceRow.Schema().Fields() {
										if srcField.Name == field.Name {
											sourceFieldIdx = idx
											break
										}
									}

									if sourceFieldIdx >= 0 {
										// Copy the value from source
										if err := c.copyValue(resultBuilders[fieldIdx], sourceRow.Column(sourceFieldIdx), 0); err != nil {
											return fmt.Errorf("failed to copy source value for insert: %w", err)
										}
									} else {
										// Field not found in source, use null
										resultBuilders[fieldIdx].AppendNull()
									}
								}
							}
						}
						inserted = true
						break // Only apply the first insert action
					}
				}

				if inserted {
					stats.NumInsertedRows++
				}
			}
		}
	}

	// Create result record
	fields := make([]arrow.Array, len(resultBuilders))
	for i, builder := range resultBuilders {
		fields[i] = builder.NewArray()
		defer fields[i].Release()
	}

	result := array.NewRecord(plan.TargetSchema, fields, int64(fields[0].Len()))
	defer result.Release()

	// Write result to target table
	if err := c.writeResult(ctx, result); err != nil {
		return fmt.Errorf("failed to write result: %w", err)
	}

	stats.DurationSeconds = time.Since(startTime).Seconds()
	return nil
}

// readSourceData reads data from the source table
func (c *DeltaConnector) readSourceData(ctx context.Context, schema *arrow.Schema) (arrow.Record, error) {
	// Get the latest version of the source table
	version := c.GetVersion()

	// Get all Parquet files for this version
	files, err := c.getParquetFiles(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get Parquet files: %w", err)
	}

	// Create record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builder.Release()

	// Read and combine data from all files
	for _, file := range files {
		record, err := readParquetFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read Parquet file %s: %w", file, err)
		}
		defer record.Release()

		// Append data to builders
		for i := 0; i < int(record.NumCols()); i++ {
			if err := c.appendArray(builder.Field(i), record.Column(i)); err != nil {
				return nil, fmt.Errorf("failed to append column %d: %w", i, err)
			}
		}
	}

	return builder.NewRecord(), nil
}

// readTargetData reads data from the target table
func (c *DeltaConnector) readTargetData(ctx context.Context, schema *arrow.Schema) (arrow.Record, error) {
	// Similar to readSourceData but for target table
	version := c.GetVersion()

	files, err := c.getParquetFiles(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get Parquet files: %w", err)
	}

	builder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builder.Release()

	for _, file := range files {
		record, err := readParquetFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read Parquet file %s: %w", file, err)
		}
		defer record.Release()

		for i := 0; i < int(record.NumCols()); i++ {
			if err := c.appendArray(builder.Field(i), record.Column(i)); err != nil {
				return nil, fmt.Errorf("failed to append column %d: %w", i, err)
			}
		}
	}

	return builder.NewRecord(), nil
}

// rowsMatch checks if source and target rows match based on conditions
func (c *DeltaConnector) rowsMatch(source, target arrow.Record, conditions []MergeCondition) bool {
	// If no conditions are specified, no rows match
	if len(conditions) == 0 {
		return false
	}

	for _, condition := range conditions {
		// Get the column indices
		sourceIdx := -1
		targetIdx := -1
		for i := 0; i < int(source.NumCols()); i++ {
			if source.ColumnName(i) == condition.LeftColumn {
				sourceIdx = i
				break // Found the match, stop searching
			}
		}
		for i := 0; i < int(target.NumCols()); i++ {
			if target.ColumnName(i) == condition.RightColumn {
				targetIdx = i
				break // Found the match, stop searching
			}
		}

		if sourceIdx == -1 || targetIdx == -1 {
			// If any column is not found, no match
			return false
		}

		// Debug the comparison
		leftArray := source.Column(sourceIdx)
		rightArray := target.Column(targetIdx)

		// For INT64 columns (like ID), directly compare values
		if leftArray.DataType().ID() == arrow.INT64 && rightArray.DataType().ID() == arrow.INT64 {
			// Special handling for ID comparisons to make exact matches
			leftInt := leftArray.(*array.Int64).Value(0)
			rightInt := rightArray.(*array.Int64).Value(0)

			switch condition.Operator {
			case "=", "==":
				if leftInt != rightInt {
					return false
				}
			case "!=", "<>":
				if leftInt == rightInt {
					return false
				}
			case "<":
				if leftInt >= rightInt {
					return false
				}
			case "<=":
				if leftInt > rightInt {
					return false
				}
			case ">":
				if leftInt <= rightInt {
					return false
				}
			case ">=":
				if leftInt < rightInt {
					return false
				}
			default:
				// Unsupported operator
				return false
			}
		} else {
			// For other types, use the generic comparison function
			if !c.compareValues(leftArray, rightArray, condition.Operator) {
				return false
			}
		}
	}

	// All conditions matched
	return true
}

// compareValues compares two values based on the operator
func (c *DeltaConnector) compareValues(left, right arrow.Array, operator string) bool {
	if left.Len() == 0 || right.Len() == 0 {
		return false
	}

	// Handle null values
	if left.IsNull(0) && right.IsNull(0) {
		return operator == "=" || operator == "=="
	} else if left.IsNull(0) || right.IsNull(0) {
		return operator == "!=" || operator == "<>"
	}

	// Compare based on type
	switch left.DataType().ID() {
	case arrow.STRING:
		leftStr := left.(*array.String).Value(0)
		rightStr := right.(*array.String).Value(0)
		switch operator {
		case "=", "==":
			return leftStr == rightStr
		case "!=", "<>":
			return leftStr != rightStr
		case "<":
			return leftStr < rightStr
		case "<=":
			return leftStr <= rightStr
		case ">":
			return leftStr > rightStr
		case ">=":
			return leftStr >= rightStr
		}

	case arrow.INT64:
		leftInt := left.(*array.Int64).Value(0)
		rightInt := right.(*array.Int64).Value(0)
		switch operator {
		case "=", "==":
			return leftInt == rightInt
		case "!=", "<>":
			return leftInt != rightInt
		case "<":
			return leftInt < rightInt
		case "<=":
			return leftInt <= rightInt
		case ">":
			return leftInt > rightInt
		case ">=":
			return leftInt >= rightInt
		}

	case arrow.FLOAT64:
		leftFloat := left.(*array.Float64).Value(0)
		rightFloat := right.(*array.Float64).Value(0)
		switch operator {
		case "=", "==":
			return leftFloat == rightFloat
		case "!=", "<>":
			return leftFloat != rightFloat
		case "<":
			return leftFloat < rightFloat
		case "<=":
			return leftFloat <= rightFloat
		case ">":
			return leftFloat > rightFloat
		case ">=":
			return leftFloat >= rightFloat
		}

	case arrow.BOOL:
		leftBool := left.(*array.Boolean).Value(0)
		rightBool := right.(*array.Boolean).Value(0)
		switch operator {
		case "=", "==":
			return leftBool == rightBool
		case "!=", "<>":
			return leftBool != rightBool
		}

	case arrow.TIMESTAMP:
		leftTs := left.(*array.Timestamp).Value(0)
		rightTs := right.(*array.Timestamp).Value(0)
		switch operator {
		case "=", "==":
			return leftTs == rightTs
		case "!=", "<>":
			return leftTs != rightTs
		case "<":
			return leftTs < rightTs
		case "<=":
			return leftTs <= rightTs
		case ">":
			return leftTs > rightTs
		case ">=":
			return leftTs >= rightTs
		}
	}

	// Default: use reflect.DeepEqual for unsupported types or operators
	switch operator {
	case "=", "==":
		return reflect.DeepEqual(left, right)
	case "!=", "<>":
		return !reflect.DeepEqual(left, right)
	}

	return false
}

// copyValue copies a value from one array to a builder
func (c *DeltaConnector) copyValue(builder array.Builder, arr arrow.Array, index int64) error {
	if arr.IsNull(int(index)) {
		builder.AppendNull()
		return nil
	}

	switch builder := builder.(type) {
	case *array.StringBuilder:
		// Handle string data
		switch arr := arr.(type) {
		case *array.String:
			str := arr.Value(int(index))
			builder.Append(str)
		case *array.Int64:
			// Convert int64 to string
			val := arr.Value(int(index))
			builder.Append(fmt.Sprintf("%d", val))
		case *array.Float64:
			// Convert float64 to string
			val := arr.Value(int(index))
			builder.Append(fmt.Sprintf("%g", val))
		case *array.Boolean:
			// Convert boolean to string
			val := arr.Value(int(index))
			builder.Append(fmt.Sprintf("%t", val))
		default:
			return fmt.Errorf("cannot convert %T to string", arr)
		}

	case *array.Int64Builder:
		switch arr := arr.(type) {
		case *array.Int64:
			val := arr.Value(int(index))
			builder.Append(val)
		case *array.String:
			// Try to parse string to int64
			str := arr.Value(int(index))
			val, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot convert string '%s' to int64: %w", str, err)
			}
			builder.Append(val)
		default:
			return fmt.Errorf("cannot convert %T to int64", arr)
		}

	case *array.Float64Builder:
		switch arr := arr.(type) {
		case *array.Float64:
			val := arr.Value(int(index))
			builder.Append(val)
		case *array.String:
			// Try to parse string to float64
			str := arr.Value(int(index))
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return fmt.Errorf("cannot convert string '%s' to float64: %w", str, err)
			}
			builder.Append(val)
		default:
			return fmt.Errorf("cannot convert %T to float64", arr)
		}

	case *array.BooleanBuilder:
		switch arr := arr.(type) {
		case *array.Boolean:
			val := arr.Value(int(index))
			builder.Append(val)
		case *array.String:
			// Try to parse string to boolean
			str := arr.Value(int(index))
			val, err := strconv.ParseBool(str)
			if err != nil {
				return fmt.Errorf("cannot convert string '%s' to boolean: %w", str, err)
			}
			builder.Append(val)
		default:
			return fmt.Errorf("cannot convert %T to boolean", arr)
		}

	case *array.TimestampBuilder:
		switch arr := arr.(type) {
		case *array.Timestamp:
			val := arr.Value(int(index))
			builder.Append(val)
		default:
			return fmt.Errorf("cannot convert %T to timestamp", arr)
		}

	default:
		return fmt.Errorf("unsupported builder type: %s", builder.Type())
	}

	return nil
}

// setValue sets a value in a builder
func (c *DeltaConnector) setValue(builder array.Builder, value string, sourceRecord, targetRecord arrow.Record) error {
	if value == "null" {
		builder.AppendNull()
		return nil
	}

	// Handle column references (e.g., "source.value" or "target.value")
	// This should be done before trying to parse as literals
	if strings.HasPrefix(value, "source.") || strings.HasPrefix(value, "target.") {
		// Extract the column name from the reference
		parts := strings.SplitN(value, ".", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid column reference format: %s", value)
		}

		tableRef := parts[0] // "source" or "target"
		colName := parts[1]  // The column name

		// Get the value from the appropriate record based on the reference
		var record arrow.Record
		switch tableRef {
		case "source":
			if sourceRecord == nil {
				return fmt.Errorf("source record not available for column reference: %s", value)
			}
			record = sourceRecord
		case "target":
			if targetRecord == nil {
				return fmt.Errorf("target record not available for column reference: %s", value)
			}
			record = targetRecord
		default:
			return fmt.Errorf("unknown table reference: %s", tableRef)
		}

		// Find the column index in the record
		colIdx := -1
		for i, field := range record.Schema().Fields() {
			if field.Name == colName {
				colIdx = i
				break
			}
		}

		if colIdx == -1 {
			return fmt.Errorf("column not found in %s: %s", tableRef, colName)
		}

		// Copy the value from the reference to the builder
		return c.copyValue(builder, record.Column(colIdx), 0)
	}

	// If not a column reference, parse as a literal value
	switch builder := builder.(type) {
	case *array.StringBuilder:
		builder.Append(value)

	case *array.Int64Builder:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int64: %w", err)
		}
		builder.Append(val)

	case *array.Float64Builder:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float64: %w", err)
		}
		builder.Append(val)

	case *array.BooleanBuilder:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("failed to parse bool: %w", err)
		}
		builder.Append(val)

	case *array.TimestampBuilder:
		t, err := time.Parse(time.RFC3339, value)
		if err != nil {
			// Try parsing as Unix timestamp (milliseconds)
			milliseconds, err2 := strconv.ParseInt(value, 10, 64)
			if err2 != nil {
				return fmt.Errorf("failed to parse timestamp (not RFC3339 or Unix timestamp): %w", err)
			}
			builder.Append(arrow.Timestamp(milliseconds * 1000000)) // Convert ms to ns
		} else {
			builder.Append(arrow.Timestamp(t.UnixNano()))
		}

	default:
		return fmt.Errorf("unsupported type: %s", builder.Type())
	}

	return nil
}

// Helper functions for value comparison
func (c *DeltaConnector) compareStrings(source, target, operator string) bool {
	switch operator {
	case "=":
		return source == target
	case "!=":
		return source != target
	case ">":
		return source > target
	case "<":
		return source < target
	case ">=":
		return source >= target
	case "<=":
		return source <= target
	default:
		return false
	}
}

func (c *DeltaConnector) compareInts(source, target int64, operator string) bool {
	switch operator {
	case "=":
		return source == target
	case "!=":
		return source != target
	case ">":
		return source > target
	case "<":
		return source < target
	case ">=":
		return source >= target
	case "<=":
		return source <= target
	default:
		return false
	}
}

func (c *DeltaConnector) compareFloats(source, target float64, operator string) bool {
	switch operator {
	case "=":
		return math.Abs(source-target) < 1e-10
	case "!=":
		return math.Abs(source-target) >= 1e-10
	case ">":
		return source > target
	case "<":
		return source < target
	case ">=":
		return source >= target
	case "<=":
		return source <= target
	default:
		return false
	}
}

func (c *DeltaConnector) compareBools(source, target bool, operator string) bool {
	switch operator {
	case "=":
		return source == target
	case "!=":
		return source != target
	default:
		return false
	}
}

func (c *DeltaConnector) compareTimestamps(source, target arrow.Timestamp, operator string) bool {
	switch operator {
	case "=":
		return source == target
	case "!=":
		return source != target
	case ">":
		return source > target
	case "<":
		return source < target
	case ">=":
		return source >= target
	case "<=":
		return source <= target
	default:
		return false
	}
}

// appendArray appends all values from one array to a builder
func (c *DeltaConnector) appendArray(builder array.Builder, arr arrow.Array) error {
	for i := 0; i < arr.Len(); i++ {
		if err := c.copyValue(builder, arr, int64(i)); err != nil {
			return err
		}
	}
	return nil
}

// getParquetFiles returns all Parquet files for a given version
func (c *DeltaConnector) getParquetFiles(version int64) ([]string, error) {
	// Read the transaction log for this version
	logPath := filepath.Join(c.TablePath, "_delta_log", fmt.Sprintf("%020d.json", version))
	data, err := os.ReadFile(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	var actions []Action
	if err := json.Unmarshal(data, &actions); err != nil {
		return nil, fmt.Errorf("failed to parse log file: %w", err)
	}

	// Get all add actions
	var files []string
	for _, action := range actions {
		if action.Add != nil {
			files = append(files, filepath.Join(c.TablePath, action.Add.Path))
		}
	}

	return files, nil
}

// readParquetFile reads a Parquet file and returns an Arrow record

// writeParquetFile writes an Arrow record to a Parquet file

// applyMatchedActions applies actions for matched rows
func (c *DeltaConnector) applyMatchedActions(source, target arrow.Record, actions []MergeAction, builders []array.Builder) error {
	applied := false

	for _, action := range actions {
		if action.Condition != nil && !c.rowsMatch(source, target, []MergeCondition{*action.Condition}) {
			continue
		}

		switch strings.ToLower(action.Type) {
		case "update":
			if err := c.applyUpdateAction(source, target, action, builders); err != nil {
				return fmt.Errorf("failed to apply update action: %w", err)
			}
			// We'll track updates in the executeMergePlan method
			applied = true

		case "delete":
			// Skip this row
			return nil
		}
	}

	// If no actions were applied, copy target row
	if !applied {
		for i, builder := range builders {
			if err := c.copyValue(builder, target.Column(i), 0); err != nil {
				return fmt.Errorf("failed to copy target row: %w", err)
			}
		}
	}

	return nil
}

// applyNotMatchedActions applies actions for unmatched rows
func (c *DeltaConnector) applyNotMatchedActions(source arrow.Record, actions []MergeAction, builders []array.Builder) error {
	if len(actions) == 0 {
		// If no actions are specified, do nothing for unmatched rows
		return nil
	}

	for _, action := range actions {
		switch strings.ToLower(action.Type) {
		case "insert":
			if action.Values == nil || len(action.Values) == 0 {
				// Insert all source columns without specifying values
				for i, field := range source.Schema().Fields() {
					name := field.Name
					// Find the field in the target schema
					found := false
					targetIdx := -1
					for j, targetField := range c.Schema.Fields() {
						if targetField.Name == name {
							found = true
							targetIdx = j
							break
						}
					}
					if found {
						if err := c.copyValue(builders[targetIdx], source.Column(i), 0); err != nil {
							return fmt.Errorf("failed to copy value for %s: %w", name, err)
						}
					}
				}
			} else {
				// Insert specified values
				for i, field := range c.Schema.Fields() {
					if value, ok := action.Values[field.Name]; ok {
						// Use value from action
						if err := c.setValue(builders[i], value, source, nil); err != nil {
							return fmt.Errorf("failed to set value for %s: %w", field.Name, err)
						}
					} else {
						// Use default value (null)
						builders[i].AppendNull()
					}
				}
			}
		}
	}

	return nil
}

// applyUpdateAction applies an update action
func (c *DeltaConnector) applyUpdateAction(source, target arrow.Record, action MergeAction, builders []array.Builder) error {
	for i, field := range target.Schema().Fields() {
		if value, ok := action.Values[field.Name]; ok {
			// Use value from action
			if err := c.setValue(builders[i], value, source, target); err != nil {
				return fmt.Errorf("failed to set value for %s: %w", field.Name, err)
			}
		} else {
			// Copy value from target
			if err := c.copyValue(builders[i], target.Column(i), 0); err != nil {
				return fmt.Errorf("failed to copy value for %s: %w", field.Name, err)
			}
		}
	}
	return nil
}

// applyInsertAction applies an insert action
func (c *DeltaConnector) applyInsertAction(source arrow.Record, action MergeAction, builders []array.Builder) error {
	for i, field := range source.Schema().Fields() {
		if value, ok := action.Values[field.Name]; ok {
			// Use value from action
			if err := c.setValue(builders[i], value, source, nil); err != nil {
				return fmt.Errorf("failed to set value for %s: %w", field.Name, err)
			}
		} else {
			// Copy value from source
			if err := c.copyValue(builders[i], source.Column(i), 0); err != nil {
				return fmt.Errorf("failed to copy value for %s: %w", field.Name, err)
			}
		}
	}
	return nil
}

// writeResult writes the result record to the target table
func (c *DeltaConnector) writeResult(ctx context.Context, result arrow.Record) error {
	// Create a new transaction using the connector's transaction log
	txn := c.log.BeginTransaction()
	if txn.CommitInfo == nil { // Should be initialized by BeginTransaction, but good practice to check
		txn.CommitInfo = &CommitInfo{}
	}
	txn.CommitInfo.Operation = "MERGE"
	// ReadVersion and Timestamp are typically set by BeginTransaction and Commit
	// txn.CommitInfo.ReadVersion = c.GetVersion() // This is set by BeginTransaction from t.log.version
	// txn.CommitInfo.Timestamp = time.Now().UnixMilli() // This is set by BeginTransaction
	// Ensure other necessary CommitInfo fields are set if not covered by BeginTransaction
	txn.CommitInfo.IsolationLevel = "Serializable" // Or whatever is appropriate
	txn.CommitInfo.IsBlindAppend = false

	// Generate a unique file path for the new data
	// timestamp := time.Now().UnixNano() / 1000000 // milliseconds - Use commit timestamp
	fileName := fmt.Sprintf("part-%d-%s.parquet", txn.CommitInfo.Timestamp, uuid.New().String())
	filePath := filepath.Join(c.TablePath, "data", fileName)

	// Ensure the data directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Write the Parquet file
	if err := c.writeParquetFile(filePath, result); err != nil {
		return fmt.Errorf("failed to write Parquet file: %w", err)
	}

	// Add add action for the new file
	addAction := Action{
		Add: &AddAction{
			Path:             fileName, // Store relative path for AddAction
			Size:             getFileSize(filePath),
			ModificationTime: txn.CommitInfo.Timestamp, // Use commit timestamp
			DataChange:       true,
			PartitionValues:  map[string]string{},
		},
	}

	// Add partition values if table is partitioned
	if len(c.PartitionColumns) > 0 {
		for _, col := range c.PartitionColumns {
			if colIdx := result.Schema().FieldIndices(col); len(colIdx) > 0 {
				arr := result.Column(colIdx[0])
				if arr.Len() > 0 && !arr.IsNull(0) {
					addAction.Add.PartitionValues[col] = c.getStringValue(arr, 0)
				}
			}
		}
	}

	// Add the action to the transaction
	if err := txn.AddAction(addAction); err != nil {
		// Potentially rollback or clean up filePath if AddAction fails
		return fmt.Errorf("failed to add action to transaction: %w", err)
	}

	// Commit the transaction
	if err := c.CommitTransaction(ctx, txn); err != nil {
		// Potentially rollback or clean up filePath if commit fails
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// writeParquetFile writes an Arrow record to a Parquet file
func (c *DeltaConnector) writeParquetFile(path string, record arrow.Record) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	writer, err := pqarrow.NewFileWriter(record.Schema(), f, parquet.NewWriterProperties(), pqarrow.NewArrowWriterProperties())
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %w", err)
	}
	defer writer.Close()

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}

// calculateStats calculates statistics for a record
func (c *DeltaConnector) calculateStats(record arrow.Record) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	stats["numRecords"] = record.NumRows()

	// Calculate column statistics
	columnStats := make(map[string]interface{})
	for i, field := range record.Schema().Fields() {
		colStats, err := c.calculateColumnStats(record.Column(i), field.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate stats for column %s: %w", field.Name, err)
		}
		columnStats[field.Name] = colStats
	}
	stats["columns"] = columnStats

	return stats, nil
}

// calculateColumnStats calculates statistics for a column
func (c *DeltaConnector) calculateColumnStats(arr arrow.Array, dataType arrow.DataType) (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"nullCount": arr.NullN(),
	}

	if arr.Len() == 0 {
		return stats, nil
	}

	switch arr := arr.(type) {
	case *array.String:
		stats["minValue"] = c.getMinString(arr)
		stats["maxValue"] = c.getMaxString(arr)

	case *array.Int64:
		stats["minValue"] = c.getMinInt64(arr)
		stats["maxValue"] = c.getMaxInt64(arr)

	case *array.Float64:
		stats["minValue"] = c.getMinFloat64(arr)
		stats["maxValue"] = c.getMaxFloat64(arr)

	case *array.Boolean:
		stats["trueCount"] = c.countBooleans(arr, true)
		stats["falseCount"] = c.countBooleans(arr, false)

	case *array.Timestamp:
		stats["minValue"] = c.getMinTimestamp(arr)
		stats["maxValue"] = c.getMaxTimestamp(arr)
	}

	return stats, nil
}

// Helper functions for calculating statistics
func (c *DeltaConnector) getMinString(arr *array.String) string {
	min := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val < min {
				min = val
			}
		}
	}
	return min
}

func (c *DeltaConnector) getMaxString(arr *array.String) string {
	max := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val > max {
				max = val
			}
		}
	}
	return max
}

func (c *DeltaConnector) getMinInt64(arr *array.Int64) int64 {
	min := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val < min {
				min = val
			}
		}
	}
	return min
}

func (c *DeltaConnector) getMaxInt64(arr *array.Int64) int64 {
	max := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val > max {
				max = val
			}
		}
	}
	return max
}

func (c *DeltaConnector) getMinFloat64(arr *array.Float64) float64 {
	min := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val < min {
				min = val
			}
		}
	}
	return min
}

func (c *DeltaConnector) getMaxFloat64(arr *array.Float64) float64 {
	max := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val > max {
				max = val
			}
		}
	}
	return max
}

func (c *DeltaConnector) getMinTimestamp(arr *array.Timestamp) arrow.Timestamp {
	min := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val < min {
				min = val
			}
		}
	}
	return min
}

func (c *DeltaConnector) getMaxTimestamp(arr *array.Timestamp) arrow.Timestamp {
	max := arr.Value(0)
	for i := 1; i < arr.Len(); i++ {
		if !arr.IsNull(i) {
			val := arr.Value(i)
			if val > max {
				max = val
			}
		}
	}
	return max
}

func (c *DeltaConnector) countBooleans(arr *array.Boolean, value bool) int64 {
	var count int64
	for i := 0; i < arr.Len(); i++ {
		if !arr.IsNull(i) && arr.Value(i) == value {
			count++
		}
	}
	return count
}

// getPartitionValues extracts partition values from a record
func (c *DeltaConnector) getPartitionValues(record arrow.Record) map[string]string {
	values := make(map[string]string)
	for _, col := range c.PartitionColumns {
		if colIdx := record.Schema().FieldIndices(col); len(colIdx) > 0 {
			arr := record.Column(colIdx[0])
			if arr.Len() > 0 && !arr.IsNull(0) {
				values[col] = c.getStringValue(arr, 0)
			}
		}
	}
	return values
}

// getStringValue converts any array value to string
func (c *DeltaConnector) getStringValue(arr arrow.Array, idx int) string {
	switch arr := arr.(type) {
	case *array.String:
		return arr.Value(idx)
	case *array.Int64:
		return strconv.FormatInt(arr.Value(idx), 10)
	case *array.Float64:
		return strconv.FormatFloat(arr.Value(idx), 'f', -1, 64)
	case *array.Boolean:
		return strconv.FormatBool(arr.Value(idx))
	case *array.Timestamp:
		return time.Unix(0, int64(arr.Value(idx))).Format(time.RFC3339)
	default:
		return ""
	}
}

// createRemoveActions creates remove actions for old files
func (c *DeltaConnector) createRemoveActions() ([]*RemoveAction, error) {
	// Get current version files
	files, err := c.getParquetFiles(c.GetVersion())
	if err != nil {
		return nil, fmt.Errorf("failed to get current files: %w", err)
	}

	var actions []*RemoveAction
	for _, file := range files {
		action := &RemoveAction{
			Path:         strings.TrimPrefix(file, c.TablePath+"/"),
			DeletionTime: time.Now().UnixNano() / 1000000,
			DataChange:   true,
		}
		actions = append(actions, action)
	}

	return actions, nil
}
