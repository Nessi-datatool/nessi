package datalake

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/nessi-dev/nessi/pkg/logging"
)

// SchemaManager handles schema operations for Delta Lake tables
type SchemaManager struct {
	tablePath       string
	metadataManager *MetadataManager
	mutex           sync.RWMutex
	historyCache    *SchemaHistory // Cache for schema history
}

// NewSchemaManager creates a new schema manager for a Delta Lake table
func NewSchemaManager(tablePath string) *SchemaManager {
	logger := logging.GetLogger()
	logger.Debug("Creating new SchemaManager", "tablePath", tablePath)

	return &SchemaManager{
		tablePath:       tablePath,
		metadataManager: NewMetadataManager(tablePath),
		historyCache:    nil,
	}
}

// SchemaFieldChange represents a change between two schema versions
type SchemaFieldChange struct {
	Type      string         // "added", "removed", "type_changed"
	FieldName string         // Name of the field that changed
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
	Version   int64               `json:"version"`
	Timestamp time.Time           `json:"timestamp"`
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
func (sm *SchemaManager) GetSchemaHistory() (*SchemaHistory, error) {
	logger := logging.GetLogger()
	logger.Debug("Getting schema history", "tablePath", sm.tablePath)

	// Check cache first
	sm.mutex.RLock()
	if sm.historyCache != nil {
		deferred := sm.historyCache
		sm.mutex.RUnlock()
		logger.Debug("Using cached schema history", "versionCount", len(deferred.Versions))
		return deferred, nil
	}
	sm.mutex.RUnlock()

	versions, err := sm.metadataManager.GetVersions()
	if err != nil {
		logger.Error("Failed to get versions", "error", err)
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	if len(versions) == 0 {
		logger.Warn("No versions found for table")
		emptyHistory := &SchemaHistory{Versions: []SchemaVersion{}}

		// Cache the empty history
		sm.mutex.Lock()
		sm.historyCache = emptyHistory
		sm.mutex.Unlock()

		return emptyHistory, nil
	}

	history := &SchemaHistory{
		Versions: make([]SchemaVersion, 0, len(versions)),
	}

	for _, version := range versions {
		table, err := sm.metadataManager.GetTableAtVersion(version)
		if err != nil {
			logger.Error("Failed to get table at version", "error", err, "version", version)
			return nil, fmt.Errorf("failed to get table at version %d: %w", version, err)
		}

		history.Versions = append(history.Versions, SchemaVersion{
			Version:      version,
			Timestamp:    table.LastModified,
			Schema:       table.Schema,
			SchemaFields: table.Schema.Fields(),
		})
	}

	// Cache the result
	sm.mutex.Lock()
	sm.historyCache = history
	sm.mutex.Unlock()

	logger.Debug("Successfully retrieved schema history", "versionCount", len(history.Versions))
	return history, nil
}

// DiffSchemas compares two schemas and returns the differences
func DiffSchemas(oldSchema, newSchema *arrow.Schema) []SchemaFieldChange {
	logger := logging.GetLogger()
	logger.Debug("Comparing schemas for differences")

	// Handle nil schemas
	if oldSchema == nil && newSchema == nil {
		logger.Warn("Both schemas are nil, no differences to report")
		return []SchemaFieldChange{}
	}

	if oldSchema == nil {
		// All fields in new schema are additions
		logger.Debug("Old schema is nil, all fields in new schema are additions")
		changes := make([]SchemaFieldChange, 0, newSchema.NumFields())
		for _, field := range newSchema.Fields() {
			changes = append(changes, SchemaFieldChange{
				Type:      "added",
				FieldName: field.Name,
				NewType:   field.Type,
			})
		}
		return changes
	}

	if newSchema == nil {
		// All fields in old schema are removals
		logger.Debug("New schema is nil, all fields in old schema are removals")
		changes := make([]SchemaFieldChange, 0, oldSchema.NumFields())
		for _, field := range oldSchema.Fields() {
			changes = append(changes, SchemaFieldChange{
				Type:      "removed",
				FieldName: field.Name,
				OldType:   field.Type,
			})
		}
		return changes
	}

	// Pre-allocate capacity for changes
	changes := make([]SchemaFieldChange, 0, max(oldSchema.NumFields(), newSchema.NumFields()))

	// Create maps for quick lookup
	oldFields := make(map[string]arrow.Field, oldSchema.NumFields())
	for _, field := range oldSchema.Fields() {
		oldFields[field.Name] = field
	}

	newFields := make(map[string]arrow.Field, newSchema.NumFields())
	for _, field := range newSchema.Fields() {
		newFields[field.Name] = field
	}

	// Find removed and changed fields
	for name, oldField := range oldFields {
		newField, exists := newFields[name]
		if !exists {
			// Field was removed
			logger.Debug("Field was removed", "field", name)
			changes = append(changes, SchemaFieldChange{
				Type:      "removed",
				FieldName: oldField.Name,
				OldType:   oldField.Type,
			})
		} else if !arrowTypesEqual(oldField.Type, newField.Type) {
			// Field type changed
			logger.Debug("Field type changed", "field", name,
				"oldType", formatArrowType(oldField.Type),
				"newType", formatArrowType(newField.Type))
			changes = append(changes, SchemaFieldChange{
				Type:      "type_changed",
				FieldName: oldField.Name,
				OldType:   oldField.Type,
				NewType:   newField.Type,
			})
		}
	}

	// Find added fields
	for name, newField := range newFields {
		if _, exists := oldFields[name]; !exists {
			// Field was added
			logger.Debug("Field was added", "field", name)
			changes = append(changes, SchemaFieldChange{
				Type:      "added",
				FieldName: newField.Name,
				NewType:   newField.Type,
			})
		}
	}

	logger.Debug("Schema comparison complete", "changeCount", len(changes))
	return changes
}

// arrowTypesEqual checks if two Arrow data types are equal
func arrowTypesEqual(a, b arrow.DataType) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Check type IDs
	if a.ID() != b.ID() {
		return false
	}

	// For some types, we need more detailed comparison
	switch a.ID() {
	case arrow.DECIMAL128:
		// For decimal types, check precision and scale
		dec1, ok1 := a.(*arrow.Decimal128Type)
		dec2, ok2 := b.(*arrow.Decimal128Type)
		if ok1 && ok2 {
			return dec1.Precision == dec2.Precision && dec1.Scale == dec2.Scale
		}
		return false

	case arrow.TIMESTAMP:
		// For timestamp types, check unit
		ats, ok1 := a.(*arrow.TimestampType)
		bts, ok2 := b.(*arrow.TimestampType)
		if !ok1 || !ok2 {
			return false
		}
		return ats.Unit == bts.Unit
	}

	// For other types, type ID equality is sufficient
	return true
}

// FormatSchemaChanges returns a human-readable representation of schema changes
func FormatSchemaChanges(changes []SchemaFieldChange) string {
	logger := logging.GetLogger()
	logger.Debug("Formatting schema changes", "changeCount", len(changes))

	if changes == nil || len(changes) == 0 {
		return "No schema changes detected."
	}

	// Pre-allocate a reasonable buffer size based on the number of changes
	// Assume average of 50 bytes per change
	var sb strings.Builder
	sb.Grow(len(changes) * 50)

	// Sort changes for consistent output
	sortedChanges := make([]SchemaFieldChange, len(changes))
	copy(sortedChanges, changes)
	sort.Slice(sortedChanges, func(i, j int) bool {
		// First sort by type: added, removed, type_changed
		if sortedChanges[i].Type != sortedChanges[j].Type {
			// Custom ordering: added first, then removed, then type_changed
			order := map[string]int{"added": 0, "removed": 1, "type_changed": 2}
			return order[sortedChanges[i].Type] < order[sortedChanges[j].Type]
		}
		// Then sort by field name
		return sortedChanges[i].FieldName < sortedChanges[j].FieldName
	})

	for _, change := range sortedChanges {
		switch change.Type {
		case "added":
			sb.WriteString(fmt.Sprintf("+ Added field: %s (%s)\n", change.FieldName, formatArrowType(change.NewType)))
		case "removed":
			sb.WriteString(fmt.Sprintf("- Removed field: %s (%s)\n", change.FieldName, formatArrowType(change.OldType)))
		case "type_changed":
			sb.WriteString(fmt.Sprintf("~ Changed type: %s from %s to %s\n",
				change.FieldName, formatArrowType(change.OldType), formatArrowType(change.NewType)))
		default:
			logger.Warn("Unknown change type", "type", change.Type, "field", change.FieldName)
			sb.WriteString(fmt.Sprintf("? Unknown change: %s (%s)\n", change.FieldName, change.Type))
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
	case arrow.LARGE_STRING:
		return "large_string"
	case arrow.BINARY:
		return "binary"
	case arrow.LARGE_BINARY:
		return "large_binary"
	case arrow.DATE32:
		return "date32"
	case arrow.DATE64:
		return "date64"
	case arrow.TIMESTAMP:
		// Include unit information for timestamp
		ts, ok := dt.(*arrow.TimestampType)
		if ok {
			return fmt.Sprintf("timestamp[%s]", ts.Unit)
		}
		return "timestamp"
	case arrow.TIME32:
		return "time32"
	case arrow.TIME64:
		return "time64"
	case arrow.DECIMAL128:
		if dec, ok := dt.(*arrow.Decimal128Type); ok {
			return fmt.Sprintf("decimal128(%d,%d)", dec.Precision, dec.Scale)
		}
		return "decimal128"
	case arrow.LIST:
		list, ok := dt.(*arrow.ListType)
		if ok && list.Elem() != nil {
			return fmt.Sprintf("list<%s>", formatArrowType(list.Elem()))
		}
		return "list"
	case arrow.STRUCT:
		return "struct"
	case arrow.MAP:
		return "map"
	case arrow.DICTIONARY:
		return "dictionary"
	case arrow.INTERVAL_MONTHS, arrow.INTERVAL_DAY_TIME, arrow.INTERVAL_MONTH_DAY_NANO:
		return "interval"
	case arrow.FIXED_SIZE_BINARY:
		return "fixed_size_binary"
	case arrow.DURATION:
		return "duration"
	default:
		return fmt.Sprintf("other(%s)", dt.Name())
	}
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ConvertToArrowSchema converts a DeltaSchema to an Arrow Schema
func ConvertToArrowSchema(schema *DeltaSchema) (*arrow.Schema, []arrow.Field, error) {
	logger := logging.GetLogger()
	logger.Debug("Converting Schema to Arrow Schema", "fields", len(schema.Fields))

	if schema == nil {
		return nil, nil, fmt.Errorf("cannot convert nil schema to arrow schema")
	}

	// Pre-allocate the slice for better performance
	arrowFields := make([]arrow.Field, len(schema.Fields))
	for i, field := range schema.Fields {
		// Convert field type to arrow data type
		var dataType arrow.DataType
		switch field.Type {
		case "integer":
			dataType = arrow.PrimitiveTypes.Int64
		case "float":
			dataType = arrow.PrimitiveTypes.Float64
		case "boolean":
			dataType = arrow.FixedWidthTypes.Boolean
		case "timestamp":
			dataType = arrow.FixedWidthTypes.Timestamp_s
		case "date":
			dataType = arrow.FixedWidthTypes.Date32
		case "binary":
			dataType = arrow.BinaryTypes.Binary
		default:
			// Default to string for unknown types
			dataType = arrow.BinaryTypes.String
		}

		// Create metadata map if needed
		var metadata map[string]string
		if field.Metadata != nil && len(field.Metadata) > 0 {
			metadata = field.Metadata
		}

		// Create arrow field
		arrowFields[i] = arrow.Field{
			Name:     field.Name,
			Type:     dataType,
			Nullable: field.Nullable,
			Metadata: arrow.MetadataFrom(metadata),
		}
	}

	// Create arrow schema
	arrowSchema := arrow.NewSchema(arrowFields, nil)
	logger.Debug("Successfully converted Schema to Arrow Schema", "fields", len(arrowFields))

	return arrowSchema, arrowFields, nil
}

// ConvertFromArrowSchema converts an Arrow Schema to a DeltaSchema
func ConvertFromArrowSchema(arrowSchema *arrow.Schema) (*DeltaSchema, error) {
	logger := logging.GetLogger()
	logger.Debug("Converting Arrow Schema to Schema", "fields", arrowSchema.NumFields())

	if arrowSchema == nil {
		return nil, fmt.Errorf("cannot convert nil arrow schema to schema")
	}

	// Create schema fields
	fields := make([]DeltaField, arrowSchema.NumFields())
	for i, arrowField := range arrowSchema.Fields() {
		// Convert arrow data type to field type
		var fieldType string
		switch arrowField.Type.ID() {
		case arrow.INT8, arrow.INT16, arrow.INT32, arrow.INT64, arrow.UINT8, arrow.UINT16, arrow.UINT32, arrow.UINT64:
			fieldType = "integer"
		case arrow.FLOAT32, arrow.FLOAT64:
			fieldType = "float"
		case arrow.BOOL:
			fieldType = "boolean"
		case arrow.TIMESTAMP:
			fieldType = "timestamp"
		case arrow.DATE32, arrow.DATE64:
			fieldType = "date"
		case arrow.BINARY, arrow.LARGE_BINARY:
			fieldType = "binary"
		default:
			// Default to string for unknown types
			fieldType = "string"
		}

		// Extract metadata
		metadata := make(map[string]string)
		if arrowField.Metadata.Len() > 0 {
			for _, key := range arrowField.Metadata.Keys() {
				metadata[key] = arrowField.Metadata.Values()[arrowField.Metadata.FindKey(key)]
			}
		}

		// Create field
		fields[i] = DeltaField{
			Name:     arrowField.Name,
			Type:     fieldType,
			Nullable: arrowField.Nullable,
			Metadata: metadata,
		}
	}

	// Create schema
	schema := &DeltaSchema{
		Fields: fields,
	}

	logger.Debug("Successfully converted Arrow Schema to Schema", "fields", len(fields))
	return schema, nil
}

// FormatSchemaHistory returns a human-readable representation of schema history
func FormatSchemaHistory(history *SchemaHistory) string {
	logger := logging.GetLogger()
	logger.Debug("Formatting schema history")

	if history == nil {
		logger.Warn("Schema history is nil")
		return "No schema history available."
	}

	if len(history.Versions) == 0 {
		logger.Debug("Schema history has no versions")
		return "Schema History: No versions available."
	}

	// Pre-allocate a reasonable buffer size based on the number of versions and fields
	// Estimate 200 bytes per version plus 50 bytes per field
	totalFields := 0
	for _, version := range history.Versions {
		totalFields += len(version.SchemaFields)
	}

	var sb strings.Builder
	sb.Grow(len(history.Versions)*200 + totalFields*50)

	sb.WriteString("Schema History:\n")
	sb.WriteString("==============\n\n")

	// Process versions in reverse order (newest first)
	for i := len(history.Versions) - 1; i >= 0; i-- {
		version := history.Versions[i]
		sb.WriteString(fmt.Sprintf("Version: %d\n", version.Version))
		sb.WriteString(fmt.Sprintf("Timestamp: %s\n", version.Timestamp.Format(time.RFC3339)))
		sb.WriteString("Fields:\n")

		// Sort fields by name for consistent output
		sortedFields := make([]arrow.Field, len(version.SchemaFields))
		copy(sortedFields, version.SchemaFields)
		sort.Slice(sortedFields, func(i, j int) bool {
			return sortedFields[i].Name < sortedFields[j].Name
		})

		for _, field := range sortedFields {
			// Check if field has description metadata
			description := ""
			if field.Metadata.Len() > 0 {
				for i := 0; i < field.Metadata.Len(); i++ {
					if field.Metadata.Keys()[i] == "description" {
						description = field.Metadata.Values()[i]
						break
					}
				}
			}

			// Include description if available
			if description != "" {
				sb.WriteString(fmt.Sprintf("  - %s (%s): %s\n",
					field.Name, formatArrowType(field.Type), description))
			} else {
				sb.WriteString(fmt.Sprintf("  - %s (%s)\n",
					field.Name, formatArrowType(field.Type)))
			}
		}

		if i > 0 {
			// Show diff with previous version
			prevVersion := history.Versions[i-1]

			// Compare the current version with the previous version
			changes := DiffSchemas(prevVersion.Schema, version.Schema)

			if len(changes) > 0 {
				sb.WriteString("\nChanges from previous version:\n")
				sb.WriteString(FormatSchemaChanges(changes))
			} else {
				sb.WriteString("\nNo schema changes from previous version.\n")
			}
		}

		sb.WriteString("\n")
	}

	logger.Debug("Successfully formatted schema history", "versionCount", len(history.Versions))
	return sb.String()
}
