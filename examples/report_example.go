package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nessi-dev/nessi/internal/report"
)

func main() {
	// Check if table path is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run report_example.go <table_path> [output_format] [output_path]")
		fmt.Println("  output_format: html, pdf, json, csv (default: html)")
		fmt.Println("  output_path: directory to save the report (default: ./reports)")
		os.Exit(1)
	}

	// Get parameters
	tablePath := os.Args[1]
	format := "html"
	outputPath := "./reports"

	if len(os.Args) > 2 {
		format = os.Args[2]
	}

	if len(os.Args) > 3 {
		outputPath = os.Args[3]
	}

	// Create a report manager
	manager, err := report.NewManager(outputPath)
	if err != nil {
		fmt.Printf("Error creating report manager: %v\n", err)
		os.Exit(1)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Register the template
	templateContent, err := os.ReadFile("examples/templates/quality_report.html")
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		os.Exit(1)
	}

	template := report.ReportTemplate{
		ID:          "quality_report",
		Name:        "Quality Report",
		Description: "A comprehensive quality report for Delta Lake tables",
		Format:      report.HTML,
		Template:    string(templateContent),
		Parameters:  []string{"table_path", "timestamp"},
	}

	if err := manager.AddTemplate(template); err != nil {
		fmt.Printf("Error adding template: %v\n", err)
		os.Exit(1)
	}

	// Set up report parameters
	params := map[string]interface{}{
		"table_path": tablePath,
		"timestamp":  fmt.Sprintf("%d", time.Now().Unix()),
		"format":     format,
	}

	// Generate the report
	fmt.Printf("Generating %s report for table: %s\n", format, tablePath)
	report, err := manager.GenerateReport(context.Background(), "quality_report", params)
	if err != nil {
		fmt.Printf("Error generating report: %v\n", err)
		os.Exit(1)
	}

	// Save the report
	err = manager.SaveReport(report)
	if err != nil {
		fmt.Printf("Error saving report: %v\n", err)
		os.Exit(1)
	}

	// Construct the file path
	filePath := fmt.Sprintf("%s/%s.%s", outputPath, report.ID, report.Format)

	fmt.Printf("Report successfully generated and saved to: %s\n", filePath)
}
