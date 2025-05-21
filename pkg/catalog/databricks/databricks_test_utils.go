package databricks

import (
	"context"
	"errors"
	"net/http"

	"github.com/nessi-dev/nessi/pkg/api/types"
)

// mockDatabricksClientWithTableDetails is a mock implementation for testing table details
type mockDatabricksClientWithTableDetails struct {
	baseURL        string
	token          string
	httpClient     *http.Client
	workspaceID    string
	deltaTablePath string // Path to the Delta table for testing
}

// GetWorkspaces implements the DatabricksAPI interface
func (c *mockDatabricksClientWithTableDetails) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	return []Workspace{{ID: "test-workspace", Name: "Test Workspace"}}, nil
}

// GetCatalogs implements the DatabricksAPI interface
func (c *mockDatabricksClientWithTableDetails) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	return []Catalog{{Name: "test-catalog", Description: "Test catalog"}}, nil
}

// GetSchemas implements the DatabricksAPI interface
func (c *mockDatabricksClientWithTableDetails) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	return []Schema{{Name: "test-schema", CatalogName: catalogName}}, nil
}

// GetTables implements the DatabricksAPI interface
func (c *mockDatabricksClientWithTableDetails) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	return []Table{{Name: "test-table", CatalogName: catalogName, SchemaName: schemaName}}, nil
}

// GetTableDetails implements the DatabricksAPI interface
func (c *mockDatabricksClientWithTableDetails) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	// Use the provided deltaTablePath if available, otherwise use a default path
	location := "/path/to/delta/table"
	if c.deltaTablePath != "" {
		location = c.deltaTablePath
	}

	return &types.TableDetails{
		Info: types.TableInfo{
			Name:        tableName,
			Type:        "delta",
			Description: "Test table for integration testing",
			Location:    location,
		},
		Schema: &types.TableSchema{
			Format:  "delta",
			Version: 1,
			Fields: []types.FieldInfo{
				{Name: "id", Type: "long", Nullable: false},
				{Name: "name", Type: "string", Nullable: true},
				{Name: "age", Type: "integer", Nullable: true},
				{Name: "active", Type: "boolean", Nullable: true},
				{Name: "salary", Type: "double", Nullable: true},
				{Name: "hire_date", Type: "timestamp", Nullable: true},
			},
		},
	}, nil
}

// UpdateTableDetails implements the DatabricksAPI interface
func (c *mockDatabricksClientWithTableDetails) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	return nil
}

// mockDatabricksClientWithErrors is a mock implementation that returns errors for testing
type mockDatabricksClientWithErrors struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	workspaceID string
	errorMode   string
}

// getError returns an error based on the error mode
func (c *mockDatabricksClientWithErrors) getError() error {
	switch c.errorMode {
	case "workspaces":
		return errors.New("failed to get workspaces")
	case "catalogs":
		return errors.New("failed to get catalogs")
	case "schemas":
		return errors.New("failed to get schemas")
	case "tables":
		return errors.New("failed to get tables")
	case "tableDetails":
		return errors.New("failed to get table details")
	case "invalidPath":
		return errors.New("invalid path")
	default:
		return nil
	}
}

// GetWorkspaces implements the DatabricksAPI interface
func (c *mockDatabricksClientWithErrors) GetWorkspaces(ctx context.Context) ([]Workspace, error) {
	if c.errorMode == "workspaces" {
		return nil, c.getError()
	}
	return []Workspace{{ID: "test-workspace", Name: "Test Workspace"}}, nil
}

// GetCatalogs implements the DatabricksAPI interface
func (c *mockDatabricksClientWithErrors) GetCatalogs(ctx context.Context, workspaceID string) ([]Catalog, error) {
	if c.errorMode == "catalogs" {
		return nil, c.getError()
	}
	return []Catalog{{Name: "test-catalog"}}, nil
}

// GetSchemas implements the DatabricksAPI interface
func (c *mockDatabricksClientWithErrors) GetSchemas(ctx context.Context, workspaceID, catalogName string) ([]Schema, error) {
	if c.errorMode == "schemas" {
		return nil, c.getError()
	}
	return []Schema{{Name: "test-schema", CatalogName: catalogName}}, nil
}

// GetTables implements the DatabricksAPI interface
func (c *mockDatabricksClientWithErrors) GetTables(ctx context.Context, workspaceID, catalogName, schemaName string) ([]Table, error) {
	if c.errorMode == "tables" {
		return nil, c.getError()
	}
	return []Table{{Name: "test-table", CatalogName: catalogName, SchemaName: schemaName}}, nil
}

// GetTableDetails implements the DatabricksAPI interface
func (c *mockDatabricksClientWithErrors) GetTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string) (*types.TableDetails, error) {
	if c.errorMode == "tableDetails" {
		return nil, c.getError()
	}
	return &types.TableDetails{
		Info: types.TableInfo{
			Name:        tableName,
			Type:        "delta",
			Description: "Test table for error testing",
			Location:    "/path/to/delta/table",
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

// UpdateTableDetails implements the DatabricksAPI interface
func (c *mockDatabricksClientWithErrors) UpdateTableDetails(ctx context.Context, workspaceID, catalogName, schemaName, tableName string, details *types.TableDetails) error {
	return c.getError()
}
