package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// DeltaTable represents a Delta table
type DeltaTable struct {
	Path            string
	Version         int64
	LastModified    time.Time
	PartitionSchema []string
	Schema          map[string]string
	Stats           *TableStats
	LogFiles        []string
	Checkpoint      *Checkpoint
}

// TableStats represents statistics about a Delta table
type TableStats struct {
	NumFiles        int64
	NumRecords      int64
	TotalSize       int64
	PartitionCounts map[string]int64
	ColumnStats     map[string]*ColumnStats
}

// ColumnStats represents statistics for a single column
type ColumnStats struct {
	NullCount     int64
	DistinctCount int64
	MinValue      interface{}
	MaxValue      interface{}
	AvgValue      float64
}

// Checkpoint represents a Delta table checkpoint
type Checkpoint struct {
	Version    int64
	Timestamp  time.Time
	FileCount  int64
	FilePaths  []string
	Schema     map[string]string
	Partitions []string
}

// Reader handles reading from Delta tables
type Reader struct {
	table  *DeltaTable
	files  []*os.File
	closed bool
	mu     sync.Mutex
}

// NewReader creates a new Delta table reader
func NewReader(tablePath string) (*Reader, error) {
	// Verify table path exists
	if _, err := os.Stat(tablePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("table path does not exist: %s", tablePath)
	}

	// Create reader
	reader := &Reader{
		table: &DeltaTable{
			Path:         tablePath,
			Schema:       make(map[string]string),
			Stats:        &TableStats{ColumnStats: make(map[string]*ColumnStats)},
			LastModified: time.Now(),
		},
		files: make([]*os.File, 0),
	}

	// Initialize reader
	if err := reader.initialize(); err != nil {
		reader.Close() // Clean up any resources
		return nil, fmt.Errorf("failed to initialize reader: %w", err)
	}

	return reader, nil
}

// initialize sets up the reader by reading table metadata
func (r *Reader) initialize() error {
	// Read _delta_log directory
	logDir := filepath.Join(r.table.Path, "_delta_log")
	if err := r.readLogDirectory(logDir); err != nil {
		return fmt.Errorf("failed to read log directory: %w", err)
	}

	// Read latest log file first to catch any corrupted logs
	if err := r.readLatestLog(); err != nil {
		return fmt.Errorf("failed to read latest log: %w", err)
	}

	// Read latest checkpoint
	if err := r.readLatestCheckpoint(); err != nil {
		return fmt.Errorf("failed to read checkpoint: %w", err)
	}

	// Update table metadata
	r.table.Version = r.table.Checkpoint.Version
	r.table.Schema = r.table.Checkpoint.Schema
	r.table.PartitionSchema = r.table.Checkpoint.Partitions
	r.table.LastModified = r.table.Checkpoint.Timestamp

	// Compute table statistics
	if err := r.computeTableStats(); err != nil {
		return fmt.Errorf("failed to compute table stats: %w", err)
	}

	return nil
}

// readLogDirectory reads the _delta_log directory
func (r *Reader) readLogDirectory(logDir string) error {
	entries, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("invalid Delta table: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			// Collect both JSON log files and checkpoint parquet files
			if filepath.Ext(entry.Name()) == ".json" || strings.HasSuffix(entry.Name(), ".checkpoint.parquet") {
				r.table.LogFiles = append(r.table.LogFiles, entry.Name())
			}
		}
	}

	return nil
}

// readLatestCheckpoint reads the latest checkpoint file
func (r *Reader) readLatestCheckpoint() error {
	// Find latest checkpoint file
	var checkpointFile string
	for _, file := range r.table.LogFiles {
		if strings.HasSuffix(filepath.Base(file), ".checkpoint.parquet") {
			checkpointFile = file
			break
		}
	}

	if checkpointFile == "" {
		return fmt.Errorf("no checkpoint file found")
	}

	// Read checkpoint file
	file, err := os.Open(filepath.Join(r.table.Path, "_delta_log", checkpointFile))
	if err != nil {
		return fmt.Errorf("failed to open checkpoint file: %w", err)
	}
	defer file.Close()

	// Parse checkpoint data
	var checkpoint Checkpoint
	if err := json.NewDecoder(file).Decode(&checkpoint); err != nil {
		return fmt.Errorf("failed to parse checkpoint: %w", err)
	}

	r.table.Checkpoint = &checkpoint
	return nil
}

// readLatestLog reads the latest log file
func (r *Reader) readLatestLog() error {
	// Find latest log file
	var latestLog string
	for _, file := range r.table.LogFiles {
		if filepath.Ext(file) == ".json" && file > latestLog {
			latestLog = file
		}
	}

	if latestLog == "" {
		return fmt.Errorf("no log files found")
	}

	// Read log file
	file, err := os.Open(filepath.Join(r.table.Path, "_delta_log", latestLog))
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Parse log data
	var logData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&logData); err != nil {
		return fmt.Errorf("invalid log entry: %w", err)
	}

	// Update table metadata from log
	if version, ok := logData["version"].(float64); ok {
		r.table.Version = int64(version)
	}
	if timestamp, ok := logData["timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
			r.table.LastModified = t
		}
	}

	return nil
}

// computeTableStats computes table statistics
func (r *Reader) computeTableStats() error {
	// Count files and compute sizes
	var totalSize int64
	partitionCounts := make(map[string]int64)

	for _, filePath := range r.table.Checkpoint.FilePaths {
		info, err := os.Stat(filepath.Join(r.table.Path, filePath))
		if err != nil {
			continue // Skip files that can't be accessed
		}

		totalSize += info.Size()
		r.table.Stats.NumFiles++

		// Update partition counts
		if partition := r.getPartitionFromPath(filePath); partition != "" {
			partitionCounts[partition]++
		}
	}

	r.table.Stats.TotalSize = totalSize
	r.table.Stats.PartitionCounts = partitionCounts

	// TODO: Implement column statistics computation
	// This would involve:
	// 1. Reading sample data from files
	// 2. Computing null counts, distinct values, etc.
	// 3. Updating ColumnStats for each column

	return nil
}

// getPartitionFromPath extracts partition information from a file path
func (r *Reader) getPartitionFromPath(path string) string {
	// TODO: Implement partition extraction from path
	// This would involve:
	// 1. Parsing the path for partition keys
	// 2. Returning the partition string
	return ""
}

// GetSchema returns the table schema
func (r *Reader) GetSchema() map[string]string {
	return r.table.Schema
}

// GetPartitionSchema returns the partition schema
func (r *Reader) GetPartitionSchema() []string {
	return r.table.PartitionSchema
}

// GetStats returns the table statistics
func (r *Reader) GetStats() *TableStats {
	return r.table.Stats
}

// GetVersion returns the current table version
func (r *Reader) GetVersion() int64 {
	if r.closed {
		// Instead of panicking, we'll return 0, but in a real implementation
		// you might want to return an error or properly handle this case
		return 0
	}

	if r.table == nil {
		return 0
	}

	return r.table.Version
}

// GetLastModified returns the last modification time
func (r *Reader) GetLastModified() time.Time {
	return r.table.LastModified
}

// GetTablePath returns the path of the table being read.
func (r *Reader) GetTablePath() string {
	// Ensure table is not nil before accessing Path, though in normal operation it shouldn't be.
	if r.table != nil {
		return r.table.Path
	}
	return ""
}

// ReadPartition reads data from a specific partition
func (r *Reader) ReadPartition(partition string) (io.ReadCloser, error) {
	if r.closed {
		return nil, fmt.Errorf("reader closed")
	}
	// Find files in partition
	var partitionFiles []string
	for _, filePath := range r.table.Checkpoint.FilePaths {
		if r.getPartitionFromPath(filePath) == partition {
			partitionFiles = append(partitionFiles, filePath)
		}
	}

	if len(partitionFiles) == 0 {
		return nil, fmt.Errorf("partition not found: %s", partition)
	}

	// Create multi-reader for partition files
	readers := make([]io.Reader, 0, len(partitionFiles))
	for _, filePath := range partitionFiles {
		file, err := os.Open(filepath.Join(r.table.Path, filePath))
		if err != nil {
			// Close any already opened files
			for _, f := range r.files {
				f.Close()
			}
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		r.files = append(r.files, file)
		readers = append(readers, file)
	}

	return io.NopCloser(io.MultiReader(readers...)), nil
}

// ReadAll reads all data from the table
func (r *Reader) ReadAll() (io.ReadCloser, error) {
	if r.closed {
		return nil, fmt.Errorf("reader closed")
	}
	// Create multi-reader for all files
	readers := make([]io.Reader, 0, len(r.table.Checkpoint.FilePaths))
	for _, filePath := range r.table.Checkpoint.FilePaths {
		file, err := os.Open(filepath.Join(r.table.Path, filePath))
		if err != nil {
			// Close any already opened files
			for _, f := range r.files {
				f.Close()
			}
			return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
		}
		r.files = append(r.files, file)
		readers = append(readers, file)
	}

	return io.NopCloser(io.MultiReader(readers...)), nil
}

// ListPartitions returns all available partitions
func (r *Reader) ListPartitions() ([]string, error) {
	partitions := make([]string, 0, len(r.table.Stats.PartitionCounts))
	for partition := range r.table.Stats.PartitionCounts {
		partitions = append(partitions, partition)
	}
	return partitions, nil
}

// GetPartitionStats returns statistics for a specific partition
func (r *Reader) GetPartitionStats(partition string) (int64, error) {
	count, exists := r.table.Stats.PartitionCounts[partition]
	if !exists {
		return 0, fmt.Errorf("partition %s not found", partition)
	}
	return count, nil
}

// ReadAllParsed reads all data from the table and parses it into a structured format
func (r *Reader) ReadAllParsed() ([]map[string]interface{}, error) {
	rawData, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	defer rawData.Close()

	// Parse the data into structured records
	records := make([]map[string]interface{}, 0)
	scanner := bufio.NewScanner(rawData)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Parse the line as JSON
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, fmt.Errorf("failed to parse line as JSON: %w", err)
		}

		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading data: %w", err)
	}

	return records, nil
}

// Close closes the reader and releases resources
func (r *Reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil // Already closed
	}

	var lastErr error
	for _, file := range r.files {
		if err := file.Close(); err != nil {
			lastErr = err
		}
	}
	r.files = nil
	r.closed = true
	return lastErr
}
