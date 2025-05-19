package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeltaFormatHandler_ReadWrite(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-read-write-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(deltaTableDir, "_delta_log"), 0755))

	// Create test data
	testData := []map[string]interface{}{
		{"id": int64(1), "name": "John Doe", "age": int32(30), "active": true},
		{"id": int64(2), "name": "Jane Smith", "age": int32(25), "active": true},
		{"id": int64(3), "name": "Bob Johnson", "age": int32(40), "active": false},
	}

	// Create schema
	schema := &Schema{
		fields: []Field{
			{Name: "id", Type: FieldTypeInt64},
			{Name: "name", Type: FieldTypeString},
			{Name: "age", Type: FieldTypeInt32},
			{Name: "active", Type: FieldTypeBool},
		},
	}

	// Initialize Delta format handler
	handler := NewDeltaFormatHandler()

	// Write data to the Delta table
	err = handler.Write(deltaTableDir, testData, schema)
	require.NoError(t, err, "Should write data to Delta table without error")

	// Verify the Delta table structure
	assert.True(t, handler.IsDeltaTable(deltaTableDir), "Should be a valid Delta table after writing")
	
	// Check that _delta_log contains at least one transaction file
	files, err := os.ReadDir(filepath.Join(deltaTableDir, "_delta_log"))
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(files), 1, "Should have at least one transaction log file")
	
	// Check that at least one parquet file was created
	parquetFiles, err := filepath.Glob(filepath.Join(deltaTableDir, "*.parquet"))
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(parquetFiles), 1, "Should have at least one parquet file")

	// Read data from the Delta table
	readData, err := handler.Read(deltaTableDir)
	require.NoError(t, err, "Should read data from Delta table without error")
	
	// Verify the data was read correctly
	require.Len(t, readData, len(testData), "Should read the same number of records")
	
	// Verify specific fields in the data
	for i, record := range readData {
		assert.Equal(t, testData[i]["id"], record["id"], "ID should match")
		assert.Equal(t, testData[i]["name"], record["name"], "Name should match")
		assert.Equal(t, testData[i]["age"], record["age"], "Age should match")
		assert.Equal(t, testData[i]["active"], record["active"], "Active status should match")
	}
}

func TestDeltaFormatHandler_ReadWithInference(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "delta-inference-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock Delta Lake table
	deltaTableDir := filepath.Join(tempDir, "delta-table")
	require.NoError(t, os.MkdirAll(deltaTableDir, 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(deltaTableDir, "_delta_log"), 0755))

	// Create test data with various data types
	testData := []map[string]interface{}{
		{
			"id":          int64(1),
			"name":        "John Doe",
			"age":         int32(30),
			"salary":      float64(75000.50),
			"active":      true,
			"department":  "Engineering",
			"hired_date":  "2022-01-15",
		},
		{
			"id":          int64(2),
			"name":        "Jane Smith",
			"age":         int32(25),
			"salary":      float64(82000.75),
			"active":      true,
			"department":  "Marketing",
			"hired_date":  "2021-08-10",
		},
	}

	// Create schema
	schema := &Schema{
		fields: []Field{
			{Name: "id", Type: FieldTypeInt64},
			{Name: "name", Type: FieldTypeString},
			{Name: "age", Type: FieldTypeInt32},
			{Name: "salary", Type: FieldTypeFloat64},
			{Name: "active", Type: FieldTypeBool},
			{Name: "department", Type: FieldTypeString},
			{Name: "hired_date", Type: FieldTypeString},
		},
	}

	// Initialize Delta format handler
	handler := NewDeltaFormatHandler()

	// Write data to the Delta table
	err = handler.Write(deltaTableDir, testData, schema)
	require.NoError(t, err, "Should write data to Delta table without error")

	// Read data with schema inference
	readData, inferredSchema, err := handler.ReadWithInference(deltaTableDir)
	require.NoError(t, err, "Should read data with schema inference without error")
	
	// Verify the data was read correctly
	require.Len(t, readData, len(testData), "Should read the same number of records")
	
	// Verify the inferred schema
	require.NotNil(t, inferredSchema, "Should infer a schema")
	
	// Check that all fields are present in the inferred schema
	schemaFields := inferredSchema.Fields()
	fieldNames := make(map[string]bool)
	for _, field := range schemaFields {
		fieldNames[field.Name] = true
	}
	
	for _, field := range schema.fields {
		assert.True(t, fieldNames[field.Name], "Inferred schema should contain field: "+field.Name)
	}
	
	// Verify specific fields in the data
	for i, record := range readData {
		assert.Equal(t, testData[i]["id"], record["id"], "ID should match")
		assert.Equal(t, testData[i]["name"], record["name"], "Name should match")
		assert.Equal(t, testData[i]["age"], record["age"], "Age should match")
		assert.Equal(t, testData[i]["salary"], record["salary"], "Salary should match")
		assert.Equal(t, testData[i]["active"], record["active"], "Active status should match")
		assert.Equal(t, testData[i]["department"], record["department"], "Department should match")
		assert.Equal(t, testData[i]["hired_date"], record["hired_date"], "Hired date should match")
	}
}
