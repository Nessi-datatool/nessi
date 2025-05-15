package freshness

import (
	"fmt"
	"sync"
	"time"
)

// BetterMockDeltaConnector is an improved mock implementation of the Delta connector interface
// that doesn't rely on file operations and can be configured with predefined responses
type BetterMockDeltaConnector struct {
	tableMetadata map[string]*TableMetadata
	tables        []string
	mutex         sync.RWMutex
}

// NewBetterMockDeltaConnector creates a new instance of BetterMockDeltaConnector
func NewBetterMockDeltaConnector() *BetterMockDeltaConnector {
	return &BetterMockDeltaConnector{
		tableMetadata: make(map[string]*TableMetadata),
		tables:        []string{},
		mutex:         sync.RWMutex{},
	}
}

// GetTableMetadata returns the metadata for a table
func (m *BetterMockDeltaConnector) GetTableMetadata(tablePath string) (*TableMetadata, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	metadata, ok := m.tableMetadata[tablePath]
	if !ok {
		return nil, fmt.Errorf("table not found: %s", tablePath)
	}
	return metadata, nil
}

// ListTables returns a list of all tables
func (m *BetterMockDeltaConnector) ListTables() ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	return m.tables, nil
}

// SetTableLastModified sets the last modified time for a table
func (m *BetterMockDeltaConnector) SetTableLastModified(tablePath string, lastModified time.Time) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if meta, ok := m.tableMetadata[tablePath]; ok {
		meta.LastModified = &lastModified
	} else {
		m.tableMetadata[tablePath] = &TableMetadata{
			LastModified: &lastModified,
			Name:         getTableNameFromPath(tablePath),
			Path:         tablePath,
		}
		m.tables = append(m.tables, tablePath)
	}
}

// AddTable adds a table to the mock connector
func (m *BetterMockDeltaConnector) AddTable(tablePath string, lastModified time.Time) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	if !contains(m.tables, tablePath) {
		m.tables = append(m.tables, tablePath)
	}
	
	m.tableMetadata[tablePath] = &TableMetadata{
		LastModified: &lastModified,
		Version:      1,
		Name:         getTableNameFromPath(tablePath),
		Path:         tablePath,
	}
}

// UpdateTableLastModified updates the last modified time for a table
func (m *BetterMockDeltaConnector) UpdateTableLastModified(tablePath string, lastModified time.Time) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	metadata, ok := m.tableMetadata[tablePath]
	if !ok {
		return fmt.Errorf("table not found: %s", tablePath)
	}
	
	metadata.LastModified = &lastModified
	return nil
}

// SetTableMetadata directly sets the metadata for a table
func (m *BetterMockDeltaConnector) SetTableMetadata(tablePath string, metadata *TableMetadata) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.tableMetadata[tablePath] = metadata
	
	// Add to tables list if not already present
	if !contains(m.tables, tablePath) {
		m.tables = append(m.tables, tablePath)
	}
}

// SetTablesList sets the list of tables
func (m *BetterMockDeltaConnector) SetTablesList(tables []string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.tables = tables
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to extract table name from path
func getTableNameFromPath(path string) string {
	// In a real implementation, this would parse the path to extract the table name
	// For simplicity, we'll just return the path as the name
	return path
}
