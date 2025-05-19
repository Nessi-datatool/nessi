package types

import "sync"

var (
	defaultFactory *CatalogFactory
	factoryOnce    sync.Once
)

// GetCatalogFactory returns the singleton instance of CatalogFactory
func GetCatalogFactory() *CatalogFactory {
	factoryOnce.Do(func() {
		defaultFactory = &CatalogFactory{
			providers: make(map[CatalogType]func() DataCatalog),
		}
	})
	return defaultFactory
}

// NewCatalogFactory creates a new catalog factory
// Deprecated: Use GetCatalogFactory() instead
func NewCatalogFactory() *CatalogFactory {
	return GetCatalogFactory()
}
