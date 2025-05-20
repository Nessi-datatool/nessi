package monitoring

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi/pkg/security"
)

// Monitor represents a monitoring system
type Monitor struct {
	mu          sync.RWMutex
	metricsPort int

	metricStore     *MetricStore
	authManager     *security.AuthManager
	certManager     *security.CertManager
	config          *MonitoringConfig
	metricRetention *MetricRetention
	stopCh          chan struct{}
	running         bool
	configPath      string
	metrics         map[string]interface{}
}

// MonitorOptions represents monitor configuration options
type MonitorOptions struct {
	MetricsPort int
}

// New creates a new Monitor instance
func New(options MonitorOptions) (*Monitor, error) {
	monitor := &Monitor{
		metricsPort: options.MetricsPort,
	}

	// Create metric store
	monitor.metricStore = NewMetricStore()
	return monitor, nil
}

// GetMetricsPort returns the metrics server port
func (m *Monitor) GetMetricsPort() int {
	return m.metricsPort
}

// Export functionality is defined in export.go

// AlertManagerOptions represents options for creating an AlertManager
type AlertManagerOptions struct {
	EnableEmail bool
	EmailConfig *alerts.EmailConfig
	// Webhook support is limited in OSS version
	EnableWebhook bool
	WebhookConfig *alerts.WebhookConfig
}

// DefaultAlertManagerOptions returns default options for creating an AlertManager
func DefaultAlertManagerOptions() AlertManagerOptions {
	return AlertManagerOptions{
		EnableEmail:   true,
		EnableWebhook: false, // Disabled by default in OSS version
		EmailConfig: &alerts.EmailConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "alerts@example.com",
			Password: "password123",
			From:     "alerts@example.com",
			FromName: "Nessi Alerts",
			UseHTML:  true,
		},

		WebhookConfig: &alerts.WebhookConfig{
			URL:           "https://example.com/webhook",
			Method:        "POST",
			Headers:       map[string]string{"Content-Type": "application/json"},
			MaxRetries:    3,
			RetryInterval: time.Second * 5,
		},
	}
}

// CreateAlertManager creates a new AlertManager instance
func CreateAlertManager(options ...AlertManagerOptions) (*alerts.AlertManager, error) {
	// Use default options if none provided
	opts := DefaultAlertManagerOptions()
	if len(options) > 0 {
		opts = options[0]
	}

	// Create data directory for alerts if it doesn't exist
	dataDir := "./data/alerts"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create alert manager
	alertManager, err := alerts.NewAlertManager(dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create alert manager: %w", err)
	}

	// Create email notifier if enabled
	if opts.EnableEmail && opts.EmailConfig != nil {
		emailNotifier, err := alerts.NewEmailNotifier(*opts.EmailConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create email notifier: %w", err)
		}
		alertManager.RegisterNotifier(emailNotifier)
	}

	// Webhook support is limited in OSS version
	// Only enabled if explicitly configured
	if opts.EnableWebhook && opts.WebhookConfig != nil {
		webhookNotifier := alerts.NewWebhookNotifier(*opts.WebhookConfig)
		alertManager.RegisterNotifier(webhookNotifier)
	}

	return alertManager, nil
}
