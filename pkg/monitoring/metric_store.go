package monitoring

import (
	"fmt"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

// PrometheusMetricStore implements the alerts.MetricStore interface using Prometheus as a backend
type PrometheusMetricStore struct {
	monitor *Monitor
}

// NewPrometheusMetricStore creates a new PrometheusMetricStore
func NewPrometheusMetricStore(monitor *Monitor) *PrometheusMetricStore {
	return &PrometheusMetricStore{
		monitor: monitor,
	}
}

// GetMetricValues retrieves historical values for a metric from Prometheus
func (p *PrometheusMetricStore) GetMetricValues(metricName string, start, end time.Time, labels map[string]string) ([]alerts.MetricDataPoint, error) {
	// Build query with labels
	query := metricName
	if len(labels) > 0 {
		query += "{"
		first := true
		for k, v := range labels {
			if !first {
				query += ","
			}
			query += fmt.Sprintf("%s=\"%s\"", k, v)
			first = false
		}
		query += "}"
	}

	// Format time range
	step := "1h" // 1 hour resolution for historical data
	if end.Sub(start) < 24*time.Hour {
		step = "5m" // 5 minute resolution for shorter ranges
	}

	// Query Prometheus
	result, err := p.monitor.queryPrometheus(query, start, end, step)
	if err != nil {
		return nil, fmt.Errorf("failed to query Prometheus: %w", err)
	}

	// Convert to MetricDataPoints
	var dataPoints []alerts.MetricDataPoint
	for _, series := range result {
		// Extract labels
		seriesLabels := make(map[string]string)
		for k, v := range series.Labels {
			seriesLabels[k] = v
		}

		// Add data points
		for _, point := range series.Points {
			dataPoints = append(dataPoints, alerts.MetricDataPoint{
				Timestamp: point.Timestamp,
				Value:     point.Value,
				Labels:    seriesLabels,
			})
		}
	}

	return dataPoints, nil
}

// GetMetricNames returns all available metric names from Prometheus
func (p *PrometheusMetricStore) GetMetricNames() ([]string, error) {
	// Query Prometheus for all metric names
	result, err := p.monitor.queryPrometheusMetricNames()
	if err != nil {
		return nil, fmt.Errorf("failed to query Prometheus metric names: %w", err)
	}

	return result, nil
}

// PrometheusSeriesResult represents a time series result from Prometheus
type PrometheusSeriesResult struct {
	Labels map[string]string
	Points []PrometheusDataPoint
}

// PrometheusDataPoint represents a single data point from Prometheus
type PrometheusDataPoint struct {
	Timestamp time.Time
	Value     float64
}

// queryPrometheus queries Prometheus for time series data
func (m *Monitor) queryPrometheus(query string, start, end time.Time, step string) ([]PrometheusSeriesResult, error) {
	// This is a simplified implementation that would normally use the Prometheus HTTP API
	// In a real implementation, this would make an HTTP request to the Prometheus server
	
	// For now, return some mock data for testing
	result := []PrometheusSeriesResult{
		{
			Labels: map[string]string{
				"instance": "localhost:9090",
				"job":      "nessi",
			},
			Points: []PrometheusDataPoint{
				{
					Timestamp: time.Now().Add(-1 * time.Hour),
					Value:     100.0,
				},
				{
					Timestamp: time.Now(),
					Value:     110.0,
				},
			},
		},
	}

	return result, nil
}

// queryPrometheusMetricNames queries Prometheus for all available metric names
func (m *Monitor) queryPrometheusMetricNames() ([]string, error) {
	// This is a simplified implementation that would normally use the Prometheus HTTP API
	// In a real implementation, this would make an HTTP request to the Prometheus server
	
	// For now, return some mock data for testing
	return []string{
		"nessi_data_quality_score",
		"nessi_record_count",
		"nessi_table_size_bytes",
		"nessi_errors_total",
		"nessi_operation_latency_seconds",
		"nessi_api_requests_total",
	}, nil
}
