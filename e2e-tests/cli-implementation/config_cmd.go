package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage Nessi configuration",
		Long:  `Commands for initializing and managing Nessi configuration.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	configCmd.AddCommand(newConfigInitCmd())
	configCmd.AddCommand(newConfigLoadCmd())

	return configCmd
}

func newConfigInitCmd() *cobra.Command {
	var configDir string

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Nessi configuration",
		Long:  `Initialize Nessi configuration in the specified directory.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create .nessi directory if it doesn't exist
			nessiDir := filepath.Join(configDir, ".nessi")
			if err := os.MkdirAll(nessiDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
				os.Exit(1)
			}

			// Create config.yaml file
			configFile := filepath.Join(nessiDir, "config.yaml")
			configContent := `# Nessi Configuration File
version: 1.0

# General settings
general:
  log_level: info
  data_dir: ./data
  temp_dir: /tmp/nessi

# Delta Lake settings
delta:
  default_path: ./delta_tables
  time_travel_enabled: true
  schema_validation: true

# Quality check settings
quality:
  completeness_threshold: 0.95
  accuracy_threshold: 0.98
  consistency_threshold: 0.90
  uniqueness_threshold: 1.0
  timeliness_threshold: 0.85

# Report settings
report:
  default_format: html
  output_dir: ./reports
  include_charts: true
  include_recommendations: true
`
			if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing config file: %s\n", err)
				os.Exit(1)
			}

			fmt.Println("Configuration initialized successfully")
		},
	}

	// Add flags
	initCmd.Flags().StringVar(&configDir, "dir", ".", "Directory to initialize configuration in")

	return initCmd
}

func newConfigLoadCmd() *cobra.Command {
	var configFile string

	loadCmd := &cobra.Command{
		Use:   "load",
		Short: "Load Nessi configuration",
		Long:  `Load Nessi configuration from the specified file.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if config file exists
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Config file %s does not exist\n", configFile)
				os.Exit(1)
			}

			fmt.Println("Configuration loaded successfully")
		},
	}

	// Add flags
	loadCmd.Flags().StringVar(&configFile, "file", ".nessi/config.yaml", "Configuration file to load")

	return loadCmd
}
