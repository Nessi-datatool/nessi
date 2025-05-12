package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/stretchr/testify/require"
)

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
	record, err := createTestRecord(schema)
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
