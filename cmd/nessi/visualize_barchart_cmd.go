package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// visualizeBarChartCmd represents the bar chart visualization command
var visualizeBarChartCmd = &cobra.Command{
	Use:   "barchart",
	Short: "Generate a CLI bar chart visualization",
	Long:  `Generate a bar chart visualization in the terminal for data quality metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the table path and metrics
		tablePath, _ := cmd.Flags().GetString("table")
		showAll, _ := cmd.Flags().GetBool("all-metrics")

		// Validate inputs
		if tablePath == "" {
			fmt.Println("Error: No table path specified.")
			fmt.Println("Usage: nessi visualize barchart --table <path_to_delta_table> [--all-metrics]")
			fmt.Println("Example: nessi visualize barchart --table s3://my-bucket/my-table")
			return
		}

		// Check if the table exists (in a real implementation)
		// For now, we'll just print a message
		fmt.Printf("Checking table path: %s...\n", tablePath)

		fmt.Printf("Generating bar chart visualization for table: %s\n\n", tablePath)

		// In a real implementation, we would fetch actual metrics data
		// For demonstration purposes, we'll create a sample bar chart
		generateBarChartVisualization(tablePath, showAll)
	},
}

// generateBarChartVisualization creates a simple ASCII bar chart visualization
func generateBarChartVisualization(tablePath string, showAll bool) {
	// Define colors for the bar chart
	goodColor := color.New(color.FgGreen).SprintFunc()
	mediumColor := color.New(color.FgYellow).SprintFunc()
	badColor := color.New(color.FgRed).SprintFunc()
	titleColor := color.New(color.FgCyan, color.Bold).SprintFunc()

	// Define metrics to display
	metrics := []string{"Completeness", "Accuracy", "Consistency", "Validity"}
	if showAll {
		metrics = append(metrics, "Timeliness", "Uniqueness", "Integrity", "Conformity")
	}

	// Print title
	fmt.Println(titleColor("Data Quality Metrics Bar Chart"))
	fmt.Println(titleColor(strings.Repeat("=", 30)))
	fmt.Println()

	// Generate random data for demonstration
	maxBarLength := 40
	for _, metric := range metrics {
		// Generate a random value between 0 and 1
		value := rand.Float64()

		// Calculate bar length
		barLength := int(value * float64(maxBarLength))

		// Determine color based on value
		barColor := badColor
		if value > 0.8 {
			barColor = goodColor
		} else if value > 0.5 {
			barColor = mediumColor
		}

		// Print metric name and value
		fmt.Printf("%-12s [%5.1f%%] ", metric, value*100)

		// Print the bar
		fmt.Print(barColor(strings.Repeat("█", barLength)))
		fmt.Println()
	}

	fmt.Println()
	fmt.Println("Legend:")
	fmt.Printf("%s Good (>80%%)  ", goodColor("█████"))
	fmt.Printf("%s Medium (50-80%%)  ", mediumColor("█████"))
	fmt.Printf("%s Poor (<50%%)\n", badColor("█████"))
	fmt.Println()
	fmt.Println("Note: This is a demonstration visualization. In a production environment,")
	fmt.Println("this would use actual data quality metrics from the specified table.")
}

func init() {
	visualizeCmd.AddCommand(visualizeBarChartCmd)

	// Add flags
	visualizeBarChartCmd.Flags().StringP("table", "t", "", "Path to the Delta Lake table (required)")
	visualizeBarChartCmd.MarkFlagRequired("table")
	visualizeBarChartCmd.Flags().BoolP("all-metrics", "a", false, "Show all available metrics")
}
