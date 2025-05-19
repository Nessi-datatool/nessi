package azure

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/cloud/common"
)

// AzureProvider implements the CloudProvider interface for Azure
type AzureProvider struct {
	connected   bool
	accountName string
	endpoint    string
}

// NewAzureProvider creates a new Azure provider
func NewAzureProvider() *AzureProvider {
	return &AzureProvider{
		connected: false,
	}
}

// Name returns the name of the cloud provider
func (p *AzureProvider) Name() string {
	return "Azure"
}

// Connect establishes a connection to Azure
func (p *AzureProvider) Connect(ctx context.Context, configMap map[string]interface{}) error {
	// Extract configuration
	accountName, _ := configMap["account_name"].(string)
	endpoint, _ := configMap["endpoint"].(string)

	if accountName == "" {
		return fmt.Errorf("account_name is required for Azure connection")
	}

	p.accountName = accountName
	if endpoint != "" {
		p.endpoint = endpoint
	} else {
		p.endpoint = fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	}

	// In a real implementation, we would establish a connection to Azure
	// For testing purposes, we'll just set connected to true
	p.connected = true

	return nil
}

// Disconnect closes the connection to Azure
func (p *AzureProvider) Disconnect(ctx context.Context) error {
	p.connected = false
	return nil
}

// ListBuckets lists all Azure storage containers
func (p *AzureProvider) ListBuckets(ctx context.Context) ([]common.BucketInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, we would list containers from Azure
	// For testing purposes, we'll just return mock data
	buckets := []common.BucketInfo{
		{
			Name:         "container1",
			CreationDate: time.Now().Format(time.RFC3339),
		},
		{
			Name:         "container2",
			CreationDate: time.Now().Format(time.RFC3339),
		},
	}

	return buckets, nil
}

// ListObjects lists all objects in an Azure storage container with the given prefix
func (p *AzureProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, we would list objects from Azure
	// For testing purposes, we'll just return mock data
	objects := []common.ObjectInfo{
		{
			Key:          "blob1.txt",
			Size:         1024,
			LastModified: time.Now().Format(time.RFC3339),
			ETag:         "etag1",
			ContentType:  "text/plain",
			Metadata:     map[string]string{"key1": "value1"},
		},
		{
			Key:          "blob2.txt",
			Size:         2048,
			LastModified: time.Now().Format(time.RFC3339),
			ETag:         "etag2",
			ContentType:  "text/plain",
			Metadata:     map[string]string{"key2": "value2"},
		},
	}

	// Filter by prefix if provided
	if prefix != "" {
		filteredObjects := []common.ObjectInfo{}
		for _, obj := range objects {
			if strings.HasPrefix(obj.Key, prefix) {
				filteredObjects = append(filteredObjects, obj)
			}
		}
		return filteredObjects, nil
	}

	return objects, nil
}

// GetObject retrieves an object from Azure storage
func (p *AzureProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, we would download the object from Azure
	// For testing purposes, we'll just return a string reader
	return io.NopCloser(strings.NewReader("mock content")), nil
}

// PutObject uploads an object to Azure storage
func (p *AzureProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	if !p.connected {
		return fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, we would upload the object to Azure
	// For testing purposes, we'll just return success
	return nil
}

// DeleteObject deletes an object from Azure storage
func (p *AzureProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	if !p.connected {
		return fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, we would delete the object from Azure
	// For testing purposes, we'll just return success
	return nil
}

// GetMetadata retrieves metadata for an Azure storage object
func (p *AzureProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, we would get the metadata from Azure
	// For testing purposes, we'll just return mock data
	metadata := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	return metadata, nil
}

// GetPresignedURL generates a SAS URL for an Azure storage object
func (p *AzureProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	if !p.connected {
		return "", fmt.Errorf("not connected to Azure")
	}

	// In a real implementation, you would generate a SAS token here using the Azure SDK
	// For testing purposes, we'll just return a placeholder URL
	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", p.accountName, bucket, key)
	return fmt.Sprintf("%s?sastoken=placeholder", blobURL), nil
}

// Factory for creating Azure providers
type AzureProviderFactory struct{}

// Create creates a new Azure provider instance
func (f *AzureProviderFactory) Create(config common.CloudConfig) (common.CloudProvider, error) {
	provider := NewAzureProvider()

	// Convert CloudConfig to map for the Connect method
	configMap := map[string]interface{}{
		"endpoint": config.EndpointOverride,
	}

	// Extract credentials if provided
	if config.Credentials != nil {
		if accountName, ok := config.Credentials["account_name"].(string); ok {
			configMap["account_name"] = accountName
		}
		if accountKey, ok := config.Credentials["account_key"].(string); ok {
			configMap["account_key"] = accountKey
		}
		if sasToken, ok := config.Credentials["sas_token"].(string); ok {
			configMap["sas_token"] = sasToken
		}
		if useAzureAD, ok := config.Credentials["use_azure_ad"].(bool); ok {
			configMap["use_azure_ad"] = useAzureAD
		}
	}

	// Add any additional options
	for k, v := range config.AdditionalOptions {
		configMap[k] = v
	}

	// Connect to Azure
	err := provider.Connect(context.Background(), configMap)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
