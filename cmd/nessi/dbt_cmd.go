package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nessi-dev/nessi-dev/pkg/dbt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dbtCmd = &cobra.Command{
	Use:   "dbt",
	Short: "Integration with dbt for data quality and profiling",
	Long: `Integration with dbt (data build tool) for data quality and profiling.

This is an opt-in feature that allows you to run Nessi.dev data quality rules
against dbt models and generate data profiles for Delta tables associated with
dbt models.

To enable this feature, set the 'enable_dbt_plugin' flag to true in your
nessi.yaml configuration file or use the --enable-dbt-plugin flag.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Check if the dbt plugin is enabled
		enabled := viper.GetBool("enable_dbt_plugin")
		if !enabled {
			enableFlag, _ := cmd.Flags().GetBool("enable-dbt-plugin")
			if !enableFlag {
				return fmt.Errorf("the dbt plugin is not enabled. Set 'enable_dbt_plugin: true' in your nessi.yaml or use --enable-dbt-plugin")
			}
		}
		return nil
	},
}

var dbtValidateCmd = &cobra.Command{
	Use:   "validate [--select model_selection]",
	Short: "Validate dbt models using Nessi.dev data quality rules",
	Long: `Validate dbt models using Nessi.dev data quality rules.

This command executes Nessi.dev data quality rules against Delta tables associated
with dbt models. It uses dbt's selection syntax to determine which models to validate.

Examples:
  # Validate all models tagged as 'daily'
  nessi dbt validate --select tag:daily

  # Validate a specific model and all its downstream dependencies
  nessi dbt validate --select my_model+downstream

  # Validate multiple models
  nessi dbt validate --select model1 model2 model3`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get command flags
		modelSelection, _ := cmd.Flags().GetStringSlice("select")
		configPath, _ := cmd.Flags().GetString("config")
		outputFormat, _ := cmd.Flags().GetString("output")
		includeDBTTests, _ := cmd.Flags().GetBool("include-dbt-tests")
		lineageAware, _ := cmd.Flags().GetBool("lineage-aware")
		failOnError, _ := cmd.Flags().GetBool("fail-on-error")
		
		// Initialize dbt validator
		validator, err := dbt.NewValidator(configPath)
		if err != nil {
			fmt.Printf("Error initializing dbt validator: %v\n", err)
			os.Exit(1)
		}
		
		// Set validator options
		validator.SetIncludeDBTTests(includeDBTTests)
		validator.SetLineageAware(lineageAware)
		
		// Run validation
		results, err := validator.Validate(modelSelection)
		if err != nil {
			fmt.Printf("Error validating dbt models: %v\n", err)
			os.Exit(1)
		}
		
		// Output results
		err = outputResults(results, outputFormat)
		if err != nil {
			fmt.Printf("Error outputting results: %v\n", err)
			os.Exit(1)
		}
		
		// Check if we should exit with error
		if failOnError && results.HasFailures() {
			os.Exit(1)
		}
	},
}

var dbtProfileCmd = &cobra.Command{
	Use:   "profile [--select model_selection]",
	Short: "Generate data profiles for dbt models",
	Long: `Generate data profiles for Delta tables associated with dbt models.

This command generates data profiles for Delta tables associated with dbt models.
It uses dbt's selection syntax to determine which models to profile.

Examples:
  # Profile all models tagged as 'daily'
  nessi dbt profile --select tag:daily

  # Profile a specific model and all its downstream dependencies
  nessi dbt profile --select my_model+downstream

  # Profile multiple models
  nessi dbt profile --select model1 model2 model3`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get command flags
		modelSelection, _ := cmd.Flags().GetStringSlice("select")
		configPath, _ := cmd.Flags().GetString("config")
		outputFormat, _ := cmd.Flags().GetString("output")
		profileType, _ := cmd.Flags().GetString("profile-type")
		failOnError, _ := cmd.Flags().GetBool("fail-on-error")
		
		// Initialize dbt profiler
		profiler, err := dbt.NewProfiler(configPath)
		if err != nil {
			fmt.Printf("Error initializing dbt profiler: %v\n", err)
			os.Exit(1)
		}
		
		// Set profiler options
		profiler.SetProfileType(profileType)
		
		// Run profiling
		results, err := profiler.Profile(modelSelection)
		if err != nil {
			fmt.Printf("Error profiling dbt models: %v\n", err)
			os.Exit(1)
		}
		
		// Output results
		err = outputResults(results, outputFormat)
		if err != nil {
			fmt.Printf("Error outputting results: %v\n", err)
			os.Exit(1)
		}
		
		// Check if we should exit with error
		if failOnError && results.HasFailures() {
			os.Exit(1)
		}
	},
}

// outputResults outputs results in the specified format
func outputResults(results interface{}, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return dbt.OutputJSON(results, os.Stdout)
	case "csv":
		return dbt.OutputCSV(results, os.Stdout)
	case "table":
		return dbt.OutputTable(results, os.Stdout)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func init() {
	rootCmd.AddCommand(dbtCmd)
	dbtCmd.AddCommand(dbtValidateCmd)
	dbtCmd.AddCommand(dbtProfileCmd)
	
	// Global dbt command flags
	dbtCmd.PersistentFlags().Bool("enable-dbt-plugin", false, "Enable the dbt plugin")
	dbtCmd.PersistentFlags().String("config", "", "Path to configuration file")
	dbtCmd.PersistentFlags().String("output", "table", "Output format (table, json, csv)")
	dbtCmd.PersistentFlags().Bool("fail-on-error", false, "Exit with non-zero code if validation/profiling fails")
	
	// Validate command flags
	dbtValidateCmd.Flags().StringSlice("select", []string{}, "dbt model selection syntax (e.g., tag:daily, model_name+downstream)")
	dbtValidateCmd.Flags().Bool("include-dbt-tests", false, "Include dbt test results in validation")
	dbtValidateCmd.Flags().Bool("lineage-aware", false, "Enable lineage-aware validation")
	
	// Profile command flags
	dbtProfileCmd.Flags().StringSlice("select", []string{}, "dbt model selection syntax (e.g., tag:daily, model_name+downstream)")
	dbtProfileCmd.Flags().String("profile-type", "basic", "Profile type (basic, enhanced)")
	
	// Register flags with viper
	viper.BindPFlag("enable_dbt_plugin", dbtCmd.PersistentFlags().Lookup("enable-dbt-plugin"))
}
