package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/project-management/models"
	"github.com/yourusername/project-management/store"
	"go.mongodb.org/mongo-driver/bson"
)

type SprintService struct {
	sprintStore   *store.SprintStore
	issueStore    *store.IssueStore
	activityStore *store.ActivityStore
	wsSvc         *WebSocketService
}

func NewSprintService(
	sprintStore *store.SprintStore,
	issueStore *store.IssueStore,
	activityStore *store.ActivityStore,
	wsSvc *WebSocketService,
) *SprintService {
	return &SprintService{
		sprintStore:   sprintStore,
		issueStore:    issueStore,
		activityStore: activityStore,
		wsSvc:         wsSvc,
	}
}

func (s *SprintService) Create(ctx context.Context, req *models.CreateSprintRequest) (*models.Sprint, error) {
	// Validate project exists
	if req.ProjectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}

	// If no dates provided, set defaults
	startDate := req.StartDate
	endDate := req.EndDate

	if startDate == nil {
		now := time.Now()
		startDate = &now
	}

	if endDate == nil && startDate != nil {
		futureDate := startDate.AddDate(0, 0, 14) // Default 2-week sprint
		endDate = &futureDate
	}

	sprint := &models.Sprint{
		BaseModel:  models.BaseModel{ID: uuid.New().String()},
		ProjectID:  req.ProjectID,
		Name:       req.Name,
		Goal:       req.Goal,
		StartDate:  startDate,
		EndDate:    endDate,
		Status:     models.SprintStatusFuture,
		IsDefault:  false,
		IssueCount: 0,
	}

	if err := s.sprintStore.Create(ctx, sprint); err != nil {
		return nil, fmt.Errorf("failed to create sprint: %w", err)
	}

	return sprint, nil
}

func (s *SprintService) Start(ctx context.Context, sprintID, userID string) (*models.Sprint, error) {
	sprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	if sprint.Status != models.SprintStatusFuture {
		return nil, fmt.Errorf("only future sprints can be started, current status: %s", sprint.Status)
	}

	// Check if another sprint is already active
	activeSprints, err := s.sprintStore.FindByProject(ctx, sprint.ProjectID)
	if err == nil {
		for _, s := range activeSprints {
			if s.Status == models.SprintStatusActive {
				return nil, fmt.Errorf("another sprint is already active: %s", s.Name)
			}
		}
	}

	// Calculate total points and issue count
	issues, err := s.issueStore.FindBySprint(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sprint issues: %w", err)
	}

	totalPoints := 0
	for _, issue := range issues {
		totalPoints += issue.StoryPoints
	}

	now := time.Now()
	update := bson.M{
		"status":       models.SprintStatusActive,
		"total_points": totalPoints,
		"issue_count":  len(issues),
		"start_date":   &now,
		"updated_at":   &now,
	}

	if err := s.sprintStore.Update(ctx, sprintID, update); err != nil {
		return nil, fmt.Errorf("failed to update sprint: %w", err)
	}

	// Log activity
	activity := &models.Activity{
		ID:        uuid.New().String(),
		ProjectID: sprint.ProjectID,
		UserID:    userID,
		Action:    models.ActivitySprintStarted,
		Changes: map[string]interface{}{
			"sprint_id":   sprintID,
			"sprint_name": sprint.Name,
			"status":      models.SprintStatusActive,
		},
	}
	_ = s.activityStore.Create(ctx, activity)

	updatedSprint, _ := s.sprintStore.FindByID(ctx, sprintID)
	s.wsSvc.BroadcastToProject(sprint.ProjectID, models.WSEvent{
		Type:      models.WSEventSprintUpdated,
		ProjectID: sprint.ProjectID,
		Data:      updatedSprint,
	})
	return updatedSprint, nil
}

func (s *SprintService) Close(ctx context.Context, sprintID string, req *models.CloseSprintRequest, userID string) (*models.Sprint, error) {
	sprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	if sprint.Status != models.SprintStatusActive {
		return nil, fmt.Errorf("only active sprints can be closed, current status: %s", sprint.Status)
	}

	// Get all issues in sprint
	issues, err := s.issueStore.FindBySprint(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sprint issues: %w", err)
	}

	// Calculate completed points and issue counts
	completedPoints := 0
	completedIssueCount := 0
	incompleteIssues := []string{}

	for _, issue := range issues {
		if issue.Status == "done" {
			completedPoints += issue.StoryPoints
			completedIssueCount++
		} else {
			incompleteIssues = append(incompleteIssues, issue.ID)
		}
	}

	// Handle carry-over - move non-carried issues to backlog
	carryOverMap := make(map[string]bool)
	for _, issueID := range req.CarryOverIssues {
		carryOverMap[issueID] = true
	}

	for _, issueID := range incompleteIssues {
		if !carryOverMap[issueID] {
			// Remove from sprint (move to backlog)
			if err := s.issueStore.MoveToBacklog(ctx, issueID, sprintID); err != nil {
				return nil, fmt.Errorf("failed to move issue %s to backlog: %w", issueID, err)
			}
		}
	}

	// Calculate velocity (completed points / sprint duration)
	velocity := 0
	if sprint.StartDate != nil && sprint.EndDate != nil {
		durationDays := sprint.EndDate.Sub(*sprint.StartDate).Hours() / 24
		if durationDays > 0 {
			velocity = int(float64(completedPoints) / durationDays * 7) // Velocity per week
		}
	}

	now := time.Now()
	update := bson.M{
		"status":               models.SprintStatusClosed,
		"completed_points":     completedPoints,
		"complete_issue_count": completedIssueCount,
		"velocity":             velocity,
		"end_date":             &now,
		"updated_at":           &now,
	}

	if err := s.sprintStore.Update(ctx, sprintID, update); err != nil {
		return nil, fmt.Errorf("failed to update sprint: %w", err)
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
			"status":           models.SprintStatusClosed,
			"completed_points": completedPoints,
			"total_points":     sprint.TotalPoints,
			"completed_issues": completedIssueCount,
			"total_issues":     len(issues),
			"velocity":         velocity,
			"carry_over_count": len(req.CarryOverIssues),
		},
	}
	_ = s.activityStore.Create(ctx, activity)

	updatedSprint, _ := s.sprintStore.FindByID(ctx, sprintID)
	s.wsSvc.BroadcastToProject(sprint.ProjectID, models.WSEvent{
		Type:      models.WSEventSprintUpdated,
		ProjectID: sprint.ProjectID,
		Data:      updatedSprint,
	})
	return updatedSprint, nil
}

// Backward compatibility - alias Complete to Close
func (s *SprintService) Complete(ctx context.Context, sprintID string, req *models.CloseSprintRequest, userID string) (*models.Sprint, error) {
	return s.Close(ctx, sprintID, req, userID)
}

func (s *SprintService) GetByID(ctx context.Context, sprintID string) (*models.Sprint, error) {
	sprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}
	return sprint, nil
}

func (s *SprintService) GetByProject(ctx context.Context, projectID string) ([]*models.Sprint, error) {
	sprints, err := s.sprintStore.FindByProject(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sprints: %w", err)
	}
	return sprints, nil
}

func (s *SprintService) Update(ctx context.Context, sprintID string, req *models.UpdateSprintRequest) (*models.Sprint, error) {
	sprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, fmt.Errorf("sprint not found: %w", err)
	}

	update := bson.M{}
	if req.Name != nil {
		update["name"] = *req.Name
	}
	if req.Goal != nil {
		update["goal"] = *req.Goal
	}
	if req.StartDate != nil {
		update["start_date"] = req.StartDate
	}
	if req.EndDate != nil {
		update["end_date"] = req.EndDate
	}

	if len(update) == 0 {
		return sprint, nil
	}

	if err := s.sprintStore.Update(ctx, sprintID, update); err != nil {
		return nil, fmt.Errorf("failed to update sprint: %w", err)
	}

	updatedSprint, err := s.sprintStore.FindByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	s.wsSvc.BroadcastToProject(updatedSprint.ProjectID, models.WSEvent{
		Type:      models.WSEventSprintUpdated,
		ProjectID: updatedSprint.ProjectID,
		Data:      updatedSprint,
	})

	return updatedSprint, nil
}
