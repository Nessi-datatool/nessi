package freshness

import (
	"fmt"
	"time"
)

// SLALevel represents the severity level for SLA violations
type SLALevel string

const (
	// SLALevelInfo represents an informational SLA level
	SLALevelInfo SLALevel = "info"
	
	// SLALevelWarning represents a warning SLA level
	SLALevelWarning SLALevel = "warning"
	
	// SLALevelCritical represents a critical SLA level
	SLALevelCritical SLALevel = "critical"
)

// SLAConfig represents the SLA configuration for a table
type SLAConfig struct {
	// TableName is the name of the table
	TableName string `json:"table_name"`
	
	// TablePath is the path to the table
	TablePath string `json:"table_path"`
	
	// ExpectedFrequency is the expected update frequency for the table
	ExpectedFrequency time.Duration `json:"expected_frequency"`
	
	// WarningThreshold is the threshold for warning alerts (as a percentage of ExpectedFrequency)
	// For example, if ExpectedFrequency is 1h and WarningThreshold is 150, a warning will be triggered
	// if the table hasn't been updated for 1.5 hours
	WarningThreshold int `json:"warning_threshold"`
	
	// CriticalThreshold is the threshold for critical alerts (as a percentage of ExpectedFrequency)
	// For example, if ExpectedFrequency is 1h and CriticalThreshold is 200, a critical alert will be triggered
	// if the table hasn't been updated for 2 hours
	CriticalThreshold int `json:"critical_threshold"`
	
	// Enabled indicates whether SLA monitoring is enabled for this table
	Enabled bool `json:"enabled"`
	
	// Description provides additional context about the SLA
	Description string `json:"description,omitempty"`
	
	// Tags are optional tags for filtering and grouping
	Tags []string `json:"tags,omitempty"`
	
	// GrafanaDashboardURL is an optional URL to a Grafana dashboard for this table
	GrafanaDashboardURL string `json:"grafana_dashboard_url,omitempty"`
}

// Validate validates the SLA configuration
func (c *SLAConfig) Validate() error {
	if c.TableName == "" {
		return fmt.Errorf("table name cannot be empty")
	}
	
	if c.TablePath == "" {
		return fmt.Errorf("table path cannot be empty")
	}
	
	if c.ExpectedFrequency <= 0 {
		return fmt.Errorf("expected frequency must be greater than zero")
	}
	
	if c.WarningThreshold <= 0 {
		return fmt.Errorf("warning threshold must be greater than zero")
	}
	
	if c.CriticalThreshold <= 0 {
		return fmt.Errorf("critical threshold must be greater than zero")
	}
	
	// Warning threshold should be at least 100% (equal to expected frequency)
	if c.WarningThreshold < 100 {
		return fmt.Errorf("warning threshold must be at least 100")
	}
	
	if c.WarningThreshold >= c.CriticalThreshold {
		return fmt.Errorf("warning threshold must be less than critical threshold")
	}
	
	return nil
}

// ParseDuration parses a duration string and returns a time.Duration
func ParseDuration(durationStr string) (time.Duration, error) {
	// Try standard Go duration parsing first
	duration, err := time.ParseDuration(durationStr)
	if err == nil {
		return duration, nil
	}
	
	// Try custom formats with numeric prefixes
	var value int
	var unit string
	
	_, err = fmt.Sscanf(durationStr, "%d%s", &value, &unit)
	if err == nil {
		switch unit {
		case "d":
			return time.Duration(value) * 24 * time.Hour, nil
		case "w":
			return time.Duration(value) * 7 * 24 * time.Hour, nil
		case "mo":
			return time.Duration(value) * 30 * 24 * time.Hour, nil
		}
	}
	
	// Try custom formats without numeric prefixes
	switch durationStr {
	case "hourly":
		return time.Hour, nil
	case "daily":
		return 24 * time.Hour, nil
	case "weekly":
		return 7 * 24 * time.Hour, nil
	case "monthly":
		return 30 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}
}

// FormatDuration formats a duration in a human-readable format
func FormatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours", int(duration.Hours()))
	} else if duration < 7*24*time.Hour {
		return fmt.Sprintf("%d days", int(duration.Hours()/24))
	} else if duration < 30*24*time.Hour {
		return fmt.Sprintf("%d weeks", int(duration.Hours()/(24*7)))
	} else {
		return fmt.Sprintf("%d months", int(duration.Hours()/(24*30)))
	}
}
