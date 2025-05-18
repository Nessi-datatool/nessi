package monitoring

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
)

// ExportFormat defines the format for metric exports
type ExportFormat string

const (
	// ExportFormatCSV represents CSV export format
	ExportFormatCSV ExportFormat = "csv"
	// ExportFormatJSON represents JSON export format
	ExportFormatJSON ExportFormat = "json"
)

// ExportOptions represents options for exporting metrics
type ExportOptions struct {
	Format     ExportFormat
	OutputPath string
	StartTime  time.Time
	EndTime    time.Time
	MetricName string
	Labels     map[string]string
}

// ExportMetricsToCSV exports metrics to a CSV file
func (r *MetricRetention) ExportMetricsToCSV(options ExportOptions) (string, error) {
	// Get metric history
	metrics, err := r.GetMetricHistory(options.MetricName, options.Labels, options.StartTime, options.EndTime)
	if err != nil {
		return "", fmt.Errorf("failed to get metric history: %w", err)
	}

	if len(metrics) == 0 {
		return "", fmt.Errorf("no metrics found for the specified criteria")
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(options.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create output file
	file, err := os.Create(options.OutputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Timestamp", "Value"}
	
	// Add label columns
	var labelKeys []string
	if len(metrics) > 0 && len(metrics[0].Labels) > 0 {
		for k := range metrics[0].Labels {
			labelKeys = append(labelKeys, k)
		}
		sort.Strings(labelKeys)
		header = append(header, labelKeys...)
	}
	
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Sort metrics by timestamp
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Timestamp.Before(metrics[j].Timestamp)
	})

	// Write metrics
	for _, metric := range metrics {
		// Create row with timestamp and value
		row := []string{
			metric.Timestamp.Format(time.RFC3339),
			strconv.FormatFloat(metric.Value, 'f', -1, 64),
		}
		
		// Add label values
		for _, key := range labelKeys {
			row = append(row, metric.Labels[key])
		}
		
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	logging.Info(fmt.Sprintf("Exported %d metrics to %s", len(metrics), options.OutputPath))
	return options.OutputPath, nil
}

// ExportMetrics exports metrics in the specified format
func (r *MetricRetention) ExportMetrics(options ExportOptions) (string, error) {
	switch options.Format {
	case ExportFormatCSV:
		return r.ExportMetricsToCSV(options)
	case ExportFormatJSON:
		// TODO: Implement JSON export
		return "", fmt.Errorf("JSON export not yet implemented")
	default:
		return "", fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// ExportMetrics exports metrics in the specified format
func (m *Monitor) ExportMetrics(options ExportOptions) (string, error) {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(options.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	
	// If metric retention is enabled, use it
	if m.config != nil && m.config.Retention.Enabled && m.metricRetention != nil {
		return m.metricRetention.ExportMetrics(options)
	}
	
	// Otherwise, use the metric store directly
	switch options.Format {
	case ExportFormatCSV:
		return m.exportMetricsToCSV(options)
	case ExportFormatJSON:
		return m.exportMetricsToJSON(options)
	default:
		return "", fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// exportMetricsToCSV exports metrics to a CSV file
func (m *Monitor) exportMetricsToCSV(options ExportOptions) (string, error) {
	// Create output file
	file, err := os.Create(options.OutputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"metric_name", "value", "timestamp"}); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Get metric names
	metricNames, err := m.metricStore.GetMetricNames()
	if err != nil {
		return "", fmt.Errorf("failed to get metric names: %w", err)
	}

	// Filter by metric name if specified
	if options.MetricName != "" {
		found := false
		for _, name := range metricNames {
			if name == options.MetricName {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("metric %s not found", options.MetricName)
		}
		metricNames = []string{options.MetricName}
	}

	// Write metrics
	for _, name := range metricNames {
		value, err := m.metricStore.GetMetric(name, options.Labels)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to get metric %s: %v", name, err))
			continue
		}

		// Write row
		row := []string{
			name,
			strconv.FormatFloat(value, 'f', -1, 64),
			time.Now().Format(time.RFC3339),
		}
		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	logging.Info(fmt.Sprintf("Exported %d metrics to %s", len(metricNames), options.OutputPath))
	return options.OutputPath, nil
}

// exportMetricsToJSON exports metrics to a JSON file
func (m *Monitor) exportMetricsToJSON(options ExportOptions) (string, error) {
	// Get metric names
	metricNames, err := m.metricStore.GetMetricNames()
	if err != nil {
		return "", fmt.Errorf("failed to get metric names: %w", err)
	}

	// Filter by metric name if specified
	if options.MetricName != "" {
		found := false
		for _, name := range metricNames {
			if name == options.MetricName {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("metric %s not found", options.MetricName)
		}
		metricNames = []string{options.MetricName}
	}

	// Create metrics map
	metrics := make(map[string]interface{})
	for _, name := range metricNames {
		value, err := m.metricStore.GetMetric(name, options.Labels)
		if err != nil {
			logging.Warn(fmt.Sprintf("Failed to get metric %s: %v", name, err))
			continue
		}

		metrics[name] = map[string]interface{}{
			"value":     value,
			"timestamp": time.Now().Format(time.RFC3339),
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal metrics to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(options.OutputPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write JSON file: %w", err)
	}

	logging.Info(fmt.Sprintf("Exported %d metrics to %s", len(metricNames), options.OutputPath))
	return options.OutputPath, nil
}
