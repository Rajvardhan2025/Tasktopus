package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
)

type NotificationService struct {
	notificationStore *store.NotificationStore
}

func NewNotificationService(notificationStore *store.NotificationStore) *NotificationService {
	return &NotificationService{
		notificationStore: notificationStore,
	}
}

func (s *NotificationService) NotifyAssignment(ctx context.Context, userID, issueID, actorID string) error {
	notification := &models.Notification{
		ID:      uuid.New().String(),
		UserID:  userID,
		Type:    models.NotificationAssigned,
		IssueID: issueID,
		ActorID: actorID,
		Message: "You have been assigned to an issue",
	}
	return s.notificationStore.Create(ctx, notification)
}

func (s *NotificationService) NotifyMention(ctx context.Context, userID, issueID, actorID string) error {
	notification := &models.Notification{
		ID:      uuid.New().String(),
		UserID:  userID,
		Type:    models.NotificationMentioned,
		IssueID: issueID,
		ActorID: actorID,
		Message: "You were mentioned in a comment",
	}
	return s.notificationStore.Create(ctx, notification)
}

func (s *NotificationService) NotifyWatcher(ctx context.Context, userID, issueID, actorID, action string) error {
	notification := &models.Notification{
		ID:      uuid.New().String(),
		UserID:  userID,
		Type:    models.NotificationWatched,
		IssueID: issueID,
		ActorID: actorID,
		Message: fmt.Sprintf("Issue you're watching: %s", action),
	}
	return s.notificationStore.Create(ctx, notification)
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string, limit int) ([]*models.Notification, error) {
	return s.notificationStore.FindByUser(ctx, userID, limit)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	return s.notificationStore.MarkAsRead(ctx, notificationID)
}
