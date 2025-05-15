package dbt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ProfileResult represents the result of profiling a model
type ProfileResult struct {
	ModelName      string                 `json:"model_name"`
	TablePath      string                 `json:"table_path"`
	ProfileType    string                 `json:"profile_type"`
	RowCount       int64                  `json:"row_count"`
	ColumnCount    int                    `json:"column_count"`
	ColumnProfiles map[string]ColumnStats `json:"column_profiles"`
	Timestamp      time.Time              `json:"timestamp"`
}

// ColumnStats represents statistics for a column
type ColumnStats struct {
	DataType      string      `json:"data_type"`
	NonNullCount  int64       `json:"non_null_count"`
	NullCount     int64       `json:"null_count"`
	NullPercent   float64     `json:"null_percent"`
	Unique        int64       `json:"unique"`
	UniquePercent float64     `json:"unique_percent"`
	Min           interface{} `json:"min,omitempty"`
	Max           interface{} `json:"max,omitempty"`
	Mean          float64     `json:"mean,omitempty"`
	Stddev        float64     `json:"stddev,omitempty"`
	Quantiles     []float64   `json:"quantiles,omitempty"`
	TopValues     []ValueFreq `json:"top_values,omitempty"`
}

// ValueFreq represents a value and its frequency
type ValueFreq struct {
	Value     interface{} `json:"value"`
	Frequency int64       `json:"frequency"`
	Percent   float64     `json:"percent"`
}

// ProfileResults represents the results of profiling multiple models
type ProfileResults struct {
	Results []ProfileResult `json:"results"`
	Summary ProfileSummary  `json:"summary"`
}

// ProfileSummary provides a summary of profiling results
type ProfileSummary struct {
	TotalModels   int     `json:"total_models"`
	TotalColumns  int     `json:"total_columns"`
	TotalRows     int64   `json:"total_rows"`
	ExecutionTime float64 `json:"execution_time_seconds"`
}

// Profiler profiles dbt models
type Profiler struct {
	configPath     string
	dbtProjectPath string
	profileType    string
	config         *Config
}

// NewProfiler creates a new dbt profiler
func NewProfiler(configPath string) (*Profiler, error) {
	p := &Profiler{
		configPath:  configPath,
		profileType: "basic", // Default profile type
	}

	// Load configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	p.config = config

	// Determine dbt project path
	if config.DBTProjectPath != "" {
		p.dbtProjectPath = config.DBTProjectPath
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
		p.dbtProjectPath = dbtProjectPath
	}

	return p, nil
}

// SetProfileType sets the profile type (basic or enhanced)
func (p *Profiler) SetProfileType(profileType string) {
	// Only accept valid profile types, default to basic for invalid types
	switch profileType {
	case "basic", "enhanced":
		p.profileType = profileType
	default:
		p.profileType = "basic"
	}
}

// Profile profiles dbt models
func (p *Profiler) Profile(modelSelection []string) (*ProfileResults, error) {
	// Parse dbt manifest to get model information
	manifest, err := parseDBTManifest(p.dbtProjectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dbt manifest: %w", err)
	}

	// Select models based on selection criteria
	selectedModels, err := selectModels(manifest, modelSelection)
	if err != nil {
		return nil, fmt.Errorf("failed to select models: %w", err)
	}

	if len(selectedModels) == 0 {
		return nil, fmt.Errorf("no models selected")
	}

	// Profile each model
	startTime := time.Now()
	results := &ProfileResults{
		Results: make([]ProfileResult, 0, len(selectedModels)),
		Summary: ProfileSummary{
			TotalModels: len(selectedModels),
		},
	}

	for _, model := range selectedModels {
		// Map dbt model to Delta table
		tablePath, err := p.mapModelToTable(model)
		if err != nil {
			return nil, fmt.Errorf("failed to map model %s to table: %w", model.Name, err)
		}

		// Profile the table
		profile, err := p.profileTable(model, tablePath)
		if err != nil {
			return nil, fmt.Errorf("failed to profile model %s: %w", model.Name, err)
		}

		results.Results = append(results.Results, *profile)
		results.Summary.TotalColumns += len(profile.ColumnProfiles)
		results.Summary.TotalRows += profile.RowCount
	}

	// Calculate execution time
	results.Summary.ExecutionTime = time.Since(startTime).Seconds()

	return results, nil
}

// HasFailures returns true if any profiling operations failed
// For profiling, we consider a table with 0 rows as a failure
func (r *ProfileResults) HasFailures() bool {
	for _, result := range r.Results {
		// Consider a table with 0 rows as a failure
		if result.RowCount == 0 {
			return true
		}
		
		// Check for columns with 100% null values
		for _, stats := range result.ColumnProfiles {
			if stats.NullPercent == 100 {
				return true
			}
		}
	}
	return false
}

// mapModelToTable maps a dbt model to its Delta table path
func (p *Profiler) mapModelToTable(model *DBTModel) (string, error) {
	// Use configured mapping if available
	for _, mapping := range p.config.ModelTableMappings {
		if mapping.ModelName == model.Name {
			return mapping.TablePath, nil
		}
	}

	// Use default mapping strategy
	// This is a simplified example - in reality, you'd need to consider
	// database, schema, catalog, etc. from the dbt project configuration
	return filepath.Join(p.config.DeltaBasePath, model.Schema, model.Name), nil
}

// profileTable profiles a Delta table
func (p *Profiler) profileTable(model *DBTModel, tablePath string) (*ProfileResult, error) {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Connect to the Delta table
	// 2. Collect statistics based on the profile type (basic or enhanced)
	// 3. Process the results

	// For now, just return a mock result
	result := &ProfileResult{
		ModelName:      model.Name,
		TablePath:      tablePath,
		ProfileType:    p.profileType,
		RowCount:       1000,
		ColumnCount:    5,
		ColumnProfiles: make(map[string]ColumnStats),
		Timestamp:      time.Now(),
	}

	// Add mock column profiles
	result.ColumnProfiles["id"] = ColumnStats{
		DataType:      "integer",
		NonNullCount:  1000,
		NullCount:     0,
		NullPercent:   0,
		Unique:        1000,
		UniquePercent: 100,
		Min:           1,
		Max:           1000,
	}

	result.ColumnProfiles["name"] = ColumnStats{
		DataType:      "string",
		NonNullCount:  950,
		NullCount:     50,
		NullPercent:   5,
		Unique:        900,
		UniquePercent: 90,
	}

	// Add more columns for enhanced profile
	if p.profileType == "enhanced" {
		result.ColumnProfiles["created_at"] = ColumnStats{
			DataType:      "timestamp",
			NonNullCount:  1000,
			NullCount:     0,
			NullPercent:   0,
			Unique:        800,
			UniquePercent: 80,
		}

		result.ColumnProfiles["amount"] = ColumnStats{
			DataType:      "double",
			NonNullCount:  980,
			NullCount:     20,
			NullPercent:   2,
			Unique:        500,
			UniquePercent: 50,
			Min:           10.5,
			Max:           999.99,
			Mean:          450.75,
			Stddev:        200.5,
			Quantiles:     []float64{100.0, 300.0, 500.0, 700.0, 900.0},
		}
	}

	return result, nil
}

// parseDBTManifest parses the dbt manifest.json file
func parseDBTManifest(dbtProjectPath string) (*DBTManifest, error) {
	// Reuse the implementation from validator.go
	manifestPath := filepath.Join(dbtProjectPath, "target", "manifest.json")
	
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
func selectModels(manifest *DBTManifest, selection []string) ([]*DBTModel, error) {
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
