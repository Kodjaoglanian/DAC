package service

import (
	"errors"
	"time"

	"dac/project-tracker/internal/domain/model"
	"dac/project-tracker/internal/domain/repository"

	"github.com/google/uuid"
)

// MemberService handles project member business logic
type MemberService struct {
	memberRepo  repository.MemberRepository
	projectRepo repository.ProjectRepository
}

// NewMemberService creates a new member service
func NewMemberService(memberRepo repository.MemberRepository, projectRepo repository.ProjectRepository) *MemberService {
	return &MemberService{
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
	}
}

// AddMember adds a user to a project
func (s *MemberService) AddMember(projectID, userID uuid.UUID, role string, addedBy uuid.UUID) (*model.ProjectMember, error) {
	requesterRole, err := s.memberRepo.GetRole(projectID, addedBy)
	if err != nil || (requesterRole != string(model.MemberRoleOwner) && requesterRole != string(model.MemberRoleManager)) {
		return nil, errors.New("insufficient permissions")
	}

	isMember, err := s.memberRepo.IsMember(projectID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, errors.New("user is already a member of this project")
	}

	member := &model.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      model.ProjectMemberRole(role),
		JoinedAt:  time.Now(),
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, err
	}

	return member, nil
}

// GetMembers returns all members of a project
func (s *MemberService) GetMembers(projectID uuid.UUID) ([]model.ProjectMember, error) {
	return s.memberRepo.FindByProject(projectID)
}

// UpdateMemberRole updates a member's role
func (s *MemberService) UpdateMemberRole(projectID, userID uuid.UUID, role string, updatedBy uuid.UUID) (*model.ProjectMember, error) {
	requesterRole, err := s.memberRepo.GetRole(projectID, updatedBy)
	if err != nil || requesterRole != string(model.MemberRoleOwner) {
		return nil, errors.New("only owner can change roles")
	}

	if userID == updatedBy {
		return nil, errors.New("cannot change your own role")
	}

	member, err := s.memberRepo.FindByProjectAndUser(projectID, userID)
	if err != nil {
		return nil, errors.New("member not found")
	}

	// Ensure at least one owner remains
	if member.Role == model.MemberRoleOwner && model.ProjectMemberRole(role) != model.MemberRoleOwner {
		ownerCount, err := s.memberRepo.CountOwners(projectID)
		if err != nil {
			return nil, err
		}
		if ownerCount <= 1 {
			return nil, errors.New("project must have at least one owner")
		}
	}

	member.Role = model.ProjectMemberRole(role)
	if err := s.memberRepo.Update(member); err != nil {
		return nil, err
	}

	return member, nil
}

// RemoveMember removes a member from a project
func (s *MemberService) RemoveMember(projectID, userID uuid.UUID, removedBy uuid.UUID) error {
	requesterRole, err := s.memberRepo.GetRole(projectID, removedBy)
	if err != nil || (requesterRole != string(model.MemberRoleOwner) && requesterRole != string(model.MemberRoleManager)) {
		return errors.New("insufficient permissions")
	}

	member, err := s.memberRepo.FindByProjectAndUser(projectID, userID)
	if err != nil {
		return errors.New("member not found")
	}

	if member.Role == model.MemberRoleOwner {
		return errors.New("cannot remove the project owner")
	}

	return s.memberRepo.Delete(projectID, userID)
}
