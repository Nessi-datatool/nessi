package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "field-metadata",
	Short: "Manage field-level metadata for Delta tables",
	Long:  `Add, update, or view metadata for fields in Delta tables, including descriptions, tags, and custom properties.`,
}

var getCmd = &cobra.Command{
	Use:   "get [table-path] [field-name]",
	Short: "Get metadata for a specific field",
	Long:  `Get metadata for a specific field in a Delta table, including description, tags, and custom properties.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tablePath := args[0]
		fieldName := args[1]

		// Create reader
		reader, err := datalake.NewReader(tablePath)
		if err != nil {
			return fmt.Errorf("failed to create reader: %w", err)
		}
		defer reader.Close()

		// Get schema
		schema := reader.GetSchema()
		if schema == nil {
			return fmt.Errorf("failed to get schema")
		}

		// Get field metadata
		metadata, err := datalake.GetFieldMetadata(schema, fieldName)
		if err != nil {
			return fmt.Errorf("failed to get field metadata: %w", err)
		}

		// Format output based on format flag
		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "json":
			// Output as JSON
			jsonBytes, err := json.MarshalIndent(metadata, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
			}
			fmt.Println(string(jsonBytes))
		default:
			// Output as human-readable text
			fmt.Println(datalake.FormatFieldMetadata(metadata))
		}

		return nil
	},
}

var getAllCmd = &cobra.Command{
	Use:   "get-all [table-path]",
	Short: "Get metadata for all fields",
	Long:  `Get metadata for all fields in a Delta table, including descriptions, tags, and custom properties.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tablePath := args[0]

		// Create reader
		reader, err := datalake.NewReader(tablePath)
		if err != nil {
			return fmt.Errorf("failed to create reader: %w", err)
		}
		defer reader.Close()

		// Get schema
		schema := reader.GetSchema()
		if schema == nil {
			return fmt.Errorf("failed to get schema")
		}

		// Get all fields metadata
		metadata, err := datalake.GetAllFieldsMetadata(schema)
		if err != nil {
			return fmt.Errorf("failed to get all fields metadata: %w", err)
		}

		// Format output based on format flag
		format, _ := cmd.Flags().GetString("format")
		switch format {
		case "json":
			// Output as JSON
			jsonBytes, err := json.MarshalIndent(metadata, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
			}
			fmt.Println(string(jsonBytes))
		default:
			// Output as human-readable text
			fmt.Println(datalake.FormatAllFieldsMetadata(metadata))
		}

		return nil
	},
}

var addCmd = &cobra.Command{
	Use:   "add [table-path] [field-name]",
	Short: "Add metadata to a field",
	Long:  `Add metadata to a field in a Delta table, including description, tags, and custom properties.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		tablePath := args[0]
		fieldName := args[1]

		// Get flags
		description, _ := cmd.Flags().GetString("description")
		tagsStr, _ := cmd.Flags().GetString("tags")
		propertiesStr, _ := cmd.Flags().GetString("properties")

		// Parse tags
		var tags []string
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}

		// Parse properties
		properties := make(map[string]string)
		if propertiesStr != "" {
			propPairs := strings.Split(propertiesStr, ",")
			for _, pair := range propPairs {
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) == 2 {
					key := strings.TrimSpace(kv[0])
					value := strings.TrimSpace(kv[1])
					properties[key] = value
				}
			}
		}

		// Create metadata manager
		metadataManager := datalake.NewMetadataManager(tablePath)

		// Read table metadata
		table, err := metadataManager.ReadTableMetadata()
		if err != nil {
			return fmt.Errorf("failed to read table metadata: %w", err)
		}

		// Add field metadata
		updatedSchema, err := datalake.AddFieldMetadata(table.Schema, fieldName, description, tags, properties)
		if err != nil {
			return fmt.Errorf("failed to add field metadata: %w", err)
		}

		// Update schema in table
		table.Schema = updatedSchema

		// Write updated metadata
		err = metadataManager.WriteTableMetadata(table)
		if err != nil {
			return fmt.Errorf("failed to write table metadata: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Successfully added metadata to field '%s'\n", fieldName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(getAllCmd)
	rootCmd.AddCommand(addCmd)

	// Add flags for get commands
	getCmd.Flags().StringP("format", "f", "text", "Output format (text or json)")
	getAllCmd.Flags().StringP("format", "f", "text", "Output format (text or json)")

	// Add flags for add command
	addCmd.Flags().StringP("description", "d", "", "Field description")
	addCmd.Flags().StringP("tags", "t", "", "Comma-separated list of tags")
	addCmd.Flags().StringP("properties", "p", "", "Comma-separated list of key=value pairs")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
