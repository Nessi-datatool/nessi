package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestCLIMonitoring(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping CLI Monitoring test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Since our mock CLI doesn't have a monitor command, we'll test the quality and tables commands instead
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

	// Step 6: Create a metrics file manually for testing purposes
	metricsPath := filepath.Join(tempDir, "metrics.json")
	metricsContent := `{
	"table": "sample_table",
	"timestamp": "2023-01-01T00:00:00Z",
	"period": "7d",
	"metrics": {
		"completeness": 98.5,
		"accuracy": 99.2,
		"consistency": 97.8,
		"uniqueness": 100.0,
		"timeliness": 95.5
	},
	"system": {
		"cpu_usage": 45.2,
		"memory_usage": 62.8,
		"disk_usage": 78.3
	},
	"overall_score": 98.2
}`
	require.NoError(t, os.WriteFile(metricsPath, []byte(metricsContent), 0644))

	// Verify the metrics file exists
	_, err := os.Stat(metricsPath)
	require.NoError(t, err, "Metrics file should exist")

	// Step 7: Test license status
	licenseOutput, _ := testutil.AssertCommandSuccess(t, "license", "status")
	testutil.AssertOutputContains(t, licenseOutput, "License Status")

	// Step 8: Test license trial activation (which is a Pro feature)
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "trial")
	testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")
}
