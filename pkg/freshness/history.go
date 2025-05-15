package freshness

import (
	"sync"
	"time"
)

// updateHistory updates the history for a table's freshness status
func (m *SLAManager) updateHistory(tableName string, status *FreshnessStatus) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	lastCheckTime, exists := m.lastCheck[tableName]

	// Add history entry if:
	// 1. This is the first check for this table, or
	// 2. At least 1 hour has passed since the last check, or
	// 3. The status has changed since the last check
	if !exists || now.Sub(lastCheckTime) >= time.Hour {
		// Create history entry
		historyEntry := FreshnessHistoryEntry{
			TableName:         status.TableName,
			TablePath:         status.TablePath,
			Timestamp:         now,
			LastUpdateTime:    status.LastUpdateTime,
			TimeSinceUpdate:   status.TimeSinceUpdate,
			ExpectedFrequency: status.ExpectedFrequency,
			Status:            string(status.Status),
		}

		// Add to history
		m.historyMgr.addHistoryEntry(historyEntry)

		// Update last check time
		m.lastCheck[tableName] = now
	}
}

// GetTableTrends gets the freshness trends for a specific table
func (m *SLAManager) GetTableTrends(tableName string) (*FreshnessTrends, error) {
	// Check if SLA exists for the table
	_, err := m.GetSLA(tableName)
	if err != nil {
		return nil, err
	}

	// Get the table status
	status, err := m.CheckFreshness(tableName)
	if err != nil {
		return nil, err
	}

	// Get the table history
	history := m.historyMgr.getTableHistory(tableName)

	// Calculate compliance
	compliance := FreshnessCompliance{
		TotalCount: 1,
	}

	switch status.Status {
	case SLALevelInfo:
		compliance.InfoCount = 1
	case SLALevelWarning:
		compliance.WarningCount = 1
	case SLALevelCritical:
		compliance.CriticalCount = 1
	}

	compliance.ComplianceRate = float64(compliance.InfoCount) / float64(compliance.TotalCount) * 100

	// Convert history entries to the expected format
	historyEntries := make([]FreshnessHistoryEntry, len(history))
	for i, entry := range history {
		historyEntries[i] = FreshnessHistoryEntry{
			TableName:         entry.TableName,
			TablePath:         entry.TablePath,
			Timestamp:         entry.Timestamp,
			LastUpdateTime:    entry.LastUpdateTime,
			TimeSinceUpdate:   entry.TimeSinceUpdate,
			ExpectedFrequency: entry.ExpectedFrequency,
			Status:            entry.Status,
		}
	}

	return &FreshnessTrends{
		History:    historyEntries,
		Compliance: compliance,
	}, nil
}

// GetAllTablesTrends gets the freshness trends for all tables
func (m *SLAManager) GetAllTablesTrends() (*FreshnessTrends, error) {
	// Get all table statuses
	statuses, err := m.CheckAllFreshness()
	if err != nil {
		return nil, err
	}

	// Get all history
	history := m.historyMgr.getAllHistory()

	// Calculate compliance
	compliance := FreshnessCompliance{
		TotalCount: len(statuses),
	}

	for _, status := range statuses {
		switch status.Status {
		case SLALevelInfo:
			compliance.InfoCount++
		case SLALevelWarning:
			compliance.WarningCount++
		case SLALevelCritical:
			compliance.CriticalCount++
		}
	}

	if compliance.TotalCount > 0 {
		compliance.ComplianceRate = float64(compliance.InfoCount) / float64(compliance.TotalCount) * 100
	}

	return &FreshnessTrends{
		History:    history,
		Compliance: compliance,
	}, nil
}

// FreshnessTrends represents the historical freshness data and compliance statistics
type FreshnessTrends struct {
	// History contains the historical freshness status entries
	History []FreshnessHistoryEntry `json:"history"`

	// Compliance contains the aggregated compliance statistics
	Compliance FreshnessCompliance `json:"compliance"`
}

// FreshnessHistoryEntry represents a single historical freshness status entry
type FreshnessHistoryEntry struct {
	// TableName is the name of the table
	TableName string `json:"table_name"`

	// TablePath is the path to the table
	TablePath string `json:"table_path"`

	// Timestamp is the time when the status was recorded
	Timestamp time.Time `json:"timestamp"`

	// LastUpdateTime is the last time the table was updated
	LastUpdateTime time.Time `json:"last_update_time"`

	// TimeSinceUpdate is the duration since the last update
	TimeSinceUpdate time.Duration `json:"time_since_update"`

	// ExpectedFrequency is the expected update frequency for the table
	ExpectedFrequency time.Duration `json:"expected_frequency"`

	// Status is the freshness status (info, warning, critical)
	Status string `json:"status"`
}

// FreshnessCompliance represents the aggregated compliance statistics
type FreshnessCompliance struct {
	// InfoCount is the number of tables with "info" status
	InfoCount int `json:"info_count"`

	// WarningCount is the number of tables with "warning" status
	WarningCount int `json:"warning_count"`

	// CriticalCount is the number of tables with "critical" status
	CriticalCount int `json:"critical_count"`

	// TotalCount is the total number of tables
	TotalCount int `json:"total_count"`

	// ComplianceRate is the percentage of tables that are compliant (info status)
	ComplianceRate float64 `json:"compliance_rate"`
}

// historyManager manages the freshness history data
type historyManager struct {
	// history is a map of table name to history entries
	history map[string][]FreshnessHistoryEntry

	// maxHistoryEntries is the maximum number of history entries to keep per table
	maxHistoryEntries int

	// mutex protects the history map
	mutex sync.RWMutex
}

// newHistoryManager creates a new history manager
func newHistoryManager(maxEntries int) *historyManager {
	if maxEntries <= 0 {
		maxEntries = 100 // Default to 100 entries per table
	}

	return &historyManager{
		history:           make(map[string][]FreshnessHistoryEntry),
		maxHistoryEntries: maxEntries,
		mutex:             sync.RWMutex{},
	}
}

// addHistoryEntry adds a new history entry for a table
func (h *historyManager) addHistoryEntry(entry FreshnessHistoryEntry) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	tableName := entry.TableName
	entries, exists := h.history[tableName]

	if !exists {
		entries = make([]FreshnessHistoryEntry, 0, h.maxHistoryEntries)
	}

	// Add new entry at the beginning
	entries = append([]FreshnessHistoryEntry{entry}, entries...)

	// Trim if exceeds max entries
	if len(entries) > h.maxHistoryEntries {
		entries = entries[:h.maxHistoryEntries]
	}

	h.history[tableName] = entries
}

// getTableHistory gets the history entries for a specific table
func (h *historyManager) getTableHistory(tableName string) []FreshnessHistoryEntry {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	entries, exists := h.history[tableName]
	if !exists {
		return []FreshnessHistoryEntry{}
	}

	// Return a copy to avoid concurrent modification
	result := make([]FreshnessHistoryEntry, len(entries))
	copy(result, entries)

	return result
}

// getAllHistory gets all history entries for all tables
func (h *historyManager) getAllHistory() []FreshnessHistoryEntry {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	var allEntries []FreshnessHistoryEntry

	for _, entries := range h.history {
		allEntries = append(allEntries, entries...)
	}

	return allEntries
}
