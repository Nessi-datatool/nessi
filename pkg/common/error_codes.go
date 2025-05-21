package common

// Error codes for Nessi
// Format: NXXX where XXX is a 3-digit number
// N1XX: General errors
// N2XX: Delta Lake errors
// N3XX: Databricks errors
// N4XX: Cloud provider errors
// N5XX: Schema and data errors
// N6XX: Configuration errors
// N7XX: Security errors
// N8XX: Report generation errors
// N9XX: Plugin and extension errors

// ErrorCode represents a unique error code in Nessi
type ErrorCode string

const (
	// General errors (N1XX)
	ErrInvalidPath       ErrorCode = "N101" // Invalid path
	ErrFileNotFound      ErrorCode = "N102" // File not found
	ErrPermissionDenied  ErrorCode = "N103" // Permission denied
	ErrInvalidArgument   ErrorCode = "N104" // Invalid argument
	ErrInvalidOperation  ErrorCode = "N105" // Invalid operation
	ErrTimeout           ErrorCode = "N106" // Operation timeout
	ErrInternal          ErrorCode = "N107" // Internal error

	// Delta Lake errors (N2XX)
	ErrNotDeltaTable     ErrorCode = "N201" // Not a Delta table
	ErrCorruptedDeltaLog ErrorCode = "N202" // Corrupted Delta log
	ErrInvalidVersion    ErrorCode = "N203" // Invalid version
	ErrEmptyTable        ErrorCode = "N204" // Empty table
	ErrDeltaLogNotFound  ErrorCode = "N205" // Delta log not found
	ErrParquetReadFailed ErrorCode = "N206" // Failed to read parquet file

	// Databricks errors (N3XX)
	ErrAuthFailed        ErrorCode = "N301" // Authentication failed
	ErrConnectionFailed  ErrorCode = "N302" // Connection failed
	ErrResourceNotFound  ErrorCode = "N303" // Resource not found
	ErrRateLimitExceeded ErrorCode = "N304" // Rate limit exceeded
	ErrServerError       ErrorCode = "N305" // Server error
	ErrWorkspaceIDEmpty  ErrorCode = "N306" // Workspace ID is empty
	ErrInvalidResponse   ErrorCode = "N307" // Invalid response from server

	// Cloud provider errors (N4XX)
	ErrAWSAuthFailed     ErrorCode = "N401" // AWS authentication failed
	ErrGCPAuthFailed     ErrorCode = "N402" // GCP authentication failed
	ErrAzureAuthFailed   ErrorCode = "N403" // Azure authentication failed
	ErrS3BucketNotFound  ErrorCode = "N404" // S3 bucket not found
	ErrGCSBucketNotFound ErrorCode = "N405" // GCS bucket not found
	ErrADLSNotFound      ErrorCode = "N406" // ADLS not found

	// Schema and data errors (N5XX)
	ErrSchemaValidation  ErrorCode = "N501" // Schema validation failed
	ErrDataTypeConversion ErrorCode = "N502" // Data type conversion failed
	ErrInvalidSchema     ErrorCode = "N503" // Invalid schema
	ErrSchemaMismatch    ErrorCode = "N504" // Schema mismatch
	ErrDataValidation    ErrorCode = "N505" // Data validation failed

	// Configuration errors (N6XX)
	ErrConfigNotFound    ErrorCode = "N601" // Configuration not found
	ErrInvalidConfig     ErrorCode = "N602" // Invalid configuration
	ErrMissingEnvVar     ErrorCode = "N603" // Missing environment variable
	ErrInvalidEnvVar     ErrorCode = "N604" // Invalid environment variable

	// Security errors (N7XX)
	ErrInvalidToken      ErrorCode = "N701" // Invalid token
	ErrTokenExpired      ErrorCode = "N702" // Token expired
	ErrUnauthorized      ErrorCode = "N703" // Unauthorized
	ErrForbidden         ErrorCode = "N704" // Forbidden

	// Report generation errors (N8XX)
	ErrReportGeneration  ErrorCode = "N801" // Report generation failed
	ErrTemplateNotFound  ErrorCode = "N802" // Template not found
	ErrInvalidTemplate   ErrorCode = "N803" // Invalid template
	ErrPDFGeneration     ErrorCode = "N804" // PDF generation failed
	ErrHTMLGeneration    ErrorCode = "N805" // HTML generation failed

	// Plugin and extension errors (N9XX)
	ErrPluginNotFound    ErrorCode = "N901" // Plugin not found
	ErrPluginLoadFailed  ErrorCode = "N902" // Plugin load failed
	ErrInvalidPlugin     ErrorCode = "N903" // Invalid plugin
	ErrPluginExecution   ErrorCode = "N904" // Plugin execution failed
)

// ErrorCodeMap maps error codes to their descriptions
var ErrorCodeMap = map[ErrorCode]string{
	// General errors
	ErrInvalidPath:       "Invalid path",
	ErrFileNotFound:      "File not found",
	ErrPermissionDenied:  "Permission denied",
	ErrInvalidArgument:   "Invalid argument",
	ErrInvalidOperation:  "Invalid operation",
	ErrTimeout:           "Operation timeout",
	ErrInternal:          "Internal error",

	// Delta Lake errors
	ErrNotDeltaTable:     "Not a Delta table",
	ErrCorruptedDeltaLog: "Corrupted Delta log",
	ErrInvalidVersion:    "Invalid version",
	ErrEmptyTable:        "Empty table",
	ErrDeltaLogNotFound:  "Delta log not found",
	ErrParquetReadFailed: "Failed to read parquet file",

	// Databricks errors
	ErrAuthFailed:        "Authentication failed",
	ErrConnectionFailed:  "Connection failed",
	ErrResourceNotFound:  "Resource not found",
	ErrRateLimitExceeded: "Rate limit exceeded",
	ErrServerError:       "Server error",
	ErrWorkspaceIDEmpty:  "Workspace ID is empty",
	ErrInvalidResponse:   "Invalid response from server",

	// Cloud provider errors
	ErrAWSAuthFailed:     "AWS authentication failed",
	ErrGCPAuthFailed:     "GCP authentication failed",
	ErrAzureAuthFailed:   "Azure authentication failed",
	ErrS3BucketNotFound:  "S3 bucket not found",
	ErrGCSBucketNotFound: "GCS bucket not found",
	ErrADLSNotFound:      "ADLS not found",

	// Schema and data errors
	ErrSchemaValidation:  "Schema validation failed",
	ErrDataTypeConversion: "Data type conversion failed",
	ErrInvalidSchema:     "Invalid schema",
	ErrSchemaMismatch:    "Schema mismatch",
	ErrDataValidation:    "Data validation failed",

	// Configuration errors
	ErrConfigNotFound:    "Configuration not found",
	ErrInvalidConfig:     "Invalid configuration",
	ErrMissingEnvVar:     "Missing environment variable",
	ErrInvalidEnvVar:     "Invalid environment variable",

	// Security errors
	ErrInvalidToken:      "Invalid token",
	ErrTokenExpired:      "Token expired",
	ErrUnauthorized:      "Unauthorized",
	ErrForbidden:         "Forbidden",

	// Report generation errors
	ErrReportGeneration:  "Report generation failed",
	ErrTemplateNotFound:  "Template not found",
	ErrInvalidTemplate:   "Invalid template",
	ErrPDFGeneration:     "PDF generation failed",
	ErrHTMLGeneration:    "HTML generation failed",

	// Plugin and extension errors
	ErrPluginNotFound:    "Plugin not found",
	ErrPluginLoadFailed:  "Plugin load failed",
	ErrInvalidPlugin:     "Invalid plugin",
	ErrPluginExecution:   "Plugin execution failed",
}

// GetErrorDescription returns the description for an error code
func GetErrorDescription(code ErrorCode) string {
	if desc, ok := ErrorCodeMap[code]; ok {
		return desc
	}
	return "Unknown error"
}
