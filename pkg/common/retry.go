package common

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RetryableFunc is a function that can be retried
type RetryableFunc func() error

// RetryConfig configures the retry behavior
type RetryConfig struct {
	MaxRetries     int           // Maximum number of retries
	InitialBackoff time.Duration // Initial backoff duration
	MaxBackoff     time.Duration // Maximum backoff duration
	BackoffFactor  float64       // Factor to increase backoff by after each retry
	Jitter         float64       // Jitter factor (0-1) to randomize backoff
}

// DefaultRetryConfig provides sensible defaults for retry behavior
var DefaultRetryConfig = RetryConfig{
	MaxRetries:     3,
	InitialBackoff: 1 * time.Second,
	MaxBackoff:     30 * time.Second,
	BackoffFactor:  2.0,
	Jitter:         0.2,
}

// IsRetryableError determines if an error should be retried
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Check if it's a NessiError
	var nessiErr *NessiError
	if errors.As(err, &nessiErr) {
		// Retry connection errors, rate limit errors, and server errors
		retryableCodes := []ErrorCode{
			ErrConnectionFailed,
			ErrRateLimitExceeded,
			ErrServerError,
			ErrTimeout,
		}
		
		for _, code := range retryableCodes {
			if nessiErr.Code == code {
				return true
			}
		}
		
		return false
	}
	
	// For non-NessiError, retry on common transient errors
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}
	
	// Check error string for common network errors
	errStr := err.Error()
	networkErrors := []string{
		"connection refused",
		"connection reset",
		"connection timed out",
		"no route to host",
		"network is unreachable",
		"i/o timeout",
		"EOF",
		"broken pipe",
		"too many open files",
	}
	
	for _, netErr := range networkErrors {
		if errors.Is(err, fmt.Errorf(netErr)) || (errStr != "" && contains(errStr, netErr)) {
			return true
		}
	}
	
	return false
}

// WithRetry retries a function with exponential backoff
func WithRetry(fn RetryableFunc, config RetryConfig) error {
	var err error
	backoff := config.InitialBackoff
	
	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Execute the function
		err = fn()
		
		// If no error or not retryable, return immediately
		if err == nil || !IsRetryableError(err) {
			return err
		}
		
		// If this was the last attempt, return the error
		if attempt == config.MaxRetries {
			return fmt.Errorf("failed after %d retries: %w", config.MaxRetries, err)
		}
		
		// Apply jitter to backoff
		jitterRange := backoff.Seconds() * config.Jitter
		jitterSeconds := (rand.Float64() * jitterRange * 2) - jitterRange
		jitteredBackoff := backoff + time.Duration(jitterSeconds*float64(time.Second))
		
		// Sleep before the next attempt
		time.Sleep(jitteredBackoff)
		
		// Increase backoff for the next attempt
		backoff = time.Duration(float64(backoff) * config.BackoffFactor)
		if backoff > config.MaxBackoff {
			backoff = config.MaxBackoff
		}
	}
	
	return err
}

// WithDefaultRetry retries a function with default retry settings
func WithDefaultRetry(fn RetryableFunc) error {
	return WithRetry(fn, DefaultRetryConfig)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && strings.Contains(s, substr)
}
