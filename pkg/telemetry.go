package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// TelemetryConfig represents telemetry settings in config
// Add this to your main config struct
//
//	Telemetry TelemetryConfig `mapstructure:"telemetry"`
type TelemetryConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	Endpoint string `mapstructure:"endpoint"`
}

// TelemetryEvent represents a single telemetry event
// Extend fields as needed
// E.g., add Command, Args, Duration, Error, etc.
type TelemetryEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Version   string    `json:"version"`
}

// TelemetryManager handles sending telemetry events
// Usage: pkg.TelemetryManager.SendEvent(...)
type TelemetryManager struct {
	Enabled  bool
	Endpoint string
}

// NewTelemetryManager creates a TelemetryManager from config
func NewTelemetryManager(cfg TelemetryConfig) *TelemetryManager {
	return &TelemetryManager{
		Enabled:  cfg.Enabled,
		Endpoint: cfg.Endpoint,
	}
}

// SendEvent sends a telemetry event if enabled
func (tm *TelemetryManager) SendEvent(event TelemetryEvent) {
	if !tm.Enabled || tm.Endpoint == "" {
		return
	}
	data, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[telemetry] failed to marshal event: %v\n", err)
		return
	}
	go func() { // send asynchronously
		req, err := http.NewRequest("POST", tm.Endpoint, bytes.NewBuffer(data))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[telemetry] failed to create request: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[telemetry] failed to send event: %v\n", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 {
			fmt.Fprintf(os.Stderr, "[telemetry] server returned status: %s\n", resp.Status)
		}
	}()
}

// Helper: No-op singleton if not enabled
var Telemetry = &TelemetryManager{Enabled: false, Endpoint: ""}
