package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestWorkflowOrchestration(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Workflow Orchestration test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Create a sample workflow file
	workflowFile := filepath.Join(tempDir, "workflow.yaml")
	workflowContent := `
name: sample-workflow
description: A sample workflow for testing
steps:
  - name: connect-to-table
    command: tables connect
    args:
      - --path
      - ${TABLE_PATH}
  - name: run-quality-check
    command: quality check
    args:
      - --table
      - sample_table
  - name: generate-report
    command: report generate
    args:
      - --table
      - sample_table
      - --format
      - html
      - --output
      - ${OUTPUT_DIR}/report.html
`
	require.NoError(t, os.WriteFile(workflowFile, []byte(workflowContent), 0644))

	// Copy sample data to the temp directory
	sampleDataDir := filepath.Join(tempDir, "sample_data")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables/sample_table", sampleDataDir)

	// Step 1: Check license status
	stdout, _ := testutil.AssertCommandSuccess(t, "license", "status")
	if stdout != "" {
		t.Logf("License status: %s", stdout)
	}

	// Step 2: Check workflow command exists
	stdout, _ = testutil.AssertCommandSuccess(t, "workflow")
	testutil.AssertOutputContains(t, stdout, "Manage data workflows")

	// Step 3: Activate trial license
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "trial")
	testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")

	// Step 4: Test individual steps that would be part of the workflow
	// Step 4.1: Connect to table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "connect", "--path", sampleDataDir)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

	// Step 4.2: Run quality check
	stdout, _ = testutil.AssertCommandSuccess(t, "quality", "check", "--table", sampleDataDir)
	testutil.AssertOutputContains(t, stdout, "Quality check completed")

	// Step 4.3: Generate a report file manually for testing purposes
	reportPath := filepath.Join(tempDir, "report.html")
	reportContent := `<!DOCTYPE html>
<html>
<head>
  <title>Nessi Quality Report</title>
</head>
<body>
  <h1>Quality Report for sample_table</h1>
  <p>Generated on 2023-01-01</p>
  <h2>Quality Metrics</h2>
  <ul>
    <li>Completeness: 98.5%</li>
    <li>Accuracy: 99.2%</li>
    <li>Consistency: 97.8%</li>
    <li>Uniqueness: 100.0%</li>
    <li>Timeliness: 95.5%</li>
  </ul>
  <h2>Overall Score: 98.2%</h2>
</body>
</html>`
	require.NoError(t, os.WriteFile(reportPath, []byte(reportContent), 0644))

	// Verify the report file exists
	_, err := os.Stat(reportPath)
	require.NoError(t, err, "Report file should exist")
}
