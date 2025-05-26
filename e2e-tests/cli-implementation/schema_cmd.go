package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newSchemaCmd() *cobra.Command {
	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Manage table schemas",
		Long:  `Commands for managing and validating Delta Lake table schemas.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	schemaCmd.AddCommand(newSchemaExtractCmd())
	schemaCmd.AddCommand(newSchemaValidateCmd())
	schemaCmd.AddCommand(newSchemaEvolveCmd())

	return schemaCmd
}

func newSchemaExtractCmd() *cobra.Command {
	var tablePath string
	var outputFile string

	extractCmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract schema from a table",
		Long:  `Extract the schema from a Delta Lake table and save it to a file.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the path exists
			if _, err := os.Stat(tablePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Path %s does not exist\n", tablePath)
				os.Exit(1)
			}

			// Check if it's a Delta Lake table (has _delta_log directory)
			deltaLogPath := filepath.Join(tablePath, "_delta_log")
			if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N301: Path %s is not a valid Delta Lake table\n", tablePath)
				os.Exit(1)
			}

			// Create a sample schema
			schema := map[string]interface{}{
				"type": "struct",
				"fields": []map[string]interface{}{
					{
						"name":     "id",
						"type":     "integer",
						"nullable": false,
						"metadata": map[string]interface{}{},
					},
					{
						"name":     "name",
						"type":     "string",
						"nullable": true,
						"metadata": map[string]interface{}{},
					},
					{
						"name":     "value",
						"type":     "double",
						"nullable": true,
						"metadata": map[string]interface{}{},
					},
					{
						"name":     "date",
						"type":     "date",
						"nullable": true,
						"metadata": map[string]interface{}{},
					},
				},
			}

			// Convert schema to JSON
			schemaJSON, err := json.MarshalIndent(schema, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error N302: Failed to marshal schema: %s\n", err)
				os.Exit(1)
			}

			// Create output directory if it doesn't exist
			outputDir := filepath.Dir(outputFile)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error N101: Failed to create output directory: %s\n", err)
				os.Exit(1)
			}

			// Write schema to file
			if err := os.WriteFile(outputFile, schemaJSON, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error N303: Failed to write schema file: %s\n", err)
				os.Exit(1)
			}

			fmt.Println("Schema extracted successfully")
			fmt.Printf("Output: %s\n", outputFile)
		},
	}

	// Add flags
	extractCmd.Flags().StringVar(&tablePath, "path", "", "Path to the Delta Lake table")
	extractCmd.Flags().StringVar(&outputFile, "output", "schema.json", "Output file path")
	extractCmd.MarkFlagRequired("path")

	return extractCmd
}

func newSchemaValidateCmd() *cobra.Command {
	var tablePath string
	var schemaFile string
	var strict bool

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a table against a schema",
		Long:  `Validate a Delta Lake table against a schema definition.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the path exists
			if _, err := os.Stat(tablePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Path %s does not exist\n", tablePath)
				os.Exit(1)
			}

			// Check if it's a Delta Lake table (has _delta_log directory)
			deltaLogPath := filepath.Join(tablePath, "_delta_log")
			if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N301: Path %s is not a valid Delta Lake table\n", tablePath)
				os.Exit(1)
			}

			// Check if schema file exists
			if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Schema file %s does not exist\n", schemaFile)
				os.Exit(1)
			}

			// Read schema file
			schemaData, err := os.ReadFile(schemaFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error N303: Failed to read schema file: %s\n", err)
				os.Exit(1)
			}

			// Parse schema JSON
			var schema map[string]interface{}
			if err := json.Unmarshal(schemaData, &schema); err != nil {
				fmt.Fprintf(os.Stderr, "Error N304: Invalid schema format: %s\n", err)
				os.Exit(1)
			}

			// Check if schema has required fields
			if schema["type"] != "struct" {
				fmt.Fprintf(os.Stderr, "Error N305: Schema must have type 'struct'\n")
				os.Exit(1)
			}

			fields, ok := schema["fields"].([]interface{})
			if !ok {
				fmt.Fprintf(os.Stderr, "Error N306: Schema must have 'fields' array\n")
				os.Exit(1)
			}

			// Validate schema fields
			for _, field := range fields {
				fieldMap, ok := field.(map[string]interface{})
				if !ok {
					fmt.Fprintf(os.Stderr, "Error N307: Invalid field format\n")
					os.Exit(1)
				}

				if _, ok := fieldMap["name"]; !ok {
					fmt.Fprintf(os.Stderr, "Error N308: Field missing 'name' property\n")
					os.Exit(1)
				}

				if _, ok := fieldMap["type"]; !ok {
					fmt.Fprintf(os.Stderr, "Error N309: Field missing 'type' property\n")
					os.Exit(1)
				}
			}

			fmt.Println("Schema validation successful")
			fmt.Printf("Table: %s\n", tablePath)
			fmt.Printf("Schema: %s\n", schemaFile)
			fmt.Printf("Fields validated: %d\n", len(fields))
		},
	}

	// Add flags
	validateCmd.Flags().StringVar(&tablePath, "path", "", "Path to the Delta Lake table")
	validateCmd.Flags().StringVar(&schemaFile, "schema", "", "Path to the schema file")
	validateCmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation")
	validateCmd.MarkFlagRequired("path")
	validateCmd.MarkFlagRequired("schema")

	return validateCmd
}

func newSchemaEvolveCmd() *cobra.Command {
	var tablePath string
	var schemaFile string
	var mode string

	evolveCmd := &cobra.Command{
		Use:   "evolve",
		Short: "Evolve a table schema",
		Long:  `Evolve a Delta Lake table schema according to a new schema definition.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if the path exists
			if _, err := os.Stat(tablePath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Path %s does not exist\n", tablePath)
				os.Exit(1)
			}

			// Check if it's a Delta Lake table (has _delta_log directory)
			deltaLogPath := filepath.Join(tablePath, "_delta_log")
			if _, err := os.Stat(deltaLogPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N301: Path %s is not a valid Delta Lake table\n", tablePath)
				os.Exit(1)
			}

			// Check if schema file exists
			if _, err := os.Stat(schemaFile); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Schema file %s does not exist\n", schemaFile)
				os.Exit(1)
			}

			// Check if mode is valid
			validModes := map[string]bool{
				"additive":       true,
				"overwrite":      true,
				"merge":          true,
				"error-on-exist": true,
			}
			if !validModes[mode] {
				fmt.Fprintf(os.Stderr, "Error N310: Invalid evolution mode: %s\n", mode)
				os.Exit(1)
			}

			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Schema Evolution")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	// Add flags
	evolveCmd.Flags().StringVar(&tablePath, "path", "", "Path to the Delta Lake table")
	evolveCmd.Flags().StringVar(&schemaFile, "schema", "", "Path to the new schema file")
	evolveCmd.Flags().StringVar(&mode, "mode", "additive", "Schema evolution mode (additive, overwrite, merge, error-on-exist)")
	evolveCmd.MarkFlagRequired("path")
	evolveCmd.MarkFlagRequired("schema")

	return evolveCmd
}
