package middleware

import (
	"dac/project-tracker/pkg/response"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware checks if the authenticated user has the required role
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("user_role")
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "Acesso negado: permissão insuficiente")
		c.Abort()
	}
}
