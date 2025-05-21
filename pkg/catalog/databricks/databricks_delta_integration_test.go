package databricks

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabricksDeltaIntegrationAlt tests the integration between Databricks catalog and Delta Lake (alternative implementation)
func TestDatabricksDeltaIntegrationAlt(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "databricks-delta-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(deltaTableDir, "_delta_log"), 0755))

	// Create test data with various data types
	testData := []map[string]interface{}{
		{
			"id":        int64(1),
			"name":      "John Doe",
			"age":       int32(30),
			"active":    true,
			"salary":    float64(75000.50),
			"hire_date": time.Now(),
		},
	}

	// Create a datalake schema for testing
	fields := []datalake.Field{
		{Name: "id", Type: datalake.FieldTypeInt64},
		{Name: "name", Type: datalake.FieldTypeString},
		{Name: "age", Type: datalake.FieldTypeInt32},
		{Name: "active", Type: datalake.FieldTypeBool},
		{Name: "salary", Type: datalake.FieldTypeFloat64},
		{Name: "hire_date", Type: datalake.FieldTypeTimestamp},
	}
	datalakeSchema := datalake.NewSchema(fields)

	// Initialize Delta format handler
	deltaHandler := datalake.NewDeltaFormatHandler()

	// Write test data to a Delta table
	err = deltaHandler.Write(deltaTableDir, testData, datalakeSchema)
	assert.NoError(t, err, "Failed to write test data to Delta table")

	// Create a mock client for testing table details
	mockClient := &mockDatabricksClientWithTableDetails{
		baseURL:     "https://test.databricks.com",
		token:       "test-token",
		httpClient:  nil,
		workspaceID: "test-workspace",
	}

	// Create catalog instance
	catalog := &DatabricksCatalog{
		client:      mockClient,
		workspaceID: "test-workspace",
	}

	// Test getting table details
	ctx := context.Background()
	tableDetails, err := mockClient.GetTableDetails(ctx, "test-workspace", "test-catalog", "test-schema", "test-table")
	assert.NoError(t, err, "Should get table details without error")
	assert.Equal(t, "test-table", tableDetails.Info.Name, "Table name should match")
	assert.Equal(t, "delta", tableDetails.Info.Type, "Table type should be delta")

	// Verify the schema fields
	assert.NotNil(t, tableDetails.Schema, "Schema should not be nil")
	assert.NotEmpty(t, tableDetails.Schema.Fields, "Schema should have fields")
	assert.Equal(t, "id", tableDetails.Schema.Fields[0].Name, "First field should be id")
	assert.Equal(t, "long", tableDetails.Schema.Fields[0].Type, "Id field should be of type long")

	// Test getting Delta table details directly from the path
	deltaDetails, err := catalog.GetDeltaTableDetails(ctx, deltaTableDir)
	assert.NoError(t, err, "Should get Delta table details without error")
	assert.NotNil(t, deltaDetails, "Delta table details should not be nil")

	// Test error handling with a non-existent Delta table
	nonExistentPath := filepath.Join(tempDir, "non-existent")
	_, err = catalog.GetDeltaTableDetails(ctx, nonExistentPath)
	assert.Error(t, err, "Should error when path is not a Delta table")
	assert.Contains(t, err.Error(), "not a Delta Lake table", "Error message should indicate it's not a Delta table")
}

// TestDatabricksErrorScenarios tests various error scenarios in the Databricks catalog
func TestDatabricksErrorScenarios(t *testing.T) {
	ctx := context.Background()

	// Test API errors
	testCases := []struct {
		name      string
		errorMode string
		testFunc  func(client DatabricksAPI) (interface{}, error)
	}{
		{
			name:      "GetWorkspaces error",
			errorMode: "workspaces",
			testFunc: func(client DatabricksAPI) (interface{}, error) {
				return client.GetWorkspaces(ctx)
			},
		},
		{
			name:      "GetCatalogs error",
			errorMode: "catalogs",
			testFunc: func(client DatabricksAPI) (interface{}, error) {
				return client.GetCatalogs(ctx, "test-workspace")
			},
		},
		{
			name:      "GetSchemas error",
			errorMode: "schemas",
			testFunc: func(client DatabricksAPI) (interface{}, error) {
				return client.GetSchemas(ctx, "test-workspace", "main")
			},
		},
		{
			name:      "GetTables error",
			errorMode: "tables",
			testFunc: func(client DatabricksAPI) (interface{}, error) {
				return client.GetTables(ctx, "test-workspace", "main", "default")
			},
		},
		{
			name:      "GetTableDetails error",
			errorMode: "tableDetails",
			testFunc: func(client DatabricksAPI) (interface{}, error) {
				return client.GetTableDetails(ctx, "test-workspace", "main", "default", "delta_table")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockDatabricksClientWithErrors{
				baseURL:     "http://localhost:8080",
				token:       "test-token",
				httpClient:  nil,
				workspaceID: "test-workspace",
				errorMode:   tc.errorMode,
			}

			_, err := tc.testFunc(client)
			assert.Error(t, err, "Should return an error when API call fails")
		})
	}

	// Test with empty workspace ID
	client := &mockDatabricksClientWithErrors{
		baseURL:     "http://localhost:8080",
		token:       "test-token",
		httpClient:  &http.Client{},
		workspaceID: "test-workspace",
		errorMode:   "catalogs",
	}

	emptyCatalog := &DatabricksCatalog{
		client:      client,
		workspaceID: "", // Empty workspace ID
	}

	// When workspace ID is empty, we should get an error
	_, err := emptyCatalog.ListCatalogs(ctx)
	assert.Error(t, err, "Should error when workspace ID is empty")
}
