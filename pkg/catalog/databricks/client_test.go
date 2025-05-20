package databricks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabricksClient_GetCatalogs_Basic tests the basic functionality of the GetCatalogs method
func TestDatabricksClient_GetCatalogs_Basic(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/2.1/unity-catalog/catalogs", r.URL.Path)

		// Check authorization header
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"catalogs": []map[string]interface{}{
				{
					"name":        "main",
					"description": "Main catalog",
					"owner":       "admin",
				},
				{
					"name":        "test",
					"description": "Test catalog",
					"owner":       "test-user",
				},
			},
		})
	}))
	defer server.Close()

	// Create a client
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Call the method
	catalogs, err := client.GetCatalogs(context.Background(), "default")

	// Check results
	require.NoError(t, err)
	require.Len(t, catalogs, 2)
	assert.Equal(t, "main", catalogs[0].Name)
	assert.Equal(t, "Main catalog", catalogs[0].Description)
	assert.Equal(t, "admin", catalogs[0].Owner)
	assert.Equal(t, "test", catalogs[1].Name)
	assert.Equal(t, "Test catalog", catalogs[1].Description)
	assert.Equal(t, "test-user", catalogs[1].Owner)
}

// TestDatabricksClient_GetSchemas_Basic tests the basic functionality of the GetSchemas method
func TestDatabricksClient_GetSchemas_Basic(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/2.1/unity-catalog/schemas", r.URL.Path)

		// Check query parameters
		assert.Equal(t, "main", r.URL.Query().Get("catalog_name"))

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"schemas": []map[string]interface{}{
				{
					"name":        "default",
					"description": "Default schema",
					"owner":       "admin",
					"created_at":  "2023-01-01T00:00:00Z",
				},
				{
					"name":        "test",
					"description": "Test schema",
					"owner":       "test-user",
					"created_at":  "2023-01-02T00:00:00Z",
				},
			},
		})
	}))
	defer server.Close()

	// Create a client
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Call the method
	schemas, err := client.GetSchemas(context.Background(), "default", "main")

	// Check results
	require.NoError(t, err)
	require.Len(t, schemas, 2)
	assert.Equal(t, "default", schemas[0].Name)
	assert.Equal(t, "Default schema", schemas[0].Description)
	assert.Equal(t, "admin", schemas[0].Owner)
	assert.Equal(t, "2023-01-01T00:00:00Z", schemas[0].CreatedAt)
	assert.Equal(t, "test", schemas[1].Name)
	assert.Equal(t, "Test schema", schemas[1].Description)
	assert.Equal(t, "test-user", schemas[1].Owner)
	assert.Equal(t, "2023-01-02T00:00:00Z", schemas[1].CreatedAt)
}

// TestDatabricksClient_GetTables_Basic tests the basic functionality of the GetTables method
func TestDatabricksClient_GetTables_Basic(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/2.1/unity-catalog/tables", r.URL.Path)

		// Check query parameters
		assert.Equal(t, "main", r.URL.Query().Get("catalog_name"))
		assert.Equal(t, "default", r.URL.Query().Get("schema_name"))

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tables": []map[string]interface{}{
				{
					"name":        "customers",
					"description": "Customer table",
					"owner":       "admin",
					"created_at":  "2023-01-01T00:00:00Z",
					"format":      "delta",
				},
				{
					"name":        "orders",
					"description": "Order table",
					"owner":       "test-user",
					"created_at":  "2023-01-02T00:00:00Z",
					"format":      "delta",
				},
			},
		})
	}))
	defer server.Close()

	// Create a client
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Call the method
	tables, err := client.GetTables(context.Background(), "default", "main", "default")

	// Check results
	require.NoError(t, err)
	require.Len(t, tables, 2)
	assert.Equal(t, "customers", tables[0].Name)
	assert.Equal(t, "Customer table", tables[0].Description)
	assert.Equal(t, "admin", tables[0].Owner)
	assert.Equal(t, "2023-01-01T00:00:00Z", tables[0].CreatedAt)
	assert.Equal(t, "delta", tables[0].Format)
	assert.Equal(t, "orders", tables[1].Name)
	assert.Equal(t, "Order table", tables[1].Description)
	assert.Equal(t, "test-user", tables[1].Owner)
	assert.Equal(t, "2023-01-02T00:00:00Z", tables[1].CreatedAt)
	assert.Equal(t, "delta", tables[1].Format)
}

// TestDatabricksClient_GetTableDetails tests the GetTableDetails method
func TestDatabricksClient_GetTableDetails(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/2.1/unity-catalog/tables/main.default.customers", r.URL.Path)

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":         "customers",
			"catalog_name": "main",
			"schema_name":  "default",
			"table_type":   "MANAGED",
			"columns": []map[string]interface{}{
				{
					"name":        "id",
					"type_text":   "bigint",
					"type_name":   "bigint",
					"description": "Customer ID",
					"nullable":    false,
				},
				{
					"name":        "name",
					"type_text":   "string",
					"type_name":   "string",
					"description": "Customer name",
					"nullable":    true,
				},
			},
			"properties":       map[string]string{"delta.lastUpdateVersion": "1"},
			"description":      "Customer table",
			"owner":            "admin",
			"format":           "delta",
			"storage_location": "dbfs:/user/hive/warehouse/main.db/customers",
		})
	}))
	defer server.Close()

	// Create a client
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Call the method
	details, err := client.GetTableDetails(context.Background(), "default", "main", "default", "customers")

	// Check results
	require.NoError(t, err)
	assert.Equal(t, "customers", details.Info.Name)
	assert.Equal(t, "delta", details.Info.Type)
	assert.Equal(t, "Customer table", details.Info.Description)
	assert.Equal(t, "dbfs:/user/hive/warehouse/main.db/customers", details.Info.Location)
	assert.Equal(t, "1", details.Info.Properties["delta.lastUpdateVersion"])
	assert.Equal(t, "admin", details.Metadata.Owner)
	assert.Len(t, details.Schema.Fields, 2)
	assert.Equal(t, "id", details.Schema.Fields[0].Name)
	assert.Equal(t, "bigint", details.Schema.Fields[0].Type)
	assert.Equal(t, "Customer ID", details.Schema.Fields[0].Description)
	assert.False(t, details.Schema.Fields[0].Nullable)
	assert.Equal(t, "name", details.Schema.Fields[1].Name)
	assert.Equal(t, "string", details.Schema.Fields[1].Type)
	assert.Equal(t, "Customer name", details.Schema.Fields[1].Description)
	assert.True(t, details.Schema.Fields[1].Nullable)
}

// TestDatabricksClient_UpdateTableDetails tests the UpdateTableDetails method
func TestDatabricksClient_UpdateTableDetails(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/2.1/unity-catalog/tables/main.default.customers", r.URL.Path)

		// Check request body
		var requestBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		require.NoError(t, err)
		assert.Equal(t, "Updated customer table", requestBody["description"])
		assert.Equal(t, map[string]interface{}{"delta.lastUpdateVersion": "2"}, requestBody["properties"])

		// Return a success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Create a client
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Create table details
	details := &types.TableDetails{
		Info: types.TableInfo{
			Name:        "customers",
			Type:        "delta",
			Description: "Updated customer table",
			Location:    "dbfs:/user/hive/warehouse/main.db/customers",
			Properties:  map[string]string{"delta.lastUpdateVersion": "2"},
		},
		Schema: &types.TableSchema{
			Format:  "delta",
			Version: 1,
			Fields: []types.FieldInfo{
				{
					Name:        "id",
					Type:        "bigint",
					Description: "Customer ID",
					Nullable:    false,
				},
				{
					Name:        "name",
					Type:        "string",
					Description: "Customer name",
					Nullable:    true,
				},
			},
		},
		Metadata: &types.TableMetadata{
			Owner:      "admin",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Properties: map[string]string{"delta.lastUpdateVersion": "2"},
		},
	}

	// Call the method
	err := client.UpdateTableDetails(context.Background(), "default", "main", "default", "customers", details)

	// Check results
	require.NoError(t, err)
}

// TestDatabricksClient_ErrorHandling_Basic tests basic error handling in the client
func TestDatabricksClient_ErrorHandling_Basic(t *testing.T) {
	// Test cases for different error scenarios
	testCases := []struct {
		name       string
		statusCode int
		response   string
		expectErr  string
	}{
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"error_code":"INVALID_TOKEN","message":"Invalid token"}`,
			expectErr:  "failed to get catalogs: API error: {\"error_code\":\"INVALID_TOKEN\",\"message\":\"Invalid token\"} (status code: 401)",
		},
		{
			name:       "Not Found",
			statusCode: http.StatusNotFound,
			response:   `{"error_code":"RESOURCE_DOES_NOT_EXIST","message":"Resource not found"}`,
			expectErr:  "failed to get catalogs: API error: {\"error_code\":\"RESOURCE_DOES_NOT_EXIST\",\"message\":\"Resource not found\"} (status code: 404)",
		},
		{
			name:       "Rate Limit",
			statusCode: http.StatusTooManyRequests,
			response:   `{"error_code":"RATE_LIMIT_EXCEEDED","message":"Rate limit exceeded"}`,
			expectErr:  "failed to get catalogs: API error: {\"error_code\":\"RATE_LIMIT_EXCEEDED\",\"message\":\"Rate limit exceeded\"} (status code: 429)",
		},
		{
			name:       "Server Error",
			statusCode: http.StatusInternalServerError,
			response:   `{"error_code":"INTERNAL_ERROR","message":"Internal server error"}`,
			expectErr:  "failed to get catalogs: API error: {\"error_code\":\"INTERNAL_ERROR\",\"message\":\"Internal server error\"} (status code: 500)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			// Create a client
			client := &DatabricksClient{
				baseURL:    server.URL,
				token:      "test-token",
				httpClient: http.DefaultClient,
			}

			// Call a method that will fail
			_, err := client.GetCatalogs(context.Background(), "default")

			// Check error
			require.Error(t, err)
			assert.Equal(t, tc.expectErr, err.Error())
		})
	}
}

// TestDatabricksClient_NetworkError tests handling of network errors
func TestDatabricksClient_NetworkError(t *testing.T) {
	// Create a client with an invalid URL
	client := &DatabricksClient{
		baseURL:    "http://invalid-url-that-does-not-exist.example.com",
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Call a method that will fail due to network error
	_, err := client.GetCatalogs(context.Background(), "default")

	// Check error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute request")
}

// TestDatabricksClient_InvalidJSON_Basic tests basic handling of invalid JSON responses
func TestDatabricksClient_InvalidJSON_Basic(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"catalogs": ["invalid"`))
	}))
	defer server.Close()

	// Create a client
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: http.DefaultClient,
	}

	// Call a method that will fail due to invalid JSON
	_, err := client.GetCatalogs(context.Background(), "default")

	// Check error
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal catalogs response")
}
