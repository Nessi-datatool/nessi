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

// MetadataManager handles Delta table metadata operations
type MetadataManager struct {
	tablePath string
}

// NewMetadataManager creates a new metadata manager
func NewMetadataManager(tablePath string) *MetadataManager {
	return &MetadataManager{
		tablePath: tablePath,
	}
}

// ReadTableMetadata reads the latest table metadata from the transaction log
func (m *MetadataManager) ReadTableMetadata() (*DeltaTable, error) {
	// Get latest version
	version, err := m.getLatestVersion()
	if err != nil {
		return nil, err
	}

	// Read transaction log
	logFile := filepath.Join(m.tablePath, "_delta_log", fmt.Sprintf("%020d.json", version))
	data, err := os.ReadFile(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction log: %w", err)
	}

	// Parse metadata
	var metadata struct {
		Version      int64                    `json:"version"`
		Timestamp    int64                    `json:"timestamp"`
		Schema       map[string]interface{}   `json:"schema"`
		Files        []string                 `json:"files"`
		Partitions   map[string][]string      `json:"partitions"`
		Stats        *TableStats              `json:"stats"`
		Metadata     map[string]interface{}   `json:"metadata"`
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Parse schema
	fields, ok := metadata.Schema["fields"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid schema format")
	}

	arrowFields := make([]arrow.Field, len(fields))
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

		var dataType arrow.DataType
		switch typeStr {
		case "int32":
			dataType = &arrow.Int32Type{}
		case "string", "utf8":
			dataType = &arrow.StringType{}
		case "double", "float64":
			dataType = &arrow.Float64Type{}
		default:
			return nil, fmt.Errorf("unsupported field type: %s", typeStr)
		}

		arrowFields[i] = arrow.Field{Name: name, Type: dataType}
	}

	// Create DeltaTable
	table := &DeltaTable{
		Path:         m.tablePath,
		Version:      metadata.Version,
		LastModified: time.Unix(metadata.Timestamp, 0),
		Schema:       arrow.NewSchema(arrowFields, nil),
		Files:        metadata.Files,
		Partitions:   metadata.Partitions,
		Stats:        metadata.Stats,
		Metadata:     metadata.Metadata,
	}

	return table, nil
}

// WriteTableMetadata writes table metadata to the transaction log
func (m *MetadataManager) WriteTableMetadata(table *DeltaTable) error {
	if table == nil {
		return fmt.Errorf("table cannot be nil")
	}

	// Create metadata entry
	fields := make([]map[string]interface{}, len(table.Schema.Fields()))
	for i, field := range table.Schema.Fields() {
		var typeStr string
		switch field.Type.(type) {
		case *arrow.Int32Type:
			typeStr = "int32"
		case *arrow.StringType:
			typeStr = "utf8"
		case *arrow.Float64Type:
			typeStr = "float64"
		default:
			return fmt.Errorf("unsupported field type: %T", field.Type)
		}

		fields[i] = map[string]interface{}{
			"name": field.Name,
			"type": typeStr,
		}
	}

	metadata := map[string]interface{}{
		"version":      table.Version,
		"timestamp":    table.LastModified.Unix(),
		"schema":       map[string]interface{}{"fields": fields},
		"files":        table.Files,
		"partitions":   table.Partitions,
		"stats":        table.Stats,
		"metadata":     table.Metadata,
	}

	// Convert to JSON
	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write to transaction log
	logFile := filepath.Join(m.tablePath, "_delta_log", fmt.Sprintf("%020d.json", table.Version))
	if err := os.WriteFile(logFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write transaction log: %w", err)
	}

	return nil
}

// getLatestVersion gets the latest version from the transaction log
func (m *MetadataManager) getLatestVersion() (int64, error) {
	// Read _delta_log directory
	logDir := filepath.Join(m.tablePath, "_delta_log")
	entries, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read _delta_log directory: %w", err)
	}

	if len(entries) == 0 {
		return 0, nil
	}

	// Sort entries by name (which contains version number)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() > entries[j].Name()
	})

	// Parse version from filename
	var version int64
	_, err = fmt.Sscanf(entries[0].Name(), "%020d.json", &version)
	if err != nil {
		return 0, fmt.Errorf("failed to parse version from filename: %w", err)
	}

	return version, nil
}

// GetVersions returns all available versions
func (m *MetadataManager) GetVersions() ([]int64, error) {
	// Read _delta_log directory
	logDir := filepath.Join(m.tablePath, "_delta_log")
	entries, err := os.ReadDir(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []int64{}, nil
		}
		return nil, fmt.Errorf("failed to read _delta_log directory: %w", err)
	}

	versions := make([]int64, 0, len(entries))
	for _, entry := range entries {
		var version int64
		if _, err := fmt.Sscanf(entry.Name(), "%020d.json", &version); err == nil {
			versions = append(versions, version)
		}
	}

	// Sort versions in descending order
	sort.Slice(versions, func(i, j int) bool {
		return versions[i] > versions[j]
	})

	return versions, nil
}

// GetTableAtVersion returns the table metadata at a specific version
func (m *MetadataManager) GetTableAtVersion(version int64) (*DeltaTable, error) {
	// Read transaction log
	logFile := filepath.Join(m.tablePath, "_delta_log", fmt.Sprintf("%020d.json", version))
	data, err := os.ReadFile(logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read transaction log: %w", err)
	}

	// Parse metadata
	var metadata struct {
		Version      int64                    `json:"version"`
		Timestamp    int64                    `json:"timestamp"`
		Schema       *arrow.Schema            `json:"schema"`
		Files        []string                 `json:"files"`
		Partitions   map[string][]string      `json:"partitions"`
		Stats        *TableStats              `json:"stats"`
		Metadata     map[string]interface{}   `json:"metadata"`
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	// Create DeltaTable
	table := &DeltaTable{
		Path:         m.tablePath,
		Version:      metadata.Version,
		LastModified: time.Unix(metadata.Timestamp, 0),
		Schema:       metadata.Schema,
		Files:        metadata.Files,
		Partitions:   metadata.Partitions,
		Stats:        metadata.Stats,
		Metadata:     metadata.Metadata,
	}

	return table, nil
}
