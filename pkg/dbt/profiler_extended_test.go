package dbt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfilerSetProfileType(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-profiler-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test manifest file
	manifestData := `{
		"nodes": {
			"model.test_project.model1": {
				"name": "model1",
				"resource_type": "model",
				"tags": ["daily"],
				"config": {
					"tags": ["daily"]
				}
			}
		}
	}`
	
	manifestDir := filepath.Join(tempDir, "target")
	err = os.MkdirAll(manifestDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create manifest directory: %v", err)
	}
	
	manifestPath := filepath.Join(manifestDir, "manifest.json")
	err = os.WriteFile(manifestPath, []byte(manifestData), 0644)
	if err != nil {
		t.Fatalf("Failed to write manifest file: %v", err)
	}

	// Create a test config file
	configData := `
enable_dbt_plugin: true
dbt_project_path: "` + tempDir + `"
delta_base_path: "/delta"

model_mappings:
  - model: "model1"
    table_path: "/delta/test/model1"
`
	
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create a profiler
	profiler, err := NewProfiler(configPath)
	if err != nil {
		t.Fatalf("NewProfiler() error = %v", err)
	}

	// Test SetProfileType
	profiler.SetProfileType("enhanced")
	if profiler.profileType != "enhanced" {
		t.Errorf("SetProfileType() failed to set profile type, got %s", profiler.profileType)
	}

	profiler.SetProfileType("basic")
	if profiler.profileType != "basic" {
		t.Errorf("SetProfileType() failed to set profile type, got %s", profiler.profileType)
	}

	// Test with invalid profile type (should default to "basic")
	profiler.SetProfileType("invalid")
	if profiler.profileType != "basic" {
		t.Errorf("SetProfileType() with invalid type should default to 'basic', got %s", profiler.profileType)
	}
}

func TestProfilerHasFailures(t *testing.T) {
	// Test with no failures
	results := &ProfileResults{
		Results: []ProfileResult{
			{
				ModelName:   "model1",
				TablePath:   "/delta/test/model1",
				ProfileType: "basic",
				RowCount:    1000,
				ColumnCount: 5,
			},
		},
		Summary: ProfileSummary{
			TotalModels:   1,
			TotalColumns:  5,
			TotalRows:     1000,
			ExecutionTime: 1.5,
		},
	}

	if results.HasFailures() {
		t.Errorf("HasFailures() should return false for profile results with no failures")
	}

	// Test with failures (in this case, we'll consider a table with 0 rows as a failure)
	results.Results[0].RowCount = 0
	
	if !results.HasFailures() {
		t.Errorf("HasFailures() should return true for profile results with failures")
	}
}

func TestProfilerEdgeCases(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-profiler-edge-cases")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file with no model mappings
	configData := `
enable_dbt_plugin: true
dbt_project_path: "` + tempDir + `"
delta_base_path: "/delta"
`
	
	configPath := filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create a test manifest file with no models
	manifestData := `{
		"nodes": {}
	}`
	
	manifestDir := filepath.Join(tempDir, "target")
	err = os.MkdirAll(manifestDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create manifest directory: %v", err)
	}
	
	manifestPath := filepath.Join(manifestDir, "manifest.json")
	err = os.WriteFile(manifestPath, []byte(manifestData), 0644)
	if err != nil {
		t.Fatalf("Failed to write manifest file: %v", err)
	}

	// Create a profiler
	profiler, err := NewProfiler(configPath)
	if err != nil {
		t.Fatalf("NewProfiler() error = %v", err)
	}

	// Test Profile with no models
	_, err = profiler.Profile([]string{})
	if err == nil {
		t.Errorf("Profile() with no models should return an error")
	}

	// Test Profile with non-existent model
	_, err = profiler.Profile([]string{"non_existent_model"})
	if err == nil {
		t.Errorf("Profile() with non-existent model should return an error")
	}

	// Test Profile with invalid tag
	_, err = profiler.Profile([]string{"tag:non_existent_tag"})
	if err == nil {
		t.Errorf("Profile() with invalid tag should return an error")
	}
}
