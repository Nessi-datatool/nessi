package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
)

// ReportMetricsCollector collects and exports metrics for reports
type ReportMetricsCollector struct {
	mu             sync.RWMutex
	exportPath     string
	exportInterval time.Duration
	enableSystem   bool
	metricsStore   map[string]interface{}
	stopCh         chan struct{}
	running        bool
}

// NewReportMetricsCollector creates a new metrics collector for reports
func NewReportMetricsCollector() *ReportMetricsCollector {
	return &ReportMetricsCollector{
		exportPath:     "./data/metrics",
		exportInterval: time.Minute * 5,
		metricsStore:   make(map[string]interface{}),
		stopCh:         make(chan struct{}),
	}
}

// SetMetricsExportPath sets the path where metrics will be exported
func (c *ReportMetricsCollector) SetMetricsExportPath(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.exportPath = path
}

// SetExportInterval sets the interval at which metrics are exported
func (c *ReportMetricsCollector) SetExportInterval(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.exportInterval = interval
}

// EnableSystemMetrics enables collection of system metrics
func (c *ReportMetricsCollector) EnableSystemMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.enableSystem = true
}

// RecordMetrics records metrics for a specific category
func (c *ReportMetricsCollector) RecordMetrics(category string, metrics map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Add timestamp to metrics
	metrics["timestamp"] = time.Now().Format(time.RFC3339)

	// Store metrics
	c.metricsStore[category] = metrics
}

// Start starts the metrics collection and export process
func (c *ReportMetricsCollector) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return
	}

	// Create export directory if it doesn't exist
	if err := os.MkdirAll(c.exportPath, 0755); err != nil {
		logging.Error("Failed to create metrics export directory", err)
	}

	// Start export loop
	go c.exportLoop()

	c.running = true
}

// Stop stops the metrics collection and export process
func (c *ReportMetricsCollector) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	// Signal export loop to stop
	close(c.stopCh)
	c.running = false
}

// exportLoop periodically exports metrics to disk
func (c *ReportMetricsCollector) exportLoop() {
	ticker := time.NewTicker(c.exportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.exportMetrics()
		case <-c.stopCh:
			return
		}
	}
}

// exportMetrics exports metrics to disk
func (c *ReportMetricsCollector) exportMetrics() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Create a timestamp for the export
	timestamp := time.Now().Format("20060102-150405")

	// Export each category of metrics
	for category, metrics := range c.metricsStore {
		// Create filename
		filename := filepath.Join(c.exportPath, fmt.Sprintf("%s-%s.json", category, timestamp))

		// Marshal metrics to JSON
		data, err := json.MarshalIndent(metrics, "", "  ")
		if err != nil {
			logging.Error("Failed to marshal metrics", err)
			continue
		}

		// Write to file
		if err := os.WriteFile(filename, data, 0644); err != nil {
			logging.Error("Failed to write metrics to file", err)
		}
	}
}
