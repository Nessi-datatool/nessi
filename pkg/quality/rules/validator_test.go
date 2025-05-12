package rules

import (
	"testing"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/stretchr/testify/assert"
)

func TestNullCheckRule(t *testing.T) {
	// Create test data
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "name", Type: arrow.BinaryTypes.String, Nullable: true},
		},
		nil,
	)

	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	// Add some data with nulls
	b.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 0, 3}, []bool{true, false, true})
	b.Field(1).(*array.StringBuilder).AppendValues([]string{"test", "", "example"}, []bool{true, false, true})

	record := b.NewRecord()
	defer record.Release()

	// Create rule
	rule := NewNullCheckRule(
		[]string{"id", "name"},
		RuleMetadata{
			ID:          "null_check",
			Name:        "Null Check",
			Description: "Check for null values",
			Severity:    "error",
		},
	)

	// Validate
	errors := rule.Validate(record)
	assert.Len(t, errors, 2) // Should find 2 null values
	assert.Equal(t, "id", errors[0].Field)
	assert.Equal(t, "name", errors[1].Field)
}

func TestRangeCheckRule(t *testing.T) {
	// Create test data
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "age", Type: arrow.PrimitiveTypes.Float64},
			{Name: "score", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	// Add some data with out-of-range values
	b.Field(0).(*array.Float64Builder).AppendValues([]float64{25, -1, 150}, nil)
	b.Field(1).(*array.Float64Builder).AppendValues([]float64{85, 101, 50}, nil)

	record := b.NewRecord()
	defer record.Release()

	// Create rule
	rule := NewRangeCheckRule(
		map[string]RangeConfig{
			"age":   {Min: 0, Max: 120},
			"score": {Min: 0, Max: 100},
		},
		RuleMetadata{
			ID:          "range_check",
			Name:        "Range Check",
			Description: "Check value ranges",
			Severity:    "warning",
		},
	)

	// Validate
	errors := rule.Validate(record)
	assert.Len(t, errors, 3) // Should find 3 out-of-range values
}

func TestRuleValidator(t *testing.T) {
	validator := NewRuleValidator()
	assert.NotNil(t, validator)

	// Add rules
	validator.AddRule(NewNullCheckRule(
		[]string{"id", "name"},
		RuleMetadata{ID: "null_check"},
	))
	validator.AddRule(NewRangeCheckRule(
		map[string]RangeConfig{"age": {Min: 0, Max: 120}},
		RuleMetadata{ID: "range_check"},
	))

	// Create test data
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Nullable: true},
			{Name: "name", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "age", Type: arrow.PrimitiveTypes.Float64},
		},
		nil,
	)

	b := array.NewRecordBuilder(pool, schema)
	defer b.Release()

	// Add data with violations
	b.Field(0).(*array.Int64Builder).AppendValues([]int64{1, 0}, []bool{true, false})
	b.Field(1).(*array.StringBuilder).AppendValues([]string{"test", ""}, []bool{true, true})
	b.Field(2).(*array.Float64Builder).AppendValues([]float64{25, 150}, nil)

	record := b.NewRecord()
	defer record.Release()

	// Validate
	errors := validator.ValidateRecord(record)
	assert.NotEmpty(t, errors)

	// Check stats
	stats := validator.GetStats()
	assert.Equal(t, int64(1), stats.TotalRecords)
	assert.Equal(t, int64(1), stats.InvalidRecords)
	assert.Greater(t, stats.ErrorsByRule["null_check"], int64(0))
	assert.Greater(t, stats.ErrorsByRule["range_check"], int64(0))
}
