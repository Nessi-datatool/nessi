package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BatchConfig represents the configuration for a batch job
type BatchConfig struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tables      []string          `json:"tables"`
	Rules       string            `json:"rules"`
	Output      string            `json:"output"`
	Format      string            `json:"format"`
	Schedule    string            `json:"schedule,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
}

// BatchResult represents the result of a batch job
type BatchResult struct {
	BatchName   string    `json:"batch_name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    string    `json:"duration"`
	TablesTotal int       `json:"tables_total"`
	TablesPassed int      `json:"tables_passed"`
	TablesFailed int      `json:"tables_failed"`
	Results     []string  `json:"results"`
}

var (
	// Batch command flags
	batchConfigFile string
	batchOutputDir  string
	batchParallel   int
	batchTags       []string
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch [config_file]",
	Short: "Run batch data quality checks",
	Long: `The batch command runs data quality checks on multiple tables using a batch configuration file.
	
Example:
  nessi batch config.json --output ./results --parallel 4
  nessi batch --tags production,daily`,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the config file
		configFile := ""
		if len(args) > 0 {
			configFile = args[0]
		} else if batchConfigFile != "" {
			configFile = batchConfigFile
		} else {
			// Look for default batch config files
			defaultFiles := []string{"batch.json", "batch.yaml", "batch.yml"}
			for _, file := range defaultFiles {
				if _, err := os.Stat(file); err == nil {
					configFile = file
					break
				}
			}
		}

		if configFile == "" {
			fmt.Println("Error: No batch configuration file specified")
			os.Exit(1)
		}

		// Load the batch configuration
		config, err := loadBatchConfig(configFile)
		if err != nil {
			fmt.Printf("Error loading batch configuration: %v\n", err)
			os.Exit(1)
		}

		// Filter by tags if specified
		if len(batchTags) > 0 && len(config.Tags) > 0 {
			// Check if any of the specified tags match the config tags
			match := false
			for _, tag := range batchTags {
				for _, configTag := range config.Tags {
					if tag == configTag {
						match = true
						break
					}
				}
				if match {
					break
				}
			}
			if !match {
				fmt.Printf("Skipping batch %s: no matching tags\n", config.Name)
				return
			}
		}

		// Create output directory if it doesn't exist
		outputDir := batchOutputDir
		if outputDir == "" {
			outputDir = config.Output
		}
		if outputDir == "" {
			outputDir = "batch_results"
		}
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		// Run the batch job
		fmt.Printf("Running batch job: %s\n", config.Name)
		fmt.Printf("Description: %s\n", config.Description)
		fmt.Printf("Tables: %d\n", len(config.Tables))

		startTime := time.Now()
		results := &BatchResult{
			BatchName:    config.Name,
			StartTime:    startTime,
			TablesTotal:  len(config.Tables),
			TablesPassed: 0,
			TablesFailed: 0,
			Results:      make([]string, 0),
		}

		// Process each table
		for _, table := range config.Tables {
			fmt.Printf("Processing table: %s\n", table)
			
			// Run the check command for this table
			outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.json", filepath.Base(table)))
			
			// Build the command arguments
			args := []string{
				"check",
				table,
				"--rules", config.Rules,
				"--output", "json",
				"--output-file", outputFile,
			}
			
			// Add any custom parameters
			for k, v := range config.Params {
				args = append(args, fmt.Sprintf("--%s", k), v)
			}
			
			// Create a new root command for execution
			cmd := &cobra.Command{Use: "nessi"}
			
			// Add the check command
			cmd.AddCommand(checkCmd)
			
			// Set the arguments
			cmd.SetArgs(args)
			
			// Execute the command
			err := cmd.Execute()
			if err != nil {
				fmt.Printf("Error processing table %s: %v\n", table, err)
				results.TablesFailed++
			} else {
				fmt.Printf("Successfully processed table: %s\n", table)
				results.TablesPassed++
			}
			
			results.Results = append(results.Results, outputFile)
		}

		// Complete the results
		endTime := time.Now()
		results.EndTime = endTime
		results.Duration = endTime.Sub(startTime).String()

		// Save the batch results
		resultsFile := filepath.Join(outputDir, "batch_summary.json")
		resultsJSON, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Error creating results JSON: %v\n", err)
		} else {
			if err := os.WriteFile(resultsFile, resultsJSON, 0644); err != nil {
				fmt.Printf("Error writing results file: %v\n", err)
			}
		}

		// Print summary
		fmt.Println("\nBatch Job Summary:")
		fmt.Printf("Name: %s\n", results.BatchName)
		fmt.Printf("Duration: %s\n", results.Duration)
		fmt.Printf("Tables Processed: %d\n", results.TablesTotal)
		fmt.Printf("Tables Passed: %d\n", results.TablesPassed)
		fmt.Printf("Tables Failed: %d\n", results.TablesFailed)
		fmt.Printf("Results saved to: %s\n", resultsFile)
	},
}

// loadBatchConfig loads a batch configuration from a file
func loadBatchConfig(file string) (*BatchConfig, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config BatchConfig
	ext := strings.ToLower(filepath.Ext(file))
	if ext == ".json" {
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		// For YAML, we'll use Viper
		viper := viper.New()
		viper.SetConfigType(strings.TrimPrefix(ext, "."))
		if err := viper.ReadConfig(strings.NewReader(string(data))); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}
	} else {
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	return &config, nil
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// Add flags to the batch command
	batchCmd.Flags().StringVar(&batchConfigFile, "config", "", "Path to batch configuration file")
	batchCmd.Flags().StringVar(&batchOutputDir, "output", "", "Directory to store batch results")
	batchCmd.Flags().IntVar(&batchParallel, "parallel", 1, "Number of tables to process in parallel")
	batchCmd.Flags().StringSliceVar(&batchTags, "tags", []string{}, "Filter batch jobs by tags")
}
