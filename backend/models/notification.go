package models

import "time"

type NotificationType string

const (
	NotificationAssigned  NotificationType = "assigned"
	NotificationMentioned NotificationType = "mentioned"
	NotificationWatched   NotificationType = "watched"
	NotificationCommented NotificationType = "commented"
)

type Notification struct {
	ID        string           `json:"id" bson:"_id"`
	UserID    string           `json:"user_id" bson:"user_id"`
	Type      NotificationType `json:"type" bson:"type"`
	IssueID   string           `json:"issue_id" bson:"issue_id"`
	ActorID   string           `json:"actor_id" bson:"actor_id"`
	Message   string           `json:"message" bson:"message"`
	Read      bool             `json:"read" bson:"read"`
	Timestamp time.Time        `json:"timestamp" bson:"timestamp"`
}
