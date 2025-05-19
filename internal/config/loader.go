package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Security    SecurityConfig    `yaml:"security"`
	Delta       DeltaConfig       `yaml:"delta"`
	Quality     QualityConfig     `yaml:"quality"`
	Monitoring  MonitoringConfig  `yaml:"monitoring"`
	Reports     ReportsConfig     `yaml:"reports"`
	Logging     LoggingConfig     `yaml:"logging"`
	Development DevelopmentConfig `yaml:"development"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	MaxConnections  int           `yaml:"max_connections"`
	CORS            CORSConfig    `yaml:"cors"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	ExposedHeaders []string `yaml:"exposed_headers"`
	MaxAge         int      `yaml:"max_age"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	JWT         JWTConfig             `yaml:"jwt"`
	TLS         TLSConfig             `yaml:"tls"`
	RateLimit   RateLimitConfig       `yaml:"rate_limit"`
	IPAllowlist IPAllowlistConfig     `yaml:"ip_allowlist"`
	Roles       map[string]RoleConfig `yaml:"roles"`
	DefaultRole string                `yaml:"default_role"`
	Headers     SecurityHeadersConfig `yaml:"headers"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	MinVersion string `yaml:"min_version"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
	Burst             int `yaml:"burst"`
}

// IPAllowlistConfig represents IP allowlist configuration
type IPAllowlistConfig struct {
	Enabled bool     `yaml:"enabled"`
	IPs     []string `yaml:"ips"`
}

// RoleConfig represents role configuration
type RoleConfig struct {
	Permissions []string `yaml:"permissions"`
}

// SecurityHeadersConfig represents security headers configuration
type SecurityHeadersConfig struct {
	HSTS                  HSTSConfig                  `yaml:"hsts"`
	ContentSecurityPolicy ContentSecurityPolicyConfig `yaml:"content_security_policy"`
	XContentTypeOptions   string                      `yaml:"x_content_type_options"`
	XFrameOptions         string                      `yaml:"x_frame_options"`
	XXSSProtection        string                      `yaml:"x_xss_protection"`
}

// HSTSConfig represents HSTS configuration
type HSTSConfig struct {
	Enabled           bool `yaml:"enabled"`
	MaxAge            int  `yaml:"max_age"`
	IncludeSubdomains bool `yaml:"include_subdomains"`
	Preload           bool `yaml:"preload"`
}

// ContentSecurityPolicyConfig represents CSP configuration
type ContentSecurityPolicyConfig struct {
	Enabled bool   `yaml:"enabled"`
	Policy  string `yaml:"policy"`
}

// DeltaConfig represents Delta Lake configuration
type DeltaConfig struct {
	BasePath           string            `yaml:"base_path"`
	TablePath          string            `yaml:"table_path"`
	LogPath            string            `yaml:"log_path"`
	CheckpointInterval int               `yaml:"checkpoint_interval"`
	MinReaderVersion   int               `yaml:"min_reader_version"`
	MinWriterVersion   int               `yaml:"min_writer_version"`
	PartitionColumns   []string          `yaml:"partition_columns"`
	Properties         map[string]string `yaml:"properties"`
}

// QualityConfig represents data quality configuration
type QualityConfig struct {
	Rules   []QualityRuleConfig  `yaml:"rules"`
	Reports QualityReportsConfig `yaml:"reports"`
}

// QualityRuleConfig represents quality rule configuration
type QualityRuleConfig struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Type        string  `yaml:"type"`
	Column      string  `yaml:"column"`
	Condition   string  `yaml:"condition"`
	Threshold   float64 `yaml:"threshold"`
	Severity    string  `yaml:"severity"`
}

// QualityReportsConfig represents quality reports configuration
type QualityReportsConfig struct {
	OutputDir     string   `yaml:"output_dir"`
	RetentionDays int      `yaml:"retention_days"`
	Formats       []string `yaml:"formats"`
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	Prometheus PrometheusConfig `yaml:"prometheus"`
	Metrics    []MetricConfig   `yaml:"metrics"`
	Alerts     []AlertConfig    `yaml:"alerts"`
}

// PrometheusConfig represents Prometheus configuration
type PrometheusConfig struct {
	Enabled     bool          `yaml:"enabled"`
	PushGateway string        `yaml:"push_gateway"`
	JobName     string        `yaml:"job_name"`
	Interval    time.Duration `yaml:"interval"`
}

// MetricConfig represents metric configuration
type MetricConfig struct {
	Name        string   `yaml:"name"`
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Labels      []string `yaml:"labels"`
}

// AlertConfig represents alert configuration
type AlertConfig struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Severity    string  `yaml:"severity"`
	Condition   string  `yaml:"condition"`
	Threshold   float64 `yaml:"threshold"`
}

// ReportsConfig represents report generation configuration
type ReportsConfig struct {
	OutputDir string           `yaml:"output_dir"`
	Templates []TemplateConfig `yaml:"templates"`
	Retention RetentionConfig  `yaml:"retention"`
}

// TemplateConfig represents report template configuration
type TemplateConfig struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Format      string   `yaml:"format"`
	Template    string   `yaml:"template"`
	Parameters  []string `yaml:"parameters"`
}

// RetentionConfig represents report retention configuration
type RetentionConfig struct {
	Days       int `yaml:"days"`
	MaxReports int `yaml:"max_reports"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string       `yaml:"level"`
	Format string       `yaml:"format"`
	Output string       `yaml:"output"`
	File   FileConfig   `yaml:"file"`
	Fields FieldsConfig `yaml:"fields"`
}

// FileConfig represents file logging configuration
type FileConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// FieldsConfig represents logging fields configuration
type FieldsConfig struct {
	Service     string `yaml:"service"`
	Environment string `yaml:"environment"`
}

// DevelopmentConfig represents development configuration
type DevelopmentConfig struct {
	HotReload    bool         `yaml:"hot_reload"`
	Debug        bool         `yaml:"debug"`
	TestCoverage bool         `yaml:"test_coverage"`
	Lint         LintConfig   `yaml:"lint"`
	Docker       DockerConfig `yaml:"docker"`
}

// LintConfig represents linting configuration
type LintConfig struct {
	Go     GoLintConfig     `yaml:"go"`
	Python PythonLintConfig `yaml:"python"`
}

// GoLintConfig represents Go linting configuration
type GoLintConfig struct {
	Enabled bool     `yaml:"enabled"`
	Tools   []string `yaml:"tools"`
}

// PythonLintConfig represents Python linting configuration
type PythonLintConfig struct {
	Enabled bool     `yaml:"enabled"`
	Tools   []string `yaml:"tools"`
}

// DockerConfig represents Docker configuration
type DockerConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Registry string `yaml:"registry"`
	Tag      string `yaml:"tag"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server configuration
	if config.Server.Port <= 0 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate security configuration
	if config.Security.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if config.Security.TLS.Enabled {
		if config.Security.TLS.CertFile == "" || config.Security.TLS.KeyFile == "" {
			return fmt.Errorf("TLS certificate and key files are required when TLS is enabled")
		}
	}

	// Validate Delta Lake configuration
	if config.Delta.BasePath == "" {
		return fmt.Errorf("Delta Lake base path is required")
	}

	// Validate quality configuration
	for _, rule := range config.Quality.Rules {
		if rule.Name == "" {
			return fmt.Errorf("quality rule name is required")
		}
		if rule.Type == "" {
			return fmt.Errorf("quality rule type is required")
		}
	}

	// Validate monitoring configuration
	if config.Monitoring.Prometheus.Enabled {
		if config.Monitoring.Prometheus.PushGateway == "" {
			return fmt.Errorf("Prometheus push gateway URL is required when monitoring is enabled")
		}
	}

	// Validate report configuration
	if config.Reports.OutputDir == "" {
		return fmt.Errorf("report output directory is required")
	}

	// Validate logging configuration
	if config.Logging.Level == "" {
		return fmt.Errorf("logging level is required")
	}

	return nil
}

// SaveConfig saves the configuration to a file
func SaveConfig(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:            "0.0.0.0",
			Port:            8080,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     120 * time.Second,
			ShutdownTimeout: 10 * time.Second,
			MaxConnections:  1000,
			CORS: CORSConfig{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"Content-Type", "Authorization"},
				ExposedHeaders: []string{"Content-Length"},
				MaxAge:         86400,
			},
		},
		Security: SecurityConfig{
			JWT: JWTConfig{
				Secret:     "your-secret-key",
				Expiration: 24 * time.Hour,
			},
			TLS: TLSConfig{
				Enabled:    false,
				MinVersion: "1.2",
			},
			RateLimit: RateLimitConfig{
				RequestsPerMinute: 60,
				Burst:             10,
			},
			IPAllowlist: IPAllowlistConfig{
				Enabled: true,
				IPs:     []string{"127.0.0.1", "::1"},
			},
			DefaultRole: "viewer",
		},
		Delta: DeltaConfig{
			BasePath:           "data/delta",
			TablePath:          "data/delta/tables",
			LogPath:            "data/delta/logs",
			CheckpointInterval: 10,
			MinReaderVersion:   1,
			MinWriterVersion:   2,
		},
		Quality: QualityConfig{
			Reports: QualityReportsConfig{
				OutputDir:     "data/quality/reports",
				RetentionDays: 30,
				Formats:       []string{"pdf", "csv", "html", "json"},
			},
		},
		Monitoring: MonitoringConfig{
			Prometheus: PrometheusConfig{
				Enabled:     true,
				PushGateway: "http://localhost:9091",
				JobName:     "nessi",
				Interval:    15 * time.Second,
			},
		},
		Reports: ReportsConfig{
			OutputDir: "data/reports",
			Retention: RetentionConfig{
				Days:       90,
				MaxReports: 1000,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
			Fields: FieldsConfig{
				Service:     "nessi",
				Environment: "development",
			},
		},
		Development: DevelopmentConfig{
			HotReload:    true,
			Debug:        true,
			TestCoverage: true,
			Lint: LintConfig{
				Go: GoLintConfig{
					Enabled: true,
					Tools:   []string{"golangci-lint"},
				},
				Python: PythonLintConfig{
					Enabled: true,
					Tools:   []string{"pylint", "black", "isort"},
				},
			},
			Docker: DockerConfig{
				Enabled:  true,
				Registry: "localhost:5000",
				Tag:      "latest",
			},
		},
	}
}
