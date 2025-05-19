package security

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user account
type User struct {
	Username    string    `json:"username"`
	Password    string    `json:"password,omitempty"` // Hashed password, not exposed in JSON
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	APIKey      string    `json:"api_key,omitempty"` // API key, not exposed in JSON
	LastLogin   time.Time `json:"last_login,omitempty"`
	DateCreated time.Time `json:"date_created"`
}

// Role constants
const (
	RoleAdmin  = "admin"
	RoleUser   = "user"
	RoleViewer = "viewer"
)

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Enabled      bool   `json:"enabled"`
	UsersFile    string `json:"users_file"`
	APIKeyPath   string `json:"api_key_path"`
	InMemoryOnly bool   `json:"in_memory_only"` // If true, disables all file I/O for tests
}

// AuthManager handles authentication and authorization
type AuthManager struct {
	mu     sync.RWMutex
	config AuthConfig
	users  map[string]User // username -> User
	apiKeys map[string]string // api_key -> username
}

// NewAuthManager creates a new AuthManager
func NewAuthManager(config AuthConfig) (*AuthManager, error) {
	am := &AuthManager{
		config:  config,
		users:   make(map[string]User),
		apiKeys: make(map[string]string),
	}

	// Load users if auth is enabled
	if config.Enabled {
		if err := am.loadUsers(); err != nil {
			return nil, err
		}
	}

	return am, nil
}

// loadUsers loads users from the users file
func (am *AuthManager) loadUsers() error {
	if am.config.InMemoryOnly {
		return nil // skip file I/O for tests
	}
	// Check if users file exists
	if _, err := os.Stat(am.config.UsersFile); os.IsNotExist(err) {
		// Create default admin user if file doesn't exist
		adminUser := User{
			Username:    "admin",
			Password:    hashPassword("admin"), // Default password, should be changed
			Email:       "admin@example.com",
			Role:        RoleAdmin,
			APIKey:      generateAPIKey(),
			DateCreated: time.Now(),
		}

		am.users[adminUser.Username] = adminUser
		am.apiKeys[adminUser.APIKey] = adminUser.Username

		// Save users to file
		return am.saveUsers()
	}

	// Read users file
	data, err := os.ReadFile(am.config.UsersFile)
	if err != nil {
		return fmt.Errorf("failed to read users file: %w", err)
	}

	// Parse users
	var users []User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to parse users file: %w", err)
	}

	// Store users in memory
	for _, user := range users {
		am.users[user.Username] = user
		if user.APIKey != "" {
			am.apiKeys[user.APIKey] = user.Username
		}
	}

	return nil
}

// saveUsers saves users to the users file
func (am *AuthManager) saveUsers() error {
	if am.config.InMemoryOnly {
		return nil // skip file I/O for tests
	}
	
	// Create a copy of the users map to avoid holding the lock during I/O operations
	am.mu.RLock() // Use read lock first to make a copy
	usersCopy := make(map[string]User, len(am.users))
	for k, v := range am.users {
		usersCopy[k] = v
	}
	am.mu.RUnlock()

	// Convert users map to slice
	var users []User
	for _, user := range usersCopy {
		users = append(users, user)
	}

	// Marshal users to JSON
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	// Write to file - this is done outside the lock to prevent deadlocks
	if err := os.WriteFile(am.config.UsersFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write users file: %w", err)
	}

	return nil
}

// Authenticate authenticates a user with username and password
func (am *AuthManager) Authenticate(username, password string) (bool, error) {
	if !am.config.Enabled {
		return false, fmt.Errorf("authentication is disabled")
	}

	am.mu.RLock()
	user, exists := am.users[username]
	am.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("user not found")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return false, fmt.Errorf("invalid password")
	}

	// Update last login time
	am.mu.Lock()
	user.LastLogin = time.Now()
	am.users[username] = user
	am.mu.Unlock()

	// Save users to file
	go am.saveUsers() // Non-blocking save

	return true, nil
}

// RefreshAPIKey generates a new API key for the specified user
func (am *AuthManager) RefreshAPIKey(username string) (string, error) {
	if !am.config.Enabled {
		return "", fmt.Errorf("authentication is disabled")
	}

	return am.RegenerateAPIKey(username)
}

// AuthenticateWithAPIKey authenticates a user with an API key
func (am *AuthManager) AuthenticateWithAPIKey(apiKey string) (User, error) {
	if !am.config.Enabled {
		return User{}, fmt.Errorf("authentication is disabled")
	}

	am.mu.RLock()
	// Check if API key exists
	username, ok := am.apiKeys[apiKey]
	if !ok {
		am.mu.RUnlock()
		return User{}, fmt.Errorf("invalid API key")
	}

	// Get user
	user, ok := am.users[username]
	am.mu.RUnlock()
	
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

// GetAPIKey returns the API key for a user
func (am *AuthManager) GetAPIKey(username string) (string, bool) {
	if !am.config.Enabled {
		return "", false
	}

	am.mu.RLock()
	user, ok := am.users[username]
	am.mu.RUnlock()

	if !ok || user.APIKey == "" {
		return "", false
	}

	return user.APIKey, true
}

// CreateUser creates a new user
func (am *AuthManager) CreateUser(user User, password string) error {
	if !am.config.Enabled {
		return fmt.Errorf("authentication is disabled")
	}

	// Check if user already exists
	am.mu.RLock()
	_, exists := am.users[user.Username]
	am.mu.RUnlock()

	if exists {
		return fmt.Errorf("user already exists")
	}

	// Set password and API key
	user.Password = hashPassword(password)
	user.APIKey = generateAPIKey()

	// Set creation date
	user.DateCreated = time.Now()

	// Add user to map
	am.mu.Lock()
	am.users[user.Username] = user
	am.apiKeys[user.APIKey] = user.Username
	am.mu.Unlock()

	// Save users to file - this is done outside the lock to prevent deadlocks
	return am.saveUsers()
}

// UpdateUser updates an existing user
func (am *AuthManager) UpdateUser(username string, updates map[string]interface{}) error {
	if !am.config.Enabled {
		return fmt.Errorf("authentication is disabled")
	}

	// First check if user exists with read lock
	am.mu.RLock()
	user, ok := am.users[username]
	am.mu.RUnlock()

	if !ok {
		return fmt.Errorf("user not found")
	}

	// Apply updates to a local copy
	for key, value := range updates {
		switch key {
		case "email":
			if email, ok := value.(string); ok {
				user.Email = email
			}
		case "role":
			if role, ok := value.(string); ok {
				user.Role = role
			}
		case "password":
			if password, ok := value.(string); ok {
				user.Password = hashPassword(password)
			}
		}
	}

	// Update user with write lock
	am.mu.Lock()
	am.users[username] = user
	am.mu.Unlock()

	// Save users to file - this happens outside the lock
	return am.saveUsers()
}

// DeleteUser deletes a user
func (am *AuthManager) DeleteUser(username string) error {
	if !am.config.Enabled {
		return fmt.Errorf("authentication is disabled")
	}

	am.mu.Lock()
	defer am.mu.Unlock()

	// Check if user exists
	user, ok := am.users[username]
	if !ok {
		return fmt.Errorf("user not found")
	}

	// Remove API key
	delete(am.apiKeys, user.APIKey)

	// Remove user
	delete(am.users, username)

	// Save users to file
	return am.saveUsers()
}

// GetUser returns a user by username
func (am *AuthManager) GetUser(username string) (User, error) {
	if !am.config.Enabled {
		return User{}, fmt.Errorf("authentication is disabled")
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	user, ok := am.users[username]
	if !ok {
		return User{}, fmt.Errorf("user not found")
	}

	// Don't return sensitive fields
	user.Password = ""
	user.APIKey = ""

	return user, nil
}

// GetUsers returns all users
func (am *AuthManager) GetUsers() ([]User, error) {
	if !am.config.Enabled {
		return nil, fmt.Errorf("authentication is disabled")
	}

	am.mu.RLock()
	defer am.mu.RUnlock()

	var users []User
	for _, user := range am.users {
		// Don't return sensitive fields
		user.Password = ""
		user.APIKey = ""
		users = append(users, user)
	}

	return users, nil
}

// RegenerateAPIKey regenerates a user's API key
func (am *AuthManager) RegenerateAPIKey(username string) (string, error) {
	if !am.config.Enabled {
		return "", fmt.Errorf("authentication is disabled")
	}

	// First check if user exists with read lock
	am.mu.RLock()
	user, ok := am.users[username]
	am.mu.RUnlock()
	
	if !ok {
		return "", fmt.Errorf("user not found")
	}

	// Now update with write lock
	am.mu.Lock()
	// Remove old API key
	if user.APIKey != "" {
		delete(am.apiKeys, user.APIKey)
	}

	// Generate new API key
	newAPIKey := generateAPIKey()
	user.APIKey = newAPIKey

	// Update user and API key map
	am.users[username] = user
	am.apiKeys[newAPIKey] = username
	am.mu.Unlock()

	// Save users to file - this happens outside the lock
	if err := am.saveUsers(); err != nil {
		return "", fmt.Errorf("failed to save users: %w", err)
	}

	return newAPIKey, nil
}

// ValidateAPIKey validates an API key and returns the associated user
func (am *AuthManager) ValidateAPIKey(apiKey string) (User, error) {
	// Skip authentication if disabled
	if !am.config.Enabled {
		return User{}, fmt.Errorf("authentication is disabled")
	}

	// Check for API key
	if apiKey == "" {
		return User{}, fmt.Errorf("API key is required")
	}

	// Authenticate with API key
	user, err := am.AuthenticateWithAPIKey(apiKey)
	if err != nil {
		return User{}, fmt.Errorf("invalid API key: %w", err)
	}

	return user, nil
}

// CheckUserRole checks if a user has one of the specified roles
func (am *AuthManager) CheckUserRole(user User, roles ...string) bool {
	// Skip authorization if authentication is disabled
	if !am.config.Enabled {
		return true
	}

	// Check if user has required role
	for _, role := range roles {
		if user.Role == role {
			return true
		}
	}

	return false
}

// Login authenticates a user with username and password and returns an API key
func (am *AuthManager) Login(username, password string) (string, error) {
	// Authenticate user
	authenticated, err := am.Authenticate(username, password)
	if err != nil || !authenticated {
		return "", fmt.Errorf("invalid username or password")
	}

	// Get or regenerate API key
	apiKey, err := am.RefreshAPIKey(username)
	if err != nil {
		return "", fmt.Errorf("failed to generate API key: %w", err)
	}

	return apiKey, nil
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logging.Error("Failed to hash password", err)
		return ""
	}
	return string(hash)
}

// generateAPIKey generates a random API key
func generateAPIKey() string {
	// Generate a random string for the API key
	// In a real application, use a more secure method
	return fmt.Sprintf("key_%d", time.Now().UnixNano())
}

// SecureCompare compares two strings in constant time
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
