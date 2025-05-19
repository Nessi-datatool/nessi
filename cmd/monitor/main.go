package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
	"github.com/nessi-dev/nessi/pkg/monitoring"
)

func main() {
	// Parse command line arguments
	exportPath := flag.String("export-path", "./data/metrics", "Path to export metrics to")
	interval := flag.Duration("interval", 1*time.Minute, "Metrics collection interval")
	enableSystemMetrics := flag.Bool("system-metrics", true, "Enable system metrics collection")
	flag.Parse()

	// Initialize metrics collector
	collector := monitoring.NewReportMetricsCollector()

	// Set export path
	collector.SetMetricsExportPath(*exportPath)

	// Set collection interval
	collector.SetExportInterval(*interval)

	// Start metrics collection
	collector.Start()
	logging.Info("Metrics collection started")

	if *enableSystemMetrics {
		logging.Info("Enabling system metrics collection")
		collector.EnableSystemMetrics()
	}

	// Generate sample metrics for demonstration
	go generateSampleMetrics(collector)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	logging.Info("Shutting down metrics collector")

	// Stop monitoring service
	collector.Stop()
	logging.Info("Metrics collection stopped")
}

func generateSampleMetrics(collector *monitoring.ReportMetricsCollector) {
	for {
		// Record some sample metrics
		collector.RecordMetrics("sample_metrics", map[string]interface{}{
			"counter": 1,
			"gauge":   42.5,
		})

		// Sleep for a while
		time.Sleep(10 * time.Second)
	}
}
