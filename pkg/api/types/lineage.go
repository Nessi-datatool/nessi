package types

// LineageInfo represents lineage information for a table
type LineageInfo struct {
	// Upstream represents tables that are read by this table
	Upstream []TableReference `json:"upstream"`
	
	// Downstream represents tables that write to this table
	Downstream []TableReference `json:"downstream"`
	
	// Properties contains additional metadata about the lineage
	Properties map[string]string `json:"properties,omitempty"`
}

// TableReference represents a reference to a table in a data catalog
type TableReference struct {
	// Catalog is the name of the data catalog
	Catalog string `json:"catalog"`
	
	// Database is the name of the database
	Database string `json:"database"`
	
	// Table is the name of the table
	Table string `json:"table"`
	
	// Properties contains additional metadata about the reference
	Properties map[string]string `json:"properties,omitempty"`
}
