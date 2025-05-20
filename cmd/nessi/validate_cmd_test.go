package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidationResult(t *testing.T) {
	// Create a validation result
	result := &ValidationResult{
		TablePath:    "/path/to/table",
		Timestamp:    time.Now(),
		Valid:        true,
		TotalChecks:  10,
		PassedChecks: 9,
		FailedChecks: 1,
		Details:      []string{"Check 1 passed", "Check 2 failed"},
		Duration:     "1.5s",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(result)
	require.NoError(t, err)

	// Deserialize from JSON
	var deserializedResult ValidationResult
	err = json.Unmarshal(jsonData, &deserializedResult)
	require.NoError(t, err)

	// Verify deserialized result
	assert.Equal(t, result.TablePath, deserializedResult.TablePath)
	assert.Equal(t, result.Timestamp.Unix(), deserializedResult.Timestamp.Unix())
	assert.Equal(t, result.Valid, deserializedResult.Valid)
	assert.Equal(t, result.TotalChecks, deserializedResult.TotalChecks)
	assert.Equal(t, result.PassedChecks, deserializedResult.PassedChecks)
	assert.Equal(t, result.FailedChecks, deserializedResult.FailedChecks)
	assert.Equal(t, result.Details, deserializedResult.Details)
	assert.Equal(t, result.Duration, deserializedResult.Duration)
}

func TestDeltaTableValidation(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "validate_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta table structure
	tablePath := filepath.Join(tempDir, "test_table")
	err = os.Mkdir(tablePath, 0755)
	require.NoError(t, err)

	// Test invalid Delta table (no _delta_log directory)
	isDelta, err := isValidDeltaTable(tablePath)
	require.NoError(t, err)
	assert.False(t, isDelta)

	// Create _delta_log directory
	deltaLogPath := filepath.Join(tablePath, "_delta_log")
	err = os.Mkdir(deltaLogPath, 0755)
	require.NoError(t, err)

	// Test valid Delta table
	isDelta, err = isValidDeltaTable(tablePath)
	require.NoError(t, err)
	assert.True(t, isDelta)

	// Test schema validation
	hasSchema, err := hasValidSchema(tablePath)
	require.NoError(t, err)
	assert.False(t, hasSchema) // No transaction log files yet

	// Create a mock transaction log file
	logFile := filepath.Join(deltaLogPath, "00000000000000000000.json")
	err = os.WriteFile(logFile, []byte("{}"), 0644)
	require.NoError(t, err)

	// Test schema validation again
	hasSchema, err = hasValidSchema(tablePath)
	require.NoError(t, err)
	assert.True(t, hasSchema)

	// Test data files validation
	hasData, err := hasDataFiles(tablePath)
	require.NoError(t, err)
	assert.False(t, hasData) // No parquet files yet

	// Create a mock parquet file
	parquetFile := filepath.Join(tablePath, "part-00000-123456.parquet")
	err = os.WriteFile(parquetFile, []byte("mock parquet data"), 0644)
	require.NoError(t, err)

	// Test data files validation again
	hasData, err = hasDataFiles(tablePath)
	require.NoError(t, err)
	assert.True(t, hasData)

	// Test transaction log validation
	hasLog, err := hasTransactionLog(tablePath)
	require.NoError(t, err)
	assert.True(t, hasLog)
}

func TestOutputValidationResult(t *testing.T) {
	// Create a validation result
	result := &ValidationResult{
		TablePath:    "/path/to/table",
		Timestamp:    time.Now(),
		Valid:        true,
		TotalChecks:  10,
		PassedChecks: 9,
		FailedChecks: 1,
		Details:      []string{"Check 1 passed", "Check 2 failed"},
		Duration:     "1.5s",
	}

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "output_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test JSON output
	jsonFile := filepath.Join(tempDir, "result.json")
	validateFormat = "json"
	validateOutputFile = jsonFile
	validateVerbose = true

	// Capture stdout
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Output the result
	outputValidationResult(result)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Verify JSON file was created
	assert.FileExists(t, jsonFile)

	// Read the JSON file
	jsonData, err := os.ReadFile(jsonFile)
	require.NoError(t, err)

	// Deserialize from JSON
	var deserializedResult ValidationResult
	err = json.Unmarshal(jsonData, &deserializedResult)
	require.NoError(t, err)

	// Verify deserialized result
	assert.Equal(t, result.TablePath, deserializedResult.TablePath)
	assert.Equal(t, result.Timestamp.Unix(), deserializedResult.Timestamp.Unix())
	assert.Equal(t, result.Valid, deserializedResult.Valid)
	assert.Equal(t, result.TotalChecks, deserializedResult.TotalChecks)
	assert.Equal(t, result.PassedChecks, deserializedResult.PassedChecks)
	assert.Equal(t, result.FailedChecks, deserializedResult.FailedChecks)
	assert.Equal(t, result.Details, deserializedResult.Details)
	assert.Equal(t, result.Duration, deserializedResult.Duration)

	// Test text output
	textFile := filepath.Join(tempDir, "result.txt")
	validateFormat = "text"
	validateOutputFile = textFile
	validateVerbose = true

	// Output the result
	outputValidationResult(result)

	// Verify text file was created
	assert.FileExists(t, textFile)

	// Read the text file
	textData, err := os.ReadFile(textFile)
	require.NoError(t, err)
	textOutput := string(textData)

	// Verify text output contains expected information
	assert.Contains(t, textOutput, result.TablePath)
	assert.Contains(t, textOutput, "Valid: true")
	assert.Contains(t, textOutput, "Checks: 10 total, 9 passed, 1 failed")
	assert.Contains(t, textOutput, "Duration: 1.5s")
	assert.Contains(t, textOutput, "Check 1 passed")
	assert.Contains(t, textOutput, "Check 2 failed")

	// Reset global variables
	validateFormat = "text"
	validateOutputFile = ""
	validateVerbose = false
}
