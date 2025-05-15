package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"github.com/nessi-dev/nessi-dev/pkg/rca"
)

// RCAHandlerClient is an interface for root cause analysis in handlers
type RCAHandlerClient interface {
	// AnalyzeAnomaly performs root cause analysis on an anomaly
	AnalyzeAnomaly(anomalyID string, config *rca.Config) (*rca.RCAResult, error)
	
	// GetRecentAnalyses returns recent RCA analyses
	GetRecentAnalyses(limit int) ([]*rca.Analysis, error)
	
	// GetAnalysisResult returns a specific RCA result
	GetAnalysisResult(anomalyID string) (*rca.RCAResult, error)
	
	// GetInsights returns aggregated insights from RCA results
	GetInsights() (*rca.Insights, error)
}

// handleRcaDashboard handles the RCA dashboard page
func (d *Dashboard) handleRcaDashboard(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Render the RCA dashboard template
	if err := d.templates.ExecuteTemplate(w, "rca.html", nil); err != nil {
		logging.Error("Failed to render RCA dashboard template", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleRcaAPI handles RCA API requests
func (d *Dashboard) handleRcaAPI(w http.ResponseWriter, r *http.Request) {
	// Set content type for all responses
	w.Header().Set("Content-Type", "application/json")
	
	// Check if RCA client is available
	if d.rcaClient == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "RCA service is not available",
		})
		return
	}
	
	// Handle different API endpoints based on the path
	path := r.URL.Path
	
	switch {
	case path == "/api/rca/recent":
		d.handleRecentRcaAPI(w, r)
	case path == "/api/rca/analyze":
		d.handleAnalyzeRcaAPI(w, r)
	case path == "/api/rca/insights":
		d.handleRcaInsightsAPI(w, r)
	case path == "/api/rca/export":
		d.handleRcaExportAPI(w, r)
	default:
		// Check if it's a specific RCA result request
		var anomalyID string
		if n, err := fmt.Sscanf(path, "/api/rca/%s", &anomalyID); n == 1 && err == nil {
			if r.URL.Path == fmt.Sprintf("/api/rca/%s/export", anomalyID) {
				d.handleRcaDetailExportAPI(w, r, anomalyID)
			} else {
				d.handleRcaDetailAPI(w, r, anomalyID)
			}
			return
		}
		
		// If no matching endpoint, return 404
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "API endpoint not found",
		})
	}
}

// handleRecentRcaAPI handles requests for recent RCA results
func (d *Dashboard) handleRecentRcaAPI(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Get limit from query parameter, default to 10
	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
			limit = 10
		}
	}
	
	// Get recent analyses
	results, err := d.rcaClient.GetRecentAnalyses(limit)
	if err != nil {
		logging.Error("Failed to get recent RCA analyses", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get recent RCA results: %v", err),
		})
		return
	}
	
	// Return results
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// handleAnalyzeRcaAPI handles requests to perform RCA
func (d *Dashboard) handleAnalyzeRcaAPI(w http.ResponseWriter, r *http.Request) {
	// Only support POST requests
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Parse request body
	var request struct {
		AnomalyID string `json:"anomaly_id"`
		Format    string `json:"format"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Invalid request body: %v", err),
		})
		return
	}
	
	// Validate anomaly ID
	if request.AnomalyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Anomaly ID is required",
		})
		return
	}
	
	// Create RCA configuration
	config := &rca.Config{
		// Use basic configuration
	}
	
	// Perform RCA
	result, err := d.rcaClient.AnalyzeAnomaly(request.AnomalyID, config)
	if err != nil {
		logging.Error("Failed to analyze anomaly", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to analyze anomaly: %v", err),
		})
		return
	}
	
	// Return result
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// handleRcaDetailAPI handles requests for a specific RCA result
func (d *Dashboard) handleRcaDetailAPI(w http.ResponseWriter, r *http.Request, anomalyID string) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Get RCA result
	result, err := d.rcaClient.GetAnalysisResult(anomalyID)
	if err != nil {
		logging.Error("Failed to get RCA analysis result", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get RCA result: %v", err),
		})
		return
	}
	
	// If result not found
	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("RCA result not found for anomaly ID: %s", anomalyID),
		})
		return
	}
	
	// Return result
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// handleRcaDetailExportAPI handles requests to export a specific RCA result
func (d *Dashboard) handleRcaDetailExportAPI(w http.ResponseWriter, r *http.Request, anomalyID string) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Get RCA result
	result, err := d.rcaClient.GetAnalysisResult(anomalyID)
	if err != nil {
		logging.Error("Failed to get RCA analysis result", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get RCA result: %v", err),
		})
		return
	}
	
	// If result not found
	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("RCA result not found for anomaly ID: %s", anomalyID),
		})
		return
	}
	
	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=rca_%s.json", anomalyID))
	w.Header().Set("Content-Type", "application/json")
	
	// Return result
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// handleRcaInsightsAPI handles requests for RCA insights
func (d *Dashboard) handleRcaInsightsAPI(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Get insights
	insights, err := d.rcaClient.GetInsights()
	if err != nil {
		logging.Error("Failed to get RCA insights", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get RCA insights: %v", err),
		})
		return
	}
	
	// Return insights
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(insights)
}

// handleRcaExportAPI handles requests to export RCA data
func (d *Dashboard) handleRcaExportAPI(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Get format from query parameter, default to CSV
	format := r.URL.Query().Get("format")
	if format != "json" && format != "csv" {
		format = "csv"
	}
	
	// Get limit from query parameter, default to 100
	limit := 100
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
			limit = 100
		}
	}
	
	// Get recent analyses
	analyses, err := d.rcaClient.GetRecentAnalyses(limit)
	if err != nil {
		logging.Error("Failed to get RCA results for export", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get RCA results: %v", err),
		})
		return
	}
	
	// Set headers for file download
	filename := fmt.Sprintf("rca_export_%s.%s", time.Now().Format("20060102_150405"), format)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	if format == "json" {
		// Export as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(analyses)
	} else {
		// Export as CSV
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		
		// Write CSV header
		fmt.Fprintf(w, "AnomalyID,AnalysisTime,Status,StartTime,EndTime,ResultID\n")
		
		// Write CSV rows for analyses
		for _, analysis := range analyses {
			fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s\n",
				analysis.AnomalyID,
				analysis.StartTime.Format(time.RFC3339),
				analysis.Status,
				analysis.StartTime.Format(time.RFC3339),
				analysis.EndTime.Format(time.RFC3339),
				analysis.ResultID,
			)
		}
	}
}
