package datalake

import (
	"encoding/json"
	"fmt"

	"github.com/apache/arrow/go/v15/arrow"
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
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	// Find field
	fieldIdx := schema.FieldIndices(fieldName)
	if len(fieldIdx) == 0 {
		return nil, fmt.Errorf("field not found: %s", fieldName)
	}

	// Get field
	field := schema.Field(fieldIdx[0])

	// Create metadata map
	metadata := map[string]string{
		"description": description,
	}

	// Add tags
	if len(tags) > 0 {
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		metadata["tags"] = string(tagsJSON)
	}

	// Add properties
	for k, v := range properties {
		metadata[k] = v
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
	return arrow.NewSchema(fields, nil), nil
}

// GetFieldMetadata gets metadata for a field in a schema
func GetFieldMetadata(schema *arrow.Schema, fieldName string) (*FieldMetadata, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	// Find field
	fieldIdx := schema.FieldIndices(fieldName)
	if len(fieldIdx) == 0 {
		return nil, fmt.Errorf("field not found: %s", fieldName)
	}

	// Get field
	field := schema.Field(fieldIdx[0])

	// Create metadata
	metadata := &FieldMetadata{
		Name:        field.Name,
		Type:        formatArrowType(field.Type),
		Description: "",
		Tags:        []string{},
		Properties:  make(map[string]string),
	}

	// Get description
	for i := 0; i < field.Metadata.Len(); i++ {
		key := field.Metadata.Keys()[i]
		value := field.Metadata.Values()[i]
		
		if key == "description" {
			metadata.Description = value
		} else if key == "tags" {
			// Parse tags JSON
			var tags []string
			if err := json.Unmarshal([]byte(value), &tags); err == nil {
				metadata.Tags = tags
			}
		} else {
			// Other properties
			metadata.Properties[key] = value
		}
	}

	return metadata, nil
}

// GetAllFieldsMetadata gets metadata for all fields in a schema
func GetAllFieldsMetadata(schema *arrow.Schema) (*SchemaWithMetadata, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	result := &SchemaWithMetadata{
		Fields: make([]FieldMetadata, schema.NumFields()),
	}

	for i := 0; i < schema.NumFields(); i++ {
		field := schema.Field(i)
		metadata, err := GetFieldMetadata(schema, field.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata for field %s: %w", field.Name, err)
		}
		result.Fields[i] = *metadata
	}

	return result, nil
}

// FormatFieldMetadata returns a human-readable representation of field metadata
func FormatFieldMetadata(metadata *FieldMetadata) string {
	result := fmt.Sprintf("Field: %s (%s)\n", metadata.Name, metadata.Type)
	if metadata.Description != "" {
		result += fmt.Sprintf("Description: %s\n", metadata.Description)
	}
	if len(metadata.Tags) > 0 {
		result += fmt.Sprintf("Tags: %v\n", metadata.Tags)
	}
	if len(metadata.Properties) > 0 {
		result += "Properties:\n"
		for k, v := range metadata.Properties {
			result += fmt.Sprintf("  %s: %s\n", k, v)
		}
	}
	return result
}

// FormatAllFieldsMetadata returns a human-readable representation of all field metadata
func FormatAllFieldsMetadata(metadata *SchemaWithMetadata) string {
	result := "Schema Fields:\n"
	for _, field := range metadata.Fields {
		result += "---\n"
		result += FormatFieldMetadata(&field)
	}
	return result
}

// Note: formatArrowType is already defined in schema.go
