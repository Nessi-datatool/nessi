package security_test

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityIntegration(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "security-integration-")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create users file path
	usersFilePath := tempDir + "/users.json"

	// Create empty users file with valid JSON array
	usersFile, err := os.Create(usersFilePath)
	require.NoError(t, err)
	_, err = usersFile.WriteString("[]")
	require.NoError(t, err)
	usersFile.Close()

	// Create certificate paths
	certFile := tempDir + "/cert.pem"
	keyFile := tempDir + "/key.pem"

	// Create auth config
	authConfig := security.AuthConfig{
		Enabled:      true,
		JWTSecret:    "integration-test-secret",
		UsersFile:    usersFilePath,
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create SSL config
	sslConfig := security.SSLConfig{
		Enabled:      true,
		CertFile:     certFile,
		KeyFile:      keyFile,
		AutoGenerate: true,
	}

	// Create auth manager
	authManager, err := security.NewAuthManager(authConfig)
	require.NoError(t, err)
	require.NotNil(t, authManager)

	// Create cert manager
	certManager := security.NewCertManager(sslConfig)
	require.NotNil(t, certManager)

	// Test user management and authentication
	t.Run("User Management", func(t *testing.T) {
		// Create test user
		testUser := security.User{
			Username: "integrationuser",
			Email:    "integration@example.com",
			Role:     security.RoleUser,
		}
		err = authManager.CreateUser(testUser, "integration123")
		require.NoError(t, err)

		// Get users
		users, err := authManager.GetUsers()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 2) // test user + default admin

		// Authenticate
		token, err := authManager.Authenticate("integrationuser", "integration123")
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate token
		claims, err := authManager.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "integrationuser", claims["username"])

		// Generate API key
		apiKey, err := authManager.RegenerateAPIKey("integrationuser")
		require.NoError(t, err)
		assert.NotEmpty(t, apiKey)

		// Authenticate with API key
		user, err := authManager.AuthenticateWithAPIKey(apiKey)
		require.NoError(t, err)
		assert.Equal(t, "integrationuser", user.Username)

		// Update user
		updates := map[string]interface{}{
			"email": "updated@example.com",
		}
		err = authManager.UpdateUser("integrationuser", updates)
		require.NoError(t, err)

		// Verify update
		user, err = authManager.GetUser("integrationuser")
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", user.Email)
	})

	// Test authentication middleware
	t.Run("Authentication Middleware", func(t *testing.T) {
		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := security.UserFromContext(r.Context())
			if !ok {
				http.Error(w, "User not found in context", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(user.Username))
		})

		// Create middleware
		middleware := authManager.AuthMiddleware(testHandler)

		// Get token
		token, err := authManager.Authenticate("integrationuser", "integration123")
		require.NoError(t, err)

		// Test with token
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "integrationuser", w.Body.String())

		// Test with API key
		apiKey, err := authManager.RegenerateAPIKey("integrationuser")
		require.NoError(t, err)

		req = httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-API-Key", apiKey)
		w = httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "integrationuser", w.Body.String())

		// Test with invalid token
		req = httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w = httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test role-based access control
	t.Run("Role-Based Access Control", func(t *testing.T) {
		// Create admin user
		adminUser := security.User{
			Username: "adminuser",
			Email:    "admin@example.com",
			Role:     security.RoleAdmin,
		}
		err = authManager.CreateUser(adminUser, "admin123")
		require.NoError(t, err)

		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Access granted"))
		})

		// Create admin middleware
		adminMiddleware := authManager.RoleMiddleware(security.RoleAdmin)(testHandler)

		// Get admin token
		adminToken, err := authManager.Authenticate("adminuser", "admin123")
		require.NoError(t, err)

		// Get user token
		userToken, err := authManager.Authenticate("integrationuser", "integration123")
		require.NoError(t, err)

		// Test with admin token
		req := httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		adminMiddleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test with user token (should be forbidden)
		req = httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w = httptest.NewRecorder()
		adminMiddleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	// Test HTTPS support
	t.Run("HTTPS Support", func(t *testing.T) {
		// Get TLS config
		tlsConfig, err := certManager.GetTLSConfig()
		require.NoError(t, err)
		require.NotNil(t, tlsConfig)

		// Verify certificate files were created
		_, err = os.Stat(certFile)
		assert.NoError(t, err)
		_, err = os.Stat(keyFile)
		assert.NoError(t, err)

		// Create test server
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("HTTPS works!"))
		}))
		defer server.Close()

		// Configure server with TLS
		server.TLS = tlsConfig
		server.StartTLS()

		// Create HTTP client that skips certificate verification
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		// Make HTTPS request
		resp, err := client.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	// Test user persistence
	t.Run("User Persistence", func(t *testing.T) {
		// Create a new auth manager with the same config (simulating restart)
		newAuthManager, err := security.NewAuthManager(authConfig)
		require.NoError(t, err)

		// Verify users still exist
		user, err := newAuthManager.GetUser("integrationuser")
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", user.Email)

		admin, err := newAuthManager.GetUser("adminuser")
		require.NoError(t, err)
		assert.Equal(t, security.RoleAdmin, admin.Role)

		// Verify authentication still works
		token, err := newAuthManager.Authenticate("integrationuser", "integration123")
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
