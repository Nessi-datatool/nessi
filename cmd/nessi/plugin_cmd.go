package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi/pkg/plugin"
	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage Nessi plugins",
	Long: `Manage Nessi plugins.

Plugins allow extending Nessi's functionality with custom validation rules,
alerting mechanisms, metrics collection, storage backends, export formats,
and UI components.`,
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed plugins",
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		plugins := manager.ListPlugins()
		if len(plugins) == 0 {
			fmt.Println("No plugins installed.")
			return
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(plugins, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting output: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonOutput))
			return
		}

		// Table format (default)
		fmt.Println("NAME\tVERSION\tTYPE\tENABLED\tDESCRIPTION")
		for _, p := range plugins {
			fmt.Printf("%s\t%s\t%s\t%v\t%s\n", p.Info.Name, p.Info.Version, p.Info.Type, p.Info.Enabled, p.Info.Description)
		}
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:   "install [plugin_path]",
	Short: "Install a plugin from a shared library file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		pluginPath := args[0]
		if !filepath.IsAbs(pluginPath) {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Printf("Error getting current directory: %v\n", err)
				os.Exit(1)
			}
			pluginPath = filepath.Join(cwd, pluginPath)
		}

		// Check if plugin file exists
		if _, err := os.Stat(pluginPath); os.IsNotExist(err) {
			fmt.Printf("Error: Plugin file %s does not exist\n", pluginPath)
			os.Exit(1)
		}

		// Install the plugin
		p, err := manager.LoadPlugin(pluginPath)
		if err != nil {
			fmt.Printf("Error installing plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin %s (version %s) installed successfully.\n", p.Info.Name, p.Info.Version)
		fmt.Printf("Type: %s\n", p.Info.Type)
		fmt.Printf("Description: %s\n", p.Info.Description)
		fmt.Printf("Author: %s\n", p.Info.Author)
	},
}

var pluginUninstallCmd = &cobra.Command{
	Use:   "uninstall [plugin_name]",
	Short: "Uninstall a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		pluginName := args[0]

		// Check if plugin is installed
		_, err := manager.GetPlugin(pluginName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Uninstall the plugin
		err = manager.UnloadPlugin(pluginName)
		if err != nil {
			fmt.Printf("Error uninstalling plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin %s uninstalled successfully.\n", pluginName)
	},
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info [plugin_name]",
	Short: "Get detailed information about a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		pluginName := args[0]

		// Get plugin information
		p, err := manager.GetPlugin(pluginName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonOutput, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				fmt.Printf("Error formatting output: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(jsonOutput))
			return
		}

		// Detailed format (default)
		fmt.Println("Name:", p.Info.Name)
		fmt.Println("Version:", p.Info.Version)
		fmt.Println("Type:", p.Info.Type)
		fmt.Println("Description:", p.Info.Description)
		fmt.Println("Author:", p.Info.Author)
		fmt.Println("Enabled:", p.Info.Enabled)
		fmt.Println("Path:", p.Path)

		// Try to get plugin options
		optionsSymbol, err := p.LookupSymbol("GetOptions")
		if err == nil {
			if optionsFunc, ok := optionsSymbol.(func() map[string]interface{}); ok {
				options := optionsFunc()
				if len(options) > 0 {
					fmt.Println("\nOptions:")
					for key, value := range options {
						fmt.Printf("  %s: %v\n", key, value)
					}
				}
			}
		}
	},
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable [plugin_name]",
	Short: "Enable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		pluginName := args[0]

		// Enable the plugin
		err := manager.EnablePlugin(pluginName)
		if err != nil {
			fmt.Printf("Error enabling plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin %s enabled successfully.\n", pluginName)
	},
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable [plugin_name]",
	Short: "Disable a plugin",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		pluginName := args[0]

		// Disable the plugin
		err := manager.DisablePlugin(pluginName)
		if err != nil {
			fmt.Printf("Error disabling plugin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Plugin %s disabled successfully.\n", pluginName)
	},
}

var pluginTypesCmd = &cobra.Command{
	Use:   "types",
	Short: "List supported plugin types",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Supported plugin types:")
		fmt.Println("- validation: Custom validation rules")
		fmt.Println("- alert: Custom alerting mechanisms")
		fmt.Println("- metric: Custom metrics collection")
		fmt.Println("- storage: Custom storage backends")
		fmt.Println("- export: Custom export formats")
		fmt.Println("- ui: Custom UI components")
	},
}

var pluginExecCmd = &cobra.Command{
	Use:   "exec [plugin_name] [function_name] [args...]",
	Short: "Execute a function in a plugin",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		manager := getPluginManager()

		pluginName := args[0]
		functionName := args[1]
		functionArgs := make([]interface{}, len(args)-2)
		for i, arg := range args[2:] {
			functionArgs[i] = arg
		}

		// Get plugin
		p, err := manager.GetPlugin(pluginName)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Execute function
		results, err := p.ExecuteFunc(functionName, functionArgs...)
		if err != nil {
			fmt.Printf("Error executing function: %v\n", err)
			os.Exit(1)
		}

		// Print results
		fmt.Println("Results:")
		for i, result := range results {
			fmt.Printf("  %d: %v\n", i, result)
		}
	},
}

func getPluginManager() *plugin.PluginManager {
	return plugin.NewPluginManager()
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
	pluginCmd.AddCommand(pluginInfoCmd)
	pluginCmd.AddCommand(pluginEnableCmd)
	pluginCmd.AddCommand(pluginDisableCmd)
	pluginCmd.AddCommand(pluginTypesCmd)
	pluginCmd.AddCommand(pluginExecCmd)

	// List command flags
	pluginListCmd.Flags().String("format", "table", "Output format (table, json)")

	// Info command flags
	pluginInfoCmd.Flags().String("format", "detailed", "Output format (detailed, json)")
}
