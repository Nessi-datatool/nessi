package dbt

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
)

// OutputJSON outputs results in JSON format
func OutputJSON(results interface{}, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}

// OutputCSV outputs results in CSV format
func OutputCSV(results interface{}, w io.Writer) error {
	csvWriter := csv.NewWriter(w)
	defer csvWriter.Flush()

	// Handle different result types
	switch r := results.(type) {
	case *ValidationResults:
		return outputValidationResultsCSV(r, csvWriter)
	case *ProfileResults:
		return outputProfileResultsCSV(r, csvWriter)
	default:
		return fmt.Errorf("unsupported result type: %T", results)
	}
}

// OutputTable outputs results in table format
func OutputTable(results interface{}, w io.Writer) error {
	// Handle different result types
	switch r := results.(type) {
	case *ValidationResults:
		return outputValidationResultsTable(r, w)
	case *ProfileResults:
		return outputProfileResultsTable(r, w)
	default:
		return fmt.Errorf("unsupported result type: %T", results)
	}
}

// outputValidationResultsCSV outputs validation results in CSV format
func outputValidationResultsCSV(results *ValidationResults, csvWriter *csv.Writer) error {
	// Write header
	header := []string{"Model", "Rule", "Status", "Message", "Table Path", "Failure Count"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write results
	for _, result := range results.Results {
		row := []string{
			result.ModelName,
			result.RuleName,
			result.Status,
			result.Message,
			result.TablePath,
			fmt.Sprintf("%d", result.FailureCount),
		}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	// Write summary
	csvWriter.Write([]string{})
	csvWriter.Write([]string{"Summary"})
	csvWriter.Write([]string{"Total Models", fmt.Sprintf("%d", results.Summary.TotalModels)})
	csvWriter.Write([]string{"Passed Models", fmt.Sprintf("%d", results.Summary.PassedModels)})
	csvWriter.Write([]string{"Failed Models", fmt.Sprintf("%d", results.Summary.FailedModels)})
	csvWriter.Write([]string{"Total Rules", fmt.Sprintf("%d", results.Summary.TotalRules)})
	csvWriter.Write([]string{"Passed Rules", fmt.Sprintf("%d", results.Summary.PassedRules)})
	csvWriter.Write([]string{"Failed Rules", fmt.Sprintf("%d", results.Summary.FailedRules)})
	csvWriter.Write([]string{"Quality Score", fmt.Sprintf("%.2f%%", results.Summary.QualityScore)})
	csvWriter.Write([]string{"Execution Time", fmt.Sprintf("%.2fs", results.Summary.ExecutionTime)})

	return nil
}

// outputProfileResultsCSV outputs profile results in CSV format
func outputProfileResultsCSV(results *ProfileResults, csvWriter *csv.Writer) error {
	// Write header
	header := []string{"Model", "Table Path", "Profile Type", "Column", "Data Type", "Non-Null Count", "Null Count", "Null %", "Unique", "Unique %"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write results
	for _, result := range results.Results {
		for column, stats := range result.ColumnProfiles {
			row := []string{
				result.ModelName,
				result.TablePath,
				result.ProfileType,
				column,
				stats.DataType,
				fmt.Sprintf("%d", stats.NonNullCount),
				fmt.Sprintf("%d", stats.NullCount),
				fmt.Sprintf("%.2f%%", stats.NullPercent),
				fmt.Sprintf("%d", stats.Unique),
				fmt.Sprintf("%.2f%%", stats.UniquePercent),
			}
			if err := csvWriter.Write(row); err != nil {
				return err
			}
		}
	}

	// Write summary
	csvWriter.Write([]string{})
	csvWriter.Write([]string{"Summary"})
	csvWriter.Write([]string{"Total Models", fmt.Sprintf("%d", results.Summary.TotalModels)})
	csvWriter.Write([]string{"Total Columns", fmt.Sprintf("%d", results.Summary.TotalColumns)})
	csvWriter.Write([]string{"Total Rows", fmt.Sprintf("%d", results.Summary.TotalRows)})
	csvWriter.Write([]string{"Execution Time", fmt.Sprintf("%.2fs", results.Summary.ExecutionTime)})

	return nil
}

// outputValidationResultsTable outputs validation results in table format
func outputValidationResultsTable(results *ValidationResults, w io.Writer) error {
	// Create tabwriter for aligned output
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Write header
	fmt.Fprintln(tw, "MODEL\tRULE\tSTATUS\tMESSAGE\tTABLE PATH\tFAILURE COUNT")
	fmt.Fprintln(tw, "-----\t----\t------\t-------\t----------\t-------------")

	// Write results
	for _, result := range results.Results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%d\n",
			result.ModelName,
			result.RuleName,
			result.Status,
			result.Message,
			result.TablePath,
			result.FailureCount,
		)
	}

	// Write summary
	fmt.Fprintln(tw)
	fmt.Fprintln(tw, "SUMMARY")
	fmt.Fprintln(tw, "-------")
	fmt.Fprintf(tw, "Total Models:\t%d\n", results.Summary.TotalModels)
	fmt.Fprintf(tw, "Passed Models:\t%d\n", results.Summary.PassedModels)
	fmt.Fprintf(tw, "Failed Models:\t%d\n", results.Summary.FailedModels)
	fmt.Fprintf(tw, "Total Rules:\t%d\n", results.Summary.TotalRules)
	fmt.Fprintf(tw, "Passed Rules:\t%d\n", results.Summary.PassedRules)
	fmt.Fprintf(tw, "Failed Rules:\t%d\n", results.Summary.FailedRules)
	fmt.Fprintf(tw, "Quality Score:\t%.2f%%\n", results.Summary.QualityScore)
	fmt.Fprintf(tw, "Execution Time:\t%.2fs\n", results.Summary.ExecutionTime)

	return tw.Flush()
}

// outputProfileResultsTable outputs profile results in table format
func outputProfileResultsTable(results *ProfileResults, w io.Writer) error {
	// Create tabwriter for aligned output
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Write header
	fmt.Fprintln(tw, "MODEL\tCOLUMN\tDATA TYPE\tNON-NULL COUNT\tNULL COUNT\tNULL %\tUNIQUE\tUNIQUE %")
	fmt.Fprintln(tw, "-----\t------\t---------\t-------------\t---------\t------\t------\t--------")

	// Write results
	for _, result := range results.Results {
		// Print model header
		fmt.Fprintf(tw, "Model: %s\tTable: %s\tProfile Type: %s\tRow Count: %d\n",
			result.ModelName,
			result.TablePath,
			result.ProfileType,
			result.RowCount,
		)
		
		// Print column stats
		for column, stats := range result.ColumnProfiles {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%.2f%%\t%d\t%.2f%%\n",
				result.ModelName,
				column,
				stats.DataType,
				stats.NonNullCount,
				stats.NullCount,
				stats.NullPercent,
				stats.Unique,
				stats.UniquePercent,
			)
			
			// Print additional stats for numeric columns
			if stats.Min != nil && stats.Max != nil {
				fmt.Fprintf(tw, "\tMin: %v\tMax: %v", stats.Min, stats.Max)
				if stats.Mean != 0 {
					fmt.Fprintf(tw, "\tMean: %.2f\tStddev: %.2f", stats.Mean, stats.Stddev)
				}
				fmt.Fprintln(tw)
			}
			
			// Print top values if available
			if len(stats.TopValues) > 0 {
				fmt.Fprintln(tw, "\tTop Values:")
				for _, vf := range stats.TopValues {
					fmt.Fprintf(tw, "\t\t%v: %d (%.2f%%)\n", vf.Value, vf.Frequency, vf.Percent)
				}
			}
		}
		
		fmt.Fprintln(tw)
	}

	// Write summary
	fmt.Fprintln(tw, "SUMMARY")
	fmt.Fprintln(tw, "-------")
	fmt.Fprintf(tw, "Total Models:\t%d\n", results.Summary.TotalModels)
	fmt.Fprintf(tw, "Total Columns:\t%d\n", results.Summary.TotalColumns)
	fmt.Fprintf(tw, "Total Rows:\t%d\n", results.Summary.TotalRows)
	fmt.Fprintf(tw, "Execution Time:\t%.2fs\n", results.Summary.ExecutionTime)

	return tw.Flush()
}

// GenerateDBTDocsArtifacts generates artifacts for dbt Docs integration
func GenerateDBTDocsArtifacts(results interface{}, artifactsPath string) error {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Generate artifacts in a format compatible with dbt Docs
	// 2. Write them to the specified path

	// For now, just return nil
	return nil
}
