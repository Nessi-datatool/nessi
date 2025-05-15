package dbt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// AlertManager manages alerts for validation results
type AlertManager struct {
	config *AlertConfig
}

// NewAlertManager creates a new alert manager
func NewAlertManager(config *AlertConfig) *AlertManager {
	return &AlertManager{
		config: config,
	}
}

// SendAlerts sends alerts for validation results
func (m *AlertManager) SendAlerts(results *ValidationResults) error {
	if !m.config.Enabled {
		return nil
	}

	// Check if there are any failures
	if !results.HasFailures() {
		return nil
	}

	// Send alerts to all configured channels
	for _, channel := range m.config.Channels {
		switch channel.Type {
		case "slack":
			if err := m.sendSlackAlert(channel, results); err != nil {
				return fmt.Errorf("failed to send Slack alert: %w", err)
			}
		case "email":
			if err := m.sendEmailAlert(channel, results); err != nil {
				return fmt.Errorf("failed to send email alert: %w", err)
			}
		default:
			return fmt.Errorf("unsupported alert channel type: %s", channel.Type)
		}
	}

	return nil
}

// sendSlackAlert sends an alert to Slack
func (m *AlertManager) sendSlackAlert(channel AlertChannel, results *ValidationResults) error {
	if channel.Webhook == "" {
		return fmt.Errorf("Slack webhook URL is required")
	}

	// Create Slack message
	message := map[string]interface{}{
		"text": fmt.Sprintf("Data Quality Alert: %d of %d rules failed (%.2f%% quality score)",
			results.Summary.FailedRules,
			results.Summary.TotalRules,
			results.Summary.QualityScore),
		"attachments": []map[string]interface{}{
			{
				"color":      "#ff0000",
				"title":      "Failed Rules",
				"title_link": "", // Could be a link to a dashboard
				"text":       m.formatFailedRules(results),
				"footer":     "Nessi.dev dbt Plugin",
				"ts":         time.Now().Unix(),
			},
		},
	}

	// Convert message to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack message: %w", err)
	}

	// Send HTTP request
	resp, err := http.Post(channel.Webhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Slack: %s", resp.Status)
	}

	return nil
}

// sendEmailAlert sends an alert via email
func (m *AlertManager) sendEmailAlert(channel AlertChannel, results *ValidationResults) error {
	if len(channel.Recipients) == 0 {
		return fmt.Errorf("email recipients are required")
	}

	// In a real implementation, you would:
	// 1. Connect to an SMTP server
	// 2. Format an email with the validation results
	// 3. Send the email to all recipients

	// For now, just log the email that would be sent
	fmt.Printf("Would send email to: %s\n", strings.Join(channel.Recipients, ", "))
	fmt.Printf("Subject: Data Quality Alert: %d of %d rules failed\n", 
		results.Summary.FailedRules, 
		results.Summary.TotalRules)
	fmt.Printf("Body: %s\n", m.formatFailedRules(results))

	return nil
}

// formatFailedRules formats failed rules for alerts
func (m *AlertManager) formatFailedRules(results *ValidationResults) string {
	var builder strings.Builder

	failureFound := false
	for _, result := range results.Results {
		if result.Status == "failed" {
			failureFound = true
			builder.WriteString(fmt.Sprintf("• Model: %s, Rule: %s\n  %s\n", 
				result.ModelName, 
				result.RuleName, 
				result.Message))
		}
	}

	if !failureFound {
		return "All data quality rules passed successfully."
	}

	return builder.String()
}

// GenerateQualityScore generates a quality score from validation results
func GenerateQualityScore(results *ValidationResults) float64 {
	if results.Summary.TotalRules == 0 {
		return 100.0
	}

	// Basic quality score: percentage of passed rules
	return float64(results.Summary.PassedRules) / float64(results.Summary.TotalRules) * 100.0
}

// QualityTrend represents a trend in quality scores over time
type QualityTrend struct {
	ModelName string      `json:"model_name"`
	Scores    []ScoreData `json:"scores"`
}

// ScoreData represents a quality score at a point in time
type ScoreData struct {
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
}

// StoreQualityScore stores a quality score for trend analysis
func StoreQualityScore(modelName string, score float64) error {
	// In a real implementation, you would:
	// 1. Connect to a database
	// 2. Store the quality score with a timestamp
	// 3. Return any errors

	// For now, just log the score
	fmt.Printf("Storing quality score for model %s: %.2f\n", modelName, score)

	return nil
}

// GetQualityTrend gets the quality trend for a model
func GetQualityTrend(modelName string, days int) (*QualityTrend, error) {
	// In a real implementation, you would:
	// 1. Connect to a database
	// 2. Retrieve quality scores for the specified time period
	// 3. Return the trend data

	// For now, just return mock data
	trend := &QualityTrend{
		ModelName: modelName,
		Scores:    make([]ScoreData, 0, days),
	}

	// Generate mock data
	now := time.Now()
	for i := 0; i < days; i++ {
		timestamp := now.AddDate(0, 0, -i)
		// Mock score between 80 and 100
		score := 80.0 + float64(i%20)
		trend.Scores = append(trend.Scores, ScoreData{
			Timestamp: timestamp,
			Score:     score,
		})
	}

	return trend, nil
}
