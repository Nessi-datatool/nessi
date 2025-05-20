package monitoring

// This file contains legacy Monitor implementations that were causing duplicate declarations
// They are moved here and renamed to avoid conflicts with the current implementation

import (
	"sync"
	"time"

	"github.com/nessi-dev/nessi/pkg/security"
	"github.com/prometheus/client_golang/prometheus"
)

// LegacyMonitor represents a monitoring instance from the old implementation
type LegacyMonitor struct {
	mu              sync.Mutex
	metrics         *Metrics
	alerts          chan Alert
	running         bool
	stopCh          chan struct{}
	collector       *prometheus.Registry
	configPath      string
	config          *MonitoringConfig
	alertThresholds map[string]AlertThresholdConfig
	silencePeriod   time.Duration
	cooldownPeriod  time.Duration
	slackConfig     *SlackNotificationConfig
	emailConfig     *EmailNotificationConfig
	webhookConfig   *WebhookNotificationConfig
	lastAlertTimes  map[string]time.Time
	metricRetention *MetricRetention
	authManager     *security.AuthManager
	certManager     *security.CertManager
}

// NewLegacyMonitor creates a new LegacyMonitor instance
func NewLegacyMonitor() *LegacyMonitor {
	return &LegacyMonitor{
		stopCh:         make(chan struct{}),
		alerts:         make(chan Alert, 100),
		lastAlertTimes: make(map[string]time.Time),
	}
}
