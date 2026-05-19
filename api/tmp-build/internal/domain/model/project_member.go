package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectMemberRole string

const (
	MemberRoleOwner   ProjectMemberRole = "owner"
	MemberRoleManager ProjectMemberRole = "manager"
	MemberRoleMember  ProjectMemberRole = "member"
)

type ProjectMember struct {
	ID        uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ProjectID uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:idx_project_user" json:"project_id"`
	Project   Project           `gorm:"foreignKey:ProjectID" json:"-"`
	UserID    uuid.UUID         `gorm:"type:uuid;not null;uniqueIndex:idx_project_user" json:"user_id"`
	User      User              `gorm:"foreignKey:UserID" json:"user"`
	Role      ProjectMemberRole `gorm:"type:varchar(50);not null;default:'member'" json:"role" validate:"oneof=owner manager member"`
	JoinedAt  time.Time         `gorm:"default:now()" json:"joined_at"`
}

func (pm *ProjectMember) BeforeCreate(tx *gorm.DB) error {
	pm.ID = uuid.New()
	return nil
}
