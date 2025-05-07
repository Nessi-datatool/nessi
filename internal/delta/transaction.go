package delta

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Transaction represents a Delta Lake transaction
type Transaction struct {
	Actions    []Action
	ReadVersion int64
	CommitInfo *CommitInfo
}

// CommitInfo represents information about a commit
type CommitInfo struct {
	Timestamp int64             `json:"timestamp"`
	UserID    string            `json:"userId"`
	UserName  string            `json:"userName"`
	Operation string            `json:"operation"`
	OperationParameters map[string]string `json:"operationParameters"`
	JobInfo   map[string]string `json:"jobInfo"`
	Notebook  map[string]string `json:"notebook"`
	ClusterID string            `json:"clusterId"`
	ReadVersion int64           `json:"readVersion"`
	IsolationLevel string       `json:"isolationLevel"`
	IsBlindAppend bool          `json:"isBlindAppend"`
}

// TransactionLog manages Delta Lake transactions
type TransactionLog struct {
	tablePath string
	mu        sync.RWMutex
	version   int64
}

// NewTransactionLog creates a new transaction log
func NewTransactionLog(tablePath string) *TransactionLog {
	return &TransactionLog{
		tablePath: tablePath,
		version:   0,
	}
}

// BeginTransaction begins a new transaction
func (t *TransactionLog) BeginTransaction() *Transaction {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return &Transaction{
		Actions:     make([]Action, 0),
		ReadVersion: t.version,
		CommitInfo: &CommitInfo{
			Timestamp: time.Now().UnixMilli(),
			IsolationLevel: "Serializable",
			IsBlindAppend: false,
		},
	}
}

// AddAction adds an action to the transaction
func (tx *Transaction) AddAction(action Action) error {
	if err := action.Validate(); err != nil {
		return fmt.Errorf("invalid action: %w", err)
	}
	tx.Actions = append(tx.Actions, action)
	return nil
}

// Commit commits the transaction
func (t *TransactionLog) Commit(ctx context.Context, tx *Transaction) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check if the read version is still valid
	if tx.ReadVersion != t.version {
		return fmt.Errorf("concurrent modification detected: expected version %d, got %d", t.version, tx.ReadVersion)
	}

	// Create the commit file
	commitPath := filepath.Join(t.tablePath, "_delta_log", fmt.Sprintf("%d.json", t.version+1))
	
	// Write the transaction to the log
	data, err := json.Marshal(tx.Actions)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	if err := os.WriteFile(commitPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	// Write the commit info
	commitInfoPath := filepath.Join(t.tablePath, "_delta_log", fmt.Sprintf("%d.commit.json", t.version+1))
	commitInfoData, err := json.Marshal(tx.CommitInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal commit info: %w", err)
	}

	if err := os.WriteFile(commitInfoPath, commitInfoData, 0644); err != nil {
		// Clean up the commit file if writing commit info fails
		os.Remove(commitPath)
		return fmt.Errorf("failed to write commit info: %w", err)
	}

	// Update the version
	t.version++

	return nil
}

// GetVersion returns the current version
func (t *TransactionLog) GetVersion() int64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.version
}

// GetActions returns all actions up to a specific version
func (t *TransactionLog) GetActions(version int64) ([]Action, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if version < 0 || version > t.version {
		return nil, fmt.Errorf("invalid version: %d", version)
	}

	var actions []Action
	for v := int64(0); v <= version; v++ {
		filePath := filepath.Join(t.tablePath, "_delta_log", fmt.Sprintf("%d.json", v))
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read log file: %w", err)
		}

		var versionActions []Action
		if err := json.Unmarshal(data, &versionActions); err != nil {
			return nil, fmt.Errorf("failed to parse log file: %w", err)
		}
		actions = append(actions, versionActions...)
	}

	return actions, nil
}

// GetCommitInfo returns the commit info for a specific version
func (t *TransactionLog) GetCommitInfo(version int64) (*CommitInfo, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if version < 0 || version > t.version {
		return nil, fmt.Errorf("invalid version: %d", version)
	}

	filePath := filepath.Join(t.tablePath, "_delta_log", fmt.Sprintf("%d.commit.json", version))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit info: %w", err)
	}

	var commitInfo CommitInfo
	if err := json.Unmarshal(data, &commitInfo); err != nil {
		return nil, fmt.Errorf("failed to parse commit info: %w", err)
	}

	return &commitInfo, nil
}

// ValidateTransaction validates a transaction
func (t *TransactionLog) ValidateTransaction(tx *Transaction) error {
	if len(tx.Actions) == 0 {
		return fmt.Errorf("transaction must have at least one action")
	}

	for _, action := range tx.Actions {
		if err := action.Validate(); err != nil {
			return fmt.Errorf("invalid action: %w", err)
		}
	}

	return nil
}

// Rollback rolls back a transaction
func (t *TransactionLog) Rollback(tx *Transaction) {
	tx.Actions = nil
	tx.CommitInfo = nil
}

// SetCommitInfo sets the commit info for a transaction
func (tx *Transaction) SetCommitInfo(info *CommitInfo) {
	tx.CommitInfo = info
}

// GetReadVersion returns the read version of the transaction
func (tx *Transaction) GetReadVersion() int64 {
	return tx.ReadVersion
}

// GetActions returns the actions in the transaction
func (tx *Transaction) GetActions() []Action {
	return tx.Actions
}

// GetCommitInfo returns the commit info of the transaction
func (tx *Transaction) GetCommitInfo() *CommitInfo {
	return tx.CommitInfo
} 