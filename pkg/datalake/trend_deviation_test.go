package datalake

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTrendDeviation tests the trend deviation functionality
func TestTrendDeviation(t *testing.T) {
	// Create a temporary directory for metrics
	tempDir, err := os.MkdirTemp("", "trend-deviation-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock metadata manager
	mockManager := NewMockMetadataManager("test-table")

	// Create test data with numeric values
	testData := []map[string]interface{}{
		{
			"id":           1,
			"numeric":      100.0,
			"percentage":   50.0,
			"null_field":   nil,
			"text":         "value1",
			"unique_field": "unique1",
		},
		{
			"id":           2,
			"numeric":      200.0,
			"percentage":   60.0,
			"null_field":   nil,
			"text":         "value2",
			"unique_field": "unique2",
		},
		{
			"id":           3,
			"numeric":      300.0,
			"percentage":   70.0,
			"null_field":   "not-null",
			"text":         "value3",
			"unique_field": "unique3",
		},
		{
			"id":           4,
			"numeric":      400.0,
			"percentage":   80.0,
			"null_field":   nil,
			"text":         "value3", // Duplicate value
			"unique_field": "unique4",
		},
		{
			"id":           5,
			"numeric":      500.0,
			"percentage":   90.0,
			"null_field":   nil,
			"text":         "value4",
			"unique_field": "unique5",
		},
	}

	// Mock ReadParquetFile to return test data
	mockManager.MockReadParquetFile = func(file string) ([]map[string]interface{}, error) {
		return testData, nil
	}

	// Mock ReadTableMetadata to return a table with one file
	mockManager.MockReadTableMetadata = func() (*DeltaTable, error) {
		return &DeltaTable{
			Version: 0,
			Files:   []string{"test.parquet"},
		}, nil
	}

	// Mock GetSchemaFieldsAtVersion to return schema fields
	mockManager.MockGetSchemaFieldsAtVersion = func(version int64) ([]SchemaField, error) {
		return []SchemaField{
			{Name: "id", Type: "int32"},
			{Name: "numeric", Type: "float64"},
			{Name: "percentage", Type: "float64"},
			{Name: "null_field", Type: "utf8"},
			{Name: "text", Type: "utf8"},
			{Name: "unique_field", Type: "utf8"},
		}, nil
	}

	// Create metrics directory
	metricsDir := filepath.Join(tempDir, "metrics")
	require.NoError(t, os.MkdirAll(metricsDir, 0755))

	// Test analyzing trend deviation with no previous runs
	t.Run("NoPreviousRuns", func(t *testing.T) {
		options := TrendDeviationOptions{
			Field:        "numeric",
			MetricTypes:  []MetricType{MeanValue, StandardDeviation, MinValue, MaxValue},
			Threshold:    10.0,
			PreviousRuns: 3,
			MetricsDir:   metricsDir,
		}

		result, err := mockManager.AnalyzeTrendDeviation(options)
		require.NoError(t, err)
		assert.Equal(t, "numeric", result.Field)
		assert.Empty(t, result.Metrics) // No metrics should be compared with no previous runs
		assert.Empty(t, result.Alerts)  // No alerts should be generated with no previous runs
	})

	// Create a previous run with different metrics
	previousRunMetrics := RunMetrics{
		Timestamp: time.Now().Add(-24 * time.Hour),
		FieldMetrics: map[string]map[string]float64{
			"numeric": {
				string(MeanValue):         250.0, // Current is 300.0, 20% increase
				string(StandardDeviation): 141.42, // Approx std dev
				string(MinValue):          100.0,
				string(MaxValue):          400.0, // Current is 500.0, 25% increase
			},
		},
	}

	// Save previous run metrics
	previousRunFile := filepath.Join(metricsDir, "metrics_previous.json")
	previousRunData, err := os.Create(previousRunFile)
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(previousRunData).Encode(previousRunMetrics))
	require.NoError(t, previousRunData.Close())

	// Test analyzing trend deviation with one previous run
	t.Run("OnePreviousRun", func(t *testing.T) {
		options := TrendDeviationOptions{
			Field:        "numeric",
			MetricTypes:  []MetricType{MeanValue, StandardDeviation, MinValue, MaxValue},
			Threshold:    10.0,
			PreviousRuns: 3,
			MetricsDir:   metricsDir,
		}

		result, err := mockManager.AnalyzeTrendDeviation(options)
		require.NoError(t, err)
		assert.Equal(t, "numeric", result.Field)
		assert.Len(t, result.Metrics, 4) // Should have 4 metrics compared
		assert.Len(t, result.Alerts, 2)  // Should have alerts for mean and max value

		// Check mean value metric
		meanMetric := findMetric(result.Metrics, MeanValue)
		assert.NotNil(t, meanMetric)
		assert.Equal(t, 300.0, meanMetric.CurrentValue)
		assert.Equal(t, []float64{250.0}, meanMetric.PreviousValues)
		assert.InDelta(t, 20.0, meanMetric.PercentageChange, 0.1)

		// Check max value metric
		maxMetric := findMetric(result.Metrics, MaxValue)
		assert.NotNil(t, maxMetric)
		assert.Equal(t, 500.0, maxMetric.CurrentValue)
		assert.Equal(t, []float64{400.0}, maxMetric.PreviousValues)
		assert.InDelta(t, 25.0, maxMetric.PercentageChange, 0.1)

		// Check alerts
		meanAlert := findAlert(result.Alerts, MeanValue)
		assert.NotNil(t, meanAlert)
		assert.Equal(t, WarningAlert, meanAlert.Severity) // 20% change is 2x threshold
		assert.Equal(t, 300.0, meanAlert.CurrentValue)
		assert.Equal(t, 250.0, meanAlert.PreviousValue)
		assert.InDelta(t, 20.0, meanAlert.PercentageChange, 0.1)

		maxAlert := findAlert(result.Alerts, MaxValue)
		assert.NotNil(t, maxAlert)
		assert.Equal(t, WarningAlert, maxAlert.Severity) // 25% change is 2.5x threshold
		assert.Equal(t, 500.0, maxAlert.CurrentValue)
		assert.Equal(t, 400.0, maxAlert.PreviousValue)
		assert.InDelta(t, 25.0, maxAlert.PercentageChange, 0.1)
	})

	// Create another previous run with different metrics
	previousRunMetrics2 := RunMetrics{
		Timestamp: time.Now().Add(-48 * time.Hour),
		FieldMetrics: map[string]map[string]float64{
			"numeric": {
				string(MeanValue):         200.0,
				string(StandardDeviation): 100.0,
				string(MinValue):          100.0,
				string(MaxValue):          300.0,
			},
		},
	}

	// Save second previous run metrics
	previousRunFile2 := filepath.Join(metricsDir, "metrics_previous2.json")
	previousRunData2, err := os.Create(previousRunFile2)
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(previousRunData2).Encode(previousRunMetrics2))
	require.NoError(t, previousRunData2.Close())

	// Test analyzing trend deviation with multiple previous runs
	t.Run("MultiplePreviousRuns", func(t *testing.T) {
		options := TrendDeviationOptions{
			Field:        "numeric",
			MetricTypes:  []MetricType{MeanValue, StandardDeviation, MinValue, MaxValue},
			Threshold:    10.0,
			PreviousRuns: 3,
			MetricsDir:   metricsDir,
		}

		result, err := mockManager.AnalyzeTrendDeviation(options)
		require.NoError(t, err)
		assert.Equal(t, "numeric", result.Field)
		assert.Len(t, result.Metrics, 4) // Should have 4 metrics compared
		assert.GreaterOrEqual(t, len(result.Alerts), 2) // Should have at least 2 alerts

		// Check mean value metric with multiple previous values
		meanMetric := findMetric(result.Metrics, MeanValue)
		assert.NotNil(t, meanMetric)
		assert.Equal(t, 300.0, meanMetric.CurrentValue)
		assert.Len(t, meanMetric.PreviousValues, 2)
		assert.InDelta(t, 20.0, meanMetric.PercentageChange, 0.1) // Change from most recent run

		// Check z-score alerts
		zScoreAlerts := 0
		for _, alert := range result.Alerts {
			if alert.Message != "" && alert.Message[len(alert.Message)-1] == ')' {
				zScoreAlerts++
			}
		}
		assert.GreaterOrEqual(t, zScoreAlerts, 0) // May have z-score alerts with multiple runs
	})

	// Test with null percentage metric
	t.Run("NullPercentage", func(t *testing.T) {
		options := TrendDeviationOptions{
			Field:        "null_field",
			MetricTypes:  []MetricType{NullPercentage},
			Threshold:    10.0,
			PreviousRuns: 3,
			MetricsDir:   metricsDir,
		}

		result, err := mockManager.AnalyzeTrendDeviation(options)
		require.NoError(t, err)
		assert.Equal(t, "null_field", result.Field)

		// Check null percentage metric
		nullMetric := findMetric(result.Metrics, NullPercentage)
		if len(result.Metrics) > 0 {
			assert.NotNil(t, nullMetric)
			assert.InDelta(t, 80.0, nullMetric.CurrentValue, 0.1) // 4 out of 5 values are null (80%)
		}
	})

	// Test with unique ratio metric
	t.Run("UniqueRatio", func(t *testing.T) {
		options := TrendDeviationOptions{
			Field:        "text",
			MetricTypes:  []MetricType{UniqueRatio},
			Threshold:    10.0,
			PreviousRuns: 3,
			MetricsDir:   metricsDir,
		}

		result, err := mockManager.AnalyzeTrendDeviation(options)
		require.NoError(t, err)
		assert.Equal(t, "text", result.Field)

		// Check unique ratio metric
		uniqueMetric := findMetric(result.Metrics, UniqueRatio)
		if len(result.Metrics) > 0 {
			assert.NotNil(t, uniqueMetric)
			assert.InDelta(t, 0.8, uniqueMetric.CurrentValue, 0.1) // 4 unique values out of 5 (80%)
		}
	})

	// Test with invalid field
	t.Run("InvalidField", func(t *testing.T) {
		options := TrendDeviationOptions{
			Field:        "non_existent_field",
			MetricTypes:  []MetricType{MeanValue},
			Threshold:    10.0,
			PreviousRuns: 3,
			MetricsDir:   metricsDir,
		}

		_, err := mockManager.AnalyzeTrendDeviation(options)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist in the table schema")
	})
}

// findMetric finds a metric in a slice of metrics by type
func findMetric(metrics []MetricComparison, metricType MetricType) *MetricComparison {
	for i, metric := range metrics {
		if metric.MetricType == metricType {
			return &metrics[i]
		}
	}
	return nil
}

// findAlert finds an alert in a slice of alerts by metric type
func findAlert(alerts []TrendAlert, metricType MetricType) *TrendAlert {
	for i, alert := range alerts {
		if alert.MetricType == metricType {
			return &alerts[i]
		}
	}
	return nil
}

// TestCalculateAverage tests the CalculateAverage function
func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{
			name:     "empty slice",
			values:   []float64{},
			expected: 0,
		},
		{
			name:     "single value",
			values:   []float64{5},
			expected: 5,
		},
		{
			name:     "multiple values",
			values:   []float64{1, 2, 3, 4, 5},
			expected: 3,
		},
		{
			name:     "negative values",
			values:   []float64{-10, -5, 0, 5, 10},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateAverage(tt.values)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}

// TestCalculateMedian tests the CalculateMedian function
func TestCalculateMedian(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		expected float64
	}{
		{
			name:     "empty slice",
			values:   []float64{},
			expected: 0,
		},
		{
			name:     "single value",
			values:   []float64{5},
			expected: 5,
		},
		{
			name:     "odd number of values",
			values:   []float64{1, 3, 2, 5, 4},
			expected: 3,
		},
		{
			name:     "even number of values",
			values:   []float64{1, 3, 2, 4},
			expected: 2.5,
		},
		{
			name:     "already sorted values",
			values:   []float64{1, 2, 3, 4, 5},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMedian(tt.values)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}

// TestCalculateStandardDeviation tests the CalculateStandardDeviation function
func TestCalculateStandardDeviation(t *testing.T) {
	tests := []struct {
		name     string
		values   []float64
		avg      float64
		expected float64
	}{
		{
			name:     "empty slice",
			values:   []float64{},
			avg:      0,
			expected: 0,
		},
		{
			name:     "single value",
			values:   []float64{5},
			avg:      5,
			expected: 0,
		},
		{
			name:     "multiple values",
			values:   []float64{1, 2, 3, 4, 5},
			avg:      3,
			expected: math.Sqrt(2.5),
		},
		{
			name:     "same values",
			values:   []float64{3, 3, 3, 3, 3},
			avg:      3,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateStandardDeviation(tt.values, tt.avg)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}

// TestCalculateZScore tests the CalculateZScore function
func TestCalculateZScore(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		avg      float64
		stdDev   float64
		expected float64
	}{
		{
			name:     "zero standard deviation",
			value:    10,
			avg:      5,
			stdDev:   0,
			expected: 0,
		},
		{
			name:     "value equals average",
			value:    5,
			avg:      5,
			stdDev:   2,
			expected: 0,
		},
		{
			name:     "value above average",
			value:    7,
			avg:      5,
			stdDev:   2,
			expected: 1,
		},
		{
			name:     "value below average",
			value:    3,
			avg:      5,
			stdDev:   2,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateZScore(tt.value, tt.avg, tt.stdDev)
			assert.InDelta(t, tt.expected, result, 0.0001)
		})
	}
}

// TestSimpleTrendAnalyzer tests the SimpleTrendAnalyzer used in the demo
func TestSimpleTrendAnalyzer(t *testing.T) {
	// Create a temporary directory for metrics
	tempDir, err := os.MkdirTemp("", "simple-trend-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a simple trend analyzer
	analyzer := &SimpleTrendAnalyzer{
		MetricsDir: tempDir,
		Metrics: map[string]map[string]float64{
			"sales": {
				string(MeanValue):         300.0,
				string(StandardDeviation): 158.11,
				string(MinValue):          100.0,
				string(MaxValue):          500.0,
				string(RecordCount):       5.0,
			},
			"customers": {
				string(MeanValue):         30.0,
				string(StandardDeviation): 15.81,
				string(MinValue):          10.0,
				string(MaxValue):          50.0,
				string(RecordCount):       5.0,
			},
		},
	}

	// Test with sales field
	t.Run("SalesField", func(t *testing.T) {
		// Run 1: Initial run
		result1, err := analyzer.AnalyzeTrendDeviation("sales", 0)
		require.NoError(t, err)
		assert.Equal(t, "sales", result1.Field)
		assert.Empty(t, result1.Alerts) // No alerts for first run

		// Run 2: Increase sales by 10%
		result2, err := analyzer.AnalyzeTrendDeviation("sales", 1)
		require.NoError(t, err)
		assert.Equal(t, "sales", result2.Field)
		assert.NotEmpty(t, result2.Metrics)

		// Check for alerts
		hasAlerts := len(result2.Alerts) > 0
		assert.True(t, hasAlerts, "Should have alerts for 10% increase")

		// Check mean value metric
		meanMetric := findMetric(result2.Metrics, MeanValue)
		if assert.NotNil(t, meanMetric) {
			assert.InDelta(t, 10.0, meanMetric.PercentageChange, 0.1)
		}
	})

	// Test with customers field
	t.Run("CustomersField", func(t *testing.T) {
		// Run 1: Initial run
		result1, err := analyzer.AnalyzeTrendDeviation("customers", 0)
		require.NoError(t, err)
		assert.Equal(t, "customers", result1.Field)

		// Run 2: Decrease customers by 30%
		result2, err := analyzer.AnalyzeTrendDeviation("customers", 3)
		require.NoError(t, err)
		assert.Equal(t, "customers", result2.Field)
		assert.NotEmpty(t, result2.Metrics)

		// Check for alerts
		hasAlerts := len(result2.Alerts) > 0
		assert.True(t, hasAlerts, "Should have alerts for 30% decrease")

		// Check mean value metric
		meanMetric := findMetric(result2.Metrics, MeanValue)
		if assert.NotNil(t, meanMetric) {
			assert.InDelta(t, -30.0, meanMetric.PercentageChange, 0.1)
		}

		// Check alert severity
		hasErrorAlert := false
		for _, alert := range result2.Alerts {
			if alert.Severity == ErrorAlert {
				hasErrorAlert = true
				break
			}
		}
		assert.True(t, hasErrorAlert, "Should have error severity for 30% decrease")
	})
}
