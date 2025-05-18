package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/webhook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage webhooks for Nessi events",
	Long: `Manage webhooks for Nessi events.

Webhooks allow external systems to receive notifications when specific events occur in Nessi.
Events include data validation results, alerts, rule changes, and more.`,
}

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		webhooks := service.ListWebhooks()
		if len(webhooks) == 0 {
			fmt.Println("No webhooks registered.")
			return
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(webhooks, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting output: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonOutput))
			return
		}

		// Table format (default)
		fmt.Println("ID\tNAME\tURL\tEVENTS\tENABLED")
		for _, w := range webhooks {
			events := strings.Join(w.Events, ",")
			if len(events) > 30 {
				events = events[:27] + "..."
			}
			fmt.Printf("%s\t%s\t%s\t%s\t%v\n", w.ID, w.Name, w.URL, events, w.Enabled)
		}
	},
}

var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new webhook",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")
		eventsStr, _ := cmd.Flags().GetString("events")
		headersStr, _ := cmd.Flags().GetString("headers")
		description, _ := cmd.Flags().GetString("description")

		if name == "" || url == "" || eventsStr == "" {
			fmt.Println("Error: name, url, and events are required")
			os.Exit(1)
		}

		events := strings.Split(eventsStr, ",")
		for i, event := range events {
			events[i] = strings.TrimSpace(event)
		}

		headers := make(map[string]string)
		if headersStr != "" {
			headerPairs := strings.Split(headersStr, ",")
			for _, pair := range headerPairs {
				parts := strings.SplitN(pair, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					headers[key] = value
				}
			}
		}

		webhook, err := service.CreateWebhook(name, url, events, headers, description)
		if err != nil {
			fmt.Printf("Error creating webhook: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Webhook created successfully with ID: %s\n", webhook.ID)
	},
}

var webhookGetCmd = &cobra.Command{
	Use:   "get [webhook-id]",
	Short: "Get details of a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		id := args[0]
		webhook, err := service.GetWebhook(id)
		if err != nil {
			fmt.Printf("Error getting webhook: %v\n", err)
			os.Exit(1)
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(webhook, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting output: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonOutput))
			return
		}

		// Detailed format (default)
		fmt.Println("ID:", webhook.ID)
		fmt.Println("Name:", webhook.Name)
		fmt.Println("URL:", webhook.URL)
		fmt.Println("Events:", strings.Join(webhook.Events, ", "))
		fmt.Println("Enabled:", webhook.Enabled)
		fmt.Println("Retry Count:", webhook.RetryCount)
		fmt.Println("Retry Delay:", webhook.RetryDelay, "seconds")
		fmt.Println("Description:", webhook.Description)
		fmt.Println("Created At:", webhook.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("Updated At:", webhook.UpdatedAt.Format("2006-01-02 15:04:05"))

		if len(webhook.Headers) > 0 {
			fmt.Println("Headers:")
			for key, value := range webhook.Headers {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	},
}

var webhookUpdateCmd = &cobra.Command{
	Use:   "update [webhook-id]",
	Short: "Update a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		id := args[0]
		updates := make(map[string]interface{})

		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			updates["name"] = name
		}

		if cmd.Flags().Changed("url") {
			url, _ := cmd.Flags().GetString("url")
			updates["url"] = url
		}

		if cmd.Flags().Changed("events") {
			eventsStr, _ := cmd.Flags().GetString("events")
			events := strings.Split(eventsStr, ",")
			for i, event := range events {
				events[i] = strings.TrimSpace(event)
			}
			updates["events"] = events
		}

		if cmd.Flags().Changed("headers") {
			headersStr, _ := cmd.Flags().GetString("headers")
			headers := make(map[string]string)
			headerPairs := strings.Split(headersStr, ",")
			for _, pair := range headerPairs {
				parts := strings.SplitN(pair, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					headers[key] = value
				}
			}
			updates["headers"] = headers
		}

		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			updates["description"] = description
		}

		if cmd.Flags().Changed("retry-count") {
			retryCount, _ := cmd.Flags().GetInt("retry-count")
			updates["retry_count"] = retryCount
		}

		if cmd.Flags().Changed("retry-delay") {
			retryDelay, _ := cmd.Flags().GetInt("retry-delay")
			updates["retry_delay"] = retryDelay
		}

		if cmd.Flags().Changed("enabled") {
			enabled, _ := cmd.Flags().GetBool("enabled")
			updates["enabled"] = enabled
		}

		if len(updates) == 0 {
			fmt.Println("No updates specified.")
			return
		}

		err = service.UpdateWebhook(id, updates)
		if err != nil {
			fmt.Printf("Error updating webhook: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Webhook updated successfully.")
	},
}

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete [webhook-id]",
	Short: "Delete a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		id := args[0]
		err = service.DeleteWebhook(id)
		if err != nil {
			fmt.Printf("Error deleting webhook: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Webhook deleted successfully.")
	},
}

var webhookEnableCmd = &cobra.Command{
	Use:   "enable [webhook-id]",
	Short: "Enable a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		id := args[0]
		err = service.EnableWebhook(id)
		if err != nil {
			fmt.Printf("Error enabling webhook: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Webhook enabled successfully.")
	},
}

var webhookDisableCmd = &cobra.Command{
	Use:   "disable [webhook-id]",
	Short: "Disable a webhook",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		id := args[0]
		err = service.DisableWebhook(id)
		if err != nil {
			fmt.Printf("Error disabling webhook: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Webhook disabled successfully.")
	},
}

var webhookEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "List supported event types",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		events := service.GetSupportedEvents()
		fmt.Println("Supported event types:")
		for _, event := range events {
			fmt.Println("-", event)
		}
	},
}

var webhookTestCmd = &cobra.Command{
	Use:   "test [webhook-id] [event-type]",
	Short: "Test a webhook by sending a test event",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		service, err := getWebhookService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		id := args[0]
		eventType := args[1]

		// Get webhook to verify it exists
		webhook, err := service.GetWebhook(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Check if webhook is subscribed to this event
		subscribed := false
		for _, e := range webhook.Events {
			if e == eventType || e == "*" {
				subscribed = true
				break
			}
		}

		if !subscribed {
			fmt.Printf("Warning: Webhook %s is not subscribed to event %s\n", id, eventType)
		}

		// Create test payload
		payload := map[string]interface{}{
			"test":        true,
			"message":     "This is a test event",
			"webhook_id":  id,
			"event_type":  eventType,
			"timestamp":   time.Now().Format(time.RFC3339),
		}

		// Trigger event
		results := service.TriggerEvent(eventType, payload)
		result, ok := results[id]
		if !ok {
			fmt.Println("Error: Webhook was not triggered")
			os.Exit(1)
		}

		if result.Success {
			fmt.Println("Webhook test successful!")
			fmt.Println("Status Code:", result.StatusCode)
			fmt.Println("Response:", result.ResponseBody)
		} else {
			fmt.Println("Webhook test failed!")
			fmt.Println("Error:", result.Error)
			fmt.Println("Attempt Count:", result.AttemptCount)
			if result.StatusCode > 0 {
				fmt.Println("Status Code:", result.StatusCode)
				fmt.Println("Response:", result.ResponseBody)
			}
		}
	},
}

func getWebhookService() (*webhook.WebhookService, error) {
	// Get webhook storage directory from config
	configDir := viper.GetString("config_dir")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %v", err)
		}
		configDir = filepath.Join(homeDir, ".nessi")
	}

	webhookDir := filepath.Join(configDir, "webhooks")
	storage, err := webhook.NewFileStorage(webhookDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook storage: %v", err)
	}

	service, err := webhook.NewWebhookService(storage)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook service: %v", err)
	}

	return service, nil
}

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(webhookListCmd)
	webhookCmd.AddCommand(webhookCreateCmd)
	webhookCmd.AddCommand(webhookGetCmd)
	webhookCmd.AddCommand(webhookUpdateCmd)
	webhookCmd.AddCommand(webhookDeleteCmd)
	webhookCmd.AddCommand(webhookEnableCmd)
	webhookCmd.AddCommand(webhookDisableCmd)
	webhookCmd.AddCommand(webhookEventsCmd)
	webhookCmd.AddCommand(webhookTestCmd)

	// List command flags
	webhookListCmd.Flags().String("format", "table", "Output format (table, json)")

	// Create command flags
	webhookCreateCmd.Flags().String("name", "", "Name of the webhook")
	webhookCreateCmd.Flags().String("url", "", "URL to send webhook events to")
	webhookCreateCmd.Flags().String("events", "", "Comma-separated list of event types to subscribe to")
	webhookCreateCmd.Flags().String("headers", "", "Comma-separated list of headers in key:value format")
	webhookCreateCmd.Flags().String("description", "", "Description of the webhook")

	// Get command flags
	webhookGetCmd.Flags().String("format", "detailed", "Output format (detailed, json)")

	// Update command flags
	webhookUpdateCmd.Flags().String("name", "", "Name of the webhook")
	webhookUpdateCmd.Flags().String("url", "", "URL to send webhook events to")
	webhookUpdateCmd.Flags().String("events", "", "Comma-separated list of event types to subscribe to")
	webhookUpdateCmd.Flags().String("headers", "", "Comma-separated list of headers in key:value format")
	webhookUpdateCmd.Flags().String("description", "", "Description of the webhook")
	webhookUpdateCmd.Flags().Int("retry-count", 0, "Number of retry attempts")
	webhookUpdateCmd.Flags().Int("retry-delay", 0, "Delay between retry attempts in seconds")
	webhookUpdateCmd.Flags().Bool("enabled", true, "Whether the webhook is enabled")
}
