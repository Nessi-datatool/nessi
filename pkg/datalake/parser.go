package datalake

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
)

// ParsedData represents parsed data from a Reader
type ParsedData struct {
	Data   []map[string]interface{} `json:"data"`
	Schema *DeltaSchema             `json:"schema"`
	Count  int                      `json:"count"`
}

// ParseReadCloser parses data from an io.ReadCloser into a structured format
// This is used to convert the raw data from ReadAll() into a format suitable for profiling
func ParseReadCloser(r io.ReadCloser, format string) (*ParsedData, error) {
	defer r.Close()

	switch strings.ToLower(format) {
	case "parquet":
		return parseParquet(r)
	case "json":
		return parseJSON(r)
	case "csv":
		return parseCSV(r)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// parseParquet parses Parquet format data from an io.ReadCloser
func parseParquet(r io.ReadCloser) (*ParsedData, error) {
	// For now, we'll return a stub implementation
	// In a real implementation, we would need to convert the io.ReadCloser to a ReaderAtSeeker
	// or use a different approach to read the Parquet file

	// Initialize result with empty data
	result := &ParsedData{
		Data: make([]map[string]interface{}, 0),
		Schema: &DeltaSchema{
			Fields: []DeltaField{
				{Name: "placeholder", Type: "string", Nullable: true},
			},
		},
		Count: 0,
	}

	return result, nil
}

// extractArrowValue extracts a value from an Arrow array at the specified index
func extractArrowValue(col arrow.Array, rowIdx int64) interface{} {
	if col.IsNull(int(rowIdx)) {
		return nil
	}

	switch arr := col.(type) {
	case *array.Int8:
		return arr.Value(int(rowIdx))
	case *array.Int16:
		return arr.Value(int(rowIdx))
	case *array.Int32:
		return arr.Value(int(rowIdx))
	case *array.Int64:
		return arr.Value(int(rowIdx))
	case *array.Uint8:
		return arr.Value(int(rowIdx))
	case *array.Uint16:
		return arr.Value(int(rowIdx))
	case *array.Uint32:
		return arr.Value(int(rowIdx))
	case *array.Uint64:
		return arr.Value(int(rowIdx))
	case *array.Float32:
		return arr.Value(int(rowIdx))
	case *array.Float64:
		return arr.Value(int(rowIdx))
	case *array.Boolean:
		return arr.Value(int(rowIdx))
	case *array.String:
		return arr.Value(int(rowIdx))
	case *array.Binary:
		return arr.Value(int(rowIdx))
	case *array.Date32:
		date := arr.Value(int(rowIdx))
		return time.Unix(int64(date)*86400, 0).Format("2006-01-02")
	case *array.Date64:
		date := arr.Value(int(rowIdx))
		return time.Unix(0, int64(date)*1000000).Format("2006-01-02")
	case *array.Timestamp:
		ts := arr.Value(int(rowIdx))
		return time.Unix(0, int64(ts)).Format(time.RFC3339)
	default:
		// For unsupported types, convert to string
		return fmt.Sprintf("%v", arr.ValueStr(int(rowIdx)))
	}
}

// parseJSON parses JSON format data from an io.ReadCloser
func parseJSON(r io.ReadCloser) (*ParsedData, error) {
	var rawData []map[string]interface{}

	// Decode JSON data
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Initialize result
	result := &ParsedData{
		Data: rawData,
		Schema: &DeltaSchema{
			Fields: []DeltaField{},
		},
		Count: len(rawData),
	}

	// Extract schema information from the first record
	if len(rawData) > 0 {
		for key, value := range rawData[0] {
			dataType := inferJSONType(value)
			result.Schema.Fields = append(result.Schema.Fields, DeltaField{
				Name:     key,
				Type:     dataType,
				Nullable: true,
			})
		}
	}

	return result, nil
}

// inferJSONType infers the data type of a JSON value
func inferJSONType(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case bool:
		return "boolean"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "integer"
	case float32, float64:
		return "float"
	case string:
		// Try to detect date/time formats
		if isDateString(v) {
			return "date"
		}
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}

// isDateString checks if a string appears to be a date
func isDateString(s string) bool {
	dateFormats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"01-02-2006",
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range dateFormats {
		if _, err := time.Parse(format, s); err == nil {
			return true
		}
	}

	return false
}

// parseCSV parses CSV format data from an io.ReadCloser
func parseCSV(r io.ReadCloser) (*ParsedData, error) {
	// For simplicity, we'll implement a basic CSV parser
	// In a production environment, you might want to use a more robust CSV parser

	// Read all data
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV data: %w", err)
	}

	// Split into lines
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("CSV must have at least a header row and one data row")
	}

	// Parse header
	header := strings.Split(lines[0], ",")
	for i, h := range header {
		header[i] = strings.TrimSpace(h)
	}

	// Initialize result
	result := &ParsedData{
		Data: make([]map[string]interface{}, 0, len(lines)-1),
		Schema: &DeltaSchema{
			Fields: []DeltaField{},
		},
		Count: len(lines) - 1,
	}

	// Create schema fields based on header
	schemaFields := make(map[string]string)

	// Process data rows
	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue // Skip empty lines
		}

		// Split the line into fields
		fields := strings.Split(line, ",")

		// Create a record
		record := make(map[string]interface{})

		// Map fields to header
		for j, field := range fields {
			if j >= len(header) {
				break // Skip extra fields
			}

			field = strings.TrimSpace(field)

			// Try to convert to appropriate type
			value := inferCSVValue(field)
			record[header[j]] = value

			// For the first row, infer schema
			if i == 1 {
				schemaFields[header[j]] = inferJSONType(value)
			}
		}

		result.Data = append(result.Data, record)
	}

	// Convert schema map to DeltaSchema
	for field, dataType := range schemaFields {
		result.Schema.Fields = append(result.Schema.Fields, DeltaField{
			Name:     field,
			Type:     dataType,
			Nullable: true,
		})
	}

	return result, nil
}

// inferCSVValue tries to convert a CSV string to an appropriate type
func inferCSVValue(s string) interface{} {
	// Check for empty or null values
	if s == "" || strings.ToLower(s) == "null" || strings.ToLower(s) == "na" || strings.ToLower(s) == "n/a" {
		return nil
	}

	// Try to parse as integer
	if i, err := parseInt(s); err == nil {
		return i
	}

	// Try to parse as float
	if f, err := parseFloat(s); err == nil {
		return f
	}

	// Try to parse as boolean
	if b, err := parseBoolean(s); err == nil {
		return b
	}

	// Default to string
	return s
}

// parseInt tries to parse a string as an integer
func parseInt(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscanf(s, "%d", &i)
	return i, err
}

// parseFloat tries to parse a string as a float
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// parseBoolean tries to parse a string as a boolean
func parseBoolean(s string) (bool, error) {
	s = strings.ToLower(s)
	if s == "true" || s == "yes" || s == "y" || s == "1" {
		return true, nil
	}
	if s == "false" || s == "no" || s == "n" || s == "0" {
		return false, nil
	}
	return false, fmt.Errorf("not a boolean: %s", s)
}
