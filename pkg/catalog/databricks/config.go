package databricks

import (
	"fmt"
	"os"
)

// Config represents the configuration for Databricks
type Config struct {
	// BaseURL is the URL of the Databricks workspace
	BaseURL string `json:"base_url" yaml:"base_url"`

	// Token is the Databricks personal access token
	Token string `json:"token" yaml:"token"`

	// WorkspaceID is the ID of the Databricks workspace
	WorkspaceID string `json:"workspace_id" yaml:"workspace_id"`

	// DefaultSchema is the default schema to use when none is specified
	DefaultSchema string `json:"default_schema" yaml:"default_schema"`

	// DefaultCatalog is the default catalog to use when none is specified
	DefaultCatalog string `json:"default_catalog" yaml:"default_catalog"`
}

// NewConfigFromEnv creates a new Databricks configuration from environment variables
func NewConfigFromEnv() (*Config, error) {
	baseURL := os.Getenv("DATABRICKS_HOST")
	token := os.Getenv("DATABRICKS_TOKEN")
	workspaceID := os.Getenv("DATABRICKS_WORKSPACE_ID")
	defaultSchema := os.Getenv("DATABRICKS_DEFAULT_SCHEMA")
	defaultCatalog := os.Getenv("DATABRICKS_DEFAULT_CATALOG")

	if baseURL == "" {
		return nil, fmt.Errorf("DATABRICKS_HOST environment variable is not set")
	}

	if token == "" {
		return nil, fmt.Errorf("DATABRICKS_TOKEN environment variable is not set")
	}

	if workspaceID == "" {
		// Use a default workspace ID if not provided
		workspaceID = "0"
	}

	if defaultSchema == "" {
		// Use a default schema if not provided
		defaultSchema = "default"
	}

	if defaultCatalog == "" {
		// Use a default catalog if not provided
		defaultCatalog = "hive_metastore"
	}

	return &Config{
		BaseURL:        baseURL,
		Token:          token,
		WorkspaceID:    workspaceID,
		DefaultSchema:  defaultSchema,
		DefaultCatalog: defaultCatalog,
	}, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("Databricks base URL is required")
	}

	if c.Token == "" {
		return fmt.Errorf("Databricks token is required")
	}

	return nil
}
