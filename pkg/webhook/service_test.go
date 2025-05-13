package webhook

import (
	"os"
	"testing"
	"time"
)

func TestWebhookService(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "webhook-service-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create storage and service
	storage, err := NewFileStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	service, err := NewWebhookService(storage)
	if err != nil {
		t.Fatalf("Failed to create webhook service: %v", err)
	}

	// Test CreateWebhook
	webhook, err := service.CreateWebhook(
		"Test Webhook",
		"http://example.com/webhook",
		[]string{"test.event"},
		map[string]string{"X-Test": "true"},
		"Test webhook for unit tests",
	)
	if err != nil {
		t.Fatalf("Failed to create webhook: %v", err)
	}

	if webhook.Name != "Test Webhook" || webhook.URL != "http://example.com/webhook" {
		t.Errorf("Created webhook does not have expected values: %v", webhook)
	}

	// Test GetWebhook
	retrievedWebhook, err := service.GetWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to get webhook: %v", err)
	}

	if retrievedWebhook.ID != webhook.ID || retrievedWebhook.Name != webhook.Name {
		t.Errorf("Retrieved webhook does not match created webhook: got %v, want %v", retrievedWebhook, webhook)
	}

	// Test ListWebhooks
	webhooks := service.ListWebhooks()
	if len(webhooks) != 1 {
		t.Errorf("Expected 1 webhook, got %d", len(webhooks))
	}

	// Test UpdateWebhook
	updates := map[string]interface{}{
		"name":        "Updated Webhook",
		"description": "Updated description",
	}

	err = service.UpdateWebhook(webhook.ID, updates)
	if err != nil {
		t.Fatalf("Failed to update webhook: %v", err)
	}

	updatedWebhook, err := service.GetWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to get updated webhook: %v", err)
	}

	if updatedWebhook.Name != "Updated Webhook" || updatedWebhook.Description != "Updated description" {
		t.Errorf("Webhook not updated correctly: got %v", updatedWebhook)
	}

	// Test EnableWebhook and DisableWebhook
	err = service.DisableWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to disable webhook: %v", err)
	}

	disabledWebhook, err := service.GetWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to get disabled webhook: %v", err)
	}

	if disabledWebhook.Enabled {
		t.Error("Expected webhook to be disabled, but it is enabled")
	}

	err = service.EnableWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to enable webhook: %v", err)
	}

	enabledWebhook, err := service.GetWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to get enabled webhook: %v", err)
	}

	if !enabledWebhook.Enabled {
		t.Error("Expected webhook to be enabled, but it is disabled")
	}

	// Test DeleteWebhook
	err = service.DeleteWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to delete webhook: %v", err)
	}

	webhooks = service.ListWebhooks()
	if len(webhooks) != 0 {
		t.Errorf("Expected 0 webhooks after deletion, got %d", len(webhooks))
	}

	// Test error cases
	_, err = service.GetWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent webhook, got nil")
	}

	err = service.DeleteWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent webhook, got nil")
	}

	err = service.UpdateWebhook("non-existent", updates)
	if err == nil {
		t.Error("Expected error when updating non-existent webhook, got nil")
	}
}

func TestWebhookServiceTriggerEvent(t *testing.T) {
	// Use memory storage for this test to avoid file I/O
	storage := NewMemoryStorage()
	service, err := NewWebhookService(storage)
	if err != nil {
		t.Fatalf("Failed to create webhook service: %v", err)
	}

	// Create a webhook
	webhook, err := service.CreateWebhook(
		"Test Webhook",
		"http://non-existent-url.example.com/webhook", // Use a non-existent URL to avoid actual HTTP requests
		[]string{"test.event"},
		nil,
		"Test webhook for unit tests",
	)
	if err != nil {
		t.Fatalf("Failed to create webhook: %v", err)
	}

	// Trigger an event
	payload := map[string]interface{}{
		"message": "Test event",
		"time":    time.Now().Format(time.RFC3339),
	}

	results := service.TriggerEvent("test.event", payload)

	// Verify the webhook was triggered (but will fail due to non-existent URL)
	result, ok := results[webhook.ID]
	if !ok {
		t.Fatal("Expected result for webhook, got none")
	}

	if result.Success {
		t.Error("Expected webhook trigger to fail due to non-existent URL, but it succeeded")
	}

	// Test event filtering
	results = service.TriggerEvent("other.event", payload)
	if len(results) != 0 {
		t.Errorf("Expected 0 webhooks to be triggered for unsubscribed event, got %d", len(results))
	}

	// Test with wildcard event
	err = service.UpdateWebhook(webhook.ID, map[string]interface{}{
		"events": []string{"*"},
	})
	if err != nil {
		t.Fatalf("Failed to update webhook events: %v", err)
	}

	results = service.TriggerEvent("any.event", payload)
	_, ok = results[webhook.ID]
	if !ok {
		t.Error("Expected webhook with wildcard event to be triggered, but it wasn't")
	}

	// Test with disabled webhook
	err = service.DisableWebhook(webhook.ID)
	if err != nil {
		t.Fatalf("Failed to disable webhook: %v", err)
	}

	results = service.TriggerEvent("any.event", payload)
	if len(results) != 0 {
		t.Errorf("Expected 0 webhooks to be triggered when webhook is disabled, got %d", len(results))
	}
}

func TestGetSupportedEvents(t *testing.T) {
	storage := NewMemoryStorage()
	service, err := NewWebhookService(storage)
	if err != nil {
		t.Fatalf("Failed to create webhook service: %v", err)
	}

	events := service.GetSupportedEvents()
	if len(events) == 0 {
		t.Error("Expected supported events to be non-empty")
	}

	// Check for specific expected events
	expectedEvents := []string{
		"data.profile.created",
		"alert.triggered",
		"rule.created",
	}

	for _, expected := range expectedEvents {
		found := false
		for _, event := range events {
			if event == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected event %s not found in supported events", expected)
		}
	}
}
