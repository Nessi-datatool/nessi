package dbt

import (
	"fmt"
)

// ModelNode represents a node in the lineage graph
type ModelNode struct {
	Name      string
	NodeID    string
	Parents   []*ModelNode
	Children  []*ModelNode
	ModelType string
	Tags      []string
}

// LineageGraph represents a lineage graph of dbt models
type LineageGraph struct {
	Nodes map[string]*ModelNode
}

// NewLineageGraph creates a new lineage graph from a dbt manifest
func NewLineageGraph(manifest *DBTManifest) (*LineageGraph, error) {
	graph := &LineageGraph{
		Nodes: make(map[string]*ModelNode),
	}

	// Create nodes for all models
	for id, node := range manifest.Nodes {
		if node.ResourceType == "model" {
			modelNode := &ModelNode{
				Name:      node.Name,
				NodeID:    id,
				Parents:   make([]*ModelNode, 0),
				Children:  make([]*ModelNode, 0),
				ModelType: node.ResourceType,
				Tags:      node.Tags,
			}
			graph.Nodes[id] = modelNode
		}
	}

	// Build parent-child relationships
	// In a real implementation, you would parse the manifest.Parent/ChildMap
	// For now, we'll just use a simplified approach

	// Add some mock relationships for demonstration
	for id, node := range graph.Nodes {
		// For each node, find potential parents based on naming conventions
		// This is a very simplified approach - in reality, you'd use the manifest
		for otherId, otherNode := range graph.Nodes {
			if id != otherId {
				// Simple heuristic: if this node's name starts with the other node's name,
				// it might be a child of that node
				if len(node.Name) > len(otherNode.Name) && 
				   node.Name[:len(otherNode.Name)] == otherNode.Name {
					node.Parents = append(node.Parents, otherNode)
					otherNode.Children = append(otherNode.Children, node)
				}
			}
		}
	}

	return graph, nil
}

// GetUpstreamModels gets all upstream models for a given model
func (g *LineageGraph) GetUpstreamModels(modelName string) ([]*ModelNode, error) {
	// Find the model node
	var startNode *ModelNode
	for _, node := range g.Nodes {
		if node.Name == modelName {
			startNode = node
			break
		}
	}

	if startNode == nil {
		return nil, fmt.Errorf("model %s not found in lineage graph", modelName)
	}

	// Traverse the graph to find all upstream models
	visited := make(map[string]bool)
	upstream := make([]*ModelNode, 0)

	var traverse func(*ModelNode)
	traverse = func(node *ModelNode) {
		if visited[node.NodeID] {
			return
		}
		visited[node.NodeID] = true
		upstream = append(upstream, node)

		for _, parent := range node.Parents {
			traverse(parent)
		}
	}

	// Start traversal from the parents of the start node
	for _, parent := range startNode.Parents {
		traverse(parent)
	}

	return upstream, nil
}

// GetDownstreamModels gets all downstream models for a given model
func (g *LineageGraph) GetDownstreamModels(modelName string) ([]*ModelNode, error) {
	// Find the model node
	var startNode *ModelNode
	for _, node := range g.Nodes {
		if node.Name == modelName {
			startNode = node
			break
		}
	}

	if startNode == nil {
		return nil, fmt.Errorf("model %s not found in lineage graph", modelName)
	}

	// Traverse the graph to find all downstream models
	visited := make(map[string]bool)
	downstream := make([]*ModelNode, 0)

	var traverse func(*ModelNode)
	traverse = func(node *ModelNode) {
		if visited[node.NodeID] {
			return
		}
		visited[node.NodeID] = true
		downstream = append(downstream, node)

		for _, child := range node.Children {
			traverse(child)
		}
	}

	// Start traversal from the children of the start node
	for _, child := range startNode.Children {
		traverse(child)
	}

	return downstream, nil
}

// PropagateValidationResults propagates validation results through the lineage graph
func (g *LineageGraph) PropagateValidationResults(results *ValidationResults) (*ValidationResults, error) {
	// Create a map of model name to validation status
	modelStatus := make(map[string]string)
	for _, result := range results.Results {
		// If the model already has a status and it's "failed", keep it as "failed"
		if status, exists := modelStatus[result.ModelName]; exists && status == "failed" {
			continue
		}
		modelStatus[result.ModelName] = result.Status
	}

	// Propagate failures upstream
	for modelName, status := range modelStatus {
		if status == "failed" {
			// Find the model node
			var modelNode *ModelNode
			for _, node := range g.Nodes {
				if node.Name == modelName {
					modelNode = node
					break
				}
			}

			if modelNode == nil {
				continue
			}

			// Get upstream models
			upstream, err := g.GetUpstreamModels(modelName)
			if err != nil {
				return nil, fmt.Errorf("failed to get upstream models for %s: %w", modelName, err)
			}

			// Mark upstream models as potentially affected
			for _, upstreamModel := range upstream {
				// Only mark as affected if it doesn't already have a status
				if _, exists := modelStatus[upstreamModel.Name]; !exists {
					modelStatus[upstreamModel.Name] = "affected"
				}
			}
		}
	}

	// Create new validation results with propagated statuses
	propagatedResults := &ValidationResults{
		Results: make([]ValidationResult, len(results.Results)),
		Summary: results.Summary,
	}

	// Copy original results
	copy(propagatedResults.Results, results.Results)

	// Add new results for affected models
	for modelName, status := range modelStatus {
		if status == "affected" {
			// Check if this model already has a result
			hasResult := false
			for _, result := range propagatedResults.Results {
				if result.ModelName == modelName {
					hasResult = true
					break
				}
			}

			if !hasResult {
				propagatedResults.Results = append(propagatedResults.Results, ValidationResult{
					ModelName: modelName,
					RuleName:  "lineage_propagation",
					Status:    "affected",
					Message:   "Model may be affected by failures in downstream models",
				})
			}
		}
	}

	// Update summary
	propagatedResults.Summary.TotalModels = len(modelStatus)
	propagatedResults.Summary.PassedModels = 0
	propagatedResults.Summary.FailedModels = 0

	for _, status := range modelStatus {
		if status == "passed" {
			propagatedResults.Summary.PassedModels++
		} else if status == "failed" || status == "affected" {
			propagatedResults.Summary.FailedModels++
		}
	}

	return propagatedResults, nil
}
