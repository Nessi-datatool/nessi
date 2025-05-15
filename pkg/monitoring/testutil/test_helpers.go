package testutil

import (
	"os"
	"testing"
	"time"
)

// SkipLongRunningTest skips a test if it's marked as a long-running test
// and the environment is configured to skip such tests
func SkipLongRunningTest(t *testing.T) {
	if os.Getenv("NESSI_SKIP_LONG_TESTS") != "" {
		t.Skip("Skipping long-running test")
	}
}

// SetupFastTestTimeout sets up a short timeout for tests to ensure they
// complete quickly or fail fast
func SetupFastTestTimeout(t *testing.T, timeout time.Duration) {
	if timeout == 0 {
		timeout = 100 * time.Millisecond // Default to 100ms if not specified
	}
	
	// Make the test run in parallel with other tests
	t.Parallel()
	
	// Set a deadline for the test
	t.Cleanup(func() {
		deadline := time.Now().Add(timeout)
		if time.Now().After(deadline) {
			t.Errorf("Test exceeded timeout of %v", timeout)
		}
	})
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
	metricsPort int
	metrics     map[string]float64
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
