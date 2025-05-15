package rca

import (
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
}

func TestRCAResultToJSON(t *testing.T) {
	// Create a simple RCA result
	result := &RCAResult{
		AnomalyID:    "test-anomaly-1",
		AnalysisTime: time.Now(),
		PrimaryRootCause: &RootCause{
			Type:        "test",
			Confidence:  0.9,
			Description: "Test root cause",
			Timestamp:   time.Now(),
		},
	}
	
	// Convert to JSON
	json, err := result.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}
	
	// Basic validation of JSON output
	if json == "" {
		t.Fatal("Expected non-empty JSON string")
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
}
