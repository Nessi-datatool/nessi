package webhook

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// WebhookStorage provides storage for webhook configurations
type WebhookStorage interface {
	// SaveWebhook saves a webhook configuration
	SaveWebhook(config *WebhookConfig) error

	// LoadWebhook loads a webhook configuration by ID
	LoadWebhook(id string) (*WebhookConfig, error)

	// DeleteWebhook deletes a webhook configuration
	DeleteWebhook(id string) error

	// ListWebhooks lists all webhook configurations
	ListWebhooks() ([]*WebhookConfig, error)
}

// FileStorage implements WebhookStorage using the file system
type FileStorage struct {
	baseDir string
	mu      sync.RWMutex
}

// NewFileStorage creates a new file storage
func NewFileStorage(baseDir string) (*FileStorage, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	return &FileStorage{
		baseDir: baseDir,
	}, nil
}

// SaveWebhook saves a webhook configuration to a file
func (s *FileStorage) SaveWebhook(config *WebhookConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Marshal webhook to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal webhook: %v", err)
	}

	// Write to file
	filePath := filepath.Join(s.baseDir, fmt.Sprintf("%s.json", config.ID))
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

// LoadWebhook loads a webhook configuration from a file
func (s *FileStorage) LoadWebhook(id string) (*WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Read file
	filePath := filepath.Join(s.baseDir, fmt.Sprintf("%s.json", id))
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Unmarshal webhook from JSON
	var config WebhookConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhook: %v", err)
	}

	return &config, nil
}

// DeleteWebhook deletes a webhook configuration file
func (s *FileStorage) DeleteWebhook(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete file
	filePath := filepath.Join(s.baseDir, fmt.Sprintf("%s.json", id))
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// ListWebhooks lists all webhook configurations from files
func (s *FileStorage) ListWebhooks() ([]*WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// List files in directory
	files, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	// Load webhooks from files
	webhooks := make([]*WebhookConfig, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Extract ID from filename
		id := file.Name()[:len(file.Name())-5] // Remove .json extension

		// Load webhook
		webhook, err := s.LoadWebhook(id)
		if err != nil {
			continue // Skip invalid webhooks
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// MemoryStorage implements WebhookStorage using in-memory storage
type MemoryStorage struct {
	webhooks map[string]*WebhookConfig
	mu       sync.RWMutex
}

// NewMemoryStorage creates a new memory storage
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		webhooks: make(map[string]*WebhookConfig),
	}
}

// SaveWebhook saves a webhook configuration to memory
func (s *MemoryStorage) SaveWebhook(config *WebhookConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a copy of the config
	webhookCopy := *config
	s.webhooks[config.ID] = &webhookCopy

	return nil
}

// LoadWebhook loads a webhook configuration from memory
func (s *MemoryStorage) LoadWebhook(id string) (*WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	webhook, exists := s.webhooks[id]
	if !exists {
		return nil, fmt.Errorf("webhook with ID %s does not exist", id)
	}

	// Create a copy of the webhook
	webhookCopy := *webhook
	return &webhookCopy, nil
}

// DeleteWebhook deletes a webhook configuration from memory
func (s *MemoryStorage) DeleteWebhook(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.webhooks[id]; !exists {
		return fmt.Errorf("webhook with ID %s does not exist", id)
	}

	delete(s.webhooks, id)
	return nil
}

// ListWebhooks lists all webhook configurations from memory
func (s *MemoryStorage) ListWebhooks() ([]*WebhookConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	webhooks := make([]*WebhookConfig, 0, len(s.webhooks))
	for _, webhook := range s.webhooks {
		// Create a copy of the webhook
		webhookCopy := *webhook
		webhooks = append(webhooks, &webhookCopy)
	}

	return webhooks, nil
}
