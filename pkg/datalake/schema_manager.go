package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
)

// SchemaManager handles Delta Lake schema operations
type SchemaManager struct {
	tablePath string
	metaDir   string
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(tablePath string) *SchemaManager {
	return &SchemaManager{
		tablePath: tablePath,
		metaDir:   filepath.Join(tablePath, "_delta_log", "schema_history"),
	}
}

// InitializeSchema initializes a schema for a table
func (sm *SchemaManager) InitializeSchema(schema *Schema, partitionBy []string, zOrderBy []string) (*SchemaVersion, error) {
	// Check if schema already exists
	if _, err := sm.GetCurrentSchema(); err == nil {
		return nil, fmt.Errorf("schema already initialized")
	}

	// Convert Schema to arrow.Schema
	arrowFields := make([]arrow.Field, len(schema.Fields))
	for i, field := range schema.Fields {
		// Simple conversion - in a real implementation, this would be more sophisticated
		var dataType arrow.DataType
		switch field.Type {
		case "integer":
			dataType = arrow.PrimitiveTypes.Int64
		case "float":
			dataType = arrow.PrimitiveTypes.Float64
		case "boolean":
			dataType = arrow.FixedWidthTypes.Boolean
		case "timestamp":
			dataType = arrow.FixedWidthTypes.Timestamp_s
		default:
			dataType = arrow.BinaryTypes.String
		}
		arrowFields[i] = arrow.Field{Name: field.Name, Type: dataType, Nullable: field.Nullable}
	}
	arrowSchema := arrow.NewSchema(arrowFields, nil)

	// Create initial schema version
	schemaVersion := &SchemaVersion{
		Version:      1,
		Timestamp:    time.Now(),
		Schema:       arrowSchema,
		SchemaFields: arrowFields,
	}

	// Create schema history
	history := &SchemaHistory{
		Versions: []SchemaVersion{*schemaVersion},
	}

	// Write schema history to file
	if err := sm.writeHistory(history); err != nil {
		return nil, err
	}

	return schemaVersion, nil
}

// GetCurrentSchema gets the current schema
func (sm *SchemaManager) GetCurrentSchema() (*SchemaVersion, error) {
	history, err := sm.readHistory()
	if err != nil {
		return nil, err
	}

	if len(history.Versions) == 0 {
		return nil, fmt.Errorf("no schema found")
	}

	// Return the latest schema version (last in the list)
	return &history.Versions[len(history.Versions)-1], nil
}

// GetSchemaVersion gets a specific schema version
func (sm *SchemaManager) GetSchemaVersion(version int) (*SchemaVersion, error) {
	history, err := sm.readHistory()
	if err != nil {
		return nil, err
	}

	for _, schema := range history.Versions {
		if int(schema.Version) == version {
			return &schema, nil
		}
	}

	return nil, fmt.Errorf("schema version %d not found", version)
}

// GetSchemaHistory gets the schema history
func (sm *SchemaManager) GetSchemaHistory() (*SchemaHistory, error) {
	return sm.readHistory()
}

// UpdateSchema updates the schema
func (sm *SchemaManager) UpdateSchema(newSchema *Schema, commitInfo map[string]string) (*SchemaVersion, error) {
	history, err := sm.readHistory()
	if err != nil {
		return nil, err
	}

	if len(history.Versions) == 0 {
		return nil, fmt.Errorf("no schema found, use InitializeSchema instead")
	}

	// Get the latest schema version
	currentSchema := history.Versions[len(history.Versions)-1]
	
	// Convert Schema to arrow.Schema
	arrowFields := make([]arrow.Field, len(newSchema.Fields))
	for i, field := range newSchema.Fields {
		// Simple conversion - in a real implementation, this would be more sophisticated
		var dataType arrow.DataType
		switch field.Type {
		case "integer":
			dataType = arrow.PrimitiveTypes.Int64
		case "float":
			dataType = arrow.PrimitiveTypes.Float64
		case "boolean":
			dataType = arrow.FixedWidthTypes.Boolean
		case "timestamp":
			dataType = arrow.FixedWidthTypes.Timestamp_s
		default:
			dataType = arrow.BinaryTypes.String
		}
		arrowFields[i] = arrow.Field{Name: field.Name, Type: dataType, Nullable: field.Nullable}
	}
	arrowSchema := arrow.NewSchema(arrowFields, nil)

	// For simplicity, we'll skip the actual change detection in this implementation
	// In a real implementation, we would compare the schemas and detect changes

	// Create new schema version
	newVersion := &SchemaVersion{
		Version:      currentSchema.Version + 1,
		Timestamp:    time.Now(),
		Schema:       arrowSchema,
		SchemaFields: arrowFields,
	}

	// Add new schema to history
	history.Versions = append(history.Versions, *newVersion)

	// Write schema history to file
	if err := sm.writeHistory(history); err != nil {
		return nil, err
	}

	return newVersion, nil
}

// Helper functions

// readHistory reads the schema history from file
func (sm *SchemaManager) readHistory() (*SchemaHistory, error) {
	// Create schema history directory if it doesn't exist
	if err := os.MkdirAll(sm.metaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create schema history directory: %w", err)
	}

	// Check if schema history file exists
	historyPath := filepath.Join(sm.metaDir, "history.json")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		// Return empty history if file doesn't exist
		return &SchemaHistory{
			Versions: []SchemaVersion{},
		}, nil
	}

	// Read schema history file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema history: %w", err)
	}

	// Parse schema history
	var history SchemaHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse schema history: %w", err)
	}

	return &history, nil
}

// writeHistory writes the schema history to file
func (sm *SchemaManager) writeHistory(history *SchemaHistory) error {
	// Create schema history directory if it doesn't exist
	if err := os.MkdirAll(sm.metaDir, 0755); err != nil {
		return fmt.Errorf("failed to create schema history directory: %w", err)
	}

	// Marshal schema history
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema history: %w", err)
	}

	// Write schema history file
	historyPath := filepath.Join(sm.metaDir, "history.json")
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema history: %w", err)
	}

	return nil
}
