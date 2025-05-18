package comparison

import (
	"github.com/nessi-dev/nessi/pkg/lineage/model"
)

// DiffGenerator generates diffs between lineage graphs
type DiffGenerator struct{}

// NewDiffGenerator creates a new diff generator
func NewDiffGenerator() *DiffGenerator {
	return &DiffGenerator{}
}

// GenerateDiff generates a diff between two lineage graphs
func (g *DiffGenerator) GenerateDiff(oldGraph, newGraph *model.LineageGraph) *model.LineageDiff {
	diff := model.NewLineageDiff(oldGraph, newGraph)
	
	// Create maps for faster lookup
	oldNodes := make(map[string]*model.Node)
	newNodes := make(map[string]*model.Node)
	oldEdges := make(map[string]*model.Edge)
	newEdges := make(map[string]*model.Edge)
	
	// Populate maps
	for _, node := range oldGraph.Nodes {
		oldNodes[node.ID] = node
	}
	
	for _, node := range newGraph.Nodes {
		newNodes[node.ID] = node
	}
	
	for _, edge := range oldGraph.Edges {
		oldEdges[edge.ID] = edge
	}
	
	for _, edge := range newGraph.Edges {
		newEdges[edge.ID] = edge
	}
	
	// Find added, removed, and modified nodes
	for id, node := range newNodes {
		if oldNode, exists := oldNodes[id]; exists {
			// Node exists in both graphs, check if modified
			if isNodeModified(oldNode, node) {
				diff.ModifiedNodes = append(diff.ModifiedNodes, node)
			}
		} else {
			// Node exists only in new graph
			diff.AddedNodes = append(diff.AddedNodes, node)
		}
	}
	
	for id, node := range oldNodes {
		if _, exists := newNodes[id]; !exists {
			// Node exists only in old graph
			diff.RemovedNodes = append(diff.RemovedNodes, node)
		}
	}
	
	// Find added, removed, and modified edges
	for id, edge := range newEdges {
		if oldEdge, exists := oldEdges[id]; exists {
			// Edge exists in both graphs, check if modified
			if isEdgeModified(oldEdge, edge) {
				diff.ModifiedEdges = append(diff.ModifiedEdges, edge)
			}
		} else {
			// Edge exists only in new graph
			diff.AddedEdges = append(diff.AddedEdges, edge)
		}
	}
	
	for id, edge := range oldEdges {
		if _, exists := newEdges[id]; !exists {
			// Edge exists only in old graph
			diff.RemovedEdges = append(diff.RemovedEdges, edge)
		}
	}
	
	return diff
}

// isNodeModified checks if a node has been modified
func isNodeModified(oldNode, newNode *model.Node) bool {
	// Check basic properties
	if oldNode.Name != newNode.Name ||
		oldNode.Type != newNode.Type ||
		oldNode.Description != newNode.Description ||
		oldNode.Catalog != newNode.Catalog ||
		oldNode.Database != newNode.Database ||
		oldNode.Schema != newNode.Schema ||
		oldNode.Table != newNode.Table ||
		oldNode.Version != newNode.Version {
		return true
	}
	
	// Check properties map
	if len(oldNode.Properties) != len(newNode.Properties) {
		return true
	}
	
	for k, v := range oldNode.Properties {
		if newV, exists := newNode.Properties[k]; !exists || newV != v {
			return true
		}
	}
	
	return false
}

// isEdgeModified checks if an edge has been modified
func isEdgeModified(oldEdge, newEdge *model.Edge) bool {
	// Check basic properties
	if oldEdge.SourceID != newEdge.SourceID ||
		oldEdge.TargetID != newEdge.TargetID ||
		oldEdge.Type != newEdge.Type ||
		oldEdge.Description != newEdge.Description ||
		oldEdge.Version != newEdge.Version {
		return true
	}
	
	// Check properties map
	if len(oldEdge.Properties) != len(newEdge.Properties) {
		return true
	}
	
	for k, v := range oldEdge.Properties {
		if newV, exists := newEdge.Properties[k]; !exists || newV != v {
			return true
		}
	}
	
	return false
}

// DiffAnalyzer analyzes lineage diffs to provide insights
type DiffAnalyzer struct{}

// NewDiffAnalyzer creates a new diff analyzer
func NewDiffAnalyzer() *DiffAnalyzer {
	return &DiffAnalyzer{}
}

// ImpactAnalysis performs impact analysis on a lineage diff
type ImpactAnalysis struct {
	AffectedDownstreamNodes []*model.Node `json:"affected_downstream_nodes"`
	AffectedUpstreamNodes   []*model.Node `json:"affected_upstream_nodes"`
	CriticalNodes           []*model.Node `json:"critical_nodes"`
	ImpactSeverity          string        `json:"impact_severity"`
	ImpactDetails           string        `json:"impact_details"`
}

// AnalyzeImpact analyzes the impact of changes in a lineage diff
func (a *DiffAnalyzer) AnalyzeImpact(diff *model.LineageDiff, graph *model.LineageGraph) *ImpactAnalysis {
	analysis := &ImpactAnalysis{
		AffectedDownstreamNodes: make([]*model.Node, 0),
		AffectedUpstreamNodes:   make([]*model.Node, 0),
		CriticalNodes:           make([]*model.Node, 0),
	}
	
	// Analyze impact of modified nodes
	for _, node := range diff.ModifiedNodes {
		// Find downstream nodes
		downstreamNodes := graph.GetDownstreamNodes(node.ID)
		for _, downstream := range downstreamNodes {
			// Check if already in the list
			found := false
			for _, existing := range analysis.AffectedDownstreamNodes {
				if existing.ID == downstream.ID {
					found = true
					break
				}
			}
			
			if !found {
				analysis.AffectedDownstreamNodes = append(analysis.AffectedDownstreamNodes, downstream)
			}
		}
		
		// Find upstream nodes
		upstreamNodes := graph.GetUpstreamNodes(node.ID)
		for _, upstream := range upstreamNodes {
			// Check if already in the list
			found := false
			for _, existing := range analysis.AffectedUpstreamNodes {
				if existing.ID == upstream.ID {
					found = true
					break
				}
			}
			
			if !found {
				analysis.AffectedUpstreamNodes = append(analysis.AffectedUpstreamNodes, upstream)
			}
		}
	}
	
	// Analyze impact of added nodes
	for _, node := range diff.AddedNodes {
		// Find downstream nodes
		downstreamNodes := graph.GetDownstreamNodes(node.ID)
		for _, downstream := range downstreamNodes {
			// Check if already in the list
			found := false
			for _, existing := range analysis.AffectedDownstreamNodes {
				if existing.ID == downstream.ID {
					found = true
					break
				}
			}
			
			if !found {
				analysis.AffectedDownstreamNodes = append(analysis.AffectedDownstreamNodes, downstream)
			}
		}
	}
	
	// Analyze impact of removed nodes
	for _, node := range diff.RemovedNodes {
		// For removed nodes, we need to check the old graph
		oldGraph, err := getOldGraph(diff)
		if err != nil {
			// If we can't get the old graph, we can't analyze the impact
			continue
		}
		
		// Find downstream nodes in the old graph
		downstreamNodes := oldGraph.GetDownstreamNodes(node.ID)
		for _, downstream := range downstreamNodes {
			// Check if the downstream node still exists in the new graph
			if newNode := graph.GetNode(downstream.ID); newNode != nil {
				// Check if already in the list
				found := false
				for _, existing := range analysis.AffectedDownstreamNodes {
					if existing.ID == newNode.ID {
						found = true
						break
					}
				}
				
				if !found {
					analysis.AffectedDownstreamNodes = append(analysis.AffectedDownstreamNodes, newNode)
				}
			}
		}
	}
	
	// Identify critical nodes
	// Critical nodes are those that have many downstream dependencies
	for _, node := range graph.Nodes {
		downstreamNodes := graph.GetDownstreamNodes(node.ID)
		if len(downstreamNodes) > 5 { // Arbitrary threshold
			analysis.CriticalNodes = append(analysis.CriticalNodes, node)
		}
	}
	
	// Determine impact severity
	if len(analysis.AffectedDownstreamNodes) > 10 {
		analysis.ImpactSeverity = "High"
		analysis.ImpactDetails = "Changes affect many downstream nodes"
	} else if len(analysis.AffectedDownstreamNodes) > 5 {
		analysis.ImpactSeverity = "Medium"
		analysis.ImpactDetails = "Changes affect several downstream nodes"
	} else {
		analysis.ImpactSeverity = "Low"
		analysis.ImpactDetails = "Changes have minimal downstream impact"
	}
	
	return analysis
}

// getOldGraph retrieves the old graph from a diff
// This is a placeholder function that would need to be implemented
// based on how graphs are stored and retrieved
func getOldGraph(diff *model.LineageDiff) (*model.LineageGraph, error) {
	// In a real implementation, this would retrieve the old graph
	// from storage using diff.OldGraphID
	return nil, nil
}

// SchemaChangeAnalysis analyzes schema changes in a lineage diff
type SchemaChangeAnalysis struct {
	AddedColumns    []string `json:"added_columns"`
	RemovedColumns  []string `json:"removed_columns"`
	ModifiedColumns []string `json:"modified_columns"`
	ImpactSeverity  string   `json:"impact_severity"`
	ImpactDetails   string   `json:"impact_details"`
}

// AnalyzeSchemaChanges analyzes schema changes in a lineage diff
func (a *DiffAnalyzer) AnalyzeSchemaChanges(diff *model.LineageDiff) *SchemaChangeAnalysis {
	// This is a simplified implementation that would need to be expanded
	// to handle actual schema changes
	analysis := &SchemaChangeAnalysis{
		AddedColumns:    make([]string, 0),
		RemovedColumns:  make([]string, 0),
		ModifiedColumns: make([]string, 0),
	}
	
	// In a real implementation, this would analyze the schema changes
	// by comparing the schemas of modified nodes
	
	return analysis
}
