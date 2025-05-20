package profile

// This file contains shared type definitions to avoid duplicate declarations

// DataQualityScore represents the quality scores for a data profile
type DataQualityScore struct {
	Overall         float64 `json:"overall"`
	Completeness    float64 `json:"completeness"`
	Consistency     float64 `json:"consistency"`
	Accuracy        float64 `json:"accuracy"`
	Uniqueness      float64 `json:"uniqueness"`
	RecommendedType string  `json:"recommended_type,omitempty"`
}

// EnhancedProfile combines a profile with quality scores and additional analysis
type EnhancedProfile struct {
	Profile           *Profile              `json:"profile"`
	QualityScore      *DataQualityScore     `json:"quality_score,omitempty"`
	Distribution      *DistributionAnalysis `json:"distribution,omitempty"`
	DetailedPatterns  []PatternInfo         `json:"detailed_patterns,omitempty"`
	FormatConsistency float64               `json:"format_consistency,omitempty"`
	TypeConsistency   float64               `json:"type_consistency,omitempty"`
}

// PatternInfo represents detailed information about a detected pattern
type PatternInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Confidence  float64  `json:"confidence"`
	Examples    []string `json:"examples"`
	Regex       string   `json:"regex,omitempty"`
}

// ValueCount represents a value and its count
type ValueCount struct {
	Value interface{} `json:"value"`
	Count int         `json:"count"`
}
