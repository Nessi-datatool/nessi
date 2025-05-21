package databricks

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

// errorDatabricksClient is a mock implementation that returns errors
type errorDatabricksClient struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	workspaceID string
	errorMode   string // Determines which error to return
}

func (c *errorDatabricksClient) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	if c.errorMode == "workspaces" {
		return nil, errors.New("failed to get workspaces: API error")
	}
	return []Workspace{{ID: "test-workspace", Name: "Test Workspace"}}, nil
}

func (c *errorDatabricksClient) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	if c.errorMode == "catalogs" {
		return nil, errors.New("failed to get catalogs: API error")
	}
	return []Catalog{{Name: "main"}}, nil
}

func (c *errorDatabricksClient) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	if c.errorMode == "schemas" {
		return nil, errors.New("failed to get schemas: API error")
	}
	return []Schema{{Name: "default", CatalogName: "main"}}, nil
}

func (c *errorDatabricksClient) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	if c.errorMode == "tables" {
		return nil, errors.New("failed to get tables: API error")
	}
	return []Table{{Name: "delta_table", CatalogName: "main", SchemaName: "default"}}, nil
}

func (c *errorDatabricksClient) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	if c.errorMode == "tableDetails" {
		return nil, errors.New("failed to get table details: API error")
	}

	if c.errorMode == "invalidPath" {
		return &types.TableDetails{
			Info: types.TableInfo{
				Name:     "delta_table",
				Type:     "delta",
				Location: "/nonexistent/path",
			},
		}, nil
	}

	return &types.TableDetails{
		Info: types.TableInfo{
			Name:        "delta_table",
			Type:        "delta",
			Location:    "/tmp/delta-table", // Will be replaced in tests
			Description: "Delta Lake table",
		},
		Schema: &types.TableSchema{
			Format:  "delta",
			Version: 1,
			Fields: []types.FieldInfo{
				{Name: "id", Type: "long", Nullable: false},
				{Name: "name", Type: "string", Nullable: true},
			},
		},
	}, nil
}

func (c *errorDatabricksClient) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	if c.errorMode == "updateDetails" {
		return errors.New("failed to update table details: API error")
	}
	return nil
}

// TestDatabricksErrorHandlingBasic tests basic error handling in the Databricks catalog
func TestDatabricksErrorHandlingBasic(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping error handling test in short mode")
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test error handling for each API call
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
			client := &errorDatabricksClient{
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
		client := &errorDatabricksClient{
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

	// Test with empty workspace ID
	t.Run("Empty workspace ID", func(t *testing.T) {
		client := &errorDatabricksClient{
			baseURL:     "http://localhost:8080",
			token:       "test-token",
			httpClient:  &http.Client{},
			workspaceID: "test-workspace",
		}

		emptyCatalog := &DatabricksCatalog{
			client:      client,
			workspaceID: "", // Empty workspace ID
		}

		_, err := emptyCatalog.ListCatalogs(ctx)
		assert.Error(t, err, "Should error when workspace ID is empty")
	})
}
