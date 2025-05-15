package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/datalake"
	"github.com/spf13/cobra"
)

func init() {
	// Time travel command
	timeTravelCmd := &cobra.Command{
		Use:   "timetravel",
		Short: "Query Delta table as it existed at a specific version or timestamp",
		Long:  "Time travel allows you to query a Delta table as it existed at a specific version or point in time",
	}

	// Query by version
	versionCmd := &cobra.Command{
		Use:   "version [table_path] [version]",
		Short: "Query Delta table as it existed at a specific version",
		Long:  "Query a Delta table as it existed at a specific version number",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid table path: %w", err)
			}

			// Parse version
			var version int64
			_, err = fmt.Sscanf(args[1], "%d", &version)
			if err != nil {
				return fmt.Errorf("invalid version: %w", err)
			}

			// Get output format
			outputFormat, _ := cmd.Flags().GetString("format")

			manager := datalake.NewMetadataManager(tablePath)
			result, err := manager.TimeTravel(datalake.DeltaTimeTravelOptions{Version: &version})
			if err != nil {
				return fmt.Errorf("time travel failed: %w", err)
			}

			// Output the result
			printTimeTravelResult(result, outputFormat)
			return nil
		},
	}

	// Query by timestamp
	timestampCmd := &cobra.Command{
		Use:   "timestamp [table_path] [timestamp]",
		Short: "Query Delta table as it existed at a specific timestamp",
		Long:  "Query a Delta table as it existed at a specific point in time (format: YYYY-MM-DDTHH:MM:SS)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tablePath, err := filepath.Abs(args[0])
			if err != nil {
				return fmt.Errorf("invalid table path: %w", err)
			}

			// Parse timestamp
			timestamp, err := time.Parse(time.RFC3339, args[1])
			if err != nil {
				return fmt.Errorf("invalid timestamp format, use YYYY-MM-DDTHH:MM:SS: %w", err)
			}

			// Get output format
			outputFormat, _ := cmd.Flags().GetString("format")

			manager := datalake.NewMetadataManager(tablePath)
			result, err := manager.TimeTravel(datalake.DeltaTimeTravelOptions{Timestamp: &timestamp})
			if err != nil {
				return fmt.Errorf("time travel failed: %w", err)
			}

			// Output the result
			printTimeTravelResult(result, outputFormat)
			return nil
		},
	}

	// Add format flag to both commands
	versionCmd.Flags().StringP("format", "f", "text", "Output format (text, json)")
	timestampCmd.Flags().StringP("format", "f", "text", "Output format (text, json)")

	timeTravelCmd.AddCommand(versionCmd, timestampCmd)
	rootCmd.AddCommand(timeTravelCmd)
}

// printTimeTravelResult prints the time travel result in the specified format
func printTimeTravelResult(result *datalake.DeltaTimeTravelResult, format string) {
	switch format {
	case "json":
		// Create a simplified result for JSON output
		jsonResult := struct {
			Version   int64     `json:"version"`
			Timestamp time.Time `json:"timestamp"`
			Schema    []struct {
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"schema"`
			Files []string `json:"files"`
		}{
			Version:   result.Version,
			Timestamp: result.Timestamp,
			Files:     result.Files,
		}

		// Convert schema fields
		jsonResult.Schema = make([]struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}, len(result.SchemaFields))

		for i, field := range result.SchemaFields {
			jsonResult.Schema[i] = struct {
				Name string `json:"name"`
				Type string `json:"type"`
			}{
				Name: field.Name,
				Type: field.Type,
			}
		}

		// Marshal to JSON
		jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
			return
		}

		fmt.Println(string(jsonData))

	default: // text format
		fmt.Printf("Time Travel Result:\n")
		fmt.Printf("Version: %d\n", result.Version)
		fmt.Printf("Timestamp: %s\n", result.Timestamp.Format(time.RFC3339))
		fmt.Printf("\nSchema (%d fields):\n", len(result.SchemaFields))
		
		for _, field := range result.SchemaFields {
			fmt.Printf("  %s: %s\n", field.Name, field.Type)
		}

		fmt.Printf("\nFiles (%d):\n", len(result.Files))
		for i, file := range result.Files {
			if i < 10 || i >= len(result.Files)-5 { // Show first 10 and last 5 files
				fmt.Printf("  %s\n", file)
			} else if i == 10 {
				fmt.Printf("  ... (%d more files) ...\n", len(result.Files)-15)
			}
		}
	}
}
