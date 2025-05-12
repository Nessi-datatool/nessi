package monitoring

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportMetricsToCSV(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "metric-export-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create retention config
	config := RetentionConfig{
		Enabled:          true,
		StoragePath:      tempDir,
		RetentionPeriod:  24 * time.Hour,
		SnapshotInterval: 1 * time.Second,
	}

	// Create metric retention
	retention := NewMetricRetention(config)
	require.NotNil(t, retention)

	// Record some metrics
	now := time.Now()
	for i := 0; i < 10; i++ {
		timestamp := now.Add(time.Duration(-i) * time.Hour)
		retention.RecordMetric("test_metric", "Test metric", MetricTypeGauge, float64(i*10), map[string]string{
			"service": "test",
			"env":     "dev",
			"region":  "us-west",
		})
		
		// Manually set timestamp for testing
		series, ok := retention.currentSnapshot.Metrics["test_metric"]
		require.True(t, ok)
		series.Values[len(series.Values)-1].Timestamp = timestamp
		retention.currentSnapshot.Metrics["test_metric"] = series
	}

	// Create export options
	exportPath := filepath.Join(tempDir, "export.csv")
	options := ExportOptions{
		Format:     ExportFormatCSV,
		OutputPath: exportPath,
		StartTime:  now.Add(-24 * time.Hour),
		EndTime:    now,
		MetricName: "test_metric",
		Labels: map[string]string{
			"service": "test",
		},
	}

	// Export metrics
	outputPath, err := retention.ExportMetrics(options)
	require.NoError(t, err)
	assert.Equal(t, exportPath, outputPath)

	// Verify exported file exists
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)

	// Read and verify CSV content
	file, err := os.Open(outputPath)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	// Verify header
	assert.Equal(t, []string{"Timestamp", "Value", "env", "region", "service"}, records[0])
	
	// Verify we have the expected number of rows (header + 10 data rows)
	assert.Equal(t, 11, len(records))
}

func TestExportMetricsInvalidFormat(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "metric-export-test-invalid")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create retention config
	config := RetentionConfig{
		Enabled:          true,
		StoragePath:      tempDir,
		RetentionPeriod:  24 * time.Hour,
		SnapshotInterval: 1 * time.Second,
	}

	// Create metric retention
	retention := NewMetricRetention(config)
	require.NotNil(t, retention)

	// Create export options with invalid format
	exportPath := filepath.Join(tempDir, "export.xyz")
	options := ExportOptions{
		Format:     ExportFormat("xyz"),
		OutputPath: exportPath,
		StartTime:  time.Now().Add(-24 * time.Hour),
		EndTime:    time.Now(),
		MetricName: "test_metric",
	}

	// Export metrics should fail
	_, err = retention.ExportMetrics(options)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported export format")
}

func TestExportMetricsNoData(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "metric-export-test-no-data")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create retention config
	config := RetentionConfig{
		Enabled:          true,
		StoragePath:      tempDir,
		RetentionPeriod:  24 * time.Hour,
		SnapshotInterval: 1 * time.Second,
	}

	// Create metric retention
	retention := NewMetricRetention(config)
	require.NotNil(t, retention)

	// Create export options for non-existent metric
	exportPath := filepath.Join(tempDir, "export.csv")
	options := ExportOptions{
		Format:     ExportFormatCSV,
		OutputPath: exportPath,
		StartTime:  time.Now().Add(-24 * time.Hour),
		EndTime:    time.Now(),
		MetricName: "nonexistent_metric",
	}

	// Export metrics should fail due to no data
	_, err = retention.ExportMetrics(options)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no metrics found")
}
