package service

import (
	"errors"
	"time"

	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
)

// TaskService handles task business logic
type TaskService struct {
	taskRepo   repository.TaskRepository
	memberRepo repository.MemberRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repository.TaskRepository, memberRepo repository.MemberRepository) *TaskService {
	return &TaskService{
		taskRepo:   taskRepo,
		memberRepo: memberRepo,
	}
}

// CreateTask creates a new task in a project
func (s *TaskService) CreateTask(projectID uuid.UUID, title, description, priority string, dueDate *time.Time, assignedTo *uuid.UUID, userID uuid.UUID) (*model.Task, error) {
	role, err := s.memberRepo.GetRole(projectID, userID)
	if err != nil || (role != string(model.MemberRoleOwner) && role != string(model.MemberRoleManager)) {
		return nil, errors.New("insufficient permissions")
	}

	if assignedTo != nil {
		isMember, err := s.memberRepo.IsMember(projectID, *assignedTo)
		if err != nil || !isMember {
			return nil, errors.New("assigned user is not a project member")
		}
	}

	task := &model.Task{
		Title:       title,
		Description: description,
		Status:      model.TaskStatusTodo,
		Priority:    model.ProjectPriority(priority),
		DueDate:     dueDate,
		ProjectID:   projectID,
		AssignedTo:  assignedTo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTaskByID returns a task by ID
func (s *TaskService) GetTaskByID(id uuid.UUID) (*model.Task, error) {
	return s.taskRepo.FindByID(id)
}

// GetTasksByProject returns tasks for a project
func (s *TaskService) GetTasksByProject(projectID uuid.UUID, filter repository.TaskFilter) ([]model.Task, int64, error) {
	return s.taskRepo.FindByProject(projectID, filter)
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(taskID uuid.UUID, title, description, status, priority string, dueDate *time.Time, assignedTo *uuid.UUID, userID uuid.UUID) (*model.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	role, err := s.memberRepo.GetRole(task.ProjectID, userID)
	if err != nil {
		return nil, err
	}

	isAssignee := task.AssignedTo != nil && *task.AssignedTo == userID
	if role != string(model.MemberRoleOwner) && role != string(model.MemberRoleManager) && !isAssignee {
		return nil, errors.New("insufficient permissions")
	}

	oldStatus := task.Status

	if title != "" {
		task.Title = title
	}
	if description != "" {
		task.Description = description
	}
	if status != "" {
		task.Status = model.TaskStatus(status)
	}
	if priority != "" {
		task.Priority = model.ProjectPriority(priority)
	}
	if dueDate != nil {
		task.DueDate = dueDate
	}
	if assignedTo != nil {
		isMember, err := s.memberRepo.IsMember(task.ProjectID, *assignedTo)
		if err != nil || !isMember {
			return nil, errors.New("assigned user is not a project member")
		}
		task.AssignedTo = assignedTo
	}
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	// Record status change
	if status != "" && oldStatus != model.TaskStatus(status) {
		history := &model.TaskHistory{
			TaskID:    taskID,
			OldStatus: string(oldStatus),
			NewStatus: status,
			ChangedBy: userID,
			ChangedAt: time.Now(),
		}
		s.taskRepo.AddTaskHistory(history)
	}

	return task, nil
}

// UpdateTaskStatus updates only the task status
func (s *TaskService) UpdateTaskStatus(taskID uuid.UUID, status string, userID uuid.UUID) (*model.Task, error) {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	isMember, err := s.memberRepo.IsMember(task.ProjectID, userID)
	if err != nil || !isMember {
		return nil, errors.New("you are not a member of this project")
	}

	oldStatus := task.Status
	task.Status = model.TaskStatus(status)
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	history := &model.TaskHistory{
		TaskID:    taskID,
		OldStatus: string(oldStatus),
		NewStatus: status,
		ChangedBy: userID,
		ChangedAt: time.Now(),
	}
	s.taskRepo.AddTaskHistory(history)

	return task, nil
}

// DeleteTask performs a soft delete on a task
func (s *TaskService) DeleteTask(taskID, userID uuid.UUID) error {
	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	role, err := s.memberRepo.GetRole(task.ProjectID, userID)
	if err != nil || (role != string(model.MemberRoleOwner) && role != string(model.MemberRoleManager)) {
		return errors.New("insufficient permissions")
	}

	return s.taskRepo.Delete(taskID)
}
