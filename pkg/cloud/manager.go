package cloud

import (
	"fmt"
	"sync"

	"github.com/nessi-dev/nessi-dev/pkg/cloud/aws"
	"github.com/nessi-dev/nessi-dev/pkg/cloud/azure"
	"github.com/nessi-dev/nessi-dev/pkg/cloud/common"
	"github.com/nessi-dev/nessi-dev/pkg/cloud/gcp"
)

// CloudManager manages cloud provider instances
type CloudManager struct {
	factories map[string]common.CloudProviderFactory
	providers map[string]common.CloudProvider
	mu        sync.RWMutex
}

// NewCloudManager creates a new CloudManager
func NewCloudManager() *CloudManager {
	manager := &CloudManager{
		factories: make(map[string]common.CloudProviderFactory),
		providers: make(map[string]common.CloudProvider),
	}

	// Register default cloud providers
	manager.RegisterProvider("aws", &aws.AWSProviderFactory{})
	manager.RegisterProvider("azure", &azure.AzureProviderFactory{})
	manager.RegisterProvider("gcp", &gcp.GCPProviderFactory{})

	return manager
}

// RegisterProvider registers a cloud provider factory
func (m *CloudManager) RegisterProvider(name string, factory common.CloudProviderFactory) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.factories[name] = factory
}

// GetProvider returns a cloud provider instance
func (m *CloudManager) GetProvider(name string) (common.CloudProvider, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	provider, exists := m.providers[name]
	return provider, exists
}

// CreateProvider creates a new cloud provider instance
func (m *CloudManager) CreateProvider(name string, config common.CloudConfig) (common.CloudProvider, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if provider already exists
	if provider, exists := m.providers[name]; exists {
		return provider, nil
	}

	// Check if factory exists
	factory, exists := m.factories[config.Provider]
	if !exists {
		return nil, fmt.Errorf("cloud provider '%s' not supported", config.Provider)
	}

	// Create provider
	provider, err := factory.Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloud provider '%s': %w", config.Provider, err)
	}

	// Store provider
	m.providers[name] = provider
	return provider, nil
}

// RemoveProvider removes a cloud provider instance
func (m *CloudManager) RemoveProvider(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[name]
	if !exists {
		return fmt.Errorf("cloud provider '%s' not found", name)
	}

	// Disconnect provider
	if err := provider.Disconnect(nil); err != nil {
		return fmt.Errorf("failed to disconnect cloud provider '%s': %w", name, err)
	}

	// Remove provider
	delete(m.providers, name)
	return nil
}

// ListProviders returns a list of registered cloud providers
func (m *CloudManager) ListProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var providers []string
	for name := range m.providers {
		providers = append(providers, name)
	}
	return providers
}

// ListSupportedProviders returns a list of supported cloud providers
func (m *CloudManager) ListSupportedProviders() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var providers []string
	for name := range m.factories {
		providers = append(providers, name)
	}
	return providers
}

// GetProviderByType returns a provider by type (aws, azure, gcp)
func (m *CloudManager) GetProviderByType(providerType string) (common.CloudProvider, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, provider := range m.providers {
		if provider.Name() == providerType {
			return provider, true
		}
	}
	return nil, false
}

// Shutdown disconnects all cloud providers
func (m *CloudManager) Shutdown() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for name, provider := range m.providers {
		if err := provider.Disconnect(nil); err != nil {
			errs = append(errs, fmt.Errorf("failed to disconnect cloud provider '%s': %w", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to shutdown cloud manager: %v", errs)
	}
	return nil
}
