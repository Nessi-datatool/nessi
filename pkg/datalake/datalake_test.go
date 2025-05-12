package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReader(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	assert.NotNil(t, reader)
	assert.False(t, reader.IsClosed())

	// Test closing the reader
	err = reader.Close()
	require.NoError(t, err)
	assert.True(t, reader.IsClosed())
}

func TestInitialize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	defer reader.Close()

	// Initialize the reader
	err = reader.Initialize()
	require.NoError(t, err)
}

func TestReadPartition(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a partition directory
	partitionDir := filepath.Join(tempDir, "partition1")
	err = os.MkdirAll(partitionDir, 0755)
	require.NoError(t, err)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	defer reader.Close()

	// Read the partition
	record, err := reader.ReadPartition("partition1")
	require.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, int64(1), record.NumRows())
}

func TestGetStats(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	defer reader.Close()

	// Get stats
	stats, err := reader.GetStats()
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, int64(100), stats.NumRecords)
	assert.Equal(t, int64(1024), stats.TotalSize)
}

func TestReadAllStructured(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	defer reader.Close()

	// Read all structured data
	data, err := reader.ReadAllStructured()
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 2)
	assert.Equal(t, int32(1), data[0]["id"])
	assert.Equal(t, "test", data[0]["name"])
	assert.Equal(t, float64(1.0), data[0]["value"])
	assert.Equal(t, int32(2), data[1]["id"])
	assert.Equal(t, "test2", data[1]["name"])
	assert.Equal(t, float64(2.0), data[1]["value"])
}

func TestGetSchema(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	defer reader.Close()

	// Test getting schema
	schema := reader.GetSchema()
	assert.NotNil(t, schema)
}

func TestConvertToStructuredData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	defer reader.Close()

	// Read a partition to get a record
	record, err := reader.ReadPartition("partition1")
	require.NoError(t, err)
	assert.NotNil(t, record)

	// Convert to structured data
	data, err := reader.ConvertToStructuredData(record)
	require.NoError(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 1)
	assert.Equal(t, int32(1), data[0]["id"])
	assert.Equal(t, "test", data[0]["name"])
	assert.Equal(t, float64(1.0), data[0]["value"])
}
