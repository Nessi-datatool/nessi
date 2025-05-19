package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/nessi-dev/nessi/pkg/quality/profile"
	"github.com/nessi-dev/nessi/pkg/quality/rules"
)

// sendJSONResponse sends a JSON response with the given status code and data
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DataQualityResponse represents the response for data quality endpoints
type DataQualityResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// RuleValidationResult represents the result of rule validation
type RuleValidationResult struct {
	RuleID      string    `json:"rule_id"`
	RuleName    string    `json:"rule_name"`
	Severity    string    `json:"severity"`
	Passed      bool      `json:"passed"`
	ErrorCount  int       `json:"error_count"`
	Details     []string  `json:"details,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	ExecutionID string    `json:"execution_id"`
}

// ProfileSummary represents a summary of a data profile
type ProfileSummary struct {
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	NullPercent    float64 `json:"null_percent"`
	QualityScore   float64 `json:"quality_score,omitempty"`
	AnomalyCount   int     `json:"anomaly_count"`
	PatternCount   int     `json:"pattern_count"`
	RecommendedType string  `json:"recommended_type,omitempty"`
}

// handleGetProfiles handles requests to get data profiles
func (d *Dashboard) handleGetProfiles(w http.ResponseWriter, r *http.Request) {
	// Get table path from query parameters
	tablePath := r.URL.Query().Get("table")
	if tablePath == "" {
		sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
			Success: false,
			Message: "Missing table parameter",
		})
		return
	}

	// Use injected profiler if available, else fallback to real implementation
	var profiler Profiler
	if d.profiler != nil {
		profiler = d.profiler
	} else {
		profiler = profile.NewAdvancedProfiler(tablePath)
	}

	// Generate profiles
	profiles, err := profiler.GenerateProfile()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, DataQualityResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to generate profiles: %v", err),
		})
		return
	}

	// Convert to enhanced profiles
	var enhancedProfiles []*profile.EnhancedProfile
	for _, p := range profiles {
		// In a real implementation, we would get the actual values
		// Here we're just creating a placeholder enhanced profile
		enhancedProfiles = append(enhancedProfiles, &profile.EnhancedProfile{
			Profile: p,
			QualityScore: &profile.DataQualityScore{
				Overall: 0.95, // Placeholder
			},
		})
	}

	sendJSONResponse(w, http.StatusOK, DataQualityResponse{
		Success: true,
		Data:    enhancedProfiles,
	})
}

// handleGetProfileSummaries handles requests to get profile summaries
func (d *Dashboard) handleGetProfileSummaries(w http.ResponseWriter, r *http.Request) {
	// Get table path from query parameters
	tablePath := r.URL.Query().Get("table")
	if tablePath == "" {
		sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
			Success: false,
			Message: "Missing table parameter",
		})
		return
	}

	// Create advanced profiler
	profiler := profile.NewAdvancedProfiler(tablePath)
	
	// Generate profiles
	profiles, err := profiler.GenerateProfile()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, DataQualityResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to generate profiles: %v", err),
		})
		return
	}

	// Convert to summaries
	var summaries []ProfileSummary
	for _, p := range profiles {
		summary := ProfileSummary{
			Name:         p.Name,
			Type:         p.Type,
			NullPercent:  p.NullPercent,
			AnomalyCount: len(p.Anomalies),
			PatternCount: len(p.Patterns),
		}
		summaries = append(summaries, summary)
	}

	sendJSONResponse(w, http.StatusOK, DataQualityResponse{
		Success: true,
		Data:    summaries,
	})
}

// handleGetRules handles requests to get validation rules
func (d *Dashboard) handleGetRules(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, we would load rules from storage
	// Here we're just creating some example rules
	nullCheckRule := rules.NewNullCheckRule([]string{"id", "name"}, rules.RuleMetadata{
		ID:          "null_check_1",
		Name:        "Required Fields Check",
		Description: "Checks that required fields are not null",
		Severity:    "error",
		Tags:        []string{"data_quality", "nulls"},
	})

	rangeRule := rules.NewRangeCheckRule(map[string]rules.RangeConfig{
		"age": {Min: 0, Max: 120},
	}, rules.RuleMetadata{
		ID:          "range_check_1",
		Name:        "Age Range Check",
		Description: "Checks that age is within valid range",
		Severity:    "warning",
		Tags:        []string{"data_quality", "range"},
	})

	// Create a list of rules
	rulesList := []rules.Rule{nullCheckRule, rangeRule}

	sendJSONResponse(w, http.StatusOK, DataQualityResponse{
		Success: true,
		Data:    rulesList,
	})
}

// handleValidateRules handles requests to validate data against rules
func (d *Dashboard) handleValidateRules(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, DataQualityResponse{
			Success: false,
			Message: "Only POST method is allowed",
		})
		return
	}

	// Get table parameter
	tablePath := r.URL.Query().Get("table")
	if tablePath == "" {
		sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
			Success: false,
			Message: "Missing table parameter",
		})
		return
	}

	// Use injected profiler if available, else fallback
	var profiler Profiler
	if d.profiler != nil {
		profiler = d.profiler
	} else {
		profiler = profile.NewAdvancedProfiler(tablePath)
	}
	profiles, err := profiler.GenerateProfile()
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, DataQualityResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to generate profiles: %v", err),
		})
		return
	}


	// Example: Validate each profile (mocked)
	var results []RuleValidationResult
	for range profiles {
		// In a real implementation, would validate records
		results = append(results, RuleValidationResult{
			RuleID:      "null_check_1",
			RuleName:    "Required Fields Check",
			Severity:    "error",
			Passed:      true,
			ErrorCount:  0,
			Timestamp:   time.Now(),
			ExecutionID: "exec-" + strconv.FormatInt(time.Now().Unix(), 10),
		})
	}

	sendJSONResponse(w, http.StatusOK, DataQualityResponse{
		Success: true,
		Data:    results,
	})
}

// handleGetRuleHistory handles requests to get rule execution history
func (d *Dashboard) handleGetRuleHistory(w http.ResponseWriter, r *http.Request) {
	// Get rule ID from query parameters
	ruleID := r.URL.Query().Get("rule_id")
	if ruleID == "" {
		sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
			Success: false,
			Message: "Missing rule_id parameter",
		})
		return
	}

	// Get limit from query parameters (default to 10)
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
				Success: false,
				Message: "Invalid limit parameter",
			})
			return
		}
	}

	// In a real implementation, we would load rule history from storage
	// Here we're just creating some example history
	var history []rules.RuleExecutionRecord
	
	// Create example history
	now := time.Now()
	for i := 0; i < limit; i++ {
		history = append(history, rules.RuleExecutionRecord{
			ExecutionID:    fmt.Sprintf("exec-%d", i),
			RuleID:         ruleID,
			Timestamp:      now.Add(time.Duration(-i) * time.Hour),
			RecordsChecked: 1000,
			Failures:       int64(i * 2), // Increasing failures over time
			ExecutionTimeMs: 150 + int64(i*10),
			DatasetID:      "dataset-1",
		})
	}

	sendJSONResponse(w, http.StatusOK, DataQualityResponse{
		Success: true,
		Data:    history,
	})
}

// handleGetExecutionTrends handles requests to get rule execution trends
func (d *Dashboard) handleGetExecutionTrends(w http.ResponseWriter, r *http.Request) {
	// Get rule ID from query parameters
	ruleID := r.URL.Query().Get("rule_id")
	if ruleID == "" {
		sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
			Success: false,
			Message: "Missing rule_id parameter",
		})
		return
	}

	// Get days from query parameters (default to 7)
	daysStr := r.URL.Query().Get("days")
	days := 7
	if daysStr != "" {
		var err error
		days, err = strconv.Atoi(daysStr)
		if err != nil || days <= 0 {
			sendJSONResponse(w, http.StatusBadRequest, DataQualityResponse{
				Success: false,
				Message: "Invalid days parameter",
			})
			return
		}
	}

	// In a real implementation, we would calculate trends from rule history
	// Here we're just creating some example trends
	var trends []rules.ExecutionTrend
	
	// Create example trends
	now := time.Now()
	for i := 0; i < days; i++ {
		trends = append(trends, rules.ExecutionTrend{
			Date:            now.Add(time.Duration(-i) * 24 * time.Hour),
			ExecutionCount:  5 + i%3,
			FailureRate:     0.05 + float64(i%5)*0.01,
			AverageFailures: float64(2 + i%4),
		})
	}

	sendJSONResponse(w, http.StatusOK, DataQualityResponse{
		Success: true,
		Data:    trends,
	})
}

// registerDataQualityHandlers registers data quality handlers
func (d *Dashboard) registerDataQualityHandlers() {
	d.mux.HandleFunc("/api/profiles", d.handleGetProfiles)
	d.mux.HandleFunc("/api/profiles/summary", d.handleGetProfileSummaries)
	d.mux.HandleFunc("/api/rules", d.handleGetRules)
	d.mux.HandleFunc("/api/rules/validate", d.handleValidateRules)
	d.mux.HandleFunc("/api/rules/history", d.handleGetRuleHistory)
	d.mux.HandleFunc("/api/rules/trends", d.handleGetExecutionTrends)
}
