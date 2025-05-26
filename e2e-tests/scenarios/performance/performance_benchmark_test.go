package performance

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPerformanceBenchmarks tests the performance of various CLI operations
// under different data volumes and conditions
func TestPerformanceBenchmarks(t *testing.T) {
	if os.Getenv("ENABLE_PERFORMANCE_TESTS") != "true" {
		t.Skip("Skipping Performance Benchmark tests. Use ENABLE_PERFORMANCE_TESTS=true to run.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-perf-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test data files of different sizes
	testDataDir := filepath.Join(tempDir, "data")
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Create small, medium, and large mock Delta tables for testing
	smallTablePath := filepath.Join(testDataDir, "small_table")
	mediumTablePath := filepath.Join(testDataDir, "medium_table")
	largeTablePath := filepath.Join(testDataDir, "large_table")

	// Create tables with different row counts
	require.NoError(t, testutil.CreateMockDeltaTableWithRows(smallTablePath, 1000), "Failed to create small table")
	require.NoError(t, testutil.CreateMockDeltaTableWithRows(mediumTablePath, 10000), "Failed to create medium table")
	require.NoError(t, testutil.CreateMockDeltaTableWithRows(largeTablePath, 100000), "Failed to create large table")

	// Benchmark table describe operation
	t.Run("Benchmark Table Describe", func(t *testing.T) {
		benchmarkCommand(t, "Small Table Describe", []string{
			"tables", "describe",
			"--table", smallTablePath,
		})

		benchmarkCommand(t, "Medium Table Describe", []string{
			"tables", "describe",
			"--table", mediumTablePath,
		})

		benchmarkCommand(t, "Large Table Describe", []string{
			"tables", "describe",
			"--table", largeTablePath,
		})
	})

	// Benchmark quality check operation
	t.Run("Benchmark Quality Check", func(t *testing.T) {
		benchmarkCommand(t, "Small Table Quality Check", []string{
			"quality", "check",
			"--table", smallTablePath,
		})

		benchmarkCommand(t, "Medium Table Quality Check", []string{
			"quality", "check",
			"--table", mediumTablePath,
		})

		benchmarkCommand(t, "Large Table Quality Check", []string{
			"quality", "check",
			"--table", largeTablePath,
		})
	})

	// Benchmark schema extraction
	t.Run("Benchmark Schema Extraction", func(t *testing.T) {
		benchmarkCommand(t, "Small Table Schema Extract", []string{
			"schema", "extract",
			"--table", smallTablePath,
		})

		benchmarkCommand(t, "Medium Table Schema Extract", []string{
			"schema", "extract",
			"--table", mediumTablePath,
		})

		benchmarkCommand(t, "Large Table Schema Extract", []string{
			"schema", "extract",
			"--table", largeTablePath,
		})
	})

	// Benchmark report generation
	t.Run("Benchmark Report Generation", func(t *testing.T) {
		outputFile := filepath.Join(tempDir, "report.html")

		benchmarkCommand(t, "Small Table Report Generation", []string{
			"report", "generate",
			"--table", smallTablePath,
			"--output", outputFile,
			"--format", "html",
		})

		benchmarkCommand(t, "Medium Table Report Generation", []string{
			"report", "generate",
			"--table", mediumTablePath,
			"--output", outputFile,
			"--format", "html",
		})

		benchmarkCommand(t, "Large Table Report Generation", []string{
			"report", "generate",
			"--table", largeTablePath,
			"--output", outputFile,
			"--format", "html",
		})
	})

	// Benchmark memory usage
	t.Run("Benchmark Memory Usage", func(t *testing.T) {
		// Run with memory profiling enabled
		outputFile := filepath.Join(tempDir, "memory_profile.out")

		cmd := testutil.NewCommand(
			"quality", "check",
			"--table", largeTablePath,
			"--memory-profile", outputFile,
		)
		output, err := cmd.Run()
		assert.NoError(t, err, "Memory profiling should succeed")
		assert.Contains(t, output, "Memory profile written to", "Output should indicate memory profile was written")

		// Verify the profile file exists
		_, err = os.Stat(outputFile)
		assert.NoError(t, err, "Memory profile file should exist")
	})

	// Benchmark concurrent operations
	t.Run("Benchmark Concurrent Operations", func(t *testing.T) {
		// Create a channel to collect results
		results := make(chan error, 3)

		// Start 3 concurrent operations
		go func() {
			cmd := testutil.NewCommand(
				"tables", "describe",
				"--table", smallTablePath,
			)
			_, err := cmd.Run()
			results <- err
		}()

		go func() {
			cmd := testutil.NewCommand(
				"quality", "check",
				"--table", mediumTablePath,
			)
			_, err := cmd.Run()
			results <- err
		}()

		go func() {
			cmd := testutil.NewCommand(
				"schema", "extract",
				"--table", largeTablePath,
			)
			_, err := cmd.Run()
			results <- err
		}()

		// Wait for all operations to complete
		for i := 0; i < 3; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent operation should succeed")
		}
	})
}

// benchmarkCommand runs a command and measures its execution time
func benchmarkCommand(t *testing.T, name string, args []string) {
	t.Helper()

	cmd := testutil.NewCommand(args...)

	start := time.Now()
	output, err := cmd.Run()
	duration := time.Since(start)

	assert.NoError(t, err, "%s should succeed", name)
	t.Logf("%s completed in %v", name, duration)

	// Verify the command produced some output
	assert.NotEmpty(t, output, "%s should produce output", name)
}
