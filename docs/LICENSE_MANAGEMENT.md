# License Management System

Nessi includes a robust license management system that controls access to premium features while providing a free trial option for evaluation.

## License Tiers

Nessi offers the following license tiers:

### Community Edition (Open Source)

The Community Edition is free and open-source under the Apache License 2.0. It includes all the core functionality needed for local Delta Lake management and basic data quality checks.

### Pro Edition (Commercial)

The Pro Edition requires a paid license and includes all Community features plus additional enterprise capabilities. A 1-month free trial is available.

### Enterprise Edition (Contact Only)

For enterprise-level support and custom features, please visit the [Nessi website](https://nessi.dev) for contact information.

## Feature Breakdown

### Core Features (Community Edition)

The following features are available in the free, open-source Community Edition:

- **Delta Lake Support**
  - Schema evolution tracking
  - Transaction log parsing and analysis
  - Version control with commit history
  - Time travel capabilities (up to 30 days)

- **Data Quality**
  - Basic validation checks
  - Data profiling and statistics
  - Anomaly detection
  - Pattern recognition

- **Reporting**
  - HTML/PDF reports for local tables
  - Basic visualizations
  - Export capabilities

- **Integrations**
  - Local file systems
  - SQLite support
  - CSV/JSON imports

- **Workflow**
  - Basic CLI commands
  - Simple automation scripts

- **Monitoring**
  - Basic metrics collection
  - CLI-based reporting

### Premium Features (Pro Edition)

The following additional features require a Pro license or active trial:

- **Delta Lake Support**
  - Cloud storage integration
  - Enhanced performance for large tables

- **Data Quality**
  - Advanced validation rules
  - Custom rule engines
  - Cross-table validation

- **Reporting**
  - Scheduled reporting

- **Integrations**
  - AWS S3 cloud storage
  - Azure Blob Storage
  - Google Cloud Storage
  - Databricks integration
  - dbt integration
  - Data catalog integration

- **Workflow**
  - Airflow integration
  - Prefect integration
  - Dagster integration

- **Support**
  - Priority support
  - SLA guarantees
  - Direct assistance

## Open-Source vs. Commercial Features

### Value Proposition

Nessi's licensing model is designed to provide significant value in both the open-source Community Edition and the commercial Pro tier:

**Community Edition (Open Source)**
* Free and open-source under the Apache License 2.0
* Complete functionality for local Delta Lake management
* Suitable for individual users, small teams, and educational purposes
* Full access to all core features and documentation

**Pro Tier (Commercial)**
* Built on top of the Community Edition
* Adds enterprise integrations and advanced capabilities
* Designed for production environments and larger teams
* Includes priority support and enterprise-grade features

### Feature Comparison

The table below provides a detailed comparison of what's included in each tier:

| Category | Community Features | Additional Pro Features |
|----------|-------------------|------------------------|
| **Delta Lake** | Schema evolution tracking<br>Transaction history<br>Time travel<br>Version control | Cloud storage integration<br>Enhanced performance for large tables |
| **Data Quality** | Basic validation rules<br>Profiling<br>Anomaly detection<br>Pattern matching | Advanced validation<br>Custom rule engines<br>Cross-table validation |
| **Reporting** | HTML/PDF reports<br>Basic visualizations<br>Export capabilities | Scheduled reporting |
| **Integrations** | Local file systems<br>SQLite support<br>CSV/JSON imports | AWS S3, Azure, GCP<br>Databricks<br>dbt<br>Data catalogs |
| **Workflow** | CLI commands<br>Basic automation | Airflow integration<br>Prefect integration<br>Dagster integration |
| **Support** | Community forums<br>Documentation<br>Issue tracker | Priority support<br>SLA guarantees<br>Direct assistance |

### Transparency in Licensing

The license management system itself is open-source and transparent, allowing for community review and contributions. The anti-tampering measures only restrict access to premium features but do not obscure the implementation. This approach ensures that:

1. Users can fully understand how the licensing system works
2. The community can contribute to all aspects of the codebase
3. Premium features are clearly separated from open-source functionality
4. The integrity of the licensing system is maintained

## Free Trial

Nessi offers a 1-month free trial that provides access to all premium features, including those in the Enterprise tier. Key aspects of the trial system:

- **Duration**: 30 days from activation
- **Limitations**: Maximum of 2 trials per machine
- **Security**: Trial information is cryptographically signed to prevent tampering

## License Commands

### Starting a Trial

```bash
# Start a free trial to access premium features
nessi license start-trial
```

### Checking License Status

```bash
# View your current license information
nessi license info
```

Example output:
```
License Status: Trial
Plan: Enterprise
Is Trial: Yes
Expires: 2025-06-22
Days Remaining: 30
Message: Free trial active. 30 days remaining.
```

### Activating a License

```bash
# Activate a purchased license
nessi license activate --key YOUR_LICENSE_KEY
```

## Technical Implementation

The license management system uses several security measures to ensure the integrity of license validation:

1. **Machine ID**: Each machine is uniquely identified using hardware characteristics
2. **Signature Verification**: All license and trial data is signed using HMAC-SHA256
3. **Anti-Tampering**: Multiple validation checks prevent binary modification
4. **Secure Storage**: License information is stored securely with access controls

## Troubleshooting

### Trial Expiration

When a trial expires, premium features will no longer be accessible. You'll need to purchase a license to continue using these features.

### Maximum Trials Reached

If you've already used 2 trials on your machine, you won't be able to start another trial. Consider purchasing a license to continue using premium features.

### Invalid License

If your license is reported as invalid, check the following:
- License key is entered correctly
- License hasn't expired
- License is for the correct machine ID

## Support

For licensing questions or issues, contact support at licensing@nessi-dev.com.
