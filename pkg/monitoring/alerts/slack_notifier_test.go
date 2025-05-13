package alerts

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSlackNotifier tests the SlackNotifier
func TestSlackNotifier(t *testing.T) {
	// Create a test alert
	alert := &Alert{
		ID:                "test-alert-id",
		Name:              "Test Alert",
		Description:       "This is a test alert for Slack notifications",
		Type:              TypeQuality,
		Severity:          SeverityCritical,
		Status:            StatusActive,
		Source:            "test-source",
		Timestamp:         time.Now(),
		LastUpdated:       time.Now(),
		Value:             95.5,
		Threshold:         90.0,
		ComparisonOperator: ">",
		Labels: map[string]string{
			"environment": "test",
			"component":   "data-quality",
		},
		Annotations: map[string]string{
			"summary": "Data quality score exceeded threshold",
			"impact":  "High",
		},
	}

	// Test successful Slack notification
	t.Run("SendSuccess", func(t *testing.T) {
		// Create a test server that simulates Slack API
		var receivedPayload []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request method
			assert.Equal(t, "POST", r.Method)
			
			// Check content type
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			
			// Read request body
			var err error
			receivedPayload, err = io.ReadAll(r.Body)
			require.NoError(t, err)
			
			// Return success response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}))
		defer server.Close()
		
		// Create Slack config
		config := SlackConfig{
			WebhookURL: server.URL,
			Channel:    "#alerts",
			Username:   "Nessi Alert Bot",
			IconEmoji:  ":warning:",
		}
		
		// Create notifier
		notifier := NewSlackNotifier(config)
		
		// Send notification
		recipient := "#data-quality"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify no error occurred
		require.NoError(t, err)
		
		// Verify notification details
		assert.Equal(t, alert.ID, notification.AlertID)
		assert.Equal(t, "slack", notification.Channel)
		assert.Equal(t, recipient, notification.Recipient)
		assert.Equal(t, "sent", notification.Status)
		
		// Parse received payload
		var message SlackMessage
		err = json.Unmarshal(receivedPayload, &message)
		require.NoError(t, err)
		
		// Verify message details
		assert.Contains(t, message.Text, "Alert: Test Alert")
		assert.Contains(t, message.Text, "[CRITICAL]")
		assert.Equal(t, recipient, message.Channel)
		assert.Equal(t, "Nessi Alert Bot", message.Username)
		assert.Equal(t, ":warning:", message.IconEmoji)
		
		// Verify attachment details
		require.Len(t, message.Attachments, 1)
		attachment := message.Attachments[0]
		assert.Equal(t, "Test Alert", attachment.Title)
		assert.Contains(t, attachment.TitleLink, "alerts/test-alert-id")
		assert.Equal(t, "This is a test alert for Slack notifications", attachment.Text)
		assert.Equal(t, "#f44336", attachment.Color) // Red for critical
		
		// Verify fields
		assert.GreaterOrEqual(t, len(attachment.Fields), 6)
		
		// Check for specific fields
		var foundSeverity, foundType, foundValue, foundThreshold bool
		for _, field := range attachment.Fields {
			switch field.Title {
			case "Severity":
				assert.Equal(t, "critical", field.Value)
				foundSeverity = true
			case "Type":
				assert.Equal(t, "quality", field.Value)
				foundType = true
			case "Value":
				assert.Equal(t, "95.50", field.Value)
				foundValue = true
			case "Threshold":
				assert.Equal(t, "90.00", field.Value)
				foundThreshold = true
			}
		}
		
		assert.True(t, foundSeverity, "Severity field not found")
		assert.True(t, foundType, "Type field not found")
		assert.True(t, foundValue, "Value field not found")
		assert.True(t, foundThreshold, "Threshold field not found")
	})
	
	// Test Slack notification failure
	t.Run("SendFailure", func(t *testing.T) {
		// Create a test server that simulates Slack API failure
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Return error response
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal_error"))
		}))
		defer server.Close()
		
		// Create Slack config
		config := SlackConfig{
			WebhookURL: server.URL,
			Channel:    "#alerts",
			Username:   "Nessi Alert Bot",
		}
		
		// Create notifier
		notifier := NewSlackNotifier(config)
		
		// Send notification
		recipient := "#data-quality"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify error occurred
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to send message")
		
		// Verify notification status
		assert.Equal(t, "failed", notification.Status)
		assert.Contains(t, notification.ErrorMessage, "Failed to send message")
	})
	
	// Test with invalid webhook URL
	t.Run("InvalidWebhookURL", func(t *testing.T) {
		// Create Slack config with invalid URL
		config := SlackConfig{
			WebhookURL: "http://invalid-url-that-does-not-exist.example.com",
			Channel:    "#alerts",
			Username:   "Nessi Alert Bot",
		}
		
		// Create notifier
		notifier := NewSlackNotifier(config)
		
		// Send notification
		recipient := "#data-quality"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify error occurred
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to send message")
		
		// Verify notification status
		assert.Equal(t, "failed", notification.Status)
		assert.Contains(t, notification.ErrorMessage, "Failed to send message")
	})
	
	// Test different severity colors
	t.Run("SeverityColors", func(t *testing.T) {
		testCases := []struct {
			severity AlertSeverity
			color    string
		}{
			{SeverityInfo, "#2196f3"},    // Blue
			{SeverityWarning, "#ff9800"}, // Orange
			{SeverityCritical, "#f44336"}, // Red
		}
		
		for _, tc := range testCases {
			// Create a test server
			var receivedPayload []byte
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedPayload, _ = io.ReadAll(r.Body)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			}))
			
			// Create alert with specific severity
			testAlert := *alert
			testAlert.Severity = tc.severity
			
			// Create Slack config
			config := SlackConfig{
				WebhookURL: server.URL,
				Channel:    "#alerts",
			}
			
			// Create notifier and send notification
			notifier := NewSlackNotifier(config)
			_, err := notifier.Send(&testAlert, "#channel")
			require.NoError(t, err)
			
			// Parse received payload
			var message SlackMessage
			err = json.Unmarshal(receivedPayload, &message)
			require.NoError(t, err)
			
			// Verify color based on severity
			require.Len(t, message.Attachments, 1)
			assert.Equal(t, tc.color, message.Attachments[0].Color)
			
			server.Close()
		}
	})
}
