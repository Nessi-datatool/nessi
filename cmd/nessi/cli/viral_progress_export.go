package cli

// ShowProgress exports the ShowProgress function for testing
func ShowProgress(message string, duration time.Duration) {
	// Call the internal function
	showProgress(message, duration)
}

// ShowSuccess exports the ShowSuccess function for testing
func ShowSuccess(message string) {
	// Call the internal function
	showSuccess(message)
}

// ShowWarning exports the ShowWarning function for testing
func ShowWarning(message string) {
	// Call the internal function
	showWarning(message)
}
