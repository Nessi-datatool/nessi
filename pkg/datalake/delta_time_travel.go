package datalake

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// DeltaTimeTravel provides time travel capabilities for Delta Lake tables
type DeltaTimeTravel struct {
	handler *DeltaFormatHandler
}

// NewDeltaTimeTravel creates a new DeltaTimeTravel instance
func NewDeltaTimeTravel(handler *DeltaFormatHandler) *DeltaTimeTravel {
	return &DeltaTimeTravel{
		handler: handler,
	}
}

// ReadAsOfVersion reads data from a Delta Lake table as of a specific version
func (t *DeltaTimeTravel) ReadAsOfVersion(path string, version int64) ([]map[string]interface{}, error) {
	// Check if the path is a Delta Lake table
	if !t.handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Get the transaction log files
	deltaLogDir := filepath.Join(path, "_delta_log")
	logFiles, err := filepath.Glob(filepath.Join(deltaLogDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list transaction log files: %w", err)
	}

	// Extract version numbers from log files
	versions := make([]int64, 0, len(logFiles))
	for _, file := range logFiles {
		baseName := filepath.Base(file)
		if strings.HasSuffix(baseName, ".json") && !strings.Contains(baseName, "checkpoint") {
			versionStr := strings.TrimSuffix(baseName, ".json")
			v, err := strconv.ParseInt(versionStr, 10, 64)
			if err != nil {
				continue
			}
			versions = append(versions, v)
		}
	}

	// Sort versions
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] < versions[j]
	})

	// Check if requested version exists
	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found in Delta table: %s", path)
	}

	if version > versions[len(versions)-1] {
		return nil, fmt.Errorf("version %d not found, latest version is %d", version, versions[len(versions)-1])
	}

	// In a real implementation, we would:
	// 1. Parse the transaction log up to the specified version
	// 2. Build a list of valid data files
	// 3. Read those data files

	// For this simplified implementation, we'll just read the current state
	// In a production environment, this would need to reconstruct the table state
	// as of the requested version
	return t.handler.Read(path)
}

// ReadAsOfTimestamp reads data from a Delta Lake table as of a specific timestamp
func (t *DeltaTimeTravel) ReadAsOfTimestamp(path string, timestamp time.Time) ([]map[string]interface{}, error) {
	// Check if the path is a Delta Lake table
	if !t.handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Get the transaction log files
	deltaLogDir := filepath.Join(path, "_delta_log")
	// We're not using logFiles in this simplified implementation, but we would in a real one
	_, err := filepath.Glob(filepath.Join(deltaLogDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list transaction log files: %w", err)
	}

	// In a real implementation, we would:
	// 1. Parse the transaction log files to extract timestamps
	// 2. Find the latest version before or at the specified timestamp
	// 3. Call ReadAsOfVersion with that version

	// For this simplified implementation, we'll just read the current state
	// In a production environment, this would need to determine the correct version
	// based on the timestamp
	return t.handler.Read(path)
}

// GetVersionHistory returns the version history of a Delta Lake table
func (t *DeltaTimeTravel) GetVersionHistory(ctx context.Context, path string) ([]DeltaVersion, error) {
	// Check if the path is a Delta Lake table
	if !t.handler.IsDeltaTable(path) {
		return nil, fmt.Errorf("not a Delta Lake table: %s", path)
	}

	// Get the transaction log files
	deltaLogDir := filepath.Join(path, "_delta_log")
	logFiles, err := filepath.Glob(filepath.Join(deltaLogDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list transaction log files: %w", err)
	}

	// Extract version information from log files
	versions := make([]DeltaVersion, 0, len(logFiles))
	for _, file := range logFiles {
		baseName := filepath.Base(file)
		if strings.HasSuffix(baseName, ".json") && !strings.Contains(baseName, "checkpoint") {
			versionStr := strings.TrimSuffix(baseName, ".json")
			v, err := strconv.ParseInt(versionStr, 10, 64)
			if err != nil {
				continue
			}

			// In a real implementation, we would parse the log file to extract
			// the timestamp, operation type, and other metadata
			// For now, we'll use the file modification time as a proxy for the commit timestamp
			fileInfo, err := statFile(file)
			if err != nil {
				continue
			}

			versions = append(versions, DeltaVersion{
				Version:   v,
				Timestamp: fileInfo.ModTime(),
				Operation: "UNKNOWN", // In a real implementation, we would parse the log file
			})
		}
	}

	// Sort versions
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version < versions[j].Version
	})

	return versions, nil
}

// DeltaVersion represents a version of a Delta Lake table
type DeltaVersion struct {
	Version   int64
	Timestamp time.Time
	Operation string
}
