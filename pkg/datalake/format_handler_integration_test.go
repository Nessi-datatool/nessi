package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestFormatDetection(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "format_detection_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a Delta Lake directory
	deltaDir := filepath.Join(tempDir, "delta_table")
	err = os.MkdirAll(filepath.Join(deltaDir, "_delta_log"), 0755)
	require.NoError(t, err)

	// Create a Parquet file
	parquetFile := filepath.Join(tempDir, "data.parquet")
	err = os.WriteFile(parquetFile, []byte("parquet data"), 0644)
	require.NoError(t, err)

	// Create a CSV file
	csvFile := filepath.Join(tempDir, "data.csv")
	err = os.WriteFile(csvFile, []byte("id,name,value\n1,test,1.23"), 0644)
	require.NoError(t, err)

	// Test format detection
	handler := NewFormatHandler()

	format, err := handler.DetectFormat(deltaDir)
	require.NoError(t, err)
	assert.Equal(t, DeltaFormat, format)

	format, err = handler.DetectFormat(parquetFile)
	require.NoError(t, err)
	assert.Equal(t, ParquetFormat, format)

	format, err = handler.DetectFormat(csvFile)
	require.NoError(t, err)
	assert.Equal(t, CSVFormat, format)

	// Test custom configuration
	customConfig := FormatHandlerConfig{
		DateFormats:           []string{"01/02/2006", "2006-01-02"},
		CSVHasHeader:          false,
		MaxRowsForInference:   500,
		MinConfidenceThreshold: 0.9,
	}
	customHandler := NewFormatHandlerWithConfig(customConfig)
	assert.Equal(t, false, customHandler.config.CSVHasHeader)
	assert.Equal(t, 500, customHandler.config.MaxRowsForInference)
	assert.Equal(t, 0.9, customHandler.config.MinConfidenceThreshold)
}

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestCSVSchemaInference(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csv_inference_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a CSV file with various data types
	csvContent := `id,name,active,score,date
1,Alice,true,95.5,2023-01-01
2,Bob,false,82.3,2023-01-02
3,Charlie,true,78.9,2023-01-03
4,David,false,91.2,2023-01-04
5,Eve,true,88.7,2023-01-05`

	csvFile := filepath.Join(tempDir, "data.csv")
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Test schema inference
	handler := NewFormatHandler()
	record, schema, err := handler.ReadWithInference(csvFile)
	require.NoError(t, err)
	require.NotNil(t, schema)
	require.NotNil(t, record)

	// Verify schema fields
	assert.Equal(t, 5, len(schema.Fields()))
	
	// Check field names
	assert.Equal(t, "id", schema.Field(0).Name)
	assert.Equal(t, "name", schema.Field(1).Name)
	assert.Equal(t, "active", schema.Field(2).Name)
	assert.Equal(t, "score", schema.Field(3).Name)
	assert.Equal(t, "date", schema.Field(4).Name)
	
	// Check field types
	assert.Equal(t, arrow.PrimitiveTypes.Int64, schema.Field(0).Type)
	assert.Equal(t, arrow.BinaryTypes.String, schema.Field(1).Type)
	assert.Equal(t, arrow.FixedWidthTypes.Boolean, schema.Field(2).Type)
	assert.Equal(t, arrow.PrimitiveTypes.Float64, schema.Field(3).Type)
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_s, schema.Field(4).Type)
	
	// Check record
	assert.Equal(t, int64(5), record.NumRows())
}

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestCSVWithoutHeader(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csv_no_header_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a CSV file without header
	csvContent := `1,Alice,true,95.5,2023-01-01
2,Bob,false,82.3,2023-01-02
3,Charlie,true,78.9,2023-01-03
4,David,false,91.2,2023-01-04
5,Eve,true,88.7,2023-01-05`

	csvFile := filepath.Join(tempDir, "data.csv")
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Test with custom configuration for no header
	config := FormatHandlerConfig{
		CSVHasHeader: false,
	}
	handler := NewFormatHandlerWithConfig(config)
	
	record, schema, err := handler.ReadWithInference(csvFile)
	require.NoError(t, err)
	require.NotNil(t, schema)
	require.NotNil(t, record)

	// Verify schema fields
	assert.Equal(t, 5, len(schema.Fields()))
	
	// Check field names (should be col1, col2, etc.)
	assert.Equal(t, "col1", schema.Field(0).Name)
	assert.Equal(t, "col2", schema.Field(1).Name)
	assert.Equal(t, "col3", schema.Field(2).Name)
	assert.Equal(t, "col4", schema.Field(3).Name)
	assert.Equal(t, "col5", schema.Field(4).Name)
}

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestDataTypeInference(t *testing.T) {
	handler := NewFormatHandler()

	// Test integer inference
	dataType := handler.detectDataTypeWithConfidence([]string{"1", "2", "3", "4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Int64, dataType)

	// Test float inference
	dataType = handler.detectDataTypeWithConfidence([]string{"1.1", "2.2", "3.3", "4.4", "5.5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)

	// Test mixed numeric inference
	dataType = handler.detectDataTypeWithConfidence([]string{"1", "2.2", "3", "4.4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)

	// Test boolean inference
	dataType = handler.detectDataTypeWithConfidence([]string{"true", "false", "true", "false", "true"})
	assert.Equal(t, arrow.FixedWidthTypes.Boolean, dataType)

	// Test date inference
	dataType = handler.detectDataTypeWithConfidence([]string{"2023-01-01", "2023-01-02", "2023-01-03"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_s, dataType)

	// Test mixed types (should default to string)
	dataType = handler.detectDataTypeWithConfidence([]string{"abc", "123", "true", "2023-01-01"})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)
}
