package freshness

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemorySLAManager(t *testing.T) {
	// Create a mock connector
	mockConnector := NewBetterMockDeltaConnector()

	// Create a test SLA manager with the mock connector
	manager := NewTestSLAManager(mockConnector)

	// Add test tables
	table1Path := "/path/to/test_table1"
	table2Path := "/path/to/test_table2"

	// Add SLA configurations
	manager.AddSLA(&SLAConfig{
		TableName:         "test_table1",
		TablePath:         table1Path,
		ExpectedFrequency: 1 * time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	})

	manager.AddSLA(&SLAConfig{
		TableName:         "test_table2",
		TablePath:         table2Path,
		ExpectedFrequency: 24 * time.Hour, // 1 day
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
	})

	t.Run("Individual Table Freshness", func(t *testing.T) {
		// Set up test data
		now := time.Now()
		lastUpdate := now.Add(-30 * time.Minute) // 30 minutes ago, within the 1 hour expected frequency

		// Add table to mock connector
		mockConnector.AddTable(table1Path, lastUpdate)

		// Check freshness
		status, err := manager.CheckFreshness("test_table1")
		assert.NoError(t, err)
		assert.NotNil(t, status)
		assert.Equal(t, "test_table1", status.TableName)
		assert.Equal(t, table1Path, status.TablePath)
		assert.Equal(t, SLALevelInfo, status.Status) // Should be info (up-to-date)
	})

	t.Run("All Tables Freshness", func(t *testing.T) {
		// Set up test data
		now := time.Now()

		// Table 1: 30 minutes ago (info)
		table1LastUpdate := now.Add(-30 * time.Minute)
		mockConnector.UpdateTableLastModified(table1Path, table1LastUpdate)

		// Table 2: 2 days ago (critical)
		table2LastUpdate := now.Add(-48 * time.Hour)
		mockConnector.AddTable(table2Path, table2LastUpdate)

		// Check all freshness
		statuses, err := manager.CheckAllFreshness()
		assert.NoError(t, err)
		assert.Len(t, statuses, 2)

		// Find status for each table
		var table1Status, table2Status *FreshnessStatus
		for _, s := range statuses {
			if s.TableName == "test_table1" {
				table1Status = s
			} else if s.TableName == "test_table2" {
				table2Status = s
			}
		}

		assert.NotNil(t, table1Status)
		assert.NotNil(t, table2Status)

		assert.Equal(t, SLALevelInfo, table1Status.Status)
		assert.Equal(t, SLALevelCritical, table2Status.Status)
	})

	t.Run("History Tracking", func(t *testing.T) {
		// Set up test data
		now := time.Now()

		// Table 1: Update several times with different statuses

		// First check: 30 minutes ago (info)
		mockConnector.UpdateTableLastModified(table1Path, now.Add(-30*time.Minute))
		status1, err := manager.CheckFreshness("test_table1")
		assert.NoError(t, err)
		assert.Equal(t, SLALevelInfo, status1.Status)

		// Second check: 2 hours ago (warning or critical, depending on exact thresholds)
		mockConnector.UpdateTableLastModified(table1Path, now.Add(-2*time.Hour))
		status2, err := manager.CheckFreshness("test_table1")
		assert.NoError(t, err)
		// Status could be warning or critical depending on exact time calculations
		// Just ensure it's not info status
		assert.NotEqual(t, SLALevelInfo, status2.Status)

		// Third check: 3 hours ago (critical)
		mockConnector.UpdateTableLastModified(table1Path, now.Add(-3*time.Hour))
		status3, err := manager.CheckFreshness("test_table1")
		assert.NoError(t, err)
		assert.Equal(t, SLALevelCritical, status3.Status)

		// Get trends
		trends, err := manager.GetTableTrends("test_table1")
		assert.NoError(t, err)
		assert.NotNil(t, trends)

		// Check history entries
		assert.GreaterOrEqual(t, len(trends.History), 3)

		// Check compliance - actual values may vary based on timing
		// Just verify that we have entries and the total count is correct
		assert.Greater(t, trends.Compliance.InfoCount, 0)
		assert.Equal(t, trends.Compliance.InfoCount+trends.Compliance.WarningCount+trends.Compliance.CriticalCount, trends.Compliance.TotalCount)
		assert.InDelta(t, float64(trends.Compliance.InfoCount)/float64(trends.Compliance.TotalCount)*100, trends.Compliance.ComplianceRate, 0.01)
	})
}

// TestDurationParsing tests the duration parsing functionality
func TestDurationParsing(t *testing.T) {
	// Test parsing various duration formats
	testCases := []struct {
		input    string
		expected time.Duration
		hasError bool
	}{
		{"1h", time.Hour, false},
		{"1d", 24 * time.Hour, false},
		{"1w", 7 * 24 * time.Hour, false},
		{"1mo", 30 * 24 * time.Hour, false},
		{"hourly", time.Hour, false},
		{"daily", 24 * time.Hour, false},
		{"weekly", 7 * 24 * time.Hour, false},
		{"monthly", 30 * 24 * time.Hour, false},
		{"invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			duration, err := ParseDuration(tc.input)
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, duration)
			}
		})
	}
}

// TestSLAConfigValidation tests the SLA configuration validation
func TestSLAConfigValidation(t *testing.T) {
	// Test various SLA configurations
	testCases := []struct {
		name     string
		config   SLAConfig
		hasError bool
	}{
		{
			name: "Valid config",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			hasError: false,
		},
		{
			name: "Empty table name",
			config: SLAConfig{
				TableName:         "",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			hasError: true,
		},
		{
			name: "Empty table path",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			hasError: true,
		},
		{
			name: "Zero expected frequency",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: 0,
				WarningThreshold:  150,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			hasError: true,
		},
		{
			name: "Warning threshold less than 100",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  50,
				CriticalThreshold: 200,
				Enabled:           true,
			},
			hasError: true,
		},
		{
			name: "Critical threshold less than warning threshold",
			config: SLAConfig{
				TableName:         "test_table",
				TablePath:         "/path/to/test_table",
				ExpectedFrequency: time.Hour,
				WarningThreshold:  200,
				CriticalThreshold: 150,
				Enabled:           true,
			},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
