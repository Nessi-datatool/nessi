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

func TestTimeTravelSimple(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "delta-timetravel-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create _delta_log directory
	logDir := filepath.Join(tempDir, "_delta_log")
	err = os.MkdirAll(logDir, 0755)
	require.NoError(t, err)

	// Create test versions with different timestamps
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	// Create version 0 (2 hours ago)
	createTestVersionWithTime(t, logDir, 0, "CREATE_TABLE", twoHoursAgo, []string{"id", "name", "value"})

	// Create version 1 (1 hour ago)
	createTestVersionWithTime(t, logDir, 1, "ALTER_TABLE", oneHourAgo, []string{"id", "name", "value", "timestamp"})

	// Create version 2 (now)
	createTestVersionWithTime(t, logDir, 2, "ALTER_TABLE", now, []string{"id", "name", "value", "timestamp", "is_active"})

	// Create metadata manager
	manager := NewMetadataManager(tempDir)

	t.Run("TimeTravel by Version", func(t *testing.T) {
		// Query version 0
		version := int64(0)
		result, err := manager.TimeTravel(TimeTravelOptions{Version: &version})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Check result
		assert.Equal(t, version, result.Version)
		assert.Equal(t, twoHoursAgo.Unix(), result.Timestamp.Unix())
		assert.Equal(t, 3, len(result.SchemaFields))
		assert.Equal(t, "id", result.SchemaFields[0].Name)
		assert.Equal(t, "name", result.SchemaFields[1].Name)
		assert.Equal(t, "value", result.SchemaFields[2].Name)

		// Query version 1
		version = 1
		result, err = manager.TimeTravel(TimeTravelOptions{Version: &version})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Check result
		assert.Equal(t, version, result.Version)
		assert.Equal(t, oneHourAgo.Unix(), result.Timestamp.Unix())
		assert.Equal(t, 4, len(result.SchemaFields))
		assert.Equal(t, "timestamp", result.SchemaFields[3].Name)

		// Query version 2
		version = 2
		result, err = manager.TimeTravel(TimeTravelOptions{Version: &version})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Check result
		assert.Equal(t, version, result.Version)
		assert.Equal(t, now.Unix(), result.Timestamp.Unix())
		assert.Equal(t, 5, len(result.SchemaFields))
		assert.Equal(t, "is_active", result.SchemaFields[4].Name)
	})

	t.Run("TimeTravel by Timestamp", func(t *testing.T) {
		// Query at a time between version 0 and version 1
		timestamp := twoHoursAgo.Add(30 * time.Minute)
		result, err := manager.TimeTravel(TimeTravelOptions{Timestamp: &timestamp})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return version 0
		assert.Equal(t, int64(0), result.Version)
		assert.Equal(t, 3, len(result.SchemaFields))

		// Query at a time between version 1 and version 2
		timestamp = oneHourAgo.Add(30 * time.Minute)
		result, err = manager.TimeTravel(TimeTravelOptions{Timestamp: &timestamp})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return version 1
		assert.Equal(t, int64(1), result.Version)
		assert.Equal(t, 4, len(result.SchemaFields))

		// Query at a time after version 2
		timestamp = now.Add(30 * time.Minute)
		result, err = manager.TimeTravel(TimeTravelOptions{Timestamp: &timestamp})
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return version 2
		assert.Equal(t, int64(2), result.Version)
		assert.Equal(t, 5, len(result.SchemaFields))
	})

	t.Run("TimeTravel with Invalid Options", func(t *testing.T) {
		// Query with no options
		_, err := manager.TimeTravel(TimeTravelOptions{})
		assert.Error(t, err)

		// Query with non-existent version
		version := int64(999)
		_, err = manager.TimeTravel(TimeTravelOptions{Version: &version})
		assert.Error(t, err)
	})

	t.Run("GetSchemaFieldsAtVersion", func(t *testing.T) {
		// Get schema fields at version 1
		schemaFields, err := manager.GetSchemaFieldsAtVersion(1)
		require.NoError(t, err)
		require.NotNil(t, schemaFields)

		// Check schema fields
		assert.Equal(t, 4, len(schemaFields))
		assert.Equal(t, "id", schemaFields[0].Name)
		assert.Equal(t, "name", schemaFields[1].Name)
		assert.Equal(t, "value", schemaFields[2].Name)
		assert.Equal(t, "timestamp", schemaFields[3].Name)
	})

	t.Run("GetSchemaFieldsAtTimestamp", func(t *testing.T) {
		// Get schema fields at a time between version 1 and version 2
		timestamp := oneHourAgo.Add(30 * time.Minute)
		schemaFields, err := manager.GetSchemaFieldsAtTimestamp(timestamp)
		require.NoError(t, err)
		require.NotNil(t, schemaFields)

		// Check schema fields
		assert.Equal(t, 4, len(schemaFields))
		assert.Equal(t, "timestamp", schemaFields[3].Name)
	})

	t.Run("GetFilesAtVersion", func(t *testing.T) {
		// Get files at version 0
		files, err := manager.GetFilesAtVersion(0)
		require.NoError(t, err)
		require.NotNil(t, files)

		// Check files
		assert.Equal(t, 1, len(files))
		assert.Equal(t, "part-00000.parquet", files[0])

		// Get files at version 2
		files, err = manager.GetFilesAtVersion(2)
		require.NoError(t, err)
		require.NotNil(t, files)

		// Check files
		assert.Equal(t, 3, len(files))
		assert.Equal(t, "part-00000.parquet", files[0])
		assert.Equal(t, "part-00001.parquet", files[1])
		assert.Equal(t, "part-00002.parquet", files[2])
	})
}

// Helper function to create a test version with a specific timestamp
func createTestVersionWithTime(t *testing.T, logDir string, version int64, operation string, timestamp time.Time, fieldNames []string) {
	// Create schema fields
	fields := make([]map[string]interface{}, len(fieldNames))
	for i, name := range fieldNames {
		var typeStr string
		switch name {
		case "id":
			typeStr = "int32"
		case "name":
			typeStr = "utf8"
		case "value":
			typeStr = "float64"
		case "timestamp":
			typeStr = "timestamp"
		case "is_active":
			typeStr = "boolean"
		default:
			typeStr = "string"
		}

		fields[i] = map[string]interface{}{
			"name": name,
			"type": typeStr,
		}
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
		"schema":     map[string]interface{}{"fields": fields},
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
		"timestamp":    timestamp.UnixMilli(),
		"operation":    operation,
		"operationParameters": map[string]string{},
		"isBlindAppend": operation == "APPEND",
		"isolationLevel": "Serializable",
	}

	// Convert to JSON
	commitData, err := json.MarshalIndent(commitInfo, "", "  ")
	require.NoError(t, err)

	// Write commit info
	commitFile := filepath.Join(logDir, fmt.Sprintf("%020d.commit.json", version))
	err = os.WriteFile(commitFile, commitData, 0644)
	require.NoError(t, err)
}
