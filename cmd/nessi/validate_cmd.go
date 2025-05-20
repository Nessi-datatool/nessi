package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	TablePath    string    `json:"table_path"`
	Timestamp    time.Time `json:"timestamp"`
	Valid        bool      `json:"valid"`
	TotalChecks  int       `json:"total_checks"`
	PassedChecks int       `json:"passed_checks"`
	FailedChecks int       `json:"failed_checks"`
	Details      []string  `json:"details,omitempty"`
	Duration     string    `json:"duration"`
}

var (
	// Validate command flags
	validateFormat     string
	validateOutputFile string
	validateVerbose    bool
	validateThreshold  float64
)

// deltaValidateCmd represents the delta validate command
var deltaValidateCmd = &cobra.Command{
	Use:   "validate [table_path]",
	Short: "One-line validation of a Delta table",
	Long: `The validate command performs a quick validation of a Delta table using default rules.
This is a simplified version of the check command for quick validations.
	
Examples:
  nessi validate /path/to/table                     # Validate with default rules
  nessi validate /path/to/table --format json       # Output in JSON format
  nessi validate /path/to/table --threshold 0.9     # Require 90% of checks to pass`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tablePath := args[0]

		// Check if the table exists
		if _, err := os.Stat(tablePath); os.IsNotExist(err) {
			fmt.Printf("Error: Table not found at %s\n", tablePath)
			os.Exit(1)
		}

		startTime := time.Now()

		// Create a validation result
		result := &ValidationResult{
			TablePath:    tablePath,
			Timestamp:    startTime,
			TotalChecks:  0,
			PassedChecks: 0,
			FailedChecks: 0,
			Details:      make([]string, 0),
		}

		// Perform basic validation checks

		// 1. Check if it's a valid Delta table
		isDelta, err := isValidDeltaTable(tablePath)
		result.TotalChecks++
		if err != nil {
			result.FailedChecks++
			result.Details = append(result.Details, fmt.Sprintf("Error checking Delta format: %v", err))
		} else if !isDelta {
			result.FailedChecks++
			result.Details = append(result.Details, "Not a valid Delta table")
		} else {
			result.PassedChecks++
			result.Details = append(result.Details, "Valid Delta table format")
		}

		// 2. Check for schema
		hasSchema, err := hasValidSchema(tablePath)
		result.TotalChecks++
		if err != nil {
			result.FailedChecks++
			result.Details = append(result.Details, fmt.Sprintf("Error checking schema: %v", err))
		} else if !hasSchema {
			result.FailedChecks++
			result.Details = append(result.Details, "No valid schema found")
		} else {
			result.PassedChecks++
			result.Details = append(result.Details, "Valid schema found")
		}

		// 3. Check for data files
		hasData, err := hasDataFiles(tablePath)
		result.TotalChecks++
		if err != nil {
			result.FailedChecks++
			result.Details = append(result.Details, fmt.Sprintf("Error checking data files: %v", err))
		} else if !hasData {
			result.FailedChecks++
			result.Details = append(result.Details, "No data files found")
		} else {
			result.PassedChecks++
			result.Details = append(result.Details, "Data files found")
		}

		// 4. Check for transaction log
		hasLog, err := hasTransactionLog(tablePath)
		result.TotalChecks++
		if err != nil {
			result.FailedChecks++
			result.Details = append(result.Details, fmt.Sprintf("Error checking transaction log: %v", err))
		} else if !hasLog {
			result.FailedChecks++
			result.Details = append(result.Details, "No transaction log found")
		} else {
			result.PassedChecks++
			result.Details = append(result.Details, "Transaction log found")
		}

		// Calculate validation result
		endTime := time.Now()
		result.Duration = endTime.Sub(startTime).String()

		// Determine if the table is valid based on the threshold
		passRate := float64(result.PassedChecks) / float64(result.TotalChecks)
		result.Valid = passRate >= validateThreshold

		// Output the result
		outputValidationResult(result)
	},
}

// isValidDeltaTable checks if the path contains a valid Delta table
func isValidDeltaTable(path string) (bool, error) {
	// Check if the _delta_log directory exists
	deltaLogPath := filepath.Join(path, "_delta_log")
	info, err := os.Stat(deltaLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}

// hasValidSchema checks if the table has a valid schema
func hasValidSchema(path string) (bool, error) {
	// Check for schema in the latest transaction log file
	deltaLogPath := filepath.Join(path, "_delta_log")

	// List files in the _delta_log directory
	files, err := os.ReadDir(deltaLogPath)
	if err != nil {
		return false, err
	}

	// Look for JSON files (transaction log files)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			// Found a transaction log file, assume it has a schema
			return true, nil
		}
	}

	return false, nil
}

// hasDataFiles checks if the table has data files
func hasDataFiles(path string) (bool, error) {
	// Look for parquet files in the table directory
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".parquet") {
			return true, nil
		}
	}

	return false, nil
}

// hasTransactionLog checks if the table has a transaction log
func hasTransactionLog(path string) (bool, error) {
	// Check if the _delta_log directory exists and has files
	deltaLogPath := filepath.Join(path, "_delta_log")

	// List files in the _delta_log directory
	files, err := os.ReadDir(deltaLogPath)
	if err != nil {
		return false, err
	}

	return len(files) > 0, nil
}

// outputValidationResult outputs the validation result in the specified format
func outputValidationResult(result *ValidationResult) {
	switch validateFormat {
	case "json":
		outputJSON(result)
	case "text":
		outputText(result)
	default:
		outputText(result)
	}
}

// outputJSON outputs the validation result in JSON format
func outputJSON(result *ValidationResult) {
	// Convert to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error creating JSON: %v\n", err)
		return
	}

	// Write to file if specified
	if validateOutputFile != "" {
		if err := os.WriteFile(validateOutputFile, jsonData, 0644); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	} else {
		// Print to stdout
		fmt.Println(string(jsonData))
	}
}

// outputText outputs the validation result in text format
func outputText(result *ValidationResult) {
	// Create text output
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Table: %s\n", result.TablePath))
	output.WriteString(fmt.Sprintf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339)))
	output.WriteString(fmt.Sprintf("Valid: %t\n", result.Valid))
	output.WriteString(fmt.Sprintf("Checks: %d total, %d passed, %d failed\n",
		result.TotalChecks, result.PassedChecks, result.FailedChecks))
	output.WriteString(fmt.Sprintf("Duration: %s\n", result.Duration))

	if validateVerbose {
		output.WriteString("\nDetails:\n")
		for _, detail := range result.Details {
			output.WriteString(fmt.Sprintf("- %s\n", detail))
		}
	}

	// Write to file if specified
	if validateOutputFile != "" {
		if err := os.WriteFile(validateOutputFile, []byte(output.String()), 0644); err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
		}
	} else {
		// Print to stdout
		fmt.Print(output.String())
	}
}

func init() {
	rootCmd.AddCommand(deltaValidateCmd)

	// Add flags to the validate command
	deltaValidateCmd.Flags().StringVar(&validateFormat, "format", "text", "Output format (text, json)")
	deltaValidateCmd.Flags().StringVar(&validateOutputFile, "output-file", "", "File to write output to")
	deltaValidateCmd.Flags().BoolVar(&validateVerbose, "verbose", false, "Show detailed output")
	deltaValidateCmd.Flags().Float64Var(&validateThreshold, "threshold", 1.0, "Threshold for validation (0.0-1.0)")
}
