package handler

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/transport/http/dto"
	"dac/project-tracker/pkg/response"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// TaskHandler handles task HTTP requests
type TaskHandler struct {
	taskService *service.TaskService
	validate    *validator.Validate
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		validate:    validator.New(),
	}
}

// List returns tasks for a project
func (h *TaskHandler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	filter := repository.TaskFilter{
		Status:     c.Query("status"),
		Priority:   c.Query("priority"),
		AssignedTo: c.Query("assigned_to"),
		Search:     c.Query("search"),
		Sort:       c.Query("sort"),
		Order:      c.Query("order"),
	}

	tasks, total, err := h.taskService.GetTasksByProject(projectID, filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var taskResponses []dto.TaskResponse
	for _, t := range tasks {
		taskResponses = append(taskResponses, toTaskResponse(&t))
	}

	response.Success(c, gin.H{
		"tasks": taskResponses,
		"total": total,
	})
}

// Get returns a task by ID
func (h *TaskHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid task ID")
		return
	}

	task, err := h.taskService.GetTaskByID(id)
	if err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	response.Success(c, toTaskResponse(task))
}

// Create creates a new task in a project
func (h *TaskHandler) Create(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var details []response.ErrorDetail
		for _, err := range err.(validator.ValidationErrors) {
			details = append(details, response.ErrorDetail{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		response.ValidationError(c, details...)
		return
	}

	userID := uuid.MustParse(c.GetString("user_id"))

	var assignedTo *uuid.UUID
	if req.AssignedTo != "" {
		id := uuid.MustParse(req.AssignedTo)
		assignedTo = &id
	}

	task, err := h.taskService.CreateTask(
		projectID,
		req.Title,
		req.Description,
		req.Priority,
		parseDate(req.DueDate),
		assignedTo,
		userID,
	)
	if err != nil {
		if err.Error() == "insufficient permissions" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, toTaskResponse(task))
}

// Update updates a task
func (h *TaskHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid task ID")
		return
	}

	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var details []response.ErrorDetail
		for _, err := range err.(validator.ValidationErrors) {
			details = append(details, response.ErrorDetail{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		response.ValidationError(c, details...)
		return
	}

	userID := uuid.MustParse(c.GetString("user_id"))

	var assignedTo *uuid.UUID
	if req.AssignedTo != "" {
		uid := uuid.MustParse(req.AssignedTo)
		assignedTo = &uid
	}

	task, err := h.taskService.UpdateTask(
		id,
		req.Title,
		req.Description,
		req.Status,
		req.Priority,
		parseDate(req.DueDate),
		assignedTo,
		userID,
	)
	if err != nil {
		if err.Error() == "insufficient permissions" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, toTaskResponse(task))
}

// UpdateStatus updates only the task status
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid task ID")
		return
	}

	var req dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		var details []response.ErrorDetail
		for _, err := range err.(validator.ValidationErrors) {
			details = append(details, response.ErrorDetail{
				Field:   err.Field(),
				Message: err.Tag(),
			})
		}
		response.ValidationError(c, details...)
		return
	}

	userID := uuid.MustParse(c.GetString("user_id"))

	fmt.Printf("HANDLER DEBUG: taskID=%s, status=%s, userID=%s\n", id, req.Status, userID)

	task, err := h.taskService.UpdateTaskStatus(id, req.Status, userID)
	if err != nil {
		fmt.Printf("HANDLER DEBUG ERROR: %v\n", err)
		if err.Error() == "you are not a member of this project" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":      task.ID.String(),
		"status":  task.Status,
		"message": "Status updated successfully",
	})
}

// Delete removes a task (soft delete)
func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid task ID")
		return
	}

	userID := uuid.MustParse(c.GetString("user_id"))

	if err := h.taskService.DeleteTask(id, userID); err != nil {
		if err.Error() == "insufficient permissions" {
			response.Forbidden(c, err.Error())
			return
		}
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Task deleted successfully"})
}

func toTaskResponse(task *model.Task) dto.TaskResponse {
	var dueDate string
	if task.DueDate != nil {
		dueDate = task.DueDate.Format("2006-01-02")
	}

	var assignedTo *dto.UserDTO
	if task.Assignee != nil {
		assignedTo = &dto.UserDTO{
			ID:        task.Assignee.ID.String(),
			Name:      task.Assignee.Name,
			Email:     task.Assignee.Email,
			Role:      task.Assignee.Role,
			AvatarURL: task.Assignee.AvatarURL,
		}
	}

	return dto.TaskResponse{
		ID:          task.ID.String(),
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		DueDate:     dueDate,
		ProjectID:   task.ProjectID.String(),
		AssignedTo:  assignedTo,
		CreatedAt:   task.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Format(time.RFC3339),
	}
}
