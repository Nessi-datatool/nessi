package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	// TelemetryConfigFile is the name of the telemetry configuration file
	TelemetryConfigFile = "telemetry.json"

	// DefaultTelemetryEnabled is the default telemetry setting
	DefaultTelemetryEnabled = true
)

// Default telemetry file path
var TelemetryFilePath = "${HOME}/.nessi/telemetry/data.json"

// TelemetryConfig represents the telemetry configuration
type TelemetryConfig struct {
	Enabled       bool      `json:"enabled"`
	InstallID     string    `json:"install_id"`
	FirstRun      time.Time `json:"first_run"`
	LastRun       time.Time `json:"last_run"`
	LastSubmitted time.Time `json:"last_submitted"`
	Version       string    `json:"version"`
}

// TelemetryData represents the telemetry data sent to the server
type TelemetryData struct {
	InstallID     string                 `json:"install_id"`
	Timestamp     time.Time              `json:"timestamp"`
	Version       string                 `json:"version"`
	OS            string                 `json:"os"`
	Arch          string                 `json:"arch"`
	CPUCores      int                    `json:"cpu_cores"`
	MemoryTotal   uint64                 `json:"memory_total"`
	Uptime        uint64                 `json:"uptime"`
	CommandStats  map[string]int         `json:"command_stats"`
	FeatureUsage  map[string]int         `json:"feature_usage"`
	ErrorCounts   map[string]int         `json:"error_counts"`
	CustomMetrics map[string]interface{} `json:"custom_metrics"`
}

// getConfigPath returns the path to the telemetry configuration file
func getConfigPath(dataDir string) string {
	return filepath.Join(dataDir, TelemetryConfigFile)
}

// loadConfig loads the telemetry configuration
func loadConfig(dataDir string) (*TelemetryConfig, error) {
	configPath := getConfigPath(dataDir)

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		config := &TelemetryConfig{
			Enabled:   DefaultTelemetryEnabled,
			InstallID: uuid.New().String(),
			FirstRun:  time.Now(),
			LastRun:   time.Now(),
			Version:   "0.1.0", // This should be updated from the version package
		}

		// Save config
		if err := saveConfig(dataDir, config); err != nil {
			return nil, fmt.Errorf("failed to save default config: %w", err)
		}

		return config, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	var config TelemetryConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// saveConfig saves the telemetry configuration
func saveConfig(dataDir string, config *TelemetryConfig) error {
	configPath := getConfigPath(dataDir)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Enable enables telemetry collection
func Enable(dataDir string) error {
	config, err := loadConfig(dataDir)
	if err != nil {
		return err
	}

	config.Enabled = true
	config.LastRun = time.Now()

	return saveConfig(dataDir, config)
}

// Disable disables telemetry collection
func Disable(dataDir string) error {
	config, err := loadConfig(dataDir)
	if err != nil {
		return err
	}

	config.Enabled = false
	config.LastRun = time.Now()

	return saveConfig(dataDir, config)
}

// GetStatus returns the current telemetry status
func GetStatus(dataDir string) (bool, error) {
	config, err := loadConfig(dataDir)
	if err != nil {
		return false, err
	}

	return config.Enabled, nil
}

// CollectTelemetryData collects telemetry data
func CollectTelemetryData(dataDir string, commandStats map[string]int, featureUsage map[string]int, errorCounts map[string]int, customMetrics map[string]interface{}) (*TelemetryData, error) {
	config, err := loadConfig(dataDir)
	if err != nil {
		return nil, err
	}

	// Skip if telemetry is disabled
	if !config.Enabled {
		return nil, nil
	}

	// Update last run time
	config.LastRun = time.Now()
	if err := saveConfig(dataDir, config); err != nil {
		return nil, err
	}

	// Collect system information
	cpuInfo, err := cpu.Info()
	if err != nil {
		// Don't fail if we can't get CPU info
		cpuInfo = []cpu.InfoStat{}
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		// Don't fail if we can't get memory info
		memInfo = &mem.VirtualMemoryStat{}
	}

	hostInfo, err := host.Info()
	if err != nil {
		// Don't fail if we can't get host info
		hostInfo = &host.InfoStat{}
	}

	// Create telemetry data
	data := &TelemetryData{
		InstallID:     config.InstallID,
		Timestamp:     time.Now(),
		Version:       config.Version,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		CPUCores:      len(cpuInfo),
		MemoryTotal:   memInfo.Total,
		Uptime:        hostInfo.Uptime,
		CommandStats:  commandStats,
		FeatureUsage:  featureUsage,
		ErrorCounts:   errorCounts,
		CustomMetrics: customMetrics,
	}

	return data, nil
}

// SendTelemetryData writes telemetry data to a file
func SendTelemetryData(dataDir string, data *TelemetryData) error {
	config, err := loadConfig(dataDir)
	if err != nil {
		return err
	}

	// Skip if telemetry is disabled
	if !config.Enabled {
		return nil
	}

	// Create telemetry directory if it doesn't exist
	telemetryDir := filepath.Join(dataDir, "telemetry")
	if err := os.MkdirAll(telemetryDir, 0755); err != nil {
		return fmt.Errorf("failed to create telemetry directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	filename := filepath.Join(telemetryDir, fmt.Sprintf("telemetry-%s.json", timestamp))

	// Marshal telemetry data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry data: %w", err)
	}

	// Write telemetry data to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write telemetry data: %w", err)
	}

	// Update last submitted time
	config.LastSubmitted = time.Now()
	return saveConfig(dataDir, config)
}

// GenerateBadge generates a GitHub usage badge for the repository
func GenerateBadge(repoURL string) (string, error) {
	if repoURL == "" {
		repoURL = "https://github.com/username/repo"
	}

	// Generate badge markdown
	badge := "[![Powered by Nessi](https://img.shields.io/badge/powered%20by-nessi-blue)](https://github.com/nessi-dev/nessi)"

	return badge, nil
}

// RecordCommand records a command execution for telemetry
func RecordCommand(dataDir string, command string) error {
	// Load command stats
	statsFile := filepath.Join(dataDir, "telemetry_commands.json")

	var stats map[string]int
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		stats = make(map[string]int)
	} else {
		data, err := os.ReadFile(statsFile)
		if err != nil {
			return fmt.Errorf("failed to read command stats: %w", err)
		}

		if err := json.Unmarshal(data, &stats); err != nil {
			// If the file is corrupted, start fresh
			stats = make(map[string]int)
		}
	}

	// Update stats
	stats[command]++

	// Save stats
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal command stats: %w", err)
	}

	if err := os.WriteFile(statsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write command stats: %w", err)
	}

	return nil
}

// RecordFeatureUsage records feature usage for telemetry
func RecordFeatureUsage(dataDir string, feature string) error {
	// Load feature usage stats
	statsFile := filepath.Join(dataDir, "telemetry_features.json")

	var stats map[string]int
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		stats = make(map[string]int)
	} else {
		data, err := os.ReadFile(statsFile)
		if err != nil {
			return fmt.Errorf("failed to read feature usage stats: %w", err)
		}

		if err := json.Unmarshal(data, &stats); err != nil {
			// If the file is corrupted, start fresh
			stats = make(map[string]int)
		}
	}

	// Update stats
	stats[feature]++

	// Save stats
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal feature usage stats: %w", err)
	}

	if err := os.WriteFile(statsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write feature usage stats: %w", err)
	}

	return nil
}

// RecordError records an error for telemetry
func RecordError(dataDir string, errorType string) error {
	// Load error stats
	statsFile := filepath.Join(dataDir, "telemetry_errors.json")

	var stats map[string]int
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		stats = make(map[string]int)
	} else {
		data, err := os.ReadFile(statsFile)
		if err != nil {
			return fmt.Errorf("failed to read error stats: %w", err)
		}

		if err := json.Unmarshal(data, &stats); err != nil {
			// If the file is corrupted, start fresh
			stats = make(map[string]int)
		}
	}

	// Update stats
	stats[errorType]++

	// Save stats
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal error stats: %w", err)
	}

	if err := os.WriteFile(statsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write error stats: %w", err)
	}

	return nil
}
