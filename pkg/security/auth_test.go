package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthManager(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "users-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Initialize the file with empty JSON array
	_, err = tempFile.WriteString("[]")
	require.NoError(t, err)
	tempFile.Close()

	// Create auth config
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create auth manager
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
	assert.Len(t, users, 2) // testuser + default admin

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
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "users-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Initialize the file with empty JSON array
	_, err = tempFile.WriteString("[]")
	require.NoError(t, err)
	tempFile.Close()

	// Create auth config
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create auth manager
	am, err := NewAuthManager(config)
	require.NoError(t, err)
	require.NotNil(t, am)

	// Create test user
	testUser := User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     RoleUser,
	}
	err = am.CreateUser(testUser, "password123")
	require.NoError(t, err)

	// Get token
	token, err := am.Authenticate("testuser", "password123")
	require.NoError(t, err)

	// Get API key
	apiKey, err := am.RegenerateAPIKey("testuser")
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
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "testuser", w.Body.String())

	// Test with API key
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", apiKey)
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
}

func TestRoleMiddleware(t *testing.T) {
	// Create temporary users file
	tempFile, err := os.CreateTemp("", "users-*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	
	// Initialize the file with empty JSON array
	_, err = tempFile.WriteString("[]")
	require.NoError(t, err)
	tempFile.Close()

	// Create auth config
	config := AuthConfig{
		Enabled:      true,
		JWTSecret:    "test-secret",
		UsersFile:    tempFile.Name(),
		TokenExpiry:  24,
		RequireHTTPS: false,
	}

	// Create auth manager
	am, err := NewAuthManager(config)
	require.NoError(t, err)
	require.NotNil(t, am)

	// Create test users
	adminUser := User{
		Username: "admin",
		Email:    "admin@example.com",
		Role:     RoleAdmin,
	}
	err = am.CreateUser(adminUser, "admin123")
	require.NoError(t, err)

	regularUser := User{
		Username: "user",
		Email:    "user@example.com",
		Role:     RoleUser,
	}
	err = am.CreateUser(regularUser, "user123")
	require.NoError(t, err)

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
}

func TestContext(t *testing.T) {
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
