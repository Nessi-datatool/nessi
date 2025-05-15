package datalake

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_manager_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a schema manager
	sm := NewDeltaSchemaManager(tempDir)

	// Test initializing a schema
	t.Run("InitializeSchema", func(t *testing.T) {
		// Create a simple schema
		schema := &DeltaSchema{
			Fields: []DeltaField{
				{Name: "id", Type: "INTEGER", Nullable: false},
				{Name: "name", Type: "STRING", Nullable: true},
				{Name: "created_at", Type: "TIMESTAMP", Nullable: false},
			},
		}

		// Initialize the schema
		_, err := sm.InitializeSchema(schema, []string{}, []string{})
		require.NoError(t, err)

		// Get the schema version
		schemaVersion, err := sm.GetCurrentSchema()
		require.NoError(t, err)

		// Check version
		if schemaVersion.Version != 1 {
			t.Errorf("Expected version 1, got %d", schemaVersion.Version)
		}

		// Check schema fields
		if len(schemaVersion.Schema.Fields()) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(schemaVersion.Schema.Fields()))
		}
	})

	// Test getting current schema
	t.Run("GetCurrentSchema", func(t *testing.T) {
		schemaVersion, err := sm.GetCurrentSchema()
		require.NoError(t, err)
		assert.Equal(t, int64(1), schemaVersion.Version)
	})

	// Test getting schema history
	t.Run("GetSchemaHistory", func(t *testing.T) {
		history, err := sm.GetSchemaHistory()
		require.NoError(t, err)
		assert.Equal(t, 1, len(history.Versions))
	})

	// Test updating schema
	t.Run("UpdateSchema", func(t *testing.T) {
		// Create a new schema with an additional field
		updatedSchema := &DeltaSchema{
			Fields: []DeltaField{
				{Name: "id", Type: "INTEGER", Nullable: false},
				{Name: "name", Type: "STRING", Nullable: true},
				{Name: "created_at", Type: "TIMESTAMP", Nullable: false},
				{Name: "active", Type: "BOOLEAN", Nullable: true},
			},
		}

		// Update the schema
		_, err = sm.UpdateSchema(updatedSchema, map[string]string{})
		require.NoError(t, err)

		// Get the updated schema version
		updatedVersion, err := sm.GetCurrentSchema()
		require.NoError(t, err)
		assert.Equal(t, int64(2), updatedVersion.Version)

		// Check schema fields
		if len(updatedVersion.Schema.Fields()) != 4 {
			t.Errorf("Expected 4 fields, got %d", len(updatedVersion.Schema.Fields()))
		}
	})

	// Test getting schema at version
	t.Run("GetSchemaAtVersion", func(t *testing.T) {
		schemaVersion, err := sm.GetSchemaVersion(1)
		require.NoError(t, err)
		assert.Equal(t, int64(1), schemaVersion.Version)
		assert.Equal(t, 3, len(schemaVersion.Schema.Fields()))
	})
}
