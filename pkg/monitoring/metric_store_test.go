package monitoring

import (
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMonitor is a mock implementation of the Monitor for testing
type mockMonitor struct {
	metricNames []string
	seriesData  []PrometheusSeriesResult
}

func newMockMonitor() *mockMonitor {
	return &mockMonitor{
		metricNames: []string{
			"test_metric_1",
			"test_metric_2",
		},
		seriesData: []PrometheusSeriesResult{
			{
				Labels: map[string]string{
					"instance": "localhost:9090",
					"job":      "nessi",
				},
				Points: []PrometheusDataPoint{
					{
						Timestamp: time.Now().Add(-2 * time.Hour),
						Value:     100.0,
					},
					{
						Timestamp: time.Now().Add(-1 * time.Hour),
						Value:     110.0,
					},
					{
						Timestamp: time.Now(),
						Value:     120.0,
					},
				},
			},
		},
	}
}

func (m *mockMonitor) queryPrometheus(query string, start, end time.Time, step string) ([]PrometheusSeriesResult, error) {
	return m.seriesData, nil
}

func (m *mockMonitor) queryPrometheusMetricNames() ([]string, error) {
	return m.metricNames, nil
}

// TestPrometheusMetricStore tests the PrometheusMetricStore implementation
func TestPrometheusMetricStore(t *testing.T) {
	mockMon := newMockMonitor()
	store := &PrometheusMetricStore{
		monitor: mockMon,
	}

	// Test GetMetricNames
	names, err := store.GetMetricNames()
	require.NoError(t, err)
	assert.Equal(t, mockMon.metricNames, names)
	assert.Equal(t, 2, len(names))
	assert.Contains(t, names, "test_metric_1")
	assert.Contains(t, names, "test_metric_2")

	// Test GetMetricValues
	now := time.Now()
	start := now.Add(-3 * time.Hour)
	end := now
	dataPoints, err := store.GetMetricValues("test_metric", start, end, map[string]string{
		"job": "nessi",
	})
	require.NoError(t, err)
	assert.Equal(t, 3, len(dataPoints))

	// Check that the data points have the correct values
	assert.InDelta(t, 100.0, dataPoints[0].Value, 0.001)
	assert.InDelta(t, 110.0, dataPoints[1].Value, 0.001)
	assert.InDelta(t, 120.0, dataPoints[2].Value, 0.001)

	// Check that the labels were preserved
	assert.Equal(t, "localhost:9090", dataPoints[0].Labels["instance"])
	assert.Equal(t, "nessi", dataPoints[0].Labels["job"])
}

// TestIntegrationWithIntelligentAlertManager tests the integration of the metric store with the intelligent alert manager
func TestIntegrationWithIntelligentAlertManager(t *testing.T) {
	// Create a mock monitor
	mockMon := newMockMonitor()
	
	// Create a metric store
	store := &PrometheusMetricStore{
		monitor: mockMon,
	}
	
	// Create an alert manager
	alertManager, err := alerts.NewAlertManager()
	require.NoError(t, err)
	
	// Create an intelligent alert manager
	config := alerts.DefaultIntelligentAlertingConfig()
	iam := alerts.NewIntelligentAlertManager(alertManager, store, config)
	
	// Verify the intelligent alert manager was created successfully
	assert.NotNil(t, iam)
	
	// Test that the metric store can be used by the intelligent alert manager
	metrics, err := store.GetMetricNames()
	require.NoError(t, err)
	assert.Equal(t, 2, len(metrics))
	
	// Test that we can get metric values
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	dataPoints, err := store.GetMetricValues(metrics[0], start, now, nil)
	require.NoError(t, err)
	assert.Equal(t, 3, len(dataPoints))
}
