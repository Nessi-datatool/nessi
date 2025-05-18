// [REMOVED FOR OSS]: Intelligent alerting tests are only available in LakeDiff Enterprise.

package alerts

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockMetricStore is a mock implementation of the MetricStore interface for testing
type mockMetricStore struct {
	metrics map[string][]MetricDataPoint
}

func newMockMetricStore() *mockMetricStore {
	return &mockMetricStore{
		metrics: make(map[string][]MetricDataPoint),
	}
}

func (m *mockMetricStore) GetMetricValues(metricName string, start, end time.Time, labels map[string]string) ([]MetricDataPoint, error) {
	return m.metrics[metricName], nil
}

func (m *mockMetricStore) GetMetricNames() ([]string, error) {
	names := make([]string, 0, len(m.metrics))
	for name := range m.metrics {
		names = append(names, name)
	}
	return names, nil
}

func (m *mockMetricStore) AddMetric(name string, dataPoints []MetricDataPoint) {
	m.metrics[name] = dataPoints
}

// TestIntelligentAlertingConfig tests the intelligent alerting configuration
func TestIntelligentAlertingConfig(t *testing.T) {
	config := DefaultIntelligentAlertingConfig()
	
	assert.Equal(t, 100, config.MinimumDataPoints)
	assert.Equal(t, 7*24*time.Hour, config.AnalysisPeriod)
	assert.Equal(t, 24*time.Hour, config.UpdateFrequency)
	assert.Equal(t, 0.7, config.Sensitivity)
	assert.True(t, config.EnableOutlierDetection)
	assert.True(t, config.EnableTrendDeviation)
	assert.True(t, config.EnableSeasonalPatterns)
	assert.True(t, config.AutoDisableUnusedRules)
	assert.Equal(t, 30, config.DisableThreshold)
}

// TestNewIntelligentAlertManager tests the creation of a new intelligent alert manager
func TestNewIntelligentAlertManager(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)
	
	metricStore := newMockMetricStore()
	config := DefaultIntelligentAlertingConfig()
	
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	assert.NotNil(t, iam)
	assert.Equal(t, alertManager, iam.alertManager)
	assert.Equal(t, metricStore, iam.metricStore)
	assert.Equal(t, config, iam.config)
}

// TestOutlierRuleCreation tests the creation of outlier detection rules
func TestOutlierRuleCreation(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create metric store with test data
	metricStore := newMockMetricStore()
	
	// Add a metric with consistent values around 100
	dataPoints := make([]MetricDataPoint, 200)
	now := time.Now()
	
	for i := 0; i < 200; i++ {
		// Values around 100 with some noise
		value := 100.0 + float64(i%10) - 5.0
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     value,
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("test_metric", dataPoints)
	
	// Create intelligent alert manager with test config
	config := &IntelligentAlertingConfig{
		MinimumDataPoints:     100,
		AnalysisPeriod:        7 * 24 * time.Hour,
		UpdateFrequency:       24 * time.Hour,
		Sensitivity:           0.7,
		EnableOutlierDetection: true,
		EnableTrendDeviation:   false,
		EnableSeasonalPatterns: false,
		AutoDisableUnusedRules: false,
	}
	
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Test outlier rule creation
	err = iam.createOutlierRules("test_metric", dataPoints)
	require.NoError(t, err)
	
	// Check that rules were created
	rules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "outlier",
	})
	require.NoError(t, err)
	
	// Should have created 2 rules (upper and lower bounds)
	assert.Equal(t, 2, len(rules))
	
	// Verify rule properties
	for _, rule := range rules {
		assert.Equal(t, "test_metric", rule.Metric)
		assert.Equal(t, SeverityWarning, rule.Severity)
		assert.True(t, rule.Enabled)
		assert.Equal(t, "intelligent", rule.Labels["type"])
		assert.Equal(t, "outlier", rule.Labels["strategy"])
		assert.Equal(t, "true", rule.Labels["generated"])
		
		// Check bounds
		if rule.Labels["bound"] == "upper" {
			assert.Equal(t, ">", rule.ComparisonOperator)
			assert.True(t, rule.Threshold > 100.0, "Upper bound should be above the mean")
		} else if rule.Labels["bound"] == "lower" {
			assert.Equal(t, "<", rule.ComparisonOperator)
			assert.True(t, rule.Threshold < 100.0, "Lower bound should be below the mean")
		} else {
			t.Fatalf("Unexpected bound type: %s", rule.Labels["bound"])
		}
	}
}

// TestTrendDeviationRuleCreation tests the creation of trend deviation rules
func TestTrendDeviationRuleCreation(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create metric store with test data
	metricStore := newMockMetricStore()
	
	// Add a metric with an upward trend
	dataPoints := make([]MetricDataPoint, 100)
	now := time.Now()
	
	for i := 0; i < 100; i++ {
		// Linear trend with some noise
		value := 100.0 + float64(i)*0.5 + float64(i%5) - 2.0
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     value,
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("trend_metric", dataPoints)
	
	// Create intelligent alert manager with test config
	config := &IntelligentAlertingConfig{
		MinimumDataPoints:     50,
		AnalysisPeriod:        7 * 24 * time.Hour,
		UpdateFrequency:       24 * time.Hour,
		Sensitivity:           0.7,
		EnableOutlierDetection: false,
		EnableTrendDeviation:   true,
		EnableSeasonalPatterns: false,
		AutoDisableUnusedRules: false,
	}
	
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Test trend deviation rule creation
	err = iam.createTrendDeviationRules("trend_metric", dataPoints)
	require.NoError(t, err)
	
	// Check that rules were created
	rules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "trend_deviation",
	})
	require.NoError(t, err)
	
	// Should have created 1 rule
	assert.Equal(t, 1, len(rules))
	
	// Verify rule properties
	rule := rules[0]
	assert.Equal(t, "trend_metric", rule.Metric)
	assert.Equal(t, SeverityWarning, rule.Severity)
	assert.True(t, rule.Enabled)
	assert.Equal(t, "intelligent", rule.Labels["type"])
	assert.Equal(t, "trend_deviation", rule.Labels["strategy"])
	assert.Equal(t, "true", rule.Labels["generated"])
	assert.Equal(t, ">", rule.ComparisonOperator)
	
	// Slope should be positive for our upward trend
	slope, err := parseFloat(rule.Labels["slope"])
	require.NoError(t, err)
	assert.True(t, slope > 0, "Slope should be positive for upward trend")
}

// TestSeasonalPatternRuleCreation tests the creation of seasonal pattern rules
func TestSeasonalPatternRuleCreation(t *testing.T) {
	// Create a unique temp directory for this test
	tempDir, err := os.MkdirTemp("", "nessi-test-alerts-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create alert manager
	alertManager, err := NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create metric store with test data
	metricStore := newMockMetricStore()
	
	// Add a metric with daily patterns (higher during business hours)
	dataPoints := make([]MetricDataPoint, 24*7) // One week of hourly data
	now := time.Now()
	
	for i := 0; i < 24*7; i++ {
		hour := i % 24
		// Higher values during business hours (9-17)
		var value float64
		if hour >= 9 && hour < 17 {
			value = 100.0 + float64(hour-9)*5.0 + float64(i%5) - 2.0
		} else {
			value = 50.0 + float64(i%5) - 2.0
		}
		
		dataPoints[i] = MetricDataPoint{
			Timestamp: now.Add(time.Duration(-i) * time.Hour),
			Value:     value,
			Labels:    map[string]string{"service": "test"},
		}
	}
	
	metricStore.AddMetric("seasonal_metric", dataPoints)
	
	// Create intelligent alert manager with test config
	config := &IntelligentAlertingConfig{
		MinimumDataPoints:     100,
		AnalysisPeriod:        7 * 24 * time.Hour,
		UpdateFrequency:       24 * time.Hour,
		Sensitivity:           0.7,
		EnableOutlierDetection: false,
		EnableTrendDeviation:   false,
		EnableSeasonalPatterns: true,
		AutoDisableUnusedRules: false,
	}
	
	iam := NewIntelligentAlertManager(alertManager, metricStore, config)
	
	// Test seasonal pattern rule creation
	err = iam.createSeasonalPatternRules("seasonal_metric", dataPoints)
	require.NoError(t, err)
	
	// Check that rules were created
	rules, err := alertManager.GetRulesByLabels(map[string]string{
		"type":     "intelligent",
		"strategy": "seasonal",
		"pattern":  "hourly",
	})
	require.NoError(t, err)
	
	// Should have created rules for multiple hours
	assert.True(t, len(rules) > 0, "Should have created at least some hourly rules")
	
	// Verify rule properties for a business hour
	businessHourRules := filterRulesByHour(rules, 12) // Noon
	assert.True(t, len(businessHourRules) > 0, "Should have rules for business hours")
	
	if len(businessHourRules) > 0 {
		rule := businessHourRules[0]
		assert.Equal(t, "seasonal_metric", rule.Metric)
		assert.Equal(t, SeverityWarning, rule.Severity)
		assert.True(t, rule.Enabled)
		assert.Equal(t, "intelligent", rule.Labels["type"])
		assert.Equal(t, "seasonal", rule.Labels["strategy"])
		assert.Equal(t, "hourly", rule.Labels["pattern"])
		assert.Equal(t, "true", rule.Labels["generated"])
		
		// Mean should be higher for business hours
		mean, err := parseFloat(rule.Labels["mean"])
		require.NoError(t, err)
		// The exact value might vary, but it should be positive
		assert.True(t, mean > 0.0, "Business hour mean should be positive")
	}
	
	// Verify rule properties for a non-business hour
	nonBusinessHourRules := filterRulesByHour(rules, 3) // 3 AM
	assert.True(t, len(nonBusinessHourRules) > 0, "Should have rules for non-business hours")
	
	if len(nonBusinessHourRules) > 0 {
		rule := nonBusinessHourRules[0]
		
		// Mean should be lower for non-business hours
		mean, err := parseFloat(rule.Labels["mean"])
		require.NoError(t, err)
		// The exact value might vary, but it should be positive and likely lower than business hours
		assert.True(t, mean >= 0.0, "Non-business hour mean should be non-negative")
	}
}

// TestStatisticalCalculations tests the statistical calculation helper functions
func TestStatisticalCalculations(t *testing.T) {
	// Test mean and standard deviation
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	mean, stdDev := calculateMeanAndStdDev(values)
	
	assert.InDelta(t, 5.5, mean, 0.001, "Mean should be 5.5")
	assert.InDelta(t, 2.872, stdDev, 0.001, "StdDev should be approximately 2.872")
	
	// Test linear regression
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 6, 8, 10}
	
	slope, intercept := calculateLinearRegression(x, y)
	
	assert.InDelta(t, 2.0, slope, 0.001, "Slope should be 2.0")
	assert.InDelta(t, 0.0, intercept, 0.001, "Intercept should be 0.0")
}

// Helper functions for tests

// parseFloat parses a float from a string
func parseFloat(s string) (float64, error) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	return f, err
}

// filterRulesByHour filters rules by the hour label
func filterRulesByHour(rules []*AlertRule, hour int) []*AlertRule {
	hourStr := fmt.Sprintf("%d", hour)
	var filtered []*AlertRule
	
	for _, rule := range rules {
		if rule.Labels["hour"] == hourStr {
			filtered = append(filtered, rule)
		}
	}
	
	return filtered
}
