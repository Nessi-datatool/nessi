package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestReportGeneration(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Report Generation test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Test table path
	tablePath := filepath.Join(sampleDataDir, "sample_table")

	// Step 1: Test connecting to the table
	stdout, _ := testutil.AssertCommandSuccess(t, "tables", "connect", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

	// Step 2: Generate HTML report
	htmlReportPath := filepath.Join(tempDir, "report.html")
	stdout, _ = testutil.AssertCommandSuccess(t, "report", "generate", "--table", tablePath, "--format", "html", "--output", htmlReportPath)
	testutil.AssertOutputContains(t, stdout, "Report generated successfully")

	// Verify the HTML report file exists
	_, err := os.Stat(htmlReportPath)
	require.NoError(t, err, "HTML report file should exist")

	// Step 3: Generate JSON report
	jsonReportPath := filepath.Join(tempDir, "report.json")
	stdout, _ = testutil.AssertCommandSuccess(t, "report", "generate", "--table", tablePath, "--format", "json", "--output", jsonReportPath)
	testutil.AssertOutputContains(t, stdout, "Report generated successfully")

	// Verify the JSON report file exists
	_, err = os.Stat(jsonReportPath)
	require.NoError(t, err, "JSON report file should exist")

	// Step 4: Generate CSV report
	csvReportPath := filepath.Join(tempDir, "report.csv")
	stdout, _ = testutil.AssertCommandSuccess(t, "report", "generate", "--table", tablePath, "--format", "csv", "--output", csvReportPath)
	testutil.AssertOutputContains(t, stdout, "Report generated successfully")

	// Verify the CSV report file exists
	_, err = os.Stat(csvReportPath)
	require.NoError(t, err, "CSV report file should exist")

	// Step 5: Generate Markdown report
	mdReportPath := filepath.Join(tempDir, "report.md")
	stdout, _ = testutil.AssertCommandSuccess(t, "report", "generate", "--table", tablePath, "--format", "md", "--output", mdReportPath)
	testutil.AssertOutputContains(t, stdout, "Report generated successfully")

	// Verify the Markdown report file exists
	_, err = os.Stat(mdReportPath)
	require.NoError(t, err, "Markdown report file should exist")

	// Step 6: Test with charts included
	chartsReportPath := filepath.Join(tempDir, "report_with_charts.html")
	stdout, _ = testutil.AssertCommandSuccess(t, "report", "generate", "--table", tablePath, "--format", "html", "--output", chartsReportPath, "--include-charts")
	testutil.AssertOutputContains(t, stdout, "Report generated successfully")

	// Verify the report with charts file exists
	_, err = os.Stat(chartsReportPath)
	require.NoError(t, err, "Report with charts file should exist")

	// Step 7: Test with invalid format (should fail)
	_, stderr := testutil.AssertCommandFailure(t, "report", "generate", "--table", tablePath, "--format", "invalid", "--output", filepath.Join(tempDir, "invalid.report"))
	testutil.AssertOutputContains(t, stderr, "Invalid report format")
}
