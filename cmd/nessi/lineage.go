package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/catalog"
	"github.com/nessi-dev/nessi-dev/pkg/lineage"
	"github.com/nessi-dev/nessi-dev/pkg/lineage/model"
	"github.com/nessi-dev/nessi-dev/pkg/lineage/visualization"
	"github.com/spf13/cobra"
)

// initLineageCmd initializes the lineage command
func initLineageCmd(rootCmd *cobra.Command) {
	// Create lineage command
	lineageCmd := &cobra.Command{
		Use:   "lineage",
		Short: "Manage data lineage",
		Long:  `Manage data lineage graphs, snapshots, and visualizations.`,
	}

	// Add lineage command to root command
	rootCmd.AddCommand(lineageCmd)

	// Add subcommands
	initLineageGraphCmd(lineageCmd)
	initLineageSnapshotCmd(lineageCmd)
	initLineageVisualizeCmd(lineageCmd)
	initLineageCompareCmd(lineageCmd)
	initLineageImportCmd(lineageCmd)
	initLineageExportCmd(lineageCmd)
}

// initLineageGraphCmd initializes the lineage graph command
func initLineageGraphCmd(parentCmd *cobra.Command) {
	// Create graph command
	graphCmd := &cobra.Command{
		Use:   "graph",
		Short: "Manage lineage graphs",
		Long:  `Create, list, get, and delete lineage graphs.`,
	}

	// Add graph command to parent command
	parentCmd.AddCommand(graphCmd)

	// Create graph list command
	graphListCmd := &cobra.Command{
		Use:   "list",
		Short: "List lineage graphs",
		Long:  `List all lineage graphs.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphs, err := service.ListGraphs(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing graphs: %v\n", err)
				os.Exit(1)
			}

			if len(graphs) == 0 {
				fmt.Println("No lineage graphs found.")
				return
			}

			fmt.Println("Lineage Graphs:")
			fmt.Println("ID\tName\tDescription\tCreated\tUpdated")
			for _, graph := range graphs {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
					graph.ID,
					graph.Name,
					graph.Description,
					graph.CreatedAt.Format(time.RFC3339),
					graph.UpdatedAt.Format(time.RFC3339))
			}
		},
	}

	// Create graph create command
	graphCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a lineage graph",
		Long:  `Create a new lineage graph.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			name, _ := cmd.Flags().GetString("name")
			description, _ := cmd.Flags().GetString("description")

			graph, err := service.CreateGraph(ctx, name, description)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating graph: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Created lineage graph: %s\n", graph.ID)
		},
	}
	graphCreateCmd.Flags().String("name", "", "Name of the lineage graph")
	graphCreateCmd.Flags().String("description", "", "Description of the lineage graph")
	graphCreateCmd.MarkFlagRequired("name")

	// Create graph get command
	graphGetCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a lineage graph",
		Long:  `Get a lineage graph by ID.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			id := args[0]
			graph, err := service.GetGraph(ctx, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting graph: %v\n", err)
				os.Exit(1)
			}

			if graph == nil {
				fmt.Fprintf(os.Stderr, "Graph not found: %s\n", id)
				os.Exit(1)
			}

			fmt.Printf("ID: %s\n", graph.ID)
			fmt.Printf("Name: %s\n", graph.Name)
			fmt.Printf("Description: %s\n", graph.Description)
			fmt.Printf("Created: %s\n", graph.CreatedAt.Format(time.RFC3339))
			fmt.Printf("Updated: %s\n", graph.UpdatedAt.Format(time.RFC3339))
			fmt.Printf("Nodes: %d\n", len(graph.Nodes))
			fmt.Printf("Edges: %d\n", len(graph.Edges))

			format, _ := cmd.Flags().GetString("format")
			if format == "json" {
				jsonData, err := json.MarshalIndent(graph, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling graph to JSON: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(string(jsonData))
			}
		},
	}
	graphGetCmd.Flags().String("format", "", "Output format (json)")

	// Create graph delete command
	graphDeleteCmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a lineage graph",
		Long:  `Delete a lineage graph by ID.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			id := args[0]
			err := service.DeleteGraph(ctx, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting graph: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Deleted lineage graph: %s\n", id)
		},
	}

	// Create graph add-node command
	graphAddNodeCmd := &cobra.Command{
		Use:   "add-node [graph-id]",
		Short: "Add a node to a lineage graph",
		Long:  `Add a node to a lineage graph.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			name, _ := cmd.Flags().GetString("name")
			nodeType, _ := cmd.Flags().GetString("type")
			description, _ := cmd.Flags().GetString("description")

			node := model.NewNode(name, model.NodeType(nodeType))
			node.Description = description

			if err := service.AddNode(ctx, graphID, node); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding node: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Added node %s to graph %s\n", node.ID, graphID)
		},
	}
	graphAddNodeCmd.Flags().String("name", "", "Name of the node")
	graphAddNodeCmd.Flags().String("type", "", "Type of the node (table, view, query, process, dataset, api, application)")
	graphAddNodeCmd.Flags().String("description", "", "Description of the node")
	graphAddNodeCmd.MarkFlagRequired("name")
	graphAddNodeCmd.MarkFlagRequired("type")

	// Create graph add-edge command
	graphAddEdgeCmd := &cobra.Command{
		Use:   "add-edge [graph-id]",
		Short: "Add an edge to a lineage graph",
		Long:  `Add an edge to a lineage graph.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			sourceID, _ := cmd.Flags().GetString("source")
			targetID, _ := cmd.Flags().GetString("target")
			edgeType, _ := cmd.Flags().GetString("type")

			edge := model.NewEdge(sourceID, targetID, model.EdgeType(edgeType))

			if err := service.AddEdge(ctx, graphID, edge); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding edge: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Added edge from %s to %s in graph %s\n", sourceID, targetID, graphID)
		},
	}
	graphAddEdgeCmd.Flags().String("source", "", "ID of the source node")
	graphAddEdgeCmd.Flags().String("target", "", "ID of the target node")
	graphAddEdgeCmd.Flags().String("type", "", "Type of the edge (read, write, depends_on, produces)")
	graphAddEdgeCmd.MarkFlagRequired("source")
	graphAddEdgeCmd.MarkFlagRequired("target")
	graphAddEdgeCmd.MarkFlagRequired("type")

	// Add subcommands to graph command
	graphCmd.AddCommand(graphListCmd)
	graphCmd.AddCommand(graphCreateCmd)
	graphCmd.AddCommand(graphGetCmd)
	graphCmd.AddCommand(graphDeleteCmd)
	graphCmd.AddCommand(graphAddNodeCmd)
	graphCmd.AddCommand(graphAddEdgeCmd)
}

// initLineageSnapshotCmd initializes the lineage snapshot command
func initLineageSnapshotCmd(parentCmd *cobra.Command) {
	// Create snapshot command
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage lineage snapshots",
		Long:  `Create, list, and get lineage snapshots.`,
	}

	// Add snapshot command to parent command
	parentCmd.AddCommand(snapshotCmd)

	// Create snapshot list command
	snapshotListCmd := &cobra.Command{
		Use:   "list [graph-id]",
		Short: "List lineage snapshots",
		Long:  `List all lineage snapshots for a graph.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			snapshots, err := service.ListSnapshots(ctx, graphID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing snapshots: %v\n", err)
				os.Exit(1)
			}

			if len(snapshots) == 0 {
				fmt.Println("No lineage snapshots found.")
				return
			}

			fmt.Println("Lineage Snapshots:")
			fmt.Println("ID\tGraph ID\tComment\tCreated")
			for _, snapshot := range snapshots {
				fmt.Printf("%s\t%s\t%s\t%s\n",
					snapshot.ID,
					snapshot.GraphID,
					snapshot.Comment,
					snapshot.CreatedAt.Format(time.RFC3339))
			}
		},
	}

	// Create snapshot create command
	snapshotCreateCmd := &cobra.Command{
		Use:   "create [graph-id]",
		Short: "Create a lineage snapshot",
		Long:  `Create a new lineage snapshot for a graph.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			comment, _ := cmd.Flags().GetString("comment")

			snapshot, err := service.CreateSnapshot(ctx, graphID, comment)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating snapshot: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Created lineage snapshot: %s\n", snapshot.ID)
		},
	}
	snapshotCreateCmd.Flags().String("comment", "", "Comment for the snapshot")

	// Create snapshot get command
	snapshotGetCmd := &cobra.Command{
		Use:   "get [id]",
		Short: "Get a lineage snapshot",
		Long:  `Get a lineage snapshot by ID.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			id := args[0]
			snapshot, err := service.GetSnapshot(ctx, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting snapshot: %v\n", err)
				os.Exit(1)
			}

			if snapshot == nil {
				fmt.Fprintf(os.Stderr, "Snapshot not found: %s\n", id)
				os.Exit(1)
			}

			fmt.Printf("ID: %s\n", snapshot.ID)
			fmt.Printf("Graph ID: %s\n", snapshot.GraphID)
			fmt.Printf("Comment: %s\n", snapshot.Comment)
			fmt.Printf("Created: %s\n", snapshot.CreatedAt.Format(time.RFC3339))
			fmt.Printf("Nodes: %d\n", len(snapshot.Graph.Nodes))
			fmt.Printf("Edges: %d\n", len(snapshot.Graph.Edges))

			format, _ := cmd.Flags().GetString("format")
			if format == "json" {
				jsonData, err := json.MarshalIndent(snapshot, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling snapshot to JSON: %v\n", err)
					os.Exit(1)
				}
				fmt.Println(string(jsonData))
			}
		},
	}
	snapshotGetCmd.Flags().String("format", "", "Output format (json)")

	// Create snapshot get-at-time command
	snapshotGetAtTimeCmd := &cobra.Command{
		Use:   "get-at-time [graph-id] [timestamp]",
		Short: "Get a lineage snapshot at a specific time",
		Long:  `Get a lineage snapshot closest to a specific time.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			timestampStr := args[1]

			timestamp, err := time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing timestamp: %v\n", err)
				os.Exit(1)
			}

			snapshot, err := service.GetSnapshotAtTime(ctx, graphID, timestamp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting snapshot: %v\n", err)
				os.Exit(1)
			}

			if snapshot == nil {
				fmt.Fprintf(os.Stderr, "No snapshot found for graph %s at time %s\n", graphID, timestampStr)
				os.Exit(1)
			}

			fmt.Printf("ID: %s\n", snapshot.ID)
			fmt.Printf("Graph ID: %s\n", snapshot.GraphID)
			fmt.Printf("Comment: %s\n", snapshot.Comment)
			fmt.Printf("Created: %s\n", snapshot.CreatedAt.Format(time.RFC3339))
			fmt.Printf("Nodes: %d\n", len(snapshot.Graph.Nodes))
			fmt.Printf("Edges: %d\n", len(snapshot.Graph.Edges))
		},
	}

	// Add subcommands to snapshot command
	snapshotCmd.AddCommand(snapshotListCmd)
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotGetCmd)
	snapshotCmd.AddCommand(snapshotGetAtTimeCmd)
}

// initLineageVisualizeCmd initializes the lineage visualize command
func initLineageVisualizeCmd(parentCmd *cobra.Command) {
	// Create visualize command
	visualizeCmd := &cobra.Command{
		Use:   "visualize [graph-id]",
		Short: "Visualize a lineage graph",
		Long:  `Generate a visualization of a lineage graph.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")
			width, _ := cmd.Flags().GetInt("width")
			height, _ := cmd.Flags().GetInt("height")
			focusNodeID, _ := cmd.Flags().GetString("focus-node")
			maxDepth, _ := cmd.Flags().GetInt("max-depth")
			layout, _ := cmd.Flags().GetString("layout")
			title, _ := cmd.Flags().GetString("title")

			// Create visualization options
			options := visualization.NewDefaultVisualizationOptions()
			options.Format = visualization.GraphFormat(format)
			options.Width = width
			options.Height = height
			options.FocusNodeID = focusNodeID
			options.MaxDepth = maxDepth
			options.LayoutAlgorithm = layout
			options.Title = title

			// Render graph
			data, err := service.RenderGraph(ctx, graphID, options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error visualizing graph: %v\n", err)
				os.Exit(1)
			}

			// Write to output file or stdout
			if output == "" {
				fmt.Println(string(data))
			} else {
				if err := os.WriteFile(output, data, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Visualization saved to %s\n", output)
			}
		},
	}
	visualizeCmd.Flags().String("format", "html", "Output format (html, json, svg, dot)")
	visualizeCmd.Flags().String("output", "", "Output file path")
	visualizeCmd.Flags().Int("width", 1200, "Width of the visualization")
	visualizeCmd.Flags().Int("height", 800, "Height of the visualization")
	visualizeCmd.Flags().String("focus-node", "", "ID of the node to focus on")
	visualizeCmd.Flags().Int("max-depth", 10, "Maximum depth from focus node")
	visualizeCmd.Flags().String("layout", "force", "Layout algorithm (force, dagre, radial)")
	visualizeCmd.Flags().String("title", "", "Title of the visualization")

	// Add visualize command to parent command
	parentCmd.AddCommand(visualizeCmd)
}

// initLineageCompareCmd initializes the lineage compare command
func initLineageCompareCmd(parentCmd *cobra.Command) {
	// Create compare command
	compareCmd := &cobra.Command{
		Use:   "compare",
		Short: "Compare lineage snapshots",
		Long:  `Compare two lineage snapshots and analyze the impact of changes.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			snapshot1, _ := cmd.Flags().GetString("snapshot1")
			snapshot2, _ := cmd.Flags().GetString("snapshot2")
			output, _ := cmd.Flags().GetString("output")

			diff, err := service.CompareSnapshots(ctx, snapshot1, snapshot2)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error comparing snapshots: %v\n", err)
				os.Exit(1)
			}

			// Print summary
			fmt.Printf("Comparison of snapshots %s and %s:\n", snapshot1, snapshot2)
			fmt.Printf("Added nodes: %d\n", len(diff.AddedNodes))
			fmt.Printf("Removed nodes: %d\n", len(diff.RemovedNodes))
			fmt.Printf("Modified nodes: %d\n", len(diff.ModifiedNodes))
			fmt.Printf("Added edges: %d\n", len(diff.AddedEdges))
			fmt.Printf("Removed edges: %d\n", len(diff.RemovedEdges))

			// Write to output file if specified
			if output != "" {
				jsonData, err := json.MarshalIndent(diff, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error marshaling diff to JSON: %v\n", err)
					os.Exit(1)
				}

				if err := os.WriteFile(output, jsonData, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Comparison saved to %s\n", output)
			}
		},
	}
	compareCmd.Flags().String("snapshot1", "", "ID of the first snapshot")
	compareCmd.Flags().String("snapshot2", "", "ID of the second snapshot")
	compareCmd.Flags().String("output", "", "Output file path for detailed comparison")
	compareCmd.MarkFlagRequired("snapshot1")
	compareCmd.MarkFlagRequired("snapshot2")

	// Add compare command to parent command
	parentCmd.AddCommand(compareCmd)
}

// initLineageImportCmd initializes the lineage import command
func initLineageImportCmd(parentCmd *cobra.Command) {
	// Create import command
	importCmd := &cobra.Command{
		Use:   "import",
		Short: "Import lineage information",
		Long:  `Import lineage information from various sources.`,
	}

	// Add import command to parent command
	parentCmd.AddCommand(importCmd)

	// Create import catalog command
	importCatalogCmd := &cobra.Command{
		Use:   "catalog",
		Short: "Import lineage from a data catalog",
		Long:  `Import lineage information from a data catalog.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			catalogName, _ := cmd.Flags().GetString("catalog")
			database, _ := cmd.Flags().GetString("database")
			table, _ := cmd.Flags().GetString("table")

			graph, err := service.ImportFromCatalog(ctx, catalogName, database, table)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error importing from catalog: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Imported lineage graph: %s\n", graph.ID)
			fmt.Printf("Name: %s\n", graph.Name)
			fmt.Printf("Nodes: %d\n", len(graph.Nodes))
			fmt.Printf("Edges: %d\n", len(graph.Edges))
		},
	}
	importCatalogCmd.Flags().String("catalog", "", "Name of the catalog")
	importCatalogCmd.Flags().String("database", "", "Name of the database")
	importCatalogCmd.Flags().String("table", "", "Name of the table")
	importCatalogCmd.MarkFlagRequired("catalog")
	importCatalogCmd.MarkFlagRequired("database")
	importCatalogCmd.MarkFlagRequired("table")

	// Add subcommands to import command
	importCmd.AddCommand(importCatalogCmd)
}

// initLineageExportCmd initializes the lineage export command
func initLineageExportCmd(parentCmd *cobra.Command) {
	// Create export command
	exportCmd := &cobra.Command{
		Use:   "export [graph-id]",
		Short: "Export lineage information",
		Long:  `Export lineage information to various formats.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			service := getLineageService()

			graphID := args[0]
			format, _ := cmd.Flags().GetString("format")
			output, _ := cmd.Flags().GetString("output")

			data, err := service.ExportLineage(ctx, graphID, format)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error exporting lineage: %v\n", err)
				os.Exit(1)
			}

			// Write to output file or stdout
			if output == "" {
				fmt.Println(string(data))
			} else {
				if err := os.WriteFile(output, data, 0644); err != nil {
					fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Lineage exported to %s\n", output)
			}
		},
	}
	exportCmd.Flags().String("format", "json", "Output format (json, dot, svg, html)")
	exportCmd.Flags().String("output", "", "Output file path")

	// Add export command to parent command
	parentCmd.AddCommand(exportCmd)
}

// getLineageService returns a lineage service instance
func getLineageService() *lineage.LineageService {
	// Get config directory
	configDir := getConfigDir()
	
	// Create storage directory
	storageDir := filepath.Join(configDir, "lineage")
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating storage directory: %v\n", err)
		os.Exit(1)
	}
	
	// Create template directory
	templateDir := filepath.Join(configDir, "templates", "lineage")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating template directory: %v\n", err)
		os.Exit(1)
	}
	
	// Create storage
	storage, err := model.NewFileLineageStorage(storageDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating lineage storage: %v\n", err)
		os.Exit(1)
	}
	
	// Create renderer
	renderer, err := visualization.NewGraphRenderer(templateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating graph renderer: %v\n", err)
		os.Exit(1)
	}
	
	// Create catalog manager
	catalogManager := catalog.NewCatalogManager()
	
	// Create lineage service
	return lineage.NewLineageService(storage, renderer, catalogManager)
}

// getConfigDir returns the configuration directory
func getConfigDir() string {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	
	// Create config directory
	configDir := filepath.Join(homeDir, ".nessi")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}
	
	return configDir
}
