package gcp

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/nessi-dev/nessi/pkg/cloud/common"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GCPProvider implements the CloudProvider interface for Google Cloud Platform
type GCPProvider struct {
	client    *storage.Client
	connected bool
	projectID string
}

// NewGCPProvider creates a new GCP provider
func NewGCPProvider() *GCPProvider {
	return &GCPProvider{
		connected: false,
	}
}

// Name returns the name of the cloud provider
func (p *GCPProvider) Name() string {
	return "GCP"
}

// Connect establishes a connection to GCP
func (p *GCPProvider) Connect(ctx context.Context, configMap map[string]interface{}) error {
	// Extract configuration
	projectID, _ := configMap["project_id"].(string)
	credentialsFile, _ := configMap["credentials_file"].(string)
	credentialsJSON, _ := configMap["credentials_json"].(string)

	p.projectID = projectID

	var opts []option.ClientOption

	// Set credentials if provided
	if credentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsFile))
	} else if credentialsJSON != "" {
		opts = append(opts, option.WithCredentialsJSON([]byte(credentialsJSON)))
	}

	// Create storage client
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return fmt.Errorf("failed to create GCP storage client: %w", err)
	}

	p.client = client
	p.connected = true

	return nil
}

// Disconnect closes the connection to GCP
func (p *GCPProvider) Disconnect(ctx context.Context) error {
	if !p.connected {
		return nil
	}

	err := p.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close GCP storage client: %w", err)
	}

	p.connected = false
	return nil
}

// ListBuckets lists all GCP storage buckets
func (p *GCPProvider) ListBuckets(ctx context.Context) ([]common.BucketInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to GCP")
	}

	var buckets []common.BucketInfo
	it := p.client.Buckets(ctx, p.projectID)

	for {
		bucketAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list buckets: %w", err)
		}

		buckets = append(buckets, common.BucketInfo{
			Name:         bucketAttrs.Name,
			CreationDate: bucketAttrs.Created.Format(time.RFC3339),
			Region:       bucketAttrs.Location,
		})
	}

	return buckets, nil
}

// ListObjects lists all objects in a GCP storage bucket with the given prefix
func (p *GCPProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to GCP")
	}

	bkt := p.client.Bucket(bucket)
	query := &storage.Query{Prefix: prefix}

	var objects []common.ObjectInfo
	it := bkt.Objects(ctx, query)

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in bucket %s: %w", bucket, err)
		}

		objects = append(objects, common.ObjectInfo{
			Key:          attrs.Name,
			Size:         attrs.Size,
			LastModified: attrs.Updated.Format(time.RFC3339),
			ETag:         attrs.Etag,
			ContentType:  attrs.ContentType,
			Metadata:     attrs.Metadata,
		})
	}

	return objects, nil
}

// GetObject retrieves an object from GCP storage
func (p *GCPProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to GCP")
	}

	bkt := p.client.Bucket(bucket)
	obj := bkt.Object(key)

	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s from bucket %s: %w", key, bucket, err)
	}

	return reader, nil
}

// PutObject uploads an object to GCP storage
func (p *GCPProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	if !p.connected {
		return fmt.Errorf("not connected to GCP")
	}

	bkt := p.client.Bucket(bucket)
	obj := bkt.Object(key)

	writer := obj.NewWriter(ctx)
	writer.Metadata = metadata

	if _, err := io.Copy(writer, data); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write data to object %s in bucket %s: %w", key, bucket, err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer for object %s in bucket %s: %w", key, bucket, err)
	}

	return nil
}

// DeleteObject deletes an object from GCP storage
func (p *GCPProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	if !p.connected {
		return fmt.Errorf("not connected to GCP")
	}

	bkt := p.client.Bucket(bucket)
	obj := bkt.Object(key)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object %s from bucket %s: %w", key, bucket, err)
	}

	return nil
}

// GetMetadata retrieves metadata for a GCP storage object
func (p *GCPProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to GCP")
	}

	bkt := p.client.Bucket(bucket)
	obj := bkt.Object(key)

	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get attributes for object %s in bucket %s: %w", key, bucket, err)
	}

	return attrs.Metadata, nil
}

// GetPresignedURL generates a signed URL for a GCP storage object
func (p *GCPProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	if !p.connected {
		return "", fmt.Errorf("not connected to GCP")
	}

	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expiration) * time.Second),
	}

	// In the latest GCP SDK, SignedURL is a function in the storage package, not a method on the object
	url, err := storage.SignedURL(bucket, key, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for object %s in bucket %s: %w", key, bucket, err)
	}

	return url, nil
}

// Factory for creating GCP providers
type GCPProviderFactory struct{}

// Create creates a new GCP provider instance
func (f *GCPProviderFactory) Create(config common.CloudConfig) (common.CloudProvider, error) {
	provider := NewGCPProvider()

	// Convert CloudConfig to map for the Connect method
	configMap := map[string]interface{}{
		"project_id": config.AdditionalOptions["project_id"],
	}

	// Extract credentials if provided
	if config.Credentials != nil {
		if credentialsFile, ok := config.Credentials["credentials_file"].(string); ok {
			configMap["credentials_file"] = credentialsFile
		}
		if credentialsJSON, ok := config.Credentials["credentials_json"].(string); ok {
			configMap["credentials_json"] = credentialsJSON
		}
	}

	// Add any additional options
	for k, v := range config.AdditionalOptions {
		configMap[k] = v
	}

	// Connect to GCP
	err := provider.Connect(context.Background(), configMap)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
