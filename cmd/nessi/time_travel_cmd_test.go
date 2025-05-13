package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeTravelCommands(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "time-travel-cmd-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test table
	tablePath := filepath.Join(tempDir, "test-table")
	err = os.MkdirAll(tablePath, 0755)
	require.NoError(t, err)

	// Create version manager and schema manager
	vm := datalake.NewVersionManager(tablePath)
	sm := datalake.NewSchemaManager(tablePath)

	// Initialize schema
	schema := &datalake.Schema{
		Fields: []datalake.Field{
			{Name: "id", Type: "integer", Nullable: false},
			{Name: "name", Type: "string", Nullable: true},
			{Name: "age", Type: "integer", Nullable: true},
		},
	}
	err = sm.InitializeSchema(schema)
	require.NoError(t, err)

	// Record initial transaction
	initialFiles := []string{"part-00000.parquet", "part-00001.parquet"}
	tx1, err := vm.RecordTransaction(
		"WRITE",
		map[string]string{"message": "Initial data load"},
		initialFiles,
		nil,
		nil,
		map[string]interface{}{"numRecords": float64(100)},
	)
	require.NoError(t, err)
	assert.Equal(t, 0, tx1.Version)

	// Wait a moment to ensure timestamps are different
	time.Sleep(10 * time.Millisecond)
	timestamp1 := time.Now()
	time.Sleep(10 * time.Millisecond)

	// Record second transaction - add a column and some files
	updatedSchema := &datalake.Schema{
		Fields: []datalake.Field{
			{Name: "id", Type: "integer", Nullable: false},
			{Name: "name", Type: "string", Nullable: true},
			{Name: "age", Type: "integer", Nullable: true},
			{Name: "email", Type: "string", Nullable: true},
		},
	}
	err = sm.UpdateSchema(updatedSchema)
	require.NoError(t, err)

	additionalFiles := []string{"part-00002.parquet", "part-00003.parquet"}
	tx2, err := vm.RecordTransaction(
		"UPDATE",
		map[string]string{"message": "Add email column"},
		additionalFiles,
		nil,
		&datalake.MetadataChange{
			SchemaChange:    true,
			AddedColumns:    []string{"email"},
			RemovedColumns:  nil,
			ModifiedColumns: nil,
		},
		map[string]interface{}{"numRecords": float64(50)},
	)
	require.NoError(t, err)
	assert.Equal(t, 1, tx2.Version)

	// Helper function to execute a command and capture its output
	executeCommand := func(cmd *cobra.Command, args ...string) (string, error) {
		buf := new(bytes.Buffer)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs(args)
		err := cmd.Execute()
		return buf.String(), err
	}

	// Test query-version command
	t.Run("QueryVersion", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(queryVersionCmd)

		// Test with text format
		output, err := executeCommand(cmd, "query-version", tablePath, "0")
		assert.NoError(t, err)
		assert.Contains(t, output, "Query result for version 0")
		assert.Contains(t, output, "Initial data load")
		assert.Contains(t, output, "Files: 2")

		// Test with JSON format
		output, err = executeCommand(cmd, "query-version", tablePath, "1", "--format", "json")
		assert.NoError(t, err)
		assert.Contains(t, output, "\"version\": 1")
		assert.Contains(t, output, "\"message\": \"Add email column\"")
		assert.Contains(t, output, "\"email\"")

		// Test with show-files flag
		output, err = executeCommand(cmd, "query-version", tablePath, "1", "--show-files")
		assert.NoError(t, err)
		assert.Contains(t, output, "Files:")
		assert.Contains(t, output, "part-00000.parquet")
		assert.Contains(t, output, "part-00002.parquet")

		// Test with show-schema flag
		output, err = executeCommand(cmd, "query-version", tablePath, "1", "--show-schema")
		assert.NoError(t, err)
		assert.Contains(t, output, "Schema:")
		assert.Contains(t, output, "id: integer (not null)")
		assert.Contains(t, output, "email: string")

		// Test with invalid version
		output, err = executeCommand(cmd, "query-version", tablePath, "99")
		assert.Error(t, err)
		assert.Contains(t, output, "Error querying version")
	})

	// Test query-timestamp command
	t.Run("QueryTimestamp", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(queryTimestampCmd)

		// Test with timestamp of first transaction
		output, err := executeCommand(cmd, "query-timestamp", tablePath, tx1.Timestamp.Format(time.RFC3339))
		assert.NoError(t, err)
		assert.Contains(t, output, "Query result for timestamp")
		assert.Contains(t, output, "Actual version: 0")
		assert.Contains(t, output, "Initial data load")

		// Test with timestamp between transactions
		output, err = executeCommand(cmd, "query-timestamp", tablePath, timestamp1.Format(time.RFC3339))
		assert.NoError(t, err)
		assert.Contains(t, output, "Actual version: 0")

		// Test with timestamp of second transaction
		output, err = executeCommand(cmd, "query-timestamp", tablePath, tx2.Timestamp.Format(time.RFC3339))
		assert.NoError(t, err)
		assert.Contains(t, output, "Actual version: 1")
		assert.Contains(t, output, "Add email column")

		// Test with invalid timestamp format
		output, err = executeCommand(cmd, "query-timestamp", tablePath, "2023-01-01")
		assert.Error(t, err)
		assert.Contains(t, output, "Error parsing timestamp")
	})

	// Test versions-in-range command
	t.Run("VersionsInRange", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(versionsInRangeCmd)

		// Test with range including both transactions
		startTime := tx1.Timestamp.Add(-1 * time.Second)
		endTime := tx2.Timestamp.Add(1 * time.Second)
		output, err := executeCommand(cmd, "versions-in-range", tablePath, 
			startTime.Format(time.RFC3339), 
			endTime.Format(time.RFC3339))
		assert.NoError(t, err)
		assert.Contains(t, output, "Versions between")
		assert.Contains(t, output, "Found 2 versions")
		assert.Contains(t, output, "Version 0")
		assert.Contains(t, output, "Version 1")

		// Test with range including only first transaction
		endTime = tx1.Timestamp.Add(1 * time.Millisecond)
		output, err = executeCommand(cmd, "versions-in-range", tablePath, 
			startTime.Format(time.RFC3339), 
			endTime.Format(time.RFC3339))
		assert.NoError(t, err)
		assert.Contains(t, output, "Found 1 versions")
		assert.Contains(t, output, "Version 0")
		assert.NotContains(t, output, "Version 1")

		// Test with empty range
		startTime = tx2.Timestamp.Add(1 * time.Second)
		endTime = tx2.Timestamp.Add(2 * time.Second)
		output, err = executeCommand(cmd, "versions-in-range", tablePath, 
			startTime.Format(time.RFC3339), 
			endTime.Format(time.RFC3339))
		assert.NoError(t, err)
		assert.Contains(t, output, "Found 0 versions")
	})

	// Test export-snapshot command
	t.Run("ExportSnapshot", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(exportSnapshotCmd)

		// Create a directory for the snapshot
		snapshotDir := filepath.Join(tempDir, "snapshot")

		// Test exporting version 1
		output, err := executeCommand(cmd, "export-snapshot", tablePath, "1", snapshotDir)
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully exported snapshot of version 1")
		
		// Check that files were created
		assert.FileExists(t, filepath.Join(snapshotDir, "metadata.json"))
		assert.FileExists(t, filepath.Join(snapshotDir, "schema.json"))
		assert.FileExists(t, filepath.Join(snapshotDir, "files.json"))
		assert.FileExists(t, filepath.Join(snapshotDir, "transaction.json"))
		assert.DirExists(t, filepath.Join(snapshotDir, "data"))

		// Test with invalid version
		invalidDir := filepath.Join(tempDir, "invalid-snapshot")
		output, err = executeCommand(cmd, "export-snapshot", tablePath, "99", invalidDir)
		assert.Error(t, err)
		assert.Contains(t, output, "Error exporting snapshot")
	})

	// Test reconstruct-state command
	t.Run("ReconstructState", func(t *testing.T) {
		// Create a new command for testing
		cmd := &cobra.Command{Use: "test"}
		cmd.AddCommand(reconstructStateCmd)

		// Create a directory for the reconstruction
		reconstructDir := filepath.Join(tempDir, "reconstruct")

		// Test reconstructing version 0
		output, err := executeCommand(cmd, "reconstruct-state", tablePath, "0", reconstructDir)
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully reconstructed state at version 0")
		
		// Check that files were created
		assert.FileExists(t, filepath.Join(reconstructDir, "metadata.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "schema.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "files.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "transaction.json"))
		assert.DirExists(t, filepath.Join(reconstructDir, "data"))

		// Create a directory for the timestamp reconstruction
		timestampDir := filepath.Join(tempDir, "timestamp-reconstruct")

		// Test reconstructing at timestamp
		output, err = executeCommand(cmd, "reconstruct-state", tablePath, 
			tx1.Timestamp.Format(time.RFC3339), timestampDir, "--timestamp")
		assert.NoError(t, err)
		assert.Contains(t, output, "Successfully reconstructed state at timestamp")
		
		// Check that files were created
		assert.FileExists(t, filepath.Join(timestampDir, "metadata.json"))
		assert.FileExists(t, filepath.Join(timestampDir, "schema.json"))
		assert.FileExists(t, filepath.Join(timestampDir, "files.json"))
		assert.FileExists(t, filepath.Join(timestampDir, "transaction.json"))
		assert.DirExists(t, filepath.Join(timestampDir, "data"))

		// Test with invalid version
		invalidDir := filepath.Join(tempDir, "invalid-reconstruct")
		output, err = executeCommand(cmd, "reconstruct-state", tablePath, "99", invalidDir)
		assert.Error(t, err)
		assert.Contains(t, output, "Error reconstructing state")
	})
}
