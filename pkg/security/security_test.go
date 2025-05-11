package security

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)
	assert.NotNil(t, s)
	assert.NotNil(t, s.tokenSecret)
	assert.Equal(t, tokenExpiration, s.tokenExpiration)
}

func TestAuthenticate(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)

	// Test with non-existent user
	token, err := s.Authenticate("nonexistent", "password")
	require.Error(t, err)
	assert.Empty(t, token)

	// Add test user
	err = s.AddUser("testuser", "password", []string{"user"})
	require.NoError(t, err)

	// Test with correct credentials
	token, err = s.Authenticate("testuser", "password")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test with incorrect password
	token, err = s.Authenticate("testuser", "wrongpassword")
	require.Error(t, err)
	assert.Empty(t, token)
}

func TestValidateToken(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)

	// Add test user
	err := s.AddUser("testuser", "password", []string{"user"})
	require.NoError(t, err)

	// Get valid token
	token, err := s.Authenticate("testuser", "password")
	require.NoError(t, err)

	// Validate token
	claims, err := s.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, "testuser", claims.Username)
	assert.Equal(t, []string{"user"}, claims.Roles)
}

func TestCreateAPIKey(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)

	// Add test user
	err := s.AddUser("testuser", "password", []string{"user"})
	require.NoError(t, err)

	// Create API key
	apiKey, err := s.CreateAPIKey("testuser")
	require.NoError(t, err)
	assert.NotEmpty(t, apiKey)

	// Validate API key
	username, err := s.ValidateAPIKey(apiKey)
	require.NoError(t, err)
	assert.Equal(t, "testuser", username)
}

func TestValidateAPIKey(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)

	// Add test user and create API key
	err := s.AddUser("testuser", "password", []string{"user"})
	require.NoError(t, err)

	apiKey, err := s.CreateAPIKey("testuser")
	require.NoError(t, err)

	// Validate existing API key
	username, err := s.ValidateAPIKey(apiKey)
	require.NoError(t, err)
	assert.Equal(t, "testuser", username)

	// Validate non-existent API key
	_, err = s.ValidateAPIKey("nonexistent-key")
	require.Error(t, err)
}

func TestAddUser(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)

	// Add user
	err := s.AddUser("testuser", "password", []string{"user"})
	require.NoError(t, err)

	// Verify user exists
	user, exists := s.users["testuser"]
	assert.True(t, exists)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, []string{"user"}, user.Roles)

	// Try adding same user again
	err = s.AddUser("testuser", "password", []string{"user"})
	require.Error(t, err)
}

func TestGenerateToken(t *testing.T) {
	tokenSecret := []byte("test-secret")
	tokenExpiration := 24 * time.Hour
	s := New(tokenSecret, tokenExpiration)

	// Add test user
	err := s.AddUser("testuser", "password", []string{"user"})
	require.NoError(t, err)

	// Get user
	user, exists := s.users["testuser"]
	require.True(t, exists)

	// Generate token
	token, err := s.generateToken(user)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Parse token
	parsedToken, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	require.NoError(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hashed := hashPassword(password)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)
}

func TestGenerateUserID(t *testing.T) {
	id1 := generateUserID()
	id2 := generateUserID()
	assert.NotEmpty(t, id1)
	assert.NotEqual(t, id1, id2)
}

func TestGenerateAPIKey(t *testing.T) {
	key1 := generateAPIKey()
	key2 := generateAPIKey()
	assert.NotEmpty(t, key1)
	assert.NotEqual(t, key1, key2)
}
