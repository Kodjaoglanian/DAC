package middleware

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"
	"dac/project-tracker/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectMemberMiddleware checks if user is a member of the project
func ProjectMemberMiddleware(memberRepo repository.MemberRepository, db *gorm.DB) gin.HandlerFunc {
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

		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/tasks/") {
			var task model.Task
			if err := db.Select("project_id").First(&task, "id = ?", projectID).Error; err != nil {
				response.BadRequest(c, "Invalid task ID")
				c.Abort()
				return
			}
			projectID = task.ProjectID
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
