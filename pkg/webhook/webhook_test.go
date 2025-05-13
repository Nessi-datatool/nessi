package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookManager(t *testing.T) {
	manager := NewWebhookManager()

	// Test RegisterWebhook
	webhook := &WebhookConfig{
		ID:          "test-webhook",
		Name:        "Test Webhook",
		URL:         "http://example.com/webhook",
		Events:      []string{"test.event"},
		Enabled:     true,
		RetryCount:  3,
		RetryDelay:  5,
		Description: "Test webhook for unit tests",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := manager.RegisterWebhook(webhook)
	if err != nil {
		t.Fatalf("Failed to register webhook: %v", err)
	}

	// Test GetWebhook
	retrievedWebhook, err := manager.GetWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to get webhook: %v", err)
	}

	if retrievedWebhook.ID != webhook.ID || retrievedWebhook.Name != webhook.Name {
		t.Errorf("Retrieved webhook does not match original: got %v, want %v", retrievedWebhook, webhook)
	}

	// Test ListWebhooks
	webhooks := manager.ListWebhooks()
	if len(webhooks) != 1 {
		t.Errorf("Expected 1 webhook, got %d", len(webhooks))
	}

	// Test UpdateWebhook
	updates := map[string]interface{}{
		"name":        "Updated Webhook",
		"description": "Updated description",
	}

	err = manager.UpdateWebhook("test-webhook", updates)
	if err != nil {
		t.Fatalf("Failed to update webhook: %v", err)
	}

	updatedWebhook, err := manager.GetWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to get updated webhook: %v", err)
	}

	if updatedWebhook.Name != "Updated Webhook" || updatedWebhook.Description != "Updated description" {
		t.Errorf("Webhook not updated correctly: got %v", updatedWebhook)
	}

	// Test UnregisterWebhook
	err = manager.UnregisterWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to unregister webhook: %v", err)
	}

	webhooks = manager.ListWebhooks()
	if len(webhooks) != 0 {
		t.Errorf("Expected 0 webhooks after unregistering, got %d", len(webhooks))
	}

	// Test error cases
	_, err = manager.GetWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent webhook, got nil")
	}

	err = manager.UnregisterWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when unregistering non-existent webhook, got nil")
	}

	err = manager.UpdateWebhook("non-existent", updates)
	if err == nil {
		t.Error("Expected error when updating non-existent webhook, got nil")
	}
}

func TestTriggerEvent(t *testing.T) {
	// Create a test server to receive webhook events
	var receivedPayload []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		if r.Header.Get("User-Agent") != "Nessi-Webhook-Client/1.0" {
			t.Errorf("Expected User-Agent: Nessi-Webhook-Client/1.0, got %s", r.Header.Get("User-Agent"))
		}

		// Read the request body
		decoder := json.NewDecoder(r.Body)
		var event WebhookEvent
		err := decoder.Decode(&event)
		if err != nil {
			t.Errorf("Failed to decode webhook event: %v", err)
		}

		// Store the payload for later verification
		receivedPayload, _ = json.Marshal(event)

		// Return a success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
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
	if err != nil {
		t.Fatalf("Failed to register webhook: %v", err)
	}

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
	if !ok {
		t.Fatal("Expected result for test-webhook, got none")
	}

	if !result.Success {
		t.Errorf("Expected webhook trigger to succeed, got failure: %v", result.Error)
	}

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}

	// Verify the payload was received correctly
	var receivedEvent WebhookEvent
	err = json.Unmarshal(receivedPayload, &receivedEvent)
	if err != nil {
		t.Fatalf("Failed to unmarshal received payload: %v", err)
	}

	if receivedEvent.ID != event.ID || receivedEvent.EventType != event.EventType {
		t.Errorf("Received event does not match sent event: got %v, want %v", receivedEvent, event)
	}

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
	if err != nil {
		t.Fatalf("Failed to register other webhook: %v", err)
	}

	// Trigger the same event
	results = manager.TriggerEvent(event)

	// Verify only the first webhook was triggered
	_, ok = results["test-webhook"]
	if !ok {
		t.Error("Expected result for test-webhook, got none")
	}

	_, ok = results["other-webhook"]
	if ok {
		t.Error("Expected no result for other-webhook, got one")
	}

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
	if err != nil {
		t.Fatalf("Failed to register wildcard webhook: %v", err)
	}

	// Trigger the event
	results = manager.TriggerEvent(event)

	// Verify both the specific and wildcard webhooks were triggered
	_, ok = results["test-webhook"]
	if !ok {
		t.Error("Expected result for test-webhook, got none")
	}

	_, ok = results["wildcard-webhook"]
	if !ok {
		t.Error("Expected result for wildcard-webhook, got none")
	}

	// Test disabled webhook
	err = manager.UpdateWebhook("test-webhook", map[string]interface{}{
		"enabled": false,
	})
	if err != nil {
		t.Fatalf("Failed to disable webhook: %v", err)
	}

	// Trigger the event
	results = manager.TriggerEvent(event)

	// Verify the disabled webhook was not triggered
	_, ok = results["test-webhook"]
	if ok {
		t.Error("Expected no result for disabled webhook, got one")
	}

	_, ok = results["wildcard-webhook"]
	if !ok {
		t.Error("Expected result for wildcard-webhook, got none")
	}
}
