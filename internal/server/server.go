// Implement an HTTP server for Nessi.dev with these requirements:
// 1. Define a Server struct with:
//    - Configuration references
//    - HTTP server instance
//    - Router (using Gin)
//    - Extension manager reference
//    - Security manager reference
// 2. New(cfg *config.Config, extManager *extensions.Manager) *Server constructor
// 3. Methods:
//    - Start() error - Start the server
//    - Stop() error - Gracefully stop the server
//    - setupRoutes() - Configure API routes
//    - setupMiddleware() - Set up common middleware
// 4. API endpoints grouped by feature:
//    - /api/v1/tables - Delta table operations
//    - /api/v1/quality - Data quality operations
//    - /api/v1/monitor - Monitoring endpoints
//    - /api/v1/extensions - Extension management
//    - /metrics - Prometheus metrics endpoint
// 5. Middleware components:
//    - Authentication middleware
//    - Logging middleware
//    - Rate limiting middleware
//    - Recovery middleware
// 6. Support for HTTPS with TLS if configured
// 7. Graceful shutdown with timeout
// Use the Gin framework for routing
// Implement proper error handling with appropriate HTTP status codes
// Set up CORS if needed

package server

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nessi-dev/nessi-dev/internal/config"
	"github.com/nessi-dev/nessi-dev/internal/extensions"
	"github.com/nessi-dev/nessi-dev/internal/quality/profile"
	"github.com/nessi-dev/nessi-dev/internal/security"
	"github.com/nessi-dev/nessi-dev/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Server represents the API server
type Server struct {
	config      *config.Config
	server      *http.Server
	router      *gin.Engine
	extManager  *extensions.Manager
	secManager  *security.SecurityManager
	logger      *zap.Logger
	metrics     *Metrics
	startTime   time.Time
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
		Enabled bool
		CertFile string
		KeyFile  string
	}
}

// NewServer creates a new API server
func NewServer(config *Config) *Server {
	if config == nil {
		config = &Config{
			Host:            "localhost",
			Port:            8080,
			ReadTimeout:     15 * time.Second,
			WriteTimeout:    15 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		}
	}

	// Load security configuration
	secConfig, err := security.LoadConfigYAML("config/security.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to load security config: %v", err))
	}

	// Initialize security manager
	secManager := security.NewSecurityManager(secConfig.JWT.Secret, secConfig.JWT.Expiration)
	
	// Apply security configuration
	if err := secManager.ApplyConfig(secConfig); err != nil {
		panic(fmt.Sprintf("failed to apply security config: %v", err))
	}

	server := &Server{
		router:      gin.New(),
		config:      config,
		secManager:  secManager,
		startTime:   time.Now(),
	}

	server.setupMiddleware()
	server.setupRoutes()
	return server
}

// setupMiddleware configures the server middleware
func (s *Server) setupMiddleware() {
	// Use security middleware
	s.router.Use(s.secManager.Middleware)

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
	path := c.Param("path")
	profiler, err := profile.NewProfiler(path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create profiler: %v", err)})
		return
	}
	defer profiler.Close()

	qualityProfile, err := profiler.GenerateProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate profile: %v", err)})
		return
	}

	c.JSON(http.StatusOK, qualityProfile)
}

// handleQualityCheck handles quality check requests
func (s *Server) handleQualityCheck(c *gin.Context) {
	path := c.Param("path")
	profiler, err := profile.NewProfiler(path, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create profiler: %v", err)})
		return
	}
	defer profiler.Close()

	quality, err := profiler.GetDataQuality()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get quality metrics: %v", err)})
		return
	}

	c.JSON(http.StatusOK, quality)
}

// handleGetStats handles table statistics requests
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
	name := c.Param("name")
	if err := s.extManager.InstallExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to install extension: %v", err)})
		return
	}
	c.Status(http.StatusCreated)
}

// handleUninstallExtension handles extension uninstallation requests
func (s *Server) handleUninstallExtension(c *gin.Context) {
	name := c.Param("name")
	if err := s.extManager.UninstallExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to uninstall extension: %v", err)})
		return
	}
	c.Status(http.StatusNoContent)
}

// handleEnableExtension handles extension enabling requests
func (s *Server) handleEnableExtension(c *gin.Context) {
	name := c.Param("name")
	if err := s.extManager.EnableExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to enable extension: %v", err)})
		return
	}
	c.Status(http.StatusOK)
}

// handleDisableExtension handles extension disabling requests
func (s *Server) handleDisableExtension(c *gin.Context) {
	name := c.Param("name")
	if err := s.extManager.DisableExtension(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to disable extension: %v", err)})
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