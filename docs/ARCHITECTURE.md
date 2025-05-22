# Nessi Architecture

This document provides a comprehensive overview of the Nessi architecture, explaining the design principles, core components, and how they interact.

## Table of Contents

- [Overview](#overview)
- [Design Principles](#design-principles)
- [System Architecture](#system-architecture)
- [Core Components](#core-components)
- [Data Flow](#data-flow)
- [Extension Points](#extension-points)
- [Technology Stack](#technology-stack)
- [Security Architecture](#security-architecture)
- [Performance Considerations](#performance-considerations)
- [Deployment Architecture](#deployment-architecture)
- [Future Architecture](#future-architecture)

## Overview

Nessi is a CLI-based data quality and observability tool designed for data engineers and analysts working with data lakes. The architecture is designed to be modular, extensible, and performant, with a focus on providing comprehensive data quality monitoring and reporting capabilities.

## Design Principles

Nessi's architecture is guided by the following design principles:

1. **Modularity**: Components are designed with clear boundaries and interfaces, allowing for independent development and testing.

2. **Extensibility**: The system is designed to be easily extended with new features, integrations, and plugins.

3. **Performance**: Performance is a key consideration, with optimizations for handling large datasets efficiently.

4. **Reliability**: Error handling, recovery mechanisms, and validation are built into all components.

5. **Usability**: The CLI interface is designed to be intuitive and consistent, with comprehensive documentation.

6. **Security**: Security is integrated into the design, with proper authentication, authorization, and data protection.

## System Architecture

Nessi follows a layered architecture with the following main layers:

1. **User Interface Layer**: CLI commands and interfaces
2. **Application Layer**: Core business logic and workflows
3. **Domain Layer**: Domain models and business rules
4. **Infrastructure Layer**: External system integrations and data access

![Nessi Architecture Diagram](https://nessi.dev/images/architecture.png)

## Core Components

### CLI Framework

The CLI framework provides the command-line interface for interacting with Nessi. It handles command parsing, validation, and execution, as well as output formatting.

**Key files**:
- `cmd/nessi/main.go`: Entry point for the CLI application
- `internal/cli/root.go`: Root command definition
- `internal/cli/commands/`: Individual command implementations

### Delta Lake Handler

The Delta Lake handler provides functionality for working with Delta Lake tables, including reading, writing, and metadata extraction.

**Key files**:
- `pkg/datalake/delta.go`: Delta Lake format handler
- `pkg/datalake/delta_metadata.go`: Delta Lake metadata extraction
- `pkg/datalake/delta_time_travel.go`: Time travel functionality

### Quality Engine

The quality engine provides data quality validation functionality, including rule definition, validation, and quality scoring.

**Key files**:
- `pkg/quality/rules.go`: Quality rule definitions
- `pkg/quality/validator.go`: Quality validation logic
- `pkg/quality/metrics.go`: Quality metrics calculation

### Schema Manager

The schema manager handles schema operations, including validation, evolution tracking, and comparison.

**Key files**:
- `pkg/schema/validator.go`: Schema validation
- `pkg/schema/evolution.go`: Schema evolution tracking
- `pkg/schema/comparison.go`: Schema comparison

### Metrics Collector

The metrics collector handles metrics collection and analysis, including quality metrics, performance metrics, and freshness metrics.

**Key files**:
- `pkg/metrics/collector.go`: Metrics collection
- `pkg/metrics/analyzer.go`: Metrics analysis
- `pkg/metrics/storage.go`: Metrics storage

### Report Generator

The report generator creates reports in various formats, including HTML, PDF, JSON, and CSV.

**Key files**:
- `pkg/report/generator.go`: Report generation
- `pkg/report/templates/`: Report templates
- `pkg/report/formatters/`: Output formatters

### Integration Framework

The integration framework provides a unified interface for integrating with external systems, such as Databricks and AWS S3.

**Key files**:
- `pkg/integration/framework.go`: Integration framework
- `pkg/integration/databricks/`: Databricks integration
- `pkg/integration/s3/`: AWS S3 integration

### License Manager

The license manager handles license validation, feature gating, and trial management.

**Key files**:
- `pkg/license/manager.go`: License management
- `pkg/license/validator.go`: License validation
- `pkg/license/features.go`: Feature gating

## Data Flow

### Quality Check Workflow

1. User initiates a quality check via the CLI
2. CLI framework parses and validates the command
3. Quality engine loads quality rules
4. Delta Lake handler reads the table data
5. Quality engine validates the data against the rules
6. Quality metrics are calculated and stored
7. Results are formatted and displayed to the user

### Report Generation Workflow

1. User initiates report generation via the CLI
2. CLI framework parses and validates the command
3. Report generator loads the report template
4. Quality metrics are retrieved from storage
5. Schema information is retrieved from the table
6. Report generator combines the data with the template
7. Report is generated in the requested format
8. Report is saved to the specified location

## Extension Points

Nessi provides several extension points for adding new functionality:

### Format Handlers

New data formats can be supported by implementing the `FormatHandler` interface:

```go
type FormatHandler interface {
    IsFormat(path string) bool
    Read(path string) (data []byte, err error)
    Write(path string, data []byte) error
    GetMetadata(path string) (Metadata, error)
}
```

### Quality Rules

New quality rules can be added by implementing the `Rule` interface:

```go
type Rule interface {
    Name() string
    Description() string
    Validate(data []byte, parameters map[string]interface{}) (Result, error)
}
```

### Report Templates

New report templates can be added by creating HTML templates in the `pkg/report/templates/` directory.

### Integrations

New integrations can be added by implementing the `Integration` interface:

```go
type Integration interface {
    Name() string
    Description() string
    Connect(parameters map[string]interface{}) error
    Disconnect() error
    ListResources() ([]Resource, error)
    GetResource(id string) (Resource, error)
}
```

## Technology Stack

Nessi is built using the following technologies:

### Core Technologies

- **Go**: Primary programming language
- **Python**: Used for some data processing and analysis
- **Arrow**: Data format for efficient data processing
- **Delta Lake**: Table format for data storage

### Libraries and Frameworks

- **Cobra**: CLI framework
- **Viper**: Configuration management
- **Logrus**: Logging
- **Testify**: Testing framework
- **GoMock**: Mocking framework
- **Apache Arrow**: In-memory data format
- **Delta Lake Go**: Delta Lake implementation for Go

### External Systems

- **Databricks**: Cloud-based data engineering platform
- **AWS S3**: Object storage service
- **GitHub**: Source code hosting and CI/CD

## Security Architecture

### Authentication

Nessi supports various authentication mechanisms for external systems:

- **Databricks**: Personal access tokens
- **AWS S3**: Access keys, IAM roles, and environment variables

### Authorization

Authorization is handled at the integration level, with proper permission checking for all operations.

### Data Protection

Sensitive data, such as authentication credentials, is handled securely:

- Credentials are never stored in plain text
- Configuration files with sensitive information are protected with appropriate permissions
- Environment variables are used for sensitive configuration when possible

### License Protection

The license management system includes several security measures:

- License keys are validated cryptographically
- Machine IDs are used to prevent unauthorized usage
- Anti-tampering measures prevent binary modification

## Performance Considerations

Nessi is designed to handle large datasets efficiently:

### Memory Management

- Streaming operations are used where possible to minimize memory usage
- Sampling is supported for large datasets
- Memory usage is monitored and optimized

### Concurrency

- Parallel processing is used for CPU-intensive operations
- Goroutines are used for concurrent execution
- Resource usage is balanced to prevent overloading

### Caching

- Results are cached where appropriate to avoid redundant calculations
- Cache invalidation is handled properly to ensure data freshness

## Deployment Architecture

Nessi can be deployed in various ways:

### Standalone Deployment

The simplest deployment is as a standalone CLI tool, installed on a user's machine.

### Containerized Deployment

Nessi can be deployed in a container, such as Docker, for consistent execution across environments.

### CI/CD Integration

Nessi can be integrated into CI/CD pipelines for automated quality checks and reporting.

### Scheduled Execution

Nessi can be scheduled to run periodically using tools like cron or Airflow.

## Future Architecture

The Nessi architecture is designed to evolve over time. Future architectural enhancements may include:

### Distributed Processing

Support for distributed processing of large datasets using frameworks like Apache Spark.

### Real-time Monitoring

Real-time monitoring of data quality metrics with alerting capabilities.

### Machine Learning Integration

Integration with machine learning frameworks for anomaly detection and predictive analytics.

### Web Interface

A web-based interface for visualizing data quality metrics and reports.

### API Layer

A comprehensive API layer for programmatic access to Nessi functionality.

---

For more information on the Nessi architecture, please contact the development team or refer to the source code documentation.
