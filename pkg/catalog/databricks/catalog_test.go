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
	// Create a new client with test configuration
	client := NewDatabricksClient("https://test.databricks.com", "test-token", "test-workspace")

	// Test GetWorkspaces
	workspaces, err := client.GetWorkspaces(context.Background())
	require.NoError(t, err)
	require.Len(t, workspaces, 1)
	assert.Equal(t, "test-workspace", workspaces[0].ID)
	assert.Equal(t, "Default Workspace", workspaces[0].Name)
	assert.Equal(t, "https://test.databricks.com", workspaces[0].URL)
}

func TestDatabricksClient_GetCatalogs(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/api/2.1/unity-catalog/catalogs", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"catalogs": [
				{
					"name": "main",
					"description": "Main catalog",
					"owner": "test-user",
					"type": "MANAGED",
					"created_at": "2023-01-01T00:00:00Z"
				},
				{
					"name": "samples",
					"description": "Sample data",
					"owner": "test-user",
					"type": "MANAGED",
					"created_at": "2023-01-01T00:00:00Z"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client that points to test server
	client := NewDatabricksClient(server.URL, "test-token", "test-workspace")

	// Test GetCatalogs
	catalogs, err := client.GetCatalogs(context.Background(), "test-workspace")
	require.NoError(t, err)
	require.Len(t, catalogs, 2)

	assert.Equal(t, "main", catalogs[0].Name)
	assert.Equal(t, "Main catalog", catalogs[0].Description)
	assert.Equal(t, "test-user", catalogs[0].Owner)

	assert.Equal(t, "samples", catalogs[1].Name)
	assert.Equal(t, "Sample data", catalogs[1].Description)
}

func TestDatabricksClient_GetSchemas(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/api/2.1/unity-catalog/schemas", r.URL.Path)
		assert.Equal(t, "main", r.URL.Query().Get("catalog_name"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"schemas": [
				{
					"name": "default",
					"catalog_name": "main",
					"description": "Default schema",
					"owner": "test-user",
					"created_at": "2023-01-01T00:00:00Z"
				},
				{
					"name": "sales",
					"catalog_name": "main",
					"description": "Sales data",
					"owner": "test-user",
					"created_at": "2023-01-01T00:00:00Z"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client that points to test server
	client := NewDatabricksClient(server.URL, "test-token", "test-workspace")

	// Test GetSchemas
	schemas, err := client.GetSchemas(context.Background(), "test-workspace", "main")
	require.NoError(t, err)
	require.Len(t, schemas, 2)

	assert.Equal(t, "default", schemas[0].Name)
	assert.Equal(t, "main", schemas[0].CatalogName)
	assert.Equal(t, "Default schema", schemas[0].Description)

	assert.Equal(t, "sales", schemas[1].Name)
	assert.Equal(t, "Sales data", schemas[1].Description)
}

func TestDatabricksClient_GetTables(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/api/2.1/unity-catalog/tables", r.URL.Path)
		assert.Equal(t, "main", r.URL.Query().Get("catalog_name"))
		assert.Equal(t, "default", r.URL.Query().Get("schema_name"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"tables": [
				{
					"name": "customers",
					"catalog_name": "main",
					"schema_name": "default",
					"description": "Customer data",
					"owner": "test-user",
					"format": "DELTA",
					"created_at": "2023-01-01T00:00:00Z"
				},
				{
					"name": "orders",
					"catalog_name": "main",
					"schema_name": "default",
					"description": "Order data",
					"owner": "test-user",
					"format": "DELTA",
					"created_at": "2023-01-01T00:00:00Z"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create client that points to test server
	client := NewDatabricksClient(server.URL, "test-token", "test-workspace")

	// Test GetTables
	tables, err := client.GetTables(context.Background(), "test-workspace", "main", "default")
	require.NoError(t, err)
	require.Len(t, tables, 2)

	assert.Equal(t, "customers", tables[0].Name)
	assert.Equal(t, "main", tables[0].CatalogName)
	assert.Equal(t, "default", tables[0].SchemaName)
	assert.Equal(t, "Customer data", tables[0].Description)
	assert.Equal(t, "DELTA", tables[0].Format)

	assert.Equal(t, "orders", tables[1].Name)
	assert.Equal(t, "Order data", tables[1].Description)
}
