package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand executes a shell command and returns its output
func RunCommand(command string) (string, error) {
	// Create the command
	cmd := exec.Command("bash", "-c", command)

	// Create buffers for stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()

	// Get the output
	output := stdout.String()

	// If there's an error, include stderr in the error message
	if err != nil {
		errOutput := stderr.String()
		if errOutput != "" {
			return output, fmt.Errorf("%w: %s", err, strings.TrimSpace(errOutput))
		}
		return output, err
	}

	return strings.TrimSpace(output), nil
}

// IsCommandAvailable checks if a command is available in the system
func IsCommandAvailable(command string) bool {
	// Extract the command name (before any arguments)
	cmdName := strings.Split(command, " ")[0]

	// Check if the command exists
	_, err := exec.LookPath(cmdName)
	return err == nil
}
