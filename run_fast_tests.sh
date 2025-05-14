#!/bin/bash

# Script to run only specific fast tests
echo "Running optimized fast tests only..."

# Set environment variables for testing
export NESSI_SKIP_LONG_TESTS=true
export NESSI_TEST_TIMEOUT_MS=5000
export NESSI_TEST_DEBUG=true
export GO_TEST=true

# Create a temporary test file for dashboard
cat > ./pkg/monitoring/dashboard/temp_test.go << 'EOF'
package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	mock := &alertsMockMonitor{
		metricsPort:  9090,
		alertManager: alertManager,
	}
	
	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}
	
	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)
	
	// Test a simple request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	dash.handleIndex(w, req)
	
	// Check response
	resp := w.Result()
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
EOF

# Create a temporary test file for monitoring
cat > ./pkg/monitoring/temp_test.go << 'EOF'
package monitoring

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFastMonitoring is a simple test for the monitoring package
func TestFastMonitoring(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "fast-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a monitor
	options := MonitorOptions{
		MetricsPort: 9090,
	}
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Record a metric
	monitor.RecordMetric("test_metric", 123.45, nil)
	
	// Test invalid format
	exportOptions := ExportOptions{
		Format:     "invalid",
		OutputPath: filepath.Join(tempDir, "test.txt"),
	}
	
	_, err = monitor.ExportMetrics(exportOptions)
	require.Error(t, err)
}
EOF

# Run the tests
echo "Running dashboard fast test..."
go test -v ./pkg/monitoring/dashboard -run TestFastDashboard -count=1 -timeout=5s

echo "Running monitoring fast test..."
go test -v ./pkg/monitoring -run TestFastMonitoring -count=1 -timeout=5s

# Clean up temporary test files
rm ./pkg/monitoring/dashboard/temp_test.go
rm ./pkg/monitoring/temp_test.go

echo "Fast tests completed!"
