package models

type User struct {
	BaseModel    `bson:",inline"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
	PasswordHash string `json:"-" bson:"password_hash"` // Never expose in JSON
	DisplayName  string `json:"display_name" bson:"display_name" validate:"required"`
	IsActive     bool   `json:"is_active" bson:"is_active"`
}

type CreateUserRequest struct {
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

// Auth models
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
}
