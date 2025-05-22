package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"nessi/pkg/security"
)

func main() {
	fmt.Println("Nessi Trial Expiration Test")
	fmt.Println("--------------------------------------------------")
	fmt.Println("This test verifies the behavior when a trial expires")
	fmt.Println()

	// Clear existing trial files
	homeDir, _ := os.UserHomeDir()
	trialPath := homeDir + "/.nessi/trial.json"
	registryPath := homeDir + "/.nessi/trial_registry.json"

	os.Remove(trialPath)
	os.Remove(registryPath)

	// Start a new trial
	fmt.Println("Starting a new trial...")
	err := security.StartTrial()
	if err != nil {
		fmt.Printf("Error starting trial: %v\n", err)
		return
	}
	fmt.Println("Trial started successfully.")

	// Get and display license info
	fmt.Println("\nCurrent license info:")
	info, err := security.GetLicenseInfo()
	if err != nil {
		fmt.Printf("Error getting license info: %v\n", err)
	} else {
		displayLicenseInfo(info)
	}

	// Test premium feature access
	fmt.Println("\nTesting premium feature access with active trial:")
	testFeatureAccess("databricks_integration")

	// Modify the trial file to make it expired (1 month and 1 day ago)
	fmt.Println("\nSimulating trial expiration...")
	expiredTime := time.Now().AddDate(0, -1, -1)
	err = security.SetTrialStartTimeForTesting(expiredTime)
	if err != nil {
		fmt.Printf("Error setting expired trial time: %v\n", err)
		return
	}
	fmt.Println("Trial has been expired for testing purposes.")

	// Get and display license info after expiration
	fmt.Println("\nLicense info after trial expiration:")
	info, err = security.GetLicenseInfo()
	if err != nil {
		fmt.Printf("Error getting license info: %v\n", err)
	} else {
		displayLicenseInfo(info)
	}

	// Test premium feature access after expiration
	fmt.Println("\nTesting premium feature access with expired trial:")
	testFeatureAccess("databricks_integration")

	fmt.Println("\nTest Complete")
}

func displayLicenseInfo(info security.LicenseInfo) {
	fmt.Printf("Status: %s\n", info.Status)
	fmt.Printf("Plan: %s\n", info.Plan)
	fmt.Printf("Is Trial: %v\n", info.IsTrial)
	if info.ExpirationDate != nil {
		fmt.Printf("Expires: %v\n", *info.ExpirationDate)
		if info.DaysRemaining != nil {
			fmt.Printf("Days Remaining: %d\n", *info.DaysRemaining)
		}
	}
	fmt.Printf("Message: %s\n", info.Message)
}

func testFeatureAccess(feature string) {
	err := security.RequireFeature(feature)
	if err != nil {
		color.Red("Feature %s: BLOCKED (%v)", feature, err)
	} else {
		color.Green("Feature %s: ACCESSIBLE", feature)
	}
}
