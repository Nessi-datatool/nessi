# Nessi Reporting Capabilities

Nessi provides comprehensive reporting capabilities through its flexible report generation system. This document outlines the available report formats and how to use them.

## Supported Report Formats

Nessi supports the following report formats:

1. **HTML Reports**: Interactive reports that can be viewed in a web browser
2. **PDF Reports**: Static reports suitable for sharing and archiving
3. **JSON Reports**: Machine-readable reports for integration with other systems
4. **CSV Reports**: Tabular data exports for analysis in spreadsheet applications

## Report Generation

Reports can be generated using the Nessi CLI or programmatically through the API:

```bash
# Generate a quality report for a table
nessi report --table path/to/table --format html

# Generate a freshness report
nessi freshness report --table my_table --format pdf
```

## Report Templates

Nessi uses templates for report generation, which can be customized to match your organization's branding and requirements. Default templates are provided for common report types:

- Quality reports
- Freshness reports
- Schema validation reports
- Performance metrics reports

## Customizing Reports

To customize report templates:

1. Copy the default templates from the `templates` directory
2. Modify the templates according to your needs
3. Specify the custom template path when generating reports:

```bash
nessi report --table path/to/table --template path/to/custom/template.html
```

## Programmatic Report Generation

Reports can also be generated programmatically using the Nessi API:

```go
import "github.com/nessi-dev/nessi/internal/report"

// Create a report manager
manager, err := report.NewManager("./reports")
if err != nil {
    // Handle error
}

// Generate a report
report, err := manager.GenerateReport(ctx, "quality_report", parameters)
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

By default, reports are stored in the configured output directory. The report manager handles report storage and retrieval, including:

- Generating unique report IDs
- Organizing reports by type and date
- Managing report retention according to configured policies

## Integration with External Systems

Reports can be integrated with external systems through:

- Webhooks for report notifications
- API endpoints for report retrieval
- Export capabilities for integration with data catalogs and documentation systems
