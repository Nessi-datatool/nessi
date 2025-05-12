package rules

import (
	"fmt"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
)

// Rule defines the interface for data quality rules
type Rule interface {
	Validate(record arrow.Record) []ValidationError
	GetMetadata() RuleMetadata
}

// RuleMetadata contains rule configuration and metadata
type RuleMetadata struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`
	Tags        []string              `json:"tags"`
	Config      map[string]interface{} `json:"config"`
}

// ValidationError represents a single validation error
type ValidationError struct {
	RuleID      string `json:"rule_id"`
	Field       string `json:"field"`
	Value       string `json:"value"`
	Message     string `json:"message"`
	RowIndex    int64  `json:"row_index"`
}

// ValidationStats tracks validation metrics
type ValidationStats struct {
	TotalRecords   int64
	ValidRecords   int64
	InvalidRecords int64
	ErrorsByRule   map[string]int64
	StartTime      time.Time
	EndTime        time.Time
}

// RuleValidator orchestrates rule validation
type RuleValidator struct {
	rules []Rule
	stats ValidationStats
}

// NewRuleValidator creates a new validator
func NewRuleValidator() *RuleValidator {
	return &RuleValidator{
		rules: make([]Rule, 0),
		stats: ValidationStats{
			ErrorsByRule: make(map[string]int64),
		},
	}
}

// AddRule adds a new validation rule
func (v *RuleValidator) AddRule(rule Rule) {
	v.rules = append(v.rules, rule)
}

// ValidateRecord applies all rules to a record
func (v *RuleValidator) ValidateRecord(record arrow.Record) []ValidationError {
	var errors []ValidationError
	v.stats.TotalRecords++

	for _, rule := range v.rules {
		ruleErrors := rule.Validate(record)
		if len(ruleErrors) > 0 {
			errors = append(errors, ruleErrors...)
			v.stats.ErrorsByRule[rule.GetMetadata().ID] += int64(len(ruleErrors))
		}
	}

	if len(errors) > 0 {
		v.stats.InvalidRecords++
	} else {
		v.stats.ValidRecords++
	}

	return errors
}

// GetStats returns current validation statistics
func (v *RuleValidator) GetStats() ValidationStats {
	return v.stats
}

// ResetStats resets validation statistics
func (v *RuleValidator) ResetStats() {
	v.stats = ValidationStats{
		ErrorsByRule: make(map[string]int64),
		StartTime:    time.Now(),
	}
}

// NullCheckRule implements a rule to check for null values
type NullCheckRule struct {
	fields []string
	config RuleMetadata
}

// NewNullCheckRule creates a new null check rule
func NewNullCheckRule(fields []string, config RuleMetadata) *NullCheckRule {
	return &NullCheckRule{
		fields: fields,
		config: config,
	}
}

func (r *NullCheckRule) Validate(record arrow.Record) []ValidationError {
	var errors []ValidationError
	
	for _, field := range r.fields {
		colIdx := record.Schema().FieldIndices(field)
		if len(colIdx) == 0 {
			continue
		}
		
		col := record.Column(colIdx[0])
		for i := 0; i < int(record.NumRows()); i++ {
			if col.IsNull(i) {
				errors = append(errors, ValidationError{
					RuleID:   r.config.ID,
					Field:    field,
					Value:    "null",
					Message:  fmt.Sprintf("Field '%s' cannot be null at row %d", field, i),
					RowIndex: int64(i),
				})
			}
		}
	}
	
	return errors
}

func (r *NullCheckRule) GetMetadata() RuleMetadata {
	return r.config
}

// RangeCheckRule implements a rule to check value ranges
type RangeCheckRule struct {
	fields map[string]RangeConfig
	config RuleMetadata
}

// RangeConfig specifies the valid range for a field
type RangeConfig struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// NewRangeCheckRule creates a new range check rule
func NewRangeCheckRule(fields map[string]RangeConfig, config RuleMetadata) *RangeCheckRule {
	return &RangeCheckRule{
		fields: fields,
		config: config,
	}
}

func (r *RangeCheckRule) Validate(record arrow.Record) []ValidationError {
	var errors []ValidationError
	
	for field, rangeConfig := range r.fields {
		colIdx := record.Schema().FieldIndices(field)
		if len(colIdx) == 0 {
			continue
		}
		
		col := record.Column(colIdx[0])
		if col.DataType().ID() != arrow.FLOAT64 {
			continue
		}

		floatCol := col.(*array.Float64)
		for i := 0; i < int(record.NumRows()); i++ {
			if col.IsNull(i) {
				continue
			}
			
			value := floatCol.Value(i)
			if value < rangeConfig.Min || value > rangeConfig.Max {
				errors = append(errors, ValidationError{
					RuleID:   r.config.ID,
					Field:    field,
					Value:    fmt.Sprintf("%f", value),
					Message:  fmt.Sprintf("Value %f is outside range [%f, %f]",
						value, rangeConfig.Min, rangeConfig.Max),
					RowIndex: int64(i),
				})
			}
		}
	}
	
	return errors
}

func (r *RangeCheckRule) GetMetadata() RuleMetadata {
	return r.config
}
