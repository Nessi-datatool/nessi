package monitoring

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExportMetricsToCSVOptimized is a fast version of the export test
func TestExportMetricsToCSVOptimized(t *testing.T) {
	// Make this test run in parallel with other tests
	t.Parallel()
	
	// Set up test timeout to prevent hanging - using a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// Use context to ensure test doesn't hang
	done := make(chan bool)
	go func() {
		select {
		case <-done:
			return
		case <-ctx.Done():
			t.Error("Test timed out")
			t.FailNow()
		}
	}()
	defer close(done)
	
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "metrics-export-fast-test")
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
	
	// Signal test completion
	done <- true
}

// TestExportMetricsNoDataOptimized is a fast version of the no data export test
func TestExportMetricsNoDataOptimized(t *testing.T) {
	// Make this test run in parallel with other tests
	t.Parallel()
	
	// Set up test timeout to prevent hanging - using a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// Use context to ensure test doesn't hang
	done := make(chan bool)
	go func() {
		select {
		case <-done:
			return
		case <-ctx.Done():
			t.Error("Test timed out")
			t.FailNow()
		}
	}()
	defer close(done)
	
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "metrics-export-fast-test")
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
	
	// Signal test completion
	done <- true
}
