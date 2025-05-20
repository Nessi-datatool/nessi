package datalake

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/arrow/go/v15/arrow"
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
	err = os.WriteFile(parquetFile, []byte("test"), 0644)
	require.NoError(t, err)

	format, err = handler.DetectFormat(parquetFile)
	require.NoError(t, err)
	assert.Equal(t, ParquetFormat, format)

	// Test CSV format detection
	csvFile := filepath.Join(tempDir, "test.csv")
	err = os.WriteFile(csvFile, []byte("test"), 0644)
	require.NoError(t, err)

	format, err = handler.DetectFormat(csvFile)
	require.NoError(t, err)
	assert.Equal(t, CSVFormat, format)
}

func TestInferCSVSchema(t *testing.T) {
	handler := NewFormatHandler()

	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "format_handler_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a CSV file
	csvContent := `id,name,active,score,date
1,Alice,true,95.5,2023-01-01
2,Bob,false,82.3,2023-01-02
3,Charlie,true,78.9,2023-01-03`

	csvFile := filepath.Join(tempDir, "test.csv")
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Test schema inference
	schema, err := handler.inferCSVSchema([]string{"id", "name", "active", "score", "date"}, [][]string{
		{"1", "Alice", "true", "95.5", "2023-01-01"},
		{"2", "Bob", "false", "82.3", "2023-01-02"},
		{"3", "Charlie", "true", "78.9", "2023-01-03"},
	})
	require.NoError(t, err)
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
}

func TestDetectDataType(t *testing.T) {
	handler := NewFormatHandler()

	// Test integer detection
	dataType := handler.detectDataType([]string{"1", "2", "3", "4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Int64, dataType)

	// Test float detection
	dataType = handler.detectDataType([]string{"1.1", "2.2", "3.3", "4.4", "5.5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)

	// Test boolean detection
	dataType = handler.detectDataType([]string{"true", "false", "true", "false", "true"})
	assert.Equal(t, arrow.FixedWidthTypes.Boolean, dataType)

	// Test date detection
	dataType = handler.detectDataType([]string{"2023-01-01", "2023-01-02", "2023-01-03"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_s, dataType)

	// Test string detection
	dataType = handler.detectDataType([]string{"abc", "def", "ghi", "jkl", "mno"})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)

	// Test mixed types (should default to string)
	dataType = handler.detectDataType([]string{"abc", "123", "true", "2023-01-01"})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)
}

func TestDetectDataTypeWithConfidence(t *testing.T) {
	handler := NewFormatHandler()

	// Test integer detection with confidence
	dataType := handler.detectDataTypeWithConfidence([]string{"1", "2", "3", "4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Int64, dataType)

	// Test float detection with confidence
	dataType = handler.detectDataTypeWithConfidence([]string{"1.1", "2.2", "3.3", "4.4", "5.5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)

	// Test mixed numeric detection with confidence
	dataType = handler.detectDataTypeWithConfidence([]string{"1", "2.2", "3", "4.4", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Float64, dataType)

	// Test boolean detection with confidence
	dataType = handler.detectDataTypeWithConfidence([]string{"true", "false", "true", "false", "true"})
	assert.Equal(t, arrow.FixedWidthTypes.Boolean, dataType)

	// Test date detection with confidence
	dataType = handler.detectDataTypeWithConfidence([]string{"2023-01-01", "2023-01-02", "2023-01-03"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_s, dataType)

	// Test mixed types with confidence (should default to string)
	dataType = handler.detectDataTypeWithConfidence([]string{"abc", "123", "true", "2023-01-01"})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)
}

func TestFormatHandlerConfig(t *testing.T) {
	// Test default configuration
	defaultConfig := DefaultFormatHandlerConfig()
	assert.True(t, defaultConfig.CSVHasHeader)
	assert.Equal(t, 100, defaultConfig.MaxRowsForInference)
	assert.Equal(t, 0.7, defaultConfig.MinConfidenceThreshold)

	// Test custom configuration
	customConfig := FormatHandlerConfig{
		CSVHasHeader:           false,
		MaxRowsForInference:    500,
		MinConfidenceThreshold: 0.9,
	}
	handler := NewFormatHandlerWithConfig(customConfig)
	assert.Equal(t, false, handler.config.CSVHasHeader)
	assert.Equal(t, 500, handler.config.MaxRowsForInference)
	assert.Equal(t, 0.9, handler.config.MinConfidenceThreshold)
}
