package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nessi-dev/nessi/pkg/security"
)

// MonitorConfig defines configuration for a monitor
type MonitorConfig struct {
	Enabled       bool          `json:"enabled"`
	Interval      time.Duration `json:"interval"`
	MetricsPrefix string        `json:"metrics_prefix"`
	OutputPath    string        `json:"output_path"`
}

// MonitoringConfig defines the monitoring configuration (renamed to avoid redeclaration)
type MonitoringConfig struct {
	Alerts struct {
		Enabled              bool
		Thresholds           map[string]AlertThresholdConfig
		SilencePeriod        time.Duration
		CooldownPeriod       time.Duration
		CheckInterval        time.Duration
		NotificationChannels []string `json:"notification_channels"`
	}
	Metrics struct {
		Enabled bool
		Port    int
	}
	Notifications struct {
		Slack   *SlackNotificationConfig
		Email   *EmailNotificationConfig
		Webhook *WebhookNotificationConfig
	}
	Retention struct {
		Enabled          bool          `json:"enabled"`
		StoragePath      string        `json:"storage_path"`
		RetentionPeriod  time.Duration `json:"retention_period"`
		SnapshotInterval time.Duration `json:"snapshot_interval"`
	}
	Security struct {
		Auth security.AuthConfig `json:"auth"`
		SSL  security.SSLConfig  `json:"ssl"`
	}
}

// AlertThresholdConfig defines thresholds for alerts (renamed to avoid redeclaration)
type AlertThresholdConfig struct {
	Warning  float64 `json:"warning"`
	Critical float64 `json:"critical"`
}

// loadConfig loads the monitoring configuration from a file
func loadConfig(configPath string) (*MonitoringConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return defaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	var config MonitoringConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// defaultConfig returns the default monitoring configuration
func defaultConfig() *MonitoringConfig {
	return &MonitoringConfig{
		Alerts: struct {
			Enabled              bool
			Thresholds           map[string]AlertThresholdConfig
			SilencePeriod        time.Duration
			CooldownPeriod       time.Duration
			CheckInterval        time.Duration
			NotificationChannels []string `json:"notification_channels"`
		}{
			Enabled:              true,
			Thresholds:           make(map[string]AlertThresholdConfig),
			SilencePeriod:        time.Hour * 24,
			CooldownPeriod:       time.Minute * 30,
			CheckInterval:        time.Minute * 5,
			NotificationChannels: []string{"email", "slack"},
		},
		Metrics: struct {
			Enabled bool
			Port    int
		}{
			Enabled: true,
			Port:    9090,
		},
		Notifications: struct {
			Slack   *SlackNotificationConfig
			Email   *EmailNotificationConfig
			Webhook *WebhookNotificationConfig
		}{
			Slack: &SlackNotificationConfig{
				WebhookURL: "https://hooks.slack.com/services/your-webhook-url",
				Username:   "Nessi Alert Bot",
				Channel:    "#alerts",
			},
			Email: &EmailNotificationConfig{
				SMTPServer: "smtp.example.com",
				Port:       587,
				Username:   "alerts@example.com",
				Password:   "password123",
				From:       "alerts@example.com",
				To:         []string{"admin@example.com"},
			},
			Webhook: &WebhookNotificationConfig{
				URL:     "https://example.com/webhook",
				Method:  "POST",
				Headers: map[string]string{"Content-Type": "application/json"},
			},
		},
		Retention: struct {
			Enabled          bool          `json:"enabled"`
			StoragePath      string        `json:"storage_path"`
			RetentionPeriod  time.Duration `json:"retention_period"`
			SnapshotInterval time.Duration `json:"snapshot_interval"`
		}{
			Enabled:          true,
			StoragePath:      "./data/metrics",
			RetentionPeriod:  time.Hour * 24 * 30, // 30 days
			SnapshotInterval: time.Hour * 6,       // Every 6 hours
		},
		Security: struct {
			Auth security.AuthConfig `json:"auth"`
			SSL  security.SSLConfig  `json:"ssl"`
		}{
			Auth: security.AuthConfig{
				Enabled:      false,
				UsersFile:    "./data/users.json",
				InMemoryOnly: false,
			},
			SSL: security.SSLConfig{
				Enabled:  false,
				CertFile: "./certs/server.crt",
				KeyFile:  "./certs/server.key",
			},
		},
	}
}

// Notification config types are defined in notifications.go
