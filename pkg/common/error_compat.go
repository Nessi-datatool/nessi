package common

import (
	"github.com/nessi-dev/nessi/pkg/errorcode"
	"github.com/nessi-dev/nessi/pkg/errors"
)

// ErrorCode is a compatibility type that wraps errorcode.ErrorCode
type ErrorCode = errorcode.ErrorCode

// NessiError is a compatibility type that wraps errors.NessiError
type NessiError = errors.NessiError

// Error code constants for backward compatibility
const (
	// Unknown error
	ErrUnknown = "N000"
	// General errors (N1XX)
	ErrInvalidPath       = errorcode.ErrInvalidPath
	ErrFileNotFound      = errorcode.ErrFileNotFound
	ErrPermissionDenied  = errorcode.ErrPermissionDenied
	ErrInvalidArgument   = errorcode.ErrInvalidArgument
	ErrInvalidOperation  = errorcode.ErrInvalidOperation
	ErrTimeout           = errorcode.ErrTimeout
	ErrInternal          = errorcode.ErrInternal
	ErrInternalError     = errorcode.ErrInternal
	ErrInvalidDeltaTable = errorcode.ErrNotDeltaTable

	// Delta Lake errors (N2XX)
	ErrNotDeltaTable     = errorcode.ErrNotDeltaTable
	ErrCorruptedDeltaLog = errorcode.ErrCorruptedDeltaLog
	ErrInvalidVersion    = errorcode.ErrInvalidVersion
	ErrEmptyTable        = errorcode.ErrEmptyTable
	ErrDeltaLogNotFound  = errorcode.ErrDeltaLogNotFound
	ErrParquetReadFailed = errorcode.ErrParquetReadFailed

	// Databricks errors (N3XX)
	ErrAuthFailed        = errorcode.ErrAuthFailed
	ErrConnectionFailed  = errorcode.ErrConnectionFailed
	ErrResourceNotFound  = errorcode.ErrResourceNotFound
	ErrRateLimitExceeded = errorcode.ErrRateLimitExceeded
	ErrServerError       = errorcode.ErrServerError
	ErrWorkspaceIDEmpty  = errorcode.ErrWorkspaceIDEmpty
	ErrInvalidResponse   = errorcode.ErrInvalidResponse

	// Cloud provider errors (N4XX)
	ErrAWSAuthFailed     = errorcode.ErrAWSAuthFailed
	ErrGCPAuthFailed     = errorcode.ErrGCPAuthFailed
	ErrAzureAuthFailed   = errorcode.ErrAzureAuthFailed
	ErrS3BucketNotFound  = errorcode.ErrS3BucketNotFound
	ErrGCSBucketNotFound = errorcode.ErrGCSBucketNotFound
	ErrADLSNotFound      = errorcode.ErrADLSNotFound

	// Schema and data errors (N5XX)
	ErrSchemaValidation   = errorcode.ErrSchemaValidation
	ErrDataTypeConversion = errorcode.ErrDataTypeConversion
	ErrInvalidSchema      = errorcode.ErrInvalidSchema
	ErrSchemaMismatch     = errorcode.ErrSchemaMismatch
	ErrDataValidation     = errorcode.ErrDataValidation

	// Configuration errors (N6XX)
	ErrConfigNotFound = errorcode.ErrConfigNotFound
	ErrInvalidConfig  = errorcode.ErrInvalidConfig
	ErrMissingEnvVar  = errorcode.ErrMissingEnvVar
	ErrInvalidEnvVar  = errorcode.ErrInvalidEnvVar

	// Security errors (N7XX)
	ErrInvalidToken = errorcode.ErrInvalidToken
	ErrTokenExpired = errorcode.ErrTokenExpired
	ErrUnauthorized = errorcode.ErrUnauthorized
	ErrForbidden    = errorcode.ErrForbidden

	// Report generation errors (N8XX)
	ErrReportGeneration = errorcode.ErrReportGeneration
	ErrTemplateNotFound = errorcode.ErrTemplateNotFound
	ErrInvalidTemplate  = errorcode.ErrInvalidTemplate
	ErrPDFGeneration    = errorcode.ErrPDFGeneration
	ErrHTMLGeneration   = errorcode.ErrHTMLGeneration

	// Plugin and extension errors (N9XX)
	ErrPluginNotFound   = errorcode.ErrPluginNotFound
	ErrPluginLoadFailed = errorcode.ErrPluginLoadFailed
	ErrInvalidPlugin    = errorcode.ErrInvalidPlugin
	ErrPluginExecution  = errorcode.ErrPluginExecution

	// Viral growth features errors (V1XX)
	ErrViralShareFailed     = errorcode.ErrViralShareFailed
	ErrViralBadgeFailed     = errorcode.ErrViralBadgeFailed
	ErrViralCommunityFailed = errorcode.ErrViralCommunityFailed
	ErrViralPluginNotFound  = errorcode.ErrViralPluginNotFound
	ErrViralInvalidInput    = errorcode.ErrViralInvalidInput
)

// NewError is a compatibility function that wraps errors.NewError
func NewError(code ErrorCode, message string, details ...string) *NessiError {
	return errors.NewError(code, message, details...)
}

// GetErrorDescription is a compatibility function that wraps errorcode.GetErrorDescription
func GetErrorDescription(code ErrorCode) string {
	return errorcode.GetErrorDescription(code)
}

// WrapError is a compatibility function that wraps errors.WrapError
func WrapError(err error, code ErrorCode, message string, details ...string) *NessiError {
	return errors.WrapError(err, code, message, details...)
}

// NewInvalidPathError is a compatibility function that wraps errors.NewInvalidPathError
func NewInvalidPathError(path string) *NessiError {
	return errors.NewInvalidPathError(path)
}

// NewNotDeltaTableError is a compatibility function that wraps errors.NewNotDeltaTableError
func NewNotDeltaTableError(path string) *NessiError {
	return errors.NewNotDeltaTableError(path)
}

// NewAuthFailedError is a compatibility function that wraps errors.NewAuthFailedError
func NewAuthFailedError(service string) *NessiError {
	return errors.NewAuthFailedError(service)
}

// NewConnectionFailedError is a compatibility function that wraps errors.NewConnectionFailedError
func NewConnectionFailedError(service string, details string) *NessiError {
	return errors.NewConnectionFailedError(service, details)
}

// NewResourceNotFoundError is a compatibility function that wraps errors.NewResourceNotFoundError
func NewResourceNotFoundError(resourceType, resourceName string) *NessiError {
	return errors.NewResourceNotFoundError(resourceType, resourceName)
}

// NewSchemaValidationError is a compatibility function that wraps errors.NewSchemaValidationError
func NewSchemaValidationError(field, expectedType, actualType string) *NessiError {
	return errors.NewSchemaValidationError(field, expectedType, actualType)
}

// NewWorkspaceIDEmptyError is a compatibility function that wraps errors.NewWorkspaceIDEmptyError
func NewWorkspaceIDEmptyError() *NessiError {
	return errors.NewWorkspaceIDEmptyError()
}

// NewRateLimitExceededError is a compatibility function that wraps errors.NewRateLimitExceededError
func NewRateLimitExceededError(service string) *NessiError {
	return errors.NewRateLimitExceededError(service)
}

// NewServerError is a compatibility function that wraps errors.NewServerError
func NewServerError(service string, statusCode int) *NessiError {
	return errors.NewServerError(service, statusCode)
}

// NewViralBadgeError creates a new viral badge error
func NewViralBadgeError(message string, details string) *NessiError {
	return NewError(ErrViralBadgeFailed, message).WithDetails(details)
}

// NewViralCommunityError creates a new viral community error
func NewViralCommunityError(message string, details string) *NessiError {
	return NewError(ErrViralCommunityFailed, message).WithDetails(details)
}

// NewViralShareError creates a new viral share error
func NewViralShareError(message string, details string) *NessiError {
	return NewError(ErrViralShareFailed, message).WithDetails(details)
}

// NewViralPluginNotFoundError creates a new viral plugin not found error
func NewViralPluginNotFoundError(pluginName string) *NessiError {
	return NewError(ErrViralPluginNotFound, "Viral plugin not found: "+pluginName)
}

// NewViralInvalidInputError creates a new viral invalid input error
func NewViralInvalidInputError(message string, details string) *NessiError {
	return NewError(ErrViralInvalidInput, message).WithDetails(details)
}
