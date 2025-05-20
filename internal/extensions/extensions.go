package extensions

import "fmt"

// Extension represents a Nessi extension
type Extension struct {
	Name        string
	Version     string
	Description string
	Enabled     bool
}

// Manager manages Nessi extensions
type Manager struct {
	extensions map[string]*Extension
}

// NewManager creates a new extension manager
func NewManager() *Manager {
	return &Manager{
		extensions: make(map[string]*Extension),
	}
}

// RegisterExtension registers a new extension
func (m *Manager) RegisterExtension(ext *Extension) error {
	if ext == nil {
		return fmt.Errorf("extension cannot be nil")
	}
	if ext.Name == "" {
		return fmt.Errorf("extension name cannot be empty")
	}
	m.extensions[ext.Name] = ext
	return nil
}

// GetExtension returns an extension by name
func (m *Manager) GetExtension(name string) (*Extension, error) {
	ext, ok := m.extensions[name]
	if !ok {
		return nil, fmt.Errorf("extension %s not found", name)
	}
	return ext, nil
}

// ListExtensions returns all registered extensions
func (m *Manager) ListExtensions() []*Extension {
	var exts []*Extension
	for _, ext := range m.extensions {
		exts = append(exts, ext)
	}
	return exts
}

// EnableExtension enables an extension
func (m *Manager) EnableExtension(name string) error {
	ext, err := m.GetExtension(name)
	if err != nil {
		return err
	}
	ext.Enabled = true
	return nil
}

// DisableExtension disables an extension
func (m *Manager) DisableExtension(name string) error {
	ext, err := m.GetExtension(name)
	if err != nil {
		return err
	}
	ext.Enabled = false
	return nil
}
