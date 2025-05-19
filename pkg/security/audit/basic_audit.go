package audit

// BasicAuditLogger provides audit logging functionality for Nessi

import (
	"log"
	"time"
)

// BasicAuditEntry represents a simple audit log entry
type BasicAuditEntry struct {
	Timestamp time.Time
	User      string
	Action    string
	Resource  string
	Details   string
}

// BasicAuditLogger implements a simplified version of the AuditLogger interface
type BasicAuditLogger struct {
	enabled bool
}

// NewBasicAuditLogger creates a new BasicAuditLogger
func NewBasicAuditLogger(enabled bool) *BasicAuditLogger {
	return &BasicAuditLogger{
		enabled: enabled,
	}
}

// LogAction logs an action performed by a user
func (l *BasicAuditLogger) LogAction(user, action, resource, details string) {
	if !l.enabled {
		return
	}

	entry := BasicAuditEntry{
		Timestamp: time.Now(),
		User:      user,
		Action:    action,
		Resource:  resource,
		Details:   details,
	}

	// In OSS version, just log to standard logger
	log.Printf("[AUDIT] User: %s, Action: %s, Resource: %s, Details: %s",
		entry.User, entry.Action, entry.Resource, entry.Details)
}

// Enable enables audit logging
func (l *BasicAuditLogger) Enable() {
	l.enabled = true
}

// Disable disables audit logging
func (l *BasicAuditLogger) Disable() {
	l.enabled = false
}

// IsEnabled returns whether audit logging is enabled
func (l *BasicAuditLogger) IsEnabled() bool {
	return l.enabled
}
