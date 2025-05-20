package monitor

import (
	"fmt"
	"sync"
	"time"
)

// SimpleManager provides a simple monitoring manager without external dependencies
type SimpleManager struct {
	metricsCollector *SimpleMetricsCollector
	config           Config
	mutex            sync.RWMutex
	isRunning        bool
}

// NewSimpleManager creates a new SimpleManager
func NewSimpleManager(config Config) *SimpleManager {
	return &SimpleManager{
		metricsCollector: NewSimpleMetricsCollector(),
		config:           config,
		isRunning:        false,
	}
}

// Start starts the monitor manager
func (m *SimpleManager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.isRunning {
		return fmt.Errorf("monitor manager is already running")
	}

	m.isRunning = true
	return nil
}

// Stop stops the monitor manager
func (m *SimpleManager) Stop() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.isRunning {
		return fmt.Errorf("monitor manager is not running")
	}

	m.isRunning = false
	return nil
}

// IsRunning returns whether the monitor manager is running
func (m *SimpleManager) IsRunning() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.isRunning
}

// GetConfig returns the monitor configuration
func (m *SimpleManager) GetConfig() Config {
	return m.config
}

// SetConfig sets the monitor configuration
func (m *SimpleManager) SetConfig(config Config) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.config = config
}

// RecordMetric records a metric
func (m *SimpleManager) RecordMetric(name, description string, value float64, labels map[string]string) {
	m.metricsCollector.recordMetric(name, description, value, labels)
}

// GetMetrics returns all metrics
func (m *SimpleManager) GetMetrics() []SimpleMetric {
	return m.metricsCollector.GetMetrics()
}

// ExportMetrics exports metrics to a file
func (m *SimpleManager) ExportMetrics(filePath string) error {
	// In a real implementation, this would export metrics to a file
	return nil
}

// RegisterMetric is a placeholder for compatibility with the old API
func (m *SimpleManager) RegisterMetric(name, metricType, description string, labels []string) error {
	// This is just a placeholder for compatibility
	return nil
}

// UpdateMetric is a placeholder for compatibility with the old API
func (m *SimpleManager) UpdateMetric(name string, value float64, labelValues map[string]string) error {
	m.RecordMetric(name, "", value, labelValues)
	return nil
}

// PushMetrics is a placeholder for compatibility with the old API
func (m *SimpleManager) PushMetrics() error {
	// This is just a placeholder for compatibility
	return nil
}

// SchedulePushMetrics is a placeholder for compatibility with the old API
func (m *SimpleManager) SchedulePushMetrics(interval time.Duration) error {
	// This is just a placeholder for compatibility
	return nil
}
