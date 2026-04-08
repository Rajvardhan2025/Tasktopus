package models

type User struct {
	BaseModel   `bson:",inline"`
	Email       string `json:"email" bson:"email" validate:"required,email"`
	DisplayName string `json:"display_name" bson:"display_name" validate:"required"`
	AvatarURL   string `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
}
