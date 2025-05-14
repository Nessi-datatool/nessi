package dbt

import (
	"time"
)

// ModelValidationResult represents the validation results for a model
type ModelValidationResult struct {
	ModelName string             `json:"model_name"`
	TablePath string             `json:"table_path"`
	Rules     []*ValidationResult `json:"rules"`
	Timestamp time.Time          `json:"timestamp"`
}

// These types are already defined in validator.go:
// - ValidationResult
// - ValidationResults
// - ValidationSummary

// We need to update the alert_test.go file to use these types correctly
