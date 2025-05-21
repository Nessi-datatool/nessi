package cli

import (
	"fmt"
	"time"

	"github.com/nessi-dev/nessi/pkg/common"
	"github.com/spf13/cobra"
)

// progressDemoCmd represents the command for demonstrating progress indicators
var progressDemoCmd = &cobra.Command{
	Use:   "progress-demo",
	Short: "Demonstrate progress indicators and bars",
	Long: `Demonstrate the various progress indicators and bars available in Nessi.
This command is useful for seeing how progress is visualized during long-running operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get flags
		showTiming, _ := cmd.Flags().GetBool("timing")
		barDemo, _ := cmd.Flags().GetBool("bar")
		spinnerDemo, _ := cmd.Flags().GetBool("spinner")

		if !barDemo && !spinnerDemo {
			// If no specific demo is requested, show both
			barDemo = true
			spinnerDemo = true
		}

		fmt.Println("Nessi Progress Visualization Demo")
		fmt.Println("=================================")

		if spinnerDemo {
			fmt.Println("\n1. Progress Indicator (Spinner)")
			fmt.Println("------------------------------")
			demoProgressIndicator(showTiming)
		}

		if barDemo {
			fmt.Println("\n2. Progress Bar")
			fmt.Println("---------------")
			demoProgressBar()
		}

		fmt.Println("\nDemo completed! These progress visualizations are used throughout Nessi for long-running operations.")
	},
}

func demoProgressIndicator(showTiming bool) {
	// Demo success case
	fmt.Println("\nSuccess case:")
	pi := common.NewProgressIndicator("Processing data...", showTiming)
	pi.Start()
	time.Sleep(2 * time.Second)
	pi.Success("Data processed successfully")

	// Demo warning case
	fmt.Println("\nWarning case:")
	pi = common.NewProgressIndicator("Validating schema...", showTiming)
	pi.Start()
	time.Sleep(1 * time.Second)
	pi.Warning("Schema validated with warnings")

	// Demo error case
	fmt.Println("\nError case:")
	pi = common.NewProgressIndicator("Connecting to server...", showTiming)
	pi.Start()
	time.Sleep(1500 * time.Millisecond)
	pi.Error("Failed to connect to server")

	// Demo message update
	fmt.Println("\nMessage update case:")
	pi = common.NewProgressIndicator("Starting operation...", showTiming)
	pi.Start()
	time.Sleep(1 * time.Second)
	pi.UpdateMessage("Phase 1: Loading data...")
	time.Sleep(1 * time.Second)
	pi.UpdateMessage("Phase 2: Processing data...")
	time.Sleep(1 * time.Second)
	pi.UpdateMessage("Phase 3: Finalizing...")
	time.Sleep(500 * time.Millisecond)
	pi.Success("Operation completed successfully")
}

func demoProgressBar() {
	// Demo simple progress bar
	fmt.Println("\nSimple progress bar:")
	total := 20
	bar := common.NewProgressBar(total, "Processing items", 40)

	for i := 0; i <= total; i++ {
		bar.Update(i)
		time.Sleep(200 * time.Millisecond)
	}

	// Demo incremental progress bar
	fmt.Println("\nIncremental progress bar:")
	total = 10
	bar = common.NewProgressBar(total, "Uploading files", 40)

	for i := 0; i < total; i++ {
		bar.Increment()
		time.Sleep(300 * time.Millisecond)
	}
}

// InitProgressDemoCmd initializes the progress demo command
func InitProgressDemoCmd() {
	// Add the progress-demo command to the root command
	CLI.RootCmd.AddCommand(progressDemoCmd)

	// Add flags
	progressDemoCmd.Flags().Bool("timing", false, "Show timing information")
	progressDemoCmd.Flags().Bool("bar", false, "Demo progress bar only")
	progressDemoCmd.Flags().Bool("spinner", false, "Demo spinner only")
}
