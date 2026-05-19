package handler

import (
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/transport/http/dto"
	"dac/project-tracker/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// MemberHandler handles project member HTTP requests
type MemberHandler struct {
	memberService *service.MemberService
	validate      *validator.Validate
}

// NewMemberHandler creates a new member handler
func NewMemberHandler(memberService *service.MemberService) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
		validate:      validator.New(),
	}
}

// List returns all members of a project
func (h *MemberHandler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	members, err := h.memberService.GetMembers(projectID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var memberResponses []dto.MemberResponse
	for _, m := range members {
		memberResponses = append(memberResponses, dto.MemberResponse{
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

	response.Success(c, memberResponses)
}

// Add adds a member to a project
func (h *MemberHandler) Add(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	var req dto.AddMemberRequest
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
	memberID := uuid.MustParse(req.UserID)

	member, err := h.memberService.AddMember(projectID, memberID, req.Role, userID)
	if err != nil {
		if err.Error() == "insufficient permissions" {
			response.Forbidden(c, err.Error())
			return
		}
		if err.Error() == "user is already a member of this project" {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, dto.MemberResponse{
		ID: member.ID.String(),
		User: dto.UserDTO{
			ID:   member.UserID.String(),
			Name: "", // Will be populated by preload
		},
		Role:     string(member.Role),
		JoinedAt: member.JoinedAt.Format(time.RFC3339),
	})
}

// UpdateRole updates a member's role
func (h *MemberHandler) UpdateRole(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req dto.UpdateMemberRoleRequest
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

	updatedBy := uuid.MustParse(c.GetString("user_id"))

	member, err := h.memberService.UpdateMemberRole(projectID, userID, req.Role, updatedBy)
	if err != nil {
		if err.Error() == "only owner can change roles" || err.Error() == "cannot change your own role" || err.Error() == "project must have at least one owner" {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, dto.MemberResponse{
		ID:   member.ID.String(),
		Role: string(member.Role),
	})
}

// Remove removes a member from a project
func (h *MemberHandler) Remove(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	removedBy := uuid.MustParse(c.GetString("user_id"))

	if err := h.memberService.RemoveMember(projectID, userID, removedBy); err != nil {
		if err.Error() == "insufficient permissions" || err.Error() == "cannot remove the project owner" {
			response.Forbidden(c, err.Error())
			return
		}
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "Member removed successfully"})
}
