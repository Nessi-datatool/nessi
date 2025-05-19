package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/spf13/cobra"
)

// schemaTreeCmd represents the schema tree command
var schemaTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Display schema as an ASCII tree",
	Long:  `Display the schema of a Delta Lake table as an ASCII tree for better visualization.`,
	Run: func(cmd *cobra.Command, args []string) {
		tablePath, _ := cmd.Flags().GetString("table")
		maxDepth, _ := cmd.Flags().GetInt("max-depth")
		showTypes, _ := cmd.Flags().GetBool("show-types")
		showNullable, _ := cmd.Flags().GetBool("show-nullable")

		// Open the Delta table
		table, err := datalake.OpenTable(tablePath)
		if err != nil {
			fmt.Printf("Error opening table: %v\n", err)
			os.Exit(1)
		}

		// Get the schema
		schema, err := table.GetSchema()
		if err != nil {
			fmt.Printf("Error getting schema: %v\n", err)
			os.Exit(1)
		}

		// Print the schema as an ASCII tree
		fmt.Printf("Schema tree for table: %s\n\n", tablePath)
		printSchemaTree(schema, maxDepth, showTypes, showNullable)
	},
}

func init() {
	schemaCmd.AddCommand(schemaTreeCmd)

	// Add flags
	schemaTreeCmd.Flags().StringP("table", "t", "", "Path to the Delta Lake table (required)")
	schemaTreeCmd.MarkFlagRequired("table")
	schemaTreeCmd.Flags().IntP("max-depth", "d", 10, "Maximum depth to display for nested fields")
	schemaTreeCmd.Flags().BoolP("show-types", "y", true, "Show field types")
	schemaTreeCmd.Flags().BoolP("show-nullable", "n", true, "Show nullable status")
}

// printSchemaTree prints the schema as an ASCII tree
func printSchemaTree(schema *datalake.Schema, maxDepth int, showTypes, showNullable bool) {
	if schema == nil || len(schema.Fields) == 0 {
		fmt.Println("Empty schema")
		return
	}

	fmt.Println("└── Table")
	for i, field := range schema.Fields {
		prefix := "    "
		if i == len(schema.Fields)-1 {
			printField(field, prefix+"└── ", prefix+"    ", 1, maxDepth, showTypes, showNullable)
		} else {
			printField(field, prefix+"├── ", prefix+"│   ", 1, maxDepth, showTypes, showNullable)
		}
	}
}

// printField prints a field and its children recursively
func printField(field *datalake.Field, prefix, childPrefix string, depth, maxDepth int, showTypes, showNullable bool) {
	// Print field name
	fieldInfo := field.Name
	
	// Add type information if requested
	if showTypes {
		fieldInfo += fmt.Sprintf(" (%s)", field.Type)
	}
	
	// Add nullable information if requested
	if showNullable && field.Nullable {
		fieldInfo += " [nullable]"
	}
	
	fmt.Println(prefix + fieldInfo)
	
	// Stop recursion if we've reached the maximum depth
	if depth >= maxDepth {
		return
	}
	
	// Print nested fields for struct types
	if field.Type == "struct" && field.Children != nil {
		for i, child := range field.Children {
			newPrefix := childPrefix
			if i == len(field.Children)-1 {
				printField(child, childPrefix+"└── ", childPrefix+"    ", depth+1, maxDepth, showTypes, showNullable)
			} else {
				printField(child, childPrefix+"├── ", childPrefix+"│   ", depth+1, maxDepth, showTypes, showNullable)
			}
		}
	}
	
	// Print array element type for array types
	if strings.HasPrefix(field.Type, "array") && field.Children != nil && len(field.Children) > 0 {
		printField(field.Children[0], childPrefix+"└── [elements] ", childPrefix+"    ", depth+1, maxDepth, showTypes, showNullable)
	}
	
	// Print key-value types for map types
	if strings.HasPrefix(field.Type, "map") && field.Children != nil && len(field.Children) >= 2 {
		printField(field.Children[0], childPrefix+"├── [keys] ", childPrefix+"│   ", depth+1, maxDepth, showTypes, showNullable)
		printField(field.Children[1], childPrefix+"└── [values] ", childPrefix+"    ", depth+1, maxDepth, showTypes, showNullable)
	}
}
