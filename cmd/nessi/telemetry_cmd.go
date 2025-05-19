package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi/pkg/telemetry"
	"github.com/spf13/cobra"
)

// telemetryCmd represents the telemetry command
var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage telemetry settings",
	Long:  `Enable or disable anonymous usage telemetry and manage GitHub usage badge.`,
}

// telemetryEnableCmd enables telemetry
var telemetryEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable anonymous telemetry",
	Long:  `Enable anonymous telemetry to help improve Nessi. No sensitive data is collected.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := telemetry.Enable()
		if err != nil {
			fmt.Printf("Error enabling telemetry: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Telemetry enabled. Thank you for helping improve Nessi!")
		fmt.Println("You can disable telemetry at any time with 'nessi telemetry disable'")
	},
}

// telemetryDisableCmd disables telemetry
var telemetryDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable anonymous telemetry",
	Long:  `Disable anonymous telemetry collection.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := telemetry.Disable()
		if err != nil {
			fmt.Printf("Error disabling telemetry: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Telemetry disabled.")
	},
}

// telemetryStatusCmd shows the current telemetry status
var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show telemetry status",
	Long:  `Display the current status of anonymous telemetry collection.`,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := telemetry.GetStatus()
		if err != nil {
			fmt.Printf("Error getting telemetry status: %v\n", err)
			os.Exit(1)
		}

		if status {
			fmt.Println("Telemetry is currently ENABLED")
		} else {
			fmt.Println("Telemetry is currently DISABLED")
		}
	},
}

// badgeCmd represents the badge command
var badgeCmd = &cobra.Command{
	Use:   "badge",
	Short: "Generate GitHub usage badge",
	Long:  `Generate a GitHub usage badge for your repository to show that you're using Nessi.`,
	Run: func(cmd *cobra.Command, args []string) {
		repoURL, _ := cmd.Flags().GetString("repo")
		outputPath, _ := cmd.Flags().GetString("output")
		
		// Generate the badge
		badge, err := telemetry.GenerateBadge(repoURL)
		if err != nil {
			fmt.Printf("Error generating badge: %v\n", err)
			os.Exit(1)
		}
		
		// If output path is provided, save the badge to a file
		if outputPath != "" {
			err = os.WriteFile(outputPath, []byte(badge), 0644)
			if err != nil {
				fmt.Printf("Error saving badge to file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Badge saved to %s\n", outputPath)
		} else {
			// Otherwise, print the badge to stdout
			fmt.Println("Add this to your README.md:")
			fmt.Println()
			fmt.Println(badge)
		}
	},
}

func init() {
	rootCmd.AddCommand(telemetryCmd)
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryStatusCmd)
	telemetryCmd.AddCommand(badgeCmd)
	
	// Add flags for the badge command
	badgeCmd.Flags().StringP("repo", "r", "", "GitHub repository URL (e.g., https://github.com/username/repo)")
	badgeCmd.Flags().StringP("output", "o", "", "Output file path for the badge markdown")
}
