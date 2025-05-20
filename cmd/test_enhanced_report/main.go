package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi/pkg/report"
)

func main() {
	// Set up paths
	templatesDir := filepath.Join("examples", "templates")
	outputDir := "reports"

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Create report generator with absolute paths to ensure templates are found
	absTemplatesDir, _ := filepath.Abs(templatesDir)
	absOutputDir, _ := filepath.Abs(outputDir)
	generator, err := report.NewReportGenerator(absTemplatesDir, absOutputDir)
	if err != nil {
		fmt.Printf("Error creating report generator: %v\n", err)
		return
	}

	// Generate test reports for each type
	generateQualityReport(generator)
	generateSchemaReport(generator)
	generateFreshnessReport(generator)
	generatePerformanceReport(generator)

	fmt.Println("All enhanced reports generated successfully!")
}

func generateQualityReport(generator *report.ReportGenerator) {
	// Generate mock data for quality report
	data := map[string]interface{}{
		"table_path":     "/path/to/delta/table",
		"timestamp":      "2025-05-20 21:00:00",
		"quality_score":  92,
		"total_rows":     15482,
		"analyzed_rows":  15482,
		"valid_rows":     14243,
		"invalid_rows":   1239,
		"completeness":   98,
		"uniqueness":     95,
		"consistency":    90,
		"accuracy":       94,
		"transaction_id": "12345abcde",
		"columns": []map[string]interface{}{
			{
				"name":         "id",
				"type":         "string",
				"completeness": 100,
				"uniqueness":   100,
				"validity":     100,
				"status":       "OK",
				"null_count":   0,
				"total_count":  15482,
			},
			{
				"name":         "name",
				"type":         "string",
				"completeness": 98,
				"uniqueness":   95,
				"validity":     99,
				"status":       "OK",
				"null_count":   310,
				"total_count":  15482,
			},
			{
				"name":         "email",
				"type":         "string",
				"completeness": 95,
				"uniqueness":   99,
				"validity":     92,
				"status":       "WARNING",
				"null_count":   774,
				"total_count":  15482,
			},
			{
				"name":         "age",
				"type":         "integer",
				"completeness": 90,
				"uniqueness":   30,
				"validity":     85,
				"status":       "WARNING",
				"null_count":   1548,
				"total_count":  15482,
			},
			{
				"name":         "created_at",
				"type":         "timestamp",
				"completeness": 100,
				"uniqueness":   99,
				"validity":     100,
				"status":       "OK",
				"null_count":   0,
				"total_count":  15482,
			},
		},
		"anomalies": []map[string]interface{}{
			{
				"column":      "email",
				"type":        "validity",
				"description": "Invalid email format",
				"count":       155,
				"percentage":  1.0,
				"examples":    []string{"john.doe@", "@example.com", "invalid-email"},
			},
			{
				"column":      "age",
				"type":        "range",
				"description": "Age out of expected range (0-120)",
				"count":       42,
				"percentage":  0.3,
				"examples":    []string{"250", "999", "-5"},
			},
			{
				"column":      "name",
				"type":        "pattern",
				"description": "Unusual name patterns detected",
				"count":       89,
				"percentage":  0.6,
				"examples":    []string{"Test123", "ABCDEFG", "User1234"},
			},
		},
	}

	// Generate report
	outputPath := filepath.Join("reports", "enhanced-quality-report-output.html")
	outputFile, err := generator.GenerateReport(data, "enhanced-quality", report.HTML, outputPath)
	if err != nil {
		fmt.Printf("Error generating quality report: %v\n", err)
		return
	}

	fmt.Printf("Enhanced quality report generated: %s\n", outputFile)
}

func generateSchemaReport(generator *report.ReportGenerator) {
	// Generate mock data for schema report
	data := map[string]interface{}{
		"table_path":     "/path/to/delta/table",
		"timestamp":      "2025-05-20 21:00:00",
		"schema_version": 5,
		"last_updated":   "2025-05-15",
		"changes_count":  3,
		"transaction_id": "12345abcde",
		"fields": []map[string]interface{}{
			{
				"name":        "id",
				"type":        "string",
				"nullable":    false,
				"description": "Primary key",
			},
			{
				"name":        "name",
				"type":        "string",
				"nullable":    true,
				"description": "User's full name",
			},
			{
				"name":        "email",
				"type":        "string",
				"nullable":    true,
				"description": "User's email address",
			},
			{
				"name":        "age",
				"type":        "integer",
				"nullable":    true,
				"description": "User's age",
			},
			{
				"name":        "created_at",
				"type":        "timestamp",
				"nullable":    false,
				"description": "Record creation timestamp",
			},
			{
				"name":        "updated_at",
				"type":        "timestamp",
				"nullable":    true,
				"description": "Record update timestamp",
			},
			{
				"name":        "status",
				"type":        "string",
				"nullable":    true,
				"description": "User account status",
			},
		},
		"schema_history": []map[string]interface{}{
			{
				"version":    5,
				"type":       "added",
				"field_name": "status",
				"details":    "Added status field",
				"timestamp":  "2025-05-15",
			},
			{
				"version":    4,
				"type":       "modified",
				"field_name": "email",
				"details":    "Changed from required to nullable",
				"timestamp":  "2025-04-20",
			},
			{
				"version":    3,
				"type":       "added",
				"field_name": "updated_at",
				"details":    "Added updated_at timestamp",
				"timestamp":  "2025-03-10",
			},
			{
				"version":    2,
				"type":       "added",
				"field_name": "age",
				"details":    "Added age field",
				"timestamp":  "2025-02-15",
			},
			{
				"version":    1,
				"type":       "created",
				"field_name": "schema",
				"details":    "Initial schema creation",
				"timestamp":  "2025-01-01",
			},
		},
	}

	// Generate report
	outputPath := filepath.Join("reports", "enhanced-schema-report-output.html")
	outputFile, err := generator.GenerateReport(data, "enhanced-schema", report.HTML, outputPath)
	if err != nil {
		fmt.Printf("Error generating schema report: %v\n", err)
		return
	}

	fmt.Printf("Enhanced schema report generated: %s\n", outputFile)
}

func generateFreshnessReport(generator *report.ReportGenerator) {
	// Generate mock data for freshness report
	data := map[string]interface{}{
		"table_path":        "/path/to/delta/table",
		"timestamp":         "2025-05-20 21:00:00",
		"freshness_score":   85,
		"last_update_age":   "6 hours",
		"update_frequency":  "Daily",
		"last_update_score": 90,
		"frequency_score":   80,
		"consistency_score": 85,
		"transaction_id":    "12345abcde",
		"partitions": []map[string]interface{}{
			{
				"name":                  "date=2025-05-20",
				"last_update":           "2025-05-20 15:30:00",
				"age_hours":             6,
				"update_frequency_days": 1,
				"freshness_score":       90,
				"status":                "OK",
			},
			{
				"name":                  "date=2025-05-19",
				"last_update":           "2025-05-19 16:45:00",
				"age_hours":             28,
				"update_frequency_days": 1,
				"freshness_score":       85,
				"status":                "OK",
			},
			{
				"name":                  "date=2025-05-18",
				"last_update":           "2025-05-18 18:20:00",
				"age_hours":             51,
				"update_frequency_days": 1,
				"freshness_score":       80,
				"status":                "OK",
			},
			{
				"name":                  "date=2025-05-17",
				"last_update":           "2025-05-17 14:10:00",
				"age_hours":             79,
				"update_frequency_days": 1,
				"freshness_score":       75,
				"status":                "WARNING",
			},
			{
				"name":                  "date=2025-05-16",
				"last_update":           "2025-05-16 10:30:00",
				"age_hours":             106,
				"update_frequency_days": 1,
				"freshness_score":       60,
				"status":                "WARNING",
			},
		},
	}

	// Generate report
	outputPath := filepath.Join("reports", "enhanced-freshness-report-output.html")
	outputFile, err := generator.GenerateReport(data, "enhanced-freshness", report.HTML, outputPath)
	if err != nil {
		fmt.Printf("Error generating freshness report: %v\n", err)
		return
	}

	fmt.Printf("Enhanced freshness report generated: %s\n", outputFile)
}

func generatePerformanceReport(generator *report.ReportGenerator) {
	// Generate mock data for performance report
	data := map[string]interface{}{
		"table_path":               "/path/to/delta/table",
		"timestamp":                "2025-05-20 21:00:00",
		"performance_score":        78,
		"file_count":               1250,
		"total_size":               "2.5 GB",
		"partition_count":          5,
		"file_size_score":          70,
		"partition_score":          85,
		"z_order_score":            80,
		"transaction_id":           "12345abcde",
		"file_size_recommendation": "Consider compacting small files. Current average file size is 2.1 MB, target should be 128+ MB.",
		"file_size_command":        "OPTIMIZE delta.`/path/to/table`",
		"partition_recommendation": "Current partitioning strategy is optimal. No changes needed.",
		"partition_command":        "N/A",
		"z_order_recommendation":   "Consider Z-Ordering by 'id' and 'created_at' fields to improve query performance.",
		"z_order_command":          "OPTIMIZE delta.`/path/to/table` ZORDER BY (id, created_at)",
		"file_stats": []map[string]interface{}{
			{
				"metric":         "Small files (<1MB)",
				"value":          "450 files (36%)",
				"recommendation": "Too many small files, consider compaction",
				"status":         "WARNING",
			},
			{
				"metric":         "Medium files (1-64MB)",
				"value":          "700 files (56%)",
				"recommendation": "Consider compaction to larger files",
				"status":         "OK",
			},
			{
				"metric":         "Large files (>64MB)",
				"value":          "100 files (8%)",
				"recommendation": "Good file size distribution",
				"status":         "OK",
			},
			{
				"metric":         "Partition skew",
				"value":          "15% max difference",
				"recommendation": "Partition distribution is balanced",
				"status":         "OK",
			},
			{
				"metric":         "Z-Order status",
				"value":          "Not optimized",
				"recommendation": "Consider Z-Ordering for query performance",
				"status":         "WARNING",
			},
		},
	}

	// Generate report
	outputPath := filepath.Join("reports", "enhanced-performance-report-output.html")
	outputFile, err := generator.GenerateReport(data, "enhanced-performance", report.HTML, outputPath)
	if err != nil {
		fmt.Printf("Error generating performance report: %v\n", err)
		return
	}

	fmt.Printf("Enhanced performance report generated: %s\n", outputFile)
}
