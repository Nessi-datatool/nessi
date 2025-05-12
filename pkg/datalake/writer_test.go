package datalake

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestRecordForWriter creates a test Arrow record with the given schema
func createTestRecordForWriter(schema *arrow.Schema) (arrow.Record, error) {
	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	// Add data based on schema fields
	for i, field := range schema.Fields() {
		switch field.Type.ID() {
		case arrow.INT32:
			builder.Field(i).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		case arrow.STRING:
			builder.Field(i).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)
		case arrow.FLOAT64:
			builder.Field(i).(*array.Float64Builder).AppendValues([]float64{1.1, 2.2, 3.3}, nil)
		default:
			return nil, fmt.Errorf("unsupported field type: %s", field.Type)
		}
	}

	return builder.NewRecord(), nil
}

func TestWriter(t *testing.T) {
	// Create temp dir for testing
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create schema for testing
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	// Create writer
	writer, err := NewWriter(tempDir, schema)
	require.NoError(t, err)
	require.NotNil(t, writer)

	// Test initialization
	err = writer.Initialize()
	require.NoError(t, err)

	// Verify directories were created
	_, err = os.Stat(tempDir)
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "_delta_log"))
	require.NoError(t, err)

	// Test writing partition
	record, err := createTestRecordForWriter(schema)
	require.NoError(t, err)

	err = writer.WritePartition("test", record)
	require.NoError(t, err)

	// Verify partition was created
	_, err = os.Stat(filepath.Join(tempDir, "test"))
	require.NoError(t, err)

	// Test commit
	err = writer.Commit()
	require.NoError(t, err)

	// Verify transaction log was created
	files, err := os.ReadDir(filepath.Join(tempDir, "_delta_log"))
	require.NoError(t, err)
	require.Len(t, files, 1)

	// Test closing
	err = writer.Close()
	require.NoError(t, err)
	require.True(t, writer.IsClosed())
}

func TestSchemaValidationOnWrite(t *testing.T) {
	// Create temp dir for testing
	tempDir, err := os.MkdirTemp("", "delta-schema-validation-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create schema for testing
	expectedSchema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	// Create writer
	writer, err := NewWriter(tempDir, expectedSchema)
	require.NoError(t, err)
	require.NotNil(t, writer)

	// Initialize writer
	err = writer.Initialize()
	require.NoError(t, err)

	t.Run("Valid schema", func(t *testing.T) {
		// Create valid record
		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, expectedSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)
		builder.Field(2).(*array.Float64Builder).AppendValues([]float64{1.1, 2.2, 3.3}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Write should succeed
		err := writer.WritePartition("valid", record)
		require.NoError(t, err)
	})

	t.Run("Invalid schema - wrong type", func(t *testing.T) {
		// Create schema with wrong type
		invalidSchema := arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int32},
				{Name: "name", Type: arrow.BinaryTypes.String},
				{Name: "value", Type: arrow.PrimitiveTypes.Int32}, // Should be Float64
			},
			nil,
		)

		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, invalidSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)
		builder.Field(2).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Write should fail with schema validation error
		err := writer.WritePartition("invalid_type", record)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "schema validation failed")
		assert.Contains(t, err.Error(), "value")
	})

	t.Run("Invalid schema - missing field", func(t *testing.T) {
		// Create schema with missing field
		invalidSchema := arrow.NewSchema(
			[]arrow.Field{
				{Name: "id", Type: arrow.PrimitiveTypes.Int32},
				{Name: "name", Type: arrow.BinaryTypes.String},
				// Missing "value" field
			},
			nil,
		)

		pool := memory.NewGoAllocator()
		builder := array.NewRecordBuilder(pool, invalidSchema)
		defer builder.Release()

		// Add data
		builder.Field(0).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		builder.Field(1).(*array.StringBuilder).AppendValues([]string{"a", "b", "c"}, nil)

		record := builder.NewRecord()
		defer record.Release()

		// Write should fail with schema validation error
		err := writer.WritePartition("invalid_missing", record)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "schema validation failed")
		assert.Contains(t, err.Error(), "value")
	})
}

func TestWriterErrors(t *testing.T) {
	// Test nil schema
	writer, err := NewWriter("test", nil)
	require.Error(t, err)
	require.Nil(t, writer)

	// Create valid writer for testing errors
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
		},
		nil,
	)
	writer, err = NewWriter("test", schema)
	require.NoError(t, err)

	// Test writing to closed writer
	writer.Close()
	err = writer.WritePartition("test", nil)
	require.Error(t, err)

	// Test writing nil record
	writer, err = NewWriter("test", schema)
	require.NoError(t, err)
	err = writer.WritePartition("test", nil)
	require.Error(t, err)
}
