package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nessi-dev/nessi/pkg/datalake"
	"github.com/spf13/cobra"
)

func init() {
	// Add consistency command
	var consistencyCmd = &cobra.Command{
		Use:   "consistency",
		Short: "Check consistency between versions of a Delta table",
		Long:  "Check consistency between versions of a Delta table, including schema, types, and fields",
	}
	rootCmd.AddCommand(consistencyCmd)

	// Add schema consistency check command
	var schemaCmd = &cobra.Command{
		Use:   "schema [table_path] [versions...]",
		Short: "Check schema consistency between versions",
		Long:  "Check schema consistency between versions of a Delta table",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			versions, err := parseVersions(args[1:])
			if err != nil {
				return err
			}

			outputFormat, _ := cmd.Flags().GetString("output")
			return runConsistencyCheck(tablePath, datalake.SchemaConsistencyCheck, versions, outputFormat)
		},
	}
	schemaCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	consistencyCmd.AddCommand(schemaCmd)

	// Add type consistency check command
	var typeCmd = &cobra.Command{
		Use:   "type [table_path] [versions...]",
		Short: "Check type consistency between versions",
		Long:  "Check type consistency between versions of a Delta table",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			versions, err := parseVersions(args[1:])
			if err != nil {
				return err
			}

			outputFormat, _ := cmd.Flags().GetString("output")
			return runConsistencyCheck(tablePath, datalake.TypeConsistencyCheck, versions, outputFormat)
		},
	}
	typeCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	consistencyCmd.AddCommand(typeCmd)

	// Add field consistency check command
	var fieldCmd = &cobra.Command{
		Use:   "field [table_path] [versions...]",
		Short: "Check field consistency between versions",
		Long:  "Check field consistency between versions of a Delta table",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			versions, err := parseVersions(args[1:])
			if err != nil {
				return err
			}

			outputFormat, _ := cmd.Flags().GetString("output")
			return runConsistencyCheck(tablePath, datalake.FieldConsistencyCheck, versions, outputFormat)
		},
	}
	fieldCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	consistencyCmd.AddCommand(fieldCmd)

	// Add all consistency checks command
	var allCmd = &cobra.Command{
		Use:   "all [table_path] [versions...]",
		Short: "Run all consistency checks between versions",
		Long:  "Run all consistency checks between versions of a Delta table",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			versions, err := parseVersions(args[1:])
			if err != nil {
				return err
			}

			outputFormat, _ := cmd.Flags().GetString("output")
			
			// Run all checks
			fmt.Println("Running schema consistency check...")
			err = runConsistencyCheck(tablePath, datalake.SchemaConsistencyCheck, versions, outputFormat)
			if err != nil {
				return err
			}
			
			fmt.Println("\nRunning type consistency check...")
			err = runConsistencyCheck(tablePath, datalake.TypeConsistencyCheck, versions, outputFormat)
			if err != nil {
				return err
			}
			
			fmt.Println("\nRunning field consistency check...")
			err = runConsistencyCheck(tablePath, datalake.FieldConsistencyCheck, versions, outputFormat)
			if err != nil {
				return err
			}
			
			return nil
		},
	}
	allCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	consistencyCmd.AddCommand(allCmd)
}

// parseVersions parses version strings to int64 slice
func parseVersions(versionStrs []string) ([]int64, error) {
	versions := make([]int64, len(versionStrs))
	for i, vStr := range versionStrs {
		v, err := strconv.ParseInt(vStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid version '%s': %w", vStr, err)
		}
		versions[i] = v
	}
	return versions, nil
}

// runConsistencyCheck runs a consistency check and outputs the result
func runConsistencyCheck(tablePath string, checkType datalake.ConsistencyCheckType, versions []int64, outputFormat string) error {
	// Create metadata manager
	manager := datalake.NewMetadataManager(tablePath)
	
	// Run consistency check
	result, err := manager.CheckConsistency(checkType, versions)
	if err != nil {
		return fmt.Errorf("consistency check failed: %w", err)
	}
	
	// Output the result
	printConsistencyResult(result, outputFormat)
	return nil
}

// printConsistencyResult prints the consistency check result in the specified format
func printConsistencyResult(result *datalake.ConsistencyCheckResult, format string) {
	switch format {
	case "json":
		// Marshal to JSON
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
		
	default: // text format
		fmt.Printf("Consistency Check Result: %s\n", result.CheckType)
		fmt.Printf("Versions Compared: %s\n", formatVersions(result.VersionsCompared))
		fmt.Printf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
		fmt.Printf("Passed: %t\n", result.Passed)
		
		if len(result.Issues) > 0 {
			fmt.Printf("\nIssues Found (%d):\n", len(result.Issues))
			
			// Group issues by severity
			errorIssues := []datalake.ConsistencyIssue{}
			warningIssues := []datalake.ConsistencyIssue{}
			infoIssues := []datalake.ConsistencyIssue{}
			
			for _, issue := range result.Issues {
				switch issue.Severity {
				case "error":
					errorIssues = append(errorIssues, issue)
				case "warning":
					warningIssues = append(warningIssues, issue)
				case "info":
					infoIssues = append(infoIssues, issue)
				}
			}
			
			// Print errors first
			if len(errorIssues) > 0 {
				fmt.Printf("\nErrors (%d):\n", len(errorIssues))
				for _, issue := range errorIssues {
					fmt.Printf("  - %s: %s (Versions: %s)\n", issue.Field, issue.Description, formatVersions(issue.Versions))
				}
			}
			
			// Then warnings
			if len(warningIssues) > 0 {
				fmt.Printf("\nWarnings (%d):\n", len(warningIssues))
				for _, issue := range warningIssues {
					fmt.Printf("  - %s: %s (Versions: %s)\n", issue.Field, issue.Description, formatVersions(issue.Versions))
				}
			}
			
			// Then info
			if len(infoIssues) > 0 {
				fmt.Printf("\nInfo (%d):\n", len(infoIssues))
				for _, issue := range infoIssues {
					fmt.Printf("  - %s: %s (Versions: %s)\n", issue.Field, issue.Description, formatVersions(issue.Versions))
				}
			}
		} else {
			fmt.Println("\nNo issues found.")
		}
	}
}

// formatVersions formats a slice of versions as a comma-separated string
func formatVersions(versions []int64) string {
	versionStrs := make([]string, len(versions))
	for i, v := range versions {
		versionStrs[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(versionStrs, ", ")
}
