package monitoring

import (
	"time"
)

// Alert represents a monitoring alert
type Alert struct {
	Name      string
	Severity  string
	Message   string
	Timestamp time.Time
	Metadata  map[string]string
}
