package store

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoDB(uri, database string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(database)

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

func (m *MongoDB) CreateIndexes(ctx context.Context) error {
	// Issues indexes
	issuesCol := m.Database.Collection("issues")
	_, err := issuesCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: map[string]interface{}{"project_id": 1, "status": 1}},
		{Keys: map[string]interface{}{"project_id": 1, "sprint_id": 1}},
		{Keys: map[string]interface{}{"assignee_id": 1}},
		{Keys: map[string]interface{}{"issue_key": 1}, Options: options.Index().SetUnique(true)},
		{Keys: map[string]interface{}{"title": "text", "description": "text"}},
	})
	if err != nil {
		return fmt.Errorf("failed to create issues indexes: %w", err)
	}

	// Activities indexes
	activitiesCol := m.Database.Collection("activities")
	_, err = activitiesCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: map[string]interface{}{"project_id": 1, "timestamp": -1}},
		{Keys: map[string]interface{}{"issue_id": 1, "timestamp": -1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create activities indexes: %w", err)
	}

	// Comments indexes
	commentsCol := m.Database.Collection("comments")
	_, err = commentsCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: map[string]interface{}{"issue_id": 1, "created_at": 1}},
		{Keys: map[string]interface{}{"parent_id": 1}},
	})
	if err != nil {
		return fmt.Errorf("failed to create comments indexes: %w", err)
	}

	// Projects indexes
	projectsCol := m.Database.Collection("projects")
	_, err = projectsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]interface{}{"key": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create projects indexes: %w", err)
	}

	return nil
}
