package common

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockReader is a mock implementation of bufio.Reader for testing
type mockReader struct {
	responses []string
	index     int
}

// ReadString implements the bufio.Reader interface for testing
func (m *mockReader) ReadString(delim byte) (string, error) {
	if m.index >= len(m.responses) {
		return "", nil
	}

	response := m.responses[m.index]
	m.index++

	return response + string(delim), nil
}

// TestInteractiveErrorResolver tests the interactive error resolver
func TestInteractiveErrorResolver(t *testing.T) {
	// Create a buffer to capture output
	output := &bytes.Buffer{}

	// Create a resolver with a mock reader and buffer
	resolver := &InteractiveErrorResolver{
		Reader: bufio.NewReader(strings.NewReader("1\n")), // Select option 1
		Writer: bufio.NewWriter(output),
	}

	// Create a resolvable error
	err := &mockResolvableError{}

	// Resolve the error
	resolvedErr := resolver.ResolveError(err)

	// Verify the error was resolved
	assert.NoError(t, resolvedErr)

	// Verify the output contains the expected text
	outputStr := output.String()
	assert.Contains(t, outputStr, "Error: mock error")
	assert.Contains(t, outputStr, "Resolution options:")
	assert.Contains(t, outputStr, "1. Option 1")
	assert.Contains(t, outputStr, "2. Option 2")
	assert.Contains(t, outputStr, "0. Do nothing")
	assert.Contains(t, outputStr, "Select an option:")
}

// mockResolvableError is a mock implementation of ResolvableError for testing
type mockResolvableError struct {
	option int
}

// GetResolutionOptions implements the ResolvableError interface for testing
func (m *mockResolvableError) GetResolutionOptions() []string {
	return []string{"Option 1", "Option 2"}
}

// ResolveWithOption implements the ResolvableError interface for testing
func (m *mockResolvableError) ResolveWithOption(option int) error {
	m.option = option
	return nil
}

// Error implements the error interface for testing
func (m *mockResolvableError) Error() string {
	return "mock error"
}

// TestResolvablePathError tests the resolvable path error
func TestResolvablePathError(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a non-existent path
	nonExistentPath := tempDir + "/non-existent"

	// Create a resolvable path error
	pathErr := NewResolvablePathError(nonExistentPath)

	// Verify the error message
	assert.Contains(t, pathErr.Error(), nonExistentPath)

	// Verify the resolution options
	options := pathErr.GetResolutionOptions()
	assert.Equal(t, 2, len(options))
	assert.Contains(t, options[0], "Create")
	assert.Contains(t, options[1], "Specify")

	// Resolve the error by creating the directory
	err = pathErr.ResolveWithOption(0)
	assert.NoError(t, err)

	// Verify the directory was created
	_, err = os.Stat(nonExistentPath)
	assert.NoError(t, err)
}

// TestResolvableDeltaTableError tests the resolvable Delta table error
func TestResolvableDeltaTableError(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a non-Delta table directory
	nonDeltaPath := tempDir + "/non-delta"
	err = os.Mkdir(nonDeltaPath, 0755)
	assert.NoError(t, err)

	// Create a resolvable Delta table error
	deltaErr := NewResolvableDeltaTableError(nonDeltaPath)

	// Verify the error message
	assert.Contains(t, deltaErr.Error(), nonDeltaPath)

	// Verify the resolution options
	options := deltaErr.GetResolutionOptions()
	assert.Equal(t, 2, len(options))
	assert.Contains(t, options[0], "Initialize")
	assert.Contains(t, options[1], "Specify")

	// Resolve the error by initializing as a Delta table
	err = deltaErr.ResolveWithOption(0)
	assert.NoError(t, err)

	// Verify the Delta table was initialized
	deltaLogPath := nonDeltaPath + "/_delta_log"
	_, err = os.Stat(deltaLogPath)
	assert.NoError(t, err)

	// Verify the transaction log was created
	transactionLogPath := deltaLogPath + "/00000000000000000000.json"
	_, err = os.Stat(transactionLogPath)
	assert.NoError(t, err)
}

// TestResolvableConfigError tests the resolvable configuration error
func TestResolvableConfigError(t *testing.T) {
	// Create a resolvable configuration error
	configErr := NewResolvableConfigError("DATABRICKS_WORKSPACE_ID", "")

	// Verify the error message
	assert.Contains(t, configErr.Error(), "DATABRICKS_WORKSPACE_ID")

	// Verify the resolution options
	options := configErr.GetResolutionOptions()
	assert.Equal(t, 2, len(options))
	assert.Contains(t, options[0], "Provide")
	assert.Contains(t, options[1], "default")

	// Resolve the error by using the default value
	err := configErr.ResolveWithOption(1)
	assert.NoError(t, err)

	// Verify the default value was used
	assert.Equal(t, "0", configErr.ConfigValue)
}
