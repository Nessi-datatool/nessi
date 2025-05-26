package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newWorkflowCmd() *cobra.Command {
	workflowCmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage data workflows",
		Long:  `Commands for managing and executing data workflows.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Add subcommands
	workflowCmd.AddCommand(newWorkflowRunCmd())
	workflowCmd.AddCommand(newWorkflowListCmd())
	workflowCmd.AddCommand(newWorkflowCreateCmd())

	return workflowCmd
}

func newWorkflowRunCmd() *cobra.Command {
	var workflowFile string

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run a workflow",
		Long:  `Run a workflow defined in a YAML file.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if workflow file exists
			if _, err := os.Stat(workflowFile); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error N101: Workflow file %s does not exist\n", workflowFile)
				os.Exit(1)
			}

			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Workflow Orchestration")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	// Add flags
	runCmd.Flags().StringVar(&workflowFile, "file", "", "Path to workflow YAML file")
	runCmd.MarkFlagRequired("file")

	return runCmd
}

func newWorkflowListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available workflows",
		Long:  `List all available workflow templates.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Workflow Orchestration")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	return listCmd
}

func newWorkflowCreateCmd() *cobra.Command {
	var outputFile string
	var template string

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workflow",
		Long:  `Create a new workflow from a template or from scratch.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Check if this is a Pro feature
			fmt.Println("Checking license for Pro feature: Workflow Orchestration")
			fmt.Println("This is a Pro Edition feature.")
			fmt.Println("Please activate a trial with 'nessi license trial' or purchase a Pro license.")
			os.Exit(1)
		},
	}

	// Add flags
	createCmd.Flags().StringVar(&outputFile, "output", "", "Output file path")
	createCmd.Flags().StringVar(&template, "template", "basic", "Workflow template to use (basic, quality, etl)")
	createCmd.MarkFlagRequired("output")

	return createCmd
}
