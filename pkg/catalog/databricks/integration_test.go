package databricks

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testDatabricksClient is a test implementation of the DatabricksAPI interface defined in catalog.go

// testDatabricksClient is a test implementation of DatabricksAPI
type testDatabricksClient struct {
	baseURL       string
	token         string
	httpClient    *http.Client
	workspaceID   string
	deltaTableDir string
}

func (c *testDatabricksClient) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	return []Workspace{{ID: "test-workspace", Name: "Test Workspace"}}, nil
}

func (c *testDatabricksClient) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	return []Catalog{{Name: "main"}}, nil
}

func (c *testDatabricksClient) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	return []Schema{{Name: "default", CatalogName: "main"}}, nil
}

func (c *testDatabricksClient) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	return []Table{{Name: "delta_table", CatalogName: "main", SchemaName: "default"}}, nil
}

func (c *testDatabricksClient) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	return &types.TableDetails{
		Info: types.TableInfo{
			Name:        "delta_table",
			Type:        "delta",
			Location:    c.deltaTableDir,
			Description: "Delta Lake table",
		},
		Schema: &types.TableSchema{
			Format:  "delta",
			Version: 1,
			Fields: []types.FieldInfo{
				{Name: "id", Type: "long", Nullable: false},
				{Name: "name", Type: "string", Nullable: true},
				{Name: "age", Type: "integer", Nullable: true},
			},
		},
		Metadata: &types.TableMetadata{
			Owner:     "test-user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Tags:      []string{"delta", "test"},
		},
	}, nil
}

func (c *testDatabricksClient) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	return nil
}

// testDatabricksCatalog is a test implementation of a Databricks catalog
type testDatabricksCatalog struct {
	client      DatabricksAPI
	workspaceID string
	catalog     string
	schema      string
}

func (c *testDatabricksCatalog) ListCatalogs(ctx context.Context) ([]string, error) {
	return []string{"main"}, nil
}

func (c *testDatabricksCatalog) ListSchemas(ctx context.Context, catalogName string) ([]string, error) {
	return []string{"default"}, nil
}

func (c *testDatabricksCatalog) ListTables(ctx context.Context, catalogName, schemaName string) ([]string, error) {
	return []string{"delta_table"}, nil
}

func (c *testDatabricksCatalog) GetTableDetails(ctx context.Context, schemaName, tableName string) (*types.TableDetails, error) {
	return c.client.GetTableDetails(ctx, c.workspaceID, c.catalog, schemaName, tableName)
}

func (c *testDatabricksCatalog) GetDeltaTableDetails(ctx context.Context, path string) (*types.TableDetails, error) {
	return datalake.GetDeltaTableMetadata(ctx, path)
}

// TestDatabricksIntegration tests the integration between Databricks catalog and Delta Lake format handler
func TestDatabricksIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "databricks-integration-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(deltaTableDir, "_delta_log"), 0755))

	// Create Databricks client
	client := &testDatabricksClient{
		baseURL:       "http://localhost:8080",
		token:         "test-token",
		httpClient:    &http.Client{},
		workspaceID:   "test-workspace",
		deltaTableDir: deltaTableDir,
	}

	// Create catalog instance
	catalog := &testDatabricksCatalog{
		client:      client,
		workspaceID: "test-workspace",
		catalog:     "main",
		schema:      "default",
	}

	// Test getting catalogs
	ctx := context.Background()
	catalogs, err := catalog.ListCatalogs(ctx)
	require.NoError(t, err)
	require.Len(t, catalogs, 1)
	assert.Equal(t, "main", catalogs[0])

	// Test getting schemas
	schemas, err := catalog.ListSchemas(ctx, "main")
	require.NoError(t, err)
	require.Len(t, schemas, 1)
	assert.Equal(t, "default", schemas[0])

	// Test getting tables
	tables, err := catalog.ListTables(ctx, "main", "default")
	require.NoError(t, err)
	require.Len(t, tables, 1)
	assert.Equal(t, "delta_table", tables[0])

	// Test getting table details
	tableDetails, err := catalog.GetTableDetails(ctx, "default", "delta_table")
	require.NoError(t, err)
	assert.Equal(t, "delta_table", tableDetails.Info.Name)
	assert.Equal(t, "delta", tableDetails.Schema.Format)
	assert.Equal(t, deltaTableDir, tableDetails.Info.Location)
	assert.Len(t, tableDetails.Schema.Fields, 3)

	// Now create some test data and write it to the Delta table
	testData := []map[string]interface{}{
		{"id": int64(1), "name": "John Doe", "age": int32(30)},
		{"id": int64(2), "name": "Jane Smith", "age": int32(25)},
	}

	// Create schema
	schema := datalake.NewSchema([]datalake.Field{
		{Name: "id", Type: datalake.FieldTypeInt64},
		{Name: "name", Type: datalake.FieldTypeString},
		{Name: "age", Type: datalake.FieldTypeInt32},
	})

	// Initialize Delta format handler
	deltaHandler := datalake.NewDeltaFormatHandler()

	// Write data to the Delta table
	err = deltaHandler.Write(deltaTableDir, testData, schema)
	require.NoError(t, err, "Should write data to Delta table without error")

	// Verify the table location is a Delta table
	assert.True(t, deltaHandler.IsDeltaTable(tableDetails.Info.Location), "Table location should be a Delta table")

	// Read data from the Delta table using the location from table details
	readData, err := deltaHandler.Read(tableDetails.Info.Location)
	require.NoError(t, err, "Should read data from Delta table without error")

	// Verify the data was read correctly
	require.Len(t, readData, len(testData), "Should read the same number of records")

	// Verify specific fields in the data
	for i, record := range readData {
		assert.Equal(t, testData[i]["id"], record["id"], "ID should match")
		assert.Equal(t, testData[i]["name"], record["name"], "Name should match")
		assert.Equal(t, testData[i]["age"], record["age"], "Age should match")
	}

	// Test getting table metadata
	metadata, err := datalake.GetDeltaTableMetadata(ctx, tableDetails.Info.Location)
	require.NoError(t, err)
	assert.Equal(t, filepath.Base(deltaTableDir), metadata.Info.Name)
	assert.Equal(t, "delta", metadata.Schema.Format)
}
