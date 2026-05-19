package service

import (
	"errors"
	"time"

	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/pkg/auth"
	"dac/project-tracker/pkg/hash"

	"github.com/google/uuid"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo         repository.UserRepository
	jwtSecret        string
	jwtExpirationHours int
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpirationHours int) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		jwtSecret:        jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}

// Register creates a new user and returns a JWT token
func (s *AuthService) Register(name, email, password string) (string, *model.User, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil && existing != nil {
		return "", nil, errors.New("email already registered")
	}

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return "", nil, err
	}

	user := &model.User{
		Name:      name,
		Email:     email,
		Password:  hashedPassword,
		Role:      "member",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", nil, err
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret, s.jwtExpirationHours)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !hash.CheckPassword(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(user.ID, user.Role, s.jwtSecret, s.jwtExpirationHours)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// GetUserByID returns a user by ID
func (s *AuthService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}
