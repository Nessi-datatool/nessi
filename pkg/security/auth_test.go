package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthManager(t *testing.T) {
	// Run in parallel for faster execution
	t.Parallel()

	// Create a clean AuthManager for this test
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-secret-" + t.Name(),
		UsersFile:    "/tmp/nonexistent-users-file.json",
		TokenExpiry:  1, // 1 hour - minimum value for faster tests
		RequireHTTPS: false,
		InMemoryOnly: true, // Skip all file I/O for better performance
	}

	// Create the auth manager
	am, err := NewAuthManager(config)
	require.NoError(t, err)
	require.NotNil(t, am)

	// Test create user
	testUser := User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     RoleUser,
	}
	err = am.CreateUser(testUser, "password123")
	require.NoError(t, err)

	// Test authenticate
	token, err := am.Authenticate("testuser", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test authenticate with wrong password
	_, err = am.Authenticate("testuser", "wrongpassword")
	assert.Error(t, err)

	// Test get user
	user, err := am.GetUser("testuser")
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, RoleUser, user.Role)
	assert.Empty(t, user.Password) // Password should not be returned

	// Test get users
	users, err := am.GetUsers()
	require.NoError(t, err)
	// With in-memory storage, we only have the user we created
	assert.Len(t, users, 1) // just testuser, no default admin in in-memory mode

	// Test update user
	updates := map[string]interface{}{
		"email": "updated@example.com",
		"role":  RoleAdmin,
	}
	err = am.UpdateUser("testuser", updates)
	require.NoError(t, err)

	// Verify update
	user, err = am.GetUser("testuser")
	require.NoError(t, err)
	assert.Equal(t, "updated@example.com", user.Email)
	assert.Equal(t, RoleAdmin, user.Role)

	// Test regenerate API key
	apiKey, err := am.RegenerateAPIKey("testuser")
	require.NoError(t, err)
	assert.NotEmpty(t, apiKey)

	// Test authenticate with API key
	user, err = am.AuthenticateWithAPIKey(apiKey)
	require.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)

	// Test validate token
	claims, err := am.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "testuser", claims["username"])

	// Test delete user
	err = am.DeleteUser("testuser")
	require.NoError(t, err)

	// Verify deletion
	_, err = am.GetUser("testuser")
	assert.Error(t, err)
}

func TestAuthMiddleware(t *testing.T) {
	// Run in parallel for faster execution
	t.Parallel()
	
	// Run with timeout to prevent hanging tests
	RunWithTimeout(t, func() {

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
		"username":  "testuser",
		"email":     "test@example.com",
		"role":      RoleUser,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"created_at": time.Now().Unix(),
	}
	
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// Sign the token
	tokenString, err := token.SignedString([]byte("test-secret"))
	require.NoError(t, err)

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok {
			http.Error(w, "User not found in context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.Username))
	})

	// Create middleware
	middleware := am.AuthMiddleware(testHandler)

	// Test with token
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "testuser", w.Body.String())

	// Test with API key
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "test-api-key")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "testuser", w.Body.String())

	// Test with invalid token
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test with invalid API key
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "invalid-key")
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Test with no authentication
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRoleMiddleware(t *testing.T) {
	// Run in parallel for faster execution
	t.Parallel()
	
	// Run with timeout to prevent hanging tests
	RunWithTimeout(t, func() {
	
	// Use the optimized test helper to create an AuthManager
	am, err := CreateTestAuthManager()
	require.NoError(t, err)

	// Create test users directly
	adminUser := User{
		Username: "admin",
		Email:    "admin@example.com",
		Role:     RoleAdmin,
	}

	regularUser := User{
		Username: "user",
		Email:    "user@example.com",
		Role:     RoleUser,
	}

	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create middleware
	adminMiddleware := am.RoleMiddleware(RoleAdmin)(testHandler)

	// Test with admin user
	req := httptest.NewRequest("GET", "/", nil)
	ctx := WithUser(req.Context(), adminUser)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()
	adminMiddleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test with regular user
	req = httptest.NewRequest("GET", "/", nil)
	ctx = WithUser(req.Context(), regularUser)
	req = req.WithContext(ctx)
	w = httptest.NewRecorder()
	adminMiddleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	// Test with no user
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	adminMiddleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestContext(t *testing.T) {
	// Run in parallel for faster execution
	t.Parallel()
	// Create test user
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     RoleUser,
	}

	// Test WithUser and UserFromContext
	ctx := context.Background()
	ctx = WithUser(ctx, user)
	retrievedUser, ok := UserFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, user, retrievedUser)

	// Test IsAuthenticated
	assert.True(t, IsAuthenticated(ctx))
	assert.False(t, IsAuthenticated(context.Background()))

	// Test HasRole
	assert.True(t, HasRole(ctx, RoleUser))
	assert.False(t, HasRole(ctx, RoleAdmin))
	assert.False(t, HasRole(context.Background(), RoleUser))

	// Test IsAdmin
	assert.False(t, IsAdmin(ctx))
	
	// Create admin user
	adminUser := User{
		Username: "admin",
		Email:    "admin@example.com",
		Role:     RoleAdmin,
	}
	
	// Test IsAdmin with admin user
	adminCtx := WithUser(context.Background(), adminUser)
	assert.True(t, IsAdmin(adminCtx))
}
