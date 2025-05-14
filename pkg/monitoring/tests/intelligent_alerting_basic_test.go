package tests

import (
	"os"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntelligentAlertingBasic tests the basic functionality of intelligent alerting
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
	for i := 0; i < 20; i++ {
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
