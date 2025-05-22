# Frequently Asked Questions (FAQ)

This document answers common questions about Nessi, its features, and usage.

## General Questions

### What is Nessi?

Nessi is an open-source CLI-only data quality and Delta Lake management tool that helps organizations maintain high-quality data and efficiently manage their Delta Lake tables through comprehensive HTML and PDF reporting.

### Is Nessi free to use?

Nessi offers a Community Edition that is completely free and open-source under the Apache License 2.0. It also offers a Pro Edition with additional features that requires a paid license.

### What programming languages is Nessi built with?

Nessi is primarily built with Go, with some Python components for specific integrations. The Go-first architecture ensures high performance and low resource usage.

## Installation & Setup

### How do I install Nessi?

You can install Nessi in several ways:
1. Using Go: `go install github.com/nessi-dev/nessi/cmd/nessi@latest`
2. Downloading a prebuilt binary from the [GitHub Releases page](https://github.com/nessi-dev/nessi/releases)
3. Using Docker: `docker pull nessi/nessi:latest`

For detailed instructions, see our [Getting Started Guide](GETTING_STARTED.md).

### What are the system requirements?

- Go 1.18 or higher (if installing from source)
- Python 3.8 or higher (for Python integrations)
- Access to Delta Lake tables
- Minimal disk space (the binary is lightweight)

### How do I configure Nessi?

Nessi uses a YAML configuration file located at `~/.nessi/config.yaml`. You can also specify a configuration file with the `--config` flag.

## Features & Usage

### What Delta Lake operations does Nessi support?

Nessi supports a wide range of Delta Lake operations, including:
- Schema evolution tracking
- Transaction log parsing and analysis
- Version control with commit history
- Time travel capabilities
- Partition management
- Table optimization

### How does Nessi handle data quality?

Nessi provides several data quality features:
- Basic validation checks
- Data profiling and statistics
- Anomaly detection
- Pattern recognition
- Custom validation rules

### Can I use Nessi with cloud storage?

Yes, with the Pro Edition, Nessi supports:
- AWS S3
- Azure Blob Storage
- Google Cloud Storage

The Community Edition supports local file systems only.

### How do I generate reports?

Use the `report` command to generate HTML or PDF reports:
```bash
nessi report /path/to/delta/table --format html --output ./reports
```

## License & Pricing

### What's included in the Community Edition?

The Community Edition includes:
- Core Delta Lake management functionality
- Basic data quality checks and profiling
- HTML/PDF reporting for local tables
- CLI monitoring for basic metrics
- Local file system support

### What additional features are in the Pro Edition?

The Pro Edition adds:
- Cloud storage integration (AWS S3, Azure, GCP)
- Databricks and dbt integration
- Data catalog integration
- Workflow orchestration with Airflow, Prefect, and Dagster
- Priority support

### How do I start a free trial?

Run the following command to start a 1-month free trial of the Pro Edition:
```bash
nessi license start-trial
```

### How many trials can I use?

You can use a maximum of 2 trials per machine.

## Troubleshooting

### Nessi can't find my Delta Lake table

Ensure the path is correct and that the directory contains a valid Delta Lake table with a `_delta_log` directory.

### I'm getting permission errors with cloud storage

Check that your environment variables for cloud authentication are properly set:
- AWS: `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
- Azure: `AZURE_STORAGE_ACCOUNT` and `AZURE_STORAGE_KEY`
- GCP: `GOOGLE_APPLICATION_CREDENTIALS`

### How do I report a bug?

Please open an issue on our [GitHub repository](https://github.com/nessi-dev/nessi/issues) with detailed information about the bug, steps to reproduce, and your environment.

## Community & Support

### How can I contribute to Nessi?

See our [Contributing Guide](../CONTRIBUTING.md) for details on how to contribute code, documentation, or report issues.

### Where can I get help?

- Check this FAQ and our documentation
- Ask in our [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions)
- Open an issue for bugs or feature requests
- For Pro Edition users, contact our support team

### Is there a community forum?

Yes, we use [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions) as our community forum.

## Advanced Usage

### Can I use Nessi in scripts or automation?

Yes, Nessi is designed to be scriptable and can be easily integrated into CI/CD pipelines, cron jobs, or other automation workflows.

### Does Nessi support custom validation rules?

Yes, you can define custom validation rules in YAML files and use them with the `validate` command.

### Can I extend Nessi's functionality?

Nessi has a plugin architecture that allows for extending its functionality. See our developer documentation for details.
