package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/nessi-dev/nessi/pkg/common"
)

// ViralErrorCodes defines error codes specific to viral growth features
const (
	ErrViralShareFailed     = "V101" // Error code for share command failures
	ErrViralBadgeFailed     = "V102" // Error code for badge command failures
	ErrViralCommunityFailed = "V103" // Error code for community command failures
	ErrViralPluginNotFound  = "V104" // Error code for plugin not found errors
	ErrViralInvalidInput    = "V105" // Error code for invalid input errors
)

// ViralErrorHandler handles errors specific to viral growth features
type ViralErrorHandler struct {
	BaseHandler *ErrorHandler
}

// NewViralErrorHandler creates a new viral error handler
func NewViralErrorHandler() *ViralErrorHandler {
	return &ViralErrorHandler{
		BaseHandler: NewErrorHandler(true),
	}
}

// HandleViralError handles errors specific to viral growth features
func HandleViralError(err error, feature string) {
	// Check if it's a NessiError
	var nessiErr *common.NessiError
	if err == nil {
		return
	}

	if errors.As(err, &nessiErr) {
		// Handle based on the feature
		switch feature {
		case "share":
			handleShareError(nessiErr)
		case "badge":
			handleBadgeError(nessiErr)
		case "community":
			handleCommunityError(nessiErr)
		default:
			// Use the base handler for unknown features
			HandleError(err)
		}
	} else {
		// Create a new NessiError based on the feature
		code := getErrorCodeForFeature(feature)
		nessiErr := &common.NessiError{
			Code:    code,
			Message: fmt.Sprintf("Error in viral %s feature: %s", feature, err.Error()),
			Details: err.Error(),
		}

		// Handle the error
		HandleError(nessiErr)
	}
}

// getErrorCodeForFeature returns the error code for a specific feature
func getErrorCodeForFeature(feature string) string {
	switch feature {
	case "share":
		return ErrViralShareFailed
	case "badge":
		return ErrViralBadgeFailed
	case "community":
		return ErrViralCommunityFailed
	default:
		return "V100" // Generic viral feature error
	}
}

// handleShareError handles errors specific to the share feature
func handleShareError(err *common.NessiError) {
	// Print error with custom formatting
	error := color.New(color.FgRed, color.Bold)
	error.Printf("\n❌ Sharing Error: %s\n", err.Message)

	// Print error code
	fmt.Printf("Code: %s\n", err.Code)

	// Print details if available
	if err.Details != "" {
		fmt.Printf("Details: %s\n", err.Details)
	}

	// Print suggestions specific to share errors
	fmt.Println("\nSuggestions:")
	fmt.Println("1. Check if the table exists and is accessible")
	fmt.Println("2. Verify that the social sharing plugin is installed correctly")
	fmt.Println("3. Try using a different output format or path")
	fmt.Println("4. Run 'nessi viral share --help' for usage information")
}

// handleBadgeError handles errors specific to the badge feature
func handleBadgeError(err *common.NessiError) {
	// Print error with custom formatting
	error := color.New(color.FgRed, color.Bold)
	error.Printf("\n❌ Badge Generation Error: %s\n", err.Message)

	// Print error code
	fmt.Printf("Code: %s\n", err.Code)

	// Print details if available
	if err.Details != "" {
		fmt.Printf("Details: %s\n", err.Details)
	}

	// Print suggestions specific to badge errors
	fmt.Println("\nSuggestions:")
	fmt.Println("1. Check if the badge plugin is installed correctly")
	fmt.Println("2. Try using a different badge format (markdown, html)")
	fmt.Println("3. Verify that the color value is valid (hex or named color)")
	fmt.Println("4. Run 'nessi viral badge --help' for usage information")
}

// handleCommunityError handles errors specific to the community feature
func handleCommunityError(err *common.NessiError) {
	// Print error with custom formatting
	error := color.New(color.FgRed, color.Bold)
	error.Printf("\n❌ Community Engagement Error: %s\n", err.Message)

	// Print error code
	fmt.Printf("Code: %s\n", err.Code)

	// Print details if available
	if err.Details != "" {
		fmt.Printf("Details: %s\n", err.Details)
	}

	// Print suggestions specific to community errors
	fmt.Println("\nSuggestions:")
	fmt.Println("1. Check your internet connection")
	fmt.Println("2. Verify that the community engagement plugin is installed correctly")
	fmt.Println("3. Try again later if the issue persists")
	fmt.Println("4. Run 'nessi viral community --help' for usage information")
}

// showProgress displays a progress indicator for viral feature operations
func showProgress(message string, duration time.Duration) {
	// Create a new spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + message
	s.Color("green")

	// Start the spinner
	s.Start()

	// Wait for the specified duration
	time.Sleep(duration)

	// Stop the spinner
	s.Stop()
}

// showSuccess displays a success message for viral feature operations
func showSuccess(message string) {
	success := color.New(color.FgGreen, color.Bold)
	success.Printf("\n✅ %s\n", message)
}

// showWarning displays a warning message for viral feature operations
func showWarning(message string) {
	warning := color.New(color.FgYellow, color.Bold)
	warning.Printf("\n⚠️ %s\n", message)
}

// PromptForConfirmation prompts the user for confirmation
func PromptForConfirmation(message string) bool {
	fmt.Printf("\n%s [y/N]: ", message)

	// Read user input
	var input string
	fmt.Scanln(&input)
	input = strings.ToLower(strings.TrimSpace(input))

	// Return true if user confirms
	return input == "y" || input == "yes"
}
