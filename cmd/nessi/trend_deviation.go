package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
)

// trendDeviationCmd represents the trend-deviation command
var trendDeviationCmd = &cobra.Command{
	Use:   "trend-deviation",
	Short: "Analyze trend deviations in data quality metrics",
	Long: `Analyze trend deviations in data quality metrics and alert on significant changes.
This command compares current metrics with previous runs and identifies significant deviations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		tablePath, _ := cmd.Flags().GetString("table")
		field, _ := cmd.Flags().GetString("field")
		metricTypesStr, _ := cmd.Flags().GetString("metrics")
		thresholdStr, _ := cmd.Flags().GetString("threshold")
		previousRunsStr, _ := cmd.Flags().GetString("previous-runs")
		metricsDir, _ := cmd.Flags().GetString("metrics-dir")
		outputFormat, _ := cmd.Flags().GetString("output")

		// Validate table path
		if tablePath == "" {
			fmt.Println("Error: table path is required")
			os.Exit(1)
		}

		// Validate field
		if field == "" {
			fmt.Println("Error: field is required")
			os.Exit(1)
		}

		// Parse metric types
		metricTypes := []datalake.MetricType{}
		if metricTypesStr != "" {
			metricTypeNames := strings.Split(metricTypesStr, ",")
			for _, name := range metricTypeNames {
				name = strings.TrimSpace(name)
				metricType := datalake.MetricType(name)
				metricTypes = append(metricTypes, metricType)
			}
		} else {
			// Default metrics
			metricTypes = []datalake.MetricType{
				datalake.NullPercentage,
				datalake.UniqueRatio,
				datalake.MinValue,
				datalake.MaxValue,
				datalake.MeanValue,
				datalake.MedianValue,
				datalake.StandardDeviation,
				datalake.RecordCount,
			}
		}

		// Parse threshold
		threshold := 10.0 // Default threshold is 10%
		if thresholdStr != "" {
			var err error
			threshold, err = strconv.ParseFloat(thresholdStr, 64)
			if err != nil {
				fmt.Printf("Error parsing threshold: %v\n", err)
				os.Exit(1)
			}
		}

		// Parse previous runs
		previousRuns := 3 // Default is 3 previous runs
		if previousRunsStr != "" {
			var err error
			previousRuns, err = strconv.Atoi(previousRunsStr)
			if err != nil {
				fmt.Printf("Error parsing previous runs: %v\n", err)
				os.Exit(1)
			}
		}

		// Set default metrics directory if not provided
		if metricsDir == "" {
			metricsDir = filepath.Join(tablePath, "_metrics")
		}

		// Create metadata manager
		manager := datalake.NewMetadataManager(tablePath)

		// Create options
		options := datalake.TrendDeviationOptions{
			Field:        field,
			MetricTypes:  metricTypes,
			Threshold:    threshold,
			PreviousRuns: previousRuns,
			MetricsDir:   metricsDir,
		}

		// Analyze trend deviation
		fmt.Printf("Analyzing trend deviation for field '%s'...\n", field)
		result, err := manager.AnalyzeTrendDeviation(options)
		if err != nil {
			fmt.Printf("Error analyzing trend deviation: %v\n", err)
			os.Exit(1)
		}

		// Output results
		switch outputFormat {
		case "json":
			// Output as JSON
			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Printf("Error marshaling result to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonData))
		default:
			// Output as text
			fmt.Printf("Trend Deviation Analysis for field '%s'\n", result.Field)
			fmt.Printf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
			fmt.Println()

			// Print metrics
			fmt.Printf("Metrics (%d):\n", len(result.Metrics))
			for _, metric := range result.Metrics {
				fmt.Printf("  %s:\n", metric.MetricType)
				fmt.Printf("    Current Value: %.4f\n", metric.CurrentValue)
				if len(metric.PreviousValues) > 0 {
					fmt.Printf("    Previous Value: %.4f\n", metric.PreviousValues[0])
					fmt.Printf("    Percentage Change: %.2f%%\n", metric.PercentageChange)
				}
				if len(metric.PreviousValues) > 1 {
					fmt.Printf("    Average Value: %.4f\n", metric.AverageValue)
					fmt.Printf("    Standard Deviation: %.4f\n", metric.StandardDeviation)
					fmt.Printf("    Z-Score: %.2f\n", metric.ZScore)
				}
				fmt.Println()
			}

			// Print alerts
			if len(result.Alerts) > 0 {
				fmt.Printf("Alerts (%d):\n", len(result.Alerts))
				for i, alert := range result.Alerts {
					fmt.Printf("  %d. [%s] %s\n", i+1, alert.Severity, alert.Message)
				}
			} else {
				fmt.Println("No alerts generated.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(trendDeviationCmd)

	// Add flags
	trendDeviationCmd.Flags().StringP("table", "t", "", "Path to the Delta table (required)")
	trendDeviationCmd.Flags().StringP("field", "f", "", "Field to analyze (required)")
	trendDeviationCmd.Flags().StringP("metrics", "m", "", "Comma-separated list of metrics to analyze (default: all)")
	trendDeviationCmd.Flags().StringP("threshold", "", "10.0", "Threshold for alerting (percentage change)")
	trendDeviationCmd.Flags().StringP("previous-runs", "p", "3", "Number of previous runs to compare with")
	trendDeviationCmd.Flags().StringP("metrics-dir", "d", "", "Directory where metrics are stored (default: <table-path>/_metrics)")
	trendDeviationCmd.Flags().StringP("output", "o", "text", "Output format (text or json)")

	// Mark required flags
	trendDeviationCmd.MarkFlagRequired("table")
	trendDeviationCmd.MarkFlagRequired("field")
}
