package catalog

import (
	"context"
	"fmt"
	"strings"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/nessi-dev/nessi/pkg/catalog/databricks"
)

// DatabricksCatalog implements the Catalog interface for Databricks
type DatabricksCatalog struct {
	client        databricks.DatabricksAPI
	workspaceID   string
	defaultSchema string
}

// NewDatabricksCatalog creates a new Databricks catalog
func NewDatabricksCatalog(baseURL, token, workspaceID, defaultSchema string) (*DatabricksCatalog, error) {
	if baseURL == "" || token == "" || workspaceID == "" {
		return nil, fmt.Errorf("baseURL, token, and workspaceID are required")
	}

	client := databricks.NewDatabricksClient(baseURL, token, workspaceID)

	return &DatabricksCatalog{
		client:        client,
		workspaceID:   workspaceID,
		defaultSchema: defaultSchema,
	}, nil
}

// GetName returns the name of the catalog
func (c *DatabricksCatalog) GetName() string {
	return "Databricks"
}

// GetType returns the type of the catalog
func (c *DatabricksCatalog) GetType() string {
	return "databricks"
}

// Connect implements the DataCatalog interface
func (c *DatabricksCatalog) Connect(ctx context.Context, config map[string]interface{}) error {
	// Already connected via the constructor
	return nil
}

// Disconnect implements the DataCatalog interface
func (c *DatabricksCatalog) Disconnect(ctx context.Context) error {
	// Nothing to disconnect in the current implementation
	return nil
}

// Name implements the DataCatalog interface
func (c *DatabricksCatalog) Name() string {
	return "databricks"
}

// GetTableMetadata implements the DataCatalog interface
func (c *DatabricksCatalog) GetTableMetadata(ctx context.Context, database, table string) (*types.TableMetadata, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return nil, err
	}
	return details.Metadata, nil
}

// UpdateTableMetadata implements the DataCatalog interface
func (c *DatabricksCatalog) UpdateTableMetadata(ctx context.Context, database, table string, metadata *types.TableMetadata) error {
	// Parse database name (catalog.schema)
	catalogName, schemaName, err := c.parseDatabaseName(database)
	if err != nil {
		return err
	}

	// Get current table details
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return err
	}

	// Update metadata
	details.Metadata = metadata

	// Update table details in Databricks
	return c.client.UpdateTableDetails(ctx, c.workspaceID, catalogName, schemaName, table, details)
}

// GetTableLineage implements the DataCatalog interface
func (c *DatabricksCatalog) GetTableLineage(ctx context.Context, database, table string) (*types.LineageInfo, error) {
	// Databricks Unity Catalog doesn't expose lineage information via API yet
	// Return empty lineage info
	return &types.LineageInfo{}, nil
}

// UpdateTableLineage implements the DataCatalog interface
func (c *DatabricksCatalog) UpdateTableLineage(ctx context.Context, database, table string, lineage *types.LineageInfo) error {
	// Databricks Unity Catalog doesn't support updating lineage information via API yet
	return fmt.Errorf("updating lineage information is not supported for Databricks tables")
}

// PublishQualityMetrics implements the DataCatalog interface
func (c *DatabricksCatalog) PublishQualityMetrics(ctx context.Context, database, table string, metrics *types.QualityMetrics) error {
	// Databricks Unity Catalog doesn't support quality metrics via API yet
	return fmt.Errorf("publishing quality metrics is not supported for Databricks tables")
}

// GetQualityMetrics implements the DataCatalog interface
func (c *DatabricksCatalog) GetQualityMetrics(ctx context.Context, database, table string) (*types.QualityMetrics, error) {
	// Databricks Unity Catalog doesn't expose quality metrics via API yet
	// Return empty quality metrics
	return &types.QualityMetrics{}, nil
}

// ListDatabases returns a list of databases in the catalog
func (c *DatabricksCatalog) ListDatabases(ctx context.Context) ([]types.DatabaseInfo, error) {
	// Get catalogs from Databricks
	catalogs, err := c.client.GetCatalogs(ctx, c.workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list catalogs: %w", err)
	}

	// Get schemas for each catalog
	var databases []types.DatabaseInfo
	for _, catalog := range catalogs {
		schemas, err := c.client.GetSchemas(ctx, c.workspaceID, catalog.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to list schemas for catalog %s: %w", catalog.Name, err)
		}

		// Format database names as catalog.schema
		for _, schema := range schemas {
			databases = append(databases, types.DatabaseInfo{
				Name:        fmt.Sprintf("%s.%s", catalog.Name, schema.Name),
				Description: schema.Description,
			})
		}
	}

	return databases, nil
}

// ListTables returns a list of tables in the specified database
func (c *DatabricksCatalog) ListTables(ctx context.Context, database string) ([]types.TableInfo, error) {
	// Parse database name (catalog.schema)
	catalogName, schemaName, err := c.parseDatabaseName(database)
	if err != nil {
		return nil, err
	}

	// Get tables from Databricks
	tables, err := c.client.GetTables(ctx, c.workspaceID, catalogName, schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	// Convert to TableInfo
	tableInfos := make([]types.TableInfo, 0, len(tables))
	for _, table := range tables {
		tableInfos = append(tableInfos, types.TableInfo{
			Name:        table.Name,
			Description: table.Description,
			Type:        table.Format,
		})
	}

	return tableInfos, nil
}

// GetTableDetails returns details of a specific table
func (c *DatabricksCatalog) GetTableDetails(ctx context.Context, database, table string) (*types.TableDetails, error) {
	// Parse database name (catalog.schema)
	catalogName, schemaName, err := c.parseDatabaseName(database)
	if err != nil {
		return nil, err
	}

	// Get table details from Databricks
	details, err := c.client.GetTableDetails(ctx, c.workspaceID, catalogName, schemaName, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get table details: %w", err)
	}

	return details, nil
}

// UpdateTableDetails updates the details of a specific table
func (c *DatabricksCatalog) UpdateTableDetails(ctx context.Context, database, table string, details *types.TableDetails) error {
	// Parse database name (catalog.schema)
	catalogName, schemaName, err := c.parseDatabaseName(database)
	if err != nil {
		return err
	}

	// Update table details in Databricks
	err = c.client.UpdateTableDetails(ctx, c.workspaceID, catalogName, schemaName, table, details)
	if err != nil {
		return fmt.Errorf("failed to update table details: %w", err)
	}

	return nil
}

// parseDatabaseName parses a database name in the format "catalog.schema"
func (c *DatabricksCatalog) parseDatabaseName(database string) (catalogName, schemaName string, err error) {
	parts := strings.Split(database, ".")
	if len(parts) != 2 {
		// If no schema is specified, use the default schema
		if len(parts) == 1 {
			return parts[0], c.defaultSchema, nil
		}
		return "", "", fmt.Errorf("invalid database name format: %s (expected 'catalog.schema')", database)
	}
	return parts[0], parts[1], nil
}

// GetTableLocation returns the location of a specific table
func (c *DatabricksCatalog) GetTableLocation(ctx context.Context, database, table string) (string, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return "", err
	}
	return details.Info.Location, nil
}

// GetTableFormat returns the format of a specific table
func (c *DatabricksCatalog) GetTableFormat(ctx context.Context, database, table string) (string, error) {
	details, err := c.GetTableDetails(ctx, database, table)
	if err != nil {
		return "", err
	}
	if details.Schema != nil {
		return details.Schema.Format, nil
	}
	return details.Info.Type, nil
}

// IsTableExternal returns whether a specific table is external
func (c *DatabricksCatalog) IsTableExternal(ctx context.Context, database, table string) (bool, error) {
	// In Databricks, we can determine if a table is external by checking its type
	// External tables typically have a location that starts with "s3://" or "dbfs:/mnt/"
	location, err := c.GetTableLocation(ctx, database, table)
	if err != nil {
		return false, err
	}

	// Tables with locations outside of DBFS managed locations are external
	return strings.HasPrefix(location, "s3://") ||
		strings.HasPrefix(location, "abfss://") ||
		strings.HasPrefix(location, "gs://") ||
		strings.HasPrefix(location, "dbfs:/mnt/"), nil
}
