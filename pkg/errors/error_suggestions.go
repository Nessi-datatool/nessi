package errors

import (
	stderrors "errors"
	"fmt"
	"strings"

	"github.com/nessi-dev/nessi/pkg/errorcode"
)

// ErrorSuggestion represents a suggestion for resolving an error
type ErrorSuggestion struct {
	ErrorCode   errorcode.ErrorCode // The error code this suggestion applies to
	Description string              // Description of the suggestion
	Solution    string              // The solution to try
	DocumentURL string              // URL to documentation for more information
}

// errorSuggestions is a map of error codes to suggestions
var errorSuggestions = map[errorcode.ErrorCode][]ErrorSuggestion{}

// RegisterErrorSuggestion registers a suggestion for an error code
func RegisterErrorSuggestion(suggestion ErrorSuggestion) {
	if _, ok := errorSuggestions[suggestion.ErrorCode]; !ok {
		errorSuggestions[suggestion.ErrorCode] = []ErrorSuggestion{}
	}
	errorSuggestions[suggestion.ErrorCode] = append(errorSuggestions[suggestion.ErrorCode], suggestion)
}

// GetSuggestionsForError returns suggestions for the given error
func GetSuggestionsForError(err error) []ErrorSuggestion {
	// Check if it's a NessiError
	var nessiErr *NessiError
	if !stderrors.As(err, &nessiErr) {
		return nil
	}

	// Return suggestions for this error code
	return errorSuggestions[nessiErr.Code]
}

// FormatSuggestion formats a suggestion for display
func FormatSuggestion(suggestion ErrorSuggestion) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📝 %s\n", suggestion.Description))
	sb.WriteString(fmt.Sprintf("🔧 Solution: %s\n", suggestion.Solution))

	if suggestion.DocumentURL != "" {
		sb.WriteString(fmt.Sprintf("📚 Documentation: %s\n", suggestion.DocumentURL))
	}

	return sb.String()
}

// FormatSuggestions formats all suggestions for an error
func FormatSuggestions(err error) string {
	suggestions := GetSuggestionsForError(err)
	if len(suggestions) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("\n🔍 Suggested solutions:\n\n")

	for i, suggestion := range suggestions {
		sb.WriteString(FormatSuggestion(suggestion))
		if i < len(suggestions)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// InitDefaultSuggestions initializes default suggestions for common errors
func InitDefaultSuggestions() {
	// Path errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrInvalidPath,
		Description: "The specified path does not exist or is not accessible",
		Solution:    "Check that the path exists and you have the necessary permissions to access it",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#path-errors",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrInvalidPath,
		Description: "Path might be using incorrect separators",
		Solution:    "Ensure you're using the correct path separators for your operating system (/ for Unix/Linux/macOS, \\ for Windows)",
		DocumentURL: "",
	})

	// Configuration errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrInvalidConfig,
		Description: "Configuration file is invalid or missing required fields",
		Solution:    "Check your configuration file for syntax errors or missing required fields",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#configuration-errors",
	})

	// Authentication errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrAuthFailed,
		Description: "Authentication failed due to invalid credentials",
		Solution:    "Check your credentials and ensure they are correctly configured",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#authentication-errors",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrInvalidToken,
		Description: "Token may have expired",
		Solution:    "Try regenerating your authentication token",
		DocumentURL: "",
	})

	// Connection errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrConnectionFailed,
		Description: "Failed to connect to the remote server",
		Solution:    "Check your network connection and ensure the server is running",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#connection-errors",
	})

	// Timeout errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrTimeout,
		Description: "Operation timed out",
		Solution:    "Try increasing the timeout value or check if the server is under heavy load",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#timeout-errors",
	})

	// Delta Lake errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrNotDeltaTable,
		Description: "The specified path is not a valid Delta Lake table",
		Solution:    "Ensure the path points to a valid Delta Lake table with _delta_log directory",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#delta-lake-errors",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrCorruptedDeltaLog,
		Description: "The Delta Lake table appears to be corrupted",
		Solution:    "Try running 'nessi repair' command on the table to fix corruption issues",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#delta-lake-errors",
	})

	// Databricks errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrServerError,
		Description: "Databricks API returned an error",
		Solution:    "Check your Databricks workspace URL and token, and ensure you have the necessary permissions",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#databricks-errors",
	})

	// Resource errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrResourceNotFound,
		Description: "The requested resource was not found",
		Solution:    "Check that the resource exists and you have the necessary permissions to access it",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#resource-errors",
	})

	// Rate limiting errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrRateLimitExceeded,
		Description: "Rate limit exceeded for API requests",
		Solution:    "Reduce the frequency of requests or implement exponential backoff",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/ERROR_HANDLING.md#rate-limiting-errors",
	})

	// Viral growth features errors
	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrViralShareFailed,
		Description: "Failed to generate shareable report",
		Solution:    "Check that the table exists and you have the necessary permissions to access it",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/VIRAL_FEATURES_IMPLEMENTATION_SUMMARY.md",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrViralBadgeFailed,
		Description: "Failed to generate badge",
		Solution:    "Check that the badge plugin is properly installed and configured",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/VIRAL_FEATURES_IMPLEMENTATION_SUMMARY.md",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrViralCommunityFailed,
		Description: "Failed to submit community feedback",
		Solution:    "Check your network connection and try again later",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/VIRAL_FEATURES_IMPLEMENTATION_SUMMARY.md",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrViralPluginNotFound,
		Description: "Viral growth plugin not found",
		Solution:    "Check that the plugin is properly installed and configured",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/VIRAL_FEATURES_IMPLEMENTATION_SUMMARY.md",
	})

	RegisterErrorSuggestion(ErrorSuggestion{
		ErrorCode:   ErrViralInvalidInput,
		Description: "Invalid input for viral command",
		Solution:    "Check the command syntax and try again",
		DocumentURL: "https://github.com/nessi-dev/nessi/blob/main/docs/VIRAL_FEATURES_IMPLEMENTATION_SUMMARY.md",
	})
}
