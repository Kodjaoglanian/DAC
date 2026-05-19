package dto

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description"`
	Priority    string `json:"priority" validate:"required,oneof=low medium high critical"`
	DueDate     string `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
	AssignedTo  string `json:"assigned_to" validate:"omitempty,uuid"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"omitempty,min=3,max=255"`
	Description string `json:"description"`
	Status      string `json:"status" validate:"omitempty,oneof=todo in_progress review done"`
	Priority    string `json:"priority" validate:"omitempty,oneof=low medium high critical"`
	DueDate     string `json:"due_date" validate:"omitempty,datetime=2006-01-02"`
	AssignedTo  string `json:"assigned_to" validate:"omitempty,uuid"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=todo in_progress review done"`
}

type TaskResponse struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Priority    string   `json:"priority"`
	DueDate     string   `json:"due_date"`
	ProjectID   string   `json:"project_id"`
	AssignedTo  *UserDTO `json:"assigned_to"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
