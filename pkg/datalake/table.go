package datalake

// Table represents a Delta table
type Table struct {
	// The version of the table
	Version int64

	// The files in the table
	Files []string

	// The partitions in the table
	Partitions map[string][]string

	// Additional metadata
	Metadata map[string]interface{}
}

// SchemaField represents a field in a Delta table schema
type SchemaField struct {
	// The name of the field
	Name string

	// The type of the field
	Type string

	// Whether the field is nullable
	Nullable bool
}
