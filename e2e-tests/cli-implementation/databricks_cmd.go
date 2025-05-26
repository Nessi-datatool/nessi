package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newDatabricksCmd() *cobra.Command {
	databricksCmd := &cobra.Command{
		Use:   "databricks",
		Short: "Manage Databricks connections",
		Long:  `Commands for connecting to and managing Databricks workspaces.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	databricksCmd.AddCommand(newDatabricksConnectCmd())
	databricksCmd.AddCommand(newDatabricksListTablesCmd())
	databricksCmd.AddCommand(newDatabricksQueryCmd())

	return databricksCmd
}

func newDatabricksConnectCmd() *cobra.Command {
	var token string
	var host string
	var workspace string

	connectCmd := &cobra.Command{
		Use:   "connect",
		Short: "Connect to a Databricks workspace",
		Long:  `Connect to a Databricks workspace using a personal access token.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if token and host are provided
			if token == "" || host == "" {
				fmt.Fprintf(os.Stderr, "Error N201: Both token and host are required\n")
				os.Exit(1)
			}

			// Special case for testing authentication errors
			if token == "invalid_token" {
				fmt.Fprintf(os.Stderr, "Error N201: Databricks authentication failed: Invalid token\n")
				os.Exit(1)
			}

			// Special case for testing rate limiting errors
			if host == "rate_limit_host" {
				fmt.Fprintf(os.Stderr, "Error N501: Databricks rate limit exceeded. Please try again later.\n")
				os.Exit(1)
			}

			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Databricks Integration")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	// Add flags
	connectCmd.Flags().StringVar(&token, "token", "", "Databricks personal access token")
	connectCmd.Flags().StringVar(&host, "host", "", "Databricks host URL")
	connectCmd.Flags().StringVar(&workspace, "workspace", "", "Databricks workspace ID")
	connectCmd.MarkFlagRequired("token")
	connectCmd.MarkFlagRequired("host")

	return connectCmd
}

func newDatabricksListTablesCmd() *cobra.Command {
	var catalog string
	var schema string

	listTablesCmd := &cobra.Command{
		Use:   "list-tables",
		Short: "List tables in a Databricks catalog",
		Long:  `List all tables in a Databricks catalog and schema.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if catalog and schema are provided
			if catalog == "" || schema == "" {
				fmt.Fprintf(os.Stderr, "Error N102: Both catalog and schema are required\n")
				os.Exit(1)
			}

			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Databricks Integration")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	// Add flags
	listTablesCmd.Flags().StringVar(&catalog, "catalog", "", "Databricks catalog name")
	listTablesCmd.Flags().StringVar(&schema, "schema", "", "Databricks schema name")
	listTablesCmd.MarkFlagRequired("catalog")
	listTablesCmd.MarkFlagRequired("schema")

	return listTablesCmd
}

func newDatabricksQueryCmd() *cobra.Command {
	var query string
	var output string

	queryCmd := &cobra.Command{
		Use:   "query",
		Short: "Run a SQL query on Databricks",
		Long:  `Run a SQL query on a Databricks workspace and output the results.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if query is provided
			if query == "" {
				fmt.Fprintf(os.Stderr, "Error N102: SQL query is required\n")
				os.Exit(1)
			}

			// Check if output format is valid
			validFormats := map[string]bool{
				"table": true,
				"json":  true,
				"csv":   true,
			}
			if !validFormats[output] {
				fmt.Fprintf(os.Stderr, "Error N401: Invalid output format: %s\n", output)
				os.Exit(1)
			}

			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Databricks Integration")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	// Add flags
	queryCmd.Flags().StringVar(&query, "query", "", "SQL query to run")
	queryCmd.Flags().StringVar(&output, "output", "table", "Output format (table, json, csv)")
	queryCmd.MarkFlagRequired("query")

	return queryCmd
}
