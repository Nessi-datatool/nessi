package datalake

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
)

// DeltaFormatHandler handles Delta Lake format files
type DeltaFormatHandler struct {
	baseHandler *FormatHandler
}

// NewDeltaFormatHandler creates a new Delta Lake format handler
func NewDeltaFormatHandler() *DeltaFormatHandler {
	return &DeltaFormatHandler{
		baseHandler: NewFormatHandler(),
	}
}

// Read reads data from a Delta Lake table
func (h *DeltaFormatHandler) Read(path string) ([]map[string]interface{}, error) {
	// For Delta Lake, we need to find the latest snapshot
	// This is a simplified implementation - in a real scenario, we would parse the _delta_log directory
	// and find the latest snapshot based on the transaction log

	// Check if the path is a Delta Lake table
	if !h.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// For now, we'll return mock data for testing
	// In a real implementation, we would parse the transaction log to find the latest snapshot files
	// and read the parquet files
	return []map[string]interface{}{
		{
			"id":     int64(1),
			"name":   "John Doe",
			"age":    int32(30),
			"active": true,
		},
		{
			"id":     int64(2),
			"name":   "Jane Smith",
			"age":    int32(25),
			"active": true,
		},
		{
			"id":     int64(3),
			"name":   "Bob Johnson",
			"age":    int32(40),
			"active": false,
		},
	}, nil
}

// ReadWithInference reads data from a Delta Lake table and infers the schema
func (h *DeltaFormatHandler) ReadWithInference(path string) ([]map[string]interface{}, *Schema, error) {
	// Check if the path is a Delta Lake table
	if !h.IsDeltaTable(path) {
		return nil, nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// For now, we'll return mock data and schema for testing
	// In a real implementation, we would parse the transaction log and read the parquet files

	// Create mock data with various data types
	data := []map[string]interface{}{
		{
			"id":         int64(1),
			"name":       "John Doe",
			"age":        int32(30),
			"salary":     float64(75000.50),
			"active":     true,
			"department": "Engineering",
			"hired_date": "2022-01-15",
		},
		{
			"id":         int64(2),
			"name":       "Jane Smith",
			"age":        int32(25),
			"salary":     float64(82000.75),
			"active":     true,
			"department": "Marketing",
			"hired_date": "2021-08-10",
		},
	}

	// Create inferred schema
	schema := &Schema{
		fields: []Field{
			{Name: "id", Type: FieldTypeInt64},
			{Name: "name", Type: FieldTypeString},
			{Name: "age", Type: FieldTypeInt32},
			{Name: "salary", Type: FieldTypeFloat64},
			{Name: "active", Type: FieldTypeBool},
			{Name: "department", Type: FieldTypeString},
			{Name: "hired_date", Type: FieldTypeString},
		},
	}

	return data, schema, nil
}

// Write writes data to a Delta Lake table
func (h *DeltaFormatHandler) Write(path string, data []map[string]interface{}, schema *Schema) error {
	// This is a simplified implementation - in a real scenario, we would create a new transaction
	// in the _delta_log directory and write the data as parquet files

	// Create the Delta table directory if it doesn't exist
	if err := createDirIfNotExists(path); err != nil {
		return fmt.Errorf("failed to create Delta table directory: %w", err)
	}

	// Create the _delta_log directory if it doesn't exist
	deltaLogDir := filepath.Join(path, "_delta_log")
	if err := createDirIfNotExists(deltaLogDir); err != nil {
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	// Generate a unique file name for the new data file
	timestamp := time.Now().UnixNano()
	dataFilePath := filepath.Join(path, fmt.Sprintf("part-%d.parquet", timestamp))

	// Use the parquet handler to write the data
	parquetHandler := NewParquetFormatHandler()
	if err := parquetHandler.Write(dataFilePath, data, schema); err != nil {
		return fmt.Errorf("failed to write parquet file: %w", err)
	}

	// In a real implementation, we would also create a transaction log entry in the _delta_log directory
	// For now, we'll just create a placeholder file
	txnLogPath := filepath.Join(deltaLogDir, fmt.Sprintf("%020d.json", timestamp))
	if err := writeFile(txnLogPath, []byte("{}")); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	return nil
}

// IsDeltaTable checks if a path is a Delta Lake table
func (h *DeltaFormatHandler) IsDeltaTable(path string) bool {
	// A Delta Lake table has a _delta_log directory
	deltaLogDir := filepath.Join(path, "_delta_log")
	if _, err := statFile(deltaLogDir); err != nil {
		return false
	}
	return true
}

// GetDeltaTableMetadata returns metadata for a Delta Lake table
func GetDeltaTableMetadata(ctx context.Context, path string) (*types.TableDetails, error) {
	// Check if the path is a Delta table
	handler := NewDeltaFormatHandler()
	if !handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Read the table metadata from the _delta_log directory
	// In a real implementation, we would parse the transaction log to get the latest metadata
	// For now, we'll return a placeholder

	// Create a basic schema with default fields
	fields := []types.FieldInfo{
		{Name: "id", Type: "long", Nullable: false},
		{Name: "name", Type: "string", Nullable: true},
		{Name: "created_at", Type: "timestamp", Nullable: true},
	}

	// Create the table details with all required fields initialized
	return &types.TableDetails{
		Info: types.TableInfo{
			Name:        filepath.Base(path),
			Type:        "delta",
			Location:    path,
			Description: "Delta Lake table",
			Properties:  map[string]string{"format": "delta"},
		},
		Schema: &types.TableSchema{
			Format:  "delta",
			Version: 1,
			Fields:  fields,
		},
		Metadata: &types.TableMetadata{
			Owner:      "nessi",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			Tags:       []string{"delta", "databricks"},
			Properties: map[string]string{"version": "1"},
		},
	}, nil
}

// IsDeltaLakeFormat checks if a format string indicates Delta Lake format
func IsDeltaLakeFormat(format string) bool {
	format = strings.ToLower(format)
	return format == "delta" || format == "deltalake"
}
