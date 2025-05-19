package databricks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabricksClient_ErrorHandling(t *testing.T) {
	// Test cases for different error scenarios
	testCases := []struct {
		name           string
		statusCode     int
		responseBody   string
		endpoint       string
		expectedErrMsg string
	}{
		{
			name:           "unauthorized access",
			statusCode:     http.StatusUnauthorized,
			responseBody:   `{"error_code": "UNAUTHORIZED", "message": "Invalid access token"}`,
			endpoint:       catalogsEndpoint,
			expectedErrMsg: "API error: {\"error_code\": \"UNAUTHORIZED\", \"message\": \"Invalid access token\"} (status code: 401)",
		},
		{
			name:           "resource not found",
			statusCode:     http.StatusNotFound,
			responseBody:   `{"error_code": "RESOURCE_DOES_NOT_EXIST", "message": "Catalog 'nonexistent' not found"}`,
			endpoint:       schemasEndpoint,
			expectedErrMsg: "API error: {\"error_code\": \"RESOURCE_DOES_NOT_EXIST\", \"message\": \"Catalog 'nonexistent' not found\"} (status code: 404)",
		},
		{
			name:           "rate limit exceeded",
			statusCode:     http.StatusTooManyRequests,
			responseBody:   `{"error_code": "RATE_LIMIT_EXCEEDED", "message": "Too many requests"}`,
			endpoint:       tablesEndpoint,
			expectedErrMsg: "API error: {\"error_code\": \"RATE_LIMIT_EXCEEDED\", \"message\": \"Too many requests\"} (status code: 429)",
		},
		{
			name:           "internal server error",
			statusCode:     http.StatusInternalServerError,
			responseBody:   `{"error_code": "INTERNAL_ERROR", "message": "An internal error occurred"}`,
			endpoint:       catalogsEndpoint,
			expectedErrMsg: "API error: {\"error_code\": \"INTERNAL_ERROR\", \"message\": \"An internal error occurred\"} (status code: 500)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server that returns the specified error
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.responseBody))
			}))
			defer server.Close()

			// Create client that points to test server
			client := NewDatabricksClient(server.URL, "test-token", "test-workspace")

			// Test error handling based on the endpoint
			var err error
			ctx := context.Background()

			switch tc.endpoint {
			case catalogsEndpoint:
				_, err = client.GetCatalogs(ctx, "test-workspace")
			case schemasEndpoint:
				_, err = client.GetSchemas(ctx, "test-workspace", "main")
			case tablesEndpoint:
				_, err = client.GetTables(ctx, "test-workspace", "main", "default")
			}

			// Verify error message
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErrMsg)
		})
	}
}

func TestDatabricksClient_NetworkErrors(t *testing.T) {
	// Create a client with an invalid URL to simulate network errors
	client := NewDatabricksClient("http://invalid-url-that-does-not-exist.example", "test-token", "test-workspace")

	// Test that network errors are properly handled
	_, err := client.GetCatalogs(context.Background(), "test-workspace")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute request")
}

func TestDatabricksClient_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"catalogs": [{"name": "main", "invalid JSON`))
	}))
	defer server.Close()

	// Create client that points to test server
	client := NewDatabricksClient(server.URL, "test-token", "test-workspace")

	// Test that JSON parsing errors are properly handled
	_, err := client.GetCatalogs(context.Background(), "test-workspace")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal")
}
