package service

import (
	"context"
	"encoding/json"
	"time"

	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/internal/infrastructure/database"

	"github.com/google/uuid"
)

// ReportService handles report business logic
type ReportService struct {
	reportRepo  repository.ReportRepository
	projectRepo repository.ProjectRepository
	taskRepo    repository.TaskRepository
	memberRepo  repository.MemberRepository
	redis       *database.RedisClient
}

// NewReportService creates a new report service
func NewReportService(reportRepo repository.ReportRepository, projectRepo repository.ProjectRepository, taskRepo repository.TaskRepository, memberRepo repository.MemberRepository, redis *database.RedisClient) *ReportService {
	return &ReportService{
		reportRepo:  reportRepo,
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
		memberRepo:  memberRepo,
		redis:       redis,
	}
}

// GetDashboard returns consolidated dashboard data
func (s *ReportService) GetDashboard(userID uuid.UUID) (*repository.DashboardData, error) {
	cacheKey := "dashboard:" + userID.String()

	// Try cache
	if s.redis != nil {
		cached, err := s.redis.Get(context.Background(), cacheKey)
		if err == nil && cached != "" {
			var data repository.DashboardData
			if err := json.Unmarshal([]byte(cached), &data); err == nil {
				return &data, nil
			}
		}
	}

	data, err := s.reportRepo.GetDashboardData(userID)
	if err != nil {
		return nil, err
	}

	// Additional aggregations
	tasksByStatus, err := s.taskRepo.CountByStatus(userID)
	if err != nil {
		return nil, err
	}
	data.TasksByStatus = tasksByStatus

	projectsByStatus, err := s.reportRepo.GetProjectsByStatus(userID)
	if err != nil {
		return nil, err
	}
	data.ProjectsByStatus = projectsByStatus

	projectsByPriority, err := s.reportRepo.GetProjectsByPriority(userID)
	if err != nil {
		return nil, err
	}
	data.ProjectsByPriority = projectsByPriority

	recentProjects, err := s.reportRepo.GetRecentProjects(userID, 5)
	if err != nil {
		return nil, err
	}
	data.RecentProjects = recentProjects

	// Cache result
	if s.redis != nil {
		jsonData, _ := json.Marshal(data)
		s.redis.Set(context.Background(), cacheKey, jsonData, 5*time.Minute)
	}

	return data, nil
}

// GetProjectReport returns detailed project report
func (s *ReportService) GetProjectReport(projectID uuid.UUID) (*repository.ProjectReportData, error) {
	return s.reportRepo.GetProjectReport(projectID)
}

// GetProjectsByStatus returns projects grouped by status
func (s *ReportService) GetProjectsByStatus(userID uuid.UUID) (map[string]int64, error) {
	return s.reportRepo.GetProjectsByStatus(userID)
}

// GetTasksByStatus returns tasks grouped by status
func (s *ReportService) GetTasksByStatus(userID uuid.UUID) (map[string]int64, error) {
	return s.taskRepo.CountByStatus(userID)
}
