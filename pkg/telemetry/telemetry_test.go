package telemetry

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelemetryConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "telemetry_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test loading default config
	config, err := loadConfig(tempDir)
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.True(t, config.Enabled)
	assert.NotEmpty(t, config.InstallID)
	assert.False(t, config.FirstRun.IsZero())
	assert.False(t, config.LastRun.IsZero())
	assert.Equal(t, "0.1.0", config.Version)

	// Test saving and loading config
	config.Enabled = false
	config.Version = "0.2.0"
	err = saveConfig(tempDir, config)
	require.NoError(t, err)

	loadedConfig, err := loadConfig(tempDir)
	require.NoError(t, err)
	assert.Equal(t, false, loadedConfig.Enabled)
	assert.Equal(t, config.InstallID, loadedConfig.InstallID)
	assert.Equal(t, "0.2.0", loadedConfig.Version)
}

func TestEnableDisable(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "telemetry_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test enabling telemetry
	err = Enable(tempDir)
	require.NoError(t, err)

	status, err := GetStatus(tempDir)
	require.NoError(t, err)
	assert.True(t, status)

	// Test disabling telemetry
	err = Disable(tempDir)
	require.NoError(t, err)

	status, err = GetStatus(tempDir)
	require.NoError(t, err)
	assert.False(t, status)
}

func TestCollectTelemetryData(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "telemetry_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Enable telemetry
	err = Enable(tempDir)
	require.NoError(t, err)

	// Create test data
	commandStats := map[string]int{
		"scan":     10,
		"validate": 5,
	}
	featureUsage := map[string]int{
		"delta_lake": 15,
		"quality":    8,
	}
	errorCounts := map[string]int{
		"connection_error": 3,
		"validation_error": 2,
	}
	customMetrics := map[string]interface{}{
		"scan_duration_ms": 1500.5,
		"memory_usage_mb": 256.0,
	}

	// Collect telemetry data
	data, err := CollectTelemetryData(tempDir, commandStats, featureUsage, errorCounts, customMetrics)
	require.NoError(t, err)
	require.NotNil(t, data)

	// Verify collected data
	assert.NotEmpty(t, data.InstallID)
	assert.False(t, data.Timestamp.IsZero())
	assert.Equal(t, "0.1.0", data.Version)
	assert.Equal(t, commandStats, data.CommandStats)
	assert.Equal(t, featureUsage, data.FeatureUsage)
	assert.Equal(t, errorCounts, data.ErrorCounts)
	assert.Equal(t, customMetrics, data.CustomMetrics)
}

func TestRecordCommandAndFeature(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "telemetry_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Record a command
	err = RecordCommand(tempDir, "scan")
	require.NoError(t, err)
	err = RecordCommand(tempDir, "scan")
	require.NoError(t, err)
	err = RecordCommand(tempDir, "validate")
	require.NoError(t, err)

	// Verify command stats
	statsFile := filepath.Join(tempDir, "telemetry_commands.json")
	data, err := os.ReadFile(statsFile)
	require.NoError(t, err)

	var stats map[string]int
	err = json.Unmarshal(data, &stats)
	require.NoError(t, err)
	assert.Equal(t, 2, stats["scan"])
	assert.Equal(t, 1, stats["validate"])

	// Record a feature
	err = RecordFeatureUsage(tempDir, "delta_lake")
	require.NoError(t, err)
	err = RecordFeatureUsage(tempDir, "quality")
	require.NoError(t, err)
	err = RecordFeatureUsage(tempDir, "delta_lake")
	require.NoError(t, err)

	// Verify feature stats
	statsFile = filepath.Join(tempDir, "telemetry_features.json")
	data, err = os.ReadFile(statsFile)
	require.NoError(t, err)

	err = json.Unmarshal(data, &stats)
	require.NoError(t, err)
	assert.Equal(t, 2, stats["delta_lake"])
	assert.Equal(t, 1, stats["quality"])
}

func TestTelemetryMiddleware(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "telemetry_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Enable telemetry
	err = Enable(tempDir)
	require.NoError(t, err)

	// Create middleware
	middleware := NewTelemetryMiddleware(tempDir)

	// Create test handler
	handler := middleware.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/error" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Create test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make test requests
	resp, err := http.Get(server.URL + "/api/v1/tables")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(server.URL + "/api/v1/tables")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Post(server.URL+"/api/v1/validate", "application/json", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	resp, err = http.Get(server.URL + "/api/v1/error")
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()

	// Verify stats
	assert.Equal(t, 2, middleware.stats.Endpoints["/api/v1/tables"])
	assert.Equal(t, 1, middleware.stats.Endpoints["/api/v1/validate"])
	assert.Equal(t, 1, middleware.stats.Endpoints["/api/v1/error"])
	assert.Equal(t, 3, middleware.stats.Methods["GET"])
	assert.Equal(t, 1, middleware.stats.Methods["POST"])
	assert.Equal(t, 3, middleware.stats.StatusCodes[http.StatusOK])
	assert.Equal(t, 1, middleware.stats.StatusCodes[http.StatusInternalServerError])
	assert.Equal(t, 1, middleware.stats.Errors["error_500"])
}

func TestReporter(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "telemetry_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Enable telemetry
	err = Enable(tempDir)
	require.NoError(t, err)

	// Save original file path and set test file path
	savedFilePath := TelemetryFilePath
	// Create a local variable for the test file path
	testFilePath := filepath.Join(tempDir, "telemetry/test-data.json")

	// Use the test file path in the test
	t.Run("SendReport", func(t *testing.T) {
		// Set the file path for this test
		TelemetryFilePath = testFilePath
		// Restore it after the test
		defer func() { TelemetryFilePath = savedFilePath }()

		// Record some test data
		err = RecordCommand(tempDir, "scan")
		require.NoError(t, err)
		err = RecordFeatureUsage(tempDir, "delta_lake")
		require.NoError(t, err)
		err = RecordError(tempDir, "connection_error")
		require.NoError(t, err)

		// Create reporter with a short interval for testing
		reporter := NewReporter(tempDir, MinReportInterval)

		// Manually trigger a report
		reporter.sendReport()

		// Check if telemetry directory exists
		telemetryDir := filepath.Join(tempDir, "telemetry")
		require.DirExists(t, telemetryDir)

		// Find the telemetry file
		files, err := os.ReadDir(telemetryDir)
		require.NoError(t, err)
		require.NotEmpty(t, files, "No telemetry files found")

		// Get the latest telemetry file
		var latestFile os.DirEntry
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), "telemetry-") {
				latestFile = file
			}
		}
		require.NotNil(t, latestFile, "No telemetry file found")

		// Read the telemetry file
		filePath := filepath.Join(telemetryDir, latestFile.Name())
		fileData, err := os.ReadFile(filePath)
		require.NoError(t, err)

		// Parse the telemetry data
		var data TelemetryData
		err = json.Unmarshal(fileData, &data)
		require.NoError(t, err)

		// Verify the telemetry data
		assert.NotEmpty(t, data.InstallID)
		assert.Equal(t, "0.1.0", data.Version)
		assert.Equal(t, 1, data.CommandStats["scan"])
		assert.Equal(t, 1, data.FeatureUsage["delta_lake"])
		assert.Equal(t, 1, data.ErrorCounts["connection_error"])
	})
}

func TestGenerateBadge(t *testing.T) {
	badge, err := GenerateBadge("https://github.com/user/repo")
	require.NoError(t, err)
	assert.Contains(t, badge, "[![Powered by Nessi]")
	assert.Contains(t, badge, "https://img.shields.io/badge/powered%20by-nessi-blue")
}
