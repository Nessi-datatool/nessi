package stresstests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestErrorHandlingComprehensive(t *testing.T) {
	// Skip if E2E tests are not enabled
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping E2E test; set ENABLE_E2E_TESTS=true to run")
	}

	// Create a temporary directory for the test
	tempDir, cleanup := testutil.SetupTestEnvironment(t)
	defer cleanup()

	// Create test directories
	configDir := filepath.Join(tempDir, "config")
	tableDir := filepath.Join(tempDir, "tables")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}
	if err := os.MkdirAll(tableDir, 0755); err != nil {
		t.Fatalf("Failed to create tables directory: %v", err)
	}

	// Create mock Delta table
	validTablePath := filepath.Join(tableDir, "valid_table")
	if err := testutil.CreateMockDeltaTable(validTablePath); err != nil {
		t.Fatalf("Failed to create mock Delta table: %v", err)
	}

	// Create invalid Delta table
	invalidTablePath := filepath.Join(tableDir, "invalid_table")
	if err := testutil.CreateInvalidDeltaTable(invalidTablePath); err != nil {
		t.Fatalf("Failed to create invalid Delta table: %v", err)
	}

	// Create directory with permission error
	permErrorPath := filepath.Join(tableDir, "perm_error")
	if err := testutil.CreatePermissionErrorDir(permErrorPath); err != nil {
		t.Fatalf("Failed to create permission error directory: %v", err)
	}

	// Create valid config file
	validConfigPath := filepath.Join(configDir, "valid_config.json")
	validConfig := `{"table_path": "` + validTablePath + `", "quality_threshold": 0.8}`
	if err := testutil.CreateMockConfigFile(validConfigPath, validConfig); err != nil {
		t.Fatalf("Failed to create valid config file: %v", err)
	}

	// Create invalid config file
	invalidConfigPath := filepath.Join(configDir, "invalid_config.json")
	if err := testutil.CreateInvalidConfigFile(invalidConfigPath); err != nil {
		t.Fatalf("Failed to create invalid config file: %v", err)
	}

	// Run subtests for different error scenarios
	t.Run("Non-existent_Table", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("table", "describe", "--path", filepath.Join(tableDir, "non_existent"))
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "Table not found")
	})

	t.Run("Invalid_Command_Flags", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("table", "describe", "--invalid-flag", "value")
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "unknown flag")
	})

	t.Run("Missing_Required_Arguments", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("table", "describe")
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "required flag")
	})

	t.Run("Invalid_File_Format", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("table", "describe", "--path", invalidTablePath)
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "invalid Delta table")
	})

	t.Run("Permission_Errors", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("table", "describe", "--path", permErrorPath)
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "permission denied")
	})

	t.Run("Invalid_Configuration_File", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("config", "validate", "--path", invalidConfigPath)
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "invalid configuration")
	})

	t.Run("Invalid_Connection_String", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("databricks", "connect", "--host", "invalid://host", "--token", "invalid-token")
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "invalid connection")
	})

	t.Run("Timeout_Handling", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("table", "describe", "--path", validTablePath, "--timeout", "1ms")
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "timeout")
	})

	t.Run("License_Validation", func(t *testing.T) {
		cmd := testutil.RunNessiCommand("license", "validate", "--key", "invalid-license-key")
		output := testutil.AssertCommandFailure(t, cmd)
		testutil.AssertOutputContains(t, output, "invalid license")
	})

	t.Run("Graceful_Shutdown", func(t *testing.T) {
		// Start a long-running command
		cmd := testutil.RunNessiCommand("table", "monitor", "--path", validTablePath, "--interval", "10s")
		
		// Start the command in the background
		if err := cmd.Start(); err != nil {
			t.Fatalf("Failed to start command: %v", err)
		}

		// Wait a bit for the command to start
		time.Sleep(2 * time.Second)

		// Send interrupt signal
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			t.Fatalf("Failed to send interrupt signal: %v", err)
		}

		// Wait for the command to exit
		err := cmd.Wait()
		
		// The command should exit with a non-zero status due to the interrupt
		if err == nil {
			t.Fatalf("Command did not exit with error after interrupt")
		}
	})
}
