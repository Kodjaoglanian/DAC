package dto

type CreateProjectRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,min=10"`
	Priority    string `json:"priority" validate:"required,oneof=low medium high critical"`
	StartDate   string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate     string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name" validate:"omitempty,min=3,max=255"`
	Description string `json:"description" validate:"omitempty,min=10"`
	Status      string `json:"status" validate:"omitempty,oneof=planning in_progress completed cancelled"`
	Priority    string `json:"priority" validate:"omitempty,oneof=low medium high critical"`
	StartDate   string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate     string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

type UpdateProjectStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=planning in_progress completed cancelled"`
}

type ProjectResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	Priority    string           `json:"priority"`
	StartDate   string           `json:"start_date"`
	EndDate     string           `json:"end_date"`
	CreatedBy   UserDTO          `json:"created_by"`
	Members     []MemberResponse `json:"members"`
	TasksCount  int              `json:"tasks_count"`
	CreatedAt   string           `json:"created_at"`
	UpdatedAt   string           `json:"updated_at"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}
