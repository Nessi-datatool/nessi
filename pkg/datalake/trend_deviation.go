package datalake

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

// MetricType defines the type of metric to compare
type MetricType string

const (
	// NullPercentage is the percentage of null values
	NullPercentage MetricType = "null_percentage"

	// UniqueRatio is the ratio of unique values
	UniqueRatio MetricType = "unique_ratio"

	// MinValue is the minimum value
	MinValue MetricType = "min_value"

	// MaxValue is the maximum value
	MaxValue MetricType = "max_value"

	// MeanValue is the mean value
	MeanValue MetricType = "mean_value"

	// MedianValue is the median value
	MedianValue MetricType = "median_value"

	// StandardDeviation is the standard deviation
	StandardDeviation MetricType = "standard_deviation"

	// RecordCount is the total number of records
	RecordCount MetricType = "record_count"

	// InvalidCount is the number of invalid values
	InvalidCount MetricType = "invalid_count"
)

// AlertSeverity defines the severity of an alert
type AlertSeverity string

const (
	// InfoAlert is an informational alert
	InfoAlert AlertSeverity = "info"

	// WarningAlert is a warning alert
	WarningAlert AlertSeverity = "warning"

	// ErrorAlert is an error alert
	ErrorAlert AlertSeverity = "error"

	// CriticalAlert is a critical alert
	CriticalAlert AlertSeverity = "critical"
)

// TrendDeviationOptions defines the options for trend deviation analysis
type TrendDeviationOptions struct {
	// The field to analyze
	Field string

	// The metric types to compare
	MetricTypes []MetricType

	// The threshold for alerting (percentage change)
	Threshold float64

	// The number of previous runs to compare with
	PreviousRuns int

	// The directory where previous run metrics are stored
	MetricsDir string
}

// TrendDeviationResult represents the result of a trend deviation analysis
type TrendDeviationResult struct {
	// The field that was analyzed
	Field string

	// The metrics that were compared
	Metrics []MetricComparison

	// The alerts that were generated
	Alerts []TrendAlert

	// The timestamp when the analysis was performed
	Timestamp time.Time
}

// MetricComparison represents a comparison of a metric between runs
type MetricComparison struct {
	// The metric type
	MetricType MetricType

	// The current value
	CurrentValue float64

	// The previous values
	PreviousValues []float64

	// The percentage change from the previous run
	PercentageChange float64

	// The average value across all previous runs
	AverageValue float64

	// The standard deviation across all previous runs
	StandardDeviation float64

	// The z-score of the current value
	ZScore float64
}

// TrendAlert represents an alert for a significant trend deviation
type TrendAlert struct {
	// The metric type
	MetricType MetricType

	// The field that triggered the alert
	Field string

	// The severity of the alert
	Severity AlertSeverity

	// The message describing the alert
	Message string

	// The current value
	CurrentValue float64

	// The previous value
	PreviousValue float64

	// The percentage change
	PercentageChange float64

	// The threshold that was exceeded
	Threshold float64
}

// RunMetrics represents the metrics for a single run
type RunMetrics struct {
	// The timestamp of the run
	Timestamp time.Time

	// The metrics for each field
	FieldMetrics map[string]map[string]float64
}

// AnalyzeTrendDeviation analyzes trend deviations for a field
func (m *MetadataManager) AnalyzeTrendDeviation(options TrendDeviationOptions) (*TrendDeviationResult, error) {
	// Create result
	result := &TrendDeviationResult{
		Field:     options.Field,
		Metrics:   []MetricComparison{},
		Alerts:    []TrendAlert{},
		Timestamp: time.Now(),
	}

	// Get current metrics
	currentMetrics, err := m.calculateFieldMetrics(options.Field, options.MetricTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate current metrics: %w", err)
	}

	// Load previous run metrics
	previousRuns, err := loadPreviousRunMetrics(options.MetricsDir, options.PreviousRuns)
	if err != nil {
		return nil, fmt.Errorf("failed to load previous run metrics: %w", err)
	}

	// Compare metrics and generate alerts
	for _, metricType := range options.MetricTypes {
		// Get current value
		currentValue, ok := currentMetrics[string(metricType)]
		if !ok {
			continue // Skip if metric is not available
		}

		// Get previous values
		previousValues := []float64{}
		for _, run := range previousRuns {
			fieldMetrics, ok := run.FieldMetrics[options.Field]
			if !ok {
				continue
			}

			value, ok := fieldMetrics[string(metricType)]
			if !ok {
				continue
			}

			previousValues = append(previousValues, value)
		}

		// Skip if no previous values
		if len(previousValues) == 0 {
			continue
		}

		// Calculate statistics
		avgValue := CalculateAverage(previousValues)
		stdDev := CalculateStandardDeviation(previousValues, avgValue)
		zScore := CalculateZScore(currentValue, avgValue, stdDev)

		// Calculate percentage change from most recent previous run
		percentageChange := 0.0
		if len(previousValues) > 0 {
			mostRecentValue := previousValues[0]
			if mostRecentValue != 0 {
				percentageChange = (currentValue - mostRecentValue) / math.Abs(mostRecentValue) * 100
			}
		}

		// Add metric comparison
		comparison := MetricComparison{
			MetricType:        metricType,
			CurrentValue:      currentValue,
			PreviousValues:    previousValues,
			PercentageChange:  percentageChange,
			AverageValue:      avgValue,
			StandardDeviation: stdDev,
			ZScore:            zScore,
		}
		result.Metrics = append(result.Metrics, comparison)

		// Check for significant deviations
		if math.Abs(percentageChange) >= options.Threshold {
			// Determine severity based on percentage change
			severity := InfoAlert
			if math.Abs(percentageChange) >= options.Threshold*2 {
				severity = WarningAlert
			}
			if math.Abs(percentageChange) >= options.Threshold*3 {
				severity = ErrorAlert
			}
			if math.Abs(percentageChange) >= options.Threshold*4 {
				severity = CriticalAlert
			}

			// Create alert
			alert := TrendAlert{
				MetricType:       metricType,
				Field:            options.Field,
				Severity:         severity,
				Message:          fmt.Sprintf("%s for field '%s' has changed by %.2f%% (threshold: %.2f%%)", metricType, options.Field, percentageChange, options.Threshold),
				CurrentValue:     currentValue,
				PreviousValue:    previousValues[0],
				PercentageChange: percentageChange,
				Threshold:        options.Threshold,
			}
			result.Alerts = append(result.Alerts, alert)
		}

		// Check for significant z-score deviations
		if math.Abs(zScore) >= 2.0 && len(previousValues) >= 3 {
			// Determine severity based on z-score
			severity := InfoAlert
			if math.Abs(zScore) >= 2.0 {
				severity = WarningAlert
			}
			if math.Abs(zScore) >= 3.0 {
				severity = ErrorAlert
			}
			if math.Abs(zScore) >= 4.0 {
				severity = CriticalAlert
			}

			// Create alert
			alert := TrendAlert{
				MetricType:       metricType,
				Field:            options.Field,
				Severity:         severity,
				Message:          fmt.Sprintf("%s for field '%s' has a z-score of %.2f (current: %.2f, avg: %.2f, std: %.2f)", metricType, options.Field, zScore, currentValue, avgValue, stdDev),
				CurrentValue:     currentValue,
				PreviousValue:    avgValue,
				PercentageChange: percentageChange,
				Threshold:        options.Threshold,
			}
			result.Alerts = append(result.Alerts, alert)
		}
	}

	// Save current metrics
	err = saveCurrentRunMetrics(options.MetricsDir, options.Field, currentMetrics)
	if err != nil {
		return nil, fmt.Errorf("failed to save current run metrics: %w", err)
	}

	return result, nil
}

// calculateFieldMetrics calculates metrics for a field
func (m *MetadataManager) calculateFieldMetrics(field string, metricTypes []MetricType) (map[string]float64, error) {
	// Read the latest table metadata
	table, err := m.ReadTableMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to read table metadata: %w", err)
	}

	// Check if the field exists in the schema
	fieldExists := false
	for _, schemaField := range table.Schema.Fields() {
		if schemaField.Name == field {
			fieldExists = true
			break
		}
	}

	if !fieldExists {
		return nil, fmt.Errorf("field '%s' does not exist in the table schema", field)
	}

	// Initialize metrics
	metrics := make(map[string]float64)

	// Read all parquet files and calculate metrics
	allValues := []interface{}{}
	nullCount := 0
	uniqueValues := make(map[interface{}]bool)

	for _, file := range table.Files {
		// Read the parquet file
		records, err := m.ReadParquetFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read parquet file '%s': %w", file, err)
		}

		// Process each record
		for _, record := range records {
			// Get the field value
			value, ok := record[field]
			if !ok || value == nil {
				nullCount++
				continue
			}

			// Add to all values
			allValues = append(allValues, value)

			// Add to unique values
			uniqueValues[value] = true
		}
	}

	// Calculate metrics
	totalCount := len(allValues) + nullCount

	// Record count
	if contains(metricTypes, RecordCount) {
		metrics[string(RecordCount)] = float64(totalCount)
	}

	// Null percentage
	if contains(metricTypes, NullPercentage) {
		if totalCount > 0 {
			metrics[string(NullPercentage)] = float64(nullCount) / float64(totalCount) * 100
		} else {
			metrics[string(NullPercentage)] = 0
		}
	}

	// Unique ratio
	if contains(metricTypes, UniqueRatio) {
		if len(allValues) > 0 {
			metrics[string(UniqueRatio)] = float64(len(uniqueValues)) / float64(len(allValues))
		} else {
			metrics[string(UniqueRatio)] = 0
		}
	}

	// Numeric metrics
	numericValues := []float64{}
	for _, value := range allValues {
		// Convert to float64
		switch v := value.(type) {
		case int:
			numericValues = append(numericValues, float64(v))
		case int32:
			numericValues = append(numericValues, float64(v))
		case int64:
			numericValues = append(numericValues, float64(v))
		case float32:
			numericValues = append(numericValues, float64(v))
		case float64:
			numericValues = append(numericValues, v)
		}
	}

	// Skip numeric metrics if no numeric values
	if len(numericValues) == 0 {
		return metrics, nil
	}

	// Min value
	if contains(metricTypes, MinValue) {
		min := numericValues[0]
		for _, value := range numericValues {
			if value < min {
				min = value
			}
		}
		metrics[string(MinValue)] = min
	}

	// Max value
	if contains(metricTypes, MaxValue) {
		max := numericValues[0]
		for _, value := range numericValues {
			if value > max {
				max = value
			}
		}
		metrics[string(MaxValue)] = max
	}

	// Mean value
	if contains(metricTypes, MeanValue) {
		metrics[string(MeanValue)] = CalculateAverage(numericValues)
	}

	// Median value
	if contains(metricTypes, MedianValue) {
		metrics[string(MedianValue)] = CalculateMedian(numericValues)
	}

	// Standard deviation
	if contains(metricTypes, StandardDeviation) {
		avg := CalculateAverage(numericValues)
		metrics[string(StandardDeviation)] = CalculateStandardDeviation(numericValues, avg)
	}

	return metrics, nil
}

// loadPreviousRunMetrics loads metrics from previous runs
func loadPreviousRunMetrics(metricsDir string, maxRuns int) ([]RunMetrics, error) {
	// Create metrics directory if it doesn't exist
	err := os.MkdirAll(metricsDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics directory: %w", err)
	}

	// Get all metric files
	files, err := filepath.Glob(filepath.Join(metricsDir, "metrics_*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list metric files: %w", err)
	}

	// Sort files by modification time (newest first)
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			fileInfoI, err := os.Stat(files[i])
			if err != nil {
				continue
			}

			fileInfoJ, err := os.Stat(files[j])
			if err != nil {
				continue
			}

			if fileInfoI.ModTime().Before(fileInfoJ.ModTime()) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	// Limit to maxRuns
	if len(files) > maxRuns {
		files = files[:maxRuns]
	}

	// Load metrics from each file
	runs := []RunMetrics{}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var run RunMetrics
		err = json.Unmarshal(data, &run)
		if err != nil {
			continue
		}

		runs = append(runs, run)
	}

	return runs, nil
}

// saveCurrentRunMetrics saves metrics from the current run
func saveCurrentRunMetrics(metricsDir string, field string, metrics map[string]float64) error {
	// Create metrics directory if it doesn't exist
	err := os.MkdirAll(metricsDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create metrics directory: %w", err)
	}

	// Create run metrics
	run := RunMetrics{
		Timestamp: time.Now(),
		FieldMetrics: map[string]map[string]float64{
			field: metrics,
		},
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(run, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal run metrics to JSON: %w", err)
	}

	// Write to file
	filename := filepath.Join(metricsDir, fmt.Sprintf("metrics_%s.json", time.Now().Format("20060102_150405")))
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write run metrics to file: %w", err)
	}

	return nil
}

// CalculateAverage calculates the average of a slice of float64 values
func CalculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, value := range values {
		sum += value
	}

	return sum / float64(len(values))
}

// CalculateMedian calculates the median of a slice of float64 values
func CalculateMedian(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Sort values
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	for i := 0; i < len(sortedValues); i++ {
		for j := i + 1; j < len(sortedValues); j++ {
			if sortedValues[i] > sortedValues[j] {
				sortedValues[i], sortedValues[j] = sortedValues[j], sortedValues[i]
			}
		}
	}

	// Calculate median
	if len(sortedValues)%2 == 0 {
		// Even number of values
		middle1 := sortedValues[len(sortedValues)/2-1]
		middle2 := sortedValues[len(sortedValues)/2]
		return (middle1 + middle2) / 2
	} else {
		// Odd number of values
		return sortedValues[len(sortedValues)/2]
	}
}

// CalculateStandardDeviation calculates the standard deviation of a slice of float64 values
func CalculateStandardDeviation(values []float64, avg float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	sumSquaredDiff := 0.0
	for _, value := range values {
		diff := value - avg
		sumSquaredDiff += diff * diff
	}

	return math.Sqrt(sumSquaredDiff / float64(len(values)-1))
}

// CalculateZScore calculates the z-score of a value
func CalculateZScore(value, avg, stdDev float64) float64 {
	if stdDev == 0 {
		return 0
	}

	return (value - avg) / stdDev
}

// SimpleTrendAnalyzer is a simplified implementation for demonstrating trend deviation
type SimpleTrendAnalyzer struct {
	Metrics    map[string]map[string]float64
	MetricsDir string
}

// NewSimpleTrendAnalyzer creates a new simple trend analyzer
func NewSimpleTrendAnalyzer(metricsDir string) *SimpleTrendAnalyzer {
	return &SimpleTrendAnalyzer{
		Metrics: map[string]map[string]float64{
			"sales": {
				string(MeanValue):         300.0,
				string(StandardDeviation): 158.11,
				string(MinValue):          100.0,
				string(MaxValue):          500.0,
				string(RecordCount):       5.0,
			},
			"revenue": {
				string(MeanValue):         3000.0,
				string(StandardDeviation): 1581.1,
				string(MinValue):          1000.0,
				string(MaxValue):          5000.0,
				string(RecordCount):       5.0,
			},
			"customers": {
				string(MeanValue):         30.0,
				string(StandardDeviation): 15.81,
				string(MinValue):          10.0,
				string(MaxValue):          50.0,
				string(RecordCount):       5.0,
			},
			"null_field": {
				string(NullPercentage): 80.0,
				string(RecordCount):    5.0,
			},
		},
		MetricsDir: metricsDir,
	}
}

// AnalyzeTrendDeviation simulates the trend deviation analysis
func (s *SimpleTrendAnalyzer) AnalyzeTrendDeviation(field string, runNumber int) (*TrendDeviationResult, error) {
	// Check if field exists
	fieldMetrics, ok := s.Metrics[field]
	if !ok {
		return nil, fmt.Errorf("field '%s' does not exist", field)
	}

	// Create a copy of the metrics and modify them based on the run number
	currentMetrics := make(map[string]float64)
	for k, v := range fieldMetrics {
		currentMetrics[k] = v
	}

	// Modify metrics based on run number
	switch {
	case field == "sales" && runNumber == 1:
		// Increase mean value by 10%
		currentMetrics[string(MeanValue)] *= 1.1
		// Increase max value by 10%
		currentMetrics[string(MaxValue)] *= 1.1
	case field == "sales" && runNumber == 2:
		// Increase mean value by another 10%
		currentMetrics[string(MeanValue)] *= 1.21 // 1.1 * 1.1
		// Increase max value by another 10%
		currentMetrics[string(MaxValue)] *= 1.21 // 1.1 * 1.1
	case field == "revenue" && runNumber == 2:
		// Increase mean value by 20%
		currentMetrics[string(MeanValue)] *= 1.2
		// Increase max value by 20%
		currentMetrics[string(MaxValue)] *= 1.2
	case field == "customers" && runNumber == 3:
		// Decrease mean value by 30%
		currentMetrics[string(MeanValue)] *= 0.7
		// Decrease max value by 30%
		currentMetrics[string(MaxValue)] *= 0.7
	}

	// Create metrics directory if it doesn't exist
	if err := os.MkdirAll(s.MetricsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create metrics directory: %w", err)
	}

	// Create previous runs data
	previousRuns := []RunMetrics{}
	for i := 0; i < runNumber; i++ {
		prevRunMetrics := RunMetrics{
			Timestamp: time.Now().Add(-time.Duration(runNumber-i) * time.Hour),
			FieldMetrics: map[string]map[string]float64{
				field: make(map[string]float64),
			},
		}

		// Copy original metrics
		for k, v := range s.Metrics[field] {
			prevRunMetrics.FieldMetrics[field][k] = v
		}

		// Modify based on run number
		if field == "sales" && i >= 1 {
			// Increase mean value by 10% for each run
			prevRunMetrics.FieldMetrics[field][string(MeanValue)] *= math.Pow(1.1, float64(i))
			// Increase max value by 10% for each run
			prevRunMetrics.FieldMetrics[field][string(MaxValue)] *= math.Pow(1.1, float64(i))
		}

		if field == "revenue" && i >= 2 {
			// Increase mean value by 20% for run 2
			prevRunMetrics.FieldMetrics[field][string(MeanValue)] *= 1.2
			// Increase max value by 20% for run 2
			prevRunMetrics.FieldMetrics[field][string(MaxValue)] *= 1.2
		}

		previousRuns = append(previousRuns, prevRunMetrics)
	}

	// Create the result
	result := &TrendDeviationResult{
		Field:     field,
		Timestamp: time.Now(),
		Metrics:   []MetricComparison{},
		Alerts:    []TrendAlert{},
	}

	// If no previous runs, return empty result
	if len(previousRuns) == 0 {
		return result, nil
	}

	// Compare metrics with previous runs
	threshold := 10.0 // Default threshold

	// Analyze each metric
	for metricType, currentValue := range currentMetrics {
		// Skip record count for now
		if metricType == string(RecordCount) {
			continue
		}

		// Get previous values
		previousValues := []float64{}
		for _, prevRun := range previousRuns {
			if prevMetrics, ok := prevRun.FieldMetrics[field]; ok {
				if prevValue, ok := prevMetrics[metricType]; ok {
					previousValues = append(previousValues, prevValue)
				}
			}
		}

		// Skip if no previous values
		if len(previousValues) == 0 {
			continue
		}

		// Calculate percentage change from most recent previous run
		percentageChange := ((currentValue - previousValues[0]) / previousValues[0]) * 100

		// Create metric comparison
		metricComparison := MetricComparison{
			MetricType:       MetricType(metricType),
			CurrentValue:     currentValue,
			PreviousValues:   previousValues,
			PercentageChange: percentageChange,
		}

		// Calculate statistics if we have multiple previous values
		if len(previousValues) > 1 {
			avgValue := CalculateAverage(previousValues)
			stdDev := CalculateStandardDeviation(previousValues, avgValue)
			zScore := CalculateZScore(currentValue, avgValue, stdDev)

			metricComparison.AverageValue = avgValue
			metricComparison.StandardDeviation = stdDev
			metricComparison.ZScore = zScore
		}

		result.Metrics = append(result.Metrics, metricComparison)

		// Check if we need to generate an alert
		if math.Abs(percentageChange) >= threshold {
			// Determine severity based on percentage change
			severity := InfoAlert
			if math.Abs(percentageChange) >= 20.0 {
				severity = WarningAlert
			}
			if math.Abs(percentageChange) >= 30.0 {
				severity = ErrorAlert
			}

			// Create alert message
			message := fmt.Sprintf("%s for field '%s' has changed by %.2f%% (threshold: %.2f%%)",
				metricType, field, percentageChange, threshold)

			// Add z-score information if available
			if len(previousValues) > 1 && math.Abs(metricComparison.ZScore) > 2.0 {
				message += fmt.Sprintf(" (z-score: %.2f)", metricComparison.ZScore)
			}

			// Create alert
			alert := TrendAlert{
				MetricType: MetricType(metricType),
				Severity:   severity,
				Message:    message,
			}

			result.Alerts = append(result.Alerts, alert)
		}
	}

	return result, nil
}

// contains checks if a slice contains a value
func contains(slice []MetricType, value MetricType) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
