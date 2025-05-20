package main

import (
	"github.com/spf13/cobra"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize data quality metrics",
	Long:  `Generate visualizations for data quality metrics in various formats.`,
}

func init() {
	rootCmd.AddCommand(visualizeCmd)
}
