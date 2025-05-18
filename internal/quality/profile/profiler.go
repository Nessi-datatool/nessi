// Implement a data profiler for Delta tables with these requirements:
// 1. Define a Profile struct containing:
//    - TableName string
//    - RowCount int
//    - ColumnProfiles map[string]ColumnProfile
//    - CreatedAt time.Time
//    - ExecutionTimeMs int64
// 2. Define a ColumnProfile struct with:
//    - Name, DataType string
//    - NullCount, NullPercentage float64
//    - UniqueCount, UniquePercentage float64
//    - Min, Max interface{} (for appropriate types)
//    - Mean, Median, StdDev float64 (for numeric types)
//    - TopValues []ValueCount - Most common values
//    - Histogram []HistogramBin (for numeric types)
//    - Patterns []PatternCount (for string types)
// 3. Implement ProfileTable function:
//    - ProfileTable(table *delta.DeltaTable) (*Profile, error)
//    - Uses sampling for large tables
//    - Handles different data types appropriately
//    - Calculates statistics in a memory-efficient way
// 4. Type-specific profiling helpers:
//    - profileNumericColumn(values []interface{}) NumericStats
//    - profileStringColumn(values []interface{}) StringStats
//    - profileDateColumn(values []interface{}) DateStats
//    - detectPatterns(values []string) []PatternCount
// 5. Helper functions:
//    - calculateHistogram(values []float64, bins int) []HistogramBin
//    - detectOutliers(values []float64) []int
//    - calculateUniquePercentage(values []interface{}) float64
// Ensure efficient memory usage for large tables
// Use goroutines to profile multiple columns concurrently
// Handle all common data types: string, int, float, bool, date/time

package profile

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/nessi-dev/nessi/pkg"
)

// Profile represents a complete data profile
type Profile struct {
	TableName      string
	RowCount       int64
	ColumnProfiles map[string]*ColumnProfile
	CreatedAt      time.Time
	ExecutionTimeMs int64
}

// ColumnProfile represents detailed statistics for a column
type ColumnProfile struct {
	Name            string
	DataType        string
	NullCount       int64
	NullPercentage  float64
	UniqueCount     int64
	UniquePercentage float64
	Min             interface{}
	Max             interface{}
	Mean            float64
	Median          float64
	StdDev          float64
	TopValues       []ValueCount
	Histogram       []HistogramBin
	Patterns        []PatternCount
}

// ValueCount represents a value and its frequency
type ValueCount struct {
	Value     interface{}
	Count     int64
	Frequency float64
}

// HistogramBin represents a bin in a histogram
type HistogramBin struct {
	Start     float64
	End       float64
	Count     int64
	Frequency float64
}

// PatternCount represents a pattern and its frequency
type PatternCount struct {
	Pattern   string
	Count     int64
	Frequency float64
}

// NumericStats represents statistics for numeric columns
type NumericStats struct {
	Min     float64
	Max     float64
	Mean    float64
	Median  float64
	StdDev  float64
	Outliers []int
}

// StringStats represents statistics for string columns
type StringStats struct {
	MinLength    int
	MaxLength    int
	AvgLength    float64
	Patterns     []PatternCount
	TopValues    []ValueCount
}

// DateStats represents statistics for date columns
type DateStats struct {
	MinDate     time.Time
	MaxDate     time.Time
	AvgInterval time.Duration
	Outliers    []int
}

// QualityProfile represents the quality profile of a dataset
type QualityProfile struct {
	TablePath     string
	Timestamp     time.Time
	TotalRows     int64
	TotalColumns  int
	ColumnProfiles map[string]*ColumnProfile
	DataQuality   *DataQuality
}

// DataQuality represents overall data quality metrics
type DataQuality struct {
	Completeness   float64
	Consistency    float64
	Uniqueness     float64
	Timeliness     float64
	Validity       float64
	OverallScore   float64
	Issues         []*QualityIssue
}

// Anomaly represents a data anomaly in a column
type Anomaly struct {
	Type        string
	Description string
	Severity    string
	Count       int64
	Examples    []interface{}
}

// QualityIssue represents a data quality issue
type QualityIssue struct {
	Column     string
	Type       string
	Description string
	Severity    string
	Impact      string
	Recommendation string
}

// OutlierDetectionMethod defines the method used for outlier detection
type OutlierDetectionMethod string

const (
	// ZScoreMethod uses z-score (standard deviations from mean) for outlier detection
	ZScoreMethod OutlierDetectionMethod = "zscore"
	// IQRMethod uses interquartile range for outlier detection
	IQRMethod OutlierDetectionMethod = "iqr"
)

// ProfilerConfig represents profiler configuration
type ProfilerConfig struct {
	SampleSize          int64
	MaxTopValues        int
	NumHistogramBins    int
	MinPatternFrequency float64
	OutlierThreshold    float64
	IQRMultiplier       float64
	OutlierMethod       OutlierDetectionMethod
	SamplingRate        float64
	MaxRows             int64
}

// Profiler handles data profiling
type Profiler struct {
	reader *pkg.Reader
	config *ProfilerConfig
	concurrency int
}

// NewProfiler creates a new data profiler
func NewProfiler(tablePath string, config *ProfilerConfig) (*Profiler, error) {
	reader, err := pkg.NewReader(tablePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}

	if config == nil {
		config = &ProfilerConfig{
			SampleSize:          10000,
			MaxTopValues:        10,
			NumHistogramBins:    20,
			MinPatternFrequency: 0.01,
			OutlierThreshold:    3.0,
			IQRMultiplier:       1.5,
			OutlierMethod:       ZScoreMethod, // Default to z-score method
			SamplingRate:        1.0,
			MaxRows:             100000,
		}
	}

	return &Profiler{
		reader: reader,
		config: config,
		concurrency: 4, // default concurrency
	}, nil
}

// ProfileTable profiles a complete Delta table.
// It reads data from the reader and applies the configuration.
func (p *Profiler) ProfileTable() (*Profile, error) {
	startTime := time.Now()

	schemaMap := p.reader.GetSchema() // Returns map[string]string
	if schemaMap == nil || len(schemaMap) == 0 { // Check if schema is empty or nil
		return nil, fmt.Errorf("failed to get a valid schema from reader")
	}

	// Use the new ReadAllParsed method to get structured records
	records, err := p.reader.ReadAllParsed()
	if err != nil {
		return nil, fmt.Errorf("failed to read and parse records: %w", err)
	}

	if len(records) == 0 { // This will currently be true
		// Handle empty table or sample (or in this case, unparsed data)
		return &Profile{
			TableName:      p.reader.GetTablePath(),
			RowCount:       0, // Actual row count would come from stats or parsed records
			ColumnProfiles: make(map[string]*ColumnProfile),
			CreatedAt:      time.Now(),
			ExecutionTimeMs: time.Since(startTime).Milliseconds(),
		}, nil
	}

	columnProfilesMap := make(map[string]*ColumnProfile)
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, p.concurrency)

	for fieldName, fieldType := range schemaMap { // Iterate over the map
		wg.Add(1)
		sem <- struct{}{}
		go func(fName string, fType string) {
			defer wg.Done()
			defer func() { <-sem }()

			// This part assumes 'records' is populated. With current changes, 'values' will be empty.
			values := make([]interface{}, len(records))
			for i, rec := range records {
				values[i] = rec[fName]
			}
			cp := p.profileColumn(fName, fType, values)
			mu.Lock()
			columnProfilesMap[fName] = cp
			mu.Unlock()
		}(fieldName, fieldType)
	}

	wg.Wait()

	return &Profile{
		TableName:      p.reader.GetTablePath(),
		RowCount:       int64(len(records)), // Will be 0 for now
		ColumnProfiles: columnProfilesMap,
		CreatedAt:      time.Now(),
		ExecutionTimeMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// GetDataQuality generates the profile and then calculates data quality metrics.
func (p *Profiler) GetDataQuality() (*QualityProfile, error) {
	profile, err := p.ProfileTable()
	if err != nil {
		return nil, fmt.Errorf("failed to generate profile for quality assessment: %w", err)
	}

	// The QualityProfile struct seems more appropriate to return here,
	// as it's designed to hold both ColumnProfiles and DataQuality.
	// calculateDataQualityMetrics returns *DataQuality.
	// We need to embed this into a QualityProfile.

	dataQuality := calculateDataQualityMetrics(profile) // This is a helper now

	// Construct the QualityProfile to return
	qp := &QualityProfile{
		TablePath:     profile.TableName, // Or p.reader.GetTablePath() directly
		Timestamp:     time.Now(),        // Or profile.CreatedAt
		TotalRows:     profile.RowCount,  // This would be 0 if data parsing is not yet implemented
		TotalColumns:  len(profile.ColumnProfiles),
		ColumnProfiles: profile.ColumnProfiles,
		DataQuality:   dataQuality,
	}

	return qp, nil
}

// profileColumn profiles a single column.
func (p *Profiler) profileColumn(name string, dataType string, values []interface{}) *ColumnProfile {
	// Read column data
	// values, err := p.readColumnData(name)
	// if err != nil {
	// 	return nil, err
	// }

	// Create column profile
	profile := &ColumnProfile{
		Name:     name,
		DataType: dataType,
	}

	// Calculate basic statistics
	profile.NullCount = countNulls(values)
	profile.NullPercentage = float64(profile.NullCount) / float64(len(values))
	profile.UniqueCount = countUnique(values)
	profile.UniquePercentage = float64(profile.UniqueCount) / float64(len(values))

	// Calculate type-specific statistics
	switch dataType {
	case "string":
		stats := p.profileStringColumn(values)
		profile.Patterns = stats.Patterns
		profile.TopValues = stats.TopValues
	case "integer", "float":
		stats := p.profileNumericColumn(values)
		profile.Min = stats.Min
		profile.Max = stats.Max
		profile.Mean = stats.Mean
		profile.Median = stats.Median
		profile.StdDev = stats.StdDev
		profile.Histogram = calculateHistogram(convertToFloat64(values), p.config.NumHistogramBins)
	case "date", "timestamp":
		stats := p.profileDateColumn(values)
		profile.Min = stats.MinDate
		profile.Max = stats.MaxDate
	}

	return profile
}

// readColumnData reads data for a specific column
func (p *Profiler) readColumnData(name string) ([]interface{}, error) {
	// TODO: Implement actual column data reading
	// This would involve:
	// 1. Reading data from Delta table
	// 2. Sampling if needed
	// 3. Converting to appropriate types
	return nil, fmt.Errorf("column data reading not implemented")
}

// profileNumericColumn profiles a numeric column
func (p *Profiler) profileNumericColumn(values []interface{}) NumericStats {
	// Convert values to float64
	floats := convertToFloat64(values)

	// Calculate basic statistics
	stats := NumericStats{
		Min:    floats[0],
		Max:    floats[0],
		Mean:   0,
		StdDev: 0,
	}

	// Calculate mean
	var sum float64
	for _, v := range floats {
		sum += v
		if v < stats.Min {
			stats.Min = v
		}
		if v > stats.Max {
			stats.Max = v
		}
	}
	stats.Mean = sum / float64(len(floats))

	// Calculate standard deviation
	var sumSquares float64
	for _, v := range floats {
		dev := v - stats.Mean
		sumSquares += dev * dev
	}
	stats.StdDev = math.Sqrt(sumSquares / float64(len(floats)))

	// Calculate median
	sort.Float64s(floats)
	mid := len(floats) / 2
	if len(floats)%2 == 0 {
		stats.Median = (floats[mid-1] + floats[mid]) / 2
	} else {
		stats.Median = floats[mid]
	}

	// Detect outliers based on the configured method
	switch p.config.OutlierMethod {
	case IQRMethod:
		stats.Outliers = detectOutliersIQR(floats, p.config.IQRMultiplier)
	default: // ZScoreMethod or any other value defaults to z-score
		stats.Outliers = detectOutliers(floats, p.config.OutlierThreshold)
	}

	return stats
}

// profileStringColumn profiles a string column
func (p *Profiler) profileStringColumn(values []interface{}) StringStats {
	// Convert values to strings
	strings := make([]string, len(values))
	for i, v := range values {
		if v != nil {
			strings[i] = fmt.Sprint(v)
		}
	}

	// Calculate length statistics
	stats := StringStats{
		MinLength: math.MaxInt32,
		MaxLength: 0,
		AvgLength: 0,
	}

	var totalLength int
	for _, s := range strings {
		if s != "" {
			length := len(s)
			if length < stats.MinLength {
				stats.MinLength = length
			}
			if length > stats.MaxLength {
				stats.MaxLength = length
			}
			totalLength += length
		}
	}
	stats.AvgLength = float64(totalLength) / float64(len(strings))

	// Detect patterns
	stats.Patterns = detectPatterns(strings)

	// Calculate top values
	stats.TopValues = calculateTopValues(values, p.config.MaxTopValues)

	return stats
}

// profileDateColumn profiles a date column
func (p *Profiler) profileDateColumn(values []interface{}) DateStats {
	// Convert values to time.Time
	dates := make([]time.Time, 0, len(values))
	for _, v := range values {
		if t, ok := v.(time.Time); ok {
			dates = append(dates, t)
		}
	}

	// Calculate date statistics
	stats := DateStats{
		MinDate: dates[0],
		MaxDate: dates[0],
	}

	// Calculate min/max dates and intervals
	var totalInterval time.Duration
	for i := 1; i < len(dates); i++ {
		if dates[i].Before(stats.MinDate) {
			stats.MinDate = dates[i]
		}
		if dates[i].After(stats.MaxDate) {
			stats.MaxDate = dates[i]
		}
		if i > 0 {
			totalInterval += dates[i].Sub(dates[i-1])
		}
	}
	stats.AvgInterval = totalInterval / time.Duration(len(dates)-1)

	// Detect outliers based on the configured method
	switch p.config.OutlierMethod {
	case IQRMethod:
		stats.Outliers = detectDateOutliersIQR(dates, p.config.IQRMultiplier)
	default: // ZScoreMethod or any other value defaults to z-score
		stats.Outliers = detectDateOutliers(dates, p.config.OutlierThreshold)
	}

	return stats
}

// calculateHistogram creates a histogram for numeric values
func calculateHistogram(values []float64, bins int) []HistogramBin {
	if len(values) == 0 {
		return nil
	}

	// Sort values
	sort.Float64s(values)

	// Calculate bin size
	min := values[0]
	max := values[len(values)-1]
	binSize := (max - min) / float64(bins)

	// Create bins
	histogram := make([]HistogramBin, bins)
	for i := range histogram {
		histogram[i] = HistogramBin{
			Start: min + float64(i)*binSize,
			End:   min + float64(i+1)*binSize,
		}
	}

	// Count values in each bin
	for _, v := range values {
		bin := int((v - min) / binSize)
		if bin >= bins {
			bin = bins - 1
		}
		histogram[bin].Count++
	}

	// Calculate frequencies
	total := float64(len(values))
	for i := range histogram {
		histogram[i].Frequency = float64(histogram[i].Count) / total
	}

	return histogram
}

// detectOutliers identifies outliers in numeric data
func detectOutliers(values []float64, threshold float64) []int {
	if len(values) < 2 {
		return nil
	}

	// Calculate mean and standard deviation
	var sum, sumSquares float64
	for _, v := range values {
		sum += v
		sumSquares += v * v
	}
	mean := sum / float64(len(values))
	stdDev := math.Sqrt(sumSquares/float64(len(values)) - mean*mean)

	// Find outliers
	outliers := make([]int, 0)
	for i, v := range values {
		zScore := math.Abs((v - mean) / stdDev)
		if zScore > threshold {
			outliers = append(outliers, i)
		}
	}

	return outliers
}

// detectDateOutliers identifies outliers in date data using z-score method
func detectDateOutliers(dates []time.Time, threshold float64) []int {
	if len(dates) < 2 {
		return nil
	}

	// Convert dates to durations from earliest date
	earliest := dates[0]
	for _, date := range dates {
		if date.Before(earliest) {
			earliest = date
		}
	}

	// Convert to durations in hours
	durations := make([]float64, len(dates))
	for i, date := range dates {
		duration := date.Sub(earliest)
		durations[i] = duration.Hours()
	}

	// Use the same outlier detection as for numeric values
	return detectOutliers(durations, threshold)
}

// detectOutliersIQR identifies outliers in numeric data using the Interquartile Range (IQR) method
// Values outside of Q1 - multiplier*IQR and Q3 + multiplier*IQR are considered outliers
// A typical multiplier value is 1.5 for outliers and 3.0 for extreme outliers
func detectOutliersIQR(values []float64, multiplier float64) []int {
	if len(values) < 4 { // Need at least 4 values to calculate meaningful quartiles
		return nil
	}

	// Make a copy and sort the values
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	// Calculate quartiles
	n := len(sorted)
	q1Index := n / 4
	q3Index := n * 3 / 4

	// Handle even number of elements
	var q1, q3 float64
	if n%4 == 0 {
		q1 = (sorted[q1Index-1] + sorted[q1Index]) / 2
		q3 = (sorted[q3Index-1] + sorted[q3Index]) / 2
	} else {
		q1 = sorted[q1Index]
		q3 = sorted[q3Index]
	}

	// Calculate IQR and bounds
	iqr := q3 - q1
	lowerBound := q1 - multiplier*iqr
	upperBound := q3 + multiplier*iqr

	// Find outliers
	outliers := make([]int, 0)
	for i, v := range values {
		if v < lowerBound || v > upperBound {
			outliers = append(outliers, i)
		}
	}

	return outliers
}

// detectDateOutliersIQR identifies outliers in date data using the IQR method
func detectDateOutliersIQR(dates []time.Time, multiplier float64) []int {
	if len(dates) < 4 {
		return nil
	}

	// Convert dates to durations from earliest date
	earliest := dates[0]
	for _, date := range dates {
		if date.Before(earliest) {
			earliest = date
		}
	}

	// Convert to durations in hours
	durations := make([]float64, len(dates))
	for i, date := range dates {
		duration := date.Sub(earliest)
		durations[i] = duration.Hours()
	}

	// Use the IQR outlier detection
	return detectOutliersIQR(durations, multiplier)
}

// detectPatterns identifies patterns in string data
func detectPatterns(values []string) []PatternCount {
	patterns := make(map[string]int64)
	for _, v := range values {
		if v != "" {
			pattern := extractPattern(v)
			patterns[pattern]++
		}
	}

	// Convert to slice and sort by frequency
	result := make([]PatternCount, 0, len(patterns))
	for pattern, count := range patterns {
		result = append(result, PatternCount{
			Pattern:   pattern,
			Count:     count,
			Frequency: float64(count) / float64(len(values)),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}

// extractPattern extracts a pattern from a string
func extractPattern(s string) string {
	// TODO: Implement pattern extraction
	// This would involve:
	// 1. Identifying character types (letter, digit, etc.)
	// 2. Creating a pattern string
	return s
}

// calculateTopValues finds the most common values
func calculateTopValues(values []interface{}, maxValues int) []ValueCount {
	counts := make(map[interface{}]int64)
	for _, v := range values {
		if v != nil {
			counts[v]++
		}
	}

	// Convert to slice and sort by count
	result := make([]ValueCount, 0, len(counts))
	for value, count := range counts {
		result = append(result, ValueCount{
			Value:     value,
			Count:     count,
			Frequency: float64(count) / float64(len(values)),
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	// Limit to maxValues
	if len(result) > maxValues {
		result = result[:maxValues]
	}

	return result
}

// Helper functions
func countNulls(values []interface{}) int64 {
	var count int64
	for _, v := range values {
		if v == nil {
			count++
		}
	}
	return count
}

func countUnique(values []interface{}) int64 {
	unique := make(map[interface{}]bool)
	for _, v := range values {
		if v != nil {
			unique[v] = true
		}
	}
	return int64(len(unique))
}

func convertToFloat64(values []interface{}) []float64 {
	result := make([]float64, 0, len(values))
	for _, v := range values {
		if v != nil {
			switch n := v.(type) {
			case float64:
				result = append(result, n)
			case float32:
				result = append(result, float64(n))
			case int:
				result = append(result, float64(n))
			case int64:
				result = append(result, float64(n))
			}
		}
	}
	return result
}

// Close closes the profiler and releases resources
func (p *Profiler) Close() error {
	return p.reader.Close()
}

// calculateDataQualityMetrics calculates overall data quality metrics
func calculateDataQualityMetrics(profile *Profile) *DataQuality {
	// TODO: Implement actual data quality calculation
	// This would involve:
	// 1. Calculating completeness, consistency, uniqueness, timeliness, and validity
	// 2. Identifying data quality issues
	return &DataQuality{
		Completeness:   1.0,
		Consistency:    1.0,
		Uniqueness:     1.0,
		Timeliness:     1.0,
		Validity:       1.0,
		OverallScore:   1.0,
		Issues:         []*QualityIssue{},
	}
}