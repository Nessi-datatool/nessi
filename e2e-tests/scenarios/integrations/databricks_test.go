package integrations

import (
	"os"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestDatabricksIntegration(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Databricks integration test. Use --enable flag to run this test.")
	}

	// Check if we have the necessary environment variables
	token := os.Getenv("DATABRICKS_TOKEN")
	host := os.Getenv("DATABRICKS_HOST")
	if token == "" || host == "" {
		t.Skip("Skipping Databricks integration test. Set DATABRICKS_TOKEN and DATABRICKS_HOST environment variables to run this test.")
	}

	// Step 1: Connect to Databricks
	stdout, _ := testutil.AssertCommandSuccess(t, "databricks", "connect", "--token", token, "--host", host)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Databricks workspace")

	// Step 2: List tables
	stdout, _ = testutil.AssertCommandSuccess(t, "databricks", "list-tables")
	testutil.AssertOutputContains(t, stdout, "Tables in Databricks workspace")

	// Step 3: Test error handling with invalid token
	stdout, stderr := testutil.AssertCommandFailure(t, "databricks", "connect", "--token", "invalid_token", "--host", host)
	testutil.AssertOutputContains(t, stderr, "authentication failed")

	// Step 4: Test license check for Databricks features
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "status")
	if stdout != "" {
		t.Logf("License status: %s", stdout)
	}
}
