package handler

import (
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/transport/http/dto"
	"dac/project-tracker/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userService *service.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// List returns all users with pagination
func (h *UserHandler) List(c *gin.Context) {
	page := 1
	pageSize := 10
	// TODO: parse query params

	users, total, err := h.userService.GetAllUsers(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var userDTOs []dto.UserDTO
	for _, u := range users {
		userDTOs = append(userDTOs, dto.UserDTO{
			ID:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			AvatarURL: u.AvatarURL,
		})
	}

	response.Success(c, gin.H{
		"users": userDTOs,
		"total": total,
		"page":  page,
	})
}

// Get returns a user by ID
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	resp := dto.UserDTO{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
	}

	response.Success(c, resp)
}

// Update updates a user's data
func (h *UserHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest
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

	user, err := h.userService.UpdateUser(id, req.Name, req.AvatarURL)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	resp := dto.UserDTO{
		ID:        user.ID.String(),
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		AvatarURL: user.AvatarURL,
	}

	response.Success(c, resp)
}

// Delete removes a user (soft delete)
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "User deleted successfully"})
}
