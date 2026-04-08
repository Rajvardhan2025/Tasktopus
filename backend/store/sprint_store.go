package store

import (
	"context"
	"time"

	"github.com/yourusername/project-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SprintStore struct {
	collection *mongo.Collection
}

func NewSprintStore(db *mongo.Database) *SprintStore {
	return &SprintStore{
		collection: db.Collection("sprints"),
	}
}

func (s *SprintStore) Create(ctx context.Context, sprint *models.Sprint) error {
	sprint.CreatedAt = time.Now()
	sprint.UpdatedAt = time.Now()
	sprint.Status = models.SprintStatusPlanned
	_, err := s.collection.InsertOne(ctx, sprint)
	return err
}

func (s *SprintStore) FindByID(ctx context.Context, id string) (*models.Sprint, error) {
	var sprint models.Sprint
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&sprint)
	if err != nil {
		return nil, err
	}
	return &sprint, nil
}

func (s *SprintStore) FindByProject(ctx context.Context, projectID string) ([]*models.Sprint, error) {
	cursor, err := s.collection.Find(ctx, bson.M{"project_id": projectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var sprints []*models.Sprint
	if err := cursor.All(ctx, &sprints); err != nil {
		return nil, err
	}
	return sprints, nil
}

func (s *SprintStore) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *SprintStore) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
