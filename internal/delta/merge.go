package delta

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/apache/arrow/go/v12/parquet"
	"github.com/apache/arrow/go/v12/parquet/file"
	"github.com/apache/arrow/go/v12/parquet/pqarrow"
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
	SourceTable string
	Conditions  []MergeCondition
	Actions     []MergeAction
	WhenMatched []MergeAction
	WhenNotMatched []MergeAction
}

// MergeStats represents statistics for a merge operation
type MergeStats struct {
	NumSourceRows    int64
	NumTargetRows    int64
	NumMatchedRows   int64
	NumNotMatchedRows int64
	NumUpdatedRows   int64
	NumDeletedRows   int64
	NumInsertedRows  int64
	DurationSeconds  float64
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
	// TODO: Implement source schema retrieval
	return nil, fmt.Errorf("not implemented")
}

// validateMergeSchemas validates the source and target schemas for a merge
func (c *DeltaConnector) validateMergeSchemas(source, target *arrow.Schema) error {
	// Check if all source columns exist in target
	for i := 0; i < source.NumFields(); i++ {
		sourceField := source.Field(i)
		found := false
		for j := 0; j < target.NumFields(); j++ {
			if target.Field(j).Name == sourceField.Name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("source column %s not found in target", sourceField.Name)
		}
	}

	// Check if all target columns referenced in actions exist
	for _, action := range c.mergeActions {
		for col := range action.Values {
			found := false
			for i := 0; i < target.NumFields(); i++ {
				if target.Field(i).Name == col {
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
		SourceSchema: sourceSchema,
		TargetSchema: targetSchema,
		Conditions:   options.Conditions,
		Actions:      options.Actions,
		WhenMatched:  options.WhenMatched,
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

	for i := 0; i < sourceSchema.NumFields(); i++ {
		if sourceSchema.Field(i).Name == condition.LeftColumn {
			sourceFound = true
			break
		}
	}

	for i := 0; i < targetSchema.NumFields(); i++ {
		if targetSchema.Field(i).Name == condition.RightColumn {
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
			for i := 0; i < targetSchema.NumFields(); i++ {
				if targetSchema.Field(i).Name == col {
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
	resultBuilders := make([]array.Builder, plan.TargetSchema.NumFields())
	for i, field := range plan.TargetSchema.Fields() {
		resultBuilders[i] = array.NewBuilder(memory.DefaultAllocator, field.Type)
		defer resultBuilders[i].Release()
	}

	// Process matched rows
	matchedRows := make(map[int64]bool)
	for i := int64(0); i < sourceData.NumRows(); i++ {
		sourceRow := sourceData.NewSlice(i, i+1)
		defer sourceRow.Release()

		for j := int64(0); j < targetData.NumRows(); j++ {
			if matchedRows[j] {
				continue
			}

			targetRow := targetData.NewSlice(j, j+1)
			defer targetRow.Release()

			// Check if rows match based on conditions
			if c.rowsMatch(sourceRow, targetRow, plan.Conditions) {
				matchedRows[j] = true
				stats.NumMatchedRows++

				// Apply matched actions
				if err := c.applyMatchedActions(sourceRow, targetRow, plan.WhenMatched, resultBuilders); err != nil {
					return fmt.Errorf("failed to apply matched actions: %w", err)
				}
				break
			}
		}
	}

	// Process unmatched target rows
	for j := int64(0); j < targetData.NumRows(); j++ {
		if !matchedRows[j] {
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

	// Process unmatched source rows
	for i := int64(0); i < sourceData.NumRows(); i++ {
		sourceRow := sourceData.NewSlice(i, i+1)
		defer sourceRow.Release()

		matched := false
		for j := int64(0); j < targetData.NumRows(); j++ {
			if matchedRows[j] {
				continue
			}

			targetRow := targetData.NewSlice(j, j+1)
			defer targetRow.Release()

			if c.rowsMatch(sourceRow, targetRow, plan.Conditions) {
				matched = true
				break
			}
		}

		if !matched {
			stats.NumNotMatchedRows++
			// Apply not matched actions
			if err := c.applyNotMatchedActions(sourceRow, plan.WhenNotMatched, resultBuilders); err != nil {
				return fmt.Errorf("failed to apply not matched actions: %w", err)
			}
		}
	}

	// Create result record
	fields := make([]arrow.Array, len(resultBuilders))
	for i, builder := range resultBuilders {
		fields[i] = builder.NewArray()
		defer fields[i].Release()
	}

	result := array.NewRecord(plan.TargetSchema, fields, fields[0].Len())
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
		record, err := c.readParquetFile(file, schema)
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
		record, err := c.readParquetFile(file, schema)
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
	for _, condition := range conditions {
		sourceCol := source.ColumnByName(condition.LeftColumn)
		targetCol := target.ColumnByName(condition.RightColumn)

		if sourceCol == nil || targetCol == nil {
			return false
		}

		if !c.compareValues(sourceCol, targetCol, condition.Operator) {
			return false
		}
	}
	return true
}

// compareValues compares two values based on the operator
func (c *DeltaConnector) compareValues(source, target arrow.Array, operator string) bool {
	if source.Len() == 0 || target.Len() == 0 {
		return false
	}

	// Handle null values
	if source.IsNull(0) || target.IsNull(0) {
		return false
	}

	switch source.DataType().ID() {
	case arrow.STRING:
		sourceStr := source.(*array.String).Value(0)
		targetStr := target.(*array.String).Value(0)
		return c.compareStrings(sourceStr, targetStr, operator)
	
	case arrow.INT64:
		sourceInt := source.(*array.Int64).Value(0)
		targetInt := target.(*array.Int64).Value(0)
		return c.compareInts(sourceInt, targetInt, operator)
	
	case arrow.FLOAT64:
		sourceFloat := source.(*array.Float64).Value(0)
		targetFloat := target.(*array.Float64).Value(0)
		return c.compareFloats(sourceFloat, targetFloat, operator)
	
	case arrow.BOOL:
		sourceBool := source.(*array.Boolean).Value(0)
		targetBool := target.(*array.Boolean).Value(0)
		return c.compareBools(sourceBool, targetBool, operator)
	
	case arrow.TIMESTAMP:
		sourceTs := source.(*array.Timestamp).Value(0)
		targetTs := target.(*array.Timestamp).Value(0)
		return c.compareTimestamps(sourceTs, targetTs, operator)
	
	default:
		return false
	}
}

// copyValue copies a value from one array to a builder
func (c *DeltaConnector) copyValue(builder array.Builder, arr arrow.Array, index int64) error {
	if arr.IsNull(int(index)) {
		builder.AppendNull()
		return nil
	}

	switch builder := builder.(type) {
	case *array.StringBuilder:
		str := arr.(*array.String).Value(int(index))
		builder.Append(str)
	
	case *array.Int64Builder:
		val := arr.(*array.Int64).Value(int(index))
		builder.Append(val)
	
	case *array.Float64Builder:
		val := arr.(*array.Float64).Value(int(index))
		builder.Append(val)
	
	case *array.BooleanBuilder:
		val := arr.(*array.Boolean).Value(int(index))
		builder.Append(val)
	
	case *array.TimestampBuilder:
		val := arr.(*array.Timestamp).Value(int(index))
		builder.Append(val)
	
	default:
		return fmt.Errorf("unsupported type: %s", builder.Type())
	}

	return nil
}

// setValue sets a value in a builder
func (c *DeltaConnector) setValue(builder array.Builder, value string) error {
	if value == "null" {
		builder.AppendNull()
		return nil
	}

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
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
		builder.Append(arrow.Timestamp(t.UnixNano()))
	
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
	logPath := filepath.Join(c.tablePath, "_delta_log", fmt.Sprintf("%020d.json", version))
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
			files = append(files, filepath.Join(c.tablePath, action.Add.Path))
		}
	}

	return files, nil
}

// readParquetFile reads a Parquet file and returns an Arrow record
func (c *DeltaConnector) readParquetFile(path string, schema *arrow.Schema) (arrow.Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader, err := pqarrow.NewFileReader(file, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	record, err := reader.ReadRowGroup(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read row group: %w", err)
	}

	return record, nil
}

// applyMatchedActions applies actions for matched rows
func (c *DeltaConnector) applyMatchedActions(source, target arrow.Record, actions []MergeAction, builders []array.Builder) error {
	for _, action := range actions {
		if action.Condition != nil && !c.rowsMatch(source, target, []MergeCondition{*action.Condition}) {
			continue
		}

		switch strings.ToLower(action.Type) {
		case "update":
			if err := c.applyUpdateAction(source, target, action, builders); err != nil {
				return fmt.Errorf("failed to apply update action: %w", err)
			}
		case "delete":
			// Skip this row
			return nil
		}
	}

	// If no actions were applied, copy target row
	for i, builder := range builders {
		if err := c.copyValue(builder, target.Column(i), 0); err != nil {
			return fmt.Errorf("failed to copy target row: %w", err)
		}
	}

	return nil
}

// applyNotMatchedActions applies actions for unmatched rows
func (c *DeltaConnector) applyNotMatchedActions(source arrow.Record, actions []MergeAction, builders []array.Builder) error {
	for _, action := range actions {
		if action.Condition != nil && !c.rowsMatch(source, source, []MergeCondition{*action.Condition}) {
			continue
		}

		switch strings.ToLower(action.Type) {
		case "insert":
			if err := c.applyInsertAction(source, action, builders); err != nil {
				return fmt.Errorf("failed to apply insert action: %w", err)
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
			if err := c.setValue(builders[i], value); err != nil {
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
			if err := c.setValue(builders[i], value); err != nil {
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
	// Create a new transaction
	txn := NewTransaction()
	txn.ReadVersion = c.GetVersion()

	// Generate a unique file path for the new data
	timestamp := time.Now().UnixNano() / 1000000 // milliseconds
	fileName := fmt.Sprintf("part-%d-%s.parquet", timestamp, uuid.New().String())
	filePath := filepath.Join(c.tablePath, "data", fileName)

	// Ensure the data directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Write the Parquet file
	if err := c.writeParquetFile(filePath, result); err != nil {
		return fmt.Errorf("failed to write Parquet file: %w", err)
	}

	// Calculate statistics for the new file
	stats, err := c.calculateStats(result)
	if err != nil {
		return fmt.Errorf("failed to calculate statistics: %w", err)
	}

	// Create add action for the new file
	addAction := &AddAction{
		Path:            filepath.Join("data", fileName),
		Size:            getFileSize(filePath),
		ModificationTime: timestamp,
		DataChange:      true,
		Stats:           stats,
		Tags:            map[string]string{"merge": "true"},
	}

	// Add partition values if table is partitioned
	if len(c.partitionColumns) > 0 {
		addAction.PartitionValues = c.getPartitionValues(result)
	}

	// Add the action to the transaction
	txn.Add(addAction)

	// Remove old files
	removeActions, err := c.createRemoveActions()
	if err != nil {
		return fmt.Errorf("failed to create remove actions: %w", err)
	}
	for _, action := range removeActions {
		txn.Add(action)
	}

	// Commit the transaction
	if err := c.CommitTransaction(ctx, txn); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// writeParquetFile writes an Arrow record to a Parquet file
func (c *DeltaConnector) writeParquetFile(path string, record arrow.Record) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Create Parquet writer properties
	writerProps := parquet.NewWriterProperties(
		parquet.WithCompression(parquet.CompressionCodec_SNAPPY),
		parquet.WithDictionaryEnabled(true),
	)

	// Create Arrow writer properties
	arrowProps := pqarrow.NewArrowWriterProperties(
		pqarrow.WithStoreSchema(),
		pqarrow.WithCompliantNestedTypes(),
	)

	// Create Parquet writer
	writer, err := pqarrow.NewFileWriter(
		record.Schema(),
		file,
		writerProps,
		arrowProps,
	)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}
	defer writer.Close()

	// Write the record
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
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
	for _, col := range c.partitionColumns {
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
			Path:            strings.TrimPrefix(file, c.tablePath+"/"),
			DeletionTime:    time.Now().UnixNano() / 1000000,
			DataChange:      true,
		}
		actions = append(actions, action)
	}

	return actions, nil
}

// mergeActions is a helper variable for schema validation
var mergeActions []MergeAction 