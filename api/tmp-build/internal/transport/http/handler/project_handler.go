package handler

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/transport/http/dto"
	"dac/project-tracker/pkg/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ProjectHandler handles project HTTP requests
type ProjectHandler struct {
	projectService *service.ProjectService
	memberRepo     repository.MemberRepository
	validate       *validator.Validate
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService *service.ProjectService, memberRepo repository.MemberRepository) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		memberRepo:     memberRepo,
		validate:       validator.New(),
	}
}

func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}
	t, _ := time.Parse("2006-01-02", dateStr)
	return &t
}

func toUserDTO(user *model.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
	}
}

// List returns projects with pagination and filters
func (h *ProjectHandler) List(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	filter := repository.ProjectFilter{
		Page:      page,
		PageSize:  pageSize,
		Status:    c.Query("status"),
		Priority:  c.Query("priority"),
		Search:    c.Query("search"),
		CreatedBy: c.Query("created_by"),
		Sort:      c.Query("sort"),
		Order:     c.Query("order"),
	}

	projects, total, err := h.projectService.GetProjects(userID, filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var projectResponses []dto.ProjectResponse
	for _, p := range projects {
		projectResponses = append(projectResponses, toProjectResponse(&p))
	}

	resp := dto.ProjectListResponse{
		Projects: projectResponses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}

	response.Success(c, resp)
}

// Get returns a project by ID
func (h *ProjectHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(id)
	if err != nil {
		response.NotFound(c, "Project not found")
		return
	}

	response.Success(c, toProjectResponse(project))
}

// Create creates a new project
func (h *ProjectHandler) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
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

	project, err := h.projectService.CreateProject(
		req.Name,
		req.Description,
		req.Priority,
		parseDate(req.StartDate),
		parseDate(req.EndDate),
		userID,
	)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, toProjectResponse(project))
}

// Update updates a project
func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req dto.UpdateProjectRequest
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

	project, err := h.projectService.UpdateProject(
		id,
		req.Name,
		req.Description,
		req.Status,
		req.Priority,
		parseDate(req.StartDate),
		parseDate(req.EndDate),
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

	response.Success(c, toProjectResponse(project))
}

// UpdateStatus updates only the project status
func (h *ProjectHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req dto.UpdateProjectStatusRequest
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

	project, err := h.projectService.UpdateProjectStatus(id, req.Status, userID)
	if err != nil {
		if err.Error() == "insufficient permissions" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":      project.ID.String(),
		"status":  project.Status,
		"message": "Status updated successfully",
	})
}

// Delete removes a project (soft delete)
func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	userID := uuid.MustParse(c.GetString("user_id"))

	if err := h.projectService.DeleteProject(id, userID); err != nil {
		if err.Error() == "only owner can delete project" {
			response.Forbidden(c, err.Error())
			return
		}
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Project deleted successfully"})
}

func toProjectResponse(project *model.Project) dto.ProjectResponse {
	var members []dto.MemberResponse
	for _, m := range project.Members {
		members = append(members, dto.MemberResponse{
			ID: m.ID.String(),
			User: dto.UserDTO{
				ID:        m.User.ID.String(),
				Name:      m.User.Name,
				Email:     m.User.Email,
				Role:      m.User.Role,
				AvatarURL: m.User.AvatarURL,
			},
			Role:     string(m.Role),
			JoinedAt: m.JoinedAt.Format(time.RFC3339),
		})
	}

	var startDate, endDate string
	if project.StartDate != nil {
		startDate = project.StartDate.Format("2006-01-02")
	}
	if project.EndDate != nil {
		endDate = project.EndDate.Format("2006-01-02")
	}

	return dto.ProjectResponse{
		ID:          project.ID.String(),
		Name:        project.Name,
		Description: project.Description,
		Status:      string(project.Status),
		Priority:    string(project.Priority),
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedBy: dto.UserDTO{
			ID:        project.Creator.ID.String(),
			Name:      project.Creator.Name,
			Email:     project.Creator.Email,
			Role:      project.Creator.Role,
			AvatarURL: project.Creator.AvatarURL,
		},
		Members:    members,
		TasksCount: len(project.Tasks),
		CreatedAt:  project.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  project.UpdatedAt.Format(time.RFC3339),
	}
}
