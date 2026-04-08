package store

import (
	"context"
	"time"

	"github.com/yourusername/project-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ActivityStore struct {
	collection *mongo.Collection
}

func NewActivityStore(db *mongo.Database) *ActivityStore {
	return &ActivityStore{
		collection: db.Collection("activities"),
	}
}

func (s *ActivityStore) Create(ctx context.Context, activity *models.Activity) error {
	activity.Timestamp = time.Now()
	_, err := s.collection.InsertOne(ctx, activity)
	return err
}

func (s *ActivityStore) FindByProject(ctx context.Context, projectID string, limit, skip int) ([]*models.Activity, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(skip))

	cursor, err := s.collection.Find(ctx, bson.M{"project_id": projectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []*models.Activity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityStore) FindByIssue(ctx context.Context, issueID string) ([]*models.Activity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := s.collection.Find(ctx, bson.M{"issue_id": issueID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []*models.Activity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityStore) FindAfterTimestamp(ctx context.Context, projectID string, after time.Time) ([]*models.Activity, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})
	cursor, err := s.collection.Find(
		ctx,
		bson.M{"project_id": projectID, "timestamp": bson.M{"$gt": after}},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []*models.Activity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}
