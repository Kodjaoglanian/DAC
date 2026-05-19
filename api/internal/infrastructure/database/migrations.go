package database

import (
	"dac/project-tracker/internal/domain/model"

	"gorm.io/gorm"
)

// RunMigrations executes all GORM auto-migrations
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Task{},
		&model.ProjectMember{},
		&model.StatusHistory{},
		&model.TaskHistory{},
	)
}
