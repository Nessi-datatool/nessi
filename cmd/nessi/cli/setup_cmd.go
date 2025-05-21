package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Interactive setup wizard for Nessi",
	Long: `Interactive setup wizard for Nessi.

This command guides you through the initial setup process for Nessi,
helping you configure essential settings and validate your environment.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if non-interactive mode is enabled
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")

		// Get config directory if specified
		configDir, _ := cmd.Flags().GetString("config-dir")

		if nonInteractive {
			return runNonInteractiveSetup(configDir)
		} else {
			return runSetupWizard(configDir)
		}
	},
}

func init() {
	CLI.RootCmd.AddCommand(setupCmd)

	// Add flags
	setupCmd.Flags().Bool("non-interactive", false, "Run setup in non-interactive mode with default values")
	setupCmd.Flags().String("config-dir", "", "Directory to store configuration files")
}

// runSetupWizard runs the interactive setup wizard
func runSetupWizard(configDirFlag string) error {
	// Create a reader for user input
	reader := bufio.NewReader(os.Stdin)

	// Create color printers
	header := color.New(color.FgCyan, color.Bold)
	success := color.New(color.FgGreen, color.Bold)
	info := color.New(color.FgYellow)

	// Welcome message
	header.Println("\n=== Welcome to Nessi Setup Wizard ===")
	fmt.Println("This wizard will help you set up Nessi for optimal use.")
	fmt.Println("Press Ctrl+C at any time to exit.")

	// Step 1: Configuration directory
	header.Println("\n=== Step 1: Configuration Directory ===")

	// Get default config directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return common.NewError(common.ErrInternalError, fmt.Sprintf("Error getting home directory: %s", err))
	}

	defaultConfigDir := filepath.Join(homeDir, ".nessi")
	fmt.Printf("Default configuration directory: %s\n", defaultConfigDir)

	// Ask if user wants to use default or custom directory
	fmt.Print("Use default configuration directory? [Y/n]: ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	configDir := defaultConfigDir
	if response == "n" || response == "no" {
		fmt.Print("Enter custom configuration directory: ")
		configDir, _ = reader.ReadString('\n')
		configDir = strings.TrimSpace(configDir)

		// Validate the directory
		validator := common.NewPathValidator()
		validator.AllowNonExistent = true
		configDir, err = validator.ValidatePath(configDir)
		if err != nil {
			fmt.Println("Error validating path:", err)
			return common.NewError(common.ErrInvalidPath, fmt.Sprintf("Error validating path: %s", err))
		}
	}

	// Create config directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		fmt.Printf("Creating configuration directory: %s\n", configDir)
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println("Error creating directory:", err)
			return common.NewError(common.ErrInvalidPath, fmt.Sprintf("Failed to create directory: %s", err))
		}
	}

	success.Println("✓ Configuration directory set up successfully!")

	// Step 2: Error handling configuration
	header.Println("\n=== Step 2: Error Handling Configuration ===")

	// Ask about interactive error resolution
	fmt.Println("Nessi can automatically attempt to resolve certain errors interactively.")
	fmt.Print("Enable interactive error resolution? [Y/n]: ")
	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	interactiveResolution := true
	if response == "n" || response == "no" {
		interactiveResolution = false
	}

	// Ask about error telemetry
	fmt.Println("\nNessi can collect anonymous error telemetry to help improve the tool.")
	fmt.Println("This data is stored locally and never sent to any server without your permission.")
	fmt.Print("Enable error telemetry? [Y/n]: ")
	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	errorTelemetry := true
	if response == "n" || response == "no" {
		errorTelemetry = false
	}

	// Ask about automatic retries
	fmt.Println("\nNessi can automatically retry certain operations that fail due to transient issues.")
	fmt.Print("Enable automatic retries? [Y/n]: ")
	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	autoRetry := true
	if response == "n" || response == "no" {
		autoRetry = false
	}

	// Step 3: Create configuration file
	header.Println("\n=== Step 3: Creating Configuration File ===")

	// Create config file path
	configFile := filepath.Join(configDir, "config.yaml")
	fmt.Printf("Creating configuration file: %s\n", configFile)

	// Create config content
	configContent := fmt.Sprintf(`# Nessi Configuration
# Generated by setup wizard

# Error handling configuration
error_handling:
  # Whether to enable interactive error resolution
  interactive_resolution: %t

  # Whether to enable error telemetry
  telemetry:
    enabled: %t
    # Maximum number of errors to store
    max_errors: 100

  # Automatic retry configuration
  retry:
    enabled: %t
    # Maximum number of retry attempts
    max_attempts: 3
    # Initial backoff duration in seconds
    initial_backoff: 1
    # Maximum backoff duration in seconds
    max_backoff: 30
    # Backoff multiplier
    backoff_factor: 2.0
    # Whether to add jitter to backoff
    jitter: true

# Logging configuration
logging:
  # Log level (debug, info, warn, error)
  level: "info"
  # Whether to use colored output
  colored: true
  # Log file path (leave empty to log to stderr only)
  file: ""
`, interactiveResolution, errorTelemetry, autoRetry)

	// Write config file
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		fmt.Println("Error writing configuration file:", err)
		return err
	}

	success.Println("✓ Configuration file created successfully!")

	// Step 4: Environment check
	header.Println("\n=== Step 4: Environment Check ===")

	// Check Java installation
	fmt.Println("Checking Java installation...")
	javaCmd := "java -version"
	javaOutput, err := runCommand(javaCmd)
	if err != nil {
		info.Println("⚠ Java not found or not in PATH")
		fmt.Println("Java is required for some Delta Lake operations.")
		fmt.Println("Please install Java and add it to your PATH.")
	} else {
		success.Println("✓ Java found:")
		fmt.Println(javaOutput)
	}

	// Check Python installation
	fmt.Println("\nChecking Python installation...")
	pythonCmd := "python --version || python3 --version"
	pythonOutput, err := runCommand(pythonCmd)
	if err != nil {
		info.Println("⚠ Python not found or not in PATH")
		fmt.Println("Python is required for some data processing operations.")
		fmt.Println("Please install Python and add it to your PATH.")
	} else {
		success.Println("✓ Python found:")
		fmt.Println(pythonOutput)
	}

	// Step 5: Completion
	header.Println("\n=== Setup Complete! ===")
	fmt.Println("Nessi has been configured successfully.")
	fmt.Printf("Configuration file: %s\n", configFile)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Run 'nessi help' to see available commands")
	fmt.Println("2. Check out the documentation at https://github.com/nessi-dev/nessi")
	fmt.Println("3. Try running 'nessi schema show --path <delta-table-path>' to analyze a Delta table")

	// Set the config file for viper
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: Could not read config file:", err)
	}

	return nil
}

// runCommand runs a shell command and returns the output
func runCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// runNonInteractiveSetup runs the setup wizard in non-interactive mode
func runNonInteractiveSetup(configDirFlag string) error {
	// Create color printers for output
	header := color.New(color.FgCyan, color.Bold)
	success := color.New(color.FgGreen, color.Bold)

	header.Println("\n=== Running Nessi Setup in Non-Interactive Mode ===")

	// Determine config directory
	configDir := configDirFlag
	if configDir == "" {
		// Get default config directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return common.NewError(common.ErrInternalError, fmt.Sprintf("Error getting home directory: %s", err))
		}
		configDir = filepath.Join(homeDir, ".nessi")
	}

	// Create config directory if it doesn't exist
	fmt.Printf("Creating configuration directory: %s\n", configDir)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return common.NewError(common.ErrInvalidPath, fmt.Sprintf("Failed to create directory: %s", err))
		}
	}

	// Create config file path
	configFile := filepath.Join(configDir, "config.yaml")
	fmt.Printf("Creating configuration file: %s\n", configFile)

	// Create default config content
	configContent := `# Nessi Configuration
# Generated by setup wizard (non-interactive mode)

# Error handling configuration
error_handling:
  # Whether to enable interactive error resolution
  interactive_resolution: true

  # Error telemetry configuration
  telemetry:
    # Whether to enable error telemetry
    enabled: true
    # Maximum number of errors to store
    max_errors: 100

  # Automatic retry configuration
  retry:
    enabled: true
    # Maximum number of retry attempts
    max_attempts: 3
    # Initial backoff duration in seconds
    initial_backoff: 1
    # Maximum backoff duration in seconds
    max_backoff: 30
    # Backoff multiplier
    backoff_factor: 2.0
    # Whether to add jitter to backoff
    jitter: true

# Logging configuration
logging:
  # Log level (debug, info, warn, error)
  level: "info"
  # Whether to use colored output
  colored: true
  # Log file path (leave empty to log to stderr only)
  file: ""
`

	// Write config file
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		return common.NewError(common.ErrInvalidPath, fmt.Sprintf("Failed to write configuration file: %s", err))
	}

	// Print success message
	success.Println("✓ Configuration file created successfully!")

	// Set the config file for viper
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Warning: Could not read config file:", err)
	}

	return nil
}
