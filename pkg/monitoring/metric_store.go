package monitoring

import (
	"fmt"
	"time"
)

// MetricStore implements the alerts.MetricStore interface
type MetricStore struct {
	gauges     map[string]interface{}
	counters   map[string]interface{}
	histograms map[string]interface{}
}

// NewMetricStore creates a new MetricStore
func NewMetricStore() *MetricStore {
	return &MetricStore{
		gauges:     make(map[string]interface{}),
		counters:   make(map[string]interface{}),
		histograms: make(map[string]interface{}),
	}
}

// RegisterGauge registers a new gauge metric
func (m *MetricStore) RegisterGauge(name, help string) interface{} {
	// In a real implementation, this would register a metric
	// For simplicity, we'll just store a placeholder
	m.gauges[name] = name
	return name
}

// RegisterCounter registers a new counter metric
func (m *MetricStore) RegisterCounter(name, help string) interface{} {
	// In a real implementation, this would register a metric
	// For simplicity, we'll just store a placeholder
	m.counters[name] = name
	return name
}

// RegisterHistogram registers a new histogram metric
func (m *MetricStore) RegisterHistogram(name, help string) interface{} {
	// In a real implementation, this would register a metric
	// For simplicity, we'll just store a placeholder
	m.histograms[name] = name
	return name
}

// GetCounterValue returns the current value of a counter
func (m *MetricStore) GetCounterValue(counter interface{}) float64 {
	// In a real implementation, this would get the value from a metric store
	// For simplicity, we'll just return a placeholder value
	return 0.0
}

// RecordMetric records a metric value
func (m *MetricStore) RecordMetric(name string, value float64, labels ...map[string]string) {
	// In a real implementation, this would update a metric store
	// For now, we'll implement a basic version that works for tests
	gaugeKey := name
	if len(labels) > 0 && labels[0] != nil {
		// Create a key that includes labels
		gaugeKey = name + "{"
		first := true
		for k, v := range labels[0] {
			if !first {
				gaugeKey += ","
			}
			gaugeKey += fmt.Sprintf("%s=\"%s\"", k, v)
			first = false
		}
		gaugeKey += "}"
	}

	gauge, exists := m.gauges[gaugeKey]
	if !exists {
		gauge = m.RegisterGauge(gaugeKey, "Auto-registered gauge for "+name)
		m.gauges[gaugeKey] = gauge
	}

	// In a real implementation, this would set the gauge value
	// For testing, we'll just store the value in the gauge name for now
	m.gauges[gaugeKey] = value
}

// RecordMetricWithTimestamp records a metric value with a timestamp
func (m *MetricStore) RecordMetricWithTimestamp(name string, value float64, timestamp time.Time, labels ...map[string]string) {
	// In a real implementation, this would update a metric store with timestamp
	// For now, we'll just call RecordMetric since we don't need timestamps for tests
	m.RecordMetric(name, value, labels...)
}

// GetMetric retrieves a metric value
func (m *MetricStore) GetMetric(name string, labels ...map[string]string) (float64, error) {
	// In a real implementation, this would query a metric store
	// For testing, we'll return the stored value if it exists
	gaugeKey := name
	if len(labels) > 0 && labels[0] != nil {
		// Create a key that includes labels
		gaugeKey = name + "{"
		first := true
		for k, v := range labels[0] {
			if !first {
				gaugeKey += ","
			}
			gaugeKey += fmt.Sprintf("%s=\"%s\"", k, v)
			first = false
		}
		gaugeKey += "}"
	}

	if gauge, exists := m.gauges[gaugeKey]; exists {
		switch v := gauge.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		default:
			return 0.0, nil
		}
	}
	return 0.0, fmt.Errorf("metric %s not found", gaugeKey)
}

// GetMetricHistory retrieves historical values for a metric
func (m *MetricStore) GetMetricHistory(name string, start, end time.Time, labels ...map[string]string) ([]float64, []time.Time, error) {
	// In a real implementation, this would query a metric store for historical data
	// For testing, we'll return some mock data
	value, err := m.GetMetric(name, labels...)
	if err != nil {
		return []float64{}, []time.Time{}, err
	}

	// Generate some test data points
	values := []float64{value * 0.9, value * 0.95, value, value * 1.05, value * 1.1}
	times := []time.Time{
		start,
		start.Add(time.Duration(int(end.Sub(start)) / 4)),
		start.Add(time.Duration(int(end.Sub(start)) / 2)),
		start.Add(time.Duration(int(end.Sub(start)) * 3 / 4)),
		end,
	}

	return values, times, nil
}

// TableSize returns the size of the metrics table
func (m *MetricStore) TableSize() int {
	// In a real implementation, this would query a metric store
	// For simplicity, we'll just return the number of registered metrics
	return len(m.gauges) + len(m.counters) + len(m.histograms)
}
