package datalake

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
)

// ... (rest of the code remains the same)
// SchemaChange represents a change between two schema versions
type SchemaChange struct {
	Type      string      // "added", "removed", "type_changed"
	FieldName string      // Name of the field that changed
	OldType   arrow.DataType // Old data type (nil for added fields)
	NewType   arrow.DataType // New data type (nil for removed fields)
}

// SchemaHistory represents the history of a schema
type SchemaHistory struct {
	Versions []SchemaVersion
}

// SchemaVersion represents a specific version of a schema
type SchemaVersion struct {
	Version      int64
	Timestamp    time.Time
	Schema       *arrow.Schema
	SchemaFields []arrow.Field
}

// SerializableField is a JSON-serializable representation of arrow.Field
type SerializableField struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Nullable bool            `json:"nullable"`
	Metadata json.RawMessage `json:"metadata,omitempty"`
}

// SerializableSchemaVersion is a JSON-serializable representation of SchemaVersion
type SerializableSchemaVersion struct {
	Version   int64              `json:"version"`
	Timestamp time.Time          `json:"timestamp"`
	Fields    []SerializableField `json:"fields"`
}

// MarshalJSON implements custom JSON marshaling for SchemaVersion
func (sv SchemaVersion) MarshalJSON() ([]byte, error) {
	serializable := SerializableSchemaVersion{
		Version:   sv.Version,
		Timestamp: sv.Timestamp,
		Fields:    make([]SerializableField, len(sv.SchemaFields)),
	}

	// Convert arrow.Field to SerializableField
	for i, field := range sv.SchemaFields {
		serializable.Fields[i] = SerializableField{
			Name:     field.Name,
			Type:     field.Type.String(),
			Nullable: field.Nullable,
		}
	}

	return json.Marshal(serializable)
}

// UnmarshalJSON implements custom JSON unmarshaling for SchemaVersion
func (sv *SchemaVersion) UnmarshalJSON(data []byte) error {
	var serializable SerializableSchemaVersion
	if err := json.Unmarshal(data, &serializable); err != nil {
		return err
	}

	// Set basic fields
	sv.Version = serializable.Version
	sv.Timestamp = serializable.Timestamp

	// Convert SerializableField to arrow.Field
	sv.SchemaFields = make([]arrow.Field, len(serializable.Fields))
	for i, field := range serializable.Fields {
		// Convert type string to arrow.DataType
		var dataType arrow.DataType
		switch field.Type {
		case "int64":
			dataType = arrow.PrimitiveTypes.Int64
		case "float64":
			dataType = arrow.PrimitiveTypes.Float64
		case "bool":
			dataType = arrow.FixedWidthTypes.Boolean
		case "timestamp[s]":
			dataType = arrow.FixedWidthTypes.Timestamp_s
		default:
			dataType = arrow.BinaryTypes.String
		}

		sv.SchemaFields[i] = arrow.Field{
			Name:     field.Name,
			Type:     dataType,
			Nullable: field.Nullable,
		}
	}

	// Create arrow.Schema from fields
	sv.Schema = arrow.NewSchema(sv.SchemaFields, nil)

	return nil
}

// GetSchemaHistory returns the schema history for a table
func (m *MetadataManager) GetSchemaHistory() (*SchemaHistory, error) {
	versions, err := m.GetVersions()
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	history := &SchemaHistory{
		Versions: make([]SchemaVersion, 0, len(versions)),
	}

	for _, version := range versions {
		table, err := m.GetTableAtVersion(version)
		if err != nil {
			return nil, fmt.Errorf("failed to get table at version %d: %w", version, err)
		}

		history.Versions = append(history.Versions, SchemaVersion{
			Version:      version,
			Timestamp:    table.LastModified,
			Schema:       table.Schema,
			SchemaFields: table.Schema.Fields(),
		})
	}

	return history, nil
}

// DiffSchemas compares two schemas and returns the differences
func DiffSchemas(oldSchema, newSchema *arrow.Schema) []SchemaChange {
	changes := make([]SchemaChange, 0)

	// Create maps for quick lookup
	oldFields := make(map[string]arrow.Field)
	for _, field := range oldSchema.Fields() {
		oldFields[field.Name] = field
	}

	newFields := make(map[string]arrow.Field)
	for _, field := range newSchema.Fields() {
		newFields[field.Name] = field
	}

	// Find removed and changed fields
	for name, oldField := range oldFields {
		newField, exists := newFields[name]
		if !exists {
			// Field was removed
			changes = append(changes, SchemaChange{
				Type:      "removed",
				FieldName: name,
				OldType:   oldField.Type,
			})
		} else if !arrowTypesEqual(oldField.Type, newField.Type) {
			// Field type changed
			changes = append(changes, SchemaChange{
				Type:      "type_changed",
				FieldName: name,
				OldType:   oldField.Type,
				NewType:   newField.Type,
			})
		}
	}

	// Find added fields
	for name, newField := range newFields {
		if _, exists := oldFields[name]; !exists {
			// Field was added
			changes = append(changes, SchemaChange{
				Type:      "added",
				FieldName: name,
				NewType:   newField.Type,
			})
		}
	}

	return changes
}

// arrowTypesEqual checks if two Arrow data types are equal
func arrowTypesEqual(a, b arrow.DataType) bool {
	return a.ID() == b.ID()
}

// FormatSchemaChanges returns a human-readable representation of schema changes
func FormatSchemaChanges(changes []SchemaChange) string {
	var sb strings.Builder

	if len(changes) == 0 {
		return "No schema changes detected."
	}

	for _, change := range changes {
		switch change.Type {
		case "added":
			sb.WriteString(fmt.Sprintf("+ Added field: %s (%s)\n", change.FieldName, formatArrowType(change.NewType)))
		case "removed":
			sb.WriteString(fmt.Sprintf("- Removed field: %s (%s)\n", change.FieldName, formatArrowType(change.OldType)))
		case "type_changed":
			sb.WriteString(fmt.Sprintf("~ Changed type: %s from %s to %s\n", 
				change.FieldName, formatArrowType(change.OldType), formatArrowType(change.NewType)))
		}
	}

	return sb.String()
}

// formatArrowType returns a human-readable representation of an Arrow data type
func formatArrowType(dt arrow.DataType) string {
	if dt == nil {
		return "unknown"
	}

	switch dt.ID() {
	case arrow.INT8:
		return "int8"
	case arrow.INT16:
		return "int16"
	case arrow.INT32:
		return "int32"
	case arrow.INT64:
		return "int64"
	case arrow.UINT8:
		return "uint8"
	case arrow.UINT16:
		return "uint16"
	case arrow.UINT32:
		return "uint32"
	case arrow.UINT64:
		return "uint64"
	case arrow.FLOAT32:
		return "float32"
	case arrow.FLOAT64:
		return "float64"
	case arrow.BOOL:
		return "boolean"
	case arrow.STRING:
		return "string"
	case arrow.BINARY:
		return "binary"
	case arrow.DATE32:
		return "date"
	case arrow.TIMESTAMP:
		return "timestamp"
	case arrow.DECIMAL128:
		return "decimal"
	case arrow.LIST:
		return "list"
	case arrow.STRUCT:
		return "struct"
	default:
		return fmt.Sprintf("other(%s)", dt.Name())
	}
}

// FormatSchemaHistory returns a human-readable representation of schema history
func FormatSchemaHistory(history *SchemaHistory) string {
	var sb strings.Builder

	sb.WriteString("Schema History:\n")
	sb.WriteString("==============\n\n")

	for i, version := range history.Versions {
		sb.WriteString(fmt.Sprintf("Version: %d\n", version.Version))
		sb.WriteString(fmt.Sprintf("Timestamp: %s\n", version.Timestamp.Format(time.RFC3339)))
		sb.WriteString("Fields:\n")

		for _, field := range version.SchemaFields {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", field.Name, formatArrowType(field.Type)))
		}

		if i < len(history.Versions)-1 {
			// Show diff with next version
			nextVersion := history.Versions[i+1]
			
			// We need to compare the current version with the previous version
			// to show what changed to get to this version
			changes := DiffSchemas(nextVersion.Schema, version.Schema)
			
			if len(changes) > 0 {
				sb.WriteString("\nChanges from previous version:\n")
				sb.WriteString(FormatSchemaChanges(changes))
			} else {
				sb.WriteString("\nNo schema changes from previous version.\n")
			}
		}

		sb.WriteString("\n")
	}

	return sb.String()
}
