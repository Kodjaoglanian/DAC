package handler

import (
	"dac/project-tracker/internal/domain/service"
	"dac/project-tracker/internal/transport/http/dto"
	"dac/project-tracker/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReportHandler handles report HTTP requests
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// Dashboard returns consolidated dashboard data
func (h *ReportHandler) Dashboard(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))

	data, err := h.reportService.GetDashboard(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	var recentProjects []dto.ProjectResponse
	for _, p := range data.RecentProjects {
		recentProjects = append(recentProjects, toProjectResponse(&p))
	}

	resp := dto.DashboardResponse{
		TotalProjects:      data.TotalProjects,
		ActiveProjects:     data.ActiveProjects,
		CompletedProjects:  data.CompletedProjects,
		TotalTasks:         data.TotalTasks,
		TasksByStatus:      data.TasksByStatus,
		ProjectsByStatus:   data.ProjectsByStatus,
		ProjectsByPriority: data.ProjectsByPriority,
		RecentProjects:     recentProjects,
	}

	response.Success(c, resp)
}

// ProjectReport returns detailed project report
func (h *ReportHandler) ProjectReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid project ID")
		return
	}

	data, err := h.reportService.GetProjectReport(id)
	if err != nil {
		response.NotFound(c, "Project not found")
		return
	}

	var progress float64
	if data.TotalTasks > 0 {
		progress = float64(data.CompletedTasks) / float64(data.TotalTasks) * 100
	}

	var statusHistory []dto.StatusHistoryDTO
	for _, h := range data.StatusHistory {
		statusHistory = append(statusHistory, dto.StatusHistoryDTO{
			OldStatus: h.OldStatus,
			NewStatus: h.NewStatus,
			ChangedBy: h.ChangedBy.String(),
			ChangedAt: h.ChangedAt.Format(time.RFC3339),
		})
	}

	resp := dto.ProjectReportResponse{
		Project:        toProjectResponse(data.Project),
		TotalTasks:     data.TotalTasks,
		CompletedTasks: data.CompletedTasks,
		Progress:       progress,
		StatusHistory:  statusHistory,
	}

	response.Success(c, resp)
}

// ProjectsByStatus returns projects grouped by status
func (h *ReportHandler) ProjectsByStatus(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))

	data, err := h.reportService.GetProjectsByStatus(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, data)
}

// TasksByStatus returns tasks grouped by status
func (h *ReportHandler) TasksByStatus(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))

	data, err := h.reportService.GetTasksByStatus(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, data)
}
