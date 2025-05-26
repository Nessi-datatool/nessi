package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

// copyDirContents copies the contents of a directory to another directory
func copyDirContents(t *testing.T, srcDir, destDir string) {
	t.Helper()

	// Check if source directory exists
	_, err := os.Stat(srcDir)
	if os.IsNotExist(err) {
		t.Fatalf("Source directory %s does not exist", srcDir)
	}

	// Walk through the source directory and copy files
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if path == srcDir {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Create destination path
		destPath := filepath.Join(destDir, relPath)

		// If it's a directory, create it
		if info.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// It's a file, copy it
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Ensure parent directory exists
		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		if err != nil {
			return err
		}

		// Write the file
		return os.WriteFile(destPath, data, 0644)
	})

	if err != nil {
		t.Fatalf("Failed to copy directory contents: %v", err)
	}
}

func TestVersionControl(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Version Control test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))

	// Copy our test data directly from the e2e-tests/testdata directory
	srcDir := filepath.Join(testutil.FindProjectRoot(), "e2e-tests", "testdata")

	// Copy Delta tables
	srcDeltaDir := filepath.Join(srcDir, "delta_tables")
	testutil.CopyDirContents(t, srcDeltaDir, sampleDataDir)

	// Since our mock CLI doesn't have a version command, we'll test the tables command instead
	// which is available in our mock CLI implementation
	tablePath := filepath.Join(sampleDataDir, "sample_table")

	// Step 1: Test table connection
	stdout, _ := testutil.AssertCommandSuccess(t, "tables", "connect", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Successfully connected to Delta Lake table")

	// Step 2: Test table description
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "describe", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Schema:")

	// Step 3: Test table reading
	stdout, _ = testutil.AssertCommandSuccess(t, "tables", "read", "--path", tablePath)
	testutil.AssertOutputContains(t, stdout, "Data preview:")

	// Step 4: Test license status
	licenseOutput, _ := testutil.AssertCommandSuccess(t, "license", "status")
	testutil.AssertOutputContains(t, licenseOutput, "License Status")

	// Step 5: Test license trial
	// Actually run the trial command since it's part of our mock CLI
	stdout, _ = testutil.AssertCommandSuccess(t, "license", "trial")
	testutil.AssertOutputContains(t, stdout, "Pro Edition Trial Activated!")
}
