package models

import "time"

type ActivityAction string

const (
	ActivityIssueCreated    ActivityAction = "issue_created"
	ActivityIssueUpdated    ActivityAction = "issue_updated"
	ActivityIssueDeleted    ActivityAction = "issue_deleted"
	ActivityStatusChanged   ActivityAction = "status_changed"
	ActivityCommentAdded    ActivityAction = "comment_added"
	ActivityCommentUpdated  ActivityAction = "comment_updated"
	ActivityCommentDeleted  ActivityAction = "comment_deleted"
	ActivitySprintStarted   ActivityAction = "sprint_started"
	ActivitySprintCompleted ActivityAction = "sprint_completed"
	ActivityAssigneeChanged ActivityAction = "assignee_changed"
	ActivityWatcherAdded    ActivityAction = "watcher_added"
	ActivityWatcherRemoved  ActivityAction = "watcher_removed"
)

type Activity struct {
	ID        string                 `json:"id" bson:"_id"`
	ProjectID string                 `json:"project_id" bson:"project_id"`
	IssueID   string                 `json:"issue_id,omitempty" bson:"issue_id,omitempty"`
	UserID    string                 `json:"user_id" bson:"user_id"`
	Action    ActivityAction         `json:"action" bson:"action"`
	Changes   map[string]interface{} `json:"changes,omitempty" bson:"changes,omitempty"`
	Timestamp time.Time              `json:"timestamp" bson:"timestamp"`
}

type ActivityWithUser struct {
	Activity `bson:",inline"`
	User     *UserRef `json:"user" bson:"-"`
}
