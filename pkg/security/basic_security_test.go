package security

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBasicSecurity tests the core security features
func TestBasicSecurity(t *testing.T) {
	// Run in parallel for faster execution
	t.Parallel()

	// Use the helper to run with timeout
	RunWithTimeout(t, func() {
	// Use the optimized test helper to create an AuthManager
	am, err := CreateTestAuthManager()
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

		// We'll manually add users to the context instead of using tokens
		// This avoids potential token validation issues in the tests

		// For testing, we need to manually add the user to the context
		// Test with admin token
		req := httptest.NewRequest("GET", "/admin", nil)
		// Manually add admin user to context for testing
		adminUser, err = am.GetUser("basicadmin")
		require.NoError(t, err)
		ctx := WithUser(req.Context(), adminUser)
		req = req.WithContext(ctx)
		
		w := httptest.NewRecorder()
		adminMiddleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// Test with user token
		req = httptest.NewRequest("GET", "/admin", nil)
		// Manually add regular user to context for testing
		regularUser, err := am.GetUser("basicuser")
		require.NoError(t, err)
		ctx = WithUser(req.Context(), regularUser)
		req = req.WithContext(ctx)
		
		w = httptest.NewRecorder()
		adminMiddleware.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
	})
}
