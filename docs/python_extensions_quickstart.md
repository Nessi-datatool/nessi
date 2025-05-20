# Python Extensions Quick Start Guide

This guide will help you quickly get started with Nessi.dev's Python extensions, focusing on practical examples for common use cases.

## Prerequisites

Before starting, ensure you have:
- Nessi.dev installed
- Python 3.8 or later
- Required Python packages installed (see [Python Extensions Installation Guide](python_extensions_installation.md))

## Enabling Python Extensions

1. Check your Python environment:
   ```bash
   nessi python-extensions check
   ```

2. Install any missing packages:
   ```bash
   cd scripts/python
   pip install -r requirements.txt
   ```

3. Enable Python extensions:
   ```bash
   nessi python-extensions enable
   ```

4. Verify that extensions are enabled:
   ```bash
   nessi python-extensions status
   ```

## ML-Based Anomaly Detection

### Detecting Outliers in Time Series Data

```bash
# Detect anomalies in a metric using Isolation Forest algorithm
nessi anomaly detect --metric cpu_usage --algorithm isolation_forest --contamination 0.05

# Detect anomalies in a CSV file
nessi anomaly detect-file --file metrics.csv --value-column value --algorithm one_class_svm
```

### Python API Example

```python
from nessi.client import NessiClient
from nessi.anomaly import AnomalyDetector

# Initialize client
client = NessiClient(base_url="http://localhost:8080", api_key="your-api-key")

# Get metrics data
metrics = client.get_metrics(name="cpu_usage", limit=100)
values = [m.value for m in metrics]

# Detect anomalies
detector = AnomalyDetector()
result = detector.detect_anomalies(
    data=values,
    algorithm="isolation_forest",
    contamination=0.05
)

# Print anomalies
print(f"Found {len(result.anomaly_indices)} anomalies:")
for i, value in zip(result.anomaly_indices, result.anomaly_values):
    print(f"  Position {i}: Value {value}")

# Display visualization (if in a notebook)
from IPython.display import Image
import base64
Image(data=base64.b64decode(result.visualization))
```

## Seasonal Decomposition

Break down time series data into trend, seasonal, and residual components:

```bash
# Decompose a metric with weekly seasonality
nessi anomaly decompose --metric weekly_sales --period 7

# Decompose a CSV file with daily seasonality
nessi anomaly decompose-file --file daily_data.csv --value-column value --period 24
```

### Python API Example

```python
from nessi.client import NessiClient
from nessi.anomaly import AnomalyDetector

# Initialize client
client = NessiClient(base_url="http://localhost:8080", api_key="your-api-key")

# Get metrics data
metrics = client.get_metrics(name="daily_sales", limit=90)  # 3 months of daily data
values = [m.value for m in metrics]

# Decompose time series
detector = AnomalyDetector()
result = detector.decompose_time_series(
    data=values,
    period=7  # Weekly seasonality
)

# Access components
trend = result.trend
seasonal = result.seasonal
residual = result.residual

# Display visualizations (if in a notebook)
from IPython.display import Image
import base64
Image(data=base64.b64decode(result.visualizations["trend"]))
Image(data=base64.b64decode(result.visualizations["seasonal"]))
Image(data=base64.b64decode(result.visualizations["residual"]))
```

## Advanced Delta Lake Features

When Python extensions are enabled, you can use advanced Delta Lake features:

```bash
# Generate schema recommendations based on data patterns
nessi delta schema-recommend --path /path/to/delta/table

# Analyze performance patterns
nessi delta performance-analyze --path /path/to/delta/table

# Detect data drift between versions
nessi delta drift-detect --path /path/to/delta/table --version1 5 --version2 10
```

## Custom Rule Execution

Execute custom Python-based validation rules:

```bash
# Run a custom Python rule
nessi rules run-custom --rule custom_validation.py --path /path/to/data

# Create a new custom rule from template
nessi rules create-custom --name my_custom_rule
```

Example custom rule (my_custom_rule.py):
```python
def validate(data, config):
    """
    Custom validation rule
    
    Args:
        data: Pandas DataFrame with the data to validate
        config: Dictionary with rule configuration
        
    Returns:
        Dictionary with validation results
    """
    column = config.get("column", "value")
    threshold = config.get("threshold", 100)
    
    # Check if any values exceed the threshold
    violations = data[data[column] > threshold]
    
    return {
        "valid": len(violations) == 0,
        "violations": len(violations),
        "violation_indices": violations.index.tolist(),
        "message": f"Found {len(violations)} values exceeding threshold {threshold}"
    }
```

## Advanced Visualization

Generate sophisticated visualizations with Python:

```bash
# Create an advanced visualization for a metric
nessi visualize advanced --metric cpu_usage --type heatmap

# Export visualization to file
nessi visualize advanced --metric network_traffic --type network_graph --output network.png
```

## Troubleshooting

### Common Issues

**Q: I get "Python extensions are disabled" error**  
A: Make sure you've enabled Python extensions with `nessi python-extensions enable`

**Q: Missing Python packages error**  
A: Install the required packages with `pip install -r scripts/python/requirements.txt`

**Q: Visualization doesn't display**  
A: Ensure matplotlib is properly installed and configured for your environment

### Checking Python Bridge Status

To check if the Python bridge is working correctly:

```bash
# Verbose check of Python environment
nessi python-extensions check --verbose

# Test Python bridge
nessi python-extensions test-bridge
```

## Next Steps

- Explore the [Python Extensions Documentation](python_extensions.md) for more details
- Check out the [Delta Lake Features Guide](delta_lake_features.md) for advanced Delta Lake operations
- See the [Intelligent Alerting Documentation](intelligent_alerting.md) for ML-powered alerts
