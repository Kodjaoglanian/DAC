package middleware

import (
	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProjectMemberMiddleware checks if user is a member of the project
func ProjectMemberMiddleware(memberRepo repository.MemberRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectIDStr := c.Param("projectId")
		if projectIDStr == "" {
			projectIDStr = c.Param("id")
		}

		userIDStr := c.GetString("user_id")
		if userIDStr == "" {
			response.Forbidden(c, "Usuário não autenticado")
			c.Abort()
			return
		}

		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			response.BadRequest(c, "Invalid project ID")
			c.Abort()
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			response.Unauthorized(c, "Invalid user ID")
			c.Abort()
			return
		}

		isMember, err := memberRepo.IsMember(projectID, userID)
		if err != nil || !isMember {
			response.Forbidden(c, "Você não é membro deste projeto")
			c.Abort()
			return
		}

		c.Next()
	}
}
