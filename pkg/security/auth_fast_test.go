package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthMiddlewareFast is a fast version of the auth middleware test
func TestAuthMiddlewareFast(t *testing.T) {
	// Make this test run in parallel with other tests
	t.Parallel()

	// Create a simple mock auth manager with InMemoryOnly to prevent file I/O
	am := &AuthManager{
		config: AuthConfig{
			Enabled:      true,
			JWTSecret:    "test-secret",
			InMemoryOnly: true,
		},
		users: map[string]User{
			"testuser": {
				Username: "testuser",
				Email:    "test@example.com",
				Role:     RoleUser,
				APIKey:   "test-api-key",
			},
		},
		apiKeys: map[string]string{
			"test-api-key": "testuser",
		},
	}

	// Create a JWT token manually instead of using the method
	claims := jwt.MapClaims{
		"username":   "testuser",
		"email":      "test@example.com",
		"role":       RoleUser,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
		"created_at": time.Now().Unix(),
	}
	
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token
	tokenString, err := token.SignedString([]byte("test-secret"))
	require.NoError(t, err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is in context
		_, ok := UserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	middleware := am.AuthMiddleware(testHandler)

	// Test with token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test with API key
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test with invalid token
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestRoleMiddlewareFast is a fast version of the role middleware test
func TestRoleMiddlewareFast(t *testing.T) {
	// Make this test run in parallel with other tests
	t.Parallel()

	// Create a simple mock auth manager with InMemoryOnly to prevent file I/O
	am := &AuthManager{
		config: AuthConfig{
			Enabled:      true,
			JWTSecret:    "test-secret",
			InMemoryOnly: true,
		},
	}

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware for admin role
	adminMiddleware := am.RoleMiddleware(RoleAdmin)(testHandler)

	// Test with admin user in context
	adminUser := User{
		Username: "admin",
		Email:    "admin@example.com",
		Role:     RoleAdmin,
	}
	req := httptest.NewRequest("GET", "/", nil)
	ctx := WithUser(req.Context(), adminUser)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	adminMiddleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test with non-admin user in context
	regularUser := User{
		Username: "user",
		Email:    "user@example.com",
		Role:     RoleUser,
	}
	req = httptest.NewRequest("GET", "/", nil)
	ctx = WithUser(req.Context(), regularUser)
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()
	adminMiddleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Test with no user in context
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	adminMiddleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
