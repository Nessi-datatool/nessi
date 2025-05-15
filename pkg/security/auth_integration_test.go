package security_test

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/nessi-dev/nessi-dev/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestAuthIntegration(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "users-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create auth config
	config := security.AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-integration-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create auth manager
	am, err := security.NewAuthManager(config)
	require.NoError(t, err)
	require.NotNil(t, am)

	// Test user creation, authentication, and middleware in sequence
	t.Run("Full Authentication Flow", func(t *testing.T) {
		// 1. Create a test user
		testUser := security.User{
			Username: "integrationuser",
			Email:    "integration@example.com",
			Role:     security.RoleUser,
		}
		err = am.CreateUser(testUser, "integration123")
		require.NoError(t, err)

		// 2. Authenticate with username/password
		token, err := am.Authenticate("integrationuser", "integration123")
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// 3. Generate API key
		apiKey, err := am.RegenerateAPIKey("integrationuser")
		require.NoError(t, err)
		assert.NotEmpty(t, apiKey)

		// 4. Test JWT middleware
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
		middleware := am.AuthMiddleware(testHandler)

		// Test with JWT token
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "integrationuser", w.Body.String())

		// 5. Test API key middleware
		req = httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-API-Key", apiKey)
		w = httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "integrationuser", w.Body.String())

		// 6. Test role middleware
		adminOnly := am.RoleMiddleware(security.RoleAdmin)(testHandler)
		userOnly := am.RoleMiddleware(security.RoleUser)(testHandler)

		// User should not have access to admin routes
		req = httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		adminOnly.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)

		// User should have access to user routes
		req = httptest.NewRequest("GET", "/user", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		userOnly.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 7. Test token validation
		claims, err := am.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "integrationuser", claims["username"])

		// 8. Test token refresh
		refreshedToken, err := am.RefreshToken("integrationuser")
		require.NoError(t, err)
		assert.NotEmpty(t, refreshedToken)
		assert.NotEqual(t, token, refreshedToken)

		// Validate refreshed token
		refreshedClaims, err := am.ValidateToken(refreshedToken)
		require.NoError(t, err)
		assert.Equal(t, "integrationuser", refreshedClaims["username"])
	})

	t.Run("Login Handler", func(t *testing.T) {
		// Create a test user
		testUser := security.User{
			Username: "loginuser",
			Email:    "login@example.com",
			Role:     security.RoleUser,
		}
		err = am.CreateUser(testUser, "login123")
		require.NoError(t, err)

		// Create login request
		loginReq := map[string]string{
			"username": "loginuser",
			"password": "login123",
		}
		reqBody, err := json.Marshal(loginReq)
		require.NoError(t, err)

		// Create request
		req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(string(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Call login handler
		am.LoginHandler(w, req)

		// Check response
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var resp map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		// Check token
		assert.NotEmpty(t, resp["token"])
	})

	t.Run("Invalid Authentication", func(t *testing.T) {
		// Test with wrong password
		_, err := am.Authenticate("integrationuser", "wrongpassword")
		assert.Error(t, err)

		// Test with non-existent user
		_, err = am.Authenticate("nonexistentuser", "password")
		assert.Error(t, err)

		// Test with invalid token
		_, err = am.ValidateToken("invalid.token.string")
		assert.Error(t, err)

		// Test with invalid API key
		_, err = am.AuthenticateWithAPIKey("invalid-api-key")
		assert.Error(t, err)
	})

	t.Run("User Management", func(t *testing.T) {
		// Create admin user
		adminUser := security.User{
			Username: "adminuser",
			Email:    "admin@example.com",
			Role:     security.RoleAdmin,
		}
		err = am.CreateUser(adminUser, "admin123")
		require.NoError(t, err)

		// Get users
		users, err := am.GetUsers()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 3) // admin + integrationuser + loginuser + default admin

		// Update user
		updates := map[string]interface{}{
			"email": "updated@example.com",
		}
		err = am.UpdateUser("integrationuser", updates)
		require.NoError(t, err)

		// Verify update
		user, err := am.GetUser("integrationuser")
		require.NoError(t, err)
		assert.Equal(t, "updated@example.com", user.Email)

		// Delete user
		err = am.DeleteUser("loginuser")
		require.NoError(t, err)

		// Verify deletion
		_, err = am.GetUser("loginuser")
		assert.Error(t, err)
	})
}

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestAuthManagerPersistence(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "users-persistence-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create auth config
	config := security.AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-persistence-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create first auth manager instance
	am1, err := security.NewAuthManager(config)
	require.NoError(t, err)

	// Create a test user
	testUser := security.User{
		Username: "persistenceuser",
		Email:    "persistence@example.com",
		Role:     security.RoleUser,
	}
	err = am1.CreateUser(testUser, "persistence123")
	require.NoError(t, err)

	// Generate API key
	apiKey, err := am1.RegenerateAPIKey("persistenceuser")
	require.NoError(t, err)
	assert.NotEmpty(t, apiKey)

	// Create second auth manager instance (simulating restart)
	am2, err := security.NewAuthManager(config)
	require.NoError(t, err)

	// Verify user exists in second instance
	user, err := am2.GetUser("persistenceuser")
	require.NoError(t, err)
	assert.Equal(t, "persistence@example.com", user.Email)

	// Verify API key works in second instance
	authUser, err := am2.AuthenticateWithAPIKey(apiKey)
	require.NoError(t, err)
	assert.Equal(t, "persistenceuser", authUser.Username)
}

func (t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
TestRequireHTTPS(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "users-https-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Create auth config with HTTPS required
	config := security.AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-https-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: true,
	}

	// Create auth manager
	am, err := security.NewAuthManager(config)
	require.NoError(t, err)

	// Create a test user
	testUser := security.User{
		Username: "httpsuser",
		Email:    "https@example.com",
		Role:     security.RoleUser,
	}
	err = am.CreateUser(testUser, "https123")
	require.NoError(t, err)

	// Get token
	token, err := am.Authenticate("httpsuser", "https123")
	require.NoError(t, err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := am.AuthMiddleware(testHandler)

	// Test with HTTP request (should fail)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Test with HTTPS request (should succeed)
	req = httptest.NewRequest("GET", "https://example.com/test", nil)
	req.TLS = &tls.ConnectionState{} // Simulate HTTPS
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
