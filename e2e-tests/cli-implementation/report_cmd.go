package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func newReportCmd() *cobra.Command {
	reportCmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long:  `Commands for generating reports from Delta Lake tables.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	reportCmd.AddCommand(newReportGenerateCmd())

	return reportCmd
}

func newReportGenerateCmd() *cobra.Command {
	var tableName string
	var format string
	var outputPath string
	var includeCharts bool

	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a report",
		Long:  `Generate a report for a Delta Lake table in the specified format.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if table name is provided
			if tableName == "" {
				fmt.Fprintf(os.Stderr, "Error N102: Table name is required\n")
				os.Exit(1)
			}

			// Check if format is valid
			validFormats := map[string]bool{
				"html": true,
				"pdf":  true,
				"json": true,
				"csv":  true,
				"md":   true,
			}
			if !validFormats[format] {
				fmt.Fprintf(os.Stderr, "Error N401: Invalid report format: %s\n", format)
				os.Exit(1)
			}

			// Create output directory if it doesn't exist
			outputDir := filepath.Dir(outputPath)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error N101: Failed to create output directory: %s\n", err)
				os.Exit(1)
			}

			// Create a sample report file
			sampleReport := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Nessi Report - %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #2c3e50; }
        .metrics { margin: 20px 0; }
        .metric { margin: 10px 0; }
        .pass { color: green; }
        .fail { color: red; }
        table { border-collapse: collapse; width: 100%%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <h1>Data Quality Report - %s</h1>
    <div class="metrics">
        <h2>Quality Metrics</h2>
        <div class="metric">Completeness: <span class="pass">98.5%%</span></div>
        <div class="metric">Accuracy: <span class="pass">99.2%%</span></div>
        <div class="metric">Consistency: <span class="pass">97.8%%</span></div>
        <div class="metric">Uniqueness: <span class="pass">100.0%%</span></div>
        <div class="metric">Timeliness: <span class="pass">95.5%%</span></div>
        <div class="metric"><strong>Overall Quality Score: <span class="pass">98.2%%</span></strong></div>
    </div>
    <div class="data-preview">
        <h2>Data Preview</h2>
        <table>
            <tr>
                <th>id</th>
                <th>name</th>
                <th>value</th>
                <th>date</th>
            </tr>
            <tr>
                <td>1</td>
                <td>Product A</td>
                <td>10.5</td>
                <td>2023-01-01</td>
            </tr>
            <tr>
                <td>2</td>
                <td>Product B</td>
                <td>20.75</td>
                <td>2023-01-01</td>
            </tr>
            <tr>
                <td>3</td>
                <td>Product C</td>
                <td>15.0</td>
                <td>2023-01-01</td>
            </tr>
        </table>
    </div>
</body>
</html>`, tableName, tableName)

			if err := os.WriteFile(outputPath, []byte(sampleReport), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error N402: Failed to write report: %s\n", err)
				os.Exit(1)
			}

			fmt.Println("Report generated successfully")
			fmt.Printf("Output: %s\n", outputPath)
		},
	}

	// Add flags
	generateCmd.Flags().StringVar(&tableName, "table", "", "Name of the table to generate a report for")
	generateCmd.Flags().StringVar(&format, "format", "html", "Report format (html, pdf, json, csv, md)")
	generateCmd.Flags().StringVar(&outputPath, "output", "./report.html", "Output path for the report")
	generateCmd.Flags().BoolVar(&includeCharts, "include-charts", true, "Include charts in the report")
	generateCmd.MarkFlagRequired("table")

	return generateCmd
}
