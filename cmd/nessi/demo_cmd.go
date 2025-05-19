package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run an interactive demo of Nessi",
	Long:  `Run an interactive demo showcasing Nessi's core features with embedded examples.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDemo()
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

// runDemo runs the interactive demo
func runDemo() {
	// Set up colors
	titleColor := color.New(color.FgHiCyan, color.Bold)
	subtitleColor := color.New(color.FgHiBlue)
	successColor := color.New(color.FgHiGreen)
	errorColor := color.New(color.FgHiRed)
	infoColor := color.New(color.FgHiYellow)
	
	// Clear the screen
	clearScreen()
	
	// Display welcome message
	titleColor.Println("╔════════════════════════════════════════════╗")
	titleColor.Println("║            NESSI QUICKSTART DEMO           ║")
	titleColor.Println("╚════════════════════════════════════════════╝")
	fmt.Println()
	subtitleColor.Println("This interactive demo will showcase Nessi's core features:")
	fmt.Println("• Delta Lake table management")
	fmt.Println("• Data quality validation")
	fmt.Println("• Schema evolution tracking")
	fmt.Println("• Monitoring and visualization")
	fmt.Println()
	infoColor.Println("Press Enter to continue...")
	waitForEnter()
	
	// Step 1: Create a demo Delta Lake table
	clearScreen()
	subtitleColor.Println("Step 1: Creating a demo Delta Lake table")
	fmt.Println()
	
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Creating demo table..."
	s.Start()
	
	// Simulate table creation
	time.Sleep(2 * time.Second)
	s.Stop()
	
	successColor.Println("✓ Demo table created at ./demo_data/users")
	fmt.Println()
	fmt.Println("Table schema:")
	fmt.Println(`
├── id: integer
├── name: string
├── email: string
├── signup_date: timestamp
└── profile: struct
    ├── age: integer
    ├── country: string
    └── preferences: array
        └── string
`)
	fmt.Println()
	infoColor.Println("Press Enter to continue...")
	waitForEnter()
	
	// Step 2: Run data quality validation
	clearScreen()
	subtitleColor.Println("Step 2: Running data quality validation")
	fmt.Println()
	
	s.Suffix = " Validating data quality..."
	s.Start()
	
	// Simulate validation
	time.Sleep(2 * time.Second)
	s.Stop()
	
	successColor.Println("✓ Validation complete")
	fmt.Println()
	fmt.Println("Quality metrics:")
	fmt.Println("• Completeness: 98.5%")
	fmt.Println("• Accuracy: 99.2%")
	fmt.Println("• Consistency: 97.8%")
	fmt.Println("• Freshness: 100%")
	fmt.Println()
	errorColor.Println("Issues found:")
	fmt.Println("• 5 null values in 'email' column")
	fmt.Println("• 3 invalid email formats")
	fmt.Println()
	infoColor.Println("Press Enter to continue...")
	waitForEnter()
	
	// Step 3: Schema evolution tracking
	clearScreen()
	subtitleColor.Println("Step 3: Schema evolution tracking")
	fmt.Println()
	
	s.Suffix = " Tracking schema evolution..."
	s.Start()
	
	// Simulate schema tracking
	time.Sleep(2 * time.Second)
	s.Stop()
	
	successColor.Println("✓ Schema history retrieved")
	fmt.Println()
	fmt.Println("Schema changes:")
	fmt.Println("• Version 1 (2025-01-01): Initial schema")
	fmt.Println("• Version 2 (2025-02-15): Added 'profile.preferences' array")
	fmt.Println("• Version 3 (2025-04-10): Added 'signup_date' column")
	fmt.Println()
	infoColor.Println("Press Enter to continue...")
	waitForEnter()
	
	// Step 4: Visualization
	clearScreen()
	subtitleColor.Println("Step 4: Data quality visualization")
	fmt.Println()
	
	s.Suffix = " Generating visualization..."
	s.Start()
	
	// Simulate visualization generation
	time.Sleep(2 * time.Second)
	s.Stop()
	
	successColor.Println("✓ Visualization generated")
	fmt.Println()
	fmt.Println("Data quality heatmap by partition:")
	fmt.Println()
	
	// Display ASCII heatmap
	blue := color.New(color.BgBlue).SprintFunc()
	yellow := color.New(color.BgYellow).SprintFunc()
	red := color.New(color.BgRed).SprintFunc()
	
	fmt.Println(blue("   ") + blue("   ") + blue("   ") + yellow("   ") + yellow("   "))
	fmt.Println(blue("   ") + blue("   ") + yellow("   ") + yellow("   ") + yellow("   "))
	fmt.Println(blue("   ") + yellow("   ") + yellow("   ") + yellow("   ") + red("   "))
	fmt.Println(yellow("   ") + yellow("   ") + yellow("   ") + red("   ") + red("   "))
	fmt.Println(yellow("   ") + yellow("   ") + red("   ") + red("   ") + red("   "))
	fmt.Println()
	fmt.Println("Legend: " + blue("   ") + " High quality  " + yellow("   ") + " Medium quality  " + red("   ") + " Low quality")
	fmt.Println()
	infoColor.Println("Press Enter to continue...")
	waitForEnter()
	
	// Final step: Summary
	clearScreen()
	titleColor.Println("╔════════════════════════════════════════════╗")
	titleColor.Println("║          NESSI DEMO COMPLETE!              ║")
	titleColor.Println("╚════════════════════════════════════════════╝")
	fmt.Println()
	subtitleColor.Println("You've successfully completed the Nessi quickstart demo!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("1. Run 'nessi help' to explore available commands")
	fmt.Println("2. Check out the documentation at https://github.com/nessi-dev/nessi/docs")
	fmt.Println("3. Try validating your own Delta Lake tables with 'nessi validate'")
	fmt.Println()
	successColor.Println("Happy data quality management with Nessi!")
}

// clearScreen clears the terminal screen
func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Print("\033[H\033[2J")
	case "windows":
		fmt.Print("\033[H\033[2J")
	default:
		// Fallback for unsupported OS
		fmt.Println("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")
	}
}

// waitForEnter waits for the user to press Enter
func waitForEnter() {
	fmt.Scanln()
}
