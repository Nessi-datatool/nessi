package dbt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatorHasFailures(t *testing.T) {
	// Test with no failures
	results := &ValidationResults{
		Results: []ValidationResult{
			{
				ModelName:    "model1",
				RuleName:     "rule1",
				Status:       "passed",
				Message:      "Rule passed",
				TablePath:    "/delta/test/model1",
				FailureCount: 0,
			},
			{
				ModelName:    "model2",
				RuleName:     "rule2",
				Status:       "passed",
				Message:      "Rule passed",
				TablePath:    "/delta/test/model2",
				FailureCount: 0,
			},
		},
		Summary: ValidationSummary{
			TotalModels:   2,
			PassedModels:  2,
			FailedModels:  0,
			TotalRules:    2,
			PassedRules:   2,
			FailedRules:   0,
			QualityScore:  100.0,
			ExecutionTime: 1.5,
		},
	}

	if results.HasFailures() {
		t.Errorf("HasFailures() should return false for validation results with no failures")
	}

	// Test with failures
	results.Results[0].Status = "failed"
	results.Results[0].FailureCount = 5
	results.Summary.PassedModels = 1
	results.Summary.FailedModels = 1
	results.Summary.PassedRules = 1
	results.Summary.FailedRules = 1
	results.Summary.QualityScore = 50.0
	
	if !results.HasFailures() {
		t.Errorf("HasFailures() should return true for validation results with failures")
	}
}

func TestFindDBTProjectPath(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "dbt-project-path-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a dbt_project.yml file in the temp directory
	dbtProjectPath := filepath.Join(tempDir, "dbt_project.yml")
	err = os.WriteFile(dbtProjectPath, []byte("name: test_project"), 0644)
	if err != nil {
		t.Fatalf("Failed to write dbt_project.yml file: %v", err)
	}

	// Create a subdirectory with another dbt_project.yml file
	subDir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	subDbtProjectPath := filepath.Join(subDir, "dbt_project.yml")
	err = os.WriteFile(subDbtProjectPath, []byte("name: sub_project"), 0644)
	if err != nil {
		t.Fatalf("Failed to write subdirectory dbt_project.yml file: %v", err)
	}

	// Test findDBTProjectPath with explicit path
	path, err := findDBTProjectPath(tempDir)
	if err != nil {
		t.Fatalf("findDBTProjectPath() error = %v", err)
	}

	if path != tempDir {
		t.Errorf("findDBTProjectPath() with explicit path returned %s, expected %s", path, tempDir)
	}

	// Test findDBTProjectPath with subdirectory path
	path, err = findDBTProjectPath(subDir)
	if err != nil {
		t.Fatalf("findDBTProjectPath() error = %v", err)
	}

	if path != subDir {
		t.Errorf("findDBTProjectPath() with subdirectory path returned %s, expected %s", path, subDir)
	}

	// Test findDBTProjectPath with non-existent path
	nonExistentPath := filepath.Join(tempDir, "non_existent")
	_, err = findDBTProjectPath(nonExistentPath)
	if err == nil {
		t.Errorf("findDBTProjectPath() with non-existent path should return an error")
	}

	// Test findDBTProjectPath with path that doesn't contain dbt_project.yml
	// Create a completely separate directory structure for this test
	isolatedTempDir, err := os.MkdirTemp("", "dbt-isolated-test")
	if err != nil {
		t.Fatalf("Failed to create isolated temp directory: %v", err)
	}
	defer os.RemoveAll(isolatedTempDir)

	// Make sure this directory is completely empty and has no parent with dbt_project.yml
	_, err = findDBTProjectPath(isolatedTempDir)
	if err == nil {
		t.Errorf("findDBTProjectPath() with path that doesn't contain dbt_project.yml should return an error")
	}

	// Test automatic discovery (current directory and parent directories)
	// This is harder to test in a unit test environment, so we'll skip it for now
}

func TestValidatorEdgeCases(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-validator-edge-cases")
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

	// Create a validator
	validator, err := NewValidator(configPath)
	if err != nil {
		t.Fatalf("NewValidator() error = %v", err)
	}

	// Test Validate with no models
	_, err = validator.Validate([]string{})
	if err == nil {
		t.Errorf("Validate() with no models should return an error")
	}

	// Test Validate with non-existent model
	_, err = validator.Validate([]string{"non_existent_model"})
	if err == nil {
		t.Errorf("Validate() with non-existent model should return an error")
	}

	// Test Validate with invalid tag
	_, err = validator.Validate([]string{"tag:non_existent_tag"})
	if err == nil {
		t.Errorf("Validate() with invalid tag should return an error")
	}

	// Test with includeDBTTests flag
	validator.SetIncludeDBTTests(true)
	if !validator.includeDBTTests {
		t.Errorf("SetIncludeDBTTests() failed to set includeDBTTests flag")
	}

	// Test with lineageAware flag
	validator.SetLineageAware(true)
	if !validator.lineageAware {
		t.Errorf("SetLineageAware() failed to set lineageAware flag")
	}
}
