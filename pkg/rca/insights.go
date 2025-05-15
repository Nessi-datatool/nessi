package rca

import (
	"sort"
	"time"
)

// Insights represents aggregated insights from RCA results
type Insights struct {
	// Distribution of root causes by type
	RootCauseDistribution map[string]int `json:"root_cause_distribution"`
	
	// Count of incidents by affected table
	AffectedTablesCount map[string]int `json:"affected_tables_count"`
	
	// Most common root causes
	CommonRootCauses []*CommonRootCause `json:"common_root_causes"`
	
	// Analysis period
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	
	// Total number of analyses
	TotalAnalyses int `json:"total_analyses"`
}

// CommonRootCause represents a common root cause with occurrence count
type CommonRootCause struct {
	Description string  `json:"description"`
	Details     string  `json:"details,omitempty"`
	Count       int     `json:"count"`
	AvgConfidence float64 `json:"avg_confidence"`
}

// GenerateInsights generates insights from a list of RCA results
func GenerateInsights(results []*RCAResult) *Insights {
	if len(results) == 0 {
		return &Insights{
			RootCauseDistribution: make(map[string]int),
			AffectedTablesCount:   make(map[string]int),
			CommonRootCauses:      []*CommonRootCause{},
			StartTime:             time.Now(),
			EndTime:               time.Now(),
			TotalAnalyses:         0,
		}
	}
	
	insights := &Insights{
		RootCauseDistribution: make(map[string]int),
		AffectedTablesCount:   make(map[string]int),
		CommonRootCauses:      []*CommonRootCause{},
		StartTime:             results[0].AnalysisTime,
		EndTime:               results[0].AnalysisTime,
		TotalAnalyses:         len(results),
	}
	
	// Map to track unique root causes
	rootCauseMap := make(map[string]*CommonRootCause)
	
	// Process each result
	for _, result := range results {
		// Update time range
		if result.AnalysisTime.Before(insights.StartTime) {
			insights.StartTime = result.AnalysisTime
		}
		if result.AnalysisTime.After(insights.EndTime) {
			insights.EndTime = result.AnalysisTime
		}
		
		// Process primary root cause
		if result.PrimaryRootCause != nil {
			// Update distribution
			causeType := result.PrimaryRootCause.Type
			insights.RootCauseDistribution[causeType]++
			
			// Update common causes
			key := result.PrimaryRootCause.Description
			if existing, ok := rootCauseMap[key]; ok {
				existing.Count++
				existing.AvgConfidence = (existing.AvgConfidence*float64(existing.Count-1) + result.PrimaryRootCause.Confidence) / float64(existing.Count)
			} else {
				rootCauseMap[key] = &CommonRootCause{
					Description:   result.PrimaryRootCause.Description,
					Details:       result.PrimaryRootCause.Details,
					Count:         1,
					AvgConfidence: result.PrimaryRootCause.Confidence,
				}
			}
		}
		
		// Process affected tables
		for _, table := range result.AffectedTables {
			insights.AffectedTablesCount[table]++
		}
	}
	
	// Convert map to slice for common root causes
	for _, cause := range rootCauseMap {
		insights.CommonRootCauses = append(insights.CommonRootCauses, cause)
	}
	
	// Sort common root causes by count (descending)
	sort.Slice(insights.CommonRootCauses, func(i, j int) bool {
		return insights.CommonRootCauses[i].Count > insights.CommonRootCauses[j].Count
	})
	
	// Limit to top 10
	if len(insights.CommonRootCauses) > 10 {
		insights.CommonRootCauses = insights.CommonRootCauses[:10]
	}
	
	return insights
}
