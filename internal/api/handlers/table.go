package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/nessi-dev/nessi/internal/delta"
	// "github.com/nessi-dev/nessi/internal/quality/profile" // Added for profile.Profile - Will be used later
	rules "github.com/nessi-dev/nessi/internal/quality/rules" // Changed for rules.Rule
)

// Table represents a Delta table
// swagger:model table
// @Id Table
type Table struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Version int64    `json:"version"`
	Columns []string `json:"columns"`
}

// ListTablesHandler handles GET /api/v1/tables
// swagger:route GET /api/v1/tables listTables
// List all Delta tables
// responses:
//   200: []Table
//   400: error
//   500: error
func ListTablesHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "/data/delta"
	}

	tables, err := delta.ListTables(context.Background(), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// GetTableHandler handles GET /api/v1/tables/{tableName}
// swagger:route GET /api/v1/tables/{tableName} getTable
// Get table details
// responses:
//   200: Table
//   400: error
//   404: error
//   500: error
func GetTableHandler(c *gin.Context) {
	tableName := c.Param("tableName")
	versionStr := c.Query("version")
	version := int64(-1)

	if versionStr != "" {
		var err error
		version, err = strconv.ParseInt(versionStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version number"})
			return
		}
	}

	table, err := delta.GetTable(context.Background(), tableName, version)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, table)
}

// GetTableProfileHandler handles GET /api/v1/tables/{tableName}/profile
// swagger:route GET /api/v1/tables/{tableName}/profile getTableProfile
// Get table profile
// responses:
//   200: profile.Profile
//   400: error
//   404: error
//   500: error
func GetTableProfileHandler(c *gin.Context) {
	tableName := c.Param("tableName")
	sampleStr := c.Query("sample")
	sample := 100.0 // default to 100% if not specified

	if sampleStr != "" {
		sample, err := strconv.ParseFloat(sampleStr, 64)
		if err != nil || sample < 0 || sample > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sample percentage must be between 0 and 100"})
			return
		}
	}

	profile, err := delta.GetTableProfile(context.Background(), tableName, sample)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, profile)
}

// CheckTableQualityHandler handles POST /api/v1/tables/{tableName}/check
// swagger:route POST /api/v1/tables/{tableName}/check checkTableQuality
// Run quality checks on table
// responses:
//   200: quality.QualityReport
//   400: error
//   404: error
//   500: error
func CheckTableQualityHandler(c *gin.Context) {
	tableName := c.Param("tableName")
	var requestRules []rules.Rule

	if err := c.ShouldBindJSON(&requestRules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rules JSON"})
		return
	}

	// Validate all rules
	for _, rule := range requestRules {
		if err := rule.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	report, err := delta.CheckTableQuality(context.Background(), tableName, requestRules)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetTableVersionsHandler handles GET /api/v1/tables/{tableName}/versions
// swagger:route GET /api/v1/tables/{tableName}/versions getTableVersions
// List table versions
// responses:
//   200: []int64
//   404: error
//   500: error
func GetTableVersionsHandler(c *gin.Context) {
	tableName := c.Param("tableName")

	versions, err := delta.GetTableVersions(context.Background(), tableName)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, versions)
}