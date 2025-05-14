package tests

import (
	"fmt"
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

// TestMonitor is a simplified version of the Monitor struct for testing
type TestMonitor struct {
	metricsPort              int
	alertManager             *alerts.AlertManager
	intelligentAlertManager  *alerts.IntelligentAlertManager
	metricStore              alerts.MetricStore
}

// GetMetricsPort returns the metrics server port
func (m *TestMonitor) GetMetricsPort() int {
	return m.metricsPort
}

// GetAlertManager returns the alert manager
func (m *TestMonitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

// GetIntelligentAlertManager returns the intelligent alert manager
func (m *TestMonitor) GetIntelligentAlertManager() *alerts.IntelligentAlertManager {
	return m.intelligentAlertManager
}

// RecordMetric records a metric value with the current timestamp
func (m *TestMonitor) RecordMetric(metricName string, value float64, labels map[string]string) error {
	return m.RecordMetricWithTimestamp(metricName, value, time.Now(), labels)
}

// RecordMetricWithTimestamp records a metric value with a specific timestamp
func (m *TestMonitor) RecordMetricWithTimestamp(metricName string, value float64, timestamp time.Time, labels map[string]string) error {
	// Create a metric data point
	dataPoint := alerts.MetricDataPoint{
		Timestamp: timestamp,
		Value:     value,
		Labels:    labels,
	}
	
	// Store the data point in the metric store
	if m.metricStore != nil {
		if mockStore, ok := m.metricStore.(*MockMetricStore); ok {
			mockStore.AddMetric(metricName, dataPoint)
			return nil
		}
	}
	
	// If we're not using a mock metric store, just log the metric
	fmt.Printf("Recorded metric %s = %f at %s with labels %v\n", 
		metricName, value, timestamp.Format(time.RFC3339), labels)
	
	return nil
}
