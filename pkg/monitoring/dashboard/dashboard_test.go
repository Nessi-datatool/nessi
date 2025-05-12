package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMonitor is a mock implementation of the Monitor interface for testing
type mockMonitor struct {
	metricsPort int
}

func (m *mockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

func TestDashboard(t *testing.T) {
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

	// Test index handler
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	dash.handleIndex(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, w.Body.String(), "Nessi Monitoring Dashboard")
}

func TestDashboardNotFound(t *testing.T) {
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

	// Test non-existent path
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()
	dash.handleIndex(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestMetricsHandler(t *testing.T) {
	// Skip this test in automated testing environments
	// as it requires a running metrics server
	t.Skip("Requires a running metrics server")

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

	// Test metrics handler
	req := httptest.NewRequest(http.MethodGet, "/api/metrics?metric=table_size", nil)
	w := httptest.NewRecorder()
	dash.handleMetrics(w, req)

	// We can't assert much here since it depends on an external server
	// Just check that the handler doesn't panic
}

func TestAlertsHandler(t *testing.T) {
	// Skip this test in automated testing environments
	// as it requires a running metrics server
	t.Skip("Requires a running metrics server")

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

	// Test alerts handler
	req := httptest.NewRequest(http.MethodGet, "/api/alerts", nil)
	w := httptest.NewRecorder()
	dash.handleAlerts(w, req)

	// We can't assert much here since it depends on an external server
	// Just check that the handler doesn't panic
}
