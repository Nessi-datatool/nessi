package datalake

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func TestParseReadCloser(t *testing.T) {
	// Test cases
	testCases := []struct {
		name          string
		format        string
		data          string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Parse JSON data",
			format:        "json",
			data:          `[{"id": 1, "name": "Test 1"}, {"id": 2, "name": "Test 2"}]`,
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "Parse CSV data",
			format:        "csv",
			data:          "id,name\n1,Test 1\n2,Test 2",
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "Unsupported format",
			format:        "unknown",
			data:          "test data",
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "Empty JSON array",
			format:        "json",
			data:          "[]",
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "Invalid JSON",
			format:        "json",
			data:          "{invalid json",
			expectedCount: 0,
			expectError:   true,
		},
		{
			name:          "CSV with header only",
			format:        "csv",
			data:          "id,name",
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a ReadCloser from the test data
			readCloser := io.NopCloser(strings.NewReader(tc.data))
			
			// Parse the data
			result, err := ParseReadCloser(readCloser, tc.format)
			
			// Check for expected error
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}
			
			// Check for unexpected error
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			// Check result
			if result == nil {
				t.Errorf("Expected non-nil result")
				return
			}
			
			// Check record count
			if len(result.Records) != tc.expectedCount {
				t.Errorf("Expected %d records, got %d", tc.expectedCount, len(result.Records))
			}
			
			// Check count field
			if result.Count != tc.expectedCount {
				t.Errorf("Expected Count %d, got %d", tc.expectedCount, result.Count)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	// Test JSON data
	jsonData := `[
		{"id": 1, "name": "Test 1", "active": true, "score": 95.5, "created_at": "2023-01-01"},
		{"id": 2, "name": "Test 2", "active": false, "score": 87.2, "created_at": "2023-01-02"}
	]`
	
	// Create a ReadCloser
	readCloser := io.NopCloser(strings.NewReader(jsonData))
	
	// Parse JSON
	result, err := parseJSON(readCloser)
	
	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	// Check result
	if result == nil {
		t.Errorf("Expected non-nil result")
		return
	}
	
	// Check record count
	if len(result.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result.Records))
	}
	
	// Check schema
	expectedTypes := map[string]string{
		"id":         "integer",
		"name":       "string",
		"active":     "boolean",
		"score":      "float",
		"created_at": "date",
	}
	
	for field, expectedType := range expectedTypes {
		if actualType, ok := result.Schema[field]; !ok {
			t.Errorf("Field %s missing from schema", field)
		} else if actualType != expectedType {
			t.Errorf("Field %s: expected type %s, got %s", field, expectedType, actualType)
		}
	}
	
	// Check values
	if id, ok := result.Records[0]["id"].(float64); !ok || id != 1 {
		t.Errorf("Expected id 1, got %v", result.Records[0]["id"])
	}
	
	if name, ok := result.Records[0]["name"].(string); !ok || name != "Test 1" {
		t.Errorf("Expected name 'Test 1', got %v", result.Records[0]["name"])
	}
	
	if active, ok := result.Records[0]["active"].(bool); !ok || !active {
		t.Errorf("Expected active true, got %v", result.Records[0]["active"])
	}
}

func TestParseCSV(t *testing.T) {
	// Test CSV data
	csvData := "id,name,active,score,created_at\n" +
		"1,Test 1,true,95.5,2023-01-01\n" +
		"2,Test 2,false,87.2,2023-01-02"
	
	// Create a ReadCloser
	readCloser := io.NopCloser(strings.NewReader(csvData))
	
	// Parse CSV
	result, err := parseCSV(readCloser)
	
	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	
	// Check result
	if result == nil {
		t.Errorf("Expected non-nil result")
		return
	}
	
	// Check record count
	if len(result.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(result.Records))
	}
	
	// Check schema
	expectedFields := []string{"id", "name", "active", "score", "created_at"}
	for _, field := range expectedFields {
		if _, ok := result.Schema[field]; !ok {
			t.Errorf("Field %s missing from schema", field)
		}
	}
	
	// Check values
	if id, ok := result.Records[0]["id"].(int64); !ok || id != 1 {
		t.Errorf("Expected id 1, got %v (type %T)", result.Records[0]["id"], result.Records[0]["id"])
	}
	
	if name, ok := result.Records[0]["name"].(string); !ok || name != "Test 1" {
		t.Errorf("Expected name 'Test 1', got %v", result.Records[0]["name"])
	}
	
	if active, ok := result.Records[0]["active"].(bool); !ok || !active {
		t.Errorf("Expected active true, got %v", result.Records[0]["active"])
	}
}

func TestInferCSVValue(t *testing.T) {
	// Test cases
	testCases := []struct {
		input        string
		expectedType string
	}{
		{"", "nil"},
		{"null", "nil"},
		{"NA", "nil"},
		{"n/a", "nil"},
		{"123", "int64"},
		{"-456", "int64"},
		{"123.45", "float64"},
		{"-67.89", "float64"},
		{"true", "bool"},
		{"false", "bool"},
		{"yes", "bool"},
		{"no", "bool"},
		{"y", "bool"},
		{"n", "bool"},
		{"1", "int64"}, // This could be bool but our implementation prioritizes int
		{"0", "int64"}, // This could be bool but our implementation prioritizes int
		{"text", "string"},
		{"2023-01-01", "string"}, // In a real implementation, this might be detected as date
	}
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			value := inferCSVValue(tc.input)
			
			// Check type
			actualType := "nil"
			if value != nil {
				actualType = typeOf(value)
			}
			
			if actualType != tc.expectedType {
				t.Errorf("Input '%s': expected type %s, got %s", tc.input, tc.expectedType, actualType)
			}
		})
	}
}

// Helper function to get type name as string
func typeOf(v interface{}) string {
	switch v.(type) {
	case nil:
		return "nil"
	case bool:
		return "bool"
	case int:
		return "int"
	case int64:
		return "int64"
	case float64:
		return "float64"
	case string:
		return "string"
	default:
		return "unknown"
	}
}

func TestDeltaReader(t *testing.T) {
	// Create a DeltaReader
	reader := NewDeltaReader("/path/to/table", "parquet")
	
	// Test ReadAllParsed
	t.Run("ReadAllParsed", func(t *testing.T) {
		// This is a basic test since the actual implementation is a stub
		_, err := reader.ReadAllParsed()
		
		// We expect an error because ReadAll returns an empty reader
		if err == nil {
			t.Errorf("Expected error from empty reader, got nil")
		}
	})
	
	// Test ReadPartitionParsed
	t.Run("ReadPartitionParsed", func(t *testing.T) {
		// This is a basic test since the actual implementation is a stub
		_, err := reader.ReadPartitionParsed("year=2023")
		
		// We expect an error because ReadPartition returns an empty reader
		if err == nil {
			t.Errorf("Expected error from empty reader, got nil")
		}
	})
	
	// Test GetSchema
	t.Run("GetSchema", func(t *testing.T) {
		schema, err := reader.GetSchema()
		
		// Check for errors
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		
		// Check schema
		if len(schema.Fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(schema.Fields))
		}
		
		// Check field names
		expectedFields := []string{"id", "name"}
		for i, field := range schema.Fields {
			if field.Name != expectedFields[i] {
				t.Errorf("Expected field %s, got %s", expectedFields[i], field.Name)
			}
		}
	})
	
	// Test GetPartitions
	t.Run("GetPartitions", func(t *testing.T) {
		partitions, err := reader.GetPartitions()
		
		// Check for errors
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}
		
		// Check partitions
		expectedPartitions := []string{"year=2023", "year=2024"}
		if len(partitions) != len(expectedPartitions) {
			t.Errorf("Expected %d partitions, got %d", len(expectedPartitions), len(partitions))
		}
		
		for i, partition := range partitions {
			if partition != expectedPartitions[i] {
				t.Errorf("Expected partition %s, got %s", expectedPartitions[i], partition)
			}
		}
	})
}

// Mock implementation for testing
func createMockParsedData() *ParsedData {
	return &ParsedData{
		Records: []map[string]interface{}{
			{
				"id":   1,
				"name": "Test 1",
			},
			{
				"id":   2,
				"name": "Test 2",
			},
		},
		Schema: map[string]string{
			"id":   "integer",
			"name": "string",
		},
		Count: 2,
	}
}

// Test JSON serialization/deserialization
func TestParsedDataJSON(t *testing.T) {
	// Create test data
	data := createMockParsedData()
	
	// Serialize to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Errorf("Failed to marshal to JSON: %v", err)
		return
	}
	
	// Deserialize from JSON
	var parsedData ParsedData
	err = json.Unmarshal(jsonBytes, &parsedData)
	if err != nil {
		t.Errorf("Failed to unmarshal from JSON: %v", err)
		return
	}
	
	// Check record count
	if len(parsedData.Records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(parsedData.Records))
	}
	
	// Check schema
	if len(parsedData.Schema) != 2 {
		t.Errorf("Expected 2 schema fields, got %d", len(parsedData.Schema))
	}
	
	// Check count
	if parsedData.Count != 2 {
		t.Errorf("Expected count 2, got %d", parsedData.Count)
	}
}
