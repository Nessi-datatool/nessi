//go:build always
// +build always

package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertsDashboardPageShort is a short version of the alerts dashboard test
func TestAlertsDashboardPageShort(t *testing.T) {
	// Initialize test environment
	TestCleanup()
	
	// Set up test timeout to prevent hanging
	SetupTestTimeout(t, 500*time.Millisecond)
	
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alerts-dashboard-test-short")
	require.NoError(t, err)
	
	// Clean up the temporary directory after the test
	t.Cleanup(func() {
		CleanupTempDir(tempDir)
	})
	
	// Create a mock alert manager with minimal configuration
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create a mock monitor
	mock := &alertsMockMonitor{
		metricsPort:  9090,
		alertManager: alertManager,
	}
	
	// Create dashboard with minimal options
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}
	
	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)
	
	// Test alerts dashboard handler with minimal request
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	w := httptest.NewRecorder()
	dash.handleAlertsDashboard(w, req)
	
	// Check basic response
	resp := w.Result()
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
