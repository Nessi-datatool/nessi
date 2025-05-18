package cli

import (
	"github.com/nessi-dev/nessi/pkg"
	"os"
	"github.com/spf13/cobra"
)


import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/ipc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/nessi-dev/nessi/pkg/quality/engine"
	"github.com/nessi-dev/nessi/pkg/security"
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
	// Add info command
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show Nessi CLI and feature info",
		Long:  "Display information about the Nessi CLI, enabled features, and upgrade options.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`\nNessi CLI %s\n`, CLI.RootCmd.Version)
			fmt.Println("Delta Lake support: enabled")
			fmt.Println("Python extensions: enabled")
			fmt.Println("Telemetry: disabled")
			fmt.Println("Team features: not available (OSS Edition)")
			fmt.Println("\n→ Want RCA dashboards, team alerts, or governance features?")
			fmt.Println("   Learn about LakeDiff: https://lakediff.com\n")
			os.Exit(0)
		},
	}
	CLI.RootCmd.AddCommand(infoCmd)

	// Add team-features command
	teamFeaturesCmd := &cobra.Command{
		Use:   "team-features",
		Short: "Show team and enterprise features available in LakeDiff",
		Long:  "Display a list of features available in LakeDiff Team & Enterprise edition.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`\n🚀 Team & Enterprise Features (via LakeDiff)\n\n✔ Slack/Teams alerts\n✔ RCA dashboards with lineage\n✔ Role-based access control (RBAC)\n✔ dbt Cloud model mapping\n✔ Trend dashboards (Grafana/Prometheus)\n→ Learn more: https://lakediff.com/features\n`)
			os.Exit(0)
		},
	}
	CLI.RootCmd.AddCommand(teamFeaturesCmd)

	CLI.Config = viper.New()
	CLI.RootCmd = &cobra.Command{
		Use:   "nessi",
		Short: "Nessi - Delta Lake Data Quality Tool",
		Long: `Nessi is a powerful tool for managing data quality in Delta Lake tables.
It provides features for validation, profiling, monitoring, and reporting.`,
		Version: "v0.10.3", // Update as needed
	} // Add Version field to the root command
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
	initSchemaTreeCmd() // Register the new schema-tree command
	// RCA command is disabled due to missing implementation
	rootCmd.AddCommand(telemetryCmd)
}

// Execute runs the CLI
func Execute() error {
	cfg := pkg.GetTelemetryConfig()
	telemetryManager := pkg.NewTelemetryManager(cfg)
	CLI.RootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Send telemetry event for every command
		telemetryManager.SendEvent(pkg.TelemetryEvent{
			Timestamp: pkg.Now(),
			Command:   cmd.Name(),
			Args:      args,
			Success:   true, // Set to false on error in future
			Error:     "",
			Version:   "v0.1.0", // TODO: set dynamically
		})
	}
	// Check for --team or --enterprise flag in os.Args
	for _, arg := range os.Args[1:] {
		if arg == "--team" || arg == "--enterprise" {
			fmt.Println(`\n\U0001F512 Team & Enterprise Features (LakeDiff)
-----------------------------------------
✔ Slack/Teams alerts
✔ RCA dashboards (with lineage views)
✔ SSO and role-based access
✔ dbt Cloud and orchestration
→ Learn more: https://lakediff.com/features\n`)
			os.Exit(0)
		}
	}
	return CLI.RootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(telemetryCmd)

	CLI.RootCmd.Version = "v0.10.3"
	CLI.RootCmd.SetVersionTemplate(`Nessi CLI {{.Version}} — OSS edition\nLooking for team features like Slack alerts and dashboards? See LakeDiff → https://lakediff.com\n`)
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
			datalake, err := datalake.NewReader(CLI.TablePath)
			if err != nil {
				return fmt.Errorf("failed to create datalake reader: %w", err)
			}
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
			datalake, err := datalake.NewReader(CLI.TablePath)
			if err != nil {
				return fmt.Errorf("failed to create datalake reader: %w", err)
			}
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
		Long:  "Monitor metrics for a Delta Lake table", // OSS: removed alert mention
		RunE: func(cmd *cobra.Command, args []string) error {
			if CLI.TablePath == "" {
				return fmt.Errorf("table path is required")
			}

			// Initialize engine and datalake reader
			engine := engine.New()
			datalake, err := datalake.NewReader(CLI.TablePath)
			if err != nil {
				return fmt.Errorf("failed to create datalake reader: %w", err)
			}
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
			datalake, err := datalake.NewReader(CLI.TablePath)
			if err != nil {
				return fmt.Errorf("failed to create datalake reader: %w", err)
			}
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

	// User management is not implemented yet
	CLI.RootCmd.AddCommand(userCmd)
}

// printSchemaTree recursively prints the schema as an ASCII tree
func printSchemaTree(schema *arrow.Schema, indent int) {
	for i, field := range schema.Fields() {
		for j := 0; j < indent; j++ {
			fmt.Print("  ")
		}
		prefix := "├─"
		if i == len(schema.Fields())-1 {
			prefix = "└─"
		}
		fmt.Printf("%s %s (%s)\n", prefix, field.Name, field.Type)
		// If the field is a struct, recurse (simplified for Arrow)
		if structType, ok := field.Type.(*arrow.StructType); ok {
			subSchema := arrow.NewSchema(structType.Fields(), nil)
			printSchemaTree(subSchema, indent+1)
		}
	}
}

// initSchemaTreeCmd initializes the schema-tree command
func initSchemaTreeCmd() {
	schemaTreeCmd := &cobra.Command{
		Use:   "schema-tree",
		Short: "Display table schema as an ASCII tree",
		Long:  "Prints the schema of a Delta Lake or Parquet table as an ASCII tree.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if CLI.TablePath == "" {
				return fmt.Errorf("table path is required")
			}
			// For simplicity, try to open the file and infer schema
			f, err := os.Open(CLI.TablePath)
			if err != nil {
				return fmt.Errorf("failed to open table: %w", err)
			}
			defer f.Close()
			// Use Arrow IPC to infer schema (could be replaced with actual Parquet schema reading)
			reader, err := ipc.NewReader(f)
			if err != nil {
				return fmt.Errorf("failed to read Arrow IPC: %w", err)
			}
			defer reader.Release()
			schema := reader.Schema()
			fmt.Printf("Schema for table: %s\n", CLI.TablePath)
			printSchemaTree(schema, 0)
			return nil
		},
	}
	schemaTreeCmd.Flags().StringVar(&CLI.TablePath, "table", "", "Path to Delta Lake/Parquet table")
	CLI.RootCmd.AddCommand(schemaTreeCmd)
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
