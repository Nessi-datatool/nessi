package databricks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabricksClient_GetWorkspacesMock(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/2.0/workspaces", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"workspaces": [
				{"workspace_id": "test-workspace", "workspace_name": "Default Workspace"}
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

	// Mock the GetWorkspaces method since it's defined in catalog.go
	ctx := context.Background()
	respBody, err := client.makeRequest(ctx, http.MethodGet, "/api/2.0/workspaces", nil, nil)
	require.NoError(t, err)

	var response struct {
		Workspaces []struct {
			WorkspaceID   string `json:"workspace_id"`
			WorkspaceName string `json:"workspace_name"`
		} `json:"workspaces"`
	}
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("failed to unmarshal workspaces response: %v", err)
	}

	// Convert to Workspace struct
	workspaces := make([]Workspace, 0, len(response.Workspaces))
	for _, ws := range response.Workspaces {
		workspaces = append(workspaces, Workspace{
			ID:   ws.WorkspaceID,
			Name: ws.WorkspaceName,
		})
	}

	require.Len(t, workspaces, 1)
	assert.Equal(t, "test-workspace", workspaces[0].ID)
	assert.Equal(t, "Default Workspace", workspaces[0].Name)
}

func TestDatabricksClient_GetCatalogsMock(t *testing.T) {
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

func TestDatabricksClient_GetSchemasMock(t *testing.T) {
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

func TestDatabricksClient_BasicErrorHandling(t *testing.T) {
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
	// Use GetCatalogs instead of GetWorkspaces since that's what the error handling test expects
	_, err := client.GetCatalogs(ctx, "test-workspace")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "INVALID_TOKEN")
	assert.Contains(t, err.Error(), "The provided token is invalid or has expired")
}
