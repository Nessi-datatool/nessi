package rules

import (
	"fmt"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/google/uuid"
)

// RuleValidator validates data against rules
type RuleValidator struct {
	rules       []ExtendedRule
	tracker     *RuleExecutionTracker
	datasetName string
}

// NewRuleValidator creates a new rule validator
func NewRuleValidator(rules []ExtendedRule, tracker *RuleExecutionTracker, datasetName string) *RuleValidator {
	return &RuleValidator{
		rules:       rules,
		tracker:     tracker,
		datasetName: datasetName,
	}
}

// ValidateRecord validates a single record against all rules
func (v *RuleValidator) ValidateRecord(record arrow.Record) ([]ValidationResult, error) {
	results := make([]ValidationResult, 0, len(v.rules))

	for _, rule := range v.rules {
		// Generate a unique ID for the rule if it doesn't have one
		ruleID := fmt.Sprintf("rule_%s_%s", rule.Type, rule.Column)

		// Evaluate the rule
		success, err := rule.Evaluate(record)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate rule '%s': %v", rule.Name, err)
		}

		// Create a validation result
		result := ValidationResult{
			RuleID:      ruleID,
			RuleName:    rule.Name,
			Column:      rule.Column,
			Type:        string(rule.Type),
			Severity:    string(rule.Severity),
			Success:     success,
			Timestamp:   time.Now(),
			RecordCount: record.NumRows(),
		}

		if !success {
			result.Message = fmt.Sprintf("Rule '%s' failed for column '%s'", rule.Name, rule.Column)
			// Count the number of records that violate the rule
			// This is a simplified implementation; in a real system, you would count actual violations
			result.ErrorCount = record.NumRows()
			result.ErrorRate = 100.0
		}

		results = append(results, result)
	}

	return results, nil
}

// ValidateAndTrack validates a record and tracks the results
func (v *RuleValidator) ValidateAndTrack(record arrow.Record) error {
	// Validate the record
	results, err := v.ValidateRecord(record)
	if err != nil {
		return err
	}

	// Create a validation history
	history := CreateValidationHistory(v.datasetName, results, record.NumRows())

	// Save the validation history
	if v.tracker != nil {
		if err := v.tracker.SaveValidationHistory(history); err != nil {
			return fmt.Errorf("failed to save validation history: %v", err)
		}
	}

	return nil
}

// ValidateFromYAML loads rules from a YAML file and validates a record
func ValidateFromYAML(yamlPath string, record arrow.Record, historyDir string, datasetName string) error {
	// Load rules from YAML
	rules, err := LoadRulesFromYAML(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to load rules from YAML: %v", err)
	}

	// Create a rule execution tracker
	tracker, err := NewRuleExecutionTracker(historyDir)
	if err != nil {
		return fmt.Errorf("failed to create rule execution tracker: %v", err)
	}

	// Create a rule validator
	validator := NewRuleValidator(rules, tracker, datasetName)

	// Validate and track
	return validator.ValidateAndTrack(record)
}

// GenerateRuleReport generates a report of rule validation trends
func GenerateRuleReport(historyDir string, datasetName string, days int) (map[string]interface{}, error) {
	// Create a rule execution tracker
	tracker, err := NewRuleExecutionTracker(historyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule execution tracker: %v", err)
	}

	// Get validation history
	histories, err := tracker.GetValidationHistory(datasetName, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get validation history: %v", err)
	}

	// Get validation trends
	trends, err := tracker.GetValidationTrends(datasetName, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get validation trends: %v", err)
	}

	// Create a report
	report := map[string]interface{}{
		"dataset_name":  datasetName,
		"report_id":     uuid.New().String(),
		"timestamp":     time.Now(),
		"days_analyzed": days,
		"total_runs":    len(histories),
		"latest_run":    nil,
		"trends":        trends,
		"rule_summary":  make(map[string]interface{}),
	}

	// Add the latest run if available
	if len(histories) > 0 {
		report["latest_run"] = histories[0]
	}

	// Create a summary of rule performance
	ruleSummary := make(map[string]interface{})
	ruleFailures := make(map[string]int)

	for _, history := range histories {
		for _, result := range history.Results {
			ruleID := result.RuleID
			if _, ok := ruleSummary[ruleID]; !ok {
				ruleSummary[ruleID] = map[string]interface{}{
					"rule_id":      ruleID,
					"rule_name":    result.RuleName,
					"column":       result.Column,
					"type":         result.Type,
					"severity":     result.Severity,
					"runs":         0,
					"failures":     0,
					"success_rate": 100.0,
				}
			}

			summary := ruleSummary[ruleID].(map[string]interface{})
			summary["runs"] = summary["runs"].(int) + 1

			if !result.Success {
				summary["failures"] = summary["failures"].(int) + 1
				ruleFailures[ruleID] = ruleFailures[ruleID] + 1
			}

			if summary["runs"].(int) > 0 {
				successRate := float64(summary["runs"].(int)-summary["failures"].(int)) / float64(summary["runs"].(int)) * 100
				summary["success_rate"] = successRate
			}
		}
	}

	report["rule_summary"] = ruleSummary

	// Identify the most frequently failing rules
	type ruleFailure struct {
		ruleID   string
		failures int
	}

	topFailures := make([]ruleFailure, 0, len(ruleFailures))
	for ruleID, failures := range ruleFailures {
		topFailures = append(topFailures, ruleFailure{
			ruleID:   ruleID,
			failures: failures,
		})
	}

	// Sort topFailures by number of failures (descending)
	for i := 0; i < len(topFailures); i++ {
		for j := i + 1; j < len(topFailures); j++ {
			if topFailures[i].failures < topFailures[j].failures {
				topFailures[i], topFailures[j] = topFailures[j], topFailures[i]
			}
		}
	}

	// Add top 5 failing rules to the report
	topFailingRules := make([]map[string]interface{}, 0, 5)
	for i := 0; i < len(topFailures) && i < 5; i++ {
		ruleID := topFailures[i].ruleID
		if summary, ok := ruleSummary[ruleID].(map[string]interface{}); ok {
			topFailingRules = append(topFailingRules, map[string]interface{}{
				"rule_id":      ruleID,
				"rule_name":    summary["rule_name"],
				"failures":     topFailures[i].failures,
				"success_rate": summary["success_rate"],
			})
		}
	}

	report["top_failing_rules"] = topFailingRules

	return report, nil
}
