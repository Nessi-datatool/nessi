package integrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestDeltaLakeIntegration(t *testing.T) {
	// Skip this test for now since we're just setting up the framework
	// Test is now enabled

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Step 1: Connect to a Delta Lake table
	tablePath := filepath.Join(sampleDataDir, "sample_table")
	stdout, _ := testutil.AssertCommandSuccess(t, "tables", "connect", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

	// Step 2: Describe the table structure
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "describe", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Schema")
	testutil.AssertOutputContains(t, stdout, "Partitioning")

	// Step 3: Read data from the table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "read", "--path", tablePath, "--limit", "10")
	testutil.AssertOutputContains(t, stdout, "Data preview")

	// Step 4: Test time travel capabilities (if available in the Pro edition)
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "status")
	if strings.Contains(stdout, "Pro Edition") {
		// Test time travel with version
		stdout, _ = testutil.AssertCommandSuccess(t, "tables", "read", "--path", tablePath, "--version", "0", "--limit", "5")
		testutil.AssertOutputContains(t, stdout, "Data preview (version 0)")

		// Test time travel with timestamp
		stdout, _ = testutil.AssertCommandSuccess(t, "tables", "read", "--path", tablePath, "--timestamp", "2023-01-01T00:00:00Z", "--limit", "5")
		testutil.AssertOutputContains(t, stdout, "Data preview (as of")
	} else {
		t.Log("Skipping time travel tests as they require Pro Edition")
	}

	// Step 5: Test schema validation
	schemaFile := filepath.Join(tempDir, "schema.json")
	stdout, _ = testutil.AssertCommandSuccess(t, "schema", "extract", "--path", tablePath, "--output", schemaFile)
	testutil.AssertOutputContains(t, stdout, "Schema extracted successfully")

	stdout, _ = testutil.AssertCommandSuccess(t, "schema", "validate", "--path", tablePath, "--schema", schemaFile)
	testutil.AssertOutputContains(t, stdout, "Schema validation successful")
}
