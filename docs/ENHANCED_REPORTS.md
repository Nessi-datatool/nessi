# Nessi Enhanced Reports

## Overview

This document provides information about the enhanced report styling implemented for the Nessi project. The enhanced reports feature a modern, responsive design with a blue color scheme, improved typography, and better layout for report components.

## Report Types

The following enhanced report templates are available:

1. **Quality Report** - Displays data quality metrics with visual indicators
2. **Schema Report** - Shows schema details and history
3. **Freshness Report** - Provides data freshness metrics and partition information
4. **Performance Report** - Displays performance metrics and optimization recommendations

## How to Use

To generate enhanced reports, use the `ReportGenerator` with the appropriate report type prefix:

```go
// Create report generator
generator, err := report.NewReportGenerator(templatesDir, outputDir)
if err != nil {
    // Handle error
}

// Generate enhanced quality report
outputFile, err := generator.GenerateReport(data, "enhanced-quality", report.HTML, outputPath)
```

The available report type prefixes are:
- `enhanced-quality`
- `enhanced-schema`
- `enhanced-freshness`
- `enhanced-performance`

## Styling

The enhanced reports use a custom CSS file located at `assets/styles/nessi-report.css`. This file contains all the necessary styles for the reports, including:

- Color scheme (blue-based)
- Typography (using Inter font)
- Layout components (cards, tables, progress bars)
- Responsive design

## Dependencies

The enhanced reports rely on the following external resources:

- **Font Awesome** - For icons (loaded from CDN)
- **Google Fonts** - For the Inter font family (loaded from CDN)

## Testing

A test program is available at `cmd/test_enhanced_report/main.go` that demonstrates how to generate all types of enhanced reports with sample data.

## Customization

To customize the reports:

1. Modify the CSS variables in `assets/styles/nessi-report.css` to change colors and styling
2. Edit the HTML templates in `examples/templates/` to change the structure
3. Update the JavaScript in the templates for interactive features

## Example

```go
// Generate mock data for quality report
data := map[string]interface{}{
    "table_path":     "/path/to/delta/table",
    "timestamp":     "2025-05-20 21:00:00",
    "quality_score": 92,
    // Add more data fields as needed
}

// Generate report
outputPath := filepath.Join("reports", "enhanced-quality-report-output.html")
outputFile, err := generator.GenerateReport(data, "enhanced-quality", report.HTML, outputPath)
```
