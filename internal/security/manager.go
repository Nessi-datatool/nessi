package security

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// Role represents a user role with permissions
type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// User represents an authenticated user
type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	APIKey   string   `json:"api_key,omitempty"`
}

// Claims represents the JWT claims for a user
type Claims struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// SecurityManager handles authentication, authorization, and rate limiting
type SecurityManager struct {
	// JWT configuration
	jwtSecret     string
	jwtExpiration time.Duration

	// Role definitions
	roles  map[string]*Role
	roleMu sync.RWMutex

	// IP allowlist
	allowlist        map[string]bool
	allowlistEnabled bool
	allowMu          sync.RWMutex

	// Rate limiting
	rateLimiters map[string]*rate.Limiter
	rateMu       sync.RWMutex

	// TLS configuration
	tlsConfig *tls.Config

	// Applied configuration
	config *SecurityConfig
}

// NewSecurityManager creates a new SecurityManager instance
func NewSecurityManager(secret string, expiration time.Duration) *SecurityManager {
	return &SecurityManager{
		jwtSecret:        secret,
		jwtExpiration:    expiration,
		roles:            make(map[string]*Role),
		allowlist:        make(map[string]bool),
		allowlistEnabled: false, // Start with allowlist disabled
		rateLimiters:     make(map[string]*rate.Limiter),
	}
}

// ConfigureTLS sets up TLS configuration
func (m *SecurityManager) ConfigureTLS(certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	m.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}

	return nil
}

// AddRole adds a new role with permissions
func (m *SecurityManager) AddRole(name string, permissions []string) {
	m.roleMu.Lock()
	defer m.roleMu.Unlock()
	m.roles[name] = &Role{
		Name:        name,
		Permissions: permissions,
	}
}

// HasPermission checks if a user has a specific permission
func (m *SecurityManager) HasPermission(user *User, permission string) bool {
	m.roleMu.RLock()
	defer m.roleMu.RUnlock()

	for _, roleName := range user.Roles {
		if role, exists := m.roles[roleName]; exists {
			for _, p := range role.Permissions {
				if p == permission {
					return true
				}
			}
		}
	}
	return false
}

// GenerateToken generates a JWT token for a user
func (m *SecurityManager) GenerateToken(user *User) (string, error) {
	claims := &Claims{
		ID:       user.ID,
		Username: user.Username,
		Roles:    user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.jwtExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.jwtSecret))
	return tokenString, err
}

// ValidateToken validates a JWT token and returns the user
func (m *SecurityManager) ValidateToken(tokenString string) (*User, error) {
	claims := &Claims{}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if token.Valid {
		return &User{
			ID:       claims.ID,
			Username: claims.Username,
			Roles:    claims.Roles,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// SetRateLimit sets the rate limit for a key
// r is requests per second as rate.Limit, b is burst size
func (m *SecurityManager) SetRateLimit(key string, r rate.Limit, b int) {
	m.rateMu.Lock()
	defer m.rateMu.Unlock()

	// Create a new limiter with the specified rate and burst
	m.rateLimiters[key] = rate.NewLimiter(r, b)
}

// AllowRequest checks if a request should be allowed based on rate limiting
func (m *SecurityManager) AllowRequest(key string) bool {
	m.rateMu.RLock()
	limiter, exists := m.rateLimiters[key]
	m.rateMu.RUnlock()

	if !exists {
		return true
	}

	return limiter.Allow()
}

// AddToAllowlist adds an IP to the allowlist
func (m *SecurityManager) AddToAllowlist(ip string) {
	m.allowMu.Lock()
	defer m.allowMu.Unlock()
	// Enable allowlist mode when first IP is added
	m.allowlistEnabled = true
	m.allowlist[ip] = true
}

// RemoveFromAllowlist removes an IP from the allowlist
func (m *SecurityManager) RemoveFromAllowlist(ip string) {
	m.allowMu.Lock()
	defer m.allowMu.Unlock()
	delete(m.allowlist, ip)

	// Check if this was the last IP - if we still want to enforce the allowlist
	// keep allowlistEnabled true even if empty
}

// IsIPAllowed checks if an IP is allowed
func (m *SecurityManager) IsIPAllowed(ip string) bool {
	m.allowMu.RLock()
	defer m.allowMu.RUnlock()

	// If allowlist is not enabled, allow all IPs
	if !m.allowlistEnabled {
		return true
	}

	// Check if the IP is explicitly allowed in the allowlist
	allowed, exists := m.allowlist[ip]
	// If the IP exists in the map and is set to true, it's allowed
	// If it doesn't exist in the map, it's not allowed when allowlist is enabled
	return exists && allowed
}

// GetTLSConfig returns the TLS configuration
func (m *SecurityManager) GetTLSConfig() *tls.Config {
	return m.tlsConfig
}

// Middleware returns a middleware function for authentication and rate limiting
func (m *SecurityManager) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check IP allowlist
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		if !m.IsIPAllowed(ip) {
			http.Error(w, "IP not allowed", http.StatusForbidden)
			return
		}

		// Check rate limit
		if !m.AllowRequest(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Check authentication
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
