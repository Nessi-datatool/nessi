//go:build !skiplong
// +build !skiplong

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

func TestExportMetricsToCSV(t *testing.T) {
	// Skip this test if we're in a mode that should skip long tests
	if os.Getenv("NESSI_SKIP_LONG_TESTS") != "" {
		
	}
	
	// Set up test timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
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
	
	// Signal test completion
	done <- true
}

func TestExportMetricsInvalidFormat(t *testing.T) {
	// Skip this test if we're in a mode that should skip long tests
	if os.Getenv("NESSI_SKIP_LONG_TESTS") != "" {
		
	}
	
	// Set up test timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
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
	
	// Signal test completion
	done <- true
}

func TestExportMetricsNoData(t *testing.T) {
	// Skip this test if we're in a mode that should skip long tests
	if os.Getenv("NESSI_SKIP_LONG_TESTS") != "" {
		
	}
	
	// Set up test timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
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
	
	// Signal test completion
	done <- true
}
