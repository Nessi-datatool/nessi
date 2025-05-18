package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi/pkg/datalake"
)

// This demo uses the SimpleTrendAnalyzer from the datalake package

func main() {
	fmt.Println("Trend Deviation Demo")
	fmt.Println("====================")
	fmt.Println()

	// Create a temporary directory for metrics
	tempDir, err := os.MkdirTemp("", "trend-deviation-demo-*")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	metricsDir := filepath.Join(tempDir, "metrics")
	err = os.MkdirAll(metricsDir, 0755)
	if err != nil {
		fmt.Printf("Error creating metrics directory: %v\n", err)
		return
	}

	fmt.Printf("Created temporary metrics directory at: %s\n", metricsDir)
	fmt.Println()

	// Create a simple trend analyzer
	analyzer := datalake.NewSimpleTrendAnalyzer(metricsDir)

	// Run 1: Initial run to establish baseline
	fmt.Println("Run 1: Establishing baseline metrics...")
	result1, err := analyzer.AnalyzeTrendDeviation("sales", 0)
	if err != nil {
		fmt.Printf("Error analyzing trend deviation: %v\n", err)
		return
	}

	fmt.Println("Baseline metrics established.")
	fmt.Printf("Metrics for %s: %d\n", result1.Field, len(result1.Metrics))
	fmt.Println("No alerts expected for the first run.")
	fmt.Println()

	// Run 2: Increase sales by 10%
	fmt.Println("Run 2: Increasing sales by 10%...")
	result2, err := analyzer.AnalyzeTrendDeviation("sales", 1)
	if err != nil {
		fmt.Printf("Error analyzing trend deviation: %v\n", err)
		return
	}

	fmt.Println("Results for Run 2:")
	printTrendDeviationResult(result2)
	fmt.Println()

	// Run 3: Increase sales by another 10% and revenue by 20%
	fmt.Println("Run 3: Increasing sales by another 10% and revenue by 20%...")
	
	// First check sales
	result3Sales, err := analyzer.AnalyzeTrendDeviation("sales", 2)
	if err != nil {
		fmt.Printf("Error analyzing trend deviation for sales: %v\n", err)
		return
	}

	fmt.Println("Results for Run 3 (Sales):")
	printTrendDeviationResult(result3Sales)
	fmt.Println()

	// Then check revenue
	result3Revenue, err := analyzer.AnalyzeTrendDeviation("revenue", 2)
	if err != nil {
		fmt.Printf("Error analyzing trend deviation for revenue: %v\n", err)
		return
	}

	fmt.Println("Results for Run 3 (Revenue):")
	printTrendDeviationResult(result3Revenue)
	fmt.Println()

	// Run 4: Decrease customers by 30%
	fmt.Println("Run 4: Decreasing customers by 30%...")
	result4Customers, err := analyzer.AnalyzeTrendDeviation("customers", 3)
	if err != nil {
		fmt.Printf("Error analyzing trend deviation for customers: %v\n", err)
		return
	}

	fmt.Println("Results for Run 4 (Customers):")
	printTrendDeviationResult(result4Customers)
	fmt.Println()

	// Run 5: Check null percentage
	fmt.Println("Run 5: Checking null percentage...")
	result5Null, err := analyzer.AnalyzeTrendDeviation("null_field", 0)
	if err != nil {
		fmt.Printf("Error analyzing trend deviation for null percentage: %v\n", err)
		return
	}

	fmt.Println("Results for Run 5 (Null Percentage):")
	printTrendDeviationResult(result5Null)
	fmt.Println()

	fmt.Println("Demo completed!")
}

// printTrendDeviationResult prints the trend deviation result
func printTrendDeviationResult(result *datalake.TrendDeviationResult) {
	fmt.Printf("Field: %s\n", result.Field)
	fmt.Printf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
	
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
