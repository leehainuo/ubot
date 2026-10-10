package iam

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrMembershipNotFound = errors.New("workspace membership not found")
)

type MemoryStore struct {
	mu          sync.RWMutex
	byID        map[string]*User
	byEmail     map[string]string
	workspaces  map[string]*Workspace
	memberships map[string]map[string]Membership
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		byID:        make(map[string]*User),
		byEmail:     make(map[string]string),
		workspaces:  make(map[string]*Workspace),
		memberships: make(map[string]map[string]Membership),
	}
}

func (m *MemoryStore) CreateUser(_ context.Context, user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := strings.ToLower(user.Email)
	if _, exists := m.byEmail[key]; exists {
		return errors.New("email already exists")
	}
	clone := *user
	m.byID[user.ID] = &clone
	m.byEmail[key] = user.ID
	return nil
}

func (m *MemoryStore) GetUserByEmail(_ context.Context, email string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	id, ok := m.byEmail[strings.ToLower(email)]
	if !ok {
		return nil, ErrUserNotFound
	}
	clone := *m.byID[id]
	return &clone, nil
}

func (m *MemoryStore) GetUser(_ context.Context, id string) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	user, ok := m.byID[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	clone := *user
	return &clone, nil
}

func (m *MemoryStore) CreateWorkspace(_ context.Context, workspace *Workspace) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.workspaces[workspace.ID]; exists {
		return errors.New("workspace already exists")
	}
	clone := *workspace
	m.workspaces[workspace.ID] = &clone
	return nil
}

func (m *MemoryStore) GetWorkspaceList(_ context.Context, userID string) ([]*Workspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]*Workspace, 0, len(m.memberships[userID]))
	for workspaceID := range m.memberships[userID] {
		workspace := m.workspaces[workspaceID]
		if workspace == nil {
			continue
		}
		clone := *workspace
		res = append(res, &clone)
	}
	return res, nil
}

func (m *MemoryStore) AddMembership(_ context.Context, membership *Membership) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byID[membership.UserID] == nil || m.workspaces[membership.WorkspaceID] == nil {
		return errors.New("membership references an unknown user or workspace")
	}
	if m.memberships[membership.UserID] == nil {
		m.memberships[membership.UserID] = make(map[string]Membership)
	}
	m.memberships[membership.UserID][membership.WorkspaceID] = *membership
	return nil
}

func (m *MemoryStore) HasMembership(_ context.Context, userID string, workspaceID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.memberships[userID][workspaceID]
	return ok, nil
}

func (m *MemoryStore) GetMembership(_ context.Context, userID string, workspaceID string) (*Membership, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	membership, ok := m.memberships[userID][workspaceID]
	if !ok {
		return nil, ErrMembershipNotFound
	}
	clone := membership
	return &clone, nil
}
