package iam

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrWorkspaceForbidden = errors.New("workspace access denied")
)

const (
	WorkspaceRoleOwner  = "owner"
	WorkspaceRoleAdmin  = "admin"
	WorkspaceRoleEditor = "editor"
	WorkspaceRoleMember = "member"
	WorkspaceRoleViewer = "viewer"
)

type Store interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUser(ctx context.Context, id string) (*User, error)
	CreateWorkspace(ctx context.Context, workspace *Workspace) error
	GetWorkspaceList(ctx context.Context, userID string) ([]*Workspace, error)
	AddMembership(ctx context.Context, membership *Membership) error
	HasMembership(ctx context.Context, userID, workspaceID string) (bool, error)
	GetMembership(ctx context.Context, userID, workspaceID string) (*Membership, error)
}

type Service struct {
	store  Store
	secret []byte
	issuer string
}

func New(store Store, secret, issuer string) *Service {
	return &Service{
		store:  store,
		secret: []byte(secret),
		issuer: issuer,
	}
}

type LoginResult struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func newID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(raw[:])
}

func validWorkspaceRole(role string) bool {
	switch role {
	case
		WorkspaceRoleOwner,
		WorkspaceRoleAdmin,
		WorkspaceRoleEditor,
		WorkspaceRoleMember,
		WorkspaceRoleViewer:
		return true
	default:
		return false
	}
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.store.GetUserByEmail(ctx, strings.TrimSpace(email))
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return nil, ErrInvalidCredentials
	}
	now := time.Now()
	claims := claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(8 * time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}
	return &LoginResult{Token: token, User: user}, nil
}

func (s *Service) Register(ctx context.Context, email, password, name string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return nil, fmt.Errorf("invalid email")
	}

	if len(password) < 8 {
		return nil, fmt.Errorf("password must contain at least 8 characters")
	}

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// 存储 User
	user := &User{
		ID:           newID(),
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.store.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// 存储 Workspace
	workspace := &Workspace{
		ID:        newID(),
		Name:      name + " Workspace",
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.CreateWorkspace(ctx, workspace); err != nil {
		return nil, fmt.Errorf("create default workspace: %w", err)
	}

	// 存储 Membership
	if err := s.store.AddMembership(ctx, &Membership{
		UserID:      user.ID,
		WorkspaceID: workspace.ID,
		Role:        WorkspaceRoleOwner,
		CreatedAt:   time.Now().UTC(),
	}); err != nil {
		return nil, fmt.Errorf("create default workspace membership: %w", err)
	}

	return user, nil
}

func (s *Service) ParseToken(raw string) (string, error) {
	parsed := &claims{}
	token, err := jwt.ParseWithClaims(raw, parsed, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method %s", token.Method.Alg())
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid || parsed.UserID == "" {
		return "", ErrInvalidCredentials
	}
	return parsed.UserID, nil
}

func (s *Service) CheckAccess(ctx context.Context, userID, workspaceID string) (string, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(workspaceID) == "" {
		return "", ErrWorkspaceForbidden
	}
	membership, err := s.store.GetMembership(ctx, userID, workspaceID)
	if err != nil || !validWorkspaceRole(membership.Role) {
		return "", ErrWorkspaceForbidden
	}
	return membership.Role, nil
}

func (s *Service) WorkspaceRole(ctx context.Context, userID, workspaceID string) (string, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(workspaceID) == "" {
		return "", ErrWorkspaceForbidden
	}
	membership, err := s.store.GetMembership(ctx, userID, workspaceID)
	if err != nil || !validWorkspaceRole(membership.Role) {
		return "", ErrWorkspaceForbidden
	}
	return membership.Role, nil
}

func (s *Service) GetWorkspaceList(ctx context.Context, userID string) ([]*Workspace, error) {
	return s.store.GetWorkspaceList(ctx, userID)
}
