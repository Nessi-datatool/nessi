package model

import (
	"time"

	"github.com/google/uuid"
)

// NodeType represents the type of a lineage node
type NodeType string

const (
	// NodeTypeTable represents a table node
	NodeTypeTable NodeType = "table"
	
	// NodeTypeView represents a view node
	NodeTypeView NodeType = "view"
	
	// NodeTypeQuery represents a query node
	NodeTypeQuery NodeType = "query"
	
	// NodeTypeProcess represents a process node (ETL job, script, etc.)
	NodeTypeProcess NodeType = "process"
	
	// NodeTypeDataset represents a generic dataset node
	NodeTypeDataset NodeType = "dataset"
	
	// NodeTypeAPI represents an API node
	NodeTypeAPI NodeType = "api"
	
	// NodeTypeApplication represents an application node
	NodeTypeApplication NodeType = "application"
)

// EdgeType represents the type of a lineage edge
type EdgeType string

const (
	// EdgeTypeRead represents a read operation
	EdgeTypeRead EdgeType = "read"
	
	// EdgeTypeWrite represents a write operation
	EdgeTypeWrite EdgeType = "write"
	
	// EdgeTypeTransform represents a transformation operation
	EdgeTypeTransform EdgeType = "transform"
	
	// EdgeTypeAggregate represents an aggregation operation
	EdgeTypeAggregate EdgeType = "aggregate"
	
	// EdgeTypeJoin represents a join operation
	EdgeTypeJoin EdgeType = "join"
	
	// EdgeTypeFilter represents a filter operation
	EdgeTypeFilter EdgeType = "filter"
)

// Node represents a node in a lineage graph
type Node struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Type        NodeType          `json:"type"`
	Description string            `json:"description"`
	Catalog     string            `json:"catalog,omitempty"`
	Database    string            `json:"database,omitempty"`
	Schema      string            `json:"schema,omitempty"`
	Table       string            `json:"table,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     int64             `json:"version"`
}

// NewNode creates a new lineage node
func NewNode(name string, nodeType NodeType) *Node {
	now := time.Now()
	return &Node{
		ID:          uuid.New().String(),
		Name:        name,
		Type:        nodeType,
		Properties:  make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}
}

// Edge represents an edge in a lineage graph
type Edge struct {
	ID          string            `json:"id"`
	SourceID    string            `json:"source_id"`
	TargetID    string            `json:"target_id"`
	Type        EdgeType          `json:"type"`
	Description string            `json:"description"`
	Properties  map[string]string `json:"properties,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     int64             `json:"version"`
}

// NewEdge creates a new lineage edge
func NewEdge(sourceID, targetID string, edgeType EdgeType) *Edge {
	now := time.Now()
	return &Edge{
		ID:          uuid.New().String(),
		SourceID:    sourceID,
		TargetID:    targetID,
		Type:        edgeType,
		Properties:  make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}
}

// LineageGraph represents a complete lineage graph
type LineageGraph struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Nodes       []*Node           `json:"nodes"`
	Edges       []*Edge           `json:"edges"`
	Properties  map[string]string `json:"properties,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     int64             `json:"version"`
}

// NewLineageGraph creates a new lineage graph
func NewLineageGraph(name string) *LineageGraph {
	now := time.Now()
	return &LineageGraph{
		ID:          uuid.New().String(),
		Name:        name,
		Nodes:       make([]*Node, 0),
		Edges:       make([]*Edge, 0),
		Properties:  make(map[string]string),
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
	}
}

// AddNode adds a node to the lineage graph
func (g *LineageGraph) AddNode(node *Node) {
	g.Nodes = append(g.Nodes, node)
	g.UpdatedAt = time.Now()
}

// AddEdge adds an edge to the lineage graph
func (g *LineageGraph) AddEdge(edge *Edge) {
	g.Edges = append(g.Edges, edge)
	g.UpdatedAt = time.Now()
}

// GetNode returns a node by ID
func (g *LineageGraph) GetNode(id string) *Node {
	for _, node := range g.Nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

// GetEdge returns an edge by ID
func (g *LineageGraph) GetEdge(id string) *Edge {
	for _, edge := range g.Edges {
		if edge.ID == id {
			return edge
		}
	}
	return nil
}

// GetNodesByType returns all nodes of a specific type
func (g *LineageGraph) GetNodesByType(nodeType NodeType) []*Node {
	var nodes []*Node
	for _, node := range g.Nodes {
		if node.Type == nodeType {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

// GetEdgesByType returns all edges of a specific type
func (g *LineageGraph) GetEdgesByType(edgeType EdgeType) []*Edge {
	var edges []*Edge
	for _, edge := range g.Edges {
		if edge.Type == edgeType {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetUpstreamNodes returns all nodes that are upstream of the given node
func (g *LineageGraph) GetUpstreamNodes(nodeID string) []*Node {
	var upstreamNodes []*Node
	
	// Find all edges where the target is the given node
	for _, edge := range g.Edges {
		if edge.TargetID == nodeID {
			// Add the source node to the upstream nodes
			sourceNode := g.GetNode(edge.SourceID)
			if sourceNode != nil {
				upstreamNodes = append(upstreamNodes, sourceNode)
			}
		}
	}
	
	return upstreamNodes
}

// GetDownstreamNodes returns all nodes that are downstream of the given node
func (g *LineageGraph) GetDownstreamNodes(nodeID string) []*Node {
	var downstreamNodes []*Node
	
	// Find all edges where the source is the given node
	for _, edge := range g.Edges {
		if edge.SourceID == nodeID {
			// Add the target node to the downstream nodes
			targetNode := g.GetNode(edge.TargetID)
			if targetNode != nil {
				downstreamNodes = append(downstreamNodes, targetNode)
			}
		}
	}
	
	return downstreamNodes
}

// GetUpstreamEdges returns all edges that are upstream of the given node
func (g *LineageGraph) GetUpstreamEdges(nodeID string) []*Edge {
	var upstreamEdges []*Edge
	
	// Find all edges where the target is the given node
	for _, edge := range g.Edges {
		if edge.TargetID == nodeID {
			upstreamEdges = append(upstreamEdges, edge)
		}
	}
	
	return upstreamEdges
}

// GetDownstreamEdges returns all edges that are downstream of the given node
func (g *LineageGraph) GetDownstreamEdges(nodeID string) []*Edge {
	var downstreamEdges []*Edge
	
	// Find all edges where the source is the given node
	for _, edge := range g.Edges {
		if edge.SourceID == nodeID {
			downstreamEdges = append(downstreamEdges, edge)
		}
	}
	
	return downstreamEdges
}

// LineageSnapshot represents a snapshot of a lineage graph at a specific point in time
type LineageSnapshot struct {
	ID        string       `json:"id"`
	GraphID   string       `json:"graph_id"`
	Graph     *LineageGraph `json:"graph"`
	Timestamp time.Time    `json:"timestamp"`
	Version   int64        `json:"version"`
	Comment   string       `json:"comment,omitempty"`
}

// NewLineageSnapshot creates a new lineage snapshot
func NewLineageSnapshot(graph *LineageGraph, comment string) *LineageSnapshot {
	return &LineageSnapshot{
		ID:        uuid.New().String(),
		GraphID:   graph.ID,
		Graph:     graph,
		Timestamp: time.Now(),
		Version:   graph.Version,
		Comment:   comment,
	}
}

// ColumnLineage represents lineage at the column level
type ColumnLineage struct {
	SourceTable  string `json:"source_table"`
	SourceColumn string `json:"source_column"`
	TargetTable  string `json:"target_table"`
	TargetColumn string `json:"target_column"`
	Transformation string `json:"transformation,omitempty"`
}

// LineageDiff represents the difference between two lineage graphs
type LineageDiff struct {
	OldGraphID     string       `json:"old_graph_id"`
	NewGraphID     string       `json:"new_graph_id"`
	AddedNodes     []*Node      `json:"added_nodes"`
	RemovedNodes   []*Node      `json:"removed_nodes"`
	ModifiedNodes  []*Node      `json:"modified_nodes"`
	AddedEdges     []*Edge      `json:"added_edges"`
	RemovedEdges   []*Edge      `json:"removed_edges"`
	ModifiedEdges  []*Edge      `json:"modified_edges"`
	CreatedAt      time.Time    `json:"created_at"`
}

// NewLineageDiff creates a new lineage diff
func NewLineageDiff(oldGraph, newGraph *LineageGraph) *LineageDiff {
	return &LineageDiff{
		OldGraphID:    oldGraph.ID,
		NewGraphID:    newGraph.ID,
		AddedNodes:    make([]*Node, 0),
		RemovedNodes:  make([]*Node, 0),
		ModifiedNodes: make([]*Node, 0),
		AddedEdges:    make([]*Edge, 0),
		RemovedEdges:  make([]*Edge, 0),
		ModifiedEdges: make([]*Edge, 0),
		CreatedAt:     time.Now(),
	}
}
