package dashboard

import (
	"github.com/nessi-dev/nessi/pkg/monitoring"
)

// mockMonitor implements a minimal mock for the monitoring.Monitor interface
// that can be used in tests
type mockMonitor struct {
	metricsPort int
}

// GetMetricsPort returns the metrics port
func (m *mockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

// RecordMetric implements the Monitor interface
func (m *mockMonitor) RecordMetric(name string, value float64, labels map[string]string) {
	// No-op for testing
}

// ExportMetrics implements the Monitor interface
func (m *mockMonitor) ExportMetrics(options monitoring.ExportOptions) (string, error) {
	return "", nil
}

// NewMockMonitor creates a new mock monitor for testing
func NewMockMonitor(metricsPort int) *monitoring.Monitor {
	options := monitoring.MonitorOptions{
		MetricsPort: metricsPort,
	}
	monitor, _ := monitoring.New(options)
	return monitor
}
