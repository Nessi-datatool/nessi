package visualization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nessi-dev/nessi/pkg/lineage/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphRenderer(t *testing.T) {
	// Create a temporary directory for templates
	tempDir, err := os.MkdirTemp("", "lineage-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a new renderer
	renderer, err := NewGraphRenderer(tempDir)
	require.NoError(t, err)
	require.NotNil(t, renderer)

	// Create a test graph
	graph := createTestGraph()

	// Test HTML rendering
	t.Run("RenderHTML", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = FormatHTML
		options.Title = "Test Graph"

		data, err := renderer.RenderGraph(graph, options)
		require.NoError(t, err)
		require.NotNil(t, data)

		// Check that the output is HTML
		assert.True(t, strings.Contains(string(data), "<!DOCTYPE html>"))
		assert.True(t, strings.Contains(string(data), "Test Graph"))
	})

	// Test JSON rendering
	t.Run("RenderJSON", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = FormatJSON

		data, err := renderer.RenderGraph(graph, options)
		require.NoError(t, err)
		require.NotNil(t, data)

		// Check that the output is valid JSON
		var d3Graph D3Graph
		err = json.Unmarshal(data, &d3Graph)
		require.NoError(t, err)

		// Check that the graph has the expected number of nodes and links
		assert.Equal(t, 3, len(d3Graph.Nodes))
		assert.Equal(t, 2, len(d3Graph.Links))
	})

	// Test DOT rendering
	t.Run("RenderDOT", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = FormatDOT

		data, err := renderer.RenderGraph(graph, options)
		require.NoError(t, err)
		require.NotNil(t, data)

		// Check that the output is DOT format
		assert.True(t, strings.Contains(string(data), "digraph G {"))
		assert.True(t, strings.Contains(string(data), "node [shape=box, style=filled];"))
	})

	// Test filtering by node type
	t.Run("FilterByNodeType", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = FormatJSON
		options.IncludeNodeTypes = map[model.NodeType]bool{
			model.NodeTypeTable: true,
		}

		data, err := renderer.RenderGraph(graph, options)
		require.NoError(t, err)
		require.NotNil(t, data)

		// Check that only the table nodes are included
		var d3Graph D3Graph
		err = json.Unmarshal(data, &d3Graph)
		require.NoError(t, err)

		assert.Equal(t, 2, len(d3Graph.Nodes))
		for _, node := range d3Graph.Nodes {
			assert.Equal(t, "table", node.Type)
		}
	})

	// Test filtering by edge type
	t.Run("FilterByEdgeType", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = FormatJSON
		options.IncludeEdgeTypes = map[model.EdgeType]bool{
			model.EdgeTypeRead: true,
		}

		data, err := renderer.RenderGraph(graph, options)
		require.NoError(t, err)
		require.NotNil(t, data)

		// Check that only the read edges are included
		var d3Graph D3Graph
		err = json.Unmarshal(data, &d3Graph)
		require.NoError(t, err)

		assert.Equal(t, 1, len(d3Graph.Links))
		for _, link := range d3Graph.Links {
			assert.Equal(t, "read", link.Type)
		}
	})

	// Test focus node and max depth
	t.Run("FocusNodeAndMaxDepth", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = FormatJSON
		options.FocusNodeID = graph.Nodes[1].ID // The process node
		options.MaxDepth = 1

		data, err := renderer.RenderGraph(graph, options)
		require.NoError(t, err)
		require.NotNil(t, data)

		// Check that only the nodes within the max depth are included
		var d3Graph D3Graph
		err = json.Unmarshal(data, &d3Graph)
		require.NoError(t, err)

		// Should include the process node and its immediate neighbors
		assert.Equal(t, 3, len(d3Graph.Nodes))
	})

	// Test unsupported format
	t.Run("UnsupportedFormat", func(t *testing.T) {
		options := NewDefaultVisualizationOptions()
		options.Format = "unsupported"

		_, err := renderer.RenderGraph(graph, options)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported format")
	})
}

// createTestGraph creates a test lineage graph
func createTestGraph() *model.LineageGraph {
	graph := model.NewLineageGraph("Test Graph")
	graph.Description = "A test graph for unit testing"

	// Create nodes
	sourceTable := model.NewNode("source_table", model.NodeTypeTable)
	sourceTable.Description = "Source data table"
	sourceTable.Properties = map[string]string{
		"schema": "public",
		"rows":   "1000",
	}

	process := model.NewNode("etl_process", model.NodeTypeProcess)
	process.Description = "ETL process"

	targetTable := model.NewNode("target_table", model.NodeTypeTable)
	targetTable.Description = "Target data table"
	targetTable.Properties = map[string]string{
		"schema": "public",
		"rows":   "950",
	}

	// Add nodes to graph
	graph.AddNode(sourceTable)
	graph.AddNode(process)
	graph.AddNode(targetTable)

	// Create edges
	readEdge := model.NewEdge(sourceTable.ID, process.ID, model.EdgeTypeRead)
	writeEdge := model.NewEdge(process.ID, targetTable.ID, model.EdgeTypeWrite)

	// Add edges to graph
	graph.AddEdge(readEdge)
	graph.AddEdge(writeEdge)

	return graph
}

func TestConvertToD3Format(t *testing.T) {
	// Create a test graph
	graph := createTestGraph()

	// Convert to D3 format
	options := NewDefaultVisualizationOptions()
	d3Graph, err := convertToD3Format(graph, options)
	require.NoError(t, err)
	require.NotNil(t, d3Graph)

	// Check nodes
	assert.Equal(t, 3, len(d3Graph.Nodes))
	for _, node := range d3Graph.Nodes {
		assert.NotEmpty(t, node.ID)
		assert.NotEmpty(t, node.Name)
		assert.NotEmpty(t, node.Type)
	}

	// Check links
	assert.Equal(t, 2, len(d3Graph.Links))
	for _, link := range d3Graph.Links {
		assert.NotEmpty(t, link.Source)
		assert.NotEmpty(t, link.Target)
		assert.NotEmpty(t, link.Type)
	}
}

func TestCalculateNodeDepths(t *testing.T) {
	// Create a test graph
	graph := createTestGraph()

	// Calculate node depths
	depths := make(map[string]int)
	calculateNodeDepths(graph, graph.Nodes[1].ID, depths, 0, 10) // Start from the process node

	// Check depths
	assert.Equal(t, 0, depths[graph.Nodes[1].ID]) // Process node at depth 0
	assert.Equal(t, 1, depths[graph.Nodes[0].ID]) // Source table at depth 1 (upstream)
	assert.Equal(t, 1, depths[graph.Nodes[2].ID]) // Target table at depth 1 (downstream)
}

func TestGetNodeColor(t *testing.T) {
	// Test known node types
	assert.Equal(t, "#4285F4", getNodeColor(model.NodeTypeTable))
	assert.Equal(t, "#34A853", getNodeColor(model.NodeTypeView))
	assert.Equal(t, "#FBBC05", getNodeColor(model.NodeTypeQuery))
	assert.Equal(t, "#EA4335", getNodeColor(model.NodeTypeProcess))

	// Test unknown node type
	assert.Equal(t, "#999999", getNodeColor("unknown"))
}

func TestCreateDefaultTemplates(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "templates")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Print temp directory for debugging
	t.Logf("Using temp directory: %s", tempDir)

	// Create default templates
	err = createDefaultTemplates(tempDir)
	if err != nil {
		t.Fatalf("Failed to create default templates: %v", err)
	}

	// Check that the HTML template was created
	htmlTemplatePath := filepath.Join(tempDir, "html_template.html")
	_, err = os.Stat(htmlTemplatePath)
	if err != nil {
		t.Fatalf("Failed to stat HTML template file: %v", err)
	}

	// Check the content of the HTML template
	content, err := os.ReadFile(htmlTemplatePath)
	if err != nil {
		t.Fatalf("Failed to read HTML template file: %v", err)
	}

	// Check template content
	if !strings.Contains(string(content), "<!DOCTYPE html>") {
		t.Errorf("HTML template does not contain DOCTYPE")
	}
	if !strings.Contains(string(content), "d3js.org/d3.v7.min.js") {
		t.Errorf("HTML template does not contain d3.js reference")
	}
}
