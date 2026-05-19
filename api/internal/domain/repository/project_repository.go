package repository

import (
	"dac/project-tracker/internal/domain/model"

	"github.com/google/uuid"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	Create(project *model.Project) error
	FindByID(id uuid.UUID) (*model.Project, error)
	FindAll(filter ProjectFilter) ([]model.Project, int64, error)
	FindByUser(userID uuid.UUID, page, pageSize int) ([]model.Project, int64, error)
	Update(project *model.Project) error
	Delete(id uuid.UUID) error
	AddStatusHistory(history *model.StatusHistory) error
}

// ProjectFilter holds filtering options for projects
type ProjectFilter struct {
	Page       int
	PageSize   int
	Status     string
	Priority   string
	Search     string
	CreatedBy  string
	Sort       string
	Order      string
	UserID     uuid.UUID
}
