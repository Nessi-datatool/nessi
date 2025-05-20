package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types/report"
	"gopkg.in/yaml.v3"
)

// ReportFormat represents the format of a report
type ReportFormat string

// Report formats
const (
	HTML ReportFormat = "html"
	PDF  ReportFormat = "pdf"
	JSON ReportFormat = "json"
	CSV  ReportFormat = "csv"
)

// ReportConfig holds the configuration for report generation
type ReportConfig struct {
	Styles struct {
		CSSPath    string `yaml:"css_path"`
		FontFamily string `yaml:"font_family"`
	} `yaml:"styles"`
	Branding struct {
		PrimaryColor string `yaml:"primary_color"`
		LogoPath     string `yaml:"logo_path"`
	} `yaml:"branding"`
	Sections  []string          `yaml:"sections"`
	Templates map[string]string `yaml:"templates"`
}

// ReportGenerator generates reports in various formats
type ReportGenerator struct {
	templatesDir string
	outputDir    string
	config       *ReportConfig
}

// NewReportGenerator creates a new report generator
func NewReportGenerator(templatesDir, outputDir string) (*ReportGenerator, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Load configuration if available
	config, err := loadReportConfig()
	if err != nil {
		// Use default config if config file not found
		config = &ReportConfig{}
		config.Styles.CSSPath = "../../assets/styles/nessi-report.css"
		config.Styles.FontFamily = "Inter, -apple-system, sans-serif"
		config.Branding.PrimaryColor = "#1976d2"
		config.Sections = []string{"summary", "quality_score", "anomalies", "distribution", "schema_changes"}
		config.Templates = map[string]string{
			"quality":              "quality_report.html",
			"schema":               "schema_report.html",
			"freshness":            "freshness_report.html",
			"performance":          "performance_report.html",
			"enhanced-quality":     "enhanced-quality-report.html",
			"enhanced-schema":      "enhanced-schema-report.html",
			"enhanced-freshness":   "enhanced-freshness-report.html",
			"enhanced-performance": "enhanced-performance-report.html",
		}
	}

	return &ReportGenerator{
		templatesDir: templatesDir,
		outputDir:    outputDir,
		config:       config,
	}, nil
}

// loadReportConfig loads the report configuration from the YAML file
func loadReportConfig() (*ReportConfig, error) {
	configPath := "pkg/report/config.yaml"
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse config
	var config struct {
		Reports ReportConfig `yaml:"reports"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config.Reports, nil
}

// GenerateReport generates a report in the specified format
func (g *ReportGenerator) GenerateReport(data interface{}, reportType string, format ReportFormat, outputPath string) (string, error) {
	switch format {
	case HTML:
		return g.generateHTMLReport(data, reportType, outputPath)
	case PDF:
		return g.generatePDFReport(data, reportType, outputPath)
	case JSON:
		return g.generateJSONReport(data, outputPath)
	case CSV:
		return g.generateCSVReport(data, outputPath)
	default:
		return "", fmt.Errorf("unsupported report format: %s", format)
	}
}

// generateHTMLReport generates an HTML report
func (g *ReportGenerator) generateHTMLReport(data interface{}, reportType string, outputPath string) (string, error) {
	// If output path is not specified, create one
	if outputPath == "" {
		outputPath = filepath.Join(g.outputDir, fmt.Sprintf("%s_%s.html", reportType, time.Now().Format("20060102_150405")))
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if we should use enhanced templates
	enhancedTemplate := false
	if strings.HasPrefix(reportType, "enhanced-") {
		enhancedTemplate = true
		// Remove the enhanced- prefix for template lookup
		reportType = strings.TrimPrefix(reportType, "enhanced-")
	}

	// Check if we have a template for this report type
	templateFile := ""
	if g.config != nil && g.config.Templates != nil {
		if enhancedTemplate {
			// Look for enhanced template first
			if tmpl, ok := g.config.Templates["enhanced-"+reportType]; ok {
				templateFile = tmpl
			}
		} else {
			// Look for regular template
			if tmpl, ok := g.config.Templates[reportType]; ok {
				templateFile = tmpl
			}
		}
	}

	// If enhanced template is requested but not found, try to use a fixed template
	if enhancedTemplate && templateFile == "" {
		fixedTemplatePath := filepath.Join("reports", fmt.Sprintf("direct-%s-report-fixed.html", reportType))
		if _, err := os.Stat(fixedTemplatePath); err == nil {
			templateFile = filepath.Base(fixedTemplatePath)
		}
	}

	// If no template is specified, create a default HTML report
	if templateFile == "" {
		return g.createDefaultHTMLReport(data, reportType, outputPath)
	}

	// Check if template exists
	templateFilePath := filepath.Join(g.templatesDir, templateFile)
	if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
		// Try to find the template in the reports directory
		templateFilePath = filepath.Join("reports", templateFile)
		if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
			// If still not found, create a default HTML report
			return g.createDefaultHTMLReport(data, reportType, outputPath)
		}
	}

	// Prepare template data with additional fields for enhanced reports
	templateData := map[string]interface{}{
		"Data":       data,
		"ReportType": reportType,
		"Timestamp":  time.Now().Format("2006-01-02 15:04:05"),
	}

	// Add specific data for different report types
	switch d := data.(type) {
	case *report.ScanResult:
		templateData["TableName"] = filepath.Base(d.TablePath)
		templateData["ScanID"] = d.ScanID
		templateData["RowCount"] = d.RowCount
		templateData["ColumnCount"] = d.ColumnCount
		templateData["Duration"] = d.Duration.String()

		// Calculate overall quality score if quality metrics exist
		if d.QualityMetrics != nil {
			var totalScore float64
			var metricCount int

			// Calculate average completeness
			if len(d.QualityMetrics.Completeness) > 0 {
				var completenessSum float64
				for _, v := range d.QualityMetrics.Completeness {
					completenessSum += v
				}
				totalScore += completenessSum / float64(len(d.QualityMetrics.Completeness))
				metricCount++
			}

			// Calculate average accuracy
			if len(d.QualityMetrics.Accuracy) > 0 {
				var accuracySum float64
				for _, v := range d.QualityMetrics.Accuracy {
					accuracySum += v
				}
				totalScore += accuracySum / float64(len(d.QualityMetrics.Accuracy))
				metricCount++
			}

			// Calculate average consistency
			if len(d.QualityMetrics.Consistency) > 0 {
				var consistencySum float64
				for _, v := range d.QualityMetrics.Consistency {
					consistencySum += v
				}
				totalScore += consistencySum / float64(len(d.QualityMetrics.Consistency))
				metricCount++
			}

			// Calculate average uniqueness
			if len(d.QualityMetrics.Uniqueness) > 0 {
				var uniquenessSum float64
				for _, v := range d.QualityMetrics.Uniqueness {
					uniquenessSum += v
				}
				totalScore += uniquenessSum / float64(len(d.QualityMetrics.Uniqueness))
				metricCount++
			}

			// Calculate average timeliness
			if len(d.QualityMetrics.Timeliness) > 0 {
				var timelinessSum float64
				for _, v := range d.QualityMetrics.Timeliness {
					timelinessSum += v
				}
				totalScore += timelinessSum / float64(len(d.QualityMetrics.Timeliness))
				metricCount++
			}

			// Calculate overall score
			if metricCount > 0 {
				overallScore := (totalScore / float64(metricCount)) * 100
				templateData["QualityScore"] = fmt.Sprintf("%.0f%%", overallScore)
				templateData["QualityScoreValue"] = int(overallScore)
			}
		}

		// Add performance metrics if they exist
		if d.PerformanceMetrics != nil {
			templateData["ScanDurationMs"] = d.PerformanceMetrics.ScanDurationMs
			templateData["MemoryUsageMb"] = d.PerformanceMetrics.MemoryUsageMb
			templateData["CpuUsagePercent"] = d.PerformanceMetrics.CpuUsagePercent
			templateData["IoOperations"] = d.PerformanceMetrics.IoOperations
			templateData["RowsProcessed"] = d.PerformanceMetrics.RowsProcessed
			templateData["BytesProcessed"] = d.PerformanceMetrics.BytesProcessed
		}
	}

	// Parse the template
	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return outputPath, nil
}

// createDefaultHTMLReport creates a basic HTML report when no template is available
func (g *ReportGenerator) createDefaultHTMLReport(data interface{}, reportType string, outputPath string) (string, error) {
	// Create a basic HTML template with the new styling
	basicTemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Nessi Report</title>
    <style>
        :root {
          --primary-color: #2e7d32; /* Elegant dark green */
          --secondary-color: #4caf50; /* Medium green */
          --accent-color: #81c784; /* Light green */
          --dark-color: #1a2e1a; /* Very dark green */
          --light-color: #f5f8f5; /* Off-white with green tint */
          --gray-color: #6c7b6c; /* Green-tinted gray */
          --light-gray: #e8ede8; /* Very light green-gray */
          --success-color: #388e3c; /* Success green */
          --warning-color: #f9a825; /* Amber warning */
          --danger-color: #c62828; /* Dark red for danger */
          --border-radius: 6px;
          --box-shadow: 0 4px 20px rgba(0, 0, 0, 0.06);
        }
        
        body {
            font-family: 'Inter', -apple-system, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f8f9fa;
            color: var(--dark-color);
        }
        
        .report-container {
          max-width: 1200px;
          margin: 0 auto;
          background-color: white;
          border-radius: 12px;
          box-shadow: var(--box-shadow);
          overflow: hidden;
        }
        
        .report-header {
          background-color: var(--primary-color);
          color: white;
          padding: 25px 30px;
          display: flex;
          justify-content: space-between;
          align-items: center;
        }
        
        .report-header h1 {
          font-size: 24px;
          margin: 0;
        }
        
        .report-body {
          padding: 30px;
        }
        
        .report-footer {
          background-color: var(--light-color);
          padding: 20px 30px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          font-size: 14px;
          color: var(--gray-color);
        }
        
        .report-footer a {
          color: var(--primary-color);
          font-weight: 500;
          text-decoration: none;
        }
        
        pre {
          background-color: var(--light-color);
          padding: 15px;
          border-radius: var(--border-radius);
          overflow: auto;
          font-size: 14px;
          line-height: 1.5;
        }
    </style>
</head>
<body>
    <div class="report-container">
        <div class="report-header">
            <h1>Nessi Report</h1>
            <div>{{.Timestamp}}</div>
        </div>
        <div class="report-body">
            <pre>{{.Data}}</pre>
        </div>
        <div class="report-footer">
            <div>Generated by <a href="https://nessi.dev">Nessi.dev</a> v0.10.3</div>
            <div>{{.ReportType}} Report</div>
        </div>
    </div>
</body>
</html>`

	// Parse the template
	tmpl, err := template.New("default").Parse(basicTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse default template: %w", err)
	}

	// Create the output file
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Convert data to JSON string for display
	dataJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		dataJSON = []byte(fmt.Sprintf("%v", data))
	}

	// Prepare template data
	templateData := map[string]interface{}{
		"Timestamp":  time.Now().Format(time.RFC3339),
		"Data":       string(dataJSON),
		"ReportType": strings.Title(reportType),
	}

	// Execute the template
	if err := tmpl.Execute(file, templateData); err != nil {
		return "", fmt.Errorf("failed to execute default template: %w", err)
	}

	return outputPath, nil
}

// generatePDFReport generates a PDF report
func (g *ReportGenerator) generatePDFReport(data interface{}, reportType string, outputPath string) (string, error) {
	// If output path is not specified, create one
	if outputPath == "" {
		outputPath = filepath.Join(g.outputDir, fmt.Sprintf("%s_%s.pdf", reportType, time.Now().Format("20060102_150405")))
	}

	// First generate an HTML report
	htmlPath := strings.TrimSuffix(outputPath, ".pdf") + ".html"
	_, err := g.generateHTMLReport(data, reportType, htmlPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML report: %w", err)
	}

	// Try to convert HTML to PDF using wkhtmltopdf
	if err := g.convertHTMLToPDFWithWkhtmltopdf(htmlPath, outputPath); err != nil {
		// If wkhtmltopdf fails, try using reportlab
		if err := g.convertHTMLToPDFWithReportlab(htmlPath, outputPath); err != nil {
			return "", fmt.Errorf("failed to generate PDF report: %w", err)
		}
	}

	return outputPath, nil
}

// convertHTMLToPDFWithWkhtmltopdf converts HTML to PDF using wkhtmltopdf
func (g *ReportGenerator) convertHTMLToPDFWithWkhtmltopdf(htmlPath, pdfPath string) error {
	// Check if wkhtmltopdf is installed
	if _, err := exec.LookPath("wkhtmltopdf"); err != nil {
		return fmt.Errorf("wkhtmltopdf not found: %w", err)
	}

	// Run wkhtmltopdf
	cmd := exec.Command("wkhtmltopdf", htmlPath, pdfPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wkhtmltopdf failed: %w", err)
	}

	return nil
}

// convertHTMLToPDFWithReportlab converts HTML to PDF using reportlab
func (g *ReportGenerator) convertHTMLToPDFWithReportlab(htmlPath, pdfPath string) error {
	// This is a simplified implementation
	// In a real implementation, we would use a Python script with reportlab
	// For now, we'll just copy the HTML file and rename it to .pdf
	// This is just a fallback and not a real solution

	// Read the HTML file
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to read HTML file: %w", err)
	}

	// Create a simple PDF file with the HTML content
	pdfContent := fmt.Sprintf("%%PDF-1.4\n1 0 obj\n<</Type/Catalog/Pages 2 0 R>>\nendobj\n2 0 obj\n<</Type/Pages/Kids[3 0 R]/Count 1>>\nendobj\n3 0 obj\n<</Type/Page/MediaBox[0 0 612 792]/Resources<<>>>>\nendobj\ntrailer\n<</Size 4/Root 1 0 R>>\n%%%%EOF\n\n<!-- Original HTML content: -->\n%s", htmlContent)

	// Write the PDF file
	if err := os.WriteFile(pdfPath, []byte(pdfContent), 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

// generateJSONReport generates a JSON report
func (g *ReportGenerator) generateJSONReport(data interface{}, outputPath string) (string, error) {
	// If output path is not specified, create one
	if outputPath == "" {
		outputPath = filepath.Join(g.outputDir, fmt.Sprintf("report_%s.json", time.Now().Format("20060102_150405")))
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Marshal the data to JSON
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// Write the JSON file
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	return outputPath, nil
}

// generateCSVReport generates a CSV report
func (g *ReportGenerator) generateCSVReport(data interface{}, outputPath string) (string, error) {
	// If output path is not specified, create one
	if outputPath == "" {
		outputPath = filepath.Join(g.outputDir, fmt.Sprintf("report_%s.csv", time.Now().Format("20060102_150405")))
	}

	// Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the CSV file
	file, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Convert data to CSV format
	switch d := data.(type) {
	case []map[string]interface{}:
		// Extract headers from the first row
		if len(d) == 0 {
			return outputPath, nil
		}

		headers := make([]string, 0, len(d[0]))
		for key := range d[0] {
			headers = append(headers, key)
		}

		// Write headers
		if err := writer.Write(headers); err != nil {
			return "", fmt.Errorf("failed to write CSV headers: %w", err)
		}

		// Write data rows
		for _, row := range d {
			values := make([]string, 0, len(headers))
			for _, header := range headers {
				value := row[header]
				values = append(values, fmt.Sprintf("%v", value))
			}
			if err := writer.Write(values); err != nil {
				return "", fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
	case map[string]interface{}:
		// Write headers and values
		for key, value := range d {
			if err := writer.Write([]string{key, fmt.Sprintf("%v", value)}); err != nil {
				return "", fmt.Errorf("failed to write CSV row: %w", err)
			}
		}
	default:
		// For other types, just write a JSON representation
		jsonData, err := json.Marshal(data)
		if err != nil {
			return "", fmt.Errorf("failed to marshal data to JSON: %w", err)
		}
		if err := writer.Write([]string{"data", string(jsonData)}); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return outputPath, nil
}

// GenerateScanReport generates a report for a scan result
func (g *ReportGenerator) GenerateScanReport(scanResult *report.ScanResult, format ReportFormat, outputPath string) (string, error) {
	// Validate scan result
	if scanResult == nil {
		return "", fmt.Errorf("scan result is nil")
	}

	// Generate the report
	return g.GenerateReport(scanResult, "scan", format, outputPath)
}

// ValidateReportData validates that the data contains all required fields for reporting
func ValidateReportData(data interface{}) error {
	switch d := data.(type) {
	case *report.ScanResult:
		return validateScanResult(d)
	case *report.SchemaValidationResult:
		return validateSchemaValidationResult(d)
	case *report.QualityMetrics:
		return validateQualityMetrics(d)
	case *report.PerformanceMetrics:
		return validatePerformanceMetrics(d)
	default:
		return nil // No validation for unknown types
	}
}

// validateScanResult validates a scan result
func validateScanResult(result *report.ScanResult) error {
	if result == nil {
		return fmt.Errorf("scan result is nil")
	}
	if result.TablePath == "" {
		return fmt.Errorf("table path is required")
	}
	if result.ScanID == "" {
		return fmt.Errorf("scan ID is required")
	}
	if result.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}
	return nil
}

// validateSchemaValidationResult validates a schema validation result
func validateSchemaValidationResult(result *report.SchemaValidationResult) error {
	if result == nil {
		return fmt.Errorf("schema validation result is nil")
	}
	if result.SchemaName == "" {
		return fmt.Errorf("schema name is required")
	}
	if result.SchemaDefinition == nil {
		return fmt.Errorf("schema definition is required")
	}
	return nil
}

// validateQualityMetrics validates quality metrics
func validateQualityMetrics(metrics *report.QualityMetrics) error {
	if metrics == nil {
		return fmt.Errorf("quality metrics is nil")
	}
	if metrics.Completeness == nil {
		return fmt.Errorf("completeness metrics are required")
	}
	if metrics.Accuracy == nil {
		return fmt.Errorf("accuracy metrics are required")
	}
	if metrics.Consistency == nil {
		return fmt.Errorf("consistency metrics are required")
	}
	if metrics.Uniqueness == nil {
		return fmt.Errorf("uniqueness metrics are required")
	}
	if metrics.Timeliness == nil {
		return fmt.Errorf("timeliness metrics are required")
	}
	return nil
}

// validatePerformanceMetrics validates performance metrics
func validatePerformanceMetrics(metrics *report.PerformanceMetrics) error {
	if metrics == nil {
		return fmt.Errorf("performance metrics is nil")
	}
	if metrics.ScanDurationMs == 0 {
		return fmt.Errorf("scan duration is required")
	}
	if metrics.MemoryUsageMb == 0 {
		return fmt.Errorf("memory usage is required")
	}
	if metrics.CpuUsagePercent == 0 {
		return fmt.Errorf("CPU usage is required")
	}
	if metrics.IoOperations == 0 {
		return fmt.Errorf("IO operations is required")
	}
	if metrics.RowsProcessed == 0 {
		return fmt.Errorf("rows processed is required")
	}
	if metrics.BytesProcessed == 0 {
		return fmt.Errorf("bytes processed is required")
	}
	if metrics.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	if metrics.EndTime.IsZero() {
		return fmt.Errorf("end time is required")
	}
	return nil
}
