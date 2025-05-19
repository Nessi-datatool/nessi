package catalog

import (
	"fmt"
	"log"

	"github.com/nessi-dev/nessi/pkg/catalog/databricks"
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
	case Databricks:
		// Load configuration from environment variables
		config, err := databricks.NewConfigFromEnv()
		if err != nil {
			log.Printf("Warning: Failed to load Databricks configuration: %v", err)
			return nil, fmt.Errorf("failed to load Databricks configuration: %w", err)
		}
		
		return NewDatabricksCatalog(config.BaseURL, config.Token, config.WorkspaceID, config.DefaultSchema)
	default:
		return nil, fmt.Errorf("unsupported catalog type: %s", catalogType)
	}
}
