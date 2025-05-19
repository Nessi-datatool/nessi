package quality

import (
	"github.com/nessi-dev/nessi/pkg/api/types"
)

// ToQualityMetrics converts a Profile and ValidationResults to QualityMetrics
func ToQualityMetrics(profile *Profile, results *ValidationResults) *types.QualityMetrics {
	metrics := &types.QualityMetrics{
		TotalRows:     profile.RowCount,
		LastUpdated:   profile.Timestamp,
		SchemaVersion: "1.0",
	}

	// Initialize column metrics
	metrics.ColumnMetrics = make(map[string]*types.ColumnQualityMetrics)

	// Add column-level metrics
	for colName, colProfile := range profile.Columns {
		if colProfile.Stats == nil {
			continue
		}

		metrics.ColumnMetrics[colName] = &types.ColumnQualityMetrics{
			NullCount:     colProfile.Stats.NullCount,
			DistinctCount: colProfile.Stats.DistinctCount,
			MinValue:      colProfile.Stats.MinValue,
			MaxValue:      colProfile.Stats.MaxValue,
			MeanValue:     colProfile.Stats.MeanValue,
			MedianValue:   colProfile.Stats.MedianValue,
		}
	}

	// Calculate overall metrics
	if results != nil {
		metrics.InvalidRows = int64(results.FailedRules)
		metrics.DataCompleteness = float64(results.PassedRules) / float64(results.TotalRules)
		metrics.DataAccuracy = calculateAccuracy(results)
		metrics.DataConsistency = calculateConsistency(results)
	}

	return metrics
}

// Helper functions to calculate quality scores

func calculateAccuracy(results *ValidationResults) float64 {
	if results == nil || len(results.RuleResults) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalRules int

	for _, result := range results.RuleResults {
		if isAccuracyRule(result.Rule.Type) {
			totalScore += result.Score
			totalRules++
		}
	}

	if totalRules == 0 {
		return 1.0
	}

	return totalScore / float64(totalRules)
}

func calculateConsistency(results *ValidationResults) float64 {
	if results == nil || len(results.RuleResults) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalRules int

	for _, result := range results.RuleResults {
		if isConsistencyRule(result.Rule.Type) {
			totalScore += result.Score
			totalRules++
		}
	}

	if totalRules == 0 {
		return 1.0
	}

	return totalScore / float64(totalRules)
}

func isAccuracyRule(ruleType string) bool {
	accuracyRules := map[string]bool{
		"range":  true,
		"enum":   true,
		"regex":  true,
		"format": true,
	}

	return accuracyRules[ruleType]
}

func isConsistencyRule(ruleType string) bool {
	consistencyRules := map[string]bool{
		"unique":                true,
		"relationship":          true,
		"referential_integrity": true,
	}

	return consistencyRules[ruleType]
}
