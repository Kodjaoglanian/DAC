package repository

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskRepositoryImpl implements repository.TaskRepository
type TaskRepositoryImpl struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository(db *gorm.DB) repository.TaskRepository {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepositoryImpl) FindByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task
	if err := r.db.Preload("Assignee").First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}
	// GORM foreignKey bug: force reload project_id from raw SQL
	var result struct{ ProjectID string }
	r.db.Raw("SELECT project_id::text as project_id FROM tasks WHERE id = ?", id).Scan(&result)
	if result.ProjectID != "" {
		pid, _ := uuid.Parse(result.ProjectID)
		task.ProjectID = pid
	}
	return &task, nil
}

func (r *TaskRepositoryImpl) FindByProject(projectID uuid.UUID, filter repository.TaskFilter) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	query := r.db.Model(&model.Task{}).Preload("Assignee").Where("project_id = ?", projectID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.AssignedTo != "" {
		query = query.Where("assigned_to = ?", filter.AssignedTo)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ?", search)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortField := filter.Sort
	if sortField == "" {
		sortField = "created_at"
	}
	order := filter.Order
	if order == "" {
		order = "desc"
	}

	if err := query.Order(sortField + " " + order).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepositoryImpl) Update(task *model.Task) error {
	return r.db.Model(task).Omit("Project", "Assignee").Updates(task).Error
}

func (r *TaskRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}

func (r *TaskRepositoryImpl) AddTaskHistory(history *model.TaskHistory) error {
	return r.db.Create(history).Error
}

func (r *TaskRepositoryImpl) CountByProject(projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Task{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

func (r *TaskRepositoryImpl) CountCompletedByProject(projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Task{}).Where("project_id = ? AND status = ?", projectID, model.TaskStatusDone).Count(&count).Error
	return count, err
}

func (r *TaskRepositoryImpl) CountByStatus(userID uuid.UUID) (map[string]int64, error) {
	var results []struct {
		Status string
		Count  int64
	}

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)
	createdQuery := r.db.Table("projects").Select("id").Where("created_by = ?", userID)
	query := r.db.Model(&model.Task{}).Select("status, COUNT(*) as count").
		Where("project_id IN (?) OR project_id IN (?)", subQuery, createdQuery).
		Group("status")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	statusMap := make(map[string]int64)
	for _, r := range results {
		statusMap[r.Status] = r.Count
	}

	return statusMap, nil
}
