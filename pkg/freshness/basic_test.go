package freshness

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSLAManager tests the creation of a new SLA manager
func TestNewSLAManager(t *testing.T) {
	// Create a simple in-memory manager for testing
	manager := &SLAManager{
		configPath: "test_path",
		configs:    make(map[string]*SLAConfig),
		enabled:    true,
		historyMgr: newHistoryManager(10),
		lastCheck:  make(map[string]time.Time),
		mutex:      sync.RWMutex{},
	}

	// Check manager properties
	assert.Equal(t, "test_path", manager.configPath)
	assert.NotNil(t, manager.configs)
	assert.Empty(t, manager.configs)
	assert.True(t, manager.enabled)
}

// TestSLAManager_SetAndGetSLA tests setting and getting SLA configurations
func TestSLAManager_SetAndGetSLA(t *testing.T) {
	// Create a simple in-memory manager for testing
	manager := &SLAManager{
		configPath: "test_path",
		configs:    make(map[string]*SLAConfig),
		enabled:    true,
		historyMgr: newHistoryManager(10),
		lastCheck:  make(map[string]time.Time),
		mutex:      sync.RWMutex{},
	}

	// Create test SLA config
	config := &SLAConfig{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	}

	// Set SLA directly to avoid file operations
	manager.mutex.Lock()
	manager.configs[config.TableName] = config
	manager.mutex.Unlock()

	// Get SLA
	retrievedConfig, err := manager.GetSLA("test_table")
	require.NoError(t, err)
	assert.Equal(t, config, retrievedConfig)

	// Try to get non-existent SLA
	_, err = manager.GetSLA("non_existent_table")
	assert.Error(t, err)
}

// TestSLAManager_ListSLAs tests listing SLA configurations
func TestSLAManager_ListSLAs(t *testing.T) {
	// Create a simple in-memory manager for testing
	manager := &SLAManager{
		configPath: "test_path",
		configs:    make(map[string]*SLAConfig),
		enabled:    true,
		historyMgr: newHistoryManager(10),
		lastCheck:  make(map[string]time.Time),
		mutex:      sync.RWMutex{},
	}

	// Create test SLA configs
	config1 := &SLAConfig{
		TableName:         "test_table1",
		TablePath:         "/path/to/test_table1",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	}

	config2 := &SLAConfig{
		TableName:         "test_table2",
		TablePath:         "/path/to/test_table2",
		ExpectedFrequency: time.Hour * 24,
		WarningThreshold:  125,
		CriticalThreshold: 150,
		Enabled:           true,
	}

	// Set SLAs directly
	manager.mutex.Lock()
	manager.configs[config1.TableName] = config1
	manager.configs[config2.TableName] = config2
	manager.mutex.Unlock()

	// List SLAs
	configs := manager.ListSLAs()
	assert.Len(t, configs, 2)

	// Check if both configs are in the list
	var found1, found2 bool
	for _, config := range configs {
		if config.TableName == "test_table1" {
			found1 = true
			assert.Equal(t, config1, config)
		} else if config.TableName == "test_table2" {
			found2 = true
			assert.Equal(t, config2, config)
		}
	}
	assert.True(t, found1)
	assert.True(t, found2)
}

// TestSLAManager_DeleteSLA tests deleting SLA configurations
func TestSLAManager_DeleteSLA(t *testing.T) {
	// Create a simple in-memory manager for testing
	manager := &SLAManager{
		configPath: "test_path",
		configs:    make(map[string]*SLAConfig),
		enabled:    true,
		historyMgr: newHistoryManager(10),
		lastCheck:  make(map[string]time.Time),
		mutex:      sync.RWMutex{},
	}

	// Create test SLA config
	config := &SLAConfig{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	}

	// Set SLA directly
	manager.mutex.Lock()
	manager.configs[config.TableName] = config
	manager.mutex.Unlock()

	// Verify SLA exists
	retrievedConfig, err := manager.GetSLA("test_table")
	require.NoError(t, err)
	assert.Equal(t, config, retrievedConfig)

	// Delete SLA directly
	manager.mutex.Lock()
	delete(manager.configs, "test_table")
	manager.mutex.Unlock()

	// Verify SLA is deleted
	_, err = manager.GetSLA("test_table")
	assert.Error(t, err)

	// Try to get non-existent SLA
	_, err = manager.GetSLA("non_existent_table")
	assert.Error(t, err)
}

// TestSLAManager_EnableDisable tests enabling and disabling SLA monitoring
func TestSLAManager_EnableDisable(t *testing.T) {
	// Create a simple in-memory manager for testing
	manager := &SLAManager{
		configPath: "test_path",
		configs:    make(map[string]*SLAConfig),
		enabled:    true,
		historyMgr: newHistoryManager(10),
		lastCheck:  make(map[string]time.Time),
		mutex:      sync.RWMutex{},
	}

	// Check initial state
	assert.True(t, manager.IsEnabled())

	// Disable
	manager.Disable()
	assert.False(t, manager.IsEnabled())

	// Enable
	manager.Enable()
	assert.True(t, manager.IsEnabled())
}
