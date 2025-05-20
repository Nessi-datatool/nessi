package mocks

import (
	"bytes"

	"github.com/spf13/cobra"
)

// ExecuteCommand executes a command for testing and returns its output
func ExecuteCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

// SetupTestCommand sets up a command for testing with output and error buffers
func SetupTestCommand(cmd *cobra.Command) (*bytes.Buffer, *bytes.Buffer) {
	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	cmd.SetOut(outBuf)
	cmd.SetErr(errBuf)
	return outBuf, errBuf
}
