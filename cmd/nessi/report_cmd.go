package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi/pkg/datalake"
	report "github.com/nessi-dev/nessi/pkg/report"
	"github.com/spf13/cobra"
)

func init() {
	var reportCmd = &cobra.Command{
		Use:   "report",
		Short: "Generate reports for data quality, schema validation, and performance metrics",
		Long: `Generate reports for data quality, schema validation, and performance metrics in various formats.

Supported formats:
- HTML: Static HTML files viewable in any browser
- PDF: Portable document format for sharing and printing
- JSON: Machine-readable format for integration with other tools
- CSV: Tabular format for analysis in spreadsheet applications

Examples:
  # Generate a quality report for a table in HTML format (default)
  nessi report --table path/to/table

  # Generate a quality report in PDF format
  nessi report --table path/to/table --format pdf

  # Generate a quality report in JSON format
  nessi report --table path/to/table --format json

  # Generate a quality report in CSV format
  nessi report --table path/to/table --format csv

  # Specify a custom output directory
  nessi report --table path/to/table --output /path/to/reports

  # Generate a freshness report
  nessi report --table path/to/table --type freshness`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReportCmd(cmd, args)
		},
	}

	// Add flags
	reportCmd.Flags().String("table", "", "Path to the table (required)")
	reportCmd.Flags().String("format", "html", "Report format (html, pdf, json, csv)")
	reportCmd.Flags().String("output", "", "Output path for the report")
	reportCmd.Flags().String("type", "quality", "Report type (quality, freshness, schema, performance)")
	reportCmd.Flags().String("template", "", "Custom template file path")

	// Mark required flags
	reportCmd.MarkFlagRequired("table")

	// Add to root command
	rootCmd.AddCommand(reportCmd)
}

func runReportCmd(cmd *cobra.Command, args []string) error {
	// Get flags
	tablePath, _ := cmd.Flags().GetString("table")
	formatStr, _ := cmd.Flags().GetString("format")
	outputPath, _ := cmd.Flags().GetString("output")
	reportType, _ := cmd.Flags().GetString("type")
	templateFile, _ := cmd.Flags().GetString("template")

	// Validate table path
	if _, err := os.Stat(tablePath); os.IsNotExist(err) {
		return fmt.Errorf("table path does not exist: %s", tablePath)
	}

	// Convert format string to ReportFormat
	var format report.ReportFormat
	switch formatStr {
	case "html":
		format = report.HTML
	case "pdf":
		format = report.PDF
	case "json":
		format = report.JSON
	case "csv":
		format = report.CSV
	default:
		return fmt.Errorf("unsupported report format: %s", formatStr)
	}

	// Create output directory if not specified
	if outputPath == "" {
		outputPath = filepath.Join(".", "reports")
	}

	// Create templates directory
	templatesDir := filepath.Join(outputPath, "templates")
	if templateFile != "" {
		// Use the directory of the template file
		templatesDir = filepath.Dir(templateFile)
	}

	// Create report generator
	generator, err := report.NewReportGenerator(templatesDir, outputPath)
	if err != nil {
		return fmt.Errorf("failed to create report generator: %w", err)
	}

	// Generate report based on type
	switch reportType {
	case "quality":
		return generateQualityReport(generator, tablePath, format, outputPath)
	case "freshness":
		return generateFreshnessReport(generator, tablePath, format, outputPath)
	case "schema":
		return generateSchemaReport(generator, tablePath, format, outputPath)
	case "performance":
		return generatePerformanceReport(generator, tablePath, format, outputPath)
	default:
		return fmt.Errorf("unsupported report type: %s", reportType)
	}
}

func generateQualityReport(generator *report.ReportGenerator, tablePath string, format report.ReportFormat, outputPath string) error {
	// Create a mock scan result for now
	// In a real implementation, we would run a scan on the table
	scanResult := map[string]interface{}{
		"scan_id":    fmt.Sprintf("scan-%d", time.Now().UnixNano()),
		"table_path": tablePath,
		"timestamp":  time.Now().Format(time.RFC3339),
		"quality_metrics": map[string]interface{}{
			"row_count":       1000,
			"null_count":      50,
			"duplicate_count": 10,
			"outlier_count":   5,
			"quality_score":   0.95,
		},
		"performance_metrics": map[string]interface{}{
			"scan_duration_ms":    5000,
			"memory_usage_mb":     100,
			"cpu_usage_percent":   50,
			"throughput_rows_sec": 200,
			"io_operations":       1000,
			"rows_processed":      1000,
			"bytes_processed":     1000000,
			"start_time":          time.Now().Add(-time.Second * 5).Format(time.RFC3339),
			"end_time":            time.Now().Format(time.RFC3339),
		},
	}

	// Generate the report
	outputFile := filepath.Join(outputPath, fmt.Sprintf("quality_report_%s.%s", time.Now().Format("20060102_150405"), format))
	resultPath, err := generator.GenerateReport(scanResult, "quality", format, outputFile)
	if err != nil {
		return fmt.Errorf("failed to generate quality report: %w", err)
	}

	fmt.Printf("Quality report generated: %s\n", resultPath)
	return nil
}

func generateFreshnessReport(generator *report.ReportGenerator, tablePath string, format report.ReportFormat, outputPath string) error {
	// Create a mock freshness report for now
	// In a real implementation, we would analyze the table for freshness
	freshnessData := map[string]interface{}{
		"table_path":  tablePath,
		"timestamp":   time.Now().Format(time.RFC3339),
		"last_update": time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
		"freshness_metrics": map[string]interface{}{
			"age_hours":             24,
			"update_frequency_days": 1,
			"freshness_score":       0.85,
			"status":                "OK",
		},
		"partitions": []map[string]interface{}{
			{
				"name":                  "partition1",
				"last_update":           time.Now().Add(-time.Hour * 12).Format(time.RFC3339),
				"age_hours":             12,
				"update_frequency_days": 1,
				"freshness_score":       0.95,
				"status":                "OK",
			},
			{
				"name":                  "partition2",
				"last_update":           time.Now().Add(-time.Hour * 36).Format(time.RFC3339),
				"age_hours":             36,
				"update_frequency_days": 1,
				"freshness_score":       0.75,
				"status":                "WARNING",
			},
		},
	}

	// Generate the report
	outputFile := filepath.Join(outputPath, fmt.Sprintf("freshness_report_%s.%s", time.Now().Format("20060102_150405"), format))
	resultPath, err := generator.GenerateReport(freshnessData, "freshness", format, outputFile)
	if err != nil {
		return fmt.Errorf("failed to generate freshness report: %w", err)
	}

	fmt.Printf("Freshness report generated: %s\n", resultPath)
	return nil
}

func generateSchemaReport(generator *report.ReportGenerator, tablePath string, format report.ReportFormat, outputPath string) error {
	// Detect if the table is a Delta Lake table
	deltaHandler := datalake.NewDeltaFormatHandler()
	isDelta := deltaHandler.IsDeltaTable(tablePath)

	// Create a mock schema report for now
	// In a real implementation, we would extract the schema from the table
	schemaData := map[string]interface{}{
		"table_path":     tablePath,
		"timestamp":      time.Now().Format(time.RFC3339),
		"format":         "delta",
		"is_delta_table": isDelta,
		"schema_version": 1,
		"fields": []map[string]interface{}{
			{
				"name":     "id",
				"type":     "long",
				"nullable": false,
			},
			{
				"name":     "name",
				"type":     "string",
				"nullable": true,
			},
			{
				"name":     "age",
				"type":     "integer",
				"nullable": true,
			},
			{
				"name":     "active",
				"type":     "boolean",
				"nullable": true,
			},
		},
	}

	// Generate the report
	outputFile := filepath.Join(outputPath, fmt.Sprintf("schema_report_%s.%s", time.Now().Format("20060102_150405"), format))
	resultPath, err := generator.GenerateReport(schemaData, "schema", format, outputFile)
	if err != nil {
		return fmt.Errorf("failed to generate schema report: %w", err)
	}

	fmt.Printf("Schema report generated: %s\n", resultPath)
	return nil
}

func generatePerformanceReport(generator *report.ReportGenerator, tablePath string, format report.ReportFormat, outputPath string) error {
	// Create a mock performance report for now
	// In a real implementation, we would collect performance metrics from the table
	performanceData := map[string]interface{}{
		"table_path": tablePath,
		"timestamp":  time.Now().Format(time.RFC3339),
		"metrics": map[string]interface{}{
			"scan_duration_ms":    5000,
			"memory_usage_mb":     100,
			"cpu_usage_percent":   50,
			"io_operations":       1000,
			"rows_processed":      1000,
			"bytes_processed":     1000000,
			"throughput_rows_sec": 200,
			"throughput_mb_sec":   0.2,
		},
		"history": []map[string]interface{}{
			{
				"timestamp":           time.Now().Add(-time.Hour * 24).Format(time.RFC3339),
				"scan_duration_ms":    4800,
				"memory_usage_mb":     95,
				"cpu_usage_percent":   48,
				"rows_processed":      1000,
				"throughput_rows_sec": 208,
			},
			{
				"timestamp":           time.Now().Add(-time.Hour * 48).Format(time.RFC3339),
				"scan_duration_ms":    5200,
				"memory_usage_mb":     105,
				"cpu_usage_percent":   52,
				"rows_processed":      1000,
				"throughput_rows_sec": 192,
			},
		},
	}

	// Generate the report
	outputFile := filepath.Join(outputPath, fmt.Sprintf("performance_report_%s.%s", time.Now().Format("20060102_150405"), format))
	resultPath, err := generator.GenerateReport(performanceData, "performance", format, outputFile)
	if err != nil {
		return fmt.Errorf("failed to generate performance report: %w", err)
	}

	fmt.Printf("Performance report generated: %s\n", resultPath)
	return nil
}
