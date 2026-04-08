package store

import (
	"context"
	"time"

	"github.com/yourusername/project-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkflowStore struct {
	collection *mongo.Collection
}

func NewWorkflowStore(db *mongo.Database) *WorkflowStore {
	return &WorkflowStore{
		collection: db.Collection("workflows"),
	}
}

func (s *WorkflowStore) Create(ctx context.Context, workflow *models.Workflow) error {
	workflow.CreatedAt = time.Now()
	workflow.UpdatedAt = time.Now()
	_, err := s.collection.InsertOne(ctx, workflow)
	return err
}

func (s *WorkflowStore) FindByID(ctx context.Context, id string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := s.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (s *WorkflowStore) FindByProject(ctx context.Context, projectID string) (*models.Workflow, error) {
	var workflow models.Workflow
	err := s.collection.FindOne(ctx, bson.M{"project_id": projectID}).Decode(&workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (s *WorkflowStore) Update(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
