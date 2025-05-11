package delta

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// MaintenanceManager handles Delta Lake maintenance operations
type MaintenanceManager struct {
	connector *DeltaConnector
}

// NewMaintenanceManager creates a new maintenance manager
func NewMaintenanceManager(connector *DeltaConnector) *MaintenanceManager {
	return &MaintenanceManager{
		connector: connector,
	}
}

// CheckpointOptions defines options for checkpoint creation
type CheckpointOptions struct {
	Version     int64
	Force       bool
	RetainFiles bool
}

// CreateCheckpoint creates a checkpoint for the Delta table
func (m *MaintenanceManager) CreateCheckpoint(ctx context.Context, options CheckpointOptions) error {
	// Get the target version
	version := options.Version
	if version == 0 {
		version = m.connector.GetVersion()
	}

	// Check if checkpoint already exists
	checkpointPath := filepath.Join(m.connector.TablePath, "_delta_log", fmt.Sprintf("%d.checkpoint.parquet", version))
	if _, err := os.Stat(checkpointPath); err == nil && !options.Force {
		return fmt.Errorf("checkpoint already exists for version %d", version)
	}

	// Read all actions up to the target version
	var actions []Action
	for v := int64(0); v <= version; v++ {
		filePath := filepath.Join(m.connector.TablePath, "_delta_log", fmt.Sprintf("%d.json", v))
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read log file: %w", err)
		}

		var versionActions []Action
		if err := json.Unmarshal(data, &versionActions); err != nil {
			return fmt.Errorf("failed to parse log file: %w", err)
		}
		actions = append(actions, versionActions...)
	}

	// Create checkpoint schema
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "version", Type: arrow.PrimitiveTypes.Int64},
		{Name: "timestamp", Type: arrow.PrimitiveTypes.Int64},
		{Name: "action", Type: arrow.BinaryTypes.String},
	}, nil)

	// Create checkpoint data
	builder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builder.Release()

	for _, action := range actions {
		builder.Field(0).(*array.Int64Builder).Append(version)
		builder.Field(1).(*array.Int64Builder).Append(time.Now().UnixMilli())
		actionJSON, _ := json.Marshal(action)
		builder.Field(2).(*array.StringBuilder).Append(string(actionJSON))
	}

	record := builder.NewRecord()
	defer record.Release()

	// Write checkpoint file
	if err := WriteRecordToParquet(checkpointPath, record); err != nil {
		return fmt.Errorf("failed to write checkpoint file: %w", err)
	}

	// Create checkpoint metadata
	metadata := struct {
		Version     int64     `json:"version"`
		Timestamp   time.Time `json:"timestamp"`
		NumActions  int       `json:"numActions"`
		FileSize    int64     `json:"fileSize"`
	}{
		Version:     version,
		Timestamp:   time.Now(),
		NumActions:  len(actions),
		FileSize:    getFileSize(checkpointPath),
	}

	// Write metadata
	metadataPath := filepath.Join(m.connector.TablePath, "_delta_log", fmt.Sprintf("%d.checkpoint.metadata.json", version))
	metadataData, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	if err := os.WriteFile(metadataPath, metadataData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// VacuumOptions defines options for vacuum operation
type VacuumOptions struct {
	RetentionHours int64
	DryRun         bool
	Force          bool
}

// VacuumStats represents statistics for a vacuum operation
type VacuumStats struct {
	DeletedFiles    []string
	DeletedSize     int64
	RetainedFiles   []string
	RetainedSize    int64
	TotalFiles      int
	TotalSize       int64
	StartTime       time.Time
	EndTime         time.Time
	DurationSeconds float64
}

// Vacuum performs a vacuum operation on the Delta table
func (m *MaintenanceManager) Vacuum(ctx context.Context, options VacuumOptions) (*VacuumStats, error) {
	stats := &VacuumStats{
		StartTime: time.Now(),
	}

	// Get all files in the table
	files, err := m.getAllFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %w", err)
	}

	// Get the cutoff time
	cutoffTime := time.Now().Add(-time.Duration(options.RetentionHours) * time.Hour)

	// Process each file
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		stats.TotalFiles++
		stats.TotalSize += info.Size()

		// Check if file should be deleted
		if info.ModTime().Before(cutoffTime) {
			stats.DeletedFiles = append(stats.DeletedFiles, file)
			stats.DeletedSize += info.Size()

			if !options.DryRun {
				if err := os.Remove(file); err != nil {
					return nil, fmt.Errorf("failed to delete file %s: %w", file, err)
				}
			}
		} else {
			stats.RetainedFiles = append(stats.RetainedFiles, file)
			stats.RetainedSize += info.Size()
		}
	}

	stats.EndTime = time.Now()
	stats.DurationSeconds = stats.EndTime.Sub(stats.StartTime).Seconds()

	return stats, nil
}

// getAllFiles returns all files in the Delta table
func (m *MaintenanceManager) getAllFiles() ([]string, error) {
	var files []string
	err := filepath.Walk(m.connector.TablePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".parquet") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return files, nil
}

// ListCheckpoints lists all checkpoints in the Delta table
func (m *MaintenanceManager) ListCheckpoints(ctx context.Context) ([]CheckpointInfo, error) {
	logPath := filepath.Join(m.connector.TablePath, "_delta_log")
	files, err := os.ReadDir(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log directory: %w", err)
	}

	var checkpoints []CheckpointInfo
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".checkpoint.parquet") {
			version, err := parseVersionFromFilename(file.Name())
			if err != nil {
				continue
			}

			// Read checkpoint metadata
			metadataPath := filepath.Join(logPath, fmt.Sprintf("%d.checkpoint.metadata.json", version))
			data, err := os.ReadFile(metadataPath)
			if err != nil {
				continue
			}

			var metadata struct {
				Version     int64     `json:"version"`
				Timestamp   time.Time `json:"timestamp"`
				NumActions  int       `json:"numActions"`
				FileSize    int64     `json:"fileSize"`
			}

			if err := json.Unmarshal(data, &metadata); err != nil {
				continue
			}

			checkpoints = append(checkpoints, CheckpointInfo{
				Version:     metadata.Version,
				Timestamp:   metadata.Timestamp,
				NumActions:  metadata.NumActions,
				FileSize:    metadata.FileSize,
				Path:        filepath.Join(logPath, file.Name()),
			})
		}
	}

	sort.Slice(checkpoints, func(i, j int) bool {
		return checkpoints[i].Version < checkpoints[j].Version
	})

	return checkpoints, nil
}

// CheckpointInfo represents information about a checkpoint
type CheckpointInfo struct {
	Version     int64
	Timestamp   time.Time
	NumActions  int
	FileSize    int64
	Path        string
}

// DeleteCheckpoint deletes a checkpoint
func (m *MaintenanceManager) DeleteCheckpoint(ctx context.Context, version int64) error {
	checkpointPath := filepath.Join(m.connector.TablePath, "_delta_log", fmt.Sprintf("%d.checkpoint.parquet", version))
	metadataPath := filepath.Join(m.connector.TablePath, "_delta_log", fmt.Sprintf("%d.checkpoint.metadata.json", version))

	// Delete checkpoint file
	if err := os.Remove(checkpointPath); err != nil {
		return fmt.Errorf("failed to delete checkpoint file: %w", err)
	}

	// Delete metadata file
	if err := os.Remove(metadataPath); err != nil {
		return fmt.Errorf("failed to delete metadata file: %w", err)
	}

	return nil
}

// Utility function to check if a file exists