// Implementation of monitoring and alerting features.

package testutil

import (
	"os"

	"github.com/nessi-dev/nessi/pkg/monitoring"
)

// CreateTestMonitor creates a monitoring.Monitor instance for testing
func CreateTestMonitor(metricsPort int) *monitoring.Monitor {
	// Create a temporary directory for the test
	_, _ = os.MkdirTemp("", "test-monitor-")

	// Create options
	options := monitoring.MonitorOptions{
		MetricsPort: metricsPort,
	}

	// Create monitor
	monitor, _ := monitoring.New(options)

	return monitor
}
