// Package main provides a test program to verify premium feature gating
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/nessi-dev/nessi/pkg/security"
)

// Test all premium features to ensure they are properly gated
func main() {
	// Define colors for output
	success := color.New(color.FgGreen).SprintFunc()
	failure := color.New(color.FgRed).SprintFunc()
	info := color.New(color.FgCyan).SprintFunc()
	header := color.New(color.Bold).SprintFunc()

	// Print header
	fmt.Println(header("Nessi Premium Feature Test"))
	fmt.Println(strings.Repeat("-", 50))
	fmt.Println(info("This test verifies that premium features are properly gated"))
	fmt.Println()

	// Get license info
	licenseManager := security.NewLicenseManager()
	licenseInfo, err := licenseManager.GetLicenseInfo()
	if err != nil {
		fmt.Printf("Error getting license info: %v\n", err)
		os.Exit(1)
	}

	// Print current license status
	fmt.Printf("Current license status: %s\n", licenseInfo.Status)
	fmt.Printf("Current plan: %s\n", licenseInfo.Plan)
	if licenseInfo.IsTrial {
		fmt.Printf("Trial active: Yes (expires in %d days)\n", licenseInfo.DaysRemaining)
	} else if licenseInfo.Status == "Active" {
		fmt.Printf("License expires in %d days\n", licenseInfo.DaysRemaining)
	}
	fmt.Println()

	// Test all premium features
	fmt.Println(header("Testing Premium Features"))
	fmt.Println(strings.Repeat("-", 50))

	// Define all premium features to test
	features := []struct {
		name        string
		featureCode string
		plan        string
	}{
		{"Databricks Integration", security.FeatureDatabricksIntegration, "Starter"},
		{"DBT Integration", security.FeatureDBTIntegration, "Pro"},
		{"AWS S3 Storage", security.FeatureCloudStorageS3, "Starter"},
		{"Azure Blob Storage", security.FeatureCloudStorageAzure, "Pro"},
		{"Google Cloud Storage", security.FeatureCloudStorageGCP, "Pro"},
		{"Data Catalog", security.FeatureDataCatalog, "Pro"},
		{"Workflow Orchestration", security.FeatureWorkflowOrchestration, "Enterprise"},
	}

	// Test each feature
	for _, feature := range features {
		err := security.RequireFeature(feature.featureCode)
		if err == nil {
			fmt.Printf("%-30s %s (included in %s plan)\n", feature.name+":", success("ACCESSIBLE"), feature.plan)
		} else {
			fmt.Printf("%-30s %s (%s)\n", feature.name+":", failure("BLOCKED"), err.Error())
		}
	}

	fmt.Println()
	fmt.Println(header("Test Complete"))

	// Provide guidance based on test results
	if licenseInfo.Status == "Not Licensed" && !licenseInfo.IsTrial {
		fmt.Println(info("No license or trial detected. Run 'nessi license trial' to start a free trial."))
	} else if licenseInfo.IsTrial {
		fmt.Println(info("Trial active. All premium features should be accessible."))
		fmt.Println(info("If any features are blocked, there may be an issue with the trial activation."))
	} else if licenseInfo.Status == "Active" {
		fmt.Println(info("License active. Features should be accessible based on your plan."))
		fmt.Println(info("If expected features are blocked, check your license details."))
	} else if licenseInfo.Status == "Expired" {
		fmt.Println(info("License expired. Premium features are blocked until renewal."))
	}
}
