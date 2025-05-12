package security

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLegacySecurityManager tests the backward compatibility of the SecurityManager
func TestLegacySecurityManager(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "security-legacy-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Initialize the file with empty JSON object
	_, err = tempFile.WriteString("{}")
	require.NoError(t, err)
	tempFile.Close()

	tokenSecret := []byte("test-legacy-secret")
	tokenExpiration := 24 * time.Hour
	
	// Create auth config directly
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    string(tokenSecret),
		TokenExpiry:  int(tokenExpiration.Hours()),
		UsersFile:    tempFile.Name(),
		RequireHTTPS: false,
	}
	
	// Create auth manager directly
	am, err := NewAuthManager(config)
	require.NoError(t, err)
	
	// Create security manager with the auth manager
	s := &SecurityManager{
		AuthManager: am,
	}
	
	assert.NotNil(t, s)
	assert.NotNil(t, s.AuthManager)

	t.Run("AddUser and Authenticate", func(t *testing.T) {
		// Add test user
		err := s.AddUser("legacyuser", "password", []string{"user"})
		require.NoError(t, err)

		// Test with correct credentials
		token, err := s.Authenticate("legacyuser", "password")
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Test with incorrect password
		token, err = s.Authenticate("legacyuser", "wrongpassword")
		require.Error(t, err)
		assert.Empty(t, token)

		// Test with non-existent user
		token, err = s.Authenticate("nonexistent", "password")
		require.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("ValidateToken", func(t *testing.T) {
		// Get valid token
		token, err := s.Authenticate("legacyuser", "password")
		require.NoError(t, err)

		// Validate token
		claims, err := s.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "legacyuser", claims.Username)
		assert.Contains(t, claims.Roles, "user")
	})

	t.Run("API Key Management", func(t *testing.T) {
		// Create API key
		apiKey, err := s.CreateAPIKey("legacyuser")
		require.NoError(t, err)
		assert.NotEmpty(t, apiKey)

		// Validate API key
		username, err := s.ValidateAPIKey(apiKey)
		require.NoError(t, err)
		assert.Equal(t, "legacyuser", username)

		// Validate non-existent API key
		_, err = s.ValidateAPIKey("nonexistent-key")
		require.Error(t, err)
	})

	t.Run("Add Duplicate User", func(t *testing.T) {
		// Try adding same user again
		err = s.AddUser("legacyuser", "password", []string{"user"})
		require.Error(t, err)
	})

	t.Run("Add Admin User", func(t *testing.T) {
		// Add admin user
		err := s.AddUser("legacyadmin", "adminpass", []string{"admin"})
		require.NoError(t, err)

		// Verify user can authenticate
		token, err := s.Authenticate("legacyadmin", "adminpass")
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate token and check role
		claims, err := s.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "legacyadmin", claims.Username)
		assert.Contains(t, claims.Roles, "admin")
	})
}
