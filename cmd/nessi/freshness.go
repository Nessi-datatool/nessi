package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/freshness"
	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"github.com/spf13/cobra"
)

var (
	// freshnessEnabled is a flag to enable/disable freshness tracking
	freshnessEnabled bool
	
	// slaOutputFormat is the output format for SLA commands
	slaOutputFormat string
	
	// slaWarningThreshold is the warning threshold for SLA violations
	slaWarningThreshold int
	
	// slaCriticalThreshold is the critical threshold for SLA violations
	slaCriticalThreshold int
	
	// slaDescription is the description for the SLA
	slaDescription string
	
	// slaTags are the tags for the SLA
	slaTags []string
	
	// slaGrafanaDashboardURL is the Grafana dashboard URL for the SLA
	slaGrafanaDashboardURL string
)

// freshnessCmd represents the freshness command
var freshnessCmd = &cobra.Command{
	Use:   "freshness",
	Short: "Check data freshness and SLA compliance",
	Long: `Check the freshness of Delta tables and their compliance with defined SLAs.
This command allows you to monitor when tables were last updated and whether they
meet the expected update frequency defined in their SLA.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// freshnessCheckCmd represents the freshness check command
var freshnessCheckCmd = &cobra.Command{
	Use:   "check [table_name]",
	Short: "Check freshness for a table or all tables",
	Long: `Check the freshness status of a specific table or all tables with defined SLAs.
If no table name is provided, all tables with SLA configurations will be checked.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if freshness tracking is enabled
		if !freshnessEnabled {
			fmt.Println("Freshness tracking is disabled. Use --enable-freshness-tracking to enable it.")
			return
		}
		
		// Create SLA manager
		configPath := filepath.Join(configDir, "sla_config.json")
		manager, err := freshness.NewSLAManager(configPath, monitoringClient)
		if err != nil {
			logging.Error("Failed to create SLA manager", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		var statuses []*freshness.FreshnessStatus
		
		// Check freshness for a specific table or all tables
		if len(args) > 0 {
			tableName := args[0]
			status, err := manager.CheckFreshness(tableName)
			if err != nil {
				logging.Error("Failed to check freshness", err, "table", tableName)
				fmt.Printf("Error checking freshness for table %s: %v\n", tableName, err)
				return
			}
			statuses = []*freshness.FreshnessStatus{status}
		} else {
			var err error
			statuses, err = manager.CheckAllFreshness()
			if err != nil {
				logging.Error("Failed to check freshness for all tables", err)
				fmt.Printf("Error checking freshness for all tables: %v\n", err)
				return
			}
		}
		
		// Output results
		if slaOutputFormat == "json" {
			// JSON output
			output, err := json.MarshalIndent(statuses, "", "  ")
			if err != nil {
				logging.Error("Failed to marshal freshness status", err)
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println(string(output))
		} else {
			// Table output
			fmt.Println("Freshness Status:")
			fmt.Println("----------------")
			fmt.Printf("%-20s %-20s %-15s %-15s %-10s\n", "Table", "Last Update", "Time Since", "Expected", "Status")
			fmt.Printf("%-20s %-20s %-15s %-15s %-10s\n", "-----", "-----------", "----------", "--------", "------")
			
			for _, status := range statuses {
				// Format last update time
				lastUpdate := status.LastUpdateTime.Format("2006-01-02 15:04:05")
				
				// Format time since update
				timeSince := freshness.FormatDuration(status.TimeSinceUpdate)
				
				// Format expected frequency
				expected := freshness.FormatDuration(status.ExpectedFrequency)
				
				// Format status with color
				statusStr := string(status.Status)
				switch status.Status {
				case freshness.SLALevelCritical:
					statusStr = "\033[31m" + statusStr + "\033[0m" // Red
				case freshness.SLALevelWarning:
					statusStr = "\033[33m" + statusStr + "\033[0m" // Yellow
				case freshness.SLALevelInfo:
					statusStr = "\033[32m" + statusStr + "\033[0m" // Green
				}
				
				fmt.Printf("%-20s %-20s %-15s %-15s %-10s\n", status.TableName, lastUpdate, timeSince, expected, statusStr)
			}
			fmt.Println()
		}
	},
}

// slaCmd represents the SLA command
var slaCmd = &cobra.Command{
	Use:   "sla",
	Short: "Manage SLA configurations for tables",
	Long: `Manage Service Level Agreement (SLA) configurations for Delta tables.
This command allows you to define, list, and delete SLA configurations that
specify the expected update frequency for tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, show help
		cmd.Help()
	},
}

// slaListCmd represents the SLA list command
var slaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all SLA configurations",
	Long:  `List all defined SLA configurations for tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if freshness tracking is enabled
		if !freshnessEnabled {
			fmt.Println("Freshness tracking is disabled. Use --enable-freshness-tracking to enable it.")
			return
		}
		
		// Create SLA manager
		configPath := filepath.Join(configDir, "sla_config.json")
		manager, err := freshness.NewSLAManager(configPath, monitoringClient)
		if err != nil {
			logging.Error("Failed to create SLA manager", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Get all SLA configurations
		configs := manager.ListSLAs()
		
		// Output results
		if slaOutputFormat == "json" {
			// JSON output
			output, err := json.MarshalIndent(configs, "", "  ")
			if err != nil {
				logging.Error("Failed to marshal SLA configurations", err)
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Println(string(output))
		} else {
			// Table output
			fmt.Println("SLA Configurations:")
			fmt.Println("------------------")
			fmt.Printf("%-20s %-15s %-10s %-10s %-8s\n", "Table", "Frequency", "Warning", "Critical", "Enabled")
			fmt.Printf("%-20s %-15s %-10s %-10s %-8s\n", "-----", "---------", "-------", "--------", "-------")
			
			for _, config := range configs {
				// Format frequency
				frequency := freshness.FormatDuration(config.ExpectedFrequency)
				
				// Format thresholds
				warning := fmt.Sprintf("%d%%", config.WarningThreshold)
				critical := fmt.Sprintf("%d%%", config.CriticalThreshold)
				
				// Format enabled
				enabled := "Yes"
				if !config.Enabled {
					enabled = "No"
				}
				
				fmt.Printf("%-20s %-15s %-10s %-10s %-8s\n", config.TableName, frequency, warning, critical, enabled)
			}
			fmt.Println()
		}
	},
}

// slaDefineCmd represents the SLA define command
var slaDefineCmd = &cobra.Command{
	Use:   "define [table_name] [table_path] [frequency]",
	Short: "Define an SLA for a table",
	Long: `Define a Service Level Agreement (SLA) for a Delta table.
The frequency can be specified in Go duration format (e.g., 1h, 30m, 24h)
or using keywords like 'hourly', 'daily', 'weekly', or 'monthly'.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if freshness tracking is enabled
		if !freshnessEnabled {
			fmt.Println("Freshness tracking is disabled. Use --enable-freshness-tracking to enable it.")
			return
		}
		
		// Parse arguments
		tableName := args[0]
		tablePath := args[1]
		frequencyStr := args[2]
		
		// Parse frequency
		frequency, err := freshness.ParseDuration(frequencyStr)
		if err != nil {
			logging.Error("Failed to parse frequency", err, "frequency", frequencyStr)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Create SLA configuration
		config := &freshness.SLAConfig{
			TableName:         tableName,
			TablePath:         tablePath,
			ExpectedFrequency: frequency,
			WarningThreshold:  slaWarningThreshold,
			CriticalThreshold: slaCriticalThreshold,
			Enabled:           true,
			Description:       slaDescription,
			Tags:              slaTags,
			GrafanaDashboardURL: slaGrafanaDashboardURL,
		}
		
		// Validate configuration
		if err := config.Validate(); err != nil {
			logging.Error("Invalid SLA configuration", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Create SLA manager
		configPath := filepath.Join(configDir, "sla_config.json")
		manager, err := freshness.NewSLAManager(configPath, monitoringClient)
		if err != nil {
			logging.Error("Failed to create SLA manager", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Set SLA configuration
		if err := manager.SetSLA(config); err != nil {
			logging.Error("Failed to set SLA configuration", err, "table", tableName)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		fmt.Printf("SLA defined for table %s with frequency %s\n", tableName, freshness.FormatDuration(frequency))
	},
}

// slaDeleteCmd represents the SLA delete command
var slaDeleteCmd = &cobra.Command{
	Use:   "delete [table_name]",
	Short: "Delete an SLA configuration",
	Long:  `Delete an SLA configuration for a table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if freshness tracking is enabled
		if !freshnessEnabled {
			fmt.Println("Freshness tracking is disabled. Use --enable-freshness-tracking to enable it.")
			return
		}
		
		// Parse arguments
		tableName := args[0]
		
		// Create SLA manager
		configPath := filepath.Join(configDir, "sla_config.json")
		manager, err := freshness.NewSLAManager(configPath, monitoringClient)
		if err != nil {
			logging.Error("Failed to create SLA manager", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Delete SLA configuration
		if err := manager.DeleteSLA(tableName); err != nil {
			logging.Error("Failed to delete SLA configuration", err, "table", tableName)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		fmt.Printf("SLA configuration deleted for table %s\n", tableName)
	},
}

// slaEnableCmd represents the SLA enable command
var slaEnableCmd = &cobra.Command{
	Use:   "enable [table_name]",
	Short: "Enable SLA monitoring for a table",
	Long:  `Enable SLA monitoring for a table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if freshness tracking is enabled
		if !freshnessEnabled {
			fmt.Println("Freshness tracking is disabled. Use --enable-freshness-tracking to enable it.")
			return
		}
		
		// Parse arguments
		tableName := args[0]
		
		// Create SLA manager
		configPath := filepath.Join(configDir, "sla_config.json")
		manager, err := freshness.NewSLAManager(configPath, monitoringClient)
		if err != nil {
			logging.Error("Failed to create SLA manager", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Get SLA configuration
		config, err := manager.GetSLA(tableName)
		if err != nil {
			logging.Error("Failed to get SLA configuration", err, "table", tableName)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Enable SLA monitoring
		config.Enabled = true
		
		// Update SLA configuration
		if err := manager.SetSLA(config); err != nil {
			logging.Error("Failed to update SLA configuration", err, "table", tableName)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		fmt.Printf("SLA monitoring enabled for table %s\n", tableName)
	},
}

// slaDisableCmd represents the SLA disable command
var slaDisableCmd = &cobra.Command{
	Use:   "disable [table_name]",
	Short: "Disable SLA monitoring for a table",
	Long:  `Disable SLA monitoring for a table.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if freshness tracking is enabled
		if !freshnessEnabled {
			fmt.Println("Freshness tracking is disabled. Use --enable-freshness-tracking to enable it.")
			return
		}
		
		// Parse arguments
		tableName := args[0]
		
		// Create SLA manager
		configPath := filepath.Join(configDir, "sla_config.json")
		manager, err := freshness.NewSLAManager(configPath, monitoringClient)
		if err != nil {
			logging.Error("Failed to create SLA manager", err)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Get SLA configuration
		config, err := manager.GetSLA(tableName)
		if err != nil {
			logging.Error("Failed to get SLA configuration", err, "table", tableName)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		// Disable SLA monitoring
		config.Enabled = false
		
		// Update SLA configuration
		if err := manager.SetSLA(config); err != nil {
			logging.Error("Failed to update SLA configuration", err, "table", tableName)
			fmt.Printf("Error: %v\n", err)
			return
		}
		
		fmt.Printf("SLA monitoring disabled for table %s\n", tableName)
	},
}

func init() {
	rootCmd.AddCommand(freshnessCmd)
	freshnessCmd.AddCommand(freshnessCheckCmd)
	
	rootCmd.AddCommand(slaCmd)
	slaCmd.AddCommand(slaListCmd)
	slaCmd.AddCommand(slaDefineCmd)
	slaCmd.AddCommand(slaDeleteCmd)
	slaCmd.AddCommand(slaEnableCmd)
	slaCmd.AddCommand(slaDisableCmd)
	
	// Add flags
	rootCmd.PersistentFlags().BoolVar(&freshnessEnabled, "enable-freshness-tracking", false, "Enable freshness tracking")
	
	freshnessCheckCmd.Flags().StringVar(&slaOutputFormat, "format", "table", "Output format (table or json)")
	
	slaListCmd.Flags().StringVar(&slaOutputFormat, "format", "table", "Output format (table or json)")
	
	slaDefineCmd.Flags().IntVar(&slaWarningThreshold, "warning", 150, "Warning threshold (percentage of expected frequency)")
	slaDefineCmd.Flags().IntVar(&slaCriticalThreshold, "critical", 200, "Critical threshold (percentage of expected frequency)")
	slaDefineCmd.Flags().StringVar(&slaDescription, "description", "", "Description of the SLA")
	slaDefineCmd.Flags().StringSliceVar(&slaTags, "tags", []string{}, "Tags for the SLA (comma-separated)")
	slaDefineCmd.Flags().StringVar(&slaGrafanaDashboardURL, "grafana-url", "", "Grafana dashboard URL")
}
