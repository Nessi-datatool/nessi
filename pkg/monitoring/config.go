package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/security"
)

// loadConfig loads the monitoring configuration from a file
func loadConfig(configPath string) (*Config, error) {
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
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// defaultConfig returns the default monitoring configuration
func defaultConfig() *Config {
	return &Config{
		Alerts: struct {
			Enabled              bool
			Thresholds           map[string]AlertThreshold
			SilencePeriod        time.Duration
			CooldownPeriod       time.Duration
			CheckInterval        time.Duration
			NotificationChannels []string `json:"notification_channels"`
		}{
			Enabled:              true,
			Thresholds:           make(map[string]AlertThreshold),
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
				JWTSecret:    "change-me-in-production",
				UsersFile:    "./data/users.json",
				TokenExpiry:  24, // hours
				RequireHTTPS: false,
			},
			SSL: security.SSLConfig{
				Enabled:      false,
				CertFile:     "./certs/server.crt",
				KeyFile:      "./certs/server.key",
				AutoGenerate: true,
			},
		},
	}
}

// Notification config types are defined in notifications.go
