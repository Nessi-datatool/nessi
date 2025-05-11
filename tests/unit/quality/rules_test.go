package quality_test

import (
	"testing"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi-dev/internal/quality/rules"
)

func TestRuleValidation(t *testing.T) {
	tests := []struct {
		name     string
		rule     rules.Rule
		expected error
	}{
		{
			name: "valid completeness rule",
		rule: rules.Rule{
				Type:      rules.RuleTypeCompleteness,
				Column:    "age",
				Threshold: 90,
				Severity:  rules.SeverityHigh,
			},
			expected: nil,
		},
		{
			name: "invalid completeness threshold",
		rule: rules.Rule{
				Type:      rules.RuleTypeCompleteness,
				Column:    "age",
				Threshold: 105,
				Severity:  rules.SeverityHigh,
			},
			expected: assert.AnError,
		},
		{
			name: "valid pattern rule",
		rule: rules.Rule{
				Type:    rules.RuleTypePattern,
				Column:  "email",
				Pattern: "@",
				Severity: rules.SeverityMedium,
			},
			expected: nil,
		},
		{
			name: "invalid pattern rule",
		rule: rules.Rule{
				Type:    rules.RuleTypePattern,
				Column:  "email",
				Pattern: "",
				Severity: rules.SeverityMedium,
			},
			expected: assert.AnError,
		},
		{
			name: "valid range rule",
		rule: rules.Rule{
				Type:    rules.RuleTypeRange,
				Column:  "age",
				Min:     0,
				Max:     120,
				Severity: rules.SeverityHigh,
			},
			expected: nil,
		},
		{
			name: "invalid range rule",
		rule: rules.Rule{
				Type:    rules.RuleTypeRange,
				Column:  "age",
				Min:     120,
				Max:     0,
				Severity: rules.SeverityHigh,
			},
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.rule.Validate()
			if tt.expected == nil {
				assert.NoError(t, actual)
			} else {
				assert.Error(t, actual)
			}
		})
	}
}

func TestRuleEvaluation(t *testing.T) {
	// Create a test record with some data
	// Create schema
	schema := arrow.NewSchema([]arrow.Field{
		{Name: "age", Type: arrow.PrimitiveTypes.Int32},
	}, nil)

	// Create record builder for no-nulls record
	builder := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builder.Release()

	// Get the field builder for the age column
	ageBuilder := builder.Field(0).(*array.Int32Builder)
	defer ageBuilder.Release()

	// Append values - first create a record with no nulls
	ageBuilder.AppendValues([]int32{1, 2, 3, 4, 5}, nil)
	record := builder.NewRecord()
	defer record.Release()

	// Create a new builder for the record with nulls
	builderWithNulls := array.NewRecordBuilder(memory.DefaultAllocator, schema)
	defer builderWithNulls.Release()
	ageBuilderWithNulls := builderWithNulls.Field(0).(*array.Int32Builder)
	defer ageBuilderWithNulls.Release()
	ageBuilderWithNulls.AppendValues([]int32{1, 2, 3, 4, 5}, []bool{false, false, false, false, true})
	recordWithNulls := builderWithNulls.NewRecord()
	defer recordWithNulls.Release()

	tests := []struct {
		name     string
		rule     rules.Rule
		expected bool
	}{
		{
			name: "completeness rule - all values present",
			rule: rules.Rule{
				Type:      rules.RuleTypeCompleteness,
				Column:    "age",
				Threshold: 100,
				Severity:  rules.SeverityHigh,
			},
			expected: true,
		},
		{
			name: "completeness rule - some nulls",
			rule: rules.Rule{
				Type:      rules.RuleTypeCompleteness,
				Column:    "age",
				Threshold: 80,
				Severity:  rules.SeverityHigh,
			},
			expected: true, // Since we have no nulls, this should pass
		},
		{
			name: "completeness rule - with nulls",
			rule: rules.Rule{
				Type:      rules.RuleTypeCompleteness,
				Column:    "age",
				Threshold: 80,
				Severity:  rules.SeverityHigh,
			},
			expected: false, // Since we have nulls, this should fail
		},
		{
			name: "uniqueness rule - all unique",
			rule: rules.Rule{
				Type:      rules.RuleTypeUniqueness,
				Column:    "age",
				Severity:  rules.SeverityHigh,
			},
			expected: true,
		},
		{
			name: "range rule - within bounds",
			rule: rules.Rule{
				Type:    rules.RuleTypeRange,
				Column:  "age",
				Min:     0,
				Max:     10,
				Severity: rules.SeverityHigh,
			},
			expected: true, // Since all values are within 0-10, this should pass
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rec arrow.Record
			switch tt.name {
			case "completeness rule - with nulls":
				rec = recordWithNulls
			default:
				rec = record
			}
			actual, err := tt.rule.Evaluate(rec)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
