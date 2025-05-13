package catalog

import (
	"context"
	"time"
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

// DatabaseInfo contains information about a database in a data catalog
type DatabaseInfo struct {
	Name        string
	Description string
	Properties  map[string]string
}

// TableInfo contains information about a table in a data catalog
type TableInfo struct {
	Name        string
	Type        string // delta, parquet, csv, etc.
	Description string
	Location    string // S3 URI, Azure Blob path, etc.
	Properties  map[string]string
}

// TableDetails contains detailed information about a table
type TableDetails struct {
	Info        TableInfo
	Schema      *TableSchema
	Metadata    *TableMetadata
	Lineage     *LineageInfo
	Statistics  *TableStatistics
}

// TableSchema represents the schema of a table
type TableSchema struct {
	Fields      []FieldInfo
	Format      string
	Version     int64
}

// FieldInfo contains information about a field in a table schema
type FieldInfo struct {
	Name        string
	Type        string
	Description string
	Nullable    bool
	Tags        []string
	Properties  map[string]string
}

// TableMetadata contains metadata about a table
type TableMetadata struct {
	Owner       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Tags        []string
	Properties  map[string]string
}

// LineageInfo contains lineage information for a table
type LineageInfo struct {
	Upstream    []TableReference
	Downstream  []TableReference
	Process     string
	ProcessDetails map[string]string
}

// TableReference references a table in a catalog
type TableReference struct {
	Database    string
	Table       string
	Catalog     string
}

// TableStatistics contains statistics about a table
type TableStatistics struct {
	RowCount    int64
	SizeBytes   int64
	LastUpdated time.Time
	ColumnStats map[string]ColumnStatistics
}

// ColumnStatistics contains statistics about a column
type ColumnStatistics struct {
	Min         string
	Max         string
	NullCount   int64
	DistinctCount int64
}

// QualityMetrics contains data quality metrics for a table
type QualityMetrics struct {
	OverallScore    float64
	Completeness    float64
	Accuracy        float64
	Consistency     float64
	Timeliness      float64
	RuleResults     []RuleResult
	LastUpdated     time.Time
}

// RuleResult contains the result of a data quality rule
type RuleResult struct {
	RuleName    string
	RuleType    string
	Passed      bool
	Score       float64
	Details     string
}
