package datalake

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeTravel(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "time-travel-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test table
	tablePath := filepath.Join(tempDir, "test-table")
	err = os.MkdirAll(tablePath, 0755)
	require.NoError(t, err)

	// Create version manager and schema manager
	vm := NewVersionManager(tablePath)
	sm := NewSchemaManager(tablePath)

	// Initialize schema
	schema := &Schema{
		Fields: []Field{
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
	updatedSchema := &Schema{
		Fields: []Field{
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
		&MetadataChange{
			SchemaChange:    true,
			AddedColumns:    []string{"email"},
			RemovedColumns:  nil,
			ModifiedColumns: nil,
		},
		map[string]interface{}{"numRecords": float64(50)},
	)
	require.NoError(t, err)
	assert.Equal(t, 1, tx2.Version)

	// Wait a moment to ensure timestamps are different
	time.Sleep(10 * time.Millisecond)
	timestamp2 := time.Now()
	time.Sleep(10 * time.Millisecond)

	// Record third transaction - remove some files
	removedFiles := []string{"part-00000.parquet"}
	tx3, err := vm.RecordTransaction(
		"UPDATE",
		map[string]string{"message": "Remove old data"},
		nil,
		removedFiles,
		nil,
		map[string]interface{}{"numRecords": float64(-25)},
	)
	require.NoError(t, err)
	assert.Equal(t, 2, tx3.Version)

	// Create time travel manager
	tt := NewTimeTravel(tablePath)

	// Test GetVersionForTimestamp
	t.Run("GetVersionForTimestamp", func(t *testing.T) {
		// Before first transaction
		beforeTx, err := tt.GetVersionForTimestamp(tx1.Timestamp.Add(-1 * time.Hour))
		assert.Error(t, err)
		assert.Nil(t, beforeTx)

		// At first transaction
		atTx1, err := tt.GetVersionForTimestamp(tx1.Timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 0, atTx1.Version)

		// Between first and second transaction
		betweenTx, err := tt.GetVersionForTimestamp(timestamp1)
		assert.NoError(t, err)
		assert.Equal(t, 0, betweenTx.Version)

		// At second transaction
		atTx2, err := tt.GetVersionForTimestamp(tx2.Timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 1, atTx2.Version)

		// Between second and third transaction
		betweenTx2, err := tt.GetVersionForTimestamp(timestamp2)
		assert.NoError(t, err)
		assert.Equal(t, 1, betweenTx2.Version)

		// At third transaction
		atTx3, err := tt.GetVersionForTimestamp(tx3.Timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 2, atTx3.Version)

		// After third transaction
		afterTx, err := tt.GetVersionForTimestamp(tx3.Timestamp.Add(1 * time.Hour))
		assert.NoError(t, err)
		assert.Equal(t, 2, afterTx.Version)
	})

	// Test GetFilesAtVersion
	t.Run("GetFilesAtVersion", func(t *testing.T) {
		// Version 0
		filesV0, err := tt.GetFilesAtVersion(0)
		assert.NoError(t, err)
		assert.ElementsMatch(t, initialFiles, filesV0)

		// Version 1
		filesV1, err := tt.GetFilesAtVersion(1)
		assert.NoError(t, err)
		expectedFilesV1 := append(initialFiles, additionalFiles...)
		assert.ElementsMatch(t, expectedFilesV1, filesV1)

		// Version 2
		filesV2, err := tt.GetFilesAtVersion(2)
		assert.NoError(t, err)
		expectedFilesV2 := []string{"part-00001.parquet", "part-00002.parquet", "part-00003.parquet"}
		assert.ElementsMatch(t, expectedFilesV2, filesV2)

		// Invalid version
		filesInvalid, err := tt.GetFilesAtVersion(99)
		assert.Error(t, err)
		assert.Nil(t, filesInvalid)
	})

	// Test GetSchemaAtVersion
	t.Run("GetSchemaAtVersion", func(t *testing.T) {
		// Version 0
		schemaV0, err := tt.GetSchemaAtVersion(0)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(schemaV0.Fields))
		assert.Equal(t, "id", schemaV0.Fields[0].Name)
		assert.Equal(t, "name", schemaV0.Fields[1].Name)
		assert.Equal(t, "age", schemaV0.Fields[2].Name)

		// Version 1
		schemaV1, err := tt.GetSchemaAtVersion(1)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(schemaV1.Fields))
		assert.Equal(t, "id", schemaV1.Fields[0].Name)
		assert.Equal(t, "name", schemaV1.Fields[1].Name)
		assert.Equal(t, "age", schemaV1.Fields[2].Name)
		assert.Equal(t, "email", schemaV1.Fields[3].Name)

		// Version 2 (same as version 1)
		schemaV2, err := tt.GetSchemaAtVersion(2)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(schemaV2.Fields))
		assert.Equal(t, "email", schemaV2.Fields[3].Name)
	})

	// Test QueryAtVersion
	t.Run("QueryAtVersion", func(t *testing.T) {
		// Version 0
		resultV0, err := tt.QueryAtVersion(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, resultV0.Transaction.Version)
		assert.ElementsMatch(t, initialFiles, resultV0.Files)
		assert.Equal(t, 3, len(resultV0.Schema.Fields))
		assert.Equal(t, "Initial data load", resultV0.Metadata["message"])

		// Version 1
		resultV1, err := tt.QueryAtVersion(1)
		assert.NoError(t, err)
		assert.Equal(t, 1, resultV1.Transaction.Version)
		expectedFilesV1 := append(initialFiles, additionalFiles...)
		assert.ElementsMatch(t, expectedFilesV1, resultV1.Files)
		assert.Equal(t, 4, len(resultV1.Schema.Fields))
		assert.Equal(t, "Add email column", resultV1.Metadata["message"])

		// Version 2
		resultV2, err := tt.QueryAtVersion(2)
		assert.NoError(t, err)
		assert.Equal(t, 2, resultV2.Transaction.Version)
		expectedFilesV2 := []string{"part-00001.parquet", "part-00002.parquet", "part-00003.parquet"}
		assert.ElementsMatch(t, expectedFilesV2, resultV2.Files)
		assert.Equal(t, 4, len(resultV2.Schema.Fields))
		assert.Equal(t, "Remove old data", resultV2.Metadata["message"])

		// Invalid version
		resultInvalid, err := tt.QueryAtVersion(99)
		assert.Error(t, err)
		assert.Nil(t, resultInvalid)
	})

	// Test QueryAtTimestamp
	t.Run("QueryAtTimestamp", func(t *testing.T) {
		// Before first transaction
		beforeResult, err := tt.QueryAtTimestamp(tx1.Timestamp.Add(-1 * time.Hour))
		assert.Error(t, err)
		assert.Nil(t, beforeResult)

		// At first transaction
		atTx1Result, err := tt.QueryAtTimestamp(tx1.Timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 0, atTx1Result.Transaction.Version)

		// Between first and second transaction
		betweenResult, err := tt.QueryAtTimestamp(timestamp1)
		assert.NoError(t, err)
		assert.Equal(t, 0, betweenResult.Transaction.Version)

		// At second transaction
		atTx2Result, err := tt.QueryAtTimestamp(tx2.Timestamp)
		assert.NoError(t, err)
		assert.Equal(t, 1, atTx2Result.Transaction.Version)

		// After third transaction
		afterResult, err := tt.QueryAtTimestamp(tx3.Timestamp.Add(1 * time.Hour))
		assert.NoError(t, err)
		assert.Equal(t, 2, afterResult.Transaction.Version)
	})

	// Test ReadParsedAtVersion
	t.Run("ReadParsedAtVersion", func(t *testing.T) {
		// Version 0
		dataV0, err := tt.ReadParsedAtVersion(0)
		assert.NoError(t, err)
		assert.Equal(t, 10, dataV0.Count)
		assert.Equal(t, 3, len(dataV0.Schema))
		assert.Contains(t, dataV0.Schema, "id")
		assert.Contains(t, dataV0.Schema, "name")
		assert.Contains(t, dataV0.Schema, "age")

		// Version 1
		dataV1, err := tt.ReadParsedAtVersion(1)
		assert.NoError(t, err)
		assert.Equal(t, 10, dataV1.Count)
		assert.Equal(t, 4, len(dataV1.Schema))
		assert.Contains(t, dataV1.Schema, "id")
		assert.Contains(t, dataV1.Schema, "name")
		assert.Contains(t, dataV1.Schema, "age")
		assert.Contains(t, dataV1.Schema, "email")

		// Check record structure
		assert.Equal(t, 10, len(dataV1.Records))
		record := dataV1.Records[0]
		assert.Contains(t, record, "id")
		assert.Contains(t, record, "name")
		assert.Contains(t, record, "age")
		assert.Contains(t, record, "email")
	})

	// Test GetVersionsInTimeRange
	t.Run("GetVersionsInTimeRange", func(t *testing.T) {
		// Get all versions
		allVersions, err := tt.GetVersionsInTimeRange(
			tx1.Timestamp.Add(-1*time.Hour),
			tx3.Timestamp.Add(1*time.Hour),
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(allVersions))
		assert.Equal(t, 0, allVersions[0].Version)
		assert.Equal(t, 1, allVersions[1].Version)
		assert.Equal(t, 2, allVersions[2].Version)

		// Get versions between first and third
		middleVersions, err := tt.GetVersionsInTimeRange(
			tx1.Timestamp.Add(1*time.Millisecond),
			tx3.Timestamp.Add(-1*time.Millisecond),
		)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(middleVersions))
		assert.Equal(t, 1, middleVersions[0].Version)

		// Get no versions
		noVersions, err := tt.GetVersionsInTimeRange(
			tx3.Timestamp.Add(1*time.Hour),
			tx3.Timestamp.Add(2*time.Hour),
		)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(noVersions))
	})

	// Test ExportVersionSnapshot
	t.Run("ExportVersionSnapshot", func(t *testing.T) {
		// Create a temporary directory for the snapshot
		snapshotDir, err := os.MkdirTemp("", "snapshot-test")
		require.NoError(t, err)
		defer os.RemoveAll(snapshotDir)

		// Export snapshot for version 1
		err = tt.ExportVersionSnapshot(1, snapshotDir)
		assert.NoError(t, err)

		// Check that files were created
		assert.FileExists(t, filepath.Join(snapshotDir, "metadata.json"))
		assert.FileExists(t, filepath.Join(snapshotDir, "schema.json"))
		assert.FileExists(t, filepath.Join(snapshotDir, "files.json"))
		assert.FileExists(t, filepath.Join(snapshotDir, "transaction.json"))
		assert.DirExists(t, filepath.Join(snapshotDir, "data"))
		assert.FileExists(t, filepath.Join(snapshotDir, "data", "README.txt"))

		// Invalid version
		invalidDir := filepath.Join(tempDir, "invalid-snapshot")
		err = tt.ExportVersionSnapshot(99, invalidDir)
		assert.Error(t, err)
		assert.NoDirExists(t, invalidDir)
	})

	// Test ReconstructStateAtVersion
	t.Run("ReconstructStateAtVersion", func(t *testing.T) {
		// Create a temporary directory for the reconstruction
		reconstructDir, err := os.MkdirTemp("", "reconstruct-test")
		require.NoError(t, err)
		defer os.RemoveAll(reconstructDir)

		// Reconstruct state for version 2
		err = tt.ReconstructStateAtVersion(2, reconstructDir)
		assert.NoError(t, err)

		// Check that files were created
		assert.FileExists(t, filepath.Join(reconstructDir, "metadata.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "schema.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "files.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "transaction.json"))
		assert.DirExists(t, filepath.Join(reconstructDir, "data"))
	})

	// Test ReconstructStateAtTimestamp
	t.Run("ReconstructStateAtTimestamp", func(t *testing.T) {
		// Create a temporary directory for the reconstruction
		reconstructDir, err := os.MkdirTemp("", "reconstruct-timestamp-test")
		require.NoError(t, err)
		defer os.RemoveAll(reconstructDir)

		// Reconstruct state at timestamp between version 1 and 2
		err = tt.ReconstructStateAtTimestamp(timestamp2, reconstructDir)
		assert.NoError(t, err)

		// Check that files were created
		assert.FileExists(t, filepath.Join(reconstructDir, "metadata.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "schema.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "files.json"))
		assert.FileExists(t, filepath.Join(reconstructDir, "transaction.json"))
		assert.DirExists(t, filepath.Join(reconstructDir, "data"))

		// Invalid timestamp (before first transaction)
		invalidDir := filepath.Join(tempDir, "invalid-timestamp")
		err = tt.ReconstructStateAtTimestamp(tx1.Timestamp.Add(-1*time.Hour), invalidDir)
		assert.Error(t, err)
		assert.NoDirExists(t, invalidDir)
	})
}
