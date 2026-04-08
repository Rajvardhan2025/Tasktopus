package models

type IssueType string

const (
	IssueTypeEpic    IssueType = "epic"
	IssueTypeStory   IssueType = "story"
	IssueTypeTask    IssueType = "task"
	IssueTypeBug     IssueType = "bug"
	IssueTypeSubtask IssueType = "subtask"
)

type Priority string

const (
	PriorityLowest  Priority = "lowest"
	PriorityLow     Priority = "low"
	PriorityMedium  Priority = "medium"
	PriorityHigh    Priority = "high"
	PriorityHighest Priority = "highest"
)

type Issue struct {
	BaseModel    `bson:",inline"`
	IssueKey     string                 `json:"issue_key" bson:"issue_key"`
	ProjectID    string                 `json:"project_id" bson:"project_id" validate:"required"`
	Type         IssueType              `json:"type" bson:"type" validate:"required"`
	Title        string                 `json:"title" bson:"title" validate:"required,min=1,max=200"`
	Description  string                 `json:"description" bson:"description"`
	Status       string                 `json:"status" bson:"status" validate:"required"`
	Priority     Priority               `json:"priority" bson:"priority"`
	AssigneeID   string                 `json:"assignee_id,omitempty" bson:"assignee_id,omitempty"`
	ReporterID   string                 `json:"reporter_id" bson:"reporter_id" validate:"required"`
	SprintID     string                 `json:"sprint_id,omitempty" bson:"sprint_id,omitempty"`
	ParentID     string                 `json:"parent_id,omitempty" bson:"parent_id,omitempty"`
	Labels       []string               `json:"labels" bson:"labels"`
	StoryPoints  int                    `json:"story_points,omitempty" bson:"story_points,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty" bson:"custom_fields,omitempty"`
	Watchers     []string               `json:"watchers" bson:"watchers"`
	Version      int                    `json:"version" bson:"version"`
}

type IssueWithRefs struct {
	Issue    `bson:",inline"`
	Assignee *UserRef `json:"assignee,omitempty" bson:"-"`
	Reporter *UserRef `json:"reporter,omitempty" bson:"-"`
	Sprint   *Sprint  `json:"sprint,omitempty" bson:"-"`
}

type CreateIssueRequest struct {
	ProjectID    string                 `json:"project_id" validate:"required"`
	Type         IssueType              `json:"type" validate:"required"`
	Title        string                 `json:"title" validate:"required,min=1,max=200"`
	Description  string                 `json:"description"`
	Priority     Priority               `json:"priority"`
	AssigneeID   string                 `json:"assignee_id,omitempty"`
	SprintID     string                 `json:"sprint_id,omitempty"`
	ParentID     string                 `json:"parent_id,omitempty"`
	Labels       []string               `json:"labels"`
	StoryPoints  int                    `json:"story_points,omitempty"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}

type UpdateIssueRequest struct {
	Title        *string                 `json:"title,omitempty" validate:"omitempty,min=1,max=200"`
	Description  *string                 `json:"description,omitempty"`
	Priority     *Priority               `json:"priority,omitempty"`
	AssigneeID   *string                 `json:"assignee_id,omitempty"`
	SprintID     *string                 `json:"sprint_id,omitempty"`
	Labels       *[]string               `json:"labels,omitempty"`
	StoryPoints  *int                    `json:"story_points,omitempty"`
	CustomFields *map[string]interface{} `json:"custom_fields,omitempty"`
	Version      int                     `json:"version" validate:"required"`
}

type TransitionRequest struct {
	ToStatus string `json:"to_status" validate:"required"`
	Version  int    `json:"version" validate:"required"`
}
