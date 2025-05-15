package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier is a simple webhook notification client
type WebhookNotifier struct {
	url         string
	method      string
	contentType string
	headers     map[string]string
	timeout     time.Duration
	client      *http.Client
}

// NewWebhookNotifier creates a new webhook notifier
func NewWebhookNotifier(url string, method string, headers map[string]string, timeout time.Duration) *WebhookNotifier {
	if method == "" {
		method = "POST"
	}
	
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	
	return &WebhookNotifier{
		url:         url,
		method:      method,
		contentType: "application/json",
		headers:     headers,
		timeout:     timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Send sends a notification to the webhook
func (n *WebhookNotifier) Send(payload interface{}) error {
	return n.SendWithContext(context.Background(), payload)
}

// SendWithContext sends a notification to the webhook with a context
func (n *WebhookNotifier) SendWithContext(ctx context.Context, payload interface{}) error {
	// Marshal payload to JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, n.method, n.url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", n.contentType)
	req.Header.Set("User-Agent", "Nessi-Webhook-Notifier/1.0")
	for key, value := range n.headers {
		req.Header.Set(key, value)
	}
	
	// Send request
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}
	
	return nil
}
