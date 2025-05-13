# Python Extensions for Nessi.dev

## Overview

Nessi.dev follows a Go-first architecture with optional Python extensions that provide advanced capabilities. This document explains how to enable and use these Python-based features.

## Why Python Extensions?

While Nessi.dev's core functionality is implemented in Go for performance and efficiency, Python offers several advantages for specific use cases:

- **Advanced Machine Learning**: Python's ecosystem (scikit-learn, TensorFlow, etc.) enables sophisticated anomaly detection and pattern recognition
- **Data Science Tools**: Libraries like pandas and numpy provide powerful data manipulation capabilities
- **Visualization**: Libraries like matplotlib and plotly offer advanced visualization options
- **Flexibility**: Python scripts allow for highly customized rule execution and data processing

## Enabling Python Extensions

Python extensions are disabled by default and can be enabled through configuration:

### Via Configuration File

In your `config.yaml` file:

```yaml
feature_flags:
  python_extensions: true
  # Enable specific Python features
  ml_anomaly_detection: true
  advanced_delta_lake: true
  advanced_visualization: true
```

### Via CLI

```bash
nessi feature-flags enable python_extensions
nessi feature-flags enable ml_anomaly_detection
```

## Prerequisites

When Python extensions are enabled, Nessi.dev requires:

1. **Python 3.8+** installed and available in PATH
2. **Required packages**:
   - pandas
   - numpy
   - scikit-learn
   - matplotlib
   - pyarrow
   - deltalake

You can install these dependencies using:

```bash
pip install -r python/requirements.txt
```

## Available Python Extensions

### 1. ML-based Anomaly Detection

Advanced anomaly detection using machine learning algorithms:

- Isolation Forest for outlier detection
- DBSCAN for density-based clustering
- Autoencoder neural networks for complex pattern detection

### 2. Advanced Delta Lake Features

Enhanced Delta Lake capabilities:

- Schema evolution tracking with detailed diff analysis
- Advanced partition optimization recommendations
- Compression ratio analysis and optimization

### 3. Advanced Statistical Analysis

Sophisticated statistical capabilities:

- Multivariate analysis
- Correlation detection between columns
- Distribution fitting and hypothesis testing

### 4. Custom Rule Execution Engine

Flexible rule execution:

- Python-based rule definitions
- Complex validation logic using Python libraries
- Custom aggregations and transformations

### 5. Advanced Visualization Components

Rich visualization options:

- Interactive dashboards
- Custom chart types
- Exportable reports with embedded visualizations

## Troubleshooting

If you encounter issues with Python extensions:

1. Verify Python is correctly installed: `python --version`
2. Check required packages: `pip list`
3. Enable debug logging: `nessi --log-level debug`
4. Check the Python environment detection logs

## Performance Considerations

Python extensions may have performance implications:

- Startup time may increase when Python extensions are enabled
- Memory usage will be higher due to the Python interpreter
- For performance-critical deployments, consider using only the Go-based features

## Security Considerations

When enabling Python extensions:

- Ensure Python and its dependencies are kept updated
- Be cautious with custom Python rules that may execute arbitrary code
- Consider using virtual environments to isolate dependencies
