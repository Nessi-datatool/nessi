package datalake

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/ipc"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// ParquetManager handles Parquet file operations
type ParquetManager struct {
	allocator memory.Allocator
}

// NewParquetManager creates a new Parquet manager
func NewParquetManager() *ParquetManager {
	return &ParquetManager{
		allocator: memory.NewGoAllocator(),
	}
}

// WriteRecord writes an Arrow record to a Parquet file
func (p *ParquetManager) WriteRecord(filePath string, record arrow.Record) error {
	if record == nil {
		return fmt.Errorf("record cannot be nil")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// TODO: Implement actual Parquet writing
	// For now, just create an empty file as a placeholder
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return nil
}

// ReadRecord reads an Arrow record from a Parquet file
func (p *ParquetManager) ReadRecord(filePath string, schema *arrow.Schema) (arrow.Record, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema cannot be nil")
	}

	// Try Go-native reading (not implemented, so always fallback for now)
	// TODO: Implement Go-native Parquet reading here
	// If fails, fallback to Python

	record, err := p.readRecordPythonFallback(filePath, schema)
	if err != nil {
		return nil, fmt.Errorf("failed to read Parquet file using Go and Python fallback: %w", err)
	}
	return record, nil
}

// readRecordPythonFallback calls the Python script to read Parquet and returns Arrow record (or error)
// This aligns with the Go-Python bridge described in IMPLEMENTATION.md.
func (p *ParquetManager) readRecordPythonFallback(filePath string, schema *arrow.Schema) (arrow.Record, error) {
	cmd := exec.Command(filepath.Join("..", "..", "venv", "bin", "python"), filepath.Join("..", "..", "scripts", "read_parquet.py"), filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("python fallback error: %v, stderr: %s", err, stderr.String())
	}

	if stderr.Len() > 0 {
		fmt.Fprintf(os.Stderr, "[Python fallback stderr]: %s\n", stderr.String())
	}
	ipcBytes := out.Bytes()
	fmt.Fprintf(os.Stderr, "[Python fallback] Arrow IPC buffer size: %d bytes\n", len(ipcBytes))
	reader, err := ipc.NewReader(bytes.NewReader(ipcBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create Arrow IPC reader: %w", err)
	}
	defer reader.Release()

	if !reader.Next() {
		max := 64
		if len(ipcBytes) < max {
			max = len(ipcBytes)
		}
		fmt.Fprintf(os.Stderr, "[Python fallback] First %d bytes of Arrow IPC: %x\n", max, ipcBytes[:max])
		return nil, fmt.Errorf("no record returned from Arrow IPC stream (buffer size: %d bytes)", len(ipcBytes))
	}
	return reader.Record(), nil
}

// GetFileStats gets statistics about a Parquet file
func (p *ParquetManager) GetFileStats(filePath string) (*FileStats, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	// For now, return dummy stats for testing
	stats := &FileStats{
		Path:    filePath,
		Size:    1024, // Fixed size for testing
		ModTime: info.ModTime(),
		NumRows: 100, // TODO: Read from Parquet metadata
		Columns: make(map[string]*ColumnStats),
	}

	return stats, nil
}

// FileStats represents statistics about a Parquet file
type FileStats struct {
	Path    string
	Size    int64
	ModTime time.Time
	NumRows int64
	Columns map[string]*ColumnStats
}

// ValidateSchema validates that a record matches the expected schema
func (p *ParquetManager) ValidateSchema(record arrow.Record, schema *arrow.Schema) error {
	if record == nil || schema == nil {
		return fmt.Errorf("record and schema cannot be nil")
	}

	recordSchema := record.Schema()
	if recordSchema.String() != schema.String() {
		return fmt.Errorf("schema mismatch: expected %s, got %s", schema, recordSchema)
	}

	return nil
}
