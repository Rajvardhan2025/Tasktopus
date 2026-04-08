package provider

import (
	"context"
	"log"

	"github.com/yourusername/project-management/config"
	"github.com/yourusername/project-management/service"
	"github.com/yourusername/project-management/store"
)

type Provider struct {
	Config *config.Config
	DB     *store.MongoDB

	// Stores
	ProjectStore      *store.ProjectStore
	IssueStore        *store.IssueStore
	SprintStore       *store.SprintStore
	UserStore         *store.UserStore
	CommentStore      *store.CommentStore
	ActivityStore     *store.ActivityStore
	WorkflowStore     *store.WorkflowStore
	NotificationStore *store.NotificationStore

	// Services
	IssueService        *service.IssueService
	WorkflowService     *service.WorkflowService
	SprintService       *service.SprintService
	NotificationService *service.NotificationService
	SearchService       *service.SearchService
	WebSocketService    *service.WebSocketService
}

func New(cfg *config.Config) (*Provider, error) {
	// Initialize MongoDB
	db, err := store.NewMongoDB(cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		return nil, err
	}

	// Create indexes
	if err := db.CreateIndexes(context.Background()); err != nil {
		log.Printf("Warning: failed to create indexes: %v", err)
	}

	// Initialize stores
	projectStore := store.NewProjectStore(db.Database)
	issueStore := store.NewIssueStore(db.Database)
	sprintStore := store.NewSprintStore(db.Database)
	userStore := store.NewUserStore(db.Database)
	commentStore := store.NewCommentStore(db.Database)
	activityStore := store.NewActivityStore(db.Database)
	workflowStore := store.NewWorkflowStore(db.Database)
	notificationStore := store.NewNotificationStore(db.Database)

	// Initialize services
	wsService := service.NewWebSocketService()
	notificationService := service.NewNotificationService(notificationStore)
	workflowService := service.NewWorkflowService(workflowStore)
	issueService := service.NewIssueService(
		issueStore,
		projectStore,
		activityStore,
		workflowService,
		notificationService,
		wsService,
	)
	sprintService := service.NewSprintService(sprintStore, issueStore, activityStore)
	searchService := service.NewSearchService(issueStore)

	return &Provider{
		Config:              cfg,
		DB:                  db,
		ProjectStore:        projectStore,
		IssueStore:          issueStore,
		SprintStore:         sprintStore,
		UserStore:           userStore,
		CommentStore:        commentStore,
		ActivityStore:       activityStore,
		WorkflowStore:       workflowStore,
		NotificationStore:   notificationStore,
		IssueService:        issueService,
		WorkflowService:     workflowService,
		SprintService:       sprintService,
		NotificationService: notificationService,
		SearchService:       searchService,
		WebSocketService:    wsService,
	}, nil
}

func (p *Provider) Close() error {
	return p.DB.Close()
}
