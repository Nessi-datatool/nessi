package catalog

import (
	"fmt"
)

// CatalogFactory creates instances of data catalogs
type CatalogFactory struct{}

// NewCatalogFactory creates a new catalog factory
func NewCatalogFactory() *CatalogFactory {
	return &CatalogFactory{}
}

// CreateCatalog creates a data catalog of the specified type
func (f *CatalogFactory) CreateCatalog(catalogType CatalogType) (DataCatalog, error) {
	switch catalogType {
	case AWSGlue:
		return NewAWSGlueCatalog(), nil
	case AzurePurview:
		return NewAzurePurviewCatalog(), nil
	case GCPDataCatalog:
		return NewGCPDataCatalog(), nil
	default:
		return nil, fmt.Errorf("unsupported catalog type: %s", catalogType)
	}
}
