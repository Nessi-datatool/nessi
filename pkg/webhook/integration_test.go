package webhook

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookIntegration(t *testing.T) {
	// Create a temporary directory for webhook storage
	tempDir, err := os.MkdirTemp("", "webhook-integration-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create storage and service
	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)
	service, err := NewWebhookService(storage)
	require.NoError(t, err)

	// Test creating a webhook
	webhook, err := service.CreateWebhook(
		"Test Webhook",
		"http://example.com/webhook",
		[]string{"data.profile.created", "alert.triggered"},
		map[string]string{"X-Test": "true"},
		"Test webhook for integration tests",
	)
	require.NoError(t, err)
	require.NotNil(t, webhook)
	assert.Equal(t, "Test Webhook", webhook.Name)
	assert.Equal(t, "http://example.com/webhook", webhook.URL)
	assert.Contains(t, webhook.Events, "data.profile.created")
	assert.Contains(t, webhook.Events, "alert.triggered")
	assert.Equal(t, "true", webhook.Headers["X-Test"])
	assert.Equal(t, "Test webhook for integration tests", webhook.Description)
	assert.True(t, webhook.Enabled)

	// Test listing webhooks
	webhooks := service.ListWebhooks()
	assert.Len(t, webhooks, 1)
	assert.Equal(t, webhook.ID, webhooks[0].ID)

	// Test getting a webhook
	retrievedWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedWebhook)
	assert.Equal(t, webhook.ID, retrievedWebhook.ID)
	assert.Equal(t, webhook.Name, retrievedWebhook.Name)

	// Test updating a webhook
	err = service.UpdateWebhook(webhook.ID, map[string]interface{}{
		"name":        "Updated Webhook",
		"description": "Updated description",
	})
	require.NoError(t, err)

	// Get the updated webhook
	updatedWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedWebhook)
	assert.Equal(t, "Updated Webhook", updatedWebhook.Name)
	assert.Equal(t, "Updated description", updatedWebhook.Description)

	// Test disabling a webhook
	err = service.DisableWebhook(webhook.ID)
	require.NoError(t, err)

	// Get the disabled webhook
	disabledWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, disabledWebhook)
	assert.False(t, disabledWebhook.Enabled)

	// Test enabling a webhook
	err = service.EnableWebhook(webhook.ID)
	require.NoError(t, err)

	// Get the enabled webhook
	enabledWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, enabledWebhook)
	assert.True(t, enabledWebhook.Enabled)

	// Test triggering an event
	payload := map[string]interface{}{
		"test":       true,
		"message":    "This is a test event",
		"webhook_id": webhook.ID,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	// Trigger event (this will not actually send HTTP requests in tests)
	results := service.TriggerEvent("data.profile.created", payload)
	result, ok := results[webhook.ID]
	assert.True(t, ok, "Webhook should be triggered")
	// In tests, the webhook delivery will fail because the URL is not reachable
	assert.False(t, result.Success, "Webhook delivery should fail in tests")

	// Test deleting a webhook
	err = service.DeleteWebhook(webhook.ID)
	require.NoError(t, err)

	// List webhooks after deletion
	webhooks = service.ListWebhooks()
	assert.Len(t, webhooks, 0)

	// Test error cases
	_, err = service.GetWebhook("non-existent")
	assert.Error(t, err)

	err = service.DeleteWebhook("non-existent")
	assert.Error(t, err)

	err = service.UpdateWebhook("non-existent", map[string]interface{}{
		"name": "Updated Webhook",
	})
	assert.Error(t, err)
}

func TestWebhookConfigSerialization(t *testing.T) {
	// Create a webhook config
	config := &WebhookConfig{
		ID:          "test-webhook",
		Name:        "Test Webhook",
		URL:         "http://example.com/webhook",
		Events:      []string{"data.profile.created", "alert.triggered"},
		Enabled:     true,
		RetryCount:  3,
		RetryDelay:  5,
		Description: "Test webhook for serialization tests",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(config)
	require.NoError(t, err)

	// Deserialize from JSON
	var deserializedConfig WebhookConfig
	err = json.Unmarshal(jsonData, &deserializedConfig)
	require.NoError(t, err)

	// Verify deserialized config
	assert.Equal(t, config.ID, deserializedConfig.ID)
	assert.Equal(t, config.Name, deserializedConfig.Name)
	assert.Equal(t, config.URL, deserializedConfig.URL)
	assert.Equal(t, config.Events, deserializedConfig.Events)
	assert.Equal(t, config.Enabled, deserializedConfig.Enabled)
	assert.Equal(t, config.RetryCount, deserializedConfig.RetryCount)
	assert.Equal(t, config.RetryDelay, deserializedConfig.RetryDelay)
	assert.Equal(t, config.Description, deserializedConfig.Description)
}

func TestSupportedEvents(t *testing.T) {
	// Create a temporary directory for webhook storage
	tempDir, err := os.MkdirTemp("", "webhook-events-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create storage and service
	storage, err := NewFileStorage(tempDir)
	require.NoError(t, err)
	service, err := NewWebhookService(storage)
	require.NoError(t, err)

	// Get supported events
	events := service.GetSupportedEvents()
	assert.NotEmpty(t, events)

	// Check for specific expected events
	expectedEvents := []string{
		"data.profile.created",
		"alert.triggered",
		"rule.created",
	}

	for _, expected := range expectedEvents {
		assert.Contains(t, events, expected)
	}
}
