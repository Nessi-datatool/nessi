package monitor

import (
	"context"
	"fmt"
	"time"

	alertmanagerClient "github.com/prometheus/alertmanager/api/v2/client"
	"github.com/prometheus/alertmanager/api/v2/client/alert"
	"github.com/prometheus/alertmanager/api/v2/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"sync"
)

// AlertConfig defines alert configuration
type AlertConfig struct {
	Endpoint    string        `json:"endpoint"`
	Timeout     time.Duration `json:"timeout"`
	RetryPeriod time.Duration `json:"retry_period"`
}

// AlertManager manages alerts
type AlertManager struct {
	config    AlertConfig
	apiClient *alertmanagerClient.AlertmanagerAPI
	alerts    chan *Alert
	shutdown  chan struct{}
	shutdownWg *sync.WaitGroup
}

// Alert represents an alert to be sent
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Status      string           `json:"status"`
	GeneratorURL string          `json:"generatorURL"`
}

// NewAlertManager creates a new alert manager
func NewAlertManager(config AlertConfig) *AlertManager {
	transport := httptransport.New(config.Endpoint, alertmanagerClient.DefaultBasePath, nil)
	apiClient := alertmanagerClient.New(transport, strfmt.Default)

	return &AlertManager{
		config:    config,
		apiClient: apiClient,
		alerts:    make(chan *Alert, 100),
		shutdown:  make(chan struct{}),
		shutdownWg: &sync.WaitGroup{},
	}
}

// Start starts the alert manager
func (m *AlertManager) Start() {
	m.shutdownWg.Add(1)
	go m.processAlerts()
}

// Stop stops the alert manager
func (m *AlertManager) Stop() {
	close(m.shutdown)
	m.shutdownWg.Wait()
}

// SendAlert sends an alert
func (m *AlertManager) SendAlert(alert *Alert) {
	m.alerts <- alert
}

// processAlerts processes alerts from the channel
func (m *AlertManager) processAlerts() {
	defer m.shutdownWg.Done()

	for {
		select {
		case alertData := <-m.alerts:
			if err := m.sendAlert(alertData); err != nil {
				fmt.Printf("Error sending alert: %v\n", err)
			}
		case <-m.shutdown:
			return
		}
	}
}

// sendAlert sends a single alert using the Alertmanager V2 client API
func (m *AlertManager) sendAlert(alertData *Alert) error {
	postableAlert := models.PostableAlert{
		Annotations: models.LabelSet(alertData.Annotations),
		// Labels:      models.LabelSet(alertData.Labels), // Deferred assignment for diagnostics
	}
	// Explicitly assign Labels here to see if error message changes or persists
	postableAlert.Labels = models.LabelSet(alertData.Labels)

	if alertData.GeneratorURL != "" {
		postableAlert.GeneratorURL = strfmt.URI(alertData.GeneratorURL)
	}

	now := time.Now()
	startsAtTime := strfmt.DateTime(now) // Value type
	endsAtTime := strfmt.DateTime(now.Add(1 * time.Hour)) // Value type

	if alertData.Status == "firing" { 
		postableAlert.StartsAt = startsAtTime // DIAGNOSTIC: Assigning as value type
		postableAlert.EndsAt = endsAtTime   // DIAGNOSTIC: Assigning as value type
	} else if alertData.Status == "resolved" {
		postableAlert.StartsAt = startsAtTime // DIAGNOSTIC: Assigning as value type
		postableAlert.EndsAt = startsAtTime   // DIAGNOSTIC: Assigning as value type (resolved now)
	}

	params := alert.NewPostAlertsParams().
		WithContext(context.Background()).
		WithTimeout(m.config.Timeout)
	params.SetAlerts([]*models.PostableAlert{&postableAlert})

	_, err := m.apiClient.Alert.PostAlerts(params)
	if err != nil {
		return fmt.Errorf("failed to post alert via API: %w", err)
	}
	return nil
}

// CreateQualityAlert creates an alert for quality issues
func CreateQualityAlert(tableName string, ruleType string, severity string, description string) *Alert {
	return &Alert{
		Labels: map[string]string{
			"alertname":    "QualityViolation",
			"tablename":    tableName,
			"rule_type":    ruleType,
			"severity":     severity,
			"alertmanager": "nessi",
		},
		Annotations: map[string]string{
			"description": description,
			"summary":     fmt.Sprintf("Quality violation in table %s", tableName),
		},
		Status:        "firing",
		GeneratorURL: "http://nessi-monitoring:9090",
	}
}

// CreateErrorAlert creates an alert for system errors
func CreateErrorAlert(component string, errorType string, message string) *Alert {
	return &Alert{
		Labels: map[string]string{
			"alertname":    "SystemError",
			"component":    component,
			"error_type":   errorType,
			"alertmanager": "nessi",
		},
		Annotations: map[string]string{
			"description": message,
			"summary":     fmt.Sprintf("Error in %s component", component),
		},
		Status:        "firing",
		GeneratorURL: "http://nessi-monitoring:9090",
	}
}
