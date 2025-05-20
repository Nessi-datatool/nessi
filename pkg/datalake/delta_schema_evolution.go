package datalake

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DeltaSchemaEvolution provides schema evolution tracking for Delta Lake tables
type DeltaSchemaEvolution struct {
	handler *DeltaFormatHandler
}

// NewDeltaSchemaEvolution creates a new DeltaSchemaEvolution instance
func NewDeltaSchemaEvolution(handler *DeltaFormatHandler) *DeltaSchemaEvolution {
	return &DeltaSchemaEvolution{
		handler: handler,
	}
}

// SchemaChangeInfo represents a change to a Delta Lake table schema
type SchemaChangeInfo struct {
	Version      int64
	Timestamp    time.Time
	Operation    string
	ChangedCols  []string
	AddedCols    []string
	RemovedCols  []string
	ModifiedCols map[string]ColumnModification
}

// ColumnModification represents a modification to a column
type ColumnModification struct {
	OldType     string
	NewType     string
	OldNullable bool
	NewNullable bool
}

// GetSchemaHistory returns the schema evolution history of a Delta Lake table
func (s *DeltaSchemaEvolution) GetSchemaHistory(ctx context.Context, path string) ([]SchemaChangeInfo, error) {
	// Check if the path is a Delta Lake table
	if !s.handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Get the transaction log files
	deltaLogDir := filepath.Join(path, "_delta_log")
	logFiles, err := filepath.Glob(filepath.Join(deltaLogDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list transaction log files: %w", err)
	}

	// Extract version information from log files
	versions := make(map[int64]string)
	for _, file := range logFiles {
		baseName := filepath.Base(file)
		if strings.HasSuffix(baseName, ".json") && !strings.Contains(baseName, "checkpoint") {
			versionStr := strings.TrimSuffix(baseName, ".json")
			v, err := strconv.ParseInt(versionStr, 10, 64)
			if err != nil {
				continue
			}
			versions[v] = file
		}
	}

	// Sort version numbers
	versionNumbers := make([]int64, 0, len(versions))
	for v := range versions {
		versionNumbers = append(versionNumbers, v)
	}
	sort.Slice(versionNumbers, func(i, j int) bool {
		return versionNumbers[i] < versionNumbers[j]
	})

	// Track schema changes
	schemaChanges := make([]SchemaChangeInfo, 0)
	prevSchema := make(map[string]interface{})

	for _, version := range versionNumbers {
		file := versions[version]

		// Read the log file
		logContent, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}

		// Parse the log file (simplified for this implementation)
		// In a real implementation, we would parse the transaction log format properly
		var logEntry struct {
			Schema    map[string]interface{} `json:"schema"`
			Operation string                 `json:"operation"`
			Timestamp int64                  `json:"timestamp"`
		}

		if err := json.Unmarshal(logContent, &logEntry); err != nil {
			// Skip invalid log entries
			continue
		}

		// If no schema in this log entry, skip
		if logEntry.Schema == nil {
			continue
		}

		// Compare with previous schema to detect changes
		change := SchemaChangeInfo{
			Version:      version,
			Timestamp:    time.Unix(0, logEntry.Timestamp*1000000),
			Operation:    logEntry.Operation,
			ChangedCols:  []string{},
			AddedCols:    []string{},
			RemovedCols:  []string{},
			ModifiedCols: make(map[string]ColumnModification),
		}

		// In a real implementation, we would compare the schemas properly
		// For now, we'll just record that there was a schema change
		if len(prevSchema) > 0 {
			// Compare schemas to detect changes
			change.ChangedCols = []string{"schema_changed"}
		}

		schemaChanges = append(schemaChanges, change)
		prevSchema = logEntry.Schema
	}

	return schemaChanges, nil
}

// GetSchemaAtVersion returns the schema of a Delta Lake table at a specific version
func (s *DeltaSchemaEvolution) GetSchemaAtVersion(ctx context.Context, path string, version int64) (*Schema, error) {
	// Check if the path is a Delta Lake table
	if !s.handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Get the transaction log file for the specified version
	deltaLogDir := filepath.Join(path, "_delta_log")
	logFile := filepath.Join(deltaLogDir, fmt.Sprintf("%020d.json", version))

	// Read the log file
	logContent, err := ioutil.ReadFile(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction log for version %d: %w", version, err)
	}

	// Parse the log file (simplified for this implementation)
	var logEntry struct {
		Schema struct {
			Fields []struct {
				Name     string `json:"name"`
				Type     string `json:"type"`
				Nullable bool   `json:"nullable"`
			} `json:"fields"`
		} `json:"schema"`
	}

	if err := json.Unmarshal(logContent, &logEntry); err != nil {
		return nil, fmt.Errorf("failed to parse transaction log for version %d: %w", version, err)
	}

	// Create inferred schema
	fields := make([]Field, 0, len(logEntry.Schema.Fields))

	for _, field := range logEntry.Schema.Fields {
		fields = append(fields, Field{
			Name: field.Name,
			Type: convertDeltaTypeToFieldType(field.Type),
		})
	}

	// Create schema with fields
	schema := NewSchema(fields)

	return schema, nil
}

// DetectSchemaDrift detects schema drift between two versions of a Delta Lake table
func (s *DeltaSchemaEvolution) DetectSchemaDrift(ctx context.Context, path string, version1, version2 int64) (*SchemaDrift, error) {
	// Get schemas at both versions
	schema1, err := s.GetSchemaAtVersion(ctx, path, version1)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema at version %d: %w", version1, err)
	}

	schema2, err := s.GetSchemaAtVersion(ctx, path, version2)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema at version %d: %w", version2, err)
	}

	// Compare schemas
	return CompareSchemas(schema1, schema2), nil
}

// SchemaDrift represents the differences between two schemas
type SchemaDrift struct {
	AddedFields    []Field
	RemovedFields  []Field
	ModifiedFields map[string]FieldModification
}

// FieldModification represents a modification to a field
type FieldModification struct {
	OldType FieldType
	NewType FieldType
}

// CompareSchemas compares two schemas and returns the differences
func CompareSchemas(schema1, schema2 *Schema) *SchemaDrift {
	drift := &SchemaDrift{
		AddedFields:    []Field{},
		RemovedFields:  []Field{},
		ModifiedFields: make(map[string]FieldModification),
	}

	// Create maps for easier comparison
	fields1 := make(map[string]Field)
	for _, field := range schema1.fields {
		fields1[field.Name] = field
	}

	fields2 := make(map[string]Field)
	for _, field := range schema2.fields {
		fields2[field.Name] = field
	}

	// Find added fields
	for name, field := range fields2 {
		if _, exists := fields1[name]; !exists {
			drift.AddedFields = append(drift.AddedFields, field)
		}
	}

	// Find removed fields
	for name, field := range fields1 {
		if _, exists := fields2[name]; !exists {
			drift.RemovedFields = append(drift.RemovedFields, field)
		}
	}

	// Find modified fields
	for name, field1 := range fields1 {
		if field2, exists := fields2[name]; exists {
			if field1.Type != field2.Type {
				drift.ModifiedFields[name] = FieldModification{
					OldType: field1.Type,
					NewType: field2.Type,
				}
			}
		}
	}

	return drift
}

// Helper function to convert Delta Lake type strings to FieldType
func convertDeltaTypeToFieldType(deltaType string) FieldType {
	switch strings.ToLower(deltaType) {
	case "string":
		return FieldTypeString
	case "integer", "int":
		return FieldTypeInt32
	case "long":
		return FieldTypeInt64
	case "float":
		return FieldTypeFloat32
	case "double":
		return FieldTypeFloat64
	case "boolean":
		return FieldTypeBool
	case "binary":
		return FieldTypeBinary
	case "timestamp":
		return FieldTypeTimestamp
	case "date":
		return FieldTypeDate
	default:
		return FieldTypeString
	}
}
