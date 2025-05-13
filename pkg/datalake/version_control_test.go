package datalake

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestVersionManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "version_manager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a version manager
	vm := NewVersionManager(tempDir)

	// Test recording a transaction
	t.Run("RecordTransaction", func(t *testing.T) {
		// Record a transaction
		commitInfo := map[string]string{
			"message": "Initial commit",
			"user":    "test",
		}
		addedFiles := []string{"file1.parquet", "file2.parquet"}
		metadataChange := &MetadataChange{
			SchemaChange: true,
			AddedColumns: []string{"id", "name"},
		}
		stats := map[string]interface{}{
			"numRecords": float64(100),
			"bytesAdded": float64(1024),
		}

		tx, err := vm.RecordTransaction(
			"WRITE",
			commitInfo,
			addedFiles,
			nil,
			metadataChange,
			stats,
		)
		if err != nil {
			t.Fatalf("Failed to record transaction: %v", err)
		}

		// Check transaction
		if tx.Version != 1 {
			t.Errorf("Expected version 1, got %d", tx.Version)
		}

		if tx.Operation != "WRITE" {
			t.Errorf("Expected operation WRITE, got %s", tx.Operation)
		}

		if len(tx.AddedFiles) != 2 {
			t.Errorf("Expected 2 added files, got %d", len(tx.AddedFiles))
		}

		if !tx.MetadataChange.SchemaChange {
			t.Errorf("Expected schema change to be true")
		}

		if len(tx.MetadataChange.AddedColumns) != 2 {
			t.Errorf("Expected 2 added columns, got %d", len(tx.MetadataChange.AddedColumns))
		}

		if numRecords, ok := tx.Stats["numRecords"].(float64); !ok || numRecords != 100 {
			t.Errorf("Expected numRecords 100, got %v", tx.Stats["numRecords"])
		}
	})

	// Test getting transaction history
	t.Run("GetTransactionHistory", func(t *testing.T) {
		history, err := vm.GetTransactionHistory()
		if err != nil {
			t.Fatalf("Failed to get transaction history: %v", err)
		}

		if len(history.Transactions) != 1 {
			t.Errorf("Expected 1 transaction, got %d", len(history.Transactions))
		}

		if history.CurrentIndex != 0 {
			t.Errorf("Expected current index 0, got %d", history.CurrentIndex)
		}
	})

	// Test getting a specific transaction
	t.Run("GetTransaction", func(t *testing.T) {
		tx, err := vm.GetTransaction(1)
		if err != nil {
			t.Fatalf("Failed to get transaction: %v", err)
		}

		if tx.Version != 1 {
			t.Errorf("Expected version 1, got %d", tx.Version)
		}

		// Test getting a non-existent transaction
		_, err = vm.GetTransaction(99)
		if err == nil {
			t.Errorf("Expected error for non-existent transaction, got nil")
		}
	})

	// Test getting the latest transaction
	t.Run("GetLatestTransaction", func(t *testing.T) {
		tx, err := vm.GetLatestTransaction()
		if err != nil {
			t.Fatalf("Failed to get latest transaction: %v", err)
		}

		if tx.Version != 1 {
			t.Errorf("Expected version 1, got %d", tx.Version)
		}
	})

	// Record more transactions for testing
	t.Run("RecordMoreTransactions", func(t *testing.T) {
		// Record a second transaction (UPDATE)
		commitInfo := map[string]string{
			"message": "Update records",
			"user":    "test",
		}
		metadataChange := &MetadataChange{
			SchemaChange:    false,
			ModifiedColumns: []string{"name"},
		}
		stats := map[string]interface{}{
			"numRecords": float64(50),
		}

		tx, err := vm.RecordTransaction(
			"UPDATE",
			commitInfo,
			nil,
			nil,
			metadataChange,
			stats,
		)
		if err != nil {
			t.Fatalf("Failed to record transaction: %v", err)
		}

		if tx.Version != 2 {
			t.Errorf("Expected version 2, got %d", tx.Version)
		}

		// Record a third transaction (DELETE)
		commitInfo = map[string]string{
			"message": "Delete records",
			"user":    "test",
		}
		stats = map[string]interface{}{
			"numRecords": float64(25),
		}

		tx, err = vm.RecordTransaction(
			"DELETE",
			commitInfo,
			nil,
			[]string{"file2.parquet"},
			nil,
			stats,
		)
		if err != nil {
			t.Fatalf("Failed to record transaction: %v", err)
		}

		if tx.Version != 3 {
			t.Errorf("Expected version 3, got %d", tx.Version)
		}

		// Record a fourth transaction (schema change)
		commitInfo = map[string]string{
			"message": "Add email column",
			"user":    "test",
		}
		metadataChange = &MetadataChange{
			SchemaChange: true,
			AddedColumns: []string{"email"},
		}

		tx, err = vm.RecordTransaction(
			"WRITE",
			commitInfo,
			[]string{"file3.parquet"},
			nil,
			metadataChange,
			nil,
		)
		if err != nil {
			t.Fatalf("Failed to record transaction: %v", err)
		}

		if tx.Version != 4 {
			t.Errorf("Expected version 4, got %d", tx.Version)
		}
	})

	// Test rollback
	t.Run("RollbackToVersion", func(t *testing.T) {
		// Rollback to version 2
		err := vm.RollbackToVersion(2)
		if err != nil {
			t.Fatalf("Failed to rollback: %v", err)
		}

		// Check current version
		history, err := vm.GetTransactionHistory()
		if err != nil {
			t.Fatalf("Failed to get transaction history: %v", err)
		}

		// The current index should point to version 2
		// But we should have 5 transactions (1-4 plus the rollback)
		if len(history.Transactions) != 5 {
			t.Errorf("Expected 5 transactions, got %d", len(history.Transactions))
		}

		// The latest transaction should be a ROLLBACK
		latestTx := history.Transactions[len(history.Transactions)-1]
		if latestTx.Operation != "ROLLBACK" {
			t.Errorf("Expected ROLLBACK operation, got %s", latestTx.Operation)
		}

		// Test rollback to non-existent version
		err = vm.RollbackToVersion(99)
		if err == nil {
			t.Errorf("Expected error for rollback to non-existent version, got nil")
		}
	})

	// Test comparing versions
	t.Run("CompareVersions", func(t *testing.T) {
		// Compare versions 1 and 3
		diff, err := vm.CompareVersions(1, 3)
		if err != nil {
			t.Fatalf("Failed to compare versions: %v", err)
		}

		if diff.FromVersion != 1 || diff.ToVersion != 3 {
			t.Errorf("Expected from version 1 to version 3, got %d to %d", diff.FromVersion, diff.ToVersion)
		}

		// We should have 2 operations (UPDATE and DELETE)
		if len(diff.Operations) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(diff.Operations))
		}

		// Check row counts
		if diff.ModifiedRows != 50 {
			t.Errorf("Expected 50 modified rows, got %d", diff.ModifiedRows)
		}

		if diff.RemovedRows != 25 {
			t.Errorf("Expected 25 removed rows, got %d", diff.RemovedRows)
		}

		// Test comparing invalid versions
		_, err = vm.CompareVersions(3, 1) // From > To
		if err == nil {
			t.Errorf("Expected error for invalid version comparison, got nil")
		}

		_, err = vm.CompareVersions(1, 99) // Non-existent version
		if err == nil {
			t.Errorf("Expected error for non-existent version, got nil")
		}
	})

	// Test getting versions in time range
	t.Run("GetVersionsInTimeRange", func(t *testing.T) {
		// Get all transactions
		history, err := vm.GetTransactionHistory()
		if err != nil {
			t.Fatalf("Failed to get transaction history: %v", err)
		}

		// Get versions in time range
		startTime := history.Transactions[0].Timestamp
		endTime := history.Transactions[len(history.Transactions)-1].Timestamp
		versions, err := vm.GetVersionsInTimeRange(startTime, endTime)
		if err != nil {
			t.Fatalf("Failed to get versions in time range: %v", err)
		}

		// We should have all versions
		if len(versions) != len(history.Transactions) {
			t.Errorf("Expected %d versions, got %d", len(history.Transactions), len(versions))
		}

		// Test with a narrower time range
		startTime = history.Transactions[1].Timestamp
		endTime = history.Transactions[2].Timestamp
		versions, err = vm.GetVersionsInTimeRange(startTime, endTime)
		if err != nil {
			t.Fatalf("Failed to get versions in time range: %v", err)
		}

		// We should have 2 versions
		if len(versions) != 2 {
			t.Errorf("Expected 2 versions, got %d", len(versions))
		}
	})

	// Test getting version at timestamp
	t.Run("GetVersionAtTimestamp", func(t *testing.T) {
		// Get all transactions
		history, err := vm.GetTransactionHistory()
		if err != nil {
			t.Fatalf("Failed to get transaction history: %v", err)
		}

		// Get version at timestamp
		timestamp := history.Transactions[2].Timestamp.Add(time.Millisecond)
		tx, err := vm.GetVersionAtTimestamp(timestamp)
		if err != nil {
			t.Fatalf("Failed to get version at timestamp: %v", err)
		}

		// We should get version 3
		if tx.Version != 3 {
			t.Errorf("Expected version 3, got %d", tx.Version)
		}

		// Test with a timestamp before any version
		timestamp = history.Transactions[0].Timestamp.Add(-time.Hour)
		_, err = vm.GetVersionAtTimestamp(timestamp)
		if err == nil {
			t.Errorf("Expected error for timestamp before any version, got nil")
		}
	})

	// Test transaction summary
	t.Run("TransactionSummary", func(t *testing.T) {
		tx, err := vm.GetTransaction(1)
		if err != nil {
			t.Fatalf("Failed to get transaction: %v", err)
		}

		summary := tx.GetSummary()
		if summary == "" {
			t.Errorf("Expected non-empty summary, got empty string")
		}

		if !contains(summary, "Version 1") {
			t.Errorf("Expected summary to contain 'Version 1', got: %s", summary)
		}

		if !contains(summary, "WRITE") {
			t.Errorf("Expected summary to contain 'WRITE', got: %s", summary)
		}

		if !contains(summary, "Initial commit") {
			t.Errorf("Expected summary to contain commit message, got: %s", summary)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && strings.Contains(s, substr)
}
