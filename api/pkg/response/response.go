package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorDetail represents a single validation error
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIResponse is the standard response format
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error response
type APIError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, status int, code, message string, details ...ErrorDetail) {
	resp := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
	if len(details) > 0 {
		resp.Error.Details = details
	}
	c.JSON(status, resp)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, details ...ErrorDetail) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// ValidationError sends a 400 Validation Error response
func ValidationError(c *gin.Context, details ...ErrorDetail) {
	Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Dados inválidos", details...)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, "FORBIDDEN", message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict sends a 409 Conflict response
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, "CONFLICT", message)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
