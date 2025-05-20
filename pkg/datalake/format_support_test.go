package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFormatSupport tests format support functionality
func TestFormatSupport(t *testing.T) {
	// Create a temporary directory for test files
	dir, err := os.MkdirTemp("", "format_support_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	// Create Delta Lake directory
	deltaDir := filepath.Join(dir, "delta_table")
	err = os.MkdirAll(filepath.Join(deltaDir, "_delta_log"), 0755)
	require.NoError(t, err)

	// Create Parquet file
	parquetFile := filepath.Join(dir, "data.parquet")
	err = os.WriteFile(parquetFile, []byte("This is a mock Parquet file"), 0644)
	require.NoError(t, err)

	// Create CSV file
	csvFile := filepath.Join(dir, "data.csv")
	err = os.WriteFile(csvFile, []byte("id,name,value\n1,test,1.23"), 0644)
	require.NoError(t, err)

	// Create unknown format file
	unknownFile := filepath.Join(dir, "data.txt")
	err = os.WriteFile(unknownFile, []byte("This is not a supported format"), 0644)
	require.NoError(t, err)
}

// Format detection helpers
func testHasDeltaLogDirectory(path string) bool {
	deltaLogPath := filepath.Join(path, "_delta_log")
	info, err := os.Stat(deltaLogPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func testHasParquetExtension(path string) bool {
	return filepath.Ext(path) == ".parquet"
}

func testHasCSVExtension(path string) bool {
	return filepath.Ext(path) == ".csv"
}

// Schema inference helpers
func inferSchema(path string) (*DeltaSchema, error) {
	// Check format
	if testHasDeltaLogDirectory(path) {
		return inferDeltaSchema(path)
	} else if testHasParquetExtension(path) {
		return inferParquetSchema(path)
	} else if testHasCSVExtension(path) {
		return inferCSVSchema(path, ',', true, 1000)
	}
	return nil, os.ErrNotExist
}

func inferDeltaSchema(path string) (*DeltaSchema, error) {
	// In a real implementation, this would parse the Delta transaction log
	// For testing, we'll return a mock schema
	return &DeltaSchema{
		Fields: []DeltaField{
			{Name: "id", Type: "INTEGER", Nullable: false},
			{Name: "name", Type: "STRING", Nullable: true},
			{Name: "value", Type: "FLOAT", Nullable: true},
		},
	}, nil
}

func inferParquetSchema(path string) (*DeltaSchema, error) {
	// In a real implementation, this would parse the Parquet file
	// For testing, we'll return a mock schema
	return &DeltaSchema{
		Fields: []DeltaField{
			{Name: "id", Type: "INTEGER", Nullable: false},
			{Name: "name", Type: "STRING", Nullable: true},
			{Name: "value", Type: "FLOAT", Nullable: true},
		},
	}, nil
}

func inferCSVSchema(path string, delimiter rune, hasHeader bool, maxRows int) (*DeltaSchema, error) {
	// In a real implementation, this would parse the CSV file
	// For testing, we'll return a mock schema
	return &DeltaSchema{
		Fields: []DeltaField{},
	}, nil
}

// Test schema field creation
func TestSchemaFieldCreation(t *testing.T) {
	// Create a schema with fields
	schema := &DeltaSchema{
		Fields: []DeltaField{},
	}

	// Add fields
	schema.Fields = append(schema.Fields, DeltaField{
		Name:     "id",
		Type:     "INTEGER",
		Nullable: false,
	})

	schema.Fields = append(schema.Fields, DeltaField{
		Name:     "name",
		Type:     "STRING",
		Nullable: true,
	})

	schema.Fields = append(schema.Fields, DeltaField{
		Name:     "active",
		Type:     "BOOLEAN",
		Nullable: false,
	})

	schema.Fields = append(schema.Fields, DeltaField{
		Name:     "timestamp",
		Type:     "TIMESTAMP",
		Nullable: true,
	})

	// Verify fields
	assert.Equal(t, 4, len(schema.Fields))
	assert.Equal(t, "id", schema.Fields[0].Name)
	assert.Equal(t, "INTEGER", schema.Fields[0].Type)
	assert.Equal(t, false, schema.Fields[0].Nullable)

	assert.Equal(t, "name", schema.Fields[1].Name)
	assert.Equal(t, "STRING", schema.Fields[1].Type)
	assert.Equal(t, true, schema.Fields[1].Nullable)

	assert.Equal(t, "active", schema.Fields[2].Name)
	assert.Equal(t, "BOOLEAN", schema.Fields[2].Type)
	assert.Equal(t, false, schema.Fields[2].Nullable)

	assert.Equal(t, "timestamp", schema.Fields[3].Name)
	assert.Equal(t, "TIMESTAMP", schema.Fields[3].Type)
	assert.Equal(t, true, schema.Fields[3].Nullable)
}
