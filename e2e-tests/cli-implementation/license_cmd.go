package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func newLicenseCmd() *cobra.Command {
	licenseCmd := &cobra.Command{
		Use:   "license",
		Short: "Manage Nessi license",
		Long:  `Commands for managing Nessi license and checking license status.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	licenseCmd.AddCommand(newLicenseStatusCmd())
	licenseCmd.AddCommand(newLicenseTrialCmd())

	return licenseCmd
}

func newLicenseStatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check license status",
		Long:  `Check the current license status and available features.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Output license status
			fmt.Println("License Status: Active")
			fmt.Println("Edition: Community Edition")
			fmt.Println("Features:")
			fmt.Println("  - Delta Lake Table Management")
			fmt.Println("  - Basic Data Quality Checks")
			fmt.Println("  - Schema Validation")
			fmt.Println("  - Basic Reporting (HTML, CSV)")
			fmt.Println("")
			fmt.Println("For Pro Edition features, activate a trial with 'nessi license trial'")
		},
	}

	return statusCmd
}

func newLicenseTrialCmd() *cobra.Command {
	trialCmd := &cobra.Command{
		Use:   "trial",
		Short: "Activate Pro Edition trial",
		Long:  `Activate a 30-day trial of Nessi Pro Edition.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Output trial activation
			now := time.Now()
			expiryDate := now.AddDate(0, 1, 0)
			
			fmt.Println("Pro Edition Trial Activated!")
			fmt.Println("Trial Period: 30 days")
			fmt.Printf("Expiry Date: %s\n", expiryDate.Format("2006-01-02"))
			fmt.Println("")
			fmt.Println("Pro Edition Features:")
			fmt.Println("  - Advanced Data Quality Checks")
			fmt.Println("  - Time Travel Capabilities")
			fmt.Println("  - Databricks Integration")
			fmt.Println("  - AWS S3 Storage Integration")
			fmt.Println("  - Advanced Reporting (PDF, Interactive)")
			fmt.Println("  - Workflow Orchestration")
			fmt.Println("")
			fmt.Println("For Enterprise features, please contact sales@nessi.dev")
		},
	}

	return trialCmd
}
