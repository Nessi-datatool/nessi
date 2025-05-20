package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// LineageStorage defines the interface for storing and retrieving lineage graphs
type LineageStorage interface {
	// SaveGraph saves a lineage graph
	SaveGraph(ctx context.Context, graph *LineageGraph) error

	// GetGraph retrieves a lineage graph by ID
	GetGraph(ctx context.Context, id string) (*LineageGraph, error)

	// ListGraphs lists all lineage graphs
	ListGraphs(ctx context.Context) ([]*LineageGraph, error)

	// DeleteGraph deletes a lineage graph by ID
	DeleteGraph(ctx context.Context, id string) error

	// SaveSnapshot saves a lineage snapshot
	SaveSnapshot(ctx context.Context, snapshot *LineageSnapshot) error

	// GetSnapshot retrieves a lineage snapshot by ID
	GetSnapshot(ctx context.Context, id string) (*LineageSnapshot, error)

	// ListSnapshots lists all snapshots for a graph
	ListSnapshots(ctx context.Context, graphID string) ([]*LineageSnapshot, error)

	// GetSnapshotAtTime retrieves the snapshot closest to the given time
	GetSnapshotAtTime(ctx context.Context, graphID string, timestamp time.Time) (*LineageSnapshot, error)

	// SaveDiff saves a lineage diff
	SaveDiff(ctx context.Context, diff *LineageDiff) error

	// GetDiff retrieves a lineage diff between two graph IDs
	GetDiff(ctx context.Context, oldGraphID, newGraphID string) (*LineageDiff, error)
}

// FileLineageStorage implements LineageStorage using the file system
type FileLineageStorage struct {
	baseDir string
}

// NewFileLineageStorage creates a new file-based lineage storage
func NewFileLineageStorage(baseDir string) (*FileLineageStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// Create subdirectories
	graphsDir := filepath.Join(baseDir, "graphs")
	snapshotsDir := filepath.Join(baseDir, "snapshots")
	diffsDir := filepath.Join(baseDir, "diffs")

	if err := os.MkdirAll(graphsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create graphs directory: %w", err)
	}

	if err := os.MkdirAll(snapshotsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create snapshots directory: %w", err)
	}

	if err := os.MkdirAll(diffsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create diffs directory: %w", err)
	}

	return &FileLineageStorage{
		baseDir: baseDir,
	}, nil
}

// SaveGraph saves a lineage graph to a file
func (s *FileLineageStorage) SaveGraph(ctx context.Context, graph *LineageGraph) error {
	// Update the graph's updated time
	graph.UpdatedAt = time.Now()

	// Convert graph to JSON
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal graph: %w", err)
	}

	// Write to file
	filename := filepath.Join(s.baseDir, "graphs", graph.ID+".json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write graph file: %w", err)
	}

	return nil
}

// GetGraph retrieves a lineage graph from a file
func (s *FileLineageStorage) GetGraph(ctx context.Context, id string) (*LineageGraph, error) {
	// Read from file
	filename := filepath.Join(s.baseDir, "graphs", id+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read graph file: %w", err)
	}

	// Parse JSON
	var graph LineageGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph: %w", err)
	}

	return &graph, nil
}

// ListGraphs lists all lineage graphs
func (s *FileLineageStorage) ListGraphs(ctx context.Context) ([]*LineageGraph, error) {
	// Get all graph files
	pattern := filepath.Join(s.baseDir, "graphs", "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list graph files: %w", err)
	}

	// Read each graph
	var graphs []*LineageGraph
	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err != nil {
			return nil, fmt.Errorf("failed to read graph file %s: %w", match, err)
		}

		var graph LineageGraph
		if err := json.Unmarshal(data, &graph); err != nil {
			return nil, fmt.Errorf("failed to unmarshal graph from %s: %w", match, err)
		}

		graphs = append(graphs, &graph)
	}

	// Sort by updated time (newest first)
	sort.Slice(graphs, func(i, j int) bool {
		return graphs[i].UpdatedAt.After(graphs[j].UpdatedAt)
	})

	return graphs, nil
}

// DeleteGraph deletes a lineage graph
func (s *FileLineageStorage) DeleteGraph(ctx context.Context, id string) error {
	// Delete graph file
	filename := filepath.Join(s.baseDir, "graphs", id+".json")
	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("failed to delete graph file: %w", err)
	}

	// Delete associated snapshots
	snapshotPattern := filepath.Join(s.baseDir, "snapshots", id+"_*.json")
	snapshotMatches, err := filepath.Glob(snapshotPattern)
	if err != nil {
		return fmt.Errorf("failed to list snapshot files: %w", err)
	}

	for _, match := range snapshotMatches {
		if err := os.Remove(match); err != nil {
			return fmt.Errorf("failed to delete snapshot file %s: %w", match, err)
		}
	}

	// Delete associated diffs
	diffPattern := filepath.Join(s.baseDir, "diffs", id+"_*.json")
	diffMatches, err := filepath.Glob(diffPattern)
	if err != nil {
		return fmt.Errorf("failed to list diff files: %w", err)
	}

	for _, match := range diffMatches {
		if err := os.Remove(match); err != nil {
			return fmt.Errorf("failed to delete diff file %s: %w", match, err)
		}
	}

	return nil
}

// SaveSnapshot saves a lineage snapshot
func (s *FileLineageStorage) SaveSnapshot(ctx context.Context, snapshot *LineageSnapshot) error {
	// Convert snapshot to JSON
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	// Write to file
	filename := filepath.Join(s.baseDir, "snapshots", snapshot.GraphID+"_"+snapshot.ID+".json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write snapshot file: %w", err)
	}

	return nil
}

// GetSnapshot retrieves a lineage snapshot
func (s *FileLineageStorage) GetSnapshot(ctx context.Context, id string) (*LineageSnapshot, error) {
	// Find the snapshot file
	pattern := filepath.Join(s.baseDir, "snapshots", "*_"+id+".json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot file: %w", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("snapshot not found: %s", id)
	}

	// Read from file
	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot file: %w", err)
	}

	// Parse JSON
	var snapshot LineageSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal snapshot: %w", err)
	}

	return &snapshot, nil
}

// ListSnapshots lists all snapshots for a graph
func (s *FileLineageStorage) ListSnapshots(ctx context.Context, graphID string) ([]*LineageSnapshot, error) {
	// Get all snapshot files for the graph
	pattern := filepath.Join(s.baseDir, "snapshots", graphID+"_*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to list snapshot files: %w", err)
	}

	// Read each snapshot
	var snapshots []*LineageSnapshot
	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err != nil {
			return nil, fmt.Errorf("failed to read snapshot file %s: %w", match, err)
		}

		var snapshot LineageSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return nil, fmt.Errorf("failed to unmarshal snapshot from %s: %w", match, err)
		}

		snapshots = append(snapshots, &snapshot)
	}

	// Sort by timestamp (newest first)
	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.After(snapshots[j].Timestamp)
	})

	return snapshots, nil
}

// GetSnapshotAtTime retrieves the snapshot closest to the given time
func (s *FileLineageStorage) GetSnapshotAtTime(ctx context.Context, graphID string, timestamp time.Time) (*LineageSnapshot, error) {
	// Get all snapshots for the graph
	snapshots, err := s.ListSnapshots(ctx, graphID)
	if err != nil {
		return nil, err
	}

	if len(snapshots) == 0 {
		return nil, fmt.Errorf("no snapshots found for graph %s", graphID)
	}

	// Find the snapshot closest to the given time
	var closestSnapshot *LineageSnapshot
	var closestDiff time.Duration

	for _, snapshot := range snapshots {
		diff := timestamp.Sub(snapshot.Timestamp)
		if diff < 0 {
			diff = -diff // Get absolute value
		}

		if closestSnapshot == nil || diff < closestDiff {
			closestSnapshot = snapshot
			closestDiff = diff
		}
	}

	return closestSnapshot, nil
}

// SaveDiff saves a lineage diff
func (s *FileLineageStorage) SaveDiff(ctx context.Context, diff *LineageDiff) error {
	// Convert diff to JSON
	data, err := json.MarshalIndent(diff, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal diff: %w", err)
	}

	// Write to file
	filename := filepath.Join(s.baseDir, "diffs", diff.OldGraphID+"_"+diff.NewGraphID+".json")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write diff file: %w", err)
	}

	return nil
}

// GetDiff retrieves a lineage diff between two graph IDs
func (s *FileLineageStorage) GetDiff(ctx context.Context, oldGraphID, newGraphID string) (*LineageDiff, error) {
	// Read from file
	filename := filepath.Join(s.baseDir, "diffs", oldGraphID+"_"+newGraphID+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read diff file: %w", err)
	}

	// Parse JSON
	var diff LineageDiff
	if err := json.Unmarshal(data, &diff); err != nil {
		return nil, fmt.Errorf("failed to unmarshal diff: %w", err)
	}

	return &diff, nil
}
