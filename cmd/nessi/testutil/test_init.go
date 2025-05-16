package testutil

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ResetCommandFlags resets all flags for a command and its subcommands
func ResetCommandFlags(cmd *cobra.Command) {
	cmd.ResetFlags()
	cmd.ResetCommands()
	
	// Reset persistent flags
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		cmd.PersistentFlags().Set(f.Name, f.DefValue)
	})
	
	// Reset local flags
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		cmd.Flags().Set(f.Name, f.DefValue)
	})
	
	// Reset flags for all subcommands
	for _, subCmd := range cmd.Commands() {
		ResetCommandFlags(subCmd)
	}
}

// InitTestCommand initializes a command for testing
func InitTestCommand(cmd *cobra.Command) {
	// Reset all flags to avoid conflicts
	ResetCommandFlags(cmd)
	
	// Create a clean environment for the command
	cmd.SetArgs([]string{})
	cmd.SetOut(nil)
	cmd.SetErr(nil)
}

// SetupTestEnv sets up the test environment
func SetupTestEnv() {
	// Set test-specific environment variables
	os.Setenv("NESSI_TEST_MODE", "true")
}

// TeardownTestEnv cleans up the test environment
func TeardownTestEnv() {
	// Unset test-specific environment variables
	os.Unsetenv("NESSI_TEST_MODE")
}

// TestMain is a helper function for test initialization
func TestMain(m *testing.M, setup func(), teardown func()) int {
	// Set up test environment
	SetupTestEnv()
	
	// Run user-provided setup
	if setup != nil {
		setup()
	}
	
	// Run tests
	code := m.Run()
	
	// Run user-provided teardown
	if teardown != nil {
		teardown()
	}
	
	// Clean up test environment
	TeardownTestEnv()
	
	return code
}
