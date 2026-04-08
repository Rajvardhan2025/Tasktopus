package models

import "time"

type SprintStatus string

const (
	SprintStatusFuture SprintStatus = "future"
	SprintStatusActive SprintStatus = "active"
	SprintStatusClosed SprintStatus = "closed"
)

type Sprint struct {
	BaseModel          `bson:",inline"`
	ProjectID          string       `json:"project_id" bson:"project_id" validate:"required"`
	Name               string       `json:"name" bson:"name" validate:"required,min=1,max=100"`
	Goal               string       `json:"goal" bson:"goal"`
	StartDate          *time.Time   `json:"start_date" bson:"start_date"`
	EndDate            *time.Time   `json:"end_date" bson:"end_date"`
	Status             SprintStatus `json:"status" bson:"status" validate:"required"`
	CompletedPoints    int          `json:"completed_points" bson:"completed_points"`
	TotalPoints        int          `json:"total_points" bson:"total_points"`
	Velocity           int          `json:"velocity" bson:"velocity"`
	IsDefault          bool         `json:"is_default" bson:"is_default"`
	IssueCount         int          `json:"issue_count" bson:"issue_count"`
	CompleteIssueCount int          `json:"complete_issue_count" bson:"complete_issue_count"`
}

type CreateSprintRequest struct {
	ProjectID string     `json:"project_id" validate:"required"`
	Name      string     `json:"name" validate:"required,min=1,max=100"`
	Goal      string     `json:"goal"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type UpdateSprintRequest struct {
	Name      *string    `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Goal      *string    `json:"goal,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
}

type CloseSprintRequest struct {
	CarryOverIssues []string `json:"carry_over_issues"`
}
