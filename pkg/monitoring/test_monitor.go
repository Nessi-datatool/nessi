package monitoring

import (
	"fmt"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

// TestMonitor is a simplified version of the Monitor struct for testing
type Monitor struct {
	metricsPort              int
	alertManager             *alerts.AlertManager
	intelligentAlertManager  *alerts.IntelligentAlertManager
	metricStore              alerts.MetricStore
}

// GetMetricsPort returns the metrics server port
func (m *Monitor) GetMetricsPort() int {
	return m.metricsPort
}

// GetAlertManager returns the alert manager
func (m *Monitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

// GetIntelligentAlertManager returns the intelligent alert manager
func (m *Monitor) GetIntelligentAlertManager() *alerts.IntelligentAlertManager {
	return m.intelligentAlertManager
}

// RecordMetric records a metric value with the current timestamp
func (m *Monitor) RecordMetric(metricName string, value float64, labels map[string]string) error {
	return m.RecordMetricWithTimestamp(metricName, value, time.Now(), labels)
}

// RecordMetricWithTimestamp records a metric value with a specific timestamp
func (m *Monitor) RecordMetricWithTimestamp(metricName string, value float64, timestamp time.Time, labels map[string]string) error {
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
