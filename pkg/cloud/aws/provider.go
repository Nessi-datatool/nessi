package aws

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nessi-dev/nessi-dev/pkg/cloud/common"
)

// AWSProvider implements the CloudProvider interface for AWS
type AWSProvider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	config        aws.Config
	connected     bool
}

// NewAWSProvider creates a new AWS provider
func NewAWSProvider() *AWSProvider {
	return &AWSProvider{
		connected: false,
	}
}

// Name returns the name of the cloud provider
func (p *AWSProvider) Name() string {
	return "AWS"
}

// Connect establishes a connection to AWS
func (p *AWSProvider) Connect(ctx context.Context, configMap map[string]interface{}) error {
	var err error
	var cfgOpts []func(*config.LoadOptions) error

	// Extract configuration
	region, _ := configMap["region"].(string)
	accessKey, _ := configMap["access_key"].(string)
	secretKey, _ := configMap["secret_key"].(string)
	sessionToken, _ := configMap["session_token"].(string)
	endpoint, _ := configMap["endpoint"].(string)
	useIAMRole, _ := configMap["use_iam_role"].(bool)

	// Set region
	if region != "" {
		cfgOpts = append(cfgOpts, config.WithRegion(region))
	}

	// Set credentials if provided
	if accessKey != "" && secretKey != "" {
		cfgOpts = append(cfgOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken),
		))
	} else if !useIAMRole {
		// Use default credential chain if not using IAM role and no explicit credentials
		cfgOpts = append(cfgOpts, config.WithSharedConfigProfile("default"))
	}

	// Set custom endpoint if provided (for testing or S3-compatible services)
	if endpoint != "" {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})
		cfgOpts = append(cfgOpts, config.WithEndpointResolverWithOptions(customResolver))
	}

	// Load AWS configuration
	p.config, err = config.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create S3 client
	p.client = s3.NewFromConfig(p.config)
	p.presignClient = s3.NewPresignClient(p.client)
	
	// Test the connection by listing buckets
	if endpoint, ok := configMap["endpoint"].(string); ok && endpoint != "" {
		// For testing purposes, if a custom endpoint is provided, verify it works
		// This is mainly for unit tests to ensure we get an error with invalid endpoints
		_, err := p.client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			return fmt.Errorf("failed to connect to endpoint %s: %w", endpoint, err)
		}
	}
	
	p.connected = true

	return nil
}

// Disconnect closes the connection to AWS
func (p *AWSProvider) Disconnect(ctx context.Context) error {
	p.connected = false
	return nil
}

// ListBuckets lists all S3 buckets
func (p *AWSProvider) ListBuckets(ctx context.Context) ([]common.BucketInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to AWS")
	}

	result, err := p.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	buckets := make([]common.BucketInfo, 0, len(result.Buckets))
	for _, bucket := range result.Buckets {
		buckets = append(buckets, common.BucketInfo{
			Name:         *bucket.Name,
			CreationDate: bucket.CreationDate.Format(time.RFC3339),
		})
	}

	return buckets, nil
}

// ListObjects lists all objects in an S3 bucket with the given prefix
func (p *AWSProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to AWS")
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	result, err := p.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket %s: %w", bucket, err)
	}

	objects := make([]common.ObjectInfo, 0, len(result.Contents))
	for _, obj := range result.Contents {
		objects = append(objects, common.ObjectInfo{
			Key:          *obj.Key,
			Size:         int64(*obj.Size),
			LastModified: obj.LastModified.Format(time.RFC3339),
			ETag:         *obj.ETag,
			ContentType:  "",
		})
	}

	return objects, nil
}

// GetObject retrieves an object from S3
func (p *AWSProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to AWS")
	}

	result, err := p.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object %s from bucket %s: %w", key, bucket, err)
	}

	return result.Body, nil
}

// PutObject uploads an object to S3
func (p *AWSProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	if !p.connected {
		return fmt.Errorf("not connected to AWS")
	}

	// Convert metadata map to AWS format
	awsMetadata := make(map[string]string)
	for k, v := range metadata {
		awsMetadata[k] = v
	}

	_, err := p.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     data,
		Metadata: awsMetadata,
	})
	if err != nil {
		return fmt.Errorf("failed to put object %s in bucket %s: %w", key, bucket, err)
	}

	return nil
}

// DeleteObject deletes an object from S3
func (p *AWSProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	if !p.connected {
		return fmt.Errorf("not connected to AWS")
	}

	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object %s from bucket %s: %w", key, bucket, err)
	}

	return nil
}

// GetMetadata retrieves metadata for an S3 object
func (p *AWSProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected to AWS")
	}

	result, err := p.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for object %s in bucket %s: %w", key, bucket, err)
	}

	metadata := make(map[string]string)
	for k, v := range result.Metadata {
		metadata[k] = v
	}

	return metadata, nil
}

// GetPresignedURL generates a presigned URL for an S3 object
func (p *AWSProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	if !p.connected {
		return "", fmt.Errorf("not connected to AWS")
	}

	request, err := p.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expiration) * time.Second
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for object %s in bucket %s: %w", key, bucket, err)
	}

	return request.URL, nil
}

// Factory for creating AWS providers
type AWSProviderFactory struct{}

// Create creates a new AWS provider instance
func (f *AWSProviderFactory) Create(config common.CloudConfig) (common.CloudProvider, error) {
	provider := NewAWSProvider()
	
	// Convert CloudConfig to map for the Connect method
	configMap := map[string]interface{}{
		"region":        config.Region,
		"endpoint":      config.EndpointOverride,
		"use_iam_role":  false,
	}
	
	// Extract credentials if provided
	if config.Credentials != nil {
		if accessKey, ok := config.Credentials["access_key"].(string); ok {
			configMap["access_key"] = accessKey
		}
		if secretKey, ok := config.Credentials["secret_key"].(string); ok {
			configMap["secret_key"] = secretKey
		}
		if sessionToken, ok := config.Credentials["session_token"].(string); ok {
			configMap["session_token"] = sessionToken
		}
		if useIAMRole, ok := config.Credentials["use_iam_role"].(bool); ok {
			configMap["use_iam_role"] = useIAMRole
		}
	}
	
	// Add any additional options
	for k, v := range config.AdditionalOptions {
		configMap[k] = v
	}
	
	// Connect to AWS
	err := provider.Connect(context.Background(), configMap)
	if err != nil {
		return nil, err
	}
	
	return provider, nil
}
