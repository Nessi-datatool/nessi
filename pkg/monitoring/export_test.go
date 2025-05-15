package monitoring

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportMetricsToCSV(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run with timeout to prevent hanging
	testutil.RunWithTimeout(t, func() {
		// Create a temporary directory for test output
		tempDir, err := os.MkdirTemp("", "metrics-export-test")
		require.NoError(t, err)
		
		// Use t.Cleanup for more reliable cleanup
		t.Cleanup(func() {
			os.RemoveAll(tempDir)
		})
		
		// Create a monitor with some test metrics
		options := MonitorOptions{
			MetricsPort: 9090,
		}
		monitor, err := New(options)
		require.NoError(t, err)
		
		// Record some test metrics
		monitor.RecordMetric("test_metric_1", 123.45, nil)
		monitor.RecordMetric("test_metric_2", 67.89, nil)
		
		// Export to CSV
		outputPath := filepath.Join(tempDir, "metrics.csv")
		exportOptions := ExportOptions{
			Format:     "csv",
			OutputPath: outputPath,
		}
		
		result, err := monitor.ExportMetrics(exportOptions)
		require.NoError(t, err)
		require.Equal(t, outputPath, result)
		
		// Verify the file exists
		_, err = os.Stat(outputPath)
		require.NoError(t, err)
	})
}

func TestExportMetricsInvalidFormat(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run with timeout to prevent hanging
	testutil.RunWithTimeout(t, func() {
		// Create a monitor
		options := MonitorOptions{
			MetricsPort: 9090,
		}
		monitor, err := New(options)
		require.NoError(t, err)
		
		// Try to export with an invalid format
		exportOptions := ExportOptions{
			Format:     "invalid_format",
			OutputPath: "test_output.txt",
		}
		
		_, err = monitor.ExportMetrics(exportOptions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported export format")
	})
}

func TestExportMetricsNoData(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run with timeout to prevent hanging
	testutil.RunWithTimeout(t, func() {
		// Create a temporary directory for test output
		tempDir, err := os.MkdirTemp("", "metrics-export-test")
		require.NoError(t, err)
		
		// Use t.Cleanup for more reliable cleanup
		t.Cleanup(func() {
			os.RemoveAll(tempDir)
		})
		
		// Create a monitor without recording any metrics
		options := MonitorOptions{
			MetricsPort: 9090,
		}
		monitor, err := New(options)
		require.NoError(t, err)
		
		// Export to CSV
		outputPath := filepath.Join(tempDir, "empty_metrics.csv")
		exportOptions := ExportOptions{
			Format:     "csv",
			OutputPath: outputPath,
		}
		
		result, err := monitor.ExportMetrics(exportOptions)
		require.NoError(t, err)
		require.Equal(t, outputPath, result)
		
		// Verify the file exists but is essentially empty (just headers)
		data, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		assert.Contains(t, string(data), "metric_name")
		assert.Contains(t, string(data), "value")
	})
}
