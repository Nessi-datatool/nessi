package common

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNessiError(t *testing.T) {
	// Test basic error creation
	err := NewError(ErrInvalidPath, "test path is invalid")
	assert.Equal(t, "[N101] Invalid path: test path is invalid", err.Error())

	// Test with details
	err = NewError(ErrInvalidPath, "test path is invalid").WithDetails("file does not exist")
	assert.Equal(t, "[N101] Invalid path: test path is invalid (file does not exist)", err.Error())

	// Test with suggestion
	err = NewError(ErrInvalidPath, "test path is invalid").WithSuggestion("check the path")
	assert.Equal(t, "[N101] Invalid path: test path is invalid. Suggestion: check the path", err.Error())

	// Test with wrapped error
	originalErr := fmt.Errorf("original error")
	err = NewError(ErrInvalidPath, "test path is invalid").WithWrappedError(originalErr)
	assert.Equal(t, "[N101] Invalid path: test path is invalid", err.Error())
	unwrappedErr := errors.Unwrap(err)
	assert.Equal(t, originalErr, unwrappedErr)

	// Test with all fields
	err = NewError(ErrInvalidPath, "test path is invalid").
		WithDetails("file does not exist").
		WithSuggestion("check the path").
		WithWrappedError(originalErr)
	assert.Equal(t, "[N101] Invalid path: test path is invalid (file does not exist). Suggestion: check the path", err.Error())
}

func TestHelperFunctions(t *testing.T) {
	// Test NewInvalidPathError
	err := NewInvalidPathError("/test/path")
	assert.Contains(t, err.Error(), "[N101] Invalid path")
	assert.Contains(t, err.Error(), "/test/path")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewNotDeltaTableError
	err = NewNotDeltaTableError("/test/path")
	assert.Contains(t, err.Error(), "[N201] Not a Delta table")
	assert.Contains(t, err.Error(), "/test/path")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewAuthFailedError
	err = NewAuthFailedError("Databricks")
	assert.Contains(t, err.Error(), "[N301] Authentication failed")
	assert.Contains(t, err.Error(), "Databricks")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewConnectionFailedError
	err = NewConnectionFailedError("Databricks", "timeout after 30s")
	assert.Contains(t, err.Error(), "[N302] Connection failed")
	assert.Contains(t, err.Error(), "Databricks")
	assert.Contains(t, err.Error(), "timeout after 30s")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewResourceNotFoundError
	err = NewResourceNotFoundError("Catalog", "main")
	assert.Contains(t, err.Error(), "[N303] Resource not found")
	assert.Contains(t, err.Error(), "Catalog 'main'")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewSchemaValidationError
	err = NewSchemaValidationError("age", "INT", "STRING")
	assert.Contains(t, err.Error(), "[N501] Schema validation failed")
	assert.Contains(t, err.Error(), "Field 'age'")
	assert.Contains(t, err.Error(), "expected type INT but got STRING")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewWorkspaceIDEmptyError
	err = NewWorkspaceIDEmptyError()
	assert.Contains(t, err.Error(), "[N306] Workspace ID is empty")
	assert.Contains(t, err.Error(), "Workspace ID is required")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewRateLimitExceededError
	err = NewRateLimitExceededError("Databricks")
	assert.Contains(t, err.Error(), "[N304] Rate limit exceeded")
	assert.Contains(t, err.Error(), "Databricks API")
	assert.Contains(t, err.Error(), "Suggestion")

	// Test NewServerError
	err = NewServerError("Databricks", 500)
	assert.Contains(t, err.Error(), "[N305] Server error")
	assert.Contains(t, err.Error(), "Databricks API returned an error (500)")
	assert.Contains(t, err.Error(), "Suggestion")
}
