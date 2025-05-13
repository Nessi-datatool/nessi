package cloud

import (
	"context"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/cloud/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCloudProvider is a mock implementation of the CloudProvider interface
type MockCloudProvider struct {
	mock.Mock
}

func (m *MockCloudProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCloudProvider) Connect(ctx context.Context, config map[string]interface{}) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockCloudProvider) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCloudProvider) ListBuckets(ctx context.Context) ([]common.BucketInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]common.BucketInfo), args.Error(1)
}

func (m *MockCloudProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	args := m.Called(ctx, bucket, prefix)
	return args.Get(0).([]common.ObjectInfo), args.Error(1)
}

func (m *MockCloudProvider) GetObject(ctx context.Context, bucket, key string) (interface{}, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0), args.Error(1)
}

func (m *MockCloudProvider) PutObject(ctx context.Context, bucket, key string, data interface{}, metadata map[string]string) error {
	args := m.Called(ctx, bucket, key, data, metadata)
	return args.Error(0)
}

func (m *MockCloudProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	args := m.Called(ctx, bucket, key)
	return args.Error(0)
}

func (m *MockCloudProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCloudProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	args := m.Called(ctx, bucket, key, expiration)
	return args.String(0), args.Error(1)
}

// MockCloudProviderFactory is a mock implementation of the CloudProviderFactory interface
type MockCloudProviderFactory struct {
	mock.Mock
}

func (m *MockCloudProviderFactory) Create(config common.CloudConfig) (common.CloudProvider, error) {
	args := m.Called(config)
	return args.Get(0).(common.CloudProvider), args.Error(1)
}

// TestCloudManager_RegisterFactory tests the RegisterFactory method
func TestCloudManager_RegisterFactory(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Create a mock factory
	factory := new(MockCloudProviderFactory)
	
	// Register the factory
	manager.RegisterFactory("test", factory)
	
	// Verify that the factory was registered
	assert.Contains(t, manager.factories, "test")
	assert.Equal(t, factory, manager.factories["test"])
}

// TestCloudManager_CreateProvider tests the CreateProvider method
func TestCloudManager_CreateProvider(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Create a mock factory and provider
	factory := new(MockCloudProviderFactory)
	provider := new(MockCloudProvider)
	
	// Set up expectations
	config := common.CloudConfig{
		Provider: "test",
		AdditionalOptions: map[string]interface{}{
			"option1": "value1",
		},
		Credentials: map[string]interface{}{
			"key": "value",
		},
		DefaultBucket: "test-bucket",
	}
	
	factory.On("Create", config).Return(provider, nil)
	provider.On("Name").Return("test")
	
	// Register the factory
	manager.RegisterFactory("test", factory)
	
	// Create a provider
	createdProvider, err := manager.CreateProvider("test-provider", config)
	
	// Verify expectations
	assert.NoError(t, err)
	assert.Equal(t, provider, createdProvider)
	assert.Contains(t, manager.providers, "test-provider")
	assert.Equal(t, provider, manager.providers["test-provider"])
	factory.AssertExpectations(t)
	provider.AssertExpectations(t)
}

// TestCloudManager_GetProvider tests the GetProvider method
func TestCloudManager_GetProvider(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Create a mock provider
	provider := new(MockCloudProvider)
	
	// Add the provider to the manager
	manager.providers["test-provider"] = provider
	
	// Get the provider
	retrievedProvider, err := manager.GetProvider("test-provider")
	
	// Verify expectations
	assert.NoError(t, err)
	assert.Equal(t, provider, retrievedProvider)
}

// TestCloudManager_GetProvider_NotFound tests the GetProvider method with a non-existent provider
func TestCloudManager_GetProvider_NotFound(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Get a non-existent provider
	retrievedProvider, err := manager.GetProvider("non-existent")
	
	// Verify expectations
	assert.Error(t, err)
	assert.Nil(t, retrievedProvider)
}

// TestCloudManager_ListProviders tests the ListProviders method
func TestCloudManager_ListProviders(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Create mock providers
	provider1 := new(MockCloudProvider)
	provider2 := new(MockCloudProvider)
	
	// Set up expectations
	provider1.On("Name").Return("test1")
	provider2.On("Name").Return("test2")
	
	// Add the providers to the manager
	manager.providers["provider1"] = provider1
	manager.providers["provider2"] = provider2
	
	// List the providers
	providers := manager.ListProviders()
	
	// Verify expectations
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "provider1")
	assert.Contains(t, providers, "provider2")
	provider1.AssertExpectations(t)
	provider2.AssertExpectations(t)
}

// TestCloudManager_RemoveProvider tests the RemoveProvider method
func TestCloudManager_RemoveProvider(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Create a mock provider
	provider := new(MockCloudProvider)
	
	// Set up expectations
	provider.On("Disconnect", mock.Anything).Return(nil)
	
	// Add the provider to the manager
	manager.providers["test-provider"] = provider
	
	// Remove the provider
	err := manager.RemoveProvider(context.Background(), "test-provider")
	
	// Verify expectations
	assert.NoError(t, err)
	assert.NotContains(t, manager.providers, "test-provider")
	provider.AssertExpectations(t)
}

// TestCloudManager_RemoveProvider_NotFound tests the RemoveProvider method with a non-existent provider
func TestCloudManager_RemoveProvider_NotFound(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Remove a non-existent provider
	err := manager.RemoveProvider(context.Background(), "non-existent")
	
	// Verify expectations
	assert.Error(t, err)
}

// TestCloudManager_Close tests the Close method
func TestCloudManager_Close(t *testing.T) {
	// Create a new CloudManager
	manager := NewCloudManager()
	
	// Create mock providers
	provider1 := new(MockCloudProvider)
	provider2 := new(MockCloudProvider)
	
	// Set up expectations
	provider1.On("Disconnect", mock.Anything).Return(nil)
	provider2.On("Disconnect", mock.Anything).Return(nil)
	
	// Add the providers to the manager
	manager.providers["provider1"] = provider1
	manager.providers["provider2"] = provider2
	
	// Close the manager
	err := manager.Close(context.Background())
	
	// Verify expectations
	assert.NoError(t, err)
	assert.Empty(t, manager.providers)
	provider1.AssertExpectations(t)
	provider2.AssertExpectations(t)
}
