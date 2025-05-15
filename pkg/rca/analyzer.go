package rca

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
)

// Config represents the configuration for the RCA engine
type Config struct {
	EnableRCA           bool   `json:"enable_rca" yaml:"enable_rca"`
	MaxHistoryDays      int    `json:"max_history_days" yaml:"max_history_days"`
	DetailLevel         string `json:"detail_level" yaml:"detail_level"` // basic, standard, detailed
	IncludeLineage      bool   `json:"include_lineage" yaml:"include_lineage"`
	IncludeSchemaChange bool   `json:"include_schema_change" yaml:"include_schema_change"`
	AlertThreshold      int    `json:"alert_threshold" yaml:"alert_threshold"` // Severity level to trigger alerts
}

// DefaultConfig returns the default configuration for the RCA engine
func DefaultConfig() *Config {
	return &Config{
		EnableRCA:           false,
		MaxHistoryDays:      7,
		DetailLevel:         "standard",
		IncludeLineage:      true,
		IncludeSchemaChange: true,
		AlertThreshold:      2, // 0=info, 1=warning, 2=error, 3=critical
	}
}

// Analyzer is responsible for performing root cause analysis on anomalies
type Analyzer struct {
	config           *Config
	monitoringClient *monitoring.Client
	deltaConnector   *datalake.Connector
}

// NewAnalyzer creates a new RCA analyzer
func NewAnalyzer(config *Config, monitoringClient *monitoring.Client, deltaConnector *datalake.Connector) *Analyzer {
	if config == nil {
		config = DefaultConfig()
	}
	return &Analyzer{
		config:           config,
		monitoringClient: monitoringClient,
		deltaConnector:   deltaConnector,
	}
}

// AnomalyInfo represents the basic information about an anomaly
type AnomalyInfo struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Metric      string    `json:"metric"`
	Value       float64   `json:"value"`
	Threshold   float64   `json:"threshold"`
	Severity    int       `json:"severity"` // 0=info, 1=warning, 2=error, 3=critical
	Description string    `json:"description"`
	TablePath   string    `json:"table_path"`
	ColumnName  string    `json:"column_name,omitempty"`
}

// RootCause represents a potential root cause for an anomaly
type RootCause struct {
	Type        string                 `json:"type"` // schema_change, data_quality, lineage, system, unknown
	Confidence  float64                `json:"confidence"` // 0.0-1.0
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// RCAResult represents the result of a root cause analysis
type RCAResult struct {
	AnomalyID      string      `json:"anomaly_id"`
	AnalysisTime   time.Time   `json:"analysis_time"`
	PrimaryRootCause *RootCause  `json:"primary_root_cause"`
	OtherCauses    []*RootCause `json:"other_causes,omitempty"`
	AffectedTables []string    `json:"affected_tables,omitempty"`
	RelatedAnomalies []string  `json:"related_anomalies,omitempty"`
	RecommendedActions []string `json:"recommended_actions,omitempty"`
	GrafanaDashboardURL string  `json:"grafana_dashboard_url,omitempty"`
}

// AnalyzeAnomaly performs root cause analysis on a specific anomaly
func (a *Analyzer) AnalyzeAnomaly(anomalyID string) (*RCAResult, error) {
	if !a.config.EnableRCA {
		return nil, fmt.Errorf("RCA is disabled in configuration")
	}

	// Get anomaly details
	anomaly, err := a.getAnomalyInfo(anomalyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get anomaly info: %w", err)
	}

	// Initialize RCA result
	result := &RCAResult{
		AnomalyID:    anomalyID,
		AnalysisTime: time.Now(),
	}

	// Collect potential root causes
	causes := []*RootCause{}

	// Check for schema changes
	if a.config.IncludeSchemaChange {
		schemaChanges, err := a.analyzeSchemaChanges(anomaly)
		if err == nil && len(schemaChanges) > 0 {
			causes = append(causes, schemaChanges...)
		}
	}

	// Check for data quality issues
	qualityIssues, err := a.analyzeDataQualityIssues(anomaly)
	if err == nil && len(qualityIssues) > 0 {
		causes = append(causes, qualityIssues...)
	}

	// Check for lineage issues
	if a.config.IncludeLineage {
		lineageIssues, err := a.analyzeLineageIssues(anomaly)
		if err == nil && len(lineageIssues) > 0 {
			causes = append(causes, lineageIssues...)
		}
	}

	// Check for system issues
	systemIssues, err := a.analyzeSystemIssues(anomaly)
	if err == nil && len(systemIssues) > 0 {
		causes = append(causes, systemIssues...)
	}

	// If no causes found, add an unknown cause
	if len(causes) == 0 {
		causes = append(causes, &RootCause{
			Type:        "unknown",
			Confidence:  0.5,
			Description: "No specific root cause could be determined",
			Timestamp:   anomaly.Timestamp,
		})
	}

	// Sort causes by confidence and set primary cause
	sortRootCausesByConfidence(causes)
	result.PrimaryRootCause = causes[0]
	if len(causes) > 1 {
		result.OtherCauses = causes[1:]
	}

	// Find affected tables
	result.AffectedTables = a.findAffectedTables(anomaly)

	// Find related anomalies
	result.RelatedAnomalies = a.findRelatedAnomalies(anomaly)

	// Generate recommended actions
	result.RecommendedActions = a.generateRecommendations(anomaly, result.PrimaryRootCause)

	// Generate Grafana dashboard URL if applicable
	result.GrafanaDashboardURL = a.generateGrafanaDashboardURL(anomaly)

	return result, nil
}

// ToJSON converts the RCA result to a JSON string
func (r *RCAResult) ToJSON() (string, error) {
	bytes, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToHTML generates an HTML representation of the RCA result
func (r *RCAResult) ToHTML() string {
	// Implementation would generate HTML for dashboard rendering
	// This is a simplified placeholder
	html := fmt.Sprintf(`
<div class="rca-result">
  <h2>Root Cause Analysis for Anomaly %s</h2>
  <div class="primary-cause">
    <h3>Primary Root Cause</h3>
    <p><strong>Type:</strong> %s</p>
    <p><strong>Confidence:</strong> %.1f%%</p>
    <p><strong>Description:</strong> %s</p>
  </div>
`, r.AnomalyID, r.PrimaryRootCause.Type, r.PrimaryRootCause.Confidence*100, r.PrimaryRootCause.Description)

	if len(r.OtherCauses) > 0 {
		html += "<div class=\"other-causes\"><h3>Other Potential Causes</h3><ul>"
		for _, cause := range r.OtherCauses {
			html += fmt.Sprintf("<li>%s (%.1f%% confidence)</li>", cause.Description, cause.Confidence*100)
		}
		html += "</ul></div>"
	}

	if len(r.RecommendedActions) > 0 {
		html += "<div class=\"recommendations\"><h3>Recommended Actions</h3><ul>"
		for _, action := range r.RecommendedActions {
			html += fmt.Sprintf("<li>%s</li>", action)
		}
		html += "</ul></div>"
	}

	if r.GrafanaDashboardURL != "" {
		html += fmt.Sprintf("<p><a href=\"%s\" target=\"_blank\">View in Grafana</a></p>", r.GrafanaDashboardURL)
	}

	html += "</div>"
	return html
}

// Helper methods (would be implemented with actual logic)
func (a *Analyzer) getAnomalyInfo(anomalyID string) (*AnomalyInfo, error) {
	// This would fetch the anomaly from the monitoring system
	// Placeholder implementation
	return &AnomalyInfo{
		ID:          anomalyID,
		Timestamp:   time.Now().Add(-1 * time.Hour),
		Metric:      "null_percentage",
		Value:       0.15,
		Threshold:   0.05,
		Severity:    2,
		Description: "High percentage of NULL values detected",
		TablePath:   "sales/transactions",
		ColumnName:  "customer_id",
	}, nil
}

func (a *Analyzer) analyzeSchemaChanges(anomaly *AnomalyInfo) ([]*RootCause, error) {
	// This would analyze recent schema changes that might have caused the anomaly
	// Placeholder implementation
	return []*RootCause{
		{
			Type:        "schema_change",
			Confidence:  0.8,
			Description: "Column type change detected prior to anomaly",
			Timestamp:   anomaly.Timestamp.Add(-2 * time.Hour),
			Details: map[string]interface{}{
				"previous_type": "string",
				"current_type":  "integer",
				"column_name":   anomaly.ColumnName,
				"change_author": "data_pipeline_job",
			},
		},
	}, nil
}

func (a *Analyzer) analyzeDataQualityIssues(anomaly *AnomalyInfo) ([]*RootCause, error) {
	// This would analyze data quality issues that might have caused the anomaly
	// Placeholder implementation
	return []*RootCause{
		{
			Type:        "data_quality",
			Confidence:  0.6,
			Description: "Upstream data quality rule failures detected",
			Timestamp:   anomaly.Timestamp.Add(-3 * time.Hour),
			Details: map[string]interface{}{
				"failed_rules": []string{"valid_customer_id", "non_empty_customer_id"},
				"table_path":   "raw/customer_data",
			},
		},
	}, nil
}

func (a *Analyzer) analyzeLineageIssues(anomaly *AnomalyInfo) ([]*RootCause, error) {
	// This would analyze lineage issues that might have caused the anomaly
	// Placeholder implementation
	return []*RootCause{
		{
			Type:        "lineage",
			Confidence:  0.7,
			Description: "Upstream table refresh failure",
			Timestamp:   anomaly.Timestamp.Add(-4 * time.Hour),
			Details: map[string]interface{}{
				"upstream_table": "raw/customer_data",
				"job_id":         "daily_customer_refresh_20250514",
				"error_message":  "Timeout waiting for source system response",
			},
		},
	}, nil
}

func (a *Analyzer) analyzeSystemIssues(anomaly *AnomalyInfo) ([]*RootCause, error) {
	// This would analyze system issues that might have caused the anomaly
	// Placeholder implementation
	return []*RootCause{
		{
			Type:        "system",
			Confidence:  0.4,
			Description: "High system load during data processing",
			Timestamp:   anomaly.Timestamp.Add(-2 * time.Hour),
			Details: map[string]interface{}{
				"cpu_utilization": 95.2,
				"memory_usage":    87.8,
				"disk_io":         "high",
			},
		},
	}, nil
}

func (a *Analyzer) findAffectedTables(anomaly *AnomalyInfo) []string {
	// This would find tables affected by the anomaly
	// Placeholder implementation
	return []string{
		anomaly.TablePath,
		"reporting/customer_metrics",
		"analytics/customer_segmentation",
	}
}

func (a *Analyzer) findRelatedAnomalies(anomaly *AnomalyInfo) []string {
	// This would find anomalies related to the current one
	// Placeholder implementation
	return []string{
		"anom-20250514-002",
		"anom-20250514-007",
	}
}

func (a *Analyzer) generateRecommendations(anomaly *AnomalyInfo, primaryCause *RootCause) []string {
	// This would generate recommendations based on the anomaly and primary cause
	// Placeholder implementation
	var recommendations []string

	switch primaryCause.Type {
	case "schema_change":
		recommendations = []string{
			"Review recent schema changes to table " + anomaly.TablePath,
			"Verify data type compatibility for column " + anomaly.ColumnName,
			"Check for missing data transformation steps in ETL pipeline",
		}
	case "data_quality":
		recommendations = []string{
			"Investigate data quality issues in upstream table " + primaryCause.Details["table_path"].(string),
			"Review and fix failed validation rules",
			"Consider adding data quality checks earlier in the pipeline",
		}
	case "lineage":
		recommendations = []string{
			"Check status of upstream data refresh job " + primaryCause.Details["job_id"].(string),
			"Verify connectivity to source system",
			"Review error logs for the failed job",
		}
	case "system":
		recommendations = []string{
			"Investigate system resource constraints during processing window",
			"Consider optimizing resource-intensive queries",
			"Review system scaling policies",
		}
	default:
		recommendations = []string{
			"Monitor the metric for further anomalies",
			"Review recent changes to data pipeline",
			"Check for seasonal patterns in the data",
		}
	}

	return recommendations
}

func (a *Analyzer) generateGrafanaDashboardURL(anomaly *AnomalyInfo) string {
	// This would generate a URL to a Grafana dashboard for the anomaly
	// Placeholder implementation
	return fmt.Sprintf("https://grafana.example.com/d/nessi-anomaly/anomaly-details?var-anomaly_id=%s&from=%d&to=%d",
		anomaly.ID,
		anomaly.Timestamp.Add(-12*time.Hour).Unix()*1000,
		anomaly.Timestamp.Add(2*time.Hour).Unix()*1000)
}

func sortRootCausesByConfidence(causes []*RootCause) {
	// This would sort root causes by confidence (descending)
	// Simplified implementation
	for i := 0; i < len(causes)-1; i++ {
		for j := i + 1; j < len(causes); j++ {
			if causes[i].Confidence < causes[j].Confidence {
				causes[i], causes[j] = causes[j], causes[i]
			}
		}
	}
}
