package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Nessi configuration",
	Long:  `Manage Nessi configuration settings.`,
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a configuration file",
	Long:  `Validate a configuration file to ensure it has the correct format and values.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the file path from the flag
		filePath, _ := cmd.Flags().GetString("file")

		// Validate the file path
		validator := common.NewPathValidator()
		validator.RequireFile = true
		filePath, err := validator.ValidatePath(filePath)
		if err != nil {
			return err
		}

		// Check if this is a dry run
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		manager := common.NewDryRunManager(dryRun)

		// Add actions that would be performed
		manager.AddAction("Validate configuration file: %s", filePath)

		// Check if we should execute
		if manager.ShouldExecute() {
			// Read the file
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return common.NewError(common.ErrInvalidPath, fmt.Sprintf("Failed to read file: %s", err))
			}

			// Parse the YAML
			var config map[string]interface{}
			if err := yaml.Unmarshal(fileData, &config); err != nil {
				return common.NewError(common.ErrInvalidFormat, fmt.Sprintf("Invalid YAML format: %s", err))
			}

			// Validate the configuration
			err = validateConfig(config)
			if err != nil {
				return err
			}

			// Print success message
			success := color.New(color.FgGreen, color.Bold)
			success.Println("✅ Configuration file is valid!")
		}

		// Print actions if in dry run mode
		manager.PrintActions()

		return nil
	},
}

// validateConfig validates the configuration structure
func validateConfig(config map[string]interface{}) error {
	// Check error handling section
	if errorHandling, ok := config["error_handling"].(map[string]interface{}); ok {
		// Check interactive resolution
		if interactiveResolution, ok := errorHandling["interactive_resolution"]; ok {
			if _, ok := interactiveResolution.(bool); !ok {
				return common.NewError(common.ErrInvalidConfig, "error_handling.interactive_resolution must be a boolean")
			}
		}

		// Check telemetry
		if telemetry, ok := errorHandling["telemetry"].(map[string]interface{}); ok {
			// Check enabled
			if enabled, ok := telemetry["enabled"]; ok {
				if _, ok := enabled.(bool); !ok {
					return common.NewError(common.ErrInvalidConfig, "error_handling.telemetry.enabled must be a boolean")
				}
			}

			// Check max_errors
			if maxErrors, ok := telemetry["max_errors"]; ok {
				if _, ok := maxErrors.(int); !ok {
					return common.NewError(common.ErrInvalidConfig, "error_handling.telemetry.max_errors must be an integer")
				}
			}
		}

		// Check retry
		if retry, ok := errorHandling["retry"].(map[string]interface{}); ok {
			// Check enabled
			if enabled, ok := retry["enabled"]; ok {
				if _, ok := enabled.(bool); !ok {
					return common.NewError(common.ErrInvalidConfig, "error_handling.retry.enabled must be a boolean")
				}
			}

			// Check max_attempts
			if maxAttempts, ok := retry["max_attempts"]; ok {
				if _, ok := maxAttempts.(int); !ok {
					return common.NewError(common.ErrInvalidConfig, "error_handling.retry.max_attempts must be an integer")
				}
			}
		}
	}

	// Check logging section
	if logging, ok := config["logging"].(map[string]interface{}); ok {
		// Check level
		if level, ok := logging["level"].(string); ok {
			validLevels := map[string]bool{
				"debug": true,
				"info":  true,
				"warn":  true,
				"error": true,
			}
			if !validLevels[level] {
				return common.NewError(common.ErrInvalidConfig, "logging.level must be one of: debug, info, warn, error")
			}
		}

		// Check colored
		if colored, ok := logging["colored"]; ok {
			if _, ok := colored.(bool); !ok {
				return common.NewError(common.ErrInvalidConfig, "logging.colored must be a boolean")
			}
		}
	}

	return nil
}

func init() {
	// Add config command to root command
	CLI.RootCmd.AddCommand(configCmd)

	// Add validate command to config command
	configCmd.AddCommand(configValidateCmd)

	// Add flags to validate command
	configValidateCmd.Flags().String("file", "", "Path to the configuration file to validate")
	configValidateCmd.MarkFlagRequired("file")
}
