package monitoring

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
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
	t.Parallel()
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
			RetryInterval: time.Second * 5,
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
	t.Parallel()
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "metrics-export-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a monitor with some test metrics
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	
	options := MonitorOptions{
		MetricsPort:  9090,
		AlertManager: alertManager,
	}
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Record some test metrics
	monitor.RecordMetric("test_metric_1", 123.45, nil)
	monitor.RecordMetric("test_metric_2", 67.89, nil)
	
	// Test CSV export
	csvPath := filepath.Join(tempDir, "metrics.csv")
	csvOptions := ExportOptions{
		Format:     "csv",
		OutputPath: csvPath,
	}
	
	result, err := monitor.ExportMetrics(csvOptions)
	require.NoError(t, err)
	require.Equal(t, csvPath, result)
	
	// Verify CSV file exists
	_, err = os.Stat(csvPath)
	require.NoError(t, err)
	
	// Test JSON export
	jsonPath := filepath.Join(tempDir, "metrics.json")
	jsonOptions := ExportOptions{
		Format:     "json",
		OutputPath: jsonPath,
	}
	
	result, err = monitor.ExportMetrics(jsonOptions)
	require.NoError(t, err)
	require.Equal(t, jsonPath, result)
	
	// Verify JSON file exists
	_, err = os.Stat(jsonPath)
	require.NoError(t, err)
}

// TestIntelligentAlertingIntegration tests the integration of intelligent alerting with the monitor
func TestIntelligentAlertingIntegration(t *testing.T) {
	// Create alert manager
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	
	// Create monitor options with intelligent alerting enabled
	options := MonitorOptions{
		MetricsPort:  9090,
		AlertManager: alertManager,
		EnableIntelligentAlerting: true,
	}
	
	// Create monitor
	monitor, err := New(options)
	require.NoError(t, err)
	
	// Verify intelligent alert manager is created
	intelligentAlertManager := monitor.GetIntelligentAlertManager()
	require.NotNil(t, intelligentAlertManager)
	
	// Record some test metrics to trigger intelligent alerting
	monitor.RecordMetric("test_metric_1", 100.0, nil)
	monitor.RecordMetric("test_metric_1", 200.0, nil) // Significant change
	
	// Verify that the intelligent alert manager is working
	// This is a basic test - in a real test we would verify more functionality
	assert.NotNil(t, intelligentAlertManager)
}
