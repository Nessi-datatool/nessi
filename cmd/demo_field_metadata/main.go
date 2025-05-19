package main

import (
	"fmt"
	"os"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/nessi-dev/nessi/pkg/datalake"
)

// This is a standalone demo that shows how field-level metadata works
// without relying on the Delta table format
func main() {
	fmt.Println("Field-level Metadata Demo")
	fmt.Println("========================")

	// Create a schema for demonstration
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	fmt.Println("\nOriginal Schema:")
	printSchema(schema)

	// Add metadata to the "id" field
	fmt.Println("\nAdding metadata to 'id' field...")
	description := "Unique identifier for each item"
	tags := []string{"primary_key", "indexed"}
	properties := map[string]string{
		"min_value": "1",
		"max_value": "1000",
	}

	updatedSchema, err := datalake.AddFieldMetadata(schema, "id", description, tags, properties)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding field metadata: %v\n", err)
		os.Exit(1)
	}

	// Get and display the metadata
	metadata, err := datalake.GetFieldMetadata(updatedSchema, "id")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting field metadata: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nField Metadata for 'id':")
	fmt.Println(datalake.FormatFieldMetadata(metadata))

	// Add metadata to the "name" field
	fmt.Println("\nAdding metadata to 'name' field...")
	updatedSchema, err = datalake.AddFieldMetadata(updatedSchema, "name", "Customer name", []string{"searchable"}, map[string]string{"max_length": "100"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding field metadata: %v\n", err)
		os.Exit(1)
	}

	// Get all field metadata
	allMetadata, err := datalake.GetAllFieldsMetadata(updatedSchema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting all field metadata: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAll Field Metadata:")
	fmt.Println(datalake.FormatAllFieldsMetadata(allMetadata))

	// Create a record with the schema
	fmt.Println("\nCreating a record with the schema...")
	record, err := createTestRecord(updatedSchema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating test record: %v\n", err)
		os.Exit(1)
	}
	defer record.Release()

	// Print record data
	fmt.Println("\nRecord Data:")
	printRecord(record)

	fmt.Println("\nDemo completed successfully!")
}

// Helper function to print schema information
func printSchema(schema *arrow.Schema) {
	fmt.Printf("Number of fields: %d\n", schema.NumFields())
	for i, field := range schema.Fields() {
		fmt.Printf("Field %d: %s (%s)\n", i, field.Name, field.Type)
	}
}

// Helper function to print record data
func printRecord(record arrow.Record) {
	fmt.Printf("Number of rows: %d\n", record.NumRows())
	fmt.Printf("Number of columns: %d\n", record.NumCols())

	for i := int64(0); i < record.NumRows(); i++ {
		fmt.Printf("Row %d: ", i)
		for j, col := range record.Columns() {
			fieldName := record.Schema().Field(j).Name
			fmt.Printf("%s=", fieldName)

			// Handle different column types
			switch col := col.(type) {
			case *array.Int32:
				fmt.Printf("%d", col.Value(int(i)))
			case *array.String:
				fmt.Printf("\"%s\"", col.Value(int(i)))
			case *array.Float64:
				fmt.Printf("%.2f", col.Value(int(i)))
			default:
				fmt.Print("<unknown>")
			}

			if j < int(record.NumCols())-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}
}

// Helper function to create a test record with the given schema
func createTestRecord(schema *arrow.Schema) (arrow.Record, error) {
	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	// Add data based on schema fields
	for i, field := range schema.Fields() {
		switch field.Type.ID() {
		case arrow.INT32:
			builder.Field(i).(*array.Int32Builder).AppendValues([]int32{1, 2, 3}, nil)
		case arrow.STRING:
			builder.Field(i).(*array.StringBuilder).AppendValues([]string{"Alice", "Bob", "Charlie"}, nil)
		case arrow.FLOAT64:
			builder.Field(i).(*array.Float64Builder).AppendValues([]float64{10.5, 20.75, 30.25}, nil)
		default:
			return nil, fmt.Errorf("unsupported field type: %s", field.Type)
		}
	}

	return builder.NewRecord(), nil
}
