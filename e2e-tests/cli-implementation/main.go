package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "nessi",
		Short: "Nessi CLI - Data quality and management tool for Delta Lake",
		Long: `Nessi CLI is a comprehensive data quality and management tool for Delta Lake tables.
It provides commands for data quality checks, schema validation, reporting, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add commands
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newTablesCmd())
	rootCmd.AddCommand(newQualityCmd())
	rootCmd.AddCommand(newReportCmd())
	rootCmd.AddCommand(newLicenseCmd())
	rootCmd.AddCommand(newWorkflowCmd())
	rootCmd.AddCommand(newCloudCmd())
	rootCmd.AddCommand(newDatabricksCmd())
	rootCmd.AddCommand(newSchemaCmd())

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
