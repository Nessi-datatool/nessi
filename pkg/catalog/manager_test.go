package catalog

import (
	"context"
	"testing"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCatalog is a mock implementation of the DataCatalog interface
// MockCatalog implements the DataCatalog interface for testing
type MockCatalog struct {
	mock.Mock
}

func (m *MockCatalog) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockCatalog) Connect(ctx context.Context, config map[string]interface{}) error {
	args := m.Called(ctx, config)
	return args.Error(0)
}

func (m *MockCatalog) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCatalog) ListDatabases(ctx context.Context) ([]types.DatabaseInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).([]types.DatabaseInfo), args.Error(1)
}

func (m *MockCatalog) ListTables(ctx context.Context, database string) ([]types.TableInfo, error) {
	args := m.Called(ctx, database)
	return args.Get(0).([]types.TableInfo), args.Error(1)
}

func (m *MockCatalog) GetTableDetails(ctx context.Context, database, table string) (*types.TableDetails, error) {
	args := m.Called(ctx, database, table)
	return args.Get(0).(*types.TableDetails), args.Error(1)
}

func (m *MockCatalog) GetTableMetadata(ctx context.Context, database, table string) (*types.TableMetadata, error) {
	args := m.Called(ctx, database, table)
	return args.Get(0).(*types.TableMetadata), args.Error(1)
}

func (m *MockCatalog) UpdateTableMetadata(ctx context.Context, database, table string, metadata *types.TableMetadata) error {
	args := m.Called(ctx, database, table, metadata)
	return args.Error(0)
}

func (m *MockCatalog) GetTableLineage(ctx context.Context, database, table string) (*types.LineageInfo, error) {
	args := m.Called(ctx, database, table)
	return args.Get(0).(*types.LineageInfo), args.Error(1)
}

func (m *MockCatalog) UpdateTableLineage(ctx context.Context, database, table string, lineage *types.LineageInfo) error {
	args := m.Called(ctx, database, table, lineage)
	return args.Error(0)
}

func (m *MockCatalog) PublishQualityMetrics(ctx context.Context, database, table string, metrics *types.QualityMetrics) error {
	args := m.Called(ctx, database, table, metrics)
	return args.Error(0)
}

func (m *MockCatalog) GetQualityMetrics(ctx context.Context, database, table string) (*types.QualityMetrics, error) {
	args := m.Called(ctx, database, table)
	return args.Get(0).(*types.QualityMetrics), args.Error(1)
}

// TestCatalogManager_RegisterCatalog tests the RegisterCatalog method
func TestCatalogManager_RegisterCatalog(t *testing.T) {
	// Create mock catalog
	mockCatalog := new(MockCatalog)
	mockCatalog.On("Name").Return("TestCatalog")
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register catalog
	manager.RegisterCatalog(mockCatalog)
	
	// Verify catalog is registered
	catalogs := manager.ListCatalogs()
	assert.Contains(t, catalogs, "TestCatalog")
	
	// Get catalog
	catalog, err := manager.GetCatalog("TestCatalog")
	assert.NoError(t, err)
	assert.Equal(t, mockCatalog, catalog)
	
	// Verify mock expectations
	mockCatalog.AssertExpectations(t)
}

// TestCatalogManager_GetCatalog tests the GetCatalog method
func TestCatalogManager_GetCatalog(t *testing.T) {
	// Create mock catalog
	mockCatalog := new(MockCatalog)
	mockCatalog.On("Name").Return("TestCatalog")
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register catalog
	manager.RegisterCatalog(mockCatalog)
	
	// Get catalog
	catalog, err := manager.GetCatalog("TestCatalog")
	assert.NoError(t, err)
	assert.Equal(t, mockCatalog, catalog)
	
	// Get non-existent catalog
	_, err = manager.GetCatalog("NonExistentCatalog")
	assert.Error(t, err)
	
	// Verify mock expectations
	mockCatalog.AssertExpectations(t)
}

// TestCatalogManager_ListCatalogs tests the ListCatalogs method
func TestCatalogManager_ListCatalogs(t *testing.T) {
	// Create mock catalogs
	mockCatalog1 := new(MockCatalog)
	mockCatalog1.On("Name").Return("TestCatalog1")
	
	mockCatalog2 := new(MockCatalog)
	mockCatalog2.On("Name").Return("TestCatalog2")
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register catalogs
	manager.RegisterCatalog(mockCatalog1)
	manager.RegisterCatalog(mockCatalog2)
	
	// List catalogs
	catalogs := manager.ListCatalogs()
	assert.Len(t, catalogs, 2)
	assert.Contains(t, catalogs, "TestCatalog1")
	assert.Contains(t, catalogs, "TestCatalog2")
	
	// Verify mock expectations
	mockCatalog1.AssertExpectations(t)
	mockCatalog2.AssertExpectations(t)
}

// TestCatalogManager_ConnectCatalog tests the ConnectCatalog method
func TestCatalogManager_ConnectCatalog(t *testing.T) {
	// Create mock catalog
	mockCatalog := new(MockCatalog)
	mockCatalog.On("Name").Return("TestCatalog")
	
	// Create test config
	config := map[string]interface{}{
		"key": "value",
	}
	
	// Set up mock expectations
	ctx := context.Background()
	mockCatalog.On("Connect", ctx, config).Return(nil)
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register catalog
	manager.RegisterCatalog(mockCatalog)
	
	// Connect catalog
	err := manager.ConnectCatalog(ctx, "TestCatalog", config)
	assert.NoError(t, err)
	
	// Verify mock expectations
	mockCatalog.AssertExpectations(t)
}

// TestCatalogManager_DisconnectCatalog tests the DisconnectCatalog method
func TestCatalogManager_DisconnectCatalog(t *testing.T) {
	// Create mock catalog
	mockCatalog := new(MockCatalog)
	mockCatalog.On("Name").Return("TestCatalog")
	
	// Set up mock expectations
	ctx := context.Background()
	mockCatalog.On("Disconnect", ctx).Return(nil)
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register catalog
	manager.RegisterCatalog(mockCatalog)
	
	// Disconnect catalog
	err := manager.DisconnectCatalog(ctx, "TestCatalog")
	assert.NoError(t, err)
	
	// Verify mock expectations
	mockCatalog.AssertExpectations(t)
}

// TestCatalogManager_DisconnectAll tests the DisconnectAll method
func TestCatalogManager_DisconnectAll(t *testing.T) {
	// Create mock catalogs
	mockCatalog1 := new(MockCatalog)
	mockCatalog1.On("Name").Return("TestCatalog1")
	
	mockCatalog2 := new(MockCatalog)
	mockCatalog2.On("Name").Return("TestCatalog2")
	
	// Set up mock expectations
	ctx := context.Background()
	mockCatalog1.On("Disconnect", ctx).Return(nil)
	mockCatalog2.On("Disconnect", ctx).Return(nil)
	
	// Create catalog manager
	manager := NewCatalogManager()
	
	// Register catalogs
	manager.RegisterCatalog(mockCatalog1)
	manager.RegisterCatalog(mockCatalog2)
	
	// Disconnect all catalogs
	manager.DisconnectAll(ctx)
	
	// Verify mock expectations
	mockCatalog1.AssertExpectations(t)
	mockCatalog2.AssertExpectations(t)
}
