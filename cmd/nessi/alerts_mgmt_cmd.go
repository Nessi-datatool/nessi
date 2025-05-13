package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
)

// alertsAckCmd represents the alerts acknowledge command
var alertsAckCmd = &cobra.Command{
	Use:     "acknowledge [id]",
	Aliases: []string{"ack"},
	Short:   "Acknowledge an alert",
	Long:    `Acknowledge an alert to indicate that it is being addressed.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert ID
		alertID := args[0]
		user, _ := cmd.Flags().GetString("user")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Acknowledge alert
		err = manager.AcknowledgeAlert(alertID, user)
		if err != nil {
			fmt.Printf("Error acknowledging alert: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert %s acknowledged by %s\n", alertID, user)
	},
}

// alertsResolveCmd represents the alerts resolve command
var alertsResolveCmd = &cobra.Command{
	Use:   "resolve [id]",
	Short: "Resolve an alert",
	Long:  `Mark an alert as resolved, indicating that the issue has been fixed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert ID
		alertID := args[0]

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Resolve alert
		err = manager.ResolveAlert(alertID)
		if err != nil {
			fmt.Printf("Error resolving alert: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert %s resolved\n", alertID)
	},
}

// alertsSilenceCmd represents the alerts silence command
var alertsSilenceCmd = &cobra.Command{
	Use:   "silence [id]",
	Short: "Silence an alert",
	Long:  `Silence an alert for a specified duration to prevent notifications.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert ID
		alertID := args[0]
		user, _ := cmd.Flags().GetString("user")
		reason, _ := cmd.Flags().GetString("reason")
		duration, _ := cmd.Flags().GetDuration("duration")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Silence alert
		err = manager.SilenceAlert(alertID, user, reason, duration)
		if err != nil {
			fmt.Printf("Error silencing alert: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert %s silenced by %s for %s\n", alertID, user, duration)
		fmt.Printf("Reason: %s\n", reason)
	},
}

// alertsDeleteCmd represents the alerts delete command
var alertsDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an alert",
	Long:  `Permanently delete an alert from the system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert ID
		alertID := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get alert to check if it exists and is not active
		alert, err := manager.GetAlert(alertID)
		if err != nil {
			fmt.Printf("Error getting alert: %v\n", err)
			os.Exit(1)
		}

		// Check if alert is active and force flag is not set
		if alert.Status == alerts.StatusActive && !force {
			fmt.Println("Cannot delete active alert. Use --force to override or resolve the alert first.")
			os.Exit(1)
		}

		// Delete alert
		err = manager.DeleteAlert(alertID)
		if err != nil {
			fmt.Printf("Error deleting alert: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert %s deleted\n", alertID)
	},
}

// alertsNotifyCmd represents the alerts notify command
var alertsNotifyCmd = &cobra.Command{
	Use:   "notify [id]",
	Short: "Send a notification for an alert",
	Long:  `Send a notification for an alert through a specified channel.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert ID
		alertID := args[0]
		channel, _ := cmd.Flags().GetString("channel")
		recipient, _ := cmd.Flags().GetString("recipient")

		// Validate required fields
		if channel == "" || recipient == "" {
			fmt.Println("Error: channel and recipient are required")
			os.Exit(1)
		}

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Configure notifiers based on channel
		switch channel {
		case "email":
			// Configure email notifier
			emailConfig := alerts.EmailConfig{
				Host:     config.SMTP.Host,
				Port:     config.SMTP.Port,
				Username: config.SMTP.Username,
				Password: config.SMTP.Password,
				From:     config.SMTP.From,
				UseHTML:  true,
			}
			emailNotifier := alerts.NewEmailNotifier(emailConfig)
			manager.RegisterNotifier(emailNotifier)
		case "slack":
			// Configure Slack notifier
			slackConfig := alerts.SlackConfig{
				WebhookURL: config.Slack.WebhookURL,
				Channel:    config.Slack.Channel,
				Username:   config.Slack.Username,
				IconEmoji:  config.Slack.IconEmoji,
			}
			slackNotifier := alerts.NewSlackNotifier(slackConfig)
			manager.RegisterNotifier(slackNotifier)
		case "webhook":
			// Configure webhook notifier
			webhookConfig := alerts.WebhookConfig{
				URL:     config.Webhook.URL,
				Method:  config.Webhook.Method,
				Headers: config.Webhook.Headers,
				Timeout: time.Duration(config.Webhook.TimeoutSeconds) * time.Second,
			}
			webhookNotifier := alerts.NewWebhookNotifier(webhookConfig)
			manager.RegisterNotifier(webhookNotifier)
		default:
			fmt.Printf("Error: unsupported notification channel: %s\n", channel)
			fmt.Println("Supported channels: email, slack, webhook")
			os.Exit(1)
		}

		// Send notification
		notification, err := manager.SendNotification(alertID, channel, recipient)
		if err != nil {
			fmt.Printf("Error sending notification: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Notification sent via %s to %s\n", channel, recipient)
		fmt.Printf("Notification ID: %s\n", notification.ID)
		fmt.Printf("Status: %s\n", notification.Status)
		if notification.ErrorMessage != "" {
			fmt.Printf("Error: %s\n", notification.ErrorMessage)
		}
	},
}

func init() {
	alertsCmd.AddCommand(alertsAckCmd)
	alertsCmd.AddCommand(alertsResolveCmd)
	alertsCmd.AddCommand(alertsSilenceCmd)
	alertsCmd.AddCommand(alertsDeleteCmd)
	alertsCmd.AddCommand(alertsNotifyCmd)

	// Flags for acknowledge command
	alertsAckCmd.Flags().String("user", "system", "User acknowledging the alert")

	// Flags for silence command
	alertsSilenceCmd.Flags().String("user", "system", "User silencing the alert")
	alertsSilenceCmd.Flags().String("reason", "Maintenance", "Reason for silencing the alert")
	alertsSilenceCmd.Flags().Duration("duration", 1*time.Hour, "Duration to silence the alert (e.g., 1h, 30m)")

	// Flags for delete command
	alertsDeleteCmd.Flags().Bool("force", false, "Force delete even if alert is active")

	// Flags for notify command
	alertsNotifyCmd.Flags().String("channel", "", "Notification channel (email, slack, webhook) (required)")
	alertsNotifyCmd.Flags().String("recipient", "", "Notification recipient (email address, Slack channel, webhook ID) (required)")
}
