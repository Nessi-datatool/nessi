package monitoring

import (
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

// MockMetricStore is a mock implementation of the alerts.MetricStore interface for testing
type MockMetricStore struct {
	metrics map[string][]alerts.MetricDataPoint
}

// NewMockMetricStore creates a new MockMetricStore
func NewMockMetricStore() *MockMetricStore {
	return &MockMetricStore{
		metrics: make(map[string][]alerts.MetricDataPoint),
	}
}

// GetMetricValues retrieves historical values for a metric
func (m *MockMetricStore) GetMetricValues(metricName string, start, end time.Time, labels map[string]string) ([]alerts.MetricDataPoint, error) {
	// Get all data points for the metric
	allPoints, ok := m.metrics[metricName]
	if !ok {
		return []alerts.MetricDataPoint{}, nil
	}
	
	// Filter by time range and labels
	var filteredPoints []alerts.MetricDataPoint
	for _, point := range allPoints {
		// Check time range
		if (point.Timestamp.Equal(start) || point.Timestamp.After(start)) &&
		   (point.Timestamp.Equal(end) || point.Timestamp.Before(end)) {
			// Check labels
			if labels == nil || matchLabels(point.Labels, labels) {
				filteredPoints = append(filteredPoints, point)
			}
		}
	}
	
	return filteredPoints, nil
}

// GetMetricNames returns all available metric names
func (m *MockMetricStore) GetMetricNames() ([]string, error) {
	names := make([]string, 0, len(m.metrics))
	for name := range m.metrics {
		names = append(names, name)
	}
	return names, nil
}

// AddMetric adds a metric data point to the store
func (m *MockMetricStore) AddMetric(metricName string, dataPoint alerts.MetricDataPoint) {
	m.metrics[metricName] = append(m.metrics[metricName], dataPoint)
}

// matchLabels checks if a set of labels matches a filter
func matchLabels(labels, filter map[string]string) bool {
	for k, v := range filter {
		if labels[k] != v {
			return false
		}
	}
	return true
}
