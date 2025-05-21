package databricks

import (
	"context"
	"errors"
	"testing"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

// mockErrorClient is a mock implementation that returns errors
type mockErrorClient struct {
	errorMode string
}

func (c *mockErrorClient) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	if c.errorMode == "workspaces" {
		return nil, errors.New("failed to get workspaces")
	}
	return []Workspace{{ID: "test-workspace", Name: "Test Workspace"}}, nil
}

func (c *mockErrorClient) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	if c.errorMode == "catalogs" {
		return nil, errors.New("failed to get catalogs")
	}
	return []Catalog{{Name: "main"}}, nil
}

func (c *mockErrorClient) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	if c.errorMode == "schemas" {
		return nil, errors.New("failed to get schemas")
	}
	return []Schema{{Name: "default", CatalogName: "main"}}, nil
}

func (c *mockErrorClient) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	if c.errorMode == "tables" {
		return nil, errors.New("failed to get tables")
	}
	return []Table{{Name: "delta_table", CatalogName: "main", SchemaName: "default"}}, nil
}

func (c *mockErrorClient) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	if c.errorMode == "tableDetails" {
		return nil, errors.New("failed to get table details")
	}
	return &types.TableDetails{
		Info: types.TableInfo{
			Name:     "delta_table",
			Type:     "delta",
			Location: "/tmp/delta-table",
		},
	}, nil
}

func (c *mockErrorClient) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	if c.errorMode == "updateDetails" {
		return errors.New("failed to update table details")
	}
	return nil
}

// TestDatabricksClientErrors tests error handling in the Databricks client
func TestDatabricksClientErrors(t *testing.T) {
	ctx := context.Background()

	// Test error handling for each API call
	testCases := []struct {
		name      string
		errorMode string
		testFunc  func(client DatabricksAPI) error
	}{
		{
			name:      "GetWorkspaces error",
			errorMode: "workspaces",
			testFunc: func(client DatabricksAPI) error {
				_, err := client.GetWorkspaces(ctx)
				return err
			},
		},
		{
			name:      "GetCatalogs error",
			errorMode: "catalogs",
			testFunc: func(client DatabricksAPI) error {
				_, err := client.GetCatalogs(ctx, "test-workspace")
				return err
			},
		},
		{
			name:      "GetSchemas error",
			errorMode: "schemas",
			testFunc: func(client DatabricksAPI) error {
				_, err := client.GetSchemas(ctx, "test-workspace", "main")
				return err
			},
		},
		{
			name:      "GetTables error",
			errorMode: "tables",
			testFunc: func(client DatabricksAPI) error {
				_, err := client.GetTables(ctx, "test-workspace", "main", "default")
				return err
			},
		},
		{
			name:      "GetTableDetails error",
			errorMode: "tableDetails",
			testFunc: func(client DatabricksAPI) error {
				_, err := client.GetTableDetails(ctx, "test-workspace", "main", "default", "delta_table")
				return err
			},
		},
		{
			name:      "UpdateTableDetails error",
			errorMode: "updateDetails",
			testFunc: func(client DatabricksAPI) error {
				return client.UpdateTableDetails(ctx, "test-workspace", "main", "default", "delta_table", nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client := &mockErrorClient{errorMode: tc.errorMode}
			err := tc.testFunc(client)
			assert.Error(t, err, "Should return an error when API call fails")
		})
	}

	// Test empty workspace ID in catalog
	t.Run("Empty workspace ID", func(t *testing.T) {
		client := &mockErrorClient{}
		catalog := &DatabricksCatalog{
			client:      client,
			workspaceID: "", // Empty workspace ID
		}

		_, err := catalog.ListCatalogs(ctx)
		assert.Error(t, err, "Should error when workspace ID is empty")
	})
}
