package gcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCPDataCatalog_Name(t *testing.T) {
	dc := &DataCatalogClient{}
	assert.Equal(t, "GCP Data Catalog", dc.Name())
}

func TestGCPDataCatalog_ListDatabases(t *testing.T) {
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the ListDatabases method when not connected
	databases, err := dc.ListDatabases(context.Background())
	
	// Since we're not connected, we expect an error
	assert.Error(t, err)
	assert.Nil(t, databases)
	assert.Contains(t, err.Error(), "not connected")
}

func TestGCPDataCatalog_ListTables(t *testing.T) {
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the ListTables method when not connected
	tables, err := dc.ListTables(context.Background(), "test-database")
	
	// Since we're not connected, we expect an error
	assert.Error(t, err)
	assert.Nil(t, tables)
	assert.Contains(t, err.Error(), "not connected")
}

func TestGCPDataCatalog_GetTableDetails(t *testing.T) {
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the GetTableDetails method when not connected
	details, err := dc.GetTableDetails(context.Background(), "test-database", "test-table")
	
	// Since we're not connected, we expect an error
	assert.Error(t, err)
	assert.Nil(t, details)
	assert.Contains(t, err.Error(), "not connected")
}

func TestGCPDataCatalog_Connect(t *testing.T) {
	// Skip this test in CI environment
	t.Skip("Skipping test that requires GCP credentials")
	
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the Connect method with invalid config
	err := dc.Connect(context.Background(), map[string]interface{}{})
	
	// We expect an error due to missing project_id
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project_id is required")
}

func TestGCPDataCatalog_Disconnect(t *testing.T) {
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the Disconnect method
	err := dc.Disconnect(context.Background())
	
	// Since we're not connected, we expect no error
	assert.NoError(t, err)
	assert.False(t, dc.connected)
}

func TestGCPDataCatalog_PublishQualityMetrics(t *testing.T) {
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the PublishQualityMetrics method when not connected
	err := dc.PublishQualityMetrics(context.Background(), "test-database", "test-table", nil)
	
	// Since we're not connected, we expect an error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not connected")
}

func TestGCPDataCatalog_GetQualityMetrics(t *testing.T) {
	// Create a mock DataCatalogClient
	dc := &DataCatalogClient{}
	
	// Test the GetQualityMetrics method when not connected
	metrics, err := dc.GetQualityMetrics(context.Background(), "test-database", "test-table")
	
	// Since we're not connected, we expect an error
	assert.Error(t, err)
	assert.Nil(t, metrics)
	assert.Contains(t, err.Error(), "not connected")
}




