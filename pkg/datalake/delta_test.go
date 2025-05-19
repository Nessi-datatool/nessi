package datalake

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeltaFormatHandler_IsDeltaTable(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table structure
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))

	// Create _delta_log directory
	deltaLogDir := filepath.Join(deltaTableDir, "_delta_log")
	require.NoError(t, os.MkdirAll(deltaLogDir, 0755))

	// Create a mock transaction log file
	txnLogFile := filepath.Join(deltaLogDir, "00000000000000000000.json")
	require.NoError(t, os.WriteFile(txnLogFile, []byte("{}"), 0644))

	// Create a non-Delta directory
	nonDeltaDir := filepath.Join(tempDir, "non-delta")
	require.NoError(t, os.MkdirAll(nonDeltaDir, 0755))

	// Test cases
	handler := NewDeltaFormatHandler()

	// Test with a valid Delta table
	assert.True(t, handler.IsDeltaTable(deltaTableDir), "Should recognize a valid Delta table")

	// Test with a non-Delta directory
	assert.False(t, handler.IsDeltaTable(nonDeltaDir), "Should not recognize a non-Delta directory as a Delta table")

	// Test with a non-existent directory
	assert.False(t, handler.IsDeltaTable(filepath.Join(tempDir, "nonexistent")), "Should not recognize a non-existent directory as a Delta table")
}

func TestGetDeltaTableMetadata(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-metadata-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table structure
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))

	// Create _delta_log directory
	deltaLogDir := filepath.Join(deltaTableDir, "_delta_log")
	require.NoError(t, os.MkdirAll(deltaLogDir, 0755))

	// Create a mock transaction log file
	txnLogFile := filepath.Join(deltaLogDir, "00000000000000000000.json")
	require.NoError(t, os.WriteFile(txnLogFile, []byte("{}"), 0644))

	// Test getting metadata from a Delta table
	metadata, err := GetDeltaTableMetadata(context.Background(), deltaTableDir)
	require.NoError(t, err)
	assert.Equal(t, "delta-table", metadata.Info.Name)
	assert.Equal(t, "delta", metadata.Schema.Format)
	assert.Equal(t, deltaTableDir, metadata.Info.Location)
}

func TestIsDeltaLakeFormat(t *testing.T) {
	// Test various format strings
	assert.True(t, IsDeltaLakeFormat("delta"), "Should recognize 'delta' as Delta Lake format")
	assert.True(t, IsDeltaLakeFormat("DELTA"), "Should recognize 'DELTA' as Delta Lake format (case-insensitive)")
	assert.True(t, IsDeltaLakeFormat("deltalake"), "Should recognize 'deltalake' as Delta Lake format")
	assert.True(t, IsDeltaLakeFormat("DeltaLake"), "Should recognize 'DeltaLake' as Delta Lake format (case-insensitive)")

	assert.False(t, IsDeltaLakeFormat("parquet"), "Should not recognize 'parquet' as Delta Lake format")
	assert.False(t, IsDeltaLakeFormat("csv"), "Should not recognize 'csv' as Delta Lake format")
	assert.False(t, IsDeltaLakeFormat(""), "Should not recognize empty string as Delta Lake format")
}
