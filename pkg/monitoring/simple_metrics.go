package monitoring

import (
	"sync"
	"time"
)

// SimpleMetric represents a basic metric without external dependencies
type SimpleMetric struct {
	Name        string
	Description string
	Labels      map[string]string
	Value       float64
	Timestamp   time.Time
}

// SimpleMetricStore provides a simple in-memory metric storage
type SimpleMetricStore struct {
	metrics map[string]SimpleMetric
	mutex   sync.RWMutex
}

// NewSimpleMetricStore creates a new SimpleMetricStore
func NewSimpleMetricStore() *SimpleMetricStore {
	return &SimpleMetricStore{
		metrics: make(map[string]SimpleMetric),
	}
}

// StoreMetric stores a metric in the store
func (s *SimpleMetricStore) StoreMetric(metric SimpleMetric) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create a key from name and labels
	key := metric.Name
	for k, v := range metric.Labels {
		key += "_" + k + "_" + v
	}

	s.metrics[key] = metric
}

// GetMetric retrieves a metric from the store
func (s *SimpleMetricStore) GetMetric(name string, labels map[string]string) (SimpleMetric, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Create a key from name and labels
	key := name
	for k, v := range labels {
		key += "_" + k + "_" + v
	}

	metric, found := s.metrics[key]
	return metric, found
}

// GetAllMetrics returns all metrics in the store
func (s *SimpleMetricStore) GetAllMetrics() []SimpleMetric {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	metrics := make([]SimpleMetric, 0, len(s.metrics))
	for _, metric := range s.metrics {
		metrics = append(metrics, metric)
	}

	return metrics
}

// RecordCounter records a counter metric
func (s *SimpleMetricStore) RecordCounter(name, description string, value float64, labels map[string]string) {
	s.StoreMetric(SimpleMetric{
		Name:        name,
		Description: description,
		Labels:      labels,
		Value:       value,
		Timestamp:   time.Now(),
	})
}

// RecordGauge records a gauge metric
func (s *SimpleMetricStore) RecordGauge(name, description string, value float64, labels map[string]string) {
	s.StoreMetric(SimpleMetric{
		Name:        name,
		Description: description,
		Labels:      labels,
		Value:       value,
		Timestamp:   time.Now(),
	})
}

// RecordHistogram records a histogram metric
func (s *SimpleMetricStore) RecordHistogram(name, description string, value float64, labels map[string]string) {
	// For simplicity, we just store the raw value
	// In a real implementation, we would track buckets, sum, and count
	s.StoreMetric(SimpleMetric{
		Name:        name,
		Description: description,
		Labels:      labels,
		Value:       value,
		Timestamp:   time.Now(),
	})
}
