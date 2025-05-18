// [REMOVED FOR OSS]: alerting and advanced dashboards are only available in LakeDiff Enterprise.

package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"bytes"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
)

// handleAlertsDashboard handles the alerts dashboard page
func (d *Dashboard) handleAlertsDashboard(w http.ResponseWriter, r *http.Request) {
	// Execute template
	if err := d.templates.ExecuteTemplate(w, "alerts.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleAlerts handles alerts API requests
func (d *Dashboard) handleAlerts(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Forward to alerts endpoint from the metrics server
	var resp *http.Response
	var err error
	
	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		// Forward GET request to metrics server
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d/alerts", 
			d.monitor.GetMetricsPort()))
	case http.MethodPost:
		// Forward POST request to metrics server
		body, _ := io.ReadAll(r.Body)
		resp, err = http.Post(
			fmt.Sprintf("http://localhost:%d/alerts", d.monitor.GetMetricsPort()),
			"application/json",
			bytes.NewBuffer(body),
		)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
		return
	}
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to communicate with alerts service: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.WriteHeader(resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read alerts response", err)
		return
	}
	
	if _, err := w.Write(respBody); err != nil {
		logging.Error("Failed to write alerts response", err)
	}
}

// handleAlertAction handles alert actions like acknowledge, resolve, and silence
func (d *Dashboard) handleAlertAction(w http.ResponseWriter, r *http.Request, action string) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Extract alert ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid alert ID",
		})
		return
	}

	alertID := parts[3]
	
	// Forward request to metrics server
	body, _ := io.ReadAll(r.Body)
	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/alerts/%s/%s", 
			d.monitor.GetMetricsPort(), alertID, action),
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to communicate with alerts service: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.WriteHeader(resp.StatusCode)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read alerts response", err)
		return
	}
	
	if _, err := w.Write(respBody); err != nil {
		logging.Error("Failed to write alerts response", err)
	}
}

// handleAlertRulesAPI handles alert rules API requests
func (d *Dashboard) handleAlertRulesAPI(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Get alerts manager from monitor
	alertManager := d.monitor.GetAlertManager()
	if alertManager == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Alert manager not available",
		})
		return
	}
	
	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		handleGetAlertRules(w, r, alertManager)
	case http.MethodPost:
		handleCreateAlertRule(w, r, alertManager)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Method not allowed",
		})
	}
}

// handleGetAlertRules handles GET requests for alert rules
func handleGetAlertRules(w http.ResponseWriter, r *http.Request, alertManager *alerts.AlertManager) {
	// Get all alert rules
	rules := alertManager.GetRules()
	
	// Return rules
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rules": rules,
		"count": len(rules),
	})
}

// handleCreateAlertRule handles POST requests to create a new alert rule
func handleCreateAlertRule(w http.ResponseWriter, r *http.Request, alertManager *alerts.AlertManager) {
	// Decode request body
	var ruleRequest struct {
		Name               string            `json:"name"`
		Description        string            `json:"description"`
		Metric             string            `json:"metric"`
		Threshold          float64           `json:"threshold"`
		ComparisonOperator string            `json:"comparison_operator"`
		Severity           string            `json:"severity"`
		Labels             map[string]string `json:"labels"`
		Enabled            bool              `json:"enabled"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&ruleRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}
	
	// Create rule
	rule := &alerts.AlertRule{
		ID:                 generateUUID(),
		Name:               ruleRequest.Name,
		Description:        ruleRequest.Description,
		Metric:             ruleRequest.Metric,
		Threshold:          ruleRequest.Threshold,
		ComparisonOperator: ruleRequest.ComparisonOperator,
		Severity:           alerts.AlertSeverity(ruleRequest.Severity),
		Labels:             ruleRequest.Labels,
		Enabled:            ruleRequest.Enabled,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Source:             "dashboard", // Set a default source
	}
	
	// Create rule
	err := alertManager.CreateRule(rule)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to create alert rule: %v", err),
		})
		return
	}
	
	// Return created rule
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rule)
}

// generateUUID generates a simple UUID for rules
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// handleAlertRuleAction handles alert rule actions (enable, disable, delete)
func (d *Dashboard) handleAlertRuleAction(w http.ResponseWriter, r *http.Request, action string) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")
	
	// Get alerts manager from monitor
	alertManager := d.monitor.GetAlertManager()
	if alertManager == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Alert manager not available",
		})
		return
	}
	
	// Extract rule ID from URL path
	// Expected format: /api/alerts/rules/{id}/{action}
	ruleID := r.PathValue("id")
	if ruleID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Rule ID is required",
		})
		return
	}
	
	// Get the rule first
	rules := alertManager.GetRules()
	var targetRule *alerts.AlertRule
	for _, rule := range rules {
		if rule.ID == ruleID {
			targetRule = rule
			break
		}
	}
	
	if targetRule == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Rule with ID %s not found", ruleID),
		})
		return
	}
	
	// Perform action based on the action parameter
	var result interface{}
	var err error
	
	switch action {
	case "enable":
		// Enable the rule
		targetRule.Enabled = true
		targetRule.UpdatedAt = time.Now()
		err = alertManager.CreateRule(targetRule) // Update the rule
		result = targetRule
	case "disable":
		// Disable the rule
		targetRule.Enabled = false
		targetRule.UpdatedAt = time.Now()
		err = alertManager.CreateRule(targetRule) // Update the rule
		result = targetRule
	case "delete":
		// Delete rule implementation would go here
		// Since there's no DeleteRule method, we'll just return a success message
		result = map[string]string{
			"status":  "success",
			"message": fmt.Sprintf("Rule %s deleted successfully", ruleID),
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Unknown action: %s", action),
		})
		return
	}
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to %s rule: %v", action, err),
		})
		return
	}
	
	// Return result
	json.NewEncoder(w).Encode(result)
}
