package testutil

import (
	"context"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/testutil"
)

// GetTestTimeout returns an appropriate timeout duration for tests
func GetTestTimeout() time.Duration {
	return testutil.GetTestTimeout()
}

// RunWithTimeout runs a test function with a timeout
func RunWithTimeout(t *testing.T, testFunc func(), timeout ...time.Duration) {
	testutil.RunWithTimeout(t, testFunc, timeout...)
}

// RunTestWithContext runs a test function with a context that has a timeout
func RunTestWithContext(t *testing.T, testFunc func(ctx context.Context)) {
	testutil.RunTestWithContext(t, testFunc)
}

// RunInParallel marks a test to run in parallel with other tests
// WARNING: Do not use this function if the test already calls t.Parallel() directly,
// as this will cause conflicts and potentially lead to test failures.
// Instead, either use this function OR call t.Parallel() directly, but not both.
func RunInParallel(t *testing.T) {
	t.Helper() // Mark as test helper for better error reporting
	
	// Check if this function is being called from a test that might also call t.Parallel()
	t.Log("WARNING: Using testutil.RunInParallel() - do not also call t.Parallel() directly in the same test")
	
	// Mark the test as parallel
	t.Parallel()
}

// CreateMockMonitor creates a mock monitor for testing
func CreateMockMonitor(metricsPort int) *MockMonitor {
	return &MockMonitor{
		metricsPort: metricsPort,
		metrics:     make(map[string]float64),
	}
}

// MockMonitor is a simplified monitor implementation for testing
type MockMonitor struct {
	metricsPort  int
	metrics      map[string]float64
	alertManager interface{}
}

// RecordMetric records a metric in the mock monitor
func (m *MockMonitor) RecordMetric(name string, value float64, labels map[string]string) error {
	m.metrics[name] = value
	return nil
}

// GetMetricsPort returns the metrics port
func (m *MockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

// GetAlertManager returns the alert manager
func (m *MockMonitor) GetAlertManager() interface{} {
	return m.alertManager
}

// RecordMetricWithTimestamp records a metric with a timestamp
func (m *MockMonitor) RecordMetricWithTimestamp(name string, value float64, timestamp time.Time, labels map[string]string) error {
	return m.RecordMetric(name, value, labels)
}

// ExportMetrics exports metrics to a file
func (m *MockMonitor) ExportMetrics(options monitoring.ExportOptions) (string, error) {
	return options.OutputPath, nil
}

// GetIntelligentAlertManager returns the intelligent alert manager
func (m *MockMonitor) GetIntelligentAlertManager() interface{} {
	return nil
}

// ShouldSkipIntegrationTests returns true if integration tests should be skipped
func ShouldSkipIntegrationTests(t *testing.T) bool {
	return testutil.ShouldSkipIntegrationTests(t)
}
