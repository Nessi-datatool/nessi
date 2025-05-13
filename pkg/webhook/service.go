package webhook

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// WebhookService provides a high-level interface for working with webhooks
type WebhookService struct {
	manager *WebhookManager
	storage WebhookStorage
	mu      sync.RWMutex
}

// NewWebhookService creates a new webhook service
func NewWebhookService(storage WebhookStorage) (*WebhookService, error) {
	service := &WebhookService{
		manager: NewWebhookManager(),
		storage: storage,
	}

	// Load webhooks from storage
	if err := service.loadWebhooks(); err != nil {
		return nil, fmt.Errorf("failed to load webhooks: %v", err)
	}

	return service, nil
}

// loadWebhooks loads all webhooks from storage into the manager
func (s *WebhookService) loadWebhooks() error {
	webhooks, err := s.storage.ListWebhooks()
	if err != nil {
		return err
	}

	for _, webhook := range webhooks {
		if err := s.manager.RegisterWebhook(webhook); err != nil {
			return err
		}
	}

	return nil
}

// CreateWebhook creates a new webhook
func (s *WebhookService) CreateWebhook(name, url string, events []string, headers map[string]string, description string) (*WebhookConfig, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create webhook config
	config := &WebhookConfig{
		ID:          uuid.New().String(),
		Name:        name,
		URL:         url,
		Events:      events,
		Headers:     headers,
		Enabled:     true,
		RetryCount:  3,
		RetryDelay:  5,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Register webhook
	if err := s.manager.RegisterWebhook(config); err != nil {
		return nil, err
	}

	// Save webhook to storage
	if err := s.storage.SaveWebhook(config); err != nil {
		// Rollback registration
		_ = s.manager.UnregisterWebhook(config.ID)
		return nil, err
	}

	return config, nil
}

// GetWebhook gets a webhook by ID
func (s *WebhookService) GetWebhook(id string) (*WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.manager.GetWebhook(id)
}

// ListWebhooks lists all webhooks
func (s *WebhookService) ListWebhooks() []*WebhookConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.manager.ListWebhooks()
}

// UpdateWebhook updates a webhook
func (s *WebhookService) UpdateWebhook(id string, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update webhook in manager
	if err := s.manager.UpdateWebhook(id, updates); err != nil {
		return err
	}

	// Get updated webhook
	webhook, err := s.manager.GetWebhook(id)
	if err != nil {
		return err
	}

	// Save updated webhook to storage
	if err := s.storage.SaveWebhook(webhook); err != nil {
		return err
	}

	return nil
}

// DeleteWebhook deletes a webhook
func (s *WebhookService) DeleteWebhook(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Unregister webhook from manager
	if err := s.manager.UnregisterWebhook(id); err != nil {
		return err
	}

	// Delete webhook from storage
	if err := s.storage.DeleteWebhook(id); err != nil {
		return err
	}

	return nil
}

// EnableWebhook enables a webhook
func (s *WebhookService) EnableWebhook(id string) error {
	return s.UpdateWebhook(id, map[string]interface{}{
		"enabled": true,
	})
}

// DisableWebhook disables a webhook
func (s *WebhookService) DisableWebhook(id string) error {
	return s.UpdateWebhook(id, map[string]interface{}{
		"enabled": false,
	})
}

// TriggerEvent triggers an event and sends it to all subscribed webhooks
func (s *WebhookService) TriggerEvent(eventType string, payload map[string]interface{}) map[string]*WebhookResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event := &WebhookEvent{
		ID:        uuid.New().String(),
		EventType: eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}

	return s.manager.TriggerEvent(event)
}

// GetSupportedEvents returns a list of supported event types
func (s *WebhookService) GetSupportedEvents() []string {
	return []string{
		"data.profile.created",
		"data.profile.updated",
		"data.validation.failed",
		"data.validation.passed",
		"alert.triggered",
		"alert.resolved",
		"rule.created",
		"rule.updated",
		"rule.deleted",
		"metric.threshold.exceeded",
		"system.startup",
		"system.shutdown",
	}
}
