package repository

import (
	"dac/project-tracker/internal/domain/model"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAll(page, pageSize int) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
}
