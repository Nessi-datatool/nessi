package azure

import (
	"context"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/cloud/common"
	"github.com/stretchr/testify/assert"
)

// No need for mock implementations since our provider implementation is simplified

// TestAzureProvider_Name tests the Name method
func TestAzureProvider_Name(t *testing.T) {
	provider := NewAzureProvider()
	assert.Equal(t, "Azure", provider.Name())
}

// TestAzureProvider_ListBuckets tests the ListBuckets method
func TestAzureProvider_ListBuckets(t *testing.T) {
	// Create provider and set it as connected
	provider := NewAzureProvider()
	provider.connected = true
	provider.accountName = "testaccount"
	
	// Call ListBuckets
	buckets, err := provider.ListBuckets(context.Background())
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, buckets, 2)
	assert.Equal(t, "container1", buckets[0].Name)
	assert.Equal(t, "container2", buckets[1].Name)
}

// TestAzureProvider_ListObjects tests the ListObjects method
func TestAzureProvider_ListObjects(t *testing.T) {
	// Create provider and set it as connected
	provider := NewAzureProvider()
	provider.connected = true
	provider.accountName = "testaccount"
	
	// Call ListObjects
	objects, err := provider.ListObjects(context.Background(), "container1", "")
	
	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	assert.Equal(t, "blob1.txt", objects[0].Key)
	assert.Equal(t, int64(1024), objects[0].Size)
	assert.Equal(t, "etag1", objects[0].ETag)
	assert.Equal(t, "text/plain", objects[0].ContentType)
	assert.Equal(t, "value1", objects[0].Metadata["key1"])
}

// TestAzureProviderFactory_Create tests the AzureProviderFactory.Create method
func TestAzureProviderFactory_Create(t *testing.T) {
	// Create a factory
	factory := &AzureProviderFactory{}
	
	// Create test config
	config := common.CloudConfig{
		Provider: "azure",
		Credentials: map[string]interface{}{
			"account_name": "teststorageaccount",
		},
		DefaultBucket: "test-container",
	}
	
	// Test creating a provider
	provider, err := factory.Create(config)
	
	// We should not get an error with our simplified implementation
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	assert.Equal(t, "Azure", provider.Name())
}

// No helper functions needed for our simplified implementation
