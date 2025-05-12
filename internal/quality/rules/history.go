package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ValidationResult represents the result of a rule validation
type ValidationResult struct {
	RuleID      string    `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	Column      string    `json:"column"`
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Success     bool      `json:"success"`
	Message     string    `json:"message,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	RecordCount int64     `json:"record_count"`
	ErrorCount  int64     `json:"error_count,omitempty"`
	ErrorRate   float64   `json:"error_rate,omitempty"`
}

// ValidationHistory represents a collection of validation results
type ValidationHistory struct {
	DatasetName  string            `json:"dataset_name"`
	RunID        string            `json:"run_id"`
	Timestamp    time.Time         `json:"timestamp"`
	TotalRules   int               `json:"total_rules"`
	PassedRules  int               `json:"passed_rules"`
	FailedRules  int               `json:"failed_rules"`
	Results      []ValidationResult `json:"results"`
	TotalRecords int64             `json:"total_records"`
}

// RuleExecutionTracker tracks rule execution history
type RuleExecutionTracker struct {
	HistoryDir string
}

// NewRuleExecutionTracker creates a new rule execution tracker
func NewRuleExecutionTracker(historyDir string) (*RuleExecutionTracker, error) {
	// Ensure the history directory exists
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %v", err)
	}

	return &RuleExecutionTracker{
		HistoryDir: historyDir,
	}, nil
}

// SaveValidationHistory saves validation history to a file
func (t *RuleExecutionTracker) SaveValidationHistory(history ValidationHistory) error {
	// Generate a filename based on dataset name and timestamp
	timestamp := history.Timestamp.Format("20060102_150405")
	filename := fmt.Sprintf("%s_%s.json", history.DatasetName, timestamp)
	filePath := filepath.Join(t.HistoryDir, filename)

	// Marshal the history to JSON
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal validation history: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write validation history: %v", err)
	}

	return nil
}

// GetValidationHistory retrieves validation history for a dataset
func (t *RuleExecutionTracker) GetValidationHistory(datasetName string, limit int) ([]ValidationHistory, error) {
	// List all history files for the dataset
	pattern := fmt.Sprintf("%s_*.json", datasetName)
	matches, err := filepath.Glob(filepath.Join(t.HistoryDir, pattern))
	if err != nil {
		return nil, fmt.Errorf("failed to list history files: %v", err)
	}

	// Sort files by modification time (newest first)
	// This is a simple implementation; in a real system, you might want to use a more efficient approach
	type fileInfo struct {
		path    string
		modTime time.Time
	}
	files := make([]fileInfo, 0, len(matches))
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		files = append(files, fileInfo{
			path:    match,
			modTime: info.ModTime(),
		})
	}

	// Sort files by modification time (newest first)
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if files[i].modTime.Before(files[j].modTime) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	// Limit the number of files to process
	if limit > 0 && limit < len(files) {
		files = files[:limit]
	}

	// Load history from files
	histories := make([]ValidationHistory, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file.path)
		if err != nil {
			continue
		}

		var history ValidationHistory
		if err := json.Unmarshal(data, &history); err != nil {
			continue
		}

		histories = append(histories, history)
	}

	return histories, nil
}

// GetValidationTrends analyzes validation history to identify trends
func (t *RuleExecutionTracker) GetValidationTrends(datasetName string, days int) (map[string][]float64, error) {
	// Get validation history for the specified number of days
	histories, err := t.GetValidationHistory(datasetName, 0)
	if err != nil {
		return nil, err
	}

	// Filter histories by date range
	cutoff := time.Now().AddDate(0, 0, -days)
	filteredHistories := make([]ValidationHistory, 0, len(histories))
	for _, history := range histories {
		if history.Timestamp.After(cutoff) {
			filteredHistories = append(filteredHistories, history)
		}
	}

	// Calculate trends
	trends := make(map[string][]float64)
	
	// Overall success rate trend
	successRates := make([]float64, 0, len(filteredHistories))
	for _, history := range filteredHistories {
		if history.TotalRules > 0 {
			successRate := float64(history.PassedRules) / float64(history.TotalRules) * 100
			successRates = append(successRates, successRate)
		}
	}
	trends["overall_success_rate"] = successRates

	// Rule-specific success rate trends
	ruleSuccessRates := make(map[string][]float64)
	for _, history := range filteredHistories {
		for _, result := range history.Results {
			ruleID := result.RuleID
			if _, ok := ruleSuccessRates[ruleID]; !ok {
				ruleSuccessRates[ruleID] = make([]float64, 0, len(filteredHistories))
			}
			if result.Success {
				ruleSuccessRates[ruleID] = append(ruleSuccessRates[ruleID], 100.0)
			} else {
				ruleSuccessRates[ruleID] = append(ruleSuccessRates[ruleID], 0.0)
			}
		}
	}
	
	// Add rule-specific trends to the overall trends
	for ruleID, rates := range ruleSuccessRates {
		trends[fmt.Sprintf("rule_%s", ruleID)] = rates
	}

	return trends, nil
}

// CreateValidationHistory creates a validation history from rule results
func CreateValidationHistory(datasetName string, results []ValidationResult, totalRecords int64) ValidationHistory {
	passedRules := 0
	failedRules := 0
	for _, result := range results {
		if result.Success {
			passedRules++
		} else {
			failedRules++
		}
	}

	return ValidationHistory{
		DatasetName:  datasetName,
		RunID:        fmt.Sprintf("run_%s", time.Now().Format("20060102_150405")),
		Timestamp:    time.Now(),
		TotalRules:   len(results),
		PassedRules:  passedRules,
		FailedRules:  failedRules,
		Results:      results,
		TotalRecords: totalRecords,
	}
}
