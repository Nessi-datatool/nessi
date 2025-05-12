package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring/dashboard"
)

func main() {
	// Parse command-line flags
	metricsPort := flag.Int("metrics-port", 9090, "Port for metrics server")
	dashboardPort := flag.Int("dashboard-port", 8080, "Port for dashboard server")
	configPath := flag.String("config", "pkg/monitoring/config/monitoring.json", "Path to monitoring configuration")
	enableDashboard := flag.Bool("dashboard", true, "Enable web dashboard")
	flag.Parse()

	// Initialize monitor
	monitor := monitoring.New()
	monitor.SetConfigPath(*configPath)

	// Load configuration
	if err := monitor.LoadConfig(); err != nil {
		logging.Error("Failed to load monitoring configuration", err)
		os.Exit(1)
	}

	// Override metrics port if specified
	if *metricsPort != 9090 {
		monitor.SetMetricsPort(*metricsPort)
	}

	// Start monitoring service
	monitor.Start()
	logging.Info("Monitoring service started")

	// Start dashboard if enabled
	if *enableDashboard {
		dashboardOpts := dashboard.DashboardOptions{
			ListenAddr: fmt.Sprintf(":%d", *dashboardPort),
		}

		dash, err := dashboard.New(monitor, dashboardOpts)
		if err != nil {
			logging.Error("Failed to create dashboard", err)
		} else {
			go func() {
				if err := dash.Start(); err != nil {
					logging.Error("Dashboard server error", err)
				}
			}()
			logging.Info(fmt.Sprintf("Dashboard available at http://localhost:%d", *dashboardPort))
		}
	}

	// Generate some sample metrics for demonstration
	go generateSampleMetrics(monitor)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logging.Info("Shutting down monitoring service")
}

// generateSampleMetrics generates sample metrics for demonstration
func generateSampleMetrics(monitor *monitoring.Monitor) {
	tables := []string{"users", "orders", "products", "transactions"}
	operations := []string{"read", "write", "update", "delete"}
	rules := []string{"completeness", "uniqueness", "range", "format"}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Record table metrics
			for _, table := range tables {
				size := 100.0 + float64(time.Now().Unix()%1000)
				count := 1000 + int(time.Now().Unix()%100)
				monitor.RecordTableMetrics(table, []map[string]interface{}{
					{"size": size, "count": count},
				})
			}

			// Record latency metrics
			for _, op := range operations {
				latency := time.Duration(50+time.Now().Unix()%200) * time.Millisecond
				monitor.RecordLatency(op, latency)
			}

			// Record some errors and rule violations
			if time.Now().Unix()%30 == 0 {
				table := tables[time.Now().Unix()%int64(len(tables))]
				monitor.RecordError(table, "validation_error")
			}

			if time.Now().Unix()%20 == 0 {
				table := tables[time.Now().Unix()%int64(len(tables))]
				rule := rules[time.Now().Unix()%int64(len(rules))]
				monitor.RecordRuleViolation(table, rule)
			}

			// Generate some alerts
			if time.Now().Unix()%60 == 0 {
				table := tables[time.Now().Unix()%int64(len(tables))]
				value := 1000.0 + float64(time.Now().Unix()%500)
				monitor.SendAlert("high_latency", table, value)
			}
		}
	}
}
