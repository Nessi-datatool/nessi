package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// WebhookConfig represents the configuration for the webhook notifier
type WebhookConfig struct {
	// URL is the webhook URL
	URL string `json:"url"`
	
	// Method is the HTTP method to use (GET, POST, PUT)
	Method string `json:"method"`
	
	// Headers are the HTTP headers to include
	Headers map[string]string `json:"headers"`
	
	// Timeout is the timeout for the request
	Timeout time.Duration `json:"timeout"`
	
	// MaxRetries is the maximum number of retries
	MaxRetries int `json:"max_retries"`
	
	// RetryInterval is the interval between retries
	RetryInterval time.Duration `json:"retry_interval"`
}

// WebhookNotifier sends notifications via webhook
type WebhookNotifier struct {
	config WebhookConfig
	client *http.Client
}

// NewWebhookNotifier creates a new webhook notifier
func NewWebhookNotifier(config WebhookConfig) *WebhookNotifier {
	// Set default method if not specified
	if config.Method == "" {
		config.Method = "POST"
	}
	
	// Set default timeout if not specified
	if config.Timeout == 0 {
		config.Timeout = 2 * time.Second
	}
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}
	
	return &WebhookNotifier{
		config: config,
		client: client,
	}
}

// Name returns the name of the notifier
func (n *WebhookNotifier) Name() string {
	return "webhook"
}

// Send sends a notification
func (n *WebhookNotifier) Send(alert *Alert, recipient string) (*AlertNotification, error) {
	// Create notification
	notification := &AlertNotification{
		ID:        uuid.New().String(),
		AlertID:   alert.ID,
		Channel:   n.Name(),
		Recipient: recipient, // recipient is used as an identifier for the webhook
		SentAt:    time.Now(),
		Status:    "sending",
	}
	
	// Create payload
	payload := map[string]interface{}{
		"alert_id":          alert.ID,
		"name":              alert.Name,
		"description":       alert.Description,
		"type":              alert.Type,
		"severity":          alert.Severity,
		"status":            alert.Status,
		"source":            alert.Source,
		"timestamp":         alert.Timestamp,
		"last_updated":      alert.LastUpdated,
		"value":             alert.Value,
		"threshold":         alert.Threshold,
		"comparison_operator": alert.ComparisonOperator,
		"labels":            alert.Labels,
		"annotations":       alert.Annotations,
		"webhook_id":        recipient,
		"notification_id":   notification.ID,
		"notification_time": notification.SentAt,
	}
	
	// Marshal payload
	payloadData, err := json.Marshal(payload)
	if err != nil {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to marshal payload: %v", err)
		return notification, fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	// Create request
	req, err := http.NewRequest(n.config.Method, n.config.URL, bytes.NewBuffer(payloadData))
	if err != nil {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to create request: %v", err)
		return notification, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Nessi-Webhook-Notifier/1.0")
	
	// Add custom headers
	for k, v := range n.config.Headers {
		req.Header.Set(k, v)
	}
	
	// Send request with retries
	var resp *http.Response
	var lastErr error
	
	for i := 0; i <= n.config.MaxRetries; i++ {
		resp, err = n.client.Do(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			break
		}
		
		if err != nil {
			lastErr = err
		} else {
			lastErr = fmt.Errorf("received status code %d", resp.StatusCode)
			resp.Body.Close()
		}
		
		if i < n.config.MaxRetries {
			time.Sleep(n.config.RetryInterval)
		}
	}
	
	if lastErr != nil {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to send webhook after %d retries: %v", n.config.MaxRetries, lastErr)
		return notification, fmt.Errorf("Failed to send webhook after %d retries: %w", n.config.MaxRetries, lastErr)
	}
	
	defer resp.Body.Close()
	
	notification.Status = "sent"
	return notification, nil
}

// WebhookPayload represents a payload for a webhook notification
type WebhookPayload struct {
	// AlertID is the ID of the alert
	AlertID string `json:"alert_id"`
	
	// Name is the name of the alert
	Name string `json:"name"`
	
	// Description is the description of the alert
	Description string `json:"description"`
	
	// Type is the type of alert
	Type string `json:"type"`
	
	// Severity is the severity of the alert
	Severity string `json:"severity"`
	
	// Status is the status of the alert
	Status string `json:"status"`
	
	// Source is the source of the alert
	Source string `json:"source"`
	
	// Timestamp is the time the alert was created
	Timestamp time.Time `json:"timestamp"`
	
	// LastUpdated is the time the alert was last updated
	LastUpdated time.Time `json:"last_updated"`
	
	// Value is the value that triggered the alert
	Value float64 `json:"value"`
	
	// Threshold is the threshold that was exceeded
	Threshold float64 `json:"threshold"`
	
	// ComparisonOperator is the comparison operator used
	ComparisonOperator string `json:"comparison_operator"`
	
	// Labels are key-value pairs for the alert
	Labels map[string]string `json:"labels"`
	
	// Annotations are additional information for the alert
	Annotations map[string]string `json:"annotations"`
	
	// WebhookID is the ID of the webhook
	WebhookID string `json:"webhook_id"`
	
	// NotificationID is the ID of the notification
	NotificationID string `json:"notification_id"`
	
	// NotificationTime is the time the notification was sent
	NotificationTime time.Time `json:"notification_time"`
}
