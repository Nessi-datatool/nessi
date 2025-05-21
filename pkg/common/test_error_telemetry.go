package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ErrorTelemetry represents the error telemetry system
type ErrorTelemetry struct {
	Enabled      bool
	Anonymous    bool
	StoragePath  string
	ErrorCounts  map[ErrorCode]int
	TotalErrors  int
	LastReported time.Time
	Logger       *TestLogger
	mu           sync.Mutex
}

// ErrorTelemetryConfig represents the configuration for error telemetry
type ErrorTelemetryConfig struct {
	Enabled     bool
	Anonymous   bool
	StoragePath string
}

// DefaultErrorTelemetryConfig returns the default configuration for error telemetry
func DefaultErrorTelemetryConfig() ErrorTelemetryConfig {
	return ErrorTelemetryConfig{
		Enabled:     true,
		Anonymous:   true,
		StoragePath: "",
	}
}

// NewErrorTelemetry creates a new error telemetry system
func NewErrorTelemetry(config ErrorTelemetryConfig) *ErrorTelemetry {
	// Set default storage path if not provided
	if config.StoragePath == "" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			config.StoragePath = fmt.Sprintf("%s/.nessi/error_telemetry.json", homeDir)
		} else {
			config.StoragePath = "error_telemetry.json"
		}
	}

	// Create telemetry system
	telemetry := &ErrorTelemetry{
		Enabled:      config.Enabled,
		Anonymous:    config.Anonymous,
		StoragePath:  config.StoragePath,
		ErrorCounts:  make(map[ErrorCode]int),
		TotalErrors:  0,
		LastReported: time.Now(),
		Logger:       NewTestLogger("info"),
	}

	// Load existing telemetry data if available
	telemetry.Load()

	return telemetry
}

// RecordError records an error in the telemetry system
func (t *ErrorTelemetry) RecordError(err error) {
	if !t.Enabled {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Increment total errors
	t.TotalErrors++

	// Check if it's a NessiError
	var nessiErr *NessiError
	if errors.As(err, &nessiErr) {
		// Increment error count for this error code
		t.ErrorCounts[nessiErr.Code]++
	}

	// Save telemetry data periodically (every 10 errors)
	if t.TotalErrors%10 == 0 {
		t.Save()
	}
}

// GetErrorStats returns the error statistics
func (t *ErrorTelemetry) GetErrorStats() map[string]interface{} {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create stats map
	stats := make(map[string]interface{})
	stats["total_errors"] = t.TotalErrors
	stats["error_counts"] = t.ErrorCounts
	stats["last_reported"] = t.LastReported

	return stats
}

// Save saves the telemetry data to disk
func (t *ErrorTelemetry) Save() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(t.StoragePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal telemetry data
	data := map[string]interface{}{
		"total_errors":  t.TotalErrors,
		"error_counts":  t.ErrorCounts,
		"last_reported": t.LastReported,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(t.StoragePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write telemetry data: %w", err)
	}

	return nil
}

// Load loads the telemetry data from disk
func (t *ErrorTelemetry) Load() error {
	// Check if file exists
	if _, err := os.Stat(t.StoragePath); os.IsNotExist(err) {
		return nil
	}

	// Read file
	jsonData, err := os.ReadFile(t.StoragePath)
	if err != nil {
		return fmt.Errorf("failed to read telemetry data: %w", err)
	}

	// Unmarshal telemetry data
	var data map[string]interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to unmarshal telemetry data: %w", err)
	}

	// Update telemetry data
	t.mu.Lock()
	defer t.mu.Unlock()

	if totalErrors, ok := data["total_errors"].(float64); ok {
		t.TotalErrors = int(totalErrors)
	}

	if lastReported, ok := data["last_reported"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, lastReported); err == nil {
			t.LastReported = parsedTime
		}
	}

	if errorCounts, ok := data["error_counts"].(map[string]interface{}); ok {
		for code, count := range errorCounts {
			if countFloat, ok := count.(float64); ok {
				t.ErrorCounts[ErrorCode(code)] = int(countFloat)
			}
		}
	}

	return nil
}

// EnableTelemetry enables error telemetry
func (t *ErrorTelemetry) EnableTelemetry() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Enabled = true
	t.Logger.Info("Error telemetry enabled")
}

// DisableTelemetry disables error telemetry
func (t *ErrorTelemetry) DisableTelemetry() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Enabled = false
	t.Logger.Info("Error telemetry disabled")
}

// SetAnonymous sets whether telemetry is anonymous
func (t *ErrorTelemetry) SetAnonymous(anonymous bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.Anonymous = anonymous
	t.Logger.Info(fmt.Sprintf("Error telemetry anonymity updated: %v", anonymous))
}

// ExportTelemetryToFile exports error telemetry data to a file
func (t *ErrorTelemetry) ExportTelemetryToFile(filePath string) error {
	if !t.Enabled {
		return fmt.Errorf("telemetry is disabled")
	}

	// Get error stats
	stats := t.GetErrorStats()

	// Create a report structure
	report := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"anonymous":    t.Anonymous,
		"total_errors": stats["total_errors"],
		"error_counts": stats["error_counts"],
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal telemetry data: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write telemetry data to file: %w", err)
	}

	return nil
}
