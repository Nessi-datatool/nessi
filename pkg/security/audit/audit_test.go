package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileAuditLogger(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "audit_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test log file path
	logFilePath := filepath.Join(tempDir, "audit.log")

	// Create a new file audit logger
	logger, err := NewFileAuditLogger(logFilePath)
	if err != nil {
		t.Fatalf("Failed to create file audit logger: %v", err)
	}

	// Test logging an entry
	t.Run("LogEntry", func(t *testing.T) {
		// Create an audit entry
		entry := &AuditEntry{
			ID:           "test-1",
			Timestamp:    time.Now(),
			UserID:       "user1",
			ActionType:   CreateAction,
			ResourceType: TableResource,
			ResourceID:   "table1",
			Details: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
			Success:      true,
			ErrorMessage: "",
			IPAddress:    "127.0.0.1",
			UserAgent:    "test-agent",
		}

		// Log the entry
		err := logger.Log(entry)
		if err != nil {
			t.Errorf("Failed to log audit entry: %v", err)
		}

		// Check that the log file exists
		if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
			t.Errorf("Expected log file to exist, but it didn't")
		}
	})

	// Test querying entries
	t.Run("QueryEntries", func(t *testing.T) {
		// Log some more entries
		entries := []*AuditEntry{
			{
				ID:           "test-2",
				Timestamp:    time.Now(),
				UserID:       "user1",
				ActionType:   ReadAction,
				ResourceType: TableResource,
				ResourceID:   "table1",
				Success:      true,
			},
			{
				ID:           "test-3",
				Timestamp:    time.Now(),
				UserID:       "user2",
				ActionType:   UpdateAction,
				ResourceType: RuleResource,
				ResourceID:   "rule1",
				Success:      true,
			},
			{
				ID:           "test-4",
				Timestamp:    time.Now(),
				UserID:       "user1",
				ActionType:   DeleteAction,
				ResourceType: TableResource,
				ResourceID:   "table2",
				Success:      false,
				ErrorMessage: "Permission denied",
			},
		}

		for _, entry := range entries {
			err := logger.Log(entry)
			if err != nil {
				t.Errorf("Failed to log audit entry: %v", err)
			}
		}

		// Query all entries
		allEntries, err := logger.Query(nil, 10, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(allEntries) != 4 {
			t.Errorf("Expected 4 entries, but got %d", len(allEntries))
		}

		// Query with filter by user
		userEntries, err := logger.Query(map[string]interface{}{
			"user_id": "user1",
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(userEntries) != 3 {
			t.Errorf("Expected 3 entries for user1, but got %d", len(userEntries))
		}

		// Query with filter by resource type
		resourceEntries, err := logger.Query(map[string]interface{}{
			"resource_type": string(TableResource),
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(resourceEntries) != 3 {
			t.Errorf("Expected 3 entries for table resource, but got %d", len(resourceEntries))
		}

		// Query with filter by action type
		actionEntries, err := logger.Query(map[string]interface{}{
			"action_type": string(UpdateAction),
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(actionEntries) != 1 {
			t.Errorf("Expected 1 entry for update action, but got %d", len(actionEntries))
		}

		// Query with filter by success
		successEntries, err := logger.Query(map[string]interface{}{
			"success": false,
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(successEntries) != 1 {
			t.Errorf("Expected 1 entry with success=false, but got %d", len(successEntries))
		}

		// Query with multiple filters
		multiFilterEntries, err := logger.Query(map[string]interface{}{
			"user_id":       "user1",
			"resource_type": string(TableResource),
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(multiFilterEntries) != 3 {
			t.Errorf("Expected 3 entries for user1 and table resource, but got %d", len(multiFilterEntries))
		}

		// Query with limit and offset
		limitedEntries, err := logger.Query(nil, 2, 0)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(limitedEntries) != 2 {
			t.Errorf("Expected 2 entries with limit=2, but got %d", len(limitedEntries))
		}

		offsetEntries, err := logger.Query(nil, 2, 2)
		if err != nil {
			t.Errorf("Failed to query audit entries: %v", err)
		}
		if len(offsetEntries) != 2 {
			t.Errorf("Expected 2 entries with limit=2 and offset=2, but got %d", len(offsetEntries))
		}
	})
}

func TestAuditService(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "audit_service_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test log file path
	logFilePath := filepath.Join(tempDir, "audit.log")

	// Create a new file audit logger
	logger, err := NewFileAuditLogger(logFilePath)
	if err != nil {
		t.Fatalf("Failed to create file audit logger: %v", err)
	}

	// Create a new audit service
	service := NewAuditService(logger)

	// Test logging an action
	t.Run("LogAction", func(t *testing.T) {
		// Log an action
		err := service.LogAction(
			"user1",
			CreateAction,
			TableResource,
			"table1",
			map[string]interface{}{
				"schema": "public",
				"columns": []string{"id", "name", "value"},
			},
			true,
			"",
		)
		if err != nil {
			t.Errorf("Failed to log action: %v", err)
		}

		// Query the log to check the action was logged
		entries, err := service.QueryLogs(map[string]interface{}{
			"user_id":       "user1",
			"action_type":   string(CreateAction),
			"resource_type": string(TableResource),
			"resource_id":   "table1",
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query logs: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("Expected 1 entry, but got %d", len(entries))
		}
	})

	// Test logging a config change
	t.Run("LogConfigChange", func(t *testing.T) {
		// Log a config change
		err := service.LogConfigChange(
			"admin",
			"server.port",
			8080,
			9090,
			true,
			"",
		)
		if err != nil {
			t.Errorf("Failed to log config change: %v", err)
		}

		// Query the log to check the config change was logged
		entries, err := service.QueryLogs(map[string]interface{}{
			"user_id":       "admin",
			"action_type":   string(ConfigChangeAction),
			"resource_type": string(ConfigResource),
		}, 10, 0)
		if err != nil {
			t.Errorf("Failed to query logs: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("Expected 1 entry, but got %d", len(entries))
		}

		entry := entries[0]
		if entry.ResourceID != "server.port" {
			t.Errorf("Expected resource_id to be 'server.port', but it was %s", entry.ResourceID)
		}

		// Check the details map directly
		if entry.Details["old_value"] != float64(8080) {
			t.Errorf("Expected old_value to be 8080, but it was %v", entry.Details["old_value"])
		}
		if entry.Details["new_value"] != float64(9090) {
			t.Errorf("Expected new_value to be 9090, but it was %v", entry.Details["new_value"])
		}
	})

	// Test getting recent logs
	t.Run("GetRecentLogs", func(t *testing.T) {
		// Log some more actions
		for i := 0; i < 5; i++ {
			err := service.LogAction(
				"user2",
				ReadAction,
				RuleResource,
				"rule1",
				nil,
				true,
				"",
			)
			if err != nil {
				t.Errorf("Failed to log action: %v", err)
			}
		}

		// Get recent logs
		entries, err := service.GetRecentLogs(3)
		if err != nil {
			t.Errorf("Failed to get recent logs: %v", err)
		}
		if len(entries) != 3 {
			t.Errorf("Expected 3 entries, but got %d", len(entries))
		}
	})

	// Test getting user logs
	t.Run("GetUserLogs", func(t *testing.T) {
		// Get user logs
		entries, err := service.GetUserLogs("user2", 10, 0)
		if err != nil {
			t.Errorf("Failed to get user logs: %v", err)
		}
		if len(entries) != 5 {
			t.Errorf("Expected 5 entries, but got %d", len(entries))
		}

		// Get user logs with limit
		entries, err = service.GetUserLogs("user2", 2, 0)
		if err != nil {
			t.Errorf("Failed to get user logs: %v", err)
		}
		if len(entries) != 2 {
			t.Errorf("Expected 2 entries, but got %d", len(entries))
		}
	})

	// Test getting resource logs
	t.Run("GetResourceLogs", func(t *testing.T) {
		// Get resource logs
		entries, err := service.GetResourceLogs(RuleResource, "rule1", 10, 0)
		if err != nil {
			t.Errorf("Failed to get resource logs: %v", err)
		}
		if len(entries) != 5 {
			t.Errorf("Expected 5 entries, but got %d", len(entries))
		}

		// Get resource logs with limit and offset
		entries, err = service.GetResourceLogs(RuleResource, "rule1", 2, 2)
		if err != nil {
			t.Errorf("Failed to get resource logs: %v", err)
		}
		if len(entries) != 2 {
			t.Errorf("Expected 2 entries, but got %d", len(entries))
		}
	})
}
