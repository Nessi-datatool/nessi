package monitoring

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
)

// Metric represents a monitoring metric
type Metric struct {
	Name   string
	Value  float64
	Labels map[string]string
}

func TestNew(t *testing.T) {
	m := New()
	assert.NotNil(t, m)
}

func TestRecordTableMetrics(t *testing.T) {
	m := New()

	// Record metrics without starting the HTTP server
	m.RecordTableMetrics("test_table", 1024, 100)

	// Verify metrics using testutil
	expected := strings.NewReader(`
		# HELP nessi_table_size_bytes Size of Delta tables in bytes
		# TYPE nessi_table_size_bytes gauge
		# HELP nessi_record_count Number of records in Delta tables
		# TYPE nessi_record_count gauge
		nessi_table_size_bytes{table_name="test_table"} 1024
		nessi_record_count{table_name="test_table"} 100
	`)
	assert.NoError(t, testutil.CollectAndCompare(m.collector, expected))
}

func TestRecordLatency(t *testing.T) {
	m := New()

	// Record latency without starting the HTTP server
	start := time.Now()
	time.Sleep(10 * time.Millisecond)
	m.RecordLatency("read", time.Since(start))

	// Verify latency
	metricFamilies, err := m.collector.Gather()
	require.NoError(t, err)

	for _, mf := range metricFamilies {
		if mf.GetName() == "nessi_operation_latency_seconds" {
			for _, metric := range mf.GetMetric() {
				for _, label := range metric.GetLabel() {
					if label.GetValue() == "read" {
						assert.Greater(t, metric.GetHistogram().GetSampleSum(), float64(0))
					}
				}
			}
		}
	}
}

func TestRecordRuleViolation(t *testing.T) {
	m := New()

	// Record violation without starting the HTTP server
	m.RecordRuleViolation("test_table", "not_null")

	// Verify violation count
	metricFamilies, err := m.collector.Gather()
	require.NoError(t, err)

	for _, mf := range metricFamilies {
		if mf.GetName() == "nessi_rule_violations_total" {
			for _, metric := range mf.GetMetric() {
				for _, label := range metric.GetLabel() {
					if label.GetValue() == "not_null" && label.GetName() == "rule_name" {
						for _, label2 := range metric.GetLabel() {
							if label2.GetValue() == "test_table" && label2.GetName() == "table_name" {
								assert.Equal(t, float64(1), metric.GetCounter().GetValue())
							}
						}
					}
				}
			}
		}
	}
}

func TestSendAlert(t *testing.T) {
	m := New()
	alert := Alert{
		Name:        "test_alert",
		Severity:    "critical",
		Message:     "threshold exceeded",
		Timestamp:   time.Now(),
		Metadata:    map[string]string{"env": "test"},
	}
	m.SendAlert(alert)

	// Verify alert was sent
	select {
	case got := <-m.alerts:
		assert.Equal(t, "test_alert", got.Name)
		assert.Equal(t, "critical", got.Severity)
		assert.Equal(t, "threshold exceeded", got.Message)
	}
	m.Stop()
}

func TestGetCollector(t *testing.T) {
	m := New()
	collector := m.GetCollector()
	assert.NotNil(t, collector)
}

func TestAlert(t *testing.T) {
	alert := Alert{
		Name:        "test_alert",
		Severity:    "critical",
		Message:     "threshold exceeded",
		Timestamp:   time.Now(),
		Metadata:    map[string]string{"env": "test"},
	}

	assert.Equal(t, "test_alert", alert.Name)
	assert.Equal(t, "critical", alert.Severity)
	assert.Equal(t, "threshold exceeded", alert.Message)
	assert.Equal(t, map[string]string{"env": "test"}, alert.Metadata)
}

func TestRecordMetric(t *testing.T) {
	m := New()
	
	// Record a table metric using the existing method
	m.RecordTableMetrics("test_table", 42.0, 10)

	// Verify metric was recorded
	assert.NotNil(t, m.metrics.TableSize, "TableSize metric should not be nil")
}
