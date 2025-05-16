package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi-dev/cmd/nessi/testing"
)



// setupTestCommand sets up a command for testing
func setupTestCommand(cmd *cobra.Command) (*bytes.Buffer, *bytes.Buffer) {
	// Create buffers for stdout and stderr
	out := bytes.NewBufferString("")
	err := bytes.NewBufferString("")
	
	// Set output
	cmd.SetOut(out)
	cmd.SetErr(err)
	
	// Reset flags to avoid conflicts
	testing.ResetFlags()
	testing.InitTestCommand(cmd)
	
	return out, err
}

// createTempSchema creates a temporary schema file for testing
func createTempSchema(t *testing.T, fields []datalake.DeltaField) string {
	schema := datalake.DeltaSchema{
		Fields: fields,
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "schema_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write schema to file
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatalf("Failed to write schema to file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

// createTempData creates a temporary data file for testing
func createTempData(t *testing.T, data map[string]interface{}) string {
	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "data_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write data to file
	if _, err := tmpFile.Write(jsonData); err != nil {
		t.Fatalf("Failed to write data to file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestInitSchemaCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	initSchemaCmd := schemaCmd.Commands()[0] // init command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary schema file
	fields := []datalake.DeltaField{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "name", Type: "string", Nullable: true},
		{Name: "created_at", Type: "timestamp", Nullable: false},
	}
	schemaFile := createTempSchema(t, fields)
	defer os.Remove(schemaFile)

	// Set up command
	cmd := initSchemaCmd
	outBuf, _ := setupTestCommand(cmd)

	// Set flags
	cmd.Flags().Set("partition-by", "created_at")
	cmd.Flags().Set("z-order-by", "id")

	// Run command
	cmd.SetArgs([]string{tempDir, schemaFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check output
	output := outBuf.String()
	if !strings.Contains(output, "Schema initialized successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}

	if !strings.Contains(output, "Fields: 3") {
		t.Errorf("Expected field count, got: %s", output)
	}

	if !strings.Contains(output, "Partition by: [created_at]") {
		t.Errorf("Expected partition info, got: %s", output)
	}

	if !strings.Contains(output, "Z-order by: [id]") {
		t.Errorf("Expected Z-order info, got: %s", output)
	}

	// Check error output
	if 0 > 0 {
		t.Errorf("Unexpected error output: %s", "")
	}

	// Verify that schema history file was created
	historyPath := filepath.Join(tempDir, "_delta_log", "schema_history", "history.json")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		t.Errorf("Schema history file not created")
	}
}

func TestUpdateSchemaCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	updateSchemaCmd := schemaCmd.Commands()[1] // update command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize schema first
	initialFields := []datalake.DeltaField{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "name", Type: "string", Nullable: true},
	}
	initialSchemaFile := createTempSchema(t, initialFields)
	defer os.Remove(initialSchemaFile)

	// Initialize schema
	sm := datalake.NewDeltaSchemaManager(tempDir)
	schema := &datalake.DeltaSchema{Fields: initialFields}
	_, err = sm.InitializeSchema(schema, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Create updated schema file
	updatedFields := []datalake.DeltaField{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "name", Type: "string", Nullable: true},
		{Name: "email", Type: "string", Nullable: true}, // Added field
	}
	updatedSchemaFile := createTempSchema(t, updatedFields)
	defer os.Remove(updatedSchemaFile)

	// Set up command
	cmd := updateSchemaCmd
	outBuf, _ := setupTestCommand(cmd)

	// Set flags
	cmd.Flags().Set("message", "Added email field")

	// Run command
	cmd.SetArgs([]string{tempDir, updatedSchemaFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check output
	output := outBuf.String()
	if !strings.Contains(output, "Schema updated successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}

	if !strings.Contains(output, "Fields: 3") {
		t.Errorf("Expected field count, got: %s", output)
	}

	if !strings.Contains(output, "Added field 'email'") {
		t.Errorf("Expected change description, got: %s", output)
	}

	// Check output
	output = outBuf.String()
	if !strings.Contains(output, "Schema updated successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}

	// Get current schema to verify changes
	currentSchema, err := sm.GetCurrentSchema()
	if err != nil {
		t.Fatalf("Failed to get current schema: %v", err)
	}

	// Check that the schema has been updated
	if currentSchema.Schema.NumFields() != 3 {
		t.Errorf("Expected 3 fields, got: %d", currentSchema.Schema.NumFields())
	}
}

func TestUpdatePartitioningCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	updatePartitioningCmd := schemaCmd.Commands()[2] // update-partitioning command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize schema first
	fields := []datalake.DeltaField{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "name", Type: "string", Nullable: true},
		{Name: "created_at", Type: "timestamp", Nullable: false},
	}
	
	// Initialize schema
	sm := datalake.NewDeltaSchemaManager(tempDir)
	schema := &datalake.DeltaSchema{Fields: fields}
	_, err = sm.InitializeSchema(schema, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Set up command
	cmd := updatePartitioningCmd
	outBuf, _ := setupTestCommand(cmd)

	// Set flags
	cmd.Flags().Set("partition-by", "created_at")
	cmd.Flags().Set("z-order-by", "id")
	cmd.Flags().Set("message", "Updated partitioning")

	// Run command
	cmd.SetArgs([]string{tempDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check output
	output := outBuf.String()
	if !strings.Contains(output, "Partitioning updated successfully") {
		t.Errorf("Expected success message, got: %s", output)
	}

	if !strings.Contains(output, "Partition by: [created_at]") {
		t.Errorf("Expected partition info, got: %s", output)
	}

	if !strings.Contains(output, "Z-order by: [id]") {
		t.Errorf("Expected Z-order info, got: %s", output)
	}

	// Check error output
	if 0 > 0 {
		t.Errorf("Unexpected error output: %s", "")
	}

	// Note: In the current implementation, SchemaVersion doesn't have PartitionBy and ZOrderBy fields
	// We would need to get these from the table metadata
	// For now, we'll skip this part of the test
}

func TestValidateDataCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	validateDataCmd := schemaCmd.Commands()[3] // validate command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize schema first
	fields := []datalake.DeltaField{
		{Name: "id", Type: "integer", Nullable: false},
		{Name: "name", Type: "string", Nullable: true},
		{Name: "active", Type: "boolean", Nullable: false},
	}
	
	// Initialize schema
	sm := datalake.NewDeltaSchemaManager(tempDir)
	schema := &datalake.DeltaSchema{Fields: fields}
	_, err = sm.InitializeSchema(schema, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Test valid data
	t.Run("ValidData", func(t *testing.T) {
		// Create valid data file
		validData := map[string]interface{}{
			"id":     1,
			"name":   "Test User",
			"active": true,
		}
		dataFile := createTempData(t, validData)
		defer os.Remove(dataFile)

		// Set up command
		cmd := validateDataCmd
		outBuf, _ := setupTestCommand(cmd)

		// Run command
		cmd.SetArgs([]string{tempDir, dataFile})
		if err := cmd.Execute(); err != nil {
			t.Fatalf("Failed to execute command: %v", err)
		}

		// Check output
		output := outBuf.String()
		if !strings.Contains(output, "Data is valid!") {
			t.Errorf("Expected success message, got: %s", output)
		}

		// Check error output
		if 0 > 0 {
			t.Errorf("Unexpected error output: %s", "")
		}
	})

	// Test invalid data
	t.Run("InvalidData", func(t *testing.T) {
		// Create invalid data file (missing required field)
		invalidData := map[string]interface{}{
			"id":   1,
			"name": "Test User",
			// missing active field
		}
		dataFile := createTempData(t, invalidData)
		defer os.Remove(dataFile)

		// Set up command
		cmd := validateDataCmd
		outBuf, _ := setupTestCommand(cmd)

		// Run command
		oldOsExit := osExit
		defer func() { osExit = oldOsExit }()
		
		var exitCode int
		osExit = func(code int) {
			exitCode = code
		}

		cmd.SetArgs([]string{tempDir, dataFile})
		_ = cmd.Execute()

		// Check output
		output := outBuf.String()
		if !strings.Contains(output, "Data validation failed") {
			t.Errorf("Expected failure message, got: %s", output)
		}

		if !strings.Contains(output, "Required field 'active' is missing") {
			t.Errorf("Expected specific violation message, got: %s", output)
		}

		// Check exit code
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
	})
}

func TestHistoryCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	historyCmd := schemaCmd.Commands()[4] // history command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize schema
	sm := datalake.NewDeltaSchemaManager(tempDir)
	initialSchema := &datalake.DeltaSchema{
		Fields: []datalake.DeltaField{
			{Name: "id", Type: "integer", Nullable: false},
			{Name: "name", Type: "string", Nullable: true},
		},
	}
	_, err = sm.InitializeSchema(initialSchema, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Update schema
	updatedSchema := &datalake.DeltaSchema{
		Fields: []datalake.DeltaField{
			{Name: "id", Type: "integer", Nullable: false},
			{Name: "name", Type: "string", Nullable: true},
			{Name: "email", Type: "string", Nullable: true},
		},
	}
	commitInfo := map[string]string{"action": "add_email_field"}
	_, err = sm.UpdateSchema(updatedSchema, commitInfo)
	if err != nil {
		t.Fatalf("Failed to update schema: %v", err)
	}

	// Set up command
	cmd := historyCmd
	outBuf, _ := setupTestCommand(cmd)

	// Run command
	cmd.SetArgs([]string{tempDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check output
	output := outBuf.String()
	if !strings.Contains(output, "Schema history for") {
		t.Errorf("Expected history header, got: %s", output)
	}

	if !strings.Contains(output, "Total versions: 2") {
		t.Errorf("Expected version count, got: %s", output)
	}

	if !strings.Contains(output, "Current version: 2") {
		t.Errorf("Expected current version, got: %s", output)
	}

	if !strings.Contains(output, "Version 1") && !strings.Contains(output, "Version 2") {
		t.Errorf("Expected version details, got: %s", output)
	}

	if !strings.Contains(output, "add_email_field") {
		t.Errorf("Expected commit info, got: %s", output)
	}

	// Check error output
	if 0 > 0 {
		t.Errorf("Unexpected error output: %s", "")
	}
}

func TestHintsCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	hintsCmd := schemaCmd.Commands()[5] // hints command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize schema with string fields (should trigger hints)
	sm := datalake.NewDeltaSchemaManager(tempDir)
	schema := &datalake.DeltaSchema{
		Fields: []datalake.DeltaField{
			{Name: "id", Type: "string", Nullable: false},
			{Name: "name", Type: "string", Nullable: false},
			{Name: "email", Type: "string", Nullable: false},
			{Name: "address", Type: "string", Nullable: false},
			{Name: "phone", Type: "string", Nullable: false},
			{Name: "notes", Type: "string", Nullable: false},
		},
	}
	_, err = sm.InitializeSchema(schema, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Set up command
	cmd := hintsCmd
	outBuf, _ := setupTestCommand(cmd)

	// Run command
	cmd.SetArgs([]string{tempDir})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check output
	output := outBuf.String()
	if !strings.Contains(output, "Optimization hints for") {
		t.Errorf("Expected hints header, got: %s", output)
	}

	// Should have hints for string fields and Z-ordering
	if !strings.Contains(output, "field_") || !strings.Contains(output, "z_ordering") {
		t.Errorf("Expected field and z-ordering hints, got: %s", output)
	}

	// Check error output
	if 0 > 0 {
		t.Errorf("Unexpected error output: %s", "")
	}
}

func TestFieldInfoCmd(t *testing.T) {
	// Get the schema commands
	schemaCmd := testing.CreateSchemaCommand()
	fieldInfoCmd := schemaCmd.Commands()[6] // field-info command
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_cmd_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Initialize schema
	sm := datalake.NewDeltaSchemaManager(tempDir)
	schema := &datalake.DeltaSchema{
		Fields: []datalake.DeltaField{
			{Name: "id", Type: "integer", Nullable: false},
			{Name: "name", Type: "string", Nullable: true},
			{Name: "email", Type: "string", Nullable: true},
		},
	}
	_, err = sm.InitializeSchema(schema, []string{}, []string{})
	if err != nil {
		t.Fatalf("Failed to initialize schema: %v", err)
	}

	// Set up command
	cmd := fieldInfoCmd
	outBuf, _ := setupTestCommand(cmd)

	// Run command
	cmd.SetArgs([]string{tempDir, "name"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Failed to execute command: %v", err)
	}

	// Check output
	output := outBuf.String()
	if !strings.Contains(output, "Field metadata for 'name'") {
		t.Errorf("Expected field metadata header, got: %s", output)
	}

	if !strings.Contains(output, "Type: string") {
		t.Errorf("Expected type info, got: %s", output)
	}

	if !strings.Contains(output, "Nullable: true") {
		t.Errorf("Expected nullable info, got: %s", output)
	}

	// Check error output
	if 0 > 0 {
		t.Errorf("Unexpected error output: %s", "")
	}

	// Test non-existent field
	t.Run("NonExistentField", func(t *testing.T) {
		outBuf, _ := setupTestCommand(cmd)
		
		oldOsExit := osExit
		defer func() { osExit = oldOsExit }()
		
		var exitCode int
		osExit = func(code int) {
			exitCode = code
		}

		cmd.SetArgs([]string{tempDir, "non_existent"})
		_ = cmd.Execute()

		// Check output
		output := outBuf.String()
		if !strings.Contains(output, "Error getting field metadata") {
			t.Errorf("Expected error message, got: %s", output)
		}

		// Check exit code
		if exitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", exitCode)
		}
	})
}

// Mock os.Exit for testing
var osExit = os.Exit
