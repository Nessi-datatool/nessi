package visualization

import (
	"fmt"
	"io"
	"os"
)

// CLIRenderer provides visualization capabilities for CLI
type CLIRenderer struct {
	output io.Writer
}

// NewCLIRenderer creates a new CLI renderer
func NewCLIRenderer(output io.Writer) *CLIRenderer {
	if output == nil {
		output = os.Stdout
	}
	return &CLIRenderer{
		output: output,
	}
}

// RenderTable renders a table to the CLI
func (r *CLIRenderer) RenderTable(headers []string, rows [][]string) error {
	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	
	// Print headers
	for i, h := range headers {
		fmt.Fprintf(r.output, "%-*s", colWidths[i]+2, h)
	}
	fmt.Fprintln(r.output)
	
	// Print separator
	for _, w := range colWidths {
		for i := 0; i < w+2; i++ {
			fmt.Fprint(r.output, "-")
		}
	}
	fmt.Fprintln(r.output)
	
	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) {
				fmt.Fprintf(r.output, "%-*s", colWidths[i]+2, cell)
			}
		}
		fmt.Fprintln(r.output)
	}
	
	return nil
}

// RenderGraph renders a graph representation to the CLI
func (r *CLIRenderer) RenderGraph(nodes []string, edges [][2]string) error {
	fmt.Fprintln(r.output, "Graph Nodes:")
	for _, node := range nodes {
		fmt.Fprintf(r.output, "- %s\n", node)
	}
	
	fmt.Fprintln(r.output, "\nGraph Edges:")
	for _, edge := range edges {
		fmt.Fprintf(r.output, "- %s -> %s\n", edge[0], edge[1])
	}
	
	return nil
}

// RenderJSON renders JSON data to the CLI
func (r *CLIRenderer) RenderJSON(data []byte) error {
	_, err := r.output.Write(data)
	return err
}
