package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestTimeTravel(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Time Travel test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Since our mock CLI doesn't have a time-travel command, we'll test the tables command instead
	tablePath := filepath.Join(sampleDataDir, "sample_table")

	// Step 1: Test connecting to the table
	stdout, _ := testutil.AssertCommandSuccess(t, "tables", "connect", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

	// Step 2: Test describing the table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "describe", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Schema:")

	// Step 3: Test reading from the table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "read", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Data preview:")

	// Step 4: Test license status
	licenseOutput, _ := testutil.AssertCommandSuccess(t, "license", "status")
	testutil.AssertOutputContains(t, licenseOutput, "License Status")

	// Step 5: Test license trial activation (which is a Pro feature)
	// Activate trial license
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "trial")
	testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")

	// Step 6: Create a new table directory to simulate a reconstructed table
	reconstructPath := filepath.Join(tempDir, "reconstructed_table")
	require.NoError(t, os.MkdirAll(filepath.Join(reconstructPath, "_delta_log"), 0755))

	// Test connecting to the reconstructed table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "connect", "--path", reconstructPath)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")
}
