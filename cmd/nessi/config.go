package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// General configuration
	DataDir string

	// SMTP configuration for email notifications
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		From     string
	}

	// Slack configuration for Slack notifications
	Slack struct {
		WebhookURL string
		Channel    string
		Username   string
		IconEmoji  string
	}

	// Webhook configuration for webhook notifications
	Webhook struct {
		URL            string
		Method         string
		Headers        map[string]string
		TimeoutSeconds int
	}
}

// LoadConfig loads the application configuration
func LoadConfig() *Config {
	config := &Config{}

	// Set default values
	configDir := GetConfigDir()
	config.DataDir = filepath.Join(configDir, "data")

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(config.DataDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating data directory: %v\n", err)
		os.Exit(1)
	}

	// Set default SMTP configuration
	config.SMTP.Host = "smtp.example.com"
	config.SMTP.Port = 587
	config.SMTP.Username = "user@example.com"
	config.SMTP.Password = "password"
	config.SMTP.From = "nessi@example.com"

	// Set default Slack configuration
	config.Slack.WebhookURL = "https://hooks.slack.com/services/xxx/yyy/zzz"
	config.Slack.Channel = "#alerts"
	config.Slack.Username = "Nessi Alert Bot"
	config.Slack.IconEmoji = ":warning:"

	// Set default webhook configuration
	config.Webhook.URL = "https://example.com/webhook"
	config.Webhook.Method = "POST"
	config.Webhook.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	config.Webhook.TimeoutSeconds = 10

	// TODO: Load configuration from file
	// This would typically load from a YAML or JSON configuration file

	return config
}

// Global configuration instance
var appConfig = LoadConfig()
