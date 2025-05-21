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

// TestNewMonitor tests the creation of a new Monitor instance
func TestNewMonitor(t *testing.T) {
	// Test with default options
	options := MonitorOptions{
		MetricsPort: 8080,
	}
	monitor, err := New(options)
	require.NoError(t, err)
	assert.NotNil(t, monitor)
	assert.Equal(t, 8080, monitor.GetMetricsPort())
	assert.NotNil(t, monitor.metricStore)
}

// TestDefaultAlertManagerOptions tests the default alert manager options
func TestDefaultAlertManagerOptions(t *testing.T) {
	options := DefaultAlertManagerOptions()

	// Verify default email settings
	assert.True(t, options.EnableEmail)
	assert.NotNil(t, options.EmailConfig)
	assert.Equal(t, "smtp.example.com", options.EmailConfig.Host)
	assert.Equal(t, 587, options.EmailConfig.Port)
	assert.Equal(t, "alerts@example.com", options.EmailConfig.Username)
	assert.Equal(t, "password123", options.EmailConfig.Password)
	assert.Equal(t, "alerts@example.com", options.EmailConfig.From)
	assert.Equal(t, "Nessi Alerts", options.EmailConfig.FromName)
	assert.True(t, options.EmailConfig.UseHTML)

	// Verify default webhook settings
	assert.False(t, options.EnableWebhook) // Should be disabled by default in OSS version
	assert.NotNil(t, options.WebhookConfig)
	assert.Equal(t, "https://example.com/webhook", options.WebhookConfig.URL)
	assert.Equal(t, "POST", options.WebhookConfig.Method)
	assert.Equal(t, map[string]string{"Content-Type": "application/json"}, options.WebhookConfig.Headers)
	assert.Equal(t, 3, options.WebhookConfig.MaxRetries)
	assert.Equal(t, time.Second*5, options.WebhookConfig.RetryInterval)
}

// TestCreateAlertManager tests the creation of an alert manager
func TestCreateAlertManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "nessi-alerts-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Save current directory and restore it after the test
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(currentDir)

	// Change to the temporary directory for the test
	os.Chdir(tempDir)

	// Test with default options
	alertManager, err := CreateAlertManager()
	require.NoError(t, err)
	assert.NotNil(t, alertManager)

	// Verify that the data directory was created
	_, err = os.Stat(filepath.Join(tempDir, "data/alerts"))
	assert.NoError(t, err)

	// Test with custom options (email disabled, webhook enabled)
	customOptions := AlertManagerOptions{
		EnableEmail:   false,
		EnableWebhook: true,
		WebhookConfig: &alerts.WebhookConfig{
			URL:           "https://custom-webhook.example.com",
			Method:        "POST",
			Headers:       map[string]string{"Authorization": "Bearer token"},
			MaxRetries:    5,
			RetryInterval: time.Second * 10,
		},
	}

	alertManager, err = CreateAlertManager(customOptions)
	require.NoError(t, err)
	assert.NotNil(t, alertManager)
}

// TestMetricsOnly tests only the metrics functionality without alerts
func TestMetricsOnly(t *testing.T) {
	// Create a monitor with metrics enabled
	options := MonitorOptions{
		MetricsPort: 9090,
	}
	monitor, err := New(options)
	require.NoError(t, err)
	assert.NotNil(t, monitor)

	// Verify that the metrics port is set correctly
	assert.Equal(t, 9090, monitor.GetMetricsPort())

	// Verify that the metric store is initialized
	assert.NotNil(t, monitor.metricStore)
}
