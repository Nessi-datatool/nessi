package rules

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
)

// RegexRule implements a rule to validate string values against a regex pattern
type RegexRule struct {
	pattern *regexp.Regexp
	field   string
	config  RuleMetadata
}

// NewRegexRule creates a new regex validation rule
func NewRegexRule(field, pattern string, config RuleMetadata) (*RegexRule, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return &RegexRule{
		pattern: re,
		field:   field,
		config:  config,
	}, nil
}

func (r *RegexRule) Validate(record arrow.Record) []ValidationError {
	var errors []ValidationError

	colIdx := record.Schema().FieldIndices(r.field)
	if len(colIdx) == 0 {
		return errors
	}

	col := record.Column(colIdx[0])
	if col.DataType().ID() != arrow.STRING {
		return errors
	}

	strCol := col.(*array.String)
	for i := 0; i < int(record.NumRows()); i++ {
		if col.IsNull(i) {
			continue
		}

		value := strCol.Value(i)
		if !r.pattern.MatchString(value) {
			errors = append(errors, ValidationError{
				RuleID:   r.config.ID,
				Field:    r.field,
				Value:    value,
				Message:  fmt.Sprintf("Value '%s' does not match pattern '%s'", value, r.pattern.String()),
				RowIndex: int64(i),
			})
		}
	}

	return errors
}

func (r *RegexRule) GetMetadata() RuleMetadata {
	return r.config
}

// EnumRule implements a rule to validate that values are within a predefined set
type EnumRule struct {
	allowedValues map[string]bool
	field         string
	caseSensitive bool
	config        RuleMetadata
}

// NewEnumRule creates a new enum validation rule
func NewEnumRule(field string, values []string, caseSensitive bool, config RuleMetadata) *EnumRule {
	allowedValues := make(map[string]bool, len(values))

	for _, v := range values {
		if caseSensitive {
			allowedValues[v] = true
		} else {
			allowedValues[strings.ToLower(v)] = true
		}
	}

	return &EnumRule{
		allowedValues: allowedValues,
		field:         field,
		caseSensitive: caseSensitive,
		config:        config,
	}
}

func (r *EnumRule) Validate(record arrow.Record) []ValidationError {
	var errors []ValidationError

	colIdx := record.Schema().FieldIndices(r.field)
	if len(colIdx) == 0 {
		return errors
	}

	col := record.Column(colIdx[0])
	if col.DataType().ID() != arrow.STRING {
		return errors
	}

	strCol := col.(*array.String)
	for i := 0; i < int(record.NumRows()); i++ {
		if col.IsNull(i) {
			continue
		}

		value := strCol.Value(i)
		checkValue := value
		if !r.caseSensitive {
			checkValue = strings.ToLower(value)
		}

		if !r.allowedValues[checkValue] {
			errors = append(errors, ValidationError{
				RuleID:   r.config.ID,
				Field:    r.field,
				Value:    value,
				Message:  fmt.Sprintf("Value '%s' is not in the allowed set", value),
				RowIndex: int64(i),
			})
		}
	}

	return errors
}

func (r *EnumRule) GetMetadata() RuleMetadata {
	return r.config
}

// LengthRule implements a rule to validate string length
type LengthRule struct {
	minLength int
	maxLength int
	field     string
	config    RuleMetadata
}

// NewLengthRule creates a new length validation rule
func NewLengthRule(field string, minLength, maxLength int, config RuleMetadata) *LengthRule {
	return &LengthRule{
		minLength: minLength,
		maxLength: maxLength,
		field:     field,
		config:    config,
	}
}

func (r *LengthRule) Validate(record arrow.Record) []ValidationError {
	var errors []ValidationError

	colIdx := record.Schema().FieldIndices(r.field)
	if len(colIdx) == 0 {
		return errors
	}

	col := record.Column(colIdx[0])
	if col.DataType().ID() != arrow.STRING {
		return errors
	}

	strCol := col.(*array.String)
	for i := 0; i < int(record.NumRows()); i++ {
		if col.IsNull(i) {
			continue
		}

		value := strCol.Value(i)
		length := len(value)

		if length < r.minLength || length > r.maxLength {
			errors = append(errors, ValidationError{
				RuleID: r.config.ID,
				Field:  r.field,
				Value:  value,
				Message: fmt.Sprintf("Length of '%s' (%d) is outside range [%d, %d]",
					value, length, r.minLength, r.maxLength),
				RowIndex: int64(i),
			})
		}
	}

	return errors
}

func (r *LengthRule) GetMetadata() RuleMetadata {
	return r.config
}

// DateFormatRule implements a rule to validate date format
type DateFormatRule struct {
	format string
	field  string
	config RuleMetadata
}

// NewDateFormatRule creates a new date format validation rule
func NewDateFormatRule(field, format string, config RuleMetadata) *DateFormatRule {
	return &DateFormatRule{
		format: format,
		field:  field,
		config: config,
	}
}

func (r *DateFormatRule) Validate(record arrow.Record) []ValidationError {
	var errors []ValidationError

	colIdx := record.Schema().FieldIndices(r.field)
	if len(colIdx) == 0 {
		return errors
	}

	col := record.Column(colIdx[0])
	if col.DataType().ID() != arrow.STRING {
		return errors
	}

	strCol := col.(*array.String)
	for i := 0; i < int(record.NumRows()); i++ {
		if col.IsNull(i) {
			continue
		}

		value := strCol.Value(i)
		_, err := time.Parse(r.format, value)
		if err != nil {
			errors = append(errors, ValidationError{
				RuleID:   r.config.ID,
				Field:    r.field,
				Value:    value,
				Message:  fmt.Sprintf("Value '%s' does not match date format '%s'", value, r.format),
				RowIndex: int64(i),
			})
		}
	}

	return errors
}

func (r *DateFormatRule) GetMetadata() RuleMetadata {
	return r.config
}

// RuleExecutionRecord tracks the execution history of a rule
type RuleExecutionRecord struct {
	ExecutionID     string
	RuleID          string
	Timestamp       time.Time
	RecordsChecked  int64
	Failures        int64
	ExecutionTimeMs int64
	DatasetID       string
}

// ExecutionTrend represents a trend in rule execution results over time
type ExecutionTrend struct {
	Date            time.Time
	ExecutionCount  int
	FailureRate     float64
	AverageFailures float64
}

// RuleWithHistory extends the Rule interface to include history tracking
type RuleWithHistory interface {
	Rule
	IsEnabled() bool
	SetEnabled(enabled bool)
	GetHistory() []RuleExecutionRecord
	AddExecutionRecord(record RuleExecutionRecord)
}

// BaseRule provides common functionality for rules with history
type BaseRule struct {
	enabled bool
	history []RuleExecutionRecord
}

// IsEnabled returns whether the rule is enabled
func (r *BaseRule) IsEnabled() bool {
	return r.enabled
}

// SetEnabled sets whether the rule is enabled
func (r *BaseRule) SetEnabled(enabled bool) {
	r.enabled = enabled
}

// GetHistory returns the execution history of the rule
func (r *BaseRule) GetHistory() []RuleExecutionRecord {
	return r.history
}

// AddExecutionRecord adds an execution record to the rule's history
func (r *BaseRule) AddExecutionRecord(record RuleExecutionRecord) {
	r.history = append(r.history, record)
}
