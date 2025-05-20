package report

import (
	"time"
)

// ScanResult represents the result of a data scan
type ScanResult struct {
	ScanID             string              `json:"scan_id"`
	TablePath          string              `json:"table_path"`
	Timestamp          time.Time           `json:"timestamp"`
	Duration           time.Duration       `json:"duration"`
	RowCount           int64               `json:"row_count"`
	ColumnCount        int                 `json:"column_count"`
	QualityMetrics     *QualityMetrics     `json:"quality_metrics,omitempty"`
	PerformanceMetrics *PerformanceMetrics `json:"performance_metrics,omitempty"`
}

// QualityMetrics represents quality metrics for a data scan
type QualityMetrics struct {
	Completeness map[string]float64 `json:"completeness"`
	Accuracy     map[string]float64 `json:"accuracy"`
	Consistency  map[string]float64 `json:"consistency"`
	Uniqueness   map[string]float64 `json:"uniqueness"`
	Timeliness   map[string]float64 `json:"timeliness"`
}

// PerformanceMetrics represents performance metrics for a data scan
type PerformanceMetrics struct {
	ScanDurationMs  int64     `json:"scan_duration_ms"`
	MemoryUsageMb   float64   `json:"memory_usage_mb"`
	CpuUsagePercent float64   `json:"cpu_usage_percent"`
	IoOperations    int64     `json:"io_operations"`
	RowsProcessed   int64     `json:"rows_processed"`
	BytesProcessed  int64     `json:"bytes_processed"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
}

// SchemaValidationResult represents the result of a schema validation
type SchemaValidationResult struct {
	SchemaName       string                 `json:"schema_name"`
	SchemaDefinition map[string]interface{} `json:"schema_definition"`
	ValidationErrors []string               `json:"validation_errors,omitempty"`
	IsValid          bool                   `json:"is_valid"`
	Timestamp        time.Time              `json:"timestamp"`
}
