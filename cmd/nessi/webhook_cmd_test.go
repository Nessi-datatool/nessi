package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/nessi-dev/nessi-dev/cmd/nessi/cli"
	"github.com/nessi-dev/nessi-dev/pkg/webhook"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/spf13/viper"
)

func TestWebhookConfig(t *testing.T) {
	// Create a webhook config
	config := &webhook.WebhookConfig{
		ID:          "test-webhook",
		Name:        "Test Webhook",
		URL:         "http://example.com/webhook",
		Events:      []string{"data.profile.created", "alert.triggered"},
		Enabled:     true,
		RetryCount:  3,
		RetryDelay:  5,
		Description: "Test webhook for unit tests",
	}

	// Serialize to JSON
	jsonData, err := json.Marshal(config)
	require.NoError(t, err)

	// Deserialize from JSON
	var deserializedConfig webhook.WebhookConfig
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

func TestWebhookService(t *testing.T) {
	// Create a temporary directory for webhook storage
	tempDir, err := os.MkdirTemp("", "webhook-cmd-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up viper to use the temporary directory
	viper.Set("config_dir", tempDir)
	defer viper.Reset()

	// Get webhook service
	service, err := getWebhookService()
	require.NoError(t, err)

	// Create a webhook
	webhook, err := service.CreateWebhook(
		"Test Webhook",
		"http://example.com/webhook",
		[]string{"data.profile.created", "alert.triggered"},
		map[string]string{"X-Test": "true"},
		"Test webhook for CLI tests",
	)
	require.NoError(t, err)
	require.NotNil(t, webhook)
	assert.Equal(t, "Test Webhook", webhook.Name)
	assert.Equal(t, "http://example.com/webhook", webhook.URL)
	assert.Contains(t, webhook.Events, "data.profile.created")
	assert.Contains(t, webhook.Events, "alert.triggered")
	assert.Equal(t, "true", webhook.Headers["X-Test"])
	assert.Equal(t, "Test webhook for CLI tests", webhook.Description)
	assert.True(t, webhook.Enabled)

	// Get the webhook
	retrievedWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedWebhook)
	assert.Equal(t, webhook.ID, retrievedWebhook.ID)
	assert.Equal(t, webhook.Name, retrievedWebhook.Name)

	// List webhooks
	webhooks := service.ListWebhooks()
	assert.Len(t, webhooks, 1)
	assert.Equal(t, webhook.ID, webhooks[0].ID)

	// Update the webhook
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

	// Disable the webhook
	err = service.DisableWebhook(webhook.ID)
	require.NoError(t, err)

	// Get the disabled webhook
	disabledWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, disabledWebhook)
	assert.False(t, disabledWebhook.Enabled)

	// Enable the webhook
	err = service.EnableWebhook(webhook.ID)
	require.NoError(t, err)

	// Get the enabled webhook
	enabledWebhook, err := service.GetWebhook(webhook.ID)
	require.NoError(t, err)
	require.NotNil(t, enabledWebhook)
	assert.True(t, enabledWebhook.Enabled)

	// Delete the webhook
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

func TestWebhookEvents(t *testing.T) {
	// Create a temporary directory for webhook storage
	tempDir, err := os.MkdirTemp("", "webhook-events-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Set up viper to use the temporary directory
	viper.Set("config_dir", tempDir)
	defer viper.Reset()

	// Get webhook service
	service, err := getWebhookService()
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
