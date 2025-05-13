package common

import (
	"context"
	"io"
)

// CloudProvider defines the interface that all cloud providers must implement
type CloudProvider interface {
	// Name returns the name of the cloud provider
	Name() string

	// Connect establishes a connection to the cloud provider
	Connect(ctx context.Context, config map[string]interface{}) error

	// Disconnect closes the connection to the cloud provider
	Disconnect(ctx context.Context) error

	// ListBuckets lists all buckets/containers in the cloud storage
	ListBuckets(ctx context.Context) ([]BucketInfo, error)

	// ListObjects lists all objects in a bucket with the given prefix
	ListObjects(ctx context.Context, bucket, prefix string) ([]ObjectInfo, error)

	// GetObject retrieves an object from cloud storage
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)

	// PutObject uploads an object to cloud storage
	PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error

	// DeleteObject deletes an object from cloud storage
	DeleteObject(ctx context.Context, bucket, key string) error

	// GetMetadata retrieves metadata for an object
	GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error)

	// GetPresignedURL generates a presigned URL for an object
	GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error)
}

// BucketInfo represents information about a bucket/container
type BucketInfo struct {
	Name         string
	CreationDate string
	Region       string
}

// ObjectInfo represents information about an object in cloud storage
type ObjectInfo struct {
	Key          string
	Size         int64
	LastModified string
	ETag         string
	ContentType  string
	Metadata     map[string]string
}

// CloudConfig contains configuration for cloud providers
type CloudConfig struct {
	Provider          string                 `json:"provider" yaml:"provider"`
	Region            string                 `json:"region" yaml:"region"`
	Credentials       map[string]interface{} `json:"credentials" yaml:"credentials"`
	DefaultBucket     string                 `json:"default_bucket" yaml:"default_bucket"`
	EndpointOverride  string                 `json:"endpoint_override,omitempty" yaml:"endpoint_override,omitempty"`
	AdditionalOptions map[string]interface{} `json:"additional_options,omitempty" yaml:"additional_options,omitempty"`
}

// CloudProviderFactory creates cloud provider instances
type CloudProviderFactory interface {
	// Create creates a new cloud provider instance
	Create(config CloudConfig) (CloudProvider, error)
}
