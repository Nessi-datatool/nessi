package datalake

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/stretchr/testify/require"
)

func serializeSchema(schema *arrow.Schema) map[string]interface{} {
	schemaFields := make([]map[string]interface{}, len(schema.Fields()))
	for i, field := range schema.Fields() {
		var typeStr string
		switch field.Type.(type) {
		case *arrow.Int32Type:
			typeStr = "int32"
		case *arrow.StringType:
			typeStr = "string"
		case *arrow.Float64Type:
			typeStr = "double"
		default:
			typeStr = field.Type.Name()
		}
		schemaFields[i] = map[string]interface{}{
			"name": field.Name,
			"type": typeStr,
		}
	}
	return map[string]interface{}{
		"fields": schemaFields,
	}
}

func setupTestTable(t *testing.T) (string, *arrow.Schema) {
	// Create temp dir
	tempDir, err := os.MkdirTemp("", "delta-test*")
	require.NoError(t, err)

	// Create _delta_log directory
	err = os.MkdirAll(filepath.Join(tempDir, "_delta_log"), 0755)
	require.NoError(t, err)

	// Create schema
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: &arrow.Int32Type{}},
		{Name: "name", Type: &arrow.StringType{}},
		{Name: "value", Type: &arrow.Float64Type{}},
	}, nil)

	// Create metadata
	schemaFields := make([]map[string]interface{}, len(schema.Fields()))
	for i, field := range schema.Fields() {
		schemaFields[i] = map[string]interface{}{
			"name": field.Name,
			"type": field.Type.String(),
		}
	}

	metadata := map[string]interface{}{
		"version":    int64(0),
		"timestamp":  time.Now().Unix(),
		"schema":     serializeSchema(schema),
		"files":      []string{},
		"partitions": map[string][]string{},
		"stats":      &TableStats{},
		"metadata":   map[string]interface{}{},
	}

	// Write metadata
	data, err := json.Marshal(metadata)
	require.NoError(t, err)

	logFile := filepath.Join(tempDir, "_delta_log", "00000000000000000000.json")
	err = os.WriteFile(logFile, data, 0644)
	require.NoError(t, err)

	return tempDir, schema
}

func TestNewReader(t *testing.T) {
	// Setup test table
	tempDir, _ := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}

func TestInitialize(t *testing.T) {
	// Setup test table
	tempDir, schema := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Test initialization
	err = reader.Initialize()
	require.NoError(t, err)

	// Verify schema
	actualSchema := reader.GetSchema()
	expected := serializeSchema(schema)
	actual := serializeSchema(actualSchema)
	require.Equal(t, expected, actual)

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}

func TestReadPartition(t *testing.T) {
	// Setup test table
	tempDir, schema := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create test partition
	partitionPath := filepath.Join(tempDir, "test")
	err := os.MkdirAll(partitionPath, 0755)
	require.NoError(t, err)

	// Create test file
	filePath := filepath.Join(partitionPath, "test.parquet")
	testData := []byte("test data")
	err = os.WriteFile(filePath, testData, 0644)
	require.NoError(t, err)

	// Update metadata
	metadata := map[string]interface{}{
		"version":    int64(1),
		"timestamp":  time.Now().Unix(),
		"schema":     serializeSchema(schema),
		"files":      []string{filePath},
		"partitions": map[string][]string{"test": {filePath}},
		"stats":      &TableStats{},
		"metadata":   map[string]interface{}{},
	}

	data, err := json.Marshal(metadata)
	require.NoError(t, err)

	logFile := filepath.Join(tempDir, "_delta_log", "00000000000000000001.json")
	err = os.WriteFile(logFile, data, 0644)
	require.NoError(t, err)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Test reading partition
	record, err := reader.ReadPartition("test")
	require.NoError(t, err)
	require.NotNil(t, record)

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}

func TestGetStats(t *testing.T) {
	// Setup test table
	tempDir, _ := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create test file
	filePath := filepath.Join(tempDir, "test.parquet")
	testData := []byte("test data")
	err := os.WriteFile(filePath, testData, 0644)
	require.NoError(t, err)

	// Update metadata
	metadata := map[string]interface{}{
		"version":    int64(1),
		"timestamp":  time.Now().Unix(),
		"schema":     map[string]interface{}{
			"fields": []map[string]interface{}{
				{"name": "id", "type": "int32"},
				{"name": "name", "type": "utf8"},
				{"name": "value", "type": "float64"},
			},
		},
		"files":      []string{filePath},
		"partitions": map[string][]string{},
		"stats":      &TableStats{NumFiles: 1, NumRecords: 100, TotalSize: 1024},
		"metadata":   map[string]interface{}{},
	}

	data, err := json.Marshal(metadata)
	require.NoError(t, err)

	logFile := filepath.Join(tempDir, "_delta_log", "00000000000000000001.json")
	err = os.WriteFile(logFile, data, 0644)
	require.NoError(t, err)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Test getting stats
	stats, err := reader.GetStats()
	require.NoError(t, err)
	require.NotNil(t, stats)
	require.Equal(t, int64(1), stats.NumFiles)
	require.Equal(t, int64(100), stats.NumRecords)
	require.Equal(t, int64(1024), stats.TotalSize)

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}

func TestReadAllStructured(t *testing.T) {
	// Setup test table
	tempDir, schema := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create test file
	filePath := filepath.Join(tempDir, "test.parquet")
	testData := []byte("test data")
	err := os.WriteFile(filePath, testData, 0644)
	require.NoError(t, err)

	// Update metadata
	schemaFields := make([]map[string]interface{}, len(schema.Fields()))
	for i, field := range schema.Fields() {
		schemaFields[i] = map[string]interface{}{
			"name": field.Name,
			"type": field.Type.String(),
		}
	}

	metadata := map[string]interface{}{
		"version":    int64(1),
		"timestamp":  time.Now().Unix(),
		"schema":     map[string]interface{}{
			"fields": schemaFields,
		},
		"files":      []string{filePath},
		"partitions": map[string][]string{},
		"stats":      &TableStats{},
		"metadata":   map[string]interface{}{},
	}

	data, err := json.Marshal(metadata)
	require.NoError(t, err)

	logFile := filepath.Join(tempDir, "_delta_log", "00000000000000000001.json")
	err = os.WriteFile(logFile, data, 0644)
	require.NoError(t, err)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Test reading all structured
	structured, err := reader.ReadAllStructured()
	require.NoError(t, err)
	require.NotNil(t, structured)
	require.Greater(t, len(structured), 0)

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}

func TestGetSchema(t *testing.T) {
	// Setup test table
	tempDir, schema := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Test schema
	actualSchema := reader.GetSchema()
	require.NotNil(t, actualSchema)
	expected := serializeSchema(schema)
	actual := serializeSchema(actualSchema)
	require.Equal(t, expected, actual)

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}

func TestConvertToStructuredData(t *testing.T) {
	// Setup test table
	tempDir, schema := setupTestTable(t)
	defer os.RemoveAll(tempDir)

	// Create test partition
	partitionPath := filepath.Join(tempDir, "test")
	err := os.MkdirAll(partitionPath, 0755)
	require.NoError(t, err)

	// Create test file
	filePath := filepath.Join(partitionPath, "test.parquet")
	testData := []byte("test data")
	err = os.WriteFile(filePath, testData, 0644)
	require.NoError(t, err)

	// Update metadata
	metadata := map[string]interface{}{
		"version":    int64(1),
		"timestamp":  time.Now().Unix(),
		"schema":     serializeSchema(schema),
		"files":      []string{filePath},
		"partitions": map[string][]string{"test": {filePath}},
		"stats":      &TableStats{},
		"metadata":   map[string]interface{}{},
	}

	data, err := json.Marshal(metadata)
	require.NoError(t, err)

	logFile := filepath.Join(tempDir, "_delta_log", "00000000000000000001.json")
	err = os.WriteFile(logFile, data, 0644)
	require.NoError(t, err)

	// Create reader
	reader, err := NewReader(tempDir)
	require.NoError(t, err)
	require.NotNil(t, reader)

	// Read partition
	record, err := reader.ReadPartition("test")
	require.NoError(t, err)
	require.NotNil(t, record)

	// Convert to structured data
	structured, err := reader.ConvertToStructuredData(record)
	require.NoError(t, err)
	require.NotNil(t, structured)
	require.Len(t, structured, 1)
	require.Equal(t, int32(1), structured[0]["id"])
	require.Equal(t, "test", structured[0]["name"])
	require.Equal(t, float64(1.0), structured[0]["value"])

	// Test closing
	err = reader.Close()
	require.NoError(t, err)
	require.True(t, reader.IsClosed())
}
