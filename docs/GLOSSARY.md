# Nessi Glossary

This glossary provides definitions for terms used throughout the Nessi documentation and application.

## A

### Access Control
A security mechanism that regulates who or what can view or use resources in a computing environment. In Nessi, access controls determine who can access tables, run quality checks, and view reports.

### Accuracy
A data quality dimension that measures how well data values reflect the true values they are meant to represent. Nessi calculates accuracy metrics based on defined quality rules.

### Alert
A notification triggered when a data quality metric falls below a defined threshold. Nessi can generate alerts for quality issues to notify stakeholders.

### Anomaly
A data point or pattern that deviates significantly from the expected behavior. Nessi can detect anomalies in data quality metrics over time.

### API (Application Programming Interface)
A set of rules and protocols that allows different software applications to communicate with each other. Nessi provides APIs for programmatic integration with other systems.

### Arrow
An in-memory columnar data format designed for efficient data processing. Nessi uses Apache Arrow for high-performance data operations.

### Audit Log
A chronological record of system activities that provides documentary evidence of the sequence of activities affecting an operation, procedure, or event. Nessi maintains audit logs for tracking operations.

## B

### Batch Processing
The processing of data in groups or batches rather than continuously. Nessi supports batch processing for data quality checks.

### Benchmark
A standard or point of reference against which data quality metrics can be compared. Nessi allows benchmarking of quality metrics over time.

### Business Rule
A rule that defines or constrains some aspect of business data. In Nessi, business rules can be implemented as quality rules.

## C

### Cardinality
The number of distinct values in a dataset or column. Nessi can calculate cardinality as part of data profiling.

### Catalog
A collection of metadata about datasets. In Nessi, a catalog provides information about available tables and their schemas.

### CLI (Command Line Interface)
A text-based interface used to run programs, manage files, and interact with the computer. Nessi provides a comprehensive CLI for all operations.

### Completeness
A data quality dimension that measures the presence of all required data. Nessi calculates completeness metrics based on the presence of non-null values.

### Compliance
Adherence to laws, regulations, guidelines, and specifications relevant to data. Nessi helps ensure data compliance through quality validation.

### Consistency
A data quality dimension that measures how well data adheres to defined formats and values across datasets. Nessi calculates consistency metrics based on defined quality rules.

### Community Edition
The open-source version of Nessi that includes core features for data quality monitoring and reporting.

## D

### Dashboard
A visual display of data quality metrics and key performance indicators. Nessi can generate dashboards for monitoring data quality.

### Data Dictionary
A centralized repository of information about data, including meaning, relationships, origin, usage, and format. Nessi can generate data dictionaries from schema information.

### Data Lake
A centralized repository that allows you to store structured and unstructured data at any scale. Nessi works with data lakes to provide quality monitoring.

### Data Lineage
The data's origin and where it moves over time, providing visibility into the data analytics pipeline. Nessi can track data lineage for quality metrics.

### Data Profiling
The process of examining data to collect statistics and information about that data. Nessi performs data profiling as part of quality checks.

### Data Quality
The measure of how well data serves its intended purpose in a specific context. Nessi provides comprehensive data quality monitoring and validation.

### Data Type
A classification that specifies which type of value a variable has. Nessi validates data types as part of schema validation.

### Data Validation
The process of ensuring that data is accurate, consistent, and meets specified criteria. Nessi provides data validation through quality rules.

### Databricks
A cloud-based data engineering platform based on Apache Spark. Nessi integrates with Databricks for data quality monitoring (Pro Edition).

### Delta Lake
An open-source storage layer that brings reliability to data lakes. Nessi has built-in support for Delta Lake tables.

### Distribution
The frequency and pattern of values in a dataset. Nessi can analyze data distributions as part of data profiling.

### Drift
Changes in data characteristics or statistics over time. Nessi can detect data drift through trend analysis of quality metrics.

## E

### ETL (Extract, Transform, Load)
A process that extracts data from source systems, transforms it to fit business needs, and loads it into a target database. Nessi can be integrated into ETL pipelines for quality validation.

### Enterprise Edition
A commercial version of Nessi that includes advanced features for enterprise use (contact-only).

### Error Code
A standardized code that identifies a specific error condition. Nessi uses error codes to provide clear information about failures.

### Exception
A data record that fails to meet quality rules or validation criteria. Nessi identifies exceptions during quality checks.

## F

### Freshness
A data quality dimension that measures how recent the data is relative to its expected update frequency. Nessi calculates freshness metrics based on timestamp columns.

### Format Handler
A component that handles specific data formats. Nessi includes format handlers for Delta Lake, Parquet, CSV, and JSON.

## G

### Governance
The overall management of data availability, usability, integrity, and security. Nessi supports data governance through quality monitoring and reporting.

## H

### Heatmap
A graphical representation of data where values are represented as colors. Nessi can generate heatmaps to visualize data quality issues.

### Histogram
A graphical representation of the distribution of data. Nessi can generate histograms as part of data profiling.

## I

### Incremental Processing
Processing only new or changed data since the last run. Nessi supports incremental processing for efficiency with large datasets.

### Integration
The connection between Nessi and external systems or data sources. Nessi provides integrations with Databricks, AWS S3, and other systems.

### Integrity
The accuracy and consistency of data over its lifecycle. Nessi helps ensure data integrity through quality validation.

## L

### License
A legal instrument governing the use or redistribution of software. Nessi is available under different license tiers (Community Edition and Pro Edition).

### Lineage
See Data Lineage.

## M

### Metadata
Data that provides information about other data. Nessi collects and analyzes metadata as part of data profiling.

### Metric
A quantitative measurement used to track and assess the status of a specific process. Nessi collects various metrics, including quality metrics, performance metrics, and freshness metrics.

### Monitoring
The process of observing and checking the quality and integrity of data over time. Nessi provides comprehensive monitoring capabilities.

## N

### Nessi
A data quality and observability tool designed for data engineers and analysts working with data lakes.

### Null Value
A special marker used to indicate that a data value does not exist. Nessi detects and reports on null values as part of completeness checks.

## O

### Observability
The ability to measure the internal state of a system by examining its outputs. Nessi provides data observability through quality metrics and reporting.

### Outlier
A data point that differs significantly from other observations. Nessi can detect outliers as part of data profiling.

## P

### Parquet
A columnar storage file format designed for efficient data processing. Nessi supports reading and analyzing Parquet files.

### Partition
A division of a logical database or its elements into distinct independent parts. Nessi supports partitioned data processing for better performance.

### Performance Metric
A measurement of the efficiency and speed of data processing operations. Nessi collects performance metrics to help optimize operations.

### Plugin
A software component that adds a specific feature to an existing program. Nessi supports plugins for extending functionality.

### Pro Edition
The commercial version of Nessi that includes premium features such as Databricks integration and AWS S3 support.

### Profiling
See Data Profiling.

## Q

### Quality Check
The process of validating data against defined quality rules. Nessi performs quality checks to ensure data meets quality standards.

### Quality Dimension
A category of data quality characteristics, such as completeness, accuracy, consistency, etc. Nessi measures multiple quality dimensions.

### Quality Metric
A quantitative measurement of a quality dimension. Nessi calculates quality metrics based on quality checks.

### Quality Rule
A definition of what constitutes quality for a specific aspect of data. Nessi allows defining custom quality rules.

### Quality Score
A numerical representation of the overall quality of data. Nessi calculates quality scores based on multiple quality metrics.

## R

### Reference Data
Data used to categorize other data or for relating data to information beyond the boundaries of the enterprise. Nessi can validate data against reference data.

### Regex (Regular Expression)
A sequence of characters that defines a search pattern. Nessi supports regex patterns in quality rules.

### Report
A structured presentation of data quality and metrics information. Nessi generates comprehensive reports in various formats.

### Rule
See Quality Rule.

## S

### Sampling
The process of selecting a subset of data for analysis. Nessi supports sampling for efficient processing of large datasets.

### Schema
The structure of data, including field names, data types, and constraints. Nessi provides tools for schema validation and evolution tracking.

### Schema Evolution
The process of changing a schema over time. Nessi tracks schema evolution to ensure backward compatibility.

### Schema Registry
A system that stores and manages schemas. Nessi can integrate with schema registries for schema validation.

### Schema Validation
The process of verifying that data conforms to a defined schema. Nessi performs schema validation as part of quality checks.

### Streaming
Processing data continuously as it is generated. Nessi supports streaming mode for processing large datasets.

## T

### Table
A collection of related data held in a structured format. Nessi works with tables in data lakes, particularly Delta Lake tables.

### Threshold
A predefined level at which an action is triggered. Nessi uses thresholds to determine when quality issues should generate alerts.

### Time Series
A series of data points indexed in time order. Nessi can analyze time series data for trend detection.

### Time Travel
The ability to access previous versions of data. Nessi supports time travel for Delta Lake tables.

### Timeliness
A data quality dimension that measures how up-to-date data is. Nessi calculates timeliness metrics based on timestamp columns.

### Trend
A general direction in which something is developing or changing. Nessi can detect trends in quality metrics over time.

## U

### Uniqueness
A data quality dimension that measures how well data values are unique within a dataset. Nessi calculates uniqueness metrics based on defined quality rules.

## V

### Validation
See Data Validation.

### Visualization
The graphical representation of data. Nessi provides various visualizations for data quality metrics and issues.

## W

### Workflow
A sequence of tasks that processes data. Nessi can be integrated into data workflows for quality validation.

### Workflow Orchestration
The automation and coordination of complex data workflows. Nessi Pro Edition supports workflow orchestration for quality checks.

---

For more information on Nessi terminology, please visit [nessi.dev](https://nessi.dev) or contact support@nessi.dev.
