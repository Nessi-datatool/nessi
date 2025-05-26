package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestDatabricksClient(t *testing.T) {
	// Skip if E2E tests are not enabled
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping E2E test; set ENABLE_E2E_TESTS=true to run")
	}

	// Create a temporary directory for the test
	tempDir, cleanup := testutil.SetupTestEnvironment(t)
	defer cleanup()

	// Create test directories
	configDir := filepath.Join(tempDir, "config")
	tableDir := filepath.Join(tempDir, "tables")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	if err := os.MkdirAll(tableDir, 0755); err != nil {
		t.Fatalf("Failed to create tables directory: %v", err)
	}

	// Create mock Databricks config
	configPath := filepath.Join(configDir, "databricks.json")
	if err := testutil.CreateMockDatabricksConfig(configPath); err != nil {
		t.Fatalf("Failed to create mock Databricks config: %v", err)
	}

	// Create mock Delta tables
	table1Path := filepath.Join(tableDir, "table1")
	if err := testutil.CreateMockDeltaTable(table1Path); err != nil {
		t.Fatalf("Failed to create mock Delta table: %v", err)
	}

	// Run subtests
	t.Run("Basic_Connection", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("databricks", "connect", "--config", configPath)
		output := testutil.AssertCommandSuccess(t, cmd)
		testutil.AssertOutputContains(t, output, "Successfully connected to Databricks")
	})

	t.Run("Table_Operations", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("databricks", "tables", "--config", configPath)
		output := testutil.AssertCommandSuccess(t, cmd)
		testutil.AssertOutputContains(t, output, "Tables in mock_catalog.mock_schema")
	})

	t.Run("Schema_Operations", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("databricks", "schema", "describe", "--config", configPath, "--table", "mock_table")
		output := testutil.AssertCommandSuccess(t, cmd)
		testutil.AssertOutputContains(t, output, "Schema for mock_table")
	})

	t.Run("Time_Travel", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("databricks", "time-travel", "--config", configPath, "--table", "mock_table", "--version", "1")
		output := testutil.AssertCommandSuccess(t, cmd)
		testutil.AssertOutputContains(t, output, "Time travel to version 1")
	})

	t.Run("Catalog_Operations", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("databricks", "catalog", "list", "--config", configPath)
		output := testutil.AssertCommandSuccess(t, cmd)
		testutil.AssertOutputContains(t, output, "Available catalogs")
	})
}
