package monitor

import (
	"sync"
	"time"
)

// SimpleMetric represents a basic metric without external dependencies
type SimpleMetric struct {
	Name        string
	Description string
	Labels      map[string]string
	Value       float64
	Timestamp   time.Time
}

// SimpleMetricsCollector provides a simple metrics collection implementation
type SimpleMetricsCollector struct {
	metrics map[string]SimpleMetric
	mutex   sync.RWMutex
}

// NewSimpleMetricsCollector creates a new SimpleMetricsCollector
func NewSimpleMetricsCollector() *SimpleMetricsCollector {
	return &SimpleMetricsCollector{
		metrics: make(map[string]SimpleMetric),
	}
}

// RecordTableTotal records the total number of tables
func (c *SimpleMetricsCollector) RecordTableTotal(catalog, schema string, count float64) {
	c.recordMetric("tables_total", "Total number of tables", count, map[string]string{
		"catalog": catalog,
		"schema":  schema,
	})
}

// RecordTableSize records the size of a table
func (c *SimpleMetricsCollector) RecordTableSize(catalog, schema, table string, sizeBytes float64) {
	c.recordMetric("table_size_bytes", "Size of table in bytes", sizeBytes, map[string]string{
		"catalog": catalog,
		"schema":  schema,
		"table":   table,
	})
}

// RecordQualityCheck records a quality check
func (c *SimpleMetricsCollector) RecordQualityCheck(catalog, schema, table, check string) {
	c.recordMetric("quality_checks_total", "Total number of quality checks", 1, map[string]string{
		"catalog": catalog,
		"schema":  schema,
		"table":   table,
		"check":   check,
	})
}

// RecordQualityViolation records a quality violation
func (c *SimpleMetricsCollector) RecordQualityViolation(catalog, schema, table, rule string) {
	c.recordMetric("quality_violations_total", "Total number of quality violations", 1, map[string]string{
		"catalog": catalog,
		"schema":  schema,
		"table":   table,
		"rule":    rule,
	})
}

// RecordQueryLatency records the latency of a query
func (c *SimpleMetricsCollector) RecordQueryLatency(catalog, schema, table, queryType string, latencyMs float64) {
	c.recordMetric("query_latency_ms", "Query latency in milliseconds", latencyMs, map[string]string{
		"catalog":    catalog,
		"schema":     schema,
		"table":      table,
		"query_type": queryType,
	})
}

// RecordError records an error
func (c *SimpleMetricsCollector) RecordError(catalog, schema, table, errorType string) {
	c.recordMetric("errors_total", "Total number of errors", 1, map[string]string{
		"catalog":    catalog,
		"schema":     schema,
		"table":      table,
		"error_type": errorType,
	})
}

// recordMetric records a metric
func (c *SimpleMetricsCollector) recordMetric(name, description string, value float64, labels map[string]string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Create a key from name and labels
	key := name
	for k, v := range labels {
		key += "_" + k + "_" + v
	}

	c.metrics[key] = SimpleMetric{
		Name:        name,
		Description: description,
		Labels:      labels,
		Value:       value,
		Timestamp:   time.Now(),
	}
}

// GetMetrics returns all metrics
func (c *SimpleMetricsCollector) GetMetrics() []SimpleMetric {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	metrics := make([]SimpleMetric, 0, len(c.metrics))
	for _, metric := range c.metrics {
		metrics = append(metrics, metric)
	}

	return metrics
}

// GetMetric returns a specific metric
func (c *SimpleMetricsCollector) GetMetric(name string, labels map[string]string) (SimpleMetric, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Create a key from name and labels
	key := name
	for k, v := range labels {
		key += "_" + k + "_" + v
	}

	metric, found := c.metrics[key]
	return metric, found
}
