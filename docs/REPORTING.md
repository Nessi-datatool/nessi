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

Nessi provides comprehensive reporting capabilities through its flexible CLI-only report generation system. This document outlines the available report formats and how to use them effectively without requiring any web interface.

## Supported Report Formats

Nessi supports the following report formats through its CLI-only interface:

1. **HTML Reports**: Interactive HTML files with rich visualizations that can be viewed in any web browser
2. **PDF Reports**: Professional-quality PDF documents suitable for sharing with stakeholders and archiving
3. **JSON Reports**: Machine-readable reports for integration with other systems and data pipelines
4. **CSV Reports**: Tabular data exports for analysis in spreadsheet applications and business intelligence tools

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

## CLI-Based Report Generation

All reports in Nessi are generated through the command-line interface, making it easy to integrate with scripts, automation tools, and CI/CD pipelines:

```bash
# Basic report generation
nessi report quality /path/to/delta/table --format html --output quality-report.html

# Generate multiple report types in batch
nessi batch-report --tables-file tables.txt --report-types quality,schema,freshness --format pdf --output-dir reports/

# Schedule regular reports using cron
0 8 * * 1 /usr/local/bin/nessi report quality /path/to/delta/table --format pdf --output /reports/weekly/quality-$(date +"%Y%m%d").pdf
```

The CLI-only approach ensures that Nessi can be used in environments without graphical interfaces, such as servers, containers, and CI/CD pipelines.

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
