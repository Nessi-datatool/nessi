package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"strings"
	"net/http"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"bytes"
	"net/smtp"
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

// Monitor represents a monitoring instance
type Monitor struct {
	mu               sync.Mutex
	metrics          *Metrics
	alerts           chan Alert
	running          bool
	stopCh           chan struct{}
	collector        *prometheus.Registry
	configPath       string
	config           *Config
	alertThresholds  map[string]AlertThreshold
	silencePeriod    time.Duration
	cooldownPeriod   time.Duration
	slackConfig      *SlackNotificationConfig
	emailConfig      *EmailNotificationConfig
	webhookConfig    *WebhookNotificationConfig
	lastAlertTimes   map[string]time.Time
}

// NewMetrics creates a new Metrics instance
func NewMetrics(reg *prometheus.Registry) *Metrics {
	return &Metrics{
		TableSize: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "nessi_table_size_bytes",
				Help: "Size of Delta tables in bytes",
			},
			[]string{"table_name"},
		),
		RecordCount: promauto.With(reg).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "nessi_record_count",
				Help: "Number of records in Delta tables",
			},
			[]string{"table_name"},
		),
		ErrorCount: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_errors_total",
				Help: "Total number of errors",
			},
			[]string{"error_type", "table_name"},
		),
		Latency: promauto.With(reg).NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "nessi_operation_latency_seconds",
				Help: "Latency of operations",
				Buckets: []float64{0.001, 0.01, 0.1, 1.0, 10.0},
			},
			[]string{"operation"},
		),
		RuleViolations: promauto.With(reg).NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_rule_violations_total",
				Help: "Total number of rule violations",
			},
			[]string{"rule_name", "table_name"},
		),
	}
}

// updateConfig updates the monitoring configuration
func (m *Monitor) updateConfig(config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Update metrics config
	m.config.Metrics = config.Metrics

	// Update alert thresholds
	m.alertThresholds = config.Alerts.Thresholds
	m.silencePeriod = config.Alerts.SilencePeriod
	m.cooldownPeriod = config.Alerts.CooldownPeriod

	// Update notification configs
	m.slackConfig = config.Notifications.Slack
	m.emailConfig = config.Notifications.Email
	m.webhookConfig = config.Notifications.Webhook

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

	// Unmarshal config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	// Update config
	m.config = &config

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

// checkAlerts checks for alert conditions
func (m *Monitor) checkAlerts() error {
	// Get current metrics
	var tableSize, recordCount float64
	
	// Get table size
	metrics, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return fmt.Errorf("failed to gather metrics: %w", err)
	}

	// Find table size
	for _, metric := range metrics {
		if metric.GetName() == "nessi_table_size_bytes" {
			for _, mf := range metric.GetMetric() {
				for _, label := range mf.GetLabel() {
					if label.GetName() == "table_name" && label.GetValue() == "users" {
						tableSize = mf.GetGauge().GetValue()
					}
				}
			}
			break
		}
		if metric.GetName() == "nessi_record_count" {
			for _, mf := range metric.GetMetric() {
				for _, label := range mf.GetLabel() {
					if label.GetName() == "table_name" && label.GetValue() == "users" {
						recordCount = mf.GetGauge().GetValue()
					}
				}
			}
			break
		}
	}

	errorRate := m.getErrorRate()

	// Check each alert condition
	for alertName, threshold := range m.alertThresholds {
		var value float64
		switch alertName {
		case "table_size":
			value = tableSize
		case "record_count":
			value = recordCount
		case "error_rate":
			value = errorRate
		}

		// Check against thresholds
		if value > threshold.Critical {
			m.SendAlert(alertName, "users", value)
		} else if value > threshold.Warning {
			m.SendAlert(alertName, "users", value)
		}
	}
	return nil
}

// SendAlert sends an alert notification
func (m *Monitor) SendAlert(name string, tableName string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Alerts.Enabled {
		return
	}

	threshold, ok := m.alertThresholds[name]
	if !ok {
		return
	}

	// Check cooldown period
	lastAlertTime, ok := m.lastAlertTimes[name]
	if ok && time.Since(lastAlertTime) < m.cooldownPeriod {
		return
	}

	// Determine severity
	severity := AlertStateOK
	if value >= threshold.Critical {
		severity = AlertStateCritical
	} else if value >= threshold.Warning {
		severity = AlertStateWarning
	} else {
		return // No alert needed
	}

	alert := Alert{
		Name:      name,
		Severity:  severity,
		Message:   fmt.Sprintf("Metric %s exceeded threshold for table %s (value: %.2f)", name, tableName, value),
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"table_name": tableName,
		},
	}

	// Send alert through configured channels
	m.alerts <- alert
	m.lastAlertTimes[name] = time.Now()
}

// Metrics represents various monitoring metrics
type Metrics struct {
	TableSize     *prometheus.GaugeVec
	RecordCount   *prometheus.GaugeVec
	ErrorCount    *prometheus.CounterVec
	Latency       *prometheus.HistogramVec
	RuleViolations *prometheus.CounterVec
}

// New creates a new Monitor instance
func New() *Monitor {
	collector := prometheus.NewRegistry()
	metrics := &Metrics{
		TableSize: promauto.With(collector).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "nessi_table_size_bytes",
				Help: "Size of Delta tables in bytes",
			},
			[]string{"table_name"},
		),
		RecordCount: promauto.With(collector).NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "nessi_record_count",
				Help: "Number of records in Delta tables",
			},
			[]string{"table_name"},
		),
		ErrorCount: promauto.With(collector).NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_error_count_total",
				Help: "Total number of errors",
			},
			[]string{"error_type", "table_name"},
		),
		Latency: promauto.With(collector).NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "nessi_operation_latency_seconds",
				Help: "Latency of operations in seconds",
				Buckets: prometheus.ExponentialBuckets(0.001, 2, 10),
			},
			[]string{"operation"},
		),
		RuleViolations: promauto.With(collector).NewCounterVec(
			prometheus.CounterOpts{
				Name: "nessi_rule_violations_total",
				Help: "Total number of rule violations",
			},
			[]string{"rule_name", "table_name"},
		),
	}

	return metrics
}

// New creates a new Monitor instance
func New() *Monitor {
	collector := prometheus.NewRegistry()
	metrics := NewMetrics(collector)
	m := &Monitor{
		metrics:         metrics,
		alerts:          make(chan Alert, 100),
		collector:       collector,
		configPath:      "pkg/monitoring/config/monitoring.json",
		alertThresholds: make(map[string]AlertThreshold),
		lastAlertTimes:  make(map[string]time.Time),
		running:        false,
		stopCh:         make(chan struct{}),
		config:         &Config{},
	}
	return m
}

// GetErrorRate returns the error rate for a table
func (m *Monitor) GetErrorRate(tableName string) float64 {
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

	m.metrics.TableSize.WithLabelValues(tableName).Set(size)
	m.metrics.RecordCount.WithLabelValues(tableName).Set(float64(count))
}

// RecordError records an error metric
func (m *Monitor) RecordError(tableName, errorType string) {
	m.metrics.ErrorCount.WithLabelValues(errorType, tableName).Inc()
}

// RecordLatency records operation latency
func (m *Monitor) RecordLatency(operation string, duration time.Duration) {
	m.metrics.Latency.WithLabelValues(operation).Observe(duration.Seconds())
}

// RecordRuleViolation records a rule violation
func (m *Monitor) RecordRuleViolation(tableName, ruleName string) {
	if !m.config.Metrics.Enabled {
		return
	}

	m.metrics.RuleViolations.WithLabelValues(ruleName, tableName).Inc()
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

	// Start HTTP server for metrics
	if m.config.Metrics.Enabled {
		go func() {
			http.Handle("/metrics", promhttp.HandlerFor(m.collector, promhttp.HandlerOpts{}))
			http.HandleFunc("/health", m.healthCheck)
			http.HandleFunc("/alerts", m.listAlerts)
			logging.Info(fmt.Sprintf("Starting metrics server on port %d", m.config.Metrics.Port))
			if err := http.ListenAndServe(fmt.Sprintf(":%d", m.config.Metrics.Port), nil); err != nil {
				logging.Error("Failed to start metrics server", err)
			}
		}()
	}

	// Start alerting if enabled
	if m.config.Alerts.Enabled {
		go m.StartAlerting()
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
				m.checkAlerts()
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

// listAlerts handles alert listing requests
func (m *Monitor) listAlerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m.alerts)
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
