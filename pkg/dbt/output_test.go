package dbt

import (
	"bytes"
	"strings"
	"testing"
)

func TestOutputJSON(t *testing.T) {
	// Create test validation results
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
				Status:       "failed",
				Message:      "Rule failed",
				TablePath:    "/delta/test/model2",
				FailureCount: 5,
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
			ExecutionTime: 1.5,
		},
	}

	// Test JSON output
	var buf bytes.Buffer
	err := OutputJSON(results, &buf)
	if err != nil {
		t.Fatalf("OutputJSON() error = %v", err)
	}

	output := buf.String()
	
	// Check that the output contains expected JSON elements
	expectedElements := []string{
		`"model_name": "model1"`,
		`"rule_name": "rule1"`,
		`"status": "passed"`,
		`"model_name": "model2"`,
		`"status": "failed"`,
		`"failure_count": 5`,
		`"total_models": 2`,
		`"quality_score": 50`,
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("OutputJSON() output does not contain expected element: %s", expected)
		}
	}
}

func TestOutputCSV(t *testing.T) {
	// Create test validation results
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
				Status:       "failed",
				Message:      "Rule failed",
				TablePath:    "/delta/test/model2",
				FailureCount: 5,
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
			ExecutionTime: 1.5,
		},
	}

	// Test CSV output
	var buf bytes.Buffer
	err := OutputCSV(results, &buf)
	if err != nil {
		t.Fatalf("OutputCSV() error = %v", err)
	}

	output := buf.String()
	
	// Check that the output contains expected CSV elements
	expectedElements := []string{
		"Model,Rule,Status,Message,Table Path,Failure Count",
		"model1,rule1,passed,Rule passed,/delta/test/model1,0",
		"model2,rule2,failed,Rule failed,/delta/test/model2,5",
		"Summary",
		"Total Models,2",
		"Passed Models,1",
		"Failed Models,1",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("OutputCSV() output does not contain expected element: %s", expected)
		}
	}
}

func TestOutputTable(t *testing.T) {
	// Create test validation results
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
				Status:       "failed",
				Message:      "Rule failed",
				TablePath:    "/delta/test/model2",
				FailureCount: 5,
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
			ExecutionTime: 1.5,
		},
	}

	// Test table output
	var buf bytes.Buffer
	err := OutputTable(results, &buf)
	if err != nil {
		t.Fatalf("OutputTable() error = %v", err)
	}

	output := buf.String()
	
	// Check that the output contains expected table elements
	expectedElements := []string{
		"MODEL",
		"RULE",
		"STATUS",
		"model1",
		"rule1",
		"passed",
		"model2",
		"rule2",
		"failed",
		"SUMMARY",
		"Total Models:",
		"Passed Models:",
		"Failed Models:",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(output, expected) {
			t.Errorf("OutputTable() output does not contain expected element: %s", expected)
		}
	}
}

func TestOutputProfileResults(t *testing.T) {
	// Create test profile results
	results := &ProfileResults{
		Results: []ProfileResult{
			{
				ModelName:   "model1",
				TablePath:   "/delta/test/model1",
				ProfileType: "basic",
				RowCount:    1000,
				ColumnCount: 2,
				ColumnProfiles: map[string]ColumnStats{
					"id": {
						DataType:      "integer",
						NonNullCount:  1000,
						NullCount:     0,
						NullPercent:   0,
						Unique:        1000,
						UniquePercent: 100,
						Min:           1,
						Max:           1000,
					},
					"name": {
						DataType:      "string",
						NonNullCount:  950,
						NullCount:     50,
						NullPercent:   5,
						Unique:        900,
						UniquePercent: 90,
					},
				},
			},
		},
		Summary: ProfileSummary{
			TotalModels:   1,
			TotalColumns:  2,
			TotalRows:     1000,
			ExecutionTime: 1.5,
		},
	}

	// Test JSON output
	var jsonBuf bytes.Buffer
	err := OutputJSON(results, &jsonBuf)
	if err != nil {
		t.Fatalf("OutputJSON() error = %v", err)
	}

	jsonOutput := jsonBuf.String()
	
	// Check that the output contains expected JSON elements
	jsonExpectedElements := []string{
		`"model_name": "model1"`,
		`"table_path": "/delta/test/model1"`,
		`"profile_type": "basic"`,
		`"row_count": 1000`,
		`"column_profiles": {`,
		`"data_type": "integer"`,
		`"non_null_count": 1000`,
		`"null_count": 0`,
	}

	for _, expected := range jsonExpectedElements {
		if !strings.Contains(jsonOutput, expected) {
			t.Errorf("OutputJSON() output does not contain expected element: %s", expected)
		}
	}

	// Test CSV output
	var csvBuf bytes.Buffer
	err = OutputCSV(results, &csvBuf)
	if err != nil {
		t.Fatalf("OutputCSV() error = %v", err)
	}

	// Test table output
	var tableBuf bytes.Buffer
	err = OutputTable(results, &tableBuf)
	if err != nil {
		t.Fatalf("OutputTable() error = %v", err)
	}
}
