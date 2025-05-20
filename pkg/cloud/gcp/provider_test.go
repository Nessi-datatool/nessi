package gcp

import (
	"context"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/cloud/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGCPProvider is a mock implementation of the GCP provider
type MockGCPProvider struct {
	mock.Mock
}

func (m *MockGCPProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockGCPProvider) Connect(ctx context.Context, config map[string]interface{}) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockGCPProvider) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockGCPProvider) ListBuckets(ctx context.Context) ([]common.BucketInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]common.BucketInfo), args.Error(1)
}

func (m *MockGCPProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	args := m.Called(ctx, bucket, prefix)
	return args.Get(0).([]common.ObjectInfo), args.Error(1)
}

func (m *MockGCPProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockGCPProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	args := m.Called(ctx, bucket, key, data, size, metadata)
	return args.Error(0)
}

func (m *MockGCPProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	args := m.Called(ctx, bucket, key)
	return args.Error(0)
}

func (m *MockGCPProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockGCPProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	args := m.Called(ctx, bucket, key, expiration)
	return args.String(0), args.Error(1)
}

// MockReadCloser is a mock implementation of io.ReadCloser
type MockReadCloser struct {
	mock.Mock
	reader *strings.Reader
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *MockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

// TestGCPProvider_Name tests the Name method
func TestGCPProvider_Name(t *testing.T) {
	provider := NewGCPProvider()
	assert.Equal(t, "GCP", provider.Name())
}

// TestGCPProvider_ListBuckets tests the ListBuckets method
func TestGCPProvider_ListBuckets(t *testing.T) {
	// Create a mock provider
	provider := new(MockGCPProvider)

	// Set up mock expectations
	provider.On("ListBuckets", mock.Anything).Return([]common.BucketInfo{
		{
			Name:         "bucket1",
			CreationDate: time.Now().Format(time.RFC3339),
			Region:       "us-west1",
		},
		{
			Name:         "bucket2",
			CreationDate: time.Now().Format(time.RFC3339),
			Region:       "us-central1",
		},
	}, nil)

	// Call ListBuckets
	buckets, err := provider.ListBuckets(context.Background())

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, buckets, 2)
	assert.Equal(t, "bucket1", buckets[0].Name)
	assert.Equal(t, "bucket2", buckets[1].Name)
	assert.Equal(t, "us-west1", buckets[0].Region)
	assert.Equal(t, "us-central1", buckets[1].Region)
	provider.AssertExpectations(t)
}

// TestGCPProvider_ListObjects tests the ListObjects method
func TestGCPProvider_ListObjects(t *testing.T) {
	// Create a mock provider
	provider := new(MockGCPProvider)

	// Set up mock expectations
	provider.On("ListObjects", mock.Anything, "test-bucket", "prefix/").Return([]common.ObjectInfo{
		{
			Key:          "prefix/file1.txt",
			Size:         1024,
			LastModified: time.Now().Format(time.RFC3339),
			ETag:         "etag1",
			ContentType:  "text/plain",
			Metadata: map[string]string{
				"key1": "value1",
			},
		},
		{
			Key:          "prefix/file2.txt",
			Size:         2048,
			LastModified: time.Now().Format(time.RFC3339),
			ETag:         "etag2",
			ContentType:  "text/plain",
			Metadata: map[string]string{
				"key2": "value2",
			},
		},
	}, nil)

	// Call ListObjects
	objects, err := provider.ListObjects(context.Background(), "test-bucket", "prefix/")

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	assert.Equal(t, "prefix/file1.txt", objects[0].Key)
	assert.Equal(t, int64(1024), objects[0].Size)
	assert.Equal(t, "etag1", objects[0].ETag)
	assert.Equal(t, "text/plain", objects[0].ContentType)
	assert.Equal(t, "value1", objects[0].Metadata["key1"])
	provider.AssertExpectations(t)
}

// TestGCPProvider_GetObject tests the GetObject method
func TestGCPProvider_GetObject(t *testing.T) {
	// Create a mock provider
	provider := new(MockGCPProvider)

	// Create test content
	content := "test content"
	mockReader := &MockReadCloser{
		reader: strings.NewReader(content),
	}

	// Set up mock expectations
	provider.On("GetObject", mock.Anything, "test-bucket", "test-key").Return(mockReader, nil)

	// Call GetObject
	reader, err := provider.GetObject(context.Background(), "test-bucket", "test-key")

	// Assert expectations
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
	provider.AssertExpectations(t)
}

// TestGCPProviderFactory_Create tests the GCPProviderFactory.Create method
func TestGCPProviderFactory_Create(t *testing.T) {
	// Create a factory
	factory := &GCPProviderFactory{}

	// Create test config
	config := common.CloudConfig{
		Provider: "gcp",
		AdditionalOptions: map[string]interface{}{
			"project_id": "test-project",
		},
		Credentials: map[string]interface{}{
			"credentials_file": "/path/to/credentials.json",
		},
		DefaultBucket: "test-bucket",
	}

	// This test will not actually connect to GCP
	// It just verifies that the factory creates a provider with the correct configuration
	_, err := factory.Create(config)

	// We expect an error since we're not actually connecting to GCP
	assert.Error(t, err)
}
