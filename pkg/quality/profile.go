package quality

import (
	"time"
)

// Profile represents a data quality profile for a table
type Profile struct {
	// General information
	TableName  string    `json:"table_name"`
	Timestamp  time.Time `json:"timestamp"`
	RowCount   int64     `json:"row_count"`
	
	// Column profiles
	Columns map[string]*ColumnProfile `json:"columns"`
}

// ColumnProfile represents a data quality profile for a column
type ColumnProfile struct {
	// Basic statistics
	Stats *ColumnStats `json:"stats"`
	
	// Data patterns
	Patterns []string `json:"patterns,omitempty"`
	
	// Value distribution
	Distribution *ValueDistribution `json:"distribution,omitempty"`
}

// ColumnStats represents basic statistics for a column
type ColumnStats struct {
	Count         int64   `json:"count"`
	NullCount     int64   `json:"null_count"`
	DistinctCount int64   `json:"distinct_count"`
	MinValue      string  `json:"min_value,omitempty"`
	MaxValue      string  `json:"max_value,omitempty"`
	MeanValue     float64 `json:"mean_value,omitempty"`
	MedianValue   float64 `json:"median_value,omitempty"`
}

// ValueDistribution represents the distribution of values in a column
type ValueDistribution struct {
	TopValues     map[string]int64 `json:"top_values"`
	Frequencies   map[string]int64 `json:"frequencies"`
	Percentiles   []float64        `json:"percentiles"`
	StandardDev   float64          `json:"standard_dev"`
	Skewness      float64          `json:"skewness"`
	Kurtosis      float64          `json:"kurtosis"`
}
