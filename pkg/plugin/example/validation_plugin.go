package main

import (
	"fmt"
	"regexp"
)

// PluginMetadata represents the structure for plugin metadata
type PluginMetadata struct {
	Name        string
	Version     string
	Description string
	Author      string
	Type        string
	Enabled     bool
}

// PluginInfo contains metadata about this plugin
var PluginInfo = &PluginMetadata{
	Name:        "email-validator",
	Version:     "1.0.0",
	Description: "Validates email addresses using regular expressions",
	Author:      "Nessi Team",
	Type:        "validation",
	Enabled:     true,
}

// Validate validates an email address
func Validate(args ...interface{}) ([]interface{}, error) {
	if len(args) < 1 {
		return []interface{}{false, "No email provided"}, nil
	}

	email, ok := args[0].(string)
	if !ok {
		return []interface{}{false, "Email must be a string"}, nil
	}

	// Simple email validation regex
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	if err != nil {
		return nil, fmt.Errorf("regex error: %v", err)
	}

	if matched {
		return []interface{}{true, "Email is valid"}, nil
	}

	return []interface{}{false, "Email is invalid"}, nil
}

// GetOptions returns the available options for this plugin
func GetOptions() map[string]interface{} {
	return map[string]interface{}{
		"strict_mode": false,
		"allow_local_domains": true,
	}
}

// main is required for Go plugins
func main() {}
