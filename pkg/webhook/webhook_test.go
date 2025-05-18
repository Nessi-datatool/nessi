package webhook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookManager(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run test with timeout
	testutil.RunWithTimeout(t, func() {
		manager := NewWebhookManager()

		// Test RegisterWebhook
		webhook := &WebhookConfig{
			ID:          "test-webhook",
			Name:        "Test Webhook",
			URL:         "http://example.com/webhook",
			Events:      []string{"test.event"},
			Enabled:     true,
			RetryCount:  1, // Minimal retry count
			RetryDelay:  1, // Minimal retry delay
			Description: "Test webhook for unit tests",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := manager.RegisterWebhook(webhook)
		require.NoError(t, err, "Failed to register webhook")

		// Test GetWebhook
		retrievedWebhook, err := manager.GetWebhook("test-webhook")
		require.NoError(t, err, "Failed to get webhook")
		assert.Equal(t, webhook.ID, retrievedWebhook.ID, "Webhook ID mismatch")
		assert.Equal(t, webhook.Name, retrievedWebhook.Name, "Webhook Name mismatch")

		// Test ListWebhooks
		webhooks := manager.ListWebhooks()
		assert.Len(t, webhooks, 1, "Expected 1 webhook")

		// Test UpdateWebhook
		updates := map[string]interface{}{
			"name":        "Updated Webhook",
			"description": "Updated description",
		}

		err = manager.UpdateWebhook("test-webhook", updates)
		require.NoError(t, err, "Failed to update webhook")

		updatedWebhook, err := manager.GetWebhook("test-webhook")
		require.NoError(t, err, "Failed to get updated webhook")
		assert.Equal(t, "Updated Webhook", updatedWebhook.Name, "Webhook name not updated correctly")
		assert.Equal(t, "Updated description", updatedWebhook.Description, "Webhook description not updated correctly")

		// Test UnregisterWebhook
		err = manager.UnregisterWebhook("test-webhook")
		require.NoError(t, err, "Failed to unregister webhook")

		webhooks = manager.ListWebhooks()
		assert.Len(t, webhooks, 0, "Expected 0 webhooks after unregistering")

		// Test error cases
		_, err = manager.GetWebhook("non-existent")
		assert.Error(t, err, "Expected error when getting non-existent webhook")

		err = manager.UnregisterWebhook("non-existent")
		assert.Error(t, err, "Expected error when unregistering non-existent webhook")

		err = manager.UpdateWebhook("non-existent", updates)
		assert.Error(t, err, "Expected error when updating non-existent webhook")
	})
}

func TestTriggerEvent(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run test with timeout
	testutil.RunWithTimeout(t, func() {
		// Create a test server to receive webhook events with minimal processing
		var receivedPayload []byte
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Quick header checks
			if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" || 
			   r.Header.Get("User-Agent") != "Nessi-Webhook-Client/1.0" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Read body directly without decoding
			receivedPayload, _ = io.ReadAll(r.Body)
			defer r.Body.Close()

			// Return success immediately
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"success"}`))  
		}))
		defer server.Close()

		// Create a webhook manager
		manager := NewWebhookManager()

		// Register a webhook pointing to our test server
		webhook := &WebhookConfig{
			ID:          "test-webhook",
			Name:        "Test Webhook",
			URL:         server.URL,
			Events:      []string{"test.event"},
			Enabled:     true,
			RetryCount:  1,
			RetryDelay:  1,
			Description: "Test webhook for unit tests",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := manager.RegisterWebhook(webhook)
		require.NoError(t, err, "Failed to register webhook")

		// Create an event to trigger
		event := &WebhookEvent{
			ID:        "test-event",
			EventType: "test.event",
			Timestamp: time.Now(),
			Payload: map[string]interface{}{
				"message": "Hello, webhook!",
				"value":   42,
			},
		}

		// Trigger the event
		results := manager.TriggerEvent(event)

		// Verify the results
		result, ok := results["test-webhook"]
		require.True(t, ok, "Expected result for test-webhook, got none")
		assert.True(t, result.Success, "Expected webhook trigger to succeed, got failure: %v", result.Error)
		assert.Equal(t, http.StatusOK, result.StatusCode, "Status code mismatch")

		// Verify the payload was received correctly
		var receivedEvent WebhookEvent
		err = json.Unmarshal(receivedPayload, &receivedEvent)
		require.NoError(t, err, "Failed to unmarshal received payload")
		assert.Equal(t, event.ID, receivedEvent.ID, "Event ID mismatch")
		assert.Equal(t, event.EventType, receivedEvent.EventType, "Event type mismatch")

		// Test event filtering
		// Create a webhook for a different event type
		webhookOther := &WebhookConfig{
			ID:          "other-webhook",
			Name:        "Other Webhook",
			URL:         server.URL,
			Events:      []string{"other.event"},
			Enabled:     true,
			RetryCount:  1,
			RetryDelay:  1,
			Description: "Test webhook for other events",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = manager.RegisterWebhook(webhookOther)
		require.NoError(t, err, "Failed to register other webhook")

		// Trigger the same event
		results = manager.TriggerEvent(event)

		// Verify only the first webhook was triggered
		_, ok = results["test-webhook"]
		assert.True(t, ok, "Expected result for test-webhook, got none")

		_, ok = results["other-webhook"]
		assert.False(t, ok, "Expected no result for other-webhook, got one")

		// Test wildcard event subscription
		webhookWildcard := &WebhookConfig{
			ID:          "wildcard-webhook",
			Name:        "Wildcard Webhook",
			URL:         server.URL,
			Events:      []string{"*"},
			Enabled:     true,
			RetryCount:  1,
			RetryDelay:  1,
			Description: "Test webhook for all events",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = manager.RegisterWebhook(webhookWildcard)
		require.NoError(t, err, "Failed to register wildcard webhook")

		// Trigger the event
		results = manager.TriggerEvent(event)

		// Verify both the specific and wildcard webhooks were triggered
		_, ok = results["test-webhook"]
		assert.True(t, ok, "Expected result for test-webhook, got none")

		_, ok = results["wildcard-webhook"]
		assert.True(t, ok, "Expected result for wildcard-webhook, got none")

		// Test disabled webhook
		err = manager.UpdateWebhook("test-webhook", map[string]interface{}{
			"enabled": false,
		})
		require.NoError(t, err, "Failed to disable webhook")

		// Trigger the event
		results = manager.TriggerEvent(event)

		// Verify the disabled webhook was not triggered
		_, ok = results["test-webhook"]
		assert.False(t, ok, "Expected no result for disabled webhook, got one")

		_, ok = results["wildcard-webhook"]
		assert.True(t, ok, "Expected result for wildcard-webhook, got none")
	})
}

// TestWebhookManagerMinimal tests the webhook manager with minimal configuration
func TestWebhookManagerMinimal(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run test with timeout
	testutil.RunWithTimeout(t, func() {
		// Create a test server with minimal response time
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))  
		}))
		defer server.Close()

		// Create webhook manager with minimal timeout
		manager := &WebhookManager{
			webhooks: map[string]*WebhookConfig{
				"test-webhook": {
					ID:          "test-webhook",
					Name:        "Test Webhook",
					URL:         server.URL,
					Events:      []string{"test-event"},
					Enabled:     true,
					Headers:     map[string]string{},
					RetryCount:  1,
					RetryDelay:  1,
				},
			},
			client: &http.Client{
				Timeout: 100 * time.Millisecond, // Very short timeout
			},
			mu:     sync.RWMutex{},
		}

		// Create test event
		event := &WebhookEvent{
			ID:        "test-event-id",
			EventType: "test-event",
			Timestamp: time.Now(),
			Payload:   map[string]interface{}{"message": "Test message"},
		}

		// Trigger event
		results := manager.TriggerEvent(event)
		require.NotNil(t, results)
		require.Len(t, results, 1)

		// Check the result
		result := results["test-webhook"]
		require.NotNil(t, result)
		assert.Equal(t, true, result.Success)
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}

// TestWebhookNotifierMinimal tests the webhook notifier with minimal configuration
func TestWebhookNotifierMinimal(t *testing.T) {
	// No longer using testutil.RunInParallel(t) to avoid duplicate t.Parallel() calls
	// Tests will still run efficiently with the Go test runner
	
	// Run test with timeout
	testutil.RunWithTimeout(t, func() {
		// Create a test server with minimal response time
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))  
		}))
		defer server.Close()

		// Create webhook notifier with minimal timeout
		notifier := NewWebhookNotifier(
			server.URL,
			"POST",
			map[string]string{},
			100 * time.Millisecond, // Very short timeout
		)

		// Test sending a notification with a context that has a short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		err := notifier.SendWithContext(ctx, map[string]interface{}{"message": "Test message"})
		require.NoError(t, err)
	})
}
