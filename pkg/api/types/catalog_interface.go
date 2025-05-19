package types

import (
	"context"
)

// DataCatalog defines the interface for interacting with data catalogs
type DataCatalog interface {
	// Core operations
	Name() string
	Connect(ctx context.Context, config map[string]interface{}) error
	Disconnect(ctx context.Context) error

	// Asset discovery
	ListDatabases(ctx context.Context) ([]DatabaseInfo, error)
	ListTables(ctx context.Context, database string) ([]TableInfo, error)
	GetTableDetails(ctx context.Context, database, table string) (*TableDetails, error)

	// Metadata operations
	GetTableMetadata(ctx context.Context, database, table string) (*TableMetadata, error)
	UpdateTableMetadata(ctx context.Context, database, table string, metadata *TableMetadata) error

	// Lineage operations
	GetTableLineage(ctx context.Context, database, table string) (*LineageInfo, error)
	UpdateTableLineage(ctx context.Context, database, table string, lineage *LineageInfo) error

	// Quality metrics
	PublishQualityMetrics(ctx context.Context, database, table string, metrics *QualityMetrics) error
	GetQualityMetrics(ctx context.Context, database, table string) (*QualityMetrics, error)
}
