package freshness

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSLAManagerFreshness tests the SLA freshness checking functionality
func TestSLAManagerFreshness(t *testing.T) {
	// Create a better mock connector
	mockConnector := NewBetterMockDeltaConnector()

	// Create a test SLA manager with the mock connector
	manager := NewTestSLAManager(mockConnector)

	// Add test SLA config
	manager.AddSLA(&SLAConfig{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150, // 150% of expected frequency (90 minutes)
		CriticalThreshold: 200, // 200% of expected frequency (120 minutes)
		Enabled:           true,
	})

	// Test case 1: Table is up-to-date (last modified 30 minutes ago)
	t.Run("Info Status", func(t *testing.T) {
		now := time.Now()
		lastModified := now.Add(-30 * time.Minute)

		// Set up mock behavior
		mockConnector.AddTable("/path/to/test_table", lastModified)

		status, err := manager.CheckFreshness("test_table")
		require.NoError(t, err)

		assert.Equal(t, "test_table", status.TableName)
		assert.Equal(t, "/path/to/test_table", status.TablePath)
		assert.Equal(t, lastModified.Unix(), status.LastUpdateTime.Unix())
		assert.Equal(t, SLALevelInfo, status.Status)
	})

	// Test case 2: Table is approaching SLA violation (last modified 100 minutes ago)
	t.Run("Warning Status", func(t *testing.T) {
		now := time.Now()
		lastModified := now.Add(-100 * time.Minute)

		// Set up mock behavior
		mockConnector.UpdateTableLastModified("/path/to/test_table", lastModified)

		status, err := manager.CheckFreshness("test_table")
		require.NoError(t, err)

		assert.Equal(t, "test_table", status.TableName)
		assert.Equal(t, "/path/to/test_table", status.TablePath)
		assert.Equal(t, lastModified.Unix(), status.LastUpdateTime.Unix())
		assert.Equal(t, SLALevelWarning, status.Status)
	})

	// Test case 3: Table has SLA violation (last modified 150 minutes ago)
	t.Run("Critical Status", func(t *testing.T) {
		now := time.Now()
		lastModified := now.Add(-150 * time.Minute)

		// Set up mock behavior
		mockConnector.UpdateTableLastModified("/path/to/test_table", lastModified)

		status, err := manager.CheckFreshness("test_table")
		require.NoError(t, err)

		assert.Equal(t, "test_table", status.TableName)
		assert.Equal(t, "/path/to/test_table", status.TablePath)
		assert.Equal(t, lastModified.Unix(), status.LastUpdateTime.Unix())
		assert.Equal(t, SLALevelCritical, status.Status)
	})
}

// TestSLAManagerAllFreshness tests checking freshness for all tables
func TestSLAManagerAllFreshness(t *testing.T) {
	// Create a better mock connector
	mockConnector := NewBetterMockDeltaConnector()

	// Create a test SLA manager with the mock connector
	manager := NewTestSLAManager(mockConnector)

	// Add test SLA configs
	manager.AddSLA(&SLAConfig{
		TableName:         "test_table1",
		TablePath:         "/path/to/test_table1",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	})

	manager.AddSLA(&SLAConfig{
		TableName:         "test_table2",
		TablePath:         "/path/to/test_table2",
		ExpectedFrequency: time.Hour * 24, // Daily table
		WarningThreshold:  125,
		CriticalThreshold: 150,
		Enabled:           true,
	})

	// Set up mock behavior
	now := time.Now()
	lastModified1 := now.Add(-30 * time.Minute) // 30 minutes ago, should be "info" status
	lastModified2 := now.Add(-36 * time.Hour)   // 36 hours ago, should be "critical" status for daily table

	// Add tables to mock connector
	mockConnector.AddTable("/path/to/test_table1", lastModified1)
	mockConnector.AddTable("/path/to/test_table2", lastModified2)

	// Check all freshness
	statuses, err := manager.CheckAllFreshness()
	require.NoError(t, err)

	// Should have 2 statuses
	assert.Len(t, statuses, 2)

	// Find status for each table
	var status1, status2 *FreshnessStatus
	for _, status := range statuses {
		if status.TableName == "test_table1" {
			status1 = status
		} else if status.TableName == "test_table2" {
			status2 = status
		}
	}

	require.NotNil(t, status1)
	require.NotNil(t, status2)

	// Check status values
	assert.Equal(t, SLALevelInfo, status1.Status)
	assert.Equal(t, SLALevelCritical, status2.Status)
}
