# Advanced Features

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
- [FAQ](faq.md)

---

## Overview

Nessi provides a range of advanced features for data quality management, Delta Lake operations, and monitoring. This page serves as an index to the detailed documentation for these advanced features.

## Delta Lake Features

[Delta Lake Features](delta_lake_features.md) provides comprehensive capabilities for working with Delta Lake tables:

- Schema evolution tracking
- Version control and time travel
- Partition management
- Transaction log parsing and analysis

## Freshness Monitoring

[Freshness Monitoring](freshness_monitoring.md) allows you to track the freshness of your data:

- Define freshness thresholds for tables
- Monitor data arrival and processing times
- Generate freshness reports
- Set up alerts for stale data

## Intelligent Alerting

[Intelligent Alerting](intelligent_alerting.md) provides advanced alerting capabilities:

- Configure alerts based on complex conditions
- Set up notification channels (email, Slack, webhooks)
- Define alert severity levels
- Implement alert suppression and grouping

See [Intelligent Alerting Examples](intelligent_alerting_examples.md) for practical examples.

## Data Lineage

[Data Lineage](lineage.md) helps you track the flow of data through your systems:

- Visualize data dependencies
- Track data transformations
- Analyze impact of changes
- Understand data provenance

## Multi-Format Support

[Multi-Format Support](multi_format_support.md) extends Nessi's capabilities beyond Delta Lake:

- Support for Parquet, CSV, JSON, and other formats
- Schema inference and validation
- Format-specific optimizations
- Cross-format compatibility checks

## RBAC System

[RBAC System](rbac_system.md) provides role-based access control for Nessi:

- Define roles and permissions
- Manage user access
- Implement fine-grained access control
- Audit access and changes

## Audit Logging

[Audit Logging](audit_logging.md) helps you track all actions performed with Nessi:

- Log all CLI commands
- Track configuration changes
- Monitor access to sensitive data
- Generate audit reports

## Feature Flags

[Feature Flags](feature_flags.md) allow you to control the availability of features:

- Enable/disable features dynamically
- Test new features in production
- Implement gradual rollouts
- Configure feature-specific settings
