package rbac

import (
	"fmt"
	"sync"
)

// Team represents a team of users
type Team struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
	Role        Role     `json:"role"`
}

// TeamManager manages teams and their permissions
type TeamManager struct {
	teams map[string]*Team
	mu    sync.RWMutex
}

// NewTeamManager creates a new team manager
func NewTeamManager() *TeamManager {
	return &TeamManager{
		teams: make(map[string]*Team),
	}
}

// CreateTeam creates a new team
func (m *TeamManager) CreateTeam(id, name, description string, role Role) (*Team, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.teams[id]; exists {
		return nil, fmt.Errorf("team with ID %s already exists", id)
	}

	team := &Team{
		ID:          id,
		Name:        name,
		Description: description,
		Members:     []string{},
		Role:        role,
	}

	m.teams[id] = team
	return team, nil
}

// GetTeam gets a team by ID
func (m *TeamManager) GetTeam(id string) (*Team, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	team, exists := m.teams[id]
	if !exists {
		return nil, fmt.Errorf("team with ID %s not found", id)
	}

	return team, nil
}

// UpdateTeam updates a team
func (m *TeamManager) UpdateTeam(id, name, description string, role Role) (*Team, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	team, exists := m.teams[id]
	if !exists {
		return nil, fmt.Errorf("team with ID %s not found", id)
	}

	team.Name = name
	team.Description = description
	team.Role = role

	return team, nil
}

// DeleteTeam deletes a team
func (m *TeamManager) DeleteTeam(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.teams[id]; !exists {
		return fmt.Errorf("team with ID %s not found", id)
	}

	delete(m.teams, id)
	return nil
}

// AddMember adds a user to a team
func (m *TeamManager) AddMember(teamID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	team, exists := m.teams[teamID]
	if !exists {
		return fmt.Errorf("team with ID %s not found", teamID)
	}

	// Check if user is already a member
	for _, member := range team.Members {
		if member == userID {
			return fmt.Errorf("user %s is already a member of team %s", userID, teamID)
		}
	}

	team.Members = append(team.Members, userID)
	return nil
}

// RemoveMember removes a user from a team
func (m *TeamManager) RemoveMember(teamID, userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	team, exists := m.teams[teamID]
	if !exists {
		return fmt.Errorf("team with ID %s not found", teamID)
	}

	// Find and remove the user
	for i, member := range team.Members {
		if member == userID {
			team.Members = append(team.Members[:i], team.Members[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user %s is not a member of team %s", userID, teamID)
}

// GetUserTeams gets all teams that a user is a member of
func (m *TeamManager) GetUserTeams(userID string) []*Team {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var userTeams []*Team
	for _, team := range m.teams {
		for _, member := range team.Members {
			if member == userID {
				userTeams = append(userTeams, team)
				break
			}
		}
	}

	return userTeams
}

// ListTeams lists all teams
func (m *TeamManager) ListTeams() []*Team {
	m.mu.RLock()
	defer m.mu.RUnlock()

	teams := make([]*Team, 0, len(m.teams))
	for _, team := range m.teams {
		teams = append(teams, team)
	}

	return teams
}

// IsMember checks if a user is a member of a team
func (m *TeamManager) IsMember(teamID, userID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	team, exists := m.teams[teamID]
	if !exists {
		return false, fmt.Errorf("team with ID %s not found", teamID)
	}

	for _, member := range team.Members {
		if member == userID {
			return true, nil
		}
	}

	return false, nil
}

// GetTeamRole gets the role of a team
func (m *TeamManager) GetTeamRole(teamID string) (Role, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	team, exists := m.teams[teamID]
	if !exists {
		return "", fmt.Errorf("team with ID %s not found", teamID)
	}

	return team.Role, nil
}

// SetTeamRole sets the role of a team
func (m *TeamManager) SetTeamRole(teamID string, role Role) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	team, exists := m.teams[teamID]
	if !exists {
		return fmt.Errorf("team with ID %s not found", teamID)
	}

	team.Role = role
	return nil
}
