package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/nessi-dev/nessi/pkg"
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
		cfg := pkg.GetTelemetryConfig()
		if cfg.Enabled {
			fmt.Println("✅ Telemetry: enabled")
			fmt.Println("Endpoint:", cfg.Endpoint)
		} else {
			fmt.Println("❌ Telemetry: disabled")
		}
	},
}

func init() {
	telemetryCmd.AddCommand(telemetryEnableCmd)
	telemetryCmd.AddCommand(telemetryDisableCmd)
	telemetryCmd.AddCommand(telemetryStatusCmd)
}
