package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
)

func init() {
	// Add format validation command
	var formatCmd = &cobra.Command{
		Use:   "format",
		Short: "Validate format of fields in a Delta table",
		Long:  "Validate format of fields in a Delta table, including date, email, URL, and custom regex patterns",
	}
	rootCmd.AddCommand(formatCmd)

	// Add date format validation command
	var dateCmd = &cobra.Command{
		Use:   "date [table_path] [field]",
		Short: "Validate date format of a field",
		Long:  "Validate date format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			dateFormat, _ := cmd.Flags().GetString("format")
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.DateFormat,
				Field:            field,
				DateFormat:       dateFormat,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	dateCmd.Flags().StringP("format", "f", time.RFC3339, "Date format to validate (default: RFC3339)")
	dateCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	dateCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(dateCmd)

	// Add email format validation command
	var emailCmd = &cobra.Command{
		Use:   "email [table_path] [field]",
		Short: "Validate email format of a field",
		Long:  "Validate email format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.EmailFormat,
				Field:            field,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	emailCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	emailCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(emailCmd)

	// Add URL format validation command
	var urlCmd = &cobra.Command{
		Use:   "url [table_path] [field]",
		Short: "Validate URL format of a field",
		Long:  "Validate URL format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.URLFormat,
				Field:            field,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	urlCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	urlCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(urlCmd)

	// Add digit-only format validation command
	var digitCmd = &cobra.Command{
		Use:   "digit [table_path] [field]",
		Short: "Validate digit-only format of a field",
		Long:  "Validate digit-only format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.DigitOnlyFormat,
				Field:            field,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	digitCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	digitCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(digitCmd)

	// Add UUID format validation command
	var uuidCmd = &cobra.Command{
		Use:   "uuid [table_path] [field]",
		Short: "Validate UUID format of a field",
		Long:  "Validate UUID format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.UUIDFormat,
				Field:            field,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	uuidCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	uuidCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(uuidCmd)

	// Add IP address format validation command
	var ipCmd = &cobra.Command{
		Use:   "ip [table_path] [field]",
		Short: "Validate IP address format of a field",
		Long:  "Validate IP address format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.IPAddressFormat,
				Field:            field,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	ipCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	ipCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(ipCmd)

	// Add phone number format validation command
	var phoneCmd = &cobra.Command{
		Use:   "phone [table_path] [field]",
		Short: "Validate phone number format of a field",
		Long:  "Validate phone number format of a field in a Delta table",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.PhoneNumberFormat,
				Field:            field,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	phoneCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	phoneCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(phoneCmd)

	// Add custom regex format validation command
	var regexCmd = &cobra.Command{
		Use:   "regex [table_path] [field] [pattern]",
		Short: "Validate custom regex pattern of a field",
		Long:  "Validate custom regex pattern of a field in a Delta table",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath := args[0]
			field := args[1]
			pattern := args[2]
			maxInvalid, _ := cmd.Flags().GetInt("max-invalid")
			outputFormat, _ := cmd.Flags().GetString("output")

			options := datalake.FormatValidationOptions{
				FormatType:       datalake.CustomRegexFormat,
				Field:            field,
				CustomRegex:      pattern,
				MaxInvalidValues: maxInvalid,
			}

			return runFormatValidation(tablePath, options, outputFormat)
		},
	}
	regexCmd.Flags().IntP("max-invalid", "m", 10, "Maximum number of invalid values to return")
	regexCmd.Flags().StringP("output", "o", "text", "Output format (text, json)")
	formatCmd.AddCommand(regexCmd)
}

// runFormatValidation runs a format validation and outputs the result
func runFormatValidation(tablePath string, options datalake.FormatValidationOptions, outputFormat string) error {
	// Create metadata manager
	manager := datalake.NewMetadataManager(tablePath)
	
	// Run format validation
	result, err := manager.ValidateFormat(options)
	if err != nil {
		return fmt.Errorf("format validation failed: %w", err)
	}
	
	// Output the result
	printFormatValidationResult(result, outputFormat)
	return nil
}

// printFormatValidationResult prints the format validation result in the specified format
func printFormatValidationResult(result *datalake.FormatValidationResult, format string) {
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
		fmt.Printf("Format Validation Result: %s\n", result.FormatType)
		fmt.Printf("Field: %s\n", result.Field)
		fmt.Printf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
		fmt.Printf("Passed: %t\n", result.Passed)
		fmt.Printf("Total Values: %d\n", result.TotalValues)
		fmt.Printf("Valid Values: %d\n", result.ValidValues)
		fmt.Printf("Invalid Values: %d\n", result.InvalidValues)
		
		if len(result.InvalidExamples) > 0 {
			fmt.Printf("\nInvalid Examples (%d):\n", len(result.InvalidExamples))
			for i, example := range result.InvalidExamples {
				fmt.Printf("  %d. %s\n", i+1, example)
			}
		}
	}
}
