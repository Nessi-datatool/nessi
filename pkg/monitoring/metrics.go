package monitoring

// registerMetrics registers all metrics with Prometheus
func (m *Monitor) registerMetrics() {
	// Initialize metrics map if it doesn't exist
	if m.metrics == nil {
		m.metrics = make(map[string]interface{})
	}

	// Register system metrics
	m.metrics["system_cpu"] = m.metricStore.RegisterGauge("system_cpu", "System CPU usage")
	m.metrics["system_memory"] = m.metricStore.RegisterGauge("system_memory", "System memory usage")
	m.metrics["system_disk"] = m.metricStore.RegisterGauge("system_disk", "System disk usage")

	// Register application metrics
	m.metrics["request_count"] = m.metricStore.RegisterCounter("request_count", "Total number of requests")
	m.metrics["error_count"] = m.metricStore.RegisterCounter("error_count", "Total number of errors")
	m.metrics["request_duration"] = m.metricStore.RegisterHistogram("request_duration", "Request duration in seconds")
}

// getErrorRate calculates the error rate based on request count and error count
func (m *Monitor) getErrorRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	requestCount, ok := m.metrics["request_count"]
	if !ok {
		return 0
	}

	errorCount, ok := m.metrics["error_count"]
	if !ok {
		return 0
	}

	// Get the actual counter values
	requestValue := m.metricStore.GetCounterValue(requestCount)
	errorValue := m.metricStore.GetCounterValue(errorCount)

	if requestValue == 0 {
		return 0
	}

	return errorValue / requestValue
}

// GetErrorRate returns the current error rate
func (m *Monitor) GetErrorRate() float64 {
	return m.getErrorRate()
}
