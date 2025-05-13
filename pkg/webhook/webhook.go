package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// WebhookConfig represents the configuration for a webhook
type WebhookConfig struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Events      []string          `json:"events"`
	Headers     map[string]string `json:"headers,omitempty"`
	Enabled     bool              `json:"enabled"`
	RetryCount  int               `json:"retry_count,omitempty"`
	RetryDelay  int               `json:"retry_delay,omitempty"` // in seconds
	Description string            `json:"description,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// WebhookEvent represents an event to be sent to a webhook
type WebhookEvent struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

// WebhookResult represents the result of a webhook delivery
type WebhookResult struct {
	Success      bool      `json:"success"`
	StatusCode   int       `json:"status_code,omitempty"`
	ResponseBody string    `json:"response_body,omitempty"`
	Error        string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	AttemptCount int       `json:"attempt_count"`
}

// WebhookManager manages webhooks and their delivery
type WebhookManager struct {
	webhooks map[string]*WebhookConfig
	client   *http.Client
	mu       sync.RWMutex
}

// NewWebhookManager creates a new webhook manager
func NewWebhookManager() *WebhookManager {
	return &WebhookManager{
		webhooks: make(map[string]*WebhookConfig),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterWebhook registers a new webhook
func (m *WebhookManager) RegisterWebhook(config *WebhookConfig) error {
	if config.ID == "" {
		return fmt.Errorf("webhook ID cannot be empty")
	}

	if config.URL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	if len(config.Events) == 0 {
		return fmt.Errorf("webhook must subscribe to at least one event")
	}

	// Set default values if not provided
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	if config.RetryDelay == 0 {
		config.RetryDelay = 5
	}

	if config.CreatedAt.IsZero() {
		config.CreatedAt = time.Now()
	}

	config.UpdatedAt = time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.webhooks[config.ID] = config
	return nil
}

// UnregisterWebhook removes a webhook
func (m *WebhookManager) UnregisterWebhook(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.webhooks[id]; !exists {
		return fmt.Errorf("webhook with ID %s does not exist", id)
	}

	delete(m.webhooks, id)
	return nil
}

// GetWebhook returns a webhook by ID
func (m *WebhookManager) GetWebhook(id string) (*WebhookConfig, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	webhook, exists := m.webhooks[id]
	if !exists {
		return nil, fmt.Errorf("webhook with ID %s does not exist", id)
	}

	return webhook, nil
}

// ListWebhooks returns all registered webhooks
func (m *WebhookManager) ListWebhooks() []*WebhookConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	webhooks := make([]*WebhookConfig, 0, len(m.webhooks))
	for _, webhook := range m.webhooks {
		webhooks = append(webhooks, webhook)
	}

	return webhooks
}

// UpdateWebhook updates an existing webhook
func (m *WebhookManager) UpdateWebhook(id string, updates map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	webhook, exists := m.webhooks[id]
	if !exists {
		return fmt.Errorf("webhook with ID %s does not exist", id)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		webhook.Name = name
	}

	if url, ok := updates["url"].(string); ok {
		webhook.URL = url
	}

	if events, ok := updates["events"].([]string); ok {
		webhook.Events = events
	}

	if headers, ok := updates["headers"].(map[string]string); ok {
		webhook.Headers = headers
	}

	if enabled, ok := updates["enabled"].(bool); ok {
		webhook.Enabled = enabled
	}

	if retryCount, ok := updates["retry_count"].(int); ok {
		webhook.RetryCount = retryCount
	}

	if retryDelay, ok := updates["retry_delay"].(int); ok {
		webhook.RetryDelay = retryDelay
	}

	if description, ok := updates["description"].(string); ok {
		webhook.Description = description
	}

	webhook.UpdatedAt = time.Now()
	return nil
}

// TriggerEvent sends an event to all webhooks that are subscribed to it
func (m *WebhookManager) TriggerEvent(event *WebhookEvent) map[string]*WebhookResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make(map[string]*WebhookResult)

	for id, webhook := range m.webhooks {
		if !webhook.Enabled {
			continue
		}

		// Check if webhook is subscribed to this event
		subscribed := false
		for _, eventType := range webhook.Events {
			if eventType == event.EventType || eventType == "*" {
				subscribed = true
				break
			}
		}

		if !subscribed {
			continue
		}

		// Send event to webhook
		result := m.sendWebhook(webhook, event)
		results[id] = result
	}

	return results
}

// sendWebhook sends an event to a single webhook with retries
func (m *WebhookManager) sendWebhook(webhook *WebhookConfig, event *WebhookEvent) *WebhookResult {
	result := &WebhookResult{
		Timestamp: time.Now(),
	}

	// Prepare payload
	payload, err := json.Marshal(event)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to marshal event: %v", err)
		return result
	}

	// Create request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payload))
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create request: %v", err)
		return result
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Nessi-Webhook-Client/1.0")
	for key, value := range webhook.Headers {
		req.Header.Set(key, value)
	}

	// Send request with retries
	var resp *http.Response
	for attempt := 1; attempt <= webhook.RetryCount; attempt++ {
		result.AttemptCount = attempt

		resp, err = m.client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success
			result.Success = true
			result.StatusCode = resp.StatusCode

			// Read response body
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(resp.Body)
			result.ResponseBody = buf.String()
			resp.Body.Close()
			break
		}

		if err != nil {
			result.Error = fmt.Sprintf("request failed: %v", err)
		} else {
			// Read response body for error details
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(resp.Body)
			result.ResponseBody = buf.String()
			resp.Body.Close()
			result.StatusCode = resp.StatusCode
			result.Error = fmt.Sprintf("received non-success status code: %d", resp.StatusCode)
		}

		// Retry after delay if not the last attempt
		if attempt < webhook.RetryCount {
			time.Sleep(time.Duration(webhook.RetryDelay) * time.Second)
		}
	}

	return result
}
