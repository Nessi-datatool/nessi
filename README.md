# Nessi Data Tool

[![License](https://img.shields.io/badge/License-Proprietary-red.svg)](LICENSE)

A powerful data processing and analysis tool built with PySpark and Delta Lake.

## Features

- Secure license management system with machine-specific binding
- Sample data generation for testing and development
- Support for multiple file formats (Delta Lake, Parquet, CSV)
- Table scanning and profiling capabilities
- Comprehensive test suite
- Easy-to-use CLI interface

## License

Nessi is provided with a 14-day trial license. During the trial period, you can evaluate all features of the software. After the trial period, you will need to purchase a license to continue using the software.

### Security Features

The license system includes several security measures:
- Machine-specific license binding
- Encrypted configuration storage
- Tamper detection and automatic recovery
- Protection against clock manipulation
- Secure trial period tracking

### License Management

The following commands are available for managing your license:

```bash
# Check current license status
nessi status

# Activate a license key
nessi activate <license_key>
```

### Trial Period

- 14-day trial period starts automatically on first use
- Full access to all features during trial
- No registration required
- Trial status can be checked using `nessi status`
- Secure trial period tracking prevents manipulation

### License Activation

After purchasing a license, activate it using:

```bash
nessi activate <your-license-key>
```

The license will be bound to your specific machine for security.

### Using Licensed Features

All main functions in Nessi require a valid license or active trial period. The license is checked automatically when using any feature. For example:

```python
from nessi import scan_table, generate_sample_data

# These functions will check for a valid license automatically
scan_table('path/to/table')
generate_sample_data('output/path')
```

If your license is invalid or expired, the functions will raise a `LicenseError` with a descriptive message.

## Installation

```bash
pip install nessi
```

### Dependencies

Nessi requires the following main dependencies:
- Python 3.8 or higher
- PySpark 3.5.0 or higher
- Delta Lake 2.4.0 or higher
- cryptography 41.0.0 or higher

## Usage

### Basic Example

```python
from nessi import scan_table, generate_sample_data

# Generate sample data
generate_sample_data('test_table', format='delta', rows=1000)

# Scan the table
report = scan_table('test_table', format='delta')
```

### Protected Features

All main features are protected by the license system using the `@require_valid_license` decorator:

```python
from nessi.decorators import require_valid_license

@require_valid_license
def my_custom_function():
    # This function will only run with a valid license
    pass
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and development process.

## Security

For security-related information and reporting vulnerabilities, please see [SECURITY.md](SECURITY.md).

## Support

For license inquiries or support:
- Email: support@nessi-datatool.com
- Website: https://nessi-datatool.com

## Documentation

For detailed documentation, please visit our [documentation site](https://docs.nessi-datatool.com). 