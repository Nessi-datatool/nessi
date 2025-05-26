package integrations

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi/e2e-tests/testutil"
)

func TestSchemaManagementComprehensive(t *testing.T) {
	// Skip this test if we're not running with --enable flag
	if os.Getenv("NESSI_E2E_ENABLED") != "true" {
		t.Skip("Skipping Comprehensive Schema Management test. Use --enable flag to run this test.")
	}

	// Create a temporary directory for the test
	tempDir := testutil.CreateTempDir(t)

	// Copy sample Delta Lake tables to the temp directory
	sampleDataDir := filepath.Join(tempDir, "delta_tables")
	require.NoError(t, os.MkdirAll(sampleDataDir, 0755))
	testutil.CopyTestData(t, "delta_tables", sampleDataDir)

	// Test table path
	tablePath := filepath.Join(sampleDataDir, "sample_table")

	// Step 1: Extract schema from the table
	schemaPath := filepath.Join(tempDir, "extracted_schema.json")
	stdout, _ := testutil.AssertCommandSuccess(t, "schema", "extract", "--path", tablePath, "--output", schemaPath)
	testutil.AssertOutputContains(t, stdout, "Schema extracted successfully")

	// Verify the schema file exists
	_, err := os.Stat(schemaPath)
	require.NoError(t, err, "Schema file should exist")

	// Verify the schema file contains valid JSON
	schemaData, err := os.ReadFile(schemaPath)
	require.NoError(t, err, "Should be able to read schema file")

	var schema map[string]interface{}
	err = json.Unmarshal(schemaData, &schema)
	require.NoError(t, err, "Schema should be valid JSON")

	// Verify schema structure
	require.Equal(t, "struct", schema["type"], "Schema type should be 'struct'")
	fields, ok := schema["fields"].([]interface{})
	require.True(t, ok, "Schema should have 'fields' array")
	require.Greater(t, len(fields), 0, "Schema should have at least one field")

	// Step 2: Validate the table against the extracted schema
	stdout, _ = testutil.AssertCommandSuccess(t, "schema", "validate", "--path", tablePath, "--schema", schemaPath)
	testutil.AssertOutputContains(t, stdout, "Schema validation successful")
	testutil.AssertOutputContains(t, stdout, "Fields validated:")

	// Step 3: Create an invalid schema file
	invalidSchemaPath := filepath.Join(tempDir, "invalid_schema.json")
	invalidSchema := `{"type": "invalid", "fields": []}`
	require.NoError(t, os.WriteFile(invalidSchemaPath, []byte(invalidSchema), 0644))

	// Validate against the invalid schema (should fail)
	_, stderr := testutil.AssertCommandFailure(t, "schema", "validate", "--path", tablePath, "--schema", invalidSchemaPath)
	testutil.AssertOutputContains(t, stderr, "invalid schema")

	// Step 4: Create a malformed JSON schema file
	malformedSchemaPath := filepath.Join(tempDir, "malformed_schema.json")
	malformedSchema := `{"type": "struct", "fields": [{"name": "id", "type": "integer"`
	require.NoError(t, os.WriteFile(malformedSchemaPath, []byte(malformedSchema), 0644))

	// Validate against the malformed schema (should fail)
	_, stderr = testutil.AssertCommandFailure(t, "schema", "validate", "--path", tablePath, "--schema", malformedSchemaPath)
	testutil.AssertOutputContains(t, stderr, "Invalid schema format")

	// Step 5: Create a schema with missing required fields
	missingFieldsSchemaPath := filepath.Join(tempDir, "missing_fields_schema.json")
	missingFieldsSchema := `{"type": "struct", "fields": [{"name": "id"}]}`
	require.NoError(t, os.WriteFile(missingFieldsSchemaPath, []byte(missingFieldsSchema), 0644))

	// Validate against the schema with missing fields (should fail)
	_, stderr = testutil.AssertCommandFailure(t, "schema", "validate", "--path", tablePath, "--schema", missingFieldsSchemaPath)
	testutil.AssertOutputContains(t, stderr, "Field missing 'type' property")

	// Step 6: Test with non-existent paths
	_, stderr = testutil.AssertCommandFailure(t, "schema", "extract", "--path", filepath.Join(tempDir, "non_existent_table"), "--output", filepath.Join(tempDir, "output.json"))
	testutil.AssertOutputContains(t, stderr, "does not exist")

	_, stderr = testutil.AssertCommandFailure(t, "schema", "validate", "--path", tablePath, "--schema", filepath.Join(tempDir, "non_existent_schema.json"))
	testutil.AssertOutputContains(t, stderr, "does not exist")

	// Step 7: Test with strict validation flag
	stdout, _ = testutil.AssertCommandSuccess(t, "schema", "validate", "--path", tablePath, "--schema", schemaPath, "--strict")
	testutil.AssertOutputContains(t, stdout, "Schema validation successful")
}
