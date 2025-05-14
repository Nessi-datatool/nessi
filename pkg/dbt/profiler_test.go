package dbt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestProfiler(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-profiler-test")
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

profile_settings:
  include_column_stats: true
  include_data_types: true
  include_sample_data: true
  include_null_counts: true
  include_unique_counts: true
  sample_size: 1000
`

	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test NewProfiler
	t.Run("NewProfiler", func(t *testing.T) {
		profiler, err := NewProfiler(configPath)
		if err != nil {
			t.Fatalf("NewProfiler() error = %v", err)
		}

		if profiler.configPath != configPath {
			t.Errorf("Expected configPath to be %s, got %s", configPath, profiler.configPath)
		}

		if profiler.dbtProjectPath != dbtProjectDir {
			t.Errorf("Expected dbtProjectPath to be %s, got %s", dbtProjectDir, profiler.dbtProjectPath)
		}
	})

	// Test Profile
	t.Run("Profile", func(t *testing.T) {
		profiler, err := NewProfiler(configPath)
		if err != nil {
			t.Fatalf("NewProfiler() error = %v", err)
		}

		// Test with no model selection (should select all models)
		results, err := profiler.Profile([]string{})
		if err != nil {
			t.Fatalf("Profile() error = %v", err)
		}

		if results.Summary.TotalModels != 2 {
			t.Errorf("Expected 2 models, got %d", results.Summary.TotalModels)
		}

		// Test with specific model selection
		results, err = profiler.Profile([]string{"model1"})
		if err != nil {
			t.Fatalf("Profile() error = %v", err)
		}

		if results.Summary.TotalModels != 1 {
			t.Errorf("Expected 1 model, got %d", results.Summary.TotalModels)
		}

		// Test with tag selection
		results, err = profiler.Profile([]string{"tag:daily"})
		if err != nil {
			t.Fatalf("Profile() error = %v", err)
		}

		if results.Summary.TotalModels != 1 {
			t.Errorf("Expected 1 model, got %d", results.Summary.TotalModels)
		}
	})

	// Test mapModelToTable
	t.Run("mapModelToTable", func(t *testing.T) {
		profiler, err := NewProfiler(configPath)
		if err != nil {
			t.Fatalf("NewProfiler() error = %v", err)
		}

		// Parse manifest manually since we can't access the private method
		manifestPath := filepath.Join(dbtProjectDir, "target", "manifest.json")
		manifestData, err := os.ReadFile(manifestPath)
		if err != nil {
			t.Fatalf("Failed to read manifest file: %v", err)
		}
		
		var manifest DBTManifest
		if err := json.Unmarshal(manifestData, &manifest); err != nil {
			t.Fatalf("Failed to parse manifest file: %v", err)
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
		tablePath, err := profiler.mapModelToTable(model1)
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
		tablePath, err = profiler.mapModelToTable(model3)
		if err != nil {
			t.Fatalf("mapModelToTable() error = %v", err)
		}

		expectedPath := filepath.Join("/delta", "test", "model3")
		if tablePath != expectedPath {
			t.Errorf("Expected table path to be '%s', got '%s'", expectedPath, tablePath)
		}
	})

	// Test Profile function instead of private profileTable method
	t.Run("Profile", func(t *testing.T) {
		profiler, err := NewProfiler(configPath)
		if err != nil {
			t.Fatalf("NewProfiler() error = %v", err)
		}

		// Profile a specific model
		results, err := profiler.Profile([]string{"model1"})
		if err != nil {
			t.Fatalf("Profile() error = %v", err)
		}

		// Check that we got results
		if len(results.Results) == 0 {
			t.Errorf("Expected profile results, got none")
		}

		// Check that the summary is populated
		if results.Summary.TotalModels != 1 {
			t.Errorf("Expected 1 model in summary, got %d", results.Summary.TotalModels)
		}

		// Check the first result
		if len(results.Results) > 0 {
			profile := results.Results[0]
			if profile.ModelName != "model1" {
				t.Errorf("Expected model name to be 'model1', got '%s'", profile.ModelName)
			}

			if profile.ColumnCount < 1 {
				t.Errorf("Expected at least 1 column, got %d", profile.ColumnCount)
			}

			if len(profile.ColumnProfiles) < 1 {
				t.Errorf("Expected at least 1 column profile, got %d", len(profile.ColumnProfiles))
			}
		}
	})
}
