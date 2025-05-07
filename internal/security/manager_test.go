package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSecurityManager(t *testing.T) {
	// Create a new security manager
	manager := NewSecurityManager("test-secret", 1*time.Hour)

	// Test role management
	t.Run("Role Management", func(t *testing.T) {
		// Add a role
		manager.AddRole("test-role", []string{"test:permission"})

		// Create a test user
		user := &User{
			ID:       "test-user",
			Username: "test",
			Roles:    []string{"test-role"},
		}

		// Test permission check
		if !manager.HasPermission(user, "test:permission") {
			t.Error("User should have test:permission")
		}

		if manager.HasPermission(user, "invalid:permission") {
			t.Error("User should not have invalid:permission")
		}
	})

	// Test JWT token generation and validation
	t.Run("JWT Token", func(t *testing.T) {
		user := &User{
			ID:       "test-user",
			Username: "test",
			Roles:    []string{"test-role"},
		}

		// Generate token
		token, err := manager.GenerateToken(user)
		if err != nil {
			t.Fatalf("Failed to generate token: %v", err)
		}

		// Validate token
		validatedUser, err := manager.ValidateToken(token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if validatedUser.ID != user.ID {
			t.Error("User ID mismatch")
		}

		if validatedUser.Username != user.Username {
			t.Error("Username mismatch")
		}
	})

	// Test rate limiting
	t.Run("Rate Limiting", func(t *testing.T) {
		// Set rate limit
		manager.SetRateLimit("test-ip", 2, 1)

		// Test rate limit
		if !manager.AllowRequest("test-ip") {
			t.Error("First request should be allowed")
		}

		if !manager.AllowRequest("test-ip") {
			t.Error("Second request should be allowed")
		}

		if manager.AllowRequest("test-ip") {
			t.Error("Third request should be rate limited")
		}
	})

	// Test IP allowlist
	t.Run("IP Allowlist", func(t *testing.T) {
		// Add IP to allowlist
		manager.AddToAllowlist("192.168.1.1")

		// Test allowed IP
		if !manager.IsIPAllowed("192.168.1.1") {
			t.Error("IP should be allowed")
		}

		// Test disallowed IP
		if manager.IsIPAllowed("192.168.1.2") {
			t.Error("IP should not be allowed")
		}

		// Remove IP from allowlist
		manager.RemoveFromAllowlist("192.168.1.1")

		// Test removed IP
		if manager.IsIPAllowed("192.168.1.1") {
			t.Error("IP should not be allowed after removal")
		}
	})

	// Test middleware
	t.Run("Middleware", func(t *testing.T) {
		// Create test handler
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Create test request
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"

		// Test without token
		w := httptest.NewRecorder()
		manager.Middleware(handler).ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Error("Should return unauthorized without token")
		}

		// Add IP to allowlist
		manager.AddToAllowlist("192.168.1.1")

		// Generate token
		user := &User{
			ID:       "test-user",
			Username: "test",
			Roles:    []string{"test-role"},
		}
		token, _ := manager.GenerateToken(user)

		// Test with token
		req.Header.Set("Authorization", token)
		w = httptest.NewRecorder()
		manager.Middleware(handler).ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Error("Should return OK with valid token")
		}
	})
}

func TestTLSConfig(t *testing.T) {
	manager := NewSecurityManager("test-secret", 1*time.Hour)

	// Test TLS configuration
	err := manager.ConfigureTLS("testdata/cert.pem", "testdata/key.pem")
	if err == nil {
		t.Error("Should fail with non-existent certificate files")
	}

	// Note: To test with real certificates, you would need to generate them first
	// This is just a basic test to ensure the function handles errors properly
} 