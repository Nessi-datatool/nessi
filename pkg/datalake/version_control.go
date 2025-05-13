package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
)

// VersionManager handles Delta Lake version control operations
type VersionManager struct {
	tablePath string
	metaDir   string
}

// TransactionHistory represents the history of transactions
type TransactionHistory struct {
	Transactions []*Transaction `json:"transactions"`
	CurrentIndex int            `json:"current_index"`
}

// Transaction represents a transaction in the Delta Lake log
type Transaction struct {
	ID            string                 `json:"id"`
	Version       int                    `json:"version"`
	Timestamp     time.Time              `json:"timestamp"`
	Operation     string                 `json:"operation"`
	CommitInfo    map[string]string      `json:"commit_info"`
	AddedFiles    []string               `json:"added_files,omitempty"`
	RemovedFiles  []string               `json:"removed_files,omitempty"`
	MetadataChange *MetadataChange       `json:"metadata_change,omitempty"`
	Stats         map[string]interface{} `json:"stats,omitempty"`
	IsolationLevel string                `json:"isolation_level,omitempty"`
	ReadVersion   int                    `json:"read_version,omitempty"`
	UserID        string                 `json:"user_id,omitempty"`
	ClientInfo    map[string]string      `json:"client_info,omitempty"`
}

// MetadataChange represents a change to table metadata
type MetadataChange struct {
	SchemaChange   bool     `json:"schema_change,omitempty"`
	PartitionChange bool     `json:"partition_change,omitempty"`
	PropertiesChange bool     `json:"properties_change,omitempty"`
	AddedColumns   []string  `json:"added_columns,omitempty"`
	RemovedColumns []string  `json:"removed_columns,omitempty"`
	ModifiedColumns []string `json:"modified_columns,omitempty"`
}

// VersionDiff represents the difference between two versions
type VersionDiff struct {
	FromVersion   int            `json:"from_version"`
	ToVersion     int            `json:"to_version"`
	AddedRows     int64          `json:"added_rows"`
	RemovedRows   int64          `json:"removed_rows"`
	ModifiedRows  int64          `json:"modified_rows"`
	SchemaChanges []SchemaChange `json:"schema_changes,omitempty"`
	Operations    []string       `json:"operations"`
	TimeSpan      time.Duration  `json:"time_span"`
	FromTimestamp time.Time      `json:"from_timestamp"`
	ToTimestamp   time.Time      `json:"to_timestamp"`
}

// NewVersionManager creates a new version manager
func NewVersionManager(tablePath string) *VersionManager {
	return &VersionManager{
		tablePath: tablePath,
		metaDir:   filepath.Join(tablePath, "_delta_log", "transactions"),
	}
}

// RecordTransaction records a transaction in the history
func (vm *VersionManager) RecordTransaction(
	operation string,
	commitInfo map[string]string,
	addedFiles, removedFiles []string,
	metadataChange *MetadataChange,
	stats map[string]interface{},
) (*Transaction, error) {
	// Create transaction history directory if it doesn't exist
	if err := os.MkdirAll(vm.metaDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create transaction history directory: %w", err)
	}

	// Get current history
	history, err := vm.readHistory()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	// Initialize history if it doesn't exist
	if history == nil {
		history = &TransactionHistory{
			Transactions: []*Transaction{},
			CurrentIndex: -1,
		}
	}

	// Determine new version number
	version := 1
	if len(history.Transactions) > 0 {
		version = history.Transactions[len(history.Transactions)-1].Version + 1
	}

	// Create new transaction
	transaction := &Transaction{
		ID:            uuid.New().String(),
		Version:       version,
		Timestamp:     time.Now(),
		Operation:     operation,
		CommitInfo:    commitInfo,
		AddedFiles:    addedFiles,
		RemovedFiles:  removedFiles,
		MetadataChange: metadataChange,
		Stats:         stats,
		IsolationLevel: "Serializable", // Default isolation level
		ReadVersion:   version - 1,
		UserID:        os.Getenv("USER"),
		ClientInfo:    map[string]string{"client": "nessi-go"},
	}

	// Add transaction to history
	history.Transactions = append(history.Transactions, transaction)
	history.CurrentIndex = len(history.Transactions) - 1

	// Write history to file
	if err := vm.writeHistory(history); err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetTransactionHistory gets the transaction history
func (vm *VersionManager) GetTransactionHistory() (*TransactionHistory, error) {
	return vm.readHistory()
}

// GetTransaction gets a specific transaction
func (vm *VersionManager) GetTransaction(version int) (*Transaction, error) {
	history, err := vm.readHistory()
	if err != nil {
		return nil, err
	}

	for _, transaction := range history.Transactions {
		if transaction.Version == version {
			return transaction, nil
		}
	}

	return nil, fmt.Errorf("transaction version %d not found", version)
}

// GetLatestTransaction gets the latest transaction
func (vm *VersionManager) GetLatestTransaction() (*Transaction, error) {
	history, err := vm.readHistory()
	if err != nil {
		return nil, err
	}

	if len(history.Transactions) == 0 {
		return nil, fmt.Errorf("no transactions found")
	}

	return history.Transactions[history.CurrentIndex], nil
}

// RollbackToVersion rolls back to a specific version
func (vm *VersionManager) RollbackToVersion(version int) error {
	history, err := vm.readHistory()
	if err != nil {
		return err
	}

	// Find the target version
	targetIndex := -1
	for i, transaction := range history.Transactions {
		if transaction.Version == version {
			targetIndex = i
			break
		}
	}

	if targetIndex == -1 {
		return fmt.Errorf("version %d not found", version)
	}

	// Check if already at target version
	if history.CurrentIndex == targetIndex {
		return nil
	}

	// Record rollback transaction
	commitInfo := map[string]string{
		"operation": "rollback",
		"target_version": fmt.Sprintf("%d", version),
		"previous_version": fmt.Sprintf("%d", history.Transactions[history.CurrentIndex].Version),
	}

	// Create metadata change for rollback
	metadataChange := &MetadataChange{
		SchemaChange: true,
		PropertiesChange: true,
	}

	// Record rollback transaction
	_, err = vm.RecordTransaction(
		"ROLLBACK",
		commitInfo,
		nil,
		nil,
		metadataChange,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to record rollback transaction: %w", err)
	}

	// Update current index
	history, err = vm.readHistory()
	if err != nil {
		return err
	}

	// Set the current index to point to the target version
	// This is a simplified approach - in a real implementation,
	// we would need to actually restore the data files
	for i, transaction := range history.Transactions {
		if transaction.Version == version {
			history.CurrentIndex = i
			break
		}
	}

	// Write history to file
	if err := vm.writeHistory(history); err != nil {
		return err
	}

	return nil
}

// CompareVersions compares two versions and returns the differences
func (vm *VersionManager) CompareVersions(fromVersion, toVersion int) (*VersionDiff, error) {
	history, err := vm.readHistory()
	if err != nil {
		return nil, err
	}

	// Find the transactions
	var fromTx, toTx *Transaction
	var fromIndex, toIndex int
	for i, tx := range history.Transactions {
		if tx.Version == fromVersion {
			fromTx = tx
			fromIndex = i
		}
		if tx.Version == toVersion {
			toTx = tx
			toIndex = i
		}
	}

	if fromTx == nil {
		return nil, fmt.Errorf("from version %d not found", fromVersion)
	}
	if toTx == nil {
		return nil, fmt.Errorf("to version %d not found", toVersion)
	}

	// Ensure fromVersion is before toVersion
	if fromIndex > toIndex {
		return nil, fmt.Errorf("from version %d is after to version %d", fromVersion, toVersion)
	}

	// Initialize diff
	diff := &VersionDiff{
		FromVersion:   fromVersion,
		ToVersion:     toVersion,
		AddedRows:     0,
		RemovedRows:   0,
		ModifiedRows:  0,
		SchemaChanges: []SchemaChange{},
		Operations:    []string{},
		TimeSpan:      toTx.Timestamp.Sub(fromTx.Timestamp),
		FromTimestamp: fromTx.Timestamp,
		ToTimestamp:   toTx.Timestamp,
	}

	// Collect operations
	for i := fromIndex + 1; i <= toIndex; i++ {
		tx := history.Transactions[i]
		diff.Operations = append(diff.Operations, tx.Operation)

		// Count row changes (simplified)
		if tx.Stats != nil {
			if numRecords, ok := tx.Stats["numRecords"].(float64); ok {
				if tx.Operation == "WRITE" || tx.Operation == "MERGE" {
					diff.AddedRows += int64(numRecords)
				} else if tx.Operation == "DELETE" {
					diff.RemovedRows += int64(numRecords)
				} else if tx.Operation == "UPDATE" {
					diff.ModifiedRows += int64(numRecords)
				}
			}
		}

		// Collect schema changes
		if tx.MetadataChange != nil && tx.MetadataChange.SchemaChange {
			// In a real implementation, we would extract the actual schema changes
			// For now, we'll just use placeholder changes
			if len(tx.MetadataChange.AddedColumns) > 0 {
				for _, col := range tx.MetadataChange.AddedColumns {
					diff.SchemaChanges = append(diff.SchemaChanges, SchemaChange{
						Type:        "add",
						Field:       col,
						Description: fmt.Sprintf("Added column '%s'", col),
					})
				}
			}
			if len(tx.MetadataChange.RemovedColumns) > 0 {
				for _, col := range tx.MetadataChange.RemovedColumns {
					diff.SchemaChanges = append(diff.SchemaChanges, SchemaChange{
						Type:        "remove",
						Field:       col,
						Description: fmt.Sprintf("Removed column '%s'", col),
					})
				}
			}
			if len(tx.MetadataChange.ModifiedColumns) > 0 {
				for _, col := range tx.MetadataChange.ModifiedColumns {
					diff.SchemaChanges = append(diff.SchemaChanges, SchemaChange{
						Type:        "modify",
						Field:       col,
						Description: fmt.Sprintf("Modified column '%s'", col),
					})
				}
			}
		}
	}

	return diff, nil
}

// GetVersionsInTimeRange gets versions within a time range
func (vm *VersionManager) GetVersionsInTimeRange(startTime, endTime time.Time) ([]*Transaction, error) {
	history, err := vm.readHistory()
	if err != nil {
		return nil, err
	}

	var versions []*Transaction
	for _, tx := range history.Transactions {
		if (tx.Timestamp.Equal(startTime) || tx.Timestamp.After(startTime)) &&
			(tx.Timestamp.Equal(endTime) || tx.Timestamp.Before(endTime)) {
			versions = append(versions, tx)
		}
	}

	return versions, nil
}

// GetVersionAtTimestamp gets the version at a specific timestamp
func (vm *VersionManager) GetVersionAtTimestamp(timestamp time.Time) (*Transaction, error) {
	history, err := vm.readHistory()
	if err != nil {
		return nil, err
	}

	// Find the version that was current at the given timestamp
	var latestTx *Transaction
	for _, tx := range history.Transactions {
		if tx.Timestamp.Before(timestamp) || tx.Timestamp.Equal(timestamp) {
			if latestTx == nil || tx.Timestamp.After(latestTx.Timestamp) {
				latestTx = tx
			}
		}
	}

	if latestTx == nil {
		return nil, fmt.Errorf("no version found at or before timestamp %s", timestamp)
	}

	return latestTx, nil
}

// readHistory reads the transaction history from file
func (vm *VersionManager) readHistory() (*TransactionHistory, error) {
	historyPath := filepath.Join(vm.metaDir, "history.json")
	
	// Check if file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return nil, err
	}
	
	// Read file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction history: %w", err)
	}
	
	// Parse JSON
	var history TransactionHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse transaction history: %w", err)
	}
	
	return &history, nil
}

// writeHistory writes the transaction history to file
func (vm *VersionManager) writeHistory(history *TransactionHistory) error {
	historyPath := filepath.Join(vm.metaDir, "history.json")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(vm.metaDir, 0755); err != nil {
		return fmt.Errorf("failed to create transaction history directory: %w", err)
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transaction history: %w", err)
	}
	
	// Write file
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction history: %w", err)
	}
	
	return nil
}

// Helper function to get a summary of a transaction
func (tx *Transaction) GetSummary() string {
	summary := fmt.Sprintf("Version %d: %s (%s)", tx.Version, tx.Operation, tx.Timestamp.Format(time.RFC3339))
	
	if tx.CommitInfo != nil {
		if message, ok := tx.CommitInfo["message"]; ok {
			summary += fmt.Sprintf(" - %s", message)
		}
	}
	
	return summary
}
