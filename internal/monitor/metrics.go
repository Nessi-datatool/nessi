package monitor

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsCollector collects and exposes metrics
//
//go:generate stringer -type=MetricCategory
type MetricsCollector struct {
	// Table metrics
	tablesTotal *prometheus.GaugeVec
	tablesSize  *prometheus.GaugeVec

	// Quality metrics
	qualityChecksTotal *prometheus.CounterVec
	qualityViolations  *prometheus.CounterVec

	// Performance metrics
	queryLatency *prometheus.HistogramVec

	// Error metrics
	errorsTotal *prometheus.CounterVec
}

// MetricCategory defines types of metrics
type MetricCategory string

const (
	MetricCategoryTable   MetricCategory = "table"
	MetricCategoryQuality MetricCategory = "quality"
	MetricCategoryQuery   MetricCategory = "query"
	MetricCategoryError   MetricCategory = "error"
)

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		tablesTotal: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "nessi_tables_total",
				Help: "Total number of Delta tables",
			},
			[]string{"path"},
		),
		tablesSize: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "nessi_tables_size_bytes",
				Help: "Size of Delta tables in bytes",
			},
			[]string{"table_name"},
		),
		qualityChecksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_quality_checks_total",
				Help: "Total number of quality checks performed",
			},
			[]string{"rule_type", "severity"},
		),
		qualityViolations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_quality_violations_total",
				Help: "Total number of quality violations",
			},
			[]string{"rule_type", "severity"},
		),
		queryLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "nessi_query_latency_seconds",
				Help:    "Query latency distribution",
				Buckets: []float64{0.1, 0.5, 1.0, 2.0, 5.0},
			},
			[]string{"operation"},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_errors_total",
				Help: "Total number of errors",
			},
			[]string{"type", "source"},
		),
	}
}

// RecordTableMetrics records table-related metrics
func (m *MetricsCollector) RecordTableMetrics(ctx context.Context, path string, tables int, size float64) {
	m.tablesTotal.WithLabelValues(path).Set(float64(tables))
	m.tablesSize.WithLabelValues(path).Set(size)
}

// RecordQualityCheck records a quality check
func (m *MetricsCollector) RecordQualityCheck(ruleType string, severity string, success bool) {
	m.qualityChecksTotal.WithLabelValues(ruleType, severity).Inc()
	if !success {
		m.qualityViolations.WithLabelValues(ruleType, severity).Inc()
	}
}

// RecordQueryLatency records query execution time
func (m *MetricsCollector) RecordQueryLatency(operation string, duration time.Duration) {
	m.queryLatency.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordError records an error occurrence
func (m *MetricsCollector) RecordError(errorType string, source string) {
	m.errorsTotal.WithLabelValues(errorType, source).Inc()
}
