package config

// Implement a configuration system for Nessi.dev with the following requirements:
// 1. Define a Config struct with these nested components:
//    - Server: host, port, tls settings
//    - Extensions: enabledExtensions list, pythonEnabled bool, pythonPath string
//    - Delta: defaultLocation, maxVersions to keep
//    - Quality: defaultRules, thresholds
//    - Security: auth settings, rbac settings
// 2. Implement a Load() function that:
//    - Searches for config files in standard locations (./config.yaml, ~/.nessi/config.yaml, etc.)
//    - Falls back to default values if no config file found
//    - Overrides settings with environment variables (NESSI_SERVER_HOST, etc.)
//    - Validates the configuration
// 3. Implement a GetDefaultConfig() function that returns sensible defaults
// 4. Add a Save() method to write the config back to a file
// Use Viper for configuration management and provide helpful error messages
