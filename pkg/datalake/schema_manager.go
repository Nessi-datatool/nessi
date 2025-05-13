package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
)

// SchemaManager handles Delta Lake schema operations
type SchemaManager struct {
	tablePath string
	metaDir   string
}

// SchemaHistory represents the history of schema changes
type SchemaHistory struct {
	Schemas      []*SchemaVersion `json:"schemas"`
	CurrentIndex int              `json:"current_index"`
}

// SchemaVersion represents a version of a schema
type SchemaVersion struct {
	ID          string            `json:"id"`
	Version     int               `json:"version"`
	Timestamp   time.Time         `json:"timestamp"`
	Schema      *Schema           `json:"schema"`
	Changes     []SchemaChange    `json:"changes"`
	CommitInfo  map[string]string `json:"commit_info"`
	PartitionBy []string          `json:"partition_by,omitempty"`
	ZOrderBy    []string          `json:"z_order_by,omitempty"`
}

// SchemaChange represents a change to a schema
type SchemaChange struct {
	Type        string `json:"type"`        // "add", "remove", "modify", "rename"
	Field       string `json:"field"`       // Field name
	OldType     string `json:"old_type,omitempty"`
	NewType     string `json:"new_type,omitempty"`
	OldName     string `json:"old_name,omitempty"`
	NewName     string `json:"new_name,omitempty"`
	Description string `json:"description"` // Human-readable description
}

// FieldMetadata represents metadata for a field
type FieldMetadata struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Nullable    bool                   `json:"nullable"`
	Description string                 `json:"description,omitempty"`
	Stats       map[string]interface{} `json:"stats,omitempty"`
	Tags        map[string]string      `json:"tags,omitempty"`
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(tablePath string) *SchemaManager {
	return &SchemaManager{
		tablePath: tablePath,
		metaDir:   filepath.Join(tablePath, "_delta_log", "schema_history"),
	}
}

// InitializeSchema initializes a schema for a table
func (sm *SchemaManager) InitializeSchema(schema *Schema, partitionBy, zOrderBy []string) (*SchemaVersion, error) {
	// Create schema history directory if it doesn't exist
	if err := os.MkdirAll(sm.metaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create schema history directory: %w", err)
	}

	// Check if schema history already exists
	historyPath := filepath.Join(sm.metaDir, "history.json")
	if _, err := os.Stat(historyPath); err == nil {
		return nil, fmt.Errorf("schema already initialized")
	}

	// Create initial schema version
	schemaVersion := &SchemaVersion{
		ID:          uuid.New().String(),
		Version:     1,
		Timestamp:   time.Now(),
		Schema:      schema,
		Changes:     []SchemaChange{},
		CommitInfo:  map[string]string{"action": "initialize"},
		PartitionBy: partitionBy,
		ZOrderBy:    zOrderBy,
	}

	// Create schema history
	history := &SchemaHistory{
		Schemas:      []*SchemaVersion{schemaVersion},
		CurrentIndex: 0,
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

	if len(history.Schemas) == 0 {
		return nil, fmt.Errorf("no schema found")
	}

	return history.Schemas[history.CurrentIndex], nil
}

// GetSchemaVersion gets a specific schema version
func (sm *SchemaManager) GetSchemaVersion(version int) (*SchemaVersion, error) {
	history, err := sm.readHistory()
	if err != nil {
		return nil, err
	}

	for _, schema := range history.Schemas {
		if schema.Version == version {
			return schema, nil
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

	if len(history.Schemas) == 0 {
		return nil, fmt.Errorf("no schema found, use InitializeSchema instead")
	}

	currentSchema := history.Schemas[history.CurrentIndex]
	
	// Detect changes
	changes, err := sm.detectSchemaChanges(currentSchema.Schema, newSchema)
	if err != nil {
		return nil, err
	}

	if len(changes) == 0 {
		return currentSchema, nil // No changes
	}

	// Create new schema version
	newVersion := &SchemaVersion{
		ID:          uuid.New().String(),
		Version:     currentSchema.Version + 1,
		Timestamp:   time.Now(),
		Schema:      newSchema,
		Changes:     changes,
		CommitInfo:  commitInfo,
		PartitionBy: currentSchema.PartitionBy, // Preserve partition info
		ZOrderBy:    currentSchema.ZOrderBy,    // Preserve Z-order info
	}

	// Add new schema to history
	history.Schemas = append(history.Schemas, newVersion)
	history.CurrentIndex = len(history.Schemas) - 1

	// Write schema history to file
	if err := sm.writeHistory(history); err != nil {
		return nil, err
	}

	return newVersion, nil
}

// UpdatePartitioning updates the partitioning scheme
func (sm *SchemaManager) UpdatePartitioning(partitionBy, zOrderBy []string, commitInfo map[string]string) (*SchemaVersion, error) {
	history, err := sm.readHistory()
	if err != nil {
		return nil, err
	}

	if len(history.Schemas) == 0 {
		return nil, fmt.Errorf("no schema found")
	}

	currentSchema := history.Schemas[history.CurrentIndex]
	
	// Check if there are any changes
	partitionChanged := !equalStringSlices(currentSchema.PartitionBy, partitionBy)
	zOrderChanged := !equalStringSlices(currentSchema.ZOrderBy, zOrderBy)
	
	if !partitionChanged && !zOrderChanged {
		return currentSchema, nil // No changes
	}

	// Create new schema version with updated partitioning
	newVersion := &SchemaVersion{
		ID:          uuid.New().String(),
		Version:     currentSchema.Version + 1,
		Timestamp:   time.Now(),
		Schema:      currentSchema.Schema, // Same schema
		Changes:     []SchemaChange{},     // No schema changes
		CommitInfo:  commitInfo,
		PartitionBy: partitionBy,
		ZOrderBy:    zOrderBy,
	}

	// Add description of partitioning changes
	if partitionChanged {
		newVersion.Changes = append(newVersion.Changes, SchemaChange{
			Type:        "partition",
			Description: fmt.Sprintf("Updated partitioning from %v to %v", currentSchema.PartitionBy, partitionBy),
		})
	}
	
	if zOrderChanged {
		newVersion.Changes = append(newVersion.Changes, SchemaChange{
			Type:        "z_order",
			Description: fmt.Sprintf("Updated Z-ordering from %v to %v", currentSchema.ZOrderBy, zOrderBy),
		})
	}

	// Add new schema to history
	history.Schemas = append(history.Schemas, newVersion)
	history.CurrentIndex = len(history.Schemas) - 1

	// Write schema history to file
	if err := sm.writeHistory(history); err != nil {
		return nil, err
	}

	return newVersion, nil
}

// ValidateSchema validates that data conforms to the schema
func (sm *SchemaManager) ValidateSchema(data map[string]interface{}) (bool, []string, error) {
	currentSchema, err := sm.GetCurrentSchema()
	if err != nil {
		return false, nil, err
	}

	schema := currentSchema.Schema
	violations := []string{}

	// Check each field in the schema
	for _, field := range schema.Fields {
		value, exists := data[field.Name]
		
		// Check for required fields
		if !exists && !field.Nullable {
			violations = append(violations, fmt.Sprintf("Required field '%s' is missing", field.Name))
			continue
		}
		
		// Skip null values for nullable fields
		if !exists || value == nil {
			if !field.Nullable {
				violations = append(violations, fmt.Sprintf("Required field '%s' is null", field.Name))
			}
			continue
		}
		
		// Validate type
		valid := sm.validateType(value, field.Type)
		if !valid {
			violations = append(violations, fmt.Sprintf("Field '%s' has invalid type, expected %s", field.Name, field.Type))
		}
	}

	// Check for extra fields not in schema
	for key := range data {
		found := false
		for _, field := range schema.Fields {
			if field.Name == key {
				found = true
				break
			}
		}
		
		if !found {
			violations = append(violations, fmt.Sprintf("Unknown field '%s' not in schema", key))
		}
	}

	return len(violations) == 0, violations, nil
}

// GetFieldMetadata gets metadata for a field
func (sm *SchemaManager) GetFieldMetadata(fieldName string) (*FieldMetadata, error) {
	currentSchema, err := sm.GetCurrentSchema()
	if err != nil {
		return nil, err
	}

	// Find the field in the schema
	for _, field := range currentSchema.Schema.Fields {
		if field.Name == fieldName {
			// Create field metadata
			metadata := &FieldMetadata{
				Name:     field.Name,
				Type:     field.Type,
				Nullable: field.Nullable,
				Stats:    make(map[string]interface{}),
				Tags:     make(map[string]string),
			}
			
			// TODO: Load additional metadata from Delta Lake statistics
			
			return metadata, nil
		}
	}

	return nil, fmt.Errorf("field '%s' not found in schema", fieldName)
}

// GetOptimizationHints gets optimization hints for the schema
func (sm *SchemaManager) GetOptimizationHints() (map[string]string, error) {
	currentSchema, err := sm.GetCurrentSchema()
	if err != nil {
		return nil, err
	}

	hints := make(map[string]string)
	
	// Check partitioning
	if len(currentSchema.PartitionBy) == 0 {
		hints["partitioning"] = "Consider adding partitioning to improve query performance"
	} else if len(currentSchema.PartitionBy) > 2 {
		hints["partitioning"] = "Too many partition columns may lead to small files, consider reducing"
	}
	
	// Check Z-ordering
	if len(currentSchema.ZOrderBy) == 0 && len(currentSchema.Schema.Fields) > 5 {
		hints["z_ordering"] = "Consider adding Z-ordering for frequently queried columns"
	}
	
	// Check field types
	for _, field := range currentSchema.Schema.Fields {
		if field.Type == "string" && !field.Nullable {
			hints[fmt.Sprintf("field_%s", field.Name)] = "Consider using a more specific type than string if possible"
		}
	}
	
	return hints, nil
}

// Helper functions

// detectSchemaChanges detects changes between two schemas
func (sm *SchemaManager) detectSchemaChanges(oldSchema, newSchema *Schema) ([]SchemaChange, error) {
	changes := []SchemaChange{}
	
	// Create maps for quick lookup
	oldFields := make(map[string]Field)
	for _, field := range oldSchema.Fields {
		oldFields[field.Name] = field
	}
	
	newFields := make(map[string]Field)
	for _, field := range newSchema.Fields {
		newFields[field.Name] = field
	}
	
	// Find added and modified fields
	for name, newField := range newFields {
		oldField, exists := oldFields[name]
		
		if !exists {
			// New field
			changes = append(changes, SchemaChange{
				Type:        "add",
				Field:       name,
				NewType:     newField.Type,
				Description: fmt.Sprintf("Added field '%s' of type '%s'", name, newField.Type),
			})
		} else if oldField.Type != newField.Type || oldField.Nullable != newField.Nullable {
			// Modified field
			changes = append(changes, SchemaChange{
				Type:        "modify",
				Field:       name,
				OldType:     oldField.Type,
				NewType:     newField.Type,
				Description: fmt.Sprintf("Modified field '%s' from '%s' to '%s'", name, oldField.Type, newField.Type),
			})
		}
	}
	
	// Find removed fields
	for name := range oldFields {
		if _, exists := newFields[name]; !exists {
			changes = append(changes, SchemaChange{
				Type:        "remove",
				Field:       name,
				OldType:     oldFields[name].Type,
				Description: fmt.Sprintf("Removed field '%s'", name),
			})
		}
	}
	
	// Sort changes for consistent output
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Field < changes[j].Field
	})
	
	return changes, nil
}

// validateType validates that a value matches the expected type
func (sm *SchemaManager) validateType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "integer":
		switch v := value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			return true
		case float64:
			return v == float64(int(v)) // Check if it's a whole number
		default:
			return false
		}
	case "float", "double":
		_, ok := value.(float64)
		return ok
	case "boolean":
		_, ok := value.(bool)
		return ok
	case "date", "timestamp":
		// For simplicity, we'll accept strings that might be dates
		// In a real implementation, we would parse and validate the date format
		_, ok := value.(string)
		return ok
	default:
		// For complex types like arrays, structs, etc.
		// In a real implementation, we would have more sophisticated validation
		return true
	}
}

// readHistory reads the schema history from file
func (sm *SchemaManager) readHistory() (*SchemaHistory, error) {
	historyPath := filepath.Join(sm.metaDir, "history.json")
	
	// Check if file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("schema history not found")
	}
	
	// Read file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema history: %w", err)
	}
	
	// Parse JSON
	var history SchemaHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse schema history: %w", err)
	}
	
	return &history, nil
}

// writeHistory writes the schema history to file
func (sm *SchemaManager) writeHistory(history *SchemaHistory) error {
	historyPath := filepath.Join(sm.metaDir, "history.json")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(sm.metaDir, 0755); err != nil {
		return fmt.Errorf("failed to create schema history directory: %w", err)
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema history: %w", err)
	}
	
	// Write file
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write schema history: %w", err)
	}
	
	return nil
}

// equalStringSlices checks if two string slices are equal
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	
	return true
}
