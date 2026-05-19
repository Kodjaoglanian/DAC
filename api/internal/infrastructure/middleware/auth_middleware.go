package middleware

import (
	"strings"

	"dac/project-tracker/internal/config"
	"dac/project-tracker/pkg/auth"
	"dac/project-tracker/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies JWT token and injects user data into context
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Token de autenticação não fornecido")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			response.Unauthorized(c, "Token inválido ou expirado")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}
