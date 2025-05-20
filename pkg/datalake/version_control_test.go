package datalake

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "version_manager_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a version manager
	vm := NewVersionManager(tempDir)

	// Test recording a transaction
	t.Run("RecordTransaction", func(t *testing.T) {
		// Record a transaction
		tx, err := vm.RecordTransaction(
			"WRITE",
			map[string]string{"message": "Initial commit"},
			[]string{"data1.parquet", "data2.parquet"},
			nil,
			nil,
			map[string]interface{}{"numRecords": float64(100)},
		)
		require.NoError(t, err)
		assert.Equal(t, 0, tx.Version)
		assert.Equal(t, "WRITE", tx.Operation)
		// Check that we have files added
		assert.Equal(t, 2, len(tx.AddedFiles))
		assert.Equal(t, 0, len(tx.RemovedFiles))
	})

	// Test getting a transaction
	t.Run("GetTransaction", func(t *testing.T) {
		tx, err := vm.GetTransaction(0)
		require.NoError(t, err)
		assert.Equal(t, 0, tx.Version)
		assert.Equal(t, "WRITE", tx.Operation)
	})

	// Test getting latest transaction
	t.Run("GetLatestTransaction", func(t *testing.T) {
		// Record another transaction
		_, err := vm.RecordTransaction(
			"UPDATE",
			map[string]string{"message": "Update data"},
			[]string{"data3.parquet"},
			[]string{"data1.parquet"},
			nil,
			map[string]interface{}{"numRecords": float64(50)},
		)
		require.NoError(t, err)

		// Get latest transaction
		tx, err := vm.GetLatestTransaction()
		require.NoError(t, err)
		assert.Equal(t, 1, tx.Version)
		assert.Equal(t, "UPDATE", tx.Operation)
	})

	// Test transaction summary
	t.Run("TransactionSummary", func(t *testing.T) {
		// Get the first transaction
		tx, err := vm.GetTransaction(0)
		require.NoError(t, err)

		// Get summary
		summary := tx.GetSummary()
		if summary == "" {
			t.Errorf("Expected non-empty summary, got empty string")
		}

		if !containsSubstring(summary, "Version 1") {
			t.Errorf("Expected summary to contain 'Version 1', got: %s", summary)
		}

		if !containsSubstring(summary, "WRITE") {
			t.Errorf("Expected summary to contain 'WRITE', got: %s", summary)
		}

		if !containsSubstring(summary, "Initial commit") {
			t.Errorf("Expected summary to contain commit message, got: %s", summary)
		}
	})
}

// Helper function to check if a string contains a substring
func containsSubstring(s, substr string) bool {
	return s != "" && strings.Contains(s, substr)
}
