package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nessi-dev/nessi/pkg/quality"
	"github.com/nessi-dev/nessi/pkg/visualization"
	"github.com/spf13/cobra"
)

// visualizeHeatmapCmd represents the heatmap visualization command
var visualizeHeatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "Generate a CLI heatmap visualization",
	Long:  `Generate a heatmap visualization in the terminal for data quality metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		tablePath, _ := cmd.Flags().GetString("table")
		metric, _ := cmd.Flags().GetString("metric")
		partitionCol, _ := cmd.Flags().GetString("partition-column")
		colorMode, _ := cmd.Flags().GetString("color-mode")
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")

		// Get quality metrics for the table
		metrics, err := quality.GetTableMetrics(tablePath)
		if err != nil {
			fmt.Printf("Error getting metrics: %v\n", err)
			os.Exit(1)
		}

		// Generate and display the heatmap
		err = visualization.DisplayHeatmap(metrics, metric, partitionCol, colorMode, width, height)
		if err != nil {
			fmt.Printf("Error displaying heatmap: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	visualizeCmd.AddCommand(visualizeHeatmapCmd)

	// Add flags
	visualizeHeatmapCmd.Flags().StringP("table", "t", "", "Path to the Delta Lake table (required)")
	visualizeHeatmapCmd.MarkFlagRequired("table")
	visualizeHeatmapCmd.Flags().StringP("metric", "m", "completeness", "Metric to visualize (completeness, accuracy, freshness)")
	visualizeHeatmapCmd.Flags().StringP("partition-column", "p", "", "Partition column to use for the heatmap")
	visualizeHeatmapCmd.Flags().StringP("color-mode", "c", "gradient", "Color mode (gradient, binary)")
	visualizeHeatmapCmd.Flags().IntP("width", "w", 80, "Width of the heatmap")
	visualizeHeatmapCmd.Flags().IntP("height", "h", 20, "Height of the heatmap")
}
