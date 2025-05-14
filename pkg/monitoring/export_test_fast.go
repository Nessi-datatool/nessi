package monitoring

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/stretchr/testify/require"
)

// TestExportMetricsToCSVFast is a fast version of the export metrics test
func TestExportMetricsToCSVFast(t *testing.T) {
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "metrics-export-test-fast")
	require.NoError(t, err)
	
	// Clean up the temporary directory after the test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	
	// Create a monitor with minimal configuration
	options := MonitorOptions{
		MetricsPort: 9090,
	}
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Record a single test metric
	monitor.RecordMetric("test_metric_fast", 100.0, nil)
	
	// Export to CSV with minimal options
	outputPath := filepath.Join(tempDir, "metrics_fast.csv")
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
}
