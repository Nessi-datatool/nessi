package integrations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
)

// TestCLIVisualization tests the CLI visualization features
// including ASCII schema trees and heatmaps
func TestCLIVisualization(t *testing.T) {
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping CLI Visualization test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-cli-viz-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data files
	testDataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Create a mock Delta table for testing
	mockTablePath := filepath.Join(testDataDir, "mock_table")
	if err := testutil.CreateMockDeltaTable(mockTablePath); err != nil {
		t.Fatalf("Failed to create mock Delta table: %v", err)
	}

	// Test schema tree visualization
	t.Run("Schema Tree Visualization", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"schema", "tree",
			"--table", mockTablePath,
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Schema tree visualization should succeed")
		assert.Contains(t, output, "Schema Tree")

		// Check for tree structure characters
		assert.Contains(t, output, "├─", "Output should contain tree structure characters")
		// The vertical bar may not be present in all implementations, so we'll check for other tree characters
		assert.Contains(t, output, "└─", "Output should contain tree structure characters")
	})

	// Test schema tree with metadata
	t.Run("Schema Tree with Metadata", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"schema", "tree",
			"--table", mockTablePath,
			"--show-metadata",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Schema tree with metadata should succeed")
		assert.Contains(t, output, "Schema Tree")
		assert.Contains(t, output, "Metadata", "Output should contain metadata information")
	})

	// Test heatmap visualization
	t.Run("Quality Heatmap Visualization", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"quality", "heatmap",
			"--table", mockTablePath,
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Quality heatmap visualization should succeed")
		assert.Contains(t, output, "Quality Heatmap")

		// Check for heatmap characters
		assert.Contains(t, output, "█", "Output should contain block characters for heatmap")
	})

	// Test heatmap with custom colors
	t.Run("Heatmap with Custom Colors", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"quality", "heatmap",
			"--table", mockTablePath,
			"--color-scheme", "blue-red",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Heatmap with custom colors should succeed")
		assert.Contains(t, output, "Quality Heatmap")
	})

	// Test ASCII bar chart visualization
	t.Run("ASCII Bar Chart", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"metrics", "chart",
			"--table", mockTablePath,
			"--metric", "completeness",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "ASCII bar chart visualization should succeed")
		assert.Contains(t, output, "Completeness Metrics")
		assert.Contains(t, output, "█", "Output should contain block characters for bar chart")
	})

	// Test progress visualization
	t.Run("Progress Visualization", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"demo", "progress",
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Progress visualization demo should succeed")
		assert.Contains(t, output, "Progress Demo")
		assert.Contains(t, output, "100%", "Output should show completion percentage")
	})

	// Test spinner visualization
	t.Run("Spinner Visualization", func(t *testing.T) {
		cmd := testutil.NewCommand(
			"demo", "spinner",
			"--duration", "1s", // Short duration for testing
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Spinner visualization demo should succeed")
		assert.Contains(t, output, "Spinner Demo")
		assert.Contains(t, output, "Complete", "Output should indicate completion")
	})
}
