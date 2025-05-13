package azure

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/nessi-dev/nessi-dev/pkg/cloud/common"
)

// AzureProvider implements the CloudProvider interface for Azure
type AzureProvider struct {
	serviceClient *service.Client
	connected     bool
	accountName   string
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
	accountKey, _ := configMap["account_key"].(string)
	sasToken, _ := configMap["sas_token"].(string)
	useAzureAD, _ := configMap["use_azure_ad"].(bool)
	endpoint, _ := configMap["endpoint"].(string)

	p.accountName = accountName

	var client *service.Client
	var err error

	if accountName == "" {
		return fmt.Errorf("account_name is required for Azure connection")
	}

	// Create service client based on authentication method
	if accountKey != "" {
		// Use account key authentication
		cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
			return fmt.Errorf("failed to create shared key credential: %w", err)
		}

		serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
		if endpoint != "" {
			serviceURL = endpoint
		}

		client, err = service.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
		if err != nil {
			return fmt.Errorf("failed to create service client with shared key: %w", err)
		}
	} else if sasToken != "" {
		// Use SAS token authentication
		serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", accountName, sasToken)
		if endpoint != "" {
			serviceURL = fmt.Sprintf("%s?%s", endpoint, sasToken)
		}

		client, err = service.NewClientWithNoCredential(serviceURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create service client with SAS token: %w", err)
		}
	} else if useAzureAD {
		// Use Azure AD authentication
		cred, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return fmt.Errorf("failed to create Azure AD credential: %w", err)
		}

		serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
		if endpoint != "" {
			serviceURL = endpoint
		}

		client, err = service.NewClient(serviceURL, cred, nil)
		if err != nil {
			return fmt.Errorf("failed to create service client with Azure AD: %w", err)
		}
	} else {
		return fmt.Errorf("no valid authentication method provided for Azure")
	}

	p.serviceClient = client
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

	pager := p.serviceClient.NewListContainersPager(nil)
	
	var buckets []common.BucketInfo
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list containers: %w", err)
		}

		for _, container := range resp.ContainerItems {
			buckets = append(buckets, common.BucketInfo{
				Name:         *container.Name,
				CreationDate: container.Properties.LastModified.Format(time.RFC3339),
			})
		}
	}

	return buckets, nil
}

// ListObjects lists all objects in an Azure storage container with the given prefix
func (p *AzureProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	containerClient, err := p.serviceClient.NewContainerClient(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client for %s: %w", bucket, err)
	}

	options := &container.ListBlobsFlatOptions{}
	if prefix != "" {
		options.Prefix = &prefix
	}

	pager := containerClient.NewListBlobsFlatPager(options)
	
	var objects []common.ObjectInfo
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list blobs in container %s: %w", bucket, err)
		}

		for _, blob := range resp.Segment.BlobItems {
			metadata := make(map[string]string)
			for k, v := range blob.Metadata {
				if v != nil {
					metadata[k] = *v
				}
			}

			objects = append(objects, common.ObjectInfo{
				Key:          *blob.Name,
				Size:         *blob.Properties.ContentLength,
				LastModified: blob.Properties.LastModified.Format(time.RFC3339),
				ETag:         *blob.Properties.ETag,
				ContentType:  *blob.Properties.ContentType,
				Metadata:     metadata,
			})
		}
	}

	return objects, nil
}

// GetObject retrieves an object from Azure storage
func (p *AzureProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	containerClient, err := p.serviceClient.NewContainerClient(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client for %s: %w", bucket, err)
	}

	blobClient, err := containerClient.NewBlobClient(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client for %s: %w", key, err)
	}

	downloadResponse, err := blobClient.Download(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob %s from container %s: %w", key, bucket, err)
	}

	return downloadResponse.Body, nil
}

// PutObject uploads an object to Azure storage
func (p *AzureProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	if !p.connected {
		return fmt.Errorf("not connected to Azure")
	}

	containerClient, err := p.serviceClient.NewContainerClient(bucket)
	if err != nil {
		return fmt.Errorf("failed to create container client for %s: %w", bucket, err)
	}

	blobClient, err := containerClient.NewBlobClient(key)
	if err != nil {
		return fmt.Errorf("failed to create blob client for %s: %w", key, err)
	}

	// Convert metadata map to Azure format
	azureMetadata := make(map[string]*string)
	for k, v := range metadata {
		value := v
		azureMetadata[k] = &value
	}

	uploadOptions := &azblob.UploadOptions{
		Metadata: azureMetadata,
	}

	_, err = blobClient.Upload(ctx, data, uploadOptions)
	if err != nil {
		return fmt.Errorf("failed to upload blob %s to container %s: %w", key, bucket, err)
	}

	return nil
}

// DeleteObject deletes an object from Azure storage
func (p *AzureProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	if !p.connected {
		return fmt.Errorf("not connected to Azure")
	}

	containerClient, err := p.serviceClient.NewContainerClient(bucket)
	if err != nil {
		return fmt.Errorf("failed to create container client for %s: %w", bucket, err)
	}

	blobClient, err := containerClient.NewBlobClient(key)
	if err != nil {
		return fmt.Errorf("failed to create blob client for %s: %w", key, err)
	}

	_, err = blobClient.Delete(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob %s from container %s: %w", key, bucket, err)
	}

	return nil
}

// GetMetadata retrieves metadata for an Azure storage object
func (p *AzureProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to Azure")
	}

	containerClient, err := p.serviceClient.NewContainerClient(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to create container client for %s: %w", bucket, err)
	}

	blobClient, err := containerClient.NewBlobClient(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create blob client for %s: %w", key, err)
	}

	props, err := blobClient.GetProperties(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get properties for blob %s in container %s: %w", key, bucket, err)
	}

	metadata := make(map[string]string)
	for k, v := range props.Metadata {
		metadata[k] = v
	}

	return metadata, nil
}

// GetPresignedURL generates a SAS URL for an Azure storage object
func (p *AzureProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	if !p.connected {
		return "", fmt.Errorf("not connected to Azure")
	}

	// For Azure, we need the account key to generate SAS tokens
	// This is a simplified implementation and may need to be adjusted based on your authentication method
	containerClient, err := p.serviceClient.NewContainerClient(bucket)
	if err != nil {
		return "", fmt.Errorf("failed to create container client for %s: %w", bucket, err)
	}

	blobClient, err := containerClient.NewBlobClient(key)
	if err != nil {
		return "", fmt.Errorf("failed to create blob client for %s: %w", key, err)
	}

	// Get the blob URL
	blobURL := blobClient.URL()

	// In a real implementation, you would generate a SAS token here
	// For now, we'll just return the blob URL with a note
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
