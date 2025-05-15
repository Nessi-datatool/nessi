package dbt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidationResult represents the result of a validation rule
type ValidationResult struct {
	ModelName    string `json:"model_name"`
	RuleName     string `json:"rule_name"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	TablePath    string `json:"table_path"`
	FailureCount int    `json:"failure_count"`
}

// ValidationResults represents the results of validating multiple models
type ValidationResults struct {
	Results []ValidationResult `json:"results"`
	Summary ValidationSummary  `json:"summary"`
}

// ValidationSummary provides a summary of validation results
type ValidationSummary struct {
	TotalModels    int     `json:"total_models"`
	PassedModels   int     `json:"passed_models"`
	FailedModels   int     `json:"failed_models"`
	TotalRules     int     `json:"total_rules"`
	PassedRules    int     `json:"passed_rules"`
	FailedRules    int     `json:"failed_rules"`
	QualityScore   float64 `json:"quality_score"`
	ExecutionTime  float64 `json:"execution_time_seconds"`
}

// Validator validates dbt models using Nessi.dev data quality rules
type Validator struct {
	configPath      string
	dbtProjectPath  string
	includeDBTTests bool
	lineageAware    bool
	config          *Config
}

// NewValidator creates a new dbt validator
func NewValidator(configPath string) (*Validator, error) {
	v := &Validator{
		configPath: configPath,
	}

	// Load configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	v.config = config

	// Determine dbt project path
	if config.DBTProjectPath != "" {
		v.dbtProjectPath = config.DBTProjectPath
	} else {
		// Try to find dbt_project.yml in current directory or parent directories
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current working directory: %w", err)
		}

		dbtProjectPath, err := findDBTProjectPath(cwd)
		if err != nil {
			return nil, fmt.Errorf("failed to find dbt project: %w", err)
		}
		v.dbtProjectPath = dbtProjectPath
	}

	return v, nil
}

// SetIncludeDBTTests sets whether to include dbt test results in validation
func (v *Validator) SetIncludeDBTTests(include bool) {
	v.includeDBTTests = include
}

// SetLineageAware sets whether to enable lineage-aware validation
func (v *Validator) SetLineageAware(lineageAware bool) {
	v.lineageAware = lineageAware
}

// Validate validates dbt models using Nessi.dev data quality rules
func (v *Validator) Validate(modelSelection []string) (*ValidationResults, error) {
	// Parse dbt manifest to get model information
	manifest, err := v.parseDBTManifest()
	if err != nil {
		return nil, fmt.Errorf("failed to parse dbt manifest: %w", err)
	}

	// Select models based on selection criteria
	selectedModels, err := v.selectModels(manifest, modelSelection)
	if err != nil {
		return nil, fmt.Errorf("failed to select models: %w", err)
	}

	if len(selectedModels) == 0 {
		return nil, fmt.Errorf("no models selected")
	}

	// Execute validation rules for each model
	results := &ValidationResults{
		Results: make([]ValidationResult, 0),
		Summary: ValidationSummary{
			TotalModels: len(selectedModels),
		},
	}

	for _, model := range selectedModels {
		// Map dbt model to Delta table
		tablePath, err := v.mapModelToTable(model)
		if err != nil {
			return nil, fmt.Errorf("failed to map model %s to table: %w", model.Name, err)
		}

		// Get rules for this model
		rules, err := v.getRulesForModel(model)
		if err != nil {
			return nil, fmt.Errorf("failed to get rules for model %s: %w", model.Name, err)
		}

		// Execute rules
		for _, rule := range rules {
			result, err := v.executeRule(model, tablePath, rule)
			if err != nil {
				return nil, fmt.Errorf("failed to execute rule %s for model %s: %w", rule.Name, model.Name, err)
			}

			results.Results = append(results.Results, *result)
			results.Summary.TotalRules++

			if result.Status == "passed" {
				results.Summary.PassedRules++
			} else {
				results.Summary.FailedRules++
			}
		}

		// Include dbt test results if enabled
		if v.includeDBTTests {
			dbtResults, err := v.getDBTTestResults(model)
			if err != nil {
				return nil, fmt.Errorf("failed to get dbt test results for model %s: %w", model.Name, err)
			}

			results.Results = append(results.Results, dbtResults...)
			results.Summary.TotalRules += len(dbtResults)

			for _, result := range dbtResults {
				if result.Status == "passed" {
					results.Summary.PassedRules++
				} else {
					results.Summary.FailedRules++
				}
			}
		}
	}

	// Calculate passed/failed models and quality score
	modelStatus := make(map[string]bool)
	for _, result := range results.Results {
		if _, exists := modelStatus[result.ModelName]; !exists {
			modelStatus[result.ModelName] = true
		}
		if result.Status != "passed" {
			modelStatus[result.ModelName] = false
		}
	}

	for _, passed := range modelStatus {
		if passed {
			results.Summary.PassedModels++
		} else {
			results.Summary.FailedModels++
		}
	}

	// Calculate quality score (percentage of passed rules)
	if results.Summary.TotalRules > 0 {
		results.Summary.QualityScore = float64(results.Summary.PassedRules) / float64(results.Summary.TotalRules) * 100
	}

	return results, nil
}

// HasFailures returns true if any validation rules failed
func (r *ValidationResults) HasFailures() bool {
	return r.Summary.FailedRules > 0
}

// parseDBTManifest parses the dbt manifest.json file
func (v *Validator) parseDBTManifest() (*DBTManifest, error) {
	manifestPath := filepath.Join(v.dbtProjectPath, "target", "manifest.json")
	
	// Check if manifest file exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dbt manifest file not found at %s, run 'dbt compile' first", manifestPath)
	}

	// Read and parse manifest file
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	var manifest DBTManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest file: %w", err)
	}

	return &manifest, nil
}

// selectModels selects models based on selection criteria
func (v *Validator) selectModels(manifest *DBTManifest, selection []string) ([]*DBTModel, error) {
	// If no selection criteria provided, select all models
	if len(selection) == 0 {
		models := make([]*DBTModel, 0, len(manifest.Nodes))
		for _, node := range manifest.Nodes {
			if node.ResourceType == "model" {
				models = append(models, node)
			}
		}
		return models, nil
	}

	// Basic implementation of dbt selection syntax
	selectedModels := make([]*DBTModel, 0)
	for _, selector := range selection {
		// Check for tag selector (tag:tagname)
		if strings.HasPrefix(selector, "tag:") {
			tagName := strings.TrimPrefix(selector, "tag:")
			for _, node := range manifest.Nodes {
				if node.ResourceType == "model" {
					// Check if the model has the specified tag
					for _, tag := range node.Tags {
						if tag == tagName {
							selectedModels = append(selectedModels, node)
							break
						}
					}
				}
			}
		} else {
			// Simple name matching
			for _, node := range manifest.Nodes {
				if node.ResourceType == "model" && node.Name == selector {
					selectedModels = append(selectedModels, node)
				}
			}
		}
	}

	// If no models were selected, return an error
	if len(selectedModels) == 0 {
		return nil, fmt.Errorf("no models selected with criteria: %v", selection)
	}

	return selectedModels, nil
}

// mapModelToTable maps a dbt model to its Delta table path
func (v *Validator) mapModelToTable(model *DBTModel) (string, error) {
	// Use configured mapping if available
	for _, mapping := range v.config.ModelTableMappings {
		if mapping.ModelName == model.Name {
			return mapping.TablePath, nil
		}
	}

	// Use default mapping strategy
	// This is a simplified example - in reality, you'd need to consider
	// database, schema, catalog, etc. from the dbt project configuration
	return filepath.Join(v.config.DeltaBasePath, model.Schema, model.Name), nil
}

// getRulesForModel gets data quality rules for a model
func (v *Validator) getRulesForModel(model *DBTModel) ([]*Rule, error) {
	// Check if there are model-specific rules
	for _, ruleSet := range v.config.RuleSets {
		if ruleSet.ModelName == model.Name {
			return ruleSet.Rules, nil
		}
	}

	// Check if there are tag-based rules
	for _, tag := range model.Tags {
		for _, ruleSet := range v.config.RuleSets {
			if ruleSet.Tag == tag {
				return ruleSet.Rules, nil
			}
		}
	}

	// Use default rules if available
	for _, ruleSet := range v.config.RuleSets {
		if ruleSet.Default {
			return ruleSet.Rules, nil
		}
	}

	// No rules found
	return []*Rule{}, nil
}

// executeRule executes a data quality rule against a Delta table
func (v *Validator) executeRule(model *DBTModel, tablePath string, rule *Rule) (*ValidationResult, error) {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Connect to the Delta table
	// 2. Execute the rule's SQL or use Nessi.dev's rule engine
	// 3. Process the results

	// For now, just return a mock result
	return &ValidationResult{
		ModelName: model.Name,
		RuleName:  rule.Name,
		Status:    "passed", // or "failed"
		Message:   "Rule executed successfully",
		TablePath: tablePath,
	}, nil
}

// getDBTTestResults gets test results from dbt
func (v *Validator) getDBTTestResults(model *DBTModel) ([]ValidationResult, error) {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Parse dbt's run_results.json
	// 2. Extract test results for the given model
	// 3. Convert them to ValidationResult format

	// For now, just return an empty slice
	return []ValidationResult{}, nil
}

// findDBTProjectPath finds the dbt project path by looking for dbt_project.yml
func findDBTProjectPath(startDir string) (string, error) {
	// First check if the startDir exists
	if _, err := os.Stat(startDir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", startDir)
	}

	// Check if dbt_project.yml exists in current directory
	projectFile := filepath.Join(startDir, "dbt_project.yml")
	if _, err := os.Stat(projectFile); err == nil {
		return startDir, nil
	}

	// Check parent directories (up to 5 levels)
	dir := startDir
	for i := 0; i < 5; i++ {
		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root directory
		}
		dir = parent
		projectFile = filepath.Join(dir, "dbt_project.yml")
		if _, err := os.Stat(projectFile); err == nil {
			return dir, nil
		}
	}

	return "", fmt.Errorf("dbt_project.yml not found in directory or parent directories: %s", startDir)
}
