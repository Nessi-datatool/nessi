package monitoring

import (
	"fmt"
	"time"
)

// RecordMetric records a metric value with the current timestamp
func (m *Monitor) RecordMetric(metricName string, value float64, labels map[string]string) error {
	return m.RecordMetricWithTimestamp(metricName, value, time.Now(), labels)
}

// RecordMetricWithTimestamp records a metric value with a specific timestamp
func (m *Monitor) RecordMetricWithTimestamp(metricName string, value float64, timestamp time.Time, labels map[string]string) error {
	// In a real implementation, this would use the Prometheus client library to record metrics
	// For now, we'll just store the metric in memory for testing purposes

	// We'll use the timestamp, value, and labels directly

	// If we have a metric retention system, record the metric there
	if m.metricRetention != nil {
		m.metricRetention.RecordMetric(metricName, "Metric from monitoring", MetricTypeGauge, value, labels)
	}

	// Log the metric for debugging
	fmt.Printf("Recorded metric %s = %f at %s with labels %v\n",
		metricName, value, timestamp.Format(time.RFC3339), labels)

	return nil
}

// MockMetricStore is defined in mock_metric_store.go
