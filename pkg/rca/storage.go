package rca

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// ResultStorage handles storage and retrieval of RCA results
type ResultStorage struct {
	storagePath string
	mutex       sync.RWMutex
	cache       map[string]*RCAResult
}

// NewResultStorage creates a new ResultStorage instance
func NewResultStorage(storagePath string) (*ResultStorage, error) {
	// Create storage directory if it doesn't exist
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create RCA storage directory: %w", err)
	}
	
	storage := &ResultStorage{
		storagePath: storagePath,
		cache:       make(map[string]*RCAResult),
	}
	
	// Load existing results into cache
	if err := storage.loadCache(); err != nil {
		return nil, fmt.Errorf("failed to load RCA cache: %w", err)
	}
	
	return storage, nil
}

// loadCache loads all existing RCA results into the cache
func (s *ResultStorage) loadCache() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Clear existing cache
	s.cache = make(map[string]*RCAResult)
	
	// Read all JSON files in the storage directory
	files, err := filepath.Glob(filepath.Join(s.storagePath, "*.json"))
	if err != nil {
		return fmt.Errorf("failed to read RCA storage directory: %w", err)
	}
	
	// Load each file into the cache
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			// Log error but continue
			fmt.Printf("Failed to read RCA result file %s: %v\n", file, err)
			continue
		}
		
		var result RCAResult
		if err := json.Unmarshal(data, &result); err != nil {
			// Log error but continue
			fmt.Printf("Failed to parse RCA result file %s: %v\n", file, err)
			continue
		}
		
		// Add to cache
		s.cache[result.AnomalyID] = &result
	}
	
	return nil
}

// SaveResult saves an RCA result to storage
func (s *ResultStorage) SaveResult(result *RCAResult) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Add to cache
	s.cache[result.AnomalyID] = result
	
	// Save to file
	filename := filepath.Join(s.storagePath, fmt.Sprintf("%s.json", result.AnomalyID))
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal RCA result: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write RCA result to file: %w", err)
	}
	
	return nil
}

// GetResult retrieves an RCA result by anomaly ID
func (s *ResultStorage) GetResult(anomalyID string) (*RCAResult, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	result, ok := s.cache[anomalyID]
	if !ok {
		return nil, nil // Not found
	}
	
	return result, nil
}

// GetRecentResults retrieves recent RCA results, sorted by analysis time (newest first)
func (s *ResultStorage) GetRecentResults(limit int) ([]*RCAResult, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// Convert cache to slice
	results := make([]*RCAResult, 0, len(s.cache))
	for _, result := range s.cache {
		results = append(results, result)
	}
	
	// Sort by analysis time (newest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].AnalysisTime.After(results[j].AnalysisTime)
	})
	
	// Apply limit
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}
	
	return results, nil
}

// GetInsights generates insights from all RCA results
func (s *ResultStorage) GetInsights() (*Insights, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// Get all results
	results := make([]*RCAResult, 0, len(s.cache))
	for _, result := range s.cache {
		results = append(results, result)
	}
	
	// Generate insights
	return GenerateInsights(results), nil
}

// DeleteResult deletes an RCA result by anomaly ID
func (s *ResultStorage) DeleteResult(anomalyID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	// Remove from cache
	delete(s.cache, anomalyID)
	
	// Remove file
	filename := filepath.Join(s.storagePath, fmt.Sprintf("%s.json", anomalyID))
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete RCA result file: %w", err)
	}
	
	return nil
}

// GetResultsByTimeRange retrieves RCA results within a time range
func (s *ResultStorage) GetResultsByTimeRange(start, end time.Time) ([]*RCAResult, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// Filter results by time range
	results := make([]*RCAResult, 0)
	for _, result := range s.cache {
		if (result.AnalysisTime.Equal(start) || result.AnalysisTime.After(start)) &&
			(result.AnalysisTime.Equal(end) || result.AnalysisTime.Before(end)) {
			results = append(results, result)
		}
	}
	
	// Sort by analysis time (newest first)
	sort.Slice(results, func(i, j int) bool {
		return results[i].AnalysisTime.After(results[j].AnalysisTime)
	})
	
	return results, nil
}
