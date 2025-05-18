package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSchemaTreeCmd_Basic(t *testing.T) {
	// Create a sample Parquet/Arrow IPC file for testing
	tempDir := t.TempDir()
	tablePath := filepath.Join(tempDir, "test_table.arrow")
	// Use echo to create a dummy file; in real test, generate a valid Arrow IPC file
	f, err := os.Create(tablePath)
	if err != nil {
		t.Fatalf("Failed to create temp table: %v", err)
	}
	f.Close()

	// Run the CLI command
	cmd := exec.Command("go", "run", "./cli.go", "schema-tree", "--table", tablePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = filepath.Dir(tablePath)
	err = cmd.Run()
	if err == nil {
		t.Errorf("Expected error for dummy file, got nil")
	}
	if !strings.Contains(out.String(), "failed to read Arrow IPC") {
		t.Errorf("Expected Arrow IPC error, got: %s", out.String())
	}
}

// TODO: Add a test with a real Arrow IPC or Parquet file for full integration coverage
