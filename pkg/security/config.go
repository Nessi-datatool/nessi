package security

// SecurityConfig defines the security configuration
type SecurityConfig struct {
	Auth AuthConfig `json:"auth"`
	SSL  SSLConfig  `json:"ssl"`
}
