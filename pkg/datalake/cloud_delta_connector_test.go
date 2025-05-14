package datalake

import (
	"context"
	"io"
	"strings"
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

func (m *MockCloudProvider) GetMetadata(ctx context.Context, bucket, key string) (map[string]string, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockCloudProvider) GetPresignedURL(ctx context.Context, bucket, key string, expiry int64) (string, error) {
	args := m.Called(ctx, bucket, key, expiry)
	return args.String(0), args.Error(1)
}

func (m *MockCloudProvider) CreateBucket(ctx context.Context, bucket string) error {
	args := m.Called(ctx, bucket)
	return args.Error(0)
}

func (m *MockCloudProvider) DeleteBucket(ctx context.Context, bucket string) error {
	args := m.Called(ctx, bucket)
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

func (m *MockCloudProvider) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockCloudProvider) PutObject(ctx context.Context, bucket, key string, data io.Reader, size int64, metadata map[string]string) error {
	args := m.Called(ctx, bucket, key, data, size, metadata)
	return args.Error(0)
}

func (m *MockCloudProvider) DeleteObject(ctx context.Context, bucket, key string) error {
	args := m.Called(ctx, bucket, key)
	return args.Error(0)
}

// MockReadCloser is a mock implementation of io.ReadCloser
type MockReadCloser struct {
	mock.Mock
	reader io.Reader
}

func (m *MockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m *MockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestCloudDeltaConnector_GetTableInfo(t *testing.T) {
	// Create mock provider
	mockProvider := new(MockCloudProvider)
	bucket := "test-bucket"
	tablePath := "data/test_table"

	// Set up mock expectations for Connect, Disconnect, GetMetadata, and Name
	mockProvider.On("Name").Return("mock-provider")
	mockProvider.On("Connect", mock.Anything, mock.Anything).Return(nil)
	mockProvider.On("Disconnect", mock.Anything).Return(nil)
	mockProvider.On("GetMetadata", mock.Anything, mock.Anything, mock.Anything).Return(map[string]string{"provider": "mock"}, nil)
	mockProvider.On("GetPresignedURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("https://example.com/presigned", nil)

	// Create mock checkpoint
	checkpointContent := `{"version":2,"size":3,"sizeInBytes":123456,"numOfAddFiles":3,"numOfRemoveFiles":1,"path":"data/test_table/_delta_log/00000000000000000002.json"}`
	checkpointReader := &MockReadCloser{
		reader: strings.NewReader(checkpointContent),
	}

	// Create mock log files - fixing the JSON format for add and remove files
	logFile0Content := `{"metaData":{"id":"12345","format":{"provider":"parquet"},"schemaString":"struct<id:int,name:string>","partitionColumns":[]}}
{"add":[{"path":"part-00000.parquet","size":1024,"modificationTime":1609459200000,"dataChange":true}]}
{"add":[{"path":"part-00001.parquet","size":2048,"modificationTime":1609459200000,"dataChange":true}]}`

	logFile1Content := `{"remove":[{"path":"part-00000.parquet","dataChange":true}]}
{"add":[{"path":"part-00002.parquet","size":4096,"modificationTime":1609545600000,"dataChange":true}]}`

	logFile2Content := `{"commitInfo":{"timestamp":1609632000000,"operation":"UPDATE"}}
{"add":[{"path":"part-00003.parquet","size":8192,"modificationTime":1609632000000,"dataChange":true}]}`

	logFile0Reader := &MockReadCloser{
		reader: strings.NewReader(logFile0Content),
	}

	logFile1Reader := &MockReadCloser{
		reader: strings.NewReader(logFile1Content),
	}

	logFile2Reader := &MockReadCloser{
		reader: strings.NewReader(logFile2Content),
	}

	// Set up mock expectations for checkpoint
	mockProvider.On("ListObjects", mock.Anything, bucket, tablePath+"/_delta_log/_last_checkpoint").
		Return([]common.ObjectInfo{{Key: tablePath + "/_delta_log/_last_checkpoint"}}, nil)

	mockProvider.On("GetObject", mock.Anything, bucket, tablePath+"/_delta_log/_last_checkpoint").
		Return(checkpointReader, nil)

	// Set up mock expectations for log files
	mockProvider.On("ListObjects", mock.Anything, bucket, tablePath+"/_delta_log/00000000000000000000.json").
		Return([]common.ObjectInfo{{Key: tablePath + "/_delta_log/00000000000000000000.json"}}, nil)

	mockProvider.On("ListObjects", mock.Anything, bucket, tablePath+"/_delta_log/00000000000000000001.json").
		Return([]common.ObjectInfo{{Key: tablePath + "/_delta_log/00000000000000000001.json"}}, nil)

	mockProvider.On("ListObjects", mock.Anything, bucket, tablePath+"/_delta_log/00000000000000000002.json").
		Return([]common.ObjectInfo{{Key: tablePath + "/_delta_log/00000000000000000002.json"}}, nil)

	mockProvider.On("GetObject", mock.Anything, bucket, tablePath+"/_delta_log/00000000000000000000.json").
		Return(logFile0Reader, nil)

	mockProvider.On("GetObject", mock.Anything, bucket, tablePath+"/_delta_log/00000000000000000001.json").
		Return(logFile1Reader, nil)

	mockProvider.On("GetObject", mock.Anything, bucket, tablePath+"/_delta_log/00000000000000000002.json").
		Return(logFile2Reader, nil)

	// Set up mock expectations for close
	checkpointReader.On("Close").Return(nil)
	logFile0Reader.On("Close").Return(nil)
	logFile1Reader.On("Close").Return(nil)
	logFile2Reader.On("Close").Return(nil)

	// Create connector
	connector := &CloudDeltaConnector{
		provider: mockProvider,
		bucket:   bucket,
	}

	// Call GetTableInfo
	tableInfo, err := connector.GetTableInfo(context.Background(), tablePath)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, tableInfo)
	assert.Equal(t, tablePath, tableInfo.Path)
	assert.Equal(t, int64(2), tableInfo.Version)
	assert.NotNil(t, tableInfo.Metadata)
	assert.Equal(t, "12345", tableInfo.Metadata.ID)
}
