package profile

import (
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
)

// Anomaly represents an unusual pattern or outlier in the data
type Anomaly struct {
	Type        string
	Value       interface{}
	Description string
}

// Profile represents statistical information about a column
type Profile struct {
	Name        string
	Type        string
	RowCount    int64
	NullCount   int64
	NullPercent float64
	Distinct    int64
	Patterns    []string
	Anomalies   []Anomaly
	Stats       Stats               `json:"stats"`
	ValueCounts map[interface{}]int `json:"value_counts,omitempty"`
}

// Stats represents different statistical measures for numeric columns
type Stats struct {
	Min       float64
	Max       float64
	Mean      float64
	StdDev    float64
	Quartiles []float64
}

// Profiler generates profiles for Delta tables
type Profiler struct {
	tablePath string
	alloc     memory.Allocator
}

// NewProfiler creates a new Profiler
func NewProfiler(tablePath string) *Profiler {
	alloc := memory.NewGoAllocator()
	return &Profiler{
		tablePath: tablePath,
		alloc:     alloc,
	}
}

// GenerateProfile generates profiles for all columns in the table
func (p *Profiler) GenerateProfile() ([]*Profile, error) {
	// Create profiles for each column
	var profiles []*Profile

	// For now, return empty profiles since we don't have the pkg package
	return profiles, nil
}

// ProfileFromReader generates profiles from an io.Reader containing Parquet data
func (p *Profiler) ProfileFromReader(reader io.Reader) ([]*Profile, error) {
	// Parse the Parquet data
	data, err := ParseParquet(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Parquet data: %w", err)
	}

	// Create profiles from the parsed data
	var profiles []*Profile
	if len(data) > 0 {
		// Get column names from the first record
		columns := make([]string, 0, len(data[0]))
		for colName := range data[0] {
			columns = append(columns, colName)
		}

		// Create a profile for each column
		for _, colName := range columns {
			profile := &Profile{
				Name:     colName,
				RowCount: int64(len(data)),
			}

			// Determine column type and count nulls
			var nullCount int64
			var colType string
			values := make([]interface{}, 0, len(data))

			for _, row := range data {
				val := row[colName]
				if val == nil {
					nullCount++
					continue
				}

				// Determine column type from first non-null value
				if colType == "" {
					switch val.(type) {
					case int32, int64, int:
						colType = "integer"
					case float32, float64:
						colType = "float"
					case string:
						colType = "string"
					case bool:
						colType = "boolean"
					default:
						colType = "unknown"
					}
				}

				values = append(values, val)
			}

			profile.Type = colType
			profile.NullCount = nullCount
			profile.NullPercent = float64(nullCount) / float64(profile.RowCount) * 100

			// Calculate statistics for numeric columns
			if colType == "integer" || colType == "float" {
				// Convert values to float64 for statistics
				floatValues := make([]float64, 0, len(values))
				for _, v := range values {
					if v == nil {
						continue
					}

					var fv float64
					switch vt := v.(type) {
					case int:
						fv = float64(vt)
					case int32:
						fv = float64(vt)
					case int64:
						fv = float64(vt)
					case float32:
						fv = float64(vt)
					case float64:
						fv = vt
					default:
						continue
					}

					floatValues = append(floatValues, fv)
				}

				// Calculate basic statistics
				if len(floatValues) > 0 {
					profile.Stats = calculateStatistics(floatValues)
				}
			}

			// Find patterns and anomalies
			patterns, anomalies := p.findPatternsForValues(values, colType)
			profile.Patterns = patterns
			profile.Anomalies = anomalies

			// Count distinct values
			distinct := make(map[interface{}]bool)
			for _, v := range values {
				if v != nil {
					distinct[v] = true
				}
			}
			profile.Distinct = int64(len(distinct))

			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

// ProfileFromGzipReader generates profiles from a Gzip-compressed Parquet reader
func (p *Profiler) ProfileFromGzipReader(reader io.Reader) ([]*Profile, error) {
	// Create a gzip reader
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gz.Close()

	// Profile the decompressed data
	return p.ProfileFromReader(gz)
}

// ProfileTable generates profiles from an Arrow record
func (p *Profiler) ProfileTable(record arrow.Record) ([]*Profile, error) {
	// Create profiles for each column
	var profiles []*Profile

	// Get the number of columns
	numCols := int(record.NumCols())
	numRows := int(record.NumRows())

	// Create a profile for each column
	for i := 0; i < numCols; i++ {
		col := record.Column(i)
		field := record.Schema().Field(i)

		profile := &Profile{
			Name:     field.Name,
			RowCount: int64(numRows),
		}

		// Determine column type
		dataType := col.DataType().String()
		profile.Type = dataType

		// Count nulls
		var nullCount int64
		for j := 0; j < numRows; j++ {
			if !col.IsValid(j) {
				nullCount++
			}
		}

		profile.NullCount = nullCount
		profile.NullPercent = float64(nullCount) / float64(numRows) * 100

		// Calculate statistics for numeric columns
		if isNumeric(col.DataType().ID()) {
			// Extract numeric values
			values := make([]float64, 0, numRows)
			for j := 0; j < numRows; j++ {
				if !col.IsValid(j) {
					continue
				}

				// Convert to float64 based on type
				var val float64
				switch col.DataType().ID() {
				case arrow.INT8:
					val = float64(array.NewInt8Data(col.Data()).Value(j))
				case arrow.INT16:
					val = float64(array.NewInt16Data(col.Data()).Value(j))
				case arrow.INT32:
					val = float64(array.NewInt32Data(col.Data()).Value(j))
				case arrow.INT64:
					val = float64(array.NewInt64Data(col.Data()).Value(j))
				case arrow.UINT8:
					val = float64(array.NewUint8Data(col.Data()).Value(j))
				case arrow.UINT16:
					val = float64(array.NewUint16Data(col.Data()).Value(j))
				case arrow.UINT32:
					val = float64(array.NewUint32Data(col.Data()).Value(j))
				case arrow.UINT64:
					val = float64(array.NewUint64Data(col.Data()).Value(j))
				case arrow.FLOAT32:
					val = float64(array.NewFloat32Data(col.Data()).Value(j))
				case arrow.FLOAT64:
					val = array.NewFloat64Data(col.Data()).Value(j)
				default:
					continue
				}

				values = append(values, val)
			}

			// Calculate statistics
			if len(values) > 0 {
				profile.Stats = calculateStatistics(values)
			}
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// isNumeric checks if a data type is numeric
func isNumeric(dtID arrow.Type) bool {
	switch dtID {
	case arrow.INT8, arrow.INT16, arrow.INT32, arrow.INT64,
		arrow.UINT8, arrow.UINT16, arrow.UINT32, arrow.UINT64,
		arrow.FLOAT32, arrow.FLOAT64:
		return true
	default:
		return false
	}
}

// calculateStatistics calculates statistics for a slice of float64 values
func calculateStatistics(values []float64) Stats {
	stats := Stats{}
	if len(values) == 0 {
		return stats
	}

	// Calculate min, max, mean
	sum := 0.0
	min := values[0]
	max := values[0]

	for _, val := range values {
		sum += val
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}

	mean := sum / float64(len(values))

	// Calculate standard deviation
	sumOfSquares := 0.0
	for _, val := range values {
		sumOfSquares += (val - mean) * (val - mean)
	}
	stdDev := math.Sqrt(sumOfSquares / float64(len(values)))

	// Calculate quartiles
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	q1Index := int(float64(len(sorted)) * 0.25)
	q2Index := int(float64(len(sorted)) * 0.5)
	q3Index := int(float64(len(sorted)) * 0.75)

	stats.Min = min
	stats.Max = max
	stats.Mean = mean
	stats.StdDev = stdDev
	stats.Quartiles = []float64{sorted[q1Index], sorted[q2Index], sorted[q3Index]}

	return stats
}

// findPatternsForValues detects common patterns in a slice of values
func (p *Profiler) findPatternsForValues(values []interface{}, colType string) ([]string, []Anomaly) {
	patterns := []string{}
	anomalies := []Anomaly{}

	// Detect patterns based on column type
	switch colType {
	case "string":
		patterns, anomalies = p.findStringPatternsForValues(values)
	case "integer", "float":
		patterns, anomalies = p.findNumericPatternsForValues(values)
	}

	return patterns, anomalies
}

// findStringPatternsForValues detects patterns in string values
func (p *Profiler) findStringPatternsForValues(values []interface{}) ([]string, []Anomaly) {
	patterns := []string{}
	anomalies := []Anomaly{}

	// Check for common patterns
	emailCount := 0
	urlCount := 0
	digitCount := 0
	totalStrings := 0

	for _, v := range values {
		if v == nil {
			continue
		}

		s, ok := v.(string)
		if !ok {
			continue
		}

		totalStrings++

		// Check for email pattern
		if matchesRegex(s, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`) {
			emailCount++
		}

		// Check for URL pattern
		if matchesRegex(s, `^(http|https)://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`) {
			urlCount++
		}

		// Check for digit-only pattern
		if isDigit(s) {
			digitCount++
		}
	}

	// Add patterns if they occur frequently
	if totalStrings > 0 {
		threshold := int(float64(totalStrings) * 0.5) // 50% threshold

		if emailCount > threshold {
			patterns = append(patterns, "email")
		}

		if urlCount > threshold {
			patterns = append(patterns, "url")
		}

		if digitCount > threshold {
			patterns = append(patterns, "digit-only")
		}
	}

	return patterns, anomalies
}

// findNumericPatternsForValues detects patterns in numeric values
func (p *Profiler) findNumericPatternsForValues(values []interface{}) ([]string, []Anomaly) {
	patterns := []string{}
	anomalies := []Anomaly{}

	// Convert to float64 for statistics
	floatValues := make([]float64, 0, len(values))
	for _, v := range values {
		if v == nil {
			continue
		}

		var fv float64
		switch vt := v.(type) {
		case int:
			fv = float64(vt)
		case int32:
			fv = float64(vt)
		case int64:
			fv = float64(vt)
		case float32:
			fv = float64(vt)
		case float64:
			fv = vt
		default:
			continue
		}

		floatValues = append(floatValues, fv)
	}

	if len(floatValues) == 0 {
		return patterns, anomalies
	}

	// Calculate statistics
	stats := calculateStatistics(floatValues)

	// Check for outliers
	outlierThreshold := stats.StdDev * 3 // 3 standard deviations

	for _, val := range floatValues {
		// Check for outliers
		if math.Abs(val-stats.Mean) > outlierThreshold {
			anomalies = append(anomalies, Anomaly{
				Type:        "outlier",
				Value:       val,
				Description: fmt.Sprintf("Value %v is more than 3 standard deviations from the mean", val),
			})
		}
	}

	// Check for patterns
	if stats.Min >= 0 && stats.Max <= 1 {
		patterns = append(patterns, "probability")
	}

	if stats.Min >= 0 && stats.Max <= 100 && stats.StdDev < 30 {
		patterns = append(patterns, "percentage")
	}

	return patterns, anomalies
}

// matchesRegex checks if a string matches a regex pattern
func matchesRegex(s, pattern string) bool {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(s)
}

// isDigit checks if a string consists only of digits
func isDigit(s string) bool {
	return matchesRegex(s, "^[0-9]+$")
}
