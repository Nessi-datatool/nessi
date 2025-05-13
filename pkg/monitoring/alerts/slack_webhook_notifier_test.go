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

// TestSlackNotifierIntegration performs more comprehensive integration tests
func TestSlackNotifierIntegration(t *testing.T) {
	// Create a test server that simulates Slack API
	var receivedRequests [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedRequests = append(receivedRequests, body)
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

	// Create test alerts with different severities
	alerts := []*Alert{
		{
			ID:                "test-info-alert",
			Name:              "Info Alert",
			Description:       "This is an info alert",
			Type:              TypeSystem,
			Severity:          SeverityInfo,
			Status:            StatusActive,
			Source:            "test-source",
			Timestamp:         time.Now(),
			Value:             75.0,
			Threshold:         80.0,
			ComparisonOperator: "<",
		},
		{
			ID:                "test-warning-alert",
			Name:              "Warning Alert",
			Description:       "This is a warning alert",
			Type:              TypeQuality,
			Severity:          SeverityWarning,
			Status:            StatusActive,
			Source:            "test-source",
			Timestamp:         time.Now(),
			Value:             85.0,
			Threshold:         80.0,
			ComparisonOperator: ">",
		},
		{
			ID:                "test-critical-alert",
			Name:              "Critical Alert",
			Description:       "This is a critical alert",
			Type:              TypeAnomaly,
			Severity:          SeverityCritical,
			Status:            StatusActive,
			Source:            "test-source",
			Timestamp:         time.Now(),
			Value:             95.0,
			Threshold:         90.0,
			ComparisonOperator: ">",
		},
	}

	// Send notifications for each alert
	for _, alert := range alerts {
		notification, err := notifier.Send(alert, "#test-channel")
		require.NoError(t, err)
		assert.Equal(t, "sent", notification.Status)
	}

	// Verify we received 3 requests
	require.Len(t, receivedRequests, 3)

	// Verify each request
	for i, requestBody := range receivedRequests {
		var message SlackMessage
		err := json.Unmarshal(requestBody, &message)
		require.NoError(t, err)

		// Verify common fields
		assert.Equal(t, "#test-channel", message.Channel)
		assert.Equal(t, "Nessi Alert Bot", message.Username)
		assert.Equal(t, ":warning:", message.IconEmoji)
		assert.Contains(t, message.Text, "Alert: "+alerts[i].Name)
		
		// Verify attachments
		require.Len(t, message.Attachments, 1)
		attachment := message.Attachments[0]
		assert.Equal(t, alerts[i].Name, attachment.Title)
		assert.Equal(t, alerts[i].Description, attachment.Text)
		
		// Verify color based on severity
		switch alerts[i].Severity {
		case SeverityInfo:
			assert.Equal(t, "#2196f3", attachment.Color) // Blue
		case SeverityWarning:
			assert.Equal(t, "#ff9800", attachment.Color) // Orange
		case SeverityCritical:
			assert.Equal(t, "#f44336", attachment.Color) // Red
		}
		
		// Verify fields
		var foundSeverity, foundType, foundValue, foundThreshold bool
		for _, field := range attachment.Fields {
			switch field.Title {
			case "Severity":
				assert.Equal(t, string(alerts[i].Severity), field.Value)
				foundSeverity = true
			case "Type":
				assert.Equal(t, string(alerts[i].Type), field.Value)
				foundType = true
			case "Value":
				foundValue = true
			case "Threshold":
				foundThreshold = true
			}
		}
		
		assert.True(t, foundSeverity, "Severity field not found")
		assert.True(t, foundType, "Type field not found")
		assert.True(t, foundValue, "Value field not found")
		assert.True(t, foundThreshold, "Threshold field not found")
	}
}

// TestWebhookNotifierIntegration performs more comprehensive integration tests
func TestWebhookNotifierIntegration(t *testing.T) {
	// Create a test server with different response codes
	responseIndex := 0
	responses := []struct {
		statusCode int
		body       string
	}{
		{http.StatusOK, `{"status":"success"}`},
		{http.StatusInternalServerError, `{"status":"error"}`},
		{http.StatusOK, `{"status":"success"}`},
	}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the current response
		response := responses[responseIndex]
		responseIndex = (responseIndex + 1) % len(responses)
		
		// Set the status code and write the response
		w.WriteHeader(response.statusCode)
		w.Write([]byte(response.body))
	}))
	defer server.Close()
	
	// Create webhook config with retry settings
	config := WebhookConfig{
		URL:           server.URL,
		Method:        "POST",
		Headers: map[string]string{
			"X-API-Key": "test-api-key",
			"X-Source":  "test-source",
		},
		Timeout:       1 * time.Second,
		MaxRetries:    2,
		RetryInterval: 100 * time.Millisecond,
	}
	
	// Create notifier
	notifier := NewWebhookNotifier(config)
	
	// Create test alert
	alert := &Alert{
		ID:                "test-webhook-alert",
		Name:              "Webhook Test Alert",
		Description:       "This is a test alert for webhook",
		Type:              TypeQuality,
		Severity:          SeverityWarning,
		Status:            StatusActive,
		Source:            "test-webhook",
		Timestamp:         time.Now(),
		Value:             85.0,
		Threshold:         80.0,
		ComparisonOperator: ">",
		Labels: map[string]string{
			"environment": "test",
			"component":   "webhook-test",
		},
	}
	
	// Test 1: First request should succeed
	notification1, err := notifier.Send(alert, "webhook-id-1")
	require.NoError(t, err)
	assert.Equal(t, "sent", notification1.Status)
	
	// Test 2: Second request should fail but retry and still fail
	notification2, err := notifier.Send(alert, "webhook-id-2")
	require.Error(t, err)
	assert.Equal(t, "failed", notification2.Status)
	assert.Contains(t, notification2.ErrorMessage, "Failed to send webhook")
	
	// Test 3: Third request should succeed
	notification3, err := notifier.Send(alert, "webhook-id-3")
	require.NoError(t, err)
	assert.Equal(t, "sent", notification3.Status)
}

// TestNotifierWithDifferentHTTPMethods tests the webhook notifier with different HTTP methods
func TestNotifierWithDifferentHTTPMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var receivedMethod string
			
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"success"}`))
			}))
			defer server.Close()
			
			// Create webhook config
			config := WebhookConfig{
				URL:    server.URL,
				Method: method,
			}
			
			// Create notifier
			notifier := NewWebhookNotifier(config)
			
			// Create test alert
			alert := &Alert{
				ID:          "test-method-alert",
				Name:        "Method Test Alert",
				Description: "Testing different HTTP methods",
				Type:        TypeSystem,
				Severity:    SeverityInfo,
				Status:      StatusActive,
				Source:      "test-method",
				Timestamp:   time.Now(),
			}
			
			// Send notification
			notification, err := notifier.Send(alert, "webhook-method-test")
			require.NoError(t, err)
			assert.Equal(t, "sent", notification.Status)
			
			// Verify method
			assert.Equal(t, method, receivedMethod)
		})
	}
}

// TestWebhookPayloadContent tests the content of the webhook payload
func TestWebhookPayloadContent(t *testing.T) {
	var receivedPayload []byte
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		receivedPayload, err = io.ReadAll(r.Body)
		require.NoError(t, err)
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()
	
	// Create webhook config
	config := WebhookConfig{
		URL:    server.URL,
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	
	// Create notifier
	notifier := NewWebhookNotifier(config)
	
	// Create test alert with comprehensive data
	now := time.Now()
	alert := &Alert{
		ID:                "test-payload-alert",
		Name:              "Payload Test Alert",
		Description:       "Testing webhook payload content",
		Type:              TypeQuality,
		Severity:          SeverityWarning,
		Status:            StatusActive,
		Source:            "test-payload",
		Timestamp:         now,
		LastUpdated:       now,
		Value:             85.5,
		Threshold:         80.0,
		ComparisonOperator: ">",
		Labels: map[string]string{
			"environment": "test",
			"component":   "payload-test",
			"version":     "1.0.0",
		},
		Annotations: map[string]string{
			"summary":     "Test alert for webhook payload",
			"description": "This is a detailed description",
			"impact":      "Medium",
		},
	}
	
	// Send notification
	webhookID := "webhook-payload-test"
	notification, err := notifier.Send(alert, webhookID)
	require.NoError(t, err)
	assert.Equal(t, "sent", notification.Status)
	
	// Parse payload
	var payload map[string]interface{}
	err = json.Unmarshal(receivedPayload, &payload)
	require.NoError(t, err)
	
	// Verify alert data
	assert.Equal(t, "test-payload-alert", payload["alert_id"])
	assert.Equal(t, "Payload Test Alert", payload["name"])
	assert.Equal(t, "Testing webhook payload content", payload["description"])
	assert.Equal(t, "quality", payload["type"])
	assert.Equal(t, "warning", payload["severity"])
	assert.Equal(t, "active", payload["status"])
	assert.Equal(t, "test-payload", payload["source"])
	assert.Equal(t, 85.5, payload["value"])
	assert.Equal(t, 80.0, payload["threshold"])
	assert.Equal(t, ">", payload["comparison_operator"])
	assert.Equal(t, webhookID, payload["webhook_id"])
	assert.Equal(t, notification.ID, payload["notification_id"])
	
	// Verify labels
	labels, ok := payload["labels"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "test", labels["environment"])
	assert.Equal(t, "payload-test", labels["component"])
	assert.Equal(t, "1.0.0", labels["version"])
	
	// Verify annotations
	annotations, ok := payload["annotations"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "Test alert for webhook payload", annotations["summary"])
	assert.Equal(t, "This is a detailed description", annotations["description"])
	assert.Equal(t, "Medium", annotations["impact"])
}
