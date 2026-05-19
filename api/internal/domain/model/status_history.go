package model

import (
	"time"

	"github.com/google/uuid"
)

type StatusHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	Project   Project   `gorm:"foreignKey:ProjectID" json:"-"`
	OldStatus string    `gorm:"type:varchar(50)" json:"old_status"`
	NewStatus string    `gorm:"type:varchar(50);not null" json:"new_status"`
	ChangedBy uuid.UUID `gorm:"type:uuid;not null" json:"changed_by"`
	Changer   User      `gorm:"foreignKey:ChangedBy" json:"-"`
	ChangedAt time.Time `gorm:"default:now()" json:"changed_at"`
}

type TaskHistory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TaskID    uuid.UUID `gorm:"type:uuid;not null;index" json:"task_id"`
	Task      Task      `gorm:"foreignKey:TaskID" json:"-"`
	OldStatus string    `gorm:"type:varchar(50)" json:"old_status"`
	NewStatus string    `gorm:"type:varchar(50);not null" json:"new_status"`
	ChangedBy uuid.UUID `gorm:"type:uuid;not null" json:"changed_by"`
	Changer   User      `gorm:"foreignKey:ChangedBy" json:"-"`
	ChangedAt time.Time `gorm:"default:now()" json:"changed_at"`
}
