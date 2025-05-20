package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DeltaTimeTravelOptions represents options for time travel queries in Delta Lake
type DeltaTimeTravelOptions struct {
	// Version to query (if specified, Timestamp is ignored)
	Version *int64

	// Timestamp to query (used if Version is nil)
	Timestamp *time.Time
}

// DeltaTimeTravelResult represents the result of a time travel query in Delta Lake
type DeltaTimeTravelResult struct {
	// The version that was queried
	Version int64

	// The timestamp that was queried (or the timestamp of the version)
	Timestamp time.Time

	// The schema fields at the specified version/timestamp
	SchemaFields []SchemaField

	// The files at the specified version/timestamp
	Files []string
}

// Using SchemaField from table.go

// TimeTravel queries a Delta table as it existed at a specific version or timestamp
func (m *MetadataManager) TimeTravel(options DeltaTimeTravelOptions) (*DeltaTimeTravelResult, error) {
	var targetVersion int64
	var targetTimestamp time.Time

	// If version is specified, use it
	if options.Version != nil {
		targetVersion = int64(*options.Version)
		return m.getTableInfoAtVersion(targetVersion)
	}

	// If timestamp is specified, find the closest version
	if options.Timestamp != nil {
		targetTimestamp = *options.Timestamp
		version, err := m.findVersionAtTimestamp(targetTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to find version at timestamp %s: %w", targetTimestamp.Format(time.RFC3339), err)
		}

		return m.getTableInfoAtVersion(version)
	}

	return nil, fmt.Errorf("either Version or Timestamp must be specified")
}

// findVersionAtTimestamp finds the version that was current at the specified timestamp
func (m *MetadataManager) findVersionAtTimestamp(timestamp time.Time) (int64, error) {
	// Get version history
	history, err := m.GetVersionHistory()
	if err != nil {
		return 0, fmt.Errorf("failed to get version history: %w", err)
	}

	if len(history) == 0 {
		return 0, fmt.Errorf("no versions found")
	}

	// Find the latest version that was created before or at the specified timestamp
	for _, entry := range history {
		entryTimestamp := time.Unix(0, entry.Timestamp*int64(time.Millisecond))
		if entryTimestamp.Before(timestamp) || entryTimestamp.Equal(timestamp) {
			return entry.Version, nil
		}
	}

	// If all versions are after the timestamp, return the oldest version
	return history[len(history)-1].Version, nil
}

// getTableInfoAtVersion gets table information at a specific version
func (m *MetadataManager) getTableInfoAtVersion(version int64) (*DeltaTimeTravelResult, error) {
	// Read transaction log
	logFile := filepath.Join(m.tablePath, "_delta_log", fmt.Sprintf("%020d.json", version))
	data, err := os.ReadFile(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction log: %w", err)
	}

	// Parse metadata
	var metadata struct {
		Version   int64                  `json:"version"`
		Timestamp int64                  `json:"timestamp"`
		Schema    map[string]interface{} `json:"schema"`
		Files     []string               `json:"files"`
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Parse schema
	fields, ok := metadata.Schema["fields"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid schema format")
	}

	schemaFields := make([]SchemaField, len(fields))
	for i, field := range fields {
		fieldMap, ok := field.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid field format")
		}

		name, ok := fieldMap["name"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid field name")
		}

		typeStr, ok := fieldMap["type"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid field type")
		}

		schemaFields[i] = SchemaField{
			Name: name,
			Type: typeStr,
		}
	}

	// Create result
	result := &DeltaTimeTravelResult{
		Version:      version,
		Timestamp:    time.Unix(metadata.Timestamp, 0),
		SchemaFields: schemaFields,
		Files:        metadata.Files,
	}

	return result, nil
}

// GetSchemaFieldsAtVersion returns the schema fields that existed at a specific version
func (m *MetadataManager) GetSchemaFieldsAtVersion(version int64) ([]SchemaField, error) {
	result, err := m.getTableInfoAtVersion(version)
	if err != nil {
		return nil, err
	}
	return result.SchemaFields, nil
}

// GetSchemaFieldsAtTimestamp returns the schema fields that existed at a specific timestamp
func (m *MetadataManager) GetSchemaFieldsAtTimestamp(timestamp time.Time) ([]SchemaField, error) {
	version, err := m.findVersionAtTimestamp(timestamp)
	if err != nil {
		return nil, err
	}
	return m.GetSchemaFieldsAtVersion(version)
}

// GetFilesAtVersion returns the files that existed at a specific version
func (m *MetadataManager) GetFilesAtVersion(version int64) ([]string, error) {
	result, err := m.getTableInfoAtVersion(version)
	if err != nil {
		return nil, err
	}
	return result.Files, nil
}
