package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nessi-dev/nessi/pkg/monitoring"
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Execute     func(args []string) error
}

func main() {
	// Define commands
	commands := []Command{
		{
			Name:        "metrics",
			Description: "Show metrics",
			Execute:     showMetrics,
		},
		{
			Name:        "alerts",
			Description: "Show alerts",
			Execute:     showAlerts,
		},
		{
			Name:        "export",
			Description: "Export metrics to CSV or JSON",
			Execute:     exportMetrics,
		},
		{
			Name:        "login",
			Description: "Authenticate and get an access token",
			Execute:     login,
		},
		{
			Name:        "user",
			Description: "User management (create, list, update, delete)",
			Execute:     manageUsers,
		},
		{
			Name:        "apikey",
			Description: "Generate or regenerate API key",
			Execute:     manageAPIKey,
		},
		{
			Name:        "help",
			Description: "Show help",
			Execute:     showHelp,
		},
	}

	// Parse command
	if len(os.Args) < 2 {
		fmt.Println("Error: No command specified")
		showHelp(nil)
		os.Exit(1)
	}

	// Find and execute command
	commandName := os.Args[1]
	for _, cmd := range commands {
		if cmd.Name == commandName {
			if err := cmd.Execute(os.Args[2:]); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			return
		}
	}

	// Command not found
	fmt.Printf("Error: Unknown command '%s'\n", commandName)
	showHelp(nil)
	os.Exit(1)
}

// showHelp shows the help message
func showHelp(args []string) error {
	fmt.Println("Nessi Monitoring CLI")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  nessi-monitor <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  metrics     Show metrics")
	fmt.Println("  alerts      Show alerts")
	fmt.Println("  export      Export metrics to CSV or JSON")
	fmt.Println("  login       Authenticate and get an access token")
	fmt.Println("  user        User management (create, list, update, delete)")
	fmt.Println("  apikey      Generate or regenerate API key")
	fmt.Println("  help        Show help")
	fmt.Println()
	fmt.Println("For command-specific help, run:")
	fmt.Println("  nessi-monitor <command> --help")
	return nil
}

// showMetrics shows metrics
func showMetrics(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("metrics", flag.ExitOnError)
	metricName := fs.String("metric", "", "Metric name to show")
	serverURL := fs.String("server", "http://localhost:9090", "Metrics server URL")
	format := fs.String("format", "table", "Output format (table, json)")
	help := fs.Bool("help", false, "Show help")
	fs.Parse(args)

	// Show help if requested
	if *help {
		fmt.Println("Usage: nessi-monitor metrics [options]")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		return nil
	}

	// Validate metric name
	if *metricName == "" {
		return fmt.Errorf("metric name is required")
	}

	// Build URL
	url := fmt.Sprintf("%s/metrics/history?metric=%s", *serverURL, *metricName)

	// Make request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to metrics server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Parse response
	var metrics []monitoring.MetricValue
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return fmt.Errorf("failed to parse metrics: %v", err)
	}

	// Display metrics
	switch *format {
	case "table":
		return displayMetricsTable(metrics, *metricName)
	case "json":
		return displayMetricsJSON(metrics)
	default:
		return fmt.Errorf("unsupported format: %s", *format)
	}
}

// displayMetricsTable displays metrics in a table format
func displayMetricsTable(metrics []monitoring.MetricValue, metricName string) error {
	if len(metrics) == 0 {
		fmt.Printf("No metrics found for '%s'\n", metricName)
		return nil
	}

	// Create tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintf(w, "Timestamp\tValue\tLabels\n")
	fmt.Fprintf(w, "---------\t-----\t------\n")

	// Print metrics (most recent first)
	for i := len(metrics) - 1; i >= 0; i-- {
		metric := metrics[i]
		
		// Format timestamp
		timestamp := metric.Timestamp.Format("2006-01-02 15:04:05")
		
		// Format labels
		var labels []string
		for k, v := range metric.Labels {
			labels = append(labels, fmt.Sprintf("%s=%s", k, v))
		}
		labelsStr := strings.Join(labels, ", ")
		
		fmt.Fprintf(w, "%s\t%.2f\t%s\n", timestamp, metric.Value, labelsStr)
	}

	return nil
}

// displayMetricsJSON displays metrics in JSON format
func displayMetricsJSON(metrics []monitoring.MetricValue) error {
	data, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %v", err)
	}
	
	fmt.Println(string(data))
	return nil
}

// showAlerts shows alerts
func showAlerts(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("alerts", flag.ExitOnError)
	serverURL := fs.String("server", "http://localhost:9090", "Metrics server URL")
	format := fs.String("format", "table", "Output format (table, json)")
	severity := fs.String("severity", "", "Filter by severity (critical, warning, ok)")
	help := fs.Bool("help", false, "Show help")
	fs.Parse(args)

	// Show help if requested
	if *help {
		fmt.Println("Usage: nessi-monitor alerts [options]")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		return nil
	}

	// Build URL
	url := fmt.Sprintf("%s/alerts", *serverURL)

	// Make request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to metrics server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Parse response
	var alerts []monitoring.Alert
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		return fmt.Errorf("failed to parse alerts: %v", err)
	}

	// Filter by severity if specified
	if *severity != "" {
		var filtered []monitoring.Alert
		for _, alert := range alerts {
			if strings.EqualFold(alert.Severity, *severity) {
				filtered = append(filtered, alert)
			}
		}
		alerts = filtered
	}

	// Display alerts
	switch *format {
	case "table":
		return displayAlertsTable(alerts)
	case "json":
		return displayAlertsJSON(alerts)
	default:
		return fmt.Errorf("unsupported format: %s", *format)
	}
}

// displayAlertsTable displays alerts in a table format
func displayAlertsTable(alerts []monitoring.Alert) error {
	if len(alerts) == 0 {
		fmt.Println("No alerts found")
		return nil
	}

	// Create tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintf(w, "Timestamp\tName\tSeverity\tMessage\tMetadata\n")
	fmt.Fprintf(w, "---------\t----\t--------\t-------\t--------\n")

	// Print alerts (most recent first)
	for i := len(alerts) - 1; i >= 0; i-- {
		alert := alerts[i]
		
		// Format timestamp
		timestamp := alert.Timestamp.Format("2006-01-02 15:04:05")
		
		// Format metadata
		var metadata []string
		for k, v := range alert.Metadata {
			metadata = append(metadata, fmt.Sprintf("%s=%s", k, v))
		}
		metadataStr := strings.Join(metadata, ", ")
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", timestamp, alert.Name, alert.Severity, alert.Message, metadataStr)
	}

	return nil
}

// displayAlertsJSON displays alerts in JSON format
func displayAlertsJSON(alerts []monitoring.Alert) error {
	data, err := json.MarshalIndent(alerts, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal alerts: %v", err)
	}
	
	fmt.Println(string(data))
	return nil
}

// exportMetrics exports metrics to a file
func exportMetrics(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	metricName := fs.String("metric", "", "Metric name to export")
	serverURL := fs.String("server", "http://localhost:9090", "Metrics server URL")
	format := fs.String("format", "csv", "Export format (csv, json)")
	output := fs.String("output", "", "Output file path")
	startTime := fs.String("start", "", "Start time (RFC3339 format)")
	endTime := fs.String("end", "", "End time (RFC3339 format)")
	help := fs.Bool("help", false, "Show help")
	fs.Parse(args)

	// Show help if requested
	if *help {
		fmt.Println("Usage: nessi-monitor export [options]")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		return nil
	}

	// Validate metric name
	if *metricName == "" {
		return fmt.Errorf("metric name is required")
	}

	// Validate output file
	if *output == "" {
		return fmt.Errorf("output file path is required")
	}

	// Parse time range
	start := time.Now().Add(-24 * time.Hour) // Default to last 24 hours
	end := time.Now()
	
	if *startTime != "" {
		parsedStart, err := time.Parse(time.RFC3339, *startTime)
		if err != nil {
			return fmt.Errorf("invalid start time: %v", err)
		}
		start = parsedStart
	}
	
	if *endTime != "" {
		parsedEnd, err := time.Parse(time.RFC3339, *endTime)
		if err != nil {
			return fmt.Errorf("invalid end time: %v", err)
		}
		end = parsedEnd
	}

	// Build URL
	url := fmt.Sprintf("%s/api/export?metric=%s&format=%s&start=%s&end=%s", 
		*serverURL, *metricName, *format, start.Format(time.RFC3339), end.Format(time.RFC3339))

	// Make request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to metrics server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Create output file
	file, err := os.Create(*output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Copy response to file
	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write output file: %v", err)
	}

	fmt.Printf("Exported metrics to %s\n", *output)
	return nil
}
