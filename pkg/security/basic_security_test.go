package security

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBasicSecurity tests the core security features
func TestBasicSecurity(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "basic-security-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Initialize with empty array
	_, err = tempFile.WriteString("[]")
	require.NoError(t, err)
	tempFile.Close()

	// Create auth config
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    "basic-test-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create auth manager
	am, err := NewAuthManager(config)
	require.NoError(t, err)
	require.NotNil(t, am)

	// Test user creation and authentication
	t.Run("User Auth", func(t *testing.T) {
		// Create test user
		testUser := User{
			Username: "basicuser",
			Email:    "basic@example.com",
			Role:     RoleUser,
		}
		err = am.CreateUser(testUser, "basic123")
		require.NoError(t, err)

		// Authenticate
		token, err := am.Authenticate("basicuser", "basic123")
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate token
		claims, err := am.ValidateToken(token)
		require.NoError(t, err)
		assert.Equal(t, "basicuser", claims["username"])
	})

	// Test middleware
	t.Run("Auth Middleware", func(t *testing.T) {
		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				http.Error(w, "User not found", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(user.Username))
		})

		// Create middleware
		middleware := am.AuthMiddleware(testHandler)

		// Get token
		token, err := am.Authenticate("basicuser", "basic123")
		require.NoError(t, err)

		// Test with valid token
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		middleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "basicuser", w.Body.String())
	})

	// Test role middleware
	t.Run("Role Middleware", func(t *testing.T) {
		// Create admin user
		adminUser := User{
			Username: "basicadmin",
			Email:    "basicadmin@example.com",
			Role:     RoleAdmin,
		}
		err = am.CreateUser(adminUser, "admin123")
		require.NoError(t, err)

		// Create test handler
		testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Create admin middleware
		adminMiddleware := am.RoleMiddleware(RoleAdmin)(testHandler)

		// Get admin token
		adminToken, err := am.Authenticate("basicadmin", "admin123")
		require.NoError(t, err)

		// Get user token
		userToken, err := am.Authenticate("basicuser", "basic123")
		require.NoError(t, err)

		// Test with admin token
		req := httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		adminMiddleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test with user token
		req = httptest.NewRequest("GET", "/admin", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w = httptest.NewRecorder()
		adminMiddleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
