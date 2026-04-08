package models

import "time"

type SprintStatus string

const (
	SprintStatusPlanned   SprintStatus = "planned"
	SprintStatusActive    SprintStatus = "active"
	SprintStatusCompleted SprintStatus = "completed"
)

type Sprint struct {
	BaseModel       `bson:",inline"`
	ProjectID       string       `json:"project_id" bson:"project_id" validate:"required"`
	Name            string       `json:"name" bson:"name" validate:"required,min=1,max=100"`
	Goal            string       `json:"goal" bson:"goal"`
	StartDate       time.Time    `json:"start_date" bson:"start_date"`
	EndDate         time.Time    `json:"end_date" bson:"end_date"`
	Status          SprintStatus `json:"status" bson:"status"`
	CompletedPoints int          `json:"completed_points" bson:"completed_points"`
	TotalPoints     int          `json:"total_points" bson:"total_points"`
}

type CreateSprintRequest struct {
	ProjectID string    `json:"project_id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
	Goal      string    `json:"goal"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateSprintRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Goal      *string    `json:"goal,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type CompleteSprintRequest struct {
	CarryOverIssues []string `json:"carry_over_issues"`
}
