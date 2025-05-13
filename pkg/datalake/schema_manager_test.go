package datalake

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSchemaManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "schema_manager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a schema manager
	sm := NewSchemaManager(tempDir)

	// Test schema initialization
	t.Run("InitializeSchema", func(t *testing.T) {
		// Create a test schema
		schema := &Schema{
			Fields: []Field{
				{Name: "id", Type: "integer", Nullable: false},
				{Name: "name", Type: "string", Nullable: true},
				{Name: "created_at", Type: "timestamp", Nullable: false},
			},
		}

		// Initialize schema
		partitionBy := []string{"created_at"}
		zOrderBy := []string{"id"}
		schemaVersion, err := sm.InitializeSchema(schema, partitionBy, zOrderBy)
		if err != nil {
			t.Fatalf("Failed to initialize schema: %v", err)
		}

		// Check schema version
		if schemaVersion.Version != 1 {
			t.Errorf("Expected version 1, got %d", schemaVersion.Version)
		}

		// Check partitioning
		if len(schemaVersion.PartitionBy) != 1 || schemaVersion.PartitionBy[0] != "created_at" {
			t.Errorf("Unexpected partition by: %v", schemaVersion.PartitionBy)
		}

		// Check Z-ordering
		if len(schemaVersion.ZOrderBy) != 1 || schemaVersion.ZOrderBy[0] != "id" {
			t.Errorf("Unexpected Z-order by: %v", schemaVersion.ZOrderBy)
		}

		// Check schema fields
		if len(schemaVersion.Schema.Fields) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(schemaVersion.Schema.Fields))
		}
	})

	// Test getting current schema
	t.Run("GetCurrentSchema", func(t *testing.T) {
		schemaVersion, err := sm.GetCurrentSchema()
		if err != nil {
			t.Fatalf("Failed to get current schema: %v", err)
		}

		// Check schema version
		if schemaVersion.Version != 1 {
			t.Errorf("Expected version 1, got %d", schemaVersion.Version)
		}

		// Check schema fields
		if len(schemaVersion.Schema.Fields) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(schemaVersion.Schema.Fields))
		}
	})

	// Test updating schema
	t.Run("UpdateSchema", func(t *testing.T) {
		// Create a new schema with changes
		newSchema := &Schema{
			Fields: []Field{
				{Name: "id", Type: "integer", Nullable: false},
				{Name: "name", Type: "string", Nullable: true},
				{Name: "created_at", Type: "timestamp", Nullable: false},
				{Name: "email", Type: "string", Nullable: true}, // Added field
			},
		}

		// Update schema
		commitInfo := map[string]string{"action": "add_email_field"}
		schemaVersion, err := sm.UpdateSchema(newSchema, commitInfo)
		if err != nil {
			t.Fatalf("Failed to update schema: %v", err)
		}

		// Check schema version
		if schemaVersion.Version != 2 {
			t.Errorf("Expected version 2, got %d", schemaVersion.Version)
		}

		// Check schema fields
		if len(schemaVersion.Schema.Fields) != 4 {
			t.Errorf("Expected 4 fields, got %d", len(schemaVersion.Schema.Fields))
		}

		// Check changes
		if len(schemaVersion.Changes) != 1 {
			t.Errorf("Expected 1 change, got %d", len(schemaVersion.Changes))
		}

		if schemaVersion.Changes[0].Type != "add" || schemaVersion.Changes[0].Field != "email" {
			t.Errorf("Unexpected change: %v", schemaVersion.Changes[0])
		}
	})

	// Test updating partitioning
	t.Run("UpdatePartitioning", func(t *testing.T) {
		// Update partitioning
		partitionBy := []string{"created_at", "name"}
		zOrderBy := []string{"id", "email"}
		commitInfo := map[string]string{"action": "update_partitioning"}
		schemaVersion, err := sm.UpdatePartitioning(partitionBy, zOrderBy, commitInfo)
		if err != nil {
			t.Fatalf("Failed to update partitioning: %v", err)
		}

		// Check schema version
		if schemaVersion.Version != 3 {
			t.Errorf("Expected version 3, got %d", schemaVersion.Version)
		}

		// Check partitioning
		if len(schemaVersion.PartitionBy) != 2 || schemaVersion.PartitionBy[0] != "created_at" || schemaVersion.PartitionBy[1] != "name" {
			t.Errorf("Unexpected partition by: %v", schemaVersion.PartitionBy)
		}

		// Check Z-ordering
		if len(schemaVersion.ZOrderBy) != 2 || schemaVersion.ZOrderBy[0] != "id" || schemaVersion.ZOrderBy[1] != "email" {
			t.Errorf("Unexpected Z-order by: %v", schemaVersion.ZOrderBy)
		}

		// Check changes
		if len(schemaVersion.Changes) != 2 {
			t.Errorf("Expected 2 changes, got %d", len(schemaVersion.Changes))
		}
	})

	// Test schema validation
	t.Run("ValidateSchema", func(t *testing.T) {
		// Valid data
		validData := map[string]interface{}{
			"id":         1,
			"name":       "Test User",
			"created_at": "2023-01-01T00:00:00Z",
			"email":      "test@example.com",
		}

		valid, violations, err := sm.ValidateSchema(validData)
		if err != nil {
			t.Fatalf("Failed to validate schema: %v", err)
		}

		if !valid {
			t.Errorf("Expected valid data, got violations: %v", violations)
		}

		// Invalid data (missing required field)
		invalidData := map[string]interface{}{
			"id":    1,
			"name":  "Test User",
			"email": "test@example.com",
			// missing created_at
		}

		valid, violations, err = sm.ValidateSchema(invalidData)
		if err != nil {
			t.Fatalf("Failed to validate schema: %v", err)
		}

		if valid {
			t.Errorf("Expected invalid data, got valid")
		}

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d: %v", len(violations), violations)
		}

		// Invalid data (wrong type)
		invalidTypeData := map[string]interface{}{
			"id":         "not an integer", // wrong type
			"name":       "Test User",
			"created_at": "2023-01-01T00:00:00Z",
			"email":      "test@example.com",
		}

		valid, violations, err = sm.ValidateSchema(invalidTypeData)
		if err != nil {
			t.Fatalf("Failed to validate schema: %v", err)
		}

		if valid {
			t.Errorf("Expected invalid data, got valid")
		}

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d: %v", len(violations), violations)
		}

		// Invalid data (extra field)
		extraFieldData := map[string]interface{}{
			"id":         1,
			"name":       "Test User",
			"created_at": "2023-01-01T00:00:00Z",
			"email":      "test@example.com",
			"extra":      "This field is not in the schema",
		}

		valid, violations, err = sm.ValidateSchema(extraFieldData)
		if err != nil {
			t.Fatalf("Failed to validate schema: %v", err)
		}

		if valid {
			t.Errorf("Expected invalid data, got valid")
		}

		if len(violations) != 1 {
			t.Errorf("Expected 1 violation, got %d: %v", len(violations), violations)
		}
	})

	// Test getting field metadata
	t.Run("GetFieldMetadata", func(t *testing.T) {
		metadata, err := sm.GetFieldMetadata("email")
		if err != nil {
			t.Fatalf("Failed to get field metadata: %v", err)
		}

		if metadata.Name != "email" {
			t.Errorf("Expected field name 'email', got '%s'", metadata.Name)
		}

		if metadata.Type != "string" {
			t.Errorf("Expected field type 'string', got '%s'", metadata.Type)
		}

		if !metadata.Nullable {
			t.Errorf("Expected field to be nullable")
		}

		// Test non-existent field
		_, err = sm.GetFieldMetadata("non_existent")
		if err == nil {
			t.Errorf("Expected error for non-existent field, got nil")
		}
	})

	// Test getting optimization hints
	t.Run("GetOptimizationHints", func(t *testing.T) {
		hints, err := sm.GetOptimizationHints()
		if err != nil {
			t.Fatalf("Failed to get optimization hints: %v", err)
		}

		// We have proper partitioning and Z-ordering, so we shouldn't get hints for those
		if _, exists := hints["partitioning"]; exists {
			t.Errorf("Unexpected partitioning hint: %s", hints["partitioning"])
		}

		if _, exists := hints["z_ordering"]; exists {
			t.Errorf("Unexpected Z-ordering hint: %s", hints["z_ordering"])
		}

		// But we should get hints for string fields
		if _, exists := hints["field_name"]; !exists {
			t.Errorf("Expected hint for string field 'name'")
		}

		if _, exists := hints["field_email"]; !exists {
			t.Errorf("Expected hint for string field 'email'")
		}
	})

	// Test getting schema history
	t.Run("GetSchemaHistory", func(t *testing.T) {
		history, err := sm.GetSchemaHistory()
		if err != nil {
			t.Fatalf("Failed to get schema history: %v", err)
		}

		if len(history.Schemas) != 3 {
			t.Errorf("Expected 3 schema versions, got %d", len(history.Schemas))
		}

		if history.CurrentIndex != 2 {
			t.Errorf("Expected current index 2, got %d", history.CurrentIndex)
		}

		// Check versions
		if history.Schemas[0].Version != 1 {
			t.Errorf("Expected version 1, got %d", history.Schemas[0].Version)
		}

		if history.Schemas[1].Version != 2 {
			t.Errorf("Expected version 2, got %d", history.Schemas[1].Version)
		}

		if history.Schemas[2].Version != 3 {
			t.Errorf("Expected version 3, got %d", history.Schemas[2].Version)
		}
	})

	// Test getting specific schema version
	t.Run("GetSchemaVersion", func(t *testing.T) {
		// Get version 1
		schemaVersion, err := sm.GetSchemaVersion(1)
		if err != nil {
			t.Fatalf("Failed to get schema version: %v", err)
		}

		if schemaVersion.Version != 1 {
			t.Errorf("Expected version 1, got %d", schemaVersion.Version)
		}

		if len(schemaVersion.Schema.Fields) != 3 {
			t.Errorf("Expected 3 fields, got %d", len(schemaVersion.Schema.Fields))
		}

		// Get version 2
		schemaVersion, err = sm.GetSchemaVersion(2)
		if err != nil {
			t.Fatalf("Failed to get schema version: %v", err)
		}

		if schemaVersion.Version != 2 {
			t.Errorf("Expected version 2, got %d", schemaVersion.Version)
		}

		if len(schemaVersion.Schema.Fields) != 4 {
			t.Errorf("Expected 4 fields, got %d", len(schemaVersion.Schema.Fields))
		}

		// Get non-existent version
		_, err = sm.GetSchemaVersion(99)
		if err == nil {
			t.Errorf("Expected error for non-existent version, got nil")
		}
	})
}

func TestDetectSchemaChanges(t *testing.T) {
	sm := NewSchemaManager("/tmp")

	testCases := []struct {
		name      string
		oldSchema *Schema
		newSchema *Schema
		expected  int // Number of expected changes
		types     []string
	}{
		{
			name: "No changes",
			oldSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
				},
			},
			newSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
				},
			},
			expected: 0,
			types:    []string{},
		},
		{
			name: "Add field",
			oldSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
				},
			},
			newSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
				},
			},
			expected: 1,
			types:    []string{"add"},
		},
		{
			name: "Remove field",
			oldSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
				},
			},
			newSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
				},
			},
			expected: 1,
			types:    []string{"remove"},
		},
		{
			name: "Modify field",
			oldSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
				},
			},
			newSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: false}, // Changed nullable
				},
			},
			expected: 1,
			types:    []string{"modify"},
		},
		{
			name: "Multiple changes",
			oldSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
					{Name: "old_field", Type: "string", Nullable: true},
				},
			},
			newSchema: &Schema{
				Fields: []Field{
					{Name: "id", Type: "long", Nullable: false}, // Changed type
					{Name: "name", Type: "string", Nullable: true},
					{Name: "new_field", Type: "boolean", Nullable: false}, // Added
					// old_field removed
				},
			},
			expected: 3,
			types:    []string{"modify", "add", "remove"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			changes, err := sm.detectSchemaChanges(tc.oldSchema, tc.newSchema)
			if err != nil {
				t.Fatalf("Failed to detect schema changes: %v", err)
			}

			if len(changes) != tc.expected {
				t.Errorf("Expected %d changes, got %d: %v", tc.expected, len(changes), changes)
			}

			// Check change types if expected > 0
			if tc.expected > 0 {
				typeMap := make(map[string]bool)
				for _, change := range changes {
					typeMap[change.Type] = true
				}

				for _, expectedType := range tc.types {
					if !typeMap[expectedType] {
						t.Errorf("Expected change type '%s' not found", expectedType)
					}
				}
			}
		})
	}
}

func TestValidateType(t *testing.T) {
	sm := NewSchemaManager("/tmp")

	testCases := []struct {
		value        interface{}
		expectedType string
		valid        bool
	}{
		{"test", "string", true},
		{123, "integer", true},
		{123.45, "float", true},
		{true, "boolean", true},
		{"2023-01-01", "date", true},
		{"not a number", "integer", false},
		{123, "string", false},
		{123.45, "integer", false},
		{"true", "boolean", false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v as %s", tc.value, tc.expectedType), func(t *testing.T) {
			valid := sm.validateType(tc.value, tc.expectedType)
			if valid != tc.valid {
				t.Errorf("Expected validateType(%v, %s) to be %v, got %v", tc.value, tc.expectedType, tc.valid, valid)
			}
		})
	}
}

func TestEqualStringSlices(t *testing.T) {
	testCases := []struct {
		a      []string
		b      []string
		equal  bool
	}{
		{[]string{}, []string{}, true},
		{[]string{"a"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"a"}, []string{"b"}, false},
		{[]string{"a", "b"}, []string{"b", "a"}, false}, // Order matters
		{[]string{"a"}, []string{}, false},
		{[]string{}, []string{"a"}, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v == %v", tc.a, tc.b), func(t *testing.T) {
			equal := equalStringSlices(tc.a, tc.b)
			if equal != tc.equal {
				t.Errorf("Expected equalStringSlices(%v, %v) to be %v, got %v", tc.a, tc.b, tc.equal, equal)
			}
		})
	}
}
