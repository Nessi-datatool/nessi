package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
)

// VersionEntry represents a single entry in the version history
type VersionEntry struct {
	Version     int64             `json:"version"`
	Timestamp   int64             `json:"timestamp"`
	Operation   string            `json:"operation"`
	UserName    string            `json:"userName,omitempty"`
	UserID      string            `json:"userId,omitempty"`
	Parameters  map[string]string `json:"parameters,omitempty"`
	Description string            `json:"description,omitempty"`
}

// VersionComparison represents the differences between two versions
type VersionComparison struct {
	OlderVersion   int64               `json:"olderVersion"`
	NewerVersion   int64               `json:"newerVersion"`
	SchemaChanges  []SchemaFieldChange `json:"schemaChanges"`
	FilesAdded     int                 `json:"filesAdded"`
	FilesRemoved   int                 `json:"filesRemoved"`
	RecordsAdded   int64               `json:"recordsAdded"`
	RecordsRemoved int64               `json:"recordsRemoved"`
}

// GetVersionHistory returns the version history of a Delta table
func (m *MetadataManager) GetVersionHistory() ([]VersionEntry, error) {
	// Get all versions
	versions, err := m.GetVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		return []VersionEntry{}, nil
	}

	// Get commit info for each version
	history := make([]VersionEntry, 0, len(versions))
	for _, version := range versions {
		// Read commit info file
		commitInfoPath := filepath.Join(m.tablePath, "_delta_log", fmt.Sprintf("%020d.commit.json", version))
		data, err := os.ReadFile(commitInfoPath)

		// If commit info doesn't exist, try to get basic info from the log file
		if os.IsNotExist(err) {
			// Read transaction log
			logFile := filepath.Join(m.tablePath, "_delta_log", fmt.Sprintf("%020d.json", version))
			data, err := os.ReadFile(logFile)
			if err != nil {
				continue // Skip this version if we can't read the log file
			}

			// Try to extract basic info
			var logData map[string]interface{}
			if err := json.Unmarshal(data, &logData); err != nil {
				continue
			}

			// Create basic entry
			entry := VersionEntry{
				Version:   version,
				Timestamp: time.Now().UnixMilli(), // Use current time as fallback
				Operation: "unknown",
			}

			// Try to extract timestamp
			if ts, ok := logData["timestamp"].(float64); ok {
				entry.Timestamp = int64(ts)
			}

			history = append(history, entry)
			continue
		} else if err != nil {
			return nil, fmt.Errorf("failed to read commit info for version %d: %w", version, err)
		}

		// Parse commit info
		var commitInfo struct {
			Timestamp           int64             `json:"timestamp"`
			UserID              string            `json:"userId"`
			UserName            string            `json:"userName"`
			Operation           string            `json:"operation"`
			OperationParameters map[string]string `json:"operationParameters"`
		}

		if err := json.Unmarshal(data, &commitInfo); err != nil {
			return nil, fmt.Errorf("failed to parse commit info for version %d: %w", version, err)
		}

		// Create version entry
		entry := VersionEntry{
			Version:    version,
			Timestamp:  commitInfo.Timestamp,
			Operation:  commitInfo.Operation,
			UserName:   commitInfo.UserName,
			UserID:     commitInfo.UserID,
			Parameters: commitInfo.OperationParameters,
		}

		history = append(history, entry)
	}

	// Sort by version (descending)
	sort.Slice(history, func(i, j int) bool {
		return history[i].Version > history[j].Version
	})

	return history, nil
}

// CompareVersions compares two versions of a Delta table
func (m *MetadataManager) CompareVersions(version1, version2 int64) (*VersionComparison, error) {
	// Ensure version1 < version2
	if version1 > version2 {
		version1, version2 = version2, version1
	}

	// Get tables at both versions
	table1, err := m.GetTableAtVersion(version1)
	if err != nil {
		return nil, fmt.Errorf("failed to get table at version %d: %w", version1, err)
	}

	table2, err := m.GetTableAtVersion(version2)
	if err != nil {
		return nil, fmt.Errorf("failed to get table at version %d: %w", version2, err)
	}

	// Compare schemas
	schemaChanges := DiffSchemas(table1.Schema, table2.Schema)

	// Compare files
	filesAdded := 0
	filesRemoved := 0

	// Create maps for faster lookup
	files1 := make(map[string]bool)
	for _, file := range table1.Files {
		files1[file] = true
	}

	files2 := make(map[string]bool)
	for _, file := range table2.Files {
		files2[file] = true
	}

	// Count added files
	for file := range files2 {
		if !files1[file] {
			filesAdded++
		}
	}

	// Count removed files
	for file := range files1 {
		if !files2[file] {
			filesRemoved++
		}
	}

	// Estimate record changes (based on file count)
	// This is a rough estimate; for accurate counts, we'd need to read all files
	recordsAdded := int64(filesAdded) * 1000 // Assuming 1000 records per file
	recordsRemoved := int64(filesRemoved) * 1000

	comparison := &VersionComparison{
		OlderVersion:   version1,
		NewerVersion:   version2,
		SchemaChanges:  schemaChanges,
		FilesAdded:     filesAdded,
		FilesRemoved:   filesRemoved,
		RecordsAdded:   recordsAdded,
		RecordsRemoved: recordsRemoved,
	}

	return comparison, nil
}

// Note: RollbackToVersion method has been moved to metadata.go to avoid duplication

// Helper function to convert schema to map
func schemaToMap(schema *arrow.Schema) map[string]interface{} {
	fields := make([]map[string]interface{}, len(schema.Fields()))
	for i, field := range schema.Fields() {
		var typeStr string
		switch field.Type.(type) {
		case *arrow.Int32Type:
			typeStr = "int32"
		case *arrow.StringType:
			typeStr = "utf8"
		case *arrow.Float64Type:
			typeStr = "float64"
		case *arrow.TimestampType:
			typeStr = "timestamp"
		case *arrow.BooleanType:
			typeStr = "boolean"
		default:
			typeStr = "string" // Default to string for unsupported types
		}

		fields[i] = map[string]interface{}{
			"name": field.Name,
			"type": typeStr,
		}
	}

	return map[string]interface{}{
		"fields": fields,
	}
}
