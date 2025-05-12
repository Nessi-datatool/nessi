package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsistencyChecks(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "delta-consistency-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create _delta_log directory
	logDir := filepath.Join(tempDir, "_delta_log")
	err = os.MkdirAll(logDir, 0755)
	require.NoError(t, err)

	// Create test versions with schema changes
	now := time.Now()

	// Version 0: Initial schema with 3 fields
	createTestVersionWithSchema(t, logDir, 0, "CREATE_TABLE", now.Add(-3*time.Hour), []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		{"name": "value", "type": "float64"},
	})

	// Version 1: Added a field
	createTestVersionWithSchema(t, logDir, 1, "ALTER_TABLE", now.Add(-2*time.Hour), []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		{"name": "value", "type": "float64"},
		{"name": "timestamp", "type": "timestamp"},
	})

	// Version 2: Changed a field type
	createTestVersionWithSchema(t, logDir, 2, "ALTER_TABLE", now.Add(-1*time.Hour), []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		{"name": "value", "type": "int32"}, // Changed from float64 to int32
		{"name": "timestamp", "type": "timestamp"},
	})

	// Version 3: Removed a field
	createTestVersionWithSchema(t, logDir, 3, "ALTER_TABLE", now, []map[string]string{
		{"name": "id", "type": "int32"},
		{"name": "name", "type": "utf8"},
		// "value" field removed
		{"name": "timestamp", "type": "timestamp"},
		{"name": "is_active", "type": "boolean"}, // Added new field
	})

	// Create metadata manager
	manager := NewMetadataManager(tempDir)

	t.Run("SchemaConsistencyCheck", func(t *testing.T) {
		// Check schema consistency across all versions
		result, err := manager.CheckConsistency(SchemaConsistencyCheck, []int64{0, 1, 2, 3})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should fail because fields were removed
		assert.False(t, result.Passed)

		// Check issues
		var addedFields, removedFields int
		for _, issue := range result.Issues {
			if issue.Type == "field_added" {
				addedFields++
			} else if issue.Type == "field_removed" {
				removedFields++
				assert.Equal(t, "warning", issue.Severity)
			}
		}

		// Should have 2 added fields (timestamp, is_active) and 1 removed field (value)
		assert.Equal(t, 2, addedFields)
		assert.Equal(t, 1, removedFields)
	})

	t.Run("TypeConsistencyCheck", func(t *testing.T) {
		// Check type consistency across all versions
		result, err := manager.CheckConsistency(TypeConsistencyCheck, []int64{0, 1, 2, 3})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should fail because a field type was changed
		assert.False(t, result.Passed)

		// Check issues
		var typeChanges int
		for _, issue := range result.Issues {
			if issue.Type == "type_changed" {
				typeChanges++
				assert.Equal(t, "error", issue.Severity)
				assert.Equal(t, "value", issue.Field)
			}
		}

		// Should have 1 type change (value: float64 -> int32)
		assert.Equal(t, 1, typeChanges)
	})

	t.Run("FieldConsistencyCheck", func(t *testing.T) {
		// Check field consistency across all versions
		result, err := manager.CheckConsistency(FieldConsistencyCheck, []int64{0, 1, 2, 3})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should fail because latest version has fields not in earlier versions
		assert.False(t, result.Passed)

		// Check issues
		var missingFields int
		for _, issue := range result.Issues {
			if issue.Type == "field_missing" {
				missingFields++
				assert.Equal(t, "warning", issue.Severity)
			}
		}

		// Should have 2 missing fields (is_active missing in versions 0, 1, 2)
		assert.Equal(t, 2, missingFields)
	})

	t.Run("InvalidCheckType", func(t *testing.T) {
		// Try with an invalid check type
		_, err := manager.CheckConsistency("invalid", []int64{0, 1})
		assert.Error(t, err)
	})

	t.Run("NotEnoughVersions", func(t *testing.T) {
		// Try with only one version
		_, err := manager.CheckConsistency(SchemaConsistencyCheck, []int64{0})
		assert.Error(t, err)
	})

	t.Run("NonExistentVersion", func(t *testing.T) {
		// Try with a non-existent version
		_, err := manager.CheckConsistency(SchemaConsistencyCheck, []int64{0, 999})
		assert.Error(t, err)
	})
}

// Helper function to create a test version with a specific schema
func createTestVersionWithSchema(t *testing.T, logDir string, version int64, operation string, timestamp time.Time, fields []map[string]string) {
	// Create schema
	schema := map[string]interface{}{
		"fields": fields,
	}

	// Create files list
	files := make([]string, version+1)
	for i := int64(0); i <= version; i++ {
		files[i] = fmt.Sprintf("part-%05d.parquet", i)
	}

	// Create metadata
	metadata := map[string]interface{}{
		"version":    version,
		"timestamp":  timestamp.Unix(),
		"schema":     schema,
		"files":      files,
		"partitions": map[string][]string{},
		"metadata":   map[string]interface{}{},
	}

	// Convert to JSON
	data, err := json.MarshalIndent(metadata, "", "  ")
	require.NoError(t, err)

	// Write to transaction log
	logFile := filepath.Join(logDir, fmt.Sprintf("%020d.json", version))
	err = os.WriteFile(logFile, data, 0644)
	require.NoError(t, err)

	// Create commit info
	commitInfo := map[string]interface{}{
		"timestamp":           timestamp.UnixMilli(),
		"operation":           operation,
		"operationParameters": map[string]string{},
		"isBlindAppend":       operation == "APPEND",
		"isolationLevel":      "Serializable",
	}

	// Convert to JSON
	commitData, err := json.MarshalIndent(commitInfo, "", "  ")
	require.NoError(t, err)

	// Write commit info
	commitFile := filepath.Join(logDir, fmt.Sprintf("%020d.commit.json", version))
	err = os.WriteFile(commitFile, commitData, 0644)
	require.NoError(t, err)
}
