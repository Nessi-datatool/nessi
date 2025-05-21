package common

import (
	"fmt"
)

// NessiError represents a structured error in Nessi
type NessiError struct {
	Code        ErrorCode // Error code
	Message     string    // Error message
	Details     string    // Additional details about the error
	Suggestion  string    // Suggestion for fixing the error
	WrappedErr  error     // Original error that was wrapped
}

// Error implements the error interface
func (e *NessiError) Error() string {
	baseMsg := fmt.Sprintf("[%s] %s: %s", e.Code, GetErrorDescription(e.Code), e.Message)
	
	if e.Details != "" {
		baseMsg = fmt.Sprintf("%s (%s)", baseMsg, e.Details)
	}
	
	if e.Suggestion != "" {
		baseMsg = fmt.Sprintf("%s. Suggestion: %s", baseMsg, e.Suggestion)
	}
	
	return baseMsg
}

// Unwrap returns the wrapped error
func (e *NessiError) Unwrap() error {
	return e.WrappedErr
}

// NewError creates a new NessiError
func NewError(code ErrorCode, message string) *NessiError {
	return &NessiError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to the error
func (e *NessiError) WithDetails(details string) *NessiError {
	e.Details = details
	return e
}

// WithSuggestion adds a suggestion to the error
func (e *NessiError) WithSuggestion(suggestion string) *NessiError {
	e.Suggestion = suggestion
	return e
}

// WithWrappedError adds a wrapped error
func (e *NessiError) WithWrappedError(err error) *NessiError {
	e.WrappedErr = err
	return e
}

// NewInvalidPathError creates a new invalid path error
func NewInvalidPathError(path string) *NessiError {
	return NewError(ErrInvalidPath, fmt.Sprintf("Path '%s' does not exist or is not accessible", path)).
		WithSuggestion("Check that the path exists and you have permission to access it")
}

// NewNotDeltaTableError creates a new not a Delta table error
func NewNotDeltaTableError(path string) *NessiError {
	return NewError(ErrNotDeltaTable, fmt.Sprintf("'%s' is not a Delta Lake table", path)).
		WithSuggestion("Ensure the path points to a valid Delta Lake table with a _delta_log directory")
}

// NewAuthFailedError creates a new authentication failed error
func NewAuthFailedError(service string) *NessiError {
	return NewError(ErrAuthFailed, fmt.Sprintf("Authentication failed for %s", service)).
		WithSuggestion("Check your credentials and ensure they are correctly configured")
}

// NewConnectionFailedError creates a new connection failed error
func NewConnectionFailedError(service string, details string) *NessiError {
	return NewError(ErrConnectionFailed, fmt.Sprintf("Could not connect to %s", service)).
		WithDetails(details).
		WithSuggestion("Check your network connection and ensure the service is available")
}

// NewResourceNotFoundError creates a new resource not found error
func NewResourceNotFoundError(resourceType, resourceName string) *NessiError {
	return NewError(ErrResourceNotFound, fmt.Sprintf("%s '%s' does not exist", resourceType, resourceName)).
		WithSuggestion(fmt.Sprintf("Verify that the %s exists and you have permission to access it", resourceType))
}

// NewSchemaValidationError creates a new schema validation error
func NewSchemaValidationError(field, expectedType, actualType string) *NessiError {
	return NewError(ErrSchemaValidation, fmt.Sprintf("Field '%s' expected type %s but got %s", field, expectedType, actualType)).
		WithSuggestion("Ensure the data conforms to the expected schema")
}

// NewWorkspaceIDEmptyError creates a new workspace ID empty error
func NewWorkspaceIDEmptyError() *NessiError {
	return NewError(ErrWorkspaceIDEmpty, "Workspace ID is required for multi-workspace environments").
		WithSuggestion("Set the DATABRICKS_WORKSPACE_ID environment variable")
}

// NewRateLimitExceededError creates a new rate limit exceeded error
func NewRateLimitExceededError(service string) *NessiError {
	return NewError(ErrRateLimitExceeded, fmt.Sprintf("Too many requests to %s API", service)).
		WithSuggestion("Reduce the frequency of requests or implement exponential backoff")
}

// NewServerError creates a new server error
func NewServerError(service string, statusCode int) *NessiError {
	return NewError(ErrServerError, fmt.Sprintf("%s API returned an error (%d)", service, statusCode)).
		WithSuggestion("Check the server logs for more information")
}
