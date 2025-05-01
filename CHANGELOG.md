# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-05-01

### Added
- Secure license system with encryption and machine binding
- 14-day trial period with automatic initialization
- License key validation and activation via CLI
- `@require_valid_license` decorator for feature protection
- Delta Lake integration and table scanning
- Sample data generation capabilities
- Comprehensive test suite
- Security features:
  - Config file encryption using Fernet
  - Machine-specific binding
  - Tamper detection and recovery
  - Protection against license sharing
  - Clock manipulation protection

### Changed
- Updated to proprietary license with trial period
- Enhanced security measures
- Improved CLI with license management commands
- Better error handling with custom LicenseError

### Security
- Implemented secure storage for license information
- Added encryption for configuration files
- Added protection against tampering and unauthorized modifications
- Implemented machine-specific license binding
- Added automatic recovery from corrupted configurations

### Documentation
- Added comprehensive API documentation
- Updated installation and usage guides
- Added security documentation
- Added contribution guidelines
- Added code of conduct

## [0.1.0] - 2024-05-01

### Added
- Initial release with trial and license system
- 14-day trial period
- License key generation and validation
- License management CLI
- Trial status checking
- License activation
- Security features for license protection

### Changed
- Updated to proprietary license
- Enhanced security measures
- Improved documentation

### Fixed
- Security vulnerabilities
- License validation issues

## [0.0.1] - 2024-04-30

### Added
- Initial project structure
- Basic data processing capabilities
- Delta Lake integration
- Table scanning functionality
- Sample data generation
- Test suite
- Documentation 