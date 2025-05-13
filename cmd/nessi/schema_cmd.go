package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Manage Delta Lake schemas",
	Long:  `Manage Delta Lake schemas, including initialization, updates, validation, and history tracking.`,
}

// initSchemaCmd represents the init-schema command
var initSchemaCmd = &cobra.Command{
	Use:   "init [table_path] [schema_file]",
	Short: "Initialize a schema for a Delta Lake table",
	Long:  `Initialize a schema for a Delta Lake table using a JSON schema definition.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		schemaFile := args[1]

		// Get partition and Z-order columns
		partitionBy, _ := cmd.Flags().GetStringSlice("partition-by")
		zOrderBy, _ := cmd.Flags().GetStringSlice("z-order-by")

		// Read schema file
		schemaData, err := os.ReadFile(schemaFile)
		if err != nil {
			fmt.Printf("Error reading schema file: %v\n", err)
			os.Exit(1)
		}

		// Parse schema
		var schema datalake.Schema
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			fmt.Printf("Error parsing schema file: %v\n", err)
			os.Exit(1)
		}

		// Initialize schema
		sm := datalake.NewSchemaManager(tablePath)
		schemaVersion, err := sm.InitializeSchema(&schema, partitionBy, zOrderBy)
		if err != nil {
			fmt.Printf("Error initializing schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Schema initialized successfully (version %d)\n", schemaVersion.Version)
		fmt.Printf("Fields: %d\n", len(schema.Fields))
		fmt.Printf("Partition by: %v\n", partitionBy)
		fmt.Printf("Z-order by: %v\n", zOrderBy)
	},
}

// updateSchemaCmd represents the update-schema command
var updateSchemaCmd = &cobra.Command{
	Use:   "update [table_path] [schema_file]",
	Short: "Update a schema for a Delta Lake table",
	Long:  `Update a schema for a Delta Lake table using a JSON schema definition.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		schemaFile := args[1]

		// Read schema file
		schemaData, err := os.ReadFile(schemaFile)
		if err != nil {
			fmt.Printf("Error reading schema file: %v\n", err)
			os.Exit(1)
		}

		// Parse schema
		var schema datalake.Schema
		if err := json.Unmarshal(schemaData, &schema); err != nil {
			fmt.Printf("Error parsing schema file: %v\n", err)
			os.Exit(1)
		}

		// Get commit message
		message, _ := cmd.Flags().GetString("message")
		commitInfo := map[string]string{
			"action":      "update_schema",
			"message":     message,
			"timestamp":   time.Now().Format(time.RFC3339),
			"user":        os.Getenv("USER"),
		}

		// Update schema
		sm := datalake.NewSchemaManager(tablePath)
		schemaVersion, err := sm.UpdateSchema(&schema, commitInfo)
		if err != nil {
			fmt.Printf("Error updating schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Schema updated successfully (version %d)\n", schemaVersion.Version)
		fmt.Printf("Fields: %d\n", len(schema.Fields))
		fmt.Println("Changes:")
		for _, change := range schemaVersion.Changes {
			fmt.Printf("  - %s\n", change.Description)
		}
	},
}

// updatePartitioningCmd represents the update-partitioning command
var updatePartitioningCmd = &cobra.Command{
	Use:   "update-partitioning [table_path]",
	Short: "Update partitioning for a Delta Lake table",
	Long:  `Update partitioning and Z-ordering for a Delta Lake table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]

		// Get partition and Z-order columns
		partitionBy, _ := cmd.Flags().GetStringSlice("partition-by")
		zOrderBy, _ := cmd.Flags().GetStringSlice("z-order-by")

		// Get commit message
		message, _ := cmd.Flags().GetString("message")
		commitInfo := map[string]string{
			"action":      "update_partitioning",
			"message":     message,
			"timestamp":   time.Now().Format(time.RFC3339),
			"user":        os.Getenv("USER"),
		}

		// Update partitioning
		sm := datalake.NewSchemaManager(tablePath)
		schemaVersion, err := sm.UpdatePartitioning(partitionBy, zOrderBy, commitInfo)
		if err != nil {
			fmt.Printf("Error updating partitioning: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Partitioning updated successfully (version %d)\n", schemaVersion.Version)
		fmt.Printf("Partition by: %v\n", partitionBy)
		fmt.Printf("Z-order by: %v\n", zOrderBy)
		fmt.Println("Changes:")
		for _, change := range schemaVersion.Changes {
			fmt.Printf("  - %s\n", change.Description)
		}
	},
}

// validateDataCmd represents the validate-data command
var validateDataCmd = &cobra.Command{
	Use:   "validate [table_path] [data_file]",
	Short: "Validate data against a schema",
	Long:  `Validate JSON data against a Delta Lake table schema.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		dataFile := args[1]

		// Read data file
		dataBytes, err := os.ReadFile(dataFile)
		if err != nil {
			fmt.Printf("Error reading data file: %v\n", err)
			os.Exit(1)
		}

		// Parse data
		var data map[string]interface{}
		if err := json.Unmarshal(dataBytes, &data); err != nil {
			fmt.Printf("Error parsing data file: %v\n", err)
			os.Exit(1)
		}

		// Validate data
		sm := datalake.NewSchemaManager(tablePath)
		valid, violations, err := sm.ValidateSchema(data)
		if err != nil {
			fmt.Printf("Error validating data: %v\n", err)
			os.Exit(1)
		}

		if valid {
			fmt.Println("Data is valid!")
		} else {
			fmt.Printf("Data validation failed with %d violations:\n", len(violations))
			for i, violation := range violations {
				fmt.Printf("  %d. %s\n", i+1, violation)
			}
			os.Exit(1)
		}
	},
}

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history [table_path]",
	Short: "Show schema history",
	Long:  `Show the history of schema changes for a Delta Lake table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]

		// Get schema history
		sm := datalake.NewSchemaManager(tablePath)
		history, err := sm.GetSchemaHistory()
		if err != nil {
			fmt.Printf("Error getting schema history: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Schema history for %s:\n", tablePath)
		fmt.Printf("Total versions: %d\n", len(history.Schemas))
		fmt.Printf("Current version: %d\n\n", history.Schemas[history.CurrentIndex].Version)

		for i, schema := range history.Schemas {
			fmt.Printf("Version %d (%s):\n", schema.Version, schema.Timestamp.Format(time.RFC3339))
			
			// Print commit info
			if len(schema.CommitInfo) > 0 {
				fmt.Println("  Commit info:")
				for k, v := range schema.CommitInfo {
					fmt.Printf("    %s: %s\n", k, v)
				}
			}
			
			// Print changes
			if len(schema.Changes) > 0 {
				fmt.Println("  Changes:")
				for _, change := range schema.Changes {
					fmt.Printf("    - %s\n", change.Description)
				}
			}
			
			// Print partitioning
			if len(schema.PartitionBy) > 0 {
				fmt.Printf("  Partition by: %s\n", strings.Join(schema.PartitionBy, ", "))
			}
			
			// Print Z-ordering
			if len(schema.ZOrderBy) > 0 {
				fmt.Printf("  Z-order by: %s\n", strings.Join(schema.ZOrderBy, ", "))
			}
			
			// Print fields
			fmt.Printf("  Fields: %d\n", len(schema.Schema.Fields))
			
			// Mark current version
			if i == history.CurrentIndex {
				fmt.Println("  (current)")
			}
			
			fmt.Println()
		}
	},
}

// hintsCmd represents the hints command
var hintsCmd = &cobra.Command{
	Use:   "hints [table_path]",
	Short: "Get optimization hints",
	Long:  `Get optimization hints for a Delta Lake table schema.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]

		// Get optimization hints
		sm := datalake.NewSchemaManager(tablePath)
		hints, err := sm.GetOptimizationHints()
		if err != nil {
			fmt.Printf("Error getting optimization hints: %v\n", err)
			os.Exit(1)
		}

		if len(hints) == 0 {
			fmt.Println("No optimization hints found. Your schema looks good!")
			return
		}

		fmt.Printf("Optimization hints for %s:\n", tablePath)
		for key, hint := range hints {
			fmt.Printf("  - %s: %s\n", key, hint)
		}
	},
}

// fieldInfoCmd represents the field-info command
var fieldInfoCmd = &cobra.Command{
	Use:   "field-info [table_path] [field_name]",
	Short: "Get field metadata",
	Long:  `Get metadata for a specific field in a Delta Lake table schema.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]
		fieldName := args[1]

		// Get field metadata
		sm := datalake.NewSchemaManager(tablePath)
		metadata, err := sm.GetFieldMetadata(fieldName)
		if err != nil {
			fmt.Printf("Error getting field metadata: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Field metadata for '%s':\n", fieldName)
		fmt.Printf("  Name: %s\n", metadata.Name)
		fmt.Printf("  Type: %s\n", metadata.Type)
		fmt.Printf("  Nullable: %v\n", metadata.Nullable)
		
		if metadata.Description != "" {
			fmt.Printf("  Description: %s\n", metadata.Description)
		}
		
		if len(metadata.Stats) > 0 {
			fmt.Println("  Statistics:")
			for k, v := range metadata.Stats {
				fmt.Printf("    %s: %v\n", k, v)
			}
		}
		
		if len(metadata.Tags) > 0 {
			fmt.Println("  Tags:")
			for k, v := range metadata.Tags {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
	
	// Add subcommands
	schemaCmd.AddCommand(initSchemaCmd)
	schemaCmd.AddCommand(updateSchemaCmd)
	schemaCmd.AddCommand(updatePartitioningCmd)
	schemaCmd.AddCommand(validateDataCmd)
	schemaCmd.AddCommand(historyCmd)
	schemaCmd.AddCommand(hintsCmd)
	schemaCmd.AddCommand(fieldInfoCmd)
	
	// Add flags
	initSchemaCmd.Flags().StringSlice("partition-by", []string{}, "Columns to partition by")
	initSchemaCmd.Flags().StringSlice("z-order-by", []string{}, "Columns to Z-order by")
	
	updateSchemaCmd.Flags().String("message", "", "Commit message for the schema update")
	updateSchemaCmd.MarkFlagRequired("message")
	
	updatePartitioningCmd.Flags().StringSlice("partition-by", []string{}, "Columns to partition by")
	updatePartitioningCmd.Flags().StringSlice("z-order-by", []string{}, "Columns to Z-order by")
	updatePartitioningCmd.Flags().String("message", "", "Commit message for the partitioning update")
	updatePartitioningCmd.MarkFlagRequired("message")
}
