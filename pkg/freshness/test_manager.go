package freshness

import (
	"fmt"
	"sync"
	"time"
)

// TestSLAManager is a simplified version of SLAManager for testing purposes
type TestSLAManager struct {
	configs        map[string]*SLAConfig
	deltaConnector interface {
		GetTableMetadata(tablePath string) (*TableMetadata, error)
		ListTables() ([]string, error)
	}
	historyMgr *historyManager
	mutex      sync.RWMutex
	enabled    bool
}

// NewTestSLAManager creates a new test SLA manager
func NewTestSLAManager(connector interface {
	GetTableMetadata(tablePath string) (*TableMetadata, error)
	ListTables() ([]string, error)
}) *TestSLAManager {
	return &TestSLAManager{
		configs:        make(map[string]*SLAConfig),
		deltaConnector: connector,
		historyMgr:     newHistoryManager(100),
		enabled:        true,
	}
}

// AddSLA adds an SLA configuration
func (m *TestSLAManager) AddSLA(config *SLAConfig) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.configs[config.TableName] = config
}

// GetSLA gets an SLA configuration
func (m *TestSLAManager) GetSLA(tableName string) (*SLAConfig, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	config, ok := m.configs[tableName]
	return config, ok
}

// CheckFreshness checks the freshness of a table
func (m *TestSLAManager) CheckFreshness(tableName string) (*FreshnessStatus, error) {
	// Get SLA configuration
	config, ok := m.GetSLA(tableName)
	if !ok {
		return nil, fmt.Errorf("no SLA configuration found for table %s", tableName)
	}
	
	// Get table metadata using the mock connector
	table, err := m.deltaConnector.GetTableMetadata(config.TablePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read table metadata: %w", err)
	}
	
	// Calculate time since last update
	now := time.Now()
	timeSinceUpdate := now.Sub(*table.LastModified)
	
	// Calculate next expected update
	nextExpectedUpdate := table.LastModified.Add(config.ExpectedFrequency)
	
	// Determine status
	var status SLALevel
	warningThreshold := time.Duration(config.WarningThreshold) * config.ExpectedFrequency / 100
	criticalThreshold := time.Duration(config.CriticalThreshold) * config.ExpectedFrequency / 100
	
	if timeSinceUpdate >= criticalThreshold {
		status = SLALevelCritical
	} else if timeSinceUpdate >= warningThreshold {
		status = SLALevelWarning
	} else {
		status = SLALevelInfo
	}
	
	// Create freshness status
	freshnessStatus := &FreshnessStatus{
		TableName:         config.TableName,
		TablePath:         config.TablePath,
		LastUpdateTime:    *table.LastModified,
		TimeSinceUpdate:   timeSinceUpdate,
		ExpectedFrequency: config.ExpectedFrequency,
		Status:            status,
		NextExpectedUpdate: nextExpectedUpdate,
		SLAConfig:         config,
	}
	
	// Update history
	m.updateHistory(tableName, freshnessStatus)
	
	return freshnessStatus, nil
}

// CheckAllFreshness checks the freshness of all tables
func (m *TestSLAManager) CheckAllFreshness() ([]*FreshnessStatus, error) {
	m.mutex.RLock()
	configs := make([]*SLAConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	m.mutex.RUnlock()
	
	statuses := make([]*FreshnessStatus, 0, len(configs))
	for _, config := range configs {
		status, err := m.CheckFreshness(config.TableName)
		if err != nil {
			// Log error but continue with other tables
			continue
		}
		statuses = append(statuses, status)
	}
	
	return statuses, nil
}

// updateHistory updates the freshness history
func (m *TestSLAManager) updateHistory(tableName string, status *FreshnessStatus) {
	entry := FreshnessHistoryEntry{
		TableName:         status.TableName,
		TablePath:         status.TablePath,
		Timestamp:         time.Now(),
		LastUpdateTime:    status.LastUpdateTime,
		TimeSinceUpdate:   status.TimeSinceUpdate,
		ExpectedFrequency: status.ExpectedFrequency,
		Status:            string(status.Status),
	}
	
	m.historyMgr.addHistoryEntry(entry)
}

// GetTableTrends gets the freshness trends for a table
func (m *TestSLAManager) GetTableTrends(tableName string) (*FreshnessTrends, error) {
	// Check if SLA exists
	_, ok := m.GetSLA(tableName)
	if !ok {
		return nil, fmt.Errorf("no SLA configuration found for table %s", tableName)
	}
	
	// Get history entries
	entries := m.historyMgr.getTableHistory(tableName)
	
	// Calculate compliance
	var infoCount, warningCount, criticalCount int
	for _, entry := range entries {
		switch entry.Status {
		case string(SLALevelInfo):
			infoCount++
		case string(SLALevelWarning):
			warningCount++
		case string(SLALevelCritical):
			criticalCount++
		}
	}
	
	totalCount := infoCount + warningCount + criticalCount
	var complianceRate float64
	if totalCount > 0 {
		complianceRate = float64(infoCount) / float64(totalCount) * 100
	}
	
	compliance := FreshnessCompliance{
		InfoCount:      infoCount,
		WarningCount:   warningCount,
		CriticalCount:  criticalCount,
		TotalCount:     totalCount,
		ComplianceRate: complianceRate,
	}
	
	trends := &FreshnessTrends{
		History:     entries,
		Compliance:  compliance,
	}
	
	return trends, nil
}

// GetAllTablesTrends gets the freshness trends for all tables
func (m *TestSLAManager) GetAllTablesTrends() (*FreshnessTrends, error) {
	m.mutex.RLock()
	configs := make([]*SLAConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}
	m.mutex.RUnlock()
	
	var infoCount, warningCount, criticalCount, totalCount int
	
	for _, config := range configs {
		trends, err := m.GetTableTrends(config.TableName)
		if err != nil {
			continue
		}
		
		infoCount += trends.Compliance.InfoCount
		warningCount += trends.Compliance.WarningCount
		criticalCount += trends.Compliance.CriticalCount
		totalCount += trends.Compliance.TotalCount
	}
	
	var complianceRate float64
	if totalCount > 0 {
		complianceRate = float64(infoCount) / float64(totalCount) * 100
	}
	
	compliance := FreshnessCompliance{
		InfoCount:      infoCount,
		WarningCount:   warningCount,
		CriticalCount:  criticalCount,
		TotalCount:     totalCount,
		ComplianceRate: complianceRate,
	}
	
	// Get all history entries
	history := m.historyMgr.getAllHistory()
	
	return &FreshnessTrends{
		History:    history,
		Compliance: compliance,
	}, nil
}
