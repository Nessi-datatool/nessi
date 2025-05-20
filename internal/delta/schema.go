package delta

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	memory "github.com/apache/arrow/go/v15/arrow/memory"
)

// SchemaManager manages Delta Lake schemas
type SchemaManager struct {
	schema *arrow.Schema
}

// NewSchemaManager creates a new schema manager
func NewSchemaManager(schema *arrow.Schema) *SchemaManager {
	return &SchemaManager{
		schema: schema,
	}
}

// GetSchema returns the current schema
func (s *SchemaManager) GetSchema() *arrow.Schema {
	return s.schema
}

// ValidateSchema validates a schema against the current schema
func (s *SchemaManager) ValidateSchema(schema *arrow.Schema) error {
	if s.schema == nil {
		return nil
	}

	// Check if all required fields are present
	for i := 0; i < len(s.schema.Fields()); i++ {
		field := s.schema.Fields()[i]
		if !field.Nullable {
			found := false
			for j := 0; j < len(schema.Fields()); j++ {
				if schema.Fields()[j].Name == field.Name {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("required field missing: %s", field.Name)
			}
		}
	}

	// Check if field types are compatible
	for i := 0; i < len(schema.Fields()); i++ {
		field := schema.Fields()[i]
		for j := 0; j < len(s.schema.Fields()); j++ {
			if s.schema.Fields()[j].Name == field.Name {
				if !isTypeCompatible(s.schema.Fields()[j].Type, field.Type) {
					return fmt.Errorf("incompatible type for field %s: expected %s, got %s",
						field.Name, s.schema.Fields()[j].Type, field.Type)
				}
				break
			}
		}
	}

	return nil
}

// isTypeCompatible checks if two Arrow types are compatible
func isTypeCompatible(expected, actual arrow.DataType) bool {
	if expected.ID() == actual.ID() {
		return true
	}

	// Handle numeric type promotions
	switch expected.ID() {
	case arrow.INT8, arrow.INT16, arrow.INT32, arrow.INT64:
		switch actual.ID() {
		case arrow.INT8, arrow.INT16, arrow.INT32, arrow.INT64:
			return true
		}
	case arrow.UINT8, arrow.UINT16, arrow.UINT32, arrow.UINT64:
		switch actual.ID() {
		case arrow.UINT8, arrow.UINT16, arrow.UINT32, arrow.UINT64:
			return true
		}
	case arrow.FLOAT32, arrow.FLOAT64:
		switch actual.ID() {
		case arrow.FLOAT32, arrow.FLOAT64:
			return true
		}
	}

	return false
}

// MergeSchemas merges two schemas
func (s *SchemaManager) MergeSchemas(schema *arrow.Schema) (*arrow.Schema, error) {
	if s.schema == nil {
		return schema, nil
	}

	// Create a map of existing fields
	fieldMap := make(map[string]arrow.Field)
	for i := 0; i < len(s.schema.Fields()); i++ {
		field := s.schema.Fields()[i]
		fieldMap[field.Name] = field
	}

	// Add new fields
	for i := 0; i < len(schema.Fields()); i++ {
		field := schema.Fields()[i]
		if existing, exists := fieldMap[field.Name]; exists {
			// Check type compatibility
			if !isTypeCompatible(existing.Type, field.Type) {
				return nil, fmt.Errorf("incompatible type for field %s: existing %s, new %s",
					field.Name, existing.Type, field.Type)
			}
			// Use the new field if it's not nullable
			if !field.Nullable {
				fieldMap[field.Name] = field
			}
		} else {
			fieldMap[field.Name] = field
		}
	}

	// Create the merged schema
	fields := make([]arrow.Field, 0, len(fieldMap))
	for _, field := range fieldMap {
		fields = append(fields, field)
	}

	return arrow.NewSchema(fields, nil), nil
}

// SchemaToJSON converts a schema to JSON
func (s *SchemaManager) SchemaToJSON() (string, error) {
	if s.schema == nil {
		return "{}", nil
	}

	type Field struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Nullable bool   `json:"nullable"`
	}

	type Schema struct {
		Fields []Field `json:"fields"`
	}

	schema := Schema{
		Fields: make([]Field, len(s.schema.Fields())),
	}

	for i := 0; i < len(s.schema.Fields()); i++ {
		field := s.schema.Field(i)
		schema.Fields[i] = Field{
			Name:     field.Name,
			Type:     field.Type.String(),
			Nullable: field.Nullable,
		}
	}

	data, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	return string(data), nil
}

// JSONToSchema converts JSON to a schema
func (s *SchemaManager) JSONToSchema(jsonStr string) (*arrow.Schema, error) {
	type Field struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		Nullable bool   `json:"nullable"`
	}

	type Schema struct {
		Fields []Field `json:"fields"`
	}

	var schema Schema
	if err := json.Unmarshal([]byte(jsonStr), &schema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	fields := make([]arrow.Field, len(schema.Fields))
	for i, field := range schema.Fields {
		arrowType, err := parseArrowType(field.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to parse type for field %s: %w", field.Name, err)
		}
		fields[i] = arrow.Field{
			Name:     field.Name,
			Type:     arrowType,
			Nullable: field.Nullable,
		}
	}

	return arrow.NewSchema(fields, nil), nil
}

// parseArrowType parses an Arrow type string
func parseArrowType(typeStr string) (arrow.DataType, error) {
	typeStr = strings.ToLower(typeStr)
	switch typeStr {
	case "int8":
		return arrow.PrimitiveTypes.Int8, nil
	case "int16":
		return arrow.PrimitiveTypes.Int16, nil
	case "int32":
		return arrow.PrimitiveTypes.Int32, nil
	case "int64":
		return arrow.PrimitiveTypes.Int64, nil
	case "uint8":
		return arrow.PrimitiveTypes.Uint8, nil
	case "uint16":
		return arrow.PrimitiveTypes.Uint16, nil
	case "uint32":
		return arrow.PrimitiveTypes.Uint32, nil
	case "uint64":
		return arrow.PrimitiveTypes.Uint64, nil
	case "float32":
		return arrow.PrimitiveTypes.Float32, nil
	case "float64":
		return arrow.PrimitiveTypes.Float64, nil
	case "string":
		return arrow.BinaryTypes.String, nil
	case "binary":
		return arrow.BinaryTypes.Binary, nil
	case "boolean":
		return arrow.FixedWidthTypes.Boolean, nil
	case "date32":
		return arrow.FixedWidthTypes.Date32, nil
	case "date64":
		return arrow.FixedWidthTypes.Date64, nil
	case "timestamp":
		return arrow.FixedWidthTypes.Timestamp_ns, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", typeStr)
	}
}

// ValidateRecord validates a record against the schema
func (s *SchemaManager) ValidateRecord(record arrow.Record) error {
	if s.schema == nil {
		return nil
	}

	if int(record.NumCols()) != len(s.schema.Fields()) {
		return fmt.Errorf("record has %d columns, schema has %d fields",
			record.NumCols(), len(s.schema.Fields()))
	}

	for i := 0; i < len(s.schema.Fields()); i++ {
		field := s.schema.Fields()[i]
		col := record.Column(i)

		// Check if the column type matches the field type
		if !isTypeCompatible(field.Type, col.DataType()) {
			return fmt.Errorf("incompatible type for column %s: expected %s, got %s",
				field.Name, field.Type, col.DataType())
		}

		// Check for null values in non-nullable fields
		if !field.Nullable {
			for j := 0; j < col.Len(); j++ {
				if col.IsNull(j) {
					return fmt.Errorf("null value found in non-nullable field %s at row %d",
						field.Name, j)
				}
			}
		}
	}

	return nil
}

// CreateEmptyRecord creates an empty record with the schema
func (s *SchemaManager) CreateEmptyRecord() arrow.Record {
	if s.schema == nil {
		return nil
	}

	// Create empty arrays for each field
	arrays := make([]arrow.Array, len(s.schema.Fields()))
	for i := 0; i < len(s.schema.Fields()); i++ {
		field := s.schema.Field(i)
		builder := array.NewBuilder(memory.DefaultAllocator, field.Type)
		arrays[i] = builder.NewArray()
	}

	return array.NewRecord(s.schema, arrays, 0)
}
