package monitoring

import (
	"testing"
	"time"
	"strings"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Metric represents a monitoring metric
type Metric struct {
	Name   string
	Value  float64
	Labels map[string]string
}

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Errorf("New() returned nil")
	}
	// Override config path for tests
	m.configPath = "pkg/monitoring/config/monitoring_test.json"
	// Load test config
	if err := m.LoadConfig(); err != nil {
		t.Errorf("Failed to load test config: %v", err)
	}
}

func TestRecordTableMetrics(t *testing.T) {
	m := New()
	m.configPath = "pkg/monitoring/config/monitoring_test.json"

	// Record metrics without starting the HTTP server
	m.RecordTableMetrics("test_table", []map[string]interface{}{{"size": 1024, "count": 100}})

	// Verify metrics using testutil
	expected := `# HELP nessi_table_size_bytes Size of Delta tables in bytes
# TYPE nessi_table_size_bytes gauge
nessi_table_size_bytes{table_name="test_table"} 1024
# HELP nessi_record_count Number of records in Delta tables
# TYPE nessi_record_count gauge
nessi_record_count{table_name="test_table"} 100
`
	if err := testutil.CollectAndCompare(m.collector, strings.NewReader(expected)); err != nil {
		t.Errorf("Metric comparison failed: %v", err)
	}
}

func TestRecordLatency(t *testing.T) {
	m := New()
	m.configPath = "pkg/monitoring/config/monitoring_test.json"
	if err := m.LoadConfig(); err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

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
	m.configPath = "pkg/monitoring/config/monitoring_test.json"
	if err := m.LoadConfig(); err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Record violation without starting the HTTP server
	m.RecordRuleViolation("test_table", "not_null")

	// Verify violation count
	metricFamilies, err := m.collector.Gather()
	if err != nil {
		t.Errorf("Failed to gather metrics: %v", err)
		return
	}

	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "nessi_rule_violations_total" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var ruleName, tableName string
				for _, label := range labels {
					if label.GetName() == "rule_name" {
						ruleName = label.GetValue()
					} else if label.GetName() == "table_name" {
						tableName = label.GetValue()
					}
				}
				if ruleName == "not_null" && tableName == "test_table" {
					if metric.GetCounter().GetValue() != 1 {
						t.Errorf("Expected violation count to be 1, got %v", metric.GetCounter().GetValue())
					}
					found = true
				}
			}
		}
	}
	if !found {
		t.Errorf("Expected rule violation metric not found in metrics")
	}
}

func TestSendAlert(t *testing.T) {
	m := New()
	m.configPath = "pkg/monitoring/config/monitoring_test.json"
	if err := m.LoadConfig(); err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}
	m.SendAlert("test_alert", "test_table", 100)
	// Verify alert was sent
	select {
	case alert := <-m.alerts:
		if alert.Name != "test_alert" {
			t.Errorf("Expected alert name 'test_alert', got '%s'", alert.Name)
		}
		if alert.Message != "Metric test_alert exceeded threshold for table test_table (value: 100)" {
			t.Errorf("Expected alert message containing '100', got '%s'", alert.Message)
		}
	default:
		t.Errorf("No alert was sent")
	}
}

func TestAlert(t *testing.T) {
	m := New()
	m.configPath = "pkg/monitoring/config/monitoring_test.json"
	if err := m.LoadConfig(); err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}
	m.alertThresholds["test_alert"] = AlertThreshold{
		Warning:  50,
		Critical: 100,
	}
	m.alerts = make(chan Alert, 1)
	go m.processAlerts([]Alert{{
		Name:        "test_alert",
		Severity:    "CRITICAL",
		Message:     "Test alert message",
		Timestamp:   time.Now(),
		Metadata:    map[string]string{"table_name": "test_table"},
	}}, m.alertThresholds)
	select {
	case alert := <-m.alerts:
		if alert.Name != "test_alert" {
			t.Errorf("Expected alert name 'test_alert', got '%s'", alert.Name)
		}
		if alert.Severity != "CRITICAL" {
			t.Errorf("Expected alert severity 'CRITICAL', got '%s'", alert.Severity)
		}
	default:
		t.Errorf("No alert was processed")
	}
}

func TestRecordMetric(t *testing.T) {
	m := New()
	m.configPath = "pkg/monitoring/config/monitoring_test.json"
	if err := m.LoadConfig(); err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Record a table metric
	m.RecordTableMetrics("test_metric", []map[string]interface{}{{"size": 100}})

	// Verify metric
	metricFamilies, err := m.collector.Gather()
	if err != nil {
		t.Errorf("Failed to gather metrics: %v", err)
		return
	}

	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "nessi_table_size_bytes" {
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var tableName string
				for _, label := range labels {
					if label.GetName() == "table_name" {
						tableName = label.GetValue()
					}
				}
				if tableName == "test_metric" {
					if metric.GetGauge().GetValue() != 100 {
						t.Errorf("Expected metric value to be 100, got %v", metric.GetGauge().GetValue())
					}
					found = true
				}
			}
		}
	}
	if !found {
		t.Errorf("Expected metric not found in metrics")
	}
}
