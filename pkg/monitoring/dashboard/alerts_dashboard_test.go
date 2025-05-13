package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertsDashboard tests the alerts dashboard page
func TestAlertsDashboard(t *testing.T) {
	// Create a mock monitor
	mock := &mockMonitor{
		metricsPort: 9090,
	}
	
	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}
	
	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)
	
	// Test alerts dashboard handler
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	w := httptest.NewRecorder()
	dash.handleAlertsDashboard(w, req)
	
	// Check response
	resp := w.Result()
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// mockMonitor implements a minimal mock for the monitoring.Monitor interface
type mockMonitor struct {
	metricsPort int
}

func (m *mockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

func (m *mockMonitor) ExportMetrics(options monitoring.ExportOptions) (string, error) {
	return options.OutputPath, nil
}
