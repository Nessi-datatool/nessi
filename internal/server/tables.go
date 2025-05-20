package server

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nessi-dev/nessi/internal/quality/profile"
	"github.com/nessi-dev/nessi/pkg"
	"go.uber.org/zap"
)

// TableHandler handles Delta table operations
type TableHandler struct {
	logger *zap.Logger
}

// NewTableHandler creates a new table handler
func NewTableHandler(logger *zap.Logger) *TableHandler {
	return &TableHandler{
		logger: logger,
	}
}

// RegisterRoutes registers table-related routes
func (h *TableHandler) RegisterRoutes(router *gin.RouterGroup) {
	tables := router.Group("/tables")
	{
		tables.POST("/:path/profile", h.handleGenerateProfile)
		tables.POST("/:path/quality", h.handleQualityCheck)
		tables.GET("/:path/stats", h.handleGetStats)
		tables.GET("/:path/schema", h.handleGetSchema)
		tables.GET("/:path/partitions", h.handleGetPartitions)
		tables.GET("/:path/partitions/:partition", h.handleReadPartition)
	}
}

// handleGenerateProfile handles profile generation requests
func (h *TableHandler) handleGenerateProfile(c *gin.Context) {
	path := c.Param("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid table path: %v", err)})
		return
	}

	// Create profiler with default config
	profiler, err := profile.NewProfiler(absPath, nil)
	if err != nil {
		h.logger.Error("failed to create profiler",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profiler"})
		return
	}
	defer profiler.Close()

	// Generate profile
	qualityProfile, err := profiler.ProfileTable()
	if err != nil {
		h.logger.Error("failed to generate profile",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate profile"})
		return
	}

	c.JSON(http.StatusOK, qualityProfile)
}

// handleQualityCheck handles quality check requests
func (h *TableHandler) handleQualityCheck(c *gin.Context) {
	path := c.Param("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid table path: %v", err)})
		return
	}

	// Create profiler with default config
	profiler, err := profile.NewProfiler(absPath, nil)
	if err != nil {
		h.logger.Error("failed to create profiler",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profiler"})
		return
	}
	defer profiler.Close()

	// Get quality metrics
	quality, err := profiler.GetDataQuality()
	if err != nil {
		h.logger.Error("failed to get quality metrics",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get quality metrics"})
		return
	}

	c.JSON(http.StatusOK, quality)
}

// handleGetStats handles table statistics requests
func (h *TableHandler) handleGetStats(c *gin.Context) {
	path := c.Param("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid table path: %v", err)})
		return
	}

	// Create reader
	reader, err := pkg.NewReader(absPath)
	if err != nil {
		h.logger.Error("failed to create reader",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reader"})
		return
	}
	defer reader.Close()

	// Get table statistics
	stats := reader.GetStats()
	c.JSON(http.StatusOK, stats)
}

// handleGetSchema handles schema retrieval requests
func (h *TableHandler) handleGetSchema(c *gin.Context) {
	path := c.Param("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid table path: %v", err)})
		return
	}

	// Create reader
	reader, err := pkg.NewReader(absPath)
	if err != nil {
		h.logger.Error("failed to create reader",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reader"})
		return
	}
	defer reader.Close()

	// Get table schema
	schema := reader.GetSchema()
	c.JSON(http.StatusOK, schema)
}

// handleGetPartitions handles partition listing requests
func (h *TableHandler) handleGetPartitions(c *gin.Context) {
	path := c.Param("path")
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid table path: %v", err)})
		return
	}

	// Create reader
	reader, err := pkg.NewReader(absPath)
	if err != nil {
		h.logger.Error("failed to create reader",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reader"})
		return
	}
	defer reader.Close()

	// Get table partitions
	partitions, err := reader.ListPartitions()
	if err != nil {
		h.logger.Error("failed to get partitions",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get partitions"})
		return
	}

	c.JSON(http.StatusOK, partitions)
}

// handleReadPartition handles partition reading requests
func (h *TableHandler) handleReadPartition(c *gin.Context) {
	path := c.Param("path")
	partition := c.Param("partition")
	absPath, err := filepath.Abs(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid table path: %v", err)})
		return
	}

	// Create reader
	reader, err := pkg.NewReader(absPath)
	if err != nil {
		h.logger.Error("failed to create reader",
			zap.String("path", absPath),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create reader"})
		return
	}
	defer reader.Close()

	// Read partition data
	data, err := reader.ReadPartition(partition)
	if err != nil {
		h.logger.Error("failed to read partition",
			zap.String("path", absPath),
			zap.String("partition", partition),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read partition"})
		return
	}

	c.JSON(http.StatusOK, data)
}
