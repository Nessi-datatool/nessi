package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/nessi-dev/nessi/internal/delta"
	"github.com/nessi-dev/nessi/internal/extensions"
	"github.com/nessi-dev/nessi/internal/monitor"
	quality "github.com/nessi-dev/nessi/internal/quality"
	"github.com/nessi-dev/nessi/internal/security"
	"github.com/nessi-dev/nessi/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultShutdownTimeout = 5 * time.Second
)

// Server represents the API server
type Server struct {
	config         *Config
	server         *http.Server
	router         *gin.Engine
	extManager     *extensions.Manager
	secManager     *security.SecurityManager
	deltaConnector *delta.DeltaConnector
	qualityManager *quality.QualityManager
	monitorManager *monitor.MonitorManager
	logger         *zap.Logger
	metrics        *Metrics
	startTime      time.Time
}

// Metrics represents server metrics
type Metrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	activeRequests  *prometheus.GaugeVec
}

// Config represents server configuration
type Config struct {
	Host            string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	TLS             struct {
		Enabled  bool
		CertFile string
		KeyFile  string
	}
}

// NewServer creates a new API server
func NewServer(config *Config, extMgr *extensions.Manager, secMgr *security.SecurityManager, deltaConn *delta.DeltaConnector, qualMgr *quality.QualityManager, monMgr *monitor.MonitorManager, logger *zap.Logger) *Server {
	if config == nil {
		// Provide a basic default if no config is passed, though ideally, config comes from main
		config = &Config{
			Host:            "localhost",
			Port:            8080,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: defaultShutdownTimeout,
		}
	}

	server := &Server{
		router:         gin.New(),
		config:         config,
		extManager:     extMgr,
		secManager:     secMgr,
		deltaConnector: deltaConn,
		qualityManager: qualMgr,
		monitorManager: monMgr,
		logger:         logger,
		startTime:      time.Now(),
	}

	// Initialize metrics, if not already handled by monitorManager construction
	server.metrics = &Metrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "HTTP request duration in seconds.",
			},
			[]string{"method", "path"},
		),
		activeRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "http_active_requests",
				Help: "Number of active HTTP requests.",
			},
			[]string{"method", "path"},
		),
	}
	prometheus.MustRegister(server.metrics.requestsTotal)
	prometheus.MustRegister(server.metrics.requestDuration)
	prometheus.MustRegister(server.metrics.activeRequests)

	server.setupMiddleware()
	server.setupRoutes()
	return server
}

// setupMiddleware configures the server middleware
func (s *Server) setupMiddleware() {
	// Use security middleware
	s.router.Use(func(c *gin.Context) {
		nextHTTP := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r // Pass any request modifications back to Gin
			c.Next()
		})
		// s.secManager.Middleware expects a http.Handler and returns a http.Handler
		handler := s.secManager.Middleware(nextHTTP)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	// Add security headers middleware
	s.router.Use(func(c *gin.Context) {
		headers := s.secManager.GetSecurityHeaders()
		for key, value := range headers {
			c.Writer.Header().Set(key, value)
		}
		c.Next()
	})

	// Add logging middleware
	s.router.Use(gin.Logger())

	// Add recovery middleware
	s.router.Use(gin.Recovery())

	// Add CORS middleware
	s.router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	// Public routes (no authentication required)
	public := s.router.Group("/")
	{
		public.GET("/health", s.handleHealthCheck)
		public.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	// Protected routes (require authentication)
	api := s.router.Group("/api/v1")
	{
		// Table operations
		tables := api.Group("/tables")
		{
			tables.POST("/:path/profile", s.handleGenerateProfile)
			tables.POST("/:path/quality", s.handleQualityCheck)
			tables.GET("/:path/stats", s.handleGetStats)
		}

		// Extension management
		extensions := api.Group("/extensions")
		{
			extensions.GET("", s.handleListExtensions)
			extensions.POST("/:name", s.handleInstallExtension)
			extensions.DELETE("/:name", s.handleUninstallExtension)
			extensions.POST("/:name/enable", s.handleEnableExtension)
			extensions.POST("/:name/disable", s.handleDisableExtension)
		}
	}
}

// Start starts the API server
func (s *Server) Start() error {
	// Configure server
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	// Configure TLS if enabled
	if s.config.TLS.Enabled {
		tlsConfig := &tls.Config{
			MinVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		}
		s.server.TLSConfig = tlsConfig
	}

	// Start server in a goroutine
	go func() {
		var err error
		if s.config.TLS.Enabled {
			err = s.server.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		} else {
			err = s.server.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	s.logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	s.logger.Info("server exited properly")
	return nil
}

// handleHealthCheck handles health check requests
func (s *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleGenerateProfile handles profile generation requests
func (s *Server) handleGenerateProfile(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Profile generation not implemented yet"})
}

// handleGetProfileStatus handles requests to get profile status
func (s *Server) handleGetProfileStatus(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Profile status not implemented yet"})
}

// handleQualityCheck handles quality check requests
func (s *Server) handleQualityCheck(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"message": "Quality check not implemented yet"})
}

// handleGetStats handles requests to get table statistics
func (s *Server) handleGetStats(c *gin.Context) {
	path := c.Param("path")
	reader, err := pkg.NewReader(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create reader: %v", err)})
		return
	}
	defer reader.Close()

	stats := reader.GetStats()
	c.JSON(http.StatusOK, stats)
}

// handleListExtensions handles extension listing requests
func (s *Server) handleListExtensions(c *gin.Context) {
	extensions := s.extManager.ListExtensions()
	c.JSON(http.StatusOK, extensions)
}

// handleInstallExtension handles extension installation requests
func (s *Server) handleInstallExtension(c *gin.Context) {
	var ext extensions.Extension
	if err := c.ShouldBindJSON(&ext); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid extension payload"})
		return
	}

	if err := s.extManager.RegisterExtension(&ext); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to install extension: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Extension installed successfully"})
}

// handleUninstallExtension handles extension uninstall requests
func (s *Server) handleUninstallExtension(c *gin.Context) {
	name := c.Param("name")
	// TODO: DisableExtension only marks the extension as disabled, it does not remove it from the manager.
	// If true uninstallation (removal) is required, Manager needs a RemoveExtension or UnregisterExtension method.
	if err := s.extManager.DisableExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to uninstall extension: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Extension uninstalled successfully"})
}

// handleEnableExtension handles extension enabling requests
func (s *Server) handleEnableExtension(c *gin.Context) {
	name := c.Param("name")
	if err := s.extManager.EnableExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to enable extension: %v", err)})
		return
	}
	c.Status(http.StatusOK)
}

// handleDisableExtension handles extension disabling requests
func (s *Server) handleDisableExtension(c *gin.Context) {
	name := c.Param("name")
	if err := s.extManager.DisableExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to disable extension: %v", err)})
		return
	}
	c.Status(http.StatusOK)
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	// TODO: Implement graceful shutdown
	// This would typically involve:
	// 1. Stopping new request acceptance
	// 2. Waiting for in-flight requests to complete
	// 3. Closing all connections
	return nil
}
