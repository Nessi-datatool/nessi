// Package main provides a debug tool for trial activation
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nessi-dev/nessi/pkg/security"
)

func main() {
	fmt.Println("Nessi Trial Debug Tool")
	fmt.Println("======================")

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Check if .nessi directory exists
	nessiDir := filepath.Join(homeDir, ".nessi")
	if _, err := os.Stat(nessiDir); os.IsNotExist(err) {
		fmt.Println("The .nessi directory does not exist. Creating it...")
		if err := os.MkdirAll(nessiDir, 0755); err != nil {
			fmt.Printf("Error creating .nessi directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Check for trial.json
	trialPath := filepath.Join(nessiDir, "trial.json")
	if _, err := os.Stat(trialPath); os.IsNotExist(err) {
		fmt.Println("No trial.json file found.")
	} else {
		fmt.Println("Found trial.json file. Reading contents...")
		data, err := os.ReadFile(trialPath)
		if err != nil {
			fmt.Printf("Error reading trial.json: %v\n", err)
		} else {
			fmt.Println("Raw file content:")
			fmt.Println(string(data))

			// Try to parse as SignedData
			var signedData security.SignedData
			if err := json.Unmarshal(data, &signedData); err != nil {
				fmt.Printf("Error parsing as SignedData: %v\n", err)

				// Try as plain TrialInfo
				var trialInfo security.TrialInfo
				if err := json.Unmarshal(data, &trialInfo); err != nil {
					fmt.Printf("Error parsing as TrialInfo: %v\n", err)
				} else {
					fmt.Println("Successfully parsed as plain TrialInfo:")
					fmt.Printf("  Start Time: %v\n", trialInfo.StartTime)
					fmt.Printf("  Active: %v\n", trialInfo.Active)
					fmt.Printf("  Machine ID: %v\n", trialInfo.MachineID)
					fmt.Printf("  Trial Count: %v\n", trialInfo.TrialCount)
				}
			} else {
				fmt.Println("Successfully parsed as SignedData")
				fmt.Printf("  Signature: %s\n", signedData.Signature)

				// Verify signature
				if security.VerifySignature(signedData.Data, signedData.Signature) {
					fmt.Println("  Signature verification: PASSED")

					// Parse the data
					var trialInfo security.TrialInfo
					if err := json.Unmarshal(signedData.Data, &trialInfo); err != nil {
						fmt.Printf("  Error parsing trial data: %v\n", err)
					} else {
						fmt.Println("  Trial Info:")
						fmt.Printf("    Start Time: %v\n", trialInfo.StartTime)
						fmt.Printf("    Active: %v\n", trialInfo.Active)
						fmt.Printf("    Machine ID: %v\n", trialInfo.MachineID)
						fmt.Printf("    Trial Count: %v\n", trialInfo.TrialCount)
					}
				} else {
					fmt.Println("  Signature verification: FAILED")
				}
			}
		}
	}

	// Check for trial_registry.json
	registryPath := filepath.Join(nessiDir, "trial_registry.json")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		fmt.Println("No trial_registry.json file found.")
	} else {
		fmt.Println("Found trial_registry.json file. Reading contents...")
		data, err := os.ReadFile(registryPath)
		if err != nil {
			fmt.Printf("Error reading trial_registry.json: %v\n", err)
		} else {
			fmt.Println("Raw file content:")
			fmt.Println(string(data))
		}
	}

	// Test machine ID generation
	machineID, err := security.GetMachineID()
	if err != nil {
		fmt.Printf("Error getting machine ID: %v\n", err)
	} else {
		fmt.Printf("Current machine ID: %s\n", machineID)
	}

	// Test trial activation
	fmt.Println("\nTesting trial activation...")
	licenseManager := security.NewLicenseManager()
	err = licenseManager.StartTrial()
	if err != nil {
		fmt.Printf("Error starting trial: %v\n", err)
	} else {
		fmt.Println("Trial successfully started!")
	}

	// Check license info
	fmt.Println("\nChecking license info...")
	info, err := licenseManager.GetLicenseInfo()
	if err != nil {
		fmt.Printf("Error getting license info: %v\n", err)
	} else {
		fmt.Printf("Status: %s\n", info.Status)
		fmt.Printf("Plan: %s\n", info.Plan)
		fmt.Printf("Is Trial: %v\n", info.IsTrial)
		if !info.Expires.IsZero() {
			fmt.Printf("Expires: %v\n", info.Expires)
			fmt.Printf("Days Remaining: %d\n", info.DaysRemaining)
		}
		if info.Message != "" {
			fmt.Printf("Message: %s\n", info.Message)
		}
	}

	// Test feature access
	fmt.Println("\nTesting premium feature access...")
	err = security.RequireFeature(security.FeatureDatabricksIntegration)
	if err == nil {
		fmt.Println("Databricks Integration: ACCESSIBLE")
	} else {
		fmt.Printf("Databricks Integration: BLOCKED (%v)\n", err)
	}
}
