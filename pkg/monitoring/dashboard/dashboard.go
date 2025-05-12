package dashboard

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/security"
)

//go:embed templates
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

// Dashboard represents a monitoring dashboard
type Dashboard struct {
	monitor      *monitoring.Monitor
	templates    *template.Template
	listenAddr   string
	authManager  *security.AuthManager
	certManager  *security.CertManager
	secureMode   bool
}

// DashboardOptions represents dashboard configuration options
type DashboardOptions struct {
	ListenAddr  string
	AuthManager *security.AuthManager
	CertManager *security.CertManager
	SecureMode  bool
}

// New creates a new Dashboard instance
func New(monitor *monitoring.Monitor, options DashboardOptions) (*Dashboard, error) {
	// Parse templates
	templates, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Dashboard{
		monitor:     monitor,
		templates:   templates,
		listenAddr:  options.ListenAddr,
		authManager: options.AuthManager,
		certManager: options.CertManager,
		secureMode:  options.SecureMode,
	}, nil
}

// Start starts the dashboard server
func (d *Dashboard) Start() error {
	// Set up routes
	mux := http.NewServeMux()
	
	// Static files
	mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	
	// Public routes
	mux.HandleFunc("/", d.handleIndex)
	mux.HandleFunc("/health", d.handleHealth)
	
	// Authentication routes
	if d.authManager != nil {
		mux.HandleFunc("/login", d.handleLogin)
		mux.HandleFunc("/auth/login", d.handleAPILogin)
	}
	
	// Protected routes
	if d.authManager != nil {
		// Apply authentication middleware to API routes
		mux.Handle("/api/metrics", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleMetrics)))
		mux.Handle("/api/alerts", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleAlerts)))
		mux.Handle("/api/export", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleExport)))
		
		// Admin routes
		adminHandler := d.authManager.RoleMiddleware(security.RoleAdmin)
		mux.Handle("/admin/users", adminHandler(http.HandlerFunc(d.handleUsers)))
	} else {
		// No authentication, routes are public
		mux.HandleFunc("/api/metrics", d.handleMetrics)
		mux.HandleFunc("/api/alerts", d.handleAlerts)
		mux.HandleFunc("/api/export", d.handleExport)
	}
	
	// Start server
	logging.Info(fmt.Sprintf("Starting dashboard server on %s", d.listenAddr))
	
	// Use HTTPS if secure mode is enabled and cert manager is available
	if d.secureMode && d.certManager != nil {
		logging.Info("Starting dashboard in secure mode (HTTPS)")
		return d.certManager.StartHTTPSServer(d.listenAddr, mux)
	}
	
	// Fallback to HTTP
	return http.ListenAndServe(d.listenAddr, mux)
}

// handleIndex handles the dashboard index page
func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Render template
	if err := d.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleMetrics handles metrics API requests
func (d *Dashboard) handleMetrics(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Forward to metrics history endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics/history?%s", 
		d.monitor.GetMetricsPort(), r.URL.RawQuery))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get metrics: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.WriteHeader(resp.StatusCode)
	if _, err := http.MaxBytesReader(w, resp.Body, 1<<20).WriteTo(w); err != nil {
		logging.Error("Failed to write metrics response", err)
	}
}

// handleAlerts handles alerts API requests
func (d *Dashboard) handleAlerts(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Forward to alerts endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/alerts", 
		d.monitor.GetMetricsPort()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get alerts: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// Copy response
	w.WriteHeader(resp.StatusCode)
	if _, err := http.MaxBytesReader(w, resp.Body, 1<<20).WriteTo(w); err != nil {
		logging.Error("Failed to write alerts response", err)
	}
}

// handleExport handles metric export requests
func (d *Dashboard) handleExport(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	// Parse query parameters
	query := r.URL.Query()
	metricName := query.Get("metric")
	if metricName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing metric parameter",
		})
		return
	}
	
	// Parse format
	format := query.Get("format")
	if format == "" {
		format = "csv" // Default to CSV
	}
	
	// Parse time range
	start := time.Now().Add(-24 * time.Hour) // Default to last 24 hours
	end := time.Now()
	
	if startStr := query.Get("start"); startStr != "" {
		if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
			start = startTime
		}
	}
	
	if endStr := query.Get("end"); endStr != "" {
		if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
			end = endTime
		}
	}
	
	// Parse labels
	labels := make(map[string]string)
	for key, values := range query {
		if key != "metric" && key != "start" && key != "end" && key != "format" && len(values) > 0 {
			labels[key] = values[0]
		}
	}
	
	// Create temporary file for export
	tempDir, err := os.MkdirTemp("", "nessi-export-")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to create temporary directory: %v", err),
		})
		return
	}
	
	// Create export options
	exportFormat := monitoring.ExportFormatCSV
	if format == "json" {
		exportFormat = monitoring.ExportFormatJSON
	}
	
	filename := fmt.Sprintf("%s_%s.%s", metricName, time.Now().Format("20060102_150405"), format)
	exportPath := filepath.Join(tempDir, filename)
	
	options := monitoring.ExportOptions{
		Format:     exportFormat,
		OutputPath: exportPath,
		StartTime:  start,
		EndTime:    end,
		MetricName: metricName,
		Labels:     labels,
	}
	
	// Export metrics
	outputPath, err := d.monitor.ExportMetrics(options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to export metrics: %v", err),
		})
		return
	}
	
	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv")
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	
	// Serve the file
	http.ServeFile(w, r, outputPath)
	
	// Clean up temporary file (in a goroutine to not block response)
	go func() {
		time.Sleep(5 * time.Minute) // Keep file around for a while in case of download issues
		os.RemoveAll(tempDir)
	}()
}
