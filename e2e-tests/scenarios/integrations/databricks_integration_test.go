package integrations

import (
	"os"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestDatabricksIntegrationComprehensive(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Comprehensive Databricks Integration test. Use --enable flag to run this test.")
	}

	// Step 1: Test successful connection with mock credentials
	stdout, _ := testutil.AssertCommandSuccess(t, "databricks", "connect",
		"--token", "mock_token",
		"--host", "mock_host.cloud.databricks.com")
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Databricks workspace")

	// Step 2: Test connection with workspace ID
	stdout, _ = testutil.AssertCommandSuccess(t, "databricks", "connect",
		"--token", "mock_token",
		"--host", "mock_host.cloud.databricks.com",
		"--workspace", "1234567890")
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Databricks workspace")

	// Step 3: Test authentication failure with invalid token
	_, stderr := testutil.AssertCommandFailure(t, "databricks", "connect",
		"--token", "invalid_token",
		"--host", "mock_host.cloud.databricks.com")
	testutil.AssertOutputContains(t, stderr, "Databricks authentication failed")
	testutil.AssertOutputContains(t, stderr, "Invalid token")

	// Step 4: Test rate limiting error
	_, stderr = testutil.AssertCommandFailure(t, "databricks", "connect",
		"--token", "mock_token",
		"--host", "rate_limit_host")
	testutil.AssertOutputContains(t, stderr, "Databricks rate limit exceeded")
	testutil.AssertOutputContains(t, stderr, "Please try again later")

	// Step 5: Test missing required parameters
	_, stderr = testutil.AssertCommandFailure(t, "databricks", "connect",
		"--token", "mock_token")
	testutil.AssertOutputContains(t, stderr, "required flag(s)")
	testutil.AssertOutputContains(t, stderr, "host")

	_, stderr = testutil.AssertCommandFailure(t, "databricks", "connect",
		"--host", "mock_host.cloud.databricks.com")
	testutil.AssertOutputContains(t, stderr, "required flag(s)")
	testutil.AssertOutputContains(t, stderr, "token")
}
