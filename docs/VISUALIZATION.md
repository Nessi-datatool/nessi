# Nessi Visualization Guide

## Navigation

- [Documentation Home](README.md)
- [Installation & Quickstart](../QUICKSTART.md)
- [Configuration](CONFIGURATION.md)
- [CLI Reference](cli/README.md)
- [Reporting](REPORTING.md)
- [Quality Rules](QUALITY_RULES.md)
- [User Guide](USER_GUIDE.md)

---

## Overview

Nessi provides a comprehensive set of visualization tools to help you understand and analyze your data quality metrics. These visualizations are available directly in your terminal, making it easy to quickly assess the quality of your data without the need for external tools.

## Available Visualizations

Nessi currently offers three types of visualizations:

1. **Heatmap**: Displays data quality metrics across different dimensions using color coding
2. **Bar Chart**: Shows comparative metrics for different quality dimensions
3. **Scatter Plot**: Allows comparison of two different metrics to identify correlations

## Command Reference

### Heatmap Visualization

The heatmap visualization displays data quality metrics across different dimensions (e.g., time periods and data partitions) using color coding to indicate quality levels.

```bash
# Basic usage
nessi visualize heatmap --table <path_to_table> --metric <metric_name>

# Example
nessi visualize heatmap --table s3://my-bucket/my-table --metric completeness
```

#### Options

| Option | Description | Default |
|--------|-------------|--------|
| `--table`, `-t` | Path to the Delta Lake table (required) | |
| `--metric`, `-m` | Metric to visualize (e.g., completeness, accuracy) | completeness |
| `--partition-column`, `-p` | Partition column to use for the heatmap | |
| `--color-mode`, `-c` | Color mode (gradient, binary) | gradient |
| `--width`, `-w` | Width of the heatmap | 80 |
| `--height`, `-h` | Height of the heatmap | 20 |

#### Example Output

```
Generating heatmap visualization for metric: completeness

  | Jan Feb Mar Apr May Jun Jul Aug Sep Oct Nov Dec
--+------------------------------------------------
A |                                               
B |                                               
C |                                               
D |                                               
E |                                               

Legend: 
   High quality (>80%)     Medium quality (50-80%)     Low quality (<50%)
```

### Bar Chart Visualization

The bar chart visualization shows comparative metrics for different quality dimensions in an easy-to-read bar format, with color coding to highlight good, medium, and poor quality areas.

```bash
# Basic usage
nessi visualize barchart --table <path_to_table>

# Show all available metrics
nessi visualize barchart --table <path_to_table> --all-metrics

# Example
nessi visualize barchart --table s3://my-bucket/my-table
```

#### Options

| Option | Description | Default |
|--------|-------------|--------|
| `--table`, `-t` | Path to the Delta Lake table (required) | |
| `--all-metrics`, `-a` | Show all available metrics | false |

#### Example Output

```
Data Quality Metrics Bar Chart
==============================

Completeness [95.2%] ████████████████████████████████████
Accuracy     [87.6%] ██████████████████████████████
Consistency  [92.1%] ███████████████████████████████████
Validity     [78.5%] ███████████████████████████

Legend:
█████ Good (>80%)   █████ Medium (50-80%)   █████ Poor (<50%)
```

### Scatter Plot Visualization

The scatter plot visualization allows comparison of two different metrics to identify correlations and patterns in your data quality measurements.

```bash
# Basic usage
nessi visualize scatter --table <path_to_table> --x-metric <metric_name> --y-metric <metric_name>

# Example
nessi visualize scatter --table s3://my-bucket/my-table --x-metric completeness --y-metric accuracy
```

#### Options

| Option | Description | Default |
|--------|-------------|--------|
| `--table`, `-t` | Path to the Delta Lake table (required) | |
| `--x-metric`, `-x` | Metric to display on x-axis | completeness |
| `--y-metric`, `-y` | Metric to display on y-axis | accuracy |

#### Example Output

```
Data Quality Metrics Scatter Plot
===================================

1.0 │                                        
    │●                                      
    │            ●                          
    │                                        
0.5 │      ●           ●                    
    │                      ●                
    │    ●      ●                          
0.0 └────────────────────────────────────>
      0.0        0.5                 1.0  completeness
```

## Best Practices

### When to Use Each Visualization

- **Heatmap**: Use when you want to see how a specific metric varies across different dimensions, such as time periods and data partitions.

- **Bar Chart**: Use when you want to compare multiple metrics side by side to get a quick overview of your data quality.

- **Scatter Plot**: Use when you want to identify correlations between two metrics, such as whether completeness and accuracy are related.

### Tips for Effective Visualization

1. **Start with the Bar Chart**: The bar chart provides a good overview of all metrics, making it a good starting point.

2. **Drill Down with Heatmap**: Once you've identified a metric of interest, use the heatmap to see how it varies across different dimensions.

3. **Look for Correlations**: Use the scatter plot to identify relationships between metrics that might indicate underlying issues.

4. **Use Color Coding**: Pay attention to the color coding in all visualizations to quickly identify areas of concern.

## Troubleshooting

### Common Issues

- **No data displayed**: Ensure that the table path is correct and that the table contains data.

- **Incorrect metrics**: Check that the metric names are spelled correctly and are available for your table.

- **Visualization too small/large**: Adjust the width and height parameters for the heatmap to fit your terminal.

### Error Messages

- **"No table path specified"**: You must provide a valid table path using the `--table` flag.

- **"No metric specified"**: For heatmap visualizations, you must specify a metric using the `--metric` flag.

- **"Both x-metric and y-metric must be specified"**: For scatter plot visualizations, you must specify both x and y metrics.

## Advanced Usage

### Combining Visualizations

For a comprehensive analysis, you can combine multiple visualizations. For example:

```bash
# Get an overview of all metrics
nessi visualize barchart --table s3://my-bucket/my-table --all-metrics

# Drill down into completeness with a heatmap
nessi visualize heatmap --table s3://my-bucket/my-table --metric completeness

# Check for correlation between completeness and accuracy
nessi visualize scatter --table s3://my-bucket/my-table --x-metric completeness --y-metric accuracy
```

### Exporting Visualizations

While the visualizations are designed for terminal display, you can capture them in a file using standard shell redirection:

```bash
nessi visualize barchart --table s3://my-bucket/my-table > barchart.txt
```

## Future Enhancements

The Nessi team is working on additional visualization types and features, including:

- Time series visualizations for tracking metrics over time
- Pie charts for showing proportional distribution of quality issues
- Interactive visualizations with drill-down capabilities
- Export to HTML and PDF formats

Stay tuned for these exciting enhancements in future releases!
