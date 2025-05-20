package webhook

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	// Test SaveWebhook and LoadWebhook
	webhook := &WebhookConfig{
		ID:          "test-webhook",
		Name:        "Test Webhook",
		URL:         "http://example.com/webhook",
		Events:      []string{"test.event"},
		Enabled:     true,
		RetryCount:  3,
		RetryDelay:  5,
		Description: "Test webhook for unit tests",
		CreatedAt:   time.Now().Truncate(time.Second),
		UpdatedAt:   time.Now().Truncate(time.Second),
	}

	err := storage.SaveWebhook(webhook)
	if err != nil {
		t.Fatalf("Failed to save webhook: %v", err)
	}

	loadedWebhook, err := storage.LoadWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to load webhook: %v", err)
	}

	if loadedWebhook.ID != webhook.ID || loadedWebhook.Name != webhook.Name {
		t.Errorf("Loaded webhook does not match original: got %v, want %v", loadedWebhook, webhook)
	}

	// Test ListWebhooks
	webhooks, err := storage.ListWebhooks()
	if err != nil {
		t.Fatalf("Failed to list webhooks: %v", err)
	}

	if len(webhooks) != 1 {
		t.Errorf("Expected 1 webhook, got %d", len(webhooks))
	}

	// Test DeleteWebhook
	err = storage.DeleteWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to delete webhook: %v", err)
	}

	webhooks, err = storage.ListWebhooks()
	if err != nil {
		t.Fatalf("Failed to list webhooks after deletion: %v", err)
	}

	if len(webhooks) != 0 {
		t.Errorf("Expected 0 webhooks after deletion, got %d", len(webhooks))
	}

	// Test error cases
	_, err = storage.LoadWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when loading non-existent webhook, got nil")
	}

	err = storage.DeleteWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent webhook, got nil")
	}
}

func TestFileStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "webhook-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	storage, err := NewFileStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create file storage: %v", err)
	}

	// Test SaveWebhook and LoadWebhook
	webhook := &WebhookConfig{
		ID:          "test-webhook",
		Name:        "Test Webhook",
		URL:         "http://example.com/webhook",
		Events:      []string{"test.event"},
		Enabled:     true,
		RetryCount:  3,
		RetryDelay:  5,
		Description: "Test webhook for unit tests",
		CreatedAt:   time.Now().Truncate(time.Second),
		UpdatedAt:   time.Now().Truncate(time.Second),
	}

	err = storage.SaveWebhook(webhook)
	if err != nil {
		t.Fatalf("Failed to save webhook: %v", err)
	}

	// Verify file was created
	filePath := filepath.Join(tempDir, "test-webhook.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected webhook file to exist at %s", filePath)
	}

	loadedWebhook, err := storage.LoadWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to load webhook: %v", err)
	}

	if loadedWebhook.ID != webhook.ID || loadedWebhook.Name != webhook.Name {
		t.Errorf("Loaded webhook does not match original: got %v, want %v", loadedWebhook, webhook)
	}

	// Test ListWebhooks
	webhooks, err := storage.ListWebhooks()
	if err != nil {
		t.Fatalf("Failed to list webhooks: %v", err)
	}

	if len(webhooks) != 1 {
		t.Errorf("Expected 1 webhook, got %d", len(webhooks))
	}

	// Test DeleteWebhook
	err = storage.DeleteWebhook("test-webhook")
	if err != nil {
		t.Fatalf("Failed to delete webhook: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("Expected webhook file to be deleted at %s", filePath)
	}

	webhooks, err = storage.ListWebhooks()
	if err != nil {
		t.Fatalf("Failed to list webhooks after deletion: %v", err)
	}

	if len(webhooks) != 0 {
		t.Errorf("Expected 0 webhooks after deletion, got %d", len(webhooks))
	}

	// Test error cases
	_, err = storage.LoadWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when loading non-existent webhook, got nil")
	}

	err = storage.DeleteWebhook("non-existent")
	if err == nil {
		t.Error("Expected error when deleting non-existent webhook, got nil")
	}
}
