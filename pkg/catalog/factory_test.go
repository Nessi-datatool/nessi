package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCatalogFactory_CreateCatalog tests the CreateCatalog method
func TestCatalogFactory_CreateCatalog(t *testing.T) {
	// Create catalog factory
	factory := NewCatalogFactory()
	
	// Test creating AWS Glue catalog
	awsGlue, err := factory.CreateCatalog(AWSGlue)
	assert.NoError(t, err)
	assert.NotNil(t, awsGlue)
	assert.Contains(t, awsGlue.Name(), "AWS Glue")
	
	// Test creating Azure Purview catalog
	azurePurview, err := factory.CreateCatalog(AzurePurview)
	assert.NoError(t, err)
	assert.NotNil(t, azurePurview)
	assert.Contains(t, azurePurview.Name(), "Azure Purview")
	
	// Test creating GCP Data Catalog
	gcpDataCatalog, err := factory.CreateCatalog(GCPDataCatalog)
	assert.NoError(t, err)
	assert.NotNil(t, gcpDataCatalog)
	assert.Contains(t, gcpDataCatalog.Name(), "Google Cloud Data Catalog")
	
	// Test creating Apache Atlas catalog (not implemented yet)
	apacheAtlas, err := factory.CreateCatalog(ApacheAtlas)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, apacheAtlas)
	
	// Test creating Collibra catalog (not implemented yet)
	collibra, err := factory.CreateCatalog(Collibra)
	assert.Equal(t, ErrNotImplemented, err)
	assert.Nil(t, collibra)
	
	// Test creating unsupported catalog type
	unsupported, err := factory.CreateCatalog("unsupported")
	assert.Equal(t, ErrUnsupportedCatalogType, err)
	assert.Nil(t, unsupported)
}

// TestCatalogFactory_RegisterAllCatalogs tests the RegisterAllCatalogs method
func TestCatalogFactory_RegisterAllCatalogs(t *testing.T) {
	// Create catalog factory
	factory := NewCatalogFactory()
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register all catalogs
	err := factory.RegisterAllCatalogs(manager)
	assert.NoError(t, err)
	
	// Verify catalogs are registered
	catalogs := manager.ListCatalogs()
	
	// We should have at least AWS Glue, Azure Purview, and GCP Data Catalog
	assert.GreaterOrEqual(t, len(catalogs), 3)
	
	// Check for specific catalog names
	var hasAWSGlue, hasAzurePurview, hasGCPDataCatalog bool
	
	for _, name := range catalogs {
		if name == "AWS Glue Data Catalog" {
			hasAWSGlue = true
		} else if name == "Azure Purview Data Catalog" {
			hasAzurePurview = true
		} else if name == "Google Cloud Data Catalog" {
			hasGCPDataCatalog = true
		}
	}
	
	assert.True(t, hasAWSGlue, "AWS Glue catalog should be registered")
	assert.True(t, hasAzurePurview, "Azure Purview catalog should be registered")
	assert.True(t, hasGCPDataCatalog, "GCP Data Catalog should be registered")
}
