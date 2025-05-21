package errors

import (
	"github.com/nessi-dev/nessi/pkg/errorcode"
)

// This file re-exports error codes from the errorcode package for backward compatibility

const (
	// General errors (N1XX)
	ErrInvalidPath      = errorcode.ErrInvalidPath
	ErrFileNotFound     = errorcode.ErrFileNotFound
	ErrPermissionDenied = errorcode.ErrPermissionDenied
	ErrInvalidArgument  = errorcode.ErrInvalidArgument
	ErrInvalidOperation = errorcode.ErrInvalidOperation
	ErrTimeout          = errorcode.ErrTimeout
	ErrInternal         = errorcode.ErrInternal

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

// GetErrorDescription is a wrapper for errorcode.GetErrorDescription
func GetErrorDescription(code ErrorCode) string {
	return errorcode.GetErrorDescription(code)
}
