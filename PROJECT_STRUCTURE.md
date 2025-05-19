# Nessi Project Structure

This document describes the structure of the Nessi open-source project.

```
nessi-dev/
├── cmd/
│   └── nessi/                    # Main CLI application
│       └── main.go               # Entry point
├── internal/
│   ├── api/                      # API Gateway implementation
│   │   ├── handlers/             # API request handlers
│   │   ├── middleware/           # API middleware components
│   │   └── router.go             # API route definitions
│   ├── config/                   # Configuration management
│   │   ├── config.go             # Configuration structures
│   │   └── loader.go             # Config loading logic
│   ├── delta/                    # Delta Lake core operations
│   │   ├── reader.go             # Delta table reading
│   │   ├── metadata.go           # Metadata handling
│   │   └── writer.go             # Delta table writing
│   ├── quality/                  # Data quality engine
│   │   ├── profile/              # Data profiling
│   │   │   └── outlier_detection.go  # IQR-based outlier detection
│   │   ├── rules/                # Rule validation
│   │   │   ├── validator.go      # Core validation logic
│   │   │   └── extended_rules.go # Extended rule types
│   │   ├── patterns/             # Pattern matching
│   │   └── scoring/              # Quality scoring
│   ├── monitor/                  # Monitoring service
│   │   ├── metrics.go            # Metrics collection
│   │   ├── prometheus.go         # Prometheus integration
│   │   └── alerts.go             # Alert management
│   ├── report/                   # Report generation
│   │   ├── html.go               # HTML report templates
│   │   ├── pdf.go                # PDF generation
│   │   └── export.go             # JSON/CSV exports
│   ├── security/                 # Security components
│   │   ├── auth.go               # Authentication
│   │   ├── rbac.go               # Access control
│   │   └── tls.go                # TLS configuration
│   └── server/                   # HTTP server implementation
│       ├── server.go             # Server setup
│       └── graceful.go           # Graceful shutdown
├── pkg/                          # Public packages
│   ├── datalake/                 # Delta Lake & multi-format operations
│   │   ├── delta.go              # Delta table interface
│   │   ├── parquet.go            # Parquet file operations
│   │   └── schema.go             # Schema management
│   ├── dbt/                      # dbt plugin integration
│   │   ├── config.go             # Plugin configuration
│   │   ├── validator.go          # dbt model validation
│   │   ├── profiler.go           # dbt model profiling
│   │   ├── lineage.go            # dbt lineage integration
│   │   ├── alert.go              # Alerting for validation results
│   │   ├── delta.go              # Delta table integration
│   │   └── output.go             # Output formatting (JSON, CSV, table)
│   ├── quality/                  # Data quality & profiling
│   │   ├── profile/              # Data profiling & statistics
│   │   │   ├── profiler.go       # Core profiling functionality
│   │   │   └── advanced_profiling.go  # Enhanced profiling capabilities
│   │   ├── rules/                # Rule validation engine
│   │   │   ├── validator.go      # Core validation logic
│   │   │   ├── extended_rules.go # Extended rule types (Regex, Enum, Length, DateFormat)
│   │   │   └── yaml_loader.go    # YAML configuration support
│   │   └── anomaly/              # Anomaly detection algorithms
│   │       └── sudden_change.go  # Sudden change detection
│   ├── monitoring/               # Metrics & monitoring
│   │   ├── monitoring.go         # Core monitoring functionality
│   │   ├── dashboard/            # Web dashboard
│   │   │   ├── dashboard.go      # Core dashboard functionality
│   │   │   ├── data_quality_dashboard.go  # Data quality UI
│   │   │   ├── data_quality_handlers.go   # Data quality API handlers
│   │   │   └── templates/        # Dashboard templates
│   │   └── alerts/               # Alerting system
│   ├── security/                 # Authentication & authorization
│   │   ├── auth.go               # Authentication system
│   │   ├── ssl.go                # SSL/TLS support
│   │   └── rbac.go               # Role-based access control
│   ├── extensions/               # Extension system
│   │   ├── manager.go            # Extension management
│   │   ├── python.go             # Python extension interface
│   │   └── registry.go           # Extension registry
│   ├── models/                   # Shared data models
│   │   ├── delta.go              # Delta table models
│   │   ├── profile.go            # Data profile models
│   │   └── quality.go            # Quality metrics models
│   └── util/                     # Utility functions
│       ├── files.go              # File operations
│       └── concurrency.go        # Concurrency helpers
├── scripts/                      # Build and deployment scripts
│   ├── build.sh                  # Build script
│   └── install.sh                # Installation script
├── docker/                       # Docker configurations
│   ├── go/                       # Go service Dockerfile
│   ├── python/                   # Python extensions Dockerfile
│   └── docker-compose.yml        # Multi-service composition
├── python/                       # Python extension modules
│   ├── delta_advanced/           # Advanced Delta Lake features
│   ├── ml_anomaly/               # ML-based anomaly detection
│   └── setup.py                  # Python package setup
├── tests/                        # Test files
│   ├── unit/                     # Unit tests
│   ├── integration/              # Integration tests
│   └── e2e/                      # End-to-end tests
├── docs/                         # Documentation
│   ├── api/                      # API documentation
│   ├── user/                     # User guides
│   └── examples/                 # Example use cases
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── README.md                     # Project overview
└── LICENSE                       # License information
```