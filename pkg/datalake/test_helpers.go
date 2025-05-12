package datalake

// MockMetadataManager is a mock implementation of MetadataManager for testing
type MockMetadataManager struct {
	*MetadataManager
	MockReadParquetFile        func(file string) ([]map[string]interface{}, error)
	MockReadTableMetadata      func() (*DeltaTable, error)
	MockGetSchemaFieldsAtVersion func(version int64) ([]SchemaField, error)
	MockGetVersionHistory     func() ([]VersionEntry, error)
}

// ReadParquetFile overrides the MetadataManager's ReadParquetFile method for testing
func (mm *MockMetadataManager) ReadParquetFile(file string) ([]map[string]interface{}, error) {
	if mm.MockReadParquetFile != nil {
		return mm.MockReadParquetFile(file)
	}
	return mm.MetadataManager.ReadParquetFile(file)
}

// ReadTableMetadata overrides the MetadataManager's ReadTableMetadata method for testing
func (mm *MockMetadataManager) ReadTableMetadata() (*DeltaTable, error) {
	if mm.MockReadTableMetadata != nil {
		return mm.MockReadTableMetadata()
	}
	return mm.MetadataManager.ReadTableMetadata()
}

// GetSchemaFieldsAtVersion overrides the MetadataManager's GetSchemaFieldsAtVersion method for testing
func (mm *MockMetadataManager) GetSchemaFieldsAtVersion(version int64) ([]SchemaField, error) {
	if mm.MockGetSchemaFieldsAtVersion != nil {
		return mm.MockGetSchemaFieldsAtVersion(version)
	}
	return mm.MetadataManager.GetSchemaFieldsAtVersion(version)
}

// GetVersionHistory overrides the MetadataManager's GetVersionHistory method for testing
func (mm *MockMetadataManager) GetVersionHistory() ([]VersionEntry, error) {
	if mm.MockGetVersionHistory != nil {
		return mm.MockGetVersionHistory()
	}
	return mm.MetadataManager.GetVersionHistory()
}

// NewMockMetadataManager creates a new MockMetadataManager
func NewMockMetadataManager(tablePath string) *MockMetadataManager {
	return &MockMetadataManager{
		MetadataManager: NewMetadataManager(tablePath),
	}
}
