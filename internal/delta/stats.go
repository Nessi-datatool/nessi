package delta

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/memory"
	parquetfile "github.com/apache/arrow/go/v15/parquet/file"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
)

// StatsManager manages Delta Lake statistics
type StatsManager struct {
	tablePath string
	schema    *arrow.Schema
}

// NewStatsManager creates a new stats manager
func NewStatsManager(tablePath string, schema *arrow.Schema) *StatsManager {
	return &StatsManager{
		tablePath: tablePath,
		schema:    schema,
	}
}

// FileStats represents statistics for a file
type FileStats struct {
	Path      string            `json:"path"`
	Size      int64             `json:"size"`
	NumRows   int64             `json:"numRows"`
	MinValues map[string]string `json:"minValues,omitempty"`
	MaxValues map[string]string `json:"maxValues,omitempty"`
	NullCount map[string]int64  `json:"nullCount,omitempty"`
}

// ColumnStats represents statistics for a column
type ColumnStats struct {
	Name      string `json:"name"`
	MinValue  string `json:"minValue,omitempty"`
	MaxValue  string `json:"maxValue,omitempty"`
	NullCount int64  `json:"nullCount"`
	Distinct  int64  `json:"distinct,omitempty"`
}

// GetFileStats returns statistics for a file
func (s *StatsManager) GetFileStats(ctx context.Context, filePath string) (*FileStats, error) {
	fullPath := filepath.Join(s.tablePath, filePath)
	if _, err := os.Stat(fullPath); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	pr, err := parquetfile.NewParquetReader(file) // Parquet file reader
	if err != nil {
		return nil, fmt.Errorf("failed to create parquet reader: %w", err)
	}
	defer pr.Close()

	stats := &FileStats{
		Path:      filePath,
		Size:      info.Size(),
		NumRows:   pr.NumRows(),
		MinValues: make(map[string]string), // Initialize maps
		MaxValues: make(map[string]string),
		NullCount: make(map[string]int64),
	}

	// No batch reading needed in GetFileStats, other functions will populate detailed stats.
	return stats, nil
}

// CalculateColumnStats calculates min, max, and null counts for columns in a Parquet file.
func CalculateColumnStats(pr *parquetfile.Reader, fileColumnStats map[string]*ColumnStats, ctx context.Context) error {
	arrowFileRdr, err := pqarrow.NewFileReader(pr, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	if err != nil {
		return fmt.Errorf("failed to create pqarrow.FileReader: %w", err)
	}

	// Get a record reader for all columns and all row groups
	recordReader, err := arrowFileRdr.GetRecordReader(ctx, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to get record reader from pqarrow.FileReader: %w", err)
	}
	defer recordReader.Release()

	// Read all batches
	for recordReader.Next() {
		batch := recordReader.Record()
		// Process batch
		for i, field := range batch.Schema().Fields() {
			col := batch.Column(i)
			// ... (rest of the logic for processing column remains the same)
			colStats := fileColumnStats[field.Name]
			if col.Len() > 0 {
				// Update min and max values
				if colStats.MinValue == "" {
					colStats.MinValue = col.ValueStr(0)
					colStats.MaxValue = col.ValueStr(0)
				}

				for j := 0; j < col.Len(); j++ {
					if col.IsNull(j) {
						colStats.NullCount++
						continue
					}

					val := col.ValueStr(j)
					if val < colStats.MinValue {
						colStats.MinValue = val
					}
					if val > colStats.MaxValue {
						colStats.MaxValue = val
					}
				}
			}
		}
		batch.Release() // Release the current batch after processing
	}
	if recordReader.Err() != nil {
		return fmt.Errorf("error reading records: %w", recordReader.Err())
	}

	return nil
}

// GetColumnStats returns statistics for a column
func (s *StatsManager) GetColumnStats(ctx context.Context, column string) (*ColumnStats, error) {
	// Find the column in the schema
	var field *arrow.Field
	for i := 0; i < len(s.schema.Fields()); i++ {
		if s.schema.Fields()[i].Name == column {
			field = &s.schema.Fields()[i]
			break
		}
	}

	if field == nil {
		return nil, fmt.Errorf("column not found: %s", column)
	}

	stats := &ColumnStats{
		Name:      column,
		NullCount: 0,
	}

	// Walk through all Parquet files to collect statistics
	err := filepath.Walk(s.tablePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			pr, err := parquetfile.NewParquetReader(file) // Parquet file reader
			if err != nil {
				return fmt.Errorf("failed to create parquet reader: %w", err)
			}
			defer pr.Close()

			fileColumnStats := map[string]*ColumnStats{column: stats}
			if err := CalculateColumnStats(pr, fileColumnStats, ctx); err != nil { // Pass context
				return fmt.Errorf("failed to calculate column stats: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get column stats: %w", err)
	}

	return stats, nil
}

func (s *StatsManager) GetTableStats(ctx context.Context) (map[string]*ColumnStats, error) {
	stats := make(map[string]*ColumnStats)

	// Initialize stats for each column
	for i := 0; i < len(s.schema.Fields()); i++ {
		field := s.schema.Fields()[i]
		stats[field.Name] = &ColumnStats{
			Name:      field.Name,
			NullCount: 0,
		}
	}

	// Walk through all Parquet files to collect statistics
	err := filepath.Walk(s.tablePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			pr, err := parquetfile.NewParquetReader(file) // Parquet file reader
			if err != nil {
				return fmt.Errorf("failed to create parquet reader: %w", err)
			}
			defer pr.Close()

			fileColumnStats := stats
			if err := CalculateColumnStats(pr, fileColumnStats, ctx); err != nil { // Pass context
				return fmt.Errorf("failed to calculate column stats: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get table stats: %w", err)
	}

	return stats, nil
}

func (s *StatsManager) GetNullCounts(ctx context.Context) (map[string]int64, error) {
	nullCounts := make(map[string]int64)

	// Walk through all Parquet files to collect null counts
	err := filepath.Walk(s.tablePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			pr, err := parquetfile.NewParquetReader(file)
			if err != nil {
				return fmt.Errorf("failed to create parquet reader: %w", err)
			}
			// pr.Close() is implicitly handled by file.Close() for underlying resources, 
			// and pqarrow.FileReader doesn't own pr to close it.

			arrowFileRdr, err := pqarrow.NewFileReader(pr, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
			if err != nil {
				return fmt.Errorf("failed to create pqarrow.FileReader for null counts: %w", err)
			}

			// Get a record reader for all columns and all row groups
			recordReader, err := arrowFileRdr.GetRecordReader(ctx, nil, nil)
			if err != nil {
				return fmt.Errorf("failed to get record reader for null counts: %w", err)
			}
			defer recordReader.Release() // Release after this file's processing

			// Read all batches
			for recordReader.Next() {
				batch := recordReader.Record()
				for i, field := range batch.Schema().Fields() {
					column := batch.Column(i)
					nullCounts[field.Name] += int64(column.NullN())
				}
				batch.Release() // Release the current batch after processing
			}
			if recordReader.Err() != nil {
                return fmt.Errorf("error reading records for null counts: %w", recordReader.Err())
            }
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get null counts: %w", err)
	}

	return nullCounts, nil
}

// WriteStats writes statistics to a file
func (s *StatsManager) WriteStats(ctx context.Context, stats interface{}, filePath string) error {
	data, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	fullPath := filepath.Join(s.tablePath, filePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write stats: %w", err)
	}

	return nil
}

// ReadStats reads statistics from a file
func (s *StatsManager) ReadStats(ctx context.Context, filePath string, stats interface{}) error {
	fullPath := filepath.Join(s.tablePath, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read stats: %w", err)
	}

	if err := json.Unmarshal(data, stats); err != nil {
		return fmt.Errorf("failed to unmarshal stats: %w", err)
	}

	return nil
}

// UpdateStats updates statistics for a file
func (s *StatsManager) UpdateStats(ctx context.Context, filePath string, record arrow.Record) error {
	stats, err := s.GetFileStats(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to get file stats: %w", err)
	}

	// Update statistics with the new record
	for i := 0; i < int(record.NumCols()); i++ { // Fixed: Cast NumCols to int
		col := record.Column(i)
		field := s.schema.Field(i)

		if col.Len() > 0 {
			// Update min and max values
			if stats.MinValues[field.Name] == "" {
				stats.MinValues[field.Name] = col.ValueStr(0)
				stats.MaxValues[field.Name] = col.ValueStr(0)
			}

			for j := 0; j < col.Len(); j++ {
				if col.IsNull(j) {
					stats.NullCount[field.Name]++
					continue
				}

				val := col.ValueStr(j)
				if val < stats.MinValues[field.Name] {
					stats.MinValues[field.Name] = val
				}
				if val > stats.MaxValues[field.Name] {
					stats.MaxValues[field.Name] = val
				}
			}
		}
	}

	// Write updated statistics
	statsPath := filepath.Join(s.tablePath, "_delta_log", "stats", filePath+".stats.json")
	return s.WriteStats(ctx, stats, statsPath)
}