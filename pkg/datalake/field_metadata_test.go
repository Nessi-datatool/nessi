package datalake

import (
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldMetadata(t *testing.T) {
	// Create a schema for testing
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
			{Name: "timestamp", Type: arrow.FixedWidthTypes.Timestamp_s},
			{Name: "is_active", Type: arrow.FixedWidthTypes.Boolean},
		},
		nil,
	)

	t.Run("Add and get field metadata", func(t *testing.T) {
		// Add metadata to a field
		description := "Customer identifier"
		tags := []string{"primary_key", "indexed"}
		properties := map[string]string{
			"min_value": "1",
			"max_value": "1000000",
		}

		updatedSchema, err := AddFieldMetadata(schema, "id", description, tags, properties)
		require.NoError(t, err)
		require.NotNil(t, updatedSchema)

		// Get metadata for the field
		metadata, err := GetFieldMetadata(updatedSchema, "id")
		require.NoError(t, err)
		require.NotNil(t, metadata)

		// Verify metadata
		assert.Equal(t, "id", metadata.Name)
		assert.Equal(t, "int32", metadata.Type)
		assert.Equal(t, description, metadata.Description)
		assert.Equal(t, tags, metadata.Tags)
		assert.Equal(t, properties, metadata.Properties)

		// Format metadata
		formatted := FormatFieldMetadata(metadata)
		assert.Contains(t, formatted, "Field: id (int32)")
		assert.Contains(t, formatted, "Description: Customer identifier")
		assert.Contains(t, formatted, "Tags: [primary_key indexed]")
		assert.Contains(t, formatted, "Properties:")
		assert.Contains(t, formatted, "min_value: 1")
		assert.Contains(t, formatted, "max_value: 1000000")
	})

	t.Run("Get all fields metadata", func(t *testing.T) {
		// Add metadata to multiple fields
		updatedSchema, err := AddFieldMetadata(schema, "id", "Customer identifier", []string{"primary_key"}, nil)
		require.NoError(t, err)

		updatedSchema, err = AddFieldMetadata(updatedSchema, "name", "Customer name", []string{"indexed"}, nil)
		require.NoError(t, err)

		// Get all fields metadata
		allMetadata, err := GetAllFieldsMetadata(updatedSchema)
		require.NoError(t, err)
		require.NotNil(t, allMetadata)

		// Verify metadata
		require.Len(t, allMetadata.Fields, 5) // We now have 5 fields in the schema
		assert.Equal(t, "id", allMetadata.Fields[0].Name)
		assert.Equal(t, "Customer identifier", allMetadata.Fields[0].Description)
		assert.Equal(t, "name", allMetadata.Fields[1].Name)
		assert.Equal(t, "Customer name", allMetadata.Fields[1].Description)

		// Format all metadata
		formatted := FormatAllFieldsMetadata(allMetadata)
		assert.Contains(t, formatted, "Schema Fields:")
		assert.Contains(t, formatted, "Field: id (int32)")
		assert.Contains(t, formatted, "Field: name (string)")
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test nil schema
		_, err := AddFieldMetadata(nil, "id", "description", nil, nil)
		assert.Error(t, err)

		// Test non-existent field
		_, err = AddFieldMetadata(schema, "non_existent", "description", nil, nil)
		assert.Error(t, err)

		// Test get metadata for nil schema
		_, err = GetFieldMetadata(nil, "id")
		assert.Error(t, err)

		// Test get metadata for non-existent field
		_, err = GetFieldMetadata(schema, "non_existent")
		assert.Error(t, err)

		// Test get all metadata for nil schema
		_, err = GetAllFieldsMetadata(nil)
		assert.Error(t, err)
	})

	t.Run("Update existing field metadata", func(t *testing.T) {
		// First add metadata
		updatedSchema, err := AddFieldMetadata(schema, "id", "Initial description", []string{"tag1"}, map[string]string{"prop1": "val1"})
		require.NoError(t, err)
		
		// Now update the metadata
		updatedSchema, err = AddFieldMetadata(updatedSchema, "id", "Updated description", []string{"tag1", "tag2"}, map[string]string{"prop1": "updated", "prop2": "val2"})
		require.NoError(t, err)

		// Verify the updated metadata
		metadata, err := GetFieldMetadata(updatedSchema, "id")
		require.NoError(t, err)
		assert.Equal(t, "Updated description", metadata.Description)
		assert.Equal(t, []string{"tag1", "tag2"}, metadata.Tags)
		assert.Equal(t, "updated", metadata.Properties["prop1"])
		assert.Equal(t, "val2", metadata.Properties["prop2"])
	})

	t.Run("Add metadata to different field types", func(t *testing.T) {
		// Test adding metadata to timestamp field
		updatedSchema, err := AddFieldMetadata(schema, "timestamp", "Event timestamp", []string{"temporal"}, map[string]string{"format": "yyyy-MM-dd HH:mm:ss"})
		require.NoError(t, err)

		// Verify timestamp field metadata
		metadata, err := GetFieldMetadata(updatedSchema, "timestamp")
		require.NoError(t, err)
		assert.Equal(t, "timestamp", metadata.Name)
		assert.Contains(t, metadata.Type, "timestamp")
		assert.Equal(t, "Event timestamp", metadata.Description)
		assert.Equal(t, []string{"temporal"}, metadata.Tags)
		assert.Equal(t, "yyyy-MM-dd HH:mm:ss", metadata.Properties["format"])

		// Test adding metadata to boolean field
		updatedSchema, err = AddFieldMetadata(updatedSchema, "is_active", "Account status", []string{"status"}, map[string]string{"true_label": "Active", "false_label": "Inactive"})
		require.NoError(t, err)

		// Verify boolean field metadata
		metadata, err = GetFieldMetadata(updatedSchema, "is_active")
		require.NoError(t, err)
		assert.Equal(t, "is_active", metadata.Name)
		assert.Equal(t, "boolean", metadata.Type)
		assert.Equal(t, "Account status", metadata.Description)
		assert.Equal(t, []string{"status"}, metadata.Tags)
		assert.Equal(t, "Active", metadata.Properties["true_label"])
		assert.Equal(t, "Inactive", metadata.Properties["false_label"])
	})

	t.Run("Schema metadata handling", func(t *testing.T) {
		// Add field metadata
		updatedSchema, err := AddFieldMetadata(schema, "id", "ID field", nil, nil)
		require.NoError(t, err)

		// Verify we can get the field metadata back
		metadata, err := GetFieldMetadata(updatedSchema, "id")
		require.NoError(t, err)
		assert.Equal(t, "ID field", metadata.Description)
	})

	t.Run("Handle empty values", func(t *testing.T) {
		// Test with empty description, tags, and properties
		updatedSchema, err := AddFieldMetadata(schema, "id", "", []string{}, map[string]string{})
		require.NoError(t, err)

		// Verify metadata
		metadata, err := GetFieldMetadata(updatedSchema, "id")
		require.NoError(t, err)
		assert.Equal(t, "", metadata.Description)
		assert.Empty(t, metadata.Tags)
		assert.Empty(t, metadata.Properties)
	})
}
