# Frequently Asked Questions (FAQ)

## Navigation

- [Documentation Home](README.md)
- [User Guide](USER_GUIDE.md)
- [Installation & Quickstart](../QUICKSTART.md)
- [Configuration](CONFIGURATION.md)
- [CLI Reference](cli/README.md)
- [Reporting](REPORTING.md)
- [Quality Rules](QUALITY_RULES.md)
- [Integrations](integration_guide.md)
- [Developer Guide](developer_experience.md)

---

## General Questions

### What is Nessi?

Nessi is an open-source data quality and Delta Lake management tool that helps organizations maintain high-quality data and efficiently manage their Delta Lake tables. It provides features for data quality checks, monitoring, reporting, and Delta Lake management through a CLI-only approach.

### Is Nessi free to use?

Yes, Nessi is open-source software licensed under the Apache License 2.0. You can use, modify, and distribute it freely.

### What programming languages does Nessi support?

Nessi is primarily written in Go, but it provides integration with Python through its Python extensions.

## Installation & Setup

### What are the system requirements for Nessi?

Nessi requires:
- Go 1.18 or higher
- Python 3.8 or higher (for Python integrations)
- Access to Delta Lake tables

### How do I install Nessi?

You can install Nessi using Go:

```bash
go install github.com/nessi-dev/nessi/cmd/nessi@latest
```

Or download prebuilt binaries from the [GitHub Releases page](https://github.com/nessi-dev/nessi/releases).

For detailed installation instructions, see the [Installation & Quickstart](../QUICKSTART.md) guide.

### Can I run Nessi in Docker?

Yes, you can run Nessi using Docker:

```bash
docker pull nessi/nessi:latest
docker run -v $(pwd):/data nessi/nessi:latest scan /data/table
```

## Features & Usage

### What report formats does Nessi support?

Nessi supports the following report formats:
- HTML
- PDF
- JSON
- CSV

For more information, see the [Reporting](REPORTING.md) documentation.

### How do I create custom quality rules?

You can create custom quality rules using YAML configuration files. For details, see the [Quality Rules](QUALITY_RULES.md) documentation and [Custom Rule Extensions](custom_rule_extensions.md).

### Does Nessi support integration with data catalogs?

Yes, Nessi integrates with data catalogs like AWS Glue, Azure Purview, Google Cloud Data Catalog, and Databricks Unity Catalog. See the [Data Catalog Integration](data_catalog_integration.md) documentation for details.

### How does Nessi handle security?

Nessi provides security features including file-based configuration for CLI security, TLS configuration, and authentication. For more information, see the Security section in the [User Guide](USER_GUIDE.md#security-features).

## Troubleshooting

### How do I report a bug?

You can report bugs by opening an issue on our [GitHub repository](https://github.com/nessi-dev/nessi/issues).

### Where can I get help with Nessi?

You can get help through:
- [GitHub Discussions](https://github.com/nessi-dev/nessi/discussions)
- [Slack Community](https://nessi-community.slack.com)
- By opening an issue on our [GitHub repository](https://github.com/nessi-dev/nessi/issues)

### How do I contribute to Nessi?

See our [Contributing Guide](../CONTRIBUTING.md) for information on how to contribute to Nessi.

## Advanced Usage

### Can Nessi be integrated with workflow orchestration tools?

Yes, Nessi can be integrated with workflow orchestration tools like Airflow. See the [Workflow Orchestration](workflow_orchestration.md) documentation for details.

### Does Nessi support multi-cloud environments?

Yes, Nessi supports AWS, Azure, and GCP through its cloud integration features. See the [Cloud Integration](cloud_integration.md) documentation for details.

### Can I extend Nessi with plugins?

Yes, Nessi provides a plugin system that allows you to extend its functionality. See the [Plugin System](plugin_system.md) documentation for details.
