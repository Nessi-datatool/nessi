package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Config command flags
	configKey   string
	configValue string
	configReset bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Nessi configuration",
	Long: `The config command allows you to view and modify Nessi configuration settings.
	
Examples:
  nessi config                           # View all configuration settings
  nessi config server.port               # View a specific configuration setting
  nessi config server.port --set 8080    # Set a configuration setting
  nessi config --reset                   # Reset configuration to defaults`,
	Run: func(cmd *cobra.Command, args []string) {
		if configReset {
			// Reset configuration to defaults
			if err := resetConfig(); err != nil {
				fmt.Printf("Error resetting configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Configuration reset to defaults")
			return
		}

		if len(args) == 0 {
			// No args, show all configuration
			showAllConfig()
			return
		}

		// Get the configuration key from args
		configKey = args[0]

		if configValue != "" {
			// Set the configuration value
			if err := setConfigValue(configKey, configValue); err != nil {
				fmt.Printf("Error setting configuration: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Set %s = %s\n", configKey, configValue)
			return
		}

		// Get the configuration value
		value := viper.Get(configKey)
		if value == nil {
			fmt.Printf("Configuration key not found: %s\n", configKey)
			os.Exit(1)
		}

		// Print the value
		fmt.Printf("%s = ", configKey)
		printValue(value)
		fmt.Println()
	},
}

// configExportCmd represents the config export command
var configExportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Export configuration to a file",
	Long: `Export the current configuration to a file in JSON or YAML format.
	
Examples:
  nessi config export config.json        # Export to JSON
  nessi config export config.yaml        # Export to YAML`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No output file specified")
			os.Exit(1)
		}

		outputFile := args[0]
		ext := strings.ToLower(filepath.Ext(outputFile))

		// Get all settings
		settings := viper.AllSettings()

		var data []byte
		var err error

		if ext == ".json" {
			// Export as JSON
			data, err = json.MarshalIndent(settings, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling configuration: %v\n", err)
				os.Exit(1)
			}
		} else if ext == ".yaml" || ext == ".yml" {
			// Export as YAML
			// Create a new viper instance to write YAML
			v := viper.New()
			for k, val := range settings {
				v.Set(k, val)
			}
			v.SetConfigType("yaml")

			// Write to a temporary file
			tempFile := filepath.Join(os.TempDir(), "nessi-config-temp.yaml")
			if err := v.WriteConfigAs(tempFile); err != nil {
				fmt.Printf("Error writing configuration: %v\n", err)
				os.Exit(1)
			}

			// Read the temporary file
			data, err = os.ReadFile(tempFile)
			if err != nil {
				fmt.Printf("Error reading temporary file: %v\n", err)
				os.Exit(1)
			}

			// Clean up
			os.Remove(tempFile)
		} else {
			fmt.Printf("Unsupported file format: %s\n", ext)
			os.Exit(1)
		}

		// Write to the output file
		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			fmt.Printf("Error writing configuration to file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration exported to %s\n", outputFile)
	},
}

// configImportCmd represents the config import command
var configImportCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import configuration from a file",
	Long: `Import configuration from a JSON or YAML file.
	
Examples:
  nessi config import config.json        # Import from JSON
  nessi config import config.yaml        # Import from YAML`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No input file specified")
			os.Exit(1)
		}

		inputFile := args[0]
		ext := strings.ToLower(filepath.Ext(inputFile))

		// Read the input file
		data, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %v\n", err)
			os.Exit(1)
		}

		// Create a new viper instance
		v := viper.New()

		if ext == ".json" {
			v.SetConfigType("json")
		} else if ext == ".yaml" || ext == ".yml" {
			v.SetConfigType("yaml")
		} else {
			fmt.Printf("Unsupported file format: %s\n", ext)
			os.Exit(1)
		}

		// Read the configuration
		if err := v.ReadConfig(strings.NewReader(string(data))); err != nil {
			fmt.Printf("Error reading configuration: %v\n", err)
			os.Exit(1)
		}

		// Get all settings from the imported config
		settings := v.AllSettings()

		// Update the main viper instance
		for k, val := range settings {
			viper.Set(k, val)
		}

		// Save the configuration
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error writing configuration: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration imported from %s\n", inputFile)
	},
}

// showAllConfig displays all configuration settings
func showAllConfig() {
	settings := viper.AllSettings()
	printSettings(settings, 0)
}

// printSettings recursively prints configuration settings
func printSettings(settings map[string]interface{}, indent int) {
	for k, v := range settings {
		// Print the key with indentation
		fmt.Printf("%s%s: ", strings.Repeat("  ", indent), k)

		// Check if the value is a nested map
		if nested, ok := v.(map[string]interface{}); ok {
			fmt.Println()
			printSettings(nested, indent+1)
		} else {
			printValue(v)
			fmt.Println()
		}
	}
}

// printValue prints a configuration value
func printValue(value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Printf("%q", v)
	case []interface{}:
		fmt.Print("[")
		for i, item := range v {
			if i > 0 {
				fmt.Print(", ")
			}
			printValue(item)
		}
		fmt.Print("]")
	default:
		fmt.Printf("%v", v)
	}
}

// setConfigValue sets a configuration value
func setConfigValue(key, value string) error {
	// Set the value in viper
	viper.Set(key, value)

	// Save the configuration
	return viper.WriteConfig()
}

// resetConfig resets the configuration to defaults
func resetConfig() error {
	// Get the home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create default configuration
	defaultConfig := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
			"tls":  false,
		},
		"extensions": map[string]interface{}{
			"enabled": []string{},
			"path":    filepath.Join(home, ".nessi", "extensions"),
		},
		"output": map[string]interface{}{
			"format": "text",
			"path":   ".",
		},
		"batch": map[string]interface{}{
			"parallel": 1,
			"output":   "batch_results",
		},
		"logging": map[string]interface{}{
			"level":  "info",
			"file":   filepath.Join(home, ".nessi", "logs", "nessi.log"),
			"stdout": true,
		},
	}

	// Reset viper
	viper.Reset()

	// Set default values
	for k, v := range defaultConfig {
		viper.Set(k, v)
	}

	// Save the configuration
	viper.SetConfigFile(filepath.Join(home, ".nessi.yaml"))
	return viper.WriteConfig()
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configExportCmd)
	configCmd.AddCommand(configImportCmd)

	// Add flags to the config command
	configCmd.Flags().StringVar(&configValue, "set", "", "Set the configuration value")
	configCmd.Flags().BoolVar(&configReset, "reset", false, "Reset configuration to defaults")
}
