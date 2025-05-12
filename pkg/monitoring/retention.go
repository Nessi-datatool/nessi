package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
)

// MetricType defines the type of metric
type MetricType string

const (
	// MetricTypeGauge represents a gauge metric
	MetricTypeGauge MetricType = "gauge"
	// MetricTypeCounter represents a counter metric
	MetricTypeCounter MetricType = "counter"
	// MetricTypeHistogram represents a histogram metric
	MetricTypeHistogram MetricType = "histogram"
)

// MetricValue represents a metric value
type MetricValue struct {
	Value     float64            `json:"value"`
	Labels    map[string]string  `json:"labels"`
	Timestamp time.Time          `json:"timestamp"`
}

// MetricSeries represents a time series of metric values
type MetricSeries struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        MetricType        `json:"type"`
	Values      []MetricValue     `json:"values"`
}

// MetricSnapshot represents a snapshot of metrics at a point in time
type MetricSnapshot struct {
	Timestamp time.Time                `json:"timestamp"`
	Metrics   map[string]MetricSeries  `json:"metrics"`
}

// RetentionConfig represents the configuration for metric retention
type RetentionConfig struct {
	Enabled         bool          `json:"enabled"`
	StoragePath     string        `json:"storage_path"`
	RetentionPeriod time.Duration `json:"retention_period"`
	SnapshotInterval time.Duration `json:"snapshot_interval"`
}

// MetricRetention manages metric retention
type MetricRetention struct {
	mu              sync.RWMutex
	config          RetentionConfig
	currentSnapshot MetricSnapshot
	stopCh          chan struct{}
	running         bool
}

// NewMetricRetention creates a new MetricRetention instance
func NewMetricRetention(config RetentionConfig) *MetricRetention {
	return &MetricRetention{
		config: config,
		currentSnapshot: MetricSnapshot{
			Timestamp: time.Now(),
			Metrics:   make(map[string]MetricSeries),
		},
		stopCh:  make(chan struct{}),
		running: false,
	}
}

// Start starts the metric retention service
func (r *MetricRetention) Start() error {
	if r.running {
		return nil
	}

	// Ensure storage directory exists
	if err := os.MkdirAll(r.config.StoragePath, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Start snapshot goroutine
	go r.snapshotLoop()

	r.running = true
	return nil
}

// Stop stops the metric retention service
func (r *MetricRetention) Stop() {
	if !r.running {
		return
	}

	close(r.stopCh)
	r.running = false
}

// RecordMetric records a metric value
func (r *MetricRetention) RecordMetric(name, description string, metricType MetricType, value float64, labels map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get or create metric series
	series, ok := r.currentSnapshot.Metrics[name]
	if !ok {
		series = MetricSeries{
			Name:        name,
			Description: description,
			Type:        metricType,
			Values:      []MetricValue{},
		}
	}

	// Add value to series
	series.Values = append(series.Values, MetricValue{
		Value:     value,
		Labels:    labels,
		Timestamp: time.Now(),
	})

	// Update snapshot
	r.currentSnapshot.Metrics[name] = series
}

// GetMetricSeries returns a metric series by name
func (r *MetricRetention) GetMetricSeries(name string) (MetricSeries, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	series, ok := r.currentSnapshot.Metrics[name]
	return series, ok
}

// GetMetricSeriesWithLabels returns a metric series by name and labels
func (r *MetricRetention) GetMetricSeriesWithLabels(name string, labels map[string]string) []MetricValue {
	r.mu.RLock()
	defer r.mu.RUnlock()

	series, ok := r.currentSnapshot.Metrics[name]
	if !ok {
		return []MetricValue{}
	}

	// Filter values by labels
	var result []MetricValue
	for _, value := range series.Values {
		match := true
		for k, v := range labels {
			if value.Labels[k] != v {
				match = false
				break
			}
		}
		if match {
			result = append(result, value)
		}
	}

	return result
}

// GetMetricsSnapshot returns the current metrics snapshot
func (r *MetricRetention) GetMetricsSnapshot() MetricSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.currentSnapshot
}

// snapshotLoop periodically takes snapshots of metrics
func (r *MetricRetention) snapshotLoop() {
	ticker := time.NewTicker(r.config.SnapshotInterval)
	defer ticker.Stop()

	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			if err := r.takeSnapshot(); err != nil {
				logging.Error("Failed to take metrics snapshot", err)
			}
			if err := r.cleanupOldSnapshots(); err != nil {
				logging.Error("Failed to clean up old snapshots", err)
			}
		}
	}
}

// takeSnapshot takes a snapshot of the current metrics
func (r *MetricRetention) takeSnapshot() error {
	r.mu.RLock()
	snapshot := r.currentSnapshot
	r.mu.RUnlock()

	// Update timestamp
	snapshot.Timestamp = time.Now()

	// Create snapshot file
	filename := fmt.Sprintf("metrics_%s.json", snapshot.Timestamp.Format("20060102_150405"))
	filePath := filepath.Join(r.config.StoragePath, filename)

	// Marshal snapshot to JSON
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// Write snapshot to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot: %w", err)
	}

	logging.Info(fmt.Sprintf("Metrics snapshot saved to %s", filePath))

	// Reset current snapshot
	r.mu.Lock()
	r.currentSnapshot = MetricSnapshot{
		Timestamp: time.Now(),
		Metrics:   make(map[string]MetricSeries),
	}
	r.mu.Unlock()

	return nil
}

// cleanupOldSnapshots removes snapshots older than the retention period
func (r *MetricRetention) cleanupOldSnapshots() error {
	// Get list of snapshot files
	files, err := filepath.Glob(filepath.Join(r.config.StoragePath, "metrics_*.json"))
	if err != nil {
		return fmt.Errorf("failed to list snapshot files: %w", err)
	}

	// Calculate cutoff time
	cutoff := time.Now().Add(-r.config.RetentionPeriod)

	// Remove old snapshots
	for _, file := range files {
		// Parse timestamp from filename
		basename := filepath.Base(file)
		if len(basename) < 20 {
			continue
		}
		timestampStr := basename[8:22]
		timestamp, err := time.Parse("20060102_150405", timestampStr)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to parse timestamp from filename %s: %v", basename, err))
			continue
		}

		// Check if snapshot is older than cutoff
		if timestamp.Before(cutoff) {
			if err := os.Remove(file); err != nil {
				logging.Warn(fmt.Sprintf("Failed to remove old snapshot %s: %v", file, err))
				continue
			}
			logging.Info(fmt.Sprintf("Removed old snapshot %s", file))
		}
	}

	return nil
}

// GetMetricHistory retrieves historical metric data
func (r *MetricRetention) GetMetricHistory(name string, labels map[string]string, start, end time.Time) ([]MetricValue, error) {
	// Get list of snapshot files
	files, err := filepath.Glob(filepath.Join(r.config.StoragePath, "metrics_*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshot files: %w", err)
	}

	// Sort files by timestamp (oldest first)
	sort.Strings(files)

	// Collect metric values from snapshots
	var result []MetricValue
	for _, file := range files {
		// Parse timestamp from filename
		basename := filepath.Base(file)
		if len(basename) < 20 {
			continue
		}
		timestampStr := basename[8:22]
		timestamp, err := time.Parse("20060102_150405", timestampStr)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to parse timestamp from filename %s: %v", basename, err))
			continue
		}

		// Check if snapshot is within time range
		if timestamp.Before(start) || timestamp.After(end) {
			continue
		}

		// Read snapshot file
		data, err := os.ReadFile(file)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to read snapshot %s: %v", file, err))
			continue
		}

		// Unmarshal snapshot
		var snapshot MetricSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			logging.Warn(fmt.Sprintf("Failed to unmarshal snapshot %s: %v", file, err))
			continue
		}

		// Get metric series
		series, ok := snapshot.Metrics[name]
		if !ok {
			continue
		}

		// Filter values by labels and time range
		for _, value := range series.Values {
			if value.Timestamp.Before(start) || value.Timestamp.After(end) {
				continue
			}

			// Check if labels match
			match := true
			for k, v := range labels {
				if value.Labels[k] != v {
					match = false
					break
				}
			}
			if match {
				result = append(result, value)
			}
		}
	}

	return result, nil
}

// GetLatestMetrics retrieves the latest metrics for each series
func (r *MetricRetention) GetLatestMetrics() (map[string]map[string]MetricValue, error) {
	// Get list of snapshot files
	files, err := filepath.Glob(filepath.Join(r.config.StoragePath, "metrics_*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshot files: %w", err)
	}

	// Sort files by timestamp (newest first)
	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	// Get latest snapshot
	if len(files) == 0 {
		return make(map[string]map[string]MetricValue), nil
	}

	// Read latest snapshot
	data, err := os.ReadFile(files[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read latest snapshot: %w", err)
	}

	// Unmarshal snapshot
	var snapshot MetricSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal latest snapshot: %w", err)
	}

	// Collect latest metric values
	result := make(map[string]map[string]MetricValue)
	for name, series := range snapshot.Metrics {
		result[name] = make(map[string]MetricValue)
		for _, value := range series.Values {
			// Create a key from labels
			key := ""
			for k, v := range value.Labels {
				key += fmt.Sprintf("%s=%s;", k, v)
			}
			result[name][key] = value
		}
	}

	return result, nil
}
