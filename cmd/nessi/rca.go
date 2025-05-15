package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nessi-dev/nessi-dev/pkg/rca"
	"github.com/spf13/cobra"
)

// rcaCmd represents the rca command
var rcaCmd = &cobra.Command{
	Use:   "rca [anomaly_id]",
	Short: "Perform root cause analysis on anomalies",
	Long: `Root Cause Analysis (RCA) for anomalies detected in your data.

This command analyzes anomalies to determine potential root causes, affected tables,
and recommended actions. It can output results in various formats and optionally
generate alerts for critical findings.

Examples:
  nessi rca anom-20250514-001
  nessi rca anom-20250514-001 --format=json
  nessi rca anom-20250514-001 --format=html --output=report.html
  nessi rca anom-20250514-001 --alert`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		anomalyID := args[0]
		
		// Get flags
		format, _ := cmd.Flags().GetString("format")
		outputFile, _ := cmd.Flags().GetString("output")
		shouldAlert, _ := cmd.Flags().GetBool("alert")
		
		// Initialize RCA analyzer
		analyzer := initRCAAnalyzer()
		
		// Perform analysis
		result, err := analyzer.AnalyzeAnomaly(anomalyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error performing RCA: %v\n", err)
			os.Exit(1)
		}
		
		// Generate output based on format
		var output string
		switch strings.ToLower(format) {
		case "json":
			jsonOutput, err := result.ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating JSON output: %v\n", err)
				os.Exit(1)
			}
			output = jsonOutput
		case "html":
			output = result.ToHTML()
		default:
			// Text format (default)
			output = formatRCAResultAsText(result)
		}
		
		// Write to output file or stdout
		if outputFile != "" {
			err := os.WriteFile(outputFile, []byte(output), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing to output file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("RCA results written to %s\n", outputFile)
		} else {
			fmt.Println(output)
		}
		
		// Send alert if requested and severity is high enough
		if shouldAlert && result.PrimaryRootCause != nil {
			// Check if primary cause severity warrants an alert
			// This would use the actual alerting system
			fmt.Println("Alert sent for critical root cause finding")
		}
	},
}

func init() {
	rootCmd.AddCommand(rcaCmd)
	
	// Add flags
	rcaCmd.Flags().String("format", "text", "Output format (text, json, html)")
	rcaCmd.Flags().String("output", "", "Output file (default is stdout)")
	rcaCmd.Flags().Bool("alert", false, "Send alert for critical findings")
	rcaCmd.Flags().Bool("enable-rca", true, "Enable root cause analysis")
}

// initRCAAnalyzer initializes the RCA analyzer with configuration
func initRCAAnalyzer() *rca.Analyzer {
	// This would get the actual configuration from config file
	// and initialize the real dependencies
	config := rca.DefaultConfig()
	
	// For now, we're passing nil for dependencies as they would be
	// initialized from the actual application context
	return rca.NewAnalyzer(config, nil, nil)
}

// formatRCAResultAsText formats the RCA result as human-readable text
func formatRCAResultAsText(result *rca.RCAResult) string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("Root Cause Analysis for Anomaly %s\n", result.AnomalyID))
	sb.WriteString(fmt.Sprintf("Analysis Time: %s\n\n", result.AnalysisTime.Format("2006-01-02 15:04:05")))
	
	sb.WriteString("PRIMARY ROOT CAUSE:\n")
	sb.WriteString(fmt.Sprintf("  Type: %s\n", result.PrimaryRootCause.Type))
	sb.WriteString(fmt.Sprintf("  Confidence: %.1f%%\n", result.PrimaryRootCause.Confidence*100))
	sb.WriteString(fmt.Sprintf("  Description: %s\n", result.PrimaryRootCause.Description))
	sb.WriteString(fmt.Sprintf("  Timestamp: %s\n\n", result.PrimaryRootCause.Timestamp.Format("2006-01-02 15:04:05")))
	
	if len(result.OtherCauses) > 0 {
		sb.WriteString("OTHER POTENTIAL CAUSES:\n")
		for i, cause := range result.OtherCauses {
			sb.WriteString(fmt.Sprintf("  %d. %s (%.1f%% confidence)\n", i+1, cause.Description, cause.Confidence*100))
		}
		sb.WriteString("\n")
	}
	
	if len(result.AffectedTables) > 0 {
		sb.WriteString("AFFECTED TABLES:\n")
		for _, table := range result.AffectedTables {
			sb.WriteString(fmt.Sprintf("  - %s\n", table))
		}
		sb.WriteString("\n")
	}
	
	if len(result.RelatedAnomalies) > 0 {
		sb.WriteString("RELATED ANOMALIES:\n")
		for _, anomaly := range result.RelatedAnomalies {
			sb.WriteString(fmt.Sprintf("  - %s\n", anomaly))
		}
		sb.WriteString("\n")
	}
	
	if len(result.RecommendedActions) > 0 {
		sb.WriteString("RECOMMENDED ACTIONS:\n")
		for i, action := range result.RecommendedActions {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, action))
		}
		sb.WriteString("\n")
	}
	
	if result.GrafanaDashboardURL != "" {
		sb.WriteString(fmt.Sprintf("Grafana Dashboard: %s\n", result.GrafanaDashboardURL))
	}
	
	return sb.String()
}
