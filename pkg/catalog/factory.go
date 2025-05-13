package catalog

import (
	"github.com/nessi-dev/nessi-dev/pkg/catalog/aws"
	"github.com/nessi-dev/nessi-dev/pkg/catalog/azure"
	"github.com/nessi-dev/nessi-dev/pkg/catalog/gcp"
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
type CatalogFactory struct{}

// NewCatalogFactory creates a new catalog factory
func NewCatalogFactory() *CatalogFactory {
	return &CatalogFactory{}
}

// CreateCatalog creates a data catalog provider of the specified type
func (f *CatalogFactory) CreateCatalog(catalogType CatalogType) (DataCatalog, error) {
	switch catalogType {
	case AWSGlue:
		return aws.NewGlueCatalog(), nil
	case AzurePurview:
		return azure.NewPurviewCatalog(), nil
	case GCPDataCatalog:
		return gcp.NewDataCatalogClient(), nil
	case ApacheAtlas:
		// TODO: Implement Apache Atlas catalog
		return nil, ErrNotImplemented
	case Collibra:
		// TODO: Implement Collibra catalog
		return nil, ErrNotImplemented
	default:
		return nil, ErrUnsupportedCatalogType
	}
}

// RegisterAllCatalogs registers all supported catalog types with the manager
func (f *CatalogFactory) RegisterAllCatalogs(manager *CatalogManager) error {
	// Register AWS Glue
	awsGlue, err := f.CreateCatalog(AWSGlue)
	if err != nil && err != ErrNotImplemented {
		return err
	}
	if err == nil {
		manager.RegisterCatalog(awsGlue)
	}
	
	// Register Azure Purview
	azurePurview, err := f.CreateCatalog(AzurePurview)
	if err != nil && err != ErrNotImplemented {
		return err
	}
	if err == nil {
		manager.RegisterCatalog(azurePurview)
	}
	
	// Register GCP Data Catalog
	gcpDataCatalog, err := f.CreateCatalog(GCPDataCatalog)
	if err != nil && err != ErrNotImplemented {
		return err
	}
	if err == nil {
		manager.RegisterCatalog(gcpDataCatalog)
	}
	
	// Register Apache Atlas (when implemented)
	apacheAtlas, err := f.CreateCatalog(ApacheAtlas)
	if err != nil && err != ErrNotImplemented {
		return err
	}
	if err == nil {
		manager.RegisterCatalog(apacheAtlas)
	}
	
	// Register Collibra (when implemented)
	collibra, err := f.CreateCatalog(Collibra)
	if err != nil && err != ErrNotImplemented {
		return err
	}
	if err == nil {
		manager.RegisterCatalog(collibra)
	}
	
	return nil
}
