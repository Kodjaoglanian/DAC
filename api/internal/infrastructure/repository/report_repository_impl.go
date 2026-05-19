package repository

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ReportRepositoryImpl implements repository.ReportRepository
type ReportRepositoryImpl struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *gorm.DB) repository.ReportRepository {
	return &ReportRepositoryImpl{db: db}
}

func (r *ReportRepositoryImpl) GetDashboardData(userID uuid.UUID) (*repository.DashboardData, error) {
	var data repository.DashboardData

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)

	// Total projects
	r.db.Model(&model.Project{}).
		Where("id IN (?) OR created_by = ?", subQuery, userID).
		Count(&data.TotalProjects)

	// Active projects
	r.db.Model(&model.Project{}).
		Where("(id IN (?) OR created_by = ?) AND status = ?", subQuery, userID, model.ProjectStatusInProgress).
		Count(&data.ActiveProjects)

	// Completed projects
	r.db.Model(&model.Project{}).
		Where("(id IN (?) OR created_by = ?) AND status = ?", subQuery, userID, model.ProjectStatusCompleted).
		Count(&data.CompletedProjects)

	// Total tasks
	r.db.Model(&model.Task{}).
		Where("project_id IN (?) OR project_id IN (?)", subQuery, r.db.Table("projects").Select("id").Where("created_by = ?", userID)).
		Count(&data.TotalTasks)

	return &data, nil
}

func (r *ReportRepositoryImpl) GetProjectReport(projectID uuid.UUID) (*repository.ProjectReportData, error) {
	var data repository.ProjectReportData

	var project model.Project
	if err := r.db.First(&project, "id = ?", projectID).Error; err != nil {
		return nil, err
	}
	data.Project = &project

	r.db.Model(&model.Task{}).Where("project_id = ?", projectID).Count(&data.TotalTasks)
	r.db.Model(&model.Task{}).Where("project_id = ? AND status = ?", projectID, model.TaskStatusDone).Count(&data.CompletedTasks)

	var history []model.StatusHistory
	if err := r.db.Where("project_id = ?", projectID).Order("changed_at desc").Find(&history).Error; err != nil {
		return nil, err
	}
	data.StatusHistory = history

	return &data, nil
}

func (r *ReportRepositoryImpl) GetProjectsByStatus(userID uuid.UUID) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)
	query := r.db.Model(&model.Project{}).Select("status, COUNT(*) as count").
		Where("id IN (?) OR created_by = ?", subQuery, userID).
		Group("status")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, res := range results {
		statusMap[res.Status] = res.Count
	}

	return statusMap, nil
}

func (r *ReportRepositoryImpl) GetProjectsByPriority(userID uuid.UUID) (map[string]int64, error) {
	var results []struct {
		Priority string
		Count    int64
	}

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)
	query := r.db.Model(&model.Project{}).Select("priority, COUNT(*) as count").
		Where("id IN (?) OR created_by = ?", subQuery, userID).
		Group("priority")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	priorityMap := make(map[string]int64)
	for _, res := range results {
		priorityMap[res.Priority] = res.Count
	}

	return priorityMap, nil
}

func (r *ReportRepositoryImpl) GetRecentProjects(userID uuid.UUID, limit int) ([]model.Project, error) {
	var projects []model.Project

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)
	if err := r.db.Preload("Creator").
		Where("id IN (?) OR created_by = ?", subQuery, userID).
		Order("created_at desc").
		Limit(limit).
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ReportRepositoryImpl) GetStatusHistory(projectID uuid.UUID) ([]model.StatusHistory, error) {
	var history []model.StatusHistory
	if err := r.db.Where("project_id = ?", projectID).Order("changed_at desc").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}
