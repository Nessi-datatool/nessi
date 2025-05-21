# Nessi Quick Start Guide

This guide will help you get started with Nessi quickly and efficiently.

## Installation

### One-Line Installation (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/nessi-dev/nessi/main/scripts/install.sh | bash
```

### Windows Installation

```powershell
Invoke-Expression (New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/nessi-dev/nessi/main/scripts/install.ps1')
```

### Using Homebrew (macOS)

```bash
brew tap nessi-dev/nessi
brew install nessi
```

### Manual Installation

Download the appropriate binary for your platform from the [GitHub Releases page](https://github.com/nessi-dev/nessi/releases).

## Basic Usage

### Check a Delta Lake Table

```bash
nessi check /path/to/delta/table
```

### Generate a Profile Report

```bash
nessi profile /path/to/delta/table --format html --output report.html
```

### View Version Information

```bash
nessi version
```

## Common Workflows

### Data Quality Workflow

1. **Check data quality**:
   ```bash
   nessi check /path/to/delta/table
   ```

2. **Generate a detailed profile**:
   ```bash
   nessi profile /path/to/delta/table --format html --output profile.html
   ```

3. **Share the report with stakeholders**:
   Open the HTML report in a browser and use the social sharing buttons.

### CI/CD Integration

Add Nessi to your CI/CD pipeline using our GitHub Actions workflow:

```yaml
name: Data Quality Check

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  quality-check:
    uses: nessi-dev/nessi/.github/workflows/nessi-quality-check.yml@main
    with:
      table_path: 's3://my-bucket/my-table'
      quality_threshold: '95'
      report_format: 'html'
      fail_on_issues: 'true'
```

## Command Reference

### Global Flags

- `--help`: Show help information
- `--config`: Specify a custom configuration file

### Check Command

```bash
nessi check [table_path] [flags]
```

Flags:
- `--threshold`: Quality threshold percentage (default: 90)
- `--format`: Output format (html, pdf, json)
- `--output`: Output file path

### Profile Command

```bash
nessi profile [table_path] [flags]
```

Flags:
- `--format`: Output format (html, pdf, json)
- `--output`: Output file path
- `--include-samples`: Include data samples in the profile

## Next Steps

- Read the [full documentation](https://github.com/nessi-dev/nessi/docs)
- Check out the [roadmap](https://github.com/nessi-dev/nessi/docs/ROADMAP.md) for upcoming features
- Join our [community](https://github.com/nessi-dev/nessi/discussions)
- Star the [GitHub repository](https://github.com/nessi-dev/nessi) if you find Nessi useful!
