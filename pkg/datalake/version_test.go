package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionHistory(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "delta-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create _delta_log directory
	logDir := filepath.Join(tempDir, "_delta_log")
	err = os.MkdirAll(logDir, 0755)
	require.NoError(t, err)

	// Create test versions
	createTestVersion(t, logDir, 0, "CREATE_TABLE", time.Now().Add(-48*time.Hour))
	createTestVersion(t, logDir, 1, "APPEND", time.Now().Add(-24*time.Hour))
	createTestVersion(t, logDir, 2, "UPDATE", time.Now())

	// Create metadata manager
	manager := NewMetadataManager(tempDir)

	t.Run("GetVersionHistory", func(t *testing.T) {
		// Get version history
		history, err := manager.GetVersionHistory()
		require.NoError(t, err)
		require.Len(t, history, 3)

		// Check versions are in descending order
		assert.Equal(t, int64(2), history[0].Version)
		assert.Equal(t, int64(1), history[1].Version)
		assert.Equal(t, int64(0), history[2].Version)

		// Check operations
		assert.Equal(t, "UPDATE", history[0].Operation)
		assert.Equal(t, "APPEND", history[1].Operation)
		assert.Equal(t, "CREATE_TABLE", history[2].Operation)
	})

	t.Run("CompareVersions", func(t *testing.T) {
		// Compare versions
		comparison, err := manager.CompareVersions(0, 2)
		require.NoError(t, err)

		// Check comparison
		assert.Equal(t, int64(0), comparison.OlderVersion)
		assert.Equal(t, int64(2), comparison.NewerVersion)
		assert.Equal(t, 1, comparison.FilesAdded)
		assert.Equal(t, 0, comparison.FilesRemoved)
	})

	t.Run("RollbackToVersion", func(t *testing.T) {
		// Rollback to version 1 with force=true to bypass the 30-day restriction
		err := manager.RollbackToVersion(1, true)
		require.NoError(t, err)

		// Check new version was created
		versions, err := manager.GetVersions()
		require.NoError(t, err)
		require.Len(t, versions, 4)
		assert.Equal(t, int64(3), versions[0])

		// Check rollback operation
		history, err := manager.GetVersionHistory()
		require.NoError(t, err)
		assert.Equal(t, "ROLLBACK", history[0].Operation)
		assert.Equal(t, "1", history[0].Parameters["targetVersion"])
	})

	t.Run("RollbackToNonExistentVersion", func(t *testing.T) {
		// Try to rollback to non-existent version
		err := manager.RollbackToVersion(999, true)
		assert.Error(t, err)
	})
}

// Helper function to create a test version
func createTestVersion(t *testing.T, logDir string, version int64, operation string, timestamp time.Time) {
	// Create schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	// Create files list
	files := []string{"part-00000.parquet"}
	if version > 0 {
		files = append(files, fmt.Sprintf("part-%05d.parquet", version))
	}

	// Create metadata
	metadata := map[string]interface{}{
		"version":    version,
		"timestamp":  timestamp.Unix(),
		"schema":     schemaToMap(schema),
		"files":      files,
		"partitions": map[string][]string{},
		"metadata":   map[string]interface{}{},
	}

	// Convert to JSON
	data, err := json.Marshal(metadata)
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
	commitData, err := json.Marshal(commitInfo)
	require.NoError(t, err)

	// Write commit info
	commitFile := filepath.Join(logDir, fmt.Sprintf("%020d.commit.json", version))
	err = os.WriteFile(commitFile, commitData, 0644)
	require.NoError(t, err)
}
