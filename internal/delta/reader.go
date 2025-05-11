// Implement a Delta Lake reader in Go with these requirements:
// 1. DeltaTable struct representing a Delta table with:
//    - Path to the table
//    - Current version
//    - Schema (as a structured type)
//    - Metadata (as a map)
// 2. Functions for Delta Log operations:
//    - OpenTable(path string) (*DeltaTable, error) - Open a Delta table
//    - GetTableMetadata(table *DeltaTable) (map[string]interface{}, error)
//    - GetTableVersion(table *DeltaTable, version int) (*DeltaTable, error) - Time travel
//    - ListVersions(table *DeltaTable) ([]int, error) - List available versions
//    - GetTableSchema(table *DeltaTable) (Schema, error) - Get table schema
// 3. Data reading functions:
//    - ReadFiles(table *DeltaTable) ([]string, error) - Get list of data files
//    - ReadParquetFile(file string) ([]map[string]interface{}, error) - Read a Parquet file
//    - ReadTable(table *DeltaTable) ([]map[string]interface{}, error) - Read entire table
// 4. Helper functions:
//    - ParseDeltaLog(logPath string) ([]Operation, error) - Parse transaction log
//    - ApplyOperations(base *DeltaTable, ops []Operation) *DeltaTable - Apply operations
//    - ValidateTable(table *DeltaTable) error - Validate table integrity
// Use only Go libraries without Python dependencies
// For Parquet files, use github.com/xitongsys/parquet-go
// Handle errors gracefully with context-specific messages

package delta

import (
	"context"
	"fmt"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// DeltaReader provides functionality to read data from a Delta table
type DeltaReader struct {
	connector *DeltaConnector
	schema    *arrow.Schema
}

// NewReader creates a new Delta table reader
func NewReader(connector *DeltaConnector) *DeltaReader {
	return &DeltaReader{
		connector: connector,
		schema:    connector.GetSchema(),
	}
}

// Read reads all data from the Delta table
func (r *DeltaReader) Read(ctx context.Context) (arrow.Record, error) {
	// Get the latest version
	version := r.connector.GetVersion()

	// Get all Parquet files for this version
	files, err := r.connector.getParquetFiles(version)
	if err != nil {
		return nil, fmt.Errorf("failed to get Parquet files: %w", err)
	}

	// Create record builder
	builder := array.NewRecordBuilder(memory.DefaultAllocator, r.schema)
	defer builder.Release()

	// Read and combine data from all files
	for _, file := range files {
		record, err := readParquetFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read Parquet file %s: %w", file, err)
		}
		defer record.Release()

		// Append data to builders
		for i := 0; i < int(record.NumCols()); i++ {
			if err := r.appendArray(builder.Field(i), record.Column(i)); err != nil {
				return nil, fmt.Errorf("failed to append column %d: %w", i, err)
			}
		}
	}

	return builder.NewRecord(), nil
}

// appendArray appends all values from one array to a builder
func (r *DeltaReader) appendArray(builder array.Builder, arr arrow.Array) error {
	for i := 0; i < arr.Len(); i++ {
		if arr.IsNull(i) {
			builder.AppendNull()
			continue
		}

		switch builder := builder.(type) {
		case *array.StringBuilder:
			str := arr.(*array.String).Value(i)
			builder.Append(str)
		case *array.Int64Builder:
			val := arr.(*array.Int64).Value(i)
			builder.Append(val)
		case *array.Float64Builder:
			val := arr.(*array.Float64).Value(i)
			builder.Append(val)
		case *array.BooleanBuilder:
			val := arr.(*array.Boolean).Value(i)
			builder.Append(val)
		case *array.TimestampBuilder:
			val := arr.(*array.Timestamp).Value(i)
			builder.Append(val)
		default:
			return fmt.Errorf("unsupported type: %s", builder.Type())
		}
	}
	return nil
}