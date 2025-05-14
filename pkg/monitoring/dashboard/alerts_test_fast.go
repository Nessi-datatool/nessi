package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertsAPIFast uses alertsMockMonitor from test_helpers.go

// TestAlertsAPIFast is a fast version of the alerts API test
func TestAlertsAPIFast(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alerts-api-test-fast")
	require.NoError(t, err)
	
	// Clean up the temporary directory
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	
	// Create a mock alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Set a strict timeout for this test
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	
	// Create a simplified mock server with fast response times
	mockMetricsServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type for all responses
		w.Header().Set("Content-Type", "application/json")
		
		// Use a simplified response pattern for faster tests
		response := map[string]interface{}{
			"id":      "test-alert-id-123",
			"status":  "active",
			"message": "Success",
		}
		
		// Success response for all requests
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	
	// Register cleanup for the mock server
	t.Cleanup(func() {
		mockMetricsServer.Close()
	})
	
	// Create a test monitor with our alert manager
	mock := testutil.CreateTestMonitorWithAlertManager(9090, alertManager)
	
	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
	}
	
	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)
	
	// Test a simple alert API call
	t.Run("GetAlerts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/alerts", nil)
		w := httptest.NewRecorder()
		
		dash.handleAlertsAPI(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	
	// Test alert rules API
	t.Run("GetAlertRules", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/alerts/rules", nil)
		w := httptest.NewRecorder()
		
		dash.handleAlertRulesAPI(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	
	// Check for context timeout
	select {
	case <-ctx.Done():
		t.Error("Test timed out")
	default:
		// Test completed within timeout
	}
}
