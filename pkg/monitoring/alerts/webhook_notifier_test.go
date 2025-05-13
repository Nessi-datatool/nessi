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

// TestWebhookNotifier tests the WebhookNotifier
func TestWebhookNotifier(t *testing.T) {
	// Create a test alert
	alert := &Alert{
		ID:                "test-alert-id",
		Name:              "Test Alert",
		Description:       "This is a test alert for webhook notifications",
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

	// Test successful webhook notification
	t.Run("SendSuccess", func(t *testing.T) {
		// Create a test server that simulates webhook endpoint
		var receivedMethod string
		var receivedHeaders http.Header
		var receivedPayload []byte
		
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Capture request details
			receivedMethod = r.Method
			receivedHeaders = r.Header
			var err error
			receivedPayload, err = io.ReadAll(r.Body)
			require.NoError(t, err)
			
			// Return success response
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		}))
		defer server.Close()
		
		// Create webhook config
		config := WebhookConfig{
			URL:    server.URL,
			Method: "POST",
			Headers: map[string]string{
				"X-API-Key":    "test-api-key",
				"X-Source":     "nessi-alerts",
				"Content-Type": "application/json",
			},
			Timeout:       5 * time.Second,
			MaxRetries:    3,
			RetryInterval: 1 * time.Second,
		}
		
		// Create notifier
		notifier := NewWebhookNotifier(config)
		
		// Send notification
		recipient := "webhook-id-123"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify no error occurred
		require.NoError(t, err)
		
		// Verify notification details
		assert.Equal(t, alert.ID, notification.AlertID)
		assert.Equal(t, "webhook", notification.Channel)
		assert.Equal(t, recipient, notification.Recipient)
		assert.Equal(t, "sent", notification.Status)
		
		// Verify request method
		assert.Equal(t, "POST", receivedMethod)
		
		// Verify headers
		assert.Equal(t, "test-api-key", receivedHeaders.Get("X-API-Key"))
		assert.Equal(t, "nessi-alerts", receivedHeaders.Get("X-Source"))
		assert.Equal(t, "application/json", receivedHeaders.Get("Content-Type"))
		assert.Equal(t, "Nessi-Webhook-Notifier/1.0", receivedHeaders.Get("User-Agent"))
		
		// Parse received payload
		var payload map[string]interface{}
		err = json.Unmarshal(receivedPayload, &payload)
		require.NoError(t, err)
		
		// Verify payload details
		assert.Equal(t, "test-alert-id", payload["alert_id"])
		assert.Equal(t, "Test Alert", payload["name"])
		assert.Equal(t, "This is a test alert for webhook notifications", payload["description"])
		assert.Equal(t, "quality", payload["type"])
		assert.Equal(t, "critical", payload["severity"])
		assert.Equal(t, "active", payload["status"])
		assert.Equal(t, "test-source", payload["source"])
		assert.Equal(t, 95.5, payload["value"])
		assert.Equal(t, 90.0, payload["threshold"])
		assert.Equal(t, ">", payload["comparison_operator"])
		assert.Equal(t, recipient, payload["webhook_id"])
		assert.NotNil(t, payload["notification_id"])
		assert.NotNil(t, payload["notification_time"])
		
		// Verify labels
		labels, ok := payload["labels"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "test", labels["environment"])
		assert.Equal(t, "data-quality", labels["component"])
		
		// Verify annotations
		annotations, ok := payload["annotations"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "Data quality score exceeded threshold", annotations["summary"])
		assert.Equal(t, "High", annotations["impact"])
	})
	
	// Test webhook notification failure
	t.Run("SendFailure", func(t *testing.T) {
		// Create a test server that simulates webhook endpoint failure
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Return error response
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"status":"error","message":"Internal server error"}`))
		}))
		defer server.Close()
		
		// Create webhook config with minimal retries for faster test
		config := WebhookConfig{
			URL:           server.URL,
			Method:        "POST",
			Timeout:       1 * time.Second,
			MaxRetries:    1,
			RetryInterval: 100 * time.Millisecond,
		}
		
		// Create notifier
		notifier := NewWebhookNotifier(config)
		
		// Send notification
		recipient := "webhook-id-123"
		notification, err := notifier.Send(alert, recipient)
		
		// Verify error occurred
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to send webhook")
		
		// Verify notification status
		assert.Equal(t, "failed", notification.Status)
		assert.Contains(t, notification.ErrorMessage, "Failed to send webhook")
	})
	
	// Test with different HTTP methods
	t.Run("DifferentMethods", func(t *testing.T) {
		testMethods := []string{"GET", "POST", "PUT"}
		
		for _, method := range testMethods {
			// Create a test server
			var receivedMethod string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedMethod = r.Method
				w.WriteHeader(http.StatusOK)
			}))
			
			// Create webhook config
			config := WebhookConfig{
				URL:    server.URL,
				Method: method,
			}
			
			// Create notifier
			notifier := NewWebhookNotifier(config)
			
			// Send notification
			_, err := notifier.Send(alert, "webhook-id")
			require.NoError(t, err)
			
			// Verify request method
			assert.Equal(t, method, receivedMethod)
			
			server.Close()
		}
	})
	
	// Test default values
	t.Run("DefaultValues", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()
		
		// Create webhook config with minimal settings
		config := WebhookConfig{
			URL: server.URL,
			// Method not specified - should default to POST
			// Timeout not specified - should default to 10 seconds
		}
		
		// Create notifier
		notifier := NewWebhookNotifier(config)
		
		// Verify default values
		assert.Equal(t, "POST", notifier.config.Method)
		assert.Equal(t, 10*time.Second, notifier.config.Timeout)
	})
	
	// Test with invalid URL
	t.Run("InvalidURL", func(t *testing.T) {
		// Create webhook config with invalid URL
		config := WebhookConfig{
			URL:        "http://invalid-url-that-does-not-exist.example.com",
			MaxRetries: 1, // Minimal retries for faster test
		}
		
		// Create notifier
		notifier := NewWebhookNotifier(config)
		
		// Send notification
		notification, err := notifier.Send(alert, "webhook-id")
		
		// Verify error occurred
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to send webhook")
		
		// Verify notification status
		assert.Equal(t, "failed", notification.Status)
		assert.Contains(t, notification.ErrorMessage, "Failed to send webhook")
	})
}
