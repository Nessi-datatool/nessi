// Package security provides license validation and security features
package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// License represents a Nessi license
type License struct {
	Token     string    `json:"token"`
	Issued    time.Time `json:"issued"`
	Expires   time.Time `json:"expires"`
	Customer  string    `json:"customer"`
	Plan      string    `json:"plan"`
	Features  []string  `json:"features"`
	Signature string    `json:"signature"`
}

// LicenseClaims represents the JWT claims in a license token
type LicenseClaims struct {
	jwt.RegisteredClaims
	Customer string   `json:"customer"`
	Plan     string   `json:"plan"`
	Features []string `json:"features"`
}

// ErrLicenseExpired is returned when the license has expired
var ErrLicenseExpired = errors.New("license has expired")

// ErrLicenseInvalid is returned when the license is invalid
var ErrLicenseInvalid = errors.New("license is invalid")

// ErrLicenseNotFound is returned when the license file is not found
var ErrLicenseNotFound = errors.New("license file not found")

// ValidateLicense validates the license file
func ValidateLicense() (*License, error) {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	// License file path
	licensePath := filepath.Join(homeDir, ".nessi", "license.json")

	// Check if license file exists
	if _, err := os.Stat(licensePath); os.IsNotExist(err) {
		return nil, ErrLicenseNotFound
	}

	// Read license file
	data, err := os.ReadFile(licensePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read license file: %w", err)
	}

	// Parse license file
	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return nil, fmt.Errorf("failed to parse license file: %w", err)
	}

	// Verify JWT token
	if err := verifyToken(license.Token); err != nil {
		return nil, err
	}

	// Check expiry
	if time.Now().After(license.Expires) {
		return &license, ErrLicenseExpired
	}

	return &license, nil
}

// verifyToken verifies the JWT token in the license
func verifyToken(tokenString string) error {
	// This is a simplified implementation
	// In a real implementation, we would verify the signature using a public key
	
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &LicenseClaims{}, func(token *jwt.Token) (interface{}, error) {
		// In a real implementation, we would verify the signing method and return the public key
		return []byte("secret"), nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid
	if !token.Valid {
		return ErrLicenseInvalid
	}

	return nil
}

// HasFeature checks if the license has a specific feature
func HasFeature(license *License, feature string) bool {
	if license == nil {
		return false
	}

	for _, f := range license.Features {
		if f == feature {
			return true
		}
	}

	return false
}