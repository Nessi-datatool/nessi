package datalake

import (
	"fmt"
	"time"

	"github.com/apache/arrow/go/v15/arrow"
	"github.com/stretchr/testify/mock"
)

// MockDeltaConnector is a mock implementation of the DeltaConnector interface for testing
type MockDeltaConnector struct {
	mock.Mock
}

// IsDeltaTable mocks the implementation of IsDeltaTable
func (m *MockDeltaConnector) IsDeltaTable(path string) bool {
	args := m.Called(path)
	return args.Bool(0)
}

// Read mocks the implementation of Read
func (m *MockDeltaConnector) Read(path string) ([]map[string]interface{}, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// ReadWithInference mocks the implementation of ReadWithInference
func (m *MockDeltaConnector) ReadWithInference(path string) ([]map[string]interface{}, *arrow.Schema, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	var schema *arrow.Schema
	if args.Get(1) != nil {
		schema = args.Get(1).(*arrow.Schema)
	}
	return args.Get(0).([]map[string]interface{}), schema, args.Error(2)
}

// Write mocks the implementation of Write
func (m *MockDeltaConnector) Write(path string, data []map[string]interface{}, schema *arrow.Schema) error {
	args := m.Called(path, data, schema)
	return args.Error(0)
}

// GetDeltaTableMetadata mocks the implementation of GetDeltaTableMetadata
func (m *MockDeltaConnector) GetDeltaTableMetadata(path string) (*DeltaTableMetadata, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*DeltaTableMetadata), args.Error(1)
}

// ReadAsOfVersion mocks the implementation of ReadAsOfVersion
func (m *MockDeltaConnector) ReadAsOfVersion(path string, version int) ([]map[string]interface{}, error) {
	args := m.Called(path, version)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// ReadAsOfTimestamp mocks the implementation of ReadAsOfTimestamp
func (m *MockDeltaConnector) ReadAsOfTimestamp(path string, timestamp time.Time) ([]map[string]interface{}, error) {
	args := m.Called(path, timestamp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// GetVersionHistory mocks the implementation of GetVersionHistory
func (m *MockDeltaConnector) GetVersionHistory(path string) ([]DeltaVersion, error) {
	args := m.Called(path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]DeltaVersion), args.Error(1)
}

// NewMockDeltaConnector creates a new mock Delta connector with default expectations
func NewMockDeltaConnector() *MockDeltaConnector {
	mockConn := &MockDeltaConnector{}

	// Set up default behaviors
	mockConn.On("IsDeltaTable", mock.AnythingOfType("string")).Return(true)

	// Default data
	defaultData := []map[string]interface{}{
		{"id": 1, "name": "Test Data", "value": 100.0},
		{"id": 2, "name": "More Test Data", "value": 200.0},
	}
	mockConn.On("Read", mock.AnythingOfType("string")).Return(defaultData, nil)

	// Default metadata
	defaultMetadata := &DeltaTableMetadata{
		Format:           "delta",
		ID:               "test-table-id",
		Name:             "test_table",
		Description:      "Test table for unit testing",
		Location:         "/path/to/test_table",
		SchemaString:     "{\"fields\":[{\"name\":\"id\",\"type\":\"integer\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"double\"}]}",
		PartitionColumns: []string{},
		CreatedTime:      time.Now().Add(-24 * time.Hour),
		LastModified:     time.Now(),
	}
	mockConn.On("GetDeltaTableMetadata", mock.AnythingOfType("string")).Return(defaultMetadata, nil)

	return mockConn
}

// WithError configures the mock to return an error for the specified method
func (m *MockDeltaConnector) WithError(method string, err error) *MockDeltaConnector {
	switch method {
	case "Read":
		m.On("Read", mock.AnythingOfType("string")).Return(nil, err)
	case "Write":
		m.On("Write", mock.AnythingOfType("string"), mock.AnythingOfType("[]map[string]interface {}"), mock.AnythingOfType("*arrow.Schema")).Return(err)
	case "GetDeltaTableMetadata":
		m.On("GetDeltaTableMetadata", mock.AnythingOfType("string")).Return(nil, err)
	case "ReadAsOfVersion":
		m.On("ReadAsOfVersion", mock.AnythingOfType("string"), mock.AnythingOfType("int")).Return(nil, err)
	case "ReadAsOfTimestamp":
		m.On("ReadAsOfTimestamp", mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil, err)
	case "GetVersionHistory":
		m.On("GetVersionHistory", mock.AnythingOfType("string")).Return(nil, err)
	default:
		panic(fmt.Sprintf("Unknown method: %s", method))
	}
	return m
}
