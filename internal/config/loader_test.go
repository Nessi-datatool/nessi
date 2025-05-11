package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a test configuration file
	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
server:
  host: "localhost"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  max_connections: 1000
  cors:
    allowed_origins: ["*"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allowed_headers: ["Content-Type", "Authorization"]
    exposed_headers: ["Content-Length"]
    max_age: 86400

security:
  jwt:
    secret: "test-secret"
    expiration: 24h
  tls:
    enabled: false
    min_version: "1.2"
  rate_limit:
    requests_per_minute: 60
    burst: 10
  ip_allowlist:
    enabled: true
    ips: ["127.0.0.1", "::1"]
  default_role: "viewer"

delta:
  base_path: "data/delta"
  table_path: "data/delta/tables"
  log_path: "data/delta/logs"
  checkpoint_interval: 10
  min_reader_version: 1
  min_writer_version: 2

quality:
  rules:
    - name: "completeness"
      description: "Check for null values"
      type: "null_check"
      column: "id"
      condition: "is_not_null"
      threshold: 0.95
      severity: "error"
  reports:
    output_dir: "data/quality/reports"
    retention_days: 30
    formats: ["pdf", "csv", "html", "json"]

monitoring:
  prometheus:
    enabled: true
    push_gateway: "http://localhost:9091"
    job_name: "nessi"
    interval: 15s
  metrics:
    - name: "request_count"
      type: "counter"
      description: "Total number of requests"
      labels: ["method", "path"]
  alerts:
    - name: "high_error_rate"
      description: "Error rate exceeds threshold"
      severity: "critical"
      condition: "error_rate > 0.1"
      threshold: 0.1

reports:
  output_dir: "data/reports"
  templates:
    - id: "quality_summary"
      name: "Quality Summary"
      description: "Summary of data quality metrics"
      format: "pdf"
      template: "templates/quality_summary.html"
      parameters: ["start_date", "end_date"]
  retention:
    days: 90
    max_reports: 1000

logging:
  level: "info"
  format: "json"
  output: "stdout"
  file:
    enabled: true
    path: "logs/nessi.log"
    max_size: 100
    max_backups: 3
    max_age: 30
    compress: true
  fields:
    service: "nessi"
    environment: "test"

development:
  hot_reload: true
  debug: true
  test_coverage: true
  lint:
    go:
      enabled: true
      tools: ["golangci-lint"]
    python:
      enabled: true
      tools: ["pylint", "black", "isort"]
  docker:
    enabled: true
    registry: "localhost:5000"
    tag: "test"
`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Test loading the configuration
	config, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, config)

	// Verify server configuration
	assert.Equal(t, "localhost", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, 30*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.Server.IdleTimeout)
	assert.Equal(t, 1000, config.Server.MaxConnections)
	assert.Equal(t, []string{"*"}, config.Server.CORS.AllowedOrigins)
	assert.Equal(t, []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, config.Server.CORS.AllowedMethods)
	assert.Equal(t, []string{"Content-Type", "Authorization"}, config.Server.CORS.AllowedHeaders)
	assert.Equal(t, []string{"Content-Length"}, config.Server.CORS.ExposedHeaders)
	assert.Equal(t, 86400, config.Server.CORS.MaxAge)

	// Verify security configuration
	assert.Equal(t, "test-secret", config.Security.JWT.Secret)
	assert.Equal(t, 24*time.Hour, config.Security.JWT.Expiration)
	assert.False(t, config.Security.TLS.Enabled)
	assert.Equal(t, "1.2", config.Security.TLS.MinVersion)
	assert.Equal(t, 60, config.Security.RateLimit.RequestsPerMinute)
	assert.Equal(t, 10, config.Security.RateLimit.Burst)
	assert.True(t, config.Security.IPAllowlist.Enabled)
	assert.Equal(t, []string{"127.0.0.1", "::1"}, config.Security.IPAllowlist.IPs)
	assert.Equal(t, "viewer", config.Security.DefaultRole)

	// Verify Delta Lake configuration
	assert.Equal(t, "data/delta", config.Delta.BasePath)
	assert.Equal(t, "data/delta/tables", config.Delta.TablePath)
	assert.Equal(t, "data/delta/logs", config.Delta.LogPath)
	assert.Equal(t, 10, config.Delta.CheckpointInterval)
	assert.Equal(t, 1, config.Delta.MinReaderVersion)
	assert.Equal(t, 2, config.Delta.MinWriterVersion)

	// Verify quality configuration
	require.Len(t, config.Quality.Rules, 1)
	rule := config.Quality.Rules[0]
	assert.Equal(t, "completeness", rule.Name)
	assert.Equal(t, "Check for null values", rule.Description)
	assert.Equal(t, "null_check", rule.Type)
	assert.Equal(t, "id", rule.Column)
	assert.Equal(t, "is_not_null", rule.Condition)
	assert.Equal(t, 0.95, rule.Threshold)
	assert.Equal(t, "error", rule.Severity)
	assert.Equal(t, "data/quality/reports", config.Quality.Reports.OutputDir)
	assert.Equal(t, 30, config.Quality.Reports.RetentionDays)
	assert.Equal(t, []string{"pdf", "csv", "html", "json"}, config.Quality.Reports.Formats)

	// Verify monitoring configuration
	assert.True(t, config.Monitoring.Prometheus.Enabled)
	assert.Equal(t, "http://localhost:9091", config.Monitoring.Prometheus.PushGateway)
	assert.Equal(t, "nessi", config.Monitoring.Prometheus.JobName)
	assert.Equal(t, 15*time.Second, config.Monitoring.Prometheus.Interval)
	require.Len(t, config.Monitoring.Metrics, 1)
	metric := config.Monitoring.Metrics[0]
	assert.Equal(t, "request_count", metric.Name)
	assert.Equal(t, "counter", metric.Type)
	assert.Equal(t, "Total number of requests", metric.Description)
	assert.Equal(t, []string{"method", "path"}, metric.Labels)
	require.Len(t, config.Monitoring.Alerts, 1)
	alert := config.Monitoring.Alerts[0]
	assert.Equal(t, "high_error_rate", alert.Name)
	assert.Equal(t, "Error rate exceeds threshold", alert.Description)
	assert.Equal(t, "critical", alert.Severity)
	assert.Equal(t, "error_rate > 0.1", alert.Condition)
	assert.Equal(t, 0.1, alert.Threshold)

	// Verify reports configuration
	assert.Equal(t, "data/reports", config.Reports.OutputDir)
	require.Len(t, config.Reports.Templates, 1)
	template := config.Reports.Templates[0]
	assert.Equal(t, "quality_summary", template.ID)
	assert.Equal(t, "Quality Summary", template.Name)
	assert.Equal(t, "Summary of data quality metrics", template.Description)
	assert.Equal(t, "pdf", template.Format)
	assert.Equal(t, "templates/quality_summary.html", template.Template)
	assert.Equal(t, []string{"start_date", "end_date"}, template.Parameters)
	assert.Equal(t, 90, config.Reports.Retention.Days)
	assert.Equal(t, 1000, config.Reports.Retention.MaxReports)

	// Verify logging configuration
	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "json", config.Logging.Format)
	assert.Equal(t, "stdout", config.Logging.Output)
	assert.True(t, config.Logging.File.Enabled)
	assert.Equal(t, "logs/nessi.log", config.Logging.File.Path)
	assert.Equal(t, 100, config.Logging.File.MaxSize)
	assert.Equal(t, 3, config.Logging.File.MaxBackups)
	assert.Equal(t, 30, config.Logging.File.MaxAge)
	assert.True(t, config.Logging.File.Compress)
	assert.Equal(t, "nessi", config.Logging.Fields.Service)
	assert.Equal(t, "test", config.Logging.Fields.Environment)

	// Verify development configuration
	assert.True(t, config.Development.HotReload)
	assert.True(t, config.Development.Debug)
	assert.True(t, config.Development.TestCoverage)
	assert.True(t, config.Development.Lint.Go.Enabled)
	assert.Equal(t, []string{"golangci-lint"}, config.Development.Lint.Go.Tools)
	assert.True(t, config.Development.Lint.Python.Enabled)
	assert.Equal(t, []string{"pylint", "black", "isort"}, config.Development.Lint.Python.Tools)
	assert.True(t, config.Development.Docker.Enabled)
	assert.Equal(t, "localhost:5000", config.Development.Docker.Registry)
	assert.Equal(t, "test", config.Development.Docker.Tag)
}

func TestLoadConfigInvalid(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test cases for invalid configurations
	testCases := []struct {
		name    string
		content string
		errMsg  string
	}{
		{
			name: "invalid port",
			content: `
server:
  port: -1
`,
			errMsg: "invalid server port: -1",
		},
		{
			name: "missing JWT secret",
			content: `
security:
  jwt:
    secret: ""
`,
			errMsg: "JWT secret is required",
		},
		{
			name: "TLS enabled without certificates",
			content: `
security:
  jwt:
    secret: "test-dummy-secret"
  tls:
    enabled: true
`,
			errMsg: "TLS certificate and key files are required when TLS is enabled",
		},
		{
			name: "missing Delta Lake base path",
			content: `
security:
  jwt:
    secret: "test-dummy-secret"
delta:
  base_path: ""
`,
			errMsg: "Delta Lake base path is required",
		},
		{
			name: "invalid quality rule",
			content: `
security:
  jwt:
    secret: "test-dummy-secret"
delta:
  base_path: "/tmp/dummy-delta-path"
quality:
  rules:
    - name: ""
      type: ""
`,
			errMsg: "quality rule name is required",
		},
		{
			name: "Prometheus enabled without push gateway",
			content: `
security:
  jwt:
    secret: "test-dummy-secret"
delta:
  base_path: "/tmp/dummy-delta-path"
monitoring:
  prometheus:
    enabled: true
`,
			errMsg: "Prometheus push gateway URL is required when monitoring is enabled",
		},
		{
			name: "missing report output directory",
			content: `
security:
  jwt:
    secret: "test-dummy-secret"
delta:
  base_path: "/tmp/dummy-delta-path"
reports:
  output_dir: ""
`,
			errMsg: "report output directory is required",
		},
		{
			name: "missing logging level",
			content: `
security:
  jwt:
    secret: "test-dummy-secret"
delta:
  base_path: "/tmp/dummy-delta-path"
reports:
  output_dir: "/tmp/dummy-reports"
logging:
  level: ""
`,
			errMsg: "logging level is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join(tmpDir, "invalid_config.yaml")
			// Prepend a valid server port for most tests, unless the test is specifically about port or empty config.
			if tc.name != "invalid port" && tc.name != "empty_config" {
				tc.content = "server:\n  port: 8080\n" + tc.content
			}

			err := os.WriteFile(configPath, []byte(tc.content), 0644)
			require.NoError(t, err)

			_, err = LoadConfig(configPath)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errMsg)
		})
	}
}

func TestSaveConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Get default configuration
	config := GetDefaultConfig()

	// Save configuration
	configPath := filepath.Join(tmpDir, "saved_config.yaml")
	err = SaveConfig(config, configPath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(configPath)
	require.NoError(t, err)

	// Load saved configuration
	savedConfig, err := LoadConfig(configPath)
	require.NoError(t, err)
	assert.NotNil(t, savedConfig)

	// Verify configuration matches
	assert.Equal(t, config.Server.Host, savedConfig.Server.Host)
	assert.Equal(t, config.Server.Port, savedConfig.Server.Port)
	assert.Equal(t, config.Security.JWT.Secret, savedConfig.Security.JWT.Secret)
	assert.Equal(t, config.Delta.BasePath, savedConfig.Delta.BasePath)
	assert.Equal(t, config.Logging.Level, savedConfig.Logging.Level)
}

func TestGetDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	assert.NotNil(t, config)

	// Verify default values
	assert.Equal(t, "0.0.0.0", config.Server.Host)
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, 30*time.Second, config.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.Server.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.Server.IdleTimeout)
	assert.Equal(t, 1000, config.Server.MaxConnections)

	assert.Equal(t, "your-secret-key", config.Security.JWT.Secret)
	assert.Equal(t, 24*time.Hour, config.Security.JWT.Expiration)
	assert.False(t, config.Security.TLS.Enabled)
	assert.Equal(t, "1.2", config.Security.TLS.MinVersion)
	assert.Equal(t, 60, config.Security.RateLimit.RequestsPerMinute)
	assert.Equal(t, 10, config.Security.RateLimit.Burst)
	assert.True(t, config.Security.IPAllowlist.Enabled)
	assert.Equal(t, []string{"127.0.0.1", "::1"}, config.Security.IPAllowlist.IPs)
	assert.Equal(t, "viewer", config.Security.DefaultRole)

	assert.Equal(t, "data/delta", config.Delta.BasePath)
	assert.Equal(t, "data/delta/tables", config.Delta.TablePath)
	assert.Equal(t, "data/delta/logs", config.Delta.LogPath)
	assert.Equal(t, 10, config.Delta.CheckpointInterval)
	assert.Equal(t, 1, config.Delta.MinReaderVersion)
	assert.Equal(t, 2, config.Delta.MinWriterVersion)

	assert.Equal(t, "data/quality/reports", config.Quality.Reports.OutputDir)
	assert.Equal(t, 30, config.Quality.Reports.RetentionDays)
	assert.Equal(t, []string{"pdf", "csv", "html", "json"}, config.Quality.Reports.Formats)

	assert.True(t, config.Monitoring.Prometheus.Enabled)
	assert.Equal(t, "http://localhost:9091", config.Monitoring.Prometheus.PushGateway)
	assert.Equal(t, "nessi", config.Monitoring.Prometheus.JobName)
	assert.Equal(t, 15*time.Second, config.Monitoring.Prometheus.Interval)

	assert.Equal(t, "data/reports", config.Reports.OutputDir)
	assert.Equal(t, 90, config.Reports.Retention.Days)
	assert.Equal(t, 1000, config.Reports.Retention.MaxReports)

	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "json", config.Logging.Format)
	assert.Equal(t, "stdout", config.Logging.Output)
	assert.Equal(t, "nessi", config.Logging.Fields.Service)
	assert.Equal(t, "development", config.Logging.Fields.Environment)

	assert.True(t, config.Development.HotReload)
	assert.True(t, config.Development.Debug)
	assert.True(t, config.Development.TestCoverage)
	assert.True(t, config.Development.Lint.Go.Enabled)
	assert.Equal(t, []string{"golangci-lint"}, config.Development.Lint.Go.Tools)
	assert.True(t, config.Development.Lint.Python.Enabled)
	assert.Equal(t, []string{"pylint", "black", "isort"}, config.Development.Lint.Python.Tools)
	assert.True(t, config.Development.Docker.Enabled)
	assert.Equal(t, "localhost:5000", config.Development.Docker.Registry)
	assert.Equal(t, "latest", config.Development.Docker.Tag)
} 