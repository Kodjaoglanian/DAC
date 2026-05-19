package repository

import (
	"dac/project-tracker/internal/domain/model"

	"github.com/google/uuid"
)

// ReportRepository defines the interface for report data access
type ReportRepository interface {
	GetDashboardData(userID uuid.UUID) (*DashboardData, error)
	GetProjectReport(projectID uuid.UUID) (*ProjectReportData, error)
	GetProjectsByStatus(userID uuid.UUID) (map[string]int64, error)
	GetProjectsByPriority(userID uuid.UUID) (map[string]int64, error)
	GetRecentProjects(userID uuid.UUID, limit int) ([]model.Project, error)
	GetStatusHistory(projectID uuid.UUID) ([]model.StatusHistory, error)
}

// DashboardData holds aggregated dashboard information
type DashboardData struct {
	TotalProjects       int64
	ActiveProjects      int64
	CompletedProjects   int64
	TotalTasks          int64
	TasksByStatus       map[string]int64
	ProjectsByStatus    map[string]int64
	ProjectsByPriority  map[string]int64
	RecentProjects      []model.Project
}

// ProjectReportData holds detailed project report information
type ProjectReportData struct {
	Project        *model.Project
	TotalTasks     int64
	CompletedTasks int64
	StatusHistory  []model.StatusHistory
}
