package security

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SecurityManager handles authentication and authorization
type SecurityManager struct {
	mu              sync.Mutex
	users           map[string]*User
	apiKeys         map[string]string
	tokenSecret     []byte
	tokenExpiration time.Duration
}

// User represents a system user
type User struct {
	ID       string
	Username string
	Password string
	Roles    []string
	APIKeys  []string
}

// TokenClaims represents JWT claims
//go:generate stringer -type=Role
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleUser     Role = "user"
	RoleReadOnly Role = "readonly"
)

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// New creates a new SecurityManager instance
func New(tokenSecret []byte, tokenExpiration time.Duration) *SecurityManager {
	return &SecurityManager{
		users:           make(map[string]*User),
		apiKeys:         make(map[string]string),
		tokenSecret:     tokenSecret,
		tokenExpiration: tokenExpiration,
	}
}

// Authenticate authenticates a user
func (s *SecurityManager) Authenticate(username, password string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[username]
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	// Verify password (in production, use a proper hashing algorithm)
	if user.Password != hashPassword(password) {
		return "", fmt.Errorf("invalid password")
	}

	return s.generateToken(user)
}

// ValidateToken validates a JWT token
func (s *SecurityManager) ValidateToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.tokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// CreateAPIKey creates a new API key for a user
func (s *SecurityManager) CreateAPIKey(username string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[username]
	if !exists {
		return "", fmt.Errorf("user not found")
	}

	apiKey := generateAPIKey()
	s.apiKeys[apiKey] = username
	user.APIKeys = append(user.APIKeys, apiKey)

	return apiKey, nil
}

// ValidateAPIKey validates an API key
func (s *SecurityManager) ValidateAPIKey(apiKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	username, exists := s.apiKeys[apiKey]
	if !exists {
		return "", fmt.Errorf("invalid api key")
	}

	return username, nil
}

// AddUser adds a new user
func (s *SecurityManager) AddUser(username, password string, roles []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[username]; exists {
		return fmt.Errorf("user already exists")
	}

	s.users[username] = &User{
		ID:       generateUserID(),
		Username: username,
		Password: hashPassword(password),
		Roles:    roles,
	}

	return nil
}

// generateToken generates a JWT token
func (s *SecurityManager) generateToken(user *User) (string, error) {
	claims := &TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.tokenSecret)
}

// hashPassword hashes a password using SHA256
func hashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// generateUserID generates a unique user ID
func generateUserID() string {
	return fmt.Sprintf("user_%x", sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
}

// generateAPIKey generates a unique API key
func generateAPIKey() string {
	return fmt.Sprintf("key_%x", sha256.Sum256([]byte(fmt.Sprintf("%d", time.Now().UnixNano()))))
}
