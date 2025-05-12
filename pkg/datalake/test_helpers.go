package datalake

// MockMetadataManager is a wrapper around MetadataManager for testing
type MockMetadataManager struct {
	*MetadataManager
	MockReadParquetFile func(file string) ([]map[string]interface{}, error)
}

// ReadParquetFile overrides the MetadataManager's ReadParquetFile method for testing
func (mm *MockMetadataManager) ReadParquetFile(file string) ([]map[string]interface{}, error) {
	if mm.MockReadParquetFile != nil {
		return mm.MockReadParquetFile(file)
	}
	return mm.MetadataManager.ReadParquetFile(file)
}

// NewMockMetadataManager creates a new MockMetadataManager
func NewMockMetadataManager(tablePath string) *MockMetadataManager {
	return &MockMetadataManager{
		MetadataManager: NewMetadataManager(tablePath),
	}
}
