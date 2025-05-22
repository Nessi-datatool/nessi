# Report Commands

This document provides detailed information about Nessi's report commands for generating, customizing, and managing data quality reports.

## `report generate`

Generates a report from scan results, quality checks, or metrics.

### Usage

```bash
nessi report generate <source> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `source` | Path to the data source, scan result, or metrics ID |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Report type (quality, scan, metrics, schema, all) |
| `--template` | Report template to use (default, enhanced, minimal, custom) |
| `--template-path` | Path to custom template directory |
| `--output` | Output format (html, pdf, json, csv) |
| `--file` | Output file path |
| `--title` | Custom report title |
| `--include-visualizations` | Include data visualizations |
| `--include-recommendations` | Include improvement recommendations |
| `--logo` | Path to custom logo image |
| `--css` | Path to custom CSS file |
| `--metadata` | Additional metadata in key=value format |

### Examples

```bash
# Generate a quality report
nessi report generate /data/customers --type quality --output html --file quality_report.html

# Generate a comprehensive report
nessi report generate /data/customers --type all --output html --file comprehensive_report.html

# Use enhanced template
nessi report generate /data/customers --type quality --template enhanced --output html --file enhanced_report.html

# Use custom template
nessi report generate /data/customers --type quality --template custom --template-path /templates/custom --output html --file custom_report.html

# Add custom title and logo
nessi report generate /data/customers --type quality --title "Customer Data Quality Report" --logo /branding/logo.png --output html --file branded_report.html

# Generate PDF report
nessi report generate /data/customers --type quality --output pdf --file quality_report.pdf
```

### Output

```
Generating report:
  Source: /data/customers
  Type: quality
  Template: enhanced
  Output Format: html
  Output File: quality_report.html
  Visualizations: included
  Recommendations: included

Report generated successfully.
```

## `report template list`

Lists available report templates.

### Usage

```bash
nessi report template list [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--type` | Filter by template type (quality, scan, metrics, schema) |
| `--format` | Filter by output format (html, pdf, json, csv) |
| `--output` | Output format (text, json) |
| `--detailed` | Show detailed template information |

### Examples

```bash
# List all templates
nessi report template list

# List quality report templates
nessi report template list --type quality

# List HTML templates
nessi report template list --format html

# Show detailed template information
nessi report template list --detailed
```

### Output

#### Text Format (Default)

```
Available Report Templates:

NAME       TYPE      FORMAT    DESCRIPTION
default    quality   html      Standard quality report template
enhanced   quality   html      Enhanced quality report with modern styling and visualizations
minimal    quality   html      Minimal quality report with essential information only
default    quality   pdf       Standard quality report template for PDF
enhanced   quality   pdf       Enhanced quality report for PDF
default    scan      html      Standard scan report template
enhanced   scan      html      Enhanced scan report with modern styling and visualizations
default    metrics   html      Standard metrics report template
enhanced   metrics   html      Enhanced metrics report with modern styling and visualizations
default    schema    html      Standard schema report template
enhanced   schema    html      Enhanced schema report with modern styling and visualizations
```

#### Detailed Text Format

```
Template: enhanced (quality, html)
  Description: Enhanced quality report with modern styling and visualizations
  Path: /templates/quality/enhanced
  Created: 2025-04-15
  Last Modified: 2025-05-01
  Features:
    - Modern blue color scheme
    - Interactive data visualizations
    - Tabbed navigation
    - Responsive design
    - Recommendation engine
  Preview: /templates/quality/enhanced/preview.png
```

## `report template create`

Creates a new report template.

### Usage

```bash
nessi report template create [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--name` | Template name |
| `--type` | Template type (quality, scan, metrics, schema) |
| `--format` | Output format (html, pdf, json, csv) |
| `--base` | Base template to extend (default, enhanced, minimal) |
| `--description` | Template description |
| `--path` | Output path for the template files |

### Examples

```bash
# Create a new template based on the enhanced template
nessi report template create --name custom --type quality --format html --base enhanced --description "Custom quality report template" --path /templates/custom

# Create a new template from scratch
nessi report template create --name minimal-dark --type quality --format html --description "Minimal dark theme quality report" --path /templates/minimal-dark
```

### Output

```
Creating template:
  Name: custom
  Type: quality
  Format: html
  Base: enhanced
  Description: Custom quality report template
  Path: /templates/custom

Template created successfully. The following files were created:
  /templates/custom/template.html
  /templates/custom/style.css
  /templates/custom/script.js
  /templates/custom/config.json

You can now customize these files to create your custom template.
```

## `report template export`

Exports a report template.

### Usage

```bash
nessi report template export <template-name> <output-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `template-name` | Name of the template to export |
| `output-path` | Path to export the template to |

### Options

| Option | Description |
|--------|-------------|
| `--type` | Template type (quality, scan, metrics, schema) |
| `--format` | Output format (html, pdf, json, csv) |
| `--include-assets` | Include all template assets |
| `--include-examples` | Include example reports |

### Examples

```bash
# Export a template
nessi report template export enhanced /exports/templates/enhanced

# Export with all assets
nessi report template export enhanced /exports/templates/enhanced --include-assets

# Export with example reports
nessi report template export enhanced /exports/templates/enhanced --include-assets --include-examples
```

### Output

```
Exporting template:
  Name: enhanced
  Type: quality
  Format: html
  Output Path: /exports/templates/enhanced
  Assets: included
  Examples: included

Template exported successfully.
```

## `report template import`

Imports a report template.

### Usage

```bash
nessi report template import <template-path> [options]
```

### Arguments

| Argument | Description |
|----------|-------------|
| `template-path` | Path to the template to import |

### Options

| Option | Description |
|--------|-------------|
| `--name` | Name for the imported template |
| `--overwrite` | Overwrite existing template with the same name |
| `--validate` | Validate template before importing |

### Examples

```bash
# Import a template
nessi report template import /exports/templates/custom --name custom-imported

# Import and overwrite existing template
nessi report template import /exports/templates/enhanced --name enhanced --overwrite

# Validate template before importing
nessi report template import /exports/templates/custom --name custom-imported --validate
```

### Output

```
Importing template:
  Path: /exports/templates/custom
  Name: custom-imported
  Validation: passed

Template imported successfully.
```

## `report schedule`

Manages scheduled reports.

### Usage

```bash
nessi report schedule [command] [options]
```

### Commands

| Command | Description |
|---------|-------------|
| `list` | List all scheduled reports |
| `create` | Create a new scheduled report |
| `update` | Update an existing scheduled report |
| `delete` | Delete a scheduled report |
| `enable` | Enable a scheduled report |
| `disable` | Disable a scheduled report |
| `history` | Show report generation history |

### Options for `list`

| Option | Description |
|--------|-------------|
| `--source` | Filter by source path |
| `--type` | Filter by report type (quality, scan, metrics, schema) |
| `--status` | Filter by status (enabled, disabled) |
| `--output` | Output format (text, json, csv) |

### Options for `create` and `update`

| Option | Description |
|--------|-------------|
| `--id` | Report schedule ID (for update only) |
| `--name` | Schedule name |
| `--source` | Path to the data source |
| `--type` | Report type (quality, scan, metrics, schema, all) |
| `--template` | Report template to use |
| `--output` | Output format (html, pdf, json, csv) |
| `--file` | Output file path pattern |
| `--schedule` | Schedule expression (cron format) |
| `--notify` | Notification settings (email, webhook, none) |
| `--description` | Schedule description |

### Options for `delete`, `enable`, and `disable`

| Option | Description |
|--------|-------------|
| `--id` | Schedule ID |
| `--force` | Force operation without confirmation (for delete) |

### Options for `history`

| Option | Description |
|--------|-------------|
| `--id` | Schedule ID |
| `--from` | Start date for history (YYYY-MM-DD) |
| `--to` | End date for history (YYYY-MM-DD) |
| `--limit` | Maximum number of history entries to show |
| `--output` | Output format (text, json, csv) |

### Examples

```bash
# List all scheduled reports
nessi report schedule list

# Create a scheduled report
nessi report schedule create --name "Daily Quality Report" --source /data/customers --type quality --template enhanced --output html --file "/reports/customers_quality_%Y%m%d.html" --schedule "0 0 * * *" --notify email --description "Daily quality report for customer data"

# Update a scheduled report
nessi report schedule update --id 12345 --schedule "0 0 * * 1" --description "Weekly quality report for customer data"

# Delete a scheduled report
nessi report schedule delete --id 12345

# Enable a scheduled report
nessi report schedule enable --id 12345

# Disable a scheduled report
nessi report schedule disable --id 12345

# Show report generation history
nessi report schedule history --id 12345 --from 2025-04-01 --to 2025-05-22
```

### Output for `list`

```
Scheduled Reports:

ID     NAME                  SOURCE            TYPE     SCHEDULE    STATUS   LAST RUN            NEXT RUN
12345  Daily Quality Report  /data/customers   quality  0 0 * * *   enabled  2025-05-22 00:00:00 2025-05-23 00:00:00
67890  Weekly Scan Report    /data/customers   scan     0 0 * * 1   enabled  2025-05-20 00:00:00 2025-05-27 00:00:00
24680  Monthly Metrics       /data/customers   metrics  0 0 1 * *   disabled 2025-05-01 00:00:00 N/A
```

### Output for `create` and `update`

```
Scheduled Report Created:
  ID: 12345
  Name: Daily Quality Report
  Source: /data/customers
  Type: quality
  Template: enhanced
  Output: html to "/reports/customers_quality_%Y%m%d.html"
  Schedule: 0 0 * * * (Daily at midnight)
  Notification: email
  Status: enabled
  Next Run: 2025-05-23 00:00:00
```

### Output for `history`

```
Report Generation History: Daily Quality Report (ID: 12345)
Period: 2025-04-01 to 2025-05-22

TIMESTAMP            STATUS    OUTPUT FILE                            SIZE      DURATION
2025-05-22 00:00:15  SUCCESS   /reports/customers_quality_20250522.html  1.2 MB    5.3s
2025-05-21 00:00:12  SUCCESS   /reports/customers_quality_20250521.html  1.2 MB    5.2s
2025-05-20 00:00:18  SUCCESS   /reports/customers_quality_20250520.html  1.2 MB    5.4s
2025-05-19 00:00:14  SUCCESS   /reports/customers_quality_20250519.html  1.2 MB    5.3s
2025-05-18 00:00:11  SUCCESS   /reports/customers_quality_20250518.html  1.2 MB    5.1s
2025-05-17 00:00:13  SUCCESS   /reports/customers_quality_20250517.html  1.2 MB    5.2s
2025-05-16 00:00:15  SUCCESS   /reports/customers_quality_20250516.html  1.2 MB    5.3s
2025-05-15 00:00:12  SUCCESS   /reports/customers_quality_20250515.html  1.2 MB    5.2s
```

## `report history`

Shows the history of generated reports.

### Usage

```bash
nessi report history [options]
```

### Options

| Option | Description |
|--------|-------------|
| `--source` | Filter by source path |
| `--type` | Filter by report type (quality, scan, metrics, schema) |
| `--from` | Start date for history (YYYY-MM-DD) |
| `--to` | End date for history (YYYY-MM-DD) |
| `--limit` | Maximum number of history entries to show |
| `--output` | Output format (text, json, csv) |
| `--detailed` | Show detailed report information |

### Examples

```bash
# Show report history
nessi report history

# Filter by source
nessi report history --source /data/customers

# Filter by report type
nessi report history --type quality

# Filter by date range
nessi report history --from 2025-04-01 --to 2025-05-22

# Show detailed report information
nessi report history --detailed
```

### Output

#### Text Format (Default)

```
Report Generation History:
Period: 2025-04-01 to 2025-05-22

TIMESTAMP            SOURCE            TYPE     TEMPLATE  FORMAT  OUTPUT FILE                        SIZE
2025-05-22 22:12:09  /data/customers   quality  enhanced  html    /reports/quality_20250522.html     1.2 MB
2025-05-22 00:00:15  /data/customers   quality  enhanced  html    /reports/customers_quality_20250522.html  1.2 MB
2025-05-21 00:00:12  /data/customers   quality  enhanced  html    /reports/customers_quality_20250521.html  1.2 MB
2025-05-20 23:00:18  /data/customers   scan     enhanced  html    /reports/customers_scan_20250520.html     2.5 MB
2025-05-20 00:00:18  /data/customers   quality  enhanced  html    /reports/customers_quality_20250520.html  1.2 MB
2025-05-19 00:00:14  /data/customers   quality  enhanced  html    /reports/customers_quality_20250519.html  1.2 MB
2025-05-18 00:00:11  /data/customers   quality  enhanced  html    /reports/customers_quality_20250518.html  1.2 MB
2025-05-15 23:00:12  /data/customers   scan     enhanced  html    /reports/customers_scan_20250515.html     2.4 MB
```

#### Detailed Text Format

```
Report: /reports/quality_20250522.html
  Source: /data/customers
  Type: quality
  Template: enhanced
  Format: html
  Timestamp: 2025-05-22 22:12:09
  Size: 1.2 MB
  Generation Duration: 5.3s
  User: john.doe
  Command: nessi report generate /data/customers --type quality --template enhanced --output html --file /reports/quality_20250522.html
  Quality Score: 0.9879 (98.79%)
  Rules Checked: 4
  Rules Passed: 3
  Rules Failed: 1
```

## Report Generation System

Nessi's report commands utilize a robust report generation system with the following features:

### Output Formats

- **HTML**: Interactive web-based reports with visualizations and navigation
- **PDF**: Printable document format generated using wkhtmltopdf or reportlab
- **JSON**: Machine-readable format with all report data
- **CSV**: Tabular format for importing into spreadsheets

### Templates

Nessi includes several built-in templates:

1. **Default Template**: Standard report with clean layout and basic styling
2. **Enhanced Template**: Modern blue color scheme with improved typography, interactive elements, and visualizations
3. **Minimal Template**: Simplified report with just the essential information

### Report Types

The report generation system supports various report types:

1. **Quality Reports**: Focus on data quality metrics and rule validation results
2. **Scan Reports**: Comprehensive data profiling and analysis
3. **Metrics Reports**: Time-series analysis of data quality metrics
4. **Schema Reports**: Schema details, validation, and evolution

### Customization

Reports can be customized in several ways:

1. **Custom Templates**: Create and modify templates with HTML, CSS, and JavaScript
2. **Custom Styling**: Apply custom CSS to existing templates
3. **Custom Branding**: Add logos and custom titles
4. **Custom Visualizations**: Add or modify data visualizations

### Templating Engine

The report generation system uses Jinja2 for HTML templates with the following features:

1. **Template Inheritance**: Extend base templates with custom content
2. **Conditional Rendering**: Show or hide sections based on data
3. **Loops and Iterations**: Iterate over data collections
4. **Filters and Transformations**: Format and transform data for display

### Report Components

Enhanced reports include the following components:

1. **Summary Cards**: Overview of key metrics and findings
2. **Progress Indicators**: Visual representation of quality scores
3. **Data Tables**: Sortable and filterable tables of detailed data
4. **Charts and Graphs**: Visual representation of data distributions and trends
5. **Recommendations**: Actionable insights for improving data quality
6. **Navigation**: Tabbed interface for accessing different report sections

For more information on the report generation system, see the [Report Generation documentation](../reports/REPORT_GENERATION.md).

## Error Handling

Report commands use the following error codes:

- `N1000`: Report generation failed
- `N1001`: Template not found
- `N1002`: Invalid template format
- `N1003`: Report scheduling failed
- `N1004`: Report history retrieval failed
- `N1005`: Template creation failed
- `N1006`: Template import/export failed
- `N1007`: Notification failed

For more information on error handling, see the [Error Handling documentation](../ERROR_HANDLING.md).
