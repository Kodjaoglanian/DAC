package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectStatus string

const (
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusCancelled  ProjectStatus = "cancelled"
)

type ProjectPriority string

const (
	PriorityLow      ProjectPriority = "low"
	PriorityMedium   ProjectPriority = "medium"
	PriorityHigh     ProjectPriority = "high"
	PriorityCritical ProjectPriority = "critical"
)

type Project struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string          `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=3,max=255"`
	Description string          `gorm:"type:text;not null" json:"description" validate:"required,min=10"`
	Status      ProjectStatus   `gorm:"type:varchar(50);not null;default:'planning'" json:"status" validate:"oneof=planning in_progress completed cancelled"`
	Priority    ProjectPriority `gorm:"type:varchar(50);not null;default:'medium'" json:"priority" validate:"oneof=low medium high critical"`
	StartDate   *time.Time      `gorm:"type:date" json:"start_date"`
	EndDate     *time.Time      `gorm:"type:date" json:"end_date"`
	CreatedBy   uuid.UUID       `gorm:"type:uuid;not null" json:"created_by"`
	Creator     User            `gorm:"foreignKey:CreatedBy" json:"creator"`
	Members     []ProjectMember `gorm:"foreignKey:ProjectID" json:"members"`
	Tasks       []Task          `gorm:"foreignKey:ProjectID" json:"tasks"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `gorm:"index" json:"-"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}
