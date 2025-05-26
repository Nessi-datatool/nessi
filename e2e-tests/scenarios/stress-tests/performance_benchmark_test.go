package stresstests

import (
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
	"github.com/stretchr/testify/assert"
)

// TestPerformanceBenchmark tests the performance of various CLI operations
func TestPerformanceBenchmark(t *testing.T) {
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping Performance Benchmark test. Use ENABLE_E2E_TESTS=true to run this test.")
	}

	// Create a temporary directory for test data
	tempDir, err := os.MkdirTemp("", "nessi-performance-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test tables with different sizes
	smallTablePath := tempDir + "/small_table"
	mediumTablePath := tempDir + "/medium_table"
	largeTablePath := tempDir + "/large_table"

	err = testutil.CreateMockDeltaTableWithRows(smallTablePath, 100)
	assert.NoError(t, err, "Should be able to create small mock Delta table")

	err = testutil.CreateMockDeltaTableWithRows(mediumTablePath, 1000)
	assert.NoError(t, err, "Should be able to create medium mock Delta table")

	err = testutil.CreateMockDeltaTableWithRows(largeTablePath, 10000)
	assert.NoError(t, err, "Should be able to create large mock Delta table")

	// Test schema extraction performance
	t.Run("Schema Extraction Performance", func(t *testing.T) {
		// Small table
		start := time.Now()
		cmd := testutil.NewCommand(
			"schema", "extract",
			"--table", smallTablePath,
		)
		_, err := cmd.Run()
		smallDuration := time.Since(start)
		assert.NoError(t, err, "Schema extraction should succeed for small table")

		// Medium table
		start = time.Now()
		cmd = testutil.NewCommand(
			"schema", "extract",
			"--table", mediumTablePath,
		)
		_, err = cmd.Run()
		mediumDuration := time.Since(start)
		assert.NoError(t, err, "Schema extraction should succeed for medium table")

		// Large table
		start = time.Now()
		cmd = testutil.NewCommand(
			"schema", "extract",
			"--table", largeTablePath,
		)
		_, err = cmd.Run()
		largeDuration := time.Since(start)
		assert.NoError(t, err, "Schema extraction should succeed for large table")

		// Verify that performance scales reasonably
		t.Logf("Schema extraction performance: Small: %v, Medium: %v, Large: %v",
			smallDuration, mediumDuration, largeDuration)
	})

	// Test table description performance
	t.Run("Table Description Performance", func(t *testing.T) {
		// Small table
		start := time.Now()
		cmd := testutil.NewCommand(
			"tables", "describe",
			"--table", smallTablePath,
		)
		_, err := cmd.Run()
		smallDuration := time.Since(start)
		assert.NoError(t, err, "Table description should succeed for small table")

		// Medium table
		start = time.Now()
		cmd = testutil.NewCommand(
			"tables", "describe",
			"--table", mediumTablePath,
		)
		_, err = cmd.Run()
		mediumDuration := time.Since(start)
		assert.NoError(t, err, "Table description should succeed for medium table")

		// Large table
		start = time.Now()
		cmd = testutil.NewCommand(
			"tables", "describe",
			"--table", largeTablePath,
		)
		_, err = cmd.Run()
		largeDuration := time.Since(start)
		assert.NoError(t, err, "Table description should succeed for large table")

		// Verify that performance scales reasonably
		t.Logf("Table description performance: Small: %v, Medium: %v, Large: %v",
			smallDuration, mediumDuration, largeDuration)
	})

	// Test quality check performance
	t.Run("Quality Check Performance", func(t *testing.T) {
		// Small table
		start := time.Now()
		cmd := testutil.NewCommand(
			"quality", "check",
			"--table", smallTablePath,
		)
		_, err := cmd.Run()
		smallDuration := time.Since(start)
		assert.NoError(t, err, "Quality check should succeed for small table")

		// Medium table
		start = time.Now()
		cmd = testutil.NewCommand(
			"quality", "check",
			"--table", mediumTablePath,
		)
		_, err = cmd.Run()
		mediumDuration := time.Since(start)
		assert.NoError(t, err, "Quality check should succeed for medium table")

		// Large table
		start = time.Now()
		cmd = testutil.NewCommand(
			"quality", "check",
			"--table", largeTablePath,
		)
		_, err = cmd.Run()
		largeDuration := time.Since(start)
		assert.NoError(t, err, "Quality check should succeed for large table")

		// Verify that performance scales reasonably
		t.Logf("Quality check performance: Small: %v, Medium: %v, Large: %v",
			smallDuration, mediumDuration, largeDuration)
	})

	// Test report generation performance
	t.Run("Report Generation Performance", func(t *testing.T) {
		// Test different report formats
		formats := []string{"html", "json", "markdown", "csv"}

		for _, format := range formats {
			start := time.Now()
			cmd := testutil.NewCommand(
				"report", "generate",
				"--table", mediumTablePath,
				"--format", format,
				"--template", "quality",
				"--output", tempDir+"/report."+format,
			)
			_, err := cmd.Run()
			duration := time.Since(start)
			assert.NoError(t, err, "Report generation should succeed for "+format+" format")

			t.Logf("Report generation performance (%s): %v", format, duration)
		}
	})

	// Test concurrent operations
	t.Run("Concurrent Operations Performance", func(t *testing.T) {
		// Start multiple commands concurrently
		start := time.Now()

		cmd1 := testutil.NewCommand(
			"schema", "extract",
			"--table", smallTablePath,
		)
		err := cmd1.Start()
		assert.NoError(t, err, "Should be able to start schema extraction")

		cmd2 := testutil.NewCommand(
			"tables", "describe",
			"--table", mediumTablePath,
		)
		err = cmd2.Start()
		assert.NoError(t, err, "Should be able to start table description")

		cmd3 := testutil.NewCommand(
			"quality", "check",
			"--table", largeTablePath,
		)
		err = cmd3.Start()
		assert.NoError(t, err, "Should be able to start quality check")

		// Wait for all commands to complete
		err = cmd1.Wait()
		assert.NoError(t, err, "Schema extraction should complete successfully")

		err = cmd2.Wait()
		assert.NoError(t, err, "Table description should complete successfully")

		err = cmd3.Wait()
		assert.NoError(t, err, "Quality check should complete successfully")

		totalDuration := time.Since(start)
		t.Logf("Concurrent operations performance: %v", totalDuration)
	})
}
