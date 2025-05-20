package security_test

import (
	"os"
	"testing"

	"github.com/nessi-dev/nessi/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityIntegration(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "security-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create security config
	config := security.Config{
		Auth: security.AuthConfig{
			Enabled:      true,
			UsersFile:    tempDir + "/users.json",
			InMemoryOnly: true,
		},
		SSL: security.SSLConfig{
			Enabled:  true,
			CertFile: tempDir + "/cert.pem",
			KeyFile:  tempDir + "/key.pem",
		},
	}

	// Create security manager
	securityManager, err := security.NewSecurityManager(config)
	require.NoError(t, err)
	require.NotNil(t, securityManager)
	require.NotNil(t, securityManager.AuthManager)
	require.NotNil(t, securityManager.SSLManager)

	// Test user creation and API key authentication
	t.Run("User Management", func(t *testing.T) {
		// Create test user
		testUser := security.User{
			Username: "securityintegration",
			Email:    "secint@example.com",
			Role:     security.RoleAdmin,
		}

		err = securityManager.AuthManager.CreateUser(testUser, "secint123")
		require.NoError(t, err)

		// Generate API key
		apiKey, err := securityManager.AuthManager.RegenerateAPIKey("securityintegration")
		require.NoError(t, err)
		assert.NotEmpty(t, apiKey)

		// Authenticate with API key
		user, err := securityManager.AuthManager.AuthenticateWithAPIKey(apiKey)
		require.NoError(t, err)
		assert.Equal(t, "securityintegration", user.Username)
		assert.Equal(t, security.RoleAdmin, user.Role)
	})

	// Test SSL certificate generation
	t.Run("SSL Certificate Generation", func(t *testing.T) {
		// Generate self-signed certificate
		err := securityManager.SSLManager.GenerateSelfSignedCertForTest()
		require.NoError(t, err)

		// Verify certificate files were created
		_, err = os.Stat(config.SSL.CertFile)
		assert.NoError(t, err)
		_, err = os.Stat(config.SSL.KeyFile)
		assert.NoError(t, err)

		// Get TLS config
		tlsConfig, err := securityManager.SSLManager.GetTLSConfig()
		require.NoError(t, err)
		assert.NotNil(t, tlsConfig)
	})
}
