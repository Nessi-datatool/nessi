// Package security provides authentication, authorization, and SSL/TLS support
package security

// This file contains legacy code that is kept for backward compatibility.
// New code should use the auth.go, context.go, and ssl.go files instead.

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecurityManager is a legacy authentication manager
// New code should use AuthManager instead
type SecurityManager struct {
	AuthManager *AuthManager
}

// LegacyTokenClaims represents JWT claims for the legacy system
type LegacyTokenClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// New creates a new SecurityManager instance that wraps AuthManager
func New(tokenSecret []byte, tokenExpiration time.Duration) *SecurityManager {
	// Create auth config
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    string(tokenSecret),
		TokenExpiry:  int(tokenExpiration.Hours()),
		RequireHTTPS: false,
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

// Authenticate authenticates a user
func (s *SecurityManager) Authenticate(username, password string) (string, error) {
	return s.AuthManager.Authenticate(username, password)
}

// ValidateToken validates a JWT token
func (s *SecurityManager) ValidateToken(tokenStr string) (*LegacyTokenClaims, error) {
	claims, err := s.AuthManager.ValidateToken(tokenStr)
	if err != nil {
		return nil, err
	}

	// Convert to legacy claims
	legacyClaims := &LegacyTokenClaims{
		Username: claims["username"].(string),
		Roles:    []string{claims["role"].(string)},
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// Safely handle exp and iat claims
	if exp, ok := claims["exp"]; ok && exp != nil {
		if expFloat, ok := exp.(float64); ok {
			legacyClaims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Unix(int64(expFloat), 0))
		}
	}

	if iat, ok := claims["iat"]; ok && iat != nil {
		if iatFloat, ok := iat.(float64); ok {
			legacyClaims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Unix(int64(iatFloat), 0))
		}
	}

	return legacyClaims, nil
}

// CreateAPIKey creates a new API key for a user
func (s *SecurityManager) CreateAPIKey(username string) (string, error) {
	return s.AuthManager.RegenerateAPIKey(username)
}

// ValidateAPIKey validates an API key
func (s *SecurityManager) ValidateAPIKey(apiKey string) (string, error) {
	user, err := s.AuthManager.AuthenticateWithAPIKey(apiKey)
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
