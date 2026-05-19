package service

import (
	"errors"
	"time"

	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
)

// UserService handles user business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// GetUserByID returns a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.userRepo.FindByID(id)
}

// GetAllUsers returns a paginated list of all users
func (s *UserService) GetAllUsers(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.FindAll(page, pageSize)
}

// UpdateUser updates user data
func (s *UserService) UpdateUser(id uuid.UUID, name, avatarURL string) (*model.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if name != "" {
		user.Name = name
	}
	if avatarURL != "" {
		user.AvatarURL = avatarURL
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser performs a soft delete on a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}
