package datalake

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeltaTimeTravel_GetVersionHistory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-time-travel-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	
	// Create _delta_log directory
	deltaLogDir := filepath.Join(deltaTableDir, "_delta_log")
	require.NoError(t, os.MkdirAll(deltaLogDir, 0755))

	// Create mock transaction log files with different timestamps
	// Version 0
	txnLogFile0 := filepath.Join(deltaLogDir, "00000000000000000000.json")
	require.NoError(t, os.WriteFile(txnLogFile0, []byte(`{"commitInfo": {"timestamp": 1620000000000, "operation": "CREATE TABLE"}}`), 0644))
	
	// Sleep briefly to ensure file timestamps are different
	time.Sleep(10 * time.Millisecond)
	
	// Version 1
	txnLogFile1 := filepath.Join(deltaLogDir, "00000000000000000001.json")
	require.NoError(t, os.WriteFile(txnLogFile1, []byte(`{"commitInfo": {"timestamp": 1620001000000, "operation": "WRITE"}}`), 0644))
	
	time.Sleep(10 * time.Millisecond)
	
	// Version 2
	txnLogFile2 := filepath.Join(deltaLogDir, "00000000000000000002.json")
	require.NoError(t, os.WriteFile(txnLogFile2, []byte(`{"commitInfo": {"timestamp": 1620002000000, "operation": "WRITE"}}`), 0644))

	// Initialize Delta format handler and time travel
	handler := NewDeltaFormatHandler()
	timeTravel := NewDeltaTimeTravel(handler)

	// Test getting version history
	ctx := context.Background()
	versions, err := timeTravel.GetVersionHistory(ctx, deltaTableDir)
	require.NoError(t, err)
	
	// Verify version history
	require.Len(t, versions, 3)
	assert.Equal(t, int64(0), versions[0].Version)
	assert.Equal(t, int64(1), versions[1].Version)
	assert.Equal(t, int64(2), versions[2].Version)
	
	// Verify timestamps are in ascending order
	assert.True(t, versions[0].Timestamp.Before(versions[1].Timestamp))
	assert.True(t, versions[1].Timestamp.Before(versions[2].Timestamp))
}

func TestDeltaTimeTravel_ReadAsOfVersion(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-time-travel-version-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	
	// Create _delta_log directory
	deltaLogDir := filepath.Join(deltaTableDir, "_delta_log")
	require.NoError(t, os.MkdirAll(deltaLogDir, 0755))

	// Create mock transaction log files
	txnLogFile0 := filepath.Join(deltaLogDir, "00000000000000000000.json")
	require.NoError(t, os.WriteFile(txnLogFile0, []byte(`{}`), 0644))
	
	txnLogFile1 := filepath.Join(deltaLogDir, "00000000000000000001.json")
	require.NoError(t, os.WriteFile(txnLogFile1, []byte(`{}`), 0644))

	// Create a mock parquet file (we won't actually write parquet data)
	require.NoError(t, os.WriteFile(filepath.Join(deltaTableDir, "part-00000.parquet"), []byte("mock parquet data"), 0644))

	// Initialize Delta format handler and time travel
	handler := NewDeltaFormatHandler()
	timeTravel := NewDeltaTimeTravel(handler)

	// Test reading as of version 1
	_, err = timeTravel.ReadAsOfVersion(deltaTableDir, 1)
	// We don't check the actual data since our implementation is simplified
	// Just verify that the function runs without error
	assert.NoError(t, err)

	// Test reading a non-existent version
	_, err = timeTravel.ReadAsOfVersion(deltaTableDir, 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "version 5 not found")

	// Test reading from a non-Delta table
	nonDeltaDir := filepath.Join(tempDir, "non-delta")
	require.NoError(t, os.MkdirAll(nonDeltaDir, 0755))
	_, err = timeTravel.ReadAsOfVersion(nonDeltaDir, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a Delta Lake table")
}

func TestDeltaTimeTravel_ReadAsOfTimestamp(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-time-travel-timestamp-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	
	// Create _delta_log directory
	deltaLogDir := filepath.Join(deltaTableDir, "_delta_log")
	require.NoError(t, os.MkdirAll(deltaLogDir, 0755))

	// Create mock transaction log files
	txnLogFile0 := filepath.Join(deltaLogDir, "00000000000000000000.json")
	require.NoError(t, os.WriteFile(txnLogFile0, []byte(`{}`), 0644))

	// Create a mock parquet file (we won't actually write parquet data)
	require.NoError(t, os.WriteFile(filepath.Join(deltaTableDir, "part-00000.parquet"), []byte("mock parquet data"), 0644))

	// Initialize Delta format handler and time travel
	handler := NewDeltaFormatHandler()
	timeTravel := NewDeltaTimeTravel(handler)

	// Test reading as of timestamp
	timestamp := time.Now().Add(-1 * time.Hour)
	_, err = timeTravel.ReadAsOfTimestamp(deltaTableDir, timestamp)
	// We don't check the actual data since our implementation is simplified
	// Just verify that the function runs without error
	assert.NoError(t, err)

	// Test reading from a non-Delta table
	nonDeltaDir := filepath.Join(tempDir, "non-delta")
	require.NoError(t, os.MkdirAll(nonDeltaDir, 0755))
	_, err = timeTravel.ReadAsOfTimestamp(nonDeltaDir, timestamp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a Delta Lake table")
}
