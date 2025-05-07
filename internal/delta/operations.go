package delta

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
)

// TimeTravelOptions defines options for time travel operations
type TimeTravelOptions struct {
	Timestamp time.Time
	Version   int64
}

// SchemaEvolutionOptions defines options for schema evolution
type SchemaEvolutionOptions struct {
	AddColumns    []arrow.Field
	DropColumns   []string
	UpdateColumns []arrow.Field
}

// DeltaOperations provides high-level operations for Delta Lake tables
type DeltaOperations struct {
	connector *DeltaConnector
}

// NewOperations creates a new DeltaOperations instance
func NewOperations(connector *DeltaConnector) *DeltaOperations {
	return &DeltaOperations{
		connector: connector,
	}
}

// TimeTravel performs time travel to a specific version or timestamp
func (o *DeltaOperations) TimeTravel(ctx context.Context, options TimeTravelOptions) error {
	if options.Timestamp.IsZero() && options.Version == 0 {
		return fmt.Errorf("either timestamp or version must be specified")
	}

	// Get the target version
	targetVersion := options.Version
	if !options.Timestamp.IsZero() {
		var err error
		targetVersion, err = o.findVersionByTimestamp(options.Timestamp)
		if err != nil {
			return fmt.Errorf("failed to find version by timestamp: %w", err)
		}
	}

	// Validate the target version
	if targetVersion < 0 || targetVersion > o.connector.GetVersion() {
		return fmt.Errorf("invalid version: %d", targetVersion)
	}

	// Create a checkpoint at the target version
	if err := o.createCheckpoint(targetVersion); err != nil {
		return fmt.Errorf("failed to create checkpoint: %w", err)
	}

	return nil
}

// findVersionByTimestamp finds the version closest to the given timestamp
func (o *DeltaOperations) findVersionByTimestamp(timestamp time.Time) (int64, error) {
	logPath := filepath.Join(o.connector.tablePath, "_delta_log")
	files, err := os.ReadDir(logPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read log directory: %w", err)
	}

	var closestVersion int64
	var closestTime time.Time
	var found bool

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			version, err := parseVersionFromFilename(file.Name())
			if err != nil {
				continue
			}

			filePath := filepath.Join(logPath, file.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			fileTime := info.ModTime()
			if !found || fileTime.Sub(timestamp).Abs() < closestTime.Sub(timestamp).Abs() {
				closestVersion = version
				closestTime = fileTime
				found = true
			}
		}
	}

	if !found {
		return 0, fmt.Errorf("no version found for timestamp: %v", timestamp)
	}

	return closestVersion, nil
}

// createCheckpoint creates a checkpoint at the specified version
func (o *DeltaOperations) createCheckpoint(version int64) error {
	checkpointPath := filepath.Join(o.connector.tablePath, "_delta_log", fmt.Sprintf("%d.checkpoint.parquet", version))
	
	// Read all actions up to the target version
	var actions []Action
	for v := int64(0); v <= version; v++ {
		filePath := filepath.Join(o.connector.tablePath, "_delta_log", fmt.Sprintf("%d.json", v))
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read log file: %w", err)
		}

		var versionActions []Action
		if err := json.Unmarshal(data, &versionActions); err != nil {
			return fmt.Errorf("failed to parse log file: %w", err)
		}
		actions = append(actions, versionActions...)
	}

	// Create a Parquet file with the checkpoint data
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "version", Type: arrow.PrimitiveTypes.Int64},
		{Name: "timestamp", Type: arrow.PrimitiveTypes.Int64},
		{Name: "action", Type: arrow.BinaryTypes.String},
	}, nil)

	builder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builder.Release()

	for _, action := range actions {
		builder.Field(0).(*array.Int64Builder).Append(version)
		builder.Field(1).(*array.Int64Builder).Append(time.Now().UnixMilli())
		actionJSON, _ := json.Marshal(action)
		builder.Field(2).(*array.StringBuilder).Append(string(actionJSON))
	}

	record := builder.NewRecord()
	defer record.Release()

	if err := writeParquetFile(checkpointPath, record); err != nil {
		return fmt.Errorf("failed to write checkpoint file: %w", err)
	}

	return nil
}

// EvolveSchema performs schema evolution on the Delta table
func (o *DeltaOperations) EvolveSchema(ctx context.Context, options SchemaEvolutionOptions) error {
	// Get current schema
	currentSchema := o.connector.GetSchema()
	if currentSchema == nil {
		return fmt.Errorf("no schema available")
	}

	// Create new schema
	newFields := make([]arrow.Field, 0, currentSchema.NumFields())
	fieldMap := make(map[string]arrow.Field)

	// Add existing fields
	for i := 0; i < currentSchema.NumFields(); i++ {
		field := currentSchema.Field(i)
		fieldMap[field.Name] = field
		newFields = append(newFields, field)
	}

	// Add new columns
	for _, field := range options.AddColumns {
		if _, exists := fieldMap[field.Name]; exists {
			return fmt.Errorf("column already exists: %s", field.Name)
		}
		newFields = append(newFields, field)
		fieldMap[field.Name] = field
	}

	// Update existing columns
	for _, field := range options.UpdateColumns {
		if _, exists := fieldMap[field.Name]; !exists {
			return fmt.Errorf("column does not exist: %s", field.Name)
		}
		for i, f := range newFields {
			if f.Name == field.Name {
				newFields[i] = field
				fieldMap[field.Name] = field
				break
			}
		}
	}

	// Remove dropped columns
	if len(options.DropColumns) > 0 {
		filteredFields := make([]arrow.Field, 0, len(newFields))
		for _, field := range newFields {
			keep := true
			for _, dropName := range options.DropColumns {
				if field.Name == dropName {
					keep = false
					break
				}
			}
			if keep {
				filteredFields = append(filteredFields, field)
			}
		}
		newFields = filteredFields
	}

	// Create new schema
	newSchema := arrow.NewSchema(newFields, nil)

	// Update metadata
	metadata, err := o.connector.GetMetadata()
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	metadata.Schema = newSchema
	metadata.UpdatedAt = time.Now()
	metadata.Version++

	// Write schema evolution transaction
	transaction := struct {
		Meta struct {
			Schema *arrow.Schema `json:"schema"`
		} `json:"meta"`
	}{
		Meta: struct {
			Schema *arrow.Schema `json:"schema"`
		}{
			Schema: newSchema,
		},
	}

	// Write the transaction to the log
	logPath := filepath.Join(o.connector.tablePath, "_delta_log", fmt.Sprintf("%d.json", metadata.Version))
	data, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	if err := os.WriteFile(logPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	return nil
}

// parseVersionFromFilename extracts the version number from a Delta log filename
func parseVersionFromFilename(filename string) (int64, error) {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	if ext != ".json" {
		return 0, fmt.Errorf("not a JSON file")
	}
	versionStr := base[:len(base)-len(ext)]
	return strconv.ParseInt(versionStr, 10, 64)
} 