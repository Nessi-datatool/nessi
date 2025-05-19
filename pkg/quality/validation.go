package quality

// ValidationResults represents the results of data quality validation
type ValidationResults struct {
	// Overall validation status
	Valid bool `json:"valid"`

	// Individual rule results
	RuleResults []*RuleResult `json:"rule_results"`

	// Summary statistics
	TotalRules   int `json:"total_rules"`
	PassedRules  int `json:"passed_rules"`
	FailedRules  int `json:"failed_rules"`
	SkippedRules int `json:"skipped_rules"`
}

// RuleResult represents the result of a single validation rule
type RuleResult struct {
	Rule    *Rule   `json:"rule"`
	Passed  bool    `json:"passed"`
	Score   float64 `json:"score"`
	Details string  `json:"details,omitempty"`
}

// Rule represents a data quality validation rule
type Rule struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Severity    string            `json:"severity"`
	Parameters  map[string]string `json:"parameters,omitempty"`
}
