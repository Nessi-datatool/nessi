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

// TestDatabricksErrorHandling tests error handling in the Databricks client
func TestDatabricksErrorHandling(t *testing.T) {
	ctx := context.Background()

	// Test cases for different error scenarios
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
				httpClient:  &http.Client{},
				workspaceID: "test-workspace",
				errorMode:   tc.errorMode,
			}

			_, err := tc.testFunc(client)
			assert.Error(t, err, "Should return an error when API call fails")
		})
	}

	// Test invalid Delta table path
	t.Run("Invalid Delta table path", func(t *testing.T) {
		client := &mockDatabricksClientWithErrors{
			baseURL:     "http://localhost:8080",
			token:       "test-token",
			httpClient:  &http.Client{},
			workspaceID: "test-workspace",
			errorMode:   "invalidPath",
		}

		catalog := &DatabricksCatalog{
			client:      client,
			workspaceID: "test-workspace",
		}

		// Try to get details for a table with an invalid path
		_, err := catalog.GetTableDetails(ctx, "default", "delta_table")
		assert.NoError(t, err, "Should not error on getting table details")

		// Try to get Delta table details for an invalid path
		_, err = catalog.GetDeltaTableDetails(ctx, "/nonexistent/path")
		assert.Error(t, err, "Should error when trying to get Delta table details for an invalid path")
	})
}

// TestDatabricksSchemaInference tests schema inference from Delta tables
func TestDatabricksSchemaInference(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "databricks-schema-test")
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

	// Create a client for error testing
	errorClient := &mockDatabricksClientWithErrors{
		baseURL:     "http://localhost:8080",
		token:       "test-token",
		httpClient:  &http.Client{},
		workspaceID: "test-workspace",
		errorMode:   "tables",
	}

	// Test error handling
	_, err = errorClient.GetTables(context.Background(), "test-workspace", "test-catalog", "test-schema")
	assert.Error(t, err, "Should return an error when in error mode")
	assert.Contains(t, err.Error(), "failed to get tables", "Error message should indicate the failure reason")

	// Create a mock client for testing table details
	mockClient := &mockDatabricksClientWithTableDetails{
		baseURL:     "https://test.databricks.com",
		token:       "test-token",
		httpClient:  &http.Client{},
		workspaceID: "test-workspace",
	}

	// Test getting table details
	tableDetails, err := mockClient.GetTableDetails(context.Background(), "test-workspace", "test-catalog", "test-schema", "test-table")
	assert.NoError(t, err, "Should get table details without error")
	assert.Equal(t, "test-table", tableDetails.Info.Name, "Table name should match")
	assert.Equal(t, "delta", tableDetails.Info.Type, "Table type should be delta")

	// Verify the schema fields
	assert.NotNil(t, tableDetails.Schema, "Schema should not be nil")
	assert.NotEmpty(t, tableDetails.Schema.Fields, "Schema should have fields")
	assert.Equal(t, "id", tableDetails.Schema.Fields[0].Name, "First field should be id")
	assert.Equal(t, "long", tableDetails.Schema.Fields[0].Type, "Id field should be of type long")

	// Verify the table metadata
	assert.NotNil(t, tableDetails.Info, "Table info should not be nil")
	assert.Equal(t, "test-table", tableDetails.Info.Name, "Table name should match")

	// Create a mock client that returns errors for nonexistent tables
	var schemaErrorClient *mockDatabricksClientWithErrors = &mockDatabricksClientWithErrors{
		baseURL:     "http://localhost:8080",
		token:       "test-token",
		httpClient:  &http.Client{},
		workspaceID: "test-workspace",
		errorMode:   "tableDetails",
	}

	// Test error handling for nonexistent table
	testCtx := context.Background()
	var testErr error
	_, testErr = schemaErrorClient.GetTableDetails(testCtx, "test-workspace", "hive_metastore", "default", "nonexistent_table")
	assert.Error(t, testErr, "Should return error for nonexistent table")

	// Create a custom mock client that returns the correct path
	customMockClient := &mockDatabricksClientWithTableDetails{
		baseURL:        "https://test.databricks.com",
		token:          "test-token",
		httpClient:     &http.Client{},
		workspaceID:    "test-workspace",
		deltaTablePath: deltaTableDir, // Use the actual temp directory path
	}

	// Create catalog instance
	catalog := &DatabricksCatalog{
		client:      customMockClient,
		workspaceID: "test-workspace",
	}

	tableDetails, err = catalog.GetTableDetails(testCtx, "default", "delta_table")
	require.NoError(t, err)
	assert.Equal(t, "delta_table", tableDetails.Info.Name)
	assert.Equal(t, deltaTableDir, tableDetails.Info.Location)

	// Test getting Delta table details directly from the path
	deltaDetails, err := catalog.GetDeltaTableDetails(testCtx, deltaTableDir)
	require.NoError(t, err)
	assert.NotNil(t, deltaDetails)
	assert.Equal(t, "delta-table", deltaDetails.Info.Name) // Should be the directory name
	assert.Equal(t, "delta", deltaDetails.Schema.Format)

	// Verify the schema fields
	require.NotNil(t, deltaDetails.Schema)
	require.NotEmpty(t, deltaDetails.Schema.Fields)

	// Check that all fields from our original schema are present
	fieldNames := make(map[string]bool)
	for _, field := range deltaDetails.Schema.Fields {
		fieldNames[field.Name] = true
	}

	// Verify all expected fields are present
	expectedFields := []string{"id", "name", "age", "active", "salary", "hire_date"}
	for _, fieldName := range expectedFields {
		assert.True(t, fieldNames[fieldName], "Schema should include field %s", fieldName)
	}
}

// TestDatabricksEdgeCases tests edge cases in the Databricks catalog
func TestDatabricksEdgeCases(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "databricks-edge-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a Delta table directory but don't initialize it as a proper Delta table
	emptyDeltaDir := filepath.Join(tempDir, "empty-delta")
	require.NoError(t, os.MkdirAll(emptyDeltaDir, 0755))

	// Create a client
	client := &mockDatabricksClientWithErrors{
		baseURL:     "http://localhost:8080",
		token:       "test-token",
		httpClient:  &http.Client{},
		workspaceID: "test-workspace",
	}

	// Create catalog instance
	catalog := &DatabricksCatalog{
		client:      client,
		workspaceID: "test-workspace",
	}

	ctx := context.Background()

	// Test with empty directory (not a Delta table)
	_, err = catalog.GetDeltaTableDetails(ctx, emptyDeltaDir)
	assert.Error(t, err, "Should error when directory is not a Delta table")

	// Test with non-existent table
	_, err = catalog.GetTableDetails(ctx, "nonexistent", "nonexistent")
	assert.NoError(t, err, "Should not error for non-existent table (mock returns data)")

	// Test with empty workspace ID
	emptyCatalog := &DatabricksCatalog{
		client:      client,
		workspaceID: "", // Empty workspace ID
	}

	_, err = emptyCatalog.ListCatalogs(ctx)
	assert.Error(t, err, "Should error when workspace ID is empty")
}
