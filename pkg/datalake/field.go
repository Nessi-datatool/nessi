package datalake

// DeltaField represents a field in a Delta Lake schema
type DeltaField struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Nullable bool              `json:"nullable"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DeltaSchema represents a collection of fields in a Delta Lake table
type DeltaSchema struct {
	Fields []DeltaField `json:"fields"`
}
