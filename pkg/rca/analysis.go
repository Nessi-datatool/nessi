package rca

import (
	"time"
)

// Analysis represents a root cause analysis job
type Analysis struct {
	ID           string    `json:"id"`
	AnomalyID    string    `json:"anomaly_id"`
	Status       string    `json:"status"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time,omitempty"`
	ResultID     string    `json:"result_id,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	Config       *Config   `json:"config,omitempty"`
}

// AnalysisStatus constants
const (
	AnalysisStatusPending   = "pending"
	AnalysisStatusRunning   = "running"
	AnalysisStatusCompleted = "completed"
	AnalysisStatusFailed    = "failed"
)

// NewAnalysis creates a new analysis
func NewAnalysis(anomalyID string, config *Config) *Analysis {
	return &Analysis{
		ID:        GenerateID(),
		AnomalyID: anomalyID,
		Status:    AnalysisStatusPending,
		StartTime: time.Now(),
		Config:    config,
	}
}

// SetRunning marks the analysis as running
func (a *Analysis) SetRunning() {
	a.Status = AnalysisStatusRunning
}

// SetCompleted marks the analysis as completed
func (a *Analysis) SetCompleted(resultID string) {
	a.Status = AnalysisStatusCompleted
	a.ResultID = resultID
	a.EndTime = time.Now()
}

// SetFailed marks the analysis as failed
func (a *Analysis) SetFailed(errorMessage string) {
	a.Status = AnalysisStatusFailed
	a.ErrorMessage = errorMessage
	a.EndTime = time.Now()
}

// GenerateID generates a unique ID for an analysis
func GenerateID() string {
	// Simple implementation for now
	return time.Now().Format("20060102150405") + "-" + RandomString(8)
}

// RandomString generates a random string of the specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
	}
	return string(result)
}
