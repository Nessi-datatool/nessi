package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/spf13/cobra"
)

// initTestErrorCmd initializes the test-error command
func init() {
	initTestErrorCmd()
}

func initTestErrorCmd() {
	testErrorCmd := &cobra.Command{
		Use:   "test-error [error-code]",
		Short: "Test error handling system",
		Long:  "Generate a test error to verify the error handling system",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get error code
			errorCode := common.ErrUnknown
			if len(args) > 0 {
				errorCode = common.ErrorCode(args[0])
			}

			// Get flags
			details, _ := cmd.Flags().GetString("details")
			suggestion, _ := cmd.Flags().GetString("suggestion")
			resolvable, _ := cmd.Flags().GetBool("resolvable")
			telemetry, _ := cmd.Flags().GetBool("telemetry")
			exportPath, _ := cmd.Flags().GetString("export")
			listCodes, _ := cmd.Flags().GetBool("list")

			// List all error codes if requested
			if listCodes {
				fmt.Println("Available error codes:")
				fmt.Println("------------------")
				for _, code := range common.GetAllErrorCodes() {
					fmt.Printf("%s: %s\n", code, common.GetErrorDescription(code))
				}
				return nil
			}

			// Create error message
			message := fmt.Sprintf("Test error for code %s", errorCode)

			// Create error
			var err error
			if resolvable {
				// Create a resolvable error
				switch errorCode {
				case common.ErrInvalidPath:
					// Create a non-existent path
					testPath := "/tmp/nessi-test-path-that-does-not-exist"
					err = common.NewResolvablePathError(testPath)
				case common.ErrInvalidConfig:
					// Create a config error
					err = common.NewResolvableConfigError("test.config.key", "invalid-value")
				case common.ErrInvalidDeltaTable:
					// Create a delta table error
					testPath := "/tmp/nessi-test-delta-table-that-does-not-exist"
					err = common.NewResolvableDeltaTableError(testPath)
				default:
					// Create a generic resolvable error
					err = common.NewResolvableError(errorCode, message, details, []string{suggestion})
				}
			} else {
				// Create a standard NessiError
				nessiErr := common.NewError(errorCode, message)
				if details != "" {
					nessiErr.Details = details
				}
				if suggestion != "" {
					nessiErr.Suggestions = []string{suggestion}
				}
				err = nessiErr
			}

			// Handle telemetry if requested
			if telemetry {
				// Get error telemetry
				errorTelemetryConfig := common.DefaultErrorTelemetryConfig()
				errorTelemetry := common.NewErrorTelemetry(errorTelemetryConfig)

				// Record the error
				errorTelemetry.RecordError(err)

				// Save telemetry data
				telemetryErr := errorTelemetry.Save()
				if telemetryErr != nil {
					fmt.Fprintf(os.Stderr, "Failed to save telemetry data: %s\n", telemetryErr)
				} else {
					fmt.Fprintf(os.Stderr, "Error recorded in telemetry at %s\n", errorTelemetry.StoragePath)
				}

				// Export telemetry if requested
				if exportPath != "" {
					// Create absolute path if relative
					if !filepath.IsAbs(exportPath) {
						cwd, _ := os.Getwd()
						exportPath = filepath.Join(cwd, exportPath)
					}

					// Export telemetry data
					exportErr := errorTelemetry.ExportTelemetryToFile(exportPath)
					if exportErr != nil {
						fmt.Fprintf(os.Stderr, "Failed to export telemetry data: %s\n", exportErr)
					} else {
						fmt.Fprintf(os.Stderr, "Telemetry data exported to %s\n", exportPath)
					}
				}
			}

			// Return the error to trigger error handling
			return err
		},
	}

	// Add flags
	testErrorCmd.Flags().String("details", "", "Error details")
	testErrorCmd.Flags().String("suggestion", "", "Error suggestion")
	testErrorCmd.Flags().Bool("resolvable", false, "Create a resolvable error")
	testErrorCmd.Flags().Bool("telemetry", false, "Record error in telemetry")
	testErrorCmd.Flags().String("export", "", "Export telemetry data to file")
	testErrorCmd.Flags().Bool("list", false, "List all available error codes")

	// Add to root command
	CLI.RootCmd.AddCommand(testErrorCmd)
}
