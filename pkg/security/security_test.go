package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLegacySecurityManager tests the backward compatibility of the SecurityManager
func TestLegacySecurityManager(t *testing.T) {
	// Use the optimized test helper to create a SecurityManager
	s := CreateTestSecurityManager()
	
	assert.NotNil(t, s)
	assert.NotNil(t, s.AuthManager)

	t.Run("AddUser and Authenticate", func(t *testing.T) {
		// Add test user
		var err error
		err = s.AddUser("legacyuser", "password", []string{"user"})
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
		var err error
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
