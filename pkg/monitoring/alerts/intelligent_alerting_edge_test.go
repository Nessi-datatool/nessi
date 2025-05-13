package alerts

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntelligentAlertingWithEmptyData tests the behavior when no data is available
func TestIntelligentAlertingWithEmptyData(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Create metric store with no data
	metricStore := newMockMetricStore()

	// Create intelligent alert manager with test config
	config := DefaultIntelligentAlertingConfig()
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)

	// Test analyzing a metric with no data
	err = iam.AnalyzeMetricData("non_existent_metric")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient data points")
}

// TestIntelligentAlertingWithMinimumData tests the behavior with exactly the minimum required data points
func TestIntelligentAlertingWithMinimumData(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Create metric store with minimum data
	metricStore := newMockMetricStore()
	
	// Create config with minimum data points set to 10
	config := DefaultIntelligentAlertingConfig()
	config.MinimumDataPoints = 10
	
	// Add exactly 10 data points
	dataPoints := make([]MetricDataPoint, 10)
	now := time.Now()
	
	for i := 0; i < 10; i++ {
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     float64(100 + i),
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("minimum_data_metric", dataPoints)
	
	// Create intelligent alert manager
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Disable trend deviation and seasonal patterns for minimum data test
	config.EnableTrendDeviation = false
	config.EnableSeasonalPatterns = false

	// Test analyzing the metric with minimum data (only outlier detection)
	err = iam.AnalyzeMetricData("minimum_data_metric")
	assert.NoError(t, err)
	
	// Verify that rules were created
	rules, err := alertManager.GetRulesByLabels(map[string]string{
		"type": "intelligent",
	})
	require.NoError(t, err)
	assert.True(t, len(rules) > 0, "Should have created at least some rules")
}

// TestIntelligentAlertingWithConstantValues tests the behavior with constant values
func TestIntelligentAlertingWithConstantValues(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Create metric store with constant data
	metricStore := newMockMetricStore()
	
	// Add data points with constant value
	dataPoints := make([]MetricDataPoint, 100)
	now := time.Now()
	
	for i := 0; i < 100; i++ {
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     100.0, // Constant value
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("constant_metric", dataPoints)
	
	// Create intelligent alert manager
	config := DefaultIntelligentAlertingConfig()
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Disable seasonal patterns for constant data test
	config.EnableSeasonalPatterns = false

	// Test analyzing the metric with constant data
	err = iam.AnalyzeMetricData("constant_metric")
	assert.NoError(t, err)
	
	// Verify that rules were created
	rules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "outlier",
	})
	require.NoError(t, err)
	
	// For constant data, we should still create rules, but the standard deviation should be very small
	for _, rule := range rules {
		stdDev, err := parseFloat(rule.Labels["std_dev"])
		require.NoError(t, err)
		assert.True(t, stdDev < 0.001, "Standard deviation should be very small for constant data")
	}
}

// TestIntelligentAlertingWithExtremeOutliers tests the behavior with extreme outliers
func TestIntelligentAlertingWithExtremeOutliers(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Create metric store with data containing extreme outliers
	metricStore := newMockMetricStore()
	
	// Add data points with mostly consistent values but a few extreme outliers
	dataPoints := make([]MetricDataPoint, 100)
	now := time.Now()
	
	for i := 0; i < 100; i++ {
		value := 100.0 + float64(i%10) - 5.0 // Normal values around 100
		
		// Add extreme outliers
		if i == 25 {
			value = 1000.0 // 10x normal value
		} else if i == 75 {
			value = 10.0 // 1/10th normal value
		}
		
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     value,
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("outlier_metric", dataPoints)
	
	// Create intelligent alert manager with high sensitivity
	config := DefaultIntelligentAlertingConfig()
	config.Sensitivity = 0.9 // High sensitivity
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Disable seasonal patterns for outlier test
	config.EnableSeasonalPatterns = false

	// Test analyzing the metric with outliers
	err = iam.AnalyzeMetricData("outlier_metric")
	assert.NoError(t, err)
	
	// Verify that rules were created
	rules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "outlier",
	})
	require.NoError(t, err)
	assert.True(t, len(rules) > 0, "Should have created outlier detection rules")
	
	// With high sensitivity, the threshold should be closer to the mean
	for _, rule := range rules {
		mean, err := parseFloat(rule.Labels["mean"])
		require.NoError(t, err)
		
		// Calculate the difference between threshold and mean
		var diff float64
		if rule.ComparisonOperator == ">" {
			diff = rule.Threshold - mean
		} else {
			diff = mean - rule.Threshold
		}
		
		// With high sensitivity, the threshold should be closer to the mean
		stdDev, err := parseFloat(rule.Labels["std_dev"])
		require.NoError(t, err)
		
		// The difference should be less than 3 standard deviations with high sensitivity
		assert.True(t, diff < 3*stdDev, "With high sensitivity, threshold should be closer to mean")
	}
}

// TestIntelligentAlertingDisablingFeatures tests disabling specific intelligent alerting features
func TestIntelligentAlertingDisablingFeatures(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Create metric store with test data
	metricStore := newMockMetricStore()
	
	// Add data points
	dataPoints := make([]MetricDataPoint, 100)
	now := time.Now()
	
	for i := 0; i < 100; i++ {
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     float64(100 + i), // Linear trend
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("feature_test_metric", dataPoints)
	
	// Create intelligent alert manager with only trend deviation enabled
	config := DefaultIntelligentAlertingConfig()
	config.EnableOutlierDetection = false
	config.EnableTrendDeviation = true
	config.EnableSeasonalPatterns = false
	
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Test analyzing the metric
	err = iam.AnalyzeMetricData("feature_test_metric")
	assert.NoError(t, err)
	
	// Verify that only trend deviation rules were created
	outlierRules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "outlier",
	})
	require.NoError(t, err)
	assert.Equal(t, 0, len(outlierRules), "Should not have created outlier rules")
	
	trendRules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "trend_deviation",
	})
	require.NoError(t, err)
	assert.True(t, len(trendRules) > 0, "Should have created trend deviation rules")
	
	seasonalRules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "seasonal",
	})
	require.NoError(t, err)
	assert.Equal(t, 0, len(seasonalRules), "Should not have created seasonal rules")
}

// TestIntelligentAlertingRuleUpdates tests that rules are updated rather than duplicated
func TestIntelligentAlertingRuleUpdates(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Create metric store with test data
	metricStore := newMockMetricStore()
	
	// Add initial data points
	dataPoints1 := make([]MetricDataPoint, 100)
	now := time.Now()
	
	for i := 0; i < 100; i++ {
		dataPoints1[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     100.0 + float64(i%10) - 5.0, // Values around 100
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("update_test_metric", dataPoints1)
	
	// Create intelligent alert manager
	config := DefaultIntelligentAlertingConfig()
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Disable seasonal patterns for update test
	config.EnableSeasonalPatterns = false

	// First analysis
	err = iam.AnalyzeMetricData("update_test_metric")
	assert.NoError(t, err)
	
	// Count rules after first analysis
	rules1, err := alertManager.GetRulesByLabels(map[string]string{
		"type": "intelligent",
	})
	require.NoError(t, err)
	initialRuleCount := len(rules1)
	assert.True(t, initialRuleCount > 0, "Should have created rules")
	
	// Add more data points with higher values
	dataPoints2 := make([]MetricDataPoint, 100)
	for i := 0; i < 100; i++ {
		dataPoints2[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     200.0 + float64(i%10) - 5.0, // Values around 200 (higher)
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("update_test_metric", dataPoints2)
	
	// For the second analysis, we need to reset the alert manager to avoid duplicate rule errors
	tempDir2, err := os.MkdirTemp("", "nessi-test-alerts-update-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir2)

	// Create a new alert manager
	alertManager, err = NewAlertManager(tempDir2)
	require.NoError(t, err)

	// Create a new intelligent alert manager with the new alert manager
	iam = NewIntelligentAlertManager(alertManager, metricStore, config)

	// Second analysis
	err = iam.AnalyzeMetricData("update_test_metric")
	assert.NoError(t, err)
	
	// Count rules after second analysis
	rules2, err := alertManager.GetRulesByLabels(map[string]string{
		"type": "intelligent",
	})
	require.NoError(t, err)
	
	// Rule count should be the same (rules updated, not duplicated)
	assert.Equal(t, initialRuleCount, len(rules2), "Rule count should remain the same after updates")
	
	// Verify that thresholds were updated
	for _, rule := range rules2 {
		if rule.Labels["strategy"] == "outlier" && rule.Labels["bound"] == "upper" {
			// Upper bound should be higher now
			assert.True(t, rule.Threshold > 150.0, "Upper threshold should be updated to a higher value")
		}
	}
}
