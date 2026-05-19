package repository

import (
	"dac/project-tracker/internal/domain/model"

	"github.com/google/uuid"
)

// TaskRepository defines the interface for task data access
type TaskRepository interface {
	Create(task *model.Task) error
	FindByID(id uuid.UUID) (*model.Task, error)
	FindByProject(projectID uuid.UUID, filter TaskFilter) ([]model.Task, int64, error)
	Update(task *model.Task) error
	Delete(id uuid.UUID) error
	AddTaskHistory(history *model.TaskHistory) error
	CountByProject(projectID uuid.UUID) (int64, error)
	CountCompletedByProject(projectID uuid.UUID) (int64, error)
	CountByStatus(userID uuid.UUID) (map[string]int64, error)
}

// TaskFilter holds filtering options for tasks
type TaskFilter struct {
	Status     string
	Priority   string
	AssignedTo string
	Search     string
	Sort       string
	Order      string
}
