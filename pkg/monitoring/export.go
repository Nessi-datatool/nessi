package monitoring

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
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
	if !m.config.Retention.Enabled || m.metricRetention == nil {
		return "", fmt.Errorf("metric retention is not enabled")
	}
	
	return m.metricRetention.ExportMetrics(options)
}
