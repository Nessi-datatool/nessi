package cli

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

// Helper functions for viral commands

// generateShareID generates a random share ID
func generateShareID() string {
	// Generate a random 8-character alphanumeric ID
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, 8)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

// generateBadgeURL generates a badge URL based on parameters
func generateBadgeURL(label, message, color, style string) string {
	// URL encode the label and message
	encodedLabel := url.QueryEscape(label)
	encodedMessage := url.QueryEscape(message)

	// Generate the badge URL
	return fmt.Sprintf("https://img.shields.io/badge/%s-%s-%s?style=%s",
		encodedLabel, encodedMessage, color, style)
}

// These error handlers are already defined in viral_error_handler.go
// No need to redefine them here
