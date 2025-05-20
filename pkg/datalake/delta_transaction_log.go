package datalake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// TransactionLogEntry represents a single entry in the Delta Lake transaction log
type TransactionLogEntry struct {
	Version             int64                  `json:"version"`
	Timestamp           time.Time              `json:"timestamp"`
	Operation           string                 `json:"operation"`
	OperationParameters map[string]interface{} `json:"operationParameters"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	Add                 []AddFileEntry         `json:"add,omitempty"`
	Remove              []RemoveFileEntry      `json:"remove,omitempty"`
	SchemaChange        *SchemaChangeEntry     `json:"schemaChange,omitempty"`
}

// AddFileEntry represents a file added in a transaction
type AddFileEntry struct {
	Path             string            `json:"path"`
	Size             int64             `json:"size"`
	ModificationTime int64             `json:"modificationTime"`
	PartitionValues  map[string]string `json:"partitionValues,omitempty"`
	DataChange       bool              `json:"dataChange"`
	Stats            string            `json:"stats,omitempty"`
}

// RemoveFileEntry represents a file removed in a transaction
type RemoveFileEntry struct {
	Path              string            `json:"path"`
	Size              int64             `json:"size,omitempty"`
	DeletionTimestamp int64             `json:"deletionTimestamp"`
	PartitionValues   map[string]string `json:"partitionValues,omitempty"`
	DataChange        bool              `json:"dataChange"`
}

// SchemaChangeEntry represents a schema change in a transaction
type SchemaChangeEntry struct {
	Schema  map[string]interface{}   `json:"schema"`
	Fields  []map[string]interface{} `json:"fields,omitempty"`
	Changes []SchemaChange           `json:"changes,omitempty"`
}

// SchemaChange represents a specific change to the schema
type SchemaChange struct {
	Type      string                 `json:"type"` // add, remove, update, etc.
	FieldName string                 `json:"fieldName"`
	FieldType string                 `json:"fieldType,omitempty"`
	Nullable  bool                   `json:"nullable,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// DeltaTransactionLog provides transaction log parsing for Delta Lake tables
type DeltaTransactionLog struct {
	handler *DeltaFormatHandler
}

// NewDeltaTransactionLog creates a new DeltaTransactionLog instance
func NewDeltaTransactionLog(handler *DeltaFormatHandler) *DeltaTransactionLog {
	return &DeltaTransactionLog{
		handler: handler,
	}
}

// GetTransactionLog returns the transaction log for a Delta Lake table
func (t *DeltaTransactionLog) GetTransactionLog(path string) ([]TransactionLogEntry, error) {
	// Check if the path is a Delta Lake table
	if !t.handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Get the log directory
	logDir := filepath.Join(path, "_delta_log")

	// Read all log files
	logFiles, err := t.getLogFiles(logDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read log files: %w", err)
	}

	// Parse the log files
	var entries []TransactionLogEntry
	for _, logFile := range logFiles {
		fileEntries, err := t.parseLogFile(filepath.Join(logDir, logFile))
		if err != nil {
			return nil, fmt.Errorf("failed to parse log file %s: %w", logFile, err)
		}
		entries = append(entries, fileEntries...)
	}

	// Sort entries by version
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Version < entries[j].Version
	})

	return entries, nil
}

// getLogFiles returns a sorted list of log files in the given directory
func (t *DeltaTransactionLog) getLogFiles(logDir string) ([]string, error) {
	// Read the directory
	files, err := ioutil.ReadDir(logDir)
	if err != nil {
		return nil, err
	}

	// Filter and sort log files
	var logFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			logFiles = append(logFiles, file.Name())
		}
	}

	// Sort log files by version number
	sort.Slice(logFiles, func(i, j int) bool {
		vi, _ := strconv.ParseInt(strings.TrimSuffix(logFiles[i], ".json"), 10, 64)
		vj, _ := strconv.ParseInt(strings.TrimSuffix(logFiles[j], ".json"), 10, 64)
		return vi < vj
	})

	return logFiles, nil
}

// parseLogFile parses a single log file and returns the transaction entries
func (t *DeltaTransactionLog) parseLogFile(filePath string) ([]TransactionLogEntry, error) {
	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Split the file into lines
	lines := strings.Split(string(data), "\n")

	// Parse each line as a JSON object
	var entries []TransactionLogEntry
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var entry TransactionLogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to parse log entry: %w", err)
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// GetLatestVersion returns the latest version of the Delta Lake table
func (t *DeltaTransactionLog) GetLatestVersion(path string) (int64, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return 0, err
	}

	// Return the latest version
	if len(entries) == 0 {
		return 0, nil
	}

	return entries[len(entries)-1].Version, nil
}

// GetSchemaAtVersion returns the schema at the specified version
func (t *DeltaTransactionLog) GetSchemaAtVersion(path string, version int64) (map[string]interface{}, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return nil, err
	}

	// Find the schema at the specified version
	var schema map[string]interface{}
	for _, entry := range entries {
		if entry.Version > version {
			break
		}

		// Check if this entry has a schema change
		if entry.SchemaChange != nil && entry.SchemaChange.Schema != nil {
			schema = entry.SchemaChange.Schema
		}
	}

	if schema == nil {
		return nil, fmt.Errorf("no schema found at version %d", version)
	}

	return schema, nil
}

// GetSchemaChanges returns the schema changes between two versions
func (t *DeltaTransactionLog) GetSchemaChanges(path string, fromVersion, toVersion int64) ([]SchemaChange, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return nil, err
	}

	// Find schema changes between the specified versions
	var changes []SchemaChange
	for _, entry := range entries {
		if entry.Version > fromVersion && entry.Version <= toVersion {
			if entry.SchemaChange != nil && len(entry.SchemaChange.Changes) > 0 {
				changes = append(changes, entry.SchemaChange.Changes...)
			}
		}
	}

	return changes, nil
}

// GetTableHistory returns the history of operations on the table
func (t *DeltaTransactionLog) GetTableHistory(path string) ([]map[string]interface{}, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return nil, err
	}

	// Convert entries to a simplified history format
	history := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		history[i] = map[string]interface{}{
			"version":   entry.Version,
			"timestamp": entry.Timestamp,
			"operation": entry.Operation,
		}

		// Add operation-specific details
		switch entry.Operation {
		case "WRITE":
			history[i]["filesAdded"] = len(entry.Add)
			history[i]["filesRemoved"] = len(entry.Remove)
		case "COMMIT":
			history[i]["filesAdded"] = len(entry.Add)
			history[i]["filesRemoved"] = len(entry.Remove)
		case "SCHEMA_CHANGE":
			if entry.SchemaChange != nil {
				history[i]["schemaChanges"] = len(entry.SchemaChange.Changes)
			}
		}
	}

	return history, nil
}

// GetFilesAtVersion returns the list of files at the specified version
func (t *DeltaTransactionLog) GetFilesAtVersion(path string, version int64) ([]string, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return nil, err
	}

	// Track files added and removed up to the specified version
	files := make(map[string]bool)
	for _, entry := range entries {
		if entry.Version > version {
			break
		}

		// Add files
		for _, add := range entry.Add {
			files[add.Path] = true
		}

		// Remove files
		for _, remove := range entry.Remove {
			delete(files, remove.Path)
		}
	}

	// Convert map to slice
	var result []string
	for file := range files {
		result = append(result, file)
	}

	// Sort the result for deterministic output
	sort.Strings(result)

	return result, nil
}

// GetVersionTimestamp returns the timestamp of a specific version
func (t *DeltaTransactionLog) GetVersionTimestamp(path string, version int64) (time.Time, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return time.Time{}, err
	}

	// Find the entry with the specified version
	for _, entry := range entries {
		if entry.Version == version {
			return entry.Timestamp, nil
		}
	}

	return time.Time{}, fmt.Errorf("version %d not found", version)
}

// GetVersionByTimestamp returns the version closest to the specified timestamp
func (t *DeltaTransactionLog) GetVersionByTimestamp(path string, timestamp time.Time) (int64, error) {
	// Get the transaction log
	entries, err := t.GetTransactionLog(path)
	if err != nil {
		return 0, err
	}

	if len(entries) == 0 {
		return 0, fmt.Errorf("no versions found")
	}

	// Find the version closest to the timestamp
	closestVersion := entries[0].Version
	closestDiff := timestamp.Sub(entries[0].Timestamp)
	if closestDiff < 0 {
		closestDiff = -closestDiff
	}

	for _, entry := range entries[1:] {
		diff := timestamp.Sub(entry.Timestamp)
		if diff < 0 {
			diff = -diff
		}

		if diff < closestDiff {
			closestVersion = entry.Version
			closestDiff = diff
		}
	}

	return closestVersion, nil
}

// GetDataFilesForVersion returns the data files for a specific version of a Delta Lake table
func (t *DeltaTransactionLog) GetDataFilesForVersion(path string, version int64) ([]string, error) {
	// This is an alias for GetFilesAtVersion for backward compatibility
	return t.GetFilesAtVersion(path, version)
}
