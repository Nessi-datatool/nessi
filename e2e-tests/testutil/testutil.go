package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Command represents a CLI command to be executed
type Command struct {
	args   []string
	cmd    *exec.Cmd
	stdout bytes.Buffer
	stderr bytes.Buffer
}

// NewCommand creates a new Command with the given arguments
func NewCommand(args ...string) *Command {
	return &Command{
		args: args,
	}
}

// Run executes the command and returns its output
func (c *Command) Run() (string, error) {
	// Use the mock CLI for testing
	mockPath := "/Users/meisi/Documents/nessi/e2e-tests/mock/main.go"

	// Check if the mock CLI exists
	if _, err := os.Stat(mockPath); os.IsNotExist(err) {
		return "", fmt.Errorf("mock CLI not found at %s", mockPath)
	}

	c.cmd = exec.Command("go", append([]string{"run", mockPath}, c.args...)...)
	c.cmd.Stdout = &c.stdout
	c.cmd.Stderr = &c.stderr

	err := c.cmd.Run()
	output := c.stdout.String() + c.stderr.String()

	if err != nil {
		return output, fmt.Errorf("command failed: %v\nOutput: %s", err, output)
	}

	return output, nil
}

// Start starts the command without waiting for it to complete
func (c *Command) Start() error {
	mockPath := "/Users/meisi/Documents/nessi/e2e-tests/mock/main.go"

	// Check if the mock CLI exists
	if _, err := os.Stat(mockPath); os.IsNotExist(err) {
		return fmt.Errorf("mock CLI not found at %s", mockPath)
	}

	c.cmd = exec.Command("go", append([]string{"run", mockPath}, c.args...)...)
	c.cmd.Stdout = &c.stdout
	c.cmd.Stderr = &c.stderr

	return c.cmd.Start()
}

// Wait waits for the command to complete
func (c *Command) Wait() error {
	return c.cmd.Wait()
}

// Process returns the underlying process
func (c *Command) Process() *os.Process {
	if c.cmd != nil {
		return c.cmd.Process
	}
	return nil
}

// CombinedOutput returns the combined stdout and stderr output
func (c *Command) CombinedOutput() string {
	return c.stdout.String() + c.stderr.String()
}

// CreateMockDeltaTable creates a mock Delta Lake table for testing
func CreateMockDeltaTable(path string) error {
	// Create the directory structure
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// Create _delta_log directory
	deltaLogDir := filepath.Join(path, "_delta_log")
	if err := os.MkdirAll(deltaLogDir, 0755); err != nil {
		return err
	}

	// Create a simple transaction log file
	txnLogFile := filepath.Join(deltaLogDir, "00000000000000000000.json")
	txnLog := map[string]interface{}{
		"commitInfo": map[string]interface{}{
			"timestamp": 1620000000000,
			"operation": "CREATE TABLE",
			"operationParameters": map[string]interface{}{
				"isManaged": true,
			},
		},
	}

	txnLogBytes, err := json.MarshalIndent(txnLog, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(txnLogFile, txnLogBytes, 0644); err != nil {
		return err
	}

	// Create a simple metadata file
	metadataFile := filepath.Join(deltaLogDir, "00000000000000000000.metadata.json")
	metadata := map[string]interface{}{
		"metaData": map[string]interface{}{
			"id": "test-table-id",
			"format": map[string]interface{}{
				"provider": "parquet",
			},
			"schemaString": `{
				"type": "struct",
				"fields": [
					{"name": "id", "type": "integer", "nullable": false},
					{"name": "name", "type": "string", "nullable": true},
					{"name": "value", "type": "double", "nullable": true},
					{"name": "timestamp", "type": "timestamp", "nullable": true}
				]
			}`,
		},
	}

	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(metadataFile, metadataBytes, 0644); err != nil {
		return err
	}

	// Create a sample parquet file
	sampleParquetFile := filepath.Join(path, "part-00000-sample.parquet")
	if err := os.WriteFile(sampleParquetFile, []byte("MOCK PARQUET DATA"), 0644); err != nil {
		return err
	}

	return nil
}

// CreateMockDeltaTableWithRows creates a mock Delta Lake table with the specified number of rows
func CreateMockDeltaTableWithRows(path string, rowCount int) error {
	// First create the basic table structure
	if err := CreateMockDeltaTable(path); err != nil {
		return err
	}

	// Update the metadata to reflect the row count
	deltaLogDir := filepath.Join(path, "_delta_log")
	statsFile := filepath.Join(deltaLogDir, "00000000000000000001.json")
	stats := map[string]interface{}{
		"add": map[string]interface{}{
			"path": "part-00000-sample.parquet",
			"size": 1024 * rowCount, // Simulate file size based on row count
			"stats": map[string]interface{}{
				"numRecords": rowCount,
				"minValues": map[string]interface{}{
					"id":    1,
					"value": 0.1,
				},
				"maxValues": map[string]interface{}{
					"id":    rowCount,
					"value": 100.0,
				},
				"nullCount": map[string]interface{}{
					"id":        0,
					"name":      rowCount / 10, // 10% null values
					"value":     rowCount / 20, // 5% null values
					"timestamp": 0,
				},
			},
		},
	}

	statsBytes, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(statsFile, statsBytes, 0644); err != nil {
		return err
	}

	return nil
}

// CopyDirContents copies the contents of a directory to another directory
func CopyDirContents(src, dst string) error {
	// Get properties of source directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	// Read the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursive copy for directories
			if err = CopyDirContents(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			if err = copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file from src to dst
func copyFile(src, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Get the source file mode
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Set the same mode on the destination file
	return os.Chmod(dst, srcInfo.Mode())
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T, prefix string) string {
	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	return tempDir
}

// CopyTestData copies test data from the testdata directory to the target directory
func CopyTestData(t *testing.T, src, dst string) {
	err := CopyDirContents(src, dst)
	if err != nil {
		t.Fatalf("Failed to copy test data: %v", err)
	}
}

// AssertCommandSuccess runs a command and asserts that it succeeds
func AssertCommandSuccess(t *testing.T, cmd *Command) string {
	output, err := cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}
	return output
}

// AssertCommandFailure runs a command and asserts that it fails
func AssertCommandFailure(t *testing.T, cmd *Command) string {
	output, err := cmd.Run()
	if err == nil {
		t.Fatalf("Command succeeded unexpectedly\nOutput: %s", output)
	}
	return output
}

// AssertOutputContains asserts that the output contains the expected string
func AssertOutputContains(t *testing.T, output, expected string) {
	if !strings.Contains(output, expected) {
		t.Fatalf("Output does not contain '%s'\nOutput: %s", expected, output)
	}
}
