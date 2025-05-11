package profile

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/parquet/file"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
)

// ParseParquet reads Parquet data from an io.Reader and returns it as []map[string]interface{}
func ParseParquet(r io.Reader) ([]map[string]interface{}, error) {
	// Since Parquet reader requires a seekable reader, we need to read the entire content
	// into memory first if it's just an io.Reader
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// Create a bytes reader which implements io.ReadSeeker
	br := bytes.NewReader(data)

	// Create a memory allocator
	alloc := memory.NewGoAllocator()

	// Create a new Parquet file reader
	pr, err := file.NewParquetReader(br)
	if err != nil {
		return nil, fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer pr.Close()

	// Create Arrow reader
	arrowReader, err := pqarrow.NewFileReader(pr, pqarrow.ArrowReadProperties{}, alloc)
	if err != nil {
		return nil, fmt.Errorf("failed to create arrow reader: %w", err)
	}

	// Read all records
	var result []map[string]interface{}

	// Read the table with context
	ctx := context.Background()
	table, err := arrowReader.ReadTable(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read table: %w", err)
	}
	defer table.Release()

	// Convert table to record batches for easier processing
	schema := table.Schema()
	
	// Create a record batch reader from the table
	rr := array.NewTableReader(table, 1024)
	defer rr.Release()
	
	// Process each batch
	for rr.Next() {
		record := rr.Record()
		
		// Process each row in the batch
		for i := 0; i < int(record.NumRows()); i++ {
			row := make(map[string]interface{})
			
			// Process each column
			for j := 0; j < int(record.NumCols()); j++ {
				col := record.Column(j)
				field := schema.Field(j)
				
				// Skip null values
				if col.IsNull(i) {
					row[field.Name] = nil
					continue
				}
				
				// Extract value based on type
				switch col.DataType().ID() {
				case arrow.INT8:
					row[field.Name] = array.NewInt8Data(col.Data()).Value(i)
				case arrow.INT16:
					row[field.Name] = array.NewInt16Data(col.Data()).Value(i)
				case arrow.INT32:
					row[field.Name] = array.NewInt32Data(col.Data()).Value(i)
				case arrow.INT64:
					row[field.Name] = array.NewInt64Data(col.Data()).Value(i)
				case arrow.UINT8:
					row[field.Name] = array.NewUint8Data(col.Data()).Value(i)
				case arrow.UINT16:
					row[field.Name] = array.NewUint16Data(col.Data()).Value(i)
				case arrow.UINT32:
					row[field.Name] = array.NewUint32Data(col.Data()).Value(i)
				case arrow.UINT64:
					row[field.Name] = array.NewUint64Data(col.Data()).Value(i)
				case arrow.FLOAT32:
					row[field.Name] = array.NewFloat32Data(col.Data()).Value(i)
				case arrow.FLOAT64:
					row[field.Name] = array.NewFloat64Data(col.Data()).Value(i)
				case arrow.STRING:
					row[field.Name] = array.NewStringData(col.Data()).Value(i)
				case arrow.BOOL:
					row[field.Name] = array.NewBooleanData(col.Data()).Value(i)
				default:
					// For unsupported types, use string representation
					row[field.Name] = col.ValueStr(i)
				}
			}
			
			result = append(result, row)
		}
	}
	
	if rr.Err() != nil {
		return nil, fmt.Errorf("error reading records: %w", rr.Err())
	}
	
	return result, nil
}

// ParseGzipParquet reads Gzip-compressed Parquet data from an io.Reader
func ParseGzipParquet(reader io.Reader) ([]map[string]interface{}, error) {
	// Create a gzip reader
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	// Parse the Parquet data
	return ParseParquet(gz)
}
