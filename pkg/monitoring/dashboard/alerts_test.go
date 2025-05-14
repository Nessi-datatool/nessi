
package dashboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAlertsAPI uses alertsMockMonitor from test_helpers.go

// TestAlertsAPI tests the alerts API endpoints
func TestAlertsAPI(t *testing.T) {
	t.Parallel()
	
	// We'll store the rule ID created during the test here
	var createdRuleID string
	
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
		
		// For rule creation, generate a rule ID
		if strings.Contains(r.URL.Path, "/rules") && r.Method == http.MethodPost {
			createdRuleID = fmt.Sprintf("rule-%d", time.Now().UnixNano())
			response["id"] = createdRuleID
		}
		
		// Success response for all requests
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	
	// Register cleanup for the mock server
	t.Cleanup(func() {
		mockMetricsServer.Close()
	})
	
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alerts-api-test")
	require.NoError(t, err)
	
	// Clean up the temporary directory
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})
	
	// Create a mock alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Extract port from mockMetricsServer.URL
	mockServerURL := mockMetricsServer.URL
	mockServerPort := 0
	if strings.HasPrefix(mockServerURL, "http://") {
		parts := strings.Split(strings.TrimPrefix(mockServerURL, "http://"), ":")
		if len(parts) == 2 {
			portStr := parts[1]
			if idx := strings.Index(portStr, "/"); idx != -1 {
				portStr = portStr[:idx]
			}
			mockServerPort, _ = strconv.Atoi(portStr)
		}
	}

	// Create a test monitor with our alert manager and the mock server's port
	mock := testutil.CreateTestMonitorWithAlertManager(mockServerPort, alertManager)
	
	// Create dashboard
	opts := DashboardOptions{
		ListenAddr: ":8080",
		Profiler:      &MockProfiler{},
		RuleValidator: &MockRuleValidator{},
	}

	dash, err := New(mock, opts)
	require.NoError(t, err)
	require.NotNil(t, dash)
	
	// Test alerts API
	t.Run("GetAlerts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/alerts", nil)
		w := httptest.NewRecorder()
		
		dash.handleAlertsAPI(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Unexpected status: %d, body: %s", resp.StatusCode, string(body))
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
	
	// Test alert rules API
	t.Run("CreateRule", func(t *testing.T) {
		ruleReq := map[string]interface{}{
			"name":                "Test Rule",
			"description":         "This is a test rule",
			"metric":              "data_quality_score",
			"threshold":           90.0,
			"comparison_operator": "<",
			"severity":            "critical",
			"enabled":             true,
		}
		
		reqBody, err := json.Marshal(ruleReq)
		require.NoError(t, err)
		
		req := httptest.NewRequest(http.MethodPost, "/api/alerts/rules", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		dash.handleAlertRulesAPI(w, req)
		
		resp := w.Result()
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("Unexpected status: %d, body: %s", resp.StatusCode, string(body))
		}
		assert.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode)
	})
	
}
