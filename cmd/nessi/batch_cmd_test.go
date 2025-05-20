package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBatchConfig(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "batch_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test JSON config
	jsonConfig := `{
		"name": "Test Batch",
		"description": "Test batch job",
		"tables": ["/path/to/table1", "/path/to/table2"],
		"rules": "default.json",
		"output": "results",
		"format": "json",
		"tags": ["test", "daily"],
		"params": {
			"threshold": "0.9",
			"verbose": "true"
		}
	}`

	jsonFile := filepath.Join(tempDir, "test_batch.json")
	err = os.WriteFile(jsonFile, []byte(jsonConfig), 0644)
	require.NoError(t, err)

	// Load the JSON config
	config, err := loadBatchConfig(jsonFile)
	require.NoError(t, err)
	assert.Equal(t, "Test Batch", config.Name)
	assert.Equal(t, "Test batch job", config.Description)
	assert.Equal(t, 2, len(config.Tables))
	assert.Equal(t, "/path/to/table1", config.Tables[0])
	assert.Equal(t, "/path/to/table2", config.Tables[1])
	assert.Equal(t, "default.json", config.Rules)
	assert.Equal(t, "results", config.Output)
	assert.Equal(t, "json", config.Format)
	assert.Equal(t, 2, len(config.Tags))
	assert.Equal(t, "test", config.Tags[0])
	assert.Equal(t, "daily", config.Tags[1])
	assert.Equal(t, 2, len(config.Params))
	assert.Equal(t, "0.9", config.Params["threshold"])
	assert.Equal(t, "true", config.Params["verbose"])

	// Test YAML config
	yamlConfig := `
name: Test Batch YAML
description: Test batch job in YAML
tables:
  - /path/to/table1
  - /path/to/table2
rules: default.yaml
output: results
format: yaml
tags:
  - test
  - yaml
params:
  threshold: "0.8"
  verbose: "false"
`

	yamlFile := filepath.Join(tempDir, "test_batch.yaml")
	err = os.WriteFile(yamlFile, []byte(yamlConfig), 0644)
	require.NoError(t, err)

	// Load the YAML config
	config, err = loadBatchConfig(yamlFile)
	require.NoError(t, err)
	assert.Equal(t, "Test Batch YAML", config.Name)
	assert.Equal(t, "Test batch job in YAML", config.Description)
	assert.Equal(t, 2, len(config.Tables))
	assert.Equal(t, "/path/to/table1", config.Tables[0])
	assert.Equal(t, "/path/to/table2", config.Tables[1])
	assert.Equal(t, "default.yaml", config.Rules)
	assert.Equal(t, "results", config.Output)
	assert.Equal(t, "yaml", config.Format)
	assert.Equal(t, 2, len(config.Tags))
	assert.Equal(t, "test", config.Tags[0])
	assert.Equal(t, "yaml", config.Tags[1])
	assert.Equal(t, 2, len(config.Params))
	assert.Equal(t, "0.8", config.Params["threshold"])
	assert.Equal(t, "false", config.Params["verbose"])

	// Test invalid file format
	invalidFile := filepath.Join(tempDir, "test_batch.txt")
	err = os.WriteFile(invalidFile, []byte("invalid"), 0644)
	require.NoError(t, err)

	_, err = loadBatchConfig(invalidFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported config file format")

	// Test invalid JSON
	invalidJSON := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidJSON, []byte("{invalid json}"), 0644)
	require.NoError(t, err)

	_, err = loadBatchConfig(invalidJSON)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON config")

	// Test invalid YAML
	invalidYAML := filepath.Join(tempDir, "invalid.yaml")
	err = os.WriteFile(invalidYAML, []byte("invalid: yaml: : : "), 0644)
	require.NoError(t, err)

	_, err = loadBatchConfig(invalidYAML)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML config")
}

func TestBatchResultSerialization(t *testing.T) {
	// Create a batch result
	result := &BatchResult{
		BatchName:    "Test Batch",
		StartTime:    parseTime(t, "2025-05-13T12:00:00Z"),
		EndTime:      parseTime(t, "2025-05-13T12:05:00Z"),
		Duration:     "5m0s",
		TablesTotal:  5,
		TablesPassed: 4,
		TablesFailed: 1,
		Results:      []string{"result1.json", "result2.json"},
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(result)
	require.NoError(t, err)

	// Deserialize from JSON
	var deserializedResult BatchResult
	err = json.Unmarshal(jsonData, &deserializedResult)
	require.NoError(t, err)

	// Verify deserialized result
	assert.Equal(t, result.BatchName, deserializedResult.BatchName)
	assert.Equal(t, result.StartTime.UTC(), deserializedResult.StartTime.UTC())
	assert.Equal(t, result.EndTime.UTC(), deserializedResult.EndTime.UTC())
	assert.Equal(t, result.Duration, deserializedResult.Duration)
	assert.Equal(t, result.TablesTotal, deserializedResult.TablesTotal)
	assert.Equal(t, result.TablesPassed, deserializedResult.TablesPassed)
	assert.Equal(t, result.TablesFailed, deserializedResult.TablesFailed)
	assert.Equal(t, result.Results, deserializedResult.Results)
}

// Helper function to parse time
func parseTime(t *testing.T, timeStr string) (result time.Time) {
	result, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err)
	return
}
