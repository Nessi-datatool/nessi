// Package security provides authentication, authorization, and SSL/TLS support
package security

// This file contains legacy code that is kept for backward compatibility.
// New code should use the auth.go, context.go, and ssl.go files instead.

import (
	"fmt"
)

// SSLManager manages SSL/TLS configuration
type SSLManager struct {
	Config SSLConfig
}

// SecurityManager manages authentication and authorization
type SecurityManager struct {
	AuthManager *AuthManager
	SSLManager  *SSLManager
}

// NewSecurityManager creates a new SecurityManager
func NewSecurityManager(config SecurityConfig) (*SecurityManager, error) {
	// Create auth manager
	authManager, err := NewAuthManager(config.Auth)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth manager: %w", err)
	}

	// Create SSL manager
	sslManager := &SSLManager{
		Config: config.SSL,
	}

	return &SecurityManager{
		AuthManager: authManager,
		SSLManager:  sslManager,
	}, nil
}

// New creates a new SecurityManager instance that wraps AuthManager
func New(apiKeyPath string) *SecurityManager {
	// Create auth config
	config := AuthConfig{
		Enabled:    true,
		APIKeyPath: apiKeyPath,
	}

	// Create auth manager
	am, err := NewAuthManager(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create AuthManager: %v", err))
	}

	return &SecurityManager{
		AuthManager: am,
	}
}

// GetAPIKeyForUser returns the API key for a user
func (s *SecurityManager) GetAPIKeyForUser(username string) (string, error) {
	apiKey, exists := s.AuthManager.GetAPIKey(username)
	if !exists {
		return "", fmt.Errorf("no API key found for user %s", username)
	}
	return apiKey, nil
}

// CreateAPIKey creates a new API key for a user
func (s *SecurityManager) CreateAPIKey(username string) (string, error) {
	return s.AuthManager.RegenerateAPIKey(username)
}

// ValidateAPIKey validates an API key
func (s *SecurityManager) ValidateAPIKey(apiKey string) (string, error) {
	user, err := s.AuthManager.ValidateAPIKey(apiKey)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

// AddUser adds a new user
func (s *SecurityManager) AddUser(username, password string, roles []string) error {
	// Convert roles to a single role
	role := RoleUser
	if len(roles) > 0 {
		switch roles[0] {
		case "admin":
			role = RoleAdmin
		case "readonly":
			role = RoleViewer
		default:
			role = RoleUser
		}
	}

	// Create user
	user := User{
		Username: username,
		Email:    username + "@example.com",
		Role:     role,
	}

	return s.AuthManager.CreateUser(user, password)
}
