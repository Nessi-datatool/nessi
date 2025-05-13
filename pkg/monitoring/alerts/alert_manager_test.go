package alerts

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlertManager(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alert-manager-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	manager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Test CreateAlert
	t.Run("CreateAlert", func(t *testing.T) {
		// Create alert
		alert := &Alert{
			ID:          uuid.New().String(),
			Name:        "Test Alert",
			Description: "This is a test alert",
			Type:        TypeQuality,
			Severity:    SeverityWarning,
			Source:      "test-table",
			Timestamp:   time.Now(),
			Value:       95.5,
			Threshold:   90.0,
			ComparisonOperator: ">",
			Labels: map[string]string{
				"environment": "test",
				"component":   "data-quality",
			},
			Annotations: map[string]string{
				"summary": "Data quality score exceeded threshold",
				"impact":  "Medium",
			},
		}

		err := manager.CreateAlert(alert)
		assert.NoError(t, err)

		// Verify alert file was created
		alertPath := filepath.Join(tempDir, "alerts", alert.ID+".json")
		assert.FileExists(t, alertPath)

		// Test GetAlert
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Equal(t, alert.ID, retrievedAlert.ID)
		assert.Equal(t, alert.Name, retrievedAlert.Name)
		assert.Equal(t, alert.Description, retrievedAlert.Description)
		assert.Equal(t, alert.Type, retrievedAlert.Type)
		assert.Equal(t, alert.Severity, retrievedAlert.Severity)
		assert.Equal(t, StatusActive, retrievedAlert.Status) // Default status should be active
		assert.Equal(t, alert.Source, retrievedAlert.Source)
		assert.Equal(t, alert.Value, retrievedAlert.Value)
		assert.Equal(t, alert.Threshold, retrievedAlert.Threshold)
		assert.Equal(t, alert.ComparisonOperator, retrievedAlert.ComparisonOperator)
		assert.Equal(t, alert.Labels, retrievedAlert.Labels)
		assert.Equal(t, alert.Annotations, retrievedAlert.Annotations)

		// Test GetAlerts
		alerts := manager.GetAlerts()
		assert.Len(t, alerts, 1)
		assert.Equal(t, alert.ID, alerts[0].ID)

		// Test GetAlertsByStatus
		activeAlerts := manager.GetAlertsByStatus(StatusActive)
		assert.Len(t, activeAlerts, 1)
		assert.Equal(t, alert.ID, activeAlerts[0].ID)

		// Test GetAlertsByType
		qualityAlerts := manager.GetAlertsByType(TypeQuality)
		assert.Len(t, qualityAlerts, 1)
		assert.Equal(t, alert.ID, qualityAlerts[0].ID)

		// Test GetAlertsBySeverity
		warningAlerts := manager.GetAlertsBySeverity(SeverityWarning)
		assert.Len(t, warningAlerts, 1)
		assert.Equal(t, alert.ID, warningAlerts[0].ID)
	})

	// Test UpdateAlert
	t.Run("UpdateAlert", func(t *testing.T) {
		// Create alert
		alert := &Alert{
			ID:          uuid.New().String(),
			Name:        "Update Test Alert",
			Description: "This is an alert to test updating",
			Type:        TypeAnomaly,
			Severity:    SeverityInfo,
			Source:      "test-metric",
			Timestamp:   time.Now(),
		}

		err := manager.CreateAlert(alert)
		assert.NoError(t, err)

		// Update alert
		alert.Description = "Updated description"
		alert.Severity = SeverityCritical
		alert.Value = 100.0
		alert.Threshold = 50.0
		alert.ComparisonOperator = ">"

		err = manager.UpdateAlert(alert)
		assert.NoError(t, err)

		// Verify alert was updated
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", retrievedAlert.Description)
		assert.Equal(t, SeverityCritical, retrievedAlert.Severity)
		assert.Equal(t, 100.0, retrievedAlert.Value)
		assert.Equal(t, 50.0, retrievedAlert.Threshold)
		assert.Equal(t, ">", retrievedAlert.ComparisonOperator)
	})

	// Test AcknowledgeAlert
	t.Run("AcknowledgeAlert", func(t *testing.T) {
		// Create alert
		alert := &Alert{
			ID:          uuid.New().String(),
			Name:        "Acknowledge Test Alert",
			Description: "This is an alert to test acknowledging",
			Type:        TypeSystem,
			Severity:    SeverityCritical,
			Source:      "test-system",
			Timestamp:   time.Now(),
		}

		err := manager.CreateAlert(alert)
		assert.NoError(t, err)

		// Acknowledge alert
		err = manager.AcknowledgeAlert(alert.ID, "test-user")
		assert.NoError(t, err)

		// Verify alert was acknowledged
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Equal(t, StatusAcknowledged, retrievedAlert.Status)
		assert.Equal(t, "test-user", retrievedAlert.AcknowledgedBy)
		assert.NotNil(t, retrievedAlert.AcknowledgedAt)
	})

	// Test ResolveAlert
	t.Run("ResolveAlert", func(t *testing.T) {
		// Create alert
		alert := &Alert{
			ID:          uuid.New().String(),
			Name:        "Resolve Test Alert",
			Description: "This is an alert to test resolving",
			Type:        TypeCustom,
			Severity:    SeverityWarning,
			Source:      "test-custom",
			Timestamp:   time.Now(),
		}

		err := manager.CreateAlert(alert)
		assert.NoError(t, err)

		// Resolve alert
		err = manager.ResolveAlert(alert.ID)
		assert.NoError(t, err)

		// Verify alert was resolved
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Equal(t, StatusResolved, retrievedAlert.Status)
		assert.NotNil(t, retrievedAlert.ResolvedAt)
	})

	// Test SilenceAlert
	t.Run("SilenceAlert", func(t *testing.T) {
		// Create alert
		alert := &Alert{
			ID:          uuid.New().String(),
			Name:        "Silence Test Alert",
			Description: "This is an alert to test silencing",
			Type:        TypeQuality,
			Severity:    SeverityWarning,
			Source:      "test-silence",
			Timestamp:   time.Now(),
		}

		err := manager.CreateAlert(alert)
		assert.NoError(t, err)

		// Silence alert
		silenceDuration := 1 * time.Hour
		err = manager.SilenceAlert(alert.ID, "test-user", "Testing silence functionality", silenceDuration)
		assert.NoError(t, err)

		// Verify alert was silenced
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Equal(t, StatusSilenced, retrievedAlert.Status)
		assert.Equal(t, "test-user", retrievedAlert.SilencedBy)
		assert.Equal(t, "Testing silence functionality", retrievedAlert.SilenceReason)
		assert.NotNil(t, retrievedAlert.SilencedUntil)
		
		// Verify silence duration
		expectedSilenceEnd := time.Now().Add(silenceDuration)
		assert.WithinDuration(t, expectedSilenceEnd, *retrievedAlert.SilencedUntil, 5*time.Second)
	})

	// Test DeleteAlert
	t.Run("DeleteAlert", func(t *testing.T) {
		// Create alert
		alert := &Alert{
			ID:          uuid.New().String(),
			Name:        "Delete Test Alert",
			Description: "This is an alert to test deleting",
			Type:        TypeSystem,
			Severity:    SeverityInfo,
			Source:      "test-delete",
			Timestamp:   time.Now(),
		}

		err := manager.CreateAlert(alert)
		assert.NoError(t, err)

		// Delete alert
		err = manager.DeleteAlert(alert.ID)
		assert.NoError(t, err)

		// Verify alert was deleted
		_, err = manager.GetAlert(alert.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "alert not found")

		// Verify alert file was deleted
		alertPath := filepath.Join(tempDir, "alerts", alert.ID+".json")
		assert.NoFileExists(t, alertPath)
	})

	// Test CreateRule
	t.Run("CreateRule", func(t *testing.T) {
		// Create rule
		rule := &AlertRule{
			ID:                 uuid.New().String(),
			Name:               "Test Rule",
			Description:        "This is a test rule",
			Type:               TypeQuality,
			Severity:           SeverityWarning,
			Source:             "test-table",
			Metric:             "quality_score",
			Threshold:          90.0,
			ComparisonOperator: ">",
			WindowSize:         1 * time.Hour,
			EvaluationInterval: 5 * time.Minute,
			Labels: map[string]string{
				"environment": "test",
				"component":   "data-quality",
			},
			Annotations: map[string]string{
				"summary": "Data quality score exceeded threshold",
				"impact":  "Medium",
			},
			NotificationChannels: []string{"email", "slack"},
			Recipients:           []string{"test@example.com", "#alerts"},
			Enabled:              true,
			CreatedAt:            time.Now(),
			CreatedBy:            "test-user",
		}

		err := manager.CreateRule(rule)
		assert.NoError(t, err)

		// Verify rule file was created
		rulePath := filepath.Join(tempDir, "rules", rule.ID+".json")
		assert.FileExists(t, rulePath)

		// Test GetRule
		retrievedRule, err := manager.GetRule(rule.ID)
		assert.NoError(t, err)
		assert.Equal(t, rule.ID, retrievedRule.ID)
		assert.Equal(t, rule.Name, retrievedRule.Name)
		assert.Equal(t, rule.Description, retrievedRule.Description)
		assert.Equal(t, rule.Type, retrievedRule.Type)
		assert.Equal(t, rule.Severity, retrievedRule.Severity)
		assert.Equal(t, rule.Source, retrievedRule.Source)
		assert.Equal(t, rule.Metric, retrievedRule.Metric)
		assert.Equal(t, rule.Threshold, retrievedRule.Threshold)
		assert.Equal(t, rule.ComparisonOperator, retrievedRule.ComparisonOperator)
		assert.Equal(t, rule.WindowSize, retrievedRule.WindowSize)
		assert.Equal(t, rule.EvaluationInterval, retrievedRule.EvaluationInterval)
		assert.Equal(t, rule.Labels, retrievedRule.Labels)
		assert.Equal(t, rule.Annotations, retrievedRule.Annotations)
		assert.Equal(t, rule.NotificationChannels, retrievedRule.NotificationChannels)
		assert.Equal(t, rule.Recipients, retrievedRule.Recipients)
		assert.Equal(t, rule.Enabled, retrievedRule.Enabled)
		assert.Equal(t, rule.CreatedBy, retrievedRule.CreatedBy)

		// Test GetRules
		rules := manager.GetRules()
		assert.Len(t, rules, 1)
		assert.Equal(t, rule.ID, rules[0].ID)

		// Test GetRulesByType
		qualityRules := manager.GetRulesByType(TypeQuality)
		assert.Len(t, qualityRules, 1)
		assert.Equal(t, rule.ID, qualityRules[0].ID)

		// Test GetEnabledRules
		enabledRules := manager.GetEnabledRules()
		assert.Len(t, enabledRules, 1)
		assert.Equal(t, rule.ID, enabledRules[0].ID)
	})

	// Test UpdateRule
	t.Run("UpdateRule", func(t *testing.T) {
		// Create rule
		rule := &AlertRule{
			ID:                 uuid.New().String(),
			Name:               "Update Test Rule",
			Description:        "This is a rule to test updating",
			Type:               TypeAnomaly,
			Severity:           SeverityInfo,
			Source:             "test-metric",
			Metric:             "anomaly_score",
			Threshold:          0.8,
			ComparisonOperator: ">",
			WindowSize:         30 * time.Minute,
			EvaluationInterval: 1 * time.Minute,
			Enabled:            true,
			CreatedAt:          time.Now(),
			CreatedBy:          "test-user",
		}

		err := manager.CreateRule(rule)
		assert.NoError(t, err)

		// Update rule
		rule.Description = "Updated description"
		rule.Severity = SeverityCritical
		rule.Threshold = 0.9
		rule.WindowSize = 1 * time.Hour
		rule.Enabled = false
		rule.UpdatedBy = "admin-user"

		err = manager.UpdateRule(rule)
		assert.NoError(t, err)

		// Verify rule was updated
		retrievedRule, err := manager.GetRule(rule.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated description", retrievedRule.Description)
		assert.Equal(t, SeverityCritical, retrievedRule.Severity)
		assert.Equal(t, 0.9, retrievedRule.Threshold)
		assert.Equal(t, 1*time.Hour, retrievedRule.WindowSize)
		assert.Equal(t, false, retrievedRule.Enabled)
		assert.Equal(t, "admin-user", retrievedRule.UpdatedBy)
	})

	// Test DeleteRule
	t.Run("DeleteRule", func(t *testing.T) {
		// Create rule
		rule := &AlertRule{
			ID:                 uuid.New().String(),
			Name:               "Delete Test Rule",
			Description:        "This is a rule to test deleting",
			Type:               TypeSystem,
			Severity:           SeverityInfo,
			Source:             "test-delete",
			Metric:             "system_metric",
			Threshold:          50.0,
			ComparisonOperator: ">",
			WindowSize:         10 * time.Minute,
			EvaluationInterval: 1 * time.Minute,
			Enabled:            true,
			CreatedAt:          time.Now(),
			CreatedBy:          "test-user",
		}

		err := manager.CreateRule(rule)
		assert.NoError(t, err)

		// Delete rule
		err = manager.DeleteRule(rule.ID)
		assert.NoError(t, err)

		// Verify rule was deleted
		_, err = manager.GetRule(rule.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rule not found")

		// Verify rule file was deleted
		rulePath := filepath.Join(tempDir, "rules", rule.ID+".json")
		assert.NoFileExists(t, rulePath)
	})
}

// MockNotifier is a mock implementation of the Notifier interface for testing
type MockNotifier struct {
	name      string
	sendFunc  func(alert *Alert, recipient string) (*AlertNotification, error)
	sentCount int
}

func NewMockNotifier(name string) *MockNotifier {
	return &MockNotifier{
		name: name,
		sendFunc: func(alert *Alert, recipient string) (*AlertNotification, error) {
			return &AlertNotification{
				ID:        uuid.New().String(),
				AlertID:   alert.ID,
				Channel:   name,
				Recipient: recipient,
				SentAt:    time.Now(),
				Status:    "sent",
			}, nil
		},
	}
}

func (n *MockNotifier) Name() string {
	return n.name
}

func (n *MockNotifier) Send(alert *Alert, recipient string) (*AlertNotification, error) {
	n.sentCount++
	return n.sendFunc(alert, recipient)
}

func TestNotifications(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alert-notifications-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create alert manager
	manager, err := NewAlertManager(tempDir)
	require.NoError(t, err)

	// Register mock notifiers
	emailNotifier := NewMockNotifier("email")
	slackNotifier := NewMockNotifier("slack")
	webhookNotifier := NewMockNotifier("webhook")

	manager.RegisterNotifier(emailNotifier)
	manager.RegisterNotifier(slackNotifier)
	manager.RegisterNotifier(webhookNotifier)

	// Create alert
	alert := &Alert{
		ID:          uuid.New().String(),
		Name:        "Notification Test Alert",
		Description: "This is an alert to test notifications",
		Type:        TypeQuality,
		Severity:    SeverityCritical,
		Source:      "test-notifications",
		Timestamp:   time.Now(),
	}

	err = manager.CreateAlert(alert)
	assert.NoError(t, err)

	// Test SendNotification
	t.Run("SendNotification", func(t *testing.T) {
		// Send email notification
		emailNotification, err := manager.SendNotification(alert.ID, "email", "test@example.com")
		assert.NoError(t, err)
		assert.Equal(t, "email", emailNotification.Channel)
		assert.Equal(t, "test@example.com", emailNotification.Recipient)
		assert.Equal(t, "sent", emailNotification.Status)
		assert.Equal(t, 1, emailNotifier.sentCount)

		// Send Slack notification
		slackNotification, err := manager.SendNotification(alert.ID, "slack", "#alerts")
		assert.NoError(t, err)
		assert.Equal(t, "slack", slackNotification.Channel)
		assert.Equal(t, "#alerts", slackNotification.Recipient)
		assert.Equal(t, "sent", slackNotification.Status)
		assert.Equal(t, 1, slackNotifier.sentCount)

		// Send webhook notification
		webhookNotification, err := manager.SendNotification(alert.ID, "webhook", "webhook-id-123")
		assert.NoError(t, err)
		assert.Equal(t, "webhook", webhookNotification.Channel)
		assert.Equal(t, "webhook-id-123", webhookNotification.Recipient)
		assert.Equal(t, "sent", webhookNotification.Status)
		assert.Equal(t, 1, webhookNotifier.sentCount)

		// Verify notifications were added to the alert
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Len(t, retrievedAlert.Notifications, 3)
	})

	// Test AddNotification
	t.Run("AddNotification", func(t *testing.T) {
		// Create notification
		notification := &AlertNotification{
			ID:        uuid.New().String(),
			AlertID:   alert.ID,
			Channel:   "custom",
			Recipient: "custom-recipient",
			SentAt:    time.Now(),
			Status:    "sent",
		}

		// Add notification
		err := manager.AddNotification(alert.ID, notification)
		assert.NoError(t, err)

		// Verify notification was added
		retrievedAlert, err := manager.GetAlert(alert.ID)
		assert.NoError(t, err)
		assert.Len(t, retrievedAlert.Notifications, 4) // 3 from previous test + 1 from this test

		// Verify notification details
		found := false
		for _, n := range retrievedAlert.Notifications {
			if n.ID == notification.ID {
				found = true
				assert.Equal(t, "custom", n.Channel)
				assert.Equal(t, "custom-recipient", n.Recipient)
				assert.Equal(t, "sent", n.Status)
				break
			}
		}
		assert.True(t, found, "Notification not found in alert")
	})

	// Test error cases
	t.Run("ErrorCases", func(t *testing.T) {
		// Test sending notification for non-existent alert
		_, err := manager.SendNotification("non-existent-id", "email", "test@example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "alert not found")

		// Test sending notification with non-existent notifier
		_, err = manager.SendNotification(alert.ID, "non-existent", "recipient")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "notifier not found")

		// Test adding notification to non-existent alert
		notification := &AlertNotification{
			ID:        uuid.New().String(),
			AlertID:   "non-existent-id",
			Channel:   "email",
			Recipient: "test@example.com",
			SentAt:    time.Now(),
			Status:    "sent",
		}
		err = manager.AddNotification("non-existent-id", notification)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "alert not found")
	})
}
