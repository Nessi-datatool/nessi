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

// alertsRulesCmd represents the alerts rules command
var alertsRulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Manage alert rules",
	Long:  `Manage alert rules for automated alert generation based on metrics and thresholds.`,
}

// alertsRulesListCmd represents the alerts rules list command
var alertsRulesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all alert rules",
	Long:  `List all alert rules in the system with optional filtering by type or enabled status.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get filters
		ruleType, _ := cmd.Flags().GetString("type")
		enabled, _ := cmd.Flags().GetBool("enabled")
		enabledFlag := cmd.Flags().Changed("enabled")
		outputFormat, _ := cmd.Flags().GetString("output")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get rules with filters
		var rulesList []*alerts.AlertRule
		if ruleType != "" {
			rulesList = manager.GetRulesByType(alerts.AlertType(ruleType))
		} else if enabledFlag {
			if enabled {
				rulesList = manager.GetEnabledRules()
			} else {
				// Get all rules and filter out enabled ones
				allRules := manager.GetRules()
				for _, rule := range allRules {
					if !rule.Enabled {
						rulesList = append(rulesList, rule)
					}
				}
			}
		} else {
			rulesList = manager.GetRules()
		}

		// Output rules
		if len(rulesList) == 0 {
			fmt.Println("No alert rules found.")
			return
		}

		if outputFormat == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(rulesList, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling rules to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		} else {
			// Output as table
			fmt.Printf("%-36s %-30s %-10s %-10s %-15s %-10s %-10s\n", 
				"ID", "NAME", "TYPE", "SEVERITY", "METRIC", "THRESHOLD", "ENABLED")
			fmt.Println(strings.Repeat("-", 130))
			for _, r := range rulesList {
				fmt.Printf("%-36s %-30s %-10s %-10s %-15s %-10.2f %-10t\n",
					r.ID,
					truncateString(r.Name, 30),
					r.Type,
					r.Severity,
					truncateString(r.Metric, 15),
					r.Threshold,
					r.Enabled,
				)
			}
		}
	},
}

// alertsRulesGetCmd represents the alerts rules get command
var alertsRulesGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get details of a specific alert rule",
	Long:  `Get detailed information about a specific alert rule by its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get rule ID
		ruleID := args[0]
		outputFormat, _ := cmd.Flags().GetString("output")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get rule
		rule, err := manager.GetRule(ruleID)
		if err != nil {
			fmt.Printf("Error getting rule: %v\n", err)
			os.Exit(1)
		}

		// Output rule
		if outputFormat == "json" {
			// Output as JSON
			jsonData, err := json.MarshalIndent(rule, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling rule to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		} else {
			// Output as formatted text
			fmt.Printf("ID:                  %s\n", rule.ID)
			fmt.Printf("Name:                %s\n", rule.Name)
			fmt.Printf("Description:         %s\n", rule.Description)
			fmt.Printf("Type:                %s\n", rule.Type)
			fmt.Printf("Severity:            %s\n", rule.Severity)
			fmt.Printf("Source:              %s\n", rule.Source)
			fmt.Printf("Metric:              %s\n", rule.Metric)
			fmt.Printf("Threshold:           %.2f\n", rule.Threshold)
			fmt.Printf("Comparison Operator: %s\n", rule.ComparisonOperator)
			fmt.Printf("Window Size:         %s\n", rule.WindowSize)
			fmt.Printf("Evaluation Interval: %s\n", rule.EvaluationInterval)
			fmt.Printf("Enabled:             %t\n", rule.Enabled)
			fmt.Printf("Created At:          %s\n", rule.CreatedAt.Format(time.RFC3339))
			fmt.Printf("Created By:          %s\n", rule.CreatedBy)
			
			if rule.UpdatedAt != nil {
				fmt.Printf("Updated At:          %s\n", rule.UpdatedAt.Format(time.RFC3339))
			}
			if rule.UpdatedBy != "" {
				fmt.Printf("Updated By:          %s\n", rule.UpdatedBy)
			}

			if len(rule.Labels) > 0 {
				fmt.Println("\nLabels:")
				for k, v := range rule.Labels {
					fmt.Printf("  %s: %s\n", k, v)
				}
			}

			if len(rule.Annotations) > 0 {
				fmt.Println("\nAnnotations:")
				for k, v := range rule.Annotations {
					fmt.Printf("  %s: %s\n", k, v)
				}
			}

			if len(rule.NotificationChannels) > 0 {
				fmt.Println("\nNotification Channels:")
				for _, channel := range rule.NotificationChannels {
					fmt.Printf("  %s\n", channel)
				}
			}

			if len(rule.Recipients) > 0 {
				fmt.Println("\nRecipients:")
				for _, recipient := range rule.Recipients {
					fmt.Printf("  %s\n", recipient)
				}
			}
		}
	},
}

// alertsRulesCreateCmd represents the alerts rules create command
var alertsRulesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new alert rule",
	Long:  `Create a new alert rule with specified properties for automated alert generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get rule properties
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		ruleType, _ := cmd.Flags().GetString("type")
		severity, _ := cmd.Flags().GetString("severity")
		source, _ := cmd.Flags().GetString("source")
		metric, _ := cmd.Flags().GetString("metric")
		threshold, _ := cmd.Flags().GetFloat64("threshold")
		operator, _ := cmd.Flags().GetString("operator")
		windowSize, _ := cmd.Flags().GetDuration("window")
		evalInterval, _ := cmd.Flags().GetDuration("interval")
		enabled, _ := cmd.Flags().GetBool("enabled")
		user, _ := cmd.Flags().GetString("user")
		labels, _ := cmd.Flags().GetStringToString("label")
		annotations, _ := cmd.Flags().GetStringToString("annotation")
		channels, _ := cmd.Flags().GetStringSlice("channel")
		recipients, _ := cmd.Flags().GetStringSlice("recipient")

		// Validate required fields
		if name == "" || description == "" || ruleType == "" || severity == "" || 
		   source == "" || metric == "" || operator == "" {
			fmt.Println("Error: name, description, type, severity, source, metric, and operator are required")
			os.Exit(1)
		}

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Create rule
		rule := &alerts.AlertRule{
			ID:                 uuid.New().String(),
			Name:               name,
			Description:        description,
			Type:               alerts.AlertType(ruleType),
			Severity:           alerts.AlertSeverity(severity),
			Source:             source,
			Metric:             metric,
			Threshold:          threshold,
			ComparisonOperator: operator,
			WindowSize:         windowSize,
			EvaluationInterval: evalInterval,
			Labels:             labels,
			Annotations:        annotations,
			NotificationChannels: channels,
			Recipients:         recipients,
			Enabled:            enabled,
			CreatedAt:          time.Now(),
			CreatedBy:          user,
		}

		// Save rule
		err = manager.CreateRule(rule)
		if err != nil {
			fmt.Printf("Error creating rule: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert rule created with ID: %s\n", rule.ID)
	},
}

// alertsRulesUpdateCmd represents the alerts rules update command
var alertsRulesUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update an existing alert rule",
	Long:  `Update properties of an existing alert rule.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get rule ID
		ruleID := args[0]
		user, _ := cmd.Flags().GetString("user")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get existing rule
		rule, err := manager.GetRule(ruleID)
		if err != nil {
			fmt.Printf("Error getting rule: %v\n", err)
			os.Exit(1)
		}

		// Update rule properties if specified
		if cmd.Flags().Changed("name") {
			name, _ := cmd.Flags().GetString("name")
			rule.Name = name
		}
		if cmd.Flags().Changed("description") {
			description, _ := cmd.Flags().GetString("description")
			rule.Description = description
		}
		if cmd.Flags().Changed("type") {
			ruleType, _ := cmd.Flags().GetString("type")
			rule.Type = alerts.AlertType(ruleType)
		}
		if cmd.Flags().Changed("severity") {
			severity, _ := cmd.Flags().GetString("severity")
			rule.Severity = alerts.AlertSeverity(severity)
		}
		if cmd.Flags().Changed("source") {
			source, _ := cmd.Flags().GetString("source")
			rule.Source = source
		}
		if cmd.Flags().Changed("metric") {
			metric, _ := cmd.Flags().GetString("metric")
			rule.Metric = metric
		}
		if cmd.Flags().Changed("threshold") {
			threshold, _ := cmd.Flags().GetFloat64("threshold")
			rule.Threshold = threshold
		}
		if cmd.Flags().Changed("operator") {
			operator, _ := cmd.Flags().GetString("operator")
			rule.ComparisonOperator = operator
		}
		if cmd.Flags().Changed("window") {
			windowSize, _ := cmd.Flags().GetDuration("window")
			rule.WindowSize = windowSize
		}
		if cmd.Flags().Changed("interval") {
			evalInterval, _ := cmd.Flags().GetDuration("interval")
			rule.EvaluationInterval = evalInterval
		}
		if cmd.Flags().Changed("enabled") {
			enabled, _ := cmd.Flags().GetBool("enabled")
			rule.Enabled = enabled
		}
		if cmd.Flags().Changed("label") {
			labels, _ := cmd.Flags().GetStringToString("label")
			rule.Labels = labels
		}
		if cmd.Flags().Changed("annotation") {
			annotations, _ := cmd.Flags().GetStringToString("annotation")
			rule.Annotations = annotations
		}
		if cmd.Flags().Changed("channel") {
			channels, _ := cmd.Flags().GetStringSlice("channel")
			rule.NotificationChannels = channels
		}
		if cmd.Flags().Changed("recipient") {
			recipients, _ := cmd.Flags().GetStringSlice("recipient")
			rule.Recipients = recipients
		}

		// Update metadata
		now := time.Now()
		rule.UpdatedAt = &now
		rule.UpdatedBy = user

		// Save updated rule
		err = manager.UpdateRule(rule)
		if err != nil {
			fmt.Printf("Error updating rule: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert rule %s updated\n", rule.ID)
	},
}

// alertsRulesDeleteCmd represents the alerts rules delete command
var alertsRulesDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete an alert rule",
	Long:  `Permanently delete an alert rule from the system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get rule ID
		ruleID := args[0]

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Delete rule
		err = manager.DeleteRule(ruleID)
		if err != nil {
			fmt.Printf("Error deleting rule: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert rule %s deleted\n", ruleID)
	},
}

// alertsRulesEnableCmd represents the alerts rules enable command
var alertsRulesEnableCmd = &cobra.Command{
	Use:   "enable [id]",
	Short: "Enable an alert rule",
	Long:  `Enable an alert rule to start generating alerts.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get rule ID
		ruleID := args[0]
		user, _ := cmd.Flags().GetString("user")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get existing rule
		rule, err := manager.GetRule(ruleID)
		if err != nil {
			fmt.Printf("Error getting rule: %v\n", err)
			os.Exit(1)
		}

		// Check if already enabled
		if rule.Enabled {
			fmt.Printf("Alert rule %s is already enabled\n", ruleID)
			return
		}

		// Enable rule
		rule.Enabled = true
		now := time.Now()
		rule.UpdatedAt = &now
		rule.UpdatedBy = user

		// Save updated rule
		err = manager.UpdateRule(rule)
		if err != nil {
			fmt.Printf("Error enabling rule: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert rule %s enabled\n", ruleID)
	},
}

// alertsRulesDisableCmd represents the alerts rules disable command
var alertsRulesDisableCmd = &cobra.Command{
	Use:   "disable [id]",
	Short: "Disable an alert rule",
	Long:  `Disable an alert rule to stop generating alerts.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get rule ID
		ruleID := args[0]
		user, _ := cmd.Flags().GetString("user")

		// Create alert manager
		manager, err := alerts.NewAlertManager(config.DataDir)
		if err != nil {
			fmt.Printf("Error creating alert manager: %v\n", err)
			os.Exit(1)
		}

		// Get existing rule
		rule, err := manager.GetRule(ruleID)
		if err != nil {
			fmt.Printf("Error getting rule: %v\n", err)
			os.Exit(1)
		}

		// Check if already disabled
		if !rule.Enabled {
			fmt.Printf("Alert rule %s is already disabled\n", ruleID)
			return
		}

		// Disable rule
		rule.Enabled = false
		now := time.Now()
		rule.UpdatedAt = &now
		rule.UpdatedBy = user

		// Save updated rule
		err = manager.UpdateRule(rule)
		if err != nil {
			fmt.Printf("Error disabling rule: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Alert rule %s disabled\n", ruleID)
	},
}

func init() {
	alertsCmd.AddCommand(alertsRulesCmd)
	alertsRulesCmd.AddCommand(alertsRulesListCmd)
	alertsRulesCmd.AddCommand(alertsRulesGetCmd)
	alertsRulesCmd.AddCommand(alertsRulesCreateCmd)
	alertsRulesCmd.AddCommand(alertsRulesUpdateCmd)
	alertsRulesCmd.AddCommand(alertsRulesDeleteCmd)
	alertsRulesCmd.AddCommand(alertsRulesEnableCmd)
	alertsRulesCmd.AddCommand(alertsRulesDisableCmd)

	// Flags for list command
	alertsRulesListCmd.Flags().String("type", "", "Filter rules by type (quality, anomaly, system, custom)")
	alertsRulesListCmd.Flags().Bool("enabled", false, "Filter rules by enabled status")
	alertsRulesListCmd.Flags().StringP("output", "o", "table", "Output format (table, json)")

	// Flags for get command
	alertsRulesGetCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")

	// Flags for create command
	alertsRulesCreateCmd.Flags().String("name", "", "Rule name (required)")
	alertsRulesCreateCmd.Flags().String("description", "", "Rule description (required)")
	alertsRulesCreateCmd.Flags().String("type", "custom", "Rule type (quality, anomaly, system, custom) (required)")
	alertsRulesCreateCmd.Flags().String("severity", "info", "Rule severity (info, warning, critical) (required)")
	alertsRulesCreateCmd.Flags().String("source", "", "Rule source (required)")
	alertsRulesCreateCmd.Flags().String("metric", "", "Metric to monitor (required)")
	alertsRulesCreateCmd.Flags().Float64("threshold", 0, "Alert threshold (required)")
	alertsRulesCreateCmd.Flags().String("operator", ">", "Comparison operator (>, <, >=, <=, ==, !=) (required)")
	alertsRulesCreateCmd.Flags().Duration("window", 1*time.Hour, "Window size for evaluation (e.g., 1h, 30m)")
	alertsRulesCreateCmd.Flags().Duration("interval", 5*time.Minute, "Evaluation interval (e.g., 5m, 1h)")
	alertsRulesCreateCmd.Flags().Bool("enabled", true, "Whether the rule is enabled")
	alertsRulesCreateCmd.Flags().String("user", "system", "User creating the rule")
	alertsRulesCreateCmd.Flags().StringToString("label", nil, "Rule labels (key=value)")
	alertsRulesCreateCmd.Flags().StringToString("annotation", nil, "Rule annotations (key=value)")
	alertsRulesCreateCmd.Flags().StringSlice("channel", nil, "Notification channels (email, slack, webhook)")
	alertsRulesCreateCmd.Flags().StringSlice("recipient", nil, "Notification recipients")

	// Flags for update command
	alertsRulesUpdateCmd.Flags().String("name", "", "Rule name")
	alertsRulesUpdateCmd.Flags().String("description", "", "Rule description")
	alertsRulesUpdateCmd.Flags().String("type", "", "Rule type (quality, anomaly, system, custom)")
	alertsRulesUpdateCmd.Flags().String("severity", "", "Rule severity (info, warning, critical)")
	alertsRulesUpdateCmd.Flags().String("source", "", "Rule source")
	alertsRulesUpdateCmd.Flags().String("metric", "", "Metric to monitor")
	alertsRulesUpdateCmd.Flags().Float64("threshold", 0, "Alert threshold")
	alertsRulesUpdateCmd.Flags().String("operator", "", "Comparison operator (>, <, >=, <=, ==, !=)")
	alertsRulesUpdateCmd.Flags().Duration("window", 0, "Window size for evaluation (e.g., 1h, 30m)")
	alertsRulesUpdateCmd.Flags().Duration("interval", 0, "Evaluation interval (e.g., 5m, 1h)")
	alertsRulesUpdateCmd.Flags().Bool("enabled", false, "Whether the rule is enabled")
	alertsRulesUpdateCmd.Flags().String("user", "system", "User updating the rule")
	alertsRulesUpdateCmd.Flags().StringToString("label", nil, "Rule labels (key=value)")
	alertsRulesUpdateCmd.Flags().StringToString("annotation", nil, "Rule annotations (key=value)")
	alertsRulesUpdateCmd.Flags().StringSlice("channel", nil, "Notification channels (email, slack, webhook)")
	alertsRulesUpdateCmd.Flags().StringSlice("recipient", nil, "Notification recipients")

	// Flags for enable/disable commands
	alertsRulesEnableCmd.Flags().String("user", "system", "User enabling the rule")
	alertsRulesDisableCmd.Flags().String("user", "system", "User disabling the rule")
}
