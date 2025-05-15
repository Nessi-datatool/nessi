package rca

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestNewAnalyzer(t *testing.T) {
	// Test with nil config
	analyzer := NewAnalyzer(nil, nil, nil)
	if analyzer == nil {
		t.Fatal("Expected non-nil analyzer with nil config")
	}
	if analyzer.config == nil {
		t.Fatal("Expected default config to be used")
	}
	
	// Verify default config values
	defaultConfig := DefaultConfig()
	if analyzer.config.EnableRCA != defaultConfig.EnableRCA ||
	   analyzer.config.MaxHistoryDays != defaultConfig.MaxHistoryDays ||
	   analyzer.config.DetailLevel != defaultConfig.DetailLevel ||
	   analyzer.config.IncludeLineage != defaultConfig.IncludeLineage ||
	   analyzer.config.IncludeSchemaChange != defaultConfig.IncludeSchemaChange ||
	   analyzer.config.AlertThreshold != defaultConfig.AlertThreshold {
		t.Fatal("Default config values not set correctly")
	}
	
	// Test with custom config
	customConfig := &Config{
		EnableRCA:           true,
		MaxHistoryDays:      14,
		DetailLevel:         "detailed",
		IncludeLineage:      false,
		IncludeSchemaChange: false,
		AlertThreshold:      3,
	}
	
	analyzer = NewAnalyzer(customConfig, nil, nil)
	if analyzer.config != customConfig {
		t.Fatal("Expected custom config to be used")
	}
	
	// Verify custom config values
	if analyzer.config.EnableRCA != true ||
	   analyzer.config.MaxHistoryDays != 14 ||
	   analyzer.config.DetailLevel != "detailed" ||
	   analyzer.config.IncludeLineage != false ||
	   analyzer.config.IncludeSchemaChange != false ||
	   analyzer.config.AlertThreshold != 3 {
		t.Fatal("Custom config values not set correctly")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	// Verify default values
	if config.EnableRCA != false {
		t.Errorf("Expected EnableRCA to be false, got %v", config.EnableRCA)
	}
	
	if config.MaxHistoryDays != 7 {
		t.Errorf("Expected MaxHistoryDays to be 7, got %d", config.MaxHistoryDays)
	}
	
	if config.DetailLevel != "standard" {
		t.Errorf("Expected DetailLevel to be 'standard', got '%s'", config.DetailLevel)
	}
	
	if config.IncludeLineage != true {
		t.Errorf("Expected IncludeLineage to be true, got %v", config.IncludeLineage)
	}
	
	if config.IncludeSchemaChange != true {
		t.Errorf("Expected IncludeSchemaChange to be true, got %v", config.IncludeSchemaChange)
	}
	
	if config.AlertThreshold != 2 {
		t.Errorf("Expected AlertThreshold to be 2, got %d", config.AlertThreshold)
	}
}

func TestRCAResultToJSON(t *testing.T) {
	// Create a sample RCA result
	testTime := time.Date(2025, 5, 15, 10, 30, 0, 0, time.UTC)
	result := &RCAResult{
		AnomalyID:    "test-anomaly-1",
		AnalysisTime: testTime,
		PrimaryRootCause: &RootCause{
			Type:        "schema_change",
			Confidence:  0.9,
			Description: "Column type change detected",
			Timestamp:   testTime.Add(-2 * time.Hour),
			Details: map[string]interface{}{
				"column_name":   "customer_id",
				"previous_type": "string",
				"current_type":  "integer",
			},
		},
		OtherCauses: []*RootCause{
			{
				Type:        "data_quality",
				Confidence:  0.6,
				Description: "Data quality rule failures",
				Timestamp:   testTime.Add(-3 * time.Hour),
			},
		},
		AffectedTables: []string{"sales/transactions", "reporting/metrics"},
		RelatedAnomalies: []string{"anom-123", "anom-456"},
		RecommendedActions: []string{"Review schema changes", "Check data quality rules"},
		GrafanaDashboardURL: "https://grafana.example.com/d/123",
	}
	
	// Convert to JSON
	jsonStr, err := result.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}
	
	// Basic validation of JSON output
	if jsonStr == "" {
		t.Fatal("Expected non-empty JSON string")
	}
	
	// Verify JSON structure by parsing it back
	var parsedResult map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsedResult); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	// Check required fields
	requiredFields := []string{"anomaly_id", "analysis_time", "primary_root_cause"}
	for _, field := range requiredFields {
		if _, exists := parsedResult[field]; !exists {
			t.Errorf("Required field '%s' missing from JSON output", field)
		}
	}
	
	// Verify anomaly ID
	if id, ok := parsedResult["anomaly_id"].(string); !ok || id != "test-anomaly-1" {
		t.Errorf("Expected anomaly_id to be 'test-anomaly-1', got '%v'", parsedResult["anomaly_id"])
	}
}

func TestRCAResultToHTML(t *testing.T) {
	// Create a sample RCA result
	testTime := time.Date(2025, 5, 15, 10, 30, 0, 0, time.UTC)
	result := &RCAResult{
		AnomalyID:    "test-anomaly-1",
		AnalysisTime: testTime,
		PrimaryRootCause: &RootCause{
			Type:        "schema_change",
			Confidence:  0.9,
			Description: "Column type change detected",
			Timestamp:   testTime.Add(-2 * time.Hour),
		},
		RecommendedActions: []string{"Review schema changes"},
		GrafanaDashboardURL: "https://grafana.example.com/d/123",
	}
	
	// Convert to HTML
	html := result.ToHTML()
	
	// Basic validation of HTML output
	if html == "" {
		t.Fatal("Expected non-empty HTML string")
	}
	
	// Check for required HTML elements
	requiredElements := []string{
		"<div class=\"rca-result\"", 
		"<h2>Root Cause Analysis for Anomaly test-anomaly-1</h2>",
		"<div class=\"primary-cause\">",
		"<p><strong>Type:</strong> schema_change</p>",
		"<p><strong>Confidence:</strong> 90.0%</p>",
	}
	
	for _, element := range requiredElements {
		if !strings.Contains(html, element) {
			t.Errorf("HTML output missing required element: %s", element)
		}
	}
	
	// Check for Grafana link
	expectedLink := fmt.Sprintf("<a href=\"%s\" target=\"_blank\">View in Grafana</a>", result.GrafanaDashboardURL)
	if !strings.Contains(html, expectedLink) {
		t.Errorf("HTML output missing Grafana link: %s", expectedLink)
	}
}

func TestSortRootCausesByConfidence(t *testing.T) {
	// Create test data
	causes := []*RootCause{
		{Type: "A", Confidence: 0.5},
		{Type: "B", Confidence: 0.8},
		{Type: "C", Confidence: 0.3},
		{Type: "D", Confidence: 0.9},
	}
	
	// Sort
	sortRootCausesByConfidence(causes)
	
	// Verify order
	if causes[0].Type != "D" || causes[1].Type != "B" || 
	   causes[2].Type != "A" || causes[3].Type != "C" {
		t.Fatalf("Incorrect sort order: %v, %v, %v, %v", 
			causes[0].Type, causes[1].Type, causes[2].Type, causes[3].Type)
	}
	
	// Test with empty slice
	emptyCauses := []*RootCause{}
	sortRootCausesByConfidence(emptyCauses)
	// Should not panic
	
	// Test with single item
	singleCause := []*RootCause{{Type: "X", Confidence: 0.7}}
	sortRootCausesByConfidence(singleCause)
	if singleCause[0].Type != "X" {
		t.Fatalf("Single item sort failed")
	}
}

func TestAnalyzeAnomalyDisabled(t *testing.T) {
	// Create analyzer with RCA disabled
	config := DefaultConfig()
	config.EnableRCA = false
	analyzer := NewAnalyzer(config, nil, nil)
	
	// Attempt to analyze an anomaly
	_, err := analyzer.AnalyzeAnomaly("test-anomaly-1")
	
	// Should return an error
	if err == nil {
		t.Fatal("Expected error when RCA is disabled, got nil")
	}
	
	// Error message should mention that RCA is disabled
	if !strings.Contains(err.Error(), "disabled") {
		t.Errorf("Expected error message to mention RCA is disabled, got: %v", err)
	}
}
