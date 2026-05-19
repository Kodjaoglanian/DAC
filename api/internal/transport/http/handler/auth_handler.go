package handler

import (
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/transport/http/dto"
	"dac/project-tracker/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
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

	token, user, err := h.authService.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already registered" {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	resp := dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
		},
	}

	response.Created(c, resp)
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
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

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	resp := dto.AuthResponse{
		Token: token,
		User: dto.UserDTO{
			ID:        user.ID.String(),
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
		},
	}

	response.Success(c, resp)
}

// Me returns the authenticated user's data
func (h *AuthHandler) Me(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	if userIDStr == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.authService.GetUserByID(uuid.MustParse(userIDStr))
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
