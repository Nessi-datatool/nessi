# Nessi Reporting Capabilities

## Navigation

- [Documentation Home](README.md)
- [User Guide](USER_GUIDE.md)
- [Installation & Quickstart](../QUICKSTART.md)
- [Configuration](CONFIGURATION.md)
- [CLI Reference](cli/README.md)
- [Quality Rules](QUALITY_RULES.md)
- [Integrations](integration_guide.md)
- [Developer Guide](developer_experience.md)
- [FAQ](faq.md)

---

Nessi provides comprehensive reporting capabilities through its flexible CLI-based report generation system. This document outlines the available report formats and how to use them.

## Supported Report Formats

Nessi supports the following report formats:

1. **HTML Reports**: Static HTML files that can be viewed in any web browser
2. **PDF Reports**: Static reports suitable for sharing and archiving
3. **JSON Reports**: Machine-readable reports for integration with other systems
4. **CSV Reports**: Tabular data exports for analysis in spreadsheet applications

## Report Generation

Reports are generated using the Nessi CLI:

```bash
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
nessi freshness report --table my_table --format pdf
```

## Command Options

The `report` command supports the following options:

| Option | Description | Default |
|--------|-------------|---------|
| `--table` | Path to the Delta Lake table (required) | - |
| `--format` | Report format (html, pdf, json, csv) | html |
| `--output` | Directory where the report will be saved | ./reports |
| `--template` | Path to a custom template file | - |

## Report Templates

Nessi uses templates for report generation, which can be customized to match your organization's branding and requirements. Default templates are provided for common report types:

- Quality reports
- Freshness reports
- Schema validation reports
- Performance metrics reports

## Customizing Reports

To customize report templates:

1. Create a new template file in HTML, JSON, or CSV format
2. Modify the template according to your needs
3. Specify the custom template path when generating reports:

```bash
nessi report --table path/to/table --template path/to/custom/template.html
```

### HTML Template Variables

HTML templates can use the following variables:

- `{{.table_path}}`: Path to the Delta Lake table
- `{{.timestamp}}`: Timestamp when the report was generated
- `{{.data}}`: The report data object containing metrics, columns, and other information

## Programmatic Report Generation

Reports can also be generated programmatically using the Nessi API:

```go
import "github.com/nessi-dev/nessi/internal/report"

// Create a report manager
manager, err := report.NewManager("./reports")
if err != nil {
    // Handle error
}

// Register a template
template := report.ReportTemplate{
    ID:          "quality_report",
    Name:        "Quality Report",
    Description: "A comprehensive quality report",
    Format:      report.HTML,
    Template:    templateContent,
    Parameters:  []string{"table_path", "timestamp"},
}
manager.AddTemplate(template)

// Set up parameters
params := map[string]interface{}{
    "table_path": "path/to/table",
    "timestamp":  time.Now().Format(time.RFC3339),
    "data":       qualityData,
}

// Generate a report
report, err := manager.GenerateReport(ctx, "quality_report", params)
if err != nil {
    // Handle error
}

// Save the report
err = manager.SaveReport(report)
if err != nil {
    // Handle error
}
```

## Report Storage

By default, reports are stored in the configured output directory (default: `./reports`). The report manager handles report storage and retrieval, including:

- Generating unique report IDs based on timestamp
- Organizing reports by format (file extension)
- Creating the output directory if it doesn't exist

## Integration with External Systems

The CLI-based reporting system is designed for easy integration with external systems:

- **Automation**: Reports can be generated automatically using cron jobs or CI/CD pipelines
- **Data Pipelines**: Include report generation as a step in your data processing pipelines
- **File-Based Monitoring**: Generate reports on a schedule and use file watchers to trigger notifications
- **Documentation**: Include reports in your data documentation systems by copying the generated files

## Examples

### Generate a Quality Report for a Production Table

```bash
nessi report --table /data/production/customer_orders --format html --output /reports/weekly
```

### Generate a JSON Report for API Integration

```bash
nessi report --table /data/api/products --format json --output /api/reports
```

### Use a Custom Template

```bash
nessi report --table /data/marketing/campaigns --template /templates/marketing_report.html
```
