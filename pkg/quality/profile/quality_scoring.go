package profile

import (
	"math"
)

// DataQualityScore represents the quality scores for a data profile
type DataQualityScore struct {
	Overall      float64 `json:"overall"`
	Completeness float64 `json:"completeness"`
	Consistency  float64 `json:"consistency"`
	Accuracy     float64 `json:"accuracy"`
	Uniqueness   float64 `json:"uniqueness"`
	// RecommendedType is populated if the detected type differs from the declared type
	RecommendedType string `json:"recommended_type,omitempty"`
}

// EnhancedProfile combines a profile with quality scores and additional analysis
type EnhancedProfile struct {
	Profile          *Profile         `json:"profile"`
	QualityScore     *DataQualityScore `json:"quality_score,omitempty"`
	Distribution     *DistributionAnalysis `json:"distribution,omitempty"`
	DetailedPatterns []PatternInfo    `json:"detailed_patterns,omitempty"`
}

// PatternInfo represents detailed information about a detected pattern
type PatternInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	Examples    []string `json:"examples"`
}

// CalculateQualityScore calculates quality scores for a profile
func CalculateQualityScore(profile *Profile) *DataQualityScore {
	score := &DataQualityScore{}
	
	// Calculate completeness score (based on null percentage)
	score.Completeness = 1.0 - (profile.NullPercent / 100.0)
	
	// Calculate consistency score
	score.Consistency = calculateConsistencyScore(profile)
	
	// Calculate accuracy score
	score.Accuracy = calculateAccuracyScore(profile)
	
	// Calculate uniqueness score
	score.Uniqueness = calculateUniquenessScore(profile)
	
	// Calculate overall score (weighted average)
	score.Overall = (score.Completeness*0.3 + score.Consistency*0.3 + 
		score.Accuracy*0.3 + score.Uniqueness*0.1)
	
	// Check if type recommendation is needed
	score.RecommendedType = getRecommendedType(profile)
	
	return score
}

// calculateConsistencyScore calculates a consistency score based on patterns and anomalies
func calculateConsistencyScore(profile *Profile) float64 {
	// Base consistency score
	baseScore := 1.0
	
	// Reduce score for each anomaly (more weight for more anomalies)
	anomalyPenalty := 0.0
	if len(profile.Anomalies) > 0 {
		// Logarithmic penalty to avoid too harsh penalties for large datasets
		anomalyPenalty = math.Min(0.5, 0.1*math.Log10(float64(len(profile.Anomalies))+1))
	}
	
	// Pattern consistency - more patterns might indicate inconsistent data
	patternPenalty := 0.0
	if len(profile.Patterns) > 1 {
		// Multiple patterns indicate some inconsistency
		patternPenalty = math.Min(0.3, 0.1*float64(len(profile.Patterns)-1))
	}
	
	return math.Max(0.0, baseScore - anomalyPenalty - patternPenalty)
}

// calculateAccuracyScore calculates an accuracy score based on data type and values
func calculateAccuracyScore(profile *Profile) float64 {
	// Base accuracy score
	baseScore := 1.0
	
	// Type-specific accuracy checks
	switch profile.Type {
	case "string":
		// For strings, check for empty strings (not nulls)
		emptyStrings := 0
		for _, vc := range profile.ValueCounts {
			if vc.Value == "" {
				emptyStrings = vc.Count
				break
			}
		}
		
		// Calculate empty string percentage
		emptyPct := 0.0
		if profile.RowCount > 0 {
			emptyPct = float64(emptyStrings) / float64(profile.RowCount)
		}
		
		// Penalize for empty strings
		baseScore -= emptyPct * 0.5
		
	case "integer", "float":
		// For numeric types, check for zeros (might indicate default values)
		zeros := 0
		for _, vc := range profile.ValueCounts {
			if vc.Value == "0" || vc.Value == "0.0" {
				zeros = vc.Count
				break
			}
		}
		
		// Calculate zero percentage
		zeroPct := 0.0
		if profile.RowCount > 0 {
			zeroPct = float64(zeros) / float64(profile.RowCount)
		}
		
		// If more than 50% zeros, might indicate default values
		if zeroPct > 0.5 {
			baseScore -= 0.2
		}
		
	case "date", "timestamp":
		// For dates, check for default dates (e.g., 1970-01-01)
		defaultDates := 0
		for _, vc := range profile.ValueCounts {
			if vc.Value == "1970-01-01" || vc.Value == "0001-01-01" {
				defaultDates = vc.Count
				break
			}
		}
		
		// Calculate default date percentage
		defaultPct := 0.0
		if profile.RowCount > 0 {
			defaultPct = float64(defaultDates) / float64(profile.RowCount)
		}
		
		// Penalize for default dates
		baseScore -= defaultPct * 0.5
	}
	
	// Penalize for anomalies (more weight for more anomalies)
	anomalyPenalty := 0.0
	if len(profile.Anomalies) > 0 {
		// Logarithmic penalty to avoid too harsh penalties for large datasets
		anomalyPenalty = math.Min(0.3, 0.1*math.Log10(float64(len(profile.Anomalies))+1))
	}
	
	return math.Max(0.0, baseScore - anomalyPenalty)
}

// calculateUniquenessScore calculates a uniqueness score based on distinct values
func calculateUniquenessScore(profile *Profile) float64 {
	// If no rows, return 0
	if profile.RowCount == 0 {
		return 0.0
	}
	
	// Calculate percentage of distinct values
	distinctPct := float64(profile.Distinct) / float64(profile.RowCount)
	
	// For ID-like columns, high uniqueness is good
	if isIdLikeColumn(profile.Name) {
		return distinctPct
	}
	
	// For non-ID columns, extremely high uniqueness might be a concern
	// (e.g., might indicate random data or too much granularity)
	if distinctPct > 0.9 && profile.RowCount > 100 {
		return 0.7 // Penalize extremely high cardinality in non-ID columns
	}
	
	return math.Min(1.0, distinctPct + 0.3) // Add a small boost for reasonable uniqueness
}

// isIdLikeColumn checks if a column name suggests it's an ID column
func isIdLikeColumn(name string) bool {
	idPatterns := []string{"id", "key", "code", "uuid", "guid"}
	
	// Convert to lowercase for case-insensitive matching
	lowerName := name
	
	// Check if the name contains any ID-like patterns
	for _, pattern := range idPatterns {
		if pattern == lowerName || 
		   (len(lowerName) > len(pattern) && 
		    (lowerName[:len(pattern)] == pattern || 
		     lowerName[len(lowerName)-len(pattern):] == pattern)) {
			return true
		}
	}
	
	return false
}

// getRecommendedType returns a recommended type if the current type seems incorrect
func getRecommendedType(profile *Profile) string {
	// If no patterns detected, can't make a recommendation
	if len(profile.Patterns) == 0 {
		return ""
	}
	
	// Check for type mismatches
	switch profile.Type {
	case "string":
		// Check if string might be a date
		for _, pattern := range profile.Patterns {
			if pattern == "date" || pattern == "datetime" || pattern == "timestamp" {
				return "date"
			}
		}
		
		// Check if string might be numeric
		numericPatterns := 0
		for _, pattern := range profile.Patterns {
			if pattern == "integer" || pattern == "decimal" {
				numericPatterns++
			}
		}
		
		if numericPatterns > 0 && numericPatterns == len(profile.Patterns) {
			return "numeric"
		}
		
	case "integer":
		// Check if integer might be a boolean
		if profile.Distinct <= 2 {
			// Check if values are 0/1
			hasZero := false
			hasOne := false
			
			for _, vc := range profile.ValueCounts {
				if vc.Value == "0" {
					hasZero = true
				} else if vc.Value == "1" {
					hasOne = true
				}
			}
			
			if hasZero && hasOne && len(profile.ValueCounts) <= 2 {
				return "boolean"
			}
		}
		
	case "float":
		// Check if float might be an integer
		isInteger := true
		for _, vc := range profile.ValueCounts {
			// Check if value has decimal part
			for i := 0; i < len(vc.Value); i++ {
				if vc.Value[i] == '.' {
					// Check if all digits after decimal are zeros
					allZeros := true
					for j := i + 1; j < len(vc.Value); j++ {
						if vc.Value[j] != '0' {
							allZeros = false
							break
						}
					}
					
					if !allZeros {
						isInteger = false
						break
					}
				}
			}
			
			if !isInteger {
				break
			}
		}
		
		if isInteger {
			return "integer"
		}
	}
	
	return ""
}

// GenerateEnhancedProfile creates an enhanced profile with quality scores and detailed analysis
func GenerateEnhancedProfile(profile *Profile) *EnhancedProfile {
	enhanced := &EnhancedProfile{
		Profile: profile,
	}
	
	// Calculate quality scores
	enhanced.QualityScore = CalculateQualityScore(profile)
	
	// Generate distribution analysis if numeric
	if profile.Type == "integer" || profile.Type == "float" {
		enhanced.Distribution = GenerateDistributionAnalysis(profile)
	}
	
	// Generate detailed pattern information
	enhanced.DetailedPatterns = GenerateDetailedPatterns(profile)
	
	return enhanced
}

// GenerateDetailedPatterns creates detailed pattern information from detected patterns
func GenerateDetailedPatterns(profile *Profile) []PatternInfo {
	var patterns []PatternInfo
	
	// Process each detected pattern
	for _, pattern := range profile.Patterns {
		info := PatternInfo{
			Name:        pattern,
			Confidence:  calculatePatternConfidence(profile, pattern),
			Examples:    getPatternExamples(profile, pattern),
		}
		
		// Add description based on pattern type
		info.Description = getPatternDescription(pattern)
		
		patterns = append(patterns, info)
	}
	
	return patterns
}

// calculatePatternConfidence estimates the confidence level for a pattern
func calculatePatternConfidence(profile *Profile, pattern string) float64 {
	// Count values matching the pattern
	matchingCount := 0
	totalCount := 0
	
	// This is a simplified approach - in a real implementation,
	// we would actually check each value against the pattern
	for _, vc := range profile.ValueCounts {
		if matchesPattern(vc.Value, pattern) {
			matchingCount += vc.Count
		}
		totalCount += vc.Count
	}
	
	if totalCount == 0 {
		return 0.0
	}
	
	return float64(matchingCount) / float64(totalCount)
}

// matchesPattern checks if a value matches a pattern (simplified implementation)
func matchesPattern(value, pattern string) bool {
	// In a real implementation, this would use regex or other pattern matching
	// For now, we'll use a simple heuristic
	switch pattern {
	case "email":
		return containsChar(value, '@')
	case "url":
		return startsWithAny(value, "http://", "https://", "www.")
	case "date":
		return containsChar(value, '-') || containsChar(value, '/')
	case "phone":
		return containsChar(value, '-') || containsChar(value, '(')
	case "integer":
		return isNumeric(value) && !containsChar(value, '.')
	case "decimal":
		return isNumeric(value) && containsChar(value, '.')
	default:
		return true // Default to matching
	}
}

// Helper functions for pattern matching
func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func startsWithAny(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

func isNumeric(s string) bool {
	// Simple check for numeric string
	for i := 0; i < len(s); i++ {
		if (s[i] < '0' || s[i] > '9') && s[i] != '.' && s[i] != '-' {
			return false
		}
	}
	return true
}

// getPatternExamples returns example values that match a pattern
func getPatternExamples(profile *Profile, pattern string) []string {
	var examples []string
	
	// Get up to 3 examples
	count := 0
	for _, vc := range profile.ValueCounts {
		if matchesPattern(vc.Value, pattern) {
			examples = append(examples, vc.Value)
			count++
			if count >= 3 {
				break
			}
		}
	}
	
	return examples
}

// getPatternDescription returns a human-readable description of a pattern
func getPatternDescription(pattern string) string {
	descriptions := map[string]string{
		"email":     "Email address format",
		"url":       "Web URL format",
		"date":      "Date format",
		"datetime":  "Date and time format",
		"timestamp": "Timestamp format",
		"phone":     "Phone number format",
		"zipcode":   "Postal/ZIP code format",
		"integer":   "Integer number format",
		"decimal":   "Decimal number format",
		"uuid":      "UUID format",
		"ipv4":      "IPv4 address format",
		"ipv6":      "IPv6 address format",
		"creditcard": "Credit card number format",
	}
	
	if desc, ok := descriptions[pattern]; ok {
		return desc
	}
	
	return "Custom pattern"
}
