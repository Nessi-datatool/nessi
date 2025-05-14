# Nessi.dev dbt Plugin

## Integrated Data Quality for dbt & Delta Lake

The Nessi.dev dbt Plugin offers dbt users the ability to seamlessly integrate Nessi.dev's powerful data quality and profiling capabilities directly into their dbt workflows for Delta Lake. This integration is offered as an opt-in feature, allowing users to explicitly enable and configure it within their Nessi.dev and/or dbt setup.

## Key Features

### DBT Integration

- Execute Nessi.dev data quality rules against dbt models using the `nessi dbt validate` command, with enhanced model selection leveraging dbt's standard syntax (e.g., `tag:daily,+downstream`).
- Generate data profiles for Delta tables associated with dbt models using the `nessi dbt profile` command.
- Optionally integrate dbt's native test results into Nessi's quality assessment for a unified view.
- Enable lineage-aware validation to propagate quality issues to upstream dbt models, facilitating root cause analysis.

### Core Functionality

- Automatic mapping between dbt models and their underlying Delta tables.
- Execution of data quality rules defined in Nessi.dev against these Delta tables.
- Generation of basic and enhanced profiling statistics.

### Extensibility & Reporting

- Parse configuration settings from YAML files.
- Export validation and profiling results in various formats (e.g., JSON, CSV, table) to integrate with external tools and reporting systems.
- Generate artifacts that can be integrated into dbt Docs, enriching documentation with data quality insights.

### Data Quality Management

- Generate an automated data quality score based on rule execution results for trend analysis and monitoring.
- Provide CI/CD-friendly output modes and exit codes for easy integration into automated pipelines.
- Offer alerting capabilities to notify users of failed validations via channels like Slack or email.

## Opt-In Configuration

Users have explicit control over whether to enable and utilize this dbt plugin through the `enable_dbt_plugin` configuration flag in the nessi.yaml file or by using the `--enable-dbt-plugin` flag with CLI commands. This ensures it doesn't impose any burden on users who do not require this integration.

## Getting Started

### Installation

The dbt plugin is included in the main Nessi.dev installation. No additional installation steps are required.

### Configuration

1. Enable the dbt plugin in your Nessi.dev configuration:

```yaml
# nessi.yaml
dbt:
  enable_dbt_plugin: true
  project_path: "/path/to/your/dbt/project"
  alert:
    slack:
      webhook_url: "https://hooks.slack.com/services/your/webhook/url"
      channel: "#data-quality"
    email:
      smtp_host: "smtp.example.com"
      smtp_port: 587
      username: "your-email@example.com"
      password: "your-password"
      recipients:
        - "team@example.com"
```

2. Map your dbt models to Delta tables:

```yaml
# nessi.yaml
dbt:
  model_mappings:
    - model: "customers"
      table_path: "/delta/warehouse/customers"
    - model: "orders"
      table_path: "/delta/warehouse/orders"
```

3. Define data quality rules for your models:

```yaml
# nessi.yaml
dbt:
  rule_sets:
    - name: "common_rules"
      rules:
        - name: "no_nulls_in_id"
          sql: "SELECT * FROM {table} WHERE id IS NULL"
          failure_threshold: 0
    - name: "customer_rules"
      rules:
        - name: "valid_email"
          sql: "SELECT * FROM {table} WHERE email NOT LIKE '%@%.%'"
          failure_threshold: 0
```

4. Map rule sets to models:

```yaml
# nessi.yaml
dbt:
  model_rule_mappings:
    - model: "customers"
      rule_sets: ["common_rules", "customer_rules"]
    - model: "orders"
      rule_sets: ["common_rules"]
```

### Usage

#### Validate Models

```bash
nessi dbt validate model_name
nessi dbt validate tag:daily
nessi dbt validate model_name+downstream
```

#### Profile Models

```bash
nessi dbt profile model_name
nessi dbt profile tag:daily
```

#### Output Formats

```bash
nessi dbt validate model_name --output=json
nessi dbt validate model_name --output=csv
nessi dbt validate model_name --output=table
```

## Integration with CI/CD

The dbt plugin is designed to work seamlessly with CI/CD pipelines. It provides appropriate exit codes based on validation results, making it easy to integrate with your CI/CD workflow.

```bash
# Example GitHub Actions workflow step
- name: Validate dbt models
  run: nessi dbt validate tag:daily --output=json
  continue-on-error: false
```

## Alerting

The dbt plugin can send alerts when validation rules fail. Configure the alerting settings in your Nessi.dev configuration file to receive notifications via Slack or email.
