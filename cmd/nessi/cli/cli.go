package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/nessi-dev/nessi-dev/pkg/quality/engine"
	"github.com/nessi-dev/nessi-dev/pkg/security"
)

// CLI represents the command-line interface
var CLI struct {
	RootCmd      *cobra.Command
	Config       *viper.Viper
	Engine       engine.QualityEngine
	Security     *security.SecurityManager
	Datalake     *datalake.Reader
	ConfigPath   string
	Token        string
	TablePath    string
	OutputFormat string
}

// Init initializes the CLI
func Init() {
	CLI.Config = viper.New()
	CLI.RootCmd = &cobra.Command{
		Use:   "nessi",
		Short: "Nessi - Delta Lake Data Quality Tool",
		Long: `Nessi is a powerful tool for managing data quality in Delta Lake tables.
It provides features for validation, profiling, monitoring, and reporting.`,
	}

	// Add global flags
	CLI.RootCmd.PersistentFlags().StringVar(&CLI.ConfigPath, "config", "", "Path to configuration file")
	CLI.RootCmd.PersistentFlags().StringVar(&CLI.Token, "token", "", "Authentication token")
	CLI.RootCmd.PersistentFlags().StringVar(&CLI.OutputFormat, "output", "json", "Output format (json|yaml|text)")

	// Initialize subcommands
	initValidateCmd()
	initProfileCmd()
	initMonitorCmd()
	initReportCmd()
	initUserCmd()
}

// Execute runs the CLI
func Execute() error {
	return CLI.RootCmd.Execute()
}

// initValidateCmd initializes the validate command
func initValidateCmd() {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate data quality rules",
		Long:  "Validate data quality rules for a Delta Lake table",
		RunE: func(cmd *cobra.Command, args []string) error {
			if CLI.TablePath == "" {
				return fmt.Errorf("table path is required")
			}

			// Initialize engine and datalake reader
			engine := engine.New()
			datalake := datalake.NewReader(CLI.TablePath)
			if err := datalake.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize datalake reader: %w", err)
			}

			// Validate table
			result, err := engine.ValidateTable(context.Background(), CLI.TablePath)
			if err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}

			// Format and print result
			return formatAndPrint(result)
		},
	}

	validateCmd.Flags().StringVar(&CLI.TablePath, "table", "", "Path to Delta Lake table")
	CLI.RootCmd.AddCommand(validateCmd)
}

// initProfileCmd initializes the profile command
func initProfileCmd() {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Generate data profile",
		Long:  "Generate a data profile for a Delta Lake table",
		RunE: func(cmd *cobra.Command, args []string) error {
			if CLI.TablePath == "" {
				return fmt.Errorf("table path is required")
			}

			// Initialize engine and datalake reader
			engine := engine.New()
			datalake := datalake.NewReader(CLI.TablePath)
			if err := datalake.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize datalake reader: %w", err)
			}

			// Profile table
			result, err := engine.ProfileTable(context.Background(), CLI.TablePath)
			if err != nil {
				return fmt.Errorf("profiling failed: %w", err)
			}

			// Format and print result
			return formatAndPrint(result)
		},
	}

	profileCmd.Flags().StringVar(&CLI.TablePath, "table", "", "Path to Delta Lake table")
	CLI.RootCmd.AddCommand(profileCmd)
}

// initMonitorCmd initializes the monitor command
func initMonitorCmd() {
	monitorCmd := &cobra.Command{
		Use:   "monitor",
		Short: "Monitor table metrics",
		Long:  "Monitor metrics and generate alerts for a Delta Lake table",
		RunE: func(cmd *cobra.Command, args []string) error {
			if CLI.TablePath == "" {
				return fmt.Errorf("table path is required")
			}

			// Initialize engine and datalake reader
			engine := engine.New()
			datalake := datalake.NewReader(CLI.TablePath)
			if err := datalake.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize datalake reader: %w", err)
			}

			// Monitor table
			result, err := engine.MonitorTable(context.Background(), CLI.TablePath)
			if err != nil {
				return fmt.Errorf("monitoring failed: %w", err)
			}

			// Format and print result
			return formatAndPrint(result)
		},
	}

	monitorCmd.Flags().StringVar(&CLI.TablePath, "table", "", "Path to Delta Lake table")
	CLI.RootCmd.AddCommand(monitorCmd)
}

// initReportCmd initializes the report command
func initReportCmd() {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate quality report",
		Long:  "Generate a comprehensive quality report for a Delta Lake table",
		RunE: func(cmd *cobra.Command, args []string) error {
			if CLI.TablePath == "" {
				return fmt.Errorf("table path is required")
			}

			// Initialize engine and datalake reader
			engine := engine.New()
			datalake := datalake.NewReader(CLI.TablePath)
			if err := datalake.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize datalake reader: %w", err)
			}

			// Generate report
			result, err := engine.GenerateReport(context.Background(), CLI.TablePath)
			if err != nil {
				return fmt.Errorf("report generation failed: %w", err)
			}

			// Format and print result
			return formatAndPrint(result)
		},
	}

	reportCmd.Flags().StringVar(&CLI.TablePath, "table", "", "Path to Delta Lake table")
	CLI.RootCmd.AddCommand(reportCmd)
}

// initUserCmd initializes user management commands
func initUserCmd() {
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management",
		Long:  "Manage users and API keys",
	}

	// Add subcommands
	initUserAddCmd(userCmd)
	initUserListCmd(userCmd)
	initApiKeyCmd(userCmd)

	CLI.RootCmd.AddCommand(userCmd)
}

// formatAndPrint formats and prints the result based on output format
func formatAndPrint(result interface{}) error {
	switch CLI.OutputFormat {
	case "json":
		jsonBytes, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}
		fmt.Println(string(jsonBytes))
	case "yaml":
		// TODO: Implement YAML formatting
	case "text":
		// TODO: Implement text formatting
	default:
		return fmt.Errorf("unsupported output format: %s", CLI.OutputFormat)
	}

	return nil
}
