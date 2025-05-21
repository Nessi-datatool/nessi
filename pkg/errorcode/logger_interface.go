package errorcode

// Logger is a minimal interface for logging that can be used by error-related packages
// without creating import cycles
type Logger interface {
	// SimpleDebug logs a debug message
	SimpleDebug(msg string)

	// SimpleInfo logs an informational message
	SimpleInfo(msg string)

	// SimpleWarn logs a warning message
	SimpleWarn(msg string)

	// SimpleError logs an error message
	SimpleError(msg string)

	// SimpleWithField adds a field to the logger
	SimpleWithField(key string, value interface{}) Logger
}
