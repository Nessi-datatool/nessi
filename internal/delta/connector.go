package delta

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/apache/arrow/go/v12/parquet"
	"github.com/apache/arrow/go/v12/parquet/file"
	"github.com/apache/arrow/go/v12/parquet/pqarrow"
)

// DeltaConnector handles Delta Lake operations
type DeltaConnector struct {
	basePath    string
	tablePath   string
	version     int64
	lastUpdated time.Time
	schema      *arrow.Schema
	partitions  []string
	metadata    *TableMetadata
}

// TableMetadata represents Delta table metadata
type TableMetadata struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Schema      *arrow.Schema     `json:"schema"`
	Partitions  []string          `json:"partitions"`
	Properties  map[string]string `json:"properties"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     int64             `json:"version"`
}

// NewConnector creates a new Delta Lake connector
func NewConnector(tablePath string) (*DeltaConnector, error) {
	if err := validateTablePath(tablePath); err != nil {
		return nil, fmt.Errorf("invalid table path: %w", err)
	}

	connector := &DeltaConnector{
		tablePath: tablePath,
		basePath:  filepath.Dir(tablePath),
	}

	// Load table metadata
	if err := connector.loadMetadata(); err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	return connector, nil
}

// loadMetadata loads the table metadata from the Delta log
func (c *DeltaConnector) loadMetadata() error {
	logPath := filepath.Join(c.tablePath, "_delta_log")
	if _, err := os.Stat(logPath); err != nil {
		return fmt.Errorf("Delta log directory not found: %w", err)
	}

	// Read all JSON files in the log directory
	files, err := os.ReadDir(logPath)
	if err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	// Sort files by version number
	var versions []int64
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			version, err := strconv.ParseInt(strings.TrimSuffix(file.Name(), ".json"), 10, 64)
			if err != nil {
				continue
			}
			versions = append(versions, version)
		}
	}
	sort.Slice(versions, func(i, j int) bool { return versions[i] < versions[j] })

	// Process each version in order
	for _, version := range versions {
		if err := c.processVersion(version); err != nil {
			return fmt.Errorf("failed to process version %d: %w", version, err)
		}
	}

	return nil
}

// processVersion processes a single version of the Delta log
func (c *DeltaConnector) processVersion(version int64) error {
	filePath := filepath.Join(c.tablePath, "_delta_log", fmt.Sprintf("%d.json", version))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	var actions []Action
	if err := json.Unmarshal(data, &actions); err != nil {
		return fmt.Errorf("failed to parse log file: %w", err)
	}

	// Process each action
	for _, action := range actions {
		if action.Meta != nil {
			c.metadata = &TableMetadata{
				ID:          action.Meta.ID,
				Name:        action.Meta.Name,
				Description: action.Meta.Description,
				Schema:      c.schema,
				Partitions:  action.Meta.PartitionColumns,
				Properties:  action.Meta.Configuration,
				CreatedAt:   time.Unix(0, action.Meta.CreatedTime*int64(time.Millisecond)),
				UpdatedAt:   time.Now(),
				Version:     version,
			}
		}
	}

	c.version = version
	return nil
}

// GetMetadata returns the current table metadata
func (c *DeltaConnector) GetMetadata() (*TableMetadata, error) {
	if c.metadata == nil {
		return nil, fmt.Errorf("no metadata available")
	}
	return c.metadata, nil
}

// ReadPartition reads data from a specific partition
func (c *DeltaConnector) ReadPartition(ctx context.Context, partition string) (arrow.Record, error) {
	partitionPath := filepath.Join(c.tablePath, partition)
	if _, err := os.Stat(partitionPath); err != nil {
		return nil, fmt.Errorf("partition not found: %w", err)
	}

	// Read all Parquet files in the partition
	files, err := os.ReadDir(partitionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read partition directory: %w", err)
	}

	var records []arrow.Record
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".parquet") {
			filePath := filepath.Join(partitionPath, file.Name())
			record, err := readParquetFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read Parquet file: %w", err)
			}
			records = append(records, record)
		}
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no data found in partition")
	}

	// Merge records if there are multiple files
	if len(records) > 1 {
		return mergeRecords(records)
	}

	return records[0], nil
}

// WritePartition writes data to a specific partition
func (c *DeltaConnector) WritePartition(ctx context.Context, partition string, record arrow.Record) error {
	partitionPath := filepath.Join(c.tablePath, partition)
	if err := os.MkdirAll(partitionPath, 0755); err != nil {
		return fmt.Errorf("failed to create partition directory: %w", err)
	}

	// Generate a unique filename
	filename := fmt.Sprintf("part-%d.parquet", time.Now().UnixNano())
	filePath := filepath.Join(partitionPath, filename)

	// Write the record to Parquet
	if err := writeParquetFile(filePath, record); err != nil {
		return fmt.Errorf("failed to write Parquet file: %w", err)
	}

	// Update the Delta log
	if err := c.updateLog(partition, filename, record.NumRows()); err != nil {
		return fmt.Errorf("failed to update Delta log: %w", err)
	}

	return nil
}

// updateLog updates the Delta log with a new transaction
func (c *DeltaConnector) updateLog(partition, filename string, numRows int64) error {
	// Create a new transaction
	transaction := struct {
		Add struct {
			Path             string            `json:"path"`
			PartitionValues  map[string]string `json:"partitionValues"`
			Size             int64             `json:"size"`
			ModificationTime int64             `json:"modificationTime"`
			DataChange       bool              `json:"dataChange"`
			Stats            string            `json:"stats,omitempty"`
		} `json:"add"`
	}{
		Add: struct {
			Path             string            `json:"path"`
			PartitionValues  map[string]string `json:"partitionValues"`
			Size             int64             `json:"size"`
			ModificationTime int64             `json:"modificationTime"`
			DataChange       bool              `json:"dataChange"`
			Stats            string            `json:"stats,omitempty"`
		}{
			Path:             filepath.Join(partition, filename),
			PartitionValues:  parsePartitionValues(partition),
			Size:             getFileSize(filepath.Join(c.tablePath, partition, filename)),
			ModificationTime: time.Now().UnixMilli(),
			DataChange:       true,
			Stats:            fmt.Sprintf(`{"numRecords":%d}`, numRows),
		},
	}

	// Write the transaction to the log
	logPath := filepath.Join(c.tablePath, "_delta_log", fmt.Sprintf("%d.json", c.version+1))
	data, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	if err := os.WriteFile(logPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	c.version++
	return nil
}

// parsePartitionValues parses partition values from a partition path
func parsePartitionValues(partition string) map[string]string {
	values := make(map[string]string)
	parts := strings.Split(partition, "/")
	for _, part := range parts {
		if strings.Contains(part, "=") {
			kv := strings.Split(part, "=")
			if len(kv) == 2 {
				values[kv[0]] = kv[1]
			}
		}
	}
	return values
}

// getFileSize returns the size of a file in bytes
func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

// mergeRecords merges multiple Arrow records into a single record
func mergeRecords(records []arrow.Record) (arrow.Record, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("no records to merge")
	}

	// Get the schema from the first record
	schema := records[0].Schema()

	// Create arrays for each column
	arrays := make([]arrow.Array, schema.NumFields())
	for i := range arrays {
		arrays[i] = array.NewBuilder(memory.DefaultAllocator, schema.Field(i).Type)
	}

	// Append data from all records
	for _, record := range records {
		for i := 0; i < record.NumCols(); i++ {
			col := record.Column(i)
			for j := 0; j < col.Len(); j++ {
				if col.IsNull(j) {
					arrays[i].AppendNull()
				} else {
					arrays[i].Append(col.ValueStr(j))
				}
			}
		}
	}

	// Create the merged record
	fields := make([]arrow.Field, schema.NumFields())
	for i := range fields {
		fields[i] = schema.Field(i)
	}

	mergedSchema := arrow.NewSchema(fields, nil)
	return array.NewRecord(mergedSchema, arrays, -1), nil
}

// GetSchema returns the current table schema
func (c *DeltaConnector) GetSchema() *arrow.Schema {
	return c.schema
}

// GetPartitions returns the current table partitions
func (c *DeltaConnector) GetPartitions() []string {
	return c.partitions
}

// GetVersion returns the current table version
func (c *DeltaConnector) GetVersion() int64 {
	return c.version
}

// GetLastUpdated returns the last update timestamp
func (c *DeltaConnector) GetLastUpdated() time.Time {
	return c.lastUpdated
}

// validateTablePath validates the Delta table path
func validateTablePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	// Check for _delta_log directory
	logPath := filepath.Join(path, "_delta_log")
	if _, err := os.Stat(logPath); err != nil {
		return fmt.Errorf("not a valid Delta table: missing _delta_log directory")
	}

	return nil
}

// readParquetFile reads a Parquet file and returns an Arrow record
func readParquetFile(path string) (arrow.Record, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader, err := file.NewParquetReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer reader.Close()

	arrowReader, err := pqarrow.NewFileReader(reader, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	if err != nil {
		return nil, fmt.Errorf("failed to create arrow reader: %w", err)
	}

	record, err := arrowReader.ReadNext()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read record: %w", err)
	}

	return record, nil
}

// writeParquetFile writes an Arrow record to a Parquet file
func writeParquetFile(path string, record arrow.Record) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer, err := pqarrow.NewFileWriter(record.Schema(), file, parquet.NewWriterProperties(), pqarrow.DefaultWriterProps())
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %w", err)
	}
	defer writer.Close()

	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	return nil
} 