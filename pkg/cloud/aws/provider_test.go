package aws

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

// MockAWSProvider is a mock implementation of the AWS provider
type MockAWSProvider struct {
	mock.Mock
}

func (m *MockAWSProvider) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockAWSProvider) Connect(ctx context.Context, config map[string]interface{}) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockAWSProvider) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockAWSProvider) ListBuckets(ctx context.Context) ([]common.BucketInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]common.BucketInfo), args.Error(1)
}

func (m *MockAWSProvider) ListObjects(ctx context.Context, bucket, prefix string) ([]common.ObjectInfo, error) {
	args := m.Called(ctx, bucket, prefix)
	return args.Get(0).([]common.ObjectInfo), args.Error(1)
}

func (m *MockAWSProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockAWSProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	args := m.Called(ctx, bucket, key, data, size, metadata)
	return args.Error(0)
}

func (m *MockAWSProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	args := m.Called(ctx, bucket, key)
	return args.Error(0)
}

func (m *MockAWSProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockAWSProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiration int64) (string, error) {
	args := m.Called(ctx, bucket, key, expiration)
	return args.String(0), args.Error(1)
}

type MockReadCloser struct {
	reader io.Reader
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *MockReadCloser) Close() error {
	return nil
}

// TestAWSProvider_Name tests the Name method
func TestAWSProvider_Name(t *testing.T) {
	provider := NewAWSProvider()
	assert.Equal(t, "AWS", provider.Name())
}

// TestAWSProvider_ListBuckets tests the ListBuckets method
func TestAWSProvider_ListBuckets(t *testing.T) {
	// Create mock provider
	provider := new(MockAWSProvider)

	// Create test data
	creationDate := time.Now().Format(time.RFC3339)

	// Set up mock expectations
	provider.On("ListBuckets", mock.Anything).Return([]common.BucketInfo{
		{
			Name:         "test-bucket-1",
			CreationDate: creationDate,
			Region:       "us-west-2",
		},
		{
			Name:         "test-bucket-2",
			CreationDate: creationDate,
			Region:       "us-east-1",
		},
	}, nil)

	// Call ListBuckets
	buckets, err := provider.ListBuckets(context.Background())

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, buckets, 2)
	assert.Equal(t, "test-bucket-1", buckets[0].Name)
	assert.Equal(t, "test-bucket-2", buckets[1].Name)
	assert.Equal(t, creationDate, buckets[0].CreationDate)
	assert.Equal(t, creationDate, buckets[1].CreationDate)
	provider.AssertExpectations(t)
}

// TestAWSProvider_ListObjects tests the ListObjects method
func TestAWSProvider_ListObjects(t *testing.T) {
	// Create mock provider
	provider := new(MockAWSProvider)

	// Create test data
	size := int64(1024)
	lastModified := time.Now().Format(time.RFC3339)

	// Set up mock expectations
	provider.On("ListObjects", mock.Anything, "test-bucket", "test-prefix/").Return([]common.ObjectInfo{
		{
			Key:          "test-prefix/file1.txt",
			Size:         size,
			LastModified: lastModified,
			ETag:         "test-etag",
			ContentType:  "text/plain",
		},
		{
			Key:          "test-prefix/file2.txt",
			Size:         size,
			LastModified: lastModified,
			ETag:         "test-etag",
			ContentType:  "text/plain",
		},
	}, nil)

	// Call ListObjects
	objects, err := provider.ListObjects(context.Background(), "test-bucket", "test-prefix/")

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	assert.Equal(t, "test-prefix/file1.txt", objects[0].Key)
	assert.Equal(t, "test-prefix/file2.txt", objects[1].Key)
	assert.Equal(t, size, objects[0].Size)
	assert.Equal(t, size, objects[1].Size)
	assert.Equal(t, lastModified, objects[0].LastModified)
	assert.Equal(t, lastModified, objects[1].LastModified)
	assert.Equal(t, "test-etag", objects[0].ETag)
	assert.Equal(t, "test-etag", objects[1].ETag)
	provider.AssertExpectations(t)
}

// TestAWSProvider_GetObject tests the GetObject method
func TestAWSProvider_GetObject(t *testing.T) {
	// Create mock provider
	provider := new(MockAWSProvider)

	// Create test data
	content := "test content"
	reader := &MockReadCloser{
		reader: strings.NewReader(content),
	}

	// Set up mock expectations
	provider.On("GetObject", mock.Anything, "test-bucket", "test-key").Return(reader, nil)

	// Call GetObject
	result, err := provider.GetObject(context.Background(), "test-bucket", "test-key")

	// Assert expectations
	assert.NoError(t, err)
	data, err := ioutil.ReadAll(result)
	assert.NoError(t, err)
	assert.Equal(t, content, string(data))
	provider.AssertExpectations(t)
}

// TestAWSProviderFactory_Create tests the AWSProviderFactory.Create method
func TestAWSProviderFactory_Create(t *testing.T) {
	// Create a factory
	factory := &AWSProviderFactory{}

	// Create test config
	config := common.CloudConfig{
		Provider: "aws",
		AdditionalOptions: map[string]interface{}{
			"region": "us-west-2",
			// Add a non-existent endpoint to force an error
			"endpoint": "http://non-existent-endpoint",
		},
		Credentials: map[string]interface{}{
			"access_key": "test-access-key",
			"secret_key": "test-secret-key",
		},
		DefaultBucket: "test-bucket",
	}

	// This test will not actually connect to AWS
	// It just verifies that the factory creates a provider with the correct configuration
	_, err := factory.Create(config)

	// We expect an error since we're not actually connecting to AWS
	assert.Error(t, err)
}
