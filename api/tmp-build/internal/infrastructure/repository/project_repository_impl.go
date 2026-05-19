package repository

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProjectRepositoryImpl implements repository.ProjectRepository
type ProjectRepositoryImpl struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &ProjectRepositoryImpl{db: db}
}

func (r *ProjectRepositoryImpl) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepositoryImpl) FindByID(id uuid.UUID) (*model.Project, error) {
	var project model.Project
	if err := r.db.Preload("Creator").Preload("Members.User").First(&project, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepositoryImpl) FindAll(filter repository.ProjectFilter) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	query := r.db.Model(&model.Project{}).Preload("Creator")

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}
	if filter.CreatedBy != "" {
		query = query.Where("created_by = ?", filter.CreatedBy)
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

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Order(sortField + " " + order).Limit(filter.PageSize).Offset(offset).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepositoryImpl) FindByUser(userID uuid.UUID, page, pageSize int) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64

	subQuery := r.db.Table("project_members").Select("project_id").Where("user_id = ?", userID)

	query := r.db.Model(&model.Project{}).Preload("Creator").
		Where("id IN (?)", subQuery).Or("created_by = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

func (r *ProjectRepositoryImpl) Update(project *model.Project) error {
	return r.db.Model(project).Omit("Creator", "Members", "Tasks").Updates(project).Error
}

func (r *ProjectRepositoryImpl) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Project{}, "id = ?", id).Error
}

func (r *ProjectRepositoryImpl) AddStatusHistory(history *model.StatusHistory) error {
	return r.db.Create(history).Error
}
