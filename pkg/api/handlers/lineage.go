package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nessi-dev/nessi/pkg/lineage"
	"github.com/nessi-dev/nessi/pkg/lineage/model"
	"github.com/nessi-dev/nessi/pkg/lineage/visualization"
)

// LineageHandler handles lineage API requests
type LineageHandler struct {
	lineageService *lineage.LineageService
}

// NewLineageHandler creates a new lineage handler
func NewLineageHandler(lineageService *lineage.LineageService) *LineageHandler {
	return &LineageHandler{
		lineageService: lineageService,
	}
}

// RegisterRoutes registers lineage API routes
func (h *LineageHandler) RegisterRoutes(router *mux.Router) {
	// Graph routes
	router.HandleFunc("/lineage/graphs", h.listGraphs).Methods("GET")
	router.HandleFunc("/lineage/graphs", h.createGraph).Methods("POST")
	router.HandleFunc("/lineage/graphs/{id}", h.getGraph).Methods("GET")
	router.HandleFunc("/lineage/graphs/{id}", h.deleteGraph).Methods("DELETE")
	router.HandleFunc("/lineage/graphs/{id}/nodes", h.addNode).Methods("POST")
	router.HandleFunc("/lineage/graphs/{id}/edges", h.addEdge).Methods("POST")
	
	// Snapshot routes
	router.HandleFunc("/lineage/graphs/{id}/snapshots", h.listSnapshots).Methods("GET")
	router.HandleFunc("/lineage/graphs/{id}/snapshots", h.createSnapshot).Methods("POST")
	router.HandleFunc("/lineage/snapshots/{id}", h.getSnapshot).Methods("GET")
	router.HandleFunc("/lineage/graphs/{id}/snapshots/time/{timestamp}", h.getSnapshotAtTime).Methods("GET")
	
	// Comparison routes
	router.HandleFunc("/lineage/compare", h.compareSnapshots).Methods("POST")
	router.HandleFunc("/lineage/diffs/{id}/impact", h.analyzeImpact).Methods("GET")
	
	// Visualization routes
	router.HandleFunc("/lineage/graphs/{id}/visualize", h.visualizeGraph).Methods("GET")
	router.HandleFunc("/lineage/graphs/{id}/export/{format}", h.exportLineage).Methods("GET")
	
	// Integration routes
	router.HandleFunc("/lineage/import/catalog", h.importFromCatalog).Methods("POST")
}

// listGraphs handles GET /lineage/graphs
func (h *LineageHandler) listGraphs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	graphs, err := h.lineageService.ListGraphs(ctx)
	if err != nil {
		http.Error(w, "Failed to list graphs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, graphs)
}

// createGraph handles POST /lineage/graphs
func (h *LineageHandler) createGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	graph, err := h.lineageService.CreateGraph(ctx, req.Name, req.Description)
	if err != nil {
		http.Error(w, "Failed to create graph: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, graph)
}

// getGraph handles GET /lineage/graphs/{id}
func (h *LineageHandler) getGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	id := vars["id"]
	
	graph, err := h.lineageService.GetGraph(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get graph: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if graph == nil {
		http.Error(w, "Graph not found", http.StatusNotFound)
		return
	}
	
	respondJSON(w, graph)
}

// deleteGraph handles DELETE /lineage/graphs/{id}
func (h *LineageHandler) deleteGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	id := vars["id"]
	
	if err := h.lineageService.DeleteGraph(ctx, id); err != nil {
		http.Error(w, "Failed to delete graph: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// addNode handles POST /lineage/graphs/{id}/nodes
func (h *LineageHandler) addNode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	
	var node model.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.lineageService.AddNode(ctx, graphID, &node); err != nil {
		http.Error(w, "Failed to add node: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

// addEdge handles POST /lineage/graphs/{id}/edges
func (h *LineageHandler) addEdge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	
	var edge model.Edge
	if err := json.NewDecoder(r.Body).Decode(&edge); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	if err := h.lineageService.AddEdge(ctx, graphID, &edge); err != nil {
		http.Error(w, "Failed to add edge: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
}

// listSnapshots handles GET /lineage/graphs/{id}/snapshots
func (h *LineageHandler) listSnapshots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	
	snapshots, err := h.lineageService.ListSnapshots(ctx, graphID)
	if err != nil {
		http.Error(w, "Failed to list snapshots: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, snapshots)
}

// createSnapshot handles POST /lineage/graphs/{id}/snapshots
func (h *LineageHandler) createSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	
	var req struct {
		Comment string `json:"comment"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	snapshot, err := h.lineageService.CreateSnapshot(ctx, graphID, req.Comment)
	if err != nil {
		http.Error(w, "Failed to create snapshot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, snapshot)
}

// getSnapshot handles GET /lineage/snapshots/{id}
func (h *LineageHandler) getSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	id := vars["id"]
	
	snapshot, err := h.lineageService.GetSnapshot(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get snapshot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if snapshot == nil {
		http.Error(w, "Snapshot not found", http.StatusNotFound)
		return
	}
	
	respondJSON(w, snapshot)
}

// getSnapshotAtTime handles GET /lineage/graphs/{id}/snapshots/time/{timestamp}
func (h *LineageHandler) getSnapshotAtTime(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	timestamp := vars["timestamp"]
	
	// Parse timestamp
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		http.Error(w, "Invalid timestamp format: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	snapshot, err := h.lineageService.GetSnapshotAtTime(ctx, graphID, ts)
	if err != nil {
		http.Error(w, "Failed to get snapshot: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if snapshot == nil {
		http.Error(w, "No snapshot found at the specified time", http.StatusNotFound)
		return
	}
	
	respondJSON(w, snapshot)
}

// compareSnapshots handles POST /lineage/compare
func (h *LineageHandler) compareSnapshots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		SnapshotID1 string `json:"snapshot_id_1"`
		SnapshotID2 string `json:"snapshot_id_2"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	diff, err := h.lineageService.CompareSnapshots(ctx, req.SnapshotID1, req.SnapshotID2)
	if err != nil {
		http.Error(w, "Failed to compare snapshots: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, diff)
}

// analyzeImpact handles GET /lineage/diffs/{id}/impact
func (h *LineageHandler) analyzeImpact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	diffID := vars["id"]
	
	impact, err := h.lineageService.AnalyzeImpact(ctx, diffID)
	if err != nil {
		http.Error(w, "Failed to analyze impact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, impact)
}

// visualizeGraph handles GET /lineage/graphs/{id}/visualize
func (h *LineageHandler) visualizeGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	
	// Parse visualization options
	options := visualization.NewDefaultVisualizationOptions()
	
	// Parse query parameters
	if format := r.URL.Query().Get("format"); format != "" {
		options.Format = visualization.GraphFormat(format)
	}
	
	if width := r.URL.Query().Get("width"); width != "" {
		if widthVal, err := strconv.Atoi(width); err == nil && widthVal > 0 {
			options.Width = widthVal
		}
	}
	
	if height := r.URL.Query().Get("height"); height != "" {
		if heightVal, err := strconv.Atoi(height); err == nil && heightVal > 0 {
			options.Height = heightVal
		}
	}
	
	if focusNodeID := r.URL.Query().Get("focus_node_id"); focusNodeID != "" {
		options.FocusNodeID = focusNodeID
	}
	
	if maxDepth := r.URL.Query().Get("max_depth"); maxDepth != "" {
		if maxDepthVal, err := strconv.Atoi(maxDepth); err == nil && maxDepthVal > 0 {
			options.MaxDepth = maxDepthVal
		}
	}
	
	if layout := r.URL.Query().Get("layout"); layout != "" {
		options.LayoutAlgorithm = layout
	}
	
	if title := r.URL.Query().Get("title"); title != "" {
		options.Title = title
	}
	
	// Render graph
	data, err := h.lineageService.RenderGraph(ctx, graphID, options)
	if err != nil {
		http.Error(w, "Failed to visualize graph: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Set content type based on format
	switch options.Format {
	case visualization.FormatHTML:
		w.Header().Set("Content-Type", "text/html")
	case visualization.FormatJSON:
		w.Header().Set("Content-Type", "application/json")
	case visualization.FormatSVG:
		w.Header().Set("Content-Type", "image/svg+xml")
	case visualization.FormatDOT:
		w.Header().Set("Content-Type", "text/plain")
	}
	
	w.Write(data)
}

// exportLineage handles GET /lineage/graphs/{id}/export/{format}
func (h *LineageHandler) exportLineage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	vars := mux.Vars(r)
	graphID := vars["id"]
	format := vars["format"]
	
	data, err := h.lineageService.ExportLineage(ctx, graphID, format)
	if err != nil {
		http.Error(w, "Failed to export lineage: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Set content type and disposition
	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=lineage.json")
	case "dot":
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=lineage.dot")
	case "svg":
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Content-Disposition", "attachment; filename=lineage.svg")
	case "html":
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Disposition", "attachment; filename=lineage.html")
	}
	
	w.Write(data)
}

// importFromCatalog handles POST /lineage/import/catalog
func (h *LineageHandler) importFromCatalog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var req struct {
		CatalogName string `json:"catalog_name"`
		Database    string `json:"database"`
		Table       string `json:"table"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	graph, err := h.lineageService.ImportFromCatalog(ctx, req.CatalogName, req.Database, req.Table)
	if err != nil {
		http.Error(w, "Failed to import from catalog: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	respondJSON(w, graph)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	}
}
