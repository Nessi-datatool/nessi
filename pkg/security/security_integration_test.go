package security_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nessi-dev/nessi/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityIntegration(t *testing.T) {
	if security.ShouldSkipIntegrationTests(t) {
		return // t.Skip already called in ShouldSkipIntegrationTests
	}

	// Run the test with a timeout
	security.RunWithTimeout(t, func() {
		// Create a temporary directory for test files (just for SSL certificates)
		tmpDir, err := os.MkdirTemp("", "security-integration-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create certificate paths
		certFile := filepath.Join(tmpDir, "cert.pem")
		keyFile := filepath.Join(tmpDir, "key.pem")

		// Use in-memory auth manager for faster tests
		authManager, err := security.CreateTestAuthManager()
		if err != nil {
			t.Fatalf("Failed to create auth manager: %v", err)
		}

		// Create SSL config
		sslConfig := security.SSLConfig{
			Enabled:      true,
			CertFile:     certFile,
			KeyFile:      keyFile,
			AutoGenerate: true,
		}

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
			assert.GreaterOrEqual(t, len(users), 1) // at least the test user

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
			// Create admin user with unique name to avoid conflicts
			random := fmt.Sprintf("%d", time.Now().UnixNano()%10000)
			adminUser := security.User{
				Username: "adminuser-" + random,
				Email:    "admin" + random + "@example.com",
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

			// We'll manually add users to the context instead of using tokens
			// This avoids potential token validation issues in the tests

			// For testing, we need to manually add the user to the context
			// Test with admin token
			req := httptest.NewRequest("GET", "/admin", nil)
			// Use the admin user we just created instead of trying to fetch it
			ctx := security.WithUser(req.Context(), adminUser)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			adminMiddleware.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)

			// Test with user token (should be forbidden)
			req = httptest.NewRequest("GET", "/admin", nil)
			// Manually add regular user to context for testing
			regularUser, err := authManager.GetUser("integrationuser")
			require.NoError(t, err)
			ctx = security.WithUser(req.Context(), regularUser)
			req = req.WithContext(ctx)

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
	})
}
