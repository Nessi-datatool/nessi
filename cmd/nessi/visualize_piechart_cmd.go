package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// visualizePieChartCmd represents the pie chart visualization command
var visualizePieChartCmd = &cobra.Command{
	Use:   "piechart",
	Short: "Generate a CLI pie chart visualization",
	Long:  `Generate a pie chart visualization in the terminal for showing the distribution of data quality issues.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the table path and category
		tablePath, _ := cmd.Flags().GetString("table")
		category, _ := cmd.Flags().GetString("category")

		// Validate inputs
		if tablePath == "" {
			fmt.Println("Error: No table path specified.")
			fmt.Println("Usage: nessi visualize piechart --table <path_to_delta_table> [--category <category>]")
			fmt.Println("Example: nessi visualize piechart --table s3://my-bucket/my-table --category issue-types")
			return
		}

		// Check if the table exists (in a real implementation)
		// For now, we'll just print a message
		fmt.Printf("Checking table path: %s...\n", tablePath)

		fmt.Printf("Generating pie chart visualization for table: %s\n", tablePath)
		if category != "" {
			fmt.Printf("Category: %s\n", category)
			fmt.Println()
		} else {
			fmt.Println("Category: issue-types (default)")
			fmt.Println()
		}

		// In a real implementation, we would fetch actual metrics data
		// For demonstration purposes, we'll create a sample pie chart
		generatePieChartVisualization(category)
	},
}

// generatePieChartVisualization creates a simple ASCII pie chart visualization
func generatePieChartVisualization(category string) {
	// Define colors for the pie chart
	colors := []*color.Color{
		color.New(color.FgRed),
		color.New(color.FgGreen),
		color.New(color.FgYellow),
		color.New(color.FgBlue),
		color.New(color.FgMagenta),
		color.New(color.FgCyan),
	}
	titleColor := color.New(color.FgCyan, color.Bold).SprintFunc()

	// Print title
	fmt.Println(titleColor("Data Quality Issues Distribution"))
	fmt.Println(titleColor(strings.Repeat("=", 35)))
	fmt.Println()

	// Define categories and values based on the selected category
	var categories []string
	var values []float64
	var symbols []string

	if category == "issue-types" || category == "" {
		// Distribution of issue types
		categories = []string{"Missing Values", "Invalid Format", "Out of Range", "Duplicates", "Inconsistencies"}
		// Generate random values that sum to 100
		values = []float64{35, 25, 20, 15, 5}
		symbols = []string{"M", "F", "R", "D", "I"}
	} else if category == "columns" {
		// Distribution of issues by column
		categories = []string{"id", "name", "email", "address", "phone", "date"}
		values = []float64{10, 15, 30, 20, 15, 10}
		symbols = []string{"I", "N", "E", "A", "P", "D"}
	} else if category == "severity" {
		// Distribution of issues by severity
		categories = []string{"Critical", "High", "Medium", "Low"}
		values = []float64{15, 25, 40, 20}
		symbols = []string{"C", "H", "M", "L"}
	} else {
		// Default to issue types if category is not recognized
		categories = []string{"Missing Values", "Invalid Format", "Out of Range", "Duplicates", "Inconsistencies"}
		values = []float64{35, 25, 20, 15, 5}
		symbols = []string{"M", "F", "R", "D", "I"}
	}

	// Randomize the values slightly to make it more realistic
	total := 0.0
	for i := range values {
		// Add or subtract up to 5% randomly
		randomFactor := 1.0 + (rand.Float64()*0.1 - 0.05)
		values[i] = values[i] * randomFactor
		total += values[i]
	}

	// Normalize to ensure they sum to 100
	for i := range values {
		values[i] = (values[i] / total) * 100
	}

	// Draw the pie chart (ASCII art style)
	drawASCIIPieChart(categories, values, symbols, colors)
}

// drawASCIIPieChart draws a simple ASCII pie chart
func drawASCIIPieChart(categories []string, values []float64, symbols []string, colors []*color.Color) {
	// Calculate total for percentages
	total := 0.0
	for _, v := range values {
		total += v
	}

	// Define the pie chart dimensions
	radius := 10
	centerX := radius
	centerY := radius

	// Create a 2D grid for the pie chart
	grid := make([][]string, radius*2+1)
	for i := range grid {
		grid[i] = make([]string, radius*2+1)
		for j := range grid[i] {
			grid[i][j] = " "
		}
	}

	// Calculate the angles for each slice
	angles := make([]float64, len(values)+1)
	angles[0] = 0
	for i := 0; i < len(values); i++ {
		angles[i+1] = angles[i] + (values[i]/100)*2*math.Pi
	}

	// Fill in the grid with pie chart segments
	for y := 0; y <= 2*radius; y++ {
		for x := 0; x <= 2*radius; x++ {
			// Calculate distance from center
			dx := float64(x - centerX)
			dy := float64(y - centerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			// If within the circle
			if distance <= float64(radius) {
				// Calculate angle
				angle := math.Atan2(dy, dx)
				if angle < 0 {
					angle += 2 * math.Pi
				}

				// Determine which segment this point belongs to
				for i := 0; i < len(values); i++ {
					if angle >= angles[i] && angle < angles[i+1] {
						grid[y][x] = symbols[i]
						break
					}
				}
			}
		}
	}

	// Print the pie chart
	for y := 0; y <= 2*radius; y++ {
		fmt.Print("  ")
		for x := 0; x <= 2*radius; x++ {
			char := grid[y][x]
			if char != " " {
				// Find the index of the symbol
				index := -1
				for i, s := range symbols {
					if s == char {
						index = i
						break
					}
				}

				// Print colored symbol
				if index >= 0 && index < len(colors) {
					fmt.Print(colors[index].Sprint(char))
				} else {
					fmt.Print(char)
				}
			} else {
				fmt.Print(char)
			}
		}
		fmt.Println()
	}

	fmt.Println()

	// Print the legend
	fmt.Println("Legend:")
	for i, category := range categories {
		percentage := values[i]
		colorFunc := colors[i%len(colors)].SprintFunc()
		fmt.Printf("%s %s: %.1f%% - %s\n", colorFunc(symbols[i]), colorFunc("■"), percentage, category)
	}

	fmt.Println()
	fmt.Println("Note: This is a demonstration visualization. In a production environment,")
	fmt.Println("this would use actual data quality issue statistics from the specified table.")
}

func init() {
	visualizeCmd.AddCommand(visualizePieChartCmd)

	// Add flags
	visualizePieChartCmd.Flags().StringP("table", "t", "", "Path to the Delta Lake table (required)")
	visualizePieChartCmd.MarkFlagRequired("table")
	visualizePieChartCmd.Flags().StringP("category", "c", "issue-types", "Category to visualize (issue-types, columns, severity)")
}
