package dbt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test configs
	tempDir, err := os.MkdirTemp("", "dbt-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test YAML config file
	yamlConfig := `
enable_dbt_plugin: true
dbt_project_path: "/path/to/dbt/project"
delta_base_path: "/path/to/delta/tables"

model_table_mappings:
  - model_name: "customers"
    table_path: "/delta/warehouse/customers"
  - model_name: "orders"
    table_path: "/delta/warehouse/orders"

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
`
	yamlPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(yamlPath, []byte(yamlConfig), 0644); err != nil {
		t.Fatalf("Failed to write YAML config: %v", err)
	}

	// Create a test JSON config file
	jsonConfig := `{
  "enable_dbt_plugin": true,
  "dbt_project_path": "/path/to/dbt/project",
  "delta_base_path": "/path/to/delta/tables",
  "model_table_mappings": [
    {
      "model_name": "customers",
      "table_path": "/delta/warehouse/customers"
    }
  ],
  "rule_sets": [
    {
      "name": "default",
      "default": true,
      "rules": [
        {
          "name": "no_nulls_in_id",
          "type": "column_null",
          "description": "ID column should not contain nulls",
          "column": "id",
          "threshold": 0,
          "severity": "error"
        }
      ]
    }
  ]
}`
	jsonPath := filepath.Join(tempDir, "config.json")
	if err := os.WriteFile(jsonPath, []byte(jsonConfig), 0644); err != nil {
		t.Fatalf("Failed to write JSON config: %v", err)
	}

	tests := []struct {
		name       string
		configPath string
		wantErr    bool
		checkFunc  func(*testing.T, *Config)
	}{
		{
			name:       "load YAML config",
			configPath: yamlPath,
			wantErr:    false,
			checkFunc: func(t *testing.T, c *Config) {
				if !c.EnableDBTPlugin {
					t.Errorf("Expected EnableDBTPlugin to be true")
				}
				if c.DBTProjectPath != "/path/to/dbt/project" {
					t.Errorf("Expected DBTProjectPath to be '/path/to/dbt/project', got '%s'", c.DBTProjectPath)
				}
				if len(c.ModelTableMappings) != 2 {
					t.Errorf("Expected 2 model table mappings, got %d", len(c.ModelTableMappings))
				}
				if len(c.RuleSets) != 1 {
					t.Errorf("Expected 1 rule set, got %d", len(c.RuleSets))
				}
				if len(c.RuleSets[0].Rules) != 1 {
					t.Errorf("Expected 1 rule, got %d", len(c.RuleSets[0].Rules))
				}
			},
		},
		{
			name:       "load JSON config",
			configPath: jsonPath,
			wantErr:    false,
			checkFunc: func(t *testing.T, c *Config) {
				if !c.EnableDBTPlugin {
					t.Errorf("Expected EnableDBTPlugin to be true")
				}
				if c.DBTProjectPath != "/path/to/dbt/project" {
					t.Errorf("Expected DBTProjectPath to be '/path/to/dbt/project', got '%s'", c.DBTProjectPath)
				}
				if len(c.ModelTableMappings) != 1 {
					t.Errorf("Expected 1 model table mapping, got %d", len(c.ModelTableMappings))
				}
			},
		},
		{
			name:       "non-existent config",
			configPath: filepath.Join(tempDir, "nonexistent.yaml"),
			wantErr:    true,
		},
		{
			name:       "default config when no path provided",
			configPath: "",
			wantErr:    false,
			checkFunc: func(t *testing.T, c *Config) {
				if c.DeltaBasePath != "/delta" {
					t.Errorf("Expected default DeltaBasePath to be '/delta', got '%s'", c.DeltaBasePath)
				}
				if len(c.RuleSets) != 1 || !c.RuleSets[0].Default {
					t.Errorf("Expected default rule set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfig(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && tt.checkFunc != nil {
				tt.checkFunc(t, config)
			}
		})
	}
}
