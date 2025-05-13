package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SlackConfig represents the configuration for the Slack notifier
type SlackConfig struct {
	// WebhookURL is the Slack webhook URL
	WebhookURL string `json:"webhook_url"`
	
	// Channel is the Slack channel
	Channel string `json:"channel"`
	
	// Username is the username to use for the bot
	Username string `json:"username"`
	
	// IconURL is the URL to an icon to use for the bot
	IconURL string `json:"icon_url"`
	
	// IconEmoji is the emoji to use for the bot
	IconEmoji string `json:"icon_emoji"`
}

// SlackNotifier sends notifications via Slack
type SlackNotifier struct {
	config SlackConfig
}

// SlackMessage represents a Slack message
type SlackMessage struct {
	// Text is the message text
	Text string `json:"text,omitempty"`
	
	// Channel is the Slack channel
	Channel string `json:"channel,omitempty"`
	
	// Username is the username to use for the bot
	Username string `json:"username,omitempty"`
	
	// IconURL is the URL to an icon to use for the bot
	IconURL string `json:"icon_url,omitempty"`
	
	// IconEmoji is the emoji to use for the bot
	IconEmoji string `json:"icon_emoji,omitempty"`
	
	// Attachments are the message attachments
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

// SlackAttachment represents a Slack message attachment
type SlackAttachment struct {
	// Fallback is the plain text summary of the attachment
	Fallback string `json:"fallback,omitempty"`
	
	// Color is the color of the attachment
	Color string `json:"color,omitempty"`
	
	// Pretext is the text that appears above the attachment
	Pretext string `json:"pretext,omitempty"`
	
	// Title is the title of the attachment
	Title string `json:"title,omitempty"`
	
	// TitleLink is the URL for the title
	TitleLink string `json:"title_link,omitempty"`
	
	// Text is the text of the attachment
	Text string `json:"text,omitempty"`
	
	// Fields are the fields of the attachment
	Fields []SlackField `json:"fields,omitempty"`
	
	// Footer is the footer text
	Footer string `json:"footer,omitempty"`
	
	// FooterIcon is the URL to an icon for the footer
	FooterIcon string `json:"footer_icon,omitempty"`
	
	// Timestamp is the timestamp for the attachment
	Timestamp int64 `json:"ts,omitempty"`
}

// SlackField represents a field in a Slack attachment
type SlackField struct {
	// Title is the title of the field
	Title string `json:"title,omitempty"`
	
	// Value is the value of the field
	Value string `json:"value,omitempty"`
	
	// Short indicates whether the field is short
	Short bool `json:"short,omitempty"`
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(config SlackConfig) *SlackNotifier {
	return &SlackNotifier{
		config: config,
	}
}

// Name returns the name of the notifier
func (n *SlackNotifier) Name() string {
	return "slack"
}

// Send sends a notification
func (n *SlackNotifier) Send(alert *Alert, recipient string) (*AlertNotification, error) {
	// Create notification
	notification := &AlertNotification{
		ID:        uuid.New().String(),
		AlertID:   alert.ID,
		Channel:   n.Name(),
		Recipient: recipient,
		SentAt:    time.Now(),
		Status:    "sending",
	}
	
	// Determine color based on severity
	color := "#2196f3" // info (blue)
	if alert.Severity == SeverityWarning {
		color = "#ff9800" // warning (orange)
	} else if alert.Severity == SeverityCritical {
		color = "#f44336" // critical (red)
	}
	
	// Create fields
	fields := []SlackField{
		{
			Title: "Severity",
			Value: string(alert.Severity),
			Short: true,
		},
		{
			Title: "Type",
			Value: string(alert.Type),
			Short: true,
		},
		{
			Title: "Source",
			Value: alert.Source,
			Short: true,
		},
		{
			Title: "Status",
			Value: string(alert.Status),
			Short: true,
		},
	}
	
	// Add value and threshold if present
	if alert.Value != 0 || alert.Threshold != 0 {
		fields = append(fields, SlackField{
			Title: "Value",
			Value: fmt.Sprintf("%.2f", alert.Value),
			Short: true,
		})
		
		fields = append(fields, SlackField{
			Title: "Threshold",
			Value: fmt.Sprintf("%.2f", alert.Threshold),
			Short: true,
		})
		
		if alert.ComparisonOperator != "" {
			fields = append(fields, SlackField{
				Title: "Comparison",
				Value: alert.ComparisonOperator,
				Short: true,
			})
		}
	}
	
	// Add labels if present
	if len(alert.Labels) > 0 {
		labelText := ""
		for k, v := range alert.Labels {
			labelText += fmt.Sprintf("• %s: %s\n", k, v)
		}
		
		fields = append(fields, SlackField{
			Title: "Labels",
			Value: labelText,
			Short: false,
		})
	}
	
	// Create attachment
	attachment := SlackAttachment{
		Fallback:   fmt.Sprintf("[%s] %s: %s", alert.Severity, alert.Type, alert.Name),
		Color:      color,
		Title:      alert.Name,
		TitleLink:  fmt.Sprintf("http://localhost:8080/alerts/%s", alert.ID),
		Text:       alert.Description,
		Fields:     fields,
		Footer:     "Nessi.dev - Data Quality Monitoring",
		FooterIcon: "https://example.com/nessi-icon.png",
		Timestamp:  alert.Timestamp.Unix(),
	}
	
	// Create message text based on alert severity and name
	messageText := fmt.Sprintf("Alert: %s [%s]", alert.Name, strings.ToUpper(string(alert.Severity)))
	
	// Create message
	message := SlackMessage{
		Text:        messageText,
		Channel:     recipient, // Use recipient as channel override
		Username:    n.config.Username,
		IconURL:     n.config.IconURL,
		IconEmoji:   n.config.IconEmoji,
		Attachments: []SlackAttachment{attachment},
	}
	
	// Marshal message
	messageData, err := json.Marshal(message)
	if err != nil {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to marshal message: %v", err)
		return notification, fmt.Errorf("Failed to marshal message: %w", err)
	}
	
	// Send message
	resp, err := http.Post(n.config.WebhookURL, "application/json", bytes.NewBuffer(messageData))
	if err != nil {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to send message: %v", err)
		return notification, fmt.Errorf("Failed to send message: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response
	if resp.StatusCode != http.StatusOK {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to send message: %s", resp.Status)
		return notification, fmt.Errorf("Failed to send message: %s", resp.Status)
	}
	
	notification.Status = "sent"
	return notification, nil
}
