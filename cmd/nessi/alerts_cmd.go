package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/nessi-dev/nessi/pkg/monitoring/alerts"
)

// alertsCmd represents the alerts command
var alertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Manage alerts and notifications",
	Long: `Manage alerts and notifications for data quality issues, anomalies, and system events.
This command allows you to create, view, update, and manage alerts, as well as configure
alert rules and send notifications.`,
}

// alertsListCmd represents the alerts list command
var alertsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all alerts",
	Long:  `List all alerts in the system with optional filtering by status, type, or severity.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get filters
		status, _ := cmd.Flags().GetString("status")
		alertType, _ := cmd.Flags().GetString("type")
		severity, _ := cmd.Flags().GetString("severity")
		outputFormat, _ := cmd.Flags().GetString("output")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get alerts with filters
		var alertsList []*alerts.Alert
		if status != "" {
			alertsList = manager.GetAlertsByStatus(alerts.AlertStatus(status))
		} else if alertType != "" {
			alertsList = manager.GetAlertsByType(alerts.AlertType(alertType))
		} else if severity != "" {
			alertsList = manager.GetAlertsBySeverity(alerts.AlertSeverity(severity))
		} else {
			alertsList = manager.GetAlerts()
		}

		// Output alerts
		if len(alertsList) == 0 {
			fmt.Println("No alerts found.")
			return
		}

		if outputFormat == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(alertsList, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling alerts to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		} else {
			// Output as table
			fmt.Printf("%-36s %-30s %-10s %-10s %-20s %-20s\n", "ID", "NAME", "SEVERITY", "STATUS", "SOURCE", "TIMESTAMP")
			fmt.Println(strings.Repeat("-", 130))
			for _, a := range alertsList {
				fmt.Printf("%-36s %-30s %-10s %-10s %-20s %-20s\n",
					a.ID,
					truncateString(a.Name, 30),
					a.Severity,
					a.Status,
					truncateString(a.Source, 20),
					a.Timestamp.Format(time.RFC3339),
				)
			}
		}
	},
}

// alertsGetCmd represents the alerts get command
var alertsGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get details of a specific alert",
	Long:  `Get detailed information about a specific alert by its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert ID
		alertID := args[0]
		outputFormat, _ := cmd.Flags().GetString("output")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get alert
		alert, err := manager.GetAlert(alertID)
		if err != nil {
			fmt.Printf("Error getting alert: %v\n", err)
			os.Exit(1)
		}

		// Output alert
		if outputFormat == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(alert, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling alert to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		} else {
			// Output as formatted text
			fmt.Printf("ID:          %s\n", alert.ID)
			fmt.Printf("Name:        %s\n", alert.Name)
			fmt.Printf("Description: %s\n", alert.Description)
			fmt.Printf("Type:        %s\n", alert.Type)
			fmt.Printf("Severity:    %s\n", alert.Severity)
			fmt.Printf("Status:      %s\n", alert.Status)
			fmt.Printf("Source:      %s\n", alert.Source)
			fmt.Printf("Timestamp:   %s\n", alert.Timestamp.Format(time.RFC3339))
			fmt.Printf("Last Updated: %s\n", alert.LastUpdated.Format(time.RFC3339))

			if alert.Value != 0 || alert.Threshold != 0 {
				fmt.Printf("Value:       %.2f\n", alert.Value)
				fmt.Printf("Threshold:   %.2f\n", alert.Threshold)
				if alert.ComparisonOperator != "" {
					fmt.Printf("Comparison:  %s\n", alert.ComparisonOperator)
				}
			}

			if len(alert.Labels) > 0 {
				fmt.Println("\nLabels:")
				for k, v := range alert.Labels {
					fmt.Printf("  %s: %s\n", k, v)
				}
			}

			if len(alert.Annotations) > 0 {
				fmt.Println("\nAnnotations:")
				for k, v := range alert.Annotations {
					fmt.Printf("  %s: %s\n", k, v)
				}
			}

			if alert.Status == alerts.StatusAcknowledged && alert.AcknowledgedBy != "" {
				fmt.Printf("\nAcknowledged by: %s\n", alert.AcknowledgedBy)
				fmt.Printf("Acknowledged at: %s\n", alert.AcknowledgedAt.Format(time.RFC3339))
			}

			if alert.Status == alerts.StatusResolved && alert.ResolvedAt != nil {
				fmt.Printf("\nResolved at: %s\n", alert.ResolvedAt.Format(time.RFC3339))
			}

			if alert.Status == alerts.StatusSilenced {
				fmt.Printf("\nSilenced by: %s\n", alert.SilencedBy)
				fmt.Printf("Silence reason: %s\n", alert.SilenceReason)
				fmt.Printf("Silenced until: %s\n", alert.SilencedUntil.Format(time.RFC3339))
			}

			if len(alert.Notifications) > 0 {
				fmt.Println("\nNotifications:")
				fmt.Printf("%-36s %-10s %-30s %-20s %-10s\n", "ID", "CHANNEL", "RECIPIENT", "SENT AT", "STATUS")
				fmt.Println(strings.Repeat("-", 110))
				for _, n := range alert.Notifications {
					fmt.Printf("%-36s %-10s %-30s %-20s %-10s\n",
						n.ID,
						n.Channel,
						truncateString(n.Recipient, 30),
						n.SentAt.Format(time.RFC3339),
						n.Status,
					)
				}
			}
		}
	},
}

// alertsCreateCmd represents the alerts create command
var alertsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new alert",
	Long:  `Create a new alert with specified properties.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get alert properties
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		alertType, _ := cmd.Flags().GetString("type")
		severity, _ := cmd.Flags().GetString("severity")
		source, _ := cmd.Flags().GetString("source")
		value, _ := cmd.Flags().GetFloat64("value")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		operator, _ := cmd.Flags().GetString("operator")
		labels, _ := cmd.Flags().GetStringToString("label")
		annotations, _ := cmd.Flags().GetStringToString("annotation")

		// Validate required fields
		if name == "" || description == "" || alertType == "" || severity == "" || source == "" {
			fmt.Println("Error: name, description, type, severity, and source are required")
			os.Exit(1)
		}

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Create alert
		alert := &alerts.Alert{
			ID:                 uuid.New().String(),
			Name:               name,
			Description:        description,
			Type:               alerts.AlertType(alertType),
			Severity:           alerts.AlertSeverity(severity),
			Status:             alerts.StatusActive,
			Source:             source,
			Timestamp:          time.Now(),
			LastUpdated:        time.Now(),
			Value:              value,
			Threshold:          threshold,
			ComparisonOperator: operator,
			Labels:             labels,
			Annotations:        annotations,
		}

		// Save alert
		err = manager.CreateAlert(alert)
		if err != nil {
			fmt.Printf("Error creating alert: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert created with ID: %s\n", alert.ID)
	},
}

// Helper function to truncate strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	rootCmd.AddCommand(alertsCmd)
	alertsCmd.AddCommand(alertsListCmd)
	alertsCmd.AddCommand(alertsGetCmd)
	alertsCmd.AddCommand(alertsCreateCmd)

	// Flags for list command
	alertsListCmd.Flags().String("status", "", "Filter alerts by status (active, acknowledged, resolved, silenced)")
	alertsListCmd.Flags().String("type", "", "Filter alerts by type (quality, anomaly, system, custom)")
	alertsListCmd.Flags().String("severity", "", "Filter alerts by severity (info, warning, critical)")
	alertsListCmd.Flags().StringP("output", "o", "table", "Output format (table, json)")

	// Flags for get command
	alertsGetCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")

	// Flags for create command
	alertsCreateCmd.Flags().String("name", "", "Alert name (required)")
	alertsCreateCmd.Flags().String("description", "", "Alert description (required)")
	alertsCreateCmd.Flags().String("type", "custom", "Alert type (quality, anomaly, system, custom) (required)")
	alertsCreateCmd.Flags().String("severity", "info", "Alert severity (info, warning, critical) (required)")
	alertsCreateCmd.Flags().String("source", "", "Alert source (required)")
	alertsCreateCmd.Flags().Float64("value", 0, "Alert value")
	alertsCreateCmd.Flags().Float64("threshold", 0, "Alert threshold")
	alertsCreateCmd.Flags().String("operator", "", "Comparison operator (>, <, >=, <=, ==, !=)")
	alertsCreateCmd.Flags().StringToString("label", nil, "Alert labels (key=value)")
	alertsCreateCmd.Flags().StringToString("annotation", nil, "Alert annotations (key=value)")
}
