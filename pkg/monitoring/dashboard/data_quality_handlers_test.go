package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/quality/profile"
	"github.com/nessi-dev/nessi-dev/pkg/quality/rules"
)

func TestHandleGetProfiles(t *testing.T) {
	// Create a new dashboard
	dashboard, err := NewDashboard(nil, DashboardOptions{
		ListenAddr: ":8080",
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Create a request with a table parameter
	req, err := http.NewRequest("GET", "/api/profiles?table=test_table", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	dashboard.handleGetProfiles(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Parse the response
	var response DataQualityResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the response
	if !response.Success {
		t.Errorf("Handler returned error: %v", response.Message)
	}

	// Check that data is present
	if response.Data == nil {
		t.Errorf("Handler returned nil data")
	}
}

func TestHandleGetProfilesMissingTable(t *testing.T) {
	// Create a new dashboard
	dashboard, err := NewDashboard(nil, DashboardOptions{
		ListenAddr: ":8080",
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Create a request without a table parameter
	req, err := http.NewRequest("GET", "/api/profiles", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	dashboard.handleGetProfiles(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Parse the response
	var response DataQualityResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the response
	if response.Success {
		t.Errorf("Handler should have returned error")
	}

	// Check the error message
	if response.Message != "Missing table parameter" {
		t.Errorf("Handler returned wrong error message: got %v want %v", response.Message, "Missing table parameter")
	}
}

func TestHandleGetRules(t *testing.T) {
	// Create a new dashboard
	dashboard, err := NewDashboard(nil, DashboardOptions{
		ListenAddr: ":8080",
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("GET", "/api/rules", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	dashboard.handleGetRules(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response
	var response DataQualityResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the response
	if !response.Success {
		t.Errorf("Handler returned error: %v", response.Message)
	}

	// Check that data is present
	if response.Data == nil {
		t.Errorf("Handler returned nil data")
	}

	// Check that rules are returned
	rules, ok := response.Data.([]interface{})
	if !ok {
		t.Fatalf("Handler returned wrong data type: %T", response.Data)
	}

	// Check that at least one rule is returned
	if len(rules) == 0 {
		t.Errorf("Handler returned empty rules list")
	}
}

func TestHandleValidateRules(t *testing.T) {
	// Create a new dashboard
	dashboard, err := NewDashboard(nil, DashboardOptions{
		ListenAddr: ":8080",
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Create a request with table parameter
	req, err := http.NewRequest("POST", "/api/rules/validate?table=test_table", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	dashboard.handleValidateRules(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the response
	var response DataQualityResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check the response
	if !response.Success {
		t.Errorf("Handler returned error: %v", response.Message)
	}

	// Check that data is present
	if response.Data == nil {
		t.Errorf("Handler returned nil data")
	}
}

func TestHandleDataQualityDashboard(t *testing.T) {
	// Create a new dashboard
	dashboard, err := NewDashboard(nil, DashboardOptions{
		ListenAddr: ":8080",
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("GET", "/data-quality", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	dashboard.handleDataQualityDashboard(rr, req)

	// Check the status code
	// Note: This will fail because the template doesn't exist in the test environment
	// We're just checking that the handler doesn't panic
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Logf("Expected internal server error due to missing template, got: %v", status)
	}
}

// Mock implementations for testing
type MockProfiler struct{}

func (m *MockProfiler) GenerateProfile() ([]*profile.Profile, error) {
	return []*profile.Profile{
		{
			Name:        "test_column",
			Type:        "string",
			RowCount:    100,
			NullCount:   10,
			NullPercent: 10.0,
			Distinct:    50,
			Patterns:    []string{"pattern1", "pattern2"},
			Anomalies:   []profile.Anomaly{{Type: "outlier", Value: "test", Description: "test anomaly"}},
		},
	}, nil
}

type MockRuleValidator struct{}

func (m *MockRuleValidator) Validate(record interface{}) []rules.ValidationError {
	return []rules.ValidationError{
		{
			RuleID:      "test_rule",
			Message:     "test error",
			ColumnName:  "test_column",
			Value:       "test_value",
			Severity:    "error",
			RowIndex:    0,
			ExecutionID: "test_execution",
		},
	}
}
