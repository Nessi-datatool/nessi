package catalog

import (
	"context"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
)

// CatalogType represents the type of data catalog
type CatalogType string

// Supported catalog types
const (
	AWSGlue       CatalogType = "aws_glue"
	AzurePurview  CatalogType = "azure_purview"
	GCPDataCatalog CatalogType = "gcp_data_catalog"
)

// DatabaseInfo represents a database in a data catalog
type DatabaseInfo struct {
	Name        string
	Description string
	Location    string
	CreateTime  time.Time
	UpdateTime  time.Time
}

// TableInfo represents a table in a data catalog
type TableInfo struct {
	Name        string
	Description string
	Type        string
	Location    string
	CreateTime  time.Time
	UpdateTime  time.Time
}

// ColumnInfo represents a column in a table
type ColumnInfo struct {
	Name        string
	Type        string
	Description string
	IsPartition bool
	IsSortKey   bool
}

// TableDetails represents detailed information about a table
type TableDetails struct {
	Table    TableInfo
	Columns  []ColumnInfo
	Metadata TableMetadata
}

// TableMetadata represents metadata for a table
type TableMetadata struct {
	Owner      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Tags       []string
	Properties map[string]string
}

// LineageInfo represents lineage information for a table
type LineageInfo struct {
	Upstream   []string
	Downstream []string
	UpdatedAt  time.Time
	UpdatedBy  string
}

// RuleResult represents the result of a data quality rule
type RuleResult struct {
	RuleName string
	RuleType string
	Passed   bool
	Score    float64
	Details  string
}

// QualityMetrics represents data quality metrics for a table
type QualityMetrics struct {
	OverallScore float64
	Completeness float64
	Accuracy     float64
	Consistency  float64
	Timeliness   float64
	LastUpdated  time.Time
	RuleResults  []RuleResult
}

// DataCatalog defines the interface for interacting with data catalogs
type DataCatalog interface {
	// Core operations
	Name() string
	Connect(ctx context.Context, config map[string]interface{}) error
	Disconnect(ctx context.Context) error
	
	// Asset discovery
	ListDatabases(ctx context.Context) ([]types.DatabaseInfo, error)
	ListTables(ctx context.Context, database string) ([]types.TableInfo, error)
	GetTableDetails(ctx context.Context, database, table string) (*types.TableDetails, error)
	
	// Metadata operations
	GetTableMetadata(ctx context.Context, database, table string) (*types.TableMetadata, error)
	UpdateTableMetadata(ctx context.Context, database, table string, metadata *types.TableMetadata) error
	
	// Lineage operations
	GetTableLineage(ctx context.Context, database, table string) (*types.LineageInfo, error)
	UpdateTableLineage(ctx context.Context, database, table string, lineage *types.LineageInfo) error
	
	// Quality metrics
	PublishQualityMetrics(ctx context.Context, database, table string, metrics *types.QualityMetrics) error
	GetQualityMetrics(ctx context.Context, database, table string) (*types.QualityMetrics, error)
}
