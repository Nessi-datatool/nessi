package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// EmailConfig represents the configuration for the email notifier
type EmailConfig struct {
	// Host is the SMTP host
	Host string `json:"host"`

	// Port is the SMTP port
	Port int `json:"port"`

	// Username is the SMTP username
	Username string `json:"username"`

	// Password is the SMTP password
	Password string `json:"password"`

	// From is the sender email address
	From string `json:"from"`

	// FromName is the sender name
	FromName string `json:"from_name"`

	// UseSSL indicates whether to use SSL
	UseSSL bool `json:"use_ssl"`

	// UseHTML indicates whether to use HTML format
	UseHTML bool `json:"use_html"`

	// Template is a custom email template
	Template string `json:"template"`

	// TemplatePath is the path to the email template
	TemplatePath string `json:"template_path"`
}

// EmailNotifier sends notifications via email
type EmailNotifier struct {
	config   EmailConfig
	template *template.Template
}

// NewEmailNotifier creates a new email notifier
func NewEmailNotifier(config EmailConfig) (*EmailNotifier, error) {
	// Determine template content
	templateContent := defaultEmailTemplate
	if config.Template != "" {
		templateContent = config.Template
	}

	// Create template
	tmpl, err := template.New("email").Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &EmailNotifier{
		config:   config,
		template: tmpl,
	}, nil
}

// Name returns the name of the notifier
func (n *EmailNotifier) Name() string {
	return "email"
}

// Variable for testing purposes
var sendMail = smtp.SendMail

// Send sends a notification
func (n *EmailNotifier) Send(alert *Alert, recipient string) (*AlertNotification, error) {
	// Create notification
	notification := &AlertNotification{
		ID:        uuid.New().String(),
		AlertID:   alert.ID,
		Channel:   n.Name(),
		Recipient: recipient,
		SentAt:    time.Now(),
		Status:    "sending",
	}

	// Determine content type and prepare message
	var messageBody string
	contentType := "text/plain"

	if n.config.UseHTML {
		// HTML format
		contentType = "text/html"

		// Prepare email data
		data := map[string]interface{}{
			"Alert":       alert,
			"Recipient":   recipient,
			"SenderName":  n.config.FromName,
			"SenderEmail": n.config.From,
			"Timestamp":   time.Now().Format(time.RFC1123),
		}

		// Render email template
		var body bytes.Buffer
		if err := n.template.Execute(&body, data); err != nil {
			notification.Status = "failed"
			notification.ErrorMessage = fmt.Sprintf("Failed to render template: %v", err)
			return notification, fmt.Errorf("failed to render template: %w", err)
		}

		messageBody = body.String()
	} else {
		// JSON format
		alertJSON, err := json.MarshalIndent(alert, "", "  ")
		if err != nil {
			notification.Status = "failed"
			notification.ErrorMessage = fmt.Sprintf("Failed to marshal alert to JSON: %v", err)
			return notification, fmt.Errorf("failed to marshal alert to JSON: %w", err)
		}

		messageBody = string(alertJSON)
		contentType = "application/json"
	}

	// Prepare email headers
	headers := make(map[string]string)
	headers["From"] = n.config.From
	if n.config.FromName != "" {
		headers["From"] = fmt.Sprintf("%s <%s>", n.config.FromName, n.config.From)
	}
	headers["To"] = recipient
	headers["Subject"] = fmt.Sprintf("[%s] %s", strings.ToUpper(string(alert.Severity)), alert.Name)
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = contentType + "; charset=UTF-8"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + messageBody

	// Connect to SMTP server
	var auth smtp.Auth
	if n.config.Username != "" {
		auth = smtp.PlainAuth("", n.config.Username, n.config.Password, n.config.Host)
	}

	addr := fmt.Sprintf("%s:%d", n.config.Host, n.config.Port)

	// Use the sendMail function that can be mocked in tests
	if err := sendMail(addr, auth, n.config.From, []string{recipient}, []byte(message)); err != nil {
		notification.Status = "failed"
		notification.ErrorMessage = fmt.Sprintf("Failed to send email: %v", err)
		return notification, fmt.Errorf("Failed to send email: %w", err)
	}

	notification.Status = "sent"
	return notification, nil
}

// Default email template
const defaultEmailTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Nessi Alert: {{.Alert.Name}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: {{if eq .Alert.Severity "critical"}}#f44336{{else if eq .Alert.Severity "warning"}}#ff9800{{else}}#2196f3{{end}};
            color: white;
            padding: 10px 20px;
            border-radius: 5px 5px 0 0;
        }
        .content {
            border: 1px solid #ddd;
            border-top: none;
            padding: 20px;
            border-radius: 0 0 5px 5px;
        }
        .footer {
            margin-top: 20px;
            font-size: 12px;
            color: #777;
            text-align: center;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        th, td {
            padding: 8px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <div class="header">
        <h2>{{.Alert.Name}}</h2>
    </div>
    <div class="content">
        <p><strong>Description:</strong> {{.Alert.Description}}</p>
        
        <table>
            <tr>
                <th>Property</th>
                <th>Value</th>
            </tr>
            <tr>
                <td>Severity</td>
                <td>{{.Alert.Severity}}</td>
            </tr>
            <tr>
                <td>Type</td>
                <td>{{.Alert.Type}}</td>
            </tr>
            <tr>
                <td>Source</td>
                <td>{{.Alert.Source}}</td>
            </tr>
            <tr>
                <td>Timestamp</td>
                <td>{{.Alert.Timestamp}}</td>
            </tr>
            {{if .Alert.Value}}
            <tr>
                <td>Value</td>
                <td>{{.Alert.Value}}</td>
            </tr>
            {{end}}
            {{if .Alert.Threshold}}
            <tr>
                <td>Threshold</td>
                <td>{{.Alert.Threshold}}</td>
            </tr>
            {{end}}
            {{if .Alert.ComparisonOperator}}
            <tr>
                <td>Comparison</td>
                <td>{{.Alert.ComparisonOperator}}</td>
            </tr>
            {{end}}
        </table>
        
        {{if .Alert.Labels}}
        <h3>Labels</h3>
        <table>
            <tr>
                <th>Key</th>
                <th>Value</th>
            </tr>
            {{range $key, $value := .Alert.Labels}}
            <tr>
                <td>{{$key}}</td>
                <td>{{$value}}</td>
            </tr>
            {{end}}
        </table>
        {{end}}
        
        {{if .Alert.Annotations}}
        <h3>Annotations</h3>
        <table>
            <tr>
                <th>Key</th>
                <th>Value</th>
            </tr>
            {{range $key, $value := .Alert.Annotations}}
            <tr>
                <td>{{$key}}</td>
                <td>{{$value}}</td>
            </tr>
            {{end}}
        </table>
        {{end}}
        
        <p>
            <a href="http://localhost:8080/alerts/{{.Alert.ID}}">View Alert</a> | 
            <a href="http://localhost:8080/alerts/{{.Alert.ID}}/acknowledge">Acknowledge</a> | 
            <a href="http://localhost:8080/alerts/{{.Alert.ID}}/silence">Silence</a>
        </p>
    </div>
    <div class="footer">
        <p>This alert was sent by Nessi.dev at {{.Timestamp}}</p>
        <p>© 2023 Nessi.dev - Data Quality Monitoring</p>
    </div>
</body>
</html>`
