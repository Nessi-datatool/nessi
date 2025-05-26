package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestConfigManagement(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Config Management test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Step 1: Initialize configuration in the temp directory
	stdout, _ := testutil.AssertCommandSuccess(t, "config", "init", "--dir", tempDir)
	testutil.AssertOutputContains(t, stdout, "Configuration initialized successfully")

	// Verify the .nessi directory and config file exist
	nessiDir := filepath.Join(tempDir, ".nessi")
	_, err := os.Stat(nessiDir)
	require.NoError(t, err, ".nessi directory should exist")

	configFile := filepath.Join(nessiDir, "config.yaml")
	_, err = os.Stat(configFile)
	require.NoError(t, err, "config.yaml file should exist")

	// Step 2: Load the configuration
	stdout, _ = testutil.AssertCommandSuccess(t, "config", "load", "--file", configFile)
	testutil.AssertOutputContains(t, stdout, "Configuration loaded successfully")

	// Step 3: Test loading a non-existent configuration file (should fail)
	nonExistentFile := filepath.Join(tempDir, "non_existent_config.yaml")
	_, stderr := testutil.AssertCommandFailure(t, "config", "load", "--file", nonExistentFile)
	testutil.AssertOutputContains(t, stderr, "Config file")
	testutil.AssertOutputContains(t, stderr, "does not exist")

	// Step 4: Create a custom configuration file
	customConfigFile := filepath.Join(tempDir, "custom_config.yaml")
	customConfigContent := `# Custom Nessi Configuration
version: 1.1
settings:
  log_level: debug
  data_directory: /tmp/nessi-data
  max_connections: 10
`
	require.NoError(t, os.WriteFile(customConfigFile, []byte(customConfigContent), 0644))

	// Load the custom configuration
	stdout, _ = testutil.AssertCommandSuccess(t, "config", "load", "--file", customConfigFile)
	testutil.AssertOutputContains(t, stdout, "Configuration loaded successfully")
}
