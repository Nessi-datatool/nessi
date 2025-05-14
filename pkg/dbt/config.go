package dbt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the configuration for the dbt plugin
type Config struct {
	// EnableDBTPlugin indicates whether the dbt plugin is enabled
	EnableDBTPlugin bool `json:"enable_dbt_plugin" yaml:"enable_dbt_plugin"`
	
	// DBTProjectPath is the path to the dbt project
	DBTProjectPath string `json:"dbt_project_path" yaml:"dbt_project_path"`
	
	// DeltaBasePath is the base path for Delta tables
	DeltaBasePath string `json:"delta_base_path" yaml:"delta_base_path"`
	
	// ModelTableMappings maps dbt models to Delta tables
	ModelTableMappings []ModelTableMapping `json:"model_table_mappings" yaml:"model_table_mappings"`
	
	// RuleSets defines sets of data quality rules
	RuleSets []RuleSet `json:"rule_sets" yaml:"rule_sets"`
	
	// AlertConfig defines alerting configuration
	AlertConfig AlertConfig `json:"alert_config" yaml:"alert_config"`
	
	// AlertSettings defines alert settings (new format)
	AlertSettings AlertSettings `json:"alert_settings" yaml:"alert_settings"`
	
	// OutputConfig defines output configuration
	OutputConfig OutputConfig `json:"output_config" yaml:"output_config"`
}

// ModelTableMapping maps a dbt model to a Delta table
type ModelTableMapping struct {
	// ModelName is the name of the dbt model
	ModelName string `json:"model_name" yaml:"model_name"`
	
	// TablePath is the path to the Delta table
	TablePath string `json:"table_path" yaml:"table_path"`
}

// RuleSet defines a set of data quality rules
type RuleSet struct {
	// Name is the name of the rule set
	Name string `json:"name" yaml:"name"`
	
	// ModelName is the name of the dbt model to apply the rules to
	// If empty, the rules are applied to all models
	ModelName string `json:"model_name" yaml:"model_name"`
	
	// Tag is the dbt tag to apply the rules to
	// If empty, the rules are applied to all models
	Tag string `json:"tag" yaml:"tag"`
	
	// Default indicates whether this is the default rule set
	Default bool `json:"default" yaml:"default"`
	
	// Rules is the list of data quality rules
	Rules []*Rule `json:"rules" yaml:"rules"`
}

// Rule defines a data quality rule
type Rule struct {
	// Name is the name of the rule
	Name string `json:"name" yaml:"name"`
	
	// Type is the type of rule (e.g., "sql", "column_null", "column_unique")
	Type string `json:"type" yaml:"type"`
	
	// Description is a description of the rule
	Description string `json:"description" yaml:"description"`
	
	// Column is the column to apply the rule to (for column-specific rules)
	Column string `json:"column" yaml:"column"`
	
	// SQL is the SQL query to execute (for SQL rules)
	SQL string `json:"sql" yaml:"sql"`
	
	// Threshold is the threshold for the rule (e.g., max null percentage)
	Threshold float64 `json:"threshold" yaml:"threshold"`
	
	// Severity is the severity of the rule (e.g., "error", "warning")
	Severity string `json:"severity" yaml:"severity"`
}

// AlertConfig defines alerting configuration
type AlertConfig struct {
	// Enabled indicates whether alerting is enabled
	Enabled bool `json:"enabled" yaml:"enabled"`
	
	// Channels defines alerting channels
	Channels []AlertChannel `json:"channels" yaml:"channels"`
}

// AlertChannel defines an alerting channel
type AlertChannel struct {
	// Type is the type of channel (e.g., "slack", "email")
	Type string `json:"type" yaml:"type"`
	
	// Webhook is the webhook URL (for Slack)
	Webhook string `json:"webhook" yaml:"webhook"`
	
	// Recipients is the list of recipients (for email)
	Recipients []string `json:"recipients" yaml:"recipients"`
}

// AlertSettings defines the settings for the alerting system
type AlertSettings struct {
	// Enabled indicates whether alerting is enabled
	Enabled bool `json:"enabled" yaml:"enabled"`
	
	// ThresholdScore is the quality score threshold for sending alerts
	ThresholdScore int `json:"threshold_score" yaml:"threshold_score"`
	
	// Slack defines Slack-specific alert settings
	Slack SlackAlertSettings `json:"slack" yaml:"slack"`
	
	// Email defines email-specific alert settings
	Email EmailAlertSettings `json:"email" yaml:"email"`
}

// SlackAlertSettings defines Slack-specific alert settings
type SlackAlertSettings struct {
	// Enabled indicates whether Slack alerting is enabled
	Enabled bool `json:"enabled" yaml:"enabled"`
	
	// WebhookURL is the Slack webhook URL
	WebhookURL string `json:"webhook_url" yaml:"webhook_url"`
	
	// Channel is the Slack channel to send alerts to
	Channel string `json:"channel" yaml:"channel"`
	
	// Username is the username to use for Slack alerts
	Username string `json:"username" yaml:"username"`
}

// EmailAlertSettings defines email-specific alert settings
type EmailAlertSettings struct {
	// Enabled indicates whether email alerting is enabled
	Enabled bool `json:"enabled" yaml:"enabled"`
	
	// SMTPHost is the SMTP host for sending emails
	SMTPHost string `json:"smtp_host" yaml:"smtp_host"`
	
	// SMTPPort is the SMTP port for sending emails
	SMTPPort int `json:"smtp_port" yaml:"smtp_port"`
	
	// SMTPUser is the SMTP username for authentication
	SMTPUser string `json:"smtp_user" yaml:"smtp_user"`
	
	// SMTPPassword is the SMTP password for authentication
	SMTPPassword string `json:"smtp_password" yaml:"smtp_password"`
	
	// From is the sender email address
	From string `json:"from" yaml:"from"`
	
	// To is the list of recipient email addresses
	To []string `json:"to" yaml:"to"`
	
	// SubjectPrefix is the prefix to add to email subjects
	SubjectPrefix string `json:"subject_prefix" yaml:"subject_prefix"`
}

// OutputConfig defines output configuration
type OutputConfig struct {
	// Format is the output format (e.g., "table", "json", "csv")
	Format string `json:"format" yaml:"format"`
	
	// Path is the output path (for file output)
	Path string `json:"path" yaml:"path"`
	
	// IncludeArtifacts indicates whether to generate dbt artifacts
	IncludeArtifacts bool `json:"include_artifacts" yaml:"include_artifacts"`
	
	// ArtifactsPath is the path for dbt artifacts
	ArtifactsPath string `json:"artifacts_path" yaml:"artifacts_path"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	// If no config path provided, look for default config files
	if configPath == "" {
		// Check for nessi_dbt.yaml in current directory
		if _, err := os.Stat("nessi_dbt.yaml"); err == nil {
			configPath = "nessi_dbt.yaml"
		} else if _, err := os.Stat("nessi_dbt.yml"); err == nil {
			configPath = "nessi_dbt.yml"
		} else if _, err := os.Stat("nessi.yaml"); err == nil {
			configPath = "nessi.yaml"
		} else if _, err := os.Stat("nessi.yml"); err == nil {
			configPath = "nessi.yml"
		} else {
			// Use default configuration
			return &Config{
				EnableDBTPlugin: true,
				DeltaBasePath: "/delta",
				RuleSets: []RuleSet{
					{
						Name:    "default",
						Default: true,
						Rules:   []*Rule{},
					},
				},
				AlertConfig: AlertConfig{
					Enabled: false,
				},
				AlertSettings: AlertSettings{
					Enabled:       false,
					ThresholdScore: 80,
					Slack: SlackAlertSettings{
						Enabled: false,
					},
					Email: EmailAlertSettings{
						Enabled: false,
					},
				},
				OutputConfig: OutputConfig{
					Format:           "table",
					IncludeArtifacts: false,
				},
			}, nil
		}
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config file based on extension
	var config Config
	ext := filepath.Ext(configPath)
	if ext == ".json" {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	}

	return &config, nil
}

// DBTManifest represents the dbt manifest.json file
type DBTManifest struct {
	Nodes    map[string]*DBTModel `json:"nodes"`
	Metadata DBTMetadata          `json:"metadata"`
}

// DBTModel represents a dbt model
type DBTModel struct {
	Name         string   `json:"name"`
	Schema       string   `json:"schema"`
	Database     string   `json:"database"`
	ResourceType string   `json:"resource_type"`
	Tags         []string `json:"tags"`
	Config       DBTConfig `json:"config"`
}

// DBTConfig represents dbt model configuration
type DBTConfig struct {
	Materialized string `json:"materialized"`
	Schema       string `json:"schema"`
	Database     string `json:"database"`
}

// DBTMetadata represents dbt metadata
type DBTMetadata struct {
	ProjectName string `json:"project_name"`
	Version     string `json:"dbt_version"`
}

// DBTRunResults represents the dbt run_results.json file
type DBTRunResults struct {
	Results []DBTTestResult `json:"results"`
}

// DBTTestResult represents a dbt test result
type DBTTestResult struct {
	Status       string `json:"status"`
	ModelName    string `json:"unique_id"`
	TestName     string `json:"name"`
	FailureCount int    `json:"failures"`
	Message      string `json:"message"`
}
