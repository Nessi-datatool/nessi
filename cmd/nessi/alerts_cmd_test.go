package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nessi-dev/nessi-dev/pkg/monitoring/alerts"
)

func setupTestAlertManager(t *testing.T) (*alerts.AlertManager, string) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "alerts-cmd-test")
	require.NoError(t, err)

	// Create alert manager
	manager, err := alerts.NewAlertManager(tempDir)
	require.NoError(t, err)

	return manager, tempDir
}

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err := root.Execute()
	return buf.String(), err
}

func TestAlertsListCommand(t *testing.T) {
	// Setup test environment
	manager, tempDir := setupTestAlertManager(t)
	defer os.RemoveAll(tempDir)

	// Set config for test
	origConfig := config
	config.DataDir = tempDir
	defer func() { config = origConfig }()

	// Create test alerts
	alert1 := &alerts.Alert{
		ID:          uuid.New().String(),
		Name:        "Test Alert 1",
		Description: "Test alert 1 description",
		Type:        alerts.TypeQuality,
		Severity:    alerts.SeverityWarning,
		Status:      alerts.StatusActive,
		Source:      "test-source-1",
		Timestamp:   time.Now(),
		LastUpdated: time.Now(),
	}

	alert2 := &alerts.Alert{
		ID:          uuid.New().String(),
		Name:        "Test Alert 2",
		Description: "Test alert 2 description",
		Type:        alerts.TypeAnomaly,
		Severity:    alerts.SeverityCritical,
		Status:      alerts.StatusActive,
		Source:      "test-source-2",
		Timestamp:   time.Now(),
		LastUpdated: time.Now(),
	}

	err := manager.CreateAlert(alert1)
	require.NoError(t, err)

	err = manager.CreateAlert(alert2)
	require.NoError(t, err)

	// Test list command
	output, err := executeCommand(rootCmd, "alerts", "list")
	require.NoError(t, err)
	assert.Contains(t, output, "Test Alert 1")
	assert.Contains(t, output, "Test Alert 2")
	assert.Contains(t, output, "warning")
	assert.Contains(t, output, "critical")

	// Test list with filter
	output, err = executeCommand(rootCmd, "alerts", "list", "--type", "quality")
	require.NoError(t, err)
	assert.Contains(t, output, "Test Alert 1")
	assert.NotContains(t, output, "Test Alert 2")

	// Test list with JSON output
	output, err = executeCommand(rootCmd, "alerts", "list", "--output", "json")
	require.NoError(t, err)
	assert.Contains(t, output, "\"id\":")
	assert.Contains(t, output, "\"name\":\"Test Alert 1\"")
	assert.Contains(t, output, "\"name\":\"Test Alert 2\"")
}

func TestAlertsGetCommand(t *testing.T) {
	// Setup test environment
	manager, tempDir := setupTestAlertManager(t)
	defer os.RemoveAll(tempDir)

	// Set config for test
	origConfig := config
	config.DataDir = tempDir
	defer func() { config = origConfig }()

	// Create test alert
	alert := &alerts.Alert{
		ID:          uuid.New().String(),
		Name:        "Test Alert",
		Description: "Test alert description",
		Type:        alerts.TypeQuality,
		Severity:    alerts.SeverityWarning,
		Status:      alerts.StatusActive,
		Source:      "test-source",
		Timestamp:   time.Now(),
		LastUpdated: time.Now(),
		Value:       95.5,
		Threshold:   90.0,
		ComparisonOperator: ">",
		Labels: map[string]string{
			"environment": "test",
			"component":   "data-quality",
		},
		Annotations: map[string]string{
			"summary": "Data quality score exceeded threshold",
		},
	}

	err := manager.CreateAlert(alert)
	require.NoError(t, err)

	// Test get command
	output, err := executeCommand(rootCmd, "alerts", "get", alert.ID)
	require.NoError(t, err)
	assert.Contains(t, output, "Test Alert")
	assert.Contains(t, output, "Test alert description")
	assert.Contains(t, output, "quality")
	assert.Contains(t, output, "warning")
	assert.Contains(t, output, "test-source")
	assert.Contains(t, output, "95.50")
	assert.Contains(t, output, "90.00")
	assert.Contains(t, output, ">")
	assert.Contains(t, output, "environment: test")
	assert.Contains(t, output, "component: data-quality")
	assert.Contains(t, output, "summary: Data quality score exceeded threshold")

	// Test get with JSON output
	output, err = executeCommand(rootCmd, "alerts", "get", alert.ID, "--output", "json")
	require.NoError(t, err)
	assert.Contains(t, output, "\"id\":\""+alert.ID+"\"")
	assert.Contains(t, output, "\"name\":\"Test Alert\"")
	assert.Contains(t, output, "\"description\":\"Test alert description\"")
}

func TestAlertsCreateCommand(t *testing.T) {
	// Setup test environment
	_, tempDir := setupTestAlertManager(t)
	defer os.RemoveAll(tempDir)

	// Set config for test
	origConfig := config
	config.DataDir = tempDir
	defer func() { config = origConfig }()

	// Test create command
	output, err := executeCommand(rootCmd, "alerts", "create",
		"--name", "New Test Alert",
		"--description", "New test alert description",
		"--type", "quality",
		"--severity", "warning",
		"--source", "test-create",
		"--value", "85.5",
		"--threshold", "90.0",
		"--operator", "<",
		"--label", "environment=test",
		"--label", "component=data-quality",
		"--annotation", "summary=Data quality score below threshold",
	)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert created with ID:")

	// Verify alert was created
	alertsDir := filepath.Join(tempDir, "alerts")
	files, err := os.ReadDir(alertsDir)
	require.NoError(t, err)
	assert.NotEmpty(t, files)

	// Test list to verify alert is there
	output, err = executeCommand(rootCmd, "alerts", "list")
	require.NoError(t, err)
	assert.Contains(t, output, "New Test Alert")
	assert.Contains(t, output, "warning")
	assert.Contains(t, output, "test-create")
}

func TestAlertManagementCommands(t *testing.T) {
	// Setup test environment
	manager, tempDir := setupTestAlertManager(t)
	defer os.RemoveAll(tempDir)

	// Set config for test
	origConfig := config
	config.DataDir = tempDir
	defer func() { config = origConfig }()

	// Create test alert
	alert := &alerts.Alert{
		ID:          uuid.New().String(),
		Name:        "Management Test Alert",
		Description: "Test alert for management commands",
		Type:        alerts.TypeQuality,
		Severity:    alerts.SeverityWarning,
		Status:      alerts.StatusActive,
		Source:      "test-management",
		Timestamp:   time.Now(),
		LastUpdated: time.Now(),
	}

	err := manager.CreateAlert(alert)
	require.NoError(t, err)

	// Test acknowledge command
	output, err := executeCommand(rootCmd, "alerts", "acknowledge", alert.ID, "--user", "test-user")
	require.NoError(t, err)
	assert.Contains(t, output, "Alert "+alert.ID+" acknowledged by test-user")

	// Verify alert was acknowledged
	acknowledgedAlert, err := manager.GetAlert(alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alerts.StatusAcknowledged, acknowledgedAlert.Status)
	assert.Equal(t, "test-user", acknowledgedAlert.AcknowledgedBy)

	// Test silence command
	output, err = executeCommand(rootCmd, "alerts", "silence", alert.ID,
		"--user", "silence-user",
		"--reason", "Testing silence",
		"--duration", "1h",
	)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert "+alert.ID+" silenced by silence-user for 1h0m0s")
	assert.Contains(t, output, "Reason: Testing silence")

	// Verify alert was silenced
	silencedAlert, err := manager.GetAlert(alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alerts.StatusSilenced, silencedAlert.Status)
	assert.Equal(t, "silence-user", silencedAlert.SilencedBy)
	assert.Equal(t, "Testing silence", silencedAlert.SilenceReason)
	assert.NotNil(t, silencedAlert.SilencedUntil)

	// Test resolve command
	output, err = executeCommand(rootCmd, "alerts", "resolve", alert.ID)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert "+alert.ID+" resolved")

	// Verify alert was resolved
	resolvedAlert, err := manager.GetAlert(alert.ID)
	require.NoError(t, err)
	assert.Equal(t, alerts.StatusResolved, resolvedAlert.Status)
	assert.NotNil(t, resolvedAlert.ResolvedAt)

	// Test delete command
	output, err = executeCommand(rootCmd, "alerts", "delete", alert.ID)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert "+alert.ID+" deleted")

	// Verify alert was deleted
	_, err = manager.GetAlert(alert.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "alert not found")
}

func TestAlertRulesCommands(t *testing.T) {
	// Setup test environment
	manager, tempDir := setupTestAlertManager(t)
	defer os.RemoveAll(tempDir)

	// Set config for test
	origConfig := config
	config.DataDir = tempDir
	defer func() { config = origConfig }()

	// Test create rule command
	output, err := executeCommand(rootCmd, "alerts", "rules", "create",
		"--name", "Test Rule",
		"--description", "Test rule description",
		"--type", "quality",
		"--severity", "warning",
		"--source", "test-rule-source",
		"--metric", "quality_score",
		"--threshold", "90.0",
		"--operator", "<",
		"--window", "1h",
		"--interval", "5m",
		"--user", "test-user",
		"--label", "environment=test",
		"--channel", "email",
		"--recipient", "test@example.com",
	)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert rule created with ID:")

	// Get the rule ID from the rules directory
	rulesDir := filepath.Join(tempDir, "rules")
	files, err := os.ReadDir(rulesDir)
	require.NoError(t, err)
	require.Len(t, files, 1)
	ruleID := files[0].Name()
	ruleID = ruleID[:len(ruleID)-5] // Remove .json extension

	// Test list rules command
	output, err = executeCommand(rootCmd, "alerts", "rules", "list")
	require.NoError(t, err)
	assert.Contains(t, output, "Test Rule")
	assert.Contains(t, output, "quality")
	assert.Contains(t, output, "warning")
	assert.Contains(t, output, "quality_score")
	assert.Contains(t, output, "90.00")
	assert.Contains(t, output, "true") // Enabled by default

	// Test get rule command
	output, err = executeCommand(rootCmd, "alerts", "rules", "get", ruleID)
	require.NoError(t, err)
	assert.Contains(t, output, "Test Rule")
	assert.Contains(t, output, "Test rule description")
	assert.Contains(t, output, "quality")
	assert.Contains(t, output, "warning")
	assert.Contains(t, output, "test-rule-source")
	assert.Contains(t, output, "quality_score")
	assert.Contains(t, output, "90.00")
	assert.Contains(t, output, "<")
	assert.Contains(t, output, "1h0m0s")
	assert.Contains(t, output, "5m0s")
	assert.Contains(t, output, "true") // Enabled
	assert.Contains(t, output, "test-user")
	assert.Contains(t, output, "environment: test")
	assert.Contains(t, output, "email")
	assert.Contains(t, output, "test@example.com")

	// Test update rule command
	output, err = executeCommand(rootCmd, "alerts", "rules", "update", ruleID,
		"--description", "Updated rule description",
		"--severity", "critical",
		"--threshold", "95.0",
		"--user", "admin-user",
	)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert rule "+ruleID+" updated")

	// Verify rule was updated
	output, err = executeCommand(rootCmd, "alerts", "rules", "get", ruleID)
	require.NoError(t, err)
	assert.Contains(t, output, "Updated rule description")
	assert.Contains(t, output, "critical")
	assert.Contains(t, output, "95.00")
	assert.Contains(t, output, "admin-user")

	// Test disable rule command
	output, err = executeCommand(rootCmd, "alerts", "rules", "disable", ruleID, "--user", "admin-user")
	require.NoError(t, err)
	assert.Contains(t, output, "Alert rule "+ruleID+" disabled")

	// Verify rule was disabled
	output, err = executeCommand(rootCmd, "alerts", "rules", "get", ruleID)
	require.NoError(t, err)
	assert.Contains(t, output, "Enabled:             false")

	// Test enable rule command
	output, err = executeCommand(rootCmd, "alerts", "rules", "enable", ruleID, "--user", "admin-user")
	require.NoError(t, err)
	assert.Contains(t, output, "Alert rule "+ruleID+" enabled")

	// Verify rule was enabled
	output, err = executeCommand(rootCmd, "alerts", "rules", "get", ruleID)
	require.NoError(t, err)
	assert.Contains(t, output, "Enabled:             true")

	// Test delete rule command
	output, err = executeCommand(rootCmd, "alerts", "rules", "delete", ruleID)
	require.NoError(t, err)
	assert.Contains(t, output, "Alert rule "+ruleID+" deleted")

	// Verify rule was deleted
	files, err = os.ReadDir(rulesDir)
	require.NoError(t, err)
	assert.Empty(t, files)
}
