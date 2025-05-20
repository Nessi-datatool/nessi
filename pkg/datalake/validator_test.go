package datalake

import (
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaValidator_ValidateRecord(t *testing.T) {
	// Create expected schema
	expectedSchema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int32},
		{Name: "name", Type: arrow.BinaryTypes.String},
		{Name: "value", Type: arrow.PrimitiveTypes.Float64},
	}, nil)

	// Create validator
	validator := NewSchemaValidator(expectedSchema)

	t.Run("Valid record", func(t *testing.T) {
		// Create a valid record
		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, expectedSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)
		builder.Field(2).(*array.Float64Builder).AppendValues([]float64{1.1, 2.2, 3.3}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Validate
		errors := validator.ValidateRecord(record)
		assert.Empty(t, errors, "Valid record should have no errors")
	})

	t.Run("Missing field", func(t *testing.T) {
		// Create a schema with missing field
		invalidSchema := arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			// Missing "value" field
		}, nil)

		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, invalidSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Validate
		errors := validator.ValidateRecord(record)
		assert.NotEmpty(t, errors, "Record with missing field should have errors")
		assert.Equal(t, "value", errors[0].FieldName, "Error should be for missing 'value' field")
		assert.Equal(t, "missing", errors[0].ActualType, "Actual type should be 'missing'")
	})

	t.Run("Wrong type", func(t *testing.T) {
		// Create a schema with wrong type
		invalidSchema := arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Int32}, // Should be Float64
		}, nil)

		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, invalidSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)
		builder.Field(2).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Validate
		errors := validator.ValidateRecord(record)
		assert.NotEmpty(t, errors, "Record with wrong type should have errors")

		// Find the type error
		var typeError *SchemaValidationError
		for _, err := range errors {
			if err.FieldName == "value" {
				typeError = &err
				break
			}
		}

		require.NotNil(t, typeError, "Should have error for 'value' field")
		assert.Equal(t, "float64", typeError.ExpectedType, "Expected type should be float64")
	})

	t.Run("Extra field", func(t *testing.T) {
		// Create a schema with extra field
		invalidSchema := arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
			{Name: "extra", Type: &arrow.BooleanType{}}, // Extra field
		}, nil)

		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, invalidSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)
		builder.Field(2).(*array.Float64Builder).AppendValues([]float64{1.1, 2.2, 3.3}, nil)
		builder.Field(3).(*array.BooleanBuilder).AppendValues([]bool{true, false, true}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Validate
		errors := validator.ValidateRecord(record)
		assert.NotEmpty(t, errors, "Record with extra field should have errors")

		// Find the extra field error
		var extraError *SchemaValidationError
		for _, err := range errors {
			if err.FieldName == "extra" {
				extraError = &err
				break
			}
		}

		require.NotNil(t, extraError, "Should have error for 'extra' field")
		assert.Equal(t, "not in schema", extraError.ExpectedType, "Expected type should be 'not in schema'")
	})
}

func TestFormatValidationErrors(t *testing.T) {
	errors := []SchemaValidationError{
		{
			FieldName:    "id",
			ExpectedType: "int32",
			ActualType:   "string",
			RowIndex:     -1,
			Value:        "",
		},
		{
			FieldName:    "value",
			ExpectedType: "float64",
			ActualType:   "invalid value",
			RowIndex:     5,
			Value:        "NaN",
		},
	}

	formatted := FormatValidationErrors(errors)
	assert.Contains(t, formatted, "Found 2 schema validation errors")
	assert.Contains(t, formatted, "Field 'id': Expected int32, got string")
	assert.Contains(t, formatted, "Field 'value' at row 5: Expected float64, got invalid value (value: NaN)")
}
