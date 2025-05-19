package rules

import (
	"fmt"
	"regexp"
	"time"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
)

// Additional rule types for extended validation
const (
	RuleTypeLength   = "length"
	RuleTypeEnum     = "enum"
	RuleTypeRegex    = "regex"
	RuleTypeDateTime = "datetime"
)

// ExtendedRule extends the basic Rule with additional validation capabilities
type ExtendedRule struct {
	Rule
	MinLength    int      `json:"min_length,omitempty"`    // For length rules
	MaxLength    int      `json:"max_length,omitempty"`    // For length rules
	EnumValues   []string `json:"enum_values,omitempty"`   // For enum rules
	RegexPattern string   `json:"regex_pattern,omitempty"` // For regex rules
	DateFormat   string   `json:"date_format,omitempty"`   // For datetime rules
}

// Validate validates the extended rule configuration
func (r *ExtendedRule) Validate() error {
	// First validate the base rule
	if err := r.Rule.Validate(); err != nil {
		return err
	}

	// Validate extended rule properties
	switch r.Type {
	case RuleTypeLength:
		if r.MinLength < 0 {
			return fmt.Errorf("min_length cannot be negative")
		}
		if r.MaxLength > 0 && r.MinLength > r.MaxLength {
			return fmt.Errorf("min_length must be less than or equal to max_length")
		}
	case RuleTypeEnum:
		if len(r.EnumValues) == 0 {
			return fmt.Errorf("enum rule must specify at least one value")
		}
	case RuleTypeRegex:
		if r.RegexPattern == "" {
			return fmt.Errorf("regex rule must specify a pattern")
		}
		// Validate that the regex pattern is valid
		_, err := regexp.Compile(r.RegexPattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %v", err)
		}
	case RuleTypeDateTime:
		if r.DateFormat == "" {
			return fmt.Errorf("datetime rule must specify a date format")
		}
		// Validate that the date format is valid
		_, err := time.Parse(r.DateFormat, time.Now().Format(r.DateFormat))
		if err != nil {
			return fmt.Errorf("invalid date format: %v", err)
		}
	}
	return nil
}

// Evaluate evaluates a record against this extended rule
func (r *ExtendedRule) Evaluate(record arrow.Record) (bool, error) {
	// Find column index by name
	colIdx := -1
	for i := 0; i < int(record.NumCols()); i++ {
		if record.ColumnName(i) == r.Column {
			colIdx = i
			break
		}
	}
	if colIdx == -1 {
		return false, fmt.Errorf("column %s not found in record", r.Column)
	}

	column := record.Column(colIdx)
	switch r.Type {
	case RuleTypeLength:
		return evaluateLength(column, r.MinLength, r.MaxLength)
	case RuleTypeEnum:
		return evaluateEnum(column, r.EnumValues)
	case RuleTypeRegex:
		return evaluateRegex(column, r.RegexPattern)
	case RuleTypeDateTime:
		return evaluateDateTime(column, r.DateFormat)
	default:
		// For base rule types, delegate to the base Rule.Evaluate
		return r.Rule.Evaluate(record)
	}
}

// evaluateLength checks if string values have a length within the specified range
func evaluateLength(column arrow.Array, minLength, maxLength int) (bool, error) {
	// Check if the column is a string type
	if column.DataType().ID() != arrow.STRING {
		return false, fmt.Errorf("length rule can only be applied to string columns")
	}

	stringArray := column.(*array.String)
	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		value := stringArray.Value(i)
		length := len(value)

		// Check minimum length if specified
		if minLength > 0 && length < minLength {
			return false, nil
		}

		// Check maximum length if specified
		if maxLength > 0 && length > maxLength {
			return false, nil
		}
	}
	return true, nil
}

// evaluateEnum checks if values are in a predefined set
func evaluateEnum(column arrow.Array, enumValues []string) (bool, error) {
	// Create a map for faster lookups
	validValues := make(map[string]bool)
	for _, v := range enumValues {
		validValues[v] = true
	}

	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		value := column.ValueStr(i)
		if !validValues[value] {
			return false, nil
		}
	}
	return true, nil
}

// evaluateRegex checks if values match a regular expression pattern
func evaluateRegex(column arrow.Array, regexPattern string) (bool, error) {
	// Compile the regex pattern
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %v", err)
	}

	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		value := column.ValueStr(i)
		if !regex.MatchString(value) {
			return false, nil
		}
	}
	return true, nil
}

// evaluateDateTime checks if values are valid dates in the specified format
func evaluateDateTime(column arrow.Array, dateFormat string) (bool, error) {
	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		value := column.ValueStr(i)
		_, err := time.Parse(dateFormat, value)
		if err != nil {
			return false, nil
		}
	}
	return true, nil
}
