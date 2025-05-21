package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorTelemetry(t *testing.T) {
	// Create a temporary directory for telemetry data
	tempDir, err := os.MkdirTemp("", "nessi-telemetry-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create telemetry config with test storage path
	config := ErrorTelemetryConfig{
		Enabled:     true,
		Anonymous:   true,
		StoragePath: filepath.Join(tempDir, "error_telemetry.json"),
	}

	// Create telemetry system
	telemetry := NewErrorTelemetry(config)

	// Test recording errors
	// Record a standard error
	stdErr := fmt.Errorf("standard error")
	telemetry.RecordError(stdErr)

	// Record a NessiError
	nessiErr := NewError(ErrInvalidPath, "test path is invalid")
	telemetry.RecordError(nessiErr)

	// Record another NessiError with the same code
	nessiErr2 := NewError(ErrInvalidPath, "another path is invalid")
	telemetry.RecordError(nessiErr2)

	// Record a different NessiError
	nessiErr3 := NewError(ErrInvalidConfig, "config is invalid")
	telemetry.RecordError(nessiErr3)

	// Test getting error stats
	stats := telemetry.GetErrorStats()
	assert.Equal(t, 4, stats["total_errors"])

	// Check error counts
	errorCounts, ok := stats["error_counts"].(map[ErrorCode]int)
	assert.True(t, ok)
	assert.Equal(t, 2, errorCounts[ErrInvalidPath])
	assert.Equal(t, 1, errorCounts[ErrInvalidConfig])

	// Test saving telemetry data
	err = telemetry.Save()
	assert.NoError(t, err)

	// Check if file exists
	_, err = os.Stat(config.StoragePath)
	assert.NoError(t, err)

	// Read the file and check its contents
	jsonData, err := os.ReadFile(config.StoragePath)
	assert.NoError(t, err)

	// Parse JSON data
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	assert.NoError(t, err)

	// Check total errors
	totalErrors, ok := data["total_errors"].(float64)
	assert.True(t, ok)
	assert.Equal(t, float64(4), totalErrors)

	// Create a new telemetry system with the same storage path
	telemetry2 := NewErrorTelemetry(config)

	// Check if data was loaded correctly
	stats2 := telemetry2.GetErrorStats()
	assert.Equal(t, 4, stats2["total_errors"])

	// Test disabling telemetry
	telemetry2.DisableTelemetry()
	assert.False(t, telemetry2.Enabled)

	// Record an error when telemetry is disabled
	nessiErr4 := NewError(ErrInvalidConfig, "another config error")
	telemetry2.RecordError(nessiErr4)

	// Check that the error was not recorded
	stats2 = telemetry2.GetErrorStats()
	assert.Equal(t, 4, stats2["total_errors"])

	// Test enabling telemetry
	telemetry2.EnableTelemetry()
	assert.True(t, telemetry2.Enabled)

	// Record an error when telemetry is enabled
	telemetry2.RecordError(nessiErr4)

	// Check that the error was recorded
	stats2 = telemetry2.GetErrorStats()
	assert.Equal(t, 5, stats2["total_errors"])

	// Test reporting telemetry
	err = telemetry2.ReportTelemetry()
	assert.NoError(t, err)
}

func TestErrorTelemetryDefaultConfig(t *testing.T) {
	// Get default config
	config := DefaultErrorTelemetryConfig()

	// Check default values
	assert.True(t, config.Enabled)
	assert.True(t, config.Anonymous)
	assert.Equal(t, "", config.StoragePath)

	// Create telemetry system with default config
	telemetry := NewErrorTelemetry(config)

	// Check that storage path was set to a default value
	assert.NotEqual(t, "", telemetry.StoragePath)

	// Disable telemetry to avoid side effects
	telemetry.DisableTelemetry()
}

func TestErrorTelemetryWithInvalidPath(t *testing.T) {
	// Create config with invalid storage path
	config := ErrorTelemetryConfig{
		Enabled:     true,
		Anonymous:   true,
		StoragePath: "/invalid/path/that/does/not/exist/error_telemetry.json",
	}

	// Create telemetry system
	telemetry := NewErrorTelemetry(config)

	// Record an error
	nessiErr := NewError(ErrInvalidPath, "test path is invalid")
	telemetry.RecordError(nessiErr)

	// Try to save telemetry data (this should fail but not crash)
	err := telemetry.Save()
	assert.Error(t, err)

	// Disable telemetry to avoid side effects
	telemetry.DisableTelemetry()
}

func TestExportTelemetryToFile(t *testing.T) {
	// Create a temporary directory for telemetry data
	tempDir, err := os.MkdirTemp("", "nessi-telemetry-export-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create telemetry config with test storage path
	config := ErrorTelemetryConfig{
		Enabled:     true,
		Anonymous:   true,
		StoragePath: filepath.Join(tempDir, "error_telemetry.json"),
	}

	// Create telemetry system
	telemetry := NewErrorTelemetry(config)

	// Record some errors
	nessiErr1 := NewError(ErrInvalidPath, "test path is invalid")
	telemetry.RecordError(nessiErr1)

	nessiErr2 := NewError(ErrInvalidConfig, "config is invalid")
	telemetry.RecordError(nessiErr2)

	// Export telemetry to file
	exportPath := filepath.Join(tempDir, "export_telemetry.json")
	err = telemetry.ExportTelemetryToFile(exportPath)
	assert.NoError(t, err)

	// Check if export file exists
	_, err = os.Stat(exportPath)
	assert.NoError(t, err)

	// Read the exported file
	jsonData, err := os.ReadFile(exportPath)
	assert.NoError(t, err)

	// Parse JSON data
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	assert.NoError(t, err)

	// Check total errors
	totalErrors, ok := data["total_errors"].(float64)
	assert.True(t, ok)
	assert.Equal(t, float64(2), totalErrors)

	// Check timestamp exists
	_, ok = data["timestamp"].(string)
	assert.True(t, ok)

	// Check anonymity flag
	anonymous, ok := data["anonymous"].(bool)
	assert.True(t, ok)
	assert.Equal(t, true, anonymous)

	// Test exporting when telemetry is disabled
	telemetry.DisableTelemetry()
	exportPath2 := filepath.Join(tempDir, "export_disabled_telemetry.json")
	err = telemetry.ExportTelemetryToFile(exportPath2)
	assert.Error(t, err)

	// Check that the file wasn't created
	_, err = os.Stat(exportPath2)
	assert.True(t, os.IsNotExist(err))
}
