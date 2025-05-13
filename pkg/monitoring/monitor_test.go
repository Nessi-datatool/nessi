package monitoring

import (
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMonitorCreation tests the creation of a new Monitor
func TestMonitorCreation(t *testing.T) {
	// Create alert manager
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	
	// Create monitor options
	options := MonitorOptions{
		MetricsPort:  9090,
		AlertManager: alertManager,
		EnableIntelligentAlerting: true,
	}
	
	// Create monitor
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Verify monitor properties
	assert.Equal(t, 9090, monitor.GetMetricsPort())
	assert.Equal(t, alertManager, monitor.GetAlertManager())
	assert.NotNil(t, monitor.GetIntelligentAlertManager())
	assert.NotNil(t, monitor.metricStore)
}

// TestMonitorWithoutIntelligentAlerting tests creating a monitor without intelligent alerting
func TestMonitorWithoutIntelligentAlerting(t *testing.T) {
	// Create alert manager
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	
	// Create monitor options
	options := MonitorOptions{
		MetricsPort:  9090,
		AlertManager: alertManager,
		EnableIntelligentAlerting: false,
	}
	
	// Create monitor
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Verify monitor properties
	assert.Equal(t, 9090, monitor.GetMetricsPort())
	assert.Equal(t, alertManager, monitor.GetAlertManager())
	assert.Nil(t, monitor.GetIntelligentAlertManager())
	assert.NotNil(t, monitor.metricStore)
}

// TestCreateAlertManager tests the creation of an alert manager with custom options
func TestCreateAlertManager(t *testing.T) {
	// Create alert manager with default options
	alertManager1, err := CreateAlertManager()
	require.NoError(t, err)
	assert.NotNil(t, alertManager1)
	
	// Create alert manager with custom options
	options := AlertManagerOptions{
		EnableEmail:   true,
		EnableSlack:   false,
		EnableWebhook: true,
		EmailConfig: &alerts.EmailConfig{
			Host:     "smtp.test.com",
			Port:     587,
			Username: "test@example.com",
			Password: "testpassword",
			From:     "test@example.com",
			FromName: "Test Alerts",
			UseHTML:  true,
		},
		WebhookConfig: &alerts.WebhookConfig{
			URL:        "https://test.com/webhook",
			Method:     "POST",
			Headers:    map[string]string{"Content-Type": "application/json"},
			MaxRetries: 3,
			RetryDelay: time.Second * 5,
		},
	}
	
	alertManager2, err := CreateAlertManager(options)
	require.NoError(t, err)
	assert.NotNil(t, alertManager2)
	
	// Verify notifiers
	notifiers := alertManager2.GetNotifiers()
	assert.Contains(t, notifiers, "email")
	assert.NotContains(t, notifiers, "slack")
	assert.Contains(t, notifiers, "webhook")
}

// TestExportMetrics tests the ExportMetrics function
func TestExportMetrics(t *testing.T) {
	// Create alert manager
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	
	// Create monitor
	options := MonitorOptions{
		MetricsPort:  9090,
		AlertManager: alertManager,
	}
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Test exporting metrics
	exportOptions := ExportOptions{
		Format:     ExportFormatCSV,
		OutputPath: "/tmp/metrics.csv",
		StartTime:  time.Now().Add(-1 * time.Hour),
		EndTime:    time.Now(),
		MetricName: "test_metric",
		Labels:     map[string]string{"service": "test"},
	}
	
	outputPath, err := monitor.ExportMetrics(exportOptions)
	require.NoError(t, err)
	assert.Equal(t, "/tmp/metrics.csv", outputPath)
}

// TestIntelligentAlertingIntegration tests the integration of intelligent alerting with the monitor
func TestIntelligentAlertingIntegration(t *testing.T) {
	// Create alert manager
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	
	// Create custom intelligent alerting config
	intelligentConfig := &alerts.IntelligentAlertingConfig{
		MinimumDataPoints:      50,
		AnalysisPeriod:         24 * time.Hour,
		UpdateFrequency:        time.Hour,
		Sensitivity:            0.8,
		EnableOutlierDetection: true,
		EnableTrendDeviation:   true,
		EnableSeasonalPatterns: false,
		AutoDisableUnusedRules: true,
		DisableThreshold:       15,
	}
	
	// Create monitor options
	options := MonitorOptions{
		MetricsPort:              9090,
		AlertManager:             alertManager,
		EnableIntelligentAlerting: true,
		IntelligentAlertingConfig: intelligentConfig,
	}
	
	// Create monitor
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Verify intelligent alert manager configuration
	iam := monitor.GetIntelligentAlertManager()
	assert.NotNil(t, iam)
	assert.Equal(t, intelligentConfig, iam.GetConfig())
	
	// Verify metric store integration
	assert.Equal(t, monitor.metricStore, iam.GetMetricStore())
	assert.Equal(t, alertManager, iam.GetAlertManager())
}
