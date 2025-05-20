package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// TelemetryCollector collects telemetry data for CLI commands
type TelemetryCollector struct {
	mu            sync.RWMutex
	dataDir       string
	events        map[string][]map[string]interface{}
	stats         map[string]int
	enabled       bool
	flushInterval time.Duration
	stopCh        chan struct{}
}

// NewTelemetryCollector creates a new TelemetryCollector
func NewTelemetryCollector(dataDir string) *TelemetryCollector {
	tc := &TelemetryCollector{
		dataDir:       dataDir,
		events:        make(map[string][]map[string]interface{}),
		stats:         make(map[string]int),
		enabled:       true, // Always enabled for now
		flushInterval: 5 * time.Minute,
		stopCh:        make(chan struct{}),
	}

	// Start background flushing if enabled
	if tc.enabled {
		go tc.flushLoop()
	}

	return tc
}

// RecordEvent records a telemetry event
func (tc *TelemetryCollector) RecordEvent(eventType string, data map[string]interface{}) {
	if !tc.enabled {
		return
	}

	// Add timestamp to event data
	eventData := make(map[string]interface{})
	for k, v := range data {
		eventData[k] = v
	}
	eventData["timestamp"] = time.Now().Unix()

	// Store event
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if _, ok := tc.events[eventType]; !ok {
		tc.events[eventType] = make([]map[string]interface{}, 0)
	}
	tc.events[eventType] = append(tc.events[eventType], eventData)
	tc.stats[eventType]++
}

// GetStats returns the current telemetry stats
func (tc *TelemetryCollector) GetStats() (map[string]int, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// Create a copy of the stats map
	stats := make(map[string]int, len(tc.stats))
	for k, v := range tc.stats {
		stats[k] = v
	}

	return stats, nil
}

// GetEvents returns the events of a specific type
func (tc *TelemetryCollector) GetEvents(eventType string, limit int) ([]map[string]interface{}, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	events, ok := tc.events[eventType]
	if !ok {
		return []map[string]interface{}{}, nil
	}

	// Return the most recent events up to the limit
	start := 0
	if len(events) > limit {
		start = len(events) - limit
	}

	result := make([]map[string]interface{}, len(events)-start)
	copy(result, events[start:])

	return result, nil
}

// Flush writes telemetry data to disk
func (tc *TelemetryCollector) Flush() error {
	if !tc.enabled {
		return nil
	}

	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// Create telemetry directory if it doesn't exist
	telemetryDir := filepath.Join(tc.dataDir, "telemetry")
	if err := os.MkdirAll(telemetryDir, 0755); err != nil {
		return fmt.Errorf("failed to create telemetry directory: %w", err)
	}

	// Write stats to file
	statsFile := filepath.Join(telemetryDir, "stats.json")
	statsData, err := json.MarshalIndent(tc.stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}

	if err := os.WriteFile(statsFile, statsData, 0644); err != nil {
		return fmt.Errorf("failed to write stats file: %w", err)
	}

	// Write events to files by type
	for eventType, events := range tc.events {
		eventsFile := filepath.Join(telemetryDir, fmt.Sprintf("%s.json", eventType))
		eventsData, err := json.MarshalIndent(events, "", "  ")
		if err != nil {
			fmt.Printf("Failed to marshal %s events: %v", eventType, err)
			continue
		}

		if err := os.WriteFile(eventsFile, eventsData, 0644); err != nil {
			fmt.Printf("Failed to write %s events file: %v", eventType, err)
			continue
		}
	}

	return nil
}

// flushLoop periodically flushes telemetry data to disk
func (tc *TelemetryCollector) flushLoop() {
	ticker := time.NewTicker(tc.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := tc.Flush(); err != nil {
				fmt.Printf("Failed to flush telemetry data: %v", err)
			}
		case <-tc.stopCh:
			return
		}
	}
}

// Stop stops the telemetry collector
func (tc *TelemetryCollector) Stop() error {
	if !tc.enabled {
		return nil
	}

	// Stop the flush loop
	close(tc.stopCh)

	// Flush data one last time
	return tc.Flush()
}
