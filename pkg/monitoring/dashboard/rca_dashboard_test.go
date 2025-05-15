package dashboard

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/rca"
)

// MockRCAClient implements the RCAClient interface for testing
type MockRCAClient struct {
	results  []*rca.RCAResult
	insights *rca.Insights
}

// AnalyzeAnomaly implements RCAClient.AnalyzeAnomaly
func (m *MockRCAClient) AnalyzeAnomaly(anomalyID string, config *rca.Config) (*rca.RCAResult, error) {
	// Return a mock result
	result := &rca.RCAResult{
		AnomalyID:    anomalyID,
		AnalysisTime: time.Now(),
		PrimaryRootCause: &rca.RootCause{
			Type:        "schema_change",
			Description: "Schema change detected in table",
			Details:     "Column 'user_id' type changed from INT to STRING",
			Confidence:  0.85,
		},
		OtherCauses: []*rca.RootCause{
			{
				Type:        "data_quality",
				Description: "Data quality issue detected",
				Details:     "Null values increased by 15%",
				Confidence:  0.65,
			},
		},
		AffectedTables:     []string{"users", "orders"},
		RelatedAnomalies:   []string{"anom-20250513-002"},
		RecommendedActions: []string{"Verify schema change", "Check data pipeline"},
	}

	// Add to mock results
	m.results = append(m.results, result)
	return result, nil
}

// GetRecentAnalyses implements RCAClient.GetRecentAnalyses
func (m *MockRCAClient) GetRecentAnalyses(limit int) ([]*rca.RCAResult, error) {
	if len(m.results) == 0 {
		// Create some mock results if none exist
		m.results = []*rca.RCAResult{
			{
				AnomalyID:    "anom-20250514-001",
				AnalysisTime: time.Now().Add(-1 * time.Hour),
				PrimaryRootCause: &rca.RootCause{
					Type:        "schema_change",
					Description: "Schema change detected in table",
					Details:     "Column 'user_id' type changed from INT to STRING",
					Confidence:  0.85,
				},
				AffectedTables: []string{"users", "orders"},
			},
			{
				AnomalyID:    "anom-20250513-002",
				AnalysisTime: time.Now().Add(-24 * time.Hour),
				PrimaryRootCause: &rca.RootCause{
					Type:        "data_quality",
					Description: "Data quality issue detected",
					Details:     "Null values increased by 15%",
					Confidence:  0.65,
				},
				AffectedTables: []string{"products"},
			},
		}
	}

	// Apply limit
	results := m.results
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	return results, nil
}

// GetAnalysisResult implements RCAClient.GetAnalysisResult
func (m *MockRCAClient) GetAnalysisResult(anomalyID string) (*rca.RCAResult, error) {
	// Find result by anomaly ID
	for _, result := range m.results {
		if result.AnomalyID == anomalyID {
			return result, nil
		}
	}
	return nil, nil
}

// GetInsights implements RCAClient.GetInsights
func (m *MockRCAClient) GetInsights() (*rca.Insights, error) {
	if m.insights == nil {
		// Create mock insights
		m.insights = &rca.Insights{
			RootCauseDistribution: map[string]int{
				"schema_change": 3,
				"data_quality":  2,
				"system_issue":  1,
			},
			AffectedTablesCount: map[string]int{
				"users":    3,
				"orders":   2,
				"products": 1,
			},
			CommonRootCauses: []*rca.CommonRootCause{
				{
					Description:   "Schema change detected in table",
					Details:       "Column type changes",
					Count:         3,
					AvgConfidence: 0.85,
				},
				{
					Description:   "Data quality issue detected",
					Details:       "Null values increased",
					Count:         2,
					AvgConfidence: 0.65,
				},
			},
			StartTime:     time.Now().Add(-7 * 24 * time.Hour),
			EndTime:       time.Now(),
			TotalAnalyses: 6,
		}
	}
	return m.insights, nil
}

// TestRcaDashboardHandler tests the RCA dashboard handler
func TestRcaDashboardHandler(t *testing.T) {
	// Create a mock RCA client
	mockRcaClient := &MockRCAClient{}

	// Create a dashboard with the mock client
	dash, err := New(nil, DashboardOptions{
		RCAClient: mockRcaClient,
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Create a request to the RCA dashboard
	req, err := http.NewRequest("GET", "/rca", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request
	dash.handleRcaDashboard(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the response contains expected content
	if !strings.Contains(rr.Body.String(), "Root Cause Analysis Dashboard") {
		t.Errorf("Handler returned unexpected body: %v", rr.Body.String())
	}
}

// TestRcaApiHandler tests the RCA API handlers
func TestRcaApiHandler(t *testing.T) {
	// Create a mock RCA client
	mockRcaClient := &MockRCAClient{}

	// Create a dashboard with the mock client
	dash, err := New(nil, DashboardOptions{
		RCAClient: mockRcaClient,
	})
	if err != nil {
		t.Fatalf("Failed to create dashboard: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		path           string
		method         string
		body           string
		expectedStatus int
		checkResponse  func(t *testing.T, body string)
	}{
		{
			name:           "Get Recent Analyses",
			path:           "/api/rca/recent",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var results []*rca.RCAResult
				if err := json.Unmarshal([]byte(body), &results); err != nil {
					t.Errorf("Failed to parse response: %v", err)
					return
				}
				if len(results) == 0 {
					t.Errorf("Expected non-empty results")
				}
			},
		},
		{
			name:           "Get Insights",
			path:           "/api/rca/insights",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var insights rca.Insights
				if err := json.Unmarshal([]byte(body), &insights); err != nil {
					t.Errorf("Failed to parse response: %v", err)
					return
				}
				if insights.TotalAnalyses == 0 {
					t.Errorf("Expected non-zero total analyses")
				}
			},
		},
		{
			name:           "Analyze Anomaly",
			path:           "/api/rca/analyze",
			method:         "POST",
			body:           `{"anomaly_id":"anom-20250515-001","format":"json"}`,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var result rca.RCAResult
				if err := json.Unmarshal([]byte(body), &result); err != nil {
					t.Errorf("Failed to parse response: %v", err)
					return
				}
				if result.AnomalyID != "anom-20250515-001" {
					t.Errorf("Expected anomaly ID anom-20250515-001, got %s", result.AnomalyID)
				}
			},
		},
		{
			name:           "Get Analysis Result",
			path:           "/api/rca/anom-20250514-001",
			method:         "GET",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body string) {
				var result rca.RCAResult
				if err := json.Unmarshal([]byte(body), &result); err != nil {
					t.Errorf("Failed to parse response: %v", err)
					return
				}
				if result.AnomalyID != "anom-20250514-001" {
					t.Errorf("Expected anomaly ID anom-20250514-001, got %s", result.AnomalyID)
				}
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			var req *http.Request
			var err error
			if tc.method == "POST" {
				req, err = http.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			} else {
				req, err = http.NewRequest(tc.method, tc.path, nil)
			}
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set content type for POST requests
			if tc.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			dash.handleRcaAPI(rr, req)

			// Check status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
			}

			// Check response
			if tc.checkResponse != nil {
				tc.checkResponse(t, rr.Body.String())
			}
		})
	}
}
