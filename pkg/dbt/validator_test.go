package dbt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestValidator(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-validator-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock dbt project directory
	dbtProjectDir := filepath.Join(tempDir, "dbt_project")
	if err := os.MkdirAll(filepath.Join(dbtProjectDir, "target"), 0755); err != nil {
		t.Fatalf("Failed to create dbt project directory: %v", err)
	}

	// Create a mock manifest.json file
	manifest := DBTManifest{
		Nodes: map[string]*DBTModel{
			"model.test.model1": {
				Name:         "model1",
				Schema:       "test",
				ResourceType: "model",
				Tags:         []string{"daily", "important"},
				Config: DBTConfig{
					Materialized: "table",
					Schema:       "test",
				},
			},
			"model.test.model2": {
				Name:         "model2",
				Schema:       "test",
				ResourceType: "model",
				Tags:         []string{"weekly"},
				Config: DBTConfig{
					Materialized: "view",
					Schema:       "test",
				},
			},
		},
		Metadata: DBTMetadata{
			ProjectName: "test_project",
			Version:     "1.0.0",
		},
	}

	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal manifest: %v", err)
	}

	manifestPath := filepath.Join(dbtProjectDir, "target", "manifest.json")
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		t.Fatalf("Failed to write manifest file: %v", err)
	}

	// Create a test config file
	configData := `
dbt_project_path: "` + dbtProjectDir + `"
delta_base_path: "/delta"

model_table_mappings:
  - model_name: "model1"
    table_path: "/delta/test/model1"
  - model_name: "model2"
    table_path: "/delta/test/model2"

rule_sets:
  - name: "default"
    default: true
    rules:
      - name: "no_nulls_in_id"
        type: "column_null"
        description: "ID column should not contain nulls"
        column: "id"
        threshold: 0
        severity: "error"
  
  - name: "daily_models"
    tag: "daily"
    rules:
      - name: "freshness_check"
        type: "freshness"
        description: "Data should be fresh"
        threshold: 24
        severity: "warning"
`

	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test NewValidator
	t.Run("NewValidator", func(t *testing.T) {
		validator, err := NewValidator(configPath)
		if err != nil {
			t.Fatalf("NewValidator() error = %v", err)
		}

		if validator.configPath != configPath {
			t.Errorf("Expected configPath to be %s, got %s", configPath, validator.configPath)
		}

		if validator.dbtProjectPath != dbtProjectDir {
			t.Errorf("Expected dbtProjectPath to be %s, got %s", dbtProjectDir, validator.dbtProjectPath)
		}
	})

	// Test Validate
	t.Run("Validate", func(t *testing.T) {
		validator, err := NewValidator(configPath)
		if err != nil {
			t.Fatalf("NewValidator() error = %v", err)
		}

		// Test with no model selection (should select all models)
		results, err := validator.Validate([]string{})
		if err != nil {
			t.Fatalf("Validate() error = %v", err)
		}

		if results.Summary.TotalModels != 2 {
			t.Errorf("Expected 2 models, got %d", results.Summary.TotalModels)
		}

		// Test with specific model selection
		results, err = validator.Validate([]string{"model1"})
		if err != nil {
			t.Fatalf("Validate() error = %v", err)
		}

		if results.Summary.TotalModels != 1 {
			t.Errorf("Expected 1 model, got %d", results.Summary.TotalModels)
		}

		// Test with tag selection
		validator.SetIncludeDBTTests(true)
		validator.SetLineageAware(true)

		results, err = validator.Validate([]string{"tag:daily"})
		if err != nil {
			t.Fatalf("Validate() error = %v", err)
		}

		if results.Summary.TotalModels != 1 {
			t.Errorf("Expected 1 model, got %d", results.Summary.TotalModels)
		}
	})

	// Test getRulesForModel
	t.Run("getRulesForModel", func(t *testing.T) {
		validator, err := NewValidator(configPath)
		if err != nil {
			t.Fatalf("NewValidator() error = %v", err)
		}

		// Get manifest
		manifest, err := validator.parseDBTManifest()
		if err != nil {
			t.Fatalf("parseDBTManifest() error = %v", err)
		}

		// Get model1 (has tag "daily")
		var model1 *DBTModel
		for _, node := range manifest.Nodes {
			if node.Name == "model1" {
				model1 = node
				break
			}
		}

		if model1 == nil {
			t.Fatalf("Failed to find model1 in manifest")
		}

		// Get rules for model1 (should get rules from "daily_models" rule set)
		rules, err := validator.getRulesForModel(model1)
		if err != nil {
			t.Fatalf("getRulesForModel() error = %v", err)
		}

		if len(rules) != 1 {
			t.Errorf("Expected 1 rule for model1, got %d", len(rules))
		}

		if rules[0].Name != "freshness_check" {
			t.Errorf("Expected rule name to be 'freshness_check', got '%s'", rules[0].Name)
		}

		// Get model2 (has no matching tag)
		var model2 *DBTModel
		for _, node := range manifest.Nodes {
			if node.Name == "model2" {
				model2 = node
				break
			}
		}

		if model2 == nil {
			t.Fatalf("Failed to find model2 in manifest")
		}

		// Get rules for model2 (should get rules from default rule set)
		rules, err = validator.getRulesForModel(model2)
		if err != nil {
			t.Fatalf("getRulesForModel() error = %v", err)
		}

		if len(rules) != 1 {
			t.Errorf("Expected 1 rule for model2, got %d", len(rules))
		}

		if rules[0].Name != "no_nulls_in_id" {
			t.Errorf("Expected rule name to be 'no_nulls_in_id', got '%s'", rules[0].Name)
		}
	})

	// Test mapModelToTable
	t.Run("mapModelToTable", func(t *testing.T) {
		validator, err := NewValidator(configPath)
		if err != nil {
			t.Fatalf("NewValidator() error = %v", err)
		}

		// Get manifest
		manifest, err := validator.parseDBTManifest()
		if err != nil {
			t.Fatalf("parseDBTManifest() error = %v", err)
		}

		// Get model1
		var model1 *DBTModel
		for _, node := range manifest.Nodes {
			if node.Name == "model1" {
				model1 = node
				break
			}
		}

		if model1 == nil {
			t.Fatalf("Failed to find model1 in manifest")
		}

		// Map model1 to table (should use mapping from config)
		tablePath, err := validator.mapModelToTable(model1)
		if err != nil {
			t.Fatalf("mapModelToTable() error = %v", err)
		}

		if tablePath != "/delta/test/model1" {
			t.Errorf("Expected table path to be '/delta/test/model1', got '%s'", tablePath)
		}

		// Create a model not in the mappings
		model3 := &DBTModel{
			Name:   "model3",
			Schema: "test",
		}

		// Map model3 to table (should use default mapping strategy)
		tablePath, err = validator.mapModelToTable(model3)
		if err != nil {
			t.Fatalf("mapModelToTable() error = %v", err)
		}

		expectedPath := filepath.Join("/delta", "test", "model3")
		if tablePath != expectedPath {
			t.Errorf("Expected table path to be '%s', got '%s'", expectedPath, tablePath)
		}
	})
}
