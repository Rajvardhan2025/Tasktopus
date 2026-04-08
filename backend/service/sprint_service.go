package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
	"go.mongodb.org/mongo-driver/bson"
)

type SprintService struct {
	sprintStore   *store.SprintStore
	issueStore    *store.IssueStore
	activityStore *store.ActivityStore
}

func NewSprintService(
	sprintStore *store.SprintStore,
	issueStore *store.IssueStore,
	activityStore *store.ActivityStore,
) *SprintService {
	return &SprintService{
		sprintStore:   sprintStore,
		issueStore:    issueStore,
		activityStore: activityStore,
	}
}

func (s *SprintService) Create(ctx context.Context, req *models.CreateSprintRequest) (*models.Sprint, error) {
	sprint := &models.Sprint{
		BaseModel: models.BaseModel{ID: uuid.New().String()},
		ProjectID: req.ProjectID,
		Name:      req.Name,
		Goal:      req.Goal,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if err := s.sprintStore.Create(ctx, sprint); err != nil {
		return nil, err
	}

	return sprint, nil
}

func (s *SprintService) Start(ctx context.Context, sprintID, userID string) (*models.Sprint, error) {
	sprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint.Status != models.SprintStatusPlanned {
		return nil, fmt.Errorf("sprint is already started or completed")
	}

	// Calculate total points
	issues, err := s.issueStore.FindBySprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	totalPoints := 0
	for _, issue := range issues {
		totalPoints += issue.StoryPoints
	}

	update := bson.M{
		"status":       models.SprintStatusActive,
		"total_points": totalPoints,
	}

	if err := s.sprintStore.Update(ctx, sprintID, update); err != nil {
		return nil, err
	}

	// Log activity
	activity := &models.Activity{
		ID:        uuid.New().String(),
		ProjectID: sprint.ProjectID,
		UserID:    userID,
		Action:    models.ActivitySprintStarted,
		Changes:   map[string]interface{}{"sprint_id": sprintID, "sprint_name": sprint.Name},
	}
	s.activityStore.Create(ctx, activity)

	return s.sprintStore.FindByID(ctx, sprintID)
}

func (s *SprintService) Complete(ctx context.Context, sprintID string, req *models.CompleteSprintRequest, userID string) (*models.Sprint, error) {
	sprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	if sprint.Status != models.SprintStatusActive {
		return nil, fmt.Errorf("sprint is not active")
	}

	// Get all issues in sprint
	issues, err := s.issueStore.FindBySprint(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	// Calculate completed points
	completedPoints := 0
	incompleteIssues := []string{}

	for _, issue := range issues {
		if issue.Status == "done" {
			completedPoints += issue.StoryPoints
		} else {
			incompleteIssues = append(incompleteIssues, issue.ID)
		}
	}

	// Handle carry-over
	carryOverMap := make(map[string]bool)
	for _, issueID := range req.CarryOverIssues {
		carryOverMap[issueID] = true
	}

	for _, issueID := range incompleteIssues {
		if !carryOverMap[issueID] {
			// Remove from sprint (move to backlog)
			s.issueStore.UpdateWithVersion(ctx, issueID, 0, bson.M{"sprint_id": ""})
		}
	}

	// Update sprint
	update := bson.M{
		"status":           models.SprintStatusCompleted,
		"completed_points": completedPoints,
	}

	if err := s.sprintStore.Update(ctx, sprintID, update); err != nil {
		return nil, err
	}

	// Log activity
	activity := &models.Activity{
		ID:        uuid.New().String(),
		ProjectID: sprint.ProjectID,
		UserID:    userID,
		Action:    models.ActivitySprintCompleted,
		Changes: map[string]interface{}{
			"sprint_id":        sprintID,
			"sprint_name":      sprint.Name,
			"completed_points": completedPoints,
			"total_points":     sprint.TotalPoints,
		},
	}
	s.activityStore.Create(ctx, activity)

	return s.sprintStore.FindByID(ctx, sprintID)
}

func (s *SprintService) GetByID(ctx context.Context, sprintID string) (*models.Sprint, error) {
	return s.sprintStore.FindByID(ctx, sprintID)
}

func (s *SprintService) GetByProject(ctx context.Context, projectID string) ([]*models.Sprint, error) {
	return s.sprintStore.FindByProject(ctx, projectID)
}
