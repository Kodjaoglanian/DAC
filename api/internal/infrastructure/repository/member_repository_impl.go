package repository

import (
	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MemberRepositoryImpl implements repository.MemberRepository
type MemberRepositoryImpl struct {
	db *gorm.DB
}

// NewMemberRepository creates a new member repository
func NewMemberRepository(db *gorm.DB) repository.MemberRepository {
	return &MemberRepositoryImpl{db: db}
}

func (r *MemberRepositoryImpl) Create(member *model.ProjectMember) error {
	return r.db.Create(member).Error
}

func (r *MemberRepositoryImpl) FindByProject(projectID uuid.UUID) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	if err := r.db.Preload("User").Where("project_id = ?", projectID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *MemberRepositoryImpl) FindByProjectAndUser(projectID, userID uuid.UUID) (*model.ProjectMember, error) {
	var member model.ProjectMember
	if err := r.db.Preload("User").Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *MemberRepositoryImpl) Update(member *model.ProjectMember) error {
	return r.db.Save(member).Error
}

func (r *MemberRepositoryImpl) Delete(projectID, userID uuid.UUID) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&model.ProjectMember{}).Error
}

func (r *MemberRepositoryImpl) IsMember(projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.ProjectMember{}).Where("project_id = ? AND user_id = ?", projectID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	// Check if user is the project creator
	err = r.db.Model(&model.Project{}).Where("id = ? AND created_by = ?", projectID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MemberRepositoryImpl) GetRole(projectID, userID uuid.UUID) (string, error) {
	var member model.ProjectMember
	if err := r.db.Select("role").Where("project_id = ? AND user_id = ?", projectID, userID).First(&member).Error; err != nil {
		// Check if user is the project creator
		var count int64
		if err := r.db.Model(&model.Project{}).Where("id = ? AND created_by = ?", projectID, userID).Count(&count).Error; err != nil {
			return "", err
		}
		if count > 0 {
			return string(model.MemberRoleOwner), nil
		}
		return "", err
	}
	return string(member.Role), nil
}

func (r *MemberRepositoryImpl) CountOwners(projectID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.ProjectMember{}).Where("project_id = ? AND role = ?", projectID, model.MemberRoleOwner).Count(&count).Error
	return count, err
}
