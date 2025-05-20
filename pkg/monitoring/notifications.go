package monitoring

// SlackNotificationConfig represents Slack notification configuration
type SlackNotificationConfig struct {
	WebhookURL string `json:"webhook_url"`
	Channel    string `json:"channel"`
	Username   string `json:"username"`
}

// EmailNotificationConfig represents email notification configuration
type EmailNotificationConfig struct {
	SMTPServer string   `json:"smtp_server"`
	Port       int      `json:"port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	From       string   `json:"from"`
	To         []string `json:"to"`
}

// WebhookNotificationConfig represents webhook notification configuration
type WebhookNotificationConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}
