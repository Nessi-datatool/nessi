package monitoring_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	mockStore := monitoring.NewMockMetricStore()
	
	// Create monitor with intelligent alerting enabled and a short update frequency
	intelligentConfig := &alerts.IntelligentAlertingConfig{
		MinimumDataPoints:      10,
		AnalysisPeriod:         24 * time.Hour,
		UpdateFrequency:        1 * time.Second, // Short frequency for testing
		Sensitivity:            0.7,
		EnableOutlierDetection: true,
		EnableTrendDeviation:   true,
		EnableSeasonalPatterns: true,
		AutoDisableUnusedRules: true,
		DisableThreshold:       15,
	}
	
	// Create the monitor
	monitor := &monitoring.Monitor{
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
	err = monitor.intelligentAlertManager.AnalyzeMetric("test_metric", map[string]string{"service": "test"})
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
	mockStore := monitoring.NewMockMetricStore()
	
	// Create the monitor
	monitor := &monitoring.Monitor{
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
	err = monitor.intelligentAlertManager.AnalyzeMetric(metricName, map[string]string{"service": "test"})
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
	mockStore := monitoring.NewMockMetricStore()
	
	// Create the monitor
	monitor := &monitoring.Monitor{
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
	err = monitor.intelligentAlertManager.AnalyzeMetric(metricName, map[string]string{"service": "test"})
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
			mockStore := monitoring.NewMockMetricStore()
			
			// Create intelligent alerting config with the test sensitivity
			intelligentConfig := &alerts.IntelligentAlertingConfig{
				MinimumDataPoints:      50,
				AnalysisPeriod:         24 * time.Hour,
				UpdateFrequency:        time.Hour,
				Sensitivity:            tc.sensitivity,
				EnableOutlierDetection: true,
				EnableTrendDeviation:   true,
				EnableSeasonalPatterns: true,
			}
			
			// Create the monitor
			monitor := &monitoring.Monitor{
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
			err := iam.AnalyzeMetric(metricName, map[string]string{"service": "test"})
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
