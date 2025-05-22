# Nessi License Management User Guide

This guide provides detailed information on how to use Nessi's license management system, including how to start a trial, check your license status, and understand the different license tiers.

## License Tiers Overview

Nessi offers the following license tiers:

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

Builds on the Community Edition with enterprise features:

| Category | Additional Features |
|----------|---------------------|
| **Delta Lake** | Cloud storage integration, Enhanced performance for large tables |
| **Data Quality** | Advanced validation, Custom rule engines, Cross-table validation |
| **Reporting** | Enhanced visualizations, Custom branding, Scheduled reporting |
| **Integrations** | AWS S3, Azure, GCP, Databricks, dbt, Data catalogs |
| **Workflow** | Airflow integration, Prefect integration, Dagster integration |
| **Support** | Priority support, SLA guarantees, Direct assistance |

### Enterprise Edition

Contact-only tier for custom features and dedicated support (available through [nessi.dev](https://nessi.dev)).

## Command Reference

### Starting a Trial

To start a free trial that gives you access to all premium features for 30 days:

```bash
nessi license start-trial
```

This will create a trial that's valid for 30 days. You can only have 2 trials per machine.

### Checking License Status

To check your current license status:

```bash
nessi license info
```

This will display information about your license, including:
- License status (Community, Licensed, Trial)
- Plan tier (Community, Starter, Pro, Enterprise)
- Whether you're using a trial
- Expiration date (for trials and time-limited licenses)
- Days remaining (for trials and time-limited licenses)

### Activating a License

When you purchase a license, you'll receive a license key. To activate it:

```bash
nessi license activate --key YOUR_LICENSE_KEY
```

### Deactivating a License

If you need to move your license to another machine:

```bash
nessi license deactivate
```

## Understanding Trial Limitations

The trial system has the following limitations:

1. **Duration**: Each trial lasts for 30 days from the activation date
2. **Number of Trials**: You can have a maximum of 2 trials per machine
3. **Features**: During the trial, you have access to all features, including Enterprise tier features

## Troubleshooting

### "Maximum number of trials already used on this machine"

This error occurs when you've already used 2 trials on your current machine. Each machine is identified by a unique hardware identifier to prevent abuse of the trial system.

**Solution**: Purchase a license to continue using premium features.

### "Invalid signature, data may have been tampered with"

This error indicates that the license or trial files have been modified. The license system uses cryptographic signatures to ensure data integrity.

**Solution**: If you're in development mode, you can set the `NESSI_DEV_MODE=true` environment variable to bypass signature verification. In production, you should restore the original files or start a new trial.

### "License file not found"

This error occurs when you try to use a premium feature without a valid license or active trial.

**Solution**: Start a trial or purchase a license to access premium features.

## Frequently Asked Questions

### How do I know which features require a license?

Features that require a license will display an error message when you try to use them without a valid license or active trial. You can also refer to the [License Management documentation](LICENSE_MANAGEMENT.md) for a complete list of premium features.

### Can I extend my trial?

The trial system is limited to 30 days per trial and a maximum of 2 trials per machine. If you need more time to evaluate Nessi, please contact our sales team.

### How do machine IDs work?

Each machine is assigned a unique identifier based on hardware characteristics. This identifier is used to track trial usage and ensure that licenses are used on authorized machines only.

### Is my license information secure?

Yes, all license and trial information is cryptographically signed to prevent tampering. The license system also includes anti-tampering measures to detect and prevent unauthorized modifications.

## Support

If you have any questions or issues with the license management system, please contact our support team at licensing@nessi-dev.com.
