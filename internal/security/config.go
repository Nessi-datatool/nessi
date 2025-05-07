package security

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the security configuration
type Config struct {
	// JWT settings
	JWTSecret     string        `json:"jwt_secret"`
	JWTExpiration time.Duration `json:"jwt_expiration"`

	// TLS settings
	TLSCertFile string `json:"tls_cert_file"`
	TLSKeyFile  string `json:"tls_key_file"`

	// Rate limiting
	DefaultRateLimit  int `json:"default_rate_limit"`
	DefaultBurstLimit int `json:"default_burst_limit"`

	// IP allowlist
	AllowedIPs []string `json:"allowed_ips"`

	// Role definitions
	Roles []RoleConfig `json:"roles"`
}

// RoleConfig defines a role and its permissions
type RoleConfig struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// SecurityConfig represents the security configuration
type SecurityConfig struct {
	JWT struct {
		Secret     string        `yaml:"secret"`
		Expiration time.Duration `yaml:"expiration"`
	} `yaml:"jwt"`

	TLS struct {
		Enabled    bool   `yaml:"enabled"`
		CertFile   string `yaml:"cert_file"`
		KeyFile    string `yaml:"key_file"`
		MinVersion string `yaml:"min_version"`
	} `yaml:"tls"`

	RateLimit struct {
		Enabled           bool `yaml:"enabled"`
		RequestsPerMinute int  `yaml:"requests_per_minute"`
		Burst            int  `yaml:"burst"`
	} `yaml:"rate_limit"`

	IPAllowlist struct {
		Enabled    bool     `yaml:"enabled"`
		AllowedIPs []string `yaml:"allowed_ips"`
	} `yaml:"ip_allowlist"`

	Roles map[string]struct {
		Permissions []string `yaml:"permissions"`
	} `yaml:"roles"`

	DefaultRole string `yaml:"default_role"`

	SecurityHeaders struct {
		Enabled bool `yaml:"enabled"`
		HSTS   struct {
			Enabled           bool `yaml:"enabled"`
			MaxAge           int  `yaml:"max_age"`
			IncludeSubdomains bool `yaml:"include_subdomains"`
			Preload          bool `yaml:"preload"`
		} `yaml:"hsts"`
		ContentSecurityPolicy struct {
			Enabled     bool   `yaml:"enabled"`
			DefaultSrc  string `yaml:"default_src"`
			ScriptSrc   string `yaml:"script_src"`
			StyleSrc    string `yaml:"style_src"`
			ImgSrc      string `yaml:"img_src"`
			ConnectSrc  string `yaml:"connect_src"`
		} `yaml:"content_security_policy"`
		XContentTypeOptions string `yaml:"x_content_type_options"`
		XFrameOptions      string `yaml:"x_frame_options"`
		XXSSProtection     string `yaml:"x_xss_protection"`
		ReferrerPolicy     string `yaml:"referrer_policy"`
	} `yaml:"security_headers"`
}

// LoadConfig loads security configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadConfig loads the security configuration from a YAML file
func LoadConfigYAML(path string) (*SecurityConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config SecurityConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// DefaultConfig returns default security configuration
func DefaultConfig() *Config {
	return &Config{
		JWTSecret:         "change-me-in-production",
		JWTExpiration:     24 * time.Hour,
		DefaultRateLimit:  100,
		DefaultBurstLimit: 200,
		Roles: []RoleConfig{
			{
				Name: "admin",
				Permissions: []string{
					"read:*",
					"write:*",
					"delete:*",
					"manage:*",
				},
			},
			{
				Name: "user",
				Permissions: []string{
					"read:*",
					"write:own",
				},
			},
			{
				Name: "viewer",
				Permissions: []string{
					"read:*",
				},
			},
		},
	}
}

// SaveConfig saves security configuration to a file
func (c *Config) SaveConfig(path string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ApplyConfig applies the security configuration to the security manager
func (m *SecurityManager) ApplyConfig(config *SecurityConfig) error {
	// Configure JWT
	m.jwtSecret = []byte(config.JWT.Secret)
	m.jwtExpiration = config.JWT.Expiration

	// Configure TLS if enabled
	if config.TLS.Enabled {
		if err := m.ConfigureTLS(config.TLS.CertFile, config.TLS.KeyFile); err != nil {
			return fmt.Errorf("failed to configure TLS: %w", err)
		}
	}

	// Configure rate limiting
	if config.RateLimit.Enabled {
		m.SetRateLimit("default", config.RateLimit.RequestsPerMinute, config.RateLimit.Burst)
	}

	// Configure IP allowlist
	if config.IPAllowlist.Enabled {
		for _, ip := range config.IPAllowlist.AllowedIPs {
			m.AddToAllowlist(ip)
		}
	}

	// Configure roles
	for role, roleConfig := range config.Roles {
		m.AddRole(role, roleConfig.Permissions)
	}

	return nil
}

// GetSecurityHeaders returns the configured security headers
func (m *SecurityManager) GetSecurityHeaders() map[string]string {
	headers := make(map[string]string)

	if m.config.SecurityHeaders.Enabled {
		// HSTS
		if m.config.SecurityHeaders.HSTS.Enabled {
			hsts := fmt.Sprintf("max-age=%d", m.config.SecurityHeaders.HSTS.MaxAge)
			if m.config.SecurityHeaders.HSTS.IncludeSubdomains {
				hsts += "; includeSubDomains"
			}
			if m.config.SecurityHeaders.HSTS.Preload {
				hsts += "; preload"
			}
			headers["Strict-Transport-Security"] = hsts
		}

		// Content Security Policy
		if m.config.SecurityHeaders.ContentSecurityPolicy.Enabled {
			csp := fmt.Sprintf("default-src %s; script-src %s; style-src %s; img-src %s; connect-src %s",
				m.config.SecurityHeaders.ContentSecurityPolicy.DefaultSrc,
				m.config.SecurityHeaders.ContentSecurityPolicy.ScriptSrc,
				m.config.SecurityHeaders.ContentSecurityPolicy.StyleSrc,
				m.config.SecurityHeaders.ContentSecurityPolicy.ImgSrc,
				m.config.SecurityHeaders.ContentSecurityPolicy.ConnectSrc,
			)
			headers["Content-Security-Policy"] = csp
		}

		// Other security headers
		headers["X-Content-Type-Options"] = m.config.SecurityHeaders.XContentTypeOptions
		headers["X-Frame-Options"] = m.config.SecurityHeaders.XFrameOptions
		headers["X-XSS-Protection"] = m.config.SecurityHeaders.XXSSProtection
		headers["Referrer-Policy"] = m.config.SecurityHeaders.ReferrerPolicy
	}

	return headers
} 