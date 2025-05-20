package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// visualizeTimeSeriesCmd represents the time series visualization command
var visualizeTimeSeriesCmd = &cobra.Command{
	Use:   "timeseries",
	Short: "Generate a CLI time series visualization",
	Long:  `Generate a time series visualization in the terminal for tracking data quality metrics over time.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the table path and metrics
		tablePath, _ := cmd.Flags().GetString("table")
		metric, _ := cmd.Flags().GetString("metric")
		days, _ := cmd.Flags().GetInt("days")

		// Validate inputs
		if tablePath == "" {
			fmt.Println("Error: No table path specified.")
			fmt.Println("Usage: nessi visualize timeseries --table <path_to_delta_table> --metric <metric_name> [--days <number_of_days>]")
			fmt.Println("Example: nessi visualize timeseries --table s3://my-bucket/my-table --metric completeness --days 30")
			return
		}

		if metric == "" {
			fmt.Println("Error: No metric specified.")
			fmt.Println("Available metrics: completeness, accuracy, freshness, consistency, validity")
			fmt.Println("Example: nessi visualize timeseries --table s3://my-bucket/my-table --metric completeness")
			return
		}

		// Check if the table exists (in a real implementation)
		// For now, we'll just print a message
		fmt.Printf("Checking table path: %s...\n", tablePath)

		fmt.Printf("Generating time series visualization for table: %s\n", tablePath)
		fmt.Printf("Metric: %s, Last %d days\n\n", metric, days)

		// In a real implementation, we would fetch actual metrics data
		// For demonstration purposes, we'll create a sample time series
		generateTimeSeriesVisualization(metric, days)
	},
}

// generateTimeSeriesVisualization creates a simple ASCII time series visualization
func generateTimeSeriesVisualization(metric string, days int) {
	// Define colors for the time series
	goodColor := color.New(color.FgGreen).SprintFunc()
	mediumColor := color.New(color.FgYellow).SprintFunc()
	badColor := color.New(color.FgRed).SprintFunc()
	titleColor := color.New(color.FgCyan, color.Bold).SprintFunc()
	axisColor := color.New(color.FgWhite).SprintFunc()

	// Print title
	fmt.Println(titleColor("Data Quality Time Series"))
	fmt.Println(titleColor(strings.Repeat("=", 30)))
	fmt.Println()

	// Generate random data for demonstration
	values := make([]float64, days)
	for i := 0; i < days; i++ {
		// Generate a value between 0.5 and 1.0 with some random fluctuation
		baseValue := 0.75 + (rand.Float64() * 0.25) - 0.125
		// Add a slight trend (improving over time)
		trendFactor := float64(i) / float64(days) * 0.1
		values[i] = baseValue + trendFactor
		if values[i] > 1.0 {
			values[i] = 1.0
		} else if values[i] < 0.0 {
			values[i] = 0.0
		}
	}

	// Determine the y-axis scale
	minValue := 0.0
	maxValue := 1.0
	height := 20

	// Print the time series
	// Y-axis labels and grid
	fmt.Printf("%s\n", metric)
	fmt.Println("^")
	for i := height - 1; i >= 0; i-- {
		// Calculate y value for current row (used for axis labels)
		yValue := minValue + (float64(i)/float64(height-1))*(maxValue-minValue)

		// Y-axis labels
		if i == height-1 {
			fmt.Printf("%.1f ", yValue)
		} else if i == 0 {
			fmt.Printf("%.1f ", yValue)
		} else if i == height/2 {
			fmt.Printf("%.1f ", yValue)
		} else {
			fmt.Print("    ")
		}

		// Grid line
		fmt.Print(axisColor("|"))

		// Plot the data points
		for j := 0; j < days; j++ {
			valueHeight := int((values[j] - minValue) / (maxValue - minValue) * float64(height-1))

			if valueHeight == i {
				// Choose color based on value
				pointChar := "*"
				if values[j] > 0.8 {
					fmt.Print(goodColor(pointChar))
				} else if values[j] > 0.5 {
					fmt.Print(mediumColor(pointChar))
				} else {
					fmt.Print(badColor(pointChar))
				}
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	// X-axis
	fmt.Print("    ")
	fmt.Print(axisColor("+"))
	fmt.Print(axisColor(strings.Repeat("-", days)))
	fmt.Println(axisColor(">"))

	// X-axis labels (dates)
	fmt.Print("      ")

	// Get current date and calculate start date
	currentDate := time.Now()
	startDate := currentDate.AddDate(0, 0, -days+1)

	// Print first date
	fmt.Print(startDate.Format("01/02"))

	// Print middle date
	middleOffset := days / 2
	middleDate := startDate.AddDate(0, 0, middleOffset)
	fmt.Print(strings.Repeat(" ", middleOffset-3))
	fmt.Print(middleDate.Format("01/02"))

	// Print end date (current date)
	fmt.Print(strings.Repeat(" ", days-middleOffset-3))
	fmt.Println(currentDate.Format("01/02"))

	fmt.Println()
	fmt.Println("Legend:")
	fmt.Printf("%s Good (>80%%)  ", goodColor("*"))
	fmt.Printf("%s Medium (50-80%%)  ", mediumColor("*"))
	fmt.Printf("%s Poor (<50%%)\n", badColor("*"))
	fmt.Println()
	fmt.Println("Note: This is a demonstration visualization. In a production environment,")
	fmt.Println("this would use actual historical data quality metrics from the specified table.")
}

func init() {
	visualizeCmd.AddCommand(visualizeTimeSeriesCmd)

	// Add flags
	visualizeTimeSeriesCmd.Flags().StringP("table", "t", "", "Path to the Delta Lake table (required)")
	visualizeTimeSeriesCmd.MarkFlagRequired("table")
	visualizeTimeSeriesCmd.Flags().StringP("metric", "m", "completeness", "Metric to visualize (e.g., completeness, accuracy)")
	visualizeTimeSeriesCmd.Flags().IntP("days", "d", 30, "Number of days to display in the time series")
}
