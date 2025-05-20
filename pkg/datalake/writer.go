package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
)

// Writer handles writing to Delta tables
type Writer struct {
	table     *DeltaTable
	closed    bool
	mu        sync.Mutex
	fileCache map[string]*os.File
}

// NewWriter creates a new Delta table writer
func NewWriter(tablePath string, schema *arrow.Schema) (*Writer, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	return &Writer{
		table: &DeltaTable{
			Path:            tablePath,
			Version:         0,
			LastModified:    time.Now(),
			Schema:          schema,
			PartitionSchema: []string{},
			Stats:           &TableStats{},
			Files:           []string{},
			Partitions:      make(map[string][]string),
			Metadata:        make(map[string]interface{}),
		},
		closed:    false,
		fileCache: make(map[string]*os.File),
	}, nil
}

// Initialize sets up the writer by creating necessary directories
func (w *Writer) Initialize() error {
	if w.closed {
		return fmt.Errorf("writer is closed")
	}

	// Create table directory if it doesn't exist
	if err := os.MkdirAll(w.table.Path, 0755); err != nil {
		return fmt.Errorf("failed to create table directory: %w", err)
	}

	// Create _delta_log directory
	logDir := filepath.Join(w.table.Path, "_delta_log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create _delta_log directory: %w", err)
	}

	return nil
}

// WritePartition writes data to a specific partition
func (w *Writer) WritePartition(partition string, record arrow.Record) error {
	if w.closed {
		return fmt.Errorf("writer is closed")
	}

	if record == nil {
		return fmt.Errorf("record cannot be nil")
	}

	// Validate schema
	validator := NewSchemaValidator(w.table.Schema)
	errors := validator.ValidateRecord(record)
	if len(errors) > 0 {
		// Log the validation errors
		formattedErrors := FormatValidationErrors(errors)
		return fmt.Errorf("schema validation failed: %s", formattedErrors)
	}

	// Create partition directory
	partitionPath := filepath.Join(w.table.Path, partition)
	if err := os.MkdirAll(partitionPath, 0755); err != nil {
		return fmt.Errorf("failed to create partition directory: %w", err)
	}

	// Generate parquet file name
	fileName := fmt.Sprintf("part-%d-%s.parquet", time.Now().UnixNano(), partition)
	filePath := filepath.Join(partitionPath, fileName)

	// Add file to table files
	w.table.Files = append(w.table.Files, filePath)

	// Update partition mapping
	if _, exists := w.table.Partitions[partition]; !exists {
		w.table.Partitions[partition] = []string{}
	}
	w.table.Partitions[partition] = append(w.table.Partitions[partition], filePath)

	// Update table stats
	w.table.Stats.NumFiles++
	w.table.Stats.NumRecords += record.NumRows()

	return nil
}

// Commit commits the changes and updates the transaction log
func (w *Writer) Commit() error {
	if w.closed {
		return fmt.Errorf("writer is closed")
	}

	// Update version and last modified time
	w.table.Version++
	w.table.LastModified = time.Now()

	// Create transaction log entry
	logEntry := map[string]interface{}{
		"version":    w.table.Version,
		"timestamp":  w.table.LastModified.Unix(),
		"operation":  "WRITE",
		"files":      w.table.Files,
		"partitions": w.table.Partitions,
		"stats":      w.table.Stats,
		"schema":     w.table.Schema,
		"metadata":   w.table.Metadata,
	}

	// Convert to JSON
	data, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Write transaction log
	logFile := filepath.Join(w.table.Path, "_delta_log", fmt.Sprintf("%020d.json", w.table.Version))
	if err := os.WriteFile(logFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	return nil
}

// Close closes the writer
func (w *Writer) Close() error {
	w.closed = true
	return nil
}

// IsClosed returns whether the writer is closed
func (w *Writer) IsClosed() bool {
	return w.closed
}

// GetSchema returns the table schema
func (w *Writer) GetSchema() *arrow.Schema {
	return w.table.Schema
}

// GetStats returns the table statistics
func (w *Writer) GetStats() *TableStats {
	return w.table.Stats
}
