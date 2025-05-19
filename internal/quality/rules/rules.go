package rules

import (
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
)

// RuleType defines different types of quality rules
const (
	RuleTypeCompleteness = "completeness"
	RuleTypeUniqueness   = "uniqueness"
	RuleTypePattern      = "pattern"
	RuleTypeRange        = "range"
)

// Rule represents a data quality rule
// This extends the basic QualityRule from manager
// with additional implementation details
//
//go:generate stringer -type=RuleType
//go:generate stringer -type=Severity
type Rule struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        RuleType `json:"type"`
	Column      string   `json:"column"`
	Condition   string   `json:"condition"`
	Threshold   float64  `json:"threshold"`
	Severity    Severity `json:"severity"`
	Pattern     string   `json:"pattern,omitempty"` // Only used for pattern rules
	Min         float64  `json:"min,omitempty"`     // Only used for range rules
	Max         float64  `json:"max,omitempty"`     // Only used for range rules
}

// RuleType defines the type of quality rule
//
//go:generate stringer -type=RuleType
type RuleType string

// Severity defines the severity level of a rule violation
//
//go:generate stringer -type=Severity
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Validate validates the rule configuration
func (r *Rule) Validate() error {
	switch r.Type {
	case RuleTypeCompleteness:
		if r.Threshold <= 0 || r.Threshold > 100 {
			return fmt.Errorf("completeness threshold must be between 0 and 100")
		}
	case RuleTypeUniqueness:
		if r.Threshold <= 0 || r.Threshold > 100 {
			return fmt.Errorf("uniqueness threshold must be between 0 and 100")
		}
	case RuleTypePattern:
		if r.Pattern == "" {
			return fmt.Errorf("pattern rule must specify a pattern")
		}
	case RuleTypeRange:
		if r.Min >= r.Max {
			return fmt.Errorf("min value must be less than max value for range rule")
		}
	}
	return nil
}

// Evaluate evaluates a record against this rule
func (r *Rule) Evaluate(record arrow.Record) (bool, error) {
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
	case RuleTypeCompleteness:
		return evaluateCompleteness(column, r.Threshold)
	case RuleTypeUniqueness:
		return evaluateUniqueness(column)
	case RuleTypePattern:
		return evaluatePattern(column, r.Pattern)
	case RuleTypeRange:
		return evaluateRange(column, r.Min, r.Max)
	default:
		return false, fmt.Errorf("unsupported rule type: %s", r.Type)
	}
}

// evaluateCompleteness checks the percentage of non-null values
func evaluateCompleteness(column arrow.Array, threshold float64) (bool, error) {
	nullCount := 0
	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			nullCount++
		}
	}

	completeness := float64(column.Len()-nullCount) / float64(column.Len()) * 100
	return completeness >= threshold, nil
}

// evaluateUniqueness checks for duplicate values
func evaluateUniqueness(column arrow.Array) (bool, error) {
	values := make(map[string]bool)
	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		value := column.ValueStr(i)
		values[value] = true
	}

	return len(values) == column.Len(), nil
}

// evaluatePattern checks if values match a pattern
func evaluatePattern(column arrow.Array, pattern string) (bool, error) {
	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		value := column.ValueStr(i)
		if !strings.Contains(value, pattern) {
			return false, nil
		}
	}
	return true, nil
}

// evaluateRange checks if values fall within a specified range
func evaluateRange(column arrow.Array, min, max float64) (bool, error) {
	for i := 0; i < column.Len(); i++ {
		if column.IsNull(i) {
			continue
		}
		// Convert to float64 based on the actual type
		switch column.DataType().ID() {
		case arrow.INT32:
			value := column.(*array.Int32).Value(i)
			if float64(value) < min || float64(value) > max {
				return false, nil
			}
		case arrow.INT64:
			value := column.(*array.Int64).Value(i)
			if float64(value) < min || float64(value) > max {
				return false, nil
			}
		case arrow.FLOAT32:
			value := column.(*array.Float32).Value(i)
			if float64(value) < min || float64(value) > max {
				return false, nil
			}
		case arrow.FLOAT64:
			value := column.(*array.Float64).Value(i)
			if value < min || value > max {
				return false, nil
			}
		default:
			return false, fmt.Errorf("unsupported data type: %s", column.DataType().String())
		}
	}
	return true, nil
}
