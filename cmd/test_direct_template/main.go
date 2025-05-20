package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
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

	// Generate test reports for each type
	generateQualityReport(templatesDir, outputDir)
	generateSchemaReport(templatesDir, outputDir)
	generateFreshnessReport(templatesDir, outputDir)
	generatePerformanceReport(templatesDir, outputDir)

	// Generate fixed versions using direct HTML files
	generateFixedQualityReport(outputDir)
	generateFixedSchemaReport(outputDir)
	generateFixedFreshnessReport(outputDir)
	generateFixedPerformanceReport(outputDir)

	fmt.Println("All enhanced reports generated successfully!")
}

func generateQualityReport(templatesDir, outputDir string) {
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
	outputPath := filepath.Join(outputDir, "direct-quality-report.html")
	generateReport(filepath.Join(templatesDir, "enhanced-quality-report.html"), data, outputPath)
}

func generateSchemaReport(templatesDir, outputDir string) {
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
	outputPath := filepath.Join(outputDir, "direct-schema-report.html")
	generateReport(filepath.Join(templatesDir, "enhanced-schema-report.html"), data, outputPath)
}

func generateFreshnessReport(templatesDir, outputDir string) {
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
	outputPath := filepath.Join(outputDir, "direct-freshness-report.html")
	generateReport(filepath.Join(templatesDir, "enhanced-freshness-report.html"), data, outputPath)
}

func generatePerformanceReport(templatesDir, outputDir string) {
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
	outputPath := filepath.Join(outputDir, "direct-performance-report.html")
	generateReport(filepath.Join(templatesDir, "enhanced-performance-report.html"), data, outputPath)
}

func generateReport(templatePath string, data map[string]interface{}, outputPath string) {
	// Parse the template
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Printf("Error parsing template %s: %v\n", templatePath, err)
		return
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputPath, err)
		return
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	fmt.Printf("Report generated: %s\n", outputPath)
}

// Functions to copy fixed HTML reports
func generateFixedQualityReport(outputDir string) {
	srcPath := filepath.Join("reports", "direct-quality-report-fixed.html")
	dstPath := filepath.Join(outputDir, "direct-quality-report-fixed.html")
	copyFixedReport(srcPath, dstPath)
}

func generateFixedSchemaReport(outputDir string) {
	srcPath := filepath.Join("reports", "direct-schema-report-fixed.html")
	dstPath := filepath.Join(outputDir, "direct-schema-report-fixed.html")
	copyFixedReport(srcPath, dstPath)
}

func generateFixedFreshnessReport(outputDir string) {
	srcPath := filepath.Join("reports", "direct-freshness-report-fixed.html")
	dstPath := filepath.Join(outputDir, "direct-freshness-report-fixed.html")
	copyFixedReport(srcPath, dstPath)
}

func generateFixedPerformanceReport(outputDir string) {
	srcPath := filepath.Join("reports", "direct-performance-report-fixed.html")
	dstPath := filepath.Join(outputDir, "direct-performance-report-fixed.html")
	copyFixedReport(srcPath, dstPath)
}

func copyFixedReport(srcPath, dstPath string) {
	// Check if source file exists
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		fmt.Printf("Source file %s does not exist, skipping\n", srcPath)
		return
	}

	// Read the source file
	content, err := os.ReadFile(srcPath)
	if err != nil {
		fmt.Printf("Error reading source file %s: %v\n", srcPath, err)
		return
	}

	// Write to destination file
	err = os.WriteFile(dstPath, content, 0644)
	if err != nil {
		fmt.Printf("Error writing to destination file %s: %v\n", dstPath, err)
		return
	}

	fmt.Printf("Fixed report copied: %s\n", dstPath)
}
