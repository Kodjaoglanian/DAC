package dto

type UserDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
}

type UpdateUserRequest struct {
	Name      string `json:"name" validate:"omitempty,min=3,max=255"`
	AvatarURL string `json:"avatar_url" validate:"omitempty,url"`
}
