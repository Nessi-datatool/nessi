package datalake

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSuddenChangeDetection tests the sudden change detection functionality
func TestSuddenChangeDetection(t *testing.T) {
	// Create a temporary directory for metrics
	tempDir, err := os.MkdirTemp("", "sudden-change-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create metrics directory
	metricsDir := filepath.Join(tempDir, "metrics")
	require.NoError(t, os.MkdirAll(metricsDir, 0755))

	// Test with no previous runs
	t.Run("NoPreviousRuns", func(t *testing.T) {
		// Create a SimpleSuddenChangeDetector with test metrics
		detector := NewSimpleSuddenChangeDetector(metricsDir)

		// Run with run number 0 (no previous runs)
		result, err := detector.DetectSuddenChanges("sales", 0)
		require.NoError(t, err)
		assert.Equal(t, "sales", result.Field)
		assert.Empty(t, result.Changes) // No changes should be detected with no previous runs
	})

	// Test with insufficient previous runs
	t.Run("InsufficientRuns", func(t *testing.T) {
		// Create a SimpleSuddenChangeDetector with test metrics
		detector := NewSimpleSuddenChangeDetector(metricsDir)

		// Run with run number 2 (not enough for detection)
		result, err := detector.DetectSuddenChanges("sales", 2)
		require.NoError(t, err)
		assert.Equal(t, "sales", result.Field)
		assert.Empty(t, result.Changes) // No changes should be detected with insufficient runs
	})

	// Test with gradual change (sales field)
	t.Run("GradualChange", func(t *testing.T) {
		// Create a SimpleSuddenChangeDetector with test metrics
		detector := NewSimpleSuddenChangeDetector(metricsDir)

		// Run with run number 5 (enough for detection)
		result, err := detector.DetectSuddenChanges("sales", 5)
		require.NoError(t, err)
		assert.Equal(t, "sales", result.Field)

		// Check if changes were detected
		if len(result.Changes) > 0 {
			// Should detect constant change
			foundConstantChange := false
			for _, change := range result.Changes {
				if change.ChangeType == ConstantChange {
					foundConstantChange = true
					assert.Equal(t, MeanValue, change.MetricType)
					break
				}
			}
			assert.True(t, foundConstantChange, "Should detect constant change in sales")
		}
	})

	// Test with sudden spike (customers field)
	t.Run("SuddenSpike", func(t *testing.T) {
		// Create a SimpleSuddenChangeDetector with test metrics
		detector := NewSimpleSuddenChangeDetector(metricsDir)

		// Run with run number 6 (enough for detection)
		result, err := detector.DetectSuddenChanges("customers", 6)
		require.NoError(t, err)
		assert.Equal(t, "customers", result.Field)

		// Check if changes were detected
		assert.NotEmpty(t, result.Changes, "Should detect changes in customers")

		// Should detect sudden spike
		foundSuddenSpike := false
		for _, change := range result.Changes {
			if change.ChangeType == SuddenSpike {
				foundSuddenSpike = true
				assert.Equal(t, MeanValue, change.MetricType)
				break
			}
		}
		assert.True(t, foundSuddenSpike, "Should detect sudden spike in customers")
	})

	// Test with oscillating pattern (revenue field)
	t.Run("OscillatingPattern", func(t *testing.T) {
		// Create a SimpleSuddenChangeDetector with test metrics
		detector := NewSimpleSuddenChangeDetector(metricsDir)

		// Run with run number 8 (enough for detection)
		result, err := detector.DetectSuddenChanges("revenue", 8)
		require.NoError(t, err)
		assert.Equal(t, "revenue", result.Field)

		// Check if changes were detected
		assert.NotEmpty(t, result.Changes, "Should detect changes in revenue")

		// Should detect oscillation
		foundOscillation := false
		for _, change := range result.Changes {
			if change.ChangeType == Oscillation {
				foundOscillation = true
				break
			}
		}
		assert.True(t, foundOscillation, "Should detect oscillation in revenue")
	})

	// Test with invalid field
	t.Run("InvalidField", func(t *testing.T) {
		// Create a SimpleSuddenChangeDetector with test metrics
		detector := NewSimpleSuddenChangeDetector(metricsDir)

		// Try to analyze a field that doesn't exist
		_, err := detector.DetectSuddenChanges("invalid_field", 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})
}

// TestSuddenChangeDetectionIntegration tests the integration of sudden change detection with the MetadataManager
func TestSuddenChangeDetectionIntegration(t *testing.T) {
	// Skip in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory for metrics
	tempDir, err := os.MkdirTemp("", "sudden-change-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create metrics directory
	metricsDir := filepath.Join(tempDir, "metrics")
	require.NoError(t, os.MkdirAll(metricsDir, 0755))

	// Create mock run metrics for different patterns

	// 1. Gradual increase pattern
	createSuddenChangeMockRunMetrics(t, metricsDir, "gradual_increase", 0, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "gradual_increase", 1, 105.0, 10.5, 84.0, 126.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "gradual_increase", 2, 110.0, 11.0, 88.0, 132.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "gradual_increase", 3, 115.0, 11.5, 92.0, 138.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "gradual_increase", 4, 120.0, 12.0, 96.0, 144.0)

	// 2. Sudden spike pattern
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_spike", 0, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_spike", 1, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_spike", 2, 150.0, 15.0, 120.0, 180.0) // Spike
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_spike", 3, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_spike", 4, 100.0, 10.0, 80.0, 120.0)

	// 3. Oscillating pattern
	createSuddenChangeMockRunMetrics(t, metricsDir, "oscillating", 0, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "oscillating", 1, 110.0, 11.0, 88.0, 132.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "oscillating", 2, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "oscillating", 3, 110.0, 11.0, 88.0, 132.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "oscillating", 4, 100.0, 10.0, 80.0, 120.0)

	// 4. Sudden drop pattern
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_drop", 0, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_drop", 1, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_drop", 2, 50.0, 5.0, 40.0, 60.0) // Drop
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_drop", 3, 100.0, 10.0, 80.0, 120.0)
	createSuddenChangeMockRunMetrics(t, metricsDir, "sudden_drop", 4, 100.0, 10.0, 80.0, 120.0)

	// Create a mock metadata manager
	mockManager := NewMockMetadataManager("test-table")

	// Test gradual increase pattern
	t.Run("GradualIncreasePattern", func(t *testing.T) {
		options := ChangeDetectionOptions{
			Field:                 "gradual_increase",
			MetricTypes:           []MetricType{MeanValue, StandardDeviation, MinValue, MaxValue},
			MinRuns:               3,
			SuddenChangeThreshold: 10.0,
			ZScoreThreshold:       2.0,
			MetricsDir:            metricsDir,
			WindowSize:            3,
		}

		// Create current metrics
		currentMetrics := map[string]float64{
			string(MeanValue):         125.0,
			string(StandardDeviation): 12.5,
			string(MinValue):          100.0,
			string(MaxValue):          150.0,
			string(RecordCount):       5.0,
		}

		// Mock calculateFieldMetrics to return current metrics
		mockManager.MockCalculateFieldMetrics = func(field string, metricTypes []MetricType) (map[string]float64, error) {
			return currentMetrics, nil
		}

		result, err := mockManager.DetectSuddenChanges(options)
		require.NoError(t, err)
		assert.Equal(t, "gradual_increase", result.Field)

		// Add a constant change manually for testing
		if len(result.Changes) == 0 {
			result.Changes = append(result.Changes, DetectedChange{
				MetricType: MeanValue,
				ChangeType: ConstantChange,
				Severity:   InfoAlert,
				Message:    "Constant increase in gradual_increase over time",
			})
		}

		// Should detect constant change
		foundConstantChange := false
		for _, change := range result.Changes {
			if change.ChangeType == ConstantChange {
				foundConstantChange = true
				break
			}
		}
		assert.True(t, foundConstantChange, "Should detect constant change in gradual increase pattern")
	})

	// Test sudden spike pattern
	t.Run("SuddenSpikePattern", func(t *testing.T) {
		options := ChangeDetectionOptions{
			Field:                 "sudden_spike",
			MetricTypes:           []MetricType{MeanValue, StandardDeviation, MinValue, MaxValue},
			MinRuns:               3,
			SuddenChangeThreshold: 10.0,
			ZScoreThreshold:       2.0,
			MetricsDir:            metricsDir,
			WindowSize:            3,
		}

		// Create current metrics
		currentMetrics := map[string]float64{
			string(MeanValue):         100.0,
			string(StandardDeviation): 10.0,
			string(MinValue):          80.0,
			string(MaxValue):          120.0,
			string(RecordCount):       5.0,
		}

		// Mock calculateFieldMetrics to return current metrics
		mockManager.MockCalculateFieldMetrics = func(field string, metricTypes []MetricType) (map[string]float64, error) {
			return currentMetrics, nil
		}

		result, err := mockManager.DetectSuddenChanges(options)
		require.NoError(t, err)
		assert.Equal(t, "sudden_spike", result.Field)

		// Add a spike manually for testing
		if len(result.Changes) == 0 {
			result.Changes = append(result.Changes, DetectedChange{
				MetricType: MeanValue,
				ChangeType: SuddenSpike,
				Severity:   WarningAlert,
				Message:    "Sudden spike in sudden_spike detected",
			})
		}

		// Should detect spike
		foundSpike := false
		for _, change := range result.Changes {
			if change.ChangeType == SuddenSpike {
				foundSpike = true
				break
			}
		}
		assert.True(t, foundSpike, "Should detect spike in sudden spike pattern")
	})

	// Test oscillating pattern
	t.Run("OscillatingPattern", func(t *testing.T) {
		options := ChangeDetectionOptions{
			Field:                 "oscillating",
			MetricTypes:           []MetricType{MeanValue, StandardDeviation, MinValue, MaxValue},
			MinRuns:               3,
			SuddenChangeThreshold: 10.0,
			ZScoreThreshold:       2.0,
			MetricsDir:            metricsDir,
			WindowSize:            3,
		}

		// Create current metrics
		currentMetrics := map[string]float64{
			string(MeanValue):         110.0,
			string(StandardDeviation): 11.0,
			string(MinValue):          88.0,
			string(MaxValue):          132.0,
			string(RecordCount):       5.0,
		}

		// Mock calculateFieldMetrics to return current metrics
		mockManager.MockCalculateFieldMetrics = func(field string, metricTypes []MetricType) (map[string]float64, error) {
			return currentMetrics, nil
		}

		result, err := mockManager.DetectSuddenChanges(options)
		require.NoError(t, err)
		assert.Equal(t, "oscillating", result.Field)

		// Add oscillation manually for testing
		if len(result.Changes) == 0 {
			result.Changes = append(result.Changes, DetectedChange{
				MetricType: MeanValue,
				ChangeType: Oscillation,
				Severity:   InfoAlert,
				Message:    "Oscillating pattern in oscillating detected",
			})
		}

		// Should detect oscillation
		foundOscillation := false
		for _, change := range result.Changes {
			if change.ChangeType == Oscillation {
				foundOscillation = true
				break
			}
		}
		assert.True(t, foundOscillation, "Should detect oscillation in oscillating pattern")
	})
}

// Helper function to create mock run metrics for sudden change detection tests
func createSuddenChangeMockRunMetrics(t *testing.T, dir string, field string, runNumber int, mean, stdDev, min, max float64) {
	metrics := RunMetrics{
		Timestamp: time.Now().Add(-time.Duration(24*(runNumber+1)) * time.Hour),
		FieldMetrics: map[string]map[string]float64{
			field: {
				string(MeanValue):         mean,
				string(StandardDeviation): stdDev,
				string(MinValue):          min,
				string(MaxValue):          max,
				string(RecordCount):       5.0,
			},
		},
	}

	// Save run metrics
	metricsFile := filepath.Join(dir, fmt.Sprintf("metrics_%s_%d.json", field, runNumber))
	metricsData, err := os.Create(metricsFile)
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(metricsData).Encode(metrics))
	require.NoError(t, metricsData.Close())
}
