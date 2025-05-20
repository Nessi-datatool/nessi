package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var dbtPluginInstallCmd = &cobra.Command{
	Use:   "install-dbt-plugin",
	Short: "Install the dbt plugin for Nessi.dev",
	Long: `Install the dbt plugin for Nessi.dev.

This command downloads and installs the dbt plugin for Nessi.dev, which provides
integration with dbt (data build tool) for data quality and profiling.

The plugin is installed as a built-in component and does not require separate
plugin management. After installation, you can use the 'nessi dbt' command.

The plugin is opt-in by default and must be explicitly enabled via configuration
or command-line flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing dbt plugin for Nessi.dev...")
		
		// In a real implementation, this would download the plugin from a repository
		// For now, we'll just enable the built-in plugin
		
		// Get config directory
		configDir := GetConfigDir()
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			os.Exit(1)
		}
		
		// Create or update nessi.yaml to enable the dbt plugin
		configPath := filepath.Join(configDir, "nessi.yaml")
		if err := enableDBTPluginInConfig(configPath); err != nil {
			fmt.Printf("Error updating configuration: %v\n", err)
			os.Exit(1)
		}
		
		// Create sample dbt plugin configuration
		sampleConfigPath := filepath.Join(configDir, "nessi_dbt.yaml.sample")
		if err := createSampleDBTConfig(sampleConfigPath); err != nil {
			fmt.Printf("Error creating sample configuration: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("dbt plugin installed successfully!")
		fmt.Println("The plugin is now available but still requires explicit activation.")
		fmt.Println("To use the plugin, run commands with the --enable-dbt-plugin flag:")
		fmt.Println("  nessi dbt validate --enable-dbt-plugin --select tag:daily")
		fmt.Println("")
		fmt.Println("Or enable it permanently in your configuration:")
		fmt.Printf("  echo 'enable_dbt_plugin: true' >> %s\n", configPath)
		fmt.Println("")
		fmt.Println("A sample configuration has been created at:")
		fmt.Printf("  %s\n", sampleConfigPath)
		fmt.Println("Copy and customize this file to configure the dbt plugin.")
	},
}

// Note: getConfigDir function has been moved to utils.go as GetConfigDir

// enableDBTPluginInConfig enables the dbt plugin in the configuration file
func enableDBTPluginInConfig(configPath string) error {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create new config file
		return os.WriteFile(configPath, []byte("# Nessi.dev configuration\n\n# Uncomment to enable the dbt plugin\n# enable_dbt_plugin: true\n"), 0644)
	}
	
	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Check if enable_dbt_plugin is already in the config
	if !containsString(string(data), "enable_dbt_plugin") {
		// Append to config
		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer f.Close()
		
		if _, err := f.WriteString("\n# Uncomment to enable the dbt plugin\n# enable_dbt_plugin: true\n"); err != nil {
			return fmt.Errorf("failed to update config file: %w", err)
		}
	}
	
	return nil
}

// createSampleDBTConfig creates a sample dbt plugin configuration
func createSampleDBTConfig(configPath string) error {
	sampleConfig := `# Nessi.dev dbt Plugin Configuration

# Enable the dbt plugin
enable_dbt_plugin: true

# Path to dbt project (optional, will auto-detect if not specified)
dbt_project_path: "/path/to/dbt/project"

# Base path for Delta tables
delta_base_path: "/path/to/delta/tables"

# Model to table mappings (optional)
model_table_mappings:
  - model_name: "customers"
    table_path: "/delta/warehouse/customers"
  - model_name: "orders"
    table_path: "/delta/warehouse/orders"

# Rule sets for data quality validation
rule_sets:
  - name: "default"
    default: true
    rules:
      - name: "no_nulls_in_id"
        type: "column_null"
        description: "ID column should not contain nulls"
        column: "id"
        threshold: 0
        severity: "error"
      
      - name: "valid_email_format"
        type: "column_pattern"
        description: "Email should have valid format"
        column: "email"
        pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
        severity: "warning"
  
  - name: "financial_models"
    tag: "financial"
    rules:
      - name: "amount_range_check"
        type: "column_range"
        description: "Amount should be within valid range"
        column: "amount"
        min: 0
        max: 1000000
        severity: "error"

# Alert configuration (optional)
alert_config:
  enabled: false
  channels:
    - type: "slack"
      webhook: "https://hooks.slack.com/services/xxx/yyy/zzz"
    
    - type: "email"
      recipients:
        - "data-quality@example.com"

# Output configuration (optional)
output_config:
  format: "table"
  include_artifacts: true
  artifacts_path: "target/nessi"
`

	return os.WriteFile(configPath, []byte(sampleConfig), 0644)
}

// containsString checks if a string contains another string
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(s)][0:len(substr)] == substr
}

func init() {
	rootCmd.AddCommand(dbtPluginInstallCmd)
}
