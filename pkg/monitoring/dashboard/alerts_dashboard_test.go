// [REMOVED FOR OSS]: alerting and advanced dashboards are only available in LakeDiff Enterprise.

package dashboard

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/monitoring"
	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi/pkg/monitoring/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertsDashboardPage tests the alerts dashboard page
func TestAlertsDashboardPage(t *testing.T) {
	// No longer skipping tests
	// Initialize test environment
	TestCleanup()
	
	// Set up test timeout to prevent hanging
	SetupTestTimeout(t, 500*time.Millisecond)
	
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alerts-dashboard-test")
	require.NoError(t, err)
	
	// Clean up the temporary directory after the test
	t.Cleanup(func() {
		CleanupTempDir(tempDir)
	})
	
	// Create a mock alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create a test monitor with our alert manager
	mock := testutil.CreateTestMonitorWithAlertManager(9090, alertManager)
	
	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
		Profiler:      &MockProfiler{},
		RuleValidator: &MockRuleValidator{},
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

// dashboardMockMonitor implements a minimal mock for the monitoring.Monitor interface
type dashboardMockMonitor struct {
	metricsPort int
	alertManager *alerts.AlertManager
}

func (m *dashboardMockMonitor) GetMetricsPort() int {
	return m.metricsPort
}

func (m *dashboardMockMonitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

func (m *dashboardMockMonitor) ExportMetrics(options monitoring.ExportOptions) (string, error) {
	return options.OutputPath, nil
}
