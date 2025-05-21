package cli

import (
	"fmt"
	"github.com/nessi-dev/nessi/pkg"
	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/spf13/cobra"
)

var telemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Manage telemetry settings",
}

var telemetryEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable telemetry",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.UpdateConfig(map[string]interface{}{"telemetry.enabled": true})
		if err != nil {
			fmt.Println("❌ Failed to enable telemetry:", err)
			return
		}
		fmt.Println("✅ Telemetry enabled.")
	},
}

var telemetryDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable telemetry",
	Run: func(cmd *cobra.Command, args []string) {
		err := pkg.UpdateConfig(map[string]interface{}{"telemetry.enabled": false})
		if err != nil {
			fmt.Println("❌ Failed to disable telemetry:", err)
			return
		}
		fmt.Println("✅ Telemetry disabled.")
	},
}

var telemetryStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show telemetry status",
	Run: func(cmd *cobra.Command, args []string) {
		// General telemetry status
		cfg := pkg.GetTelemetryConfig()
		if cfg.Enabled {
			fmt.Println("✅ General Telemetry: enabled")
			fmt.Println("Endpoint:", cfg.Endpoint)
		} else {
			fmt.Println("❌ General Telemetry: disabled")
		}
		
		// Error telemetry status
		errorTelemetryConfig := common.DefaultErrorTelemetryConfig()
		errorTelemetry := common.NewErrorTelemetry(errorTelemetryConfig)
		stats := errorTelemetry.GetErrorStats()
		
		fmt.Println()
		if errorTelemetry.Enabled {
			fmt.Println("✅ Error Telemetry: enabled")
		} else {
			fmt.Println("❌ Error Telemetry: disabled")
		}
		
		// Show anonymity status
		if errorTelemetry.Anonymous {
			fmt.Println("✅ Error Telemetry is anonymous")
		} else {
			fmt.Println("❌ Error Telemetry is not anonymous")
		}
		
		// Show error stats
		fmt.Println("Total errors recorded:", stats["total_errors"])
		
		// Show storage path
		fmt.Println("Error telemetry storage path:", errorTelemetry.StoragePath)
	},
}

var errorTelemetryEnableCmd = &cobra.Command{
	Use:   "enable-error",
	Short: "Enable error telemetry",
	Run: func(cmd *cobra.Command, args []string) {
		// Get error telemetry
		errorTelemetryConfig := common.DefaultErrorTelemetryConfig()
		errorTelemetry := common.NewErrorTelemetry(errorTelemetryConfig)
		
		// Enable error telemetry
		errorTelemetry.EnableTelemetry()
		
		// Save telemetry data
		err := errorTelemetry.Save()
		if err != nil {
			fmt.Println("❌ Failed to save error telemetry settings:", err)
			return
		}
		
		fmt.Println("✅ Error telemetry enabled.")
		fmt.Println("Anonymous error statistics will be collected to help improve Nessi.")
		fmt.Println("No personal data or table contents will be collected.")
	},
}

var errorTelemetryDisableCmd = &cobra.Command{
	Use:   "disable-error",
	Short: "Disable error telemetry",
	Run: func(cmd *cobra.Command, args []string) {
		// Get error telemetry
		errorTelemetryConfig := common.DefaultErrorTelemetryConfig()
		errorTelemetry := common.NewErrorTelemetry(errorTelemetryConfig)
		
		// Disable error telemetry
		errorTelemetry.DisableTelemetry()
		
		// Save telemetry data
		err := errorTelemetry.Save()
		if err != nil {
			fmt.Println("❌ Failed to save error telemetry settings:", err)
			return
		}
		
		fmt.Println("✅ Error telemetry disabled.")
	},
}

var errorTelemetryReportCmd = &cobra.Command{
	Use:   "report-errors",
	Short: "Report error telemetry statistics",
	Run: func(cmd *cobra.Command, args []string) {
		// Get error telemetry
		errorTelemetryConfig := common.DefaultErrorTelemetryConfig()
		errorTelemetry := common.NewErrorTelemetry(errorTelemetryConfig)
		
		// Report telemetry data
		err := errorTelemetry.ReportTelemetry()
		if err != nil {
			fmt.Println("❌ Failed to report error telemetry:", err)
			return
		}
		
		fmt.Println("✅ Error telemetry reported successfully.")
		
		// Show error stats
		stats := errorTelemetry.GetErrorStats()
		fmt.Println("Total errors recorded:", stats["total_errors"])
		
		// Show error counts by error code if available
		errorCounts, ok := stats["error_counts"].(map[common.ErrorCode]int)
		if ok && len(errorCounts) > 0 {
			fmt.Println("
Error counts by error code:")
			for code, count := range errorCounts {
				fmt.Printf("  %s (%s): %d
", code, common.GetErrorDescription(code), count)
			}
		}
	},
}

var errorTelemetryExportCmd = &cobra.Command{
	Use:   "export-errors [file]",
	Short: "Export error telemetry to a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get file path
		filePath := args[0]
		
		// Get error telemetry
		errorTelemetryConfig := common.DefaultErrorTelemetryConfig()
		errorTelemetry := common.NewErrorTelemetry(errorTelemetryConfig)
		
		// Export telemetry data to file
		err := errorTelemetry.ExportTelemetryToFile(filePath)
		if err != nil {
			fmt.Println("❌ Failed to export error telemetry:", err)
			return
		}
		
		fmt.Println("✅ Error telemetry exported successfully to", filePath)
		
		// Show error stats
		stats := errorTelemetry.GetErrorStats()
		fmt.Println("Total errors recorded:", stats["total_errors"])
	},
}

func init() {
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryStatusCmd)
	telemetryCmd.AddCommand(errorTelemetryEnableCmd)
	telemetryCmd.AddCommand(errorTelemetryDisableCmd)
	telemetryCmd.AddCommand(errorTelemetryReportCmd)
	telemetryCmd.AddCommand(errorTelemetryExportCmd)
}
