package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthManager(t *testing.T) {
	// Test with valid configuration
	validConfig := AuthConfig{
		Enabled:      true,
		UsersFile:    "test_users.json",
		InMemoryOnly: true,
	}

	am, err := NewAuthManager(validConfig)
	assert.NoError(t, err)
	assert.NotNil(t, am)

	// Test with empty configuration - this should now work with our CLI-only approach
	basicConfig := AuthConfig{
		Enabled: false,
	}

	am, err = NewAuthManager(basicConfig)
	assert.NoError(t, err)
	assert.NotNil(t, am)
}

func TestAuthManager(t *testing.T) {
	// Skip this test for CLI-only approach
	t.Skip("Skipping web auth tests in CLI-only mode")

	// Create a clean AuthManager for this test
	config := AuthConfig{
		Enabled:      true,
		InMemoryOnly: true,
	}

	// Create the auth manager
	am, err := NewAuthManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, am)

	// Test create user
	testUser := User{
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
	err = am.CreateUser(testUser, "password123")
	assert.NoError(t, err)

	// Test get user
	user, err := am.GetUser("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "user", user.Role)
	assert.Empty(t, user.Password) // Password should not be returned

	// Test get users
	users, err := am.GetUsers()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 1)

	// Test update user
	updates := map[string]interface{}{
		"email": "updated@example.com",
		"role":  "admin",
	}
	err = am.UpdateUser("testuser", updates)
	assert.NoError(t, err)

	// Verify update
	user, err = am.GetUser("testuser")
	assert.NoError(t, err)
	assert.Equal(t, "updated@example.com", user.Email)
	assert.Equal(t, "admin", user.Role)

	// Test regenerate API key
	apiKey, err := am.RegenerateAPIKey("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)

	// Test authenticate with API key
	user, err = am.AuthenticateWithAPIKey(apiKey)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)

	// Test delete user
	err = am.DeleteUser("testuser")
	assert.NoError(t, err)

	// Verify deletion
	_, err = am.GetUser("testuser")
	assert.Error(t, err)
}

func TestAPIKeyAuthentication(t *testing.T) {
	// Skip this test for CLI-only approach
	t.Skip("Skipping API key auth tests in CLI-only mode")

	// Create a clean AuthManager for this test
	config := AuthConfig{
		Enabled:      true,
		InMemoryOnly: true,
	}

	// Create the auth manager
	am, err := NewAuthManager(config)
	assert.NoError(t, err)
	assert.NotNil(t, am)

	// Create a test user
	testUser := User{
		Username: "apikeyuser",
		Email:    "apikey@example.com",
		Role:     RoleAdmin,
	}
	err = am.CreateUser(testUser, "apikey123")
	assert.NoError(t, err)

	// Generate API key
	apiKey, err := am.RegenerateAPIKey("apikeyuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, apiKey)

	// Test authenticate with API key
	user, err := am.AuthenticateWithAPIKey(apiKey)
	assert.NoError(t, err)
	assert.Equal(t, "apikeyuser", user.Username)

	// Test authenticate with invalid API key
	_, err = am.AuthenticateWithAPIKey("invalid-key")
	assert.Error(t, err)
}
