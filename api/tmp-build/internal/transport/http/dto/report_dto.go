package dto

type DashboardResponse struct {
	TotalProjects      int64             `json:"total_projects"`
	ActiveProjects     int64             `json:"active_projects"`
	CompletedProjects  int64             `json:"completed_projects"`
	TotalTasks         int64             `json:"total_tasks"`
	TasksByStatus      map[string]int64  `json:"tasks_by_status"`
	ProjectsByStatus   map[string]int64  `json:"projects_by_status"`
	ProjectsByPriority map[string]int64  `json:"projects_by_priority"`
	RecentProjects     []ProjectResponse `json:"recent_projects"`
}

type ProjectReportResponse struct {
	Project        ProjectResponse    `json:"project"`
	TotalTasks     int64              `json:"total_tasks"`
	CompletedTasks int64              `json:"completed_tasks"`
	Progress       float64            `json:"progress"`
	StatusHistory  []StatusHistoryDTO `json:"status_history"`
}

type StatusHistoryDTO struct {
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	ChangedBy string `json:"changed_by"`
	ChangedAt string `json:"changed_at"`
}
