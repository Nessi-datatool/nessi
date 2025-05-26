package integrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestCloudIntegration(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Cloud Integration test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Check license status
	licenseOutput, _ := testutil.AssertCommandSuccess(t, "license", "status")
	hasProEdition := strings.Contains(licenseOutput, "Pro Edition")
	hasTrial := strings.Contains(licenseOutput, "trial")

	// Activate trial license if needed
	if !hasProEdition && !hasTrial {
		stdout, _ := testutil.AssertCommandSuccess(t, "license", "trial")
		testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")
	}

	// Test local cloud provider since we don't have actual cloud credentials in the test environment
	t.Run("Local Provider", func(t *testing.T) {
		// Step 1: Connect to local provider
		stdout, _ := testutil.AssertCommandSuccess(t, "cloud", "connect", "--provider", "local", "--endpoint", "/tmp/delta-tables")
		testutil.AssertOutputContains(t, stdout, "Successfully connected to local cloud storage")

		// Step 2: Test connecting to a table
		sampleDataDir := filepath.Join(tempDir, "delta_tables")
		require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
		testutil.CopyTestData(t, "delta_tables", sampleDataDir)

		tablePath := filepath.Join(sampleDataDir, "sample_table")
		stdout, _ = testutil.AssertCommandSuccess(t, "tables", "connect", "--path", tablePath)
		testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")
	})

	// Test AWS S3 integration with mock implementation
	t.Run("AWS S3", func(t *testing.T) {
		// Step 1: Connect to AWS S3
		stdout, _ := testutil.AssertCommandSuccess(t, "cloud", "connect", "--provider", "aws", "--endpoint", "s3.amazonaws.com", "--region", "us-west-2")
		testutil.AssertOutputContains(t, stdout, "Successfully connected to aws cloud storage")
	})

	// Test Azure Blob Storage integration with mock implementation
	t.Run("Azure Blob Storage", func(t *testing.T) {
		// Step 1: Connect to Azure Blob Storage
		stdout, _ := testutil.AssertCommandSuccess(t, "cloud", "connect", "--provider", "azure", "--endpoint", "account.blob.core.windows.net")
		testutil.AssertOutputContains(t, stdout, "Successfully connected to azure cloud storage")
	})

	// Test Google Cloud Storage integration with mock implementation
	t.Run("Google Cloud Storage", func(t *testing.T) {
		// Step 1: Connect to Google Cloud Storage
		stdout, _ := testutil.AssertCommandSuccess(t, "cloud", "connect", "--provider", "gcp", "--endpoint", "storage.googleapis.com")
		testutil.AssertOutputContains(t, stdout, "Successfully connected to gcp cloud storage")
	})

	// Test retry mechanism
	t.Run("Retry Mechanism", func(t *testing.T) {
		// Test retry mechanism
		stdout, _ := testutil.AssertCommandSuccess(t, "cloud", "connect", "--provider", "aws", "--endpoint", "retry_endpoint", "--retry")
		testutil.AssertOutputContains(t, stdout, "Successfully connected after retries")
	})
}
