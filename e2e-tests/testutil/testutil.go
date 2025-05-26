package testutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "nessi-test-")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tempDir
}

// CopyTestData copies test data from the testdata directory
func CopyTestData(t *testing.T, sourcePath, destPath string) {
	cmd := exec.Command("cp", "-r", sourcePath, destPath)
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to copy test data: %v", err)
	}
}

// AssertCommandSuccess asserts that a command succeeds
func AssertCommandSuccess(t *testing.T, cmd *exec.Cmd) string {
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, outputStr)
	}
	return outputStr
}

// AssertCommandFailure asserts that a command fails
func AssertCommandFailure(t *testing.T, cmd *exec.Cmd) string {
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	if err == nil {
		t.Fatalf("Command succeeded when it should have failed\nOutput: %s", outputStr)
	}
	return outputStr
}

// AssertOutputContains asserts that output contains expected text
func AssertOutputContains(t *testing.T, output, expected string) {
	if !strings.Contains(output, expected) {
		t.Fatalf("Output does not contain expected text\nExpected: %s\nOutput: %s", expected, output)
	}
}

// CreateMockDeltaTable creates a mock Delta Lake table for testing
func CreateMockDeltaTable(path string) error {
	// Create the directory structure
	if err := os.MkdirAll(filepath.Join(path, "_delta_log"), 0755); err != nil {
		return err
	}

	// Create a basic transaction log file
	transactionLog := map[string]interface{}{
		"commitInfo": map[string]interface{}{
			"timestamp": 1620000000000,
			"operation": "CREATE TABLE",
			"operationParameters": map[string]interface{}{
				"isManaged": true,
				"description": "Mock Delta table for testing",
			},
		},
		"protocol": map[string]interface{}{
			"minReaderVersion": 1,
			"minWriterVersion": 2,
		},
		"metaData": map[string]interface{}{
			"id": "mock-delta-table",
			"format": map[string]interface{}{
				"provider": "parquet",
				"options": map[string]interface{}{},
			},
			"schemaString": `{
				"type": "struct",
				"fields": [
					{"name": "id", "type": "integer", "nullable": false, "metadata": {}},
					{"name": "name", "type": "string", "nullable": true, "metadata": {}},
					{"name": "value", "type": "double", "nullable": true, "metadata": {}}
				]
			}`,
			"partitionColumns": []string{},
			"configuration": map[string]interface{}{},
			"createdTime": 1620000000000,
		},
	}

	logBytes, err := json.Marshal(transactionLog)
	if err != nil {
		return err
	}

	logPath := filepath.Join(path, "_delta_log", "00000000000000000000.json")
	if err := ioutil.WriteFile(logPath, logBytes, 0644); err != nil {
		return err
	}

	// Create a sample parquet file (empty for testing purposes)
	parquetPath := filepath.Join(path, "part-00000-mock.snappy.parquet")
	if err := ioutil.WriteFile(parquetPath, []byte("MOCK PARQUET FILE"), 0644); err != nil {
		return err
	}

	return nil
}

// CreateMockDeltaTableWithRows creates a mock Delta Lake table with a specified number of rows
func CreateMockDeltaTableWithRows(path string, rowCount int) error {
	// Create basic table structure
	if err := CreateMockDeltaTable(path); err != nil {
		return err
	}

	// Update metadata to reflect row count
	logPath := filepath.Join(path, "_delta_log", "00000000000000000000.json")
	logBytes, err := ioutil.ReadFile(logPath)
	if err != nil {
		return err
	}

	var transactionLog map[string]interface{}
	if err := json.Unmarshal(logBytes, &transactionLog); err != nil {
		return err
	}

	// Add stats to the transaction log
	stats := map[string]interface{}{
		"numRecords": rowCount,
		"minValues": map[string]interface{}{
			"id": 1,
			"value": 0.1,
		},
		"maxValues": map[string]interface{}{
			"id": rowCount,
			"value": float64(rowCount) * 1.5,
		},
		"nullCount": map[string]interface{}{
			"id": 0,
			"name": rowCount / 10,
			"value": rowCount / 20,
		},
	}

	// Add the stats to the transaction log
	if _, ok := transactionLog["add"]; !ok {
		transactionLog["add"] = map[string]interface{}{
			"path": "part-00000-mock.snappy.parquet",
			"size": 1024 * rowCount,
			"partitionValues": map[string]interface{}{},
			"modificationTime": 1620000000000,
			"dataChange": true,
			"stats": stats,
		}
	}

	// Write updated transaction log
	updatedLogBytes, err := json.Marshal(transactionLog)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(logPath, updatedLogBytes, 0644); err != nil {
		return err
	}

	return nil
}

// RunNessiCommand runs a Nessi CLI command
func RunNessiCommand(args ...string) *exec.Cmd {
	// Look for the Nessi binary in the PATH or in the project's bin directory
	nessiPath, err := exec.LookPath("nessi")
	if err != nil {
		// If not found in PATH, try the project's bin directory
		projectRoot := findProjectRoot()
		if projectRoot != "" {
			binPath := filepath.Join(projectRoot, "bin", "nessi")
			if _, err := os.Stat(binPath); err == nil {
				nessiPath = binPath
			} else {
				// For testing, use the mock CLI
				mockPath := filepath.Join(projectRoot, "e2e-tests", "mock", "mock-nessi")
				if _, err := os.Stat(mockPath); err == nil {
					nessiPath = mockPath
				}
			}
		}
	}

	if nessiPath == "" {
		nessiPath = "nessi" // Fall back to PATH lookup at runtime
	}

	return exec.Command(nessiPath, args...)
}

// findProjectRoot attempts to find the root of the Nessi project
func findProjectRoot() string {
	// Start from the current directory and go up until we find a go.mod file
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root directory without finding go.mod
			return ""
		}
		dir = parent
	}
}

// SetupTestEnvironment sets up the environment for a test
func SetupTestEnvironment(t *testing.T) (string, func()) {
	// Skip if E2E tests are not enabled
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping E2E test; set ENABLE_E2E_TESTS=true to run")
	}

	// Create a temporary directory for the test
	tempDir := CreateTempDir(t)

	// Return the temp directory and a cleanup function
	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

// CreateInvalidDeltaTable creates an invalid Delta Lake table for testing error handling
func CreateInvalidDeltaTable(path string) error {
	// Create the directory structure but with invalid content
	if err := os.MkdirAll(filepath.Join(path, "_delta_log"), 0755); err != nil {
		return err
	}

	// Create an invalid transaction log file
	invalidLog := []byte("This is not valid JSON")
	logPath := filepath.Join(path, "_delta_log", "00000000000000000000.json")
	if err := ioutil.WriteFile(logPath, invalidLog, 0644); err != nil {
		return err
	}

	return nil
}

// CreatePermissionErrorDir creates a directory with restricted permissions for testing
func CreatePermissionErrorDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// Make the directory read-only
	return os.Chmod(path, 0500)
}

// CreateMockConfigFile creates a mock configuration file for testing
func CreateMockConfigFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(path, []byte(content), 0644)
}

// CreateInvalidConfigFile creates an invalid configuration file for testing
func CreateInvalidConfigFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	invalidContent := `
	This is not valid YAML or JSON
	- item1
	  * invalid indentation
	`

	return ioutil.WriteFile(path, []byte(invalidContent), 0644)
}

// CreateMockDatabricksConfig creates a mock Databricks configuration for testing
func CreateMockDatabricksConfig(path string) error {
	config := map[string]interface{}{
		"host": "https://mock-databricks.cloud.databricks.com",
		"token": "mock-token",
		"catalog": "mock_catalog",
		"schema": "mock_schema",
		"warehouse_id": "mock_warehouse_id",
	}

	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(path, configBytes, 0644)
}

// RunCommandWithTimeout runs a command with a timeout
func RunCommandWithTimeout(cmd *exec.Cmd, timeoutSeconds int) (string, error) {
	// This is a simplified version for testing
	// In a real implementation, you would use context with timeout
	output, err := cmd.CombinedOutput()
	return string(output), err
}
