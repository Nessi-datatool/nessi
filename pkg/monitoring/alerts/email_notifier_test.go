package alerts

import (
	"errors"
	"net/smtp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSendMail is a mock function for smtp.SendMail
type mockSendMail func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// TestEmailNotifier tests the EmailNotifier
func TestEmailNotifier(t *testing.T) {
	// Create a test alert
	alert := &Alert{
		ID:                "test-alert-id",
		Name:              "Test Alert",
		Description:       "This is a test alert for email notifications",
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

	// Test successful email sending
	t.Run("SendSuccess", func(t *testing.T) {
		// Create a mock SendMail function that always succeeds
		originalSendMail := sendMail
		defer func() { sendMail = originalSendMail }()
		
		var capturedAddr string
		var capturedAuth smtp.Auth
		var capturedFrom string
		var capturedTo []string
		var capturedMsg []byte
		
		// Override the sendMail function to capture parameters and avoid actual SMTP connection
		sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			capturedAddr = addr
			capturedAuth = a
			capturedFrom = from
			capturedTo = to
			capturedMsg = msg
			return nil
		}
		
		// Create email config
		config := EmailConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "alerts@example.com",
			UseHTML:  true,
		}
		
		// Create notifier with mocked sendMail function
		notifier, err := NewEmailNotifier(config)
		require.NoError(t, err)
		
		// Send notification
		recipient := "user@example.com"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify no error occurred
		require.NoError(t, err)
		
		// Verify notification details
		assert.Equal(t, alert.ID, notification.AlertID)
		assert.Equal(t, "email", notification.Channel)
		assert.Equal(t, recipient, notification.Recipient)
		assert.Equal(t, "sent", notification.Status)
		
		// Verify SMTP parameters
		assert.Equal(t, "smtp.example.com:587", capturedAddr)
		assert.NotNil(t, capturedAuth)
		assert.Equal(t, "alerts@example.com", capturedFrom)
		assert.Equal(t, []string{recipient}, capturedTo)
		
		// Verify email content
		emailContent := string(capturedMsg)
		assert.Contains(t, emailContent, "From: alerts@example.com")
		assert.Contains(t, emailContent, "To: user@example.com")
		assert.Contains(t, emailContent, "Subject: [CRITICAL] Test Alert")
		assert.Contains(t, emailContent, "MIME-Version: 1.0")
		assert.Contains(t, emailContent, "Content-Type: text/html")
		assert.Contains(t, emailContent, "<html>")
		assert.Contains(t, emailContent, "Test Alert")
		assert.Contains(t, emailContent, "This is a test alert for email notifications")
		assert.Contains(t, emailContent, "CRITICAL")
	})
	
	// Test email sending with JSON format
	t.Run("SendJSON", func(t *testing.T) {
		// Create a mock SendMail function that always succeeds
		originalSendMail := sendMail
		defer func() { sendMail = originalSendMail }()
		
		var capturedMsg []byte
		
		sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			capturedMsg = msg
			return nil
		}
		
		// Create email config with UseHTML=false
		config := EmailConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "alerts@example.com",
			UseHTML:  false, // Use JSON format
		}
		
		// Create notifier with mocked sendMail function
		notifier, err := NewEmailNotifier(config)
		require.NoError(t, err)
		
		// Send notification
		recipient := "user@example.com"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify no error occurred
		require.NoError(t, err)
		
		// Verify notification details
		assert.Equal(t, "sent", notification.Status)
		
		// Verify email content
		emailContent := string(capturedMsg)
		assert.Contains(t, emailContent, "From: alerts@example.com")
		assert.Contains(t, emailContent, "To: user@example.com")
		assert.Contains(t, emailContent, "Subject: [CRITICAL] Test Alert")
		assert.Contains(t, emailContent, "Content-Type: application/json")
		
		// Check for JSON fields without assuming specific formatting
		assert.Contains(t, emailContent, "test-alert-id")
		assert.Contains(t, emailContent, "Test Alert")
		assert.Contains(t, emailContent, "This is a test alert for email notifications")
		assert.Contains(t, emailContent, "critical")
	})
	
	// Test email sending failure
	t.Run("SendFailure", func(t *testing.T) {
		// Create a mock SendMail function that always fails
		originalSendMail := sendMail
		defer func() { sendMail = originalSendMail }()
		
		sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return errors.New("server timeout")
		}
		
		// Create email config
		config := EmailConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "alerts@example.com",
		}
		
		// Create notifier
		notifier, err := NewEmailNotifier(config)
		require.NoError(t, err)
		
		// Send notification
		recipient := "user@example.com"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify error occurred
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to send email")
		
		// Verify notification status
		assert.Equal(t, "failed", notification.Status)
		assert.Contains(t, notification.ErrorMessage, "Failed to send email")
	})
	
	// Test email template rendering
	t.Run("TemplateRendering", func(t *testing.T) {
		// Create a mock SendMail function
		originalSendMail := sendMail
		defer func() { sendMail = originalSendMail }()
		
		var capturedMsg []byte
		
		sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			capturedMsg = msg
			return nil
		}
		
		// Create email config with custom template
		config := EmailConfig{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "test@example.com",
			Password: "password123",
			From:     "alerts@example.com",
			UseHTML:  true,
			Template: `
				<html>
				<body>
					<h1>Custom Alert: {{.Alert.Name}}</h1>
					<p>{{.Alert.Description}}</p>
					<p>Value: {{.Alert.Value}}, Threshold: {{.Alert.Threshold}}</p>
					<p>Status: {{.Alert.Status}}</p>
				</body>
				</html>
			`,
		}
		
		// Create notifier with mocked sendMail function
		notifier, err := NewEmailNotifier(config)
		require.NoError(t, err)
		
		// Send notification
		recipient := "user@example.com"
		_, err = notifier.Send(alert, recipient)
		
		// Verify no error occurred
		require.NoError(t, err)
		
		// Verify custom template was used
		emailContent := string(capturedMsg)
		assert.Contains(t, emailContent, "Custom Alert: Test Alert")
		assert.Contains(t, emailContent, "Value: 95.5, Threshold: 90")
		assert.Contains(t, emailContent, "Status: active")
	})
}
