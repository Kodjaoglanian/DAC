package repository

import (
	"dac/project-tracker/internal/domain/model"

	"github.com/google/uuid"
)

// MemberRepository defines the interface for project member data access
type MemberRepository interface {
	Create(member *model.ProjectMember) error
	FindByProject(projectID uuid.UUID) ([]model.ProjectMember, error)
	FindByProjectAndUser(projectID, userID uuid.UUID) (*model.ProjectMember, error)
	Update(member *model.ProjectMember) error
	Delete(projectID, userID uuid.UUID) error
	IsMember(projectID, userID uuid.UUID) (bool, error)
	GetRole(projectID, userID uuid.UUID) (string, error)
	CountOwners(projectID uuid.UUID) (int64, error)
}
