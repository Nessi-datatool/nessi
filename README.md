# Nessi

<p align="center">
  <img src="docs/images/nessi-logo.png" alt="Nessi Logo" width="200"/>
</p>

Nessi is an open-source CLI-only data quality and Delta Lake management tool that helps organizations maintain high-quality data and efficiently manage their Delta Lake tables through comprehensive HTML and PDF reporting. Nessi offers both community and premium features with a flexible licensing system.

[![Go Report Card](https://goreportcard.com/badge/github.com/nessi-dev/nessi)](https://goreportcard.com/report/github.com/nessi-dev/nessi)
[![Build Status](https://github.com/nessi-dev/nessi/workflows/CI/badge.svg)](https://github.com/nessi-dev/nessi/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/nessi-dev/nessi.svg)](https://pkg.go.dev/github.com/nessi-dev/nessi)

## Features

- **Delta Lake Support**
  - Delta Lake table management
  - Transaction log parsing and analysis
  - Schema evolution tracking
  - Partition management
  - Version control and time travel
  - Multi-format support (Delta, Parquet, CSV)
  - Schema inference and validation
  - Cloud integration with AWS, Azure, and GCP
  - **Data Catalog Integration**: Integration with data catalogs like AWS Glue, Azure Purview, Google Cloud Data Catalog, and Databricks Unity Catalog

- **Data Quality Intelligence**
  - Data quality checks
  - Custom quality rules
  - Data profiling and statistics
  - Anomaly detection

- **Comprehensive Error Handling System**
  - Standardized error codes (N1XX-N9XX) for different error categories
  - Interactive error resolution for common errors
  - Contextual error suggestions with actionable guidance
  - Error telemetry for tracking error statistics
  - Automatic retry with exponential backoff for transient failures
  - Detailed troubleshooting documentation
  - Recovery mechanisms for corrupted tables
  - Authentication and connection error handling
  - Test tools for verifying error handling
  
- **Advanced Usability Features**
  - Input validation and auto-correction for paths and configuration values
  - Interactive setup wizard for first-time users
  - Dry run mode to preview operations without making changes
  - Smart defaults based on context
  - Comprehensive configuration system with clear precedence
  - Detailed debugging and logging capabilities
  - Consistent error messages with actionable suggestions
  
- **Monitoring & Alerting**
  - File-based metrics collection
  - Email alerting for failed checks
  - Configurable alert thresholds
  - Alert history in log files
  
- **Reporting & Visualization**
  - CLI-generated HTML and PDF reports
  - Comprehensive data quality dashboards in HTML format
  - Exportable PDF reports for sharing with stakeholders
  - Schema tree visualization with detailed type information
  - Data quality metrics with heatmap visualization
  - Trend analysis via report comparison
  - Quality scoring and metrics
  - JSON and CSV export for integration with other tools
  - **dbt Integration**: Seamless integration with dbt for validating and profiling models directly in your dbt workflow

- **Progress Visualization**
  - Real-time feedback during long-running operations
  - Spinner-style indicators for operations without measurable progress
  - Progress bars for operations with measurable percentage completion
  - Success, warning, and error messages with clear status indicators
  - Demo command to showcase progress visualization features

## Installation

### Prerequisites

- Go 1.18 or higher
- Python 3.8 or higher (for Python integrations)
- Access to Delta Lake tables

### Binary Installation

#### Using Go

```bash
go install github.com/nessi-dev/nessi/cmd/nessi@latest
```

#### Using Prebuilt Binaries

Download the latest release from the [GitHub Releases page](https://github.com/nessi-dev/nessi/releases).

```bash
# Linux/macOS
chmod +x nessi
./nessi --version

# Add to your PATH for easier access
mv nessi /usr/local/bin/
```

### Docker Installation

```bash
# Pull the latest image
docker pull nessi/nessi:latest

# Run Nessi
docker run -v $(pwd):/data nessi/nessi:latest scan /data/table
```

## Quick Start

### Basic Usage

```bash
# Scan a Delta Lake table
nessi scan s3://my-bucket/my-table

# Generate a data quality report
nessi report s3://my-bucket/my-table --format html --output report.html

# Validate a table against quality rules
nessi validate s3://my-bucket/my-table --rules rules.yaml

# Time travel to a specific version
nessi timetravel s3://my-bucket/my-table --version 5

# View schema evolution
nessi schema s3://my-bucket/my-table --history

# Display schema as ASCII tree
nessi schema-tree --table s3://my-bucket/my-table

# Generate a heatmap visualization for data quality metrics
nessi visualize heatmap --table s3://my-bucket/my-table --metric completeness

# Generate a bar chart visualization for all metrics
nessi visualize barchart --table s3://my-bucket/my-table --all-metrics

# Generate a scatter plot comparing two metrics
nessi visualize scatter --table s3://my-bucket/my-table --x-metric completeness --y-metric accuracy

# Generate a time series visualization for a metric over time
nessi visualize timeseries --table s3://my-bucket/my-table --metric completeness --days 30

# Generate a pie chart showing distribution of data quality issues
nessi visualize piechart --table s3://my-bucket/my-table --category issue-types
```

### Configuration

Create a `nessi.yaml` configuration file:

```yaml
storage:
  type: s3
  region: us-west-2
  bucket: my-delta-tables

monitoring:
  enabled: true
  metrics_retention: 30d

quality:
  default_rules: path/to/rules.yaml
  threshold: 0.95
```

### Using the API

Nessi provides a Go API for programmatic access:

```go
package main

import (
	"fmt"
	"github.com/nessi-dev/nessi/pkg/datalake"
)

func main() {
	// Open a Delta Lake table
	table, err := datalake.OpenTable("s3://my-bucket/my-table")
	if err != nil {
		panic(err)
	}
	
	// Get table metadata
	meta, err := table.Metadata()
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Table: %s, Version: %d\n", meta.Name, meta.Version)
	
	// Run quality checks
	results, err := table.ValidateQuality(nil)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Quality Score: %.2f%%\n", results.OverallScore*100)
}
```

## Documentation

For detailed documentation, visit our [Documentation Home](docs/README.md).

### Key Documentation Pages:
- [User Guide](docs/USER_GUIDE.md) - Comprehensive guide to using Nessi
- [CLI Reference](docs/cli/README.md) - Command-line interface documentation
- [Configuration Options](docs/CONFIGURATION.md) - Configure Nessi for your environment
- [Quality Rules](docs/QUALITY_RULES.md) - Define and manage quality rules
- [Reporting](docs/REPORTING.md) - Generate reports in various formats
- [Advanced Features](docs/advanced_features.md) - Delta Lake features, monitoring, and more
- [Integrations](docs/integration_guide.md) - Cloud, Databricks, and workflow integrations
- [Error Handling](docs/ERROR_HANDLING.md) - Comprehensive error handling documentation
- [Usability Features](docs/USABILITY_FEATURES_SUMMARY.md) - Summary of all usability features
- [FAQ](docs/faq.md) - Frequently asked questions

- **Monitoring and Alerting (OSS + Optional)**
  - CLI-based metrics reporting
  - Email notifications for failed checks (OSS)
  - Performance tracking
  - Resource utilization monitoring
  - Health checks
  - **Freshness SLA Monitoring**: Track data freshness with configurable SLAs and compliance tracking

- **Report Generation**
  - Customizable report templates
  - Multiple output formats (CSV, JSON)
  - Report generation via CLI
  - Report archiving and retention

- **Security**
  - API key management
  - TLS encryption for secure connections
  - File-based credential storage

## Contributing

We welcome contributions from the community! Please see our [Contributing Guide](CONTRIBUTING.md) for details on how to get started.

### Development Setup

```bash
# Clone the repository
git clone git@github.com:nessi-dev/nessi.git
cd nessi

# Install dependencies
make deps

# Build the project
make build

# Run tests
make test
```

### Testing

We aim for high test coverage. Run the test suite with:

```bash
./generate-test-coverage.sh
```

This will generate a coverage report in the `./coverage` directory.

## Community

- [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions) - Ask questions and share ideas
- [Issue Tracker](https://github.com/nessi-dev/nessi/issues) - Report bugs or request features
- [Slack Community](https://nessi-community.slack.com) - Join our Slack for real-time discussions

## Roadmap

See our [project roadmap](https://github.com/nessi-dev/nessi/projects/1) for upcoming features and enhancements.

## License

Nessi is licensed under the [Apache License 2.0](LICENSE).

## License Management System

Nessi includes a robust license management system that provides the following features:

- **Community Edition**: Free and open-source with core functionality
- **Premium Features**: Advanced features available with a valid license or during a free trial
- **Free Trial**: 1-month free trial of all premium features
- **Machine-Based Licensing**: Licenses are tied to specific machines
- **Trial Limitations**: Maximum of 2 trials per machine
- **Security**: All license and trial data is cryptographically signed to prevent tampering

### License Tiers

Nessi offers both open-source and commercial features:

### Community Edition (Open Source)

The free, open-source edition includes all core functionality:

| Category | Included Features |
|----------|-------------------|
| **Delta Lake** | Schema evolution tracking, Transaction history, Time travel, Version control |
| **Data Quality** | Basic validation rules, Profiling, Anomaly detection, Pattern matching |
| **Reporting** | HTML/PDF reports, Basic visualizations, Export capabilities |
| **Integrations** | Local file systems, SQLite support, CSV/JSON imports |
| **Workflow** | CLI commands, Basic automation |
| **Support** | Community forums, Documentation, Issue tracker |

### Pro Edition (Commercial)

Builds on the Community Edition with enterprise features. To start a free 1-month trial, run `nessi license start-trial`.

| Category | Additional Features |
|----------|---------------------|
| **Delta Lake** | Cloud storage integration, Enhanced performance for large tables |
| **Data Quality** | Advanced validation, Custom rule engines, Cross-table validation |
| **Reporting** | Scheduled reporting |
| **Integrations** | AWS S3, Azure, GCP, Databricks, dbt, Data catalogs |
| **Workflow** | Airflow integration, Prefect integration, Dagster integration |
| **Support** | Priority support, SLA guarantees, Direct assistance |

> **Enterprise Tier**: For enterprise-level support and custom features, please visit [nessi.dev](https://nessi.dev) for contact information.

### Starting a Trial

```bash
# Start a free trial to access premium features
nessi license start-trial

# Check your license status
nessi license info
```

## Acknowledgments

- The Delta Lake community for their excellent work on the Delta Lake format
- All our contributors and users who provide valuable feedback

- **Testing Framework**
  - Comprehensive unit and integration tests
  - In-memory testing for performance-critical components
  - Mock connectors for external dependencies
  - Test data generation utilities
  - Configurable test runners

## Developer Experience

- **CLI-First Approach**
  - One-line validation commands
  - Batch processing scripts
  - Comprehensive configuration management
  - Intuitive command structure

- **Developer Experience**
  - Intuitive CLI interface
  - Comprehensive documentation
  - Python API for programmatic access
  - Extensible architecture with plugin system
  - Custom extension support via plugins
  - Embeddable in Airflow, GitHub Actions, Kubernetes, and custom scripts
  - Interactive examples and integration guides
  - Observability features with metrics collection, structured logging, and distributed tracing
  - Authentication enhancements with OAuth 2.0 and token refresh capabilities
  - Environment variable-based configuration for secure deployment

- **Integration Capabilities**
  - File-based metrics collection
  - CLI-based reporting with HTML and PDF output
  - Comprehensive report generation system
  - Email notifications
  - Cloud provider integration (AWS, Azure, GCP)
  - **dbt Integration**: Opt-in plugin for executing data quality rules against dbt models and generating data profiles
  - **Workflow Orchestration Integration**: Integration with Apache Airflow, Prefect, and Dagster for incorporating data quality checks into data pipelines
  - **Error Handling and Retries**: Configurable retry mechanisms for operations and detailed error reporting
  - **Observability**: Structured logging and file-based metrics for monitoring

- **Documentation**
  - Comprehensive user guides
  - API reference
  - Example scripts
  - Best practices
  - Integration guides for dbt, Airflow, and other tools

## Requirements

- Go 1.21 or later
- Docker (optional)
- Make
- OpenSSL (for certificate generation)

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:nessi-dev/nessi.git
   cd nessi-dev
   ```

2. Install development tools:
   ```bash
   make install-tools
   ```

3. Initialize the development environment:
   ```bash
   make init
   ```

4. Build the application:
   ```bash
   make build
   ```

## Configuration

The application is configured using a YAML file located at `config/config.yaml`. The configuration includes settings for:

- Security settings (API keys, TLS for CLI tools)
- Delta Lake paths and properties
- Data quality rules and thresholds
- Monitoring and metrics export
- Email notifications
- Report generation
- Logging

See `config/config.yaml` for detailed configuration options.

## Usage

### Running the CLI

1. Run the CLI tool:
   ```bash
   nessi --help
   ```

2. Or run with Docker:
   ```bash
   docker run -v $(pwd):/data nessi/nessi:latest --help
   ```

### CLI Commands

The application provides the following CLI commands:

- **Table Management**
  - `nessi tables list` - List Delta tables
  - `nessi tables info <table_path>` - Get table details

- **Quality Management**
  - `nessi quality check <table_path>` - Run quality checks
  - `nessi quality metrics <table_path>` - Get quality metrics
  - `nessi quality profile <table_path>` - Generate data profile

- **Reporting**
  - `nessi report generate --table <table_path> --format <format>` - Generate a report (html, pdf, json, csv)
  - `nessi report list` - List available reports
  - `nessi report view <report_id>` - View a specific report
  - `nessi report export <report_id> --output <path>` - Export a report to a file

- **Catalog Integration**
  - `nessi catalog list --type <catalog_type>` - List available catalogs
  - `nessi catalog tables --type <catalog_type> --database <db>` - List tables in a catalog
  - `nessi catalog describe --type <catalog_type> --database <db> --table <table>` - Get table details from a catalog

- **Freshness Monitoring**
  - `nessi freshness check <table_path>` - Check table freshness
  - `nessi freshness report <table_path> --format <format>` - Generate freshness report

### Databricks Integration

Nessi supports integration with Databricks, allowing you to work with Databricks Unity Catalog and Delta Lake tables.

#### Configuration

To use the Databricks integration, set the following environment variables:

```bash
export DATABRICKS_HOST="https://your-workspace.cloud.databricks.com"
export DATABRICKS_TOKEN="your-personal-access-token"
export DATABRICKS_WORKSPACE_ID="your-workspace-id"  # Optional, defaults to "0"
export DATABRICKS_DEFAULT_SCHEMA="default"          # Optional, defaults to "default"
export DATABRICKS_DEFAULT_CATALOG="hive_metastore"  # Optional, defaults to "hive_metastore"
```

#### Example Commands

```bash
# List catalogs in Databricks
nessi catalog list --type databricks

# List tables in a specific catalog and schema
nessi catalog tables --type databricks --database main.default

# Get details of a specific table
nessi catalog describe --type databricks --database main.default --table customers

# Work with Delta Lake tables
nessi tables info /path/to/delta/table
```

### Development

1. Run tests:
   ```bash
   make test
   ```

2. Run linters:
   ```bash
   make lint
   ```

3. Clean build artifacts:
   ```bash
   make clean
   ```

## Security

### TLS Configuration

1. Generate self-signed certificates:
   ```bash
   make generate-certs
   ```

2. Update the configuration in `config/config.yaml`:
   ```yaml
   security:
     tls:
       enabled: true
       cert_file: "certs/server.crt"
       key_file: "certs/server.key"
   ```

### Authentication

The application uses file-based configuration for CLI security. Configure the security settings in `config/config.yaml`:

```yaml
security:
  credentials:
    enabled: true
    file: "config/credentials.yaml"
```

You can also set credentials using environment variables for CLI commands:

```bash
export NESSI_CREDENTIALS_FILE="path/to/credentials.yaml"
nessi tables list
```

## Reporting and Metrics Export

Nessi provides a comprehensive reporting system with support for HTML, PDF, JSON, and CSV formats. Configure the reporting settings in `config/config.yaml`:

```yaml
reporting:
  enabled: true
  output_path: "./data/reports"
  default_format: "html"
  templates_path: "./templates"
  retention:
    enabled: true
    days: 30
  metrics:
    export_path: "./data/metrics"
    export_format: "json"
  logging:
    structured: true
    format: "json"
    level: "info"
```

See the [Reporting Documentation](docs/REPORTING.md) for more details on the reporting capabilities.

### Error Handling System

Nessi includes a comprehensive error handling system with standardized error codes, interactive resolution, and telemetry. Configure the error handling settings in `config/config.yaml`:

```yaml
error_handling:
  interactive_resolution: true  # Enable interactive error resolution
  telemetry:
    enabled: true     # Enable error telemetry
    anonymous: true   # Collect only anonymous data
    max_errors: 100   # Maximum number of errors to store
  retry:
    enabled: true     # Enable automatic retries
    max_attempts: 3
    initial_backoff: 1s
    max_backoff: 30s
    backoff_factor: 2.0
    jitter: true      # Add jitter to backoff
  suggestions:
    enabled: true     # Enable error suggestions
    max_suggestions: 3
    show_documentation: true
  logging:
    level: info       # Log level for errors (debug, info, warn, error)
    colored: true     # Use colored output for errors
```

#### Testing Error Handling

Nessi provides a command to test the error handling system:

```bash
# List all available error codes
nessi test-error --list

# Generate a specific error to test handling
nessi test-error N101  # Path error

# Test interactive resolution
nessi test-error N201 --resolvable

# Test error telemetry
nessi test-error N301 --telemetry

# Run comprehensive error handling tests
./scripts/test_error_handling.sh
```

### Usability Features

Nessi includes several usability features to make it more error-free and easier to use:

#### Setup Wizard

First-time users can use the setup wizard to configure Nessi:

```bash
# Run the interactive setup wizard
nessi setup

# Run in non-interactive mode with default values
nessi setup --non-interactive --config-dir ~/.nessi
```

#### Dry Run Mode

Preview operations without making changes:

```bash
# Show what would happen without making changes
nessi schema create --path /path/to/table --dry-run

# Preview repair operations
nessi repair --table /path/to/table --dry-run
```

#### Input Validation

Nessi validates and normalizes input before executing commands:

- Path validation (expansion of `~`, relative to absolute conversion)
- Configuration validation (required fields, type checking)
- Output format validation (html, pdf, json, csv, text)

#### Command Line Flags

Common flags available for all commands:

```bash
# Enable interactive mode
nessi <command> --interactive

# Preview operations without making changes
nessi <command> --dry-run

# Specify configuration file
nessi <command> --config /path/to/config.yaml

# Enable verbose output
nessi <command> --verbose
```

## Test Strategy

We use a unified test script that supports multiple test modes:

### Fast/unit tests
Run locally for rapid feedback:

```sh
./scripts/run_tests.sh --fast
```
This skips all integration/long-running tests.

### Short tests
Run specific test files:

```sh
./scripts/run_tests.sh --short
```

### Standard tests
Run tests that complete in under 30 seconds (default):

```sh
./scripts/run_tests.sh
```

### Integration/slow tests
Run all tests including long-running ones:

```sh
./scripts/run_tests.sh --all
```

All integration/slow tests use the Go convention:
```go
if testing.Short() {
    t.Skip("Skipping integration test in short mode")
}
```

For more options, see the help:
```sh
./scripts/run_tests.sh --help
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support, please open an issue in the GitHub repository or contact the maintainers. 