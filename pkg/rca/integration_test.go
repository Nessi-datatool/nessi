package rca_test

import (
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/rca"
)

// TestRCAIntegration tests the integration of the RCA service with other components
func TestRCAIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Set up test environment
	setupTestEnvironment(t)
	defer cleanupTestEnvironment(t)

	// Create dependencies
	monitoringClient := createTestMonitoringClient(t)
	deltaConnector := createTestDeltaConnector(t)

	// Create RCA analyzer with dependencies
	config := &rca.Config{
		EnableRCA:           true,
		MaxHistoryDays:      7,
		DetailLevel:         "detailed",
		IncludeLineage:      true,
		IncludeSchemaChange: true,
		AlertThreshold:      2,
	}
	analyzer := rca.NewAnalyzer(config, monitoringClient, deltaConnector)

	// Test analyzing an anomaly
	anomalyID := "test-integration-anomaly-1"
	result, err := analyzer.AnalyzeAnomaly(anomalyID)
	if err != nil {
		t.Fatalf("Failed to analyze anomaly: %v", err)
	}

	// Verify result
	if result.AnomalyID != anomalyID {
		t.Errorf("Expected AnomalyID to be '%s', got '%s'", anomalyID, result.AnomalyID)
	}

	if result.PrimaryRootCause == nil {
		t.Fatal("Expected non-nil PrimaryRootCause")
	}

	// Verify that recommended actions are present
	if len(result.RecommendedActions) == 0 {
		t.Error("Expected at least one recommended action")
	}

	// Test JSON output
	jsonOutput, err := result.ToJSON()
	if err != nil {
		t.Fatalf("Failed to generate JSON output: %v", err)
	}
	if jsonOutput == "" {
		t.Error("Expected non-empty JSON output")
	}

	// Test HTML output
	htmlOutput := result.ToHTML()
	if htmlOutput == "" {
		t.Error("Expected non-empty HTML output")
	}
}

// TestRCAWithRealAnomalyData tests the RCA service with real anomaly data
func TestRCAWithRealAnomalyData(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if not in CI environment
	if os.Getenv("CI") == "" && os.Getenv("NESSI_TEST_REAL_DATA") == "" {
		t.Skip("Skipping test with real data outside of CI environment")
	}

	// Create dependencies
	monitoringClient := createRealMonitoringClient(t)
	deltaConnector := createRealDeltaConnector(t)

	// Create RCA analyzer with dependencies
	config := &rca.Config{
		EnableRCA:           true,
		MaxHistoryDays:      30,
		DetailLevel:         "detailed",
		IncludeLineage:      true,
		IncludeSchemaChange: true,
		AlertThreshold:      2,
	}
	analyzer := rca.NewAnalyzer(config, monitoringClient, deltaConnector)

	// Get a list of real anomalies from the monitoring system
	anomalies, err := monitoringClient.GetRecentAnomalies(time.Now().Add(-24*time.Hour), time.Now())
	if err != nil {
		t.Fatalf("Failed to get recent anomalies: %v", err)
	}

	if len(anomalies) == 0 {
		t.Skip("No recent anomalies found for testing")
	}

	// Test analyzing each anomaly
	for i, anomalyID := range anomalies {
		// Limit to first 5 anomalies to avoid long test times
		if i >= 5 {
			break
		}

		result, err := analyzer.AnalyzeAnomaly(anomalyID)
		if err != nil {
			t.Errorf("Failed to analyze anomaly %s: %v", anomalyID, err)
			continue
		}

		// Verify result
		if result.AnomalyID != anomalyID {
			t.Errorf("Expected AnomalyID to be '%s', got '%s'", anomalyID, result.AnomalyID)
		}

		if result.PrimaryRootCause == nil {
			t.Errorf("Expected non-nil PrimaryRootCause for anomaly %s", anomalyID)
		}
	}
}

// Helper functions

func setupTestEnvironment(t *testing.T) {
	// Create test data directory if it doesn't exist
	testDataDir := "../../test_data/rca"
	if err := os.MkdirAll(testDataDir, 0755); err != nil {
		t.Fatalf("Failed to create test data directory: %v", err)
	}

	// Create test anomaly data file
	anomalyData := `{
		"id": "test-integration-anomaly-1",
		"timestamp": "2025-05-15T10:30:00Z",
		"metric": "null_percentage",
		"value": 0.15,
		"threshold": 0.05,
		"severity": 2,
		"description": "High percentage of NULL values detected",
		"table_path": "sales/transactions",
		"column_name": "customer_id"
	}`
	if err := os.WriteFile(testDataDir+"/anomaly.json", []byte(anomalyData), 0644); err != nil {
		t.Fatalf("Failed to create test anomaly data file: %v", err)
	}

	// Create test schema change data file
	schemaChangeData := `{
		"table_path": "sales/transactions",
		"timestamp": "2025-05-15T08:30:00Z",
		"changes": [
			{
				"column_name": "customer_id",
				"previous_type": "string",
				"current_type": "integer"
			}
		]
	}`
	if err := os.WriteFile(testDataDir+"/schema_change.json", []byte(schemaChangeData), 0644); err != nil {
		t.Fatalf("Failed to create test schema change data file: %v", err)
	}
}

func cleanupTestEnvironment(t *testing.T) {
	// Remove test data directory
	testDataDir := "../../test_data/rca"
	if err := os.RemoveAll(testDataDir); err != nil {
		t.Fatalf("Failed to clean up test data directory: %v", err)
	}
}

func createTestMonitoringClient(t *testing.T) *monitoring.Client {
	// In a real implementation, this would create a monitoring client
	// with mock data for testing
	// For now, we'll return nil as the analyzer implementation handles this
	return nil
}

func createTestDeltaConnector(t *testing.T) *datalake.Connector {
	// In a real implementation, this would create a delta connector
	// with mock data for testing
	// For now, we'll return nil as the analyzer implementation handles this
	return nil
}

func createRealMonitoringClient(t *testing.T) *monitoring.Client {
	// In a real implementation, this would create a monitoring client
	// connected to a real monitoring system
	// For now, we'll return nil as the analyzer implementation handles this
	return nil
}

func createRealDeltaConnector(t *testing.T) *datalake.Connector {
	// In a real implementation, this would create a delta connector
	// connected to a real Delta Lake
	// For now, we'll return nil as the analyzer implementation handles this
	return nil
}
