package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestWorkflowIntegration(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Workflow Integration test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Test license trial activation for Pro features
	t.Run("License Trial Activation", func(t *testing.T) {
		// Activate trial license
		stdout, _ := testutil.AssertCommandSuccess(t, "license", "trial")
		testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")
		testutil.AssertOutputContains(t, stdout, "Workflow Orchestration")

		// Verify license status shows Pro features
		stdout, _ = testutil.AssertCommandSuccess(t, "license", "status")
		testutil.AssertOutputContains(t, stdout, "Pro Edition")
	})

	// Test tables commands with our sample data
	t.Run("Tables Commands", func(t *testing.T) {
		// Create a sample table directory
		sampleTablePath := filepath.Join(sampleDataDir, "sample_table")
		require.NoError(t, os.MkdirAll(filepath.Join(sampleTablePath, "_delta_log"), 0755))

		// Test connecting to the table
		stdout, _ := testutil.AssertCommandSuccess(t, "tables", "connect", "--path", sampleTablePath)
		testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

		// Test describing the table
		stdout, _ = testutil.AssertCommandSuccess(t, "tables", "describe", "--path", sampleTablePath)
		testutil.AssertOutputContains(t, stdout, "Schema:")

		// Test reading from the table
		stdout, _ = testutil.AssertCommandSuccess(t, "tables", "read", "--path", sampleTablePath)
		testutil.AssertOutputContains(t, stdout, "Data preview:")
	})

	// Test workflow command exists
	t.Run("Workflow Command Exists", func(t *testing.T) {
		// Just check if the workflow command exists
		stdout, _ := testutil.AssertCommandSuccess(t, "workflow")
		testutil.AssertOutputContains(t, stdout, "Manage data workflows")
		testutil.AssertOutputContains(t, stdout, "run")
	})
}
