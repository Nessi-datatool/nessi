package delta

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/memory"
	parquetfile "github.com/apache/arrow/go/v15/parquet/file"
	"github.com/apache/arrow/go/v15/parquet/pqarrow"
)

// PartitionManager manages Delta Lake partitions
type PartitionManager struct {
	tablePath    string
	partitions   []string
	partitionMap map[string]map[string]string
}

// NewPartitionManager creates a new partition manager
func NewPartitionManager(tablePath string, partitions []string) *PartitionManager {
	return &PartitionManager{
		tablePath:    tablePath,
		partitions:   partitions,
		partitionMap: make(map[string]map[string]string),
	}
}

// GetPartitions returns the current partitions
func (p *PartitionManager) GetPartitions() []string {
	return p.partitions
}

// GetPartitionValues returns the partition values for a path
func (p *PartitionManager) GetPartitionValues(path string) (map[string]string, error) {
	if values, exists := p.partitionMap[path]; exists {
		return values, nil
	}

	// Parse partition values from path
	values := make(map[string]string)
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.Contains(part, "=") {
			kv := strings.Split(part, "=")
			if len(kv) == 2 {
				values[kv[0]] = kv[1]
			}
		}
	}

	p.partitionMap[path] = values
	return values, nil
}

// ListPartitions lists all partitions in the table
func (p *PartitionManager) ListPartitions(ctx context.Context) ([]string, error) {
	var partitions []string
	err := filepath.Walk(p.tablePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != p.tablePath {
			relPath, err := filepath.Rel(p.tablePath, path)
			if err != nil {
				return err
			}

			// Check if the path contains partition values
			if strings.Contains(relPath, "=") {
				partitions = append(partitions, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list partitions: %w", err)
	}

	sort.Strings(partitions)
	return partitions, nil
}

// GetPartitionStats returns statistics for a partition
func (p *PartitionManager) GetPartitionStats(ctx context.Context, partition string) (*PartitionStats, error) {
	partitionPath := filepath.Join(p.tablePath, partition)
	if _, err := os.Stat(partitionPath); err != nil {
		return nil, fmt.Errorf("partition not found: %w", err)
	}

	stats := &PartitionStats{
		Path:      partition,
		NumFiles:  0,
		SizeBytes: 0,
		NumRows:   0,
	}

	err := filepath.Walk(partitionPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			stats.NumFiles++
			stats.SizeBytes += info.Size()

			// Read Parquet file to get row count
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			reader, err := parquetfile.NewParquetReader(file)
			if err != nil {
				return fmt.Errorf("failed to create parquet reader: %w", err)
			}
			defer reader.Close()

			stats.NumRows += reader.NumRows()
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get partition stats: %w", err)
	}

	return stats, nil
}

// PartitionStats represents statistics for a partition
type PartitionStats struct {
	Path      string
	NumFiles  int64
	SizeBytes int64
	NumRows   int64
}

// PrunePartitions prunes partitions based on a filter
func (p *PartitionManager) PrunePartitions(ctx context.Context, filter PartitionFilter) ([]string, error) {
	partitions, err := p.ListPartitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list partitions: %w", err)
	}

	var prunedPartitions []string
	for _, partition := range partitions {
		values, err := p.GetPartitionValues(partition)
		if err != nil {
			return nil, fmt.Errorf("failed to get partition values: %w", err)
		}

		if filter.Match(values) {
			prunedPartitions = append(prunedPartitions, partition)
		}
	}

	return prunedPartitions, nil
}

// PartitionFilter represents a filter for pruning partitions
type PartitionFilter interface {
	Match(values map[string]string) bool
}

// AndFilter combines multiple filters with AND
type AndFilter struct {
	Filters []PartitionFilter
}

// Match implements PartitionFilter
func (f *AndFilter) Match(values map[string]string) bool {
	for _, filter := range f.Filters {
		if !filter.Match(values) {
			return false
		}
	}
	return true
}

// OrFilter combines multiple filters with OR
type OrFilter struct {
	Filters []PartitionFilter
}

// Match implements PartitionFilter
func (f *OrFilter) Match(values map[string]string) bool {
	for _, filter := range f.Filters {
		if filter.Match(values) {
			return true
		}
	}
	return false
}

// EqualFilter matches when a partition value equals a specific value
type EqualFilter struct {
	Column string
	Value  string
}

// Match implements PartitionFilter
func (f *EqualFilter) Match(values map[string]string) bool {
	return values[f.Column] == f.Value
}

// RangeFilter matches when a partition value is within a range
type RangeFilter struct {
	Column    string
	MinValue  string
	MaxValue  string
	Inclusive bool
}

// Match implements PartitionFilter
func (f *RangeFilter) Match(values map[string]string) bool {
	value := values[f.Column]
	if f.Inclusive {
		return value >= f.MinValue && value <= f.MaxValue
	}
	return value > f.MinValue && value < f.MaxValue
}

// InFilter matches when a partition value is in a set of values
type InFilter struct {
	Column string
	Values []string
}

// Match implements PartitionFilter
func (f *InFilter) Match(values map[string]string) bool {
	value := values[f.Column]
	for _, v := range f.Values {
		if value == v {
			return true
		}
	}
	return false
}

// NotFilter negates another filter
type NotFilter struct {
	Filter PartitionFilter
}

// Match implements PartitionFilter
func (f *NotFilter) Match(values map[string]string) bool {
	return !f.Filter.Match(values)
}

// CreatePartition creates a new partition
func (p *PartitionManager) CreatePartition(ctx context.Context, partition string) error {
	partitionPath := filepath.Join(p.tablePath, partition)
	if err := os.MkdirAll(partitionPath, 0755); err != nil {
		return fmt.Errorf("failed to create partition directory: %w", err)
	}

	// Parse and validate partition values
	values, err := p.GetPartitionValues(partition)
	if err != nil {
		return fmt.Errorf("failed to get partition values: %w", err)
	}

	// Validate that all partition columns are present
	for _, col := range p.partitions {
		if _, exists := values[col]; !exists {
			return fmt.Errorf("missing partition column: %s", col)
		}
	}

	return nil
}

// DeletePartition deletes a partition
func (p *PartitionManager) DeletePartition(ctx context.Context, partition string) error {
	partitionPath := filepath.Join(p.tablePath, partition)
	if err := os.RemoveAll(partitionPath); err != nil {
		return fmt.Errorf("failed to delete partition: %w", err)
	}

	delete(p.partitionMap, partition)
	return nil
}

// GetPartitionSchema returns the schema for a partition
func (p *PartitionManager) GetPartitionSchema(ctx context.Context, partition string) (*arrow.Schema, error) {
	partitionPath := filepath.Join(p.tablePath, partition)
	if _, err := os.Stat(partitionPath); err != nil {
		return nil, fmt.Errorf("partition not found: %w", err)
	}

	// Find the first Parquet file in the partition
	var schema *arrow.Schema
	err := filepath.Walk(partitionPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			reader, err := parquetfile.NewParquetReader(file)
			if err != nil {
				return fmt.Errorf("failed to create parquet reader: %w", err)
			}
			defer reader.Close()

			arrowReader, err := pqarrow.NewFileReader(reader, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
			if err != nil {
				return fmt.Errorf("failed to create arrow reader from parquet: %w", err)
			}
			
			schema, err = arrowReader.Schema()
			if err != nil {
				return fmt.Errorf("failed to get schema from arrow reader: %w", err)
			}
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get partition schema: %w", err)
	}

	if schema == nil {
		return nil, fmt.Errorf("no schema found in partition")
	}

	return schema, nil
}

// GetPartitionFiles returns all files in a partition
func (p *PartitionManager) GetPartitionFiles(ctx context.Context, partition string) ([]string, error) {
	partitionPath := filepath.Join(p.tablePath, partition)
	if _, err := os.Stat(partitionPath); err != nil {
		return nil, fmt.Errorf("partition not found: %w", err)
	}

	var files []string
	err := filepath.Walk(partitionPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			relPath, err := filepath.Rel(p.tablePath, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get partition files: %w", err)
	}

	sort.Strings(files)
	return files, nil
}