package dto

// LoginRequest — corpo da requisição de login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest — corpo da requisição de registro
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// AuthResponse — resposta com token JWT
type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
