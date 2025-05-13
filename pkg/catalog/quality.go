package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/quality"
)

// QualityMetricsPublisher handles publishing quality metrics to data catalogs
type QualityMetricsPublisher struct {
	manager *CatalogManager
}

// NewQualityMetricsPublisher creates a new quality metrics publisher
func NewQualityMetricsPublisher(manager *CatalogManager) *QualityMetricsPublisher {
	return &QualityMetricsPublisher{
		manager: manager,
	}
}

// PublishQualityMetrics publishes data quality metrics to a data catalog
func (p *QualityMetricsPublisher) PublishQualityMetrics(
	ctx context.Context, 
	catalogName string, 
	database string, 
	table string, 
	profile *quality.Profile, 
	results *quality.ValidationResults,
) error {
	// Get catalog
	catalog, err := p.manager.GetCatalog(catalogName)
	if err != nil {
		return fmt.Errorf("failed to get catalog %s: %w", catalogName, err)
	}
	
	// Convert profile and validation results to QualityMetrics
	metrics := &QualityMetrics{
		OverallScore: calculateOverallScore(profile, results),
		Completeness: calculateCompleteness(profile),
		Accuracy: calculateAccuracy(results),
		Consistency: calculateConsistency(results),
		Timeliness: calculateTimeliness(profile),
		LastUpdated: time.Now(),
	}
	
	// Add rule results
	for _, result := range results.RuleResults {
		metrics.RuleResults = append(metrics.RuleResults, RuleResult{
			RuleName: result.Rule.Name,
			RuleType: result.Rule.Type,
			Passed: result.Passed,
			Score: result.Score,
			Details: result.Details,
		})
	}
	
	// Publish to catalog
	return catalog.PublishQualityMetrics(ctx, database, table, metrics)
}

// PublishQualityMetricsToAll publishes data quality metrics to all connected catalogs
func (p *QualityMetricsPublisher) PublishQualityMetricsToAll(
	ctx context.Context, 
	database string, 
	table string, 
	profile *quality.Profile, 
	results *quality.ValidationResults,
) []error {
	var errors []error
	
	// Get all catalogs
	catalogs := p.manager.ListCatalogs()
	
	// Publish to each catalog
	for _, catalogName := range catalogs {
		err := p.PublishQualityMetrics(ctx, catalogName, database, table, profile, results)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to publish to %s: %w", catalogName, err))
		}
	}
	
	return errors
}

// Helper functions to calculate quality scores

// calculateOverallScore calculates the overall quality score
func calculateOverallScore(profile *quality.Profile, results *quality.ValidationResults) float64 {
	// Calculate weighted average of completeness, accuracy, consistency, and timeliness
	completeness := calculateCompleteness(profile)
	accuracy := calculateAccuracy(results)
	consistency := calculateConsistency(results)
	timeliness := calculateTimeliness(profile)
	
	// Apply weights (can be adjusted based on importance)
	const (
		completenessWeight = 0.25
		accuracyWeight     = 0.35
		consistencyWeight  = 0.25
		timelinessWeight   = 0.15
	)
	
	overallScore := (completeness * completenessWeight) +
		(accuracy * accuracyWeight) +
		(consistency * consistencyWeight) +
		(timeliness * timelinessWeight)
	
	return overallScore
}

// calculateCompleteness calculates the completeness score
func calculateCompleteness(profile *quality.Profile) float64 {
	if profile == nil || len(profile.Columns) == 0 {
		return 0.0
	}
	
	var totalCompleteness float64
	var totalColumns int
	
	for _, col := range profile.Columns {
		if col.Stats == nil {
			continue
		}
		
		// Calculate completeness as (total - null) / total
		if col.Stats.Count > 0 {
			nullCount := col.Stats.NullCount
			completeness := 1.0 - (float64(nullCount) / float64(col.Stats.Count))
			totalCompleteness += completeness
			totalColumns++
		}
	}
	
	if totalColumns == 0 {
		return 0.0
	}
	
	return totalCompleteness / float64(totalColumns)
}

// calculateAccuracy calculates the accuracy score based on validation results
func calculateAccuracy(results *quality.ValidationResults) float64 {
	if results == nil || len(results.RuleResults) == 0 {
		return 0.0
	}
	
	var totalScore float64
	var totalRules int
	
	for _, result := range results.RuleResults {
		// Only consider accuracy-related rules
		if isAccuracyRule(result.Rule.Type) {
			totalScore += result.Score
			totalRules++
		}
	}
	
	if totalRules == 0 {
		return 1.0 // No accuracy rules means perfect accuracy (by default)
	}
	
	return totalScore / float64(totalRules)
}

// isAccuracyRule determines if a rule type is related to accuracy
func isAccuracyRule(ruleType string) bool {
	accuracyRules := map[string]bool{
		"range": true,
		"enum": true,
		"regex": true,
		"format": true,
	}
	
	return accuracyRules[ruleType]
}

// calculateConsistency calculates the consistency score
func calculateConsistency(results *quality.ValidationResults) float64 {
	if results == nil || len(results.RuleResults) == 0 {
		return 0.0
	}
	
	var totalScore float64
	var totalRules int
	
	for _, result := range results.RuleResults {
		// Only consider consistency-related rules
		if isConsistencyRule(result.Rule.Type) {
			totalScore += result.Score
			totalRules++
		}
	}
	
	if totalRules == 0 {
		return 1.0 // No consistency rules means perfect consistency (by default)
	}
	
	return totalScore / float64(totalRules)
}

// isConsistencyRule determines if a rule type is related to consistency
func isConsistencyRule(ruleType string) bool {
	consistencyRules := map[string]bool{
		"unique": true,
		"relationship": true,
		"referential_integrity": true,
	}
	
	return consistencyRules[ruleType]
}

// calculateTimeliness calculates the timeliness score
func calculateTimeliness(profile *quality.Profile) float64 {
	if profile == nil {
		return 0.0
	}
	
	// Check if profile has timestamp information
	if profile.Timestamp.IsZero() {
		return 1.0 // No timestamp means we can't calculate timeliness
	}
	
	// Calculate timeliness based on how recent the data is
	// This is a simple implementation that can be customized
	now := time.Now()
	ageHours := now.Sub(profile.Timestamp).Hours()
	
	// Define thresholds for timeliness
	const (
		perfectTimeliness = 24.0  // Less than 1 day old
		goodTimeliness    = 168.0 // Less than 1 week old
		fairTimeliness    = 720.0 // Less than 1 month old
	)
	
	if ageHours <= perfectTimeliness {
		return 1.0
	} else if ageHours <= goodTimeliness {
		return 0.8
	} else if ageHours <= fairTimeliness {
		return 0.6
	} else {
		return 0.4
	}
}
