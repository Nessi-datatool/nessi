package datalake

// FieldType represents the type of a field in a schema
type FieldType string

const (
	// FieldTypeString represents a string field
	FieldTypeString FieldType = "string"

	// FieldTypeInt32 represents a 32-bit integer field
	FieldTypeInt32 FieldType = "int32"

	// FieldTypeInt64 represents a 64-bit integer field
	FieldTypeInt64 FieldType = "int64"

	// FieldTypeFloat32 represents a 32-bit floating point field
	FieldTypeFloat32 FieldType = "float32"

	// FieldTypeFloat64 represents a 64-bit floating point field
	FieldTypeFloat64 FieldType = "float64"

	// FieldTypeBool represents a boolean field
	FieldTypeBool FieldType = "bool"

	// FieldTypeTimestamp represents a timestamp field
	FieldTypeTimestamp FieldType = "timestamp"

	// FieldTypeDate represents a date field
	FieldTypeDate FieldType = "date"

	// FieldTypeBinary represents a binary field
	FieldTypeBinary FieldType = "binary"
)

// Field represents a field in a schema
type Field struct {
	Name string
	Type FieldType
}

// Schema represents a schema for a data source
type Schema struct {
	fields []Field
}

// NewSchema creates a new schema with the specified fields
func NewSchema(fields []Field) *Schema {
	return &Schema{
		fields: fields,
	}
}

// Fields returns the fields in the schema
func (s *Schema) Fields() []Field {
	return s.fields
}

// AddField adds a field to the schema
func (s *Schema) AddField(field Field) {
	s.fields = append(s.fields, field)
}

// GetField returns a field by name
func (s *Schema) GetField(name string) (Field, bool) {
	for _, field := range s.fields {
		if field.Name == name {
			return field, true
		}
	}
	return Field{}, false
}
