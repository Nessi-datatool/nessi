package profile

import (
	"testing"
)

func TestCalculateQualityScore(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		profile        *Profile
		expectedScore  *DataQualityScore
		expectedType   string
		expectedRanges map[string][2]float64 // field -> [min, max]
	}{
		{
			name: "Complete high-quality profile",
			profile: &Profile{
				Name:         "high_quality_column",
				Type:         "string",
				RowCount:     1000,
				Distinct:     950,
				NullCount:    0,
				NullPercent:  0,
				Patterns:     []string{"email"},
				Anomalies:    []Anomaly{},
				ValueCounts:  map[interface{}]int{"example@example.com": 1},
			},
			expectedRanges: map[string][2]float64{
				"Overall":      {0.9, 1.0},
				"Completeness": {0.99, 1.0},
				"Consistency":  {0.9, 1.0},
				"Accuracy":     {0.9, 1.0},
			},
		},
		{
			name: "Profile with nulls",
			profile: &Profile{
				Name:         "partial_nulls",
				Type:         "string",
				RowCount:     1000,
				Distinct:     800,
				NullCount:    200,
				NullPercent:  20,
				Patterns:     []string{"text"},
				Anomalies:    []Anomaly{},
				ValueCounts:  map[interface{}]int{"example": 1},
			},
			expectedRanges: map[string][2]float64{
				"Overall":      {0.7, 1.0},
				"Completeness": {0.79, 0.81},
				"Consistency":  {0.9, 1.0},
				"Accuracy":     {0.9, 1.0},
			},
		},
		{
			name: "Profile with anomalies",
			profile: &Profile{
				Name:         "with_anomalies",
				Type:         "integer",
				RowCount:     1000,
				Distinct:     900,
				NullCount:    0,
				NullPercent:  0,
				Patterns:     []string{"integer"},
				Anomalies: []Anomaly{
					{Type: "outlier", Value: 1, Description: "Outlier"},
					{Type: "outlier", Value: 2, Description: "Outlier"},
					{Type: "outlier", Value: 3, Description: "Outlier"},
					{Type: "outlier", Value: 4, Description: "Outlier"},
					{Type: "outlier", Value: 5, Description: "Outlier"},
				},
				ValueCounts: map[interface{}]int{42: 1},
			},
			expectedRanges: map[string][2]float64{
				"Overall":      {0.7, 1.0},
				"Completeness": {0.99, 1.0},
				"Consistency":  {0.7, 1.0},
				"Accuracy":     {0.7, 1.0},
			},
		},
		{
			name: "Profile with multiple patterns",
			profile: &Profile{
				Name:         "mixed_patterns",
				Type:         "string",
				RowCount:     1000,
				Distinct:     950,
				NullCount:    0,
				NullPercent:  0,
				Patterns:     []string{"integer", "decimal", "text"},
				Anomalies:    []Anomaly{},
				ValueCounts:  map[interface{}]int{"mixed": 1},
			},
			expectedRanges: map[string][2]float64{
				"Overall":      {0.7, 1.0},
				"Completeness": {0.99, 1.0},
				"Consistency":  {0.7, 1.0},
				"Accuracy":     {0.9, 1.0},
			},
			expectedType: "numeric",
		},
		{
			name: "Boolean-like integer profile",
			profile: &Profile{
				Name:         "bool_like",
				Type:         "integer",
				RowCount:     1000,
				Distinct:     2,
				NullCount:    0,
				NullPercent:  0,
				Patterns:     []string{"integer"},
				Anomalies:    []Anomaly{},
				ValueCounts: map[interface{}]int{
					0: 500,
					1: 500,
				},
			},
			expectedRanges: map[string][2]float64{
				"Overall":      {0.7, 1.0},
				"Completeness": {0.99, 1.0},
				"Consistency":  {0.9, 1.0},
				"Accuracy":     {0.7, 1.0},
			},
			expectedType: "boolean",
		},
		{
			name: "Date-like string profile",
			profile: &Profile{
				Name:         "date_like",
				Type:         "string",
				RowCount:     1000,
				Distinct:     365,
				NullCount:    0,
				NullPercent:  0,
				Patterns:     []string{"date"},
				Anomalies:    []Anomaly{},
				ValueCounts:  map[interface{}]int{"2023-01-01": 1},
			},
			expectedRanges: map[string][2]float64{
				"Overall":      {0.8, 1.0},
				"Completeness": {0.99, 1.0},
				"Consistency":  {0.9, 1.0},
				"Accuracy":     {0.9, 1.0},
			},
			expectedType: "date",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate quality score
			score := CalculateQualityScore(tc.profile)
			
			// Check if score is not nil
			if score == nil {
				t.Fatal("Expected non-nil quality score")
			}
			
			// Check recommended type if expected
			if tc.expectedType != "" && score.RecommendedType != tc.expectedType {
				t.Errorf("Expected recommended type %s, got %s", tc.expectedType, score.RecommendedType)
			}
			
			// Check score ranges
			for field, expectedRange := range tc.expectedRanges {
				var actualScore float64
				
				switch field {
				case "Overall":
					actualScore = score.Overall
				case "Completeness":
					actualScore = score.Completeness
				case "Consistency":
					actualScore = score.Consistency
				case "Accuracy":
					actualScore = score.Accuracy
				case "Uniqueness":
					actualScore = score.Uniqueness
				}
				
				min, max := expectedRange[0], expectedRange[1]
				if actualScore < min || actualScore > max {
					t.Errorf("%s score %f outside expected range [%f, %f]", field, actualScore, min, max)
				}
			}
		})
	}
}

func TestGenerateEnhancedProfile(t *testing.T) {
	// Create a test profile
	profile := &Profile{
		Name:         "test_column",
		Type:         "integer",
		RowCount:     1000,
		Distinct:     100,
		NullCount:    10,
		NullPercent:  1.0,
		Patterns:     []string{"integer"},
		Anomalies:    []Anomaly{
			{Type: "outlier", Value: 101, Description: "Outlier"},
			{Type: "outlier", Value: 102, Description: "Outlier"},
		},
		ValueCounts:  map[interface{}]int{42: 10},
		Stats: Stats{
			Min:       1.0,
			Max:       100.0,
			Mean:      50.5,
			StdDev:    28.87,
			Quartiles: []float64{25.0, 50.0, 75.0},
		},
	}
	
	// Generate enhanced profile
	enhanced := GenerateEnhancedProfile(profile)
	
	// Check if enhanced profile is not nil
	if enhanced == nil {
		t.Fatal("Expected non-nil enhanced profile")
	}
	
	// Check if original profile is preserved
	if enhanced.Profile != profile {
		t.Error("Original profile not preserved in enhanced profile")
	}
	
	// Check if quality score is generated
	if enhanced.QualityScore == nil {
		t.Error("Quality score not generated")
	}
	
	// Check if distribution is generated for numeric type
	if enhanced.Distribution == nil {
		t.Error("Distribution not generated for numeric type")
	}
	
	// Check if detailed patterns are generated
	if len(enhanced.DetailedPatterns) == 0 {
		t.Error("Detailed patterns not generated")
	}
	
	// Check if pattern details match the original patterns
	if len(enhanced.DetailedPatterns) != len(profile.Patterns) {
		t.Errorf("Expected %d detailed patterns, got %d", len(profile.Patterns), len(enhanced.DetailedPatterns))
	}
	
	// Check if pattern has required fields
	pattern := enhanced.DetailedPatterns[0]
	if pattern.Name != "integer" {
		t.Errorf("Expected pattern name 'integer', got '%s'", pattern.Name)
	}
	
	if pattern.Description == "" {
		t.Error("Pattern description is empty")
	}
	
	if pattern.Confidence <= 0 || pattern.Confidence > 1 {
		t.Errorf("Pattern confidence %f outside valid range (0, 1]", pattern.Confidence)
	}
}

func TestIsIdLikeColumn(t *testing.T) {
	// Test cases for ID-like columns
	idColumns := []string{
		"id",
		"user_id",
		"customer_id",
		"order_id",
		"uuid",
		"guid",
		"primary_key",
	}
	
	// Test cases for non-ID columns
	nonIdColumns := []string{
		"name",
		"description",
		"address",
		"phone",
		"email",
		"hidden",
		"valid",
	}
	
	// Test ID-like columns
	for _, col := range idColumns {
		if !isIdLikeColumn(col) {
			t.Errorf("Column '%s' should be identified as ID-like", col)
		}
	}
	
	// Test non-ID columns
	for _, col := range nonIdColumns {
		if isIdLikeColumn(col) {
			t.Errorf("Column '%s' should not be identified as ID-like", col)
		}
	}
}

func TestGetRecommendedType(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		profile        *Profile
		expectedType   string
	}{
		{
			name: "String that should be date",
			profile: &Profile{
				Name:     "date_column",
				Type:     "string",
				Patterns: []string{"date"},
			},
			expectedType: "date",
		},
		{
			name: "String that should be numeric",
			profile: &Profile{
				Name:     "numeric_column",
				Type:     "string",
				Patterns: []string{"integer"},
			},
			expectedType: "numeric",
		},
		{
			name: "Integer that should be boolean",
			profile: &Profile{
				Name:     "bool_column",
				Type:     "integer",
				Distinct: 2,
				ValueCounts: map[interface{}]int{
					0: 50,
					1: 50,
				},
			},
			expectedType: "boolean",
		},
		{
			name: "Float that should be integer",
			profile: &Profile{
				Name:     "int_column",
				Type:     "float",
				ValueCounts: map[interface{}]int{
					1.0: 10,
					2.0: 20,
					3.0: 30,
				},
			},
			expectedType: "integer",
		},
		{
			name: "Correct type",
			profile: &Profile{
				Name:     "correct_column",
				Type:     "string",
				Patterns: []string{"text"},
			},
			expectedType: "",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recommendedType := getRecommendedType(tc.profile)
			
			if recommendedType != tc.expectedType {
				t.Errorf("Expected recommended type '%s', got '%s'", tc.expectedType, recommendedType)
			}
		})
	}
}
