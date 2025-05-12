package rules

import (
	"testing"

	"github.com/apache/arrow/go/v12/arrow"
	"github.com/apache/arrow/go/v12/arrow/array"
	"github.com/apache/arrow/go/v12/arrow/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtendedRules(t *testing.T) {
	// Create a test Arrow record with various data types
	pool := memory.NewGoAllocator()
	
	// Create schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "username", Type: arrow.BinaryTypes.String},
			{Name: "email", Type: arrow.BinaryTypes.String},
			{Name: "age", Type: arrow.PrimitiveTypes.Int32},
			{Name: "status", Type: arrow.BinaryTypes.String},
			{Name: "created_at", Type: arrow.BinaryTypes.String},
		},
		nil,
	)
	
	// Create builders
	idBuilder := array.NewInt32Builder(pool)
	defer idBuilder.Release()
	
	usernameBuilder := array.NewStringBuilder(pool)
	defer usernameBuilder.Release()
	
	emailBuilder := array.NewStringBuilder(pool)
	defer emailBuilder.Release()
	
	ageBuilder := array.NewInt32Builder(pool)
	defer ageBuilder.Release()
	
	statusBuilder := array.NewStringBuilder(pool)
	defer statusBuilder.Release()
	
	createdAtBuilder := array.NewStringBuilder(pool)
	defer createdAtBuilder.Release()
	
	// Add data
	idBuilder.AppendValues([]int32{1, 2, 3, 4, 5}, nil)
	
	usernameBuilder.AppendValues([]string{
		"user1",
		"user2",
		"user3",
		"user4",
		"user_with_very_long_name",
	}, nil)
	
	emailBuilder.AppendValues([]string{
		"user1@example.com",
		"user2@example.com",
		"invalid-email",
		"user4@example.com",
		"user5@example.com",
	}, nil)
	
	ageBuilder.AppendValues([]int32{25, 30, 150, 40, 35}, nil)
	
	statusBuilder.AppendValues([]string{
		"active",
		"inactive",
		"pending",
		"unknown",
		"active",
	}, nil)
	
	createdAtBuilder.AppendValues([]string{
		"2023-01-01",
		"2023-02-15",
		"invalid-date",
		"2023-04-10",
		"2023-05-20",
	}, nil)
	
	// Build arrays
	idArray := idBuilder.NewArray()
	defer idArray.Release()
	
	usernameArray := usernameBuilder.NewArray()
	defer usernameArray.Release()
	
	emailArray := emailBuilder.NewArray()
	defer emailArray.Release()
	
	ageArray := ageBuilder.NewArray()
	defer ageArray.Release()
	
	statusArray := statusBuilder.NewArray()
	defer statusArray.Release()
	
	createdAtArray := createdAtBuilder.NewArray()
	defer createdAtArray.Release()
	
	// Create record
	record := array.NewRecord(schema, []arrow.Array{
		idArray, usernameArray, emailArray, ageArray, statusArray, createdAtArray,
	}, 5)
	defer record.Release()
	
	// Test length rule
	t.Run("LengthRule", func(t *testing.T) {
		rule := ExtendedRule{
			Rule: Rule{
				Name:        "Username Length Check",
				Description: "Checks that username length is between 3 and 10 characters",
				Type:        RuleTypeLength,
				Column:      "username",
				Severity:    SeverityMedium,
			},
			MinLength: 3,
			MaxLength: 10,
		}
		
		success, err := rule.Evaluate(record)
		require.NoError(t, err)
		assert.False(t, success, "Rule should fail due to username 'user_with_very_long_name'")
	})
	
	// Test enum rule
	t.Run("EnumRule", func(t *testing.T) {
		rule := ExtendedRule{
			Rule: Rule{
				Name:        "Status Enum Check",
				Description: "Checks that status values are valid",
				Type:        RuleTypeEnum,
				Column:      "status",
				Severity:    SeverityMedium,
			},
			EnumValues: []string{"active", "inactive", "pending"},
		}
		
		success, err := rule.Evaluate(record)
		require.NoError(t, err)
		assert.False(t, success, "Rule should fail due to status 'unknown'")
	})
	
	// Test regex rule
	t.Run("RegexRule", func(t *testing.T) {
		rule := ExtendedRule{
			Rule: Rule{
				Name:        "Email Format Check",
				Description: "Checks that email values match a valid email format",
				Type:        RuleTypeRegex,
				Column:      "email",
				Severity:    SeverityHigh,
			},
			RegexPattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		}
		
		success, err := rule.Evaluate(record)
		require.NoError(t, err)
		assert.False(t, success, "Rule should fail due to invalid email format")
	})
	
	// Test datetime rule
	t.Run("DateTimeRule", func(t *testing.T) {
		rule := ExtendedRule{
			Rule: Rule{
				Name:        "Date Format Check",
				Description: "Checks that date values are in the correct format",
				Type:        RuleTypeDateTime,
				Column:      "created_at",
				Severity:    SeverityMedium,
			},
			DateFormat: "2006-01-02",
		}
		
		success, err := rule.Evaluate(record)
		require.NoError(t, err)
		assert.False(t, success, "Rule should fail due to invalid date format")
	})
	
	// Test range rule
	t.Run("RangeRule", func(t *testing.T) {
		rule := ExtendedRule{
			Rule: Rule{
				Name:        "Age Range Check",
				Description: "Checks that age values are within a valid range",
				Type:        RuleTypeRange,
				Column:      "age",
				Severity:    SeverityMedium,
				Min:         0,
				Max:         120,
			},
		}
		
		success, err := rule.Evaluate(record)
		require.NoError(t, err)
		assert.False(t, success, "Rule should fail due to age 150")
	})
}

func TestRuleValidatorAndHistory(t *testing.T) {
	// Create a temporary directory for history files
	tempDir := t.TempDir()
	
	// Create a test Arrow record
	pool := memory.NewGoAllocator()
	
	// Create schema
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int32},
			{Name: "name", Type: arrow.BinaryTypes.String},
		},
		nil,
	)
	
	// Create builders
	idBuilder := array.NewInt32Builder(pool)
	defer idBuilder.Release()
	
	nameBuilder := array.NewStringBuilder(pool)
	defer nameBuilder.Release()
	
	// Add data
	idBuilder.AppendValues([]int32{1, 2, 3}, nil)
	nameBuilder.AppendValues([]string{"Alice", "Bob", "Charlie"}, nil)
	
	// Build arrays
	idArray := idBuilder.NewArray()
	defer idArray.Release()
	
	nameArray := nameBuilder.NewArray()
	defer nameArray.Release()
	
	// Create record
	record := array.NewRecord(schema, []arrow.Array{idArray, nameArray}, 3)
	defer record.Release()
	
	// Create rules
	rules := []ExtendedRule{
		{
			Rule: Rule{
				Name:        "ID Not Null Check",
				Description: "Checks that ID values are not null",
				Type:        RuleTypeCompleteness,
				Column:      "id",
				Threshold:   100.0,
				Severity:    SeverityCritical,
			},
		},
		{
			Rule: Rule{
				Name:        "Name Not Null Check",
				Description: "Checks that name values are not null",
				Type:        RuleTypeCompleteness,
				Column:      "name",
				Threshold:   100.0,
				Severity:    SeverityHigh,
			},
		},
	}
	
	// Create a rule execution tracker
	tracker, err := NewRuleExecutionTracker(tempDir)
	require.NoError(t, err)
	
	// Create a rule validator
	validator := NewRuleValidator(rules, tracker, "test_dataset")
	
	// Validate and track
	err = validator.ValidateAndTrack(record)
	require.NoError(t, err)
	
	// Get validation history
	histories, err := tracker.GetValidationHistory("test_dataset", 1)
	require.NoError(t, err)
	require.Len(t, histories, 1, "Should have 1 validation history")
	
	history := histories[0]
	assert.Equal(t, "test_dataset", history.DatasetName)
	assert.Equal(t, 2, history.TotalRules)
	assert.Equal(t, 2, history.PassedRules)
	assert.Equal(t, 0, history.FailedRules)
	assert.Len(t, history.Results, 2)
	
	// Generate a report
	report, err := GenerateRuleReport(tempDir, "test_dataset", 7)
	require.NoError(t, err)
	assert.Equal(t, "test_dataset", report["dataset_name"])
	assert.Equal(t, 7, report["days_analyzed"])
	assert.Equal(t, 1, report["total_runs"])
}

func TestYAMLRuleLoading(t *testing.T) {
	// Create a temporary YAML file
	tempDir := t.TempDir()
	yamlPath := tempDir + "/test_rules.yaml"
	
	// Create sample rules
	rules := []ExtendedRule{
		{
			Rule: Rule{
				Name:        "ID Not Null Check",
				Description: "Checks that ID values are not null",
				Type:        RuleTypeCompleteness,
				Column:      "id",
				Threshold:   100.0,
				Severity:    SeverityCritical,
			},
		},
		{
			Rule: Rule{
				Name:        "Email Format Check",
				Description: "Checks that email values match a valid email format",
				Type:        RuleTypeRegex,
				Column:      "email",
				Severity:    SeverityHigh,
			},
			RegexPattern: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
		},
	}
	
	// Save rules to YAML
	err := SaveRulesToYAML(rules, yamlPath)
	require.NoError(t, err)
	
	// Load rules from YAML
	loadedRules, err := LoadRulesFromYAML(yamlPath)
	require.NoError(t, err)
	require.Len(t, loadedRules, 2, "Should have 2 rules")
	
	// Verify loaded rules
	assert.Equal(t, "ID Not Null Check", loadedRules[0].Name)
	assert.Equal(t, string(RuleTypeCompleteness), string(loadedRules[0].Type))
	assert.Equal(t, 100.0, loadedRules[0].Threshold)
	
	assert.Equal(t, "Email Format Check", loadedRules[1].Name)
	assert.Equal(t, string(RuleTypeRegex), string(loadedRules[1].Type))
	assert.Equal(t, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, loadedRules[1].RegexPattern)
}
