package models

import "time"

type BaseModel struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CustomField struct {
	Name    string   `json:"name" bson:"name"`
	Type    string   `json:"type" bson:"type"` // text, number, dropdown, date
	Options []string `json:"options,omitempty" bson:"options,omitempty"`
}

type UserRef struct {
	UserID      string `json:"user_id" bson:"user_id"`
	DisplayName string `json:"display_name" bson:"display_name"`
}
