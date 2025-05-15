package rca

import (
	"fmt"
	"time"
)

// Client implements the RCA client interface for the dashboard
type Client struct {
	analyzer *Analyzer
	storage  *ResultStorage
}

// NewClient creates a new RCA client
func NewClient(analyzer *Analyzer, storagePath string) (*Client, error) {
	storage, err := NewResultStorage(storagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create RCA storage: %w", err)
	}

	return &Client{
		analyzer: analyzer,
		storage:  storage,
	}, nil
}

// AnalyzeAnomaly performs root cause analysis on an anomaly
func (c *Client) AnalyzeAnomaly(anomalyID string, config *Config) (*RCAResult, error) {
	// Check if analysis already exists
	existing, err := c.storage.GetResult(anomalyID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing analysis: %w", err)
	}

	// If analysis exists and is recent (less than 1 hour old), return it
	if existing != nil && time.Since(existing.AnalysisTime) < time.Hour {
		return existing, nil
	}

	// Perform analysis
	result, err := c.analyzer.AnalyzeAnomaly(anomalyID)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze anomaly: %w", err)
	}

	// Save result
	if err := c.storage.SaveResult(result); err != nil {
		return nil, fmt.Errorf("failed to save analysis result: %w", err)
	}

	return result, nil
}

// GetRecentAnalyses returns recent RCA results
func (c *Client) GetRecentAnalyses(limit int) ([]*RCAResult, error) {
	return c.storage.GetRecentResults(limit)
}

// GetAnalysisResult returns a specific RCA result
func (c *Client) GetAnalysisResult(anomalyID string) (*RCAResult, error) {
	return c.storage.GetResult(anomalyID)
}

// GetInsights returns aggregated insights from RCA results
func (c *Client) GetInsights() (*Insights, error) {
	return c.storage.GetInsights()
}

// DeleteAnalysisResult deletes a specific RCA result
func (c *Client) DeleteAnalysisResult(anomalyID string) error {
	return c.storage.DeleteResult(anomalyID)
}

// GetAnalysesByTimeRange returns RCA results within a time range
func (c *Client) GetAnalysesByTimeRange(start, end time.Time) ([]*RCAResult, error) {
	return c.storage.GetResultsByTimeRange(start, end)
}
