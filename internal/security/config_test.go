package security

import (
	"os"
	"testing"
	"time"
)

func TestSecurityConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
jwt:
  secret: "test-secret"
  expiration: 1h

tls:
  enabled: false # Set to false for testing to avoid loading certificates
  cert_file: "testdata/cert.pem"
  key_file: "testdata/key.pem"
  min_version: "TLS12"

rate_limit:
  enabled: true
  requests_per_minute: 60
  burst: 10

ip_allowlist:
  enabled: true
  allowed_ips:
    - "127.0.0.1"
    - "::1"

roles:
  admin:
    permissions:
      - "test:permission"
      - "test:admin"

security_headers:
  enabled: true
  hsts:
    enabled: true
    max_age: 31536000
    include_subdomains: true
    preload: true
  content_security_policy:
    enabled: true
    default_src: "'self'"
    script_src: "'self'"
    style_src: "'self'"
    img_src: "'self'"
    connect_src: "'self'"
  x_content_type_options: "nosniff"
  x_frame_options: "DENY"
  x_xss_protection: "1; mode=block"
  referrer_policy: "strict-origin-when-cross-origin"
`

	tmpFile := "testdata/temp_config.yaml"
	if err := os.MkdirAll("testdata", 0755); err != nil {
		t.Fatalf("Failed to create testdata directory: %v", err)
	}
	if err := os.WriteFile(tmpFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}
	defer os.RemoveAll("testdata")

	// Test loading configuration
	t.Run("Load Config", func(t *testing.T) {
		config, err := LoadConfigYAML(tmpFile)
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}

		// Verify JWT settings
		if config.JWT.Secret != "test-secret" {
			t.Error("JWT secret mismatch")
		}
		if config.JWT.Expiration != time.Hour {
			t.Error("JWT expiration mismatch")
		}

		// Verify TLS settings
		if config.TLS.Enabled {
			t.Error("TLS should NOT be enabled for testing")
		}
		if config.TLS.CertFile != "testdata/cert.pem" {
			t.Error("TLS cert file path mismatch")
		}

		// Verify rate limit settings
		if !config.RateLimit.Enabled {
			t.Error("Rate limit should be enabled")
		}
		if config.RateLimit.RequestsPerMinute != 60 {
			t.Error("Rate limit requests per minute mismatch")
		}

		// Verify IP allowlist
		if !config.IPAllowlist.Enabled {
			t.Error("IP allowlist should be enabled")
		}
		if len(config.IPAllowlist.AllowedIPs) != 2 {
			t.Error("IP allowlist count mismatch")
		}

		// Verify roles
		if len(config.Roles) != 1 {
			t.Error("Role count mismatch")
		}
		if len(config.Roles["admin"].Permissions) != 2 {
			t.Error("Admin role permissions count mismatch")
		}

		// Verify security headers
		if !config.SecurityHeaders.Enabled {
			t.Error("Security headers should be enabled")
		}
		if !config.SecurityHeaders.HSTS.Enabled {
			t.Error("HSTS should be enabled")
		}
		if !config.SecurityHeaders.ContentSecurityPolicy.Enabled {
			t.Error("Content Security Policy should be enabled")
		}
	})

	// Test applying configuration
	t.Run("Apply Config", func(t *testing.T) {
		config, _ := LoadConfigYAML(tmpFile)
		manager := NewSecurityManager("test-secret", time.Hour)

		// Apply configuration
		if err := manager.ApplyConfig(config); err != nil {
			t.Fatalf("Failed to apply config: %v", err)
		}

		// Verify role configuration
		user := &User{
			ID:       "test-user",
			Username: "test",
			Roles:    []string{"admin"},
		}
		if !manager.HasPermission(user, "test:permission") {
			t.Error("User should have test:permission")
		}
		if !manager.HasPermission(user, "test:admin") {
			t.Error("User should have test:admin")
		}

		// Verify IP allowlist
		if !manager.IsIPAllowed("127.0.0.1") {
			t.Error("127.0.0.1 should be allowed")
		}
		if !manager.IsIPAllowed("::1") {
			t.Error("::1 should be allowed")
		}
		if manager.IsIPAllowed("192.168.1.1") {
			t.Error("192.168.1.1 should not be allowed")
		}

		// Verify security headers
		headers := manager.GetSecurityHeaders()
		if headers["Strict-Transport-Security"] == "" {
			t.Error("HSTS header should be set")
		}
		if headers["Content-Security-Policy"] == "" {
			t.Error("Content Security Policy header should be set")
		}
		if headers["X-Content-Type-Options"] != "nosniff" {
			t.Error("X-Content-Type-Options header mismatch")
		}
	})
} 