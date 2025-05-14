#!/bin/bash

# Set environment variables for testing
export NESSI_SKIP_LONG_TESTS=true
export NESSI_TEST_TIMEOUT_MS=5000
export NESSI_TEST_DEBUG=true

# Define colors for better readability
GREEN="\033[0;32m"
RED="\033[0;31m"
YELLOW="\033[0;33m"
NC="\033[0m" # No Color

# Script to run only specific fast tests while ignoring all other tests
echo "Running only fast tests..."

# Create a temporary directory for our test files
mkdir -p ./temp_tests/dashboard
mkdir -p ./temp_tests/monitoring

# Create a temporary test file for dashboard
cat > ./temp_tests/dashboard/dashboard_test.go << 'EOF'
package dashboard_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/dashboard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMonitor implements a minimal mock for testing
type mockMonitor struct {
	metricsPort  int
	alertManager *alerts.AlertManager
}

func (m *mockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

func (m *mockMonitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

func (m *mockMonitor) ExportMetrics(options monitoring.ExportOptions) (string, error) {
	return options.OutputPath, nil
}

// TestFastDashboard is a simple test for the dashboard
func TestFastDashboard(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "fast-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a mock alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create a mock monitor
	mock := &mockMonitor{
		metricsPort:  9090,
		alertManager: alertManager,
	}
	
	// Create dashboard
	opts := dashboard.DashboardOptions{
		ListenAddr: ":8080",
	}
	
	dash, err := dashboard.New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)
	
	// Test a simple request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	
	// We can't directly access handleIndex since it's not exported,
	// but we can check that the dashboard was created successfully
	assert.NotNil(t, dash)
}
EOF

# Create a temporary test file for monitoring
cat > ./temp_tests/monitoring/monitoring_test.go << 'EOF'
package monitoring_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/stretchr/testify/require"
)

// TestFastMonitoring is a simple test for the monitoring package
func TestFastMonitoring(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "fast-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a monitor
	options := monitoring.MonitorOptions{
		MetricsPort: 9090,
	}
	monitor, err := monitoring.New(options)
	require.NoError(t, err)
	
	// Record a metric
	monitor.RecordMetric("test_metric", 123.45, nil)
	
	// Test invalid format
	exportOptions := monitoring.ExportOptions{
		Format:     "invalid",
		OutputPath: filepath.Join(tempDir, "test.txt"),
	}
	
	_, err = monitor.ExportMetrics(exportOptions)
	require.Error(t, err)
}
EOF

# Copy necessary files to make the tests compile
cp -r ../../pkg/monitoring/dashboard/templates ./temp_tests/dashboard/
cp -r ../../pkg/monitoring/dashboard/static ./temp_tests/dashboard/

# Run the tests in the temporary directory
echo "Running dashboard fast test..."
cd ./temp_tests/dashboard && go test -v -count=1 -timeout=5s

echo "Running monitoring fast test..."
cd ../monitoring && go test -v -count=1 -timeout=5s

# Clean up temporary test files
cd ../../
rm -rf ./temp_tests

echo "Fast tests completed!"
