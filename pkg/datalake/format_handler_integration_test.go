package datalake

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatHandlerIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "format_handler_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test data files
	createTestDataFiles(t, tempDir)

	// Create a format handler
	handler := NewFormatHandler()

	// Test Delta format detection
	deltaPath := filepath.Join(tempDir, "delta_table")
	format, err := handler.DetectFormat(deltaPath)
	require.NoError(t, err)
	assert.Equal(t, DeltaFormat, format)

	// Test Parquet format detection
	parquetPath := filepath.Join(tempDir, "test.parquet")
	format, err = handler.DetectFormat(parquetPath)
	require.NoError(t, err)
	assert.Equal(t, ParquetFormat, format)

	// Test CSV format detection
	csvPath := filepath.Join(tempDir, "test.csv")
	format, err = handler.DetectFormat(csvPath)
	require.NoError(t, err)
	assert.Equal(t, CSVFormat, format)

	// Test custom configuration
	customConfig := &FormatConfig{
		DateFormats:          []string{"01/02/2006", "2006-01-02"},
		CSVDelimiter:         ';',
		CSVHasHeader:         false,
		MaxRowsForInference:  500,
		MinConfidenceThreshold: 0.9,
	}
	customHandler := NewFormatHandler().WithConfig(customConfig)
	assert.Equal(t, ';', customHandler.config.CSVDelimiter)
	assert.Equal(t, false, customHandler.config.CSVHasHeader)
	assert.Equal(t, 500, customHandler.config.MaxRowsForInference)
	assert.Equal(t, 0.9, customHandler.config.MinConfidenceThreshold)
}

func TestCSVSchemaInference(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csv_inference_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a CSV file with various data types
	csvContent := `id,name,age,active,created_at,score,tags
1,John Doe,30,true,2023-01-01,95.5,"tag1,tag2"
2,Jane Smith,25,false,2023-01-02,87.3,"tag3,tag4"
3,Bob Johnson,40,true,2023-01-03,92.1,"tag5,tag6"
4,Alice Brown,35,false,2023-01-04,88.9,"tag7,tag8"
5,Charlie Davis,28,true,2023-01-05,90.2,"tag9,tag10"
`
	csvPath := filepath.Join(tempDir, "test.csv")
	err = os.WriteFile(csvPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Create a format handler
	handler := NewFormatHandler()

	// Test CSV schema inference
	_, schema, err := handler.ReadWithInference(csvPath)
	require.NoError(t, err)
	require.NotNil(t, schema)

	// Verify inferred schema
	assert.Equal(t, 7, len(schema.Fields()))
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
	assert.Equal(t, "score", schema.Field(5).Name)
	assert.Equal(t, arrow.PrimitiveTypes.Float64, schema.Field(5).Type)
	assert.Equal(t, "tags", schema.Field(6).Name)
	assert.Equal(t, arrow.BinaryTypes.String, schema.Field(6).Type)
}

func TestCSVWithMissingValues(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "csv_missing_values_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a CSV file with missing values
	csvContent := `id,name,age,active,created_at
1,John Doe,,true,2023-01-01
2,,25,false,
,Bob Johnson,40,,2023-01-03
4,Alice Brown,35,false,2023-01-04
5,,,true,
`
	csvPath := filepath.Join(tempDir, "missing_values.csv")
	err = os.WriteFile(csvPath, []byte(csvContent), 0644)
	require.NoError(t, err)

	// Create a format handler
	handler := NewFormatHandler()

	// Test CSV schema inference with missing values
	record, schema, err := handler.ReadWithInference(csvPath)
	require.NoError(t, err)
	require.NotNil(t, schema)
	require.NotNil(t, record)

	// Verify inferred schema
	assert.Equal(t, 5, len(schema.Fields()))
	assert.Equal(t, "id", schema.Field(0).Name)
	assert.Equal(t, "name", schema.Field(1).Name)
	assert.Equal(t, "age", schema.Field(2).Name)
	assert.Equal(t, "active", schema.Field(3).Name)
	assert.Equal(t, "created_at", schema.Field(4).Name)

	// Verify record
	assert.Equal(t, int64(5), record.NumRows())
	assert.Equal(t, 5, record.NumCols())
}

func TestDataTypeInference(t *testing.T) {
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

	// Test mixed types (should default to string)
	dataType, confidence = handler.inferDataType([]string{"abc", "123", "true", "2023-01-01"})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test with missing values
	dataType, confidence = handler.inferDataType([]string{"1", "", "3", "", "5"})
	assert.Equal(t, arrow.PrimitiveTypes.Int64, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	// Test with all missing values
	dataType, confidence = handler.inferDataType([]string{"", "", ""})
	assert.Equal(t, arrow.BinaryTypes.String, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)
}

func TestCustomDateFormats(t *testing.T) {
	// Create a format handler with custom date formats
	config := &FormatConfig{
		DateFormats: []string{
			"01/02/2006",
			"02-Jan-2006",
			"January 2, 2006",
		},
		CSVDelimiter:         ',',
		CSVHasHeader:         true,
		MaxRowsForInference:  1000,
		MinConfidenceThreshold: 0.8,
	}
	handler := NewFormatHandler().WithConfig(config)

	// Test date inference with custom formats
	dataType, confidence := handler.inferDataType([]string{"01/15/2023", "02/20/2023", "03/25/2023"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_ms, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	dataType, confidence = handler.inferDataType([]string{"15-Jan-2023", "20-Feb-2023", "25-Mar-2023"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_ms, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)

	dataType, confidence = handler.inferDataType([]string{"January 15, 2023", "February 20, 2023", "March 25, 2023"})
	assert.Equal(t, arrow.FixedWidthTypes.Timestamp_ms, dataType)
	assert.InDelta(t, 1.0, confidence, 0.01)
}

// Helper function to create test data files
func createTestDataFiles(t *testing.T, dir string) {
	// Create a Delta table directory
	deltaDir := filepath.Join(dir, "delta_table")
	err := os.MkdirAll(filepath.Join(deltaDir, "_delta_log"), 0755)
	require.NoError(t, err)

	// Create a transaction log file
	logFile := filepath.Join(deltaDir, "_delta_log", "00000000000000000000.json")
	err = os.WriteFile(logFile, []byte("{}"), 0644)
	require.NoError(t, err)

	// Create a Parquet file
	parquetFile := filepath.Join(dir, "test.parquet")
	err = os.WriteFile(parquetFile, []byte("mock parquet data"), 0644)
	require.NoError(t, err)

	// Create a CSV file
	csvContent := "id,name,age\n1,John,30\n2,Jane,25\n3,Bob,40\n"
	csvFile := filepath.Join(dir, "test.csv")
	err = os.WriteFile(csvFile, []byte(csvContent), 0644)
	require.NoError(t, err)
}
