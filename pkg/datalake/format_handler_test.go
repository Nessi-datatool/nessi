package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectFormat(t *testing.T) {
	handler := NewFormatHandler()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "format_handler_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test Delta Lake format detection
	deltaDir := filepath.Join(tempDir, "delta_table")
	err = os.MkdirAll(filepath.Join(deltaDir, "_delta_log"), 0755)
	require.NoError(t, err)

	format, err := handler.DetectFormat(deltaDir)
	require.NoError(t, err)
	assert.Equal(t, DeltaFormat, format)

	// Test Parquet format detection
	parquetFile := filepath.Join(tempDir, "test.parquet")
	_, err = os.Create(parquetFile)
	require.NoError(t, err)

	format, err = handler.DetectFormat(parquetFile)
	require.NoError(t, err)
	assert.Equal(t, ParquetFormat, format)

	// Test CSV format detection
	csvFile := filepath.Join(tempDir, "test.csv")
	_, err = os.Create(csvFile)
	require.NoError(t, err)

	format, err = handler.DetectFormat(csvFile)
	require.NoError(t, err)
	assert.Equal(t, CSVFormat, format)

	// Test unknown format detection
	unknownFile := filepath.Join(tempDir, "test.txt")
	_, err = os.Create(unknownFile)
	require.NoError(t, err)

	format, err = handler.DetectFormat(unknownFile)
	assert.Error(t, err)
	assert.Equal(t, UnknownFormat, format)
}

func TestInferDataType(t *testing.T) {
	handler := NewFormatHandler()

	// Test integer inference
	dataType, confidence := handler.inferDataType([]string{"1", "2", "3", "4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Int64, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test float inference
	dataType, confidence = handler.inferDataType([]string{"1.1", "2.2", "3.3", "4.4", "5.5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test mixed numeric inference
	dataType, confidence = handler.inferDataType([]string{"1", "2.2", "3", "4.4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test boolean inference
	dataType, confidence = handler.inferDataType([]string{"true", "false", "true", "false", "true"})
	assert.Equal(t, arrow.FixedWidthTypes.Boolean, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test date inference
	dataType, confidence = handler.inferDataType([]string{"2023-01-01", "2023-01-02", "2023-01-03"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_ms, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test string fallback
	dataType, confidence = handler.inferDataType([]string{"abc", "def", "123", "true", "2023-01-01"})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)
}

func TestInferCSVSchema(t *testing.T) {
	handler := NewFormatHandler()

	// Test schema inference
	header := []string{"id", "name", "age", "active", "created_at"}
	rows := [][]string{
		{"1", "John", "30", "true", "2023-01-01"},
		{"2", "Jane", "25", "false", "2023-01-02"},
		{"3", "Bob", "40", "true", "2023-01-03"},
	}

	schema, err := handler.inferCSVSchema(header, rows)
	require.NoError(t, err)

	// Check schema fields
	assert.Equal(t, 5, len(schema.Fields()))
	assert.Equal(t, "id", schema.Field(0).Name)
	assert.Equal(t, arrow.PrimitiveTypes.Int64, schema.Field(0).Type)
	assert.Equal(t, "name", schema.Field(1).Name)
	assert.Equal(t, arrow.BinaryTypes.String, schema.Field(1).Type)
	assert.Equal(t, "age", schema.Field(2).Name)
	assert.Equal(t, arrow.PrimitiveTypes.Int64, schema.Field(2).Type)
	assert.Equal(t, "active", schema.Field(3).Name)
	assert.Equal(t, arrow.FixedWidthTypes.Boolean, schema.Field(3).Type)
	assert.Equal(t, "created_at", schema.Field(4).Name)
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_ms, schema.Field(4).Type)
}

func TestCSVToArrowRecord(t *testing.T) {
	handler := NewFormatHandler()

	// Create test data
	header := []string{"id", "name", "age"}
	rows := [][]string{
		{"1", "John", "30"},
		{"2", "Jane", "25"},
		{"3", "Bob", "40"},
	}

	// Create schema
	fields := []arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
		{Name: "name", Type: arrow.BinaryTypes.String, Nullable: true},
		{Name: "age", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
	}
	schema := arrow.NewSchema(fields, nil)

	// Convert to Arrow record
	record, err := handler.csvToArrowRecord(schema, header, rows)
	require.NoError(t, err)

	// Check record
	assert.Equal(t, int64(3), record.NumRows())
	assert.Equal(t, 3, record.NumCols())
}

func TestCSVWithMissingValues(t *testing.T) {
	handler := NewFormatHandler()

	// Create test data with missing values
	header := []string{"id", "name", "age"}
	rows := [][]string{
		{"1", "John", ""},
		{"2", "", "25"},
		{"", "Bob", "40"},
	}

	// Infer schema
	schema, err := handler.inferCSVSchema(header, rows)
	require.NoError(t, err)

	// Convert to Arrow record
	record, err := handler.csvToArrowRecord(schema, header, rows)
	require.NoError(t, err)

	// Check record
	assert.Equal(t, int64(3), record.NumRows())
	assert.Equal(t, 3, record.NumCols())
}

func TestFormatConfig(t *testing.T) {
	// Create a custom config
	config := &FormatConfig{
		DateFormats: []string{"01/02/2006", "2006-01-02"},
		CSVDelimiter: ';',
		CSVHasHeader: false,
		MaxRowsForInference: 500,
		MinConfidenceThreshold: 0.9,
	}

	// Create handler with custom config
	handler := NewFormatHandler().WithConfig(config)

	// Verify config was set correctly
	assert.Equal(t, ';', handler.config.CSVDelimiter)
	assert.Equal(t, false, handler.config.CSVHasHeader)
	assert.Equal(t, 500, handler.config.MaxRowsForInference)
	assert.Equal(t, 0.9, handler.config.MinConfidenceThreshold)
	assert.Equal(t, 2, len(handler.config.DateFormats))
}
