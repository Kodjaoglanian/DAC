package dto

type AddMemberRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
	Role   string `json:"role" validate:"required,oneof=owner manager member"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=owner manager member"`
}

type MemberResponse struct {
	ID       string  `json:"id"`
	User     UserDTO `json:"user"`
	Role     string  `json:"role"`
	JoinedAt string  `json:"joined_at"`
}
