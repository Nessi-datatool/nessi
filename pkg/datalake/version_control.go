// Package datalake provides Delta Lake table management functionality
package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nessi-dev/nessi/pkg/logging"
)

// VersionManager handles Delta Lake version control operations including
// transaction recording, version history management, and time travel capabilities.
type VersionManager struct {
	tablePath string
	mutex     sync.RWMutex
	cache     *TransactionHistory // In-memory cache of transaction history
}

// TransactionHistory represents the history of transactions in a Delta Lake table
type TransactionHistory struct {
	Transactions []*Transaction `json:"transactions"`
	CurrentIndex int            `json:"current_index"`
}

// Transaction represents a transaction in the Delta Lake log
type Transaction struct {
	ID             string                 `json:"id"`
	Version        int                    `json:"version"`
	Timestamp      int64                  `json:"timestamp"` // Milliseconds since epoch
	Operation      string                 `json:"operation"`
	CommitInfo     map[string]string      `json:"commit_info"`
	AddedFiles     []string               `json:"added_files,omitempty"`
	RemovedFiles   []string               `json:"removed_files,omitempty"`
	MetadataChange *MetadataChange        `json:"metadata_change,omitempty"`
	Stats          map[string]interface{} `json:"stats,omitempty"`
	IsolationLevel string                 `json:"isolation_level,omitempty"`
	ReadVersion    int                    `json:"read_version,omitempty"`
	UserID         string                 `json:"user_id,omitempty"`
	ClientInfo     map[string]string      `json:"client_info,omitempty"`
}

// MetadataChange represents a change to table metadata
type MetadataChange struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Format      string            `json:"format,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	SchemaID    string            `json:"schema_id,omitempty"`
}

// TransactionSummary represents a summary of transactions between two versions
type TransactionSummary struct {
	FromVersion   int            `json:"from_version"`
	ToVersion     int            `json:"to_version"`
	AddedFiles    int            `json:"added_files"`
	RemovedFiles  int            `json:"removed_files"`
	ModifiedRows  int            `json:"modified_rows"`
	SchemaChanges []SchemaChange `json:"schema_changes"`
	Operations    []string       `json:"operations"`
	TimeSpan      int64          `json:"time_span_ms"` // Time span in milliseconds
	FromTimestamp int64          `json:"from_timestamp"`
	ToTimestamp   int64          `json:"to_timestamp"`
}

// NewVersionManager creates a new version manager for the specified Delta Lake table path.
// It initializes the manager with an empty cache that will be populated on first use.
func NewVersionManager(tablePath string) *VersionManager {
	return &VersionManager{
		tablePath: tablePath,
		mutex:     sync.RWMutex{},
		cache:     nil,
	}
}

// RecordTransaction records a new transaction in the Delta Lake log.
// It creates a new transaction entry with the provided details and appends it to the transaction history.
// The transaction is assigned a new version number and timestamp.
// Returns the created transaction and any error that occurred.
func (vm *VersionManager) RecordTransaction(
	operation string,
	commitInfo map[string]string,
	addedFiles []string,
	removedFiles []string,
	metadataChange *MetadataChange,
	stats map[string]interface{},
) (*Transaction, error) {
	if operation == "" {
		return nil, fmt.Errorf("operation cannot be empty")
	}

	logger := logging.GetLogger()
	logger.WithFields(map[string]interface{}{
		"tablePath":    vm.tablePath,
		"operation":    operation,
		"addedFiles":   len(addedFiles),
		"removedFiles": len(removedFiles),
	}).Info("Recording transaction")

	// Get current transaction history
	history, err := vm.getTransactionHistory()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	// Determine next version
	version := 0
	if len(history.Transactions) > 0 {
		version = history.Transactions[len(history.Transactions)-1].Version + 1
	}

	// Create new transaction
	transaction := &Transaction{
		ID:             uuid.New().String(),
		Version:        version,
		Timestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		Operation:      operation,
		CommitInfo:     commitInfo,
		AddedFiles:     addedFiles,
		RemovedFiles:   removedFiles,
		MetadataChange: metadataChange,
		Stats:          stats,
	}

	// Add transaction to history
	history.Transactions = append(history.Transactions, transaction)
	history.CurrentIndex = len(history.Transactions) - 1

	// Save transaction history
	err = vm.saveTransactionHistory(history)
	if err != nil {
		return nil, fmt.Errorf("failed to save transaction history: %w", err)
	}

	// Use the existing logger from the beginning of the function
	logging.GetLogger().WithFields(map[string]interface{}{
		"tablePath": vm.tablePath,
		"version":   version,
		"txId":      transaction.ID,
	}).Info("Transaction recorded successfully")

	return transaction, nil
}

// GetTransaction gets a transaction by version
func (vm *VersionManager) GetTransaction(version int) (*Transaction, error) {
	// Get transaction history
	history, err := vm.getTransactionHistory()
	if err != nil {
		return nil, err
	}

	// Find transaction with the specified version
	for _, tx := range history.Transactions {
		if tx.Version == version {
			return tx, nil
		}
	}

	return nil, fmt.Errorf("transaction with version %d not found", version)
}

// GetLatestTransaction gets the latest transaction
func (vm *VersionManager) GetLatestTransaction() (*Transaction, error) {
	// Get transaction history
	history, err := vm.getTransactionHistory()
	if err != nil {
		return nil, err
	}

	if len(history.Transactions) == 0 {
		return nil, fmt.Errorf("no transactions found")
	}

	return history.Transactions[len(history.Transactions)-1], nil
}

// GetTransactionSummary gets a summary of transactions between two versions
func (vm *VersionManager) GetTransactionSummary(fromVersion, toVersion int) (*TransactionSummary, error) {
	// Get transaction history
	history, err := vm.getTransactionHistory()
	if err != nil {
		return nil, err
	}

	// Find transactions with the specified versions
	var fromIndex, toIndex int
	var fromTx, toTx *Transaction
	for i, tx := range history.Transactions {
		if tx.Version == fromVersion {
			fromIndex = i
			fromTx = tx
		}
		if tx.Version == toVersion {
			toIndex = i
			toTx = tx
		}
	}

	if fromTx == nil {
		return nil, fmt.Errorf("transaction with version %d not found", fromVersion)
	}
	if toTx == nil {
		return nil, fmt.Errorf("transaction with version %d not found", toVersion)
	}

	// Create summary
	summary := &TransactionSummary{
		FromVersion:   fromVersion,
		ToVersion:     toVersion,
		AddedFiles:    0,
		RemovedFiles:  0,
		ModifiedRows:  0,
		SchemaChanges: []SchemaChange{},
		Operations:    []string{},
		TimeSpan:      toTx.Timestamp - fromTx.Timestamp, // Time span in milliseconds
		FromTimestamp: fromTx.Timestamp,
		ToTimestamp:   toTx.Timestamp,
	}

	// Collect operations
	for i := fromIndex + 1; i <= toIndex; i++ {
		tx := history.Transactions[i]
		summary.Operations = append(summary.Operations, tx.Operation)
		summary.AddedFiles += len(tx.AddedFiles)
		summary.RemovedFiles += len(tx.RemovedFiles)

		// Collect stats
		if tx.Stats != nil {
			if numRecords, ok := tx.Stats["numRecords"].(float64); ok {
				summary.ModifiedRows += int(numRecords)
			}
		}
	}

	return summary, nil
}

// GetTransactionAtTimestamp gets the transaction that was current at the specified timestamp
func (vm *VersionManager) GetTransactionAtTimestamp(timestamp time.Time) (*Transaction, error) {
	// Get transaction history
	history, err := vm.getTransactionHistory()
	if err != nil {
		return nil, err
	}

	if len(history.Transactions) == 0 {
		return nil, fmt.Errorf("no transactions found")
	}

	// Convert timestamp to milliseconds since epoch
	timestampMillis := timestamp.UnixNano() / int64(time.Millisecond)

	// Find the latest transaction that was created before or at the specified timestamp
	var latestTx *Transaction
	for _, tx := range history.Transactions {
		if tx.Timestamp <= timestampMillis {
			latestTx = tx
		} else {
			break
		}
	}

	if latestTx == nil {
		// If no transaction was found, return the earliest transaction
		return history.Transactions[0], nil
	}

	return latestTx, nil
}

// GetTransactionsBetween gets transactions between two timestamps
func (vm *VersionManager) GetTransactionsBetween(fromTimestamp, toTimestamp time.Time) ([]*Transaction, error) {
	// Get transaction history
	history, err := vm.getTransactionHistory()
	if err != nil {
		return nil, err
	}

	// Convert timestamps to milliseconds since epoch
	fromMillis := fromTimestamp.UnixNano() / int64(time.Millisecond)
	toMillis := toTimestamp.UnixNano() / int64(time.Millisecond)

	// Find transactions between the specified timestamps
	var transactions []*Transaction
	for _, tx := range history.Transactions {
		if tx.Timestamp >= fromMillis && tx.Timestamp <= toMillis {
			transactions = append(transactions, tx)
		}
	}

	return transactions, nil
}

// GetSummary gets a summary of a transaction
func (tx *Transaction) GetSummary() string {
	// Convert timestamp from milliseconds to time.Time
	timestamp := time.Unix(0, tx.Timestamp*int64(time.Millisecond))

	summary := fmt.Sprintf("Version %d (%s) - %s\n", tx.Version+1, timestamp.Format(time.RFC3339), tx.Operation)
	if tx.CommitInfo != nil {
		if message, ok := tx.CommitInfo["message"]; ok {
			summary += fmt.Sprintf("Message: %s\n", message)
		}
	}
	summary += fmt.Sprintf("Added files: %d, Removed files: %d\n", len(tx.AddedFiles), len(tx.RemovedFiles))
	return summary
}

// getTransactionHistory gets the transaction history from disk or cache.
// It uses a read-write mutex to ensure thread safety when accessing the cache.
func (vm *VersionManager) getTransactionHistory() (*TransactionHistory, error) {
	// Check cache first with a read lock
	vm.mutex.RLock()
	if vm.cache != nil {
		defer vm.mutex.RUnlock()
		return vm.cache, nil
	}
	vm.mutex.RUnlock()

	// Cache miss, acquire write lock
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Double-check cache after acquiring write lock
	if vm.cache != nil {
		return vm.cache, nil
	}

	// Create transaction history file path
	historyPath := filepath.Join(vm.tablePath, "_delta_log", "transaction_history.json")
	logger := logging.GetLogger()
	logger.WithField("path", historyPath).Debug("Reading transaction history")

	// Check if file exists
	_, err := os.Stat(historyPath)
	if os.IsNotExist(err) {
		// Create empty history
		vm.cache = &TransactionHistory{
			Transactions: []*Transaction{},
			CurrentIndex: -1,
		}
		return vm.cache, nil
	} else if err != nil {
		logger := logging.GetLogger()
		logger.WithField("path", historyPath).WithError(err).Error("Failed to stat transaction history file")
		return nil, fmt.Errorf("failed to access transaction history: %w", err)
	}

	// Read file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		logger := logging.GetLogger()
		logger.WithField("path", historyPath).WithError(err).Error("Failed to read transaction history file")
		return nil, fmt.Errorf("failed to read transaction history: %w", err)
	}

	// Parse JSON
	var history TransactionHistory
	err = json.Unmarshal(data, &history)
	if err != nil {
		logger := logging.GetLogger()
		logger.WithField("path", historyPath).WithError(err).Error("Failed to parse transaction history")
		return nil, fmt.Errorf("failed to parse transaction history: %w", err)
	}

	// Update cache
	vm.cache = &history
	return vm.cache, nil
}

// saveTransactionHistory saves the transaction history to disk and updates the cache.
// It acquires a write lock to ensure thread safety when updating the cache.
func (vm *VersionManager) saveTransactionHistory(history *TransactionHistory) error {
	vm.mutex.Lock()
	defer vm.mutex.Unlock()

	// Create transaction history file path
	historyPath := filepath.Join(vm.tablePath, "_delta_log", "transaction_history.json")
	logger := logging.GetLogger()
	logger.WithFields(map[string]interface{}{
		"path":     historyPath,
		"txCount":  len(history.Transactions),
		"curIndex": history.CurrentIndex,
	}).Debug("Saving transaction history")

	// Create directory if it doesn't exist
	logDir := filepath.Join(vm.tablePath, "_delta_log")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		logger := logging.GetLogger()
		logger.WithField("dir", logDir).WithError(err).Error("Failed to create _delta_log directory")
		return fmt.Errorf("failed to create transaction log directory: %w", err)
	}

	// Convert to JSON
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		logger := logging.GetLogger()
		logger.WithError(err).Error("Failed to marshal transaction history")
		return fmt.Errorf("failed to serialize transaction history: %w", err)
	}

	// Write file atomically by writing to a temporary file first and then renaming
	tempFile := historyPath + ".tmp"
	err = os.WriteFile(tempFile, data, 0644)
	if err != nil {
		logger := logging.GetLogger()
		logger.WithField("path", tempFile).WithError(err).Error("Failed to write transaction history")
		return fmt.Errorf("failed to write transaction history: %w", err)
	}

	// Rename temp file to actual file (atomic operation)
	err = os.Rename(tempFile, historyPath)
	if err != nil {
		logger := logging.GetLogger()
		logger.WithFields(map[string]interface{}{
			"from": tempFile,
			"to":   historyPath,
		}).WithError(err).Error("Failed to rename transaction history file")
		return fmt.Errorf("failed to finalize transaction history: %w", err)
	}

	// Update cache
	vm.cache = history
	return nil
}
