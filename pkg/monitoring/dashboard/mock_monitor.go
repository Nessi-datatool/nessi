package dashboard

import (
	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

// mockMonitor implements a minimal mock for the monitoring.Monitor interface
// that can be used in tests
type mockMonitor struct {
	metricsPort  int
	alertManager *alerts.AlertManager
}

// GetMetricsPort returns the metrics port
func (m *mockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

// GetAlertManager returns the alert manager
func (m *mockMonitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

// GetIntelligentAlertManager returns the intelligent alert manager
func (m *mockMonitor) GetIntelligentAlertManager() *alerts.IntelligentAlertManager {
	return nil
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
func NewMockMonitor(metricsPort int, alertManager *alerts.AlertManager) *monitoring.Monitor {
	// Create a real Monitor instance with our mock values
	options := monitoring.MonitorOptions{
		MetricsPort:  metricsPort,
		AlertManager: alertManager,
	}
	
	monitor, _ := monitoring.New(options)
	return monitor
}
