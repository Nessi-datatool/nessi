package dashboard

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/meisi/nessi-dev/pkg/freshness"
	"github.com/meisi/nessi-dev/pkg/logger"
)

// FreshnessHandler handles the freshness dashboard route
func (d *Dashboard) FreshnessHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "freshness.html", gin.H{
		"Title": "Data Freshness Dashboard",
	})
}

// FreshnessStatusAPI handles the API endpoint for freshness status
func (d *Dashboard) FreshnessStatusAPI(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	tableName := c.Query("table")

	// Get freshness manager
	freshnessManager, err := d.getFreshnessManager()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get freshness manager")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness manager"})
		return
	}

	var statuses []freshness.TableFreshnessStatus
	if tableName != "" {
		// Get status for specific table
		status, err := freshnessManager.GetTableStatus(tableName)
		if err != nil {
			log.Error().Err(err).Str("table", tableName).Msg("Failed to get table freshness status")
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found or error retrieving status"})
			return
		}
		statuses = []freshness.TableFreshnessStatus{status}
	} else {
		// Get status for all tables
		statuses, err = freshnessManager.GetAllTableStatuses()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get all table freshness statuses")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness statuses"})
			return
		}
	}

	c.JSON(http.StatusOK, statuses)
}

// FreshnessSLAAPI handles the API endpoint for SLA configurations
func (d *Dashboard) FreshnessSLAAPI(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	tableName := c.Query("table")

	// Get freshness manager
	freshnessManager, err := d.getFreshnessManager()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get freshness manager")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness manager"})
		return
	}

	// Handle different HTTP methods
	switch c.Request.Method {
	case http.MethodGet:
		// GET: Retrieve SLA configurations
		if tableName != "" {
			// Get SLA for specific table
			sla, err := freshnessManager.GetSLAConfig(tableName)
			if err != nil {
				log.Error().Err(err).Str("table", tableName).Msg("Failed to get SLA configuration")
				c.JSON(http.StatusNotFound, gin.H{"error": "SLA configuration not found"})
				return
			}
			c.JSON(http.StatusOK, []freshness.SLAConfig{sla})
		} else {
			// Get all SLA configurations
			slas, err := freshnessManager.GetAllSLAConfigs()
			if err != nil {
				log.Error().Err(err).Msg("Failed to get all SLA configurations")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get SLA configurations"})
				return
			}
			c.JSON(http.StatusOK, slas)
		}

	case http.MethodPost:
		// POST: Create or update SLA configuration
		var slaConfig freshness.SLAConfig
		if err := c.ShouldBindJSON(&slaConfig); err != nil {
			log.Error().Err(err).Msg("Invalid SLA configuration format")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SLA configuration format"})
			return
		}

		// Parse expected frequency if it's a string
		if frequencyStr, ok := c.PostForm("expected_frequency").(string); ok {
			// Handle predefined frequencies
			switch frequencyStr {
			case "hourly":
				slaConfig.ExpectedFrequency = time.Hour
			case "daily":
				slaConfig.ExpectedFrequency = 24 * time.Hour
			case "weekly":
				slaConfig.ExpectedFrequency = 7 * 24 * time.Hour
			case "monthly":
				slaConfig.ExpectedFrequency = 30 * 24 * time.Hour
			default:
				// Try to parse custom duration
				duration, err := time.ParseDuration(frequencyStr)
				if err != nil {
					log.Error().Err(err).Str("frequency", frequencyStr).Msg("Invalid frequency format")
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid frequency format"})
					return
				}
				slaConfig.ExpectedFrequency = duration
			}
		}

		// Save SLA configuration
		if err := freshnessManager.SetSLAConfig(slaConfig); err != nil {
			log.Error().Err(err).Interface("config", slaConfig).Msg("Failed to save SLA configuration")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save SLA configuration"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})

	case http.MethodDelete:
		// DELETE: Remove SLA configuration
		if tableName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Table name is required"})
			return
		}

		if err := freshnessManager.DeleteSLAConfig(tableName); err != nil {
			log.Error().Err(err).Str("table", tableName).Msg("Failed to delete SLA configuration")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SLA configuration"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})

	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

// FreshnessTrendsAPI handles the API endpoint for freshness trends
func (d *Dashboard) FreshnessTrendsAPI(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	tableName := c.Query("table")

	// Get freshness manager
	freshnessManager, err := d.getFreshnessManager()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get freshness manager")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness manager"})
		return
	}

	// Get trends data
	var trendsData *freshness.FreshnessTrends
	if tableName != "" && tableName != "all" {
		// Get trends for specific table
		trendsData, err = freshnessManager.GetTableTrends(tableName)
	} else {
		// Get trends for all tables
		trendsData, err = freshnessManager.GetAllTablesTrends()
	}

	if err != nil {
		log.Error().Err(err).Str("table", tableName).Msg("Failed to get freshness trends")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness trends"})
		return
	}

	c.JSON(http.StatusOK, trendsData)
}

// FreshnessExportAPI handles the API endpoint for exporting freshness data
func (d *Dashboard) FreshnessExportAPI(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	format := c.Query("format")
	if format == "" {
		format = "csv" // Default format
	}

	// Get freshness manager
	freshnessManager, err := d.getFreshnessManager()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get freshness manager")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness manager"})
		return
	}

	// Get all table statuses
	statuses, err := freshnessManager.GetAllTableStatuses()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get all table freshness statuses")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get freshness statuses"})
		return
	}

	// Handle different export formats
	switch format {
	case "json":
		// Export as JSON
		c.Header("Content-Disposition", "attachment; filename=freshness_data.json")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, statuses)

	case "csv":
		// Export as CSV
		c.Header("Content-Disposition", "attachment; filename=freshness_data.csv")
		c.Header("Content-Type", "text/csv")

		// Create CSV writer
		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		// Write header
		header := []string{
			"Table Name",
			"Table Path",
			"Last Update Time",
			"Time Since Update (seconds)",
			"Expected Frequency (seconds)",
			"Next Expected Update",
			"Status",
			"Warning Threshold (%)",
			"Critical Threshold (%)",
			"Enabled",
		}
		if err := writer.Write(header); err != nil {
			log.Error().Err(err).Msg("Failed to write CSV header")
			return
		}

		// Write data rows
		for _, status := range statuses {
			row := []string{
				status.TableName,
				status.TablePath,
				status.LastUpdateTime.Format(time.RFC3339),
				strconv.FormatInt(int64(status.TimeSinceUpdate.Seconds()), 10),
				strconv.FormatInt(int64(status.ExpectedFrequency.Seconds()), 10),
				status.NextExpectedUpdate.Format(time.RFC3339),
				status.Status,
				strconv.Itoa(status.SLAConfig.WarningThreshold),
				strconv.Itoa(status.SLAConfig.CriticalThreshold),
				strconv.FormatBool(status.SLAConfig.Enabled),
			}
			if err := writer.Write(row); err != nil {
				log.Error().Err(err).Msg("Failed to write CSV row")
				return
			}
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported export format"})
		return
	}
}

// getFreshnessManager returns the freshness manager from the context
func (d *Dashboard) getFreshnessManager() (*freshness.Manager, error) {
	// Get freshness manager from service provider
	freshnessManager, err := d.serviceProvider.GetFreshnessManager()
	if err != nil {
		return nil, err
	}

	return freshnessManager, nil
}

// RegisterFreshnessRoutes registers the freshness dashboard routes
func (d *Dashboard) RegisterFreshnessRoutes(router *gin.Engine) {
	// Dashboard route
	router.GET("/freshness", d.authRequired(), d.FreshnessHandler)

	// API routes
	api := router.Group("/api/freshness", d.authRequired())
	{
		api.GET("/status", d.FreshnessStatusAPI)
		api.GET("/sla", d.FreshnessSLAAPI)
		api.POST("/sla", d.FreshnessSLAAPI)
		api.DELETE("/sla", d.FreshnessSLAAPI)
		api.GET("/trends", d.FreshnessTrendsAPI)
		api.GET("/export", d.FreshnessExportAPI)
	}
}
