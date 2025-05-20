package monitor

import (
	"time"
)

// Config defines the monitor configuration
type Config struct {
	Enabled       bool          `json:"enabled"`
	Interval      time.Duration `json:"interval"`
	MetricsPrefix string        `json:"metrics_prefix"`
	OutputPath    string        `json:"output_path"`
}

// MetricData represents a simple metric (renamed to avoid redeclaration)
type MetricData struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Value       float64           `json:"value"`
	Labels      map[string]string `json:"labels"`
	Timestamp   time.Time         `json:"timestamp"`
}
