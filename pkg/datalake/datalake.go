// Package datalake provides functionality for reading and writing Delta tables
package datalake

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
)

// DeltaTable represents a Delta table
type DeltaTable struct {
	Path            string
	Version         int64
	LastModified    time.Time
	Schema          *arrow.Schema
	PartitionSchema []string
	Stats           *TableStats
	Checkpoint      *Checkpoint
	Files           []string
	Partitions      map[string][]string
	Metadata        map[string]interface{}
	Errors          []error
}

// TableStats represents statistics about a Delta table
type TableStats struct {
	NumFiles        int64
	NumRecords      int64
	TotalSize       int64
	PartitionCounts map[string]int64
	ColumnStats     map[string]*ColumnStats
}

// ColumnStats represents statistics for a single column
type ColumnStats struct {
	NullCount     int64
	DistinctCount int64
	MinValue      interface{}
	MaxValue      interface{}
	AvgValue      float64
}

// Checkpoint represents a Delta table checkpoint
type Checkpoint struct {
	Version    int64
	Timestamp  time.Time
	FileCount  int64
	FilePaths  []string
	Schema     *arrow.Schema
	Partitions []string
}

// Reader handles reading from Delta tables
type Reader struct {
	table     *DeltaTable
	closed    bool
	mu        sync.Mutex
	fileCache map[string]*os.File
}

// NewReader creates a new Delta table reader
func NewReader(tablePath string) (*Reader, error) {
	// Create metadata manager
	metadataManager := NewMetadataManager(tablePath)

	// Read table metadata
	table, err := metadataManager.ReadTableMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to read table metadata: %w", err)
	}

	return &Reader{
		table:     table,
		closed:    false,
		fileCache: make(map[string]*os.File),
	}, nil
}

// Initialize sets up the reader by reading table metadata
func (r *Reader) Initialize() error {
	if r.closed {
		return fmt.Errorf("reader is closed")
	}

	// Create metadata manager
	metadataManager := NewMetadataManager(r.table.Path)

	// Read table metadata
	table, err := metadataManager.ReadTableMetadata()
	if err != nil {
		return fmt.Errorf("failed to read table metadata: %w", err)
	}

	// Update table
	r.table = table

	return nil
}

// ReadPartition reads data from a specific partition
func (r *Reader) ReadPartition(partition string) (arrow.Record, error) {
	if r.closed {
		return nil, fmt.Errorf("reader is closed")
	}

	// Check if partition exists
	files, ok := r.table.Partitions[partition]
	if !ok {
		return nil, fmt.Errorf("partition not found: %s", partition)
	}

	// Create Parquet manager
	parquetManager := NewParquetManager()

	// Read first file in partition
	if len(files) == 0 {
		return nil, fmt.Errorf("no files in partition: %s", partition)
	}

	return parquetManager.ReadRecord(files[0], r.table.Schema)
}

// GetStats returns the table statistics
func (r *Reader) GetStats() (*TableStats, error) {
	if r.closed {
		return nil, fmt.Errorf("reader is closed")
	}

	// Create Parquet manager
	parquetManager := NewParquetManager()

	// Calculate stats from all files
	stats := &TableStats{
		NumFiles:        0,
		NumRecords:      0,
		TotalSize:       0,
		PartitionCounts: make(map[string]int64),
		ColumnStats:     make(map[string]*ColumnStats),
	}

	// Process each file
	for _, file := range r.table.Files {
		fileStats, err := parquetManager.GetFileStats(file)
		if err != nil {
			return nil, fmt.Errorf("failed to get file stats: %w", err)
		}

		stats.NumFiles++
		stats.NumRecords += fileStats.NumRows
		stats.TotalSize += fileStats.Size

		// Update partition counts
		partition := filepath.Base(filepath.Dir(file))
		stats.PartitionCounts[partition]++

		// Merge column stats
		for colName, colStats := range fileStats.Columns {
			if _, ok := stats.ColumnStats[colName]; !ok {
				stats.ColumnStats[colName] = &ColumnStats{}
			}
			// Update column stats (simplified for now)
			stats.ColumnStats[colName].NullCount += colStats.NullCount
			stats.ColumnStats[colName].DistinctCount += colStats.DistinctCount
		}
	}

	return stats, nil
}

// ReadAllStructured reads all data from the table and returns it as structured records
// This is needed for the internal/quality/profile/Profiler which expects []map[string]interface{}
func (r *Reader) ReadAllStructured() ([]map[string]interface{}, error) {
	if r.closed {
		return nil, fmt.Errorf("reader is closed")
	}

	// Create Parquet manager
	parquetManager := NewParquetManager()

	// Read all files and combine records
	var result []map[string]interface{}

	for _, file := range r.table.Files {
		// Read record from file
		record, err := parquetManager.ReadRecord(file, r.table.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to read record from file %s: %w", file, err)
		}

		// Convert record to structured data
		structured, err := r.ConvertToStructuredData(record)
		if err != nil {
			return nil, fmt.Errorf("failed to convert record to structured data: %w", err)
		}

		// Append to result
		result = append(result, structured...)
	}

	return result, nil
}

// Close closes the reader
func (r *Reader) Close() error {
	r.closed = true
	return nil
}

// IsClosed returns whether the reader is closed
func (r *Reader) IsClosed() bool {
	return r.closed
}

// GetSchema returns the table schema
func (r *Reader) GetSchema() *arrow.Schema {
	return r.table.Schema
}

// ConvertToStructuredData converts an Arrow record to structured data
func (r *Reader) ConvertToStructuredData(record arrow.Record) ([]map[string]interface{}, error) {
	if record == nil {
		return nil, fmt.Errorf("record is nil")
	}

	result := make([]map[string]interface{}, record.NumRows())

	// Create a map for each row
	for i := int64(0); i < record.NumRows(); i++ {
		row := make(map[string]interface{})

		// Add each column's value to the map
		for j, col := range record.Columns() {
			fieldName := record.Schema().Field(j).Name

			// Handle null values
			if col.IsNull(int(i)) {
				row[fieldName] = nil
				continue
			}

			// Extract value based on type
			switch col := col.(type) {
			case *array.Int8:
				row[fieldName] = col.Value(int(i))
			case *array.Int16:
				row[fieldName] = col.Value(int(i))
			case *array.Int32:
				row[fieldName] = col.Value(int(i))
			case *array.Int64:
				row[fieldName] = col.Value(int(i))
			case *array.Uint8:
				row[fieldName] = col.Value(int(i))
			case *array.Uint16:
				row[fieldName] = col.Value(int(i))
			case *array.Uint32:
				row[fieldName] = col.Value(int(i))
			case *array.Uint64:
				row[fieldName] = col.Value(int(i))
			case *array.Float32:
				row[fieldName] = col.Value(int(i))
			case *array.Float64:
				row[fieldName] = col.Value(int(i))
			case *array.String:
				row[fieldName] = col.Value(int(i))
			case *array.LargeString:
				row[fieldName] = col.Value(int(i))
			case *array.Binary:
				row[fieldName] = string(col.Value(int(i)))
			case *array.LargeBinary:
				row[fieldName] = string(col.Value(int(i)))
			case *array.Timestamp:
				row[fieldName] = col.Value(int(i)) // returns int64 nanoseconds since epoch
			case *array.Decimal128:
				val := col.Value(int(i)).ToFloat64(4) // 4 = scale, matches schema
				row[fieldName] = val
			case *array.Struct:
				structMap := make(map[string]interface{})
				for k := 0; k < col.NumField(); k++ {
					field := col.Field(k)
					name := col.DataType().(*arrow.StructType).Field(k).Name
					if field.IsNull(int(i)) {
						structMap[name] = nil
						continue
					}
					switch f := field.(type) {
					case *array.String:
						structMap[name] = f.Value(int(i))
					case *array.Int32:
						structMap[name] = f.Value(int(i))
					case *array.Int64:
						structMap[name] = f.Value(int(i))
					default:
						structMap[name] = nil
					}
				}
				row[fieldName] = structMap
			case *array.List:
				// Only handle list<string> for now
				if col.DataType().(*arrow.ListType).Elem().ID() == arrow.STRING {
					valArr := col.ListValues().(*array.String)
					start := col.Offsets()[int(i)]
					end := col.Offsets()[int(i)+1]
					var vals []string
					for idx := start; idx < end; idx++ {
						vals = append(vals, valArr.Value(int(idx)))
					}
					row[fieldName] = vals
				} else {
					row[fieldName] = nil
				}
			default:
				fmt.Fprintf(os.Stderr, "[ConvertToStructuredData] Unsupported column type: %T\n", col)
				row[fieldName] = nil
			}
		}

		result[i] = row
	}

	return result, nil
}
