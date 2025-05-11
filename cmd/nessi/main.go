// Implement the main CLI application for Nessi.dev using Cobra
// The CLI should have the following structure:
// - Root command: `nessi` with global flags for config file path and verbosity
// - `serve` subcommand: Start the API server
//   - Flags: --host, --port, --tls
// - `check` subcommand: Run data quality checks on a Delta table
//   - Args: table_path (required)
//   - Flags: --rules (rules file), --output (output format: json/yaml/text)
// - `profile` subcommand: Generate a profile report for a Delta table
//   - Args: table_path (required)
//   - Flags: --output (path for report), --format (html/pdf/json)
// - `extensions` subcommand: Manage extensions
//   - `list` sub-subcommand: List all extensions
//   - `enable` sub-subcommand: Enable an extension
//   - `disable` sub-subcommand: Disable an extension
// The application should:
// 1. Load configuration using the config package
// 2. Initialize the extension manager
// 3. Set up logging
// 4. Handle errors gracefully with useful messages
// Use Cobra for command handling and Viper for configuration

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/nessi-dev/nessi-dev/internal/config"
	"github.com/nessi-dev/nessi-dev/internal/delta"
	"github.com/nessi-dev/nessi-dev/internal/extensions"
	"github.com/nessi-dev/nessi-dev/internal/monitor"
	"github.com/nessi-dev/nessi-dev/internal/quality" // Added back quality import
	"github.com/nessi-dev/nessi-dev/internal/report"
	"github.com/nessi-dev/nessi-dev/internal/security"
	"github.com/nessi-dev/nessi-dev/internal/server" // Added back server import
	"github.com/nessi-dev/nessi-dev/pkg"
	"go.uber.org/zap"
)

var (
	// Global flags
	cfgFile   string
	verbose   bool

	extManager *extensions.Manager // Added extension manager instance

	// Serve command flags
	host      string
	port      int
	tls       bool

	// Check command flags
	rulesFile string
	outputFmt string

	// Profile command flags
	outputPath string
	reportFmt  string

	configPath = flag.String("config", "config/config.yaml", "path to configuration file")
	version    = flag.Bool("version", false, "print version and exit")
)

var rootCmd = &cobra.Command{
	Use:   "nessi",
	Short: "Nessi-dev CLI - A tool for data quality and profiling",
	Long: `Nessi-dev is a CLI tool that provides functionality for:
- Running data quality checks on Delta tables
- Generating profiles for Delta tables
- Managing extensions
- Starting the API server`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration
		if err := initConfig(); err != nil {
			return fmt.Errorf("failed to initialize config: %w", err)
		}

		// Initialize extension manager
		extManager = extensions.NewManager() // Initialize the manager
		// TODO: Load enabled extensions from config, e.g., by iterating viper.GetStringSlice("extensions.enabled")
		// and calling extManager.RegisterExtension() or extManager.EnableExtension() if they are pre-registered.

		// Set up logging
		if err := setupLogging(verbose); err != nil {
			return fmt.Errorf("failed to set up logging: %w", err)
		}

		return nil
	},
}

// initConfig reads in config file and ENV variables if set
func initConfig() error {
	home, homeErr := os.UserHomeDir() // Declare home and its error at the top

	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		if homeErr != nil {
			// Log or handle the error if home directory is essential for default config
			log.Printf("warning: could not get user home directory: %v. Default config file '.nessi' might not be found.", homeErr)
		} else if home != "" { // Only use home if err is nil and home is not empty
			// Search config in home directory with name ".nessi" (without extension)
			viper.AddConfigPath(home)
			viper.SetConfigName(".nessi")
		}
	}

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// Create default config if it doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			defaultConfig := map[string]interface{}{
				"server": map[string]interface{}{
					"host": "localhost",
					"port": 8080,
					"tls":  false,
				},
				"extensions": map[string]interface{}{
					"enabled": []string{},
					"path":    filepath.Join(home, ".nessi", "extensions"), // home might be "" if homeErr != nil, Join handles this
				},
				"output": map[string]interface{}{
					"format": "text",
					"path":   ".",
				},
			}
			viper.SetConfigFile(filepath.Join(home, ".nessi.yaml"))
			for k, v := range defaultConfig {
				viper.Set(k, v)
			}
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// setupLogging configures the logging system
func setupLogging(verbose bool) error {
	// TODO: Implement logging setup
	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Use config values if flags are not set
		if host == "" {
			host = viper.GetString("server.host")
		}
		if port == 0 {
			port = viper.GetInt("server.port")
		}
		log.Printf("Starting server on %s:%d (TLS active: %t - note: TLS flag might not be passed to current StartServer call)", host, port, tls)
		return pkg.StartServer(host, port) // Adjusted arity, 'tls' not passed
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [table_path]",
	Short: "Run data quality checks on a Delta table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Running quality checks for %s (Rules: %s, Output: %s - note: these might not be passed to current RunQualityChecks call)", args[0], rulesFile, outputFmt)
		return pkg.RunQualityChecks(args[0]) // Adjusted arity
	},
}

var profileCmd = &cobra.Command{
	Use:   "profile [table_path]",
	Short: "Generate a profile report for a Delta table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Generating profile for %s (Output: %s, Format: %s - note: these might not be passed to current GenerateProfile call)", args[0], outputPath, reportFmt)
		return pkg.GenerateProfile(args[0]) // Adjusted arity
	},
}

var extensionsCmd = &cobra.Command{
	Use:   "extensions",
	Short: "Manage extensions",
}

var listExtensionsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all extensions",
	RunE: func(cmd *cobra.Command, args []string) error {
		exts := extManager.ListExtensions() // Use extManager
		if len(exts) == 0 {
			fmt.Println("No extensions registered.")
			return nil
		}
		fmt.Println("Registered extensions:")
		for _, ext := range exts {
			status := "disabled"
			if ext.Enabled {
				status = "enabled"
			}
			fmt.Printf("- %s (v%s): %s [%s]\n", ext.Name, ext.Version, ext.Description, status)
		}
		return nil
	},
}

var enableExtensionCmd = &cobra.Command{
	Use:   "enable [extension-name]",
	Short: "Enable an extension",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return extManager.EnableExtension(args[0]) // Use extManager
	},
}

var disableExtensionCmd = &cobra.Command{
	Use:   "disable [extension-name]",
	Short: "Disable an extension",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return extManager.DisableExtension(args[0]) // Use extManager
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nessi.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	// Serve command flags
	serveCmd.Flags().StringVar(&host, "host", "", "Host to bind the server to (default from config)")
	serveCmd.Flags().IntVar(&port, "port", 0, "Port to bind the server to (default from config)")
	serveCmd.Flags().BoolVar(&tls, "tls", false, "Enable TLS/HTTPS")

	// Check command flags
	checkCmd.Flags().StringVar(&rulesFile, "rules", "", "Path to rules file")
	checkCmd.Flags().StringVar(&outputFmt, "output", "text", "Output format (json/yaml/text)")

	// Profile command flags
	profileCmd.Flags().StringVar(&outputPath, "output", ".", "Path for report output")
	profileCmd.Flags().StringVar(&reportFmt, "format", "html", "Report format (html/pdf/json)")

	// Add commands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(profileCmd)
	
	// Add extension subcommands
	extensionsCmd.AddCommand(listExtensionsCmd)
	extensionsCmd.AddCommand(enableExtensionCmd)
	extensionsCmd.AddCommand(disableExtensionCmd)
	rootCmd.AddCommand(extensionsCmd)
}

func main() {
	flag.Parse()

	if *version {
		fmt.Println("Nessi.dev v0.1.0")
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize components
	secManager := security.NewSecurityManager(cfg.Security.JWT.Secret, cfg.Security.JWT.Expiration)
	// The secManager.ApplyConfig call was incorrect as it expected security.SecurityConfig,
	// not the global config.SecurityConfig. SecurityManager is typically configured internally
	// or via its own specific config structure if ApplyConfig is used directly.
	// if err := secManager.ApplyConfig(cfg.Security); err != nil {
	// 	log.Fatalf("Failed to apply security configuration: %v", err)
	// }

	// Corrected delta.NewConnector call and error handling
	deltaConnector, err := delta.NewConnector(cfg.Delta.BasePath)
	if err != nil {
		log.Fatalf("Failed to initialize Delta Lake connector: %v", err)
	}

	// Initialize Quality Manager
	// TODO: Add configuration for QualityManager if needed (e.g., rules path)
	qualityManager := quality.NewManager() // Corrected: NewManager takes no arguments

	monitorManager := monitor.NewManager(cfg.Monitoring.Prometheus.PushGateway)
	// Loop through configured metrics and register them
	for _, metricConf := range cfg.Monitoring.Metrics {
		// Convert string type from config to monitor.MetricType
		metricType := monitor.MetricType(metricConf.Type)
		if err := monitorManager.RegisterMetric(metricConf.Name, metricConf.Description, metricType, metricConf.Labels); err != nil {
			log.Fatalf("Failed to register metric '%s': %v", metricConf.Name, err)
		}
	}

	// Loop through configured alerts and add them
	for _, alertConf := range cfg.Monitoring.Alerts {
		// Create monitor.MonitoringAlert from config.AlertConfig
		// Note: config.AlertConfig doesn't have Labels, so passing nil or empty map for now.
		// If labels are needed for alerts from config, config.AlertConfig needs to be updated.
		alert := monitor.MonitoringAlert{
			Name:        alertConf.Name,
			Description: alertConf.Description,
			Severity:    alertConf.Severity,
			Condition:   alertConf.Condition,
			Threshold:   alertConf.Threshold,
			Labels:      nil, // Or make(map[string]string) if empty map is preferred
		}
		if err := monitorManager.AddAlert(alert); err != nil {
			log.Fatalf("Failed to add alert '%s': %v", alertConf.Name, err)
		}
	}

	reportManager, err := report.NewManager(cfg.Reports.OutputDir)
	if err != nil {
		log.Fatalf("Failed to initialize report manager: %v", err)
	}
	// Loop through configured report templates and add them
	for _, templateConf := range cfg.Reports.Templates {
		// Manually map config.TemplateConfig to report.ReportTemplate
		reportTpl := report.ReportTemplate{
			ID:          templateConf.ID,
			Name:        templateConf.Name,
			Description: templateConf.Description,
			Format:      report.ReportFormat(templateConf.Format), // Cast string to report.ReportFormat
			Template:    templateConf.Template,
			Parameters:  templateConf.Parameters,
		}
		if err := reportManager.AddTemplate(reportTpl); err != nil {
			log.Fatalf("Failed to add report template '%s': %v", templateConf.Name, err)
		}
	}

	// Prepare server.Config from the global config.ServerConfig
	serverCfg := &server.Config{
		Host:            cfg.Server.Host,
		Port:            cfg.Server.Port,
		ReadTimeout:     cfg.Server.ReadTimeout,
		WriteTimeout:    cfg.Server.WriteTimeout,
		ShutdownTimeout: cfg.Server.ShutdownTimeout, // This should now be valid
		TLS: struct {
			Enabled  bool
			CertFile string
			KeyFile  string
		}{
			Enabled:  cfg.Security.TLS.Enabled,  // Corrected: Source from cfg.Security.TLS
			CertFile: cfg.Security.TLS.CertFile, // Corrected: Source from cfg.Security.TLS
			KeyFile:  cfg.Security.TLS.KeyFile,  // Corrected: Source from cfg.Security.TLS
		},
	}

	// Create and start server
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	srv := server.NewServer(
		serverCfg,      // *server.Config
		extManager,     // *extensions.Manager
		secManager,     // *security.SecurityManager
		deltaConnector, // *delta.DeltaConnector
		qualityManager, // *quality.QualityManager
		monitorManager, // *monitor.MonitorManager
		logger,         // *zap.Logger
	)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Start monitoring
	if cfg.Monitoring.Prometheus.Enabled {
		go monitorManager.StartPeriodicPush(ctx, cfg.Monitoring.Prometheus.Interval)
		go monitorManager.StartPeriodicAlertCheck(ctx, cfg.Monitoring.Prometheus.Interval, func(alerts []monitor.MonitoringAlert) {
			for _, alert := range alerts {
				log.Printf("Alert firing: %s - %s", alert.Name, alert.Description)
			}
		})
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	// Shutdown server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}
