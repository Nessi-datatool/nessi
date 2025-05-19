package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
	"github.com/nessi-dev/nessi/pkg/monitoring/freshness"
)

// FreshnessHandlerManager is an interface for freshness management in handlers
type FreshnessHandlerManager interface {
	// GetTableStatus gets the table status for a table
	GetTableStatus(tableName string) (freshness.TableFreshnessStatus, error)

	// GetAllTableStatuses lists all table statuses
	GetAllTableStatuses() ([]freshness.TableFreshnessStatus, error)

	// GetSLAConfig gets the SLA configuration for a table
	GetSLAConfig(tableName string) (freshness.SLAConfig, error)

	// GetAllSLAConfigs lists all SLA configurations
	GetAllSLAConfigs() ([]freshness.SLAConfig, error)

	// SetSLAConfig sets the SLA configuration for a table
	SetSLAConfig(config freshness.SLAConfig) error

	// DeleteSLAConfig deletes the SLA configuration for a table
	DeleteSLAConfig(tableName string) error

	// GetTableTrends gets the freshness trends for a specific table
	GetTableTrends(tableName string) (*freshness.FreshnessTrends, error)

	// GetAllTablesTrends gets the freshness trends for all tables
	GetAllTablesTrends() (*freshness.FreshnessTrends, error)
}

// handleFreshnessDashboard handles the freshness dashboard page
func (d *Dashboard) handleFreshnessDashboard(w http.ResponseWriter, r *http.Request) {
	// Check if freshness manager is available
	if d.freshnessManager == nil {
		http.Error(w, "Freshness manager not available", http.StatusInternalServerError)
		return
	}

	// Render the freshness dashboard template
	if err := d.templates.ExecuteTemplate(w, "freshness.html", map[string]interface{}{
		"Title": "Data Freshness Dashboard",
	}); err != nil {
		logging.Error("Failed to render freshness dashboard template", err)
		http.Error(w, "Failed to render dashboard", http.StatusInternalServerError)
	}
}

// handleFreshnessStatusAPI handles the API endpoint for freshness status
func (d *Dashboard) handleFreshnessStatusAPI(w http.ResponseWriter, r *http.Request) {
	// Check if freshness manager is available
	if d.freshnessManager == nil {
		http.Error(w, "Freshness manager not available", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	tableName := r.URL.Query().Get("table")

	var statuses []freshness.TableFreshnessStatus
	var err error

	if tableName != "" {
		// Get status for specific table
		status, err := d.freshnessManager.GetTableStatus(tableName)
		if err != nil {
			logging.Error("Failed to get table freshness status", err)
			http.Error(w, fmt.Sprintf("Table not found or error retrieving status: %v", err), http.StatusNotFound)
			return
		}
		statuses = []freshness.TableFreshnessStatus{status}
	} else {
		// Get status for all tables
		statuses, err = d.freshnessManager.GetAllTableStatuses()
		if err != nil {
			logging.Error("Failed to get all table freshness statuses", err)
			http.Error(w, fmt.Sprintf("Failed to get freshness statuses: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		logging.Error("Failed to encode freshness statuses", err)
	}
}

// handleFreshnessSLAAPI handles the API endpoint for SLA configurations
func (d *Dashboard) handleFreshnessSLAAPI(w http.ResponseWriter, r *http.Request) {
	// Check if freshness manager is available
	if d.freshnessManager == nil {
		http.Error(w, "Freshness manager not available", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	tableName := r.URL.Query().Get("table")

	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		// GET: Retrieve SLA configurations
		if tableName != "" {
			// Get SLA for specific table
			sla, err := d.freshnessManager.GetSLAConfig(tableName)
			if err != nil {
				logging.Error("Failed to get SLA configuration", err)
				http.Error(w, fmt.Sprintf("SLA configuration not found: %v", err), http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]freshness.SLAConfig{sla})
		} else {
			// Get all SLA configurations
			slas, err := d.freshnessManager.GetAllSLAConfigs()
			if err != nil {
				logging.Error("Failed to get all SLA configurations", err)
				http.Error(w, fmt.Sprintf("Failed to get SLA configurations: %v", err), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(slas)
		}

	case http.MethodPost:
		// POST: Create or update SLA configuration
		var slaConfig freshness.SLAConfig
		if err := json.NewDecoder(r.Body).Decode(&slaConfig); err != nil {
			logging.Error("Invalid SLA configuration format", err)
			http.Error(w, fmt.Sprintf("Invalid SLA configuration format: %v", err), http.StatusBadRequest)
			return
		}

		// Parse expected frequency from query parameters
		freqStr := r.URL.Query().Get("expected_frequency")

		if freqStr != "" {
			// Handle predefined frequencies
			switch freqStr {
			case "hourly":
				slaConfig.ExpectedFrequency = time.Hour
			case "daily":
				slaConfig.ExpectedFrequency = 24 * time.Hour
			case "weekly":
				slaConfig.ExpectedFrequency = 7 * 24 * time.Hour
			case "monthly":
				slaConfig.ExpectedFrequency = 30 * 24 * time.Hour
			default:
				// Try to parse custom duration
				duration, err := time.ParseDuration(freqStr)
				if err != nil {
					logging.Error("Failed to parse frequency", err)
					http.Error(w, fmt.Sprintf("Invalid frequency format: %v", err), http.StatusBadRequest)
					return
				}
				slaConfig.ExpectedFrequency = duration
			}
		}

		// Set the SLA configuration
		if err := d.freshnessManager.SetSLAConfig(slaConfig); err != nil {
			logging.Error("Failed to save SLA configuration", err)
			http.Error(w, fmt.Sprintf("Failed to save SLA configuration: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	case http.MethodDelete:
		// DELETE: Remove SLA configuration
		if tableName == "" {
			http.Error(w, "Table name is required", http.StatusBadRequest)
			return
		}

		if err := d.freshnessManager.DeleteSLAConfig(tableName); err != nil {
			logging.Error("Failed to delete SLA configuration", err)
			http.Error(w, fmt.Sprintf("Failed to delete SLA configuration: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleFreshnessTrendsAPI handles the API endpoint for freshness trends
func (d *Dashboard) handleFreshnessTrendsAPI(w http.ResponseWriter, r *http.Request) {
	// Check if freshness manager is available
	if d.freshnessManager == nil {
		http.Error(w, "Freshness manager not available", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	tableName := r.URL.Query().Get("table")

	// Get trends data
	var trendsData *freshness.FreshnessTrends
	var err error

	if tableName != "" && tableName != "all" {
		// Get trends for specific table
		trendsData, err = d.freshnessManager.GetTableTrends(tableName)
	} else {
		// Get trends for all tables
		trendsData, err = d.freshnessManager.GetAllTablesTrends()
	}

	if err != nil {
		logging.Error("Failed to get freshness trends", err)
		http.Error(w, fmt.Sprintf("Failed to get freshness trends: %v", err), http.StatusInternalServerError)
		return
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trendsData); err != nil {
		logging.Error("Failed to encode freshness trends", err)
	}
}

// handleFreshnessExportAPI handles the API endpoint for exporting freshness data
func (d *Dashboard) handleFreshnessExportAPI(w http.ResponseWriter, r *http.Request) {
	// Check if freshness manager is available
	if d.freshnessManager == nil {
		http.Error(w, "Freshness manager not available", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv" // Default format
	}

	// Get all table statuses
	statuses, err := d.freshnessManager.GetAllTableStatuses()
	if err != nil {
		logging.Error("Failed to get all table freshness statuses", err)
		http.Error(w, fmt.Sprintf("Failed to get freshness statuses: %v", err), http.StatusInternalServerError)
		return
	}

	// Handle different export formats
	switch format {
	case "json":
		// Export as JSON
		w.Header().Set("Content-Disposition", "attachment; filename=freshness_data.json")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)

	case "csv":
		// Export as CSV
		w.Header().Set("Content-Disposition", "attachment; filename=freshness_data.csv")
		w.Header().Set("Content-Type", "text/csv")

		// Create CSV writer
		writer := http.ResponseWriter(w)
		writer.Write([]byte("Table Name,Table Path,Last Update Time,Time Since Update (seconds),Expected Frequency (seconds),Next Expected Update,Status,Warning Threshold (%),Critical Threshold (%),Enabled\n"))

		// Write data rows
		for _, status := range statuses {
			row := fmt.Sprintf("%s,%s,%s,%d,%d,%s,%s,%d,%d,%t\n",
				status.TableName,
				status.TablePath,
				status.LastUpdateTime.Format(time.RFC3339),
				int64(status.TimeSinceUpdate.Seconds()),
				int64(status.ExpectedFrequency.Seconds()),
				status.NextExpectedUpdate.Format(time.RFC3339),
				status.Status,
				status.SLAConfig.WarningThreshold,
				status.SLAConfig.CriticalThreshold,
				status.SLAConfig.Enabled,
			)
			writer.Write([]byte(row))
		}

	default:
		http.Error(w, "Unsupported export format", http.StatusBadRequest)
		return
	}
}
