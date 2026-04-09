package store

import (
	"context"
	"time"

	"github.com/yourusername/project-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProjectStore struct {
	collection *mongo.Collection
}

func NewProjectStore(db *mongo.Database) *ProjectStore {
	return &ProjectStore{
		collection: db.Collection("projects"),
	}
}

func (s *ProjectStore) Create(ctx context.Context, project *models.Project) error {
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, project)
	return err
}

func (s *ProjectStore) FindByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectStore) FindByKey(ctx context.Context, key string) (*models.Project, error) {
	var project models.Project
	err := s.collection.FindOne(ctx, bson.M{"key": key}).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectStore) List(ctx context.Context) ([]*models.Project, error) {
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var projects []*models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *ProjectStore) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *ProjectStore) Delete(ctx context.Context, id string) error {
	_, err := s.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *ProjectStore) AddMember(ctx context.Context, projectID, userID string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": projectID},
		bson.M{"$addToSet": bson.M{"members": userID}},
	)
	return err
}

func (s *ProjectStore) RemoveMember(ctx context.Context, projectID, userID string) error {
	_, err := s.collection.UpdateOne(
		ctx,
		bson.M{"_id": projectID},
		bson.M{"$pull": bson.M{"members": userID}},
	)
	return err
}
