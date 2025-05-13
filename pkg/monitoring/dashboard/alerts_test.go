package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMonitor implements a mock Monitor for testing
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

// TestAlertsDashboard tests the alerts dashboard page
func TestAlertsDashboard(t *testing.T) {
	// Create a mock alert manager
	alertManager, err := alerts.NewAlertManager()
	require.NoError(t, err)
	
	// Create a mock monitor
	mock := &mockMonitor{
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
	
	// Test alerts dashboard handler
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	w := httptest.NewRecorder()
	dash.handleAlertsDashboard(w, req)
	
	// Check response
	resp := w.Result()
	defer resp.Body.Close()
	
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestAlertsAPI tests the alerts API endpoints
func TestAlertsAPI(t *testing.T) {
	// Create a mock alert manager
	alertManager, err := alerts.NewAlertManager()
	require.NoError(t, err)
	
	// Create a mock monitor
	mock := &mockMonitor{
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
	
	// Test creating an alert
	t.Run("CreateAlert", func(t *testing.T) {
		// Create alert request
		alertReq := map[string]interface{}{
			"name":        "Test Alert",
			"description": "This is a test alert",
			"type":        "quality",
			"severity":    "critical",
			"source":      "test",
			"labels": map[string]string{
				"environment": "test",
			},
			"annotations": map[string]string{
				"summary": "Test alert summary",
			},
		}
		
		reqBody, err := json.Marshal(alertReq)
		require.NoError(t, err)
		
		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/alerts", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		// Call handler
		dash.handleAlerts(w, req)
		
		// Check response
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		// Decode response
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		
		// Check alert properties
		assert.Equal(t, "Test Alert", result["name"])
		assert.Equal(t, "This is a test alert", result["description"])
		assert.Equal(t, "quality", result["type"])
		assert.Equal(t, "critical", result["severity"])
		assert.Equal(t, "active", result["status"])
		assert.Equal(t, "test", result["source"])
		
		// Store alert ID for later tests
		alertID := result["id"].(string)
		
		// Test getting alerts
		t.Run("GetAlerts", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/alerts", nil)
			w := httptest.NewRecorder()
			
			dash.handleAlerts(w, req)
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			alerts := result["alerts"].([]interface{})
			assert.Len(t, alerts, 1)
			
			alert := alerts[0].(map[string]interface{})
			assert.Equal(t, alertID, alert["id"])
		})
		
		// Test acknowledging an alert
		t.Run("AcknowledgeAlert", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/alerts/"+alertID+"/acknowledge", nil)
			w := httptest.NewRecorder()
			
			dash.handleAlertAction(w, req, "acknowledge")
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			assert.Equal(t, alertID, result["id"])
			assert.Equal(t, "acknowledged", result["status"])
		})
		
		// Test resolving an alert
		t.Run("ResolveAlert", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/alerts/"+alertID+"/resolve", nil)
			w := httptest.NewRecorder()
			
			dash.handleAlertAction(w, req, "resolve")
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			assert.Equal(t, alertID, result["id"])
			assert.Equal(t, "resolved", result["status"])
		})
	})
	
	// Test alert rules API
	t.Run("AlertRules", func(t *testing.T) {
		// Create rule request
		ruleReq := map[string]interface{}{
			"name":                "Test Rule",
			"description":         "This is a test rule",
			"metric":              "data_quality_score",
			"threshold":           90.0,
			"comparison_operator": "<",
			"severity":            "critical",
			"labels": map[string]string{
				"environment": "test",
			},
			"enabled": true,
		}
		
		reqBody, err := json.Marshal(ruleReq)
		require.NoError(t, err)
		
		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/alerts/rules", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		
		// Call handler
		dash.handleAlertRules(w, req)
		
		// Check response
		resp := w.Result()
		defer resp.Body.Close()
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		// Decode response
		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)
		
		// Check rule properties
		assert.Equal(t, "Test Rule", result["name"])
		assert.Equal(t, "This is a test rule", result["description"])
		assert.Equal(t, "data_quality_score", result["metric"])
		assert.Equal(t, 90.0, result["threshold"])
		assert.Equal(t, "<", result["comparison_operator"])
		assert.Equal(t, "critical", result["severity"])
		assert.Equal(t, true, result["enabled"])
		
		// Store rule ID for later tests
		ruleID := result["id"].(string)
		
		// Test getting rules
		t.Run("GetRules", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/alerts/rules", nil)
			w := httptest.NewRecorder()
			
			dash.handleAlertRules(w, req)
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			rules := result["rules"].([]interface{})
			assert.Len(t, rules, 1)
			
			rule := rules[0].(map[string]interface{})
			assert.Equal(t, ruleID, rule["id"])
		})
		
		// Test disabling a rule
		t.Run("DisableRule", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/alerts/rules/"+ruleID+"/disable", nil)
			w := httptest.NewRecorder()
			
			dash.handleAlertRuleAction(w, req, "disable")
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			assert.Equal(t, ruleID, result["id"])
			assert.Equal(t, false, result["enabled"])
		})
		
		// Test enabling a rule
		t.Run("EnableRule", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/alerts/rules/"+ruleID+"/enable", nil)
			w := httptest.NewRecorder()
			
			dash.handleAlertRuleAction(w, req, "enable")
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			assert.Equal(t, ruleID, result["id"])
			assert.Equal(t, true, result["enabled"])
		})
		
		// Test deleting a rule
		t.Run("DeleteRule", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/alerts/rules/"+ruleID, nil)
			w := httptest.NewRecorder()
			
			dash.handleAlertRuleAction(w, req, "delete")
			
			resp := w.Result()
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			var result map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&result)
			require.NoError(t, err)
			
			assert.Equal(t, "success", result["status"])
		})
	})
}
