package rbac

import (
	"errors"
	"fmt"
	"sync"
)

// Role represents a user role
type Role string

const (
	// AdminRole has full access to all features
	AdminRole Role = "admin"
	// UserRole has access to most features but not administrative ones
	UserRole Role = "user"
	// ViewerRole has read-only access
	ViewerRole Role = "viewer"
	// CustomRole is a user-defined role
	CustomRole Role = "custom"
)

// Permission represents a permission to perform an action
type Permission string

const (
	// ReadPermission allows reading data
	ReadPermission Permission = "read"
	// WritePermission allows writing data
	WritePermission Permission = "write"
	// ExecutePermission allows executing operations
	ExecutePermission Permission = "execute"
	// AdminPermission allows administrative operations
	AdminPermission Permission = "admin"
)

// Resource represents a resource that can be accessed
type Resource string

const (
	// TableResource represents a table resource
	TableResource Resource = "table"
	// RuleResource represents a rule resource
	RuleResource Resource = "rule"
	// UserResource represents a user resource
	UserResource Resource = "user"
	// WebhookResource represents a webhook resource
	WebhookResource Resource = "webhook"
	// PluginResource represents a plugin resource
	PluginResource Resource = "plugin"
	// ConfigResource represents a configuration resource
	ConfigResource Resource = "config"
	// ReportResource represents a report resource
	ReportResource Resource = "report"
)

// AccessControl represents an access control entry
type AccessControl struct {
	Role       Role
	Resource   Resource
	Permission Permission
}

// RBACManager manages role-based access control
type RBACManager struct {
	accessControls []AccessControl
	userRoles      map[string]Role
	teamRoles      map[string]Role
	customRoles    map[string][]AccessControl
	mu             sync.RWMutex
}

// NewRBACManager creates a new RBAC manager
func NewRBACManager() *RBACManager {
	manager := &RBACManager{
		accessControls: []AccessControl{},
		userRoles:      make(map[string]Role),
		teamRoles:      make(map[string]Role),
		customRoles:    make(map[string][]AccessControl),
	}

	// Initialize default access controls
	manager.initializeDefaultAccessControls()

	return manager
}

// initializeDefaultAccessControls initializes default access controls
func (m *RBACManager) initializeDefaultAccessControls() {
	// Admin role has all permissions on all resources
	resources := []Resource{TableResource, RuleResource, UserResource, WebhookResource, PluginResource, ConfigResource, ReportResource}
	permissions := []Permission{ReadPermission, WritePermission, ExecutePermission, AdminPermission}

	for _, resource := range resources {
		for _, permission := range permissions {
			m.accessControls = append(m.accessControls, AccessControl{
				Role:       AdminRole,
				Resource:   resource,
				Permission: permission,
			})
		}
	}

	// User role has read, write, and execute permissions on most resources, but not admin
	for _, resource := range resources {
		if resource == UserResource || resource == ConfigResource {
			// Users can only read these resources
			m.accessControls = append(m.accessControls, AccessControl{
				Role:       UserRole,
				Resource:   resource,
				Permission: ReadPermission,
			})
		} else {
			// Users can read, write, and execute on other resources
			for _, permission := range []Permission{ReadPermission, WritePermission, ExecutePermission} {
				m.accessControls = append(m.accessControls, AccessControl{
					Role:       UserRole,
					Resource:   resource,
					Permission: permission,
				})
			}
		}
	}

	// Viewer role has only read permissions
	for _, resource := range resources {
		m.accessControls = append(m.accessControls, AccessControl{
			Role:       ViewerRole,
			Resource:   resource,
			Permission: ReadPermission,
		})
	}
}

// SetUserRole sets a user's role
func (m *RBACManager) SetUserRole(userID string, role Role) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.userRoles[userID] = role
}

// GetUserRole gets a user's role
func (m *RBACManager) GetUserRole(userID string) (Role, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	role, exists := m.userRoles[userID]
	return role, exists
}

// SetTeamRole sets a team's role
func (m *RBACManager) SetTeamRole(teamID string, role Role) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.teamRoles[teamID] = role
}

// GetTeamRole gets a team's role
func (m *RBACManager) GetTeamRole(teamID string) (Role, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	role, exists := m.teamRoles[teamID]
	return role, exists
}

// CreateCustomRole creates a custom role
func (m *RBACManager) CreateCustomRole(roleName string, accessControls []AccessControl) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.customRoles[roleName]; exists {
		return fmt.Errorf("custom role %s already exists", roleName)
	}

	m.customRoles[roleName] = accessControls
	return nil
}

// DeleteCustomRole deletes a custom role
func (m *RBACManager) DeleteCustomRole(roleName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.customRoles[roleName]; !exists {
		return fmt.Errorf("custom role %s does not exist", roleName)
	}

	delete(m.customRoles, roleName)
	return nil
}

// CheckAccess checks if a user has permission to access a resource
func (m *RBACManager) CheckAccess(userID string, resource Resource, permission Permission) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Get the user's role
	role, exists := m.userRoles[userID]
	if !exists {
		return false, errors.New("user role not found")
	}

	// Check if the user has the required permission
	if role == CustomRole {
		// For custom roles, check the custom access controls
		customRoleName := fmt.Sprintf("custom:%s", userID)
		accessControls, exists := m.customRoles[customRoleName]
		if !exists {
			return false, fmt.Errorf("custom role for user %s not found", userID)
		}

		for _, ac := range accessControls {
			if ac.Resource == resource && ac.Permission == permission {
				return true, nil
			}
		}
		return false, nil
	}

	// For built-in roles, check the default access controls
	for _, ac := range m.accessControls {
		if ac.Role == role && ac.Resource == resource && ac.Permission == permission {
			return true, nil
		}
	}

	return false, nil
}

// CheckTeamAccess checks if a team has permission to access a resource
func (m *RBACManager) CheckTeamAccess(teamID string, resource Resource, permission Permission) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Get the team's role
	role, exists := m.teamRoles[teamID]
	if !exists {
		return false, errors.New("team role not found")
	}

	// Check if the team has the required permission
	if role == CustomRole {
		// For custom roles, check the custom access controls
		customRoleName := fmt.Sprintf("team:%s", teamID)
		accessControls, exists := m.customRoles[customRoleName]
		if !exists {
			return false, fmt.Errorf("custom role for team %s not found", teamID)
		}

		for _, ac := range accessControls {
			if ac.Resource == resource && ac.Permission == permission {
				return true, nil
			}
		}
		return false, nil
	}

	// For built-in roles, check the default access controls
	for _, ac := range m.accessControls {
		if ac.Role == role && ac.Resource == resource && ac.Permission == permission {
			return true, nil
		}
	}

	return false, nil
}

// ListUsers lists all users and their roles
func (m *RBACManager) ListUsers() map[string]Role {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	users := make(map[string]Role)
	for k, v := range m.userRoles {
		users[k] = v
	}

	return users
}

// ListTeams lists all teams and their roles
func (m *RBACManager) ListTeams() map[string]Role {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	teams := make(map[string]Role)
	for k, v := range m.teamRoles {
		teams[k] = v
	}

	return teams
}

// ListCustomRoles lists all custom roles
func (m *RBACManager) ListCustomRoles() map[string][]AccessControl {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	roles := make(map[string][]AccessControl)
	for k, v := range m.customRoles {
		roles[k] = v
	}

	return roles
}
