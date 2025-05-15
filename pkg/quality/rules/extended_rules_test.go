package rules

import (
	"testing"
	
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
)

func createTestRecord(colName string, dataType arrow.DataType, values interface{}, validBits []bool) arrow.Record {
	pool := memory.NewGoAllocator()
	
	// Create builder based on data type
	var builder array.Builder
	switch dataType.(type) {
	case *arrow.StringType:
		builder = array.NewStringBuilder(pool)
		strValues := values.([]string)
		for i, v := range strValues {
			if validBits[i] {
				builder.(*array.StringBuilder).Append(v)
			} else {
				builder.(*array.StringBuilder).AppendNull()
			}
		}
	case *arrow.Int64Type:
		builder = array.NewInt64Builder(pool)
		intValues := values.([]int64)
		for i, v := range intValues {
			if validBits[i] {
				builder.(*array.Int64Builder).Append(v)
			} else {
				builder.(*array.Int64Builder).AppendNull()
			}
		}
	case *arrow.Float64Type:
		builder = array.NewFloat64Builder(pool)
		floatValues := values.([]float64)
		for i, v := range floatValues {
			if validBits[i] {
				builder.(*array.Float64Builder).Append(v)
			} else {
				builder.(*array.Float64Builder).AppendNull()
			}
		}
	case *arrow.TimestampType:
		builder = array.NewTimestampBuilder(pool, dataType.(*arrow.TimestampType))
		timeValues := values.([]arrow.Timestamp)
		for i, v := range timeValues {
			if validBits[i] {
				builder.(*array.TimestampBuilder).Append(v)
			} else {
				builder.(*array.TimestampBuilder).AppendNull()
			}
		}
	default:
		panic("Unsupported data type for test")
	}
	
	// Build array
	arr := builder.NewArray()
	defer arr.Release()
	
	// Create schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: colName, Type: dataType, Nullable: true},
		},
		nil,
	)
	
	// Create record
	columns := []arrow.Array{arr}
	record := array.NewRecord(schema, columns, int64(len(validBits)))
	
	return record
}

func TestRegexRule(t *testing.T) {
	// Create test data
	colName := "email"
	values := []string{
		"user@example.com",
		"invalid-email",
		"another@test.com",
		"not-an-email",
	}
	validBits := []bool{true, true, true, true}
	
	// Create test record
	record := createTestRecord(colName, arrow.BinaryTypes.String, values, validBits)
	defer record.Release()
	
	// Create regex rule
	rule, err := NewRegexRule(colName, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, RuleMetadata{
		ID:          "email_regex",
		Name:        "Email Format",
		Description: "Validates email format",
		Severity:    "error",
	})
	
	if err != nil {
		t.Fatalf("Failed to create regex rule: %v", err)
	}
	
	// Validate
	errors := rule.Validate(record)
	
	// Check results
	if len(errors) != 2 {
		t.Errorf("Expected 2 validation errors, got %d", len(errors))
	}
	
	// Check specific errors
	for _, err := range errors {
		if err.Value != "invalid-email" && err.Value != "not-an-email" {
			t.Errorf("Unexpected error value: %s", err.Value)
		}
	}
}

func TestEnumRule(t *testing.T) {
	// Create test data
	colName := "status"
	values := []string{
		"active",
		"inactive",
		"pending",
		"unknown",
	}
	validBits := []bool{true, true, true, true}
	
	// Create test record
	record := createTestRecord(colName, arrow.BinaryTypes.String, values, validBits)
	defer record.Release()
	
	// Create enum rule
	rule := NewEnumRule(colName, []string{"active", "inactive", "pending"}, true, RuleMetadata{
		ID:          "status_enum",
		Name:        "Status Values",
		Description: "Validates status values",
		Severity:    "error",
	})
	
	// Validate
	errors := rule.Validate(record)
	
	// Check results
	if len(errors) != 1 {
		t.Errorf("Expected 1 validation error, got %d", len(errors))
	}
	
	// Check specific error
	if len(errors) > 0 {
		if errors[0].Value != "unknown" {
			t.Errorf("Expected error for 'unknown', got %s", errors[0].Value)
		}
	}
}

func TestLengthRule(t *testing.T) {
	// Create test data
	colName := "username"
	values := []string{
		"user1",
		"a",
		"verylongusername",
		"ok",
	}
	validBits := []bool{true, true, true, true}
	
	// Create test record
	record := createTestRecord(colName, arrow.BinaryTypes.String, values, validBits)
	defer record.Release()
	
	// Create length rule
	rule := NewLengthRule(colName, 3, 10, RuleMetadata{
		ID:          "username_length",
		Name:        "Username Length",
		Description: "Validates username length",
		Severity:    "error",
	})
	
	// Validate
	errors := rule.Validate(record)
	
	// Check results
	if len(errors) != 3 {
		t.Errorf("Expected 3 validation errors, got %d", len(errors))
	}
	
	// Check specific errors
	hasA := false
	hasLong := false
	for _, err := range errors {
		if err.Value == "a" {
			hasA = true
		} else if err.Value == "verylongusername" {
			hasLong = true
		}
	}
	
	if !hasA {
		t.Errorf("Missing error for value 'a'")
	}
	if !hasLong {
		t.Errorf("Missing error for value 'verylongusername'")
	}
}

func TestDateFormatRule(t *testing.T) {
	// Create test data
	colName := "date_string"
	values := []string{
		"2023-01-01",
		"01/01/2023",
		"2023-13-01", // Invalid month
		"not-a-date",
	}
	validBits := []bool{true, true, true, true}
	
	// Create test record
	record := createTestRecord(colName, arrow.BinaryTypes.String, values, validBits)
	defer record.Release()
	
	// Create date format rule for ISO format
	rule := NewDateFormatRule(colName, "2006-01-02", RuleMetadata{
		ID:          "date_format",
		Name:        "Date Format",
		Description: "Validates date format",
		Severity:    "error",
	})
	
	// Validate
	errors := rule.Validate(record)
	
	// Check results
	if len(errors) != 3 {
		t.Errorf("Expected 3 validation errors, got %d", len(errors))
	}
	
	// Check specific errors
	errorValues := make(map[string]bool)
	for _, err := range errors {
		errorValues[err.Value] = true
	}
	
	if !errorValues["01/01/2023"] {
		t.Errorf("Expected error for '01/01/2023'")
	}
	if !errorValues["2023-13-01"] {
		t.Errorf("Expected error for '2023-13-01'")
	}
	if !errorValues["not-a-date"] {
		t.Errorf("Expected error for 'not-a-date'")
	}
}

func TestYAMLLoader(t *testing.T) {
	// Load rules from YAML
	loader := NewRuleLoader()
	
	// Load rules from the test file
	err := loader.LoadFromYAML("/Users/meisi/Documents/nessi-dev/test_data/rules/test_rules.yaml")
	if err != nil {
		t.Fatalf("Failed to load rules from YAML: %v", err)
	}
	
	// Get the validator with loaded rules
	validator := loader.GetValidator()
	rules := validator.rules
	
	// Check for errors
	if err != nil {
		t.Fatalf("Failed to load rules from YAML: %v", err)
	}
	
	// Check number of rules
	if len(rules) != 4 {
		t.Errorf("Expected 4 rules, got %d", len(rules))
	}
	
	// Check rule types
	ruleTypes := make(map[string]bool)
	for _, rule := range rules {
		switch rule.(type) {
		case *RegexRule:
			ruleTypes["regex"] = true
		case *EnumRule:
			ruleTypes["enum"] = true
		case *LengthRule:
			ruleTypes["length"] = true
		case *DateFormatRule:
			ruleTypes["date_format"] = true
		}
	}
	
	// Verify all rule types are present
	if !ruleTypes["regex"] {
		t.Errorf("Missing RegexRule")
	}
	if !ruleTypes["enum"] {
		t.Errorf("Missing EnumRule")
	}
	if !ruleTypes["length"] {
		t.Errorf("Missing LengthRule")
	}
	if !ruleTypes["date_format"] {
		t.Errorf("Missing DateFormatRule")
	}
}
