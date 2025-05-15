package dashboard

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nessi-dev/nessi-dev/pkg/logging"
	"github.com/nessi-dev/nessi-dev/pkg/monitoring"
	"github.com/nessi-dev/nessi-dev/pkg/security"
	"github.com/nessi-dev/nessi-dev/pkg/quality/profile"
	"github.com/nessi-dev/nessi-dev/pkg/quality/rules"
)

//go:embed templates
var templateFS embed.FS

//go:embed static
var staticFS embed.FS

// Profiler is an interface for data profiling (for dependency injection)
type Profiler interface {
	GenerateProfile() ([]*profile.Profile, error)
}

// RuleValidator is an interface for rule validation (for dependency injection)
type RuleValidator interface {
	Validate(record interface{}) []rules.ValidationError
}

// Dashboard represents a monitoring dashboard
// Now supports dependency injection for Profiler and RuleValidator
// so tests can inject fast mocks.
type Dashboard struct {
	monitor      *monitoring.Monitor
	templates    *template.Template
	listenAddr   string
	authManager  *security.AuthManager
	certManager  *security.CertManager
	secureMode   bool
	mux          *http.ServeMux

	profiler        Profiler
	ruleValidator   RuleValidator
	rcaClient       RCAClient
	freshnessManager FreshnessManager
}

// DashboardOptions represents dashboard configuration options
type DashboardOptions struct {
	ListenAddr  string
	AuthManager *security.AuthManager
	CertManager *security.CertManager
	SecureMode  bool

	// Dependency injection for tests
	Profiler        Profiler
	RuleValidator   RuleValidator
	RCAClient       RCAClient
	FreshnessManager FreshnessManager
}

// New creates a new Dashboard instance
func New(monitor *monitoring.Monitor, options DashboardOptions) (*Dashboard, error) {
	// Parse templates
	templates, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	// Use injected profiler/ruleValidator if provided, else default to nil (handlers must check!)
	dash := &Dashboard{
		monitor:      monitor,
		templates:    templates,
		listenAddr:   options.ListenAddr,
		authManager:  options.AuthManager,
		certManager:  options.CertManager,
		secureMode:   options.SecureMode,
		mux:          http.NewServeMux(),
		profiler:     options.Profiler,
		ruleValidator: options.RuleValidator,
		rcaClient:     options.RCAClient,
		freshnessManager: options.FreshnessManager,
	}
	return dash, nil
}

// Start starts the dashboard server
func (d *Dashboard) Start() error {
	// Set up routes
	d.mux = http.NewServeMux()
	
	// Static files
	d.mux.Handle("/static/", http.FileServer(http.FS(staticFS)))
	
	// Public routes
	d.mux.HandleFunc("/", d.handleIndex)
	d.mux.HandleFunc("/health", d.handleHealth)
	d.mux.HandleFunc("/data-quality", d.handleDataQualityDashboard)
	d.mux.HandleFunc("/alerts", d.handleAlertsDashboard)
	d.mux.HandleFunc("/rca", d.handleRcaDashboard)
	d.mux.HandleFunc("/freshness", d.handleFreshnessDashboard)
	
	// Authentication routes
	if d.authManager != nil {
		d.mux.HandleFunc("/login", d.handleLogin)
		d.mux.HandleFunc("/auth/login", d.handleAPILogin)
	}
	
	// Protected routes
	if d.authManager != nil {
		// Apply authentication middleware to API routes
		d.mux.Handle("/api/metrics", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleMetrics)))
		d.mux.Handle("/api/alerts", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleAlerts)))
		d.mux.Handle("/api/alerts/rules", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleAlertRulesAPI)))
		
		// RCA API routes
		d.mux.Handle("/api/rca/", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleRcaAPI)))
		d.mux.Handle("/api/rca/recent", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleRcaAPI)))
		d.mux.Handle("/api/rca/analyze", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleRcaAPI)))
		d.mux.Handle("/api/rca/insights", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleRcaAPI)))
		d.mux.Handle("/api/rca/export", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleRcaAPI)))
		
		// Freshness API routes
		d.mux.Handle("/api/freshness/status", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleFreshnessStatusAPI)))
		d.mux.Handle("/api/freshness/sla", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleFreshnessSLAAPI)))
		d.mux.Handle("/api/freshness/trends", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleFreshnessTrendsAPI)))
		d.mux.Handle("/api/freshness/export", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleFreshnessExportAPI)))
		d.mux.Handle("/api/alerts/{id}/acknowledge", d.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertAction(w, r, "acknowledge")
		})))
		d.mux.Handle("/api/alerts/{id}/resolve", d.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertAction(w, r, "resolve")
		})))
		d.mux.Handle("/api/alerts/{id}/silence", d.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertAction(w, r, "silence")
		})))
		d.mux.Handle("/api/alerts/rules/{id}/enable", d.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertRuleAction(w, r, "enable")
		})))
		d.mux.Handle("/api/alerts/rules/{id}/disable", d.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertRuleAction(w, r, "disable")
		})))
		d.mux.Handle("/api/alerts/rules/{id}", d.authManager.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				d.handleAlertRuleAction(w, r, "delete")
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})))
		d.mux.Handle("/api/export", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleExport)))
		
		// Data quality routes with authentication
		d.mux.Handle("/api/profiles", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleGetProfiles)))
		d.mux.Handle("/api/profiles/summary", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleGetProfileSummaries)))
		d.mux.Handle("/api/rules", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleGetRules)))
		d.mux.Handle("/api/rules/validate", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleValidateRules)))
		d.mux.Handle("/api/rules/history", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleGetRuleHistory)))
		d.mux.Handle("/api/rules/trends", d.authManager.AuthMiddleware(http.HandlerFunc(d.handleGetExecutionTrends)))
		
		// Admin routes
		adminHandler := d.authManager.RoleMiddleware(security.RoleAdmin)
		d.mux.Handle("/admin/users", adminHandler(http.HandlerFunc(d.handleUsers)))
	} else {
		// No authentication, routes are public
		d.mux.HandleFunc("/api/metrics", d.handleMetrics)
		d.mux.HandleFunc("/api/alerts", d.handleAlertsAPI)
		d.mux.HandleFunc("/api/alerts/rules", d.handleAlertRulesAPI)
		d.mux.HandleFunc("/api/alerts/{id}/acknowledge", func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertAction(w, r, "acknowledge")
		})
		d.mux.HandleFunc("/api/alerts/{id}/resolve", func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertAction(w, r, "resolve")
		})
		d.mux.HandleFunc("/api/alerts/{id}/silence", func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertAction(w, r, "silence")
		})
		d.mux.HandleFunc("/api/alerts/rules/{id}/enable", func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertRuleAction(w, r, "enable")
		})
		d.mux.HandleFunc("/api/alerts/rules/{id}/disable", func(w http.ResponseWriter, r *http.Request) {
			d.handleAlertRuleAction(w, r, "disable")
		})
		d.mux.HandleFunc("/api/alerts/rules/{id}", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				d.handleAlertRuleAction(w, r, "delete")
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})
		d.mux.HandleFunc("/api/export", d.handleExport)
		
		// Data quality routes without authentication
		d.mux.HandleFunc("/api/profiles", d.handleGetProfiles)
		d.mux.HandleFunc("/api/profiles/summary", d.handleGetProfileSummaries)
		d.mux.HandleFunc("/api/rules", d.handleGetRules)
		d.mux.HandleFunc("/api/rules/validate", d.handleValidateRules)
		d.mux.HandleFunc("/api/rules/history", d.handleGetRuleHistory)
		d.mux.HandleFunc("/api/rules/trends", d.handleGetExecutionTrends)
	}
	
	// Register data quality handlers
	d.registerDataQualityHandlers()
	
	// Start server
	logging.Info(fmt.Sprintf("Starting dashboard server on %s", d.listenAddr))
	
	// Use HTTPS if secure mode is enabled and cert manager is available
	if d.secureMode && d.certManager != nil {
		logging.Info("Starting dashboard in secure mode (HTTPS)")
		return d.certManager.StartHTTPSServer(d.listenAddr, d.mux)
	}
	
	// Fallback to HTTP
	return http.ListenAndServe(d.listenAddr, d.mux)
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read metrics response", err)
		return
	}
	
	if _, err := w.Write(respBody); err != nil {
		logging.Error("Failed to write metrics response", err)
	}
}

// handleAlertsAPI handles alerts API requests
func (d *Dashboard) handleAlertsAPI(w http.ResponseWriter, r *http.Request) {
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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read alerts response", err)
		return
	}
	
	if _, err := w.Write(respBody); err != nil {
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
	
	// Clean up temporary file
	if isTestMode() {
		// In test mode, clean up immediately to avoid hanging tests
		os.RemoveAll(tempDir)
	} else {
		// In production mode, use a delay to keep files around in case of download issues
		go func() {
			time.Sleep(5 * time.Second)
			os.RemoveAll(tempDir)
		}()
	}
}

// isTestMode returns true if the code is running in a test environment
func isTestMode() bool {
	// Check if the program was started with the "go test" command
	for _, arg := range os.Args {
		if strings.Contains(arg, "test") {
			return true
		}
	}
	
	// Also check for common testing environment variables
	_, testEnvSet := os.LookupEnv("GO_TEST")
	
	// Check if the program is being run by the testing package
	if strings.HasSuffix(os.Args[0], ".test") {
		return true
	}
	
	return testEnvSet
}
