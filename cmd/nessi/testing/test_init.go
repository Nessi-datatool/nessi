package testing

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ResetFlags resets all global flags to avoid conflicts
func ResetFlags() {
	// Reset viper
	viper.Reset()

	// Reset pflag.CommandLine
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

// InitTestCommand initializes a command for testing
func InitTestCommand(cmd *cobra.Command) {
	// Reset the command's flags
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
		InitTestCommand(subCmd)
	}

	// Create a clean environment for the command
	cmd.SetArgs([]string{})
	cmd.SetOut(nil)
	cmd.SetErr(nil)
}

// ExecuteCommand executes a command with the given arguments and returns the output
func ExecuteCommand(cmd *cobra.Command, args ...string) (string, error) {
	// Reset flags to avoid conflicts
	ResetFlags()

	// Initialize the command for testing
	InitTestCommand(cmd)

	// Set up command arguments
	cmd.SetArgs(args)

	// Capture output
	output := bytes.NewBufferString("")
	cmd.SetOut(output)
	cmd.SetErr(output)

	// Execute the command
	err := cmd.Execute()

	// Return the output and error
	return output.String(), err
}

// TestMain is a helper function for test initialization
func TestMain(m *testing.M, setup func(), teardown func()) int {
	// Set test environment variable
	os.Setenv("NESSI_TEST_MODE", "true")

	// Reset flags
	ResetFlags()

	// Run setup
	if setup != nil {
		setup()
	}

	// Run tests
	code := m.Run()

	// Run teardown
	if teardown != nil {
		teardown()
	}

	// Unset test environment variable
	os.Unsetenv("NESSI_TEST_MODE")

	return code
}
