package service

import (
	"errors"
	"time"

	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
)

// ProjectService handles project business logic
type ProjectService struct {
	projectRepo repository.ProjectRepository
	memberRepo  repository.MemberRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo repository.ProjectRepository, memberRepo repository.MemberRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
	}
}

// CreateProject creates a new project and adds the creator as owner
func (s *ProjectService) CreateProject(name, description, priority string, startDate, endDate *time.Time, createdBy uuid.UUID) (*model.Project, error) {
	if endDate != nil && startDate != nil && !endDate.After(*startDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	project := &model.Project{
		Name:        name,
		Description: description,
		Status:      model.ProjectStatusPlanning,
		Priority:    model.ProjectPriority(priority),
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}

	// Add creator as owner
	member := &model.ProjectMember{
		ProjectID: project.ID,
		UserID:    createdBy,
		Role:      model.MemberRoleOwner,
		JoinedAt:  time.Now(),
	}
	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return project, nil
}

// GetProjectByID returns a project by ID
func (s *ProjectService) GetProjectByID(id uuid.UUID) (*model.Project, error) {
	return s.projectRepo.FindByID(id)
}

// GetProjects returns a paginated list of projects for a user
func (s *ProjectService) GetProjects(userID uuid.UUID, filter repository.ProjectFilter) ([]model.Project, int64, error) {
	filter.UserID = userID
	return s.projectRepo.FindByUser(userID, filter.Page, filter.PageSize)
}

// UpdateProject updates a project
func (s *ProjectService) UpdateProject(projectID uuid.UUID, name, description, status, priority string, startDate, endDate *time.Time, userID uuid.UUID) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	role, err := s.memberRepo.GetRole(projectID, userID)
	if err != nil || (role != string(model.MemberRoleOwner) && role != string(model.MemberRoleManager)) {
		return nil, errors.New("insufficient permissions")
	}

	if endDate != nil && startDate != nil && !endDate.After(*startDate) {
		return nil, errors.New("end_date must be after start_date")
	}

	oldStatus := project.Status

	if name != "" {
		project.Name = name
	}
	if description != "" {
		project.Description = description
	}
	if status != "" {
		project.Status = model.ProjectStatus(status)
	}
	if priority != "" {
		project.Priority = model.ProjectPriority(priority)
	}
	if startDate != nil {
		project.StartDate = startDate
	}
	if endDate != nil {
		project.EndDate = endDate
	}
	project.UpdatedAt = time.Now()

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	// Record status change
	if status != "" && oldStatus != model.ProjectStatus(status) {
		history := &model.StatusHistory{
			ProjectID: projectID,
			OldStatus: string(oldStatus),
			NewStatus: status,
			ChangedBy: userID,
			ChangedAt: time.Now(),
		}
		s.projectRepo.AddStatusHistory(history)
	}

	return project, nil
}

// UpdateProjectStatus updates only the project status
func (s *ProjectService) UpdateProjectStatus(projectID uuid.UUID, status string, userID uuid.UUID) (*model.Project, error) {
	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project not found")
	}

	role, err := s.memberRepo.GetRole(projectID, userID)
	if err != nil || (role != string(model.MemberRoleOwner) && role != string(model.MemberRoleManager)) {
		return nil, errors.New("insufficient permissions")
	}

	oldStatus := project.Status
	project.Status = model.ProjectStatus(status)
	project.UpdatedAt = time.Now()

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}

	history := &model.StatusHistory{
		ProjectID: projectID,
		OldStatus: string(oldStatus),
		NewStatus: status,
		ChangedBy: userID,
		ChangedAt: time.Now(),
	}
	s.projectRepo.AddStatusHistory(history)

	return project, nil
}

// DeleteProject performs a soft delete
func (s *ProjectService) DeleteProject(projectID, userID uuid.UUID) error {
	role, err := s.memberRepo.GetRole(projectID, userID)
	if err != nil || role != string(model.MemberRoleOwner) {
		return errors.New("only owner can delete project")
	}

	return s.projectRepo.Delete(projectID)
}
