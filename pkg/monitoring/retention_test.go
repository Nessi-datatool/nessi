package monitoring

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricRetention(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "metric-retention-test")
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

	// Start retention service
	err = retention.Start()
	require.NoError(t, err)
	defer retention.Stop()

	// Record some metrics
	retention.RecordMetric("test_gauge", "Test gauge metric", MetricTypeGauge, 42.0, map[string]string{
		"service": "test",
		"env":     "dev",
	})

	retention.RecordMetric("test_counter", "Test counter metric", MetricTypeCounter, 10.0, map[string]string{
		"service": "test",
		"env":     "dev",
	})

	// Get metric series
	series, ok := retention.GetMetricSeries("test_gauge")
	assert.True(t, ok)
	assert.Equal(t, "test_gauge", series.Name)
	assert.Equal(t, "Test gauge metric", series.Description)
	assert.Equal(t, MetricTypeGauge, series.Type)
	assert.Len(t, series.Values, 1)
	assert.Equal(t, 42.0, series.Values[0].Value)
	assert.Equal(t, "test", series.Values[0].Labels["service"])
	assert.Equal(t, "dev", series.Values[0].Labels["env"])

	// Get metrics with labels
	values := retention.GetMetricSeriesWithLabels("test_gauge", map[string]string{
		"service": "test",
	})
	assert.Len(t, values, 1)
	assert.Equal(t, 42.0, values[0].Value)

	// Get metrics with non-matching labels
	values = retention.GetMetricSeriesWithLabels("test_gauge", map[string]string{
		"service": "nonexistent",
	})
	assert.Len(t, values, 0)

	// Wait for snapshot to be taken
	time.Sleep(2 * time.Second)

	// Check if snapshot file was created
	files, err := filepath.Glob(filepath.Join(tempDir, "metrics_*.json"))
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(files), 1)

	// Record more metrics after snapshot
	retention.RecordMetric("test_gauge", "Test gauge metric", MetricTypeGauge, 84.0, map[string]string{
		"service": "test",
		"env":     "prod",
	})

	// Get current snapshot
	snapshot := retention.GetMetricsSnapshot()
	assert.NotNil(t, snapshot)
	assert.Len(t, snapshot.Metrics, 1)

	// Test cleanup of old snapshots
	// Create an old snapshot file
	oldTimestamp := time.Now().Add(-48 * time.Hour)
	oldFilename := filepath.Join(tempDir, fmt.Sprintf("metrics_%s.json", oldTimestamp.Format("20060102_150405")))
	err = os.WriteFile(oldFilename, []byte("{}"), 0644)
	require.NoError(t, err)

	// Wait for cleanup
	time.Sleep(2 * time.Second)

	// Check if old snapshot was removed
	_, err = os.Stat(oldFilename)
	assert.True(t, os.IsNotExist(err))

	// Test GetMetricHistory
	start := time.Now().Add(-1 * time.Hour)
	end := time.Now().Add(1 * time.Hour)
	history, err := retention.GetMetricHistory("test_gauge", map[string]string{
		"service": "test",
	}, start, end)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(history), 1)

	// Test GetLatestMetrics
	latestMetrics, err := retention.GetLatestMetrics()
	require.NoError(t, err)
	assert.NotNil(t, latestMetrics)
}

func TestMetricRetentionWithoutStart(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "metric-retention-test-no-start")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create retention config
	config := RetentionConfig{
		Enabled:          true,
		StoragePath:      tempDir,
		RetentionPeriod:  24 * time.Hour,
		SnapshotInterval: 1 * time.Second,
	}

	// Create metric retention without starting it
	retention := NewMetricRetention(config)
	require.NotNil(t, retention)

	// Record some metrics
	retention.RecordMetric("test_gauge", "Test gauge metric", MetricTypeGauge, 42.0, map[string]string{
		"service": "test",
	})

	// Get metric series
	series, ok := retention.GetMetricSeries("test_gauge")
	assert.True(t, ok)
	assert.Equal(t, "test_gauge", series.Name)
	assert.Len(t, series.Values, 1)

	// No snapshots should be created since the service wasn't started
	files, err := filepath.Glob(filepath.Join(tempDir, "metrics_*.json"))
	require.NoError(t, err)
	assert.Len(t, files, 0)
}
