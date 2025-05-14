// Implement a Python extension interface with these requirements:
// 1. Define a PythonExtension struct that implements the Extension interface:
//    - Embeds a BaseExtension for common functionality
//    - Adds fields for Python-specific attributes:
//      - ProcessHandle for the Python process
//      - Communication channel (pipe or socket)
//      - Path to Python script
// 2. Implement methods for the Extension interface:
//    - Start() error - Start the Python process
//    - Stop() error - Stop the Python process gracefully
//    - Status() string - Get process status
//    - IsRunning() bool - Check if process is running
// 3. Communication functions:
//    - SendMessage(msg interface{}) error - Send a message to Python process
//    - ReceiveMessage() (interface{}, error) - Receive a message from Python
//    - ExecuteFunction(name string, args ...interface{}) (interface{}, error)
// 4. Helper functions:
//    - startPythonProcess() error - Start Python interpreter
//    - setupCommunicationChannel() error - Set up IPC
//    - monitorProcess() - Monitor process health
//    - handleCrash() - Handle unexpected termination
// Use JSON for data serialization between Go and Python
// Handle process lifecycle and cleanup properly
// Include timeout mechanisms for operations
// Provide detailed error information when Python operations fail

package extensions

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// BaseExtension provides common fields for extensions.
// It is intended to be embedded in specific extension types.
type BaseExtension struct {
	Name        string
	Version     string
	Description string
}

// PythonExtension represents a Python-based extension
type PythonExtension struct {
	BaseExtension
	process     *exec.Cmd
	socket      net.Conn
	scriptPath  string
	health      *ExtensionHealth
	stopChan    chan struct{}
	errorChan   chan error
	mu          sync.RWMutex
}

// ExtensionHealth tracks extension process health
type ExtensionHealth struct {
	LastHeartbeat time.Time
	Status        string
	Error         error
	mu            sync.RWMutex
}

// PythonPlugin represents a Python-based extension plugin
type PythonPlugin struct {
	Name         string
	Version      string
	Description  string
	EntryPoint   string
	Requirements []string
	Enabled      bool
	Process      *exec.Cmd
	Socket       net.Conn
	Health       *PluginHealth
}

// PluginHealth tracks plugin process health
type PluginHealth struct {
	LastHeartbeat time.Time
	Status        string
	Error         error
	mu            sync.RWMutex
}

// PythonManager handles Python extension management
type PythonManager struct {
	pluginsDir string
	venvDir    string
	pythonPath string
	plugins    map[string]*PythonPlugin
	mu         sync.RWMutex
}

// NewPythonManager creates a new Python extension manager
func NewPythonManager(pluginsDir string) (*PythonManager, error) {
	// Create plugins directory if it doesn't exist
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugins directory: %w", err)
	}

	// Create virtual environment directory
	venvDir := filepath.Join(pluginsDir, "venv")
	if err := os.MkdirAll(venvDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create venv directory: %w", err)
	}

	// Find Python executable
	pythonPath, err := findPython()
	if err != nil {
		return nil, fmt.Errorf("failed to find Python: %w", err)
	}

	return &PythonManager{
		pluginsDir: pluginsDir,
		venvDir:    venvDir,
		pythonPath: pythonPath,
		plugins:    make(map[string]*PythonPlugin),
	}, nil
}

// findPython locates the Python executable
func findPython() (string, error) {
	// Try python3 first
	if path, err := exec.LookPath("python3"); err == nil {
		return path, nil
	}

	// Fall back to python
	if path, err := exec.LookPath("python"); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("python executable not found")
}

// LoadPlugin loads a Python extension plugin
func (m *PythonManager) LoadPlugin(name string) (*PythonPlugin, error) {
	pluginDir := filepath.Join(m.pluginsDir, name)
	manifestPath := filepath.Join(pluginDir, "plugin.json")

	// Read plugin manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin manifest: %w", err)
	}

	var plugin PythonPlugin
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, fmt.Errorf("failed to parse plugin manifest: %w", err)
	}

	// Validate plugin
	if err := m.validatePlugin(&plugin); err != nil {
		return nil, fmt.Errorf("invalid plugin: %w", err)
	}

	return &plugin, nil
}

// validatePlugin validates a Python plugin
func (m *PythonManager) validatePlugin(plugin *PythonPlugin) error {
	// Check required fields
	if plugin.Name == "" {
		return fmt.Errorf("plugin name is required")
	}
	if plugin.Version == "" {
		return fmt.Errorf("plugin version is required")
	}
	if plugin.EntryPoint == "" {
		return fmt.Errorf("plugin entry point is required")
	}

	// Check entry point file exists
	entryPath := filepath.Join(m.pluginsDir, plugin.Name, plugin.EntryPoint)
	if _, err := os.Stat(entryPath); err != nil {
		return fmt.Errorf("entry point file not found: %w", err)
	}

	return nil
}

// InstallPlugin installs a Python extension plugin
func (m *PythonManager) InstallPlugin(name, source string) error {
	pluginDir := filepath.Join(m.pluginsDir, name)

	// Create plugin directory
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Clone or copy plugin files
	if strings.HasPrefix(source, "git+") {
		if err := m.clonePlugin(name, strings.TrimPrefix(source, "git+")); err != nil {
			return fmt.Errorf("failed to clone plugin: %w", err)
		}
	} else {
		if err := m.copyPlugin(name, source); err != nil {
			return fmt.Errorf("failed to copy plugin: %w", err)
		}
	}

	// Install dependencies
	if err := m.installDependencies(name); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

// clonePlugin clones a plugin from a Git repository
func (m *PythonManager) clonePlugin(name, repo string) error {
	pluginDir := filepath.Join(m.pluginsDir, name)
	cmd := exec.Command("git", "clone", repo, pluginDir)
	return cmd.Run()
}

// copyPlugin copies a plugin from a local directory
func (m *PythonManager) copyPlugin(name, source string) error {
	pluginDir := filepath.Join(m.pluginsDir, name)
	cmd := exec.Command("cp", "-r", source, pluginDir)
	return cmd.Run()
}

// installDependencies installs plugin dependencies
func (m *PythonManager) installDependencies(name string) error {
	plugin, err := m.LoadPlugin(name)
	if err != nil {
		return err
	}

	// Create requirements.txt if not exists
	reqPath := filepath.Join(m.pluginsDir, name, "requirements.txt")
	if len(plugin.Requirements) > 0 {
		reqData := strings.Join(plugin.Requirements, "\n")
		if err := os.WriteFile(reqPath, []byte(reqData), 0644); err != nil {
			return fmt.Errorf("failed to write requirements.txt: %w", err)
		}
	}

	// Install dependencies using pip
	cmd := exec.Command(m.pythonPath, "-m", "pip", "install", "-r", reqPath)
	return cmd.Run()
}

// UninstallPlugin uninstalls a Python extension plugin
func (m *PythonManager) UninstallPlugin(name string) error {
	pluginDir := filepath.Join(m.pluginsDir, name)
	return os.RemoveAll(pluginDir)
}

// EnablePlugin enables a Python extension plugin
func (m *PythonManager) EnablePlugin(name string) error {
	plugin, err := m.LoadPlugin(name)
	if err != nil {
		return err
	}

	plugin.Enabled = true
	return m.savePlugin(plugin)
}

// DisablePlugin disables a Python extension plugin
func (m *PythonManager) DisablePlugin(name string) error {
	plugin, err := m.LoadPlugin(name)
	if err != nil {
		return err
	}

	plugin.Enabled = false
	return m.savePlugin(plugin)
}

// savePlugin saves plugin configuration
func (m *PythonManager) savePlugin(plugin *PythonPlugin) error {
	manifestPath := filepath.Join(m.pluginsDir, plugin.Name, "plugin.json")
	data, err := json.MarshalIndent(plugin, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal plugin: %w", err)
	}

	return os.WriteFile(manifestPath, data, 0644)
}

// ListPlugins lists all installed Python plugins
func (m *PythonManager) ListPlugins() ([]*PythonPlugin, error) {
	entries, err := os.ReadDir(m.pluginsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %w", err)
	}

	var plugins []*PythonPlugin
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		plugin, err := m.LoadPlugin(entry.Name())
		if err != nil {
			// Skip invalid plugins
			continue
		}

		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

// StartPlugin starts a Python plugin with IPC
func (m *PythonManager) StartPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, err := m.LoadPlugin(name)
	if err != nil {
		return err
	}

	// Create Unix domain socket for IPC
	socketPath := filepath.Join(m.pluginsDir, name, "plugin.sock")
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to create IPC socket: %w", err)
	}

	// Start plugin process
	entryPath := filepath.Join(m.pluginsDir, name, plugin.EntryPoint)
	cmd := exec.Command(m.pythonPath, entryPath, "--socket", socketPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		listener.Close()
		return fmt.Errorf("failed to start plugin: %w", err)
	}

	// Accept connection from plugin
	conn, err := listener.Accept()
	if err != nil {
		cmd.Process.Kill()
		listener.Close()
		return fmt.Errorf("failed to accept plugin connection: %w", err)
	}
	listener.Close()

	// Initialize plugin health monitoring
	plugin.Process = cmd
	plugin.Socket = conn
	plugin.Health = &PluginHealth{
		LastHeartbeat: time.Now(),
		Status:        "running",
	}

	// Start health monitoring
	go m.monitorPluginHealth(plugin)

	m.plugins[name] = plugin
	return nil
}

// StopPlugin stops a running Python plugin
func (m *PythonManager) StopPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// Send shutdown signal
	if plugin.Socket != nil {
		shutdown := map[string]string{"action": "shutdown"}
		if err := json.NewEncoder(plugin.Socket).Encode(shutdown); err != nil {
			// Log error but continue with process termination
			fmt.Printf("Failed to send shutdown signal: %v\n", err)
		}
		plugin.Socket.Close()
	}

	// Terminate process
	if plugin.Process != nil && plugin.Process.Process != nil {
		if err := plugin.Process.Process.Signal(os.Interrupt); err != nil {
			// Force kill if graceful shutdown fails
			plugin.Process.Process.Kill()
		}
	}

	delete(m.plugins, name)
	return nil
}

// monitorPluginHealth monitors plugin process health
func (m *PythonManager) monitorPluginHealth(plugin *PythonPlugin) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Check if process is still running
		if plugin.Process.ProcessState != nil && plugin.Process.ProcessState.Exited() {
			plugin.Health.mu.Lock()
			plugin.Health.Status = "stopped"
			plugin.Health.Error = fmt.Errorf("process exited with code %d", plugin.Process.ProcessState.ExitCode())
			plugin.Health.mu.Unlock()
			return
		}

		// Send heartbeat request
		heartbeat := map[string]string{"action": "heartbeat"}
		if err := json.NewEncoder(plugin.Socket).Encode(heartbeat); err != nil {
			plugin.Health.mu.Lock()
			plugin.Health.Status = "error"
			plugin.Health.Error = fmt.Errorf("heartbeat failed: %w", err)
			plugin.Health.mu.Unlock()
			return
		}

		// Wait for response
		var response map[string]interface{}
		if err := json.NewDecoder(plugin.Socket).Decode(&response); err != nil {
			plugin.Health.mu.Lock()
			plugin.Health.Status = "error"
			plugin.Health.Error = fmt.Errorf("heartbeat response failed: %w", err)
			plugin.Health.mu.Unlock()
			return
		}

		plugin.Health.mu.Lock()
		plugin.Health.LastHeartbeat = time.Now()
		plugin.Health.Status = "running"
		plugin.Health.Error = nil
		plugin.Health.mu.Unlock()
	}
}

// ExecutePlugin executes a Python plugin with data exchange
func (m *PythonManager) ExecutePlugin(ctx context.Context, name string, data interface{}) (interface{}, error) {
	m.mu.RLock()
	plugin, exists := m.plugins[name]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	// Check plugin health
	plugin.Health.mu.RLock()
	if plugin.Health.Status != "running" {
		plugin.Health.mu.RUnlock()
		return nil, fmt.Errorf("plugin %s is not healthy: %v", name, plugin.Health.Error)
	}
	plugin.Health.mu.RUnlock()

	// Prepare request
	request := map[string]interface{}{
		"action": "execute",
		"data":   data,
	}

	// Send request with timeout
	done := make(chan error, 1)
	go func() {
		done <- json.NewEncoder(plugin.Socket).Encode(request)
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Read response with timeout
	var response map[string]interface{}
	done = make(chan error, 1)
	go func() {
		done <- json.NewDecoder(plugin.Socket).Decode(&response)
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Check for error in response
	if err, ok := response["error"].(string); ok {
		return nil, fmt.Errorf("plugin error: %s", err)
	}

	return response["result"], nil
}

// NewPythonExtension creates a new Python extension
func NewPythonExtension(name, version, description, scriptPath string) *PythonExtension {
	return &PythonExtension{
		BaseExtension: BaseExtension{
			Name:        name,
			Version:     version,
			Description: description,
		},
		scriptPath: scriptPath,
		health: &ExtensionHealth{
			Status: "initialized",
		},
		stopChan:  make(chan struct{}),
		errorChan: make(chan error, 1),
	}
}

// Start initializes and starts the Python extension
func (e *PythonExtension) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.process != nil {
		return fmt.Errorf("extension already running")
	}

	// Create Unix domain socket for IPC
	socketPath := filepath.Join(os.TempDir(), fmt.Sprintf("nessi-%s.sock", e.Name))
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return fmt.Errorf("failed to create socket: %w", err)
	}
	defer os.Remove(socketPath)

	// Start Python process
	cmd := exec.Command("python3", e.scriptPath)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("NESSI_SOCKET=%s", socketPath),
		fmt.Sprintf("NESSI_EXTENSION_NAME=%s", e.Name),
	)

	// Set up process pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	// Accept connection from Python process
	conn, err := listener.Accept()
	if err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to accept connection: %w", err)
	}

	e.process = cmd
	e.socket = conn

	// Start monitoring goroutines
	go e.monitorProcess()
	go e.handleOutput(stdout, stderr)

	return nil
}

// Stop gracefully stops the Python extension
func (e *PythonExtension) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.process == nil {
		return nil
	}

	// Send shutdown signal
	close(e.stopChan)

	// Give process time to clean up
	done := make(chan error, 1)
	go func() {
		done <- e.process.Wait()
	}()

	select {
	case err := <-done:
		e.cleanup()
		return err
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown fails
		e.process.Process.Signal(syscall.SIGKILL)
		e.cleanup()
		return fmt.Errorf("process did not stop gracefully")
	}
}

// Status returns the current status of the extension
func (e *PythonExtension) Status() string {
	e.health.mu.RLock()
	defer e.health.mu.RUnlock()
	return e.health.Status
}

// IsRunning checks if the extension process is running
func (e *PythonExtension) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.process != nil && e.process.Process != nil
}

// SendMessage sends a message to the Python process
func (e *PythonExtension) SendMessage(msg interface{}) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.socket == nil {
		return fmt.Errorf("not connected to Python process")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Add message length prefix
	header := make([]byte, 4)
	header[0] = byte(len(data) >> 24)
	header[1] = byte(len(data) >> 16)
	header[2] = byte(len(data) >> 8)
	header[3] = byte(len(data))

	if _, err := e.socket.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	if _, err := e.socket.Write(data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// ReceiveMessage receives a message from the Python process
func (e *PythonExtension) ReceiveMessage() (interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.socket == nil {
		return nil, fmt.Errorf("not connected to Python process")
	}

	// Read message length
	header := make([]byte, 4)
	if _, err := io.ReadFull(e.socket, header); err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	length := int(header[0])<<24 | int(header[1])<<16 | int(header[2])<<8 | int(header[3])
	data := make([]byte, length)

	if _, err := io.ReadFull(e.socket, data); err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var msg interface{}
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	return msg, nil
}

// ExecuteFunction calls a Python function with arguments
func (e *PythonExtension) ExecuteFunction(name string, args ...interface{}) (interface{}, error) {
	msg := map[string]interface{}{
		"type": "function_call",
		"name": name,
		"args": args,
	}

	if err := e.SendMessage(msg); err != nil {
		return nil, err
	}

	// Wait for response with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	responseChan := make(chan interface{}, 1)
	errorChan := make(chan error, 1)

	go func() {
		response, err := e.ReceiveMessage()
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- response
	}()

	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("function call timed out")
	}
}

// monitorProcess monitors the health of the Python process
func (e *PythonExtension) monitorProcess() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !e.IsRunning() {
				e.handleCrash()
				return
			}
			e.updateHealth()
		case <-e.stopChan:
			return
		}
	}
}

// handleOutput processes stdout and stderr from the Python process
func (e *PythonExtension) handleOutput(stdout, stderr io.ReadCloser) {
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			// Process stdout output
			log.Printf("[%s] %s", e.Name, scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			// Process stderr output
			log.Printf("[%s] ERROR: %s", e.Name, scanner.Text())
		}
	}()
}

// updateHealth updates the extension health status
func (e *PythonExtension) updateHealth() {
	e.health.mu.Lock()
	defer e.health.mu.Unlock()

	if !e.IsRunning() {
		e.health.Status = "stopped"
		return
	}

	// Send heartbeat message
	if err := e.SendMessage(map[string]string{"type": "heartbeat"}); err != nil {
		e.health.Status = "error"
		e.health.Error = err
		return
	}

	e.health.Status = "running"
	e.health.LastHeartbeat = time.Now()
}

// handleCrash handles unexpected process termination
func (e *PythonExtension) handleCrash() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.health.mu.Lock()
	e.health.Status = "crashed"
	e.health.Error = fmt.Errorf("process terminated unexpectedly")
	e.health.mu.Unlock()

	e.cleanup()
}

// cleanup releases resources
func (e *PythonExtension) cleanup() {
	if e.socket != nil {
		e.socket.Close()
		e.socket = nil
	}
	if e.process != nil {
		e.process = nil
	}
}