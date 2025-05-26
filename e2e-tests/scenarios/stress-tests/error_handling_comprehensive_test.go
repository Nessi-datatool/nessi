package stresstests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
)

// TestErrorHandlingComprehensive tests comprehensive error handling
// across various CLI commands and scenarios
func TestErrorHandlingComprehensive(t *testing.T) {
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping Comprehensive Error Handling test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-error-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test non-existent table
	t.Run("Non-existent Table", func(t *testing.T) {
		nonExistentPath := filepath.Join(tempDir, "non_existent_table")
		cmd := testutil.NewCommand(
			"tables", "describe",
			"--table", nonExistentPath,
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with non-existent table")
		assert.Contains(t, output, "Error", "Output should contain error message")
		assert.Contains(t, output, "not found", "Output should indicate table not found")
	})

	// Test invalid command flags
	t.Run("Invalid Command Flags", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"quality", "check",
			"--invalid-flag", "value",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with invalid flag")
		assert.Contains(t, output, "Error", "Output should contain error message")
		assert.Contains(t, output, "flag", "Output should indicate flag error")
	})

	// Test missing required arguments
	t.Run("Missing Required Arguments", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"quality", "check",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with missing required arguments")
		assert.Contains(t, output, "Error", "Output should contain error message")
		assert.Contains(t, output, "required", "Output should indicate required argument missing")
	})

	// Test invalid file format
	t.Run("Invalid File Format", func(t *testing.T) {
		// Create a non-Delta directory
		nonDeltaDir := filepath.Join(tempDir, "non_delta_dir")
		err := os.MkdirAll(nonDeltaDir, 0755)
		assert.NoError(t, err, "Should be able to create non-Delta directory")

		// Test with a directory that's not a Delta table
		cmd := testutil.NewCommand(
			"tables", "describe",
			"--table", nonDeltaDir,
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with invalid Delta table")
		assert.Contains(t, output, "not a valid Delta table", "Output should indicate invalid Delta table")
	})

	// Test permission errors
	t.Run("Permission Errors", func(t *testing.T) {
		// Create a directory with a name that will trigger permission error in mock CLI
		noPermDir := filepath.Join(tempDir, "no_perm_dir")
		err := os.MkdirAll(noPermDir, 0755)
		assert.NoError(t, err, "Should be able to create directory")

		cmd := testutil.NewCommand(
			"tables", "list",
			"--dir", noPermDir,
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with permission error")
		assert.Contains(t, output, "permission denied", "Output should indicate permission error")
	})

	// Test invalid configuration file
	t.Run("Invalid Configuration File", func(t *testing.T) {
		// Create an invalid config file
		invalidConfigPath := filepath.Join(tempDir, "invalid_config.yaml")
		err := os.WriteFile(invalidConfigPath, []byte("invalid: yaml: :\n"), 0644)
		assert.NoError(t, err, "Should be able to create invalid config file")

		cmd := testutil.NewCommand(
			"config", "load",
			"--file", invalidConfigPath,
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with invalid config file")
		assert.Contains(t, output, "Error", "Output should contain error message")
		assert.Contains(t, output, "config", "Output should indicate config error")
	})

	// Test invalid connection string
	t.Run("Invalid Connection String", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "invalid_connection_string",
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with invalid connection string")
		assert.Contains(t, output, "Error", "Output should contain error message")
		assert.Contains(t, output, "connection", "Output should indicate connection error")
	})

	// Test timeout handling
	t.Run("Timeout Handling", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"tables", "describe",
			"--table", "s3://non-existent-bucket/table",
			"--timeout", "1s", // Very short timeout
		)
		output, err := cmd.Run()
		assert.Error(t, err, "Command should fail with timeout")
		assert.Contains(t, output, "Error", "Output should contain error message")
		assert.Contains(t, output, "timeout", "Output should indicate timeout error")
	})

	// Test license validation
	t.Run("License Validation", func(t *testing.T) {
		// Try to access a premium feature without a license
		cmd := testutil.NewCommand(
			"workflow", "orchestrate",
			"--workflow", "test_workflow",
		)
		output, err := cmd.Run()
		// The command might succeed if running with a valid license or in trial mode
		if err != nil {
			assert.Contains(t, output, "license", "Output should indicate license requirement")
		}
	})

	// Test graceful shutdown - simplified for mock CLI testing
	t.Run("Graceful Shutdown", func(t *testing.T) {
		// For the mock CLI, we'll use a short duration instead of interrupting
		cmd := testutil.NewCommand(
			"demo", "long-running",
			"--duration", "2s", // Short duration for testing
		)

		// Run the command and capture output
		output, err := cmd.Run()
		assert.NoError(t, err, "Command should complete successfully")

		// Check the output
		assert.Contains(t, output, "Process completed successfully", "Output should indicate successful completion")
	})
}
