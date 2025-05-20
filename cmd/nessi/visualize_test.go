package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestVisualizationCommands tests the basic functionality of visualization commands
func TestVisualizationCommands(t *testing.T) {
	// Skip this test in CI environments or when running short tests
	if testing.Short() {
		t.Skip("Skipping visualization tests in short mode")
	}

	// Test that all visualization commands are registered correctly
	t.Run("CommandRegistration", func(t *testing.T) {
		// Check that all visualization commands are registered with the visualize command
		subCommands := visualizeCmd.Commands()
		commandNames := make([]string, 0, len(subCommands))
		for _, cmd := range subCommands {
			commandNames = append(commandNames, cmd.Name())
		}

		// Check for the presence of each visualization command
		expectedCommands := []string{"heatmap", "barchart", "scatter", "timeseries", "piechart"}
		for _, expected := range expectedCommands {
			found := false
			for _, name := range commandNames {
				if name == expected {
					found = true
					break
				}
			}
			assert.True(t, found, "Visualization command '%s' should be registered", expected)
		}
	})

	// Test that each visualization command has the required flags
	t.Run("CommandFlags", func(t *testing.T) {
		// Check heatmap command flags
		assert.NotNil(t, visualizeHeatmapCmd.Flags().Lookup("table"), "Heatmap command should have 'table' flag")
		assert.NotNil(t, visualizeHeatmapCmd.Flags().Lookup("metric"), "Heatmap command should have 'metric' flag")

		// Check barchart command flags
		assert.NotNil(t, visualizeBarChartCmd.Flags().Lookup("table"), "Bar chart command should have 'table' flag")
		assert.NotNil(t, visualizeBarChartCmd.Flags().Lookup("all-metrics"), "Bar chart command should have 'all-metrics' flag")

		// Check scatter command flags
		assert.NotNil(t, visualizeScatterCmd.Flags().Lookup("table"), "Scatter plot command should have 'table' flag")
		assert.NotNil(t, visualizeScatterCmd.Flags().Lookup("x-metric"), "Scatter plot command should have 'x-metric' flag")
		assert.NotNil(t, visualizeScatterCmd.Flags().Lookup("y-metric"), "Scatter plot command should have 'y-metric' flag")

		// Check timeseries command flags
		assert.NotNil(t, visualizeTimeSeriesCmd.Flags().Lookup("table"), "Time series command should have 'table' flag")
		assert.NotNil(t, visualizeTimeSeriesCmd.Flags().Lookup("metric"), "Time series command should have 'metric' flag")
		assert.NotNil(t, visualizeTimeSeriesCmd.Flags().Lookup("days"), "Time series command should have 'days' flag")

		// Check piechart command flags
		assert.NotNil(t, visualizePieChartCmd.Flags().Lookup("table"), "Pie chart command should have 'table' flag")
		assert.NotNil(t, visualizePieChartCmd.Flags().Lookup("category"), "Pie chart command should have 'category' flag")
	})

	// Test the command run functions with mock arguments
	t.Run("CommandRunFunctions", func(t *testing.T) {
		// Test each command's Run function with mock arguments
		testCases := []struct {
			name string
			run  func(cmd *cobra.Command, args []string)
		}{
			{"HeatmapVisualization", visualizeHeatmapCmd.Run},
			{"BarChartVisualization", visualizeBarChartCmd.Run},
			{"ScatterPlotVisualization", visualizeScatterCmd.Run},
			{"TimeSeriesVisualization", visualizeTimeSeriesCmd.Run},
			{"PieChartVisualization", visualizePieChartCmd.Run},
		}

		for _, tc := range testCases {
			// Create a mock command with the required flags
			cmd := &cobra.Command{}
			cmd.Flags().String("table", "s3://test-bucket/test-table", "")
			cmd.Flags().String("metric", "completeness", "")
			cmd.Flags().String("x-metric", "completeness", "")
			cmd.Flags().String("y-metric", "accuracy", "")
			cmd.Flags().Bool("all-metrics", true, "")
			cmd.Flags().String("category", "issue-types", "")
			cmd.Flags().Int("days", 14, "")

			// Redirect stdout to capture output
			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute the command's Run function
			tc.run(cmd, []string{})

			// Restore stdout
			w.Close()
			os.Stdout = stdout
			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Check that the command produced some output
			output := buf.String()
			assert.NotEmpty(t, output, "Command '%s' should produce output", tc.name)
		}
	})
}

// TestVisualizationCommandsErrorHandling tests the error handling of visualization commands
func TestVisualizationCommandsErrorHandling(t *testing.T) {
	// Skip this test in CI environments or when running short tests
	if testing.Short() {
		t.Skip("Skipping visualization error handling tests in short mode")
	}

	// Test error handling for missing required arguments
	t.Run("MissingRequiredArgs", func(t *testing.T) {
		// Test cases for commands with missing required arguments
		testCases := []struct {
			name     string
			run      func(cmd *cobra.Command, args []string)
			expected string
		}{
			{
				name:     "HeatmapNoTable",
				run:      visualizeHeatmapCmd.Run,
				expected: "Error: No table path specified",
			},
			{
				name:     "BarChartNoTable",
				run:      visualizeBarChartCmd.Run,
				expected: "Error: No table path specified",
			},
			{
				name:     "ScatterNoTable",
				run:      visualizeScatterCmd.Run,
				expected: "Error: No table path specified",
			},
			{
				name:     "TimeSeriesNoTable",
				run:      visualizeTimeSeriesCmd.Run,
				expected: "Error: No table path specified",
			},
			{
				name:     "PieChartNoTable",
				run:      visualizePieChartCmd.Run,
				expected: "Error: No table path specified",
			},
		}

		for _, tc := range testCases {
			// Create a mock command with empty flags
			cmd := &cobra.Command{}
			cmd.Flags().String("table", "", "")
			cmd.Flags().String("metric", "", "")
			cmd.Flags().String("x-metric", "", "")
			cmd.Flags().String("y-metric", "", "")
			cmd.Flags().Bool("all-metrics", false, "")
			cmd.Flags().String("category", "", "")
			cmd.Flags().Int("days", 0, "")

			// Redirect stdout to capture output
			stdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Execute the command's Run function
			tc.run(cmd, []string{})

			// Restore stdout
			w.Close()
			os.Stdout = stdout
			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Check that the output contains the expected error message
			output := buf.String()
			assert.Contains(t, output, tc.expected, "Command '%s' should return error message '%s'", tc.name, tc.expected)
		}
	})
}
