package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"text/tabwriter"

	"golang.org/x/term"
)

// TokenResponse represents the response from a login request
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// UserResponse represents a user response
type UserResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	LastLogin string `json:"last_login,omitempty"`
}

// APIKeyResponse represents an API key response
type APIKeyResponse struct {
	Username string `json:"username"`
	APIKey   string `json:"api_key"`
}

// Config represents the CLI configuration
type Config struct {
	ServerURL string `json:"server_url"`
	Token     string `json:"token"`
	Username  string `json:"username"`
	APIKey    string `json:"api_key"`
}

// getConfigPath returns the path to the config file
func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".nessi-monitor.json"
	}
	return filepath.Join(homeDir, ".nessi-monitor.json")
}

// loadConfig loads the CLI configuration
func loadConfig() (*Config, error) {
	configPath := getConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config
		return &Config{
			ServerURL: "http://localhost:9090",
		}, nil
	}

	// Read config file
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse config
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// saveConfig saves the CLI configuration
func saveConfig(config *Config) error {
	configPath := getConfigPath()

	// Marshal config to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := ioutil.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// login handles authentication
func login(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("login", flag.ExitOnError)
	serverURL := fs.String("server", "http://localhost:9090", "Server URL")
	username := fs.String("username", "", "Username")
	password := fs.String("password", "", "Password (not recommended, use interactive mode)")
	help := fs.Bool("help", false, "Show help")
	fs.Parse(args)

	// Show help if requested
	if *help {
		fmt.Println("Usage: nessi-monitor login [options]")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		return nil
	}

	// Get username if not provided
	if *username == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Username: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read username: %v", err)
		}
		*username = strings.TrimSpace(input)
	}

	// Get password if not provided
	var passwordStr string
	if *password == "" {
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		fmt.Println() // Add newline after password input
		passwordStr = string(passwordBytes)
	} else {
		passwordStr = *password
	}

	// Prepare login request
	reqBody := map[string]string{
		"username": *username,
		"password": passwordStr,
	}
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make login request
	url := fmt.Sprintf("%s/auth/login", *serverURL)
	resp, err := http.Post(url, "application/json", strings.NewReader(string(reqData)))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: %s", string(body))
	}

	// Parse response
	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Save token to config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	config.ServerURL = *serverURL
	config.Token = tokenResp.Token
	config.Username = tokenResp.Username

	if err := saveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	fmt.Printf("Logged in as %s (role: %s)\n", tokenResp.Username, tokenResp.Role)
	fmt.Printf("Token expires at: %s\n", tokenResp.ExpiresAt)
	return nil
}

// manageUsers handles user management
func manageUsers(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("user", flag.ExitOnError)
	action := fs.String("action", "list", "Action to perform (list, create, update, delete)")
	serverURL := fs.String("server", "", "Server URL (defaults to saved config)")
	username := fs.String("username", "", "Username (required for create, update, delete)")
	email := fs.String("email", "", "Email (for create, update)")
	password := fs.String("password", "", "Password (for create, update)")
	role := fs.String("role", "", "Role (for create, update)")
	token := fs.String("token", "", "Authentication token (defaults to saved config)")
	help := fs.Bool("help", false, "Show help")
	fs.Parse(args)

	// Show help if requested
	if *help {
		fmt.Println("Usage: nessi-monitor user [options]")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		return nil
	}

	// Load config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Use config values if not provided
	if *serverURL == "" {
		*serverURL = config.ServerURL
	}
	if *token == "" {
		*token = config.Token
	}

	// Check if token is available
	if *token == "" {
		return fmt.Errorf("authentication token required, please login first")
	}

	// Perform action
	switch *action {
	case "list":
		return listUsers(*serverURL, *token)
	case "create":
		return createUser(*serverURL, *token, *username, *email, *password, *role)
	case "update":
		return updateUser(*serverURL, *token, *username, *email, *password, *role)
	case "delete":
		return deleteUser(*serverURL, *token, *username)
	default:
		return fmt.Errorf("invalid action: %s", *action)
	}
}

// listUsers lists all users
func listUsers(serverURL, token string) error {
	// Make request
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/users", serverURL), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to list users: %s", string(body))
	}

	// Parse response
	var users []UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Display users
	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	// Create tabwriter
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print header
	fmt.Fprintf(w, "Username\tEmail\tRole\tCreated At\tLast Login\n")
	fmt.Fprintf(w, "--------\t-----\t----\t----------\t----------\n")

	// Print users
	for _, user := range users {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", user.Username, user.Email, user.Role, user.CreatedAt, user.LastLogin)
	}

	return nil
}

// createUser creates a new user
func createUser(serverURL, token, username, email, password, role string) error {
	// Validate required fields
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if password == "" {
		// Prompt for password
		fmt.Print("Password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		fmt.Println() // Add newline after password input
		password = string(passwordBytes)
	}

	// Prepare request
	reqBody := map[string]string{
		"username": username,
		"password": password,
	}
	if email != "" {
		reqBody["email"] = email
	}
	if role != "" {
		reqBody["role"] = role
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make request
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/admin/users", serverURL), strings.NewReader(string(reqData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create user: %s", string(body))
	}

	// Parse response
	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	fmt.Printf("User '%s' created successfully with role '%s'\n", user.Username, user.Role)
	return nil
}

// updateUser updates a user
func updateUser(serverURL, token, username, email, password, role string) error {
	// Validate required fields
	if username == "" {
		return fmt.Errorf("username is required")
	}

	// Prepare request
	reqBody := make(map[string]string)
	if email != "" {
		reqBody["email"] = email
	}
	if password != "" {
		reqBody["password"] = password
	}
	if role != "" {
		reqBody["role"] = role
	}

	// Check if there's anything to update
	if len(reqBody) == 0 {
		return fmt.Errorf("no update parameters provided")
	}

	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make request
	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/admin/users?username=%s", serverURL, username), strings.NewReader(string(reqData)))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to update user: %s", string(body))
	}

	// Parse response
	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	fmt.Printf("User '%s' updated successfully\n", user.Username)
	return nil
}

// deleteUser deletes a user
func deleteUser(serverURL, token, username string) error {
	// Validate required fields
	if username == "" {
		return fmt.Errorf("username is required")
	}

	// Confirm deletion
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to delete user '%s'? (y/n): ", username)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}
	if strings.ToLower(strings.TrimSpace(input)) != "y" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	// Make request
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/admin/users?username=%s", serverURL, username), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete user: %s", string(body))
	}

	fmt.Printf("User '%s' deleted successfully\n", username)
	return nil
}

// manageAPIKey handles API key management
func manageAPIKey(args []string) error {
	// Parse flags
	fs := flag.NewFlagSet("apikey", flag.ExitOnError)
	serverURL := fs.String("server", "", "Server URL (defaults to saved config)")
	username := fs.String("username", "", "Username (defaults to current user)")
	token := fs.String("token", "", "Authentication token (defaults to saved config)")
	help := fs.Bool("help", false, "Show help")
	fs.Parse(args)

	// Show help if requested
	if *help {
		fmt.Println("Usage: nessi-monitor apikey [options]")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		return nil
	}

	// Load config
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Use config values if not provided
	if *serverURL == "" {
		*serverURL = config.ServerURL
	}
	if *token == "" {
		*token = config.Token
	}
	if *username == "" {
		*username = config.Username
	}

	// Check if token is available
	if *token == "" {
		return fmt.Errorf("authentication token required, please login first")
	}

	// Check if username is available
	if *username == "" {
		return fmt.Errorf("username is required")
	}

	// Make request
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/admin/users/%s/apikey", *serverURL, *username), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to generate API key: %s", string(body))
	}

	// Parse response
	var apiKeyResp APIKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiKeyResp); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	// Save API key to config if it's for the current user
	if apiKeyResp.Username == config.Username {
		config.APIKey = apiKeyResp.APIKey
		if err := saveConfig(config); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}
	}

	fmt.Printf("API key for user '%s' generated successfully\n", apiKeyResp.Username)
	fmt.Printf("API Key: %s\n", apiKeyResp.APIKey)
	fmt.Println("\nIMPORTANT: Save this API key securely. It will not be shown again.")
	return nil
}
