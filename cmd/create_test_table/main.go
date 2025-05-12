package main

import (
	"fmt"
	"os"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/array"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/nessi-dev/nessi-dev/pkg/datalake"
)

// CreateTestTable creates a sample Delta table for testing
func CreateTestTable() error {
	// Define schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
			{Name: "value", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	// Create test data
	pool := memory.NewGoAllocator()
	builder := array.NewRecordBuilder(pool, schema)
	defer builder.Release()

	// Add data
	idBuilder := builder.Field(0).(*array.Int32Builder)
	nameBuilder := builder.Field(1).(*array.StringBuilder)
	valueBuilder := builder.Field(2).(*array.Float64Builder)

	// Add 10 rows of data
	for i := 0; i < 10; i++ {
		idBuilder.Append(int32(i + 1))
		nameBuilder.Append(fmt.Sprintf("Item %d", i+1))
		valueBuilder.Append(float64(i) * 1.5)
	}

	record := builder.NewRecord()
	defer record.Release()

	// Create Delta table
	tablePath := "test_data/sample_table"
	
	// Remove existing table if it exists
	os.RemoveAll(tablePath)
	
	// Ensure directories exist
	if err := os.MkdirAll(tablePath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create transaction log directory
	logDir := fmt.Sprintf("%s/_delta_log", tablePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create transaction log directory: %w", err)
	}

	// Create a metadata manager to manually set up the table
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

	// Write initial metadata
	if err := metadataManager.WriteTableMetadata(table); err != nil {
		return fmt.Errorf("failed to write initial metadata: %w", err)
	}

	// Create writer
	writer, err := datalake.NewWriter(tablePath, schema)
	if err != nil {
		return fmt.Errorf("failed to create writer: %w", err)
	}

	// Write data
	if err := writer.WritePartition("data", record); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	// Commit changes
	if err := writer.Commit(); err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	fmt.Printf("Created test Delta table at %s\n", tablePath)
	return nil
}

func main() {
	if err := CreateTestTable(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
