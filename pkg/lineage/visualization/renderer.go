package visualization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/nessi-dev/nessi-dev/pkg/lineage/model"
)

// GraphFormat represents the output format for lineage graph visualization
type GraphFormat string

const (
	// FormatHTML represents HTML output format
	FormatHTML GraphFormat = "html"
	
	// FormatJSON represents JSON output format
	FormatJSON GraphFormat = "json"
	
	// FormatSVG represents SVG output format
	FormatSVG GraphFormat = "svg"
	
	// FormatPNG represents PNG output format
	FormatPNG GraphFormat = "png"
	
	// FormatDOT represents DOT (Graphviz) output format
	FormatDOT GraphFormat = "dot"
)

// VisualizationOptions represents options for lineage graph visualization
type VisualizationOptions struct {
	Format           GraphFormat         `json:"format"`
	IncludeNodeTypes map[model.NodeType]bool `json:"include_node_types,omitempty"`
	IncludeEdgeTypes map[model.EdgeType]bool `json:"include_edge_types,omitempty"`
	MaxDepth         int                `json:"max_depth"`
	FocusNodeID      string             `json:"focus_node_id,omitempty"`
	ShowProperties   bool               `json:"show_properties"`
	ColorScheme      string             `json:"color_scheme,omitempty"`
	LayoutAlgorithm  string             `json:"layout_algorithm,omitempty"`
	Width            int                `json:"width"`
	Height           int                `json:"height"`
	Title            string             `json:"title,omitempty"`
}

// NewDefaultVisualizationOptions creates default visualization options
func NewDefaultVisualizationOptions() *VisualizationOptions {
	return &VisualizationOptions{
		Format:           FormatHTML,
		IncludeNodeTypes: make(map[model.NodeType]bool),
		IncludeEdgeTypes: make(map[model.EdgeType]bool),
		MaxDepth:         10,
		ShowProperties:   true,
		ColorScheme:      "default",
		LayoutAlgorithm:  "force",
		Width:            1200,
		Height:           800,
	}
}

// GraphRenderer renders lineage graphs
type GraphRenderer struct {
	templateDir string
}

// NewGraphRenderer creates a new graph renderer
func NewGraphRenderer(templateDir string) (*GraphRenderer, error) {
	// Create template directory if it doesn't exist
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create template directory: %w", err)
	}
	
	// Create default templates if they don't exist
	if err := createDefaultTemplates(templateDir); err != nil {
		return nil, fmt.Errorf("failed to create default templates: %w", err)
	}
	
	return &GraphRenderer{
		templateDir: templateDir,
	}, nil
}

// createDefaultTemplates creates default visualization templates
func createDefaultTemplates(templateDir string) error {
	// Create HTML template
	htmlTemplate := `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        #graph-container {
            width: {{.Width}}px;
            height: {{.Height}}px;
            margin: 20px auto;
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .node {
            cursor: pointer;
            stroke: #fff;
            stroke-width: 2px;
        }
        .link {
            stroke: #999;
            stroke-opacity: 0.6;
        }
        .node-label {
            font-size: 12px;
            fill: #333;
            text-anchor: middle;
            pointer-events: none;
        }
        .tooltip {
            position: absolute;
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 3px;
            padding: 10px;
            font-size: 12px;
            pointer-events: none;
            opacity: 0;
            transition: opacity 0.3s;
            max-width: 300px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .controls {
            margin: 10px;
            padding: 10px;
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 5px;
        }
    </style>
    <script src="https://d3js.org/d3.v7.min.js"></script>
</head>
<body>
    <div class="controls">
        <h2>{{.Title}}</h2>
        <div>
            <label>Layout: </label>
            <select id="layout-select">
                <option value="force" {{if eq .LayoutAlgorithm "force"}}selected{{end}}>Force-Directed</option>
                <option value="dagre" {{if eq .LayoutAlgorithm "dagre"}}selected{{end}}>Hierarchical</option>
                <option value="radial" {{if eq .LayoutAlgorithm "radial"}}selected{{end}}>Radial</option>
            </select>
            <button id="reset-zoom">Reset Zoom</button>
            <button id="export-svg">Export SVG</button>
            <button id="export-png">Export PNG</button>
        </div>
    </div>
    <div id="graph-container"></div>
    <div class="tooltip" id="tooltip"></div>

    <script>
        // Graph data
        const graphData = {{.GraphData}};
        
        // Create the graph visualization
        createGraph(graphData);
        
        function createGraph(data) {
            const container = document.getElementById('graph-container');
            const width = {{.Width}};
            const height = {{.Height}};
            
            // Clear previous graph
            container.innerHTML = '';
            
            // Create SVG
            const svg = d3.select(container)
                .append('svg')
                .attr('width', width)
                .attr('height', height)
                .attr('viewBox', [0, 0, width, height])
                .call(d3.zoom().on('zoom', (event) => {
                    g.attr('transform', event.transform);
                }));
            
            const g = svg.append('g');
            
            // Create links
            const links = g.selectAll('.link')
                .data(data.links)
                .enter()
                .append('line')
                .attr('class', 'link')
                .attr('stroke-width', d => d.value || 1)
                .attr('marker-end', 'url(#arrow)');
            
            // Create arrow marker
            svg.append('defs').append('marker')
                .attr('id', 'arrow')
                .attr('viewBox', '0 -5 10 10')
                .attr('refX', 15)
                .attr('refY', 0)
                .attr('markerWidth', 6)
                .attr('markerHeight', 6)
                .attr('orient', 'auto')
                .append('path')
                .attr('d', 'M0,-5L10,0L0,5')
                .attr('fill', '#999');
            
            // Create nodes
            const nodes = g.selectAll('.node')
                .data(data.nodes)
                .enter()
                .append('circle')
                .attr('class', 'node')
                .attr('r', d => d.size || 10)
                .attr('fill', d => getNodeColor(d.type))
                .on('mouseover', showTooltip)
                .on('mouseout', hideTooltip)
                .on('click', focusNode);
            
            // Create node labels
            const labels = g.selectAll('.node-label')
                .data(data.nodes)
                .enter()
                .append('text')
                .attr('class', 'node-label')
                .attr('dy', 20)
                .text(d => d.name);
            
            // Create force simulation
            const layout = d3.select('#layout-select').property('value');
            
            if (layout === 'force') {
                const simulation = d3.forceSimulation(data.nodes)
                    .force('link', d3.forceLink(data.links).id(d => d.id).distance(100))
                    .force('charge', d3.forceManyBody().strength(-300))
                    .force('center', d3.forceCenter(width / 2, height / 2))
                    .on('tick', () => {
                        links
                            .attr('x1', d => d.source.x)
                            .attr('y1', d => d.source.y)
                            .attr('x2', d => d.target.x)
                            .attr('y2', d => d.target.y);
                        
                        nodes
                            .attr('cx', d => d.x)
                            .attr('cy', d => d.y);
                        
                        labels
                            .attr('x', d => d.x)
                            .attr('y', d => d.y);
                    });
                
                // Drag behavior
                nodes.call(d3.drag()
                    .on('start', (event, d) => {
                        if (!event.active) simulation.alphaTarget(0.3).restart();
                        d.fx = d.x;
                        d.fy = d.y;
                    })
                    .on('drag', (event, d) => {
                        d.fx = event.x;
                        d.fy = event.y;
                    })
                    .on('end', (event, d) => {
                        if (!event.active) simulation.alphaTarget(0);
                        d.fx = null;
                        d.fy = null;
                    }));
            } else if (layout === 'dagre') {
                // Hierarchical layout (simplified)
                const levels = {};
                
                // Assign levels based on depth
                data.nodes.forEach(node => {
                    const depth = node.depth || 0;
                    if (!levels[depth]) {
                        levels[depth] = [];
                    }
                    levels[depth].push(node);
                });
                
                // Position nodes based on levels
                const levelCount = Object.keys(levels).length;
                const levelHeight = height / (levelCount + 1);
                
                Object.keys(levels).forEach((level, i) => {
                    const nodesInLevel = levels[level];
                    const levelWidth = width / (nodesInLevel.length + 1);
                    
                    nodesInLevel.forEach((node, j) => {
                        node.x = levelWidth * (j + 1);
                        node.y = levelHeight * (parseInt(level) + 1);
                    });
                });
                
                // Update positions
                links
                    .attr('x1', d => d.source.x)
                    .attr('y1', d => d.source.y)
                    .attr('x2', d => d.target.x)
                    .attr('y2', d => d.target.y);
                
                nodes
                    .attr('cx', d => d.x)
                    .attr('cy', d => d.y);
                
                labels
                    .attr('x', d => d.x)
                    .attr('y', d => d.y);
            } else if (layout === 'radial') {
                // Radial layout (simplified)
                const centerX = width / 2;
                const centerY = height / 2;
                const radius = Math.min(width, height) / 3;
                
                data.nodes.forEach((node, i) => {
                    const angle = (i / data.nodes.length) * 2 * Math.PI;
                    node.x = centerX + radius * Math.cos(angle);
                    node.y = centerY + radius * Math.sin(angle);
                });
                
                // Update positions
                links
                    .attr('x1', d => d.source.x)
                    .attr('y1', d => d.source.y)
                    .attr('x2', d => d.target.x)
                    .attr('y2', d => d.target.y);
                
                nodes
                    .attr('cx', d => d.x)
                    .attr('cy', d => d.y);
                
                labels
                    .attr('x', d => d.x)
                    .attr('y', d => d.y);
            }
            
            // Layout change handler
            d3.select('#layout-select').on('change', function() {
                createGraph(data);
            });
            
            // Reset zoom handler
            d3.select('#reset-zoom').on('click', function() {
                svg.transition().duration(750).call(
                    d3.zoom().transform,
                    d3.zoomIdentity
                );
            });
            
            // Export handlers
            d3.select('#export-svg').on('click', exportSVG);
            d3.select('#export-png').on('click', exportPNG);
            
            // Tooltip functions
            function showTooltip(event, d) {
                const tooltip = d3.select('#tooltip');
                
                let content = '<div><strong>' + d.name + '</strong></div>';
                content += '<div>Type: ' + d.type + '</div>';
                
                if (d.description) {
                    content += '<div>Description: ' + d.description + '</div>';
                }
                
                if (d.properties && {{.ShowProperties}}) {
                    content += '<div><strong>Properties:</strong></div>';
                    for (const [key, value] of Object.entries(d.properties)) {
                        content += '<div>' + key + ': ' + value + '</div>';
                    }
                }
                
                tooltip.html(content)
                    .style('left', (event.pageX + 10) + 'px')
                    .style('top', (event.pageY - 10) + 'px')
                    .style('opacity', 1);
            }
            
            function hideTooltip() {
                d3.select('#tooltip').style('opacity', 0);
            }
            
            function focusNode(event, d) {
                // Highlight the selected node and its connections
                nodes.attr('stroke-width', node => node === d ? 4 : 2);
                
                links.attr('stroke', link => {
                    if (link.source === d || link.target === d) {
                        return '#ff7700';
                    }
                    return '#999';
                }).attr('stroke-width', link => {
                    if (link.source === d || link.target === d) {
                        return 3;
                    }
                    return 1;
                });
            }
            
            // Export functions
            function exportSVG() {
                const svgData = new XMLSerializer().serializeToString(svg.node());
                const blob = new Blob([svgData], { type: 'image/svg+xml' });
                const url = URL.createObjectURL(blob);
                
                const link = document.createElement('a');
                link.href = url;
                link.download = 'lineage_graph.svg';
                link.click();
                
                URL.revokeObjectURL(url);
            }
            
            function exportPNG() {
                const svgData = new XMLSerializer().serializeToString(svg.node());
                const canvas = document.createElement('canvas');
                canvas.width = width;
                canvas.height = height;
                
                const ctx = canvas.getContext('2d');
                const img = new Image();
                
                img.onload = function() {
                    ctx.drawImage(img, 0, 0);
                    const pngUrl = canvas.toDataURL('image/png');
                    
                    const link = document.createElement('a');
                    link.href = pngUrl;
                    link.download = 'lineage_graph.png';
                    link.click();
                };
                
                img.src = 'data:image/svg+xml;base64,' + btoa(unescape(encodeURIComponent(svgData)));
            }
        }
        
        // Node color function
        function getNodeColor(type) {
            const colorMap = {
                'table': '#4285F4',     // Google Blue
                'view': '#34A853',      // Google Green
                'query': '#FBBC05',     // Google Yellow
                'process': '#EA4335',   // Google Red
                'dataset': '#673AB7',   // Purple
                'api': '#FF6D00',       // Orange
                'application': '#2196F3' // Blue
            };
            
            return colorMap[type] || '#999999';
        }
    </script>
</body>
</html>`
	
	// Ensure the template directory exists
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}
	
	htmlTemplatePath := filepath.Join(templateDir, "html_template.html")
	if err := os.WriteFile(htmlTemplatePath, []byte(htmlTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write HTML template: %w", err)
	}
	
	return nil
}

// RenderGraph renders a lineage graph
func (r *GraphRenderer) RenderGraph(graph *model.LineageGraph, options *VisualizationOptions) ([]byte, error) {
	switch options.Format {
	case FormatHTML:
		return r.renderHTML(graph, options)
	case FormatJSON:
		return r.renderJSON(graph, options)
	case FormatSVG:
		return r.renderSVG(graph, options)
	case FormatDOT:
		return r.renderDOT(graph, options)
	default:
		return nil, fmt.Errorf("unsupported format: %s", options.Format)
	}
}

// renderHTML renders a lineage graph as HTML
func (r *GraphRenderer) renderHTML(graph *model.LineageGraph, options *VisualizationOptions) ([]byte, error) {
	// Read HTML template
	templatePath := filepath.Join(r.templateDir, "html_template.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML template: %w", err)
	}
	
	// Convert graph to D3 format
	d3Graph, err := convertToD3Format(graph, options)
	if err != nil {
		return nil, fmt.Errorf("failed to convert graph to D3 format: %w", err)
	}
	
	d3GraphJSON, err := json.Marshal(d3Graph)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal D3 graph: %w", err)
	}
	
	// Set title if not provided
	if options.Title == "" {
		options.Title = graph.Name
	}
	
	// Create template data
	data := struct {
		Title           string
		Width           int
		Height          int
		ShowProperties  bool
		LayoutAlgorithm string
		GraphData       template.JS
	}{
		Title:           options.Title,
		Width:           options.Width,
		Height:          options.Height,
		ShowProperties:  options.ShowProperties,
		LayoutAlgorithm: options.LayoutAlgorithm,
		GraphData:       template.JS(d3GraphJSON),
	}
	
	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute HTML template: %w", err)
	}
	
	return buf.Bytes(), nil
}

// renderJSON renders a lineage graph as JSON
func (r *GraphRenderer) renderJSON(graph *model.LineageGraph, options *VisualizationOptions) ([]byte, error) {
	// Convert graph to D3 format
	d3Graph, err := convertToD3Format(graph, options)
	if err != nil {
		return nil, fmt.Errorf("failed to convert graph to D3 format: %w", err)
	}
	
	// Marshal to JSON
	return json.MarshalIndent(d3Graph, "", "  ")
}

// renderSVG renders a lineage graph as SVG
func (r *GraphRenderer) renderSVG(graph *model.LineageGraph, options *VisualizationOptions) ([]byte, error) {
	// This is a placeholder that would need to be implemented
	// with a proper SVG rendering library
	return nil, fmt.Errorf("SVG rendering not implemented")
}

// renderDOT renders a lineage graph as DOT (Graphviz)
func (r *GraphRenderer) renderDOT(graph *model.LineageGraph, options *VisualizationOptions) ([]byte, error) {
	var buf bytes.Buffer
	
	// Write DOT header
	buf.WriteString("digraph G {\n")
	buf.WriteString("  rankdir=LR;\n")
	buf.WriteString("  node [shape=box, style=filled];\n")
	
	// Write nodes
	for _, node := range graph.Nodes {
		// Skip node if its type is not included
		if len(options.IncludeNodeTypes) > 0 {
			if !options.IncludeNodeTypes[node.Type] {
				continue
			}
		}
		
		// Get node color
		color := getNodeColor(node.Type)
		
		// Escape node name for DOT
		name := strings.ReplaceAll(node.Name, "\"", "\\\"")
		
		// Write node
		buf.WriteString(fmt.Sprintf("  \"%s\" [label=\"%s\\n(%s)\", fillcolor=\"%s\"];\n",
			node.ID, name, node.Type, color))
	}
	
	// Write edges
	for _, edge := range graph.Edges {
		// Skip edge if its type is not included
		if len(options.IncludeEdgeTypes) > 0 {
			if !options.IncludeEdgeTypes[edge.Type] {
				continue
			}
		}
		
		// Write edge
		buf.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\" [label=\"%s\"];\n",
			edge.SourceID, edge.TargetID, edge.Type))
	}
	
	// Write DOT footer
	buf.WriteString("}\n")
	
	return buf.Bytes(), nil
}

// D3Node represents a node in D3 format
type D3Node struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Description string            `json:"description,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	Size        int               `json:"size"`
	Depth       int               `json:"depth,omitempty"`
}

// D3Link represents a link in D3 format
type D3Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Value  int    `json:"value"`
}

// D3Graph represents a graph in D3 format
type D3Graph struct {
	Nodes []D3Node `json:"nodes"`
	Links []D3Link `json:"links"`
}

// convertToD3Format converts a lineage graph to D3 format
func convertToD3Format(graph *model.LineageGraph, options *VisualizationOptions) (*D3Graph, error) {
	d3Graph := &D3Graph{
		Nodes: make([]D3Node, 0, len(graph.Nodes)),
		Links: make([]D3Link, 0, len(graph.Edges)),
	}
	
	// Process nodes
	nodeDepths := make(map[string]int)
	if options.FocusNodeID != "" {
		// Calculate node depths from focus node
		calculateNodeDepths(graph, options.FocusNodeID, nodeDepths, 0, options.MaxDepth)
	}
	
	for _, node := range graph.Nodes {
		// Skip node if its type is not included
		if len(options.IncludeNodeTypes) > 0 {
			if !options.IncludeNodeTypes[node.Type] {
				continue
			}
		}
		
		// Skip node if it's beyond the max depth
		if options.FocusNodeID != "" {
			depth, exists := nodeDepths[node.ID]
			if !exists || depth > options.MaxDepth {
				continue
			}
		}
		
		// Create D3 node
		d3Node := D3Node{
			ID:          node.ID,
			Name:        node.Name,
			Type:        string(node.Type),
			Description: node.Description,
			Size:        10, // Default size
		}
		
		// Include properties if requested
		if options.ShowProperties {
			d3Node.Properties = node.Properties
		}
		
		// Set node depth if available
		if depth, exists := nodeDepths[node.ID]; exists {
			d3Node.Depth = depth
		}
		
		// Add node to D3 graph
		d3Graph.Nodes = append(d3Graph.Nodes, d3Node)
	}
	
	// Process edges
	for _, edge := range graph.Edges {
		// Skip edge if its type is not included
		if len(options.IncludeEdgeTypes) > 0 {
			if !options.IncludeEdgeTypes[edge.Type] {
				continue
			}
		}
		
		// Skip edge if source or target node is not included
		sourceIncluded := false
		targetIncluded := false
		
		for _, node := range d3Graph.Nodes {
			if node.ID == edge.SourceID {
				sourceIncluded = true
			}
			if node.ID == edge.TargetID {
				targetIncluded = true
			}
		}
		
		if !sourceIncluded || !targetIncluded {
			continue
		}
		
		// Create D3 link
		d3Link := D3Link{
			Source: edge.SourceID,
			Target: edge.TargetID,
			Type:   string(edge.Type),
			Value:  1, // Default value
		}
		
		// Add link to D3 graph
		d3Graph.Links = append(d3Graph.Links, d3Link)
	}
	
	return d3Graph, nil
}

// calculateNodeDepths calculates the depth of each node from a focus node
func calculateNodeDepths(graph *model.LineageGraph, focusNodeID string, depths map[string]int, currentDepth, maxDepth int) {
	// Stop if we've reached the maximum depth
	if currentDepth > maxDepth {
		return
	}
	
	// Set depth for the current node
	depths[focusNodeID] = currentDepth
	
	// Process upstream nodes
	upstreamNodes := graph.GetUpstreamNodes(focusNodeID)
	for _, node := range upstreamNodes {
		// Skip if we've already processed this node with a lower depth
		if existingDepth, exists := depths[node.ID]; exists && existingDepth <= currentDepth+1 {
			continue
		}
		
		// Calculate depth for upstream node
		calculateNodeDepths(graph, node.ID, depths, currentDepth+1, maxDepth)
	}
	
	// Process downstream nodes
	downstreamNodes := graph.GetDownstreamNodes(focusNodeID)
	for _, node := range downstreamNodes {
		// Skip if we've already processed this node with a lower depth
		if existingDepth, exists := depths[node.ID]; exists && existingDepth <= currentDepth+1 {
			continue
		}
		
		// Calculate depth for downstream node
		calculateNodeDepths(graph, node.ID, depths, currentDepth+1, maxDepth)
	}
}

// getNodeColor returns a color for a node type
func getNodeColor(nodeType model.NodeType) string {
	colorMap := map[model.NodeType]string{
		model.NodeTypeTable:       "#4285F4", // Google Blue
		model.NodeTypeView:        "#34A853", // Google Green
		model.NodeTypeQuery:       "#FBBC05", // Google Yellow
		model.NodeTypeProcess:     "#EA4335", // Google Red
		model.NodeTypeDataset:     "#673AB7", // Purple
		model.NodeTypeAPI:         "#FF6D00", // Orange
		model.NodeTypeApplication: "#2196F3", // Blue
	}
	
	if color, exists := colorMap[nodeType]; exists {
		return color
	}
	
	return "#999999" // Default gray
}
