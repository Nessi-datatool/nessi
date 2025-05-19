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
	baseHandler FormatHandler
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
	
	// For now, we'll use the parquet handler to read the latest snapshot
	// In a real implementation, we would parse the transaction log to find the latest snapshot files
	parquetHandler := NewParquetFormatHandler()
	
	// Find parquet files in the Delta table directory
	parquetFiles, err := filepath.Glob(filepath.Join(path, "*.parquet"))
	if err != nil {
		return nil, fmt.Errorf("failed to find parquet files in Delta table: %w", err)
	}
	
	if len(parquetFiles) == 0 {
		return nil, fmt.Errorf("no parquet files found in Delta table: %s", path)
	}
	
	// Read the latest parquet file (this is a simplification)
	// In a real implementation, we would use the transaction log to determine which files to read
	latestFile := parquetFiles[len(parquetFiles)-1]
	
	return parquetHandler.Read(latestFile)
}

// ReadWithInference reads data from a Delta Lake table and infers the schema
func (h *DeltaFormatHandler) ReadWithInference(path string) ([]map[string]interface{}, *Schema, error) {
	// Check if the path is a Delta Lake table
	if !h.IsDeltaTable(path) {
		return nil, nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}
	
	// For now, we'll use the parquet handler to read the latest snapshot
	parquetHandler := NewParquetFormatHandler()
	
	// Find parquet files in the Delta table directory
	parquetFiles, err := filepath.Glob(filepath.Join(path, "*.parquet"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find parquet files in Delta table: %w", err)
	}
	
	if len(parquetFiles) == 0 {
		return nil, nil, fmt.Errorf("no parquet files found in Delta table: %s", path)
	}
	
	// Read the latest parquet file (this is a simplification)
	latestFile := parquetFiles[len(parquetFiles)-1]
	
	return parquetHandler.ReadWithInference(latestFile)
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
	// Check if the path is a Delta Lake table
	handler := NewDeltaFormatHandler()
	if !handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}
	
	// Read the table metadata from the _delta_log directory
	// In a real implementation, we would parse the transaction log to get the latest metadata
	// For now, we'll return a placeholder
	
	// Try to infer the schema from a data file
	parquetFiles, err := filepath.Glob(filepath.Join(path, "*.parquet"))
	if err != nil || len(parquetFiles) == 0 {
		return &types.TableDetails{
			Name:        filepath.Base(path),
			Format:      "delta",
			Location:    path,
			Description: "Delta Lake table",
		}, nil
	}
	
	// Read the schema from a parquet file
	parquetHandler := NewParquetFormatHandler()
	_, schema, err := parquetHandler.ReadWithInference(parquetFiles[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read schema from parquet file: %w", err)
	}
	
	// Convert the schema to column details
	columns := make([]types.ColumnDetails, 0, len(schema.Fields()))
	for _, field := range schema.Fields() {
		columns = append(columns, types.ColumnDetails{
			Name:        field.Name,
			Type:        field.Type.String(),
			Description: "",
			Nullable:    true,
		})
	}
	
	return &types.TableDetails{
		Name:        filepath.Base(path),
		Format:      "delta",
		Location:    path,
		Description: "Delta Lake table",
		Columns:     columns,
	}, nil
}

// IsDeltaLakeFormat checks if a format string indicates Delta Lake format
func IsDeltaLakeFormat(format string) bool {
	format = strings.ToLower(format)
	return format == "delta" || format == "deltalake"
}
