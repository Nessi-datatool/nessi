package userjourneys

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestNewUserOnboarding(t *testing.T) {
	// Skip this test for now since we're just setting up the framework
	// Test is now enabled

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Step 1: Initialize configuration
	stdout, _ := testutil.AssertCommandSuccess(t, "config", "init", "--dir", tempDir)
	testutil.AssertOutputContains(t, stdout, "Configuration initialized successfully")

	// Verify config file was created
	configFile := filepath.Join(tempDir, ".nessi", "config.yaml")
	_, err := os.Stat(configFile)
	require.NoError(t, err, "Config file should exist")

	// Step 2: Connect to a sample Delta Lake table
	// First, copy sample data to the temp directory
	sampleDataDir := filepath.Join(tempDir, "sample_data")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables/sample_table", sampleDataDir)

	// Connect to the table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "connect", "--path", sampleDataDir)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

	// Step 3: Run a basic quality check
	stdout, _ = testutil.AssertCommandSuccess(t, "quality", "check", "--table", "sample_table")
	testutil.AssertOutputContains(t, stdout, "Quality check completed")

	// Step 4: Generate a report
	reportPath := filepath.Join(tempDir, "report.html")
	stdout, _ = testutil.AssertCommandSuccess(t, "report", "generate", "--table", "sample_table", "--format", "html", "--output", reportPath)
	testutil.AssertOutputContains(t, stdout, "Report generated successfully")

	// Verify report file was created
	_, err = os.Stat(reportPath)
	require.NoError(t, err, "Report file should exist")

	// Step 5: Check license status (should be community edition by default)
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "status")
	testutil.AssertOutputContains(t, stdout, "Community Edition")
}
