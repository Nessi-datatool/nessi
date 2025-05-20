package freshness

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// MockDeltaConnector is a mock implementation of the Delta connector interface
type MockDeltaConnector struct {
	mock.Mock
}

// GetTableMetadata mocks the GetTableMetadata method
func (m *MockDeltaConnector) GetTableMetadata(tablePath string) (*TableMetadata, error) {
	args := m.Called(tablePath)
	return args.Get(0).(*TableMetadata), args.Error(1)
}

// ListTables mocks the ListTables method
func (m *MockDeltaConnector) ListTables() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

// TableMetadata represents metadata for a Delta table
type TableMetadata struct {
	// LastModified is the timestamp when the table was last modified
	LastModified *time.Time

	// Version is the version of the table
	Version int64

	// Name is the name of the table
	Name string

	// Path is the path to the table
	Path string
}
