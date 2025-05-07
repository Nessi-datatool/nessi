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
	"github.com/nessi-dev/nessi-dev/internal/monitor"
	"github.com/nessi-dev/nessi-dev/internal/quality"
	"github.com/nessi-dev/nessi-dev/internal/report"
	"github.com/nessi-dev/nessi-dev/internal/security"
	"github.com/nessi-dev/nessi-dev/internal/server"
	"github.com/nessi-dev/nessi-dev/pkg"
)

var (
	// Global flags
	cfgFile   string
	verbose   bool

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
		if err := pkg.InitExtensionManager(); err != nil {
			return fmt.Errorf("failed to initialize extension manager: %w", err)
		}

		// Set up logging
		if err := setupLogging(verbose); err != nil {
			return fmt.Errorf("failed to set up logging: %w", err)
		}

		return nil
	},
}

// initConfig reads in config file and ENV variables if set
func initConfig() error {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}

		// Search config in home directory with name ".nessi" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(".nessi")
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
					"path":    filepath.Join(home, ".nessi", "extensions"),
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
		return pkg.StartServer(host, port, tls)
	},
}

var checkCmd = &cobra.Command{
	Use:   "check [table-path]",
	Short: "Run data quality checks on a Delta table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkg.RunQualityChecks(args[0], rulesFile, outputFmt)
	},
}

var profileCmd = &cobra.Command{
	Use:   "profile [table-path]",
	Short: "Generate a profile report for a Delta table",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkg.GenerateProfile(args[0], outputPath, reportFmt)
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
		extensions, err := pkg.ListExtensions()
		if err != nil {
			return err
		}
		for _, ext := range extensions {
			fmt.Printf("%s (v%s): %s\n", ext.Name, ext.Version, ext.Description)
		}
		return nil
	},
}

var enableExtensionCmd = &cobra.Command{
	Use:   "enable [extension-name]",
	Short: "Enable an extension",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkg.EnableExtension(args[0])
	},
}

var disableExtensionCmd = &cobra.Command{
	Use:   "disable [extension-name]",
	Short: "Disable an extension",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkg.DisableExtension(args[0])
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
	secManager := security.NewManager(cfg.Security.JWT.Secret, cfg.Security.JWT.Expiration)
	if err := secManager.ApplyConfig(cfg.Security); err != nil {
		log.Fatalf("Failed to apply security configuration: %v", err)
	}

	deltaConnector := delta.NewConnector(cfg.Delta.BasePath, cfg.Delta.TablePath)
	if err := deltaConnector.Validate(); err != nil {
		log.Fatalf("Failed to initialize Delta Lake connector: %v", err)
	}

	qualityManager := quality.NewManager()
	if err := qualityManager.AddRules(cfg.Quality.Rules); err != nil {
		log.Fatalf("Failed to add quality rules: %v", err)
	}

	monitorManager := monitor.NewManager(cfg.Monitoring.Prometheus.PushGateway)
	if err := monitorManager.RegisterMetrics(cfg.Monitoring.Metrics); err != nil {
		log.Fatalf("Failed to register metrics: %v", err)
	}
	if err := monitorManager.AddAlerts(cfg.Monitoring.Alerts); err != nil {
		log.Fatalf("Failed to add alerts: %v", err)
	}

	reportManager := report.NewManager(cfg.Reports.OutputDir)
	if err := reportManager.AddTemplates(cfg.Reports.Templates); err != nil {
		log.Fatalf("Failed to add report templates: %v", err)
	}

	// Create and start server
	srv := server.NewServer(cfg.Server, secManager, deltaConnector, qualityManager, monitorManager, reportManager)
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
			cancel()
		}
	}()

	// Start monitoring
	if cfg.Monitoring.Prometheus.Enabled {
		go monitorManager.StartPeriodicPush(ctx, cfg.Monitoring.Prometheus.Interval)
		go monitorManager.StartPeriodicAlertCheck(ctx, cfg.Monitoring.Prometheus.Interval)
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
