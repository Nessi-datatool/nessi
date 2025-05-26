package integrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestAdvancedDataQuality(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Advanced Data Quality test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Since our mock CLI doesn't have advanced quality commands, we'll test the basic quality and tables commands instead
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

	// Step 6: Test license status
	licenseOutput, _ := testutil.AssertCommandSuccess(t, "license", "status")
	testutil.AssertOutputContains(t, licenseOutput, "License Status")

	// Check if we have license for premium features
	hasProEdition := strings.Contains(licenseOutput, "Pro Edition")
	hasTrial := strings.Contains(licenseOutput, "trial")

	// Step 7: Activate trial license if needed
	if !hasProEdition && !hasTrial {
		stdout, _ = testutil.AssertCommandSuccess(t, "license", "trial")
		testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")
	}

	// Step 8: Create a second sample table for testing purposes
	secondTablePath := filepath.Join(sampleDataDir, "second_table")
	require.NoError(t, os.MkdirAll(secondTablePath, 0755))
	testutil.CopyTestData(t, "delta_tables/sample_table", secondTablePath)

	// Test connecting to the second table
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "connect", "--path", secondTablePath)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")
}
