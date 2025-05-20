package freshness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/nessi-dev/nessi/pkg/logging"
	"github.com/nessi-dev/nessi/pkg/monitoring"
)

// SLAManager handles SLA monitoring for tables
type SLAManager struct {
	// configPath is the path to the SLA configuration file
	configPath string

	// configs is a map of table name to SLA configuration
	configs map[string]*SLAConfig

	// deltaConnector is the connector to the Delta Lake
	deltaConnector interface {
		GetTableMetadata(tablePath string) (*TableMetadata, error)
		ListTables() ([]string, error)
	}

	// mutex protects the configs map
	mutex sync.RWMutex

	// monitor is the monitoring system
	monitor *monitoring.Monitor

	// enabled indicates whether SLA monitoring is enabled
	enabled bool

	// historyMgr manages the freshness history data
	historyMgr *historyManager

	// lastCheck tracks the last time each table was checked
	lastCheck map[string]time.Time
}

// NewSLAManager creates a new SLA manager
func NewSLAManager(configPath string, monitor *monitoring.Monitor) (*SLAManager, error) {
	manager := &SLAManager{
		configPath: configPath,
		configs:    make(map[string]*SLAConfig),
		monitor:    monitor,
		enabled:    true,
		historyMgr: newHistoryManager(100), // Keep 100 history entries per table
		lastCheck:  make(map[string]time.Time),
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load existing configurations
	if err := manager.loadConfigs(); err != nil {
		// If the file doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load SLA configurations: %w", err)
		}
	}

	return manager, nil
}

// loadConfigs loads SLA configurations from the config file
func (m *SLAManager) loadConfigs() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Read config file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	// Parse config file
	var configs []*SLAConfig
	if err := json.Unmarshal(data, &configs); err != nil {
		return fmt.Errorf("failed to parse SLA configurations: %w", err)
	}

	// Store configs
	for _, config := range configs {
		m.configs[config.TableName] = config
	}

	return nil
}

// saveConfigs saves SLA configurations to the config file
func (m *SLAManager) saveConfigs() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Convert configs map to slice
	configs := make([]*SLAConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}

	// Marshal configs to JSON
	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal SLA configurations: %w", err)
	}

	// Write to file
	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write SLA configurations: %w", err)
	}

	return nil
}

// SetSLA sets the SLA configuration for a table
func (m *SLAManager) SetSLA(config *SLAConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid SLA configuration: %w", err)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Store config
	m.configs[config.TableName] = config

	// Save configs
	return m.saveConfigs()
}

// GetSLA gets the SLA configuration for a table
func (m *SLAManager) GetSLA(tableName string) (*SLAConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	config, ok := m.configs[tableName]
	if !ok {
		return nil, fmt.Errorf("no SLA configuration found for table %s", tableName)
	}

	return config, nil
}

// ListSLAs lists all SLA configurations
func (m *SLAManager) ListSLAs() []*SLAConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	configs := make([]*SLAConfig, 0, len(m.configs))
	for _, config := range m.configs {
		configs = append(configs, config)
	}

	return configs
}

// DeleteSLA deletes the SLA configuration for a table
func (m *SLAManager) DeleteSLA(tableName string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if config exists
	if _, ok := m.configs[tableName]; !ok {
		return fmt.Errorf("no SLA configuration found for table %s", tableName)
	}

	// Delete config
	delete(m.configs, tableName)

	// Save configs
	return m.saveConfigs()
}

// Enable enables SLA monitoring
func (m *SLAManager) Enable() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.enabled = true
}

// Disable disables SLA monitoring
func (m *SLAManager) Disable() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.enabled = false
}

// IsEnabled returns whether SLA monitoring is enabled
func (m *SLAManager) IsEnabled() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.enabled
}

// FreshnessStatus represents the freshness status of a table
type FreshnessStatus struct {
	// TableName is the name of the table
	TableName string `json:"table_name"`

	// TablePath is the path to the table
	TablePath string `json:"table_path"`

	// LastUpdateTime is the time of the last update
	LastUpdateTime time.Time `json:"last_update_time"`

	// TimeSinceUpdate is the time since the last update
	TimeSinceUpdate time.Duration `json:"time_since_update"`

	// ExpectedFrequency is the expected update frequency
	ExpectedFrequency time.Duration `json:"expected_frequency"`

	// Status is the freshness status
	Status SLALevel `json:"status"`

	// NextExpectedUpdate is the time of the next expected update
	NextExpectedUpdate time.Time `json:"next_expected_update"`

	// SLAConfig is the SLA configuration
	SLAConfig *SLAConfig `json:"sla_config"`
}

// CheckFreshness checks the freshness of a table
func (m *SLAManager) CheckFreshness(tableName string) (*FreshnessStatus, error) {
	// Check if SLA monitoring is enabled
	if !m.IsEnabled() {
		return nil, fmt.Errorf("SLA monitoring is disabled")
	}

	// Get SLA configuration
	config, err := m.GetSLA(tableName)
	if err != nil {
		return nil, err
	}

	// Check if SLA monitoring is enabled for this table
	if !config.Enabled {
		return nil, fmt.Errorf("SLA monitoring is disabled for table %s", tableName)
	}

	// Get table metadata
	metadataManager := datalake.NewMetadataManager(config.TablePath)
	table, err := metadataManager.ReadTableMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to read table metadata: %w", err)
	}

	// Calculate time since last update
	now := time.Now()
	timeSinceUpdate := now.Sub(table.LastModified)

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
		TableName:          config.TableName,
		TablePath:          config.TablePath,
		LastUpdateTime:     table.LastModified,
		TimeSinceUpdate:    timeSinceUpdate,
		ExpectedFrequency:  config.ExpectedFrequency,
		Status:             status,
		NextExpectedUpdate: nextExpectedUpdate,
		SLAConfig:          config,
	}

	// Record metrics if monitoring is enabled
	if m.monitor != nil {
		m.recordMetrics(freshnessStatus)
	}

	// Update history
	m.updateHistory(tableName, freshnessStatus)

	return freshnessStatus, nil
}

// CheckAllFreshness checks the freshness of all tables with SLA configurations
func (m *SLAManager) CheckAllFreshness() ([]*FreshnessStatus, error) {
	// Check if SLA monitoring is enabled
	if !m.IsEnabled() {
		return nil, fmt.Errorf("SLA monitoring is disabled")
	}

	// Get all SLA configurations
	configs := m.ListSLAs()

	// Check freshness for each table
	statuses := make([]*FreshnessStatus, 0, len(configs))
	for _, config := range configs {
		// Skip disabled configurations
		if !config.Enabled {
			continue
		}

		// Check freshness
		status, err := m.CheckFreshness(config.TableName)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to check freshness for table %s: %v", config.TableName, err))
			continue
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

// recordMetrics records freshness metrics
func (m *SLAManager) recordMetrics(status *FreshnessStatus) {
	// Record time since last update
	m.monitor.RecordMetric("freshness_time_since_update", float64(status.TimeSinceUpdate.Seconds()), map[string]string{
		"table":  status.TableName,
		"status": string(status.Status),
	})

	// Record expected frequency
	m.monitor.RecordMetric("freshness_expected_frequency", float64(status.ExpectedFrequency.Seconds()), map[string]string{
		"table": status.TableName,
	})

	// Record status as a numeric value (0 = info, 1 = warning, 2 = critical)
	var statusValue float64
	switch status.Status {
	case SLALevelInfo:
		statusValue = 0
	case SLALevelWarning:
		statusValue = 1
	case SLALevelCritical:
		statusValue = 2
	}

	m.monitor.RecordMetric("freshness_status", statusValue, map[string]string{
		"table": status.TableName,
	})
}

// StartMonitoring starts periodic SLA monitoring
func (m *SLAManager) StartMonitoring(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			// Skip if disabled
			if !m.IsEnabled() {
				continue
			}

			// Check freshness for all tables
			statuses, err := m.CheckAllFreshness()
			if err != nil {
				logging.Error("Failed to check freshness for all tables", err)
				continue
			}

			// Log results
			for _, status := range statuses {
				switch status.Status {
				case SLALevelCritical:
					logging.Error(fmt.Sprintf("Critical SLA violation for table %s: Last update was %s ago, expected every %s",
						status.TableName, FormatDuration(status.TimeSinceUpdate), FormatDuration(status.ExpectedFrequency)), nil)
				case SLALevelWarning:
					logging.Warn(fmt.Sprintf("Warning SLA violation for table %s: Last update was %s ago, expected every %s",
						status.TableName, FormatDuration(status.TimeSinceUpdate), FormatDuration(status.ExpectedFrequency)))
				case SLALevelInfo:
					logging.Info(fmt.Sprintf("Table %s is up-to-date: Last update was %s ago, expected every %s",
						status.TableName, FormatDuration(status.TimeSinceUpdate), FormatDuration(status.ExpectedFrequency)))
				}
			}
		}
	}()
}
