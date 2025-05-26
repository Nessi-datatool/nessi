package stresstests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestErrorHandling(t *testing.T) {
	// Skip this test for now since we're just setting up the framework
	// Test is now enabled

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Test 1: Invalid path error
	stdout, stderr := testutil.AssertCommandFailure(t, "tables", "connect", "--path", "/non/existent/path")
	testutil.AssertOutputContains(t, stderr, "N1") // Path error code should start with N1
	testutil.AssertOutputContains(t, stderr, "does not exist")

	// Test 2: Invalid schema error
	// Create an invalid schema file
	invalidSchemaFile := filepath.Join(tempDir, "invalid_schema.json")
	err := os.WriteFile(invalidSchemaFile, []byte("{\"invalid\": \"schema\"}"), 0644)
	require.NoError(t, err, "Failed to create invalid schema file")

	// Copy sample data to the temp directory
	sampleDataDir := filepath.Join(tempDir, "sample_data")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables/sample_table", sampleDataDir)

	// Validate with invalid schema
	stdout, stderr = testutil.AssertCommandFailure(t, "schema", "validate", "--path", sampleDataDir, "--schema", invalidSchemaFile)
	testutil.AssertOutputContains(t, stderr, "N3") // Schema error code should start with N3
	testutil.AssertOutputContains(t, stderr, "invalid schema")

	// Test 3: Authentication error (for Databricks)
	stdout, stderr = testutil.AssertCommandFailure(t, "databricks", "connect", "--token", "invalid_token", "--host", "invalid_host")
	testutil.AssertOutputContains(t, stderr, "N2") // Authentication error code should start with N2
	testutil.AssertOutputContains(t, stderr, "authentication failed")

	// Test 4: Rate limiting error simulation
	stdout, stderr = testutil.AssertCommandFailure(t, "databricks", "connect", "--token", "rate_limit_token", "--host", "rate_limit_host")
	testutil.AssertOutputContains(t, stderr, "N5") // Rate limiting error code should start with N5
	testutil.AssertOutputContains(t, stderr, "rate limit exceeded")

	// Test 5: Test retry mechanism
	// This should eventually succeed after retries
	stdout, stderr = testutil.AssertCommandSuccess(t, "cloud", "connect", "--provider", "aws", "--endpoint", "retry_endpoint", "--retry", "true")
	testutil.AssertOutputContains(t, stdout, "Successfully connected after retries")
	testutil.AssertOutputContains(t, stderr, "Retrying") // Should show retry messages in stderr
}
