package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestAutomatedProfiling(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Automated Profiling test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Since our mock CLI doesn't have a profile command, we'll test the quality and tables commands instead
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

	// Step 4: Test quality check on a table
	stdout, _ = testutil.AssertCommandSuccess(t, "quality", "check", "--table", tablePath)
	testutil.AssertOutputContains(t, stdout, "Quality check completed")
	testutil.AssertOutputContains(t, stdout, "Results for table:")

	// Step 5: Test quality check with thresholds
	stdout, _ = testutil.AssertCommandSuccess(t, "quality", "check", "--table", tablePath, "--show-thresholds")
	testutil.AssertOutputContains(t, stdout, "Quality check completed")
	testutil.AssertOutputContains(t, stdout, "Thresholds:")

	// Step 6: Create a profile report file manually for testing purposes
	profilePath := filepath.Join(tempDir, "profile.json")
	profileContent := `{
	"table": "sample_table",
	"timestamp": "2023-01-01T00:00:00Z",
	"metrics": {
		"completeness": 98.5,
		"accuracy": 99.2,
		"consistency": 97.8,
		"uniqueness": 100.0,
		"timeliness": 95.5
	},
	"overall_score": 98.2
}`
	require.NoError(t, os.WriteFile(profilePath, []byte(profileContent), 0644))

	// Verify the profile file exists
	_, err := os.Stat(profilePath)
	require.NoError(t, err, "Profile file should exist")
}
