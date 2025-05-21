package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/nessi-dev/nessi/pkg/plugins"
)

// All plugin-related types and interfaces are now imported from the plugins package.
// Remove local definitions to avoid duplication and import cycles.

var pluginsCmd *cobra.Command
var pluginManager *plugins.PluginManager

// newPluginManager creates a new plugin manager
func newPluginManager(pluginsDir string) *plugins.PluginManager {
	return plugins.NewPluginManager(pluginsDir)
}

func initPluginsCmd() {
	pluginManager = newPluginManager(getPluginsDir())

	pluginsCmd = &cobra.Command{
		Use:   "plugins",
		Short: "Manage Nessi plugins",
		Long:  "Manage Nessi plugins - list, install, and execute plugin commands",
		Run: func(cmd *cobra.Command, args []string) {
			// Default to list if no subcommand is provided
			listPlugins(cmd, args)
		},
	}

	// List plugins subcommand
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		Run:   listPlugins,
	}

	// Info subcommand
	infoCmd := &cobra.Command{
		Use:   "info [plugin_name]",
		Short: "Show information about a plugin",
		Args:  cobra.ExactArgs(1),
		Run:   showPluginInfo,
	}

	// Refresh subcommand
	refreshCmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh the plugin list",
		Run:   refreshPlugins,
	}

	// Install subcommand
	installCmd := &cobra.Command{
		Use:   "install [plugin_path]",
		Short: "Install a plugin from a local file",
		Args:  cobra.ExactArgs(1),
		Run:   installPlugin,
	}

	// Debug subcommand
	debugCmd := &cobra.Command{
		Use:   "debug [plugin_path]",
		Short: "Debug a plugin",
		Args:  cobra.ExactArgs(1),
		Run:   debugPlugin,
	}

	// Add subcommands
	pluginsCmd.AddCommand(listCmd)
	pluginsCmd.AddCommand(infoCmd)
	pluginsCmd.AddCommand(refreshCmd)
	pluginsCmd.AddCommand(installCmd)
	pluginsCmd.AddCommand(debugCmd)

	// Add flags
	listCmd.Flags().BoolP("verbose", "v", false, "Show detailed plugin information")
	debugCmd.Flags().BoolP("verbose", "v", false, "Show detailed debug information")

	// The plugins command will be added to the root command in the Init() function

	// Initialize plugin manager
	err := pluginManager.Initialize()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize plugin manager: %v\n", err)
	}

	// Add dynamic plugin commands
	addPluginCommands()
}

// getPluginsDir returns the plugins directory path
func getPluginsDir() string {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Warning: Failed to get user home directory: %v\n", err)
		return "./plugins"
	}

	// Create the plugins directory path
	pluginsDir := filepath.Join(homeDir, ".nessi", "plugins")

	// Create the directory if it doesn't exist
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		err := os.MkdirAll(pluginsDir, 0755)
		if err != nil {
			fmt.Printf("Warning: Failed to create plugins directory: %v\n", err)
		}
	}

	return pluginsDir
}

// addPluginCommands adds dynamic commands for each plugin
func addPluginCommands() {
	for _, p := range pluginManager.ListPlugins() {
		// Create a command for the plugin
		pluginCmd := &cobra.Command{
			Use:   p.Name,
			Short: p.Description,
			Run: func(cmd *cobra.Command, args []string) {
				// If no subcommand is provided, show plugin info
				if len(args) == 0 {
					showPluginInfo(cmd, []string{cmd.Use})
					return
				}

				// Execute the plugin command
				executePluginCommand(cmd.Use, args[0], args[1:])
			},
		}

		// Add subcommands for each plugin command
		for _, cmd := range p.Commands {
			commandName := cmd.Name
			commandDesc := cmd.Description

			pluginCmd.AddCommand(&cobra.Command{
				Use:   commandName,
				Short: commandDesc,
				Run: func(c *cobra.Command, args []string) {
					executePluginCommand(p.Name, commandName, args)
				},
			})
		}

		// Add the plugin command to the plugins command
		pluginsCmd.AddCommand(pluginCmd)
	}
}

// listPlugins lists all installed plugins
func listPlugins(cmd *cobra.Command, args []string) {
	// Get verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Initialize the plugin manager if not already initialized
	if !pluginManager.Initialized {
		err := pluginManager.Initialize()
		if err != nil {
			fmt.Printf("Error: Failed to initialize plugin manager: %v\n", err)
			return
		}
	}

	// Get the list of plugins
	plugins := pluginManager.ListPlugins()

	// Print the list of plugins
	if len(plugins) == 0 {
		fmt.Println("No plugins installed.")
		fmt.Println("\nTo install a plugin, run:")
		fmt.Println("  nessi plugins install <plugin_path>")
		return
	}

	// Create a new tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Print header
	header := color.New(color.FgCyan, color.Bold)
	if verbose {
		header.Fprintln(w, "NAME\tVERSION\tAUTHOR\tDESCRIPTION\tCOMMANDS")
	} else {
		header.Fprintln(w, "NAME\tVERSION\tDESCRIPTION")
	}

	// Print plugins
	for _, p := range plugins {
		if verbose {
			// Get commands as a string
			commands := make([]string, len(p.Commands))
			for i, cmd := range p.Commands {
				commands[i] = cmd.Name
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", p.Name, p.Version, p.Author, p.Description, strings.Join(commands, ", "))
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\n", p.Name, p.Version, p.Description)
		}
	}

	// Flush the tabwriter
	w.Flush()

	// Print help text
	fmt.Println("\nTo get more information about a plugin, run:")
	fmt.Println("  nessi plugins info <plugin_name>")
}

// showPluginInfo shows information about a plugin
func showPluginInfo(cmd *cobra.Command, args []string) {
	// Get the plugin name
	pluginName := args[0]

	// Get the plugin
	plugin, exists := pluginManager.GetPlugin(pluginName)
	if !exists {
		fmt.Printf("Error: Plugin '%s' not found\n", pluginName)
		return
	}

	// Print plugin information
	header := color.New(color.FgCyan, color.Bold)
	header.Println("Plugin Information:")
	fmt.Printf("Name: %s\n", plugin.Name)
	fmt.Printf("Version: %s\n", plugin.Version)
	fmt.Printf("Author: %s\n", plugin.Author)
	fmt.Printf("Description: %s\n", plugin.Description)
	fmt.Printf("Path: %s\n", plugin.Path)

	// Print commands
	if len(plugin.Commands) > 0 {
		header.Println("\nCommands:")
		for _, cmd := range plugin.Commands {
			fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
			fmt.Printf("    Usage: %s\n", cmd.Usage)
		}
	}

	// Print hooks
	if len(plugin.Hooks) > 0 {
		header.Println("\nHooks:")
		for hook, method := range plugin.Hooks {
			fmt.Printf("  %s: %s\n", hook, method)
		}
	}
}

// refreshPlugins refreshes the plugin list
func refreshPlugins(cmd *cobra.Command, args []string) {
	// Re-initialize the plugin manager
	err := pluginManager.Initialize()
	if err != nil {
		fmt.Printf("Error: Failed to refresh plugins: %v\n", err)
		return
	}

	// Re-add plugin commands
	addPluginCommands()

	fmt.Println("Plugins refreshed successfully.")

	// List plugins
	listPlugins(cmd, args)
}

// installPlugin installs a plugin from a local file
func installPlugin(cmd *cobra.Command, args []string) {
	// Get the plugin path
	pluginPath := args[0]

	// Check if the plugin file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		fmt.Printf("Error: Plugin file '%s' not found\n", pluginPath)
		return
	}

	// Check if the plugin is a shared library
	if !strings.HasSuffix(pluginPath, ".so") && !strings.HasSuffix(pluginPath, ".dll") {
		fmt.Printf("Error: Plugin file '%s' is not a shared library (.so or .dll)\n", pluginPath)
		return
	}

	// Check if the metadata file exists
	metadataPath := strings.TrimSuffix(pluginPath, filepath.Ext(pluginPath)) + ".json"
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		fmt.Printf("Error: Plugin metadata file '%s' not found\n", metadataPath)
		return
	}

	// Get the plugins directory
	pluginsDir := getPluginsDir()

	// Get the plugin filename
	pluginFilename := filepath.Base(pluginPath)
	metadataFilename := filepath.Base(metadataPath)

	// Copy the plugin file to the plugins directory
	destPath := filepath.Join(pluginsDir, pluginFilename)
	srcFile, err := os.ReadFile(pluginPath)
	if err != nil {
		fmt.Printf("Error: Failed to read plugin file: %v\n", err)
		return
	}

	err = os.WriteFile(destPath, srcFile, 0644)
	if err != nil {
		fmt.Printf("Error: Failed to write plugin file: %v\n", err)
		return
	}

	// Copy the metadata file to the plugins directory
	destMetadataPath := filepath.Join(pluginsDir, metadataFilename)
	srcMetadataFile, err := os.ReadFile(metadataPath)
	if err != nil {
		fmt.Printf("Error: Failed to read metadata file: %v\n", err)
		return
	}

	err = os.WriteFile(destMetadataPath, srcMetadataFile, 0644)
	if err != nil {
		fmt.Printf("Error: Failed to write metadata file: %v\n", err)
		return
	}

	fmt.Printf("Plugin '%s' installed successfully.\n", pluginFilename)

	// Refresh plugins
	refreshPlugins(cmd, args)
}

// debugPlugin debugs a plugin
func debugPlugin(cmd *cobra.Command, args []string) {
	// Get the plugin path
	pluginPath := args[0]

	// Get verbose flag
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Check if the plugin file exists
	if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
		fmt.Printf("Error: Plugin file '%s' not found\n", pluginPath)
		return
	}

	// Check if the plugin is a shared library
	if !strings.HasSuffix(pluginPath, ".so") && !strings.HasSuffix(pluginPath, ".dll") {
		fmt.Printf("Error: Plugin file '%s' is not a shared library (.so or .dll)\n", pluginPath)
		return
	}

	// Check if the metadata file exists
	metadataPath := strings.TrimSuffix(pluginPath, filepath.Ext(pluginPath)) + ".json"
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		fmt.Printf("Error: Plugin metadata file '%s' not found\n", metadataPath)
		return
	}

	// Create a temporary plugin manager for debugging
	tempManager := newPluginManager(filepath.Dir(pluginPath))
	tempManager.SetDebugMode(true)
	tempManager.SetVerbose(verbose)

	// Initialize the temporary plugin manager
	err := tempManager.Initialize()
	if err != nil {
		fmt.Printf("Error: Failed to initialize plugin manager: %v\n", err)
		return
	}

	// Get the plugin name from the metadata file
	pluginName := filepath.Base(pluginPath)
	pluginName = strings.TrimSuffix(pluginName, filepath.Ext(pluginName))

	// Get the plugin
	plugin, exists := tempManager.GetPlugin(pluginName)
	if !exists {
		fmt.Printf("Error: Failed to load plugin '%s'\n", pluginName)
		return
	}

	// Print plugin information
	header := color.New(color.FgGreen, color.Bold)
	header.Println("Plugin loaded successfully!")
	fmt.Printf("Name: %s\n", plugin.Name)
	fmt.Printf("Version: %s\n", plugin.Version)
	fmt.Printf("Author: %s\n", plugin.Author)
	fmt.Printf("Description: %s\n", plugin.Description)

	// Print commands
	if len(plugin.Commands) > 0 {
		header.Println("\nCommands:")
		for _, cmd := range plugin.Commands {
			fmt.Printf("  %s: %s\n", cmd.Name, cmd.Description)
			fmt.Printf("    Usage: %s\n", cmd.Usage)
		}
	}

	// Print hooks
	if len(plugin.Hooks) > 0 {
		header.Println("\nHooks:")
		for hook, method := range plugin.Hooks {
			fmt.Printf("  %s: %s\n", hook, method)
		}
	}

	// Print debug information
	header.Println("\nDebug Information:")
	fmt.Printf("Plugin Path: %s\n", plugin.Path)
	fmt.Printf("Metadata Path: %s\n", metadataPath)

	// Print success message
	header.Println("\nPlugin is valid and can be installed.")
	fmt.Println("To install this plugin, run:")
	fmt.Printf("  nessi plugins install %s\n", pluginPath)
}

// executePluginCommand executes a plugin command
func executePluginCommand(pluginName, commandName string, args []string) {
	// Get the plugin
	plugin, exists := pluginManager.GetPlugin(pluginName)
	if !exists {
		fmt.Printf("Error: Plugin '%s' not found\n", pluginName)
		return
	}

	// Execute the command
	err := pluginManager.ExecuteCommand(plugin, commandName, args)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
