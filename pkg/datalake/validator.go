package datalake

import (
	"fmt"
	"math"
	"strings"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
)

// SchemaValidationError represents a schema validation error
type SchemaValidationError struct {
	FieldName    string
	ExpectedType string
	ActualType   string
	RowIndex     int64
	Value        string
}

// SchemaValidator validates data against a schema
type SchemaValidator struct {
	schema *arrow.Schema
}

// NewSchemaValidator creates a new schema validator
func NewSchemaValidator(schema *arrow.Schema) *SchemaValidator {
	return &SchemaValidator{
		schema: schema,
	}
}

// ValidateRecord validates an Arrow record against the schema
func (v *SchemaValidator) ValidateRecord(record arrow.Record) []SchemaValidationError {
	var errors []SchemaValidationError

	// Check if schema matches
	if !v.schema.Equal(record.Schema()) {
		// Schema doesn't match, validate field by field
		for _, expectedField := range v.schema.Fields() {
			// Find matching field in record
			fieldIdx := record.Schema().FieldIndices(expectedField.Name)
			if len(fieldIdx) == 0 {
				// Field missing in record
				errors = append(errors, SchemaValidationError{
					FieldName:    expectedField.Name,
					ExpectedType: formatArrowType(expectedField.Type),
					ActualType:   "missing",
					RowIndex:     -1,
					Value:        "",
				})
				continue
			}

			// Check if types match
			actualField := record.Schema().Field(fieldIdx[0])
			if expectedField.Type.ID() != actualField.Type.ID() {
				// Types don't match
				errors = append(errors, SchemaValidationError{
					FieldName:    expectedField.Name,
					ExpectedType: formatArrowType(expectedField.Type),
					ActualType:   formatArrowType(actualField.Type),
					RowIndex:     -1,
					Value:        "",
				})
				continue
			}
		}

		// Check for extra fields in record
		for _, actualField := range record.Schema().Fields() {
			fieldIdx := v.schema.FieldIndices(actualField.Name)
			if len(fieldIdx) == 0 {
				// Extra field in record
				errors = append(errors, SchemaValidationError{
					FieldName:    actualField.Name,
					ExpectedType: "not in schema",
					ActualType:   formatArrowType(actualField.Type),
					RowIndex:     -1,
					Value:        "",
				})
			}
		}
	}

	// Validate data values
	for i := 0; i < int(record.NumRows()); i++ {
		for _, field := range v.schema.Fields() {
			// Find field in record
			fieldIdx := record.Schema().FieldIndices(field.Name)
			if len(fieldIdx) == 0 {
				continue // Already reported missing field
			}

			// Get column
			col := record.Column(fieldIdx[0])
			if col.IsNull(i) {
				// Skip null values
				continue
			}

			// Validate value based on type
			if err := validateValue(col, i, field.Type); err != nil {
				errors = append(errors, SchemaValidationError{
					FieldName:    field.Name,
					ExpectedType: formatArrowType(field.Type),
					ActualType:   "invalid value",
					RowIndex:     int64(i),
					Value:        fmt.Sprintf("%v", getValueAsString(col, i)),
				})
			}
		}
	}

	return errors
}

// validateValue validates a value against its expected type
func validateValue(col arrow.Array, rowIdx int, expectedType arrow.DataType) error {
	switch expectedType.ID() {
	case arrow.INT32:
		// For int32, just check that we can get the value
		if _, ok := col.(*array.Int32); !ok {
			return fmt.Errorf("expected Int32 array")
		}
	case arrow.FLOAT64:
		// For float64, check for NaN or Inf
		if f64, ok := col.(*array.Float64); ok {
			val := f64.Value(rowIdx)
			if isNaN(val) || isInf(val) {
				return fmt.Errorf("invalid float value")
			}
		} else {
			return fmt.Errorf("expected Float64 array")
		}
	case arrow.STRING:
		// For string, check if it's valid UTF-8
		if str, ok := col.(*array.String); ok {
			val := str.Value(rowIdx)
			if !isValidUTF8(val) {
				return fmt.Errorf("invalid UTF-8 string")
			}
		} else {
			return fmt.Errorf("expected String array")
		}
	}
	return nil
}

// isNaN checks if a float64 is NaN
func isNaN(f float64) bool {
	return f != f
}

// isInf checks if a float64 is Inf
func isInf(f float64) bool {
	return f == math.Inf(1) || f == math.Inf(-1)
}

// isValidUTF8 checks if a string is valid UTF-8
func isValidUTF8(s string) bool {
	return strings.ToValidUTF8(s, "") == s
}

// getValueAsString returns a string representation of a value
func getValueAsString(col arrow.Array, rowIdx int) string {
	switch col := col.(type) {
	case *array.Int32:
		return fmt.Sprintf("%d", col.Value(rowIdx))
	case *array.Float64:
		return fmt.Sprintf("%f", col.Value(rowIdx))
	case *array.String:
		return col.Value(rowIdx)
	case *array.Boolean:
		return fmt.Sprintf("%t", col.Value(rowIdx))
	default:
		return "unknown"
	}
}

// FormatValidationErrors returns a human-readable representation of validation errors
func FormatValidationErrors(errors []SchemaValidationError) string {
	var sb strings.Builder

	if len(errors) == 0 {
		return "No schema validation errors."
	}

	sb.WriteString(fmt.Sprintf("Found %d schema validation errors:\n", len(errors)))

	for i, err := range errors {
		if err.RowIndex >= 0 {
			sb.WriteString(fmt.Sprintf("%d. Field '%s' at row %d: Expected %s, got %s (value: %s)\n",
				i+1, err.FieldName, err.RowIndex, err.ExpectedType, err.ActualType, err.Value))
		} else {
			sb.WriteString(fmt.Sprintf("%d. Field '%s': Expected %s, got %s\n",
				i+1, err.FieldName, err.ExpectedType, err.ActualType))
		}
	}

	return sb.String()
}
