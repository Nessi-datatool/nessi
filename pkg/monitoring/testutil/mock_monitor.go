package testutil

import (
	"os"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

// CreateTestMonitor creates a monitoring.Monitor instance for testing
func CreateTestMonitor(metricsPort int) *monitoring.Monitor {
	// Create a temporary directory for the test
	tempDir, _ := os.MkdirTemp("", "test-monitor-")
	
	// Create alert manager
	alertManager, _ := alerts.NewAlertManager(tempDir)
	
	// Create options
	options := monitoring.MonitorOptions{
		MetricsPort:  metricsPort,
		AlertManager: alertManager,
	}
	
	// Create monitor
	monitor, _ := monitoring.New(options)
	
	return monitor
}

// CreateTestMonitorWithAlertManager creates a monitoring.Monitor instance with a specific AlertManager
func CreateTestMonitorWithAlertManager(metricsPort int, alertManager *alerts.AlertManager) *monitoring.Monitor {
	// Create options
	options := monitoring.MonitorOptions{
		MetricsPort:  metricsPort,
		AlertManager: alertManager,
	}
	
	// Create monitor
	monitor, _ := monitoring.New(options)
	
	return monitor
}
