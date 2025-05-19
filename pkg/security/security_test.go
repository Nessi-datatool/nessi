package security

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSecurityManager tests the SecurityManager for CLI-only approach
func TestSecurityManager(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "security-test-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test config
	config := Config{
		Auth: AuthConfig{
			Enabled:      true,
			UsersFile:    tempDir + "/users.json",
			InMemoryOnly: true,
		},
		SSL: SSLConfig{
			CertFile: tempDir + "/cert.pem",
			KeyFile:  tempDir + "/key.pem",
		},
	}

	// Create security manager
	s, err := NewSecurityManager(config)
	require.NoError(t, err)
	assert.NotNil(t, s)
	assert.NotNil(t, s.AuthManager)
	assert.NotNil(t, s.SSLManager)

	t.Run("User Management", func(t *testing.T) {
		// Create test user
		testUser := User{
			Username: "securityuser",
			Email:    "security@example.com",
			Role:     RoleUser,
		}
		err = s.AuthManager.CreateUser(testUser, "password123")
		require.NoError(t, err)

		// Get user
		user, err := s.AuthManager.GetUser("securityuser")
		require.NoError(t, err)
		assert.Equal(t, "securityuser", user.Username)
		assert.Equal(t, "security@example.com", user.Email)
		assert.Equal(t, RoleUser, user.Role)
	})

	t.Run("API Key Management", func(t *testing.T) {
		// Generate API key
		apiKey, err := s.AuthManager.RegenerateAPIKey("securityuser")
		require.NoError(t, err)
		assert.NotEmpty(t, apiKey)

		// Authenticate with API key
		user, err := s.AuthManager.AuthenticateWithAPIKey(apiKey)
		require.NoError(t, err)
		assert.Equal(t, "securityuser", user.Username)

		// Test with invalid API key
		_, err = s.AuthManager.AuthenticateWithAPIKey("invalid-key")
		assert.Error(t, err)
	})
}
