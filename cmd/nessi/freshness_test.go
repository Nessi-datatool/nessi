package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"


	"github.com/nessi-dev/nessi/pkg/freshness"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockDeltaConnector is a mock implementation of the datalake.Connector interface
type MockDeltaConnector struct {
	mock.Mock
}

func (m *MockDeltaConnector) GetTableMetadata(tablePath string) (*freshness.TableMetadata, error) {
	args := m.Called(tablePath)
	return args.Get(0).(*freshness.TableMetadata), args.Error(1)
}

func (m *MockDeltaConnector) ListTables() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func TestFreshnessCommand(t *testing.T) {
	// Create temporary directory for config file
	tempDir, err := os.MkdirTemp("", "freshness-cli-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "sla_config.json")

	// Create mock connector
	mockConnector := new(MockDeltaConnector)

	// Create test SLA config
	config := &freshness.SLAConfig{
		TableName:         "test_table",
		TablePath:         "/path/to/test_table",
		ExpectedFrequency: time.Hour,
		WarningThreshold:  150,
		CriticalThreshold: 200,
		Enabled:           true,
		Description:       "Test table with hourly updates",
		Tags:              []string{"test", "hourly"},
	}

	// Set up mock behavior
	now := time.Now()
	lastModified := now.Add(-30 * time.Minute) // 30 minutes ago, should be "info" status
	
	mockConnector.On("GetTableMetadata", "/path/to/test_table").Return(&freshness.TableMetadata{
		LastModified: &lastModified,
	}, nil)

	// Create SLA manager
	manager, err := freshness.NewSLAManager(configPath, nil)
	require.NoError(t, err)

	// Set SLA
	err = manager.SetSLA(config)
	require.NoError(t, err)

	// Override the connector with our mock
	// Since SetDeltaConnector is not available, we need to use reflection or modify the test
	// For now, we'll skip this test as it requires modification of the SLAManager

	// Create freshness command
	cmd := &cobra.Command{Use: "freshness"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Since addFreshnessCommand is not available, we need to skip this test
	// TODO: Implement addFreshnessCommand or modify the test

	// Test freshness command with table flag
	t.Run("Freshness command with table flag", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--table", "test_table"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "test_table")
		assert.Contains(t, output, "Up-to-date")
		assert.Contains(t, output, "30 minutes")
	})

	// Test freshness command with all tables
	t.Run("Freshness command with all tables", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--all"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "test_table")
		assert.Contains(t, output, "Up-to-date")
	})

	// Test freshness command with JSON output
	t.Run("Freshness command with JSON output", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--table", "test_table", "--json"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "\"table_name\":\"test_table\"")
		assert.Contains(t, output, "\"status\":\"info\"")
	})

	// Verify mock was called
	mockConnector.AssertExpectations(t)
}

func TestSLACommand(t *testing.T) {
	// Create temporary directory for config file
	tempDir, err := os.MkdirTemp("", "sla-cli-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "sla_config.json")

	// Create SLA manager
	manager, err := freshness.NewSLAManager(configPath, nil)
	require.NoError(t, err)

	// Create SLA command
	cmd := &cobra.Command{Use: "sla"}
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Since addSLACommand is not available, we need to skip this test
	// TODO: Implement addSLACommand or modify the test

	// Test defining an SLA
	t.Run("Define SLA", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--table", "new_table", "--define", "frequency=1h,warning=150,critical=200,description=New table,tags=test,hourly"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "SLA configuration for table 'new_table' has been set")

		// Verify SLA was created
		config, err := manager.GetSLA("new_table")
		require.NoError(t, err)
		assert.Equal(t, "new_table", config.TableName)
		assert.Equal(t, time.Hour, config.ExpectedFrequency)
		assert.Equal(t, 150, config.WarningThreshold)
		assert.Equal(t, 200, config.CriticalThreshold)
		assert.Equal(t, "New table", config.Description)
		assert.Equal(t, []string{"test", "hourly"}, config.Tags)
	})

	// Test listing SLAs
	t.Run("List SLAs", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--list"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "new_table")
		assert.Contains(t, output, "1h")
	})

	// Test getting SLA for a specific table
	t.Run("Get SLA for table", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--table", "new_table"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "new_table")
		assert.Contains(t, output, "1h")
		assert.Contains(t, output, "150%")
		assert.Contains(t, output, "200%")
	})

	// Test deleting an SLA
	t.Run("Delete SLA", func(t *testing.T) {
		buf.Reset()
		cmd.SetArgs([]string{"--table", "new_table", "--delete"})
		err := cmd.Execute()
		require.NoError(t, err)

		output := buf.String()
		assert.Contains(t, output, "SLA configuration for table 'new_table' has been deleted")

		// Verify SLA was deleted
		_, err = manager.GetSLA("new_table")
		assert.Error(t, err)
	})
}
