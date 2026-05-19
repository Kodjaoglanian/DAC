package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusReview     TaskStatus = "review"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title       string          `gorm:"type:varchar(255);not null" json:"title" validate:"required,min=3,max=255"`
	Description string          `gorm:"type:text" json:"description"`
	Status      TaskStatus      `gorm:"type:varchar(50);not null;default:'todo'" json:"status" validate:"oneof=todo in_progress review done"`
	Priority    ProjectPriority `gorm:"type:varchar(50);not null;default:'medium'" json:"priority" validate:"oneof=low medium high critical"`
	DueDate     *time.Time      `gorm:"type:date" json:"due_date"`
	ProjectID   uuid.UUID       `gorm:"type:uuid;not null;index" json:"project_id"`
	Project     Project         `gorm:"foreignKey:ProjectID" json:"-"`
	AssignedTo  *uuid.UUID      `gorm:"type:uuid;index" json:"assigned_to"`
	Assignee    *User           `gorm:"foreignKey:AssignedTo" json:"assignee"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	return nil
}
