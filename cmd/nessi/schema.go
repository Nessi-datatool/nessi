package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
)

var (
	// Schema command flags
	outputFormat string
	fromVersion  int64
	toVersion    int64
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage and inspect Delta table schemas",
	Long: `Schema commands allow you to:
- View schema history for a Delta table
- Compare schemas between versions
- Validate schema compatibility`,
}

var historyCmd = &cobra.Command{
	Use:   "history [table_path]",
	Short: "Show schema history for a Delta table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tablePath := args[0]
		
		// Create metadata manager
		metaManager := datalake.NewMetadataManager(tablePath)
		
		// Get schema history
		history, err := metaManager.GetSchemaHistory()
		if err != nil {
			return fmt.Errorf("failed to get schema history: %w", err)
		}
		
		if outputFormat == "json" {
			// TODO: Implement JSON output
			fmt.Println("JSON output not yet implemented")
		} else {
			// Text output
			fmt.Println(datalake.FormatSchemaHistory(history))
		}
		
		return nil
	},
}

var diffCmd = &cobra.Command{
	Use:   "diff [table_path] [version1] [version2]",
	Short: "Show schema differences between two versions",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		tablePath := args[0]
		
		// Parse versions
		v1, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid version1: %w", err)
		}
		
		v2, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid version2: %w", err)
		}
		
		// Create metadata manager
		metaManager := datalake.NewMetadataManager(tablePath)
		
		// Get tables at versions
		table1, err := metaManager.GetTableAtVersion(v1)
		if err != nil {
			return fmt.Errorf("failed to get table at version %d: %w", v1, err)
		}
		
		table2, err := metaManager.GetTableAtVersion(v2)
		if err != nil {
			return fmt.Errorf("failed to get table at version %d: %w", v2, err)
		}
		
		// Compare schemas
		changes := datalake.DiffSchemas(table1.Schema, table2.Schema)
		
		// Print results
		fmt.Printf("Schema differences between version %d (%s) and version %d (%s):\n\n",
			v1, table1.LastModified.Format(time.RFC3339),
			v2, table2.LastModified.Format(time.RFC3339))
		
		if len(changes) == 0 {
			fmt.Println("No schema changes detected.")
		} else {
			fmt.Println(datalake.FormatSchemaChanges(changes))
		}
		
		return nil
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate [table_path] [schema_file]",
	Short: "Validate schema compatibility",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tablePath := args[0]
		schemaFile := args[1]
		
		// Create metadata manager
		metaManager := datalake.NewMetadataManager(tablePath)
		
		// Get current schema
		table, err := metaManager.ReadTableMetadata()
		if err != nil {
			return fmt.Errorf("failed to read table metadata: %w", err)
		}
		
		// Read schema file
		schemaData, err := os.ReadFile(schemaFile)
		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}
		
		// TODO: Parse schema file and validate compatibility
		fmt.Printf("Schema validation not yet implemented. Current schema has %d fields.\n", len(table.Schema.Fields()))
		fmt.Printf("Schema file: %s (%d bytes)\n", schemaFile, len(schemaData))
		
		return nil
	},
}

func init() {
	// Schema command flags
	schemaCmd.PersistentFlags().StringVar(&outputFormat, "format", "text", "Output format (text/json)")
	
	// Add subcommands
	schemaCmd.AddCommand(historyCmd)
	schemaCmd.AddCommand(diffCmd)
	schemaCmd.AddCommand(validateCmd)
	
	// Add to root command
	rootCmd.AddCommand(schemaCmd)
}
