package dbt

// This file contains mock implementations of the Delta Lake interfaces
// needed for testing the dbt plugin without requiring the actual Delta package.

// Mock Delta interfaces and types

type mockDeltaTable struct {
	path     string
	version  int64
	metadata *mockTableMetadata
	schema   *mockSchema
	stats    *mockTableStats
}

type mockTableMetadata struct {
	Format mockFormat
}

type mockFormat struct {
	Provider string
}

type mockSchema struct {
	Fields []*mockField
}

type mockField struct {
	Name     string
	Type     mockDataType
	Nullable bool
}

type mockDataType struct {
	typeName string
}

func (t mockDataType) String() string {
	return t.typeName
}

type mockTableStats struct {
	Files     []string
	SizeBytes int64
	NumRows   int64
}

type mockReader struct {
	table *mockDeltaTable
}

// Mock Delta functions to replace the actual Delta package

// OpenTable opens a Delta table at the given path
func delta_OpenTable(path string) (*mockDeltaTable, error) {
	// Create a mock table with some default values
	table := &mockDeltaTable{
		path:    path,
		version: 1,
		metadata: &mockTableMetadata{
			Format: mockFormat{
				Provider: "parquet",
			},
		},
		schema: &mockSchema{
			Fields: []*mockField{
				{
					Name:     "id",
					Type:     mockDataType{typeName: "integer"},
					Nullable: false,
				},
				{
					Name:     "name",
					Type:     mockDataType{typeName: "string"},
					Nullable: true,
				},
				{
					Name:     "value",
					Type:     mockDataType{typeName: "double"},
					Nullable: true,
				},
				{
					Name:     "created_at",
					Type:     mockDataType{typeName: "timestamp"},
					Nullable: false,
				},
			},
		},
		stats: &mockTableStats{
			Files:     []string{"file1.parquet", "file2.parquet"},
			SizeBytes: 1024 * 1024, // 1MB
			NumRows:   1000,
		},
	}
	return table, nil
}

// Close closes the table
func (t *mockDeltaTable) Close() error {
	return nil
}

// Version returns the table version
func (t *mockDeltaTable) Version() (int64, error) {
	return t.version, nil
}

// Metadata returns the table metadata
func (t *mockDeltaTable) Metadata() (*mockTableMetadata, error) {
	return t.metadata, nil
}

// Schema returns the table schema
func (t *mockDeltaTable) Schema() (*mockSchema, error) {
	return t.schema, nil
}

// Stats returns the table statistics
func (t *mockDeltaTable) Stats() (*mockTableStats, error) {
	return t.stats, nil
}

// NewReader creates a new reader for the table
func (t *mockDeltaTable) NewReader() (*mockReader, error) {
	return &mockReader{table: t}, nil
}

// Close closes the reader
func (r *mockReader) Close() error {
	return nil
}

// Replace Delta package functions with our mock implementations in the delta.go file

// These functions are used to replace the imports in delta.go
// They have the same signatures as the actual Delta package functions

func delta_OpenTable_replacement(path string) (*mockDeltaTable, error) {
	return delta_OpenTable(path)
}
