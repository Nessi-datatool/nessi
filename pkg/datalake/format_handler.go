package datalake

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// FormatType represents the type of a data format
type FormatType string

const (
	// DeltaFormat represents the Delta Lake format
	DeltaFormat FormatType = "delta"

	// ParquetFormat represents the Parquet format
	ParquetFormat FormatType = "parquet"

	// CSVFormat represents the CSV format
	CSVFormat FormatType = "csv"

	// UnknownFormat represents an unknown format
	UnknownFormat FormatType = "unknown"
)

// FormatHandlerConfig represents configuration for the format handler
type FormatHandlerConfig struct {
	// CSVHasHeader indicates whether CSV files have a header row
	CSVHasHeader bool

	// DateFormats is a list of date formats to try when parsing dates
	DateFormats []string

	// MaxRowsForInference is the maximum number of rows to use for schema inference
	MaxRowsForInference int

	// MinConfidenceThreshold is the minimum confidence threshold for type inference
	MinConfidenceThreshold float64
}

// DefaultFormatHandlerConfig returns the default configuration for the format handler
func DefaultFormatHandlerConfig() FormatHandlerConfig {
	return FormatHandlerConfig{
		CSVHasHeader:           true,
		DateFormats:            []string{time.RFC3339, "2006-01-02", "2006/01/02"},
		MaxRowsForInference:    100,
		MinConfidenceThreshold: 0.7,
	}
}

// FormatHandler handles different data formats
type FormatHandler struct {
	config FormatHandlerConfig
}

// NewFormatHandler creates a new format handler with default configuration
func NewFormatHandler() *FormatHandler {
	return &FormatHandler{
		config: DefaultFormatHandlerConfig(),
	}
}

// NewFormatHandlerWithConfig creates a new format handler with the specified configuration
func NewFormatHandlerWithConfig(config FormatHandlerConfig) *FormatHandler {
	return &FormatHandler{
		config: config,
	}
}

// DetectFormat detects the format of a file
func (h *FormatHandler) DetectFormat(path string) (FormatType, error) {
	// Check if it's a Delta Lake table
	if hasDeltaLogDirectory(path) {
		return DeltaFormat, nil
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".parquet":
		return ParquetFormat, nil
	case ".csv":
		return CSVFormat, nil
	default:
		return UnknownFormat, fmt.Errorf("unknown file format: %s", ext)
	}
}

// ReadWithInference reads a file with automatic schema inference
func (h *FormatHandler) ReadWithInference(path string) (arrow.Record, *arrow.Schema, error) {
	format, err := h.DetectFormat(path)
	if err != nil {
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("error detecting format: %w", err)
	}

	switch format {
	case DeltaFormat:
		return h.readDelta(path)
	case ParquetFormat:
		return h.readParquet(path)
	case CSVFormat:
		return h.readCSV(path)
	default:
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("unsupported format: %s", format)
	}
}

// readDelta reads data from a Delta Lake table
func (h *FormatHandler) readDelta(path string) (arrow.Record, *arrow.Schema, error) {
	// For now, return a placeholder implementation
	// In a real implementation, we would use the Delta Lake library
	emptySchema := arrow.NewSchema([]arrow.Field{}, nil)
	return nil, emptySchema, fmt.Errorf("delta lake reading not fully implemented")
}

// readParquet reads data from a Parquet file
func (h *FormatHandler) readParquet(path string) (arrow.Record, *arrow.Schema, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("error opening parquet file: %w", err)
	}
	defer file.Close()

	// For now, return a placeholder implementation
	// In a real implementation, we would use the parquet library to read the file
	emptySchema := arrow.NewSchema([]arrow.Field{}, nil)
	return nil, emptySchema, fmt.Errorf("parquet reading not fully implemented")
}

// readCSV reads data from a CSV file
func (h *FormatHandler) readCSV(path string) (arrow.Record, *arrow.Schema, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read all rows first
	allRows, err := reader.ReadAll()
	if err != nil {
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("error reading CSV data: %w", err)
	}

	// Check if we have any data
	if len(allRows) == 0 {
		// Return an empty schema if there are no rows
		emptySchema := arrow.NewSchema([]arrow.Field{}, nil)
		return nil, emptySchema, nil
	}

	// Handle header and data rows
	var header []string
	var dataRows [][]string

	if h.config.CSVHasHeader {
		// First row is header, rest is data
		header = allRows[0]
		dataRows = allRows[1:]
	} else {
		// All rows are data, generate column names
		dataRows = allRows
		header = make([]string, len(allRows[0]))
		for i := range header {
			header[i] = fmt.Sprintf("col%d", i+1)
		}
	}

	// Limit rows for inference if needed
	inferenceRows := dataRows
	if len(dataRows) > h.config.MaxRowsForInference {
		inferenceRows = dataRows[:h.config.MaxRowsForInference]
	}

	// Infer schema
	schema, err := h.inferCSVSchema(header, inferenceRows)
	if err != nil {
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("error inferring CSV schema: %w", err)
	}

	// Convert data to Arrow record
	record, err := h.csvToArrowRecord(schema, header, dataRows)
	if err != nil {
		return nil, arrow.NewSchema([]arrow.Field{}, nil), fmt.Errorf("error converting CSV to Arrow: %w", err)
	}

	return record, schema, nil
}

// csvToArrowRecord converts CSV data to an Arrow record
func (h *FormatHandler) csvToArrowRecord(schema *arrow.Schema, header []string, rows [][]string) (arrow.Record, error) {
	// Create builders for each column
	builders := make([]array.Builder, len(header))
	for i, field := range schema.Fields() {
		builders[i] = array.NewBuilder(memory.DefaultAllocator, field.Type)
	}

	// Add data to builders
	for _, row := range rows {
		for i, value := range row {
			if i >= len(builders) {
				continue
			}

			// Skip empty values
			if value == "" {
				builders[i].AppendNull()
				continue
			}

			// Add value based on the field type
			switch builders[i].Type().ID() {
			case arrow.INT64:
				val, _ := strconv.ParseInt(value, 10, 64)
				builders[i].(*array.Int64Builder).Append(val)
			case arrow.FLOAT64:
				val, _ := strconv.ParseFloat(value, 64)
				builders[i].(*array.Float64Builder).Append(val)
			case arrow.BOOL:
				val, _ := strconv.ParseBool(value)
				builders[i].(*array.BooleanBuilder).Append(val)
			case arrow.TIMESTAMP:
				val, _ := time.Parse(time.RFC3339, value)
				builders[i].(*array.TimestampBuilder).Append(arrow.Timestamp(val.UnixNano()))
			default:
				builders[i].(*array.StringBuilder).Append(value)
			}
		}
	}

	// Create arrays from builders
	arrays := make([]arrow.Array, len(builders))
	for i, builder := range builders {
		arrays[i] = builder.NewArray()
		defer arrays[i].Release()
	}

	// Create record batch
	return array.NewRecord(schema, arrays, int64(len(rows))), nil
}

// inferCSVSchema infers the schema from CSV data
func (h *FormatHandler) inferCSVSchema(header []string, rows [][]string) (*arrow.Schema, error) {
	// Create fields for the schema
	fields := make([]arrow.Field, len(header))

	// Infer type for each column
	for i, name := range header {
		// Extract values for this column
		values := make([]string, 0, len(rows))
		for _, row := range rows {
			if i < len(row) {
				values = append(values, row[i])
			}
		}

		// Infer the type
		dataType := h.detectDataType(values)

		// Use string as fallback if confidence is too low
		if dataType == arrow.BinaryTypes.String {
			dataType = h.detectDataTypeWithConfidence(values)
		}

		fields[i] = arrow.Field{
			Name:     name,
			Type:     dataType,
			Nullable: true,
		}
	}

	// Create schema from fields
	return arrow.NewSchema(fields, nil), nil
}

// detectDataType detects the data type of a value
func (h *FormatHandler) detectDataType(values []string) arrow.DataType {
	if len(values) == 0 {
		return arrow.BinaryTypes.String
	}

	// Try to parse as integer
	if _, err := strconv.ParseInt(values[0], 10, 64); err == nil {
		return arrow.PrimitiveTypes.Int64
	}

	// Try to parse as float
	if _, err := strconv.ParseFloat(values[0], 64); err == nil {
		return arrow.PrimitiveTypes.Float64
	}

	// Try to parse as boolean
	value := strings.ToLower(values[0])
	if value == "true" || value == "false" || value == "yes" || value == "no" || value == "1" || value == "0" {
		return arrow.FixedWidthTypes.Boolean
	}

	// Try to parse as date
	// This is a simplified check, in a real implementation we would check against multiple date formats
	if strings.Contains(values[0], "-") && strings.Count(values[0], "-") == 2 {
		return arrow.FixedWidthTypes.Timestamp_s
	}

	// Default to string
	return arrow.BinaryTypes.String
}

// detectDataTypeWithConfidence detects the data type of a value with confidence
func (h *FormatHandler) detectDataTypeWithConfidence(values []string) arrow.DataType {
	// Count occurrences of different types
	var (
		intCount     int
		floatCount   int
		boolCount    int
		dateCount    int
		stringCount  int
		nonNullCount int
	)

	// Check each value
	for _, v := range values {
		if v == "" {
			continue
		}

		nonNullCount++

		// Try integer
		if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			intCount++
			continue
		}

		// Try float
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			floatCount++
			continue
		}

		// Try boolean
		value := strings.ToLower(v)
		if value == "true" || value == "false" || value == "yes" || value == "no" || value == "1" || value == "0" {
			boolCount++
			continue
		}

		// Try date
		// This is a simplified check, in a real implementation we would check against multiple date formats
		if strings.Contains(v, "-") && strings.Count(v, "-") == 2 {
			dateCount++
			continue
		}

		// Default to string
		stringCount++
	}

	// Determine the most likely type
	if nonNullCount == 0 {
		return arrow.BinaryTypes.String
	}

	// Calculate confidence for each type
	boolConfidence := float64(boolCount) / float64(nonNullCount)
	dateConfidence := float64(dateCount) / float64(nonNullCount)

	// Calculate combined numeric confidence (integers can be represented as floats)
	numericConfidence := float64(intCount+floatCount) / float64(nonNullCount)

	// Choose the type with the highest confidence
	if numericConfidence > 0.9 {
		// If we have any floats, use float type to represent all numbers
		if floatCount > 0 {
			return arrow.PrimitiveTypes.Float64
		}
		return arrow.PrimitiveTypes.Int64
	}
	if boolConfidence > 0.9 {
		return arrow.FixedWidthTypes.Boolean
	}
	if dateConfidence > 0.9 {
		return arrow.FixedWidthTypes.Timestamp_s
	}

	// Default to string
	return arrow.BinaryTypes.String
}

// Helper functions for format detection

// hasDeltaLogDirectory checks if a path has a _delta_log directory
func hasDeltaLogDirectory(path string) bool {
	// Check if the path is a directory
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}

	// Check if the _delta_log directory exists
	deltaLogPath := filepath.Join(path, "_delta_log")
	info, err = os.Stat(deltaLogPath)
	return err == nil && info.IsDir()
}

// hasParquetExtension checks if a path has a .parquet extension
func hasParquetExtension(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".parquet"
}

// hasCSVExtension checks if a path has a .csv extension
func hasCSVExtension(path string) bool {
	return strings.ToLower(filepath.Ext(path)) == ".csv"
}
