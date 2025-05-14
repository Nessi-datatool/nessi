package monitoring

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
	"github.com/nessi-dev/nessi-dev/pkg/security"
)

// Monitor represents a monitoring system
type Monitor struct {
	mu                     sync.RWMutex
	metricsPort            int
	alertManager           *alerts.AlertManager
	intelligentAlertManager *alerts.IntelligentAlertManager
	metricStore            *PrometheusMetricStore
	authManager            *security.AuthManager
	certManager            *security.CertManager
	config                 *Config
	metricRetention        *MetricRetention
	alertThresholds        map[string]AlertThreshold
	silencePeriod          time.Duration
	cooldownPeriod         time.Duration
	slackConfig            *SlackNotificationConfig
	emailConfig            *EmailNotificationConfig
	webhookConfig          *WebhookNotificationConfig
	lastAlertTimes         map[string]time.Time
	stopCh                 chan struct{}
	running                bool
	alerts                 chan Alert
	configPath             string
	metrics                map[string]interface{}
}

// MonitorOptions represents monitor configuration options
type MonitorOptions struct {
	MetricsPort  int
	AlertManager *alerts.AlertManager
	EnableIntelligentAlerting bool
	IntelligentAlertingConfig *alerts.IntelligentAlertingConfig
}

// New creates a new Monitor instance
func New(options MonitorOptions) (*Monitor, error) {
	monitor := &Monitor{
		metricsPort:  options.MetricsPort,
		alertManager: options.AlertManager,
	}
	
	// Create metric store
	monitor.metricStore = NewPrometheusMetricStore(monitor)
	
	// Initialize intelligent alerting if enabled
	if options.EnableIntelligentAlerting {
		config := options.IntelligentAlertingConfig
		if config == nil {
			config = alerts.DefaultIntelligentAlertingConfig()
		}
		
		monitor.intelligentAlertManager = alerts.NewIntelligentAlertManager(
			monitor.alertManager,
			monitor.metricStore,
			config,
		)
		
		// Start background analysis
		go monitor.intelligentAlertManager.StartAnalysis()
	}
	
	return monitor, nil
}

// GetMetricsPort returns the metrics server port
func (m *Monitor) GetMetricsPort() int {
	return m.metricsPort
}

// GetAlertManager returns the alert manager
func (m *Monitor) GetAlertManager() *alerts.AlertManager {
	return m.alertManager
}

// GetIntelligentAlertManager returns the intelligent alert manager
func (m *Monitor) GetIntelligentAlertManager() *alerts.IntelligentAlertManager {
	return m.intelligentAlertManager
}

// Export functionality is defined in export.go

// AlertManagerOptions represents options for creating an AlertManager
type AlertManagerOptions struct {
	EnableEmail   bool
	EnableSlack   bool
	EnableWebhook bool
	EmailConfig   *alerts.EmailConfig
	SlackConfig   *alerts.SlackConfig
	WebhookConfig *alerts.WebhookConfig
}

// DefaultAlertManagerOptions returns default options for creating an AlertManager
func DefaultAlertManagerOptions() AlertManagerOptions {
	return AlertManagerOptions{
		EnableEmail:   true,
		EnableSlack:   true,
		EnableWebhook: true,
		EmailConfig: &alerts.EmailConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "alerts@example.com",
			Password: "password123",
			From:     "alerts@example.com",
			FromName: "Nessi Alerts",
			UseHTML:  true,
		},
		SlackConfig: &alerts.SlackConfig{
			WebhookURL: "https://hooks.slack.com/services/your-webhook-url",
			Username:   "Nessi Alert Bot",
			IconEmoji:  ":warning:",
		},
		WebhookConfig: &alerts.WebhookConfig{
			URL:        "https://example.com/webhook",
			Method:     "POST",
			Headers:    map[string]string{"Content-Type": "application/json"},
			MaxRetries: 3,
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

	// Create Slack notifier if enabled
	if opts.EnableSlack && opts.SlackConfig != nil {
		slackNotifier := alerts.NewSlackNotifier(*opts.SlackConfig)
		alertManager.RegisterNotifier(slackNotifier)
	}

	// Create webhook notifier if enabled
	if opts.EnableWebhook && opts.WebhookConfig != nil {
		webhookNotifier := alerts.NewWebhookNotifier(*opts.WebhookConfig)
		alertManager.RegisterNotifier(webhookNotifier)
	}

	return alertManager, nil
}
