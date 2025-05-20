# Remove Unused Grafana/Prometheus Integrations and Update to CLI-Only Approach

## Overview
This PR removes the unused Grafana and Prometheus integrations from the Nessi project and updates the documentation to reflect the CLI-only approach. The dashboard component has been removed, and the reporting functionality has been streamlined to focus on HTML and PDF report generation through the CLI.

## Changes

### Removed Components
- Removed the `pkg/monitoring/dashboard` directory and all its contents
- Removed Grafana configuration files in `config/grafana`
- Removed Prometheus configuration files in `config/prometheus`
- Removed dashboard-related documentation

### Added Components
- Added a new `docs/REPORTING.md` file documenting the reporting capabilities
- Added a new `docs/alerts_cli.md` file documenting the CLI-based alerting system
- Added simple metrics implementation that doesn't rely on external dependencies

### Updated Components
- Updated `README.md` to focus on CLI-based reporting and monitoring
- Updated `REQUIREMENTS.md` to emphasize CLI-only monitoring and reporting
- Updated `IMPLEMENTATION.md` to remove dashboard components
- Updated `freshness_monitoring.md` to replace dashboard sections with CLI reporting
- Updated `intelligent_alerting.md` to focus on CLI-based reporting
- Updated `USER_GUIDE.md` to remove web dashboard references

## Testing
The changes have been tested to ensure that the core functionality of HTML and PDF report generation is maintained. All tests pass with the updated implementation.

## Documentation
Documentation has been updated to reflect the CLI-only approach, with comprehensive examples of CLI commands for generating reports in various formats.

## Related Issues
This PR addresses the need to clean up unused integrations and streamline the codebase, focusing on the core functionality of HTML and PDF report generation.
