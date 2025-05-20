package types

import (
	"context"
	"fmt"
	"sync"
)

// CatalogManager manages data catalog providers
type CatalogManager struct {
	catalogs map[string]DataCatalog
	mu       sync.RWMutex
}

// NewCatalogManager creates a new catalog manager
func NewCatalogManager() *CatalogManager {
	return &CatalogManager{
		catalogs: make(map[string]DataCatalog),
	}
}

// RegisterCatalog registers a data catalog provider
func (m *CatalogManager) RegisterCatalog(catalog DataCatalog) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := catalog.Name()
	if _, exists := m.catalogs[name]; exists {
		return fmt.Errorf("catalog %s already registered", name)
	}

	m.catalogs[name] = catalog
	return nil
}

// UnregisterCatalog unregisters a data catalog provider
func (m *CatalogManager) UnregisterCatalog(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.catalogs[name]; !exists {
		return fmt.Errorf("catalog %s not found", name)
	}

	delete(m.catalogs, name)
	return nil
}

// GetCatalog gets a data catalog provider by name
func (m *CatalogManager) GetCatalog(name string) (DataCatalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	catalog, exists := m.catalogs[name]
	if !exists {
		return nil, fmt.Errorf("catalog %s not found", name)
	}

	return catalog, nil
}

// ListCatalogs lists all registered data catalog providers
func (m *CatalogManager) ListCatalogs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.catalogs))
	for name := range m.catalogs {
		names = append(names, name)
	}

	return names
}

// ConnectAll connects to all registered catalogs
func (m *CatalogManager) ConnectAll(ctx context.Context, configs map[string]map[string]interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, catalog := range m.catalogs {
		config, ok := configs[name]
		if !ok {
			continue
		}

		if err := catalog.Connect(ctx, config); err != nil {
			return fmt.Errorf("failed to connect to catalog %s: %w", name, err)
		}
	}

	return nil
}

// DisconnectAll disconnects from all registered catalogs
func (m *CatalogManager) DisconnectAll(ctx context.Context) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, catalog := range m.catalogs {
		if err := catalog.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from catalog %s: %w", name, err)
		}
	}

	return nil
}
