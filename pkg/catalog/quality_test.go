package catalog

import (
	"context"
	"testing"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/api/types"
	"github.com/nessi-dev/nessi-dev/pkg/quality"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestQualityMetricsPublisher_PublishMetrics tests the PublishMetrics method
func TestQualityMetricsPublisher_PublishMetrics(t *testing.T) {
	// Create mock catalog
	mockCatalog := new(MockCatalog)
	mockCatalog.On("Name").Return("TestCatalog")
	mockCatalog.On("PublishQualityMetrics", mock.Anything, "test_db", "test_table", mock.Anything).Return(nil)

	// Create catalog manager
	manager := NewCatalogManager()
	manager.RegisterCatalog(mockCatalog)

	// Create quality metrics publisher
	publisher := NewQualityMetricsPublisher(manager)

	// Create test profile
	profile := &quality.Profile{
		Timestamp: time.Now(),
		Columns: map[string]*quality.ColumnProfile{
			"id": {
				Stats: &quality.ColumnStats{
					Count:     1000,
					NullCount: 0,
				},
			},
			"name": {
				Stats: &quality.ColumnStats{
					Count:     1000,
					NullCount: 5,
				},
			},
		},
	}

	// Create test validation results
	results := &quality.ValidationResults{
		RuleResults: []*quality.RuleResult{
			{
				Rule: &quality.Rule{
					Type: "range",
				},
				Score: 0.98,
			},
			{
				Rule: &quality.Rule{
					Type: "unique",
				},
				Score: 0.97,
			},
		},
	}

	// Publish metrics
	ctx := context.Background()
	err := publisher.PublishQualityMetrics(ctx, "TestCatalog", "test_db", "test_table", profile, results)
	assert.NoError(t, err)

	// Verify mock expectations
	mockCatalog.AssertExpectations(t)

	// Set up mock expectations for PublishQualityMetrics
	mockCatalog.On("PublishQualityMetrics", ctx, "test_db", "test_table", mock.MatchedBy(func(m *types.QualityMetrics) bool {
		// Verify metrics
		return m != nil
	})).Return(nil)

	// Call PublishMetrics
	err = publisher.PublishQualityMetrics(ctx, "TestCatalog", "test_db", "test_table", profile, results)
	assert.NoError(t, err)

	// Verify mock expectations
	mockCatalog.AssertExpectations(t)
}

// TestQualityMetricsPublisher_PublishQualityMetricsToAll tests the PublishQualityMetricsToAll method
func TestQualityMetricsPublisher_PublishQualityMetricsToAll(t *testing.T) {
	// Create mock catalogs
	mockCatalog1 := new(MockCatalog)
	mockCatalog1.On("Name").Return("TestCatalog1")
	mockCatalog1.On("PublishQualityMetrics", mock.Anything, "test_db", "test_table", mock.Anything).Return(nil)

	mockCatalog2 := new(MockCatalog)
	mockCatalog2.On("Name").Return("TestCatalog2")
	mockCatalog2.On("PublishQualityMetrics", mock.Anything, "test_db", "test_table", mock.Anything).Return(nil)
	
	// Create catalog manager
	manager := NewCatalogManager()
	manager.RegisterCatalog(mockCatalog1)
	manager.RegisterCatalog(mockCatalog2)
	
	// Create quality metrics publisher
	publisher := NewQualityMetricsPublisher(manager)
	
	// Create test profile
	profile := &quality.Profile{
		Columns: map[string]*quality.ColumnProfile{
			"col1": {
				Stats: &quality.ColumnStats{
					Count:     100,
					NullCount: 10,
				},
			},
		},
	}

	// Create test validation results
	results := &quality.ValidationResults{
		RuleResults: []*quality.RuleResult{
			{
				Rule: &quality.Rule{
					Type: "range",
				},
				Score: 0.95,
			},
		},
	}

	// Publish metrics
	ctx := context.Background()
	errs := publisher.PublishQualityMetricsToAll(ctx, "test_db", "test_table", profile, results)
	assert.Empty(t, errs)

	// Verify mock expectations
	mockCatalog1.AssertExpectations(t)
	mockCatalog2.AssertExpectations(t)

	mockCatalog1.On("PublishQualityMetrics", ctx, "testdb", "testtable", mock.Anything).Return(nil)
	mockCatalog2.On("PublishQualityMetrics", ctx, "testdb", "testtable", mock.Anything).Return(nil)
	
	// Call PublishQualityMetricsToAll
	errors := publisher.PublishQualityMetricsToAll(ctx, "testdb", "testtable", profile, results)
	
	// Verify results
	assert.Empty(t, errors)
	
	// Verify mock expectations
	mockCatalog1.AssertExpectations(t)
	mockCatalog2.AssertExpectations(t)
}

// Test helper functions

// TestCalculateCompleteness tests the calculateCompleteness function
func TestCalculateCompleteness(t *testing.T) {
	// Test with nil profile
	completeness := calculateCompleteness(nil)
	assert.Equal(t, 0.0, completeness)
	
	// Test with empty profile
	completeness = calculateCompleteness(&quality.Profile{})
	assert.Equal(t, 0.0, completeness)
	
	// Test with profile containing columns
	profile := &quality.Profile{
		Columns: map[string]*quality.ColumnProfile{
			"col1": {
				Stats: &quality.ColumnStats{
					Count:     100,
					NullCount: 10,
				},
			},
			"col2": {
				Stats: &quality.ColumnStats{
					Count:     100,
					NullCount: 20,
				},
			},
		},
	}
	
	completeness = calculateCompleteness(profile)
	assert.InDelta(t, 0.85, completeness, 0.001) // (0.9 + 0.8) / 2 = 0.85
}

// TestCalculateAccuracy tests the calculateAccuracy function
func TestCalculateAccuracy(t *testing.T) {
	// Test with nil results
	accuracy := calculateAccuracy(nil)
	assert.Equal(t, 0.0, accuracy)
	
	// Test with empty results
	accuracy = calculateAccuracy(&quality.ValidationResults{})
	assert.Equal(t, 0.0, accuracy)
	
	// Test with results containing accuracy rules
	results := &quality.ValidationResults{
		RuleResults: []*quality.RuleResult{
			{
				Rule: &quality.Rule{
					Name: "rule1",
					Type: "range",
				},
				Passed:  true,
				Score:   1.0,
				Details: "Rule passed",
			},
			{
				Rule: &quality.Rule{
					Name: "rule2",
					Type: "enum",
				},
				Passed:  false,
				Score:   0.8,
				Details: "Rule partially passed",
			},
			{
				Rule: &quality.Rule{
					Name: "rule3",
					Type: "unique", // Not an accuracy rule
				},
				Passed:  true,
				Score:   1.0,
				Details: "Rule passed",
			},
		},
	}
	
	accuracy = calculateAccuracy(results)
	assert.InDelta(t, 0.9, accuracy, 0.001) // (1.0 + 0.8) / 2 = 0.9
}

// TestCalculateConsistency tests the calculateConsistency function
func TestCalculateConsistency(t *testing.T) {
	// Test with nil results
	consistency := calculateConsistency(nil)
	assert.Equal(t, 0.0, consistency)
	
	// Test with empty results
	consistency = calculateConsistency(&quality.ValidationResults{})
	assert.Equal(t, 0.0, consistency)
	
	// Test with results containing consistency rules
	results := &quality.ValidationResults{
		RuleResults: []*quality.RuleResult{
			{
				Rule: &quality.Rule{
					Name: "rule1",
					Type: "unique",
				},
				Passed:  true,
				Score:   1.0,
				Details: "Rule passed",
			},
			{
				Rule: &quality.Rule{
					Name: "rule2",
					Type: "relationship",
				},
				Passed:  false,
				Score:   0.7,
				Details: "Rule partially passed",
			},
			{
				Rule: &quality.Rule{
					Name: "rule3",
					Type: "range", // Not a consistency rule
				},
				Passed:  true,
				Score:   1.0,
				Details: "Rule passed",
			},
		},
	}
	
	consistency = calculateConsistency(results)
	assert.InDelta(t, 0.85, consistency, 0.001) // (1.0 + 0.7) / 2 = 0.85
}

// TestCalculateTimeliness tests the calculateTimeliness function
func TestCalculateTimeliness(t *testing.T) {
	// Test with nil profile
	timeliness := calculateTimeliness(nil)
	assert.Equal(t, 0.0, timeliness)
	
	// Test with profile without timestamp
	timeliness = calculateTimeliness(&quality.Profile{})
	assert.Equal(t, 1.0, timeliness)
	
	// Test with profile with recent timestamp (less than 1 day old)
	profile := &quality.Profile{
		Timestamp: time.Now().Add(-12 * time.Hour),
	}
	timeliness = calculateTimeliness(profile)
	assert.Equal(t, 1.0, timeliness)
	
	// Test with profile with older timestamp (less than 1 week old)
	profile = &quality.Profile{
		Timestamp: time.Now().Add(-72 * time.Hour),
	}
	timeliness = calculateTimeliness(profile)
	assert.Equal(t, 0.8, timeliness)
	
	// Test with profile with old timestamp (less than 1 month old)
	profile = &quality.Profile{
		Timestamp: time.Now().Add(-500 * time.Hour),
	}
	timeliness = calculateTimeliness(profile)
	assert.Equal(t, 0.6, timeliness)
	
	// Test with profile with very old timestamp (more than 1 month old)
	profile = &quality.Profile{
		Timestamp: time.Now().Add(-1000 * time.Hour),
	}
	timeliness = calculateTimeliness(profile)
	assert.Equal(t, 0.4, timeliness)
}

// TestCalculateOverallScore tests the calculateOverallScore function
func TestCalculateOverallScore(t *testing.T) {
	// Create test profile with columns
	profile := &quality.Profile{
		Columns: map[string]*quality.ColumnProfile{
			"col1": {
				Stats: &quality.ColumnStats{
					Count:     100,
					NullCount: 10,
				},
			},
		},
	}
	
	results := &quality.ValidationResults{
		RuleResults: []*quality.RuleResult{
			{
				Rule: &quality.Rule{
					Name: "rule1",
					Type: "range",
				},
				Passed:  true,
				Score:   1.0,
				Details: "Rule passed",
			},
			{
				Rule: &quality.Rule{
					Name: "rule2",
					Type: "unique",
				},
				Passed:  true,
				Score:   1.0,
				Details: "Rule passed",
			},
		},
	}
	
	// Calculate overall score
	overallScore := calculateOverallScore(profile, results)
	
	// Verify score is between 0 and 1
	assert.GreaterOrEqual(t, overallScore, 0.0)
	assert.LessOrEqual(t, overallScore, 1.0)
}
