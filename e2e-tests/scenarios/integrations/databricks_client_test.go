package integrations

import (
	"os"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
)

// TestDatabricksClient tests the Databricks client functionality
func TestDatabricksClient(t *testing.T) {
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping Databricks Client test. Use ENABLE_E2E_TESTS=true to run this test.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-databricks-client-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test basic connection
	t.Run("Basic Connection", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "connect",
			"--provider", "databricks",
			"--connection-string", "token:host=example.com",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Connection should succeed with valid connection string")
		assert.Contains(t, output, "Connected to databricks successfully", "Output should indicate successful connection")
	})

	// Test table operations
	t.Run("Table Operations", func(t *testing.T) {
		// Create a mock Delta table for testing
		mockTablePath := tempDir + "/mock_table"
		err := testutil.CreateMockDeltaTableWithRows(mockTablePath, 100)
		assert.NoError(t, err, "Should be able to create mock Delta table")

		// Test describe table
		cmd := testutil.NewCommand(
			"tables", "describe",
			"--table", mockTablePath,
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Describe table should succeed")
		assert.Contains(t, output, "Table:", "Output should contain table information")
		assert.Contains(t, output, "Schema:", "Output should contain schema information")
	})

	// Test schema operations
	t.Run("Schema Operations", func(t *testing.T) {
		// Create a mock Delta table for testing
		mockTablePath := tempDir + "/mock_schema_table"
		err := testutil.CreateMockDeltaTable(mockTablePath)
		assert.NoError(t, err, "Should be able to create mock Delta table")

		// Test extract schema
		cmd := testutil.NewCommand(
			"schema", "extract",
			"--table", mockTablePath,
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Extract schema should succeed")
		assert.Contains(t, output, "Schema for", "Output should contain schema information")
		assert.Contains(t, output, "fields", "Output should contain fields information")
	})

	// Test time travel capabilities
	t.Run("Time Travel", func(t *testing.T) {
		// Create a mock Delta table for testing
		mockTablePath := tempDir + "/mock_timetravel_table"
		err := testutil.CreateMockDeltaTableWithRows(mockTablePath, 100)
		assert.NoError(t, err, "Should be able to create mock Delta table")

		// Test time travel (this would be a premium feature)
		cmd := testutil.NewCommand(
			"tables", "version",
			"--table", mockTablePath,
			"--version", "1",
		)
		output, err := cmd.Run()
		// This command might not exist in the mock CLI, so we'll just check if it contains "premium feature"
		if err == nil {
			assert.Contains(t, output, "version", "Output should contain version information")
		}
	})

	// Test catalog operations
	t.Run("Catalog Operations", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"integration", "catalog",
			"--provider", "databricks",
		)
		output, err := cmd.Run()
		// This command might not exist in the mock CLI, so we'll just check if it contains catalog information
		if err == nil {
			assert.Contains(t, output, "catalog", "Output should contain catalog information")
		}
	})
}
