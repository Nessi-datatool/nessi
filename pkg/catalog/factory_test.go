package catalog_test

import (
	"context"
	"testing"

	nessitypes "github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

// Define mock catalog
type mockCatalog struct{}

func (m *mockCatalog) Connect(ctx context.Context, config map[string]interface{}) error { return nil }
func (m *mockCatalog) Disconnect(ctx context.Context) error { return nil }
func (m *mockCatalog) ListDatabases(ctx context.Context) ([]nessitypes.DatabaseInfo, error) { return nil, nil }
func (m *mockCatalog) GetDatabase(ctx context.Context, name string) (*nessitypes.DatabaseInfo, error) { return nil, nil }
func (m *mockCatalog) ListTables(ctx context.Context, databaseName string) ([]nessitypes.TableInfo, error) { return nil, nil }
func (m *mockCatalog) GetTable(ctx context.Context, databaseName, tableName string) (*nessitypes.TableInfo, error) { return nil, nil }
func (m *mockCatalog) GetTableSchema(ctx context.Context, databaseName, tableName string) (*nessitypes.TableSchema, error) { return nil, nil }
func (m *mockCatalog) GetTableMetadata(ctx context.Context, databaseName, tableName string) (*nessitypes.TableMetadata, error) { return nil, nil }
func (m *mockCatalog) GetTableLineage(ctx context.Context, databaseName, tableName string) (*nessitypes.LineageInfo, error) { return nil, nil }
func (m *mockCatalog) PublishQualityMetrics(ctx context.Context, databaseName, tableName string, metrics *nessitypes.QualityMetrics) error { return nil }
func (m *mockCatalog) GetQualityMetrics(ctx context.Context, databaseName, tableName string) (*nessitypes.QualityMetrics, error) { return nil, nil }
func (m *mockCatalog) GetTableDetails(ctx context.Context, databaseName, tableName string) (*nessitypes.TableDetails, error) { return nil, nil }
func (m *mockCatalog) Name() string { return "mock" }
func (m *mockCatalog) UpdateTableLineage(ctx context.Context, databaseName, tableName string, lineage *nessitypes.LineageInfo) error { return nil }
func (m *mockCatalog) UpdateTableMetadata(ctx context.Context, databaseName, tableName string, metadata *nessitypes.TableMetadata) error { return nil }

func TestCatalogFactory(t *testing.T) {
	factory := nessitypes.GetCatalogFactory()
	assert.NotNil(t, factory)

	// Test registering a mock provider
	mockProvider := func() nessitypes.DataCatalog {
		return &mockCatalog{}
	}

	// Register mock provider
	factory.RegisterProvider(nessitypes.CatalogType("mock"), mockProvider)

	// Test creating a catalog
	catalog, err := factory.CreateCatalog(nessitypes.CatalogType("mock"))
	assert.NoError(t, err)
	assert.NotNil(t, catalog)
}
