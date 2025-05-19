package databricks

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/nessi-dev/nessi/pkg/datalake"
)

// DatabricksCatalog provides access to Databricks Unity Catalog
type DatabricksCatalog struct {
	client      DatabricksAPI
	workspaceID string
	catalog     string
	schema      string
}

// NewDatabricksCatalog creates a new Databricks catalog
func NewDatabricksCatalog(client DatabricksAPI) *DatabricksCatalog {
	// Get configuration from environment variables
	workspaceID := os.Getenv("DATABRICKS_WORKSPACE_ID")
	catalog := os.Getenv("DATABRICKS_DEFAULT_CATALOG")
	schema := os.Getenv("DATABRICKS_DEFAULT_SCHEMA")

	// Use defaults if not provided
	if workspaceID == "" {
		workspaceID = "default"
	}
	if catalog == "" {
		catalog = "main"
	}
	if schema == "" {
		schema = "default"
	}

	return &DatabricksCatalog{
		client:      client,
		workspaceID: workspaceID,
		catalog:     catalog,
		schema:      schema,
	}
}

// ListCatalogs returns a list of catalogs
func (c *DatabricksCatalog) ListCatalogs(ctx context.Context) ([]string, error) {
	catalogs, err := c.client.GetCatalogs(ctx, c.workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list catalogs: %w", err)
	}

	catalogNames := make([]string, 0, len(catalogs))
	for _, catalog := range catalogs {
		catalogNames = append(catalogNames, catalog.Name)
	}

	return catalogNames, nil
}

// ListSchemas returns a list of schemas in a catalog
func (c *DatabricksCatalog) ListSchemas(ctx context.Context, catalogName string) ([]string, error) {
	schemas, err := c.client.GetSchemas(ctx, c.workspaceID, catalogName)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}

	schemaNames := make([]string, 0, len(schemas))
	for _, schema := range schemas {
		schemaNames = append(schemaNames, schema.Name)
	}

	return schemaNames, nil
}

// ListDatabases returns a list of databases in the catalog
func (c *DatabricksCatalog) ListDatabases(ctx context.Context) ([]types.DatabaseInfo, error) {
	schemas, err := c.client.GetSchemas(ctx, c.workspaceID, c.catalog)
	if err != nil {
		return nil, fmt.Errorf("failed to list schemas: %w", err)
	}

	databases := make([]types.DatabaseInfo, 0, len(schemas))
	for _, schema := range schemas {
		databases = append(databases, types.DatabaseInfo{
			Name:        schema.Name,
			Description: schema.Description,
			Properties: map[string]string{
				"owner":      schema.Owner,
				"created_at": schema.CreatedAt,
			},
		})
	}

	return databases, nil
}

// ListTables returns a list of tables in a catalog and schema
func (c *DatabricksCatalog) ListTables(ctx context.Context, catalogName, schemaName string) ([]string, error) {
	tables, err := c.client.GetTables(ctx, c.workspaceID, catalogName, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	tableNames := make([]string, 0, len(tables))
	for _, table := range tables {
		tableNames = append(tableNames, table.Name)
	}

	return tableNames, nil
}

// ListTableInfos returns a list of tables with details in a database
func (c *DatabricksCatalog) ListTableInfos(ctx context.Context, database string) ([]types.TableInfo, error) {
	tables, err := c.client.GetTables(ctx, c.workspaceID, c.catalog, database)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	tableInfos := make([]types.TableInfo, 0, len(tables))
	for _, table := range tables {
		tableInfos = append(tableInfos, types.TableInfo{
			Name:        table.Name,
			Type:        table.Format,
			Description: table.Description,
			Properties: map[string]string{
				"owner":      table.Owner,
				"created_at": table.CreatedAt,
			},
		})
	}

	return tableInfos, nil
}

// GetTableDetails returns details for a table
func (c *DatabricksCatalog) GetTableDetails(ctx context.Context, database, table string) (*types.TableDetails, error) {
	details, err := c.client.GetTableDetails(ctx, c.workspaceID, c.catalog, database, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get table details: %w", err)
	}

	return details, nil
}

// GetDeltaTableDetails returns details for a Delta Lake table
func (c *DatabricksCatalog) GetDeltaTableDetails(ctx context.Context, path string) (*types.TableDetails, error) {
	// Use the Delta Lake format handler to get table metadata
	return datalake.GetDeltaTableMetadata(ctx, path)
}

// SetCatalog sets the default catalog
func (c *DatabricksCatalog) SetCatalog(catalog string) {
	c.catalog = catalog
}

// SetSchema sets the default schema
func (c *DatabricksCatalog) SetSchema(schema string) {
	c.schema = schema
}

// GetFullTableName returns the fully qualified table name
func (c *DatabricksCatalog) GetFullTableName(database, table string) string {
	if database == "" {
		database = c.schema
	}
	return fmt.Sprintf("%s.%s.%s", c.catalog, database, table)
}

// ParseFullTableName parses a fully qualified table name into catalog, schema, and table
func (c *DatabricksCatalog) ParseFullTableName(fullName string) (string, string, string, error) {
	parts := strings.Split(fullName, ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid table name format: %s, expected catalog.schema.table", fullName)
	}
	return parts[0], parts[1], parts[2], nil
}
