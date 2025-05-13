package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ActionType represents the type of action being audited
type ActionType string

const (
	// CreateAction represents a create action
	CreateAction ActionType = "create"
	// ReadAction represents a read action
	ReadAction ActionType = "read"
	// UpdateAction represents an update action
	UpdateAction ActionType = "update"
	// DeleteAction represents a delete action
	DeleteAction ActionType = "delete"
	// ExecuteAction represents an execute action
	ExecuteAction ActionType = "execute"
	// LoginAction represents a login action
	LoginAction ActionType = "login"
	// LogoutAction represents a logout action
	LogoutAction ActionType = "logout"
	// ConfigChangeAction represents a configuration change action
	ConfigChangeAction ActionType = "config_change"
)

// ResourceType represents the type of resource being audited
type ResourceType string

const (
	// TableResource represents a table resource
	TableResource ResourceType = "table"
	// RuleResource represents a rule resource
	RuleResource ResourceType = "rule"
	// UserResource represents a user resource
	UserResource ResourceType = "user"
	// WebhookResource represents a webhook resource
	WebhookResource ResourceType = "webhook"
	// PluginResource represents a plugin resource
	PluginResource ResourceType = "plugin"
	// ConfigResource represents a configuration resource
	ConfigResource ResourceType = "config"
	// ReportResource represents a report resource
	ReportResource ResourceType = "report"
)

// AuditEntry represents an entry in the audit log
type AuditEntry struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	UserID       string                 `json:"user_id"`
	ActionType   ActionType             `json:"action_type"`
	ResourceType ResourceType           `json:"resource_type"`
	ResourceID   string                 `json:"resource_id"`
	Details      map[string]interface{} `json:"details,omitempty"`
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	IPAddress    string                 `json:"ip_address,omitempty"`
	UserAgent    string                 `json:"user_agent,omitempty"`
}

// AuditLogger is an interface for audit logging
type AuditLogger interface {
	Log(entry *AuditEntry) error
	Query(filter map[string]interface{}, limit, offset int) ([]*AuditEntry, error)
}

// FileAuditLogger is an audit logger that logs to a file
type FileAuditLogger struct {
	filePath string
	mu       sync.Mutex
}

// NewFileAuditLogger creates a new file audit logger
func NewFileAuditLogger(filePath string) (*FileAuditLogger, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create audit log directory: %v", err)
	}

	return &FileAuditLogger{
		filePath: filePath,
	}, nil
}

// Log logs an audit entry to the file
func (l *FileAuditLogger) Log(entry *AuditEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Generate a unique ID if not provided
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Set timestamp if not provided
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	// Marshal entry to JSON
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal audit entry: %v", err)
	}

	// Append to file
	f, err := os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open audit log file: %v", err)
	}
	defer f.Close()

	if _, err := f.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write audit entry: %v", err)
	}

	return nil
}

// Query queries the audit log file for entries matching the filter
func (l *FileAuditLogger) Query(filter map[string]interface{}, limit, offset int) ([]*AuditEntry, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Read the entire file
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*AuditEntry{}, nil
		}
		return nil, fmt.Errorf("failed to read audit log file: %v", err)
	}

	// Split into lines
	lines := splitLines(string(data))

	// Parse each line as a JSON object
	var entries []*AuditEntry
	for _, line := range lines {
		if line == "" {
			continue
		}

		var entry AuditEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal audit entry: %v", err)
		}

		// Apply filter
		if matchesFilter(&entry, filter) {
			entries = append(entries, &entry)
		}
	}

	// Apply offset and limit
	if offset >= len(entries) {
		return []*AuditEntry{}, nil
	}

	end := offset + limit
	if end > len(entries) || limit <= 0 {
		end = len(entries)
	}

	return entries[offset:end], nil
}

// splitLines splits a string into lines
func splitLines(s string) []string {
	var lines []string
	var line string
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, line)
			line = ""
		} else {
			line += string(c)
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

// matchesFilter checks if an entry matches a filter
func matchesFilter(entry *AuditEntry, filter map[string]interface{}) bool {
	for key, value := range filter {
		switch key {
		case "user_id":
			if entry.UserID != value.(string) {
				return false
			}
		case "action_type":
			if string(entry.ActionType) != value.(string) {
				return false
			}
		case "resource_type":
			if string(entry.ResourceType) != value.(string) {
				return false
			}
		case "resource_id":
			if entry.ResourceID != value.(string) {
				return false
			}
		case "success":
			if entry.Success != value.(bool) {
				return false
			}
		case "timestamp_after":
			after, ok := value.(time.Time)
			if !ok {
				return false
			}
			if !entry.Timestamp.After(after) {
				return false
			}
		case "timestamp_before":
			before, ok := value.(time.Time)
			if !ok {
				return false
			}
			if !entry.Timestamp.Before(before) {
				return false
			}
		}
	}
	return true
}

// AuditService provides audit logging services
type AuditService struct {
	logger AuditLogger
}

// NewAuditService creates a new audit service
func NewAuditService(logger AuditLogger) *AuditService {
	return &AuditService{
		logger: logger,
	}
}

// LogAction logs an action
func (s *AuditService) LogAction(userID string, actionType ActionType, resourceType ResourceType, resourceID string, details map[string]interface{}, success bool, errorMessage string) error {
	entry := &AuditEntry{
		Timestamp:    time.Now(),
		UserID:       userID,
		ActionType:   actionType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      details,
		Success:      success,
		ErrorMessage: errorMessage,
	}
	return s.logger.Log(entry)
}

// LogConfigChange logs a configuration change
func (s *AuditService) LogConfigChange(userID string, configKey string, oldValue, newValue interface{}, success bool, errorMessage string) error {
	details := map[string]interface{}{
		"config_key": configKey,
		"old_value":  oldValue,
		"new_value":  newValue,
	}
	return s.LogAction(userID, ConfigChangeAction, ConfigResource, configKey, details, success, errorMessage)
}

// QueryLogs queries the audit logs
func (s *AuditService) QueryLogs(filter map[string]interface{}, limit, offset int) ([]*AuditEntry, error) {
	return s.logger.Query(filter, limit, offset)
}

// GetRecentLogs gets recent audit logs
func (s *AuditService) GetRecentLogs(limit int) ([]*AuditEntry, error) {
	return s.logger.Query(nil, limit, 0)
}

// GetUserLogs gets audit logs for a specific user
func (s *AuditService) GetUserLogs(userID string, limit, offset int) ([]*AuditEntry, error) {
	filter := map[string]interface{}{
		"user_id": userID,
	}
	return s.logger.Query(filter, limit, offset)
}

// GetResourceLogs gets audit logs for a specific resource
func (s *AuditService) GetResourceLogs(resourceType ResourceType, resourceID string, limit, offset int) ([]*AuditEntry, error) {
	filter := map[string]interface{}{
		"resource_type": string(resourceType),
		"resource_id":   resourceID,
	}
	return s.logger.Query(filter, limit, offset)
}
