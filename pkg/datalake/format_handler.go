package datalake

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
)

// FileFormat represents the supported file formats
type FileFormat string

const (
	// DeltaFormat represents Delta Lake format
	DeltaFormat FileFormat = "delta"
	
	// ParquetFormat represents Parquet format
	ParquetFormat FileFormat = "parquet"
	
	// CSVFormat represents CSV format
	CSVFormat FileFormat = "csv"
	
	// UnknownFormat represents an unknown format
	UnknownFormat FileFormat = "unknown"
)

// FormatHandler handles different file formats
type FormatHandler struct {
	// Default configuration for schema inference
	config *FormatConfig
}

// FormatConfig contains configuration options for format handling
type FormatConfig struct {
	// DateFormats is a list of date formats to try when inferring schema
	DateFormats []string
	
	// CSVDelimiter is the delimiter to use for CSV files
	CSVDelimiter rune
	
	// CSVHasHeader indicates if CSV files have a header row
	CSVHasHeader bool
	
	// MaxRowsForInference is the maximum number of rows to read for schema inference
	MaxRowsForInference int
	
	// MinConfidenceThreshold is the minimum confidence threshold for type inference
	MinConfidenceThreshold float64
}

// NewFormatHandler creates a new FormatHandler with default configuration
func NewFormatHandler() *FormatHandler {
	return &FormatHandler{
		config: &FormatConfig{
			DateFormats: []string{
				time.RFC3339,
				"2006-01-02",
				"2006/01/02",
				"01/02/2006",
				"02-Jan-2006",
			},
			CSVDelimiter:         ',',
			CSVHasHeader:         true,
			MaxRowsForInference:  1000,
			MinConfidenceThreshold: 0.8,
		},
	}
}

// WithConfig sets a custom configuration for the FormatHandler
func (h *FormatHandler) WithConfig(config *FormatConfig) *FormatHandler {
	h.config = config
	return h
}

// DetectFormat detects the format of a file or directory
func (h *FormatHandler) DetectFormat(path string) (FileFormat, error) {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return UnknownFormat, fmt.Errorf("error accessing path: %w", err)
	}
	
	// If it's a directory, check for Delta Lake format
	if info.IsDir() {
		// Check for _delta_log directory which is characteristic of Delta Lake
		deltaLogPath := filepath.Join(path, "_delta_log")
		if _, err := os.Stat(deltaLogPath); err == nil {
			return DeltaFormat, nil
		}
		return UnknownFormat, fmt.Errorf("unknown directory format: %s", path)
	}
	
	// If it's a file, check the extension
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
func (h *FormatHandler) ReadWithInference(path string) (arrow.Record, arrow.Schema, error) {
	format, err := h.DetectFormat(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error detecting format: %w", err)
	}
	
	switch format {
	case DeltaFormat:
		return h.readDelta(path)
	case ParquetFormat:
		return h.readParquet(path)
	case CSVFormat:
		return h.readCSV(path)
	default:
		return nil, nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// readDelta reads data from a Delta Lake table
func (h *FormatHandler) readDelta(path string) (arrow.Record, arrow.Schema, error) {
	// This would typically use a Delta Lake library like go-delta
	// For now, we'll return a placeholder error
	return nil, nil, errors.New("delta lake reading is implemented via the DeltaTable interface")
}

// readParquet reads data from a Parquet file
func (h *FormatHandler) readParquet(path string) (arrow.Record, arrow.Schema, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening parquet file: %w", err)
	}
	defer file.Close()
	
	// Create a parquet reader
	parquetReader, err := file.NewParquetReader(file)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating parquet reader: %w", err)
	}
	defer parquetReader.Close()
	
	// Create an Arrow reader
	arrowReader, err := pqarrow.NewFileReader(parquetReader, pqarrow.ArrowReadProperties{}, memory.DefaultAllocator)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating arrow reader: %w", err)
	}
	
	// Read the schema
	schema := arrowReader.Schema()
	
	// Read the record batch
	record, err := arrowReader.ReadRecordBatch(0)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading record batch: %w", err)
	}
	
	return record, schema, nil
}

// readCSV reads data from a CSV file with schema inference
func (h *FormatHandler) readCSV(path string) (arrow.Record, arrow.Schema, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening CSV file: %w", err)
	}
	defer file.Close()
	
	// Create a CSV reader
	reader := csv.NewReader(file)
	reader.Comma = h.config.CSVDelimiter
	
	// Read header if present
	var header []string
	if h.config.CSVHasHeader {
		header, err = reader.Read()
		if err != nil {
			return nil, nil, fmt.Errorf("error reading CSV header: %w", err)
		}
	}
	
	// Read rows for schema inference
	var rows [][]string
	for i := 0; i < h.config.MaxRowsForInference; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("error reading CSV row: %w", err)
		}
		rows = append(rows, row)
	}
	
	// If no header was provided, generate column names
	if header == nil {
		header = make([]string, len(rows[0]))
		for i := range header {
			header[i] = fmt.Sprintf("col%d", i+1)
		}
	}
	
	// Infer schema
	schema, err := h.inferCSVSchema(header, rows)
	if err != nil {
		return nil, nil, fmt.Errorf("error inferring CSV schema: %w", err)
	}
	
	// Convert data to Arrow record
	record, err := h.csvToArrowRecord(schema, header, rows)
	if err != nil {
		return nil, nil, fmt.Errorf("error converting CSV to Arrow: %w", err)
	}
	
	return record, schema, nil
}

// inferCSVSchema infers the schema from CSV data
func (h *FormatHandler) inferCSVSchema(header []string, rows [][]string) (arrow.Schema, error) {
	// Create fields for the schema
	fields := make([]arrow.Field, len(header))
	
	for i, name := range header {
		// Get all values for this column
		values := make([]string, len(rows))
		for j, row := range rows {
			if i < len(row) {
				values[j] = row[i]
			}
		}
		
		// Infer the type
		dataType, confidence := h.inferDataType(values)
		
		// Use string as fallback if confidence is too low
		if confidence < h.config.MinConfidenceThreshold {
			dataType = arrow.BinaryTypes.String
		}
		
		fields[i] = arrow.Field{
			Name: name,
			Type: dataType,
			Nullable: true,
		}
	}
	
	return arrow.NewSchema(fields, nil), nil
}

// inferDataType infers the data type from a set of values
func (h *FormatHandler) inferDataType(values []string) (arrow.DataType, float64) {
	// Count occurrences of different types
	var (
		intCount     int
		floatCount   int
		boolCount    int
		dateCount    int
		nullCount    int
		nonNullCount int
	)
	
	for _, v := range values {
		// Skip empty values
		if v == "" {
			nullCount++
			continue
		}
		
		nonNullCount++
		
		// Try integer
		if isInteger(v) {
			intCount++
			continue
		}
		
		// Try float
		if isFloat(v) {
			floatCount++
			continue
		}
		
		// Try boolean
		if isBoolean(v) {
			boolCount++
			continue
		}
		
		// Try date
		if isDate(v, h.config.DateFormats) {
			dateCount++
			continue
		}
	}
	
	// Determine the most likely type
	if nonNullCount == 0 {
		return arrow.BinaryTypes.String, 1.0
	}
	
	// Calculate confidence for each type
	intConfidence := float64(intCount) / float64(nonNullCount)
	floatConfidence := float64(floatCount+intCount) / float64(nonNullCount) // Integers can be floats
	boolConfidence := float64(boolCount) / float64(nonNullCount)
	dateConfidence := float64(dateCount) / float64(nonNullCount)
	
	// Choose the type with the highest confidence
	if intConfidence > 0.9 {
		return arrow.PrimitiveTypes.Int64, intConfidence
	}
	if floatConfidence > 0.9 {
		return arrow.PrimitiveTypes.Float64, floatConfidence
	}
	if boolConfidence > 0.9 {
		return arrow.FixedWidthTypes.Boolean, boolConfidence
	}
	if dateConfidence > 0.9 {
		return arrow.FixedWidthTypes.Timestamp_ms, dateConfidence
	}
	
	// Default to string
	return arrow.BinaryTypes.String, 1.0
}

// csvToArrowRecord converts CSV data to an Arrow record
func (h *FormatHandler) csvToArrowRecord(schema arrow.Schema, header []string, rows [][]string) (arrow.Record, error) {
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
				val, _ := parseIntegerValue(value)
				builders[i].(*array.Int64Builder).Append(val)
			case arrow.FLOAT64:
				val, _ := parseFloatValue(value)
				builders[i].(*array.Float64Builder).Append(val)
			case arrow.BOOL:
				val, _ := parseBooleanValue(value)
				builders[i].(*array.BooleanBuilder).Append(val)
			case arrow.TIMESTAMP:
				val, _ := parseDateValue(value, h.config.DateFormats)
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

// Helper functions for type checking and parsing

func isInteger(s string) bool {
	_, err := parseIntegerValue(s)
	return err == nil
}

func parseIntegerValue(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

func isFloat(s string) bool {
	_, err := parseFloatValue(s)
	return err == nil
}

func parseFloatValue(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

func isBoolean(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "false" || s == "yes" || s == "no" || s == "1" || s == "0"
}

func parseBooleanValue(s string) (bool, error) {
	s = strings.ToLower(s)
	switch s {
	case "true", "yes", "1":
		return true, nil
	case "false", "no", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}

func isDate(s string, formats []string) bool {
	_, err := parseDateValue(s, formats)
	return err == nil
}

func parseDateValue(s string, formats []string) (time.Time, error) {
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}
