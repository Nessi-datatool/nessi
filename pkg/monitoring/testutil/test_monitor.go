package testutil

import (
	"github.com/nessi-dev/nessi/pkg/monitoring"
)

// CreateOSSTestMonitor creates a monitor for testing without alert manager
func CreateOSSTestMonitor(metricsPort int) *monitoring.Monitor {
	options := monitoring.MonitorOptions{
		MetricsPort: metricsPort,
	}

	monitor, _ := monitoring.New(options)
	return monitor
}
