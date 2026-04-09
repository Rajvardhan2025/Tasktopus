package store

import (
	"context"
	"regexp"
	"time"

	"github.com/yourusername/project-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentStore struct {
	collection *mongo.Collection
}

func NewCommentStore(db *mongo.Database) *CommentStore {
	return &CommentStore{
		collection: db.Collection("comments"),
	}
}

func (s *CommentStore) Create(ctx context.Context, comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, comment)
	return err
}

func (s *CommentStore) FindByID(ctx context.Context, id string) (*models.Comment, error) {
	var comment models.Comment
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentStore) FindByIssue(ctx context.Context, issueID string) ([]*models.Comment, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cursor, err := s.collection.Find(ctx, bson.M{"issue_id": issueID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*models.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentStore) Update(ctx context.Context, id string, content string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"content": content, "updated_at": time.Now()}},
	)
	return err
}

func (s *CommentStore) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *CommentStore) SearchIssueIDsByContent(ctx context.Context, text, projectID string, limit int) ([]string, error) {
	if text == "" {
		return []string{}, nil
	}

	if limit <= 0 {
		limit = 100
	}

	matchStage := bson.D{{
		Key: "$match",
		Value: bson.M{
			"content": bson.M{
				"$regex":   regexp.QuoteMeta(text),
				"$options": "i",
			},
		},
	}}

	lookupStage := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         "issues",
			"localField":   "issue_id",
			"foreignField": "_id",
			"as":           "issue",
		},
	}}

	unwindStage := bson.D{{Key: "$unwind", Value: "$issue"}}

	pipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage}
	if projectID != "" {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M{"issue.project_id": projectID}}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$group", Value: bson.M{"_id": "$issue_id"}}},
		bson.D{{Key: "$limit", Value: limit}},
	)

	cursor, err := s.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rows []struct {
		ID string `bson:"_id"`
	}
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}

	return ids, nil
}
