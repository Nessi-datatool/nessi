package dbt

import (
	"testing"
)

// TestDeltaTableReader tests the DeltaTableReader
func TestDeltaTableReader(t *testing.T) {
	// Create a reader
	reader := NewDeltaTableReader("/delta")

	// Test ReadTable
	t.Run("ReadTable", func(t *testing.T) {
		info, err := reader.ReadTable("test/model1")
		if err != nil {
			t.Fatalf("ReadTable() error = %v", err)
		}

		if info.Path != "/delta/test/model1" {
			t.Errorf("Expected path to be '/delta/test/model1', got '%s'", info.Path)
		}

		if info.Version != 1 {
			t.Errorf("Expected version to be 1, got %d", info.Version)
		}

		if info.Format != "parquet" {
			t.Errorf("Expected format to be 'parquet', got '%s'", info.Format)
		}

		if info.NumRows != 1000 {
			t.Errorf("Expected numRows to be 1000, got %d", info.NumRows)
		}

		if len(info.Schema) != 4 {
			t.Errorf("Expected 4 columns in schema, got %d", len(info.Schema))
		}
	})

	// Test ExecuteQuery
	t.Run("ExecuteQuery", func(t *testing.T) {
		result, err := reader.ExecuteQuery("test/model1", "SELECT * FROM test.model1")
		if err != nil {
			t.Fatalf("ExecuteQuery() error = %v", err)
		}

		if len(result.Columns) != 3 {
			t.Errorf("Expected 3 columns, got %d", len(result.Columns))
		}

		if len(result.Rows) != 3 {
			t.Errorf("Expected 3 rows, got %d", len(result.Rows))
		}
	})

	// Test GetColumnStats
	t.Run("GetColumnStats", func(t *testing.T) {
		stats, err := reader.GetColumnStats("test/model1", "id")
		if err != nil {
			t.Fatalf("GetColumnStats() error = %v", err)
		}

		if stats.DataType != "integer" {
			t.Errorf("Expected dataType to be 'integer', got '%s'", stats.DataType)
		}

		if stats.NonNullCount != 1000 {
			t.Errorf("Expected nonNullCount to be 1000, got %d", stats.NonNullCount)
		}
	})

	// Test ValidateRule
	t.Run("ValidateRule", func(t *testing.T) {
		rule := &Rule{
			Name:        "test_rule",
			Type:        "column_null",
			Description: "Test rule",
			Column:      "id",
			Threshold:   0,
			Severity:    "error",
		}

		result, err := reader.ValidateRule("test/model1", rule)
		if err != nil {
			t.Fatalf("ValidateRule() error = %v", err)
		}

		if result.Status != "passed" {
			t.Errorf("Expected status to be 'passed', got '%s'", result.Status)
		}

		if result.RuleName != "test_rule" {
			t.Errorf("Expected ruleName to be 'test_rule', got '%s'", result.RuleName)
		}
	})
}
