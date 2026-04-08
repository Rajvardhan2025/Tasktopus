package models

type WSEventType string

const (
	WSEventIssueCreated  WSEventType = "issue_created"
	WSEventIssueUpdated  WSEventType = "issue_updated"
	WSEventIssueMoved    WSEventType = "issue_moved"
	WSEventCommentAdded  WSEventType = "comment_added"
	WSEventSprintUpdated WSEventType = "sprint_updated"
	WSEventPresence      WSEventType = "presence"
)

type WSEvent struct {
	Type      WSEventType `json:"type"`
	ProjectID string      `json:"project_id"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

type PresenceData struct {
	UserID  string `json:"user_id"`
	IssueID string `json:"issue_id,omitempty"`
	BoardID string `json:"board_id,omitempty"`
	Action  string `json:"action"` // joined, left
}
