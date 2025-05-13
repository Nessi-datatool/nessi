package catalog

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
func (m *CatalogManager) RegisterCatalog(catalog DataCatalog) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.catalogs[catalog.Name()] = catalog
}

// GetCatalog returns a data catalog provider by name
func (m *CatalogManager) GetCatalog(name string) (DataCatalog, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	catalog, ok := m.catalogs[name]
	if !ok {
		return nil, fmt.Errorf("catalog %s not found", name)
	}
	
	return catalog, nil
}

// ListCatalogs returns a list of registered catalog providers
func (m *CatalogManager) ListCatalogs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var names []string
	for name := range m.catalogs {
		names = append(names, name)
	}
	
	return names
}

// ConnectCatalog connects to a data catalog
func (m *CatalogManager) ConnectCatalog(ctx context.Context, name string, config map[string]interface{}) error {
	catalog, err := m.GetCatalog(name)
	if err != nil {
		return err
	}
	
	return catalog.Connect(ctx, config)
}

// DisconnectCatalog disconnects from a data catalog
func (m *CatalogManager) DisconnectCatalog(ctx context.Context, name string) error {
	catalog, err := m.GetCatalog(name)
	if err != nil {
		return err
	}
	
	return catalog.Disconnect(ctx)
}

// DisconnectAll disconnects from all data catalogs
func (m *CatalogManager) DisconnectAll(ctx context.Context) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	for _, catalog := range m.catalogs {
		_ = catalog.Disconnect(ctx)
	}
}
