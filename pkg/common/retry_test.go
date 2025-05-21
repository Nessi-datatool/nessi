package common

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsRetryableError(t *testing.T) {
	// Test nil error
	assert.False(t, IsRetryableError(nil))

	// Test retryable NessiError
	retryableNessiErrors := []ErrorCode{
		ErrConnectionFailed,
		ErrRateLimitExceeded,
		ErrServerError,
		ErrTimeout,
	}

	for _, code := range retryableNessiErrors {
		err := NewError(code, "test error")
		assert.True(t, IsRetryableError(err), "Error with code %s should be retryable", code)
	}

	// Test non-retryable NessiError
	nonRetryableNessiErrors := []ErrorCode{
		ErrInvalidPath,
		ErrFileNotFound,
		ErrNotDeltaTable,
		ErrAuthFailed,
		ErrResourceNotFound,
	}

	for _, code := range nonRetryableNessiErrors {
		err := NewError(code, "test error")
		assert.False(t, IsRetryableError(err), "Error with code %s should not be retryable", code)
	}

	// Test context errors
	assert.True(t, IsRetryableError(context.DeadlineExceeded))
	assert.True(t, IsRetryableError(context.Canceled))

	// Test regular errors
	assert.False(t, IsRetryableError(errors.New("generic error")))
}

func TestWithRetry(t *testing.T) {
	// Test successful function
	callCount := 0
	successFn := func() error {
		callCount++
		return nil
	}

	config := RetryConfig{
		MaxRetries:     3,
		InitialBackoff: 1 * time.Millisecond,
		MaxBackoff:     10 * time.Millisecond,
		BackoffFactor:  2.0,
		Jitter:         0.1,
	}

	err := WithRetry(successFn, config)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount, "Function should be called exactly once")

	// Test function that fails with non-retryable error
	callCount = 0
	nonRetryableFn := func() error {
		callCount++
		return NewError(ErrInvalidPath, "test error")
	}

	err = WithRetry(nonRetryableFn, config)
	assert.Error(t, err)
	assert.Equal(t, 1, callCount, "Function should be called exactly once")

	// Test function that fails with retryable error but eventually succeeds
	callCount = 0
	eventualSuccessFn := func() error {
		callCount++
		if callCount < 3 {
			return NewError(ErrConnectionFailed, "test error")
		}
		return nil
	}

	err = WithRetry(eventualSuccessFn, config)
	assert.NoError(t, err)
	assert.Equal(t, 3, callCount, "Function should be called exactly 3 times")

	// Test function that always fails with retryable error
	callCount = 0
	alwaysFailsFn := func() error {
		callCount++
		return NewError(ErrConnectionFailed, "test error")
	}

	err = WithRetry(alwaysFailsFn, config)
	assert.Error(t, err)
	assert.Equal(t, config.MaxRetries+1, callCount, "Function should be called MaxRetries+1 times")
	assert.Contains(t, err.Error(), fmt.Sprintf("failed after %d retries", config.MaxRetries))
}

func TestWithDefaultRetry(t *testing.T) {
	// Test successful function
	callCount := 0
	successFn := func() error {
		callCount++
		return nil
	}

	err := WithDefaultRetry(successFn)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount, "Function should be called exactly once")

	// Test function that fails with retryable error but eventually succeeds
	callCount = 0
	eventualSuccessFn := func() error {
		callCount++
		if callCount < 2 {
			return NewError(ErrConnectionFailed, "test error")
		}
		return nil
	}

	err = WithDefaultRetry(eventualSuccessFn)
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount, "Function should be called exactly 2 times")
}
