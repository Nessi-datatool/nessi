package main

import (
	"fmt"
	"os"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/nessi-dev/nessi/pkg/datalake"
)

func main() {
	// Path to the test Delta table
	tablePath := "test_data/sample_table"

	// Create a schema for testing
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	// Create metadata manager
	metadataManager := datalake.NewMetadataManager(tablePath)

	// Create a Delta table with schema
	table := &datalake.DeltaTable{
		Path:       tablePath,
		Version:    0,
		Schema:     schema,
		Partitions: make(map[string][]string),
		Files:      []string{},
		Metadata:   make(map[string]interface{}),
	}

	// Add field metadata
	fieldName := "id"
	description := "Unique identifier for each item"
	tags := []string{"primary_key", "indexed"}
	properties := map[string]string{
		"min_value": "1",
		"max_value": "1000",
	}

	// Add field metadata
	updatedSchema, err := datalake.AddFieldMetadata(schema, fieldName, description, tags, properties)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding field metadata: %v\n", err)
		os.Exit(1)
	}

	// Update schema in table
	table.Schema = updatedSchema

	// Write updated metadata
	err = metadataManager.WriteTableMetadata(table)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing table metadata: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added metadata to field '%s'\n", fieldName)

	// Now let's read the metadata back to verify it worked
	reader, err := datalake.NewReader(tablePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating reader: %v\n", err)
		os.Exit(1)
	}
	defer reader.Close()

	// Get schema
	readSchema := reader.GetSchema()
	if readSchema == nil {
		fmt.Fprintf(os.Stderr, "Failed to get schema\n")
		os.Exit(1)
	}

	// Get field metadata
	metadata, err := datalake.GetFieldMetadata(readSchema, fieldName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting field metadata: %v\n", err)
		os.Exit(1)
	}

	// Print metadata
	fmt.Println("\nField Metadata:")
	fmt.Println(datalake.FormatFieldMetadata(metadata))
}
