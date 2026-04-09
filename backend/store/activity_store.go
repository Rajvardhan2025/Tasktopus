package store

import (
	"context"
	"fmt"
	"strings"
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

func (s *ActivityStore) FindByProjectCursor(ctx context.Context, projectID string, action, userID string, limit int, cursor string) ([]*models.Activity, string, error) {
	if limit <= 0 {
		limit = 50
	}

	filters := bson.M{"project_id": projectID}
	if action != "" {
		filters["action"] = action
	}
	if userID != "" {
		filters["user_id"] = userID
	}

	query := filters
	if cursor != "" {
		parts := strings.SplitN(cursor, "|", 2)
		if len(parts) == 2 {
			ts, err := time.Parse(time.RFC3339Nano, parts[0])
			if err == nil {
				cursorFilter := bson.M{
					"$or": []bson.M{
						{"timestamp": bson.M{"$lt": ts}},
						{
							"$and": []bson.M{
								{"timestamp": ts},
								{"_id": bson.M{"$lt": parts[1]}},
							},
						},
					},
				}
				query = bson.M{"$and": []bson.M{filters, cursorFilter}}
			}
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}, {Key: "_id", Value: -1}}).
		SetLimit(int64(limit + 1))

	cursorResult, err := s.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursorResult.Close(ctx)

	var activities []*models.Activity
	if err := cursorResult.All(ctx, &activities); err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(activities) > limit {
		last := activities[limit-1]
		nextCursor = fmt.Sprintf("%s|%s", last.Timestamp.UTC().Format(time.RFC3339Nano), last.ID)
		activities = activities[:limit]
	}

	return activities, nextCursor, nil
}
