package tests

import (
	"testing"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"

	// Import providers to register them
	_ "github.com/nessi-dev/nessi/pkg/catalog/aws"
	_ "github.com/nessi-dev/nessi/pkg/catalog/azure"
	_ "github.com/nessi-dev/nessi/pkg/catalog/gcp"
)

// TestCatalogFactory_CreateCatalog tests the CreateCatalog method
func TestCatalogFactory_CreateCatalog(t *testing.T) {
	// Get catalog factory
	factory := types.GetCatalogFactory()
	
	// Test creating AWS Glue catalog
	catalog, err := factory.CreateCatalog(types.AWSGlue)
	assert.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Contains(t, catalog.Name(), "AWS Glue")
	
	// Test creating Azure Purview catalog
	catalog, err = factory.CreateCatalog(types.AzurePurview)
	assert.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Contains(t, catalog.Name(), "Azure Purview")
	
	// Test creating GCP Data Catalog
	catalog, err = factory.CreateCatalog(types.GCPDataCatalog)
	assert.NoError(t, err)
	assert.NotNil(t, catalog)
	assert.Contains(t, catalog.Name(), "GCP Data Catalog")
	
	// Test creating unsupported catalog type
	catalog, err = factory.CreateCatalog("invalid")
	assert.Error(t, err)
	assert.Nil(t, catalog)
	assert.Contains(t, err.Error(), "unsupported catalog type")
}
