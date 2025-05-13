# Python Extensions Installation Guide

Nessi.dev provides optional Python-based features that extend the core Go-based functionality. These extensions enable advanced capabilities like ML-based anomaly detection, sophisticated Delta Lake operations, and advanced visualizations.

## Prerequisites

To use the Python extensions, you need:

1. Python 3.8 or later installed on your system
2. Pip package manager
3. Required Python packages (listed below)

## Installation Steps

### 1. Check Your Environment

First, check if your environment is ready for Python extensions:

```bash
nessi python-extensions check
```

This command will:
- Verify Python is installed and accessible
- Check for required packages
- Report any missing dependencies

### 2. Install Required Packages

If you're missing any required packages, install them using pip:

```bash
# Navigate to the Nessi.dev scripts/python directory
cd scripts/python

# Install all required packages
pip install -r requirements.txt
```

Alternatively, you can install packages individually:

```bash
pip install pandas numpy scikit-learn matplotlib pyarrow deltalake statsmodels
```

### 3. Enable Python Extensions

Once all dependencies are installed, enable the Python extensions:

```bash
nessi python-extensions enable
```

This will:
- Verify your Python environment
- Enable Python extensions in the configuration
- Make Python-based features available in Nessi.dev

### 4. Verify Installation

To verify that Python extensions are properly enabled:

```bash
nessi python-extensions status
```

This will show you:
- Whether Python extensions are enabled
- Which specific Python-based features are available
- The status of your Python environment

## Available Python Extensions

When enabled, the following features become available:

1. **ML-based Anomaly Detection**
   - Advanced outlier detection using machine learning algorithms
   - Time series decomposition and pattern recognition
   - Visualization of anomalies and patterns

2. **Advanced Delta Lake Features**
   - Enhanced schema evolution tracking
   - Advanced metadata analysis
   - Performance optimization recommendations

3. **Advanced Visualization Components**
   - Interactive data visualizations
   - Custom dashboards with Python-powered charts
   - Export capabilities for reports

## Configuration

Python extensions can be configured in your Nessi.dev configuration file:

```yaml
python_extensions:
  enabled: true
  python_path: "python"  # Path to Python executable
  required_packages:
    - pandas
    - numpy
    - scikit-learn
    - matplotlib
    - pyarrow
    - deltalake
    - statsmodels
  features:
    ml_anomaly_detection: true
    advanced_delta_lake: true
    advanced_visualization: true
```

## Troubleshooting

### Python Not Found

If the `python` command is not found, ensure Python is installed and in your PATH. You can specify a custom Python path in the configuration file.

### Missing Packages

If you encounter errors about missing packages:

1. Run `nessi python-extensions check` to see which packages are missing
2. Install the missing packages using pip
3. Try enabling Python extensions again

### Performance Issues

If you experience performance issues with Python extensions:

1. Ensure you're using the latest versions of Python packages
2. Consider using a virtual environment with optimized packages
3. For large datasets, increase the memory allocation for Python processes

## Disabling Python Extensions

If you need to disable Python extensions:

```bash
nessi python-extensions disable
```

This will revert Nessi.dev to using only the Go-based features.
