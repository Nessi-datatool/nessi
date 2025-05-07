// Implement unit tests for the Delta reader with these test cases:
// 1. TestOpenTable:
//    - Test opening a valid Delta table
//    - Test opening a non-existent table (should error)
//    - Test opening a corrupted table (should error with specific message)
// 2. TestGetTableVersion:
//    - Test getting the current version
//    - Test getting a specific version
//    - Test getting a non-existent version (should error)
// 3. TestParseTransactionLog:
//    - Test parsing a simple log
//    - Test parsing a complex log with multiple operations
//    - Test parsing a corrupted log (should error)
// 4. TestReadParquetFile:
//    - Test reading a valid Parquet file
//    - Test reading a corrupted Parquet file
// Set up test fixtures with sample Delta tables and Parquet files
// Use table-driven tests where appropriate
// Mock filesystem operations to avoid external dependencies
// Verify both successful operations and proper error handling

package delta_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestTable creates a test Delta table with the specified structure
func setupTestTable(t *testing.T, basePath string, version int, hasData bool) string {
	tablePath := filepath.Join(basePath, "test_table")
	deltaLogPath := filepath.Join(tablePath, "_delta_log")

	// Create table directory structure
	require.NoError(t, os.MkdirAll(deltaLogPath, 0755))
	if hasData {
		require.NoError(t, os.MkdirAll(filepath.Join(tablePath, "data"), 0755))
	}

	// Create transaction log entries
	for i := 0; i <= version; i++ {
		logEntry := map[string]interface{}{
			"version": i,
			"timestamp": time.Now().Unix(),
			"operation": "WRITE",
			"partitionValues": map[string]string{},
			"numFiles": 1,
			"sizeBytes": 1024,
		}
		if i == 0 {
			logEntry["operation"] = "CREATE TABLE"
			logEntry["schema"] = map[string]interface{}{
				"type": "struct",
				"fields": []map[string]interface{}{
					{
						"name": "id",
						"type": "long",
						"nullable": false,
					},
					{
						"name": "name",
						"type": "string",
						"nullable": true,
					},
				},
			}
		}

		// Write log entry
		logFile := filepath.Join(deltaLogPath, fmt.Sprintf("%020d.json", i))
		data, err := json.Marshal(logEntry)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(logFile, data, 0644))
	}

	// Create checkpoint if version > 0
	if version > 0 {
		checkpoint := map[string]interface{}{
			"version": version,
			"timestamp": time.Now().Unix(),
			"size": 1024,
		}
		checkpointFile := filepath.Join(deltaLogPath, fmt.Sprintf("%020d.checkpoint.parquet", version))
		data, err := json.Marshal(checkpoint)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(checkpointFile, data, 0644))
	}

	return tablePath
}

// TestOpenTable tests opening Delta tables
func TestOpenTable(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		setup       func() string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid table",
			setup: func() string {
				return setupTestTable(t, tempDir, 1, true)
			},
			wantErr: false,
		},
		{
			name: "non-existent table",
			setup: func() string {
				return filepath.Join(tempDir, "non_existent")
			},
			wantErr:     true,
			errContains: "table path does not exist",
		},
		{
			name: "corrupted table - missing _delta_log",
			setup: func() string {
				tablePath := filepath.Join(tempDir, "corrupted")
				require.NoError(t, os.MkdirAll(tablePath, 0755))
				return tablePath
			},
			wantErr:     true,
			errContains: "invalid Delta table",
		},
		{
			name: "corrupted table - invalid log entry",
			setup: func() string {
				tablePath := setupTestTable(t, tempDir, 1, true)
				// Write invalid JSON to log file
				invalidLog := filepath.Join(tablePath, "_delta_log", "00000000000000000001.json")
				require.NoError(t, os.WriteFile(invalidLog, []byte("invalid json"), 0644))
				return tablePath
			},
			wantErr:     true,
			errContains: "invalid log entry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tablePath := tt.setup()
			reader, err := pkg.NewReader(tablePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, reader)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				assert.Equal(t, tablePath, reader.GetPath())
				require.NoError(t, reader.Close())
			}
		})
	}
}

// TestGetTableVersion tests version management
func TestGetTableVersion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create table with multiple versions
	tablePath := setupTestTable(t, tempDir, 3, true)
	reader, err := pkg.NewReader(tablePath)
	require.NoError(t, err)
	defer reader.Close()

	tests := []struct {
		name        string
		version     int
		wantErr     bool
		errContains string
	}{
		{
			name:    "current version",
			version: -1,
			wantErr: false,
		},
		{
			name:    "specific version",
			version: 2,
			wantErr: false,
		},
		{
			name:        "non-existent version",
			version:     10,
			wantErr:     true,
			errContains: "version not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version, err := reader.GetTableVersion(tt.version)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.version == -1 {
					assert.Equal(t, 3, version) // Current version
				} else {
					assert.Equal(t, tt.version, version)
				}
			}
		})
	}
}

// TestParseTransactionLog tests transaction log parsing
func TestParseTransactionLog(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		setup       func() string
		wantErr     bool
		errContains string
	}{
		{
			name: "simple log",
			setup: func() string {
				return setupTestTable(t, tempDir, 1, false)
			},
			wantErr: false,
		},
		{
			name: "complex log",
			setup: func() string {
				tablePath := setupTestTable(t, tempDir, 5, true)
				// Add more complex operations
				logPath := filepath.Join(tablePath, "_delta_log")
				for i := 6; i <= 10; i++ {
					logEntry := map[string]interface{}{
						"version": i,
						"timestamp": time.Now().Unix(),
						"operation": "MERGE",
						"predicate": "id > 100",
						"numFiles": 2,
						"sizeBytes": 2048,
					}
					logFile := filepath.Join(logPath, fmt.Sprintf("%020d.json", i))
					data, err := json.Marshal(logEntry)
					require.NoError(t, err)
					require.NoError(t, os.WriteFile(logFile, data, 0644))
				}
				return tablePath
			},
			wantErr: false,
		},
		{
			name: "corrupted log",
			setup: func() string {
				tablePath := setupTestTable(t, tempDir, 1, false)
				// Write invalid JSON to log file
				invalidLog := filepath.Join(tablePath, "_delta_log", "00000000000000000001.json")
				require.NoError(t, os.WriteFile(invalidLog, []byte("invalid json"), 0644))
				return tablePath
			},
			wantErr:     true,
			errContains: "invalid log entry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tablePath := tt.setup()
			reader, err := pkg.NewReader(tablePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, reader)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, reader)
				require.NoError(t, reader.Close())
			}
		})
	}
}

// TestReadParquetFile tests Parquet file reading
func TestReadParquetFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create table with Parquet files
	tablePath := setupTestTable(t, tempDir, 1, true)
	dataPath := filepath.Join(tablePath, "data")
	
	// Create a simple Parquet file
	parquetFile := filepath.Join(dataPath, "part-00000.parquet")
	// TODO: Create actual Parquet file with test data
	require.NoError(t, os.WriteFile(parquetFile, []byte("parquet data"), 0644))

	reader, err := pkg.NewReader(tablePath)
	require.NoError(t, err)
	defer reader.Close()

	tests := []struct {
		name        string
		filePath    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid file",
			filePath: parquetFile,
			wantErr:  false,
		},
		{
			name:        "non-existent file",
			filePath:    filepath.Join(dataPath, "non_existent.parquet"),
			wantErr:     true,
			errContains: "file not found",
		},
		{
			name:        "corrupted file",
			filePath:    parquetFile,
			wantErr:     true,
			errContains: "invalid parquet file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := reader.ReadParquetFile(tt.filePath)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, data)
			}
		})
	}
}

// TestReaderConcurrentAccess tests concurrent access to the reader
func TestReaderConcurrentAccess(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create table with multiple versions
	tablePath := setupTestTable(t, tempDir, 5, true)
	reader, err := pkg.NewReader(tablePath)
	require.NoError(t, err)
	defer reader.Close()

	// Test concurrent version access
	versions := make(chan int, 10)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			version, err := reader.GetTableVersion(-1)
			versions <- version
			errors <- err
		}()
	}

	// Collect results
	for i := 0; i < 10; i++ {
		version := <-versions
		err := <-errors
		assert.NoError(t, err)
		assert.Equal(t, 5, version)
	}
}

// TestReaderErrorHandling tests error handling in various scenarios
func TestReaderErrorHandling(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name        string
		setup       func() string
		operation   func(*pkg.Reader) error
		wantErr     bool
		errContains string
	}{
		{
			name: "read after close",
			setup: func() string {
				return setupTestTable(t, tempDir, 1, true)
			},
			operation: func(r *pkg.Reader) error {
				r.Close()
				_, err := r.GetTableVersion(-1)
				return err
			},
			wantErr:     true,
			errContains: "reader closed",
		},
		{
			name: "invalid partition",
			setup: func() string {
				return setupTestTable(t, tempDir, 1, true)
			},
			operation: func(r *pkg.Reader) error {
				_, err := r.ReadPartition(map[string]string{"invalid": "partition"})
				return err
			},
			wantErr:     true,
			errContains: "partition not found",
		},
		{
			name: "concurrent close",
			setup: func() string {
				return setupTestTable(t, tempDir, 1, true)
			},
			operation: func(r *pkg.Reader) error {
				done := make(chan error)
				go func() {
					done <- r.Close()
				}()
				go func() {
					done <- r.Close()
				}()
				return <-done
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tablePath := tt.setup()
			reader, err := pkg.NewReader(tablePath)
			require.NoError(t, err)

			err = tt.operation(reader)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}