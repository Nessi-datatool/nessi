package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
)

// TestDatabricksErrorHandling tests comprehensive error handling for Databricks integration
func TestDatabricksErrorHandling(t *testing.T) {
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping Databricks Error Handling test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-databricks-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test authentication errors
	t.Run("Authentication Errors", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "invalid_token:host=example.com",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with authentication error")
		assert.Contains(t, output, "invalid", "Output should indicate invalid authentication")
	})

	// Test invalid workspace ID
	t.Run("Invalid Workspace ID", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "token:invalid_workspace_id",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with invalid workspace ID")
		assert.Contains(t, output, "invalid", "Output should indicate invalid workspace ID")
	})

	// Test resource not found errors
	t.Run("Resource Not Found", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"tables", "describe",
			"--provider", "databricks",
			"--catalog", "non_existent_catalog",
			"--schema", "non_existent_schema",
			"--table", "non_existent_table",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with resource not found error")
		assert.Contains(t, output, "not found", "Output should indicate resource not found")
	})

	// Test rate limiting errors
	t.Run("Rate Limiting", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "rate_limited:host=example.com",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with rate limiting error")
		assert.Contains(t, output, "rate", "Output should indicate rate limiting")
	})

	// Test server errors
	t.Run("Server Errors", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "server_error:host=example.com",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with server error")
		assert.Contains(t, output, "server", "Output should indicate server error")
	})

	// Test network errors
	t.Run("Network Errors", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "network_error:host=nonexistent.example.com",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with network error")
		assert.Contains(t, output, "network", "Output should indicate network error")
	})

	// Test malformed JSON response
	t.Run("Malformed JSON Response", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "malformed_json:host=example.com",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with malformed JSON error")
		assert.Contains(t, output, "JSON", "Output should indicate JSON parsing error")
	})

	// Test Delta Lake integration errors
	t.Run("Delta Lake Integration Errors", func(t *testing.T) {
		// Create an invalid Delta table for testing
		invalidDeltaPath := filepath.Join(tempDir, "invalid_delta")
		err := os.MkdirAll(filepath.Join(invalidDeltaPath, "_delta_log"), 0755)
		assert.NoError(t, err, "Should be able to create invalid Delta directory")

		cmd := testutil.NewCommand(
			"tables", "describe",
			"--provider", "databricks",
			"--table", invalidDeltaPath,
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with invalid Delta table")
		assert.Contains(t, output, "not a valid Delta table", "Output should indicate invalid Delta table")
	})

	// Test schema validation errors
	t.Run("Schema Validation Errors", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"schema", "validate",
			"--provider", "databricks",
			"--catalog", "test_catalog",
			"--schema", "test_schema",
			"--table", "test_table",
			"--schema-file", "nonexistent_schema.json",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with schema validation error")
		assert.Contains(t, output, "schema", "Output should indicate schema validation error")
	})

	// Test recovery from corrupted Delta tables
	t.Run("Recovery from Corrupted Delta Tables", func(t *testing.T) {
		// Create a corrupted Delta table for testing
		corruptedDeltaPath := filepath.Join(tempDir, "corrupted_delta")
		err := os.MkdirAll(filepath.Join(corruptedDeltaPath, "_delta_log"), 0755)
		assert.NoError(t, err, "Should be able to create corrupted Delta directory")

		// Create a corrupted transaction log file
		corruptedLogFile := filepath.Join(corruptedDeltaPath, "_delta_log", "00000000000000000000.json")
		err = os.WriteFile(corruptedLogFile, []byte("This is not valid JSON"), 0644)
		assert.NoError(t, err, "Should be able to create corrupted log file")

		cmd := testutil.NewCommand(
			"tables", "repair",
			"--table", corruptedDeltaPath,
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with corrupted Delta table")
		assert.Contains(t, output, "corrupted", "Output should indicate corrupted Delta table")
	})
}
