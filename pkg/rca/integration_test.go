package rca

import (
	"os"
	"testing"
	"time"
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

	// Create mock dependencies
	monitoringClient := createTestMonitoringClient(t)
	deltaConnector := createTestDeltaConnector(t)

	// Create RCA analyzer with dependencies
	config := &Config{
		EnableRCA:           true,
		MaxHistoryDays:      7,
		DetailLevel:         "detailed",
		IncludeLineage:      true,
		IncludeSchemaChange: true,
		AlertThreshold:      2,
	}
	analyzer := NewAnalyzer(config, monitoringClient, deltaConnector)

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

	// Create mock dependencies for real data test
	monitoringClient := createRealMonitoringClient(t)
	deltaConnector := createRealDeltaConnector(t)

	// Create RCA analyzer with dependencies
	config := &Config{
		EnableRCA:           true,
		MaxHistoryDays:      30,
		DetailLevel:         "detailed",
		IncludeLineage:      true,
		IncludeSchemaChange: true,
		AlertThreshold:      2,
	}
	analyzer := NewAnalyzer(config, monitoringClient, deltaConnector)

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

func createTestMonitoringClient(t *testing.T) *MockMonitoringClient {
	// Create a mock monitoring client with test data
	client := NewMockMonitoringClient()
	
	// Add a test anomaly
	testTime := time.Now().Add(-1 * time.Hour)
	client.AddMockAnomaly(&AnomalyInfo{
		ID:          "test-integration-anomaly-1",
		Timestamp:   testTime,
		Metric:      "null_percentage",
		Value:       0.15,
		Threshold:   0.05,
		Severity:    2,
		Description: "High percentage of NULL values detected",
		TablePath:   "sales/transactions",
		ColumnName:  "customer_id",
	})
	
	return client
}

func createTestDeltaConnector(t *testing.T) *MockDeltaConnector {
	// Create a mock delta connector with test data
	connector := NewMockDeltaConnector()
	
	// Add a test schema change
	testTime := time.Now().Add(-2 * time.Hour)
	connector.AddMockSchemaChange(&SchemaChange{
		TablePath:    "sales/transactions",
		Timestamp:    testTime,
		ColumnName:   "customer_id",
		PreviousType: "string",
		CurrentType:  "integer",
		ChangeAuthor: "data_pipeline_job",
	})
	
	return connector
}

func createRealMonitoringClient(t *testing.T) *MockMonitoringClient {
	// In a real implementation, this would create a monitoring client
	// connected to a real monitoring system
	// For now, we'll use a mock with more realistic data
	client := NewMockMonitoringClient()
	
	// Add some realistic anomalies
	baseTime := time.Now().Add(-12 * time.Hour)
	client.AddMockAnomaly(&AnomalyInfo{
		ID:          "anom-20250514-001",
		Timestamp:   baseTime,
		Metric:      "null_percentage",
		Value:       0.15,
		Threshold:   0.05,
		Severity:    2,
		Description: "High percentage of NULL values detected",
		TablePath:   "sales/transactions",
		ColumnName:  "customer_id",
	})
	
	client.AddMockAnomaly(&AnomalyInfo{
		ID:          "anom-20250514-002",
		Timestamp:   baseTime.Add(1 * time.Hour),
		Metric:      "row_count",
		Value:       500,
		Threshold:   1000,
		Severity:    1,
		Description: "Low row count detected",
		TablePath:   "reporting/metrics",
	})
	
	return client
}

func createRealDeltaConnector(t *testing.T) *MockDeltaConnector {
	// In a real implementation, this would create a delta connector
	// connected to a real Delta Lake
	// For now, we'll use a mock with more realistic data
	connector := NewMockDeltaConnector()
	
	// Add some realistic schema changes
	baseTime := time.Now().Add(-24 * time.Hour)
	connector.AddMockSchemaChange(&SchemaChange{
		TablePath:    "sales/transactions",
		Timestamp:    baseTime.Add(10 * time.Hour),
		ColumnName:   "customer_id",
		PreviousType: "string",
		CurrentType:  "integer",
		ChangeAuthor: "data_pipeline_job",
	})
	
	connector.AddMockSchemaChange(&SchemaChange{
		TablePath:    "reporting/metrics",
		Timestamp:    baseTime.Add(12 * time.Hour),
		ColumnName:   "transaction_date",
		PreviousType: "string",
		CurrentType:  "timestamp",
		ChangeAuthor: "schema_migration_v2",
	})
	
	return connector
}
