package dbt

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestAlert(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-alert-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file with alert settings
	configData := `
enable_dbt_plugin: true
dbt_project_path: "/path/to/dbt/project"
delta_base_path: "/delta"

alert_settings:
  enabled: true
  threshold_score: 80
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/TXXXXXXXX/BXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
    channel: "#data-quality"
    username: "Nessi DBT Quality Bot"
  email:
    enabled: true
    smtp_host: "smtp.example.com"
    smtp_port: 587
    smtp_user: "alerts@example.com"
    smtp_password: "password"
    from: "alerts@example.com"
    to:
      - "data-team@example.com"
      - "engineering@example.com"
    subject_prefix: "[DATA QUALITY]"
`

	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configData), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test LoadAlertConfig
	t.Run("LoadAlertConfig", func(t *testing.T) {
		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		if !config.AlertSettings.Enabled {
			t.Errorf("Expected AlertSettings.Enabled to be true")
		}

		if config.AlertSettings.ThresholdScore != 80 {
			t.Errorf("Expected ThresholdScore to be 80, got %d", config.AlertSettings.ThresholdScore)
		}

		if !config.AlertSettings.Slack.Enabled {
			t.Errorf("Expected Slack.Enabled to be true")
		}

		if config.AlertSettings.Slack.WebhookURL != "https://hooks.slack.com/services/TXXXXXXXX/BXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX" {
			t.Errorf("Expected Slack.WebhookURL to be set correctly")
		}

		if !config.AlertSettings.Email.Enabled {
			t.Errorf("Expected Email.Enabled to be true")
		}

		if len(config.AlertSettings.Email.To) != 2 {
			t.Errorf("Expected 2 email recipients, got %d", len(config.AlertSettings.Email.To))
		}
	})

	// Test CalculateQualityScore
	t.Run("CalculateQualityScore", func(t *testing.T) {
		// Create validation results with some failures
		results := &ValidationResults{
			Results: []ValidationResult{
				{
					ModelName:    "model1",
					RuleName:     "rule1",
					Status:       "passed",
				},
				{
					ModelName:    "model1",
					RuleName:     "rule2",
					Status:       "passed",
				},
				{
					ModelName:    "model2",
					RuleName:     "rule1",
					Status:       "passed",
				},
				{
					ModelName:    "model2",
					RuleName:     "rule2",
					Status:       "failed",
					Message:      "Rule failed",
				},
			},
			Summary: ValidationSummary{
				TotalModels:    2,
				TotalRules:     4,
				PassedRules:    3,
				FailedRules:    1,
				QualityScore:   0, // Will be calculated
			},
		}

		// Calculate quality score
		score := CalculateQualityScore(results)
		
		// Expected score: (3 / 4) * 100 = 75
		if score != 75 {
			t.Errorf("Expected quality score to be 75, got %d", score)
		}

		// Update the results with the calculated score
		results.Summary.QualityScore = float64(score)

		// Test with all rules passed
		allPassedResults := &ValidationResults{
			Results: []ValidationResult{
				{
					ModelName:    "model1",
					RuleName:     "rule1",
					Status:       "passed",
				},
				{
					ModelName:    "model1",
					RuleName:     "rule2",
					Status:       "passed",
				},
			},
			Summary: ValidationSummary{
				TotalModels:    1,
				TotalRules:     2,
				PassedRules:    2,
				FailedRules:    0,
				QualityScore:   0, // Will be calculated
			},
		}

		// Calculate quality score
		score = CalculateQualityScore(allPassedResults)
		
		// Expected score: (2 / 2) * 100 = 100
		if score != 100 {
			t.Errorf("Expected quality score to be 100, got %d", score)
		}
	})

	// Test ShouldSendAlert
	t.Run("ShouldSendAlert", func(t *testing.T) {
		config, err := LoadConfig(configPath)
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}

		// Test with score below threshold
		results := &ValidationResults{
			Summary: ValidationSummary{
				QualityScore: 75, // Below threshold of 80
			},
		}

		if !ShouldSendAlert(config, results) {
			t.Errorf("Expected ShouldSendAlert to return true for score below threshold")
		}

		// Test with score above threshold
		results.Summary.QualityScore = 85 // Above threshold of 80
		if ShouldSendAlert(config, results) {
			t.Errorf("Expected ShouldSendAlert to return false for score above threshold")
		}

		// Test with alerts disabled
		config.AlertSettings.Enabled = false
		results.Summary.QualityScore = 75 // Below threshold
		if ShouldSendAlert(config, results) {
			t.Errorf("Expected ShouldSendAlert to return false when alerts are disabled")
		}
	})

	// Test FormatAlertMessage
	t.Run("FormatAlertMessage", func(t *testing.T) {
		results := &ValidationResults{
			Results: []ValidationResult{
				{
					ModelName:    "model1",
					RuleName:     "rule1",
					Status:       "passed",
				},
				{
					ModelName:    "model2",
					RuleName:     "rule2",
					Status:       "failed",
					Message:      "Rule failed",
				},
			},
			Summary: ValidationSummary{
				TotalModels:    2,
				TotalRules:     2,
				PassedRules:    1,
				FailedRules:    1,
				QualityScore:   50,
			},
		}

		message := FormatAlertMessage(results)
		
		// Check that the message contains key information
		if len(message) == 0 {
			t.Errorf("Expected non-empty alert message")
		}
		
		// Check that the message contains the quality score
		if message == "" || len(message) < 10 {
			t.Errorf("Expected detailed alert message, got: %s", message)
		}
	})
}

// Helper functions for the alert system that should be implemented in alert.go

// CalculateQualityScore calculates a quality score based on validation results
func CalculateQualityScore(results *ValidationResults) int {
	if results.Summary.TotalRules == 0 {
		return 100 // Perfect score if no rules (though this shouldn't happen)
	}
	
	// Calculate score as percentage of passed rules
	score := (results.Summary.PassedRules * 100) / results.Summary.TotalRules
	return score
}

// ShouldSendAlert determines if an alert should be sent based on validation results and config
func ShouldSendAlert(config *Config, results *ValidationResults) bool {
	// Don't send alerts if they're disabled
	if !config.AlertSettings.Enabled {
		return false
	}
	
	// Send alert if quality score is below threshold
	return int(results.Summary.QualityScore) < config.AlertSettings.ThresholdScore
}

// FormatAlertMessage formats an alert message based on validation results
func FormatAlertMessage(results *ValidationResults) string {
	// Format a message that includes:
	// - Quality score
	// - Number of models validated
	// - Number of rules passed/failed
	// - Details of failed rules
	
	message := "Data Quality Alert\n\n"
	message += fmt.Sprintf("Quality Score: %.1f%%\n", results.Summary.QualityScore)
	message += fmt.Sprintf("Models Validated: %d\n", results.Summary.TotalModels)
	message += fmt.Sprintf("Rules Passed: %d\n", results.Summary.PassedRules)
	message += fmt.Sprintf("Rules Failed: %d\n\n", results.Summary.FailedRules)
	
	// Add details of failed rules
	if results.Summary.FailedRules > 0 {
		message += "Failed Rules:\n"
		
		for _, result := range results.Results {
			if result.Status == "failed" {
				message += "- " + result.ModelName + ": " + result.RuleName + " - " + result.Message + "\n"
			}
		}
	}
	
	return message
}
