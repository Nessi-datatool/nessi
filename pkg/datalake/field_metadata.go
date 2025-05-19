package datalake

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/nessi-dev/nessi/pkg/logging"
)

// FieldMetadata represents metadata for a field in a schema
type FieldMetadata struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Tags        []string          `json:"tags"`
	Properties  map[string]string `json:"properties"`
}

// SchemaWithMetadata represents a schema with field-level metadata
type SchemaWithMetadata struct {
	Fields []FieldMetadata `json:"fields"`
}

// AddFieldMetadata adds metadata to a field in a schema
func AddFieldMetadata(schema *arrow.Schema, fieldName string, description string, tags []string, properties map[string]string) (*arrow.Schema, error) {
	logger := logging.GetLogger()
	logger.Debug("Adding field metadata", "field", fieldName)

	if schema == nil {
		logger.Error("Schema cannot be nil")
		return nil, fmt.Errorf("schema cannot be nil")
	}

	if fieldName == "" {
		logger.Error("Field name cannot be empty")
		return nil, fmt.Errorf("field name cannot be empty")
	}

	// Find field
	fieldIdx := schema.FieldIndices(fieldName)
	if len(fieldIdx) == 0 {
		logger.Error("Field not found", "field", fieldName)
		return nil, fmt.Errorf("field not found: %s", fieldName)
	}

	// Get field
	field := schema.Field(fieldIdx[0])

	// Pre-allocate capacity for metadata map
	metadataCapacity := 1 // For description
	if len(tags) > 0 {
		metadataCapacity++
	}
	if properties != nil {
		metadataCapacity += len(properties)
	}

	// Create metadata map
	metadata := make(map[string]string, metadataCapacity)
	if description != "" {
		metadata["description"] = description
	}

	// Add tags
	if len(tags) > 0 {
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			logger.Error("Failed to marshal tags", "error", err, "tags", tags)
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		metadata["tags"] = string(tagsJSON)
	}

	// Add properties
	if properties != nil {
		for k, v := range properties {
			if k != "" && k != "description" && k != "tags" {
				metadata[k] = v
			}
		}
	}

	// Create new field with metadata
	metadataKeys := make([]string, 0, len(metadata))
	metadataValues := make([]string, 0, len(metadata))
	for k, v := range metadata {
		metadataKeys = append(metadataKeys, k)
		metadataValues = append(metadataValues, v)
	}

	// Create new field with metadata
	newField := arrow.Field{
		Name:     field.Name,
		Type:     field.Type,
		Nullable: field.Nullable,
		Metadata: arrow.NewMetadata(metadataKeys, metadataValues),
	}

	// Create new schema with updated field
	fields := make([]arrow.Field, schema.NumFields())
	for i := 0; i < schema.NumFields(); i++ {
		if i == fieldIdx[0] {
			fields[i] = newField
		} else {
			fields[i] = schema.Field(i)
		}
	}

	// Create new schema with the same metadata as the original
	schemaMetadata := schema.Metadata()
	newSchema := arrow.NewSchema(fields, &schemaMetadata)
	logger.Debug("Successfully added field metadata", "field", fieldName)
	return newSchema, nil
}

// GetFieldMetadata gets metadata for a field in a schema
func GetFieldMetadata(schema *arrow.Schema, fieldName string) (*FieldMetadata, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting field metadata", "field", fieldName)

	if schema == nil {
		logger.Error("Schema cannot be nil")
		return nil, fmt.Errorf("schema cannot be nil")
	}

	if fieldName == "" {
		logger.Error("Field name cannot be empty")
		return nil, fmt.Errorf("field name cannot be empty")
	}

	// Find field
	fieldIdx := schema.FieldIndices(fieldName)
	if len(fieldIdx) == 0 {
		logger.Error("Field not found", "field", fieldName)
		return nil, fmt.Errorf("field not found: %s", fieldName)
	}

	// Get field
	field := schema.Field(fieldIdx[0])

	// Estimate properties capacity
	propertiesCapacity := field.Metadata.Len()
	if propertiesCapacity > 2 { // Subtract space for description and tags
		propertiesCapacity -= 2
	}

	// Create metadata
	metadata := &FieldMetadata{
		Name:        field.Name,
		Type:        formatArrowType(field.Type),
		Description: "",
		Tags:        []string{},
		Properties:  make(map[string]string, propertiesCapacity),
	}

	// Get metadata fields
	keys := field.Metadata.Keys()
	values := field.Metadata.Values()

	// Process metadata
	for i := 0; i < field.Metadata.Len(); i++ {
		key := keys[i]
		value := values[i]

		switch key {
		case "description":
			metadata.Description = value
		case "tags":
			// Parse tags JSON
			var tags []string
			if err := json.Unmarshal([]byte(value), &tags); err == nil {
				metadata.Tags = tags
			} else {
				logger.Warn("Failed to parse tags JSON", "error", err, "value", value)
			}
		default:
			// Other properties
			metadata.Properties[key] = value
		}
	}

	logger.Debug("Successfully retrieved field metadata", "field", fieldName)
	return metadata, nil
}

// GetAllFieldsMetadata gets metadata for all fields in a schema
func GetAllFieldsMetadata(schema *arrow.Schema) (*SchemaWithMetadata, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting metadata for all fields")

	if schema == nil {
		logger.Error("Schema cannot be nil")
		return nil, fmt.Errorf("schema cannot be nil")
	}

	numFields := schema.NumFields()
	if numFields == 0 {
		logger.Debug("Schema has no fields")
		return &SchemaWithMetadata{Fields: []FieldMetadata{}}, nil
	}

	result := &SchemaWithMetadata{
		Fields: make([]FieldMetadata, numFields),
	}

	// Process all fields
	for i := 0; i < numFields; i++ {
		field := schema.Field(i)
		metadata, err := GetFieldMetadata(schema, field.Name)
		if err != nil {
			logger.Error("Failed to get metadata for field", "field", field.Name, "error", err)
			return nil, fmt.Errorf("failed to get metadata for field %s: %w", field.Name, err)
		}
		result.Fields[i] = *metadata
	}

	logger.Debug("Successfully retrieved metadata for all fields", "count", numFields)
	return result, nil
}

// FormatFieldMetadata returns a human-readable representation of field metadata
func FormatFieldMetadata(metadata *FieldMetadata) string {
	if metadata == nil {
		return "<nil>"
	}

	// Use strings.Builder for more efficient string concatenation
	var sb strings.Builder

	// Pre-allocate a reasonable buffer size to avoid reallocations
	sb.Grow(256)

	// Write field name and type
	sb.WriteString(fmt.Sprintf("Field: %s (%s)\n", metadata.Name, metadata.Type))

	// Write description if present
	if metadata.Description != "" {
		sb.WriteString(fmt.Sprintf("Description: %s\n", metadata.Description))
	}

	// Write tags if present
	if len(metadata.Tags) > 0 {
		sb.WriteString(fmt.Sprintf("Tags: %v\n", metadata.Tags))
	}

	// Write properties if present
	if len(metadata.Properties) > 0 {
		sb.WriteString("Properties:\n")

		// Sort keys for consistent output
		keys := make([]string, 0, len(metadata.Properties))
		for k := range metadata.Properties {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s: %s\n", k, metadata.Properties[k]))
		}
	}

	return sb.String()
}

// FormatAllFieldsMetadata returns a human-readable representation of all field metadata
func FormatAllFieldsMetadata(metadata *SchemaWithMetadata) string {
	if metadata == nil {
		return "<nil>"
	}

	// Use strings.Builder for more efficient string concatenation
	var sb strings.Builder

	// Pre-allocate a reasonable buffer size based on the number of fields
	// Assume average of 200 bytes per field
	sb.Grow(200*len(metadata.Fields) + 20)

	sb.WriteString("Schema Fields:\n")

	for i, field := range metadata.Fields {
		// Add separator between fields (except for the first one)
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("---\n")
		sb.WriteString(FormatFieldMetadata(&field))
	}

	return sb.String()
}

// Note: formatArrowType is already defined in schema.go
