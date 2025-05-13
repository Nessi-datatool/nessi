# Nessi.dev Workflow Orchestration Integrations

This directory contains integrations for popular workflow orchestration tools, allowing you to incorporate Nessi.dev data quality checks, profiling, validation, and lineage tracking into your data pipelines.

## Available Integrations

- **[Apache Airflow](./airflow/)**: Custom operators and sensors for Airflow 2.0+
- **[Prefect](./prefect/)**: Custom tasks and flows for Prefect 2.0+
- **[Dagster](./dagster/)**: Custom ops, resources, and jobs for Dagster 1.0+

## Common Features

All workflow orchestration integrations provide the following capabilities:

- **Data Quality Checks**: Run quality checks on your data and wait for results
- **Data Profiling**: Generate profiles of your data to understand its characteristics
- **Data Validation**: Validate your data against custom rules
- **Lineage Tracking**: Capture and visualize data lineage information
- **Failure Handling**: Configure whether to fail workflows based on quality scores or rule failures
- **Authentication**: Secure communication with Nessi.dev API using API keys

## Documentation

For detailed documentation on how to use these integrations, please refer to the [Workflow Orchestration Integration Guide](../docs/workflow_orchestration.md).

## Examples

Each integration directory contains example workflows that demonstrate how to use Nessi.dev with the respective orchestration tool:

- **Airflow**: Example DAGs in `airflow/examples/`
- **Prefect**: Example flows in `prefect/examples/`
- **Dagster**: Example jobs in `dagster/examples/`

## Installation

Each integration can be installed separately:

```bash
# For Apache Airflow
pip install -e ./airflow

# For Prefect
pip install -e ./prefect

# For Dagster
pip install -e ./dagster
```

## Contributing

If you'd like to contribute to these integrations or add support for additional workflow orchestration tools, please follow the [contribution guidelines](../CONTRIBUTING.md).
