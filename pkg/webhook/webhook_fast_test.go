package webhook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebhookManagerFast is a fast version of the webhook manager test
func TestWebhookManagerFast(t *testing.T) {
	// Make this test run in parallel with other tests
	t.Parallel()

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
		}
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
}

// TestWebhookNotifierFast is a fast version of the webhook notifier test
func TestWebhookNotifierFast(t *testing.T) {
	// Make this test run in parallel with other tests
	t.Parallel()

	// Create a test server with minimal response time
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer server.Close()

	// Create webhook notifier with minimal timeout
	notifier := &WebhookNotifier{
		url:         server.URL,
		method:      "POST",
		contentType: "application/json",
		headers:     map[string]string{},
		timeout:     100 * time.Millisecond, // Very short timeout
	}

	// Test sending a notification with a context that has a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := notifier.SendWithContext(ctx, map[string]interface{}{"message": "Test message"})
	require.NoError(t, err)
}
