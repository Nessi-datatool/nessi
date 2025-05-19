package types

import "time"

// DatabaseInfo contains information about a database in a data catalog
type DatabaseInfo struct {
	// Name is the name of the database
	Name string `json:"name"`

	// Description is the description of the database
	Description string `json:"description,omitempty"`

	// Properties contains additional metadata about the database
	Properties map[string]string `json:"properties,omitempty"`
}

// TableInfo contains basic information about a table in a data catalog
type TableInfo struct {
	// Name is the name of the table
	Name string `json:"name"`

	// Type is the type of the table (e.g., "delta", "parquet", "csv")
	Type string `json:"type"`

	// Description is the description of the table
	Description string `json:"description,omitempty"`

	// Location is the physical location of the table data
	Location string `json:"location,omitempty"`

	// Properties contains additional metadata about the table
	Properties map[string]string `json:"properties,omitempty"`
}

// TableDetails contains detailed information about a table
type TableDetails struct {
	// Info contains basic table information
	Info TableInfo `json:"info"`

	// Schema contains the table schema
	Schema *TableSchema `json:"schema,omitempty"`

	// Metadata contains table metadata
	Metadata *TableMetadata `json:"metadata,omitempty"`
}

// TableSchema contains schema information for a table
type TableSchema struct {
	// Format is the format of the table (e.g., "parquet", "csv")
	Format string `json:"format"`

	// Version is the schema version
	Version int64 `json:"version"`

	// Fields contains the list of fields in the schema
	Fields []FieldInfo `json:"fields"`
}

// FieldInfo contains information about a field in a table schema
type FieldInfo struct {
	// Name is the name of the field
	Name string `json:"name"`

	// Type is the data type of the field
	Type string `json:"type"`

	// Description is the description of the field
	Description string `json:"description,omitempty"`

	// Nullable indicates if the field can be null
	Nullable bool `json:"nullable"`

	// Properties contains additional metadata about the field
	Properties map[string]string `json:"properties,omitempty"`
}

// TableMetadata contains metadata about a table
type TableMetadata struct {
	// Owner is the owner of the table
	Owner string `json:"owner"`

	// CreatedAt is when the table was created
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt is when the table was last updated
	UpdatedAt time.Time `json:"updated_at"`

	// Tags contains tags associated with the table
	Tags []string `json:"tags"`

	// Properties contains additional metadata about the table
	Properties map[string]string `json:"properties,omitempty"`
}

// QualityMetrics represents data quality metrics for a table
type QualityMetrics struct {
	// General metrics
	TotalRows     int64     `json:"total_rows"`
	LastUpdated   time.Time `json:"last_updated"`
	SchemaVersion string    `json:"schema_version"`

	// Data quality metrics
	NullRows         int64   `json:"null_rows"`
	DuplicateRows    int64   `json:"duplicate_rows"`
	InvalidRows      int64   `json:"invalid_rows"`
	DataCompleteness float64 `json:"data_completeness"`
	DataAccuracy     float64 `json:"data_accuracy"`
	DataConsistency  float64 `json:"data_consistency"`

	// Column-level metrics
	ColumnMetrics map[string]*ColumnQualityMetrics `json:"column_metrics"`
}

// ColumnQualityMetrics represents data quality metrics for a single column
type ColumnQualityMetrics struct {
	// Basic statistics
	NullCount     int64   `json:"null_count"`
	DistinctCount int64   `json:"distinct_count"`
	MinValue      string  `json:"min_value,omitempty"`
	MaxValue      string  `json:"max_value,omitempty"`
	MeanValue     float64 `json:"mean_value,omitempty"`
	MedianValue   float64 `json:"median_value,omitempty"`

	// Data quality
	InvalidCount int64   `json:"invalid_count"`
	OutlierCount int64   `json:"outlier_count"`
	PatternMatch float64 `json:"pattern_match"`

	// Historical trends
	TrendDirection string  `json:"trend_direction"`
	TrendStrength  float64 `json:"trend_strength"`

	// Custom metrics
	CustomMetrics map[string]interface{} `json:"custom_metrics,omitempty"`
}
