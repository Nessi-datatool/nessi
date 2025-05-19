package databricks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabricksClient_GetWorkspaces(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/2.0/workspaces", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"workspaces": [
				{"workspace_id": "123", "workspace_name": "Test Workspace"}
			]
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: &http.Client{},
	}

	// Test GetWorkspaces
	ctx := context.Background()
	workspaces, err := client.GetWorkspaces(ctx)
	require.NoError(t, err)
	require.Len(t, workspaces, 1)
	assert.Equal(t, "123", workspaces[0].ID)
	assert.Equal(t, "Test Workspace", workspaces[0].Name)
}

func TestDatabricksClient_GetCatalogs(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/2.1/unity-catalog/catalogs", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"catalogs": [
				{"name": "main", "comment": "Main catalog"}
			]
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: &http.Client{},
	}

	// Test GetCatalogs
	ctx := context.Background()
	catalogs, err := client.GetCatalogs(ctx, "123")
	require.NoError(t, err)
	require.Len(t, catalogs, 1)
	assert.Equal(t, "main", catalogs[0].Name)
}

func TestDatabricksClient_GetSchemas(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/2.1/unity-catalog/schemas", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.RawQuery, "catalog_name=main")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				{"name": "default", "catalog_name": "main"}
			]
		}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "test-token",
		httpClient: &http.Client{},
	}

	// Test GetSchemas
	ctx := context.Background()
	schemas, err := client.GetSchemas(ctx, "123", "main")
	require.NoError(t, err)
	require.Len(t, schemas, 1)
	assert.Equal(t, "default", schemas[0].Name)
	assert.Equal(t, "main", schemas[0].CatalogName)
}

func TestDatabricksClient_ErrorHandling(t *testing.T) {
	// Create a mock server that returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error_code": "INVALID_TOKEN", "message": "The provided token is invalid or has expired"}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &DatabricksClient{
		baseURL:    server.URL,
		token:      "invalid-token",
		httpClient: &http.Client{},
	}

	// Test error handling
	ctx := context.Background()
	_, err := client.GetWorkspaces(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "INVALID_TOKEN")
	assert.Contains(t, err.Error(), "The provided token is invalid or has expired")
}
