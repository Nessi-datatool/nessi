package delta

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	parquetfile "github.com/apache/arrow/go/v15/parquet/file"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
)

// DeltaConnector represents a connection to a Delta table
type DeltaConnector struct {
	TablePath        string
	Config           *Config
	Schema           *arrow.Schema
	Version          int64
	PartitionColumns []string
	MergeOptions     *MergeOptions
	basePath         string
	lastUpdated      time.Time
	partitions       []string
	metadata         *TableMetadata
	log              *TransactionLog
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
		TablePath: tablePath,
		basePath:  filepath.Dir(tablePath),
	}
	connector.log = NewTransactionLog(connector.TablePath)

	// Load table metadata
	if err := connector.loadMetadata(); err != nil {
		return nil, fmt.Errorf("failed to load metadata: %w", err)
	}

	return connector, nil
}

// loadMetadata loads the table metadata from the Delta log
func (c *DeltaConnector) loadMetadata() error {
	logPath := filepath.Join(c.TablePath, "_delta_log")
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
	filePath := filepath.Join(c.TablePath, "_delta_log", fmt.Sprintf("%020d.json", version)) // Padded version number
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	// The log file can contain multiple actions, one per line, or a single JSON array of actions.
	// For simplicity, we'll assume a single JSON array for now, as produced by CommitTransaction.
	// More robust parsing would handle line-delimited JSON objects.
	var actions []Action
	if !strings.HasPrefix(string(data), "[") {
		// Attempt to parse as line-delimited JSON if not starting with an array bracket
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			var singleAction Action
			if err := json.Unmarshal(scanner.Bytes(), &singleAction); err != nil {
				// Potentially log this error or handle incomplete lines
				continue // Skip malformed lines
			}
			actions = append(actions, singleAction)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("failed to scan log file lines: %w", err)
		}
	} else {
		// Parse as a single JSON array
		if err := json.Unmarshal(data, &actions); err != nil {
			return fmt.Errorf("failed to parse log file as JSON array: %w", err)
		}
	}

	// Process each action
	for _, action := range actions {
		if action.Meta != nil {
			// Create a basic schema if SchemaString is provided
			if action.Meta.SchemaString != "" {
				// Instead of trying to unmarshal directly into arrow.Schema,
				// parse the schema fields from the SchemaString and create a schema
				var schemaObj map[string]interface{}
				if err := json.Unmarshal([]byte(action.Meta.SchemaString), &schemaObj); err != nil {
					return fmt.Errorf("failed to parse schema string: %w", err)
				}

				// Extract schema fields from the JSON representation
				fieldsArr, ok := schemaObj["fields"].([]interface{})
				if ok && len(fieldsArr) > 0 {
					fields := make([]arrow.Field, 0, len(fieldsArr))
					for _, fieldObj := range fieldsArr {
						fieldMap, ok := fieldObj.(map[string]interface{})
						if !ok {
							continue
						}

						name, _ := fieldMap["name"].(string)
						typeStr, _ := fieldMap["type"].(string)
						nullable, _ := fieldMap["nullable"].(bool)

						// Map string type representation to arrow.DataType
						var dataType arrow.DataType
						switch typeStr {
						case "int", "integer", "int32":
							dataType = arrow.PrimitiveTypes.Int32
						case "bigint", "int64":
							dataType = arrow.PrimitiveTypes.Int64
						case "float", "float32":
							dataType = arrow.PrimitiveTypes.Float32
						case "double", "float64":
							dataType = arrow.PrimitiveTypes.Float64
						case "boolean", "bool":
							dataType = arrow.FixedWidthTypes.Boolean
						case "string", "varchar", "char":
							dataType = arrow.BinaryTypes.String
						case "binary":
							dataType = arrow.BinaryTypes.Binary
						case "timestamp", "timestamp_ms":
							dataType = arrow.FixedWidthTypes.Timestamp_ms
						default:
							// Default to string for unknown types
							dataType = arrow.BinaryTypes.String
						}

						fields = append(fields, arrow.Field{
							Name:     name,
							Type:     dataType,
							Nullable: nullable,
						})
					}

					// Create a new schema with the extracted fields
					if len(fields) > 0 {
						c.Schema = arrow.NewSchema(fields, nil)
					}
				}
			}

			c.metadata = &TableMetadata{
				ID:          action.Meta.ID,
				Name:        action.Meta.Name,
				Description: action.Meta.Description,
				Schema:      c.Schema, // Now uses the potentially parsed schema
				Partitions:  action.Meta.PartitionColumns,
				Properties:  action.Meta.Configuration,
				CreatedAt:   time.Unix(0, action.Meta.CreatedTime*int64(time.Millisecond)),
				UpdatedAt:   time.Now(), // Should ideally be from metadata if available
				Version:     version,
			}
		}
		// TODO: Handle other action types like Add, Remove, etc. to update table state (e.g., file list)
	}

	c.Version = version
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
	partitionPath := filepath.Join(c.TablePath, partition)
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
	partitionPath := filepath.Join(c.TablePath, partition)
	if err := os.MkdirAll(partitionPath, 0755); err != nil {
		return fmt.Errorf("failed to create partition directory: %w", err)
	}

	// Generate a unique filename
	filename := fmt.Sprintf("part-%d.parquet", time.Now().UnixNano())
	filePath := filepath.Join(partitionPath, filename)

	// Write the record to Parquet
	if err := WriteRecordToParquet(filePath, record); err != nil {
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
	// Create an AddAction
	addAction := &AddAction{
		Path:             filepath.Join(partition, filename),
		PartitionValues:  parsePartitionValues(partition),
		Size:             getFileSize(filepath.Join(c.TablePath, partition, filename)),
		ModificationTime: time.Now().UnixMilli(),
		DataChange:       true,
		Stats:            fmt.Sprintf(`{"numRecords":%d}`, numRows), // Consider making this a struct if further processing is needed
	}

	// Create an Action containing the AddAction
	action := Action{
		Add: addAction,
	}

	// Delta logs expect an array of actions, even if it's just one
	actionsToWrite := []Action{action}

	// Write the transaction to the log
	logPath := filepath.Join(c.TablePath, "_delta_log", fmt.Sprintf("%020d.json", c.Version+1)) // Padded version number
	data, err := json.Marshal(actionsToWrite)                                                   // Marshal the slice of actions
	if err != nil {
		return fmt.Errorf("failed to marshal transaction actions: %w", err)
	}

	if err := os.WriteFile(logPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	c.Version++
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

	// Create builders for each column
	builders := make([]array.Builder, len(schema.Fields()))
	for i := range builders {
		builders[i] = array.NewBuilder(memory.DefaultAllocator, schema.Field(i).Type)
		defer builders[i].Release()
	}

	// Append data from all records
	for _, record := range records {
		for i := 0; i < int(record.NumCols()); i++ {
			col := record.Column(i)
			for j := 0; j < col.Len(); j++ {
				if col.IsNull(j) {
					builders[i].AppendNull()
					continue
				}

				switch builder := builders[i].(type) {
				case *array.StringBuilder:
					str := col.(*array.String).Value(j)
					builder.Append(str)
				case *array.Int64Builder:
					val := col.(*array.Int64).Value(j)
					builder.Append(val)
				case *array.Float64Builder:
					val := col.(*array.Float64).Value(j)
					builder.Append(val)
				case *array.BooleanBuilder:
					val := col.(*array.Boolean).Value(j)
					builder.Append(val)
				case *array.TimestampBuilder:
					val := col.(*array.Timestamp).Value(j)
					builder.Append(val)
				default:
					return nil, fmt.Errorf("unsupported type: %s", builder.Type())
				}
			}
		}
	}

	// Create arrays from builders
	arrays := make([]arrow.Array, len(builders))
	for i, builder := range builders {
		arrays[i] = builder.NewArray()
		defer arrays[i].Release()
	}

	// Create the merged record
	return array.NewRecord(schema, arrays, -1), nil
}

// GetSchema returns the current table schema
func (c *DeltaConnector) GetSchema() *arrow.Schema {
	return c.Schema
}

// GetPartitions returns the current table partitions
func (c *DeltaConnector) GetPartitions() []string {
	return c.partitions
}

// GetVersion returns the current table version
func (c *DeltaConnector) GetVersion() int64 {
	return c.Version
}

// GetLastUpdated returns the last update timestamp
func (c *DeltaConnector) GetLastUpdated() time.Time {
	return c.lastUpdated
}

// CommitTransaction commits the given transaction to the Delta log.
func (c *DeltaConnector) CommitTransaction(ctx context.Context, tx *Transaction) error {
	if c.log == nil {
		return fmt.Errorf("TransactionLog not initialized in DeltaConnector")
	}
	return c.log.Commit(ctx, tx)
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

	reader, err := parquetfile.NewParquetReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer reader.Close()

	arrowReader, err := pqarrow.NewFileReader(reader, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	if err != nil {
		return nil, fmt.Errorf("failed to create arrow reader: %w", err)
	}

	table, err := arrowReader.ReadTable(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to read table: %w", err)
	}

	// Convert table to record
	builder := array.NewRecordBuilder(memory.DefaultAllocator, table.Schema())
	for i := 0; i < int(table.NumCols()); i++ {
		col := table.Column(i)
		builderField := builder.Field(i)

		// Handle different data types
		switch builder := builderField.(type) {
		case *array.StringBuilder:
			if strCol, ok := col.Data().Chunk(0).(*array.String); ok {
				for j := 0; j < strCol.Len(); j++ {
					if strCol.IsValid(j) {
						builder.Append(strCol.Value(j))
					} else {
						builder.AppendNull()
					}
				}
			} else {
				return nil, fmt.Errorf("failed to convert column to string array")
			}
		case *array.Int64Builder:
			if intCol, ok := col.Data().Chunk(0).(*array.Int64); ok {
				for j := 0; j < intCol.Len(); j++ {
					if intCol.IsValid(j) {
						builder.Append(intCol.Value(j))
					} else {
						builder.AppendNull()
					}
				}
			} else {
				return nil, fmt.Errorf("failed to convert column to int64 array")
			}
		case *array.Float64Builder:
			if floatCol, ok := col.Data().Chunk(0).(*array.Float64); ok {
				for j := 0; j < floatCol.Len(); j++ {
					if floatCol.IsValid(j) {
						builder.Append(floatCol.Value(j))
					} else {
						builder.AppendNull()
					}
				}
			} else {
				return nil, fmt.Errorf("failed to convert column to float64 array")
			}
		case *array.BooleanBuilder:
			if boolCol, ok := col.Data().Chunk(0).(*array.Boolean); ok {
				for j := 0; j < boolCol.Len(); j++ {
					if boolCol.IsValid(j) {
						builder.Append(boolCol.Value(j))
					} else {
						builder.AppendNull()
					}
				}
			} else {
				return nil, fmt.Errorf("failed to convert column to boolean array")
			}
		case *array.TimestampBuilder:
			if tsCol, ok := col.Data().Chunk(0).(*array.Timestamp); ok {
				for j := 0; j < tsCol.Len(); j++ {
					if tsCol.IsValid(j) {
						builder.Append(tsCol.Value(j))
					} else {
						builder.AppendNull()
					}
				}
			} else {
				return nil, fmt.Errorf("failed to convert column to timestamp array")
			}
		default:
			return nil, fmt.Errorf("unsupported column type: %T", builder)
		}
	}
	record := builder.NewRecord()
	table.Release()
	builder.Release()
	if err != nil {
		return nil, fmt.Errorf("failed to read record: %w", err)
	}

	return record, nil
}

// Initialize initializes the Delta table connector
func (c *DeltaConnector) Initialize() error {
	// Create _delta_log directory if it doesn't exist
	logPath := filepath.Join(c.TablePath, "_delta_log")
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return fmt.Errorf("failed to create Delta log directory: %w", err)
	}

	// Load metadata
	if err := c.loadMetadata(); err != nil {
		return fmt.Errorf("failed to load metadata: %w", err)
	}

	return nil
}
