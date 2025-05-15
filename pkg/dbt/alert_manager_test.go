package dbt

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewAlertManager(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-alert-manager-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test config file with alert settings
	configData := `
enable_dbt_plugin: true
dbt_project_path: "/path/to/dbt/project"
delta_base_path: "/delta"

alert_config:
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

	// Test NewAlertManager
	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	alertConfig := &config.AlertConfig
	alertManager := NewAlertManager(alertConfig)
	if alertManager == nil {
		t.Fatalf("NewAlertManager() returned nil")
	}

	if alertManager.config != alertConfig {
		t.Errorf("Expected config to be set correctly")
	}
}

func TestSendAlerts(t *testing.T) {
	// Create a mock server for Slack webhook
	slackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer slackServer.Close()

	// Create AlertConfig with mock server URL
	alertConfig := &AlertConfig{
		Enabled: true,
		Channels: []AlertChannel{
			{
				Type:    "slack",
				Webhook: slackServer.URL,
			},
			{
				Type:       "email",
				Recipients: []string{"data-team@example.com"},
			},
		},
	}

	alertManager := NewAlertManager(alertConfig)

	// Create test validation results
	results := &ValidationResults{
		Results: []ValidationResult{
			{
				ModelName:    "test_model_1",
				RuleName:     "not_null",
				Status:       "failed",
				Message:      "Column 'id' contains null values",
				TablePath:    "/delta/test/model1",
				FailureCount: 5,
			},
			{
				ModelName:    "test_model_2",
				RuleName:     "not_null",
				Status:       "passed",
				Message:      "All columns are non-null",
				TablePath:    "/delta/test/model2",
				FailureCount: 0,
			},
		},
		Summary: ValidationSummary{
			TotalModels:   2,
			PassedModels:  1,
			FailedModels:  1,
			TotalRules:    2,
			PassedRules:   1,
			FailedRules:   1,
			QualityScore:  50.0,
			ExecutionTime: 1.0,
		},
	}

	// Test SendAlerts
	err := alertManager.SendAlerts(results)
	if err != nil {
		t.Errorf("SendAlerts() error = %v", err)
	}

	// Test with disabled alerts
	alertConfig.Enabled = false
	err = alertManager.SendAlerts(results)
	if err != nil {
		t.Errorf("SendAlerts() with disabled alerts should not return error, got %v", err)
	}

	// Test with no failures
	alertConfig.Enabled = true
	results.Summary.FailedRules = 0
	results.Summary.FailedModels = 0
	results.Summary.PassedRules = 2
	results.Summary.PassedModels = 2
	results.Results[0].Status = "passed"
	err = alertManager.SendAlerts(results)
	if err != nil {
		t.Errorf("SendAlerts() with no failures should not return error, got %v", err)
	}
}

func TestFormatFailedRules(t *testing.T) {
	// Create an AlertManager
	alertConfig := &AlertConfig{
		Enabled: true,
	}
	alertManager := NewAlertManager(alertConfig)

	// Create test validation results with no failures
	noFailures := &ValidationResults{
		Results: []ValidationResult{
			{
				ModelName:    "test_model_1",
				RuleName:     "not_null",
				Status:       "passed",
				Message:      "All columns are non-null",
				TablePath:    "/delta/test/model1",
				FailureCount: 0,
			},
			{
				ModelName:    "test_model_2",
				RuleName:     "unique",
				Status:       "passed",
				Message:      "All values are unique",
				TablePath:    "/delta/test/model2",
				FailureCount: 0,
			},
		},
	}

	// Test formatFailedRules with no failures
	formatted := alertManager.formatFailedRules(noFailures)
	if formatted == "" {
		t.Errorf("formatFailedRules() with no failures should return a message, got empty string")
	}

	// Create test validation results with failures
	withFailures := &ValidationResults{
		Results: []ValidationResult{
			{
				ModelName:    "test_model_1",
				RuleName:     "not_null",
				Status:       "failed",
				Message:      "Column 'id' contains null values",
				TablePath:    "/delta/test/model1",
				FailureCount: 5,
			},
			{
				ModelName:    "test_model_2",
				RuleName:     "unique",
				Status:       "passed",
				Message:      "All values are unique",
				TablePath:    "/delta/test/model2",
				FailureCount: 0,
			},
			{
				ModelName:    "test_model_3",
				RuleName:     "valid_values",
				Status:       "failed",
				Message:      "Column 'status' contains invalid values",
				TablePath:    "/delta/test/model3",
				FailureCount: 10,
			},
		},
	}

	// Test formatFailedRules with failures
	formatted = alertManager.formatFailedRules(withFailures)
	
	// Check that the formatted string contains the failed rule information
	expectedContents := []string{
		"test_model_1", "not_null", "Column 'id' contains null values",
		"test_model_3", "valid_values", "Column 'status' contains invalid values",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(formatted, expected) {
			t.Errorf("formatFailedRules() does not contain %q in output: %q", expected, formatted)
		}
	}

	// Verify it doesn't contain passed rule information
	if strings.Contains(formatted, "test_model_2") || strings.Contains(formatted, "All values are unique") {
		t.Errorf("formatFailedRules() should not contain passed rules, got: %q", formatted)
	}
}

func TestQualityScoreTracking(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "dbt-quality-score-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// We don't need an AlertManager instance for this test as we're testing standalone functions

	// Test StoreQualityScore function
	modelName := "test_model"
	score := 75.5

	// Store a quality score
	err = StoreQualityScore(modelName, score)
	if err != nil {
		t.Fatalf("StoreQualityScore() error = %v", err)
	}

	// Test GetQualityTrend function
	trend, err := GetQualityTrend(modelName, 7)
	if err != nil {
		t.Fatalf("GetQualityTrend() error = %v", err)
	}

	// Verify the trend data
	if trend == nil {
		t.Fatalf("GetQualityTrend() returned nil trend")
	}

	if trend.ModelName != modelName {
		t.Errorf("Expected trend.ModelName = %s, got %s", modelName, trend.ModelName)
	}

	if len(trend.Scores) == 0 {
		t.Errorf("Expected trend to contain scores, got empty slice")
	}

	// Test with validation results
	results := &ValidationResults{
		Results: []ValidationResult{
			{
				ModelName:    "test_model_1",
				RuleName:     "not_null",
				Status:       "failed",
				Message:      "Column 'id' contains null values",
				TablePath:    "/delta/test/model1",
				FailureCount: 5,
			},
		},
		Summary: ValidationSummary{
			TotalModels:   1,
			PassedModels:  0,
			FailedModels:  1,
			TotalRules:    1,
			PassedRules:   0,
			FailedRules:   1,
			QualityScore:  0.0,
			ExecutionTime: 1.0,
		},
	}

	// Test GenerateQualityScore
	qualityScore := GenerateQualityScore(results)
	if qualityScore != 0.0 {
		t.Errorf("Expected quality score to be 0.0, got %f", qualityScore)
	}

	// Test with all passed rules
	results.Summary.PassedRules = 1
	results.Summary.FailedRules = 0
	qualityScore = GenerateQualityScore(results)
	if qualityScore != 100.0 {
		t.Errorf("Expected quality score to be 100.0, got %f", qualityScore)
	}

	// Test with no rules
	results.Summary.TotalRules = 0
	results.Summary.PassedRules = 0
	qualityScore = GenerateQualityScore(results)
	if qualityScore != 100.0 {
		t.Errorf("Expected quality score with no rules to be 100.0, got %f", qualityScore)
	}
}


