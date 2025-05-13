package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockMetricStore is a mock implementation of the alerts.MetricStore interface for testing
type MockMetricStore struct {
	metrics map[string][]alerts.MetricDataPoint
}

// NewMockMetricStore creates a new MockMetricStore
func NewMockMetricStore() *MockMetricStore {
	return &MockMetricStore{
		metrics: make(map[string][]alerts.MetricDataPoint),
	}
}

// GetMetricValues retrieves historical values for a metric
func (m *MockMetricStore) GetMetricValues(metricName string, start, end time.Time, labels map[string]string) ([]alerts.MetricDataPoint, error) {
	// Get all data points for the metric
	allPoints, ok := m.metrics[metricName]
	if !ok {
		return []alerts.MetricDataPoint{}, nil
	}
	
	// Filter by time range and labels
	var filteredPoints []alerts.MetricDataPoint
	for _, point := range allPoints {
		// Check time range
		if (point.Timestamp.Equal(start) || point.Timestamp.After(start)) &&
		   (point.Timestamp.Equal(end) || point.Timestamp.Before(end)) {
			// Check labels
			if labels == nil || matchLabels(point.Labels, labels) {
				filteredPoints = append(filteredPoints, point)
			}
		}
	}
	
	return filteredPoints, nil
}

// GetMetricNames returns all available metric names
func (m *MockMetricStore) GetMetricNames() ([]string, error) {
	names := make([]string, 0, len(m.metrics))
	for name := range m.metrics {
		names = append(names, name)
	}
	return names, nil
}

// AddMetric adds a metric data point to the store
func (m *MockMetricStore) AddMetric(metricName string, dataPoint alerts.MetricDataPoint) {
	m.metrics[metricName] = append(m.metrics[metricName], dataPoint)
}

// matchLabels checks if a set of labels matches a filter
func matchLabels(labels, filter map[string]string) bool {
	for k, v := range filter {
		if labels[k] != v {
			return false
		}
	}
	return true
}

// TestMonitor is a simplified version of the Monitor struct for testing
type TestMonitor struct {
	metricsPort              int
	alertManager             *alerts.AlertManager
	intelligentAlertManager  *alerts.IntelligentAlertManager
	metricStore              alerts.MetricStore
}

// GetMetricsPort returns the metrics server port
func (m *TestMonitor) GetMetricsPort() int {
	return m.metricsPort
}

// GetAlertManager returns the alert manager
func (m *TestMonitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

// GetIntelligentAlertManager returns the intelligent alert manager
func (m *TestMonitor) GetIntelligentAlertManager() *alerts.IntelligentAlertManager {
	return m.intelligentAlertManager
}

// RecordMetric records a metric value with the current timestamp
func (m *TestMonitor) RecordMetric(metricName string, value float64, labels map[string]string) error {
	return m.RecordMetricWithTimestamp(metricName, value, time.Now(), labels)
}

// RecordMetricWithTimestamp records a metric value with a specific timestamp
func (m *TestMonitor) RecordMetricWithTimestamp(metricName string, value float64, timestamp time.Time, labels map[string]string) error {
	// Create a metric data point
	dataPoint := alerts.MetricDataPoint{
		Timestamp: timestamp,
		Value:     value,
		Labels:    labels,
	}
	
	// Store the data point in the metric store
	if m.metricStore != nil {
		if mockStore, ok := m.metricStore.(*MockMetricStore); ok {
			mockStore.AddMetric(metricName, dataPoint)
			return nil
		}
	}
	
	// If we're not using a mock metric store, just log the metric
	fmt.Printf("Recorded metric %s = %f at %s with labels %v\n", 
		metricName, value, timestamp.Format(time.RFC3339), labels)
	
	return nil
}

// TestIntelligentAlertingBasic tests the scheduling of intelligent alerting analysis
func TestIntelligentAlertingBasic(t *testing.T) {
	// Create a temporary directory for alert rules
	tempDir, err := os.MkdirTemp("", "intelligent_alerting_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create a mock metric store
	mockStore := NewMockMetricStore()
	
	// Create monitor with intelligent alerting enabled and a short update frequency
	intelligentConfig := &alerts.IntelligentAlertingConfig{
		MinimumDataPoints:      5, // Lower for testing
		AnalysisPeriod:         24 * time.Hour,
		UpdateFrequency:        1 * time.Second, // Short frequency for testing
		Sensitivity:            0.7,
		EnableOutlierDetection: true,
		EnableTrendDeviation:   false, // Disable trend detection as it requires more data points
		EnableSeasonalPatterns: false, // Disable seasonal patterns as it requires more data points
		AutoDisableUnusedRules: true,
		DisableThreshold:       15,
	}
	
	// Create the monitor
	monitor := &TestMonitor{
		metricsPort:  9090,
		alertManager: alertManager,
		metricStore:  mockStore,
	}
	
	// Initialize intelligent alerting
	monitor.intelligentAlertManager = alerts.NewIntelligentAlertManager(
		alertManager,
		mockStore,
		intelligentConfig,
	)
	
	// Add some test metrics
	err = monitor.RecordMetric("test_metric", 100.0, map[string]string{"service": "test"})
	require.NoError(t, err)
	
	// Add more data points to meet minimum requirements
	for i := 0; i < intelligentConfig.MinimumDataPoints; i++ {
		value := 100.0 + float64(i%5) // Small variations
		err = monitor.RecordMetricWithTimestamp(
			"test_metric",
			value,
			time.Now().Add(-time.Duration(i)*time.Hour),
			map[string]string{"service": "test"},
		)
		require.NoError(t, err)
	}
	
	// Trigger analysis manually
	err = monitor.intelligentAlertManager.AnalyzeMetricData("test_metric")
	require.NoError(t, err)
	
	// Verify that rules were created
	rules := alertManager.GetRules()
	assert.NotEmpty(t, rules, "Expected rules to be created")
}

// TestMetricStoreIntegration tests that metrics recorded via the Monitor are available for intelligent alerting
func TestMetricStoreIntegration(t *testing.T) {
	// Create a temporary directory for alert rules
	tempDir, err := os.MkdirTemp("", "intelligent_alerting_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create a mock metric store
	mockStore := NewMockMetricStore()
	
	// Create the monitor
	monitor := &TestMonitor{
		metricsPort:  9090,
		alertManager: alertManager,
		metricStore:  mockStore,
	}
	
	// Initialize intelligent alerting with default config
	config := alerts.DefaultIntelligentAlertingConfig()
	monitor.intelligentAlertManager = alerts.NewIntelligentAlertManager(
		alertManager,
		mockStore,
		config,
	)
	
	// Record multiple metrics
	metricName := "integration_test_metric"
	for i := 0; i < 20; i++ {
		value := 100.0 + float64(i%5) // Small variations
		timestamp := time.Now().Add(-time.Duration(i) * time.Hour)
		err := monitor.RecordMetricWithTimestamp(
			metricName,
			value,
			timestamp,
			map[string]string{"service": "test"},
		)
		require.NoError(t, err)
	}
	
	// Verify metrics are available in the store
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	metrics, err := mockStore.GetMetricValues(metricName, start, end, map[string]string{"service": "test"})
	require.NoError(t, err)
	assert.NotEmpty(t, metrics, "Expected metrics to be available in the store")
	assert.GreaterOrEqual(t, len(metrics), 20, "Expected at least 20 data points")
	
	// Verify that intelligent alerting can access these metrics
	err = monitor.intelligentAlertManager.AnalyzeMetricData(metricName)
	require.NoError(t, err)
	
	// Check if rules were created
	rules := alertManager.GetRules()
	assert.NotEmpty(t, rules, "Expected rules to be created based on recorded metrics")
}

// TestAlertTriggering tests that intelligent alerts can be triggered
func TestAlertTriggering(t *testing.T) {
	// Create a temporary directory for alert rules
	tempDir, err := os.MkdirTemp("", "intelligent_alerting_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Create a mock metric store
	mockStore := NewMockMetricStore()
	
	// Create the monitor
	monitor := &TestMonitor{
		metricsPort:  9090,
		alertManager: alertManager,
		metricStore:  mockStore,
	}
	
	// Initialize intelligent alerting with default config
	config := alerts.DefaultIntelligentAlertingConfig()
	monitor.intelligentAlertManager = alerts.NewIntelligentAlertManager(
		alertManager,
		mockStore,
		config,
	)
	
	// Record metrics with a clear pattern
	metricName := "trigger_test_metric"
	for i := 0; i < 50; i++ {
		value := 100.0 // Constant value
		timestamp := time.Now().Add(-time.Duration(i) * time.Hour)
		err := monitor.RecordMetricWithTimestamp(
			metricName,
			value,
			timestamp,
			map[string]string{"service": "test"},
		)
		require.NoError(t, err)
	}
	
	// Analyze the metric to create rules
	err = monitor.intelligentAlertManager.AnalyzeMetricData(metricName)
	require.NoError(t, err)
	
	// Get the rules that were created
	rules := alertManager.GetRules()
	assert.NotEmpty(t, rules, "Expected rules to be created")
	
	// Find an upper bound rule for our metric
	var upperRule *alerts.AlertRule
	for _, rule := range rules {
		if rule.Metric == metricName && rule.ComparisonOperator == ">" {
			upperRule = rule
			break
		}
	}
	
	assert.NotNil(t, upperRule, "Expected to find an upper bound rule")
	
	// Record a value that exceeds the threshold to trigger the alert
	err = monitor.RecordMetric(
		metricName,
		upperRule.Threshold + 10.0, // Ensure it's above the threshold
		map[string]string{"service": "test"},
	)
	require.NoError(t, err)
	
	// Check if an alert was created
	time.Sleep(1 * time.Second) // Give a little time for alert processing
	
	// Get active alerts
	activeAlerts := alertManager.GetAlerts()
	
	// Verify that our alert was triggered
	var found bool
	for _, alert := range activeAlerts {
		if alert.Source == metricName && alert.Status == "active" {
			found = true
			break
		}
	}
	
	assert.True(t, found, "Expected to find a firing alert for the rule")
}

// TestSensitivityLevels tests different sensitivity levels for intelligent alerting
func TestSensitivityLevels(t *testing.T) {
	// Create a temporary directory for alert rules
	tempDir, err := os.MkdirTemp("", "intelligent_alerting_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	alertManager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)
	
	// Test different sensitivity levels
	testCases := []struct {
		name        string
		sensitivity float64
		expected    string
	}{
		{"LowSensitivity", 0.1, "less sensitive"},
		{"MediumSensitivity", 0.5, "moderately sensitive"},
		{"HighSensitivity", 0.9, "highly sensitive"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock metric store
			mockStore := NewMockMetricStore()
			
			// Create intelligent alerting config with the test sensitivity
			intelligentConfig := &alerts.IntelligentAlertingConfig{
				MinimumDataPoints:      5, // Lower for testing
				AnalysisPeriod:         24 * time.Hour,
				UpdateFrequency:        time.Hour,
				Sensitivity:            tc.sensitivity,
				EnableOutlierDetection: true,
				EnableTrendDeviation:   false, // Disable trend detection as it requires more data points
				EnableSeasonalPatterns: false, // Disable seasonal patterns as it requires more data points
			}
			
			// Create the monitor
			monitor := &TestMonitor{
				metricsPort:  9090,
				alertManager: alertManager,
				metricStore:  mockStore,
			}
			
			// Initialize intelligent alerting
			monitor.intelligentAlertManager = alerts.NewIntelligentAlertManager(
				alertManager,
				mockStore,
				intelligentConfig,
			)
			
			// Verify the configuration was set correctly
			iam := monitor.GetIntelligentAlertManager()
			assert.NotNil(t, iam)
			assert.Equal(t, tc.sensitivity, iam.GetConfig().Sensitivity)
			
			// Add test data
			metricName := fmt.Sprintf("sensitivity_test_%s", tc.name)
			for i := 0; i < 100; i++ {
				value := 100.0 + float64(i%10) // Small variations
				timestamp := time.Now().Add(-time.Duration(i) * time.Hour)
				err := monitor.RecordMetricWithTimestamp(
					metricName,
					value,
					timestamp,
					map[string]string{"service": "test"},
				)
				require.NoError(t, err)
			}
			
			// Analyze the metric
			err := iam.AnalyzeMetricData(metricName)
			require.NoError(t, err)
			
			// Get the rules that were created
			rules := alertManager.GetRules()
			
			// Filter rules for our metric
			var metricRules []*alerts.AlertRule
			for _, rule := range rules {
				if rule.Metric == metricName {
					metricRules = append(metricRules, rule)
				}
			}
			
			// Verify that rules were created
			assert.NotEmpty(t, metricRules, "Expected rules to be created for the metric")
			
			// Log the thresholds for manual verification
			t.Logf("%s: Created %d rules with sensitivity %.1f", 
				tc.name, len(metricRules), tc.sensitivity)
			
			for i, rule := range metricRules {
				t.Logf("  Rule %d: %s %.2f", i+1, rule.ComparisonOperator, rule.Threshold)
			}
		})
	}
}
