package store

import (
	"context"
	"time"

	"github.com/yourusername/project-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationStore struct {
	collection *mongo.Collection
}

func NewNotificationStore(db *mongo.Database) *NotificationStore {
	return &NotificationStore{
		collection: db.Collection("notifications"),
	}
}

func (s *NotificationStore) Create(ctx context.Context, notification *models.Notification) error {
	notification.Timestamp = time.Now()
	notification.Read = false
	_, err := s.collection.InsertOne(ctx, notification)
	return err
}

func (s *NotificationStore) FindByUser(ctx context.Context, userID string, limit int) ([]*models.Notification, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := s.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, err
	}
	return notifications, nil
}

func (s *NotificationStore) MarkAsRead(ctx context.Context, id string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"read": true}},
	)
	return err
}

func (s *NotificationStore) MarkAllAsRead(ctx context.Context, userID string) error {
	_, err := s.collection.UpdateMany(
		ctx,
		bson.M{"user_id": userID, "read": false},
		bson.M{"$set": bson.M{"read": true}},
	)
	return err
}
