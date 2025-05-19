package telemetry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// DefaultReportInterval is the default interval for sending telemetry reports
	DefaultReportInterval = 24 * time.Hour
	
	// MinReportInterval is the minimum allowed interval for sending telemetry reports
	MinReportInterval = 1 * time.Hour
)

// Reporter is responsible for collecting and sending telemetry data
type Reporter struct {
	dataDir        string
	reportInterval time.Duration
	client         *http.Client
	stopChan       chan struct{}
	wg             sync.WaitGroup
	mutex          sync.Mutex
	running        bool
}

// NewReporter creates a new telemetry reporter
func NewReporter(dataDir string, reportInterval time.Duration) *Reporter {
	if reportInterval < MinReportInterval {
		reportInterval = DefaultReportInterval
	}
	
	return &Reporter{
		dataDir:        dataDir,
		reportInterval: reportInterval,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		stopChan: make(chan struct{}),
	}
}

// Start starts the telemetry reporter
func (r *Reporter) Start() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if r.running {
		return fmt.Errorf("reporter is already running")
	}
	
	// Check if telemetry is enabled
	enabled, err := GetStatus(r.dataDir)
	if err != nil {
		return fmt.Errorf("failed to get telemetry status: %w", err)
	}
	
	if !enabled {
		return fmt.Errorf("telemetry is disabled")
	}
	
	r.running = true
	r.wg.Add(1)
	
	go r.reportLoop()
	
	return nil
}

// Stop stops the telemetry reporter
func (r *Reporter) Stop() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if !r.running {
		return
	}
	
	close(r.stopChan)
	r.wg.Wait()
	r.running = false
}

// reportLoop periodically collects and sends telemetry data
func (r *Reporter) reportLoop() {
	defer r.wg.Done()
	
	// Send initial report
	r.sendReport()
	
	ticker := time.NewTicker(r.reportInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			r.sendReport()
		case <-r.stopChan:
			return
		}
	}
}

// sendReport collects and sends telemetry data
func (r *Reporter) sendReport() {
	// Check if telemetry is enabled
	enabled, err := GetStatus(r.dataDir)
	if err != nil || !enabled {
		return
	}
	
	// Load command stats
	commandStats, err := r.loadStats("telemetry_commands.json")
	if err != nil {
		// Don't fail if we can't load stats
		commandStats = make(map[string]int)
	}
	
	// Load feature usage stats
	featureStats, err := r.loadStats("telemetry_features.json")
	if err != nil {
		// Don't fail if we can't load stats
		featureStats = make(map[string]int)
	}
	
	// Load error stats
	errorStats, err := r.loadStats("telemetry_errors.json")
	if err != nil {
		// Don't fail if we can't load stats
		errorStats = make(map[string]int)
	}
	
	// Collect telemetry data
	data, err := CollectTelemetryData(r.dataDir, commandStats, featureStats, errorStats, nil)
	if err != nil {
		// Log error but don't fail
		fmt.Printf("Failed to collect telemetry data: %v\n", err)
		return
	}
	
	if data == nil {
		// Telemetry is disabled
		return
	}
	
	// Send telemetry data
	if err := r.sendTelemetryData(data); err != nil {
		// Log error but don't fail
		fmt.Printf("Failed to send telemetry data: %v\n", err)
		return
	}
	
	// Update last submitted time
	config, err := loadConfig(r.dataDir)
	if err != nil {
		return
	}
	
	config.LastSubmitted = time.Now()
	saveConfig(r.dataDir, config)
	
	// Reset stats after successful submission
	r.resetStats("telemetry_commands.json")
	r.resetStats("telemetry_features.json")
	r.resetStats("telemetry_errors.json")
}

// loadStats loads stats from a file
func (r *Reporter) loadStats(filename string) (map[string]int, error) {
	statsFile := filepath.Join(r.dataDir, filename)
	
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		return make(map[string]int), nil
	}
	
	data, err := os.ReadFile(statsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read stats file: %w", err)
	}
	
	var stats map[string]int
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats file: %w", err)
	}
	
	return stats, nil
}

// resetStats resets stats in a file
func (r *Reporter) resetStats(filename string) error {
	statsFile := filepath.Join(r.dataDir, filename)
	
	// Create empty stats
	stats := make(map[string]int)
	
	// Marshal stats
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal stats: %w", err)
	}
	
	// Write stats file
	if err := os.WriteFile(statsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write stats file: %w", err)
	}
	
	return nil
}

// sendTelemetryData writes telemetry data to a file
func (r *Reporter) sendTelemetryData(data *TelemetryData) error {
	// Create telemetry directory if it doesn't exist
	telemetryDir := filepath.Join(r.dataDir, "telemetry")
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
	
	fmt.Printf("Telemetry data written to %s\n", filename)
	return nil
}
