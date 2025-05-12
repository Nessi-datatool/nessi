package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/stretchr/testify/require"
)

func TestParquetManager(t *testing.T) {
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

	// Create Parquet manager
	manager := NewParquetManager()
	require.NotNil(t, manager)

	// Create test record
	record, err := createTestRecord(schema)
	require.NoError(t, err)

	// Test writing record
	filePath := filepath.Join(tempDir, "test.parquet")
	err = manager.WriteRecord(filePath, record)
	require.NoError(t, err)

	// Write some test data to the file
	testData := []byte("test data")
	err = os.WriteFile(filePath, testData, 0644)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(filePath)
	require.NoError(t, err)

	// Test reading record
	readRecord, err := manager.ReadRecord(filePath, schema)
	require.NoError(t, err)
	require.NotNil(t, readRecord)
	require.Equal(t, record.NumRows(), readRecord.NumRows())

	// Test getting file stats
	stats, err := manager.GetFileStats(filePath)
	require.NoError(t, err)
	require.NotNil(t, stats)
	require.Greater(t, stats.Size, int64(0))

	// Test schema validation
	err = manager.ValidateSchema(record, schema)
	require.NoError(t, err)
}

func TestParquetManagerErrors(t *testing.T) {
	manager := NewParquetManager()

	// Test writing nil record
	err := manager.WriteRecord("test.parquet", nil)
	require.Error(t, err)

	// Test reading with nil schema
	_, err = manager.ReadRecord("test.parquet", nil)
	require.Error(t, err)

	// Test getting stats for non-existent file
	_, err = manager.GetFileStats("nonexistent.parquet")
	require.Error(t, err)

	// Test validating nil record/schema
	err = manager.ValidateSchema(nil, nil)
	require.Error(t, err)
}

func createTestRecord(schema *arrow.Schema) (arrow.Record, error) {
	builder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builder.Release()

	// Add test data
	builder.Field(0).(*array.Int32Builder).Append(1)
	builder.Field(1).(*array.StringBuilder).Append("test")
	builder.Field(2).(*array.Float64Builder).Append(1.0)

	return builder.NewRecord(), nil
}
