package monitoring

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/nessi-dev/nessi/pkg/logging"
	"github.com/nessi-dev/nessi/pkg/security"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
}

// RefreshRequest represents a token refresh request
type RefreshRequest struct {
	Token string `json:"token"`
}

// UserResponse represents a user response
type UserResponse struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	LastLogin time.Time `json:"last_login,omitempty"`
}

// handleLogin handles login requests
func (m *Monitor) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if authentication is enabled
	if m.authManager == nil {
		http.Error(w, "Authentication is not enabled", http.StatusServiceUnavailable)
		return
	}

	// Parse request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate
	token, err := m.authManager.Authenticate(req.Username, req.Password)
	if err != nil {
		logging.Error("Authentication failed", err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Get user details
	user, err := m.authManager.GetUser(req.Username)
	if err != nil {
		logging.Error("Failed to get user details", err)
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	// Calculate expiry time
	claims, err := m.authManager.ValidateToken(token)
	if err != nil {
		logging.Error("Failed to validate token", err)
		http.Error(w, "Failed to validate token", http.StatusInternalServerError)
		return
	}

	// Get expiry time from claims
	var expiresAt time.Time
	if exp, ok := claims["exp"].(float64); ok {
		expiresAt = time.Unix(int64(exp), 0)
	} else {
		// Default to 24 hours if not found
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	// Return token
	resp := LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Username:  user.Username,
		Role:      user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleRefreshToken handles token refresh requests
func (m *Monitor) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if authentication is enabled
	if m.authManager == nil {
		http.Error(w, "Authentication is not enabled", http.StatusServiceUnavailable)
		return
	}

	// Parse request
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate token
	claims, err := m.authManager.ValidateToken(req.Token)
	if err != nil {
		logging.Error("Token validation failed", err)
		http.Error(w, "Token validation failed", http.StatusUnauthorized)
		return
	}

	// Get username from claims
	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Generate new token
	newToken, err := m.authManager.RefreshToken(username)
	if err != nil {
		logging.Error("Token refresh failed", err)
		http.Error(w, "Token refresh failed", http.StatusInternalServerError)
		return
	}

	// Get user details
	user, err := m.authManager.GetUser(username)
	if err != nil {
		logging.Error("Failed to get user details", err)
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	// Calculate expiry time
	newClaims, err := m.authManager.ValidateToken(newToken)
	if err != nil {
		logging.Error("Failed to validate new token", err)
		http.Error(w, "Failed to validate new token", http.StatusInternalServerError)
		return
	}

	// Get expiry time from claims
	var expiresAt time.Time
	if exp, ok := newClaims["exp"].(float64); ok {
		expiresAt = time.Unix(int64(exp), 0)
	} else {
		// Default to 24 hours if not found
		expiresAt = time.Now().Add(24 * time.Hour)
	}

	// Return new token
	resp := LoginResponse{
		Token:     newToken,
		ExpiresAt: expiresAt,
		Username:  user.Username,
		Role:      user.Role,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleUsers handles user management requests
func (m *Monitor) handleUsers(w http.ResponseWriter, r *http.Request) {
	// Check if authentication is enabled
	if m.authManager == nil {
		http.Error(w, "Authentication is not enabled", http.StatusServiceUnavailable)
		return
	}

	// Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		m.listUsers(w, r)
	case http.MethodPost:
		m.createUser(w, r)
	case http.MethodPut:
		m.updateUser(w, r)
	case http.MethodDelete:
		m.deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listUsers lists all users
func (m *Monitor) listUsers(w http.ResponseWriter, r *http.Request) {
	// Get users
	users, err := m.authManager.GetUsers()
	if err != nil {
		logging.Error("Failed to get users", err)
		http.Error(w, "Failed to get users", http.StatusInternalServerError)
		return
	}

	// Convert to response format
	var resp []UserResponse
	for _, user := range users {
		resp = append(resp, UserResponse{
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.DateCreated,
			LastLogin: user.LastLogin,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// createUser creates a new user
func (m *Monitor) createUser(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var user security.User
	var password string

	// Parse JSON request
	type CreateUserRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Set user fields
	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role
	password = req.Password

	// Default role if not provided
	if user.Role == "" {
		user.Role = security.RoleUser
	}

	// Create user
	if err := m.authManager.CreateUser(user, password); err != nil {
		logging.Error("Failed to create user", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Get created user
	createdUser, err := m.authManager.GetUser(user.Username)
	if err != nil {
		logging.Error("Failed to get created user", err)
		http.Error(w, "User created but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Return user
	resp := UserResponse{
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Role:      createdUser.Role,
		CreatedAt: createdUser.DateCreated,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// updateUser updates a user
func (m *Monitor) updateUser(w http.ResponseWriter, r *http.Request) {
	// Get username from query
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Parse request
	type UpdateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create updates map
	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}

	// Update user
	if err := m.authManager.UpdateUser(username, updates); err != nil {
		logging.Error("Failed to update user", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Get updated user
	updatedUser, err := m.authManager.GetUser(username)
	if err != nil {
		logging.Error("Failed to get updated user", err)
		http.Error(w, "User updated but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Return user
	resp := UserResponse{
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
		CreatedAt: updatedUser.DateCreated,
		LastLogin: updatedUser.LastLogin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// deleteUser deletes a user
func (m *Monitor) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Get username from query
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Delete user
	if err := m.authManager.DeleteUser(username); err != nil {
		logging.Error("Failed to delete user", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
