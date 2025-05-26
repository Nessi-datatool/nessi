package stresstests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestPerformanceBenchmark(t *testing.T) {
	// Skip if E2E tests are not enabled
	if os.Getenv("ENABLE_E2E_TESTS") != "true" {
		t.Skip("Skipping E2E test; set ENABLE_E2E_TESTS=true to run")
	}

	// Create a temporary directory for the test
	tempDir, cleanup := testutil.SetupTestEnvironment(t)
	defer cleanup()

	// Create test directories
	tablesDir := filepath.Join(tempDir, "tables")
	reportsDir := filepath.Join(tempDir, "reports")
	if err := os.MkdirAll(tablesDir, 0755); err != nil {
		t.Fatalf("Failed to create tables directory: %v", err)
	}
	if err := os.MkdirAll(reportsDir, 0755); err != nil {
		t.Fatalf("Failed to create reports directory: %v", err)
	}

	// Create mock Delta tables with different sizes
	smallTablePath := filepath.Join(tablesDir, "small_table")
	mediumTablePath := filepath.Join(tablesDir, "medium_table")
	largeTablePath := filepath.Join(tablesDir, "large_table")

	if err := testutil.CreateMockDeltaTableWithRows(smallTablePath, 1000); err != nil {
		t.Fatalf("Failed to create small mock Delta table: %v", err)
	}
	if err := testutil.CreateMockDeltaTableWithRows(mediumTablePath, 10000); err != nil {
		t.Fatalf("Failed to create medium mock Delta table: %v", err)
	}
	if err := testutil.CreateMockDeltaTableWithRows(largeTablePath, 100000); err != nil {
		t.Fatalf("Failed to create large mock Delta table: %v", err)
	}

	// Run performance benchmarks
	t.Run("Schema_Extraction_Performance", func(t *testing.T) {
		// Measure schema extraction performance for different table sizes
		tables := []struct {
			name string
			path string
		}{
			{"Small", smallTablePath},
			{"Medium", mediumTablePath},
			{"Large", largeTablePath},
		}

		results := make(map[string]time.Duration)

		for _, table := range tables {
			// Run multiple iterations to get an average
			var totalDuration time.Duration
			iterations := 5

			for i := 0; i < iterations; i++ {
				start := time.Now()
				cmd := testutil.RunNessiCommand("schema", "extract", "--path", table.path)
				testutil.AssertCommandSuccess(t, cmd)
				duration := time.Since(start)
				totalDuration += duration
			}

			// Calculate average duration
			results[table.name] = totalDuration / time.Duration(iterations)
		}

		// Log results
		t.Logf("Schema extraction performance: Small: %v, Medium: %v, Large: %v",
			results["Small"], results["Medium"], results["Large"])
	})

	t.Run("Table_Description_Performance", func(t *testing.T) {
		// Measure table description performance for different table sizes
		tables := []struct {
			name string
			path string
		}{
			{"Small", smallTablePath},
			{"Medium", mediumTablePath},
			{"Large", largeTablePath},
		}

		results := make(map[string]time.Duration)

		for _, table := range tables {
			// Run multiple iterations to get an average
			var totalDuration time.Duration
			iterations := 5

			for i := 0; i < iterations; i++ {
				start := time.Now()
				cmd := testutil.RunNessiCommand("table", "describe", "--path", table.path)
				testutil.AssertCommandSuccess(t, cmd)
				duration := time.Since(start)
				totalDuration += duration
			}

			// Calculate average duration
			results[table.name] = totalDuration / time.Duration(iterations)
		}

		// Log results
		t.Logf("Table description performance: Small: %v, Medium: %v, Large: %v",
			results["Small"], results["Medium"], results["Large"])
	})

	t.Run("Quality_Check_Performance", func(t *testing.T) {
		// Measure quality check performance for different table sizes
		tables := []struct {
			name string
			path string
		}{
			{"Small", smallTablePath},
			{"Medium", mediumTablePath},
			{"Large", largeTablePath},
		}

		results := make(map[string]time.Duration)

		for _, table := range tables {
			// Run multiple iterations to get an average
			var totalDuration time.Duration
			iterations := 5

			for i := 0; i < iterations; i++ {
				start := time.Now()
				cmd := testutil.RunNessiCommand("quality", "check", "--path", table.path)
				testutil.AssertCommandSuccess(t, cmd)
				duration := time.Since(start)
				totalDuration += duration
			}

			// Calculate average duration
			results[table.name] = totalDuration / time.Duration(iterations)
		}

		// Log results
		t.Logf("Quality check performance: Small: %v, Medium: %v, Large: %v",
			results["Small"], results["Medium"], results["Large"])
	})

	t.Run("Report_Generation_Performance", func(t *testing.T) {
		// Measure report generation performance for different formats
		formats := []string{"html", "json", "markdown", "csv"}

		for _, format := range formats {
			// Run multiple iterations to get an average
			var totalDuration time.Duration
			iterations := 5

			for i := 0; i < iterations; i++ {
				reportPath := filepath.Join(reportsDir, "report."+format)
				start := time.Now()
				cmd := testutil.RunNessiCommand("report", "generate", "--path", mediumTablePath, "--output", reportPath, "--format", format)
				testutil.AssertCommandSuccess(t, cmd)
				duration := time.Since(start)
				totalDuration += duration
			}

			// Calculate average duration
			avgDuration := totalDuration / time.Duration(iterations)
			t.Logf("Report generation performance (%s): %v", format, avgDuration)
		}
	})

	t.Run("Concurrent_Operations_Performance", func(t *testing.T) {
		// Measure performance of concurrent operations
		start := time.Now()

		// Create a channel to collect results
		done := make(chan bool, 3)

		// Run schema extraction, table description, and quality check concurrently
		go func() {
			cmd := testutil.RunNessiCommand("schema", "extract", "--path", mediumTablePath)
			testutil.AssertCommandSuccess(t, cmd)
			done <- true
		}()

		go func() {
			cmd := testutil.RunNessiCommand("table", "describe", "--path", mediumTablePath)
			testutil.AssertCommandSuccess(t, cmd)
			done <- true
		}()

		go func() {
			cmd := testutil.RunNessiCommand("quality", "check", "--path", mediumTablePath)
			testutil.AssertCommandSuccess(t, cmd)
			done <- true
		}()

		// Wait for all operations to complete
		for i := 0; i < 3; i++ {
			<-done
		}

		duration := time.Since(start)
		t.Logf("Concurrent operations performance: %v", duration)
	})
}
