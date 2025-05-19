package lineage

import (
	"context"
	"fmt"
	"time"

	"github.com/nessi-dev/nessi/pkg/api/types"
	"github.com/nessi-dev/nessi/pkg/lineage/comparison"
	"github.com/nessi-dev/nessi/pkg/lineage/model"
	"github.com/nessi-dev/nessi/pkg/lineage/visualization"
)

// LineageService provides high-level operations for data lineage
type LineageService struct {
	storage        model.LineageStorage
	renderer       *visualization.GraphRenderer
	diffGenerator  *comparison.DiffGenerator
	diffAnalyzer   *comparison.DiffAnalyzer
	catalogManager *types.CatalogManager
}

// NewLineageService creates a new lineage service
func NewLineageService(
	storage model.LineageStorage,
	renderer *visualization.GraphRenderer,
	catalogManager *types.CatalogManager,
) *LineageService {
	return &LineageService{
		storage:        storage,
		renderer:       renderer,
		diffGenerator:  comparison.NewDiffGenerator(),
		diffAnalyzer:   comparison.NewDiffAnalyzer(),
		catalogManager: catalogManager,
	}
}

// CreateGraph creates a new lineage graph
func (s *LineageService) CreateGraph(ctx context.Context, name, description string) (*model.LineageGraph, error) {
	graph := model.NewLineageGraph(name)
	graph.Description = description

	if err := s.storage.SaveGraph(ctx, graph); err != nil {
		return nil, fmt.Errorf("failed to save graph: %w", err)
	}

	return graph, nil
}

// GetGraph retrieves a lineage graph by ID
func (s *LineageService) GetGraph(ctx context.Context, id string) (*model.LineageGraph, error) {
	return s.storage.GetGraph(ctx, id)
}

// ListGraphs lists all lineage graphs
func (s *LineageService) ListGraphs(ctx context.Context) ([]*model.LineageGraph, error) {
	return s.storage.ListGraphs(ctx)
}

// DeleteGraph deletes a lineage graph by ID
func (s *LineageService) DeleteGraph(ctx context.Context, id string) error {
	return s.storage.DeleteGraph(ctx, id)
}

// AddNode adds a node to a lineage graph
func (s *LineageService) AddNode(ctx context.Context, graphID string, node *model.Node) error {
	graph, err := s.storage.GetGraph(ctx, graphID)
	if err != nil {
		return fmt.Errorf("failed to get graph: %w", err)
	}

	graph.AddNode(node)

	if err := s.storage.SaveGraph(ctx, graph); err != nil {
		return fmt.Errorf("failed to save graph: %w", err)
	}

	return nil
}

// AddEdge adds an edge to a lineage graph
func (s *LineageService) AddEdge(ctx context.Context, graphID string, edge *model.Edge) error {
	graph, err := s.storage.GetGraph(ctx, graphID)
	if err != nil {
		return fmt.Errorf("failed to get graph: %w", err)
	}

	graph.AddEdge(edge)

	if err := s.storage.SaveGraph(ctx, graph); err != nil {
		return fmt.Errorf("failed to save graph: %w", err)
	}

	return nil
}

// CreateSnapshot creates a snapshot of a lineage graph
func (s *LineageService) CreateSnapshot(ctx context.Context, graphID, comment string) (*model.LineageSnapshot, error) {
	graph, err := s.storage.GetGraph(ctx, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph: %w", err)
	}

	snapshot := model.NewLineageSnapshot(graph, comment)

	if err := s.storage.SaveSnapshot(ctx, snapshot); err != nil {
		return nil, fmt.Errorf("failed to save snapshot: %w", err)
	}

	return snapshot, nil
}

// GetSnapshot retrieves a lineage snapshot by ID
func (s *LineageService) GetSnapshot(ctx context.Context, id string) (*model.LineageSnapshot, error) {
	return s.storage.GetSnapshot(ctx, id)
}

// ListSnapshots lists all snapshots for a graph
func (s *LineageService) ListSnapshots(ctx context.Context, graphID string) ([]*model.LineageSnapshot, error) {
	return s.storage.ListSnapshots(ctx, graphID)
}

// GetSnapshotAtTime retrieves the snapshot closest to the given time
func (s *LineageService) GetSnapshotAtTime(ctx context.Context, graphID string, timestamp time.Time) (*model.LineageSnapshot, error) {
	return s.storage.GetSnapshotAtTime(ctx, graphID, timestamp)
}

// CompareSnapshots compares two lineage snapshots
func (s *LineageService) CompareSnapshots(ctx context.Context, snapshotID1, snapshotID2 string) (*model.LineageDiff, error) {
	snapshot1, err := s.storage.GetSnapshot(ctx, snapshotID1)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot 1: %w", err)
	}

	snapshot2, err := s.storage.GetSnapshot(ctx, snapshotID2)
	if err != nil {
		return nil, fmt.Errorf("failed to get snapshot 2: %w", err)
	}

	diff := s.diffGenerator.GenerateDiff(snapshot1.Graph, snapshot2.Graph)

	if err := s.storage.SaveDiff(ctx, diff); err != nil {
		return nil, fmt.Errorf("failed to save diff: %w", err)
	}

	return diff, nil
}

// AnalyzeImpact analyzes the impact of changes in a lineage diff
func (s *LineageService) AnalyzeImpact(ctx context.Context, diffID string) (*comparison.ImpactAnalysis, error) {
	// Get the diff
	oldGraphID, newGraphID, err := parseDiffID(diffID)
	if err != nil {
		return nil, fmt.Errorf("invalid diff ID: %w", err)
	}

	diff, err := s.storage.GetDiff(ctx, oldGraphID, newGraphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}

	// Get the new graph
	graph, err := s.storage.GetGraph(ctx, newGraphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph: %w", err)
	}

	// Analyze impact
	impact := s.diffAnalyzer.AnalyzeImpact(diff, graph)

	return impact, nil
}

// parseDiffID parses a diff ID into old and new graph IDs
func parseDiffID(diffID string) (string, string, error) {
	// In a real implementation, this would parse the diff ID
	// based on the format used by the storage implementation
	return "", "", fmt.Errorf("not implemented")
}

// RenderGraph renders a lineage graph
func (s *LineageService) RenderGraph(ctx context.Context, graphID string, options *visualization.VisualizationOptions) ([]byte, error) {
	graph, err := s.storage.GetGraph(ctx, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph: %w", err)
	}

	return s.renderer.RenderGraph(graph, options)
}

// ImportFromCatalog imports lineage information from a data catalog
func (s *LineageService) ImportFromCatalog(ctx context.Context, catalogName, database, table string) (*model.LineageGraph, error) {
	// Get the catalog
	cat, err := s.catalogManager.GetCatalog(catalogName)
	if err != nil {
		return nil, fmt.Errorf("failed to get catalog: %w", err)
	}

	// Get lineage information from the catalog
	lineageInfo, err := cat.GetTableLineage(ctx, database, table)
	if err != nil {
		return nil, fmt.Errorf("failed to get lineage information: %w", err)
	}

	// Create a new lineage graph
	graph := model.NewLineageGraph(fmt.Sprintf("%s.%s.%s", catalogName, database, table))
	graph.Description = fmt.Sprintf("Lineage for %s.%s.%s", catalogName, database, table)

	// Add the target table as a node
	targetNode := model.NewNode(table, model.NodeTypeTable)
	targetNode.Catalog = catalogName
	targetNode.Database = database
	targetNode.Table = table
	graph.AddNode(targetNode)

	// Add upstream nodes and edges
	for _, ref := range lineageInfo.Upstream {
		upstreamNode := model.NewNode(ref.Table, model.NodeTypeTable)
		upstreamNode.Catalog = ref.Catalog
		upstreamNode.Database = ref.Database
		upstreamNode.Table = ref.Table
		upstreamNode.Properties = ref.Properties
		graph.AddNode(upstreamNode)

		edge := model.NewEdge(upstreamNode.ID, targetNode.ID, model.EdgeTypeRead)
		graph.AddEdge(edge)
	}

	// Add downstream nodes and edges
	for _, ref := range lineageInfo.Downstream {
		downstreamNode := model.NewNode(ref.Table, model.NodeTypeTable)
		downstreamNode.Catalog = ref.Catalog
		downstreamNode.Database = ref.Database
		downstreamNode.Table = ref.Table
		downstreamNode.Properties = ref.Properties
		graph.AddNode(downstreamNode)

		edge := model.NewEdge(targetNode.ID, downstreamNode.ID, model.EdgeTypeWrite)
		graph.AddEdge(edge)
	}

	// Save the graph
	if err := s.storage.SaveGraph(ctx, graph); err != nil {
		return nil, fmt.Errorf("failed to save graph: %w", err)
	}

	return graph, nil
}

// ExportLineage exports lineage information to various formats
func (s *LineageService) ExportLineage(ctx context.Context, graphID string, format string) ([]byte, error) {
	graph, err := s.storage.GetGraph(ctx, graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph: %w", err)
	}

	options := visualization.NewDefaultVisualizationOptions()

	switch format {
	case "json":
		options.Format = visualization.FormatJSON
	case "dot":
		options.Format = visualization.FormatDOT
	case "svg":
		options.Format = visualization.FormatSVG
	case "html":
		options.Format = visualization.FormatHTML
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return s.renderer.RenderGraph(graph, options)
}
