package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics represents various monitoring metrics
type Metrics struct {
	TableSize      *prometheus.GaugeVec
	RecordCount    *prometheus.GaugeVec
	ErrorCount     *prometheus.CounterVec
	Latency        *prometheus.HistogramVec
	RuleViolations *prometheus.CounterVec
}
