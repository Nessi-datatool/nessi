package types

import (
	"fmt"
)

// CatalogType represents the type of data catalog
type CatalogType string

const (
	// AWSGlue represents AWS Glue Data Catalog
	AWSGlue CatalogType = "aws_glue"
	
	// AzurePurview represents Azure Purview
	AzurePurview CatalogType = "azure_purview"
	
	// GCPDataCatalog represents Google Cloud Data Catalog
	GCPDataCatalog CatalogType = "gcp_data_catalog"
	
	// ApacheAtlas represents Apache Atlas
	ApacheAtlas CatalogType = "apache_atlas"
	
	// Collibra represents Collibra Data Catalog
	Collibra CatalogType = "collibra"
)

// CatalogFactory creates data catalog providers
type CatalogFactory struct {
	providers map[CatalogType]func() DataCatalog
}


// RegisterProvider registers a provider function for a catalog type
func (f *CatalogFactory) RegisterProvider(catalogType CatalogType, provider func() DataCatalog) {
	f.providers[catalogType] = provider
}

// CreateCatalog creates a data catalog provider of the specified type
func (f *CatalogFactory) CreateCatalog(catalogType CatalogType) (DataCatalog, error) {
	provider, ok := f.providers[catalogType]
	if !ok {
		return nil, fmt.Errorf("unsupported catalog type: %s", catalogType)
	}
	
	return provider(), nil
}

// RegisterAllCatalogs registers all supported catalog types with the manager
func (f *CatalogFactory) RegisterAllCatalogs(manager *CatalogManager) error {
	for catalogType := range f.providers {
		catalog, err := f.CreateCatalog(catalogType)
		if err != nil {
			return fmt.Errorf("failed to create catalog %s: %w", catalogType, err)
		}
		
		if err := manager.RegisterCatalog(catalog); err != nil {
			return fmt.Errorf("failed to register catalog %s: %w", catalogType, err)
		}
	}
	
	return nil
}
