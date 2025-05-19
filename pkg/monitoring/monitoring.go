package monitoring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/nessi-dev/nessi/pkg/logging"
	"github.com/nessi-dev/nessi/pkg/security"
)

// Config represents the monitoring configuration
type Config struct {
	Alerts struct {
		Enabled           bool
		Thresholds        map[string]AlertThreshold
		SilencePeriod     time.Duration
		CooldownPeriod    time.Duration
		CheckInterval     time.Duration
		NotificationChannels []string `json:"notification_channels"`
	} `json:"alerts"`
	Metrics struct {
		Enabled bool
		Port    int
	} `json:"metrics"`
	Notifications struct {
		Slack *SlackNotificationConfig
		Email *EmailNotificationConfig
		Webhook *WebhookNotificationConfig
	} `json:"notifications"`
	Retention struct {
		Enabled          bool          `json:"enabled"`
		StoragePath      string        `json:"storage_path"`
		RetentionPeriod  time.Duration `json:"retention_period"`
		SnapshotInterval time.Duration `json:"snapshot_interval"`
	} `json:"retention"`
	Security struct {
		Auth security.AuthConfig `json:"auth"`
		SSL  security.SSLConfig  `json:"ssl"`
	} `json:"security"`
}

// AlertThreshold represents threshold values
type AlertThreshold struct {
	Warning  float64
	Critical float64
}

// AlertState represents alert states
const (
	AlertStateOK        = "OK"
	AlertStateWarning   = "WARNING"
	AlertStateCritical  = "CRITICAL"
)

// AlertConfig represents alert configuration
type AlertConfig struct {
	Thresholds        map[string]AlertThreshold
	SilencePeriod     time.Duration
	NotificationChannels []string
}

// Alert represents a monitoring alert
type Alert struct {
	Name        string
	Severity    string
	Message     string
	Timestamp   time.Time
	Metadata    map[string]string
}

// Monitor struct is defined in monitor.go

// NewMonitoringSystem creates a new Monitor instance
func NewMonitoringSystem() (*Monitor, error) {
	// Load configuration
	config, err := loadConfig("./config/monitoring.json")
	if err != nil {
		logging.Error("Failed to load monitoring config", err)
		config = defaultConfig()
	}


	// Create monitor
	m := &Monitor{
		metricsPort: config.Metrics.Port,
		config: config,
	}

	// Initialize metric store
	m.metricStore = NewMetricStore()

	// Initialize metric retention if enabled
	if config.Retention.Enabled {
		retentionConfig := RetentionConfig{
			Enabled: true,
			StoragePath: config.Retention.StoragePath,
			RetentionPeriod: config.Retention.RetentionPeriod,
			SnapshotInterval: config.Retention.SnapshotInterval,
		}
		m.metricRetention = NewMetricRetention(retentionConfig)
	}


	// Initialize security components if enabled
	if config.Security.Auth.Enabled {
		m.authManager, err = security.NewAuthManager(config.Security.Auth)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize auth manager: %w", err)
		}
		logging.Info("Authentication system initialized")
	}

	if config.Security.SSL.Enabled {
		m.certManager = security.NewCertManager(config.Security.SSL)
		logging.Info("SSL certificate manager initialized")
	}

	// Register metrics with Prometheus
	m.registerMetrics()

	return m, nil
}

// updateConfig updates the monitoring configuration
func (m *Monitor) updateConfig(config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update metrics config
	m.config.Metrics = config.Metrics


	// Update retention config
	if m.metricRetention != nil && config.Retention.Enabled {
		// Stop existing retention service if running
		m.metricRetention.Stop()
		
		// Create new retention service with updated config
		retentionConfig := RetentionConfig{
			Enabled:          config.Retention.Enabled,
			StoragePath:      config.Retention.StoragePath,
			RetentionPeriod:  config.Retention.RetentionPeriod,
			SnapshotInterval: config.Retention.SnapshotInterval,
		}
		m.metricRetention = NewMetricRetention(retentionConfig)
		
		// Start retention service if monitor is running
		if m.running {
			if err := m.metricRetention.Start(); err != nil {
				logging.Error("Failed to start metric retention service", err)
			}
		}
	}

	return nil
}

// LoadConfig loads the monitoring configuration
func (m *Monitor) LoadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load config from file
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	// Parse config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Update config
	m.config = &config

	// Initialize metric retention if enabled
	if config.Retention.Enabled {
		retentionConfig := RetentionConfig{
			Enabled:          config.Retention.Enabled,
			StoragePath:      config.Retention.StoragePath,
			RetentionPeriod:  config.Retention.RetentionPeriod * time.Hour,
			SnapshotInterval: config.Retention.SnapshotInterval * time.Second,
		}
		m.metricRetention = NewMetricRetention(retentionConfig)
	}

	return nil
}

// watchConfig watches for configuration changes
func (m *Monitor) watchConfig() error {
	// Watch config file for changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logging.Error("Failed to create watcher", err)
		return err
	}
	defer watcher.Close()

	// Add config file to watch
	if err := watcher.Add(m.configPath); err != nil {
		logging.Error("Failed to watch config file", err)
		return err
	}

	// Process events
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if err := m.LoadConfig(); err != nil {
					logging.Error("Failed to reload config", err)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			logging.Error("Watcher error", err)
		}
	}
}


// Metrics represents various monitoring metrics
type Metrics struct {
	TableSize     *prometheus.GaugeVec
	RecordCount   *prometheus.GaugeVec
	ErrorCount    *prometheus.CounterVec
	Latency       *prometheus.HistogramVec
	RuleViolations *prometheus.CounterVec
}

// SetMetricsPort sets the metrics server port
func (m *Monitor) SetMetricsPort(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.config.Metrics.Port = port
}

// SetConfigPath sets the configuration file path
func (m *Monitor) SetConfigPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configPath = path
}

// GetTableErrorRate returns the error rate for a table
func (m *Monitor) GetTableErrorRate(tableName string) float64 {
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return 0
	}

	var errorCount float64
	for _, metric := range metrics {
		if metric.GetName() == "nessi_error_count_total" {
			for _, mf := range metric.GetMetric() {
				for _, label := range mf.GetLabel() {
					if label.GetName() == "table_name" && label.GetValue() == tableName {
						errorCount = mf.GetCounter().GetValue()
					}
				}
			}
			break
		}
	}

	return errorCount
}

// RecordTableMetrics records metrics for a Delta table
func (m *Monitor) RecordTableMetrics(tableName string, records []map[string]interface{}) {
	if !m.config.Metrics.Enabled {
		return
	}

	var size float64
	var count int
	for _, record := range records {
		if sizeVal, ok := record["size"]; ok {
			size += sizeVal.(float64)
		}
		count++
	}

	// Record metrics using the metricStore
	m.metricStore.RecordMetric("table_size_"+tableName, size)
	m.metricStore.RecordMetric("record_count_"+tableName, float64(count))
	
	// Record metrics in retention system if enabled
	if m.config.Retention.Enabled && m.metricRetention != nil {
		labels := map[string]string{"table_name": tableName}
		m.metricRetention.RecordMetric("table_size", "Size of table in bytes", MetricTypeGauge, size, labels)
		m.metricRetention.RecordMetric("record_count", "Number of records in table", MetricTypeGauge, float64(count), labels)
	}
}

// RecordError records an error metric
func (m *Monitor) RecordError(tableName, errorType string) {
	m.metricStore.RecordMetric("error_count_"+tableName, 1.0)
	
	// Record error in retention system if enabled
	if m.config.Retention.Enabled && m.metricRetention != nil {
		labels := map[string]string{
			"table_name": tableName,
			"error_type": errorType,
		}
		m.metricRetention.RecordMetric("error_count", "Number of errors", MetricTypeCounter, 1.0, labels)
	}
}

// RecordLatency records operation latency
func (m *Monitor) RecordLatency(operation string, duration time.Duration) {
	m.metricStore.RecordMetric("latency_"+operation, duration.Seconds())
	
	// Record latency in retention system if enabled
	if m.config.Retention.Enabled && m.metricRetention != nil {
		labels := map[string]string{"operation": operation}
		m.metricRetention.RecordMetric("operation_latency", "Operation latency in seconds", MetricTypeHistogram, duration.Seconds(), labels)
	}
}

// RecordRuleViolation records a rule violation
func (m *Monitor) RecordRuleViolation(tableName, ruleName string) {
	if !m.config.Metrics.Enabled {
		return
	}

	m.metricStore.RecordMetric("rule_violations_"+ruleName+"_"+tableName, 1.0)
	
	// Record rule violation in retention system if enabled
	if m.config.Retention.Enabled && m.metricRetention != nil {
		labels := map[string]string{
			"table_name": tableName,
			"rule_name": ruleName,
		}
		m.metricRetention.RecordMetric("rule_violations", "Number of rule violations", MetricTypeCounter, 1.0, labels)
	}
}

// Start starts the monitoring service
func (m *Monitor) Start() {
	if m.running {
		return
	}

	if err := m.LoadConfig(); err != nil {
		logging.Error("Failed to load config", err)
		return
	}

	// Create router and set up routes
	router := http.NewServeMux()

	// Set up handlers
	// Use the default Prometheus registry
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/health", m.healthCheck)

	// Apply authentication middleware if enabled
	if m.authManager != nil {
		// Protected routes
		router.Handle("/metrics/history", m.authManager.AuthMiddleware(http.HandlerFunc(m.metricsHistory)))
		
		// Add authentication endpoints
		router.HandleFunc("/auth/login", m.handleLogin)
		router.HandleFunc("/auth/refresh", m.handleRefreshToken)
		
		// Admin routes
		router.Handle("/admin/users", m.authManager.RoleMiddleware(security.RoleAdmin)(http.HandlerFunc(m.handleUsers)))
	} else {
		// No authentication, routes are public
		router.HandleFunc("/metrics/history", m.metricsHistory)
	}

	// Start metrics server
	if m.config.Metrics.Enabled {
		go func() {
			addr := fmt.Sprintf(":%d", m.config.Metrics.Port)
			logging.Info(fmt.Sprintf("Starting metrics server on port %d", m.config.Metrics.Port))
			
			// Use SSL if enabled
			if m.certManager != nil && m.config.Security.SSL.Enabled {
				logging.Info("Starting HTTPS server")
				if err := m.certManager.StartHTTPSServer(addr, router); err != nil {
					logging.Error("Failed to start HTTPS server", err)
				}
			} else {
				// Fallback to HTTP
				if err := http.ListenAndServe(addr, router); err != nil {
					logging.Error("Failed to start HTTP server", err)
				}
			}
		}()
	}

	// Start alerting if enabled
	if m.config.Alerts.Enabled {
		go m.StartAlerting()
	}

	// Start metric retention if enabled
	if m.config.Retention.Enabled && m.metricRetention != nil {
		if err := m.metricRetention.Start(); err != nil {
			logging.Error("Failed to start metric retention service", err)
		}
		logging.Info("Started metric retention service")
	}

	m.running = true
}

// monitor is the monitoring loop
func (m *Monitor) monitor() {
	// Start monitoring loop
	lastCheck := time.Now()
	checkInterval := time.Second * 5

	for {
		select {
		case <-m.stopCh:
			return
		default:
			if time.Since(lastCheck) >= checkInterval {
				m.checkTables()
				lastCheck = time.Now()
			}
		}
	}
}

// healthCheck handles health check requests
func (m *Monitor) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}


// metricsHistory handles metrics history requests
func (m *Monitor) metricsHistory(w http.ResponseWriter, r *http.Request) {
	if !m.config.Retention.Enabled || m.metricRetention == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Metric retention is not enabled"))
		return
	}
	
	// Parse query parameters
	query := r.URL.Query()
	metricName := query.Get("metric")
	if metricName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing metric parameter"))
		return
	}
	
	// Parse time range
	start := time.Now().Add(-24 * time.Hour) // Default to last 24 hours
	end := time.Now()
	
	if startStr := query.Get("start"); startStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = startTime
		}
	}
	
	if endStr := query.Get("end"); endStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = endTime
		}
	}
	
	// Parse labels
	labels := make(map[string]string)
	for key, values := range query {
		if key != "metric" && key != "start" && key != "end" && len(values) > 0 {
			labels[key] = values[0]
		}
	}
	
	// Get metric history
	history, err := m.metricRetention.GetMetricHistory(metricName, labels, start, end)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to get metric history: %v", err)))
		return
	}
	
	// Return history as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// StartAlerting starts the alert monitoring loop
func (m *Monitor) StartAlerting() {
	go m.monitor()
}

// checkTables checks tables
func (m *Monitor) checkTables() {
	// TODO: Implement actual table checking
	// For now, just record some mock metrics
	m.RecordTableMetrics("users", []map[string]interface{}{{"size": 100, "count": 1000}})
	m.RecordTableMetrics("orders", []map[string]interface{}{{"size": 50, "count": 500}})
}

// processAlerts processes alerts
func (m *Monitor) processAlerts(alerts []Alert, thresholds map[string]AlertThreshold) {
	// TODO: Implement actual alert processing
	logging.Info("Processing alerts")
}
