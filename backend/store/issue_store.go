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

type IssueStore struct {
	collection *mongo.Collection
}

func NewIssueStore(db *mongo.Database) *IssueStore {
	return &IssueStore{
		collection: db.Collection("issues"),
	}
}

func (s *IssueStore) Create(ctx context.Context, issue *models.Issue) error {
	issue.CreatedAt = time.Now()
	issue.UpdatedAt = time.Now()
	issue.Version = 1
	_, err := s.collection.InsertOne(ctx, issue)
	return err
}

func (s *IssueStore) FindByID(ctx context.Context, id string) (*models.Issue, error) {
	var issue models.Issue
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (s *IssueStore) FindByKey(ctx context.Context, key string) (*models.Issue, error) {
	var issue models.Issue
	err := s.collection.FindOne(ctx, bson.M{"issue_key": key}).Decode(&issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

func (s *IssueStore) FindByProject(ctx context.Context, projectID string) ([]*models.Issue, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var issues []*models.Issue
	if err := cursor.All(ctx, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (s *IssueStore) FindBySprint(ctx context.Context, sprintID string) ([]*models.Issue, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"sprint_id": sprintID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var issues []*models.Issue
	if err := cursor.All(ctx, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

// UpdateWithVersion implements optimistic locking
func (s *IssueStore) UpdateWithVersion(ctx context.Context, id string, version int, update bson.M) error {
	update["updated_at"] = time.Now()
	update["version"] = version + 1

	result, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "version": version},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("version conflict: issue has been modified")
	}

	return nil
}

func (s *IssueStore) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *IssueStore) AddWatcher(ctx context.Context, issueID, userID string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": issueID},
		bson.M{"$addToSet": bson.M{"watchers": userID}},
	)
	return err
}

func (s *IssueStore) RemoveWatcher(ctx context.Context, issueID, userID string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": issueID},
		bson.M{"$pull": bson.M{"watchers": userID}},
	)
	return err
}

func (s *IssueStore) UpdateFields(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": update,
			"$inc": bson.M{"version": 1},
		},
	)
	return err
}

func (s *IssueStore) MoveToBacklog(ctx context.Context, issueID, sprintID string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": issueID, "sprint_id": sprintID},
		bson.M{
			"$set": bson.M{
				"sprint_id":  "",
				"updated_at": time.Now(),
			},
			"$inc": bson.M{"version": 1},
		},
	)
	return err
}

func (s *IssueStore) GetNextIssueNumber(ctx context.Context, projectKey string) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "issue_key", Value: -1}})
	var issue models.Issue
	err := s.collection.FindOne(
		ctx,
		bson.M{"issue_key": bson.M{"$regex": fmt.Sprintf("^%s-", projectKey)}},
		opts,
	).Decode(&issue)

	if err == mongo.ErrNoDocuments {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}

	var num int
	fmt.Sscanf(issue.IssueKey, projectKey+"-%d", &num)
	return num + 1, nil
}

func (s *IssueStore) Search(ctx context.Context, query string, filters bson.M, limit, skip int) ([]*models.Issue, error) {
	if query != "" {
		filters["$text"] = bson.M{"$search": query}
	}

	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip))
	cursor, err := s.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var issues []*models.Issue
	if err := cursor.All(ctx, &issues); err != nil {
		return nil, err
	}
	return issues, nil
}

func (s *IssueStore) SearchWithCursor(ctx context.Context, filters bson.M, limit int, cursor string) ([]*models.Issue, string, error) {
	if limit <= 0 {
		limit = 50
	}

	query := filters
	if query == nil {
		query = bson.M{}
	}

	if cursor != "" {
		parts := strings.SplitN(cursor, "|", 2)
		if len(parts) == 2 {
			ts, err := time.Parse(time.RFC3339Nano, parts[0])
			if err == nil {
				cursorFilter := bson.M{
					"$or": []bson.M{
						{"updated_at": bson.M{"$lt": ts}},
						{
							"$and": []bson.M{
								{"updated_at": ts},
								{"_id": bson.M{"$lt": parts[1]}},
							},
						},
					},
				}

				if len(query) == 0 {
					query = cursorFilter
				} else {
					query = bson.M{"$and": []bson.M{query, cursorFilter}}
				}
			}
		}
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}, {Key: "_id", Value: -1}}).
		SetLimit(int64(limit + 1))

	cursorResult, err := s.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, "", err
	}
	defer cursorResult.Close(ctx)

	var issues []*models.Issue
	if err := cursorResult.All(ctx, &issues); err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(issues) > limit {
		last := issues[limit-1]
		nextCursor = fmt.Sprintf("%s|%s", last.UpdatedAt.UTC().Format(time.RFC3339Nano), last.ID)
		issues = issues[:limit]
	}

	return issues, nextCursor, nil
}
