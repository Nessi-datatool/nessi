// Package main provides a tool to start a free trial
package main

import (
	"fmt"
	"os"

	"github.com/nessi-dev/nessi/pkg/security"
)

func main() {
	// Create license manager
	licenseManager := security.NewLicenseManager()

	// Start trial
	err := licenseManager.StartTrial()
	if err != nil {
		fmt.Printf("Error starting trial: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Trial successfully started!")
	fmt.Println("You now have access to all premium features for 1 month.")

	// Print license info
	if err := licenseManager.PrintLicenseInfo(); err != nil {
		fmt.Printf("Error printing license info: %v\n", err)
	}
}
